package game

const (
	WorldWidth  float64 = 5000.00
	WorldHeight float64 = 5000.00

	TickRate int64 = 20 // ticks / sec

	// Player
	BaseSpeed     float64 = 6.00 // world units per tick at StartMass
	MinSpeed      float64 = 1.50
	SpeedExponent float64 = 0.20
	TurnSpeed     float64 = 0.45 // 0.0 = slow, 1.0 = fast

	StartMass   int64   = 200
	RadiusScale float64 = 4.00 // mass to radius constant

	EatDistanceThreshold float64 = 0.40 // overlap of radii required for eating
	EatMassThreshold     float64 = 1.20 // mass ratio required for eating

	SplitMinMass       int64   = 20
	SplitSpeed         float64 = 15.00 // momentum magnitude after split
	SplitVelocityDecay float64 = 0.90

	// Pellet
	MinPellets    int64 = 300
	MaxPellets    int64 = 500
	MinPelletMass int64 = 1
	MaxPelletMass int64 = 3

	// Virus
	MaxViruses int64 = 12
)
