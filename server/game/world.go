package game

import (
	_ "log"
	"math"
	"time"

	"github.com/google/uuid"
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
				cell, ok := w.Cells[id]
				if !ok {
					continue // cell doesn't exit
				}

				// Cell may have post-split momentum
				if cell.Momentum.MagnitudeSq() != 0.0 {
					cell.Position = cell.Position.Add(cell.Momentum)
					cell.Momentum = cell.Momentum.Scale(SplitVelocityDecay)
					if cell.Momentum.MagnitudeSq() < 0.01 {
						cell.Momentum = Vec2{}
					}
				}

				// Calculate new direction and velocity
				desired := p.Target.Sub(cell.Position).Normalize()
				cell.Direction = cell.Direction.Add(
					desired.Sub(cell.Direction).Scale(TurnSpeed), // smooth turning towards desired dir
				).Normalize()
				cell.Position = cell.Position.Add(
					cell.Direction.Scale(
						w.calculateSpeed(cell.Mass), // velocity in the direction of cell
					),
				)
			}
		}

		// Collisions
		for _, p := range w.Players {
			for i, idA := range p.CellIDs {
				cellA := w.Cells[idA]
				for _, idB := range p.CellIDs[i+1:] {
					cellB := w.Cells[idB]
					diff := cellB.Position.Sub(cellA.Position)
					dist := diff.Magnitude()
					minDist := cellA.Radius() + cellB.Radius()
					if dist < minDist && dist > 0 {
						// Push apart
						overlap := minDist - dist
						push := diff.Normalize().Scale(overlap * 0.5)

						cellA.Position = cellA.Position.Sub(push)
						cellB.Position = cellB.Position.Add(push)
					}
				}
			}
		}

		// Clamp to world bounds
		for _, p := range w.Players {
			for _, id := range p.CellIDs {
				cell, ok := w.Cells[id]
				if !ok {
					continue // cell doesn't exist
				}

				cell.Position.X = Clamp(cell.Position.X, 0.0, WorldWidth)
				cell.Position.Y = Clamp(cell.Position.Y, 0.0, WorldHeight)

				// TODO: Decay split cooldown

				// TODO: Recombine

				for _, pellet := range w.Pellets {
					if cell.Position.DistanceTo(pellet.Position) < cell.Radius() {
						// Eat the pellet
						cell.Mass += pellet.Mass * 4
						delete(w.Pellets, pellet.ID)
					}
				}
				for _, otherCell := range w.Cells {
					if otherCell.OwnerID == p.ID {
						continue // we don't eat our own cells
					}

					if float64(cell.Mass) <= float64(otherCell.Mass)*EatMassThreshold {
						continue // we are not big enough
					}

					if cell.Position.DistanceTo(otherCell.Position) >=
						(cell.Radius() - otherCell.Radius()*EatDistanceThreshold) {
						continue // we are not overlapping enough
					}

					// Eat the cell
					cell.Mass += otherCell.Mass
					victim := w.Players[otherCell.OwnerID]
					victim.removeCell(otherCell.ID)
					delete(w.Cells, otherCell.ID)

					if len(victim.CellIDs) == 0 {
						w.handlePlayerRemove(victim.ID) // dead
					}
				}

				// TODO: Cell hits virus
			}
		}

		w.spawnPellets()

		// TODO: Spawn viruses

		// TODO: Remove dead players

		// Broadcast state
		w.broadcastState()
	}
}

func (w *World) spawnPellets() {
	currNumPellets := int64(len(w.Pellets))
	if currNumPellets < MinPellets {
		// Spawn more
		additionalPelletNum := RandIntRange(0, MaxPellets-currNumPellets)
		for range additionalPelletNum {
			randMass := RandIntRange(MinPelletMass, MaxPelletMass)
			pellet := Pellet{
				ID: uuid.New(),

				Position: Vec2{
					X: RandFloatRange(0.0, WorldWidth),
					Y: RandFloatRange(0.0, WorldHeight),
				},

				Mass:  randMass,
				Color: w.findStartColor(),
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
			Me:      p.ID,
			Cells:   cells,
			Pellets: pellets,
			Viruses: viruses,
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
		X: RandFloatRange(0.0, WorldWidth),
		Y: RandFloatRange(0.0, WorldHeight),
	}
}

/** findStartColor
 * TODO: ideally this would generate as few
 *  color collisions as possible
 */
func (w *World) findStartColor() [3]uint8 {
	return getRandomColor()
}

func getRandomColor() [3]uint8 {
	return [3]uint8{
		uint8(RandIntRange(0, 256)),
		uint8(RandIntRange(0, 256)),
		uint8(RandIntRange(0, 256)),
	}
}

/** calculateSpeed
 * Cell speed is a function of mass
 */
func (w *World) calculateSpeed(mass int64) float64 {
	speed := BaseSpeed * math.Pow(float64(StartMass)/float64(mass), SpeedExponent)
	return math.Max(speed, MinSpeed)
}
