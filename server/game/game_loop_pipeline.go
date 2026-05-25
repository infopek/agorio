package game

import (
	"log"
	"math"

	"github.com/google/uuid"
)

/** World.applyPhysics
 *
 * Loops through all the players, their cells, and applies
 *  potential split momentum, velocity physics
 *
 * Applies momentum to all ejected mass
 *
 */
func (w *World) applyPhysics() {
	for _, p := range w.Players {
		for _, id := range p.CellIDs {
			cell, ok := w.Cells[id]
			if !ok {
				continue // cell doesn't exit
			}

			// Cell may have post-split momentum
			if cell.Momentum.MagnitudeSq() > 0.0 {
				cell.Position = cell.Position.Add(cell.Momentum)
				cell.Momentum = cell.Momentum.Scale(SplitMomentumDecay)
				if cell.Momentum.MagnitudeSq() < 0.01 {
					cell.Momentum = Vec2{} // zero it
				}
			}

			diff := p.Target.Sub(cell.Position)
			dist := diff.Magnitude()
			if dist < cell.Radius()/CellCenterThreshold {
				continue // cursor is basically on center of cell
			}

			// Speed scales with distance
			speedFactor := min(dist/(cell.Radius()*SpeedFactorRange), 1.0)

			// Calculate new direction and velocity
			desired := diff.Scale(1.0 / dist) // normalize
			cell.Direction = cell.Direction.Add(
				desired.Sub(cell.Direction).Scale(TurnSpeed), // smooth turning towards desired dir
			).Normalize()
			cell.Position = cell.Position.Add(
				cell.Direction.Scale(
					w.calculateSpeed(cell.Mass) * speedFactor, // velocity in the direction of cell
				),
			)
		}
	}

	for _, e := range w.Ejects {
		if e.Momentum.MagnitudeSq() <= 0.0 {
			continue // stationary mass
		}

		e.Position = e.Position.Add(e.Momentum)
		e.Momentum = e.Momentum.Scale(EjectMomentumDecay)
		if e.Momentum.MagnitudeSq() < 0.01 {
			e.Momentum = Vec2{} // zero it
		}
	}
}

/** World.decayMass
 *
 * Loops through all the players, their cells, and decays their mass
 *
 * Larger cells lose mass faster
 *
 * Formula:
 *  loss = r * m ^ exp
 *
 * Where
 *  loss: resulting mass loss every second
 *  r:    rate of loss
 *  m:    mass of cell
 *  exp:  how much mass contributes to loss
 */
func (w *World) decayMass() {
	for _, p := range w.Players {
		for _, id := range p.CellIDs {
			cell, ok := w.Cells[id]
			if !ok || cell.Mass <= MinDecayMass {
				continue // cell doesn't exit or reached min decay threshold
			}

			lossPerSec := math.Pow(cell.Mass, MassDecayExponent) * MassDecayRate
			lossPerTick := lossPerSec / float64(TickRate)
			cell.Mass -= lossPerTick
			if cell.Mass < MinDecayMass {
				cell.Mass = MinDecayMass
			}
		}
	}
}

/** World.resolveCollisions
 *
 * Loops through all the players, their cells, and collides
 *  them against each other
 *
 */
func (w *World) resolveCollisions() {
	for _, p := range w.Players {
		for i, idA := range p.CellIDs {
			cellA, ok := w.Cells[idA]
			if !ok {
				continue // cell doesn't exist
			}

			for _, idB := range p.CellIDs[i+1:] {
				cellB, ok := w.Cells[idB]
				if !ok {
					continue // cell doesn't exist
				}
				if cellA.MergeTimer == 0.0 && cellB.MergeTimer == 0.0 {
					continue // both cells can merge, let them overlap
				}
				if cellA.Momentum.MagnitudeSq() > MomentumThreshold*MomentumThreshold ||
					cellB.Momentum.MagnitudeSq() > MomentumThreshold*MomentumThreshold {
					continue // split has just happened, don't push apart initially
				}

				diff := cellB.Position.Sub(cellA.Position)
				dist := diff.Magnitude()
				minDist := cellA.Radius() + cellB.Radius()
				if dist < minDist && dist > 0.0 {
					// Push apart
					overlap := minDist - dist
					push := diff.Normalize().Scale(overlap * 0.5)

					cellA.Position = cellA.Position.Sub(push)
					cellB.Position = cellB.Position.Add(push)
				}
			}
		}
	}
}

/** World.clampToWorldBounds
 *
 * Loops through all the players, their cells, and clamps
 *  them against the world border
 *
 */
func (w *World) clampToWorldBounds() {
	for _, p := range w.Players {
		for _, id := range p.CellIDs {
			cell, ok := w.Cells[id]
			if !ok {
				continue // cell doesn't exist
			}

			cell.Position.X = Clamp(cell.Position.X, 0.0, WorldWidth)
			cell.Position.Y = Clamp(cell.Position.Y, 0.0, WorldHeight)
		}
	}
}

/** World.decrementMergeTimer
 *
 * Loops through all the players, their cells, and decrements
 *  the cells' merge timer
 *
 */
func (w *World) decrementMergeTimer() {
	for _, p := range w.Players {
		for _, id := range p.CellIDs {
			cell, ok := w.Cells[id]
			if !ok {
				continue // cell doesn't exist
			}

			cell.MergeTimer = max(0.0, cell.MergeTimer-1.0/float64(TickRate))
		}
	}
}

/** World.recombineCells
 *
 * Loops through all the players, their cells, and checks the merge timer
 *
 * If the merge timer is zero for two cells and they overlap sufficiently,
 *  the larger will absorb the smaller cell (order undefined if same mass)
 *
 */
func (w *World) recombineCells() {
	toRemove := map[uuid.UUID]bool{}
	for _, p := range w.Players {
		for i, idA := range p.CellIDs {
			if toRemove[idA] {
				continue // cell will be removed
			}
			cellA, ok := w.Cells[idA]
			if !ok || cellA.MergeTimer != 0.0 {
				continue // cell doesn't exist or can't merge
			}

			for _, idB := range p.CellIDs[i+1:] {
				if toRemove[idB] {
					continue // cell will be removed
				}
				cellB, ok := w.Cells[idB]
				if !ok || cellB.MergeTimer != 0.0 {
					continue // cell doesn't exist or can't merge
				}

				dist := cellA.Position.DistanceTo(cellB.Position)
				if dist < max(cellA.Radius(), cellB.Radius())*MergeMinOverlap {
					// Merge
					if cellA.Mass >= cellB.Mass {
						cellA.Mass += cellB.Mass
						cellA.MergeTimer += MergeCooldownSeconds
						toRemove[idB] = true
					} else {
						cellB.Mass += cellA.Mass
						cellB.MergeTimer += MergeCooldownSeconds
						toRemove[idA] = true
						break // cellA is gone, stop
					}
				}
			}
		}
	}

	for id := range toRemove {
		w.removeCell(id)
	}
}

/** World.eatPellets
 *
 * Loops through all the players, their cells, and checks
 *  if they eat any pellets in the world
 *
 */
func (w *World) eatPellets() {
	for _, p := range w.Players {
		for _, id := range p.CellIDs {
			cell, ok := w.Cells[id]
			if !ok {
				continue // cell doesn't exist
			}

			for _, pellet := range w.Pellets {
				if cell.Position.DistanceTo(pellet.Position) < cell.Radius() {
					// Eat the pellet
					cell.Mass += pellet.Mass
					delete(w.Pellets, pellet.ID)
				}
			}
		}
	}
}

/** World.eatEjects
 *
 * Loops through all the players, their cells, and checks
 *  if they eat any ejected mass in the world
 *
 */
func (w *World) eatEjects() {
	for _, p := range w.Players {
		for _, id := range p.CellIDs {
			cell, ok := w.Cells[id]
			if !ok {
				continue // cell doesn't exist
			}

			for _, eject := range w.Ejects {
				if cell.Position.DistanceTo(eject.Position) < cell.Radius() {
					// Eat the eject
					cell.Mass += eject.Mass
					delete(w.Ejects, eject.ID)
				}
			}
		}
	}
}

/** World.eatPlayers
 *
 * Loops through all the players, their cells, and checks
 *  if they eat any other cells in the world
 *
 */
func (w *World) eatPlayers() {
	deadPlayers := []uuid.UUID{}
	for _, p := range w.Players {
		for _, id := range p.CellIDs {
			cell, ok := w.Cells[id]
			if !ok {
				continue // cell doesn't exist
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
				victim := w.Players[otherCell.OwnerID]

				cell.Mass += otherCell.Mass
				w.removeCell(otherCell.ID)

				if len(victim.CellIDs) == 0 {
					deadPlayers = append(deadPlayers, victim.ID)
				}
			}
		}
	}

	// Clean up the dead
	for _, id := range deadPlayers {
		w.removePlayer(id)
	}
}

/** World.eatVirus
 *
 * TODO: implement
 *
 */
func (w *World) eatVirus() {

}

/** World.spawnPellets
 *
 * Spawns pellets randomly around the world,
 *  keeping in mind the current amount, maximum amount
 *  of pellets in the world
 *
 */
func (w *World) spawnPellets() {
	currNumPellets := int64(len(w.Pellets))
	if currNumPellets < MinPellets {
		// Spawn more
		additionalPelletNum := RandIntRange(0, MaxPellets-currNumPellets)
		for range additionalPelletNum {
			randMass := RandFloatRange(MinPelletMass, MaxPelletMass)
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

/** World.spawnViruses
 *
 * TODO: implement
 *
 */
func (w *World) spawnViruses() {

}

/** World.removePlayer
 *
 * If the player is still in the world, removes their cells,
 *  sends them their death note, finally removes them
 *  from the world
 *
 */
func (w *World) removePlayer(playerID uuid.UUID) {
	p, ok := w.Players[playerID]
	if !ok {
		return // player already removed
	}

	// Remove player's remaining cells
	for _, id := range p.CellIDs {
		w.removeCell(id)
	}

	// Notify player of his death
	p.OutputChan <- DeathEvent{}
	delete(w.Players, playerID)

	log.Printf("player %v removed\n", playerID)
}

/** World.removeCell
 *
 * Removes the cell from the world, and then from
 *  the player
 *
 */
func (w *World) removeCell(cellID uuid.UUID) {
	c, ok := w.Cells[cellID]
	if !ok {
		return
	}

	// Remove from world
	delete(w.Cells, c.ID)

	p, ok := w.Players[c.OwnerID]
	if !ok {
		return
	}

	// Remove from player
	for i, id := range p.CellIDs {
		if id == cellID {
			p.CellIDs[i] = p.CellIDs[len(p.CellIDs)-1]
			p.CellIDs = p.CellIDs[:len(p.CellIDs)-1]
			return
		}
	}
}

/** World.broadcastState
 *
 * Sends a current snapshot of the world to all clients
 *
 * TODO: each player should get a clipped snapshot,
 *  only containing the necessary state in their vicinity
 *
 */
func (w *World) broadcastState() {
	cells := make([]CellView, 0, len(w.Cells))
	pellets := make([]Pellet, 0, len(w.Pellets))
	ejects := make([]Eject, 0, len(w.Ejects))
	viruses := make([]Virus, 0, len(w.Viruses))

	for _, c := range w.Cells {
		cells = append(cells, CellView{
			ID:      c.ID,
			OwnerID: c.OwnerID,
			Name:    w.Players[c.OwnerID].Name,

			X:      c.Position.X,
			Y:      c.Position.Y,
			Radius: c.Radius(),

			Mass:  c.Mass,
			Color: c.Color,
		})
	}
	for _, p := range w.Pellets {
		pellets = append(pellets, *p)
	}
	for _, e := range w.Ejects {
		ejects = append(ejects, *e)
	}
	for _, v := range w.Viruses {
		viruses = append(viruses, *v)
	}

	for _, p := range w.Players {
		snapshot := TickSnapshot{
			Me:      p.ID,
			Cells:   cells,
			Pellets: pellets,
			Ejects:  ejects,
			Viruses: viruses,
			Score:   p.Score,
			Tick:    w.tick,
		}

		p.OutputChan <- snapshot
	}
}
