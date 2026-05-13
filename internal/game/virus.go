package game

import (
	"github.com/infopek/agorio/internal/math"
	"github.com/infopek/agorio/internal/types"
)

type Virus struct {
	Position math.Vector2  `json:"position"`
	Mass     types.Integer `json:"mass"`
	Stage    types.Integer `json:"stage"`
}
