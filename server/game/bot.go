package game

import (
	"math"

	"github.com/google/uuid"
)

type Bot struct {
	ID          uuid.UUID
	Target      Vec2
	UpdateTimer float64 // ticks until next decision
}

func (w *World) spawnBot() {
	id := uuid.New()
	name := BotNames[RandIntRange(0, int64(len(BotNames)))]
	color := DefaultColors[RandIntRange(0, int64(len(DefaultColors)))]

	cell := Cell{
		ID:      uuid.New(),
		OwnerID: id,

		Position: w.findEmptySpace(Radius(StartMass)),

		Mass:  StartMass,
		Color: color,
	}

	player := Player{
		ID:   id,
		Name: name,

		CellIDs: []uuid.UUID{cell.ID},
		Target:  cell.Position,
		IsBot:   true,
	}

	w.Players[id] = &player
	w.Cells[cell.ID] = &cell
	w.Bots[id] = &Bot{ID: id, Target: cell.Position}
}

func (w *World) updateBots() {
	for _, bot := range w.Bots {
		bot.UpdateTimer--
		if bot.UpdateTimer <= 0 {
			bot.UpdateTimer = RandFloatRange(10.0, 40.0)
			p, ok := w.Players[bot.ID]
			if !ok {
				delete(w.Bots, bot.ID)
				continue
			}
			bot.Target = w.botDecide(p)
			p.Target = bot.Target
		}
	}
}

func (w *World) botDecide(p *Player) Vec2 {
	c := w.Cells[p.CellIDs[0]]

	nearest := Vec2{
		X: RandFloatRange(0.0, WorldWidth),
		Y: RandFloatRange(0.0, WorldHeight),
	}
	bestDist := math.MaxFloat64

	// Chase pellets
	neighborPellets := w.pelletGrid.GetNeighbors(c.Position, c.Radius())
	for _, p := range neighborPellets {
		dist := c.Position.DistanceTo(p.Position)
		if dist < bestDist {
			bestDist = dist
			nearest = p.Position
		}
	}

	// Chase smaller cells (if close enough)
	neighborCells := w.cellGrid.GetNeighbors(c.Position, c.Radius())
	for _, oc := range neighborCells {
		if oc.OwnerID == p.ID {
			continue // don't chase our own tail
		}

		if oc.Mass >= c.Mass {
			continue // we are smaller
		}

		dist := c.Position.DistanceTo(oc.Position)
		if dist < bestDist && dist < 500.0 {
			bestDist = dist
			nearest = oc.Position
		}
	}

	// Flee from larger cells
	for _, oc := range neighborCells {
		if oc.OwnerID == p.ID {
			continue // don't be scared of ourselves
		}

		if oc.Mass < c.Mass*EatMassThreshold {
			continue // we are bigger
		}

		dist := c.Position.DistanceTo(oc.Position)
		if dist < 300.0 {
			// Run away
			flee := c.Position.Sub(oc.Position).Normalize().Scale(500.0)
			return c.Position.Add(flee)
		}
	}

	return nearest
}
