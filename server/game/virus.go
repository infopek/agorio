package game

import (
	"math"

	"github.com/google/uuid"
)

type Virus struct {
	ID uuid.UUID

	Position Vec2

	Mass  int64
	Color [3]uint8
}

func (v *Virus) Radius() float64 {
	return RadiusScale * math.Sqrt(float64(v.Mass))
}
