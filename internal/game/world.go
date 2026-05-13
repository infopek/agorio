package game

import (
	"sync"
	_"log"

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
		Cells: []Cell{
			Cell{
				Body: physics.NewBody(
					&startingCell,
					w.getPlayerStartPosition(),
				),
			}, // starter cell
		},
		Target: math.Vector2{
			X: w.config.Width / 2.0,
			Y: w.config.Height / 2.0,
		}, // starter target points to middle of world
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
			dir := math.Sub(p.Target, c.Body.Position)
			if dir.LengthSq() > 0.0 {
				dir = dir.Normalized()

				speed := constants.PlayerBaseSpeed / math.Sqrt(c.Body.Mass)
				c.Body.Velocity = math.Mul(dir, speed)
			} else {
				c.Body.Velocity = math.Vector2{X: 0.0, Y: 0.0}
			}
		}
	}

	// Euler integration
	for _, p := range w.players {
		for _, c := range p.Cells {
			c.Body.Position.Addi(math.Mul(c.Body.Velocity, dt))
			w.clampPosition(&c)
		}
	}
}

func (w *World) UpdateTarget(playerID uuid.UUID, x, y types.Real) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if p, ok := w.players[playerID]; ok {
		p.Target.X = x
		p.Target.Y = y
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

