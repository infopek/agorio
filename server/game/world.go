package game

import (
	_ "log"
	"math"
	"time"

	"github.com/google/uuid"

	"github.com/infopek/agorio/server/utils"
)

type World struct {
	Cells     map[uuid.UUID]*Cell
	Players   map[uuid.UUID]*Player
	InputChan chan PlayerMessage // shared across all players
}

func NewWorld() *World {
	return &World{}
}

func (w *World) Tick() {
	ticker := time.NewTicker(time.Second / time.Duration(TickRate))
	defer ticker.Stop()

	for range ticker.C {
		w.processInputs()

		// Move cells
		for _, p := range w.Players {
			for _, id := range p.CellIDs {
				cell := w.Cells[id]
				if cell.Momentum.MagnitudeSq() != 0.0 {
					// TODO: apply post-split momentum, then decay
				} else {
					dir := p.Target.Sub(cell.Position).Normalize()
					vel := dir.Scale(w.calculateSpeed(cell.Mass))
					cell.Position = cell.Position.Add(vel)
				}
				// TODO: clamp to world bounds
			}
		}

		// TODO: Decay split cooldown

		// TODO: Recombine

		// TODO: Check eating
		//  - Cell eats pellet (food)
		//  - Cell eats smaller cell
		//  - Cell hits virus

		// TODO: Spawn pellets

		// TODO: Spawn viruses

		// TODO: Remove dead players

		// Broadcast state
		w.broadcastState()
	}
}

func (w *World) addPlayer(name string, outputChan chan<- TickSnapshot) {
	playerID := uuid.New()
	cellID := uuid.New()

	cell := Cell{
		ID:      cellID,
		OwnerID: playerID,

		Position: w.findStartPosition(),
		Momentum: Vec2{},

		Radius: w.calculateRadius(StartMass),
		Mass:   StartMass,
		Color:  w.findStartColor(),
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
}

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
				w.addPlayer(m.Name, make(chan TickSnapshot))
			case MoveMessage:
				w.Players[m.PlayerID].Target = m.Target
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
			w.splitPlayer(m.PlayerID)
		case FeedMessage:
			w.ejectMass(m.PlayerID)
		default:
			return // shouldn't happen
		}
	}
}

func (w *World) broadcastState() error {

	return nil
}

/** findStartPosition
 * TODO: find an empty place in the map
 *  for the player to spawn
 */
func (w *World) findStartPosition() Vec2 {
	return Vec2{
		X: float64(utils.RandIntRange(0, WorldWidth)),
		Y: float64(utils.RandIntRange(0, WorldHeight)),
	}
}

/** findStartColor
 * TODO: ideally this would generate as few
 *  color collisions as possible
 */
func (w *World) findStartColor() [3]uint8 {
	return [3]uint8{
		uint8(utils.RandIntRange(0, 256)),
		uint8(utils.RandIntRange(0, 256)),
		uint8(utils.RandIntRange(0, 256)),
	}
}

/** calculateRadius
 * Cell radius is a function of mass
 *
 * radius = radiusScale * sqrt(mass)
 */
func (w *World) calculateRadius(mass int64) float64 {
	return RadiusScale * math.Sqrt(float64(mass))
}

/** calculateSpeed
 * Cell speed is a function of mass
 *
 * sp = baseSp * (1.0 / sqrt(mass / startMass))
 */
func (w *World) calculateSpeed(mass int64) float64 {
	return BaseSpeed * (1.0 / math.Sqrt(float64(mass)/float64(StartMass)))
}
