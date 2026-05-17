package game

import (
	"github.com/google/uuid"
)

type ServerMessage interface {
	serverMessage()
}

/** TickSnapshot
 * Everything a client needs to know per tick
 */
type TickSnapshot struct {
	Cells   []Cell
	Pellets []Pellet
	Viruses []Virus
	Me      []uuid.UUID
	Tick    int64
	Score   uint32
}

type DeathEvent struct {

}

func (snapshot TickSnapshot) serverMessage() {}
func (deathEvent DeathEvent) serverMessage() {}
