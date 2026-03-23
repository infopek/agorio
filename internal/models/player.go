package models

import (
	"github.com/google/uuid"

	"github.com/infopek/agorio/internal/math"
)

const (
	minMass int = 1
	maxMass int = 60

	maxCells int = 16
)

type Player struct {
	ID     uuid.UUID
	Name   string
	Mass   int
	Cells  []Cell
	Pos    math.Vector2
	Target math.Vector2
}

func NewPlayer(name string, pos math.Vector2) Player {
	return Player{
		ID:    uuid.New(),
		Name:  name,
		Mass:  math.RandRange(minMass, maxMass),
		Cells: make([]Cell, maxCells),
		Pos:   pos,
		Target: math.Vector2{
			X: 0,
			Y: 0,
		},
	}
}
