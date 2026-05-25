package game

import (
	_ "log"
	"math"
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

	tick uint64

	InputChan chan PlayerEvent // shared across all players
}

func NewWorld() *World {
	return &World{
		Cells:     make(map[uuid.UUID]*Cell),
		Pellets:   make(map[uuid.UUID]*Pellet),
		Ejects:    make(map[uuid.UUID]*Eject),
		Viruses:   make(map[uuid.UUID]*Virus),
		Players:   make(map[uuid.UUID]*Player),
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

		w.applyPhysics()
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

/** World.findStartPosition
 *
 * Finds a place for a new player to spawn in
 *
 * TODO: find an empty place in the map
 *
 */
func (w *World) findStartPosition() Vec2 {
	return Vec2{
		X: RandFloatRange(0.0, WorldWidth),
		Y: RandFloatRange(0.0, WorldHeight),
	}
}

/** World.findStartColor
 *
 * Finds a color for a new player
 *
 * TODO: ideally this would generate as few
 *  color collisions as possible
 *
 */
func (w *World) findStartColor() [3]uint8 {
	return getRandomColor()
}

/** World.calculateSpeed
 *
 * Calculates the speed of the cell given its mass
 *
 * Formula:
 *  sp = base_sp * (sm / m) ^ exp
 *
 * Where
 *  sp:      resulting speed
 *  base_sp: base speed
 *  sm:      starting mass of a new player
 *  m:       current mass of cell
 *  exp:     rate of speed penalty per unit of growth
 *
 */
func (w *World) calculateSpeed(mass float64) float64 {
	speed := BaseSpeed * math.Pow(StartMass/mass, SpeedExponent)
	return math.Max(speed, MinSpeed)
}

/** getRandomColor
 *
 * Util function for a random color
 *
 */
func getRandomColor() [3]uint8 {
	return DefaultColors[RandIntRange(0, int64(len(DefaultColors)))]
}
