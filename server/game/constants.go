package game

const (
	WorldWidth  float64 = 5000.0
	WorldHeight float64 = 5000.0

	TickRate int64 = 20 // ticks / sec

	// Player
	MaxNameLength int64 = 16

	BaseSpeed     float64 = 6.0 // world units per tick at StartMass
	MinSpeed      float64 = 1.50
	SpeedExponent float64 = 0.20 // 0.0 = big cells stay a bit faster, 1.0 = big cells slow down more
	TurnSpeed     float64 = 0.45 // 0.0 = slow, 1.0 = fast

	CellCenterThreshold float64 = 7.0 // radius / {val} is considered the center, used for movement
	SpeedFactorRange    float64 = 0.5 // speed factor is 1.0 if cursor if farther than radius * {val} from center of cell

	StartMass   float64 = 400.0
	RadiusScale float64 = 4.5 // mass to radius constant

	MassDecayRate     float64 = 0.000015
	MassDecayExponent float64 = 1.35 // 1.0 = big cells decay linearly, 2.0 = quadratic decay relative to mass
	MinDecayMass      float64 = 20.0 // can't decay below or at this mass

	EatDistanceThreshold float64 = 0.40 // overlap of radii required for eating
	EatMassThreshold     float64 = 1.20 // mass ratio required for eating

	SplitMaxCells      int64   = 16 // max amount of cells you can have with splitting
	SplitMinMass       float64 = 20.0
	SplitMomentum      float64 = 18.0 // momentum magnitude after split
	MomentumThreshold  float64 = 11.0 // below this value, collisions kick in
	SplitMomentumDecay float64 = 0.95 // split momentum decays by this amount every tick

	FeedMinMass      float64 = 22.0 // can't eject mass below or at this level
	FeedMassMomentum float64 = 18.0 // ejected mass' momentum
	FeedMassAmount   float64 = 4.0
	FeedMassPenalty  float64 = 4.8

	MergeTimerStartSeconds float64 = 5.0 // seconds until merge can happen
	MergeMinOverlap        float64 = 0.65
	MergeCooldownSeconds   float64 = 1.0 // extra seconds after merging to merge again

	// Pellet
	MinPellets    int64   = 300
	MaxPellets    int64   = 500
	MinPelletMass float64 = 1.0
	MaxPelletMass float64 = 3.0

	// Virus
	MaxViruses int64 = 12
)

var (
	DefaultColors = [...][3]uint8{
		{255, 70, 70},   // red
		{70, 130, 255},  // blue
		{50, 200, 80},   // green
		{255, 165, 0},   // orange
		{180, 70, 255},  // purple
		{0, 200, 200},   // teal
		{255, 105, 180}, // pink
		{100, 100, 255}, // indigo
		{255, 200, 0},   // gold
		{0, 180, 130},   // emerald
	}

	DefaultNames = [...]string{
		"Hungry Blob",
		"Sir EatsALot",
		"Nom Nom",
		"Chonker",
	}
)
