package game

import (
	"math"

	"github.com/google/uuid"
)

type Eject struct {
	ID uuid.UUID

	Position Vec2
	Momentum Vec2

	Mass  float64
	Color [3]uint8
}

func (e *Eject) Radius() float64 {
	return RadiusScale * math.Sqrt(e.Mass)
}
