package game

import (
	"log"
	"sync"

	"github.com/google/uuid"

	"github.com/infopek/agorio/internal/constants"
	"github.com/infopek/agorio/internal/math"
	"github.com/infopek/agorio/internal/physics"
	"github.com/infopek/agorio/internal/types"
)

type WorldState struct {
	Players map[string]Player `json:"players"`
	Viruses []Virus           `json:"viruses"`
	Pellets []Pellet          `json:"pellets"`
}

type WorldConfig struct {
	Width  types.Real
	Height types.Real
}

type World struct {
	config     WorldConfig
	iterations types.Integer // solver iters

	players map[uuid.UUID]*Player
	pellets []Pellet
	viruses []Virus
	bodies  []*physics.Body // all physics bodies

	mu sync.RWMutex
}

func NewWorld(config WorldConfig) *World {
	return &World{
		config:  config,
		players: make(map[uuid.UUID]*Player, constants.MaxPlayers),
		pellets: make([]Pellet, constants.PelletStartCount),
		viruses: make([]Virus, constants.VirusStartCount),
	}
}

func (w *World) AddPlayer(name string) *Player {
	w.mu.Lock()
	defer w.mu.Unlock()

	startingCell := physics.NewCircle(constants.PlayerStartMass)
	player := &Player{
		ID:   uuid.New(),
		Name: name,
		Cells: []*Cell{
			&Cell{
				Body: physics.NewBody(
					&startingCell,
					w.getPlayerStartPosition(),
				),
				Direction: math.Vector2{},
			}, // starter cell
		},
	}
	w.players[player.ID] = player

	return player
}

func (w *World) RemovePlayer(playerID uuid.UUID) {
	w.mu.Lock()
	defer w.mu.Unlock()

	delete(w.players, playerID)
}

func (w *World) Update(dt types.Real) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Apply input
	for _, p := range w.players {
		for _, c := range p.Cells {
			if c.Direction.LengthSq() > 0.0 {
				speed := constants.PlayerBaseSpeed / math.Sqrt(c.Body.Mass)
				c.Body.Velocity = math.Mul(c.Direction, speed)
			} else {
				c.Body.Velocity = math.Vector2{X: 0.0, Y: 0.0}
			}
			c.Body.Velocity = math.Vector2{X: 100.0, Y: 0.0}
		}
	}

	// Euler integration
	for _, p := range w.players {
		for _, c := range p.Cells {
			c.Body.Position.Addi(math.Mul(c.Body.Velocity, dt))
			w.clampPosition(c)
		}
	}
}

func (w *World) UpdateDirection(playerID uuid.UUID, x, y types.Real) {
	w.mu.Lock()
	defer w.mu.Unlock()

	p, ok := w.players[playerID]
	if !ok {
		return
	}

	targetPos := math.Vector2{X: x, Y: y}
	for _, cell := range p.Cells {
		dir := math.Sub(targetPos, cell.Body.Position)
		log.Printf("Cell pos: (%.1f, %.1f) | Target world: (%.1f, %.1f) | Delta: (%.3f, %.3f)",
			cell.Body.Position.X, cell.Body.Position.Y,
			targetPos.X, targetPos.Y,
			dir.X, dir.Y)
		cell.Direction = dir.Normalized()
		log.Printf("UpdateDirection: cell %p, direction set to %v", cell, cell.Direction)
	}
}

func (w *World) Snapshot() WorldState {
	w.mu.RLock()
	defer w.mu.RUnlock()

	players := make(map[string]Player, len(w.players))
	for _, p := range w.players {
		players[p.ID.String()] = *p
	}

	pellets := make([]Pellet, len(w.pellets))
	copy(pellets, w.pellets)

	viruses := make([]Virus, len(w.viruses))
	copy(viruses, w.viruses)

	return WorldState{
		Players: players,
		Pellets: pellets,
		Viruses: viruses,
	}
}

func (w *World) getPlayerStartPosition() math.Vector2 {
	// TODO: search for empty space for starting pos
	return math.RandVector2(
		0,
		int(w.config.Width),
		0,
		int(w.config.Height),
	)
}

func (w *World) clampPosition(c *Cell) {
	// Make sure cell is not outside world borders
}
