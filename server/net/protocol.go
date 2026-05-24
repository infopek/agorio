package net

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/infopek/agorio/server/game"
)

type CellDTO struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`

	X float64 `json:"x"`
	Y float64 `json:"y"`

	Radius float64 `json:"radius"`
	Mass   int64   `json:"mass"`
	Color  string  `json:"color"`
}

type PelletDTO struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`

	Radius float64 `json:"radius"`
	Mass   int64   `json:"mass"`
	Color  string  `json:"color"`
}

type VirusDTO struct {
	Position game.Vec2 `json:"position"`

	Radius float64 `json:"radius"`
	Mass   int64   `json:"mass"`
	Color  string  `json:"color"`
}

type TickSnapshotDTO struct {
	Type    string      `json:"t"`
	Cells   []CellDTO   `json:"cells"`
	Pellets []PelletDTO `json:"pellets"`
	Viruses []VirusDTO  `json:"viruses"`
	Me      string      `json:"me"`
	Tick    uint64      `json:"tick"`
	Score   uint32      `json:"score"`
}

type DeathEventDTO struct {
	Type string `json:"t"`
}

func tickSnapshotToTickSnapshotDTO(snapshot game.TickSnapshot) TickSnapshotDTO {
	snapshotDTO := TickSnapshotDTO{
		Type:    "snapshot",
		Cells:   make([]CellDTO, len(snapshot.Cells)),
		Pellets: make([]PelletDTO, len(snapshot.Pellets)),
		Viruses: make([]VirusDTO, len(snapshot.Viruses)),
		Me:      snapshot.Me.String(),
		Tick:    snapshot.Tick,
		Score:   snapshot.Score,
	}

	for i, c := range snapshot.Cells {
		snapshotDTO.Cells[i] = CellDTO{
			ID:      c.ID.String(),
			OwnerID: c.OwnerID.String(),

			X: c.Position.X,
			Y: c.Position.Y,

			Radius: c.Radius(),
			Mass:   c.Mass,
			Color:  colorToHex(c.Color),
		}
	}

	for i, p := range snapshot.Pellets {
		snapshotDTO.Pellets[i] = PelletDTO{
			X: p.Position.X,
			Y: p.Position.Y,

			Radius: p.Radius(),
			Mass:   p.Mass,
			Color:  colorToHex(p.Color),
		}
	}

	for i, v := range snapshot.Viruses {
		snapshotDTO.Viruses[i] = VirusDTO{
			Position: v.Position,

			Radius: v.Radius(),
			Mass:   v.Mass,
			Color:  colorToHex(v.Color),
		}
	}

	return snapshotDTO
}

func deathEventToDeathEventDTO(_ game.DeathEvent) DeathEventDTO {
	return DeathEventDTO{
		Type: "death",
	}
}

func toJSON(msg game.ServerMessage) ([]byte, error) {
	switch m := msg.(type) {
	case game.TickSnapshot:
		dto := tickSnapshotToTickSnapshotDTO(m)
		data, err := json.Marshal(dto)
		if err != nil {
			return nil, err
		}
		return data, nil
	case game.DeathEvent:
		dto := deathEventToDeathEventDTO(m)
		data, err := json.Marshal(dto)
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		return nil, nil
	}
}

func parsePlayerMessage(msg []byte) (game.PlayerMessage, error) {
	var raw map[string]any
	err := json.Unmarshal(msg, &raw)
	if err != nil {
		return nil, err
	}

	switch raw["t"] {
	case "join":
		return game.JoinMessage{
			Name: raw["name"].(string),
		}, nil
	case "move":
		return game.MoveMessage{
			Target: game.Vec2{
				X: raw["x"].(float64),
				Y: raw["y"].(float64),
			},
		}, nil
	case "split":
		return game.SplitMessage{}, nil
	case "feed":
		return game.FeedMessage{}, nil
	default:
		return nil, errors.New("unknown input message")
	}
}

func colorToHex(c [3]uint8) string {
	return fmt.Sprintf("#%02x%02x%02x", c[0], c[1], c[2])
}
