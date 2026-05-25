package game

import (
	"github.com/google/uuid"
)

type Player struct {
	ID      uuid.UUID
	CellIDs []uuid.UUID // cell ids in the world

	Name   string
	Target Vec2 // mouse coords in world space

	Score uint32

	OutputChan chan<- ServerEvent
}
