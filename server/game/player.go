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

	OutputChan chan<- ServerMessage
}

func (p *Player) removeCell(id uuid.UUID) {
	for i, cellID := range p.CellIDs {
		if cellID == id {
			p.CellIDs[i] = p.CellIDs[len(p.CellIDs)-1]
			p.CellIDs = p.CellIDs[:len(p.CellIDs)-1]
			return
		}
	}
}
