package models

import (
	"github.com/infopek/agorio/internal/math"
	"github.com/infopek/agorio/internal/types"
)

type Cell struct {
	Position math.Vector2 `json:"position"`
	Mass     types.Real   `json:"mass"`
}
