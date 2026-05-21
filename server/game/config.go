package game

const (
	WorldWidth  float64 = 5000.0
	WorldHeight float64 = 5000.0

	TickRate int64 = 20 // ticks / sec

	// Player
	BaseSpeed   float64 = 6.0 // world units per tick at StartMass
	StartMass   int64   = 10
	RadiusScale float64 = 4.0

	// Pellet
	MinPellets int64 = 300
	MaxPellets int64 = 500
	MinPelletMass int64 = 1
	MaxPelletMass int64 = 3

	MaxViruses int64 = 12
)
