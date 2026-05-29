package game

import (
	_ "log"
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
	Players map[uuid.UUID]*Player

	cellGrid   *SpatialGrid[*Cell]
	pelletGrid *SpatialGrid[*Pellet]
	ejectGrid  *SpatialGrid[*Eject]
	virusGrid  *SpatialGrid[*Virus]

	tick uint64

	InputChan chan PlayerEvent // shared across all players
}

func NewWorld() *World {
	cellGridSize := Radius(ForceSplitMassThreshold) * 2.0
	pelletGridSize := Radius(MaxPelletMass) * 2.0
	ejectGridSize := Radius(EjectAmount) * 2.0
	virusGridSize := Radius(VirusStartMass) * 2.0
	return &World{
		Cells:   make(map[uuid.UUID]*Cell),
		Pellets: make(map[uuid.UUID]*Pellet),
		Ejects:  make(map[uuid.UUID]*Eject),
		Viruses: make(map[uuid.UUID]*Virus),
		Players: make(map[uuid.UUID]*Player),

		cellGrid:   NewSpatialGrid[*Cell](cellGridSize),
		pelletGrid: NewSpatialGrid[*Pellet](pelletGridSize),
		ejectGrid:  NewSpatialGrid[*Eject](ejectGridSize),
		virusGrid:  NewSpatialGrid[*Virus](virusGridSize),

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

		w.processInputs()

		w.rebuildGrids()

		w.applyPhysics()
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

		w.spawnPellets()
		w.spawnViruses()

		w.broadcastState()
	}
}

/** World.findEmptySpace
 *
 * Finds a place for a new player to spawn in
 *
 * TODO: find an empty place in the map
 *
 */
func (w *World) findEmptySpace() Vec2 {
	return Vec2{
		X: RandFloatRange(0.0, WorldWidth),
		Y: RandFloatRange(0.0, WorldHeight),
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
