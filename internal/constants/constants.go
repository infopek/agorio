package constants

import (
	"github.com/infopek/agorio/internal/types"
)

const (
	// App
	Port         types.Integer = 8080
	WorldWidth   types.Real    = 2000.0
	WorldHeight  types.Real    = 2000.0
	CanvasWidth  types.Real    = 800.0
	CanvasHeight types.Real    = 600.0
	TickRate     types.Real    = 60.0 // Hz
	MaxPlayers   types.Integer = 10

	// Game
	PlayerStartMass types.Real    = 3.0
	PlayerMaxCells  types.Integer = 16
	PlayerBaseSpeed types.Real    = 400.0

	VirusMinStage   types.Integer = 0
	VirusMaxStage   types.Integer = 7
	VirusStartMass  types.Integer = 102
	VirusStartCount types.Integer = 15
	VirusMaxCount   types.Integer = 25

	PelletStartCount types.Integer = 20

	// Math / Physics
	Dt                   types.Real    = 1.0 / types.Real(TickRate)
	CollisionIterations  types.Integer = 3
	PosCorrectionPercent types.Real    = 0.8
	PosCorrectionSlop    types.Real    = 0.01
)
