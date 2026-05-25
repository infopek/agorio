package game

import (
	"github.com/google/uuid"
)

type ServerEvent interface {
	serverEvent()
}

/** CellView
 *
 * Necessary cell data for client
 *
 */
type CellView struct {
	ID      uuid.UUID
	OwnerID uuid.UUID
	Name    string // owner name

	X      float64
	Y      float64
	Radius float64

	Mass  float64
	Color [3]uint8
}

/** TickSnapshot
 *
 * Everything a client needs to know per tick
 *
 */
type TickSnapshot struct {
	Me      uuid.UUID // client player id
	Cells   []CellView
	Pellets []Pellet
	Viruses []Virus
	Tick    uint64
	Score   uint32
}

type DeathEvent struct{}

func (snapshot TickSnapshot) serverEvent() {}
func (deathEvent DeathEvent) serverEvent() {}
