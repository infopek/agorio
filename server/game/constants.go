package game

const (
	WorldWidth  float64 = 10000.0
	WorldHeight float64 = 10000.0

	TickRate uint64 = 20 // ticks / sec

	CollisionResolutionPasses int64   = 3
	EmptySpaceMaxAttempts     int64   = 30
	EmptySpaceRadiusQuery     float64 = 200.0
	EmptySpaceRadiusLeeway    float64 = 50.0

	LeaderboardSize int64 = 10

	// Player
	MaxNameLength int64 = 16
	TargetPlayers int64 = 20

	BaseSpeed     float64 = 8.0 // world units per tick at StartMass
	MinSpeed      float64 = 1.5
	SpeedExponent float64 = 0.2  // 0.0 = big cells stay a bit faster, 1.0 = big cells slow down more
	TurnSpeed     float64 = 0.45 // 0.0 = slow, 1.0 = fast

	CellCenterThreshold float64 = 7.0 // radius / {val} is considered the center, used for movement
	SpeedFactorRange    float64 = 3.5 // speed factor is 1.0 if cursor if farther than radius * {val} from center of cell

	StartMass   float64 = 100.0
	RadiusScale float64 = 3.6 // mass to radius constant

	MassDecayRate     float64 = 0.000015
	MassDecayExponent float64 = 1.4  // 1.0 = big cells decay linearly, 2.0 = quadratic decay relative to mass
	MinDecayMass      float64 = 20.0 // can't decay below or at this mass

	EatDistanceThreshold float64 = 0.4 // overlap of radii required for eating
	EatMassThreshold     float64 = 1.2 // mass ratio required for eating

	SplitMaxCells             int64   = 16 // max amount of cells you can have with splitting
	SplitMinMass              float64 = 20.0
	MaxCellMass               float64 = 20000.0 // auto split at or above this mass
	SplitMomentumFactor       float64 = 20.0    // momentum magnitude after split
	SplitMomentumRadiusFactor float64 = 0.5
	MomentumThreshold         float64 = 7.0  // below this value, collisions kick in
	SplitMomentumDecay        float64 = 0.83 // split momentum decays by this amount every tick

	EjectMinMass       float64 = 30.0 // can't eject mass below or at this level
	EjectMomentum      float64 = 27.0 // ejected mass' momentum
	EjectMomentumDecay float64 = 0.91
	EjectAmount        float64 = 12.0
	EjectPenalty       float64 = 1.30 // the cell loses the amount * penalty on feeding

	MergeTimerStartSeconds float64 = 13.0 // seconds until merge can happen
	MergeTimerMassFactor   float64 = 0.02 // how much the mass influences extra time
	MaxMergeTimer          float64 = 30.0
	MergeMinOverlap        float64 = 0.3
	MergeCooldownSeconds   float64 = 0.5 // extra seconds after merging to merge again

	// Pellet
	MaxPellets        int64   = 800
	PelletSpawnAmount int64   = 3   // per tick
	PelletSpawnChance float64 = 0.2 // per tick
	MinPelletMass     float64 = 1.0
	MaxPelletMass     float64 = 2.5

	// Virus
	MinViruses              int64   = 50
	MaxViruses              int64   = 80   // from MinViruses, we could spawn this many via feeding
	VirusSpawnAmount        int64   = 1    // per tick
	VirusSpawnChance        float64 = 0.15 // per tick
	VirusStartMass          float64 = 100.0
	VirusPopEqualSplitRange float64 = 5.0 // from SplitMinMass, used for a virus pop rule
	VirusPopDampenFactor    float64 = 0.8 // Momentum dampener

	VirusFeedToShoot   int64   = 7
	VirusLaunchSpeed   float64 = 30.0
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
		"Jumbo",
		"Agar Youtube",
		"Chonker",
		"Clanker",
	}

	BotNames = [...]string{
		"Latvia",
		"9gag",
		"4chan",
		"Hungary",
		"Poland",
		"Reddit",
		"X",
		"YouTube",
		"Google",
	}
)
