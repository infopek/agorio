package game

import (
	"github.com/google/uuid"
)

type Cell struct {
	ID      uuid.UUID
	OwnerID uuid.UUID

	Position Vec2
	Momentum Vec2

	Radius float64 // calculated from mass
	Mass   int64
	Color  [3]uint8
}
