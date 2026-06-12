package game

import (
	_ "log"
	"math"
	"sort"

	"github.com/google/uuid"
)

/** World.rebuildGrids
 *
 * Each tick, the spatial grids need to be rebuilt for every
 *  physical entity
 *
 */
func (w *World) rebuildGrids() {
	//log.Printf("rebuildGrids")
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

/** World.moveCells
 *
 * Loops through all the players, their cells, and applies
 *  potential split momentum, velocity physics
 *
 * Applies momentum to all ejected mass
 *
 */
func (w *World) moveCells() {
	// Cells
	//log.Printf("moveCells")
	for _, p := range w.Players {
		for _, id := range p.CellIDs {
			cell := w.Cells[id]

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
			slowdownDist := cell.Radius() * SpeedFactorRange
			speedFactor := min(dist/slowdownDist, 1.0)

			// Calculate new direction and velocity
			speed := BaseSpeed * math.Pow(StartMass/cell.Mass, SpeedExponent)
			speed = math.Max(speed, MinSpeed) * speedFactor
			desired := diff.Scale(1.0 / dist) // normalize
			actual := desired.Sub(cell.Direction)

			// Smooth turning towards desired dir
			cell.Direction = cell.Direction.Add(actual.Scale(TurnSpeed)).Normalize()
			cell.Position = cell.Position.Add(cell.Direction.Scale(speed))
		}
	}

	// Ejects
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

	// Viruses
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
 */
func (w *World) decayMass() {
	//log.Printf("decayMass")
	for _, p := range w.Players {
		for _, id := range p.CellIDs {
			cell := w.Cells[id]
			if cell.Mass <= MinDecayMass {
				continue // cell reached min decay threshold
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
	//log.Printf("resolveCollisions")
	for range CollisionResolutionPasses {
		for _, p := range w.Players {
			for _, idA := range p.CellIDs {
				cellA := w.Cells[idA]

				neighbors := w.cellGrid.GetNeighbors(cellA.Position, cellA.Radius())
				for _, cellB := range neighbors {
					_, ok := w.Cells[cellB.ID]
					if !ok {
						continue // cell doesn't exist
					}
					if cellA.ID == cellB.ID || cellB.OwnerID != cellA.OwnerID {
						continue // skip self and other players' cells
					}
					if cellA.MergeTimer == 0.0 && cellB.MergeTimer == 0.0 &&
						cellA.Mass+cellB.Mass < MaxCellMass {
						continue // both cells can merge without autosplit, let them overlap
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
			_, ok := w.Ejects[ejectB.ID]
			if !ok {
				continue // eject doesn't exist
			}

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
	//log.Printf("clampToWorldBounds")
	// Cells
	for _, p := range w.Players {
		for _, id := range p.CellIDs {
			cell := w.Cells[id]
			cell.Position.X = Clamp(cell.Position.X, 0.0, WorldWidth)
			cell.Position.Y = Clamp(cell.Position.Y, 0.0, WorldHeight)
		}
	}

	// Pellets (bouncy-bounce)
	for _, e := range w.Ejects {
		if e.Position.X < e.Radius() {
			e.Position.X = e.Radius()
			e.Momentum.X = -e.Momentum.X
		}
		if e.Position.X > WorldWidth-e.Radius() {
			e.Position.X = WorldWidth - e.Radius()
			e.Momentum.X = -e.Momentum.X
		}
		if e.Position.Y < e.Radius() {
			e.Position.Y = e.Radius()
			e.Momentum.Y = -e.Momentum.Y
		}
		if e.Position.Y > WorldHeight-e.Radius() {
			e.Position.Y = WorldHeight - e.Radius()
			e.Momentum.Y = -e.Momentum.Y
		}
	}

	// Viruses (they be bouncing too)
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

/** World.decrementMergeTimer
 *
 * Loops through all the players, their cells, and decrements
 *  the cells' merge timer
 *
 */
func (w *World) decrementMergeTimer() {
	//log.Printf("decrementMergeTimer")
	for _, p := range w.Players {
		for _, id := range p.CellIDs {
			cell := w.Cells[id]
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
			cellA := w.Cells[idA]
			if cellA.MergeTimer != 0.0 {
				continue // can't merge
			}

			for _, idB := range p.CellIDs[i+1:] {
				if toRemove[idB] {
					continue // cell will be removed, skip
				}

				cellB := w.Cells[idB]
				if cellB.MergeTimer != 0.0 {
					continue // can't merge
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
						break // reference cell to check against gone, get next
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
	//log.Printf("eatPellets")
	for _, p := range w.Players {
		for _, id := range p.CellIDs {
			cell := w.Cells[id]
			neighbors := w.pelletGrid.GetNeighbors(cell.Position, cell.Radius())
			for _, pellet := range neighbors {
				if _, ok := w.Pellets[pellet.ID]; !ok {
					continue // pellet no longer exists
				}

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
	//log.Printf("eatEjects")
	for _, p := range w.Players {
		for _, id := range p.CellIDs {
			cell := w.Cells[id]
			neighbors := w.ejectGrid.GetNeighbors(cell.Position, cell.Radius())
			for _, eject := range neighbors {
				if _, ok := w.Ejects[eject.ID]; !ok {
					continue // eject no longer exists
				}

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
	//log.Printf("eatPlayers")
	deadPlayers := map[uuid.UUID]bool{}
	for _, p := range w.Players {
		for _, id := range p.CellIDs {
			cellA, ok := w.Cells[id]
			if !ok || deadPlayers[cellA.OwnerID] {
				continue // cell doesn't exist or player already dead
			}

			neighbors := w.cellGrid.GetNeighbors(cellA.Position, cellA.Radius())
			for _, cellB := range neighbors {
				_, ok = w.Cells[cellB.ID]
				if !ok || deadPlayers[cellB.OwnerID] {
					continue // cell doesn't exist or player already dead
				}

				if cellB.OwnerID == p.ID {
					continue // we don't eat our own cells
				}

				if float64(cellA.Mass) <= float64(cellB.Mass)*EatMassThreshold {
					continue // we are not big enough
				}

				if cellA.Position.DistanceTo(cellB.Position) >=
					(cellA.Radius() - cellB.Radius()*EatDistanceThreshold) {
					continue // we are not overlapping enough
				}

				// Eat the cell
				victim := w.Players[cellB.OwnerID]

				cellA.Mass += cellB.Mass
				w.removeCell(cellB.ID)

				if len(victim.CellIDs) == 0 {
					deadPlayers[victim.ID] = true
				}
			}
		}
	}

	// Clean up the dead
	for id := range deadPlayers {
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
	//log.Printf("eatVirus")
	for _, p := range w.Players {
		for _, id := range p.CellIDs {
			cell, ok := w.Cells[id]
			if !ok {
				continue // cell doesn't exist
			}

			neighbors := w.virusGrid.GetNeighbors(cell.Position, cell.Radius())
			for _, virus := range neighbors {
				_, ok := w.Viruses[virus.ID]
				if !ok {
					continue // virus doesn't exist
				}

				if cell.Mass <= virus.Mass*EatMassThreshold {
					continue // we are not big enough
				}

				if cell.Position.DistanceTo(virus.Position) >=
					(cell.Radius() - virus.Radius()*EatDistanceThreshold) {
					continue // we are not overlapping enough
				}

				// Eat virus
				cell.Mass += virus.Mass
				delete(w.Viruses, virus.ID)

				w.popCell(cell)
			}
		}
	}
}

/** World.popCell
 *
 * Preconditions:
 *  - cell exists
 *
 * Helper for popping a cell when it eats a virus
 *
 * Uses the same mechanism as normal agar
 * Thanks: https://www.youtube.com/watch?v=MjlVSVvxADk
 *
 */
func (w *World) popCell(cell *Cell) {
	//log.Printf("popCell")
	p, ok := w.Players[cell.OwnerID]
	if !ok {
		return // cell owner dead / disconnected
	}

	available := SplitMaxCells - uint64(len(p.CellIDs))
	if available <= 0 {
		return // can't split more
	}

	totalMass := cell.Mass
	maxPieces := int64(min(float64(available), math.Floor(totalMass/SplitMinMass)-1))
	massPerPiece := cell.Mass / float64(maxPieces+1)

	// Can we split into max pieces evenly?
	if massPerPiece <= SplitMinMass+VirusPopEqualSplitRange {
		// Even split into max pieces
		cell.Mass = massPerPiece
		for range maxPieces {
			w.createPopPiece(cell, massPerPiece)
		}
	} else {
		// Recursive halving
		remainingMass := totalMass
		slotsUsed := uint64(0)

		for remainingMass > SplitMinMass*2.0 && slotsUsed < available {
			// Take half for a sibling
			siblingMass := remainingMass / 2.0
			remainingMass -= siblingMass

			slotsLeft := available - slotsUsed

			// Can the sibling be evenly split into slotsLeft pieces?
			evenPieces := min(uint64(math.Floor(siblingMass/SplitMinMass)), slotsLeft)
			evenMass := siblingMass / float64(evenPieces)

			if evenMass <= SplitMinMass+VirusPopEqualSplitRange || slotsLeft <= 2 {
				// Force even split, use all slots
				evenMass = siblingMass / float64(slotsLeft)
				for range slotsLeft {
					w.createPopPiece(cell, evenMass)
					slotsUsed++
				}
				break
			} else {
				// Create one sibling with this mass, continue halving
				w.createPopPiece(cell, siblingMass)
				slotsUsed++
			}
		}

		cell.Mass = remainingMass
	}

	cell.MergeTimer = min(MaxMergeTimer, MergeTimerStartSeconds+(cell.Mass*MergeTimerMassFactor))
}

/** World.createPopPiece
 *
 * Preconditions:
 *  - parent cell exists
 *
 * Helper function for creating a cell
 *  from parent cell pc in the world for
 *  when popped by a virus
 *
 * Uses the parent cell (pc) a lot
 *
 */
func (w *World) createPopPiece(pc *Cell, mass float64) {
	//log.Printf("createPopPiece")
	angle := RandFloatRange(0.0, 2.0*math.Pi)
	dir := Vec2{
		X: math.Cos(angle),
		Y: math.Sin(angle),
	}

	p, ok := w.Players[pc.OwnerID]
	if !ok {
		return // owner dead or disconnected
	}

	newCell := Cell{
		ID:      uuid.New(),
		OwnerID: pc.OwnerID,

		Position:  pc.Position.Add(dir.Scale(pc.Radius() * 0.3)),
		Direction: dir,
		// TODO: figure something else out
		Momentum: dir.Scale(SplitMomentumFactor*VirusPopDampenFactor +
			math.Sqrt(Radius(mass))*SplitMomentumRadiusFactor),

		MergeTimer: min(MaxMergeTimer, MergeTimerStartSeconds+(mass*MergeTimerMassFactor)),
		Mass:       mass,
		Color:      pc.Color,
	}

	p.CellIDs = append(p.CellIDs, newCell.ID)
	w.Cells[newCell.ID] = &newCell
}

/** World.autoSplit
 *
 * When a cell hits the maximum allowed mass, we
 *  split it automatically in a random direction
 *
 */
func (w *World) autoSplit() {
	for _, p := range w.Players {
		cellIDs := make([]uuid.UUID, len(p.CellIDs))
		copy(cellIDs, p.CellIDs)

		for _, id := range cellIDs {
			cell := w.Cells[id]
			if cell.Mass < MaxCellMass {
				continue // not at limit yet
			}

			if uint64(len(p.CellIDs)) >= SplitMaxCells {
				cell.Mass = MaxCellMass // cap it, can't split
				continue
			}

			// Original cell gets the extra
			newMass := MaxCellMass / 2.0
			excess := cell.Mass - MaxCellMass
			cell.Mass = newMass + excess
			cell.MergeTimer += MergeCooldownSeconds

			dir := Vec2{
				X: RandFloatRange(0.0, 1.0),
				Y: RandFloatRange(0.0, 1.0),
			}
			newCell := Cell{
				ID:      uuid.New(),
				OwnerID: cell.OwnerID,

				Position:  cell.Position.Add(cell.Direction.Scale(cell.Radius() / 2.0)),
				Direction: dir,
				Momentum:  dir.Scale(SplitMomentumFactor + math.Sqrt(cell.Radius())*SplitMomentumRadiusFactor),

				MergeTimer: min(MaxMergeTimer, MergeTimerStartSeconds+(newMass*MergeTimerMassFactor)),

				Mass:  newMass,
				Color: cell.Color,
			}

			p.CellIDs = append(p.CellIDs, newCell.ID)
			w.Cells[newCell.ID] = &newCell
		}
	}
}

/** World.feedViruses
 *
 * When an eject hits a virus from a direction, the virus absorbs it,
 *  and if it reaches its limit, it shoots another virus from itself
 *  in the opposite direction the eject came from
 *
 */
func (w *World) feedViruses() {
	toRemove := map[uuid.UUID]bool{}
	//log.Printf("feedViruses")

	for _, eject := range w.Ejects {
		if toRemove[eject.ID] {
			continue // eject already absorbed by a virus
		}

		neighbors := w.virusGrid.GetNeighbors(eject.Position, eject.Radius())
		for _, virus := range neighbors {
			_, ok := w.Viruses[virus.ID]
			if !ok {
				continue // virus doesn't exist
			}

			dist := eject.Position.DistanceTo(virus.Position)
			if dist < virus.Radius() {
				virus.Mass += eject.Mass
				virus.FedCount++
				toRemove[eject.ID] = true

				var dir Vec2
				if eject.Momentum.MagnitudeSq() < 0.01 {
					dir = virus.Position.Sub(eject.Position).Normalize()
				} else {
					dir = eject.Momentum.Normalize()
				}

				if virus.FedCount >= VirusFeedToShoot {
					w.shootVirus(virus, dir)
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
 * Preconditions:
 *  - virus exists
 *
 * Helper for World.feedViruses
 *
 * Spawns a new virus from parent virus (pv) in
 *  the direction specified
 *
 */
func (w *World) shootVirus(pv *Virus, dir Vec2) {
	//log.Printf("shootVirus")
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
	//log.Printf("spawnPellets")
	if uint64(len(w.Pellets)) >= MaxPellets {
		return // we have enough already
	}

	if RandFloatRange(0, 1) >= PelletSpawnChance {
		return // unlucky
	}

	// Spawn more
	for range PelletSpawnAmount {
		pellet := Pellet{
			ID: uuid.New(),

			Position: w.findEmptySpace(MaxPelletMass),

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
 */
func (w *World) spawnViruses() {
	//log.Printf("spawnViruses")
	currNumViruses := uint64(len(w.Viruses))
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

			Position: w.findEmptySpace(VirusStartMass),

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
	//log.Printf("removePlayer")
	p, ok := w.Players[playerID]
	if !ok {
		return // player already removed
	}

	// Remove player's remaining cells
	for _, id := range p.CellIDs {
		w.removeCell(id)
	}

	// Notify player of his death
	if !p.IsBot {
		sendServerEvent(p.OutputChan, DeathEvent{})
	}
	delete(w.Players, playerID)
}

/** World.removeCell
 *
 * Removes the cell from the world, and then from
 *  the player
 *
 */
func (w *World) removeCell(cellID uuid.UUID) {
	//log.Printf("removeCell")
	c, ok := w.Cells[cellID]
	if !ok {
		return // cell already removed
	}

	// Remove from world
	delete(w.Cells, c.ID)

	p, ok := w.Players[c.OwnerID]
	if !ok {
		return // player dead or disconnected
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
 */
func (w *World) broadcastState() {
	//log.Printf("broadcastState")
	for _, p := range w.Players {
		if p.IsBot {
			continue // don't broadcast to bots
		}

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

		sendServerEvent(p.OutputChan, snapshot)
	}
}

/** World.broadcastLeaderboard
 *
 * Sends the current leaderboard to all players
 *
 */
func (w *World) broadcastLeaderboard() {
	leaderboard := Leaderboard{
		Entries: make([]LeaderboardEntry, 0, len(w.Players)),
	}
	for id, p := range w.Players {
		leaderboard.Entries = append(leaderboard.Entries, LeaderboardEntry{
			ID:    id,
			Name:  p.Name,
			Score: uint64(w.playerTotalMass(p)),
		})
	}

	sort.Slice(leaderboard.Entries, func(i, j int) bool {
		return leaderboard.Entries[i].Score > leaderboard.Entries[j].Score
	})

	cutoff := min(uint64(len(leaderboard.Entries)), LeaderboardSize)
	leaderboard.Entries = leaderboard.Entries[:cutoff]

	for id, p := range w.Players {
		if p.IsBot {
			continue // don't broadcast to bots
		}

		leaderboard.Me = id

		sendServerEvent(p.OutputChan, leaderboard)
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
func (w *World) playerCenter(player *Player) (Vec2, float64) {
	//log.Printf("playerCenter")
	if len(player.CellIDs) == 0 {
		return Vec2{}, 0.0
	}

	center := Vec2{}
	for _, id := range player.CellIDs {
		c := w.Cells[id]
		center.X += c.Position.X * c.Mass
		center.Y += c.Position.Y * c.Mass
	}

	totalMass := w.playerTotalMass(player)
	center.X /= totalMass
	center.Y /= totalMass
	return center, totalMass
}

/** World.playerTotalMass
 *
 * Helper method for determining the sum of mass
 *  for a player's cells
 *
 */
func (w *World) playerTotalMass(player *Player) float64 {
	if len(player.CellIDs) == 0 {
		return 0.0
	}

	totalMass := 0.0
	for _, id := range player.CellIDs {
		c := w.Cells[id]
		totalMass += c.Mass
	}

	return totalMass
}

/** world.maintainBots
 *
 * Fill the world with bots until the
 *  target player number is met
 *
 */
func (w *World) maintainBots() {
	//log.Printf("maintainBots")
	curr := len(w.Players)
	for range int(MaxPlayers) - curr {
		w.spawnBot()
	}
}
