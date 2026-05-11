package models

import (
	"github.com/infopek/agorio/internal/impulse"
	"github.com/infopek/agorio/internal/types"
)

type Cell struct {
	Body *impulse.Body `json:"body"`
	Mass types.Real `json:"mass"` // different from its physical mass
}
