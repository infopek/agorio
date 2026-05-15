package game

import (
	"github.com/infopek/agorio/internal/physics"
	"github.com/infopek/agorio/internal/math"
)

type Cell struct {
	Body      *physics.Body `json:"body"`
	Direction math.Vector2  `json:"-"`
}
