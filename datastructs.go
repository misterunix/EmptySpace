package main

// Main data structure for the players.
type Planet struct {
	ID int // SQL ID and In systm global Id
	//GlobalID         int    // Id of planet in the universe. Should be unique. System use only.
	SystemID         int    // Id of system this planet is in.
	PType            int    //  0 - 9
	PlayerID         int    // Id of player that controls this planet.
	Name             string // Name of the planet.
	Population       int    // Population of the planet. in 100,000's
	TechLevel        int    // Tech level of the planet.
	Goverment        int    // Goverment type of the planet.
	Manufacturing    int    // Manufacturing on the planet. Ore Processing.
	Farming          int    // Farming on the planet. Food Processing.
	Research         int    // Research on the planet. Research.
	Crystals         int    // Needed for Shield Generators
	RawOre           int    // Raw Ore on the planet.
	RawFood          int    // Raw Food on the planet.
	ProcessedOre     int    // Processed Ore on the planet.
	ProcessedFood    int    // Processed Food on the planet.
	SublightFuel     int    // Fuel for Sublight Engines
	JumpFuel         int    // Fuel for Jump Engines
	SubLightEngines  int    // Number of Sublight Engines for ships.
	JumpEngines      int    // Number of Jump Engines for ships.
	ShieldGenerators int    // Number of Shield Generators for ships.
	SpacePort        int    // The Level of the Space Port
	LightHulls       int    // Number of Light Hulls
	MediumHulls      int    // Number of Medium Hulls
	HeavyHulls       int    // Number of Heavy Hulls
	ColonyShips      int    // Number of Colony Ships
	ScoutShips       int    // Number of Scout Ships
	Fighters         int    // Number of Fighters
	Marines          int    // Number of Marines
	WeaponTech       int    // Weapon Tech Level
	ArmorTech        int    // Armor Tech Level
	ShieldTech       int    // Shield Tech Level

}

// Struct for holding system positions
type System struct {
	ID    int    // ID of the system.
	Name  string // Name of the system.
	Q     int    // Hex coordinates
	R     int    // Hex coordinates
	Owner int    // ID of the player that owns this system.
}

// Struct for holding player data
type DBPlayer struct {
	ID           int
	Name         string
	Username     string
	Password     string // hash of the password
	Email        string
	HomeWorldID  int // ID of the players home world. This can change if the player wishes it.
	HomeSystemID int // ID of the players home system. This can change if the player wishes it.
	AI           bool
	Race         int
	TechLevel    int // Overall tech level of the player.
}

// Struct for holding ships. Nothing is ever removed from this table. Ships are just marked as destroyed.
type Ship struct {
	ID         int // ID of the ship. This is the SQL ID.
	PlayerID   int // ID of the player that owns this ship.
	PlanetID   int // ID of the planet this ship is on. -1 if in space.
	ShipType   int // Type of ship.
	CargoSpace int // Cargo space of the ship.
	SubLight   int // Sublight speed of the ship.
	Jump       int // Jump speed of the ship.
	Shield     int // Shield strength of the ship.
	Weapons    int // Weapons strength of the ship.
	Armor      int // Armor strength of the ship.
	Colonists  int // Number of colonists on the ship.
	Marines    int // Number of marines on the ship.
	Fighters   int // Number of fighters on the ship.
}
