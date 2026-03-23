package models

import (
	"agorio/internal/math"

	"github.com/google/uuid"
)

const (
	minMass int = 1
	maxMass int = 60

	maxCells int = 16
)

type Player struct {
	ID        uuid.UUID
	Name      string
	Mass      int
	Cells     []Cell
	Pos       math.Vector2
	TargetPos math.Vector2
}

func NewPlayer(name string, pos math.Vector2) Player {
	return Player{
		ID:    uuid.New(),
		Name:  name,
		Mass:  math.RandRange(minMass, maxMass),
		Cells: make([]Cell, maxCells),
		Pos:   pos,
	}
}
