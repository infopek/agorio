package game

import (
	"log"
	"math"

	"github.com/google/uuid"
)

/** World.rebuildGrids
 *
 * Each tick, the spatial grids need to be rebuilt for every
 *  physical entity
 *
 */
func (w *World) rebuildGrids() {
	w.cellGrid.Clear()
	for _, c := range w.Cells {
		w.cellGrid.Insert(c.Position, c)
	}

	w.pelletGrid.Clear()
	for _, p := range w.Pellets {
		w.pelletGrid.Insert(p.Position, p)
	}

	w.ejectGrid.Clear()
	for _, e := range w.Ejects {
		w.ejectGrid.Insert(e.Position, e)
	}

	w.virusGrid.Clear()
	for _, v := range w.Viruses {
		w.virusGrid.Insert(v.Position, v)
	}
}

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

			// Speed scales with distance, same with the slowdown distance from cursor
			slowdownDist := cell.Radius() * (SpeedFactorRange / math.Sqrt(cell.Mass/StartMass))
			speedFactor := min(dist/slowdownDist, 1.0)

			// Calculate new direction and velocity
			speed := BaseSpeed * math.Pow(StartMass/cell.Mass, SpeedExponent)
			speed = math.Max(speed, MinSpeed) * speedFactor
			desired := diff.Scale(1.0 / dist) // normalize

			cell.Direction = cell.Direction.Add(
				desired.Sub(cell.Direction).Scale(TurnSpeed), // smooth turning towards desired dir
			).Normalize()
			cell.Position = cell.Position.Add(
				cell.Direction.Scale(speed),
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

	for _, v := range w.Viruses {
		if v.Momentum.MagnitudeSq() <= 0.0 {
			continue // stationary virus
		}

		v.Position = v.Position.Add(v.Momentum)
		v.Momentum = v.Momentum.Scale(VirusMomentumDecay)
		if v.Momentum.MagnitudeSq() < 0.01 {
			v.Momentum = Vec2{} // zero it
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
	for range CollisionResolutionPasses {
		for _, p := range w.Players {
			for _, idA := range p.CellIDs {
				cellA, ok := w.Cells[idA]
				if !ok {
					continue // cell doesn't exist
				}

				neighbors := w.cellGrid.GetNeighbors(cellA.Position, cellA.Radius())
				for _, cellB := range neighbors {
					if cellA.ID == cellB.ID || cellB.OwnerID != cellA.OwnerID {
						continue // skip self and other players' cells
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
						// Push apart (mass relative)
						overlap := minDist - dist
						totalMass := cellA.Mass + cellB.Mass

						pushA := overlap * (cellB.Mass / totalMass)
						pushB := overlap * (cellA.Mass / totalMass)

						dir := diff.Normalize()
						cellA.Position = cellA.Position.Sub(dir.Scale(pushA))
						cellB.Position = cellB.Position.Add(dir.Scale(pushB))
					}
				}
			}
		}
	}

	// Ejects can collide with themselves
	for _, ejectA := range w.Ejects {
		neighbors := w.ejectGrid.GetNeighbors(ejectA.Position, ejectA.Radius())
		for _, ejectB := range neighbors {
			diff := ejectB.Position.Sub(ejectA.Position)
			dist := diff.Magnitude()
			minDist := ejectA.Radius() + ejectB.Radius()

			if dist < minDist && dist > 0.0 {
				overlap := minDist - dist
				dir := diff.Normalize()
				push := dir.Scale(overlap * 0.5)

				ejectA.Position = ejectA.Position.Sub(push)
				ejectB.Position = ejectB.Position.Add(push)
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
	// Cells
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

	// Pellets
	for _, e := range w.Ejects {
		e.Position.X = Clamp(e.Position.X, 0.0, WorldWidth)
		e.Position.Y = Clamp(e.Position.Y, 0.0, WorldHeight)
	}

	// Viruses (they be bouncing)
	for _, v := range w.Viruses {
		if v.Position.X < v.Radius() {
			v.Position.X = v.Radius()
			v.Momentum.X = -v.Momentum.X
		}
		if v.Position.X > WorldWidth-v.Radius() {
			v.Position.X = WorldWidth - v.Radius()
			v.Momentum.X = -v.Momentum.X
		}
		if v.Position.Y < v.Radius() {
			v.Position.Y = v.Radius()
			v.Momentum.Y = -v.Momentum.Y
		}
		if v.Position.Y > WorldHeight-v.Radius() {
			v.Position.Y = WorldHeight - v.Radius()
			v.Momentum.Y = -v.Momentum.Y
		}
	}
}

/*
  - World.bounceEntity
    *
    *

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

				var larger, smaller *Cell
				if cellA.Mass >= cellB.Mass {
					larger = cellA
					smaller = cellB
				} else {
					larger = cellB
					smaller = cellA
				}

				dist := cellA.Position.DistanceTo(cellB.Position)
				if dist+smaller.Radius()*MergeMinOverlap < larger.Radius() {
					// Merge
					larger.Mass += smaller.Mass
					larger.MergeTimer = MergeCooldownSeconds
					toRemove[smaller.ID] = true
					if smaller == cellA {
						break
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
 * Loops through all the players, their cells, and checks
 *  if they were popped by a virus
 *
 */
func (w *World) eatVirus() {
	for _, p := range w.Players {
		for _, id := range p.CellIDs {
			cell, ok := w.Cells[id]
			if !ok {
				continue // cell doesn't exist
			}

			for _, v := range w.Viruses {
				if float64(cell.Mass) <= float64(v.Mass)*EatMassThreshold {
					continue // we are not big enough
				}

				if cell.Position.DistanceTo(v.Position) >=
					(cell.Radius() - v.Radius()*EatDistanceThreshold) {
					continue // we are not overlapping enough
				}

				// Eat virus
				cell.Mass += v.Mass
				delete(w.Viruses, v.ID)

				w.popCell(cell)
			}
		}
	}
}

func (w *World) popCell(cell *Cell) {
	p, ok := w.Players[cell.OwnerID]
	if !ok {
		return // cell doesn't belong to anyone
	}

	available := SplitMaxCells - int64(len(p.CellIDs))
	if available <= 0 {
		return // can't split more
	}

	totalMass := cell.Mass
	maxPieces := int64(min(float64(available), math.Floor(totalMass/SplitMinMass)-1))
	massPerPiece := cell.Mass / float64(maxPieces+1)

	// Rule 1: can we split into max pieces evenly?
	if massPerPiece <= SplitMinMass+VirusPopEqualSplitRange {
		// Even split into max pieces
		cell.Mass = massPerPiece
		for range maxPieces {
			w.createPopPiece(cell, p, massPerPiece)
		}
	} else {
		// Recursive halving
		remainingMass := totalMass
		slotsUsed := int64(0)

		for remainingMass > SplitMinMass*2.0 && slotsUsed < available {
			// Take half for a sibling
			siblingMass := remainingMass / 2.0
			remainingMass -= siblingMass

			slotsLeft := available - slotsUsed

			// Can the sibling be evenly split into slotsLeft pieces?
			evenPieces := min(int64(math.Floor(siblingMass/SplitMinMass)), slotsLeft)
			evenMass := siblingMass / float64(evenPieces)

			if evenMass <= SplitMinMass+VirusPopEqualSplitRange || slotsLeft <= 2 {
				// Force even split, use all slots
				evenMass = siblingMass / float64(slotsLeft)
				for range slotsLeft {
					w.createPopPiece(cell, p, evenMass)
					slotsUsed++
				}
				break
			} else {
				// Create one sibling with this mass, continue halving
				w.createPopPiece(cell, p, siblingMass)
				slotsUsed++
			}
		}

		cell.Mass = remainingMass
	}

	cell.MergeTimer = MergeTimerStartSeconds
}

/** World.createPopPiece
 *
 * Helper function for creating a cell in the world for
 *  player p when popped by a virus
 *
 * Uses the parent cell (pc) a lot
 *
 */
func (w *World) createPopPiece(pc *Cell, p *Player, mass float64) {
	angle := RandFloatRange(0.0, 2.0*math.Pi)
	dir := Vec2{
		X: math.Cos(angle),
		Y: math.Sin(angle),
	}

	newCell := Cell{
		ID:      uuid.New(),
		OwnerID: pc.OwnerID,

		Position:  pc.Position.Add(dir.Scale(pc.Radius() * 0.3)),
		Direction: dir,
		Momentum: dir.Scale(SplitMomentumFactor*VirusPopDampenFactor +
			math.Sqrt(Radius(mass))*SplitMomentumRadiusFactor), // bit weird to use Radius() here

		MergeTimer: MergeTimerStartSeconds,
		Mass:       mass,
		Color:      pc.Color,
	}

	p.CellIDs = append(p.CellIDs, newCell.ID)
	w.Cells[newCell.ID] = &newCell
}

/** World.feedViruses
 *
 * When a pellet hits a virus from a direction, the virus absorbs it,
 *  and if it reaches its limit, it shoots another virus from itself in the direction
 *  the pellet came from
 */
func (w *World) feedViruses() {
	toRemove := map[uuid.UUID]bool{}

	for _, e := range w.Ejects {
		if toRemove[e.ID] {
			continue // pellet already absorbed by a virus
		}
		neighbors := w.virusGrid.GetNeighbors(e.Position, e.Radius())
		for _, v := range neighbors {
			dist := e.Position.DistanceTo(v.Position)
			if dist < v.Radius() {
				v.Mass += e.Mass
				v.FedCount++
				toRemove[e.ID] = true

				var dir Vec2
				if e.Momentum.MagnitudeSq() < 0.01 {
					dir = v.Position.Sub(e.Position).Normalize()
				} else {
					dir = e.Momentum.Normalize()
				}

				if v.FedCount >= VirusFeedToShoot {
					w.shootVirus(v, dir)
				}
			}
		}
	}

	for id := range toRemove {
		delete(w.Ejects, id)
	}
}

/** world.shootVirus
 *
 * Helper for World.feedViruses
 *
 * Spawns a new virus from parent virus (pv) in
 *  the direction specified
 *
 */
func (w *World) shootVirus(pv *Virus, dir Vec2) {
	newVirus := Virus{
		ID: uuid.New(),

		Position: pv.Position,
		Momentum: dir.Scale(VirusLaunchSpeed),

		Mass:  VirusStartMass,
		Color: pv.Color,
	}

	w.Viruses[newVirus.ID] = &newVirus

	// Reset parent
	pv.FedCount = 0
	pv.Mass = VirusStartMass
}

/** World.spawnPellets
 *
 * Spawns pellets randomly around the world,
 *  keeping in mind the current amount, maximum amount
 *  of pellets in the world
 *
 */
func (w *World) spawnPellets() {
	if int64(len(w.Pellets)) >= MaxPellets {
		return // we have enough already
	}

	if RandFloatRange(0, 1) >= PelletSpawnChance {
		return // unlucky
	}

	// Spawn more
	for range PelletSpawnAmount {
		pellet := Pellet{
			ID: uuid.New(),

			Position: w.findEmptySpace(),

			Mass:  RandFloatRange(MinPelletMass, MaxPelletMass),
			Color: w.getRandomColor(),
		}

		w.Pellets[pellet.ID] = &pellet
	}
}

/** World.spawnViruses
 *
 * Spawns viruses randomly around the world,
 *  keeping in mind the current amount, maximum amount
 *  of viruses in the world
 *
 * Finds a spot which is not occupied by a player for
 *  the new ones
 *
 */
func (w *World) spawnViruses() {
	currNumViruses := int64(len(w.Viruses))
	if currNumViruses >= MinViruses {
		return // don't spawn more naturally
	}

	if RandFloatRange(0, 1) >= VirusSpawnChance {
		return // unlucky
	}

	// Spawn more
	for range VirusSpawnAmount {
		virus := Virus{
			ID: uuid.New(),

			Position: w.findEmptySpace(),

			Mass:  VirusStartMass,
			Color: VirusColor,
		}

		w.Viruses[virus.ID] = &virus
	}
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
	for _, p := range w.Players {
		center, totalMass := w.playerCenter(p)

		halfW := 1000.0 + math.Sqrt(totalMass)*10.0
		halfH := 700.0 + math.Sqrt(totalMass)*10.0

		minX, minY := center.X-halfW, center.Y-halfH
		maxX, maxY := center.X+halfW, center.Y+halfH

		cells := w.cellGrid.GetInRect(minX, minY, maxX, maxY)
		pellets := w.pelletGrid.GetInRect(minX, minY, maxX, maxY)
		ejects := w.ejectGrid.GetInRect(minX, minY, maxX, maxY)
		viruses := w.virusGrid.GetInRect(minX, minY, maxX, maxY)

		cellViews := make([]CellView, 0, len(cells))
		pelletViews := make([]Pellet, 0, len(pellets))
		ejectViews := make([]Eject, 0, len(ejects))
		virusViews := make([]Virus, 0, len(viruses))

		for _, c := range cells {
			cellViews = append(cellViews, CellView{
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
		for _, p := range pellets {
			pelletViews = append(pelletViews, *p)
		}
		for _, e := range ejects {
			ejectViews = append(ejectViews, *e)
		}
		for _, v := range viruses {
			virusViews = append(virusViews, *v)
		}

		snapshot := TickSnapshot{
			Me:      p.ID,
			Cells:   cellViews,
			Pellets: pelletViews,
			Ejects:  ejectViews,
			Viruses: virusViews,
			Score:   p.Score,
			Tick:    w.tick,
		}

		p.OutputChan <- snapshot
	}
}

/** World.playerCenter
 *
 * Helper method for World.broadcastState
 *
 * Finds the weighted center of the player, and return the position
 *  with the totalMass
 *
 */
func (w *World) playerCenter(p *Player) (Vec2, float64) {
	if len(p.CellIDs) == 0 {
		return Vec2{}, 0.0
	}

	totalMass := 0.0
	center := Vec2{}
	for _, id := range p.CellIDs {
		c := w.Cells[id]
		center.X += c.Position.X * c.Mass
		center.Y += c.Position.Y * c.Mass
		totalMass += c.Mass
	}

	center.X /= totalMass
	center.Y /= totalMass
	return center, totalMass
}
