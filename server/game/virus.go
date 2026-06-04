package game

import (
	"github.com/google/uuid"
)

type Virus struct {
	ID uuid.UUID

	Position Vec2
	Momentum Vec2 // for shooting virus

	Mass     float64
	FedCount uint64
	Color    [3]uint8
}

func (v *Virus) Radius() float64 {
	return Radius(v.Mass)
}
