package game

import (
	"github.com/google/uuid"
)

type Pellet struct {
	ID uuid.UUID

	Position Vec2

	Radius float64
	Mass   int64
	Color  [3]uint8
}
