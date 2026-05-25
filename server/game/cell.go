package game

import (
	"math"

	"github.com/google/uuid"
)

type Cell struct {
	ID      uuid.UUID
	OwnerID uuid.UUID

	Position  Vec2
	Direction Vec2 // actual direction
	Momentum  Vec2 // a velocity vector

	MergeTimer float64 // time until merge is possible

	Mass  float64
	Color [3]uint8
}

func (c *Cell) Radius() float64 {
	return RadiusScale * math.Sqrt(c.Mass)
}
