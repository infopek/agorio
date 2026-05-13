package application

import (
	"github.com/infopek/agorio/internal/types"
)

type ClientMessage struct {
	Type string
	X    types.Real
	Y    types.Real
}
