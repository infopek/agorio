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
	Me      uuid.UUID // client player id
	Cells   []Cell
	Pellets []Pellet
	Viruses []Virus
	Tick    uint64
	Score   uint32
}

type DeathEvent struct{}

func (snapshot TickSnapshot) serverMessage() {}
func (deathEvent DeathEvent) serverMessage() {}
