package main

// Main data structure for the players.
type Planet struct {
	ID            int    // SQL Id and In systm Id
	SystemID      int    // Id of system this planet is in.
	PType         int    //  0 - 9
	PlayerID      int    // Id of player that controls this planet.
	Name          string // Name of the planet.
	Population    int
	TechLevel     int
	Goverment     int
	RawOre        int
	ProcessedOre  int
	RawFood       int
	ProcessedFood int
	Manufacturing int
	JumpFuel      int
	SublightFuel  int
	SpacePort     int
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
	HomeWorldID int
	AI          bool
	Race        int
}
