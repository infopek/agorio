package models

import (
	"github.com/infopek/agorio/internal/math"
	"github.com/infopek/agorio/internal/types"
)

// Represents the stationary pellet on the ground
type Pellet struct {
	Position math.Vector2
	Mass     types.Integer
}
