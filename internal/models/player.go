package models

import (
	"github.com/google/uuid"

	"github.com/infopek/agorio/internal/math"
)

type Player struct {
	ID     uuid.UUID    `json:"-"`
	Name   string       `json:"name"`
	Cells  []Cell       `json:"cells"`
	Target math.Vector2 `json:"-"`
}
