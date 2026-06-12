package game

import (
	"log"
	"math"
	"sort"
	"strings"

	"github.com/google/uuid"
)

/** processInputs
 *
 * First, it consumes the move messages, so that all the
 *  actions are working with up-to-date target vectors
 */
func (w *World) processInputs() {
	var actions []PlayerEvent

	// Drain channel, why tf do I need labels
DrainLoop:
	for {
		select {
		case msg := <-w.InputChan:
			switch m := msg.(type) {
			case PlayerJoinEvent:
				w.handlePlayerConnect(m.PlayerID, m.Name, m.OutputChan)
			case PlayerDisconnectEvent:
				w.handlePlayerDisconnect(m.PlayerID)
			case PlayerMoveEvent:
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
		case PlayerSplitEvent:
			w.handlePlayerSplit(m.PlayerID)
		case PlayerFeedEvent:
			w.handlePlayerEject(m.PlayerID)
		default:
			return // shouldn't happen
		}
	}
}

/** handlePlayerMove
 *
 * Sets the target of the player to the new one
 */
func (w *World) handlePlayerMove(playerID uuid.UUID, newTarget Vec2) {
	p, ok := w.Players[playerID]
	if !ok {
		return // player no longer exists
	}

	p.Target = newTarget
}

/** World.handlePlayerDisconnect
 *
 * Removes the player from the world, and closes their output channel
 */
func (w *World) handlePlayerDisconnect(playerID uuid.UUID) {
	p, ok := w.Players[playerID]
	if !ok {
		return // already gone
	}

	w.removePlayer(p.ID)
	close(p.OutputChan)

	log.Printf("player %v disconnected\n", playerID)
}

/** handlePlayerConnect
 *
 * Creates a new instance for the client with the given name
 *  and assigns the session's id to it, also creates a starter cell
 */
func (w *World) handlePlayerConnect(playerID uuid.UUID, name string, outputChan chan<- ServerEvent) {
	if _, ok := w.Players[playerID]; ok {
		return // already joined
	}

	cellID := uuid.New()
	name = sanitizeName(name)

	c := Cell{
		ID:      cellID,
		OwnerID: playerID,

		Position:  w.findEmptySpace(Radius(StartMass)),
		Direction: Vec2{},
		Momentum:  Vec2{},

		MergeTimer: MergeTimerStartSeconds,

		Mass:  StartMass * 20.0,
		Color: w.getRandomColor(),
	}
	p := Player{
		ID:      playerID,
		CellIDs: []uuid.UUID{cellID},

		Name:   name,
		Target: Vec2{},

		IsBot:      false,
		OutputChan: outputChan,
	}

	w.Players[playerID] = &p
	w.Cells[cellID] = &c

	log.Printf("player %v added with starter cell %v\n", p.ID, c.ID)
}

/** World.handlePlayerSplit
 *
 * Splits the player if possible, the cells with the greater mass are prioritized
 *  when approaching the max cell limit
 */
func (w *World) handlePlayerSplit(playerID uuid.UUID) {
	p, ok := w.Players[playerID]
	if !ok {
		return // player not in pool
	}

	cellIDs := make([]uuid.UUID, len(p.CellIDs))
	copy(cellIDs, p.CellIDs)
	sort.Slice(cellIDs, func(i, j int) bool {
		return w.Cells[cellIDs[i]].Mass > w.Cells[cellIDs[j]].Mass // cell with greater mass gets split first
	})

	for _, cellID := range cellIDs {
		if uint64(len(p.CellIDs)) >= SplitMaxCells {
			break // too many cells
		}

		cell, ok := w.Cells[cellID]
		if !ok || cell.Mass < SplitMinMass {
			continue // cell doesn't exist or is not big enough
		}

		// Split
		cell.Mass /= 2.0
		cell.MergeTimer += MergeCooldownSeconds
		dir := p.Target.Sub(cell.Position).Normalize()
		newCell := Cell{
			ID:      uuid.New(),
			OwnerID: playerID,

			Position:  cell.Position.Add(cell.Direction.Scale(cell.Radius() / 3.0)),
			Direction: dir,
			Momentum:  dir.Scale(SplitMomentumFactor + math.Sqrt(cell.Radius())*SplitMomentumRadiusFactor),

			MergeTimer: min(MaxMergeTimer, MergeTimerStartSeconds+(cell.Mass*MergeTimerMassFactor)),

			Mass:  cell.Mass,
			Color: cell.Color,
		}

		p.CellIDs = append(p.CellIDs, newCell.ID)
		w.Cells[newCell.ID] = &newCell
	}
}

/** World.handlePlayerEject
 *
 *
 */
func (w *World) handlePlayerEject(playerID uuid.UUID) {
	p, ok := w.Players[playerID]
	if !ok {
		return // player not in pool
	}

	for _, cellID := range p.CellIDs {
		cell, ok := w.Cells[cellID]
		if !ok || cell.Mass < EjectMinMass {
			continue // cell doesn't exist or is not big enough
		}

		// Eject mass
		cell.Mass -= EjectAmount * EjectPenalty
		dir := p.Target.Sub(cell.Position).Normalize()
		newEject := Eject{
			ID: uuid.New(),

			Position: cell.Position.Add(cell.Direction.Scale(cell.Radius())),
			Momentum: dir.Scale(EjectMomentum),

			Mass:  EjectAmount,
			Color: cell.Color,
		}

		w.Ejects[newEject.ID] = &newEject
	}
}

/** sanitizeName
 *
 * Cleans the name given by a user, uses a random haha
 *  name if empty
 */
func sanitizeName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return -1 // strip control chars
		}
		return r
	}, name)
	if len(name) == 0 {
		return DefaultNames[RandIntRange(0, int64(len(DefaultNames)))]
	}

	runes := []rune(name)
	if uint64(len(runes)) > MaxNameLength {
		name = string(runes[:MaxNameLength])
	}

	return name
}
