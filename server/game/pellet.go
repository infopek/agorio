package game

import (
	"github.com/google/uuid"
)

type Pellet struct {
	ID uuid.UUID

	Position Vec2

	Mass  float64
	Color [3]uint8
}

func (p *Pellet) Radius() float64 {
	return Radius(p.Mass)
}
