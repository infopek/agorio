package game

const (
	WorldWidth  float64 = 10000.0
	WorldHeight float64 = 10000.0

	TickRate int64 = 20 // ticks / sec

	CollisionResolutionPasses int64 = 3

	// Player
	MaxNameLength int64 = 16

	BaseSpeed     float64 = 8.0 // world units per tick at StartMass
	MinSpeed      float64 = 1.50
	SpeedExponent float64 = 0.20 // 0.0 = big cells stay a bit faster, 1.0 = big cells slow down more
	TurnSpeed     float64 = 0.45 // 0.0 = slow, 1.0 = fast

	CellCenterThreshold float64 = 7.0 // radius / {val} is considered the center, used for movement
	SpeedFactorRange    float64 = 3.5 // speed factor is 1.0 if cursor if farther than radius * {val} from center of cell

	StartMass   float64 = 100.0
	RadiusScale float64 = 3.6 // mass to radius constant

	MassDecayRate     float64 = 0.000015
	MassDecayExponent float64 = 1.40 // 1.0 = big cells decay linearly, 2.0 = quadratic decay relative to mass
	MinDecayMass      float64 = 20.0 // can't decay below or at this mass

	EatDistanceThreshold float64 = 0.40 // overlap of radii required for eating
	EatMassThreshold     float64 = 1.20 // mass ratio required for eating

	SplitMaxCells             int64   = 16 // max amount of cells you can have with splitting
	SplitMinMass              float64 = 20.0
	ForceSplitMassThreshold   float64 = 15000.0
	SplitMomentumFactor       float64 = 20.0 // momentum magnitude after split
	SplitMomentumRadiusFactor float64 = 0.5
	MomentumThreshold         float64 = 7.0  // below this value, collisions kick in
	SplitMomentumDecay        float64 = 0.83 // split momentum decays by this amount every tick

	EjectMinMass       float64 = 30.0 // can't eject mass below or at this level
	EjectMomentum      float64 = 27.0 // ejected mass' momentum
	EjectMomentumDecay float64 = 0.91
	EjectAmount        float64 = 12.0
	EjectPenalty       float64 = 1.30 // the cell loses the amount * penalty on feeding

	MergeTimerStartSeconds float64 = 13.0 // seconds until merge can happen
	MergeMinOverlap        float64 = 0.4
	MergeCooldownSeconds   float64 = 1.0 // extra seconds after merging to merge again

	// Pellet
	MaxPellets        int64   = 800
	PelletSpawnAmount int64   = 2   // per tick
	PelletSpawnChance float64 = 0.9 // per tick
	MinPelletMass     float64 = 1.0
	MaxPelletMass     float64 = 2.5

	// Virus
	MinViruses              int64   = 50
	MaxViruses              int64   = 80  // from MinViruses, we could spawn this many via feeding
	VirusSpawnAmount        int64   = 1   // per tick
	VirusSpawnChance        float64 = 0.1 // per tick
	VirusStartMass          float64 = 100.0
	VirusPopEqualSplitRange float64 = 5.0 // from SplitMinMass, used for a virus pop rule
	VirusPopDampenFactor    float64 = 0.8 // Momentum dampener

	VirusFeedToShoot   int64   = 7
	VirusLaunchSpeed   float64 = 20.0
	VirusMomentumDecay float64 = 0.93
)

var (
	// Colors
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

	VirusColor = [3]uint8{50, 168, 54} // greenish

	// Names
	DefaultNames = [...]string{
		"Hungry Blob",
		"Sir EatALot",
		"Nom Nom",
		"Wun Wun",
		"Chonker",
		"Clanker",
	}
)
