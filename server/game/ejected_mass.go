package game

import (
	"github.com/google/uuid"
)

type EjectedMass struct {
	ID uuid.UUID

	Position Vec2
	Momentum Vec2

	Mass  float64
	Color [3]uint8
}
