package net

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/infopek/agorio/server/game"
)

type CellDTO struct {
	ID        string `json:"id"`
	OwnerID   string `json:"owner_id"`
	OwnerName string `json:"owner_name,omitempty"`

	X float64 `json:"x"`
	Y float64 `json:"y"`

	Radius uint64  `json:"radius"`
	Mass   uint64  `json:"mass"`
	Color  string `json:"color"`
}

type PelletDTO struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`

	Radius uint64  `json:"radius"`
	Mass   uint64  `json:"mass"`
	Color  string `json:"color"`
}

type EjectDTO struct {
	ID string `json:"id"`

	X float64 `json:"x"`
	Y float64 `json:"y"`

	Radius uint64  `json:"radius"`
	Mass   uint64  `json:"mass"`
	Color  string `json:"color"`
}

type VirusDTO struct {
	ID string `json:"id"`

	X float64 `json:"x"`
	Y float64 `json:"y"`

	FedCount uint64 `json:"fed_count"`
	Radius   uint64 `json:"radius"`
	Mass     uint64 `json:"mass"`
	Color    string `json:"color"`
}

type TickSnapshotDTO struct {
	Type    string      `json:"t"`
	Cells   []CellDTO   `json:"cells"`
	Pellets []PelletDTO `json:"pellets"`
	Ejects  []EjectDTO  `json:"ejects"`
	Viruses []VirusDTO  `json:"viruses"`
	Me      string      `json:"me"`
	Tick    uint64      `json:"tick"`
	Score   uint64      `json:"score"`
}

type DeathEventDTO struct {
	Type string `json:"t"`
}

type LeaderboardEntryDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Score uint64 `json:"score"`
}

type LeaderboardDTO struct {
	Type    string                `json:"t"`
	Me      string                `json:"me"`
	Entries []LeaderboardEntryDTO `json:"entries"`
}

func tickSnapshotToTickSnapshotDTO(snapshot game.TickSnapshot) TickSnapshotDTO {
	snapshotDTO := TickSnapshotDTO{
		Type:    "snapshot",
		Cells:   make([]CellDTO, len(snapshot.Cells)),
		Pellets: make([]PelletDTO, len(snapshot.Pellets)),
		Ejects:  make([]EjectDTO, len(snapshot.Ejects)),
		Viruses: make([]VirusDTO, len(snapshot.Viruses)),
		Me:      snapshot.Me.String(),
		Tick:    snapshot.Tick,
		Score:   snapshot.Score,
	}

	for i, c := range snapshot.Cells {
		snapshotDTO.Cells[i] = CellDTO{
			ID:        c.ID.String(),
			OwnerID:   c.OwnerID.String(),
			OwnerName: c.Name,

			X: c.X,
			Y: c.Y,

			Radius: uint64(c.Radius),
			Mass:   uint64(c.Mass),
			Color:  colorToHex(c.Color),
		}
	}

	for i, p := range snapshot.Pellets {
		snapshotDTO.Pellets[i] = PelletDTO{
			X: p.Position.X,
			Y: p.Position.Y,

			Radius: uint64(p.Radius()),
			Mass:   uint64(p.Mass),
			Color:  colorToHex(p.Color),
		}
	}

	for i, e := range snapshot.Ejects {
		snapshotDTO.Ejects[i] = EjectDTO{
			ID: e.ID.String(),

			X: e.Position.X,
			Y: e.Position.Y,

			Radius: uint64(e.Radius()),
			Mass:   uint64(e.Mass),
			Color:  colorToHex(e.Color),
		}
	}

	for i, v := range snapshot.Viruses {
		snapshotDTO.Viruses[i] = VirusDTO{
			ID: v.ID.String(),

			X: v.Position.X,
			Y: v.Position.Y,

			FedCount: v.FedCount,
			Radius:   uint64(v.Radius()),
			Mass:     uint64(v.Mass),
			Color:    colorToHex(v.Color),
		}
	}

	return snapshotDTO
}

func deathEventToDeathEventDTO(_ game.DeathEvent) DeathEventDTO {
	return DeathEventDTO{
		Type: "death",
	}
}

func leaderboardToLeaderboardDTO(lb game.Leaderboard) LeaderboardDTO {
	dto := LeaderboardDTO{
		Type:    "leaderboard",
		Me:      lb.Me.String(),
		Entries: make([]LeaderboardEntryDTO, len(lb.Entries)),
	}
	for i := range len(dto.Entries) {
		dto.Entries[i].ID = lb.Entries[i].ID.String()
		dto.Entries[i].Name = lb.Entries[i].Name
		dto.Entries[i].Score = lb.Entries[i].Score
	}

	return dto
}

func toJSON(msg game.ServerEvent) ([]byte, error) {
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
	case game.Leaderboard:
		dto := leaderboardToLeaderboardDTO(m)
		data, err := json.Marshal(dto)
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		return nil, nil
	}
}

func parsePlayerMessage(msg []byte) (game.PlayerEvent, error) {
	var raw map[string]any
	err := json.Unmarshal(msg, &raw)
	if err != nil {
		return nil, err
	}

	switch raw["t"] {
	case "join":
		return game.PlayerJoinEvent{
			Name: raw["name"].(string),
		}, nil
	case "move":
		return game.PlayerMoveEvent{
			Target: game.Vec2{
				X: raw["x"].(float64),
				Y: raw["y"].(float64),
			},
		}, nil
	case "split":
		return game.PlayerSplitEvent{}, nil
	case "feed":
		return game.PlayerFeedEvent{}, nil
	default:
		return nil, errors.New("unknown input message")
	}
}

func colorToHex(c [3]uint8) string {
	return fmt.Sprintf("#%02x%02x%02x", c[0], c[1], c[2])
}
