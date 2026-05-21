package game

import (
	"log"
	"math"
	"time"

	"github.com/google/uuid"

	"github.com/infopek/agorio/server/utils"
)

type World struct {
	Cells   map[uuid.UUID]*Cell
	Pellets map[uuid.UUID]*Pellet
	Viruses map[uuid.UUID]*Virus
	Players map[uuid.UUID]*Player

	tick uint64

	InputChan chan PlayerMessage // shared across all players
}

func NewWorld() *World {
	return &World{
		Cells:     make(map[uuid.UUID]*Cell),
		Pellets:   make(map[uuid.UUID]*Pellet),
		Viruses:   make(map[uuid.UUID]*Virus),
		Players:   make(map[uuid.UUID]*Player),
		InputChan: make(chan PlayerMessage, 256),
	}
}

func (w *World) Tick() {
	ticker := time.NewTicker(time.Second / time.Duration(TickRate))
	defer ticker.Stop()

	for range ticker.C {
		w.tick++
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
		w.spawnPellets()

		// TODO: Spawn viruses

		// TODO: Remove dead players

		// Broadcast state
		w.broadcastState()
	}
}

func (w *World) addPlayer(playerID uuid.UUID, name string, outputChan chan<- ServerMessage) {
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
	log.Printf("added player %v with cell %v\n", player.ID, cell.ID)
}

func (w *World) removePlayer(id uuid.UUID) {
	p := w.Players[id]
	if p == nil {
		return
	}

	// Remove player's cells
	for _, cellID := range p.CellIDs {
		delete(w.Cells, cellID)
	}

	close(p.OutputChan)
	delete(w.Players, id)
	log.Printf("removed player %v\n", id)
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
				w.addPlayer(m.PlayerID, m.Name, m.OutputChan)
			case DisconnectMessage:
				w.removePlayer(m.PlayerID)
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

func (w *World) spawnPellets() {
	currNumPellets := int64(len(w.Pellets))
	if currNumPellets < MinPellets {
		// Spawn more
		additionalPelletNum := utils.RandIntRange(0, MaxPellets-currNumPellets)
		for range additionalPelletNum {
			randMass := utils.RandIntRange(MinPelletMass, MaxPelletMass)
			pellet := Pellet{
				ID: uuid.New(),

				Position: Vec2{
					X: utils.RandFloatRange(0.0, WorldWidth),
					Y: utils.RandFloatRange(0.0, WorldHeight),
				},

				Mass:   randMass,
				Radius: w.calculateRadius(randMass),
				Color: utils.GetRandomColor(),
			}
			w.Pellets[pellet.ID] = &pellet
		}
	}
}

func (w *World) broadcastState() {
	cells := make([]Cell, 0, len(w.Cells))
	pellets := make([]Pellet, 0, len(w.Pellets))
	viruses := make([]Virus, 0, len(w.Viruses))

	for _, c := range w.Cells {
		cells = append(cells, *c)
	}
	for _, p := range w.Pellets {
		pellets = append(pellets, *p)
	}
	for _, v := range w.Viruses {
		viruses = append(viruses, *v)
	}

	for _, p := range w.Players {
		snapshot := TickSnapshot{
			Cells:   cells,
			Pellets: pellets,
			Viruses: viruses,
			Me:      p.CellIDs,
			Score:   p.Score,
			Tick:    w.tick,
		}

		p.OutputChan <- snapshot
	}
}

/** findStartPosition
 * TODO: find an empty place in the map
 *  for the player to spawn
 */
func (w *World) findStartPosition() Vec2 {
	return Vec2{
		X: utils.RandFloatRange(0.0, WorldWidth),
		Y: utils.RandFloatRange(0.0, WorldHeight),
	}
}

/** findStartColor
 * TODO: ideally this would generate as few
 *  color collisions as possible
 */
func (w *World) findStartColor() [3]uint8 {
	return utils.GetRandomColor()
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
