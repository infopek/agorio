package constants

import (
	"github.com/infopek/agorio/internal/types"
)

const (
	// App
	Port               types.Integer = 8080
	DefaultWorldWidth  types.Real    = 2000.0
	DefaultWorldHeight types.Real    = 2000.0
	TickRate           types.Real    = 60.0 // Hz
	MaxPlayers         types.Integer = 10

	// Game
	PlayerStartMass types.Real    = 10.0
	PlayerMaxCells  types.Integer = 16

	VirusMinStage   types.Integer = 0
	VirusMaxStage   types.Integer = 7
	VirusStartMass  types.Integer = 102
	VirusStartCount types.Integer = 15
	VirusMaxCount   types.Integer = 25

	PelletStartCount types.Integer = 20

	// Math / Physics
	DefaultRestitution  types.Real    = 0.2
	Dt                  types.Real    = 1.0 / types.Real(TickRate)
	CollisionIterations int = 3
)
