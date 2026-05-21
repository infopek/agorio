package game

import (
	"github.com/google/uuid"
)

type Virus struct {
	ID uuid.UUID

	Position Vec2

	Radius float64 // calculated from mass
	Mass   int64
	Color  [3]uint8
}
