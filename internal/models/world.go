package models

import (

	"github.com/infopek/agorio/internal/math"
)

type WorldState struct {
	PlayerPositions []math.Vector2	`json:"player_positions"`
}
