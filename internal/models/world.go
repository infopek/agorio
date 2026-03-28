package models

import (
	"sync"

	"github.com/google/uuid"

	"github.com/infopek/agorio/internal/constants"
	"github.com/infopek/agorio/internal/math"
	"github.com/infopek/agorio/internal/types"
)

type WorldState struct {
	Players []Player `json:"players"`
	Viruses []Virus  `json:"viruses"`
	Pellets []Pellet `json:"pellets"`
}

type WorldConfig struct {
	Width  types.Real
	Height types.Real
}

type World struct {
	config  WorldConfig
	players map[uuid.UUID]*Player
	pellets []Pellet
	viruses []Virus
	physics math.PhysicsEngine
	mu      sync.RWMutex
}

func NewWorld(config WorldConfig) *World {
	return &World{
		config:  config,
		players: make(map[uuid.UUID]*Player, constants.MaxPlayers),
		pellets: make([]Pellet, constants.PelletStartCount),
		viruses: make([]Virus, constants.VirusStartCount),
		physics: math.NewPhysicsEngine(),
	}
}

func (w *World) AddPlayer(name string) *Player {
	w.mu.Lock()
	defer w.mu.Unlock()

	player := &Player{
		ID:   uuid.New(),
		Name: name,
		Cells: []Cell{
			Cell{
				Position: math.RandVector2(0, int(w.config.Width), 0, int(w.config.Height)),
				Mass:     constants.PlayerStartMass,
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

func (w *World) Update() error {
	return nil
}

func (w *World) UpdateTarget(playerID uuid.UUID, x, y int) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if p, ok := w.players[playerID]; ok {
		p.Target.X = types.Real(x)
		p.Target.Y = types.Real(y)
	}
}

func (w World) Snapshot() WorldState {
	w.mu.RLock()
	defer w.mu.RUnlock()

	players := make([]Player, 0, len(w.players))
	for _, p := range w.players {
		players = append(players, *p)
	}

	pellets := make([]Pellet, len(w.pellets))
	for i, p := range w.pellets {
		pellets[i] = p
	}

	viruses := make([]Virus, len(w.viruses))
	for i, v := range w.viruses {
		viruses[i] = v
	}

	return WorldState{
		Players: players,
		Pellets: pellets,
		Viruses: viruses,
	}
}
