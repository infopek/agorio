package game

import (
	"github.com/google/uuid"
)

type Player struct {
	ID    uuid.UUID `json:"-"`
	Name  string    `json:"name"`
	Cells []*Cell   `json:"cells"`
}
