package net

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"

	"github.com/infopek/agorio/server/game"
)

type CellDTO struct {
	OwnerID uuid.UUID `json:"owner_id"`

	Position game.Vec2 `json:"position"`

	Radius float64  `json:"radius"`
	Mass   int64    `json:"mass"`
	Color  [3]uint8 `json:"color"`
}

type PelletDTO struct {
	Position game.Vec2 `json:"position"`

	Radius float64  `json:"radius"`
	Mass   int64    `json:"mass"`
	Color  [3]uint8 `json:"color"`
}

type VirusDTO struct {
	Position game.Vec2 `json:"position"`

	Radius float64  `json:"radius"`
	Mass   int64    `json:"mass"`
	Color  [3]uint8 `json:"color"`
}

type TickSnapshotDTO struct {
	Cells   []CellDTO   `json:"cells"`
	Pellets []PelletDTO `json:"pellets"`
	Viruses []VirusDTO  `json:"viruses"`
	Me      []uuid.UUID `json:"me"`
	Tick    int64       `json:"tick"`
	Score   uint32      `json:"score"`
}

func tickSnapshotToTickSnapshotDTO(snapshot game.TickSnapshot) TickSnapshotDTO {
	snapshotDTO := TickSnapshotDTO{
		Cells:   make([]CellDTO, len(snapshot.Cells)),
		Pellets: make([]PelletDTO, len(snapshot.Pellets)),
		Viruses: make([]VirusDTO, len(snapshot.Viruses)),
		Me:      make([]uuid.UUID, len(snapshot.Me)),
		Tick:    snapshot.Tick,
		Score:   snapshot.Score,
	}

	for i, c := range snapshot.Cells {
		snapshotDTO.Cells[i] = CellDTO{
			OwnerID: c.OwnerID,

			Position: c.Position,

			Radius: c.Radius,
			Mass:   c.Mass,
			Color:  c.Color,
		}
	}

	for i, p := range snapshot.Pellets {
		snapshotDTO.Pellets[i] = PelletDTO{
			Position: p.Position,

			Radius: p.Radius,
			Mass:   p.Mass,
			Color:  p.Color,
		}
	}

	for i, v := range snapshot.Viruses {
		snapshotDTO.Viruses[i] = VirusDTO{
			Position: v.Position,

			Radius: v.Radius,
			Mass:   v.Mass,
			Color:  v.Color,
		}
	}

	copy(snapshotDTO.Me, snapshot.Me)

	return snapshotDTO
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
		return nil, nil
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
