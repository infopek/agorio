package game

import (
	"time"

	"github.com/google/uuid"
)

/** World
 *
 * Represents the entire Agor world in all its glory
 *
 */
type World struct {
	Cells   map[uuid.UUID]*Cell
	Pellets map[uuid.UUID]*Pellet
	Ejects  map[uuid.UUID]*Eject
	Viruses map[uuid.UUID]*Virus

	cellGrid   *SpatialGrid[*Cell]
	pelletGrid *SpatialGrid[*Pellet]
	ejectGrid  *SpatialGrid[*Eject]
	virusGrid  *SpatialGrid[*Virus]

	Players map[uuid.UUID]*Player
	Bots    map[uuid.UUID]*Bot

	tick uint64

	InputChan chan PlayerEvent // shared across all players
}

func NewWorld() *World {
	cellGridSize := Radius(MaxCellMass) * 2.0
	pelletGridSize := Radius(MaxPelletMass) * 2.0
	ejectGridSize := Radius(EjectAmount) * 2.0
	virusGridSize := Radius(VirusStartMass) * 2.0
	return &World{
		Cells:   make(map[uuid.UUID]*Cell),
		Pellets: make(map[uuid.UUID]*Pellet),
		Ejects:  make(map[uuid.UUID]*Eject),
		Viruses: make(map[uuid.UUID]*Virus),

		cellGrid:   NewSpatialGrid[*Cell](cellGridSize),
		pelletGrid: NewSpatialGrid[*Pellet](pelletGridSize),
		ejectGrid:  NewSpatialGrid[*Eject](ejectGridSize),
		virusGrid:  NewSpatialGrid[*Virus](virusGridSize),

		Players: make(map[uuid.UUID]*Player),
		Bots:    make(map[uuid.UUID]*Bot),

		InputChan: make(chan PlayerEvent, 256),
	}
}

/** World.Tick
 *
 * Basically the main game loop
 *
 */
func (w *World) Tick() {
	ticker := time.NewTicker(time.Second / time.Duration(TickRate))
	defer ticker.Stop()

	for range ticker.C {
		w.tick++

		w.maintainBots()
		w.updateBots()

		w.processInputs()

		w.rebuildGrids()

		w.moveCells()
		w.feedViruses()
		w.decayMass()
		w.resolveCollisions()
		w.clampToWorldBounds()

		w.decrementMergeTimer()
		w.recombineCells()

		w.eatPellets()
		w.eatEjects()
		w.eatPlayers()
		w.eatVirus()

		w.autoSplit()

		w.spawnPellets()
		w.spawnViruses()

		w.rebuildGrids() // sync updates before sending anything
		w.broadcastState()

		if w.tick%TickRate == 0 {
			w.broadcastLeaderboard()
		}
	}
}

/** World.findEmptySpace
 *
 * Finds a place for a new player to spawn in
 *
 * TODO: find an empty place in the map
 *
 */
func (w *World) findEmptySpace(radius float64) Vec2 {
	for range EmptySpaceMaxAttempts {
		pos := Vec2{
			X: RandFloatRange(radius, WorldWidth-radius),
			Y: RandFloatRange(radius, WorldHeight-radius),
		}

		neighbors := w.cellGrid.GetNeighbors(pos, radius+EmptySpaceRadiusQuery)
		occupied := false
		for _, c := range neighbors {
			if pos.DistanceTo(c.Position) < c.Radius()+EmptySpaceRadiusLeeway {
				occupied = true
				break
			}
		}

		if !occupied {
			return pos
		}
	}

	// We couldn't find unoccupied space, return random
	return Vec2{
		X: RandFloatRange(radius, WorldWidth-radius),
		Y: RandFloatRange(radius, WorldHeight-radius),
	}
}

/** World.getRandomColor
 *
 * TODO: ideally this would generate as few
 *  color collisions as possible
 *
 */
func (w *World) getRandomColor() [3]uint8 {
	return DefaultColors[RandIntRange(0, int64(len(DefaultColors)))]
}
