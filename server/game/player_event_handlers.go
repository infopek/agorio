package game

import (
	"log"

	"github.com/google/uuid"
)

/** processInputs
 * First, it consumes the move messages, so that all the
 *  actions are working with up-to-date target vectors
 */
func (w *World) processInputs() {
	var actions []PlayerMessage

	// Drain channel, why tf do I need labels
DrainLoop:
	for {
		select {
		case msg := <-w.InputChan:
			switch m := msg.(type) {
			case JoinMessage:
				w.handlePlayerAdd(m.PlayerID, m.Name, m.OutputChan)
			case DisconnectMessage:
				w.handlePlayerDisconnect(m.PlayerID)
			case MoveMessage:
				w.handlePlayerMove(m.PlayerID, m.Target)
			default:
				actions = append(actions, msg)
			}
		default:
			break DrainLoop
		}
	}

	// Execute actions
	for _, msg := range actions {
		switch m := msg.(type) {
		case SplitMessage:
			w.handlePlayerSplit(m.PlayerID)
		case FeedMessage:
			w.handlePlayerMassEject(m.PlayerID)
		default:
			return // shouldn't happen
		}
	}
}

func (w *World) handlePlayerMove(playerID uuid.UUID, newTarget Vec2) {
	p, ok := w.Players[playerID]
	if !ok {
		return // player no longer exists
	}

	p.Target = newTarget
}

func (w *World) handlePlayerDisconnect(id uuid.UUID) {
	p, ok := w.Players[id]
	if !ok {
		return // already gone
	}

	w.handlePlayerRemove(p.ID)
	close(p.OutputChan)

	log.Printf("player %v disconnected\n", id)
}

func (w *World) handlePlayerAdd(playerID uuid.UUID, name string, outputChan chan<- ServerMessage) {
	cellID := uuid.New()

	cell := Cell{
		ID:      cellID,
		OwnerID: playerID,

		Position: w.findStartPosition(),
		Momentum: Vec2{},

		Mass:  StartMass,
		Color: w.findStartColor(),
	}
	player := Player{
		ID:      playerID,
		CellIDs: []uuid.UUID{cellID},

		Name:   name,
		Target: Vec2{},

		OutputChan: outputChan,
	}

	w.Players[playerID] = &player
	w.Cells[cellID] = &cell

	log.Printf("player %v added with starter cell %v\n", player.ID, cell.ID)
}

func (w *World) handlePlayerRemove(playerID uuid.UUID) {
	p, ok := w.Players[playerID]
	if !ok {
		return // player already removed
	}

	// Remove player's remaining cells
	for _, cellID := range p.CellIDs {
		delete(w.Cells, cellID)
	}

	// Notify player of his death
	p.OutputChan <- DeathEvent{}

	delete(w.Players, playerID)

	log.Printf("player %v removed\n", playerID)
}

func (w *World) handlePlayerSplit(playerID uuid.UUID) {
	p, ok := w.Players[playerID]
	if !ok {
		return // player not in pool
	}

	for _, cellID := range p.CellIDs {
		cell, ok := w.Cells[cellID]
		if !ok || cell.Mass < SplitMinMass {
			continue // cell doesn't exist or is not big enough
		}

		// Split
		cell.Mass /= 2
		newCell := Cell{
			ID:      uuid.New(),
			OwnerID: playerID,

			Position:  cell.Position,
			Direction: cell.Direction,
			Momentum:  cell.Direction.Scale(SplitSpeed),

			Mass:  cell.Mass,
			Color: cell.Color,
		}

		p.CellIDs = append(p.CellIDs, newCell.ID)
		w.Cells[newCell.ID] = &newCell
	}
}

func (w *World) handlePlayerMassEject(id uuid.UUID) {
	// TODO: implement
}
