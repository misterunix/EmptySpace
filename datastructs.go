package main

// Main data structure for the players.
type Planet struct {
	ID               int    // SQL Id and In systm Id
	GlobalID         int    // Id of planet in the universe. Should be unique. System use only.
	SystemID         int    // Id of system this planet is in.
	PType            int    //  0 - 9
	PlayerID         int    // Id of player that controls this planet.
	Name             string // Name of the planet.
	Population       int    // Population of the planet. in 100,000's
	TechLevel        int    // Tech level of the planet.
	Goverment        int    // Goverment type of the planet.
	RawOre           int    // Raw Ore on the planet.
	ProcessedOre     int    // Processed Ore on the planet.
	RawFood          int    // Raw Food on the planet.
	ProcessedFood    int    // Processed Food on the planet.
	Manufacturing    int    // Manufacturing on the planet.
	JumpFuel         int    // Fuel for Jump Engines
	JumpEngines      int    // Number of Jump Engines for ships.
	SublightFuel     int    // Fuel for Sublight Engines
	SubLightEngines  int    // Number of Sublight Engines for ships.
	Crystals         int    // Needed for Shield Generators
	ShieldGenerators int    // Number of Shield Generators for ships.
	SpacePort        int    // The Level of the Space Port
	LightHulls       int    // Number of Light Hulls
	MediumHulls      int    // Number of Medium Hulls
	HeavyHulls       int    // Number of Heavy Hulls
	ColonyShips      int    // Number of Colony Ships
	ScoutShips       int    // Number of Scout Ships

}

// Struct for holding system positions
type System struct {
	ID   int
	Name string
	Q    int
	R    int
}

// Struct for holding player data
type DBPlayer struct {
	ID          int
	Name        string
	Username    string
	Password    string
	Email       string
	HomeWorldID int // ID of the players home world. This can change if the player sishes it.
	AI          bool
	Race        int
}
