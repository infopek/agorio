package game

import (
	"github.com/infopek/agorio/internal/physics"
)

type Cell struct {
	Body *physics.Body `json:"body"`
}
