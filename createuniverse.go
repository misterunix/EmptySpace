package main

import (
	"crypto/sha256"
	"database/sql"
	"log"
	"math"
	"math/rand"

	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/misterunix/hexgrid"
)

var database *sql.DB   // Database connection
var Systems []System   // All systems in the universe
var Planets []Planet   // All planets in the universe
var Players []DBPlayer // All players in the universe. Including NPCs

func NewPlanet() Planet {
	p := Planet{}
	p.Name = "Unknown"
	p.SystemID = 0
	p.PlayerID = 0
	p.PType = 0
	p.Population = 0
	p.TechLevel = 0
	p.Goverment = 0
	p.RawOre = 0
	p.ProcessedOre = 0
	p.RawFood = 0
	p.ProcessedFood = 0
	p.Manufacturing = 0
	p.JumpFuel = 0
	p.SublightFuel = 0
	p.SpacePort = 0
	return p
}

func NewSystem() System {
	s := System{}
	s.Name = "Unknown"
	s.Q = 0
	s.R = 0
	return s
}

func NewDBPlayer() DBPlayer {
	p := DBPlayer{}
	return p
}

// Create systems
func CreateSystems(universesystems int) {
	fmt.Println("universesystems:", universesystems)

	tr := rand.Float64() * 0.1
	universesize := int(float64(universesystems) * (0.4 + tr))
	fmt.Println("Universe Q,R size:", universesize)

	universeMinDistance := int(float64(universesize) * 0.03) // How far away systems can be from each other
	fmt.Println("Universe Fog:", universeMinDistance)

	// Loop through the number of systems to create
	for i := 0; i < universesystems; i++ {
		s := NewSystem() // Create a new system
		s.ID = i         // Set the ID to the current loop number

		for { // Loop until we find a good location for the system

			// Set a temporary Q,R location
			s.Q = rnd.Intn(universesize) + 1
			s.R = rnd.Intn(universesize) + 1
			A := hexgrid.NewHex(s.Q, s.R) // Create a new hexgrid location

			var bad bool // bad location flag, set to true if we find a bad location
			if i > 0 {
				for j := 0; j < i; j++ {
					if i == j { // Cant compare to itself
						continue
					}

					B := hexgrid.NewHex(Systems[j].Q, Systems[j].R) // Create a new hexgrid location, based on the system we are comparing to. Limited by the current loop number
					C := hexgrid.HexDistance(A, B)                  // Get the distance between the two hexgrid locations

					if C > universeMinDistance { // If the distance is greater than the fog limit, then we are good
						bad = false
						break // break out of the loop
					} else {
						bad = true // to close, try again
					}
				}
			}
			if !bad { // If we are not too close to another system, then we are good
				break
			}
		}
		name := "Unknown System " + strconv.Itoa(i) // Set the name of the system to a generic name.
		s.Name = name
		s.Owner = -1                 // No one owns the system
		Systems = append(Systems, s) // Add the system to the Systems array
	}

	fmt.Println("System Generation Complete.")

}

// Create planets
func CreatePlanets(universesystems int) int {

	var planetcount int // Planet counter

	// Loop through the number of systems to create planets
	for i := 0; i < universesystems; i++ {

		numberofplanets := rnd.Intn(7) + 1 // Number of planets in the system

		for j := 0; j < numberofplanets; j++ {
			pl := NewPlanet() // Create a new planet
			pl.ID = planetcount
			pl.SystemID = i
			name := "System " + strconv.Itoa(i) + " Planet " + strconv.Itoa(j) // Generic name for the planet
			pl.Name = name
			pl.PType = rnd.Intn(9) + 1    // Random planet type
			pl.PlayerID = -1              // No player owns the planet
			Planets = append(Planets, pl) // Add the planet to the Planets array
			planetcount++                 // Increment the planet counter
		}
	}

	fmt.Println("Planet Generation Complete.")
	fmt.Println("Number of Planets:", planetcount)
	return planetcount
}

// Create Players
func CreatePlayers(universesystems int, planetcount int) {
	fmt.Println("Creating Players")
	minDistance := int(math.Sqrt(float64(universesystems)))

	// Human player
	player := NewDBPlayer()
	player.ID = 0
	player.Name = "Unknown"
	player.HomeWorldID = rnd.Intn(planetcount) // Random home world. This is the GlobalID of the planet.
	tmpSystemID := GetSystemIDFromGlobalID(player.HomeWorldID)
	player.HomeSystemID = tmpSystemID
	Systems[tmpSystemID].Owner = player.ID
	player.Username = "Unknown"
	player.Password = "Unknown"
	player.Race = 0
	player.Email = "no@nope.com"
	player.AI = false
	Players = append(Players, player)
	Planets[player.HomeWorldID].PlayerID = player.ID

	// AI players
	pc := 25 // rnd.Intn(3) + 5
	fmt.Println("Number of NPCs:", pc)
	for i := 1; i < pc; i++ {
		fmt.Print(i, " ")
		var tmpRandPlanet int
		var tmpSystemID int
		// Start with a random planet and check if is the minDistance away from all other players.
		for {
			tmpRandPlanet = rnd.Intn(planetcount)
			tmpSystemID = GetSystemIDFromGlobalID(tmpRandPlanet)
			// make sure the system is not owned by another player
			if Systems[tmpSystemID].Owner != -1 {
				continue
			}

			A := hexgrid.NewHex(Systems[tmpSystemID].Q, Systems[tmpSystemID].R) // Create a new hexgrid location
			var abort bool = false

			for _, tmpPlayer := range Players {
				t1 := tmpPlayer.HomeWorldID
				t2 := GetSystemIDFromGlobalID(t1)
				B := hexgrid.NewHex(Systems[t2].Q, Systems[t2].R)
				C := hexgrid.HexDistance(A, B)
				//fmt.Println("Distance:", C)
				if C >= minDistance {
					abort = true
				} else {
					abort = false
					break
				}
			}
			if abort {
				break
			}
		}

		np := NewDBPlayer() // Create a new player
		np.ID = i
		np.Name = "AI_" + strconv.Itoa(i) // Generic name
		// Since systems are created first, we can use the system ID to get the global ID of the planet.
		np.HomeWorldID = tmpRandPlanet
		np.HomeSystemID = tmpSystemID
		np.Username = "AI_" + strconv.Itoa(i) // Generic username
		np.Race = rnd.Intn(9) + 1
		np.Email = "ai@skynet.net" // Silly little email address. Not used for anything.
		sha := sha256.New()
		sha.Write([]byte("BadPassword"))
		cp := fmt.Sprintf("%x", sha.Sum(nil))
		np.Password = cp
		Players = append(Players, np)        // Add the player to the slice of players.
		Planets[np.HomeWorldID].PlayerID = i // Set the planet owner to the player ID.
		fmt.Println(np)

	}
}

// Create the universe in memory to make it faster and easier to check for collisions.
// Populate items in the database.
// Creates a new database if one does not exist.
// Number of systems and planets are created based on the input 'universesystems'
func CreateUniverse(universesystems int) {
	fmt.Println("Creating systems in memory")

	CreateSystems(universesystems)                // Create the systems in memory
	planetCount := CreatePlanets(universesystems) // Create the planets in memory
	CreatePlayers(universesystems, planetCount)   // Create the players in memory

	fmt.Println("Saving systems to the database.")
	InsertSystems() // Small little sub function to loop over the slice of systems.
	InsertPlanets() // Small little sub function to loop over the slice of planets.
	InsertPlayers() // Small little sub function to loop over the slice of players.

}

// Returns the SystemID of a planet based on the GlobalID
func GetSystemIDFromGlobalID(gid int) int {
	fmt.Println("Getting system ID from global ID:", gid)
	for _, v := range Planets {
		if v.ID == gid {
			fmt.Println("System ID:", v.SystemID)
			return v.SystemID
		}
	}
	return -1
}

// Insert systems into the database. This sub function is needed to loop over the slice of systems.
// Easier to do it this way than to try to loop over the slice in the main function.
func InsertSystems() {

	fmt.Println("Inserting systems into the database.")

	// Begin and Commit are needed to speed up the insert process.
	o := "BEGIN;\n"
	beginstatement, err := database.Prepare(o)
	if err != nil {
		log.Panicln(err)
	}
	_, err = beginstatement.Exec()
	if err != nil {
		log.Panicln(err)
	}

	// Loop over the slice of systems and call the InsertIntoTable function to create the SQL statement.
	for _, v := range Systems {
		tmpsql := InsertIntoTable("system", v)
		//fmt.Println(tmpsql)
		statement, _ := database.Prepare(tmpsql)
		statement.Exec()
	}

	o = "COMMIT;\n"
	commitstatement, err := database.Prepare(o)
	if err != nil {
		log.Panicln(err)
	}
	_, err = commitstatement.Exec()
	if err != nil {
		log.Panicln(err)
	}

}

// Insert planets into the database. This sub function is needed to loop over the slice of planets.
// Easier to do it this way than to try to loop over the slice in the main function.
func InsertPlanets() {

	fmt.Println("Inserting planets into the database.")

	// Begin and Commit are needed to speed up the insert process.
	o := "BEGIN;\n"
	beginstatement, err := database.Prepare(o)
	if err != nil {
		log.Panicln(err)
	}
	_, err = beginstatement.Exec()
	if err != nil {
		log.Panicln(err)
	}

	// Loop over the slice of planets and call the InsertIntoTable function to create the SQL statement.
	for _, v := range Planets {
		tmpsql := InsertIntoTable("planet", v)
		fmt.Println(tmpsql)
		statement, _ := database.Prepare(tmpsql)
		statement.Exec()
	}

	o = "COMMIT;\n"
	commitstatement, err := database.Prepare(o)
	if err != nil {
		log.Panicln(err)
	}
	_, err = commitstatement.Exec()
	if err != nil {
		log.Panicln(err)
	}

}

// Insert players into the database. This sub function is needed to loop over the slice of players.
// Easier to do it this way than to try to loop over the slice in the main function.
func InsertPlayers() {

	fmt.Println("Inserting players into the database.")
	fmt.Println("Count:", len(Players))

	// Begin and Commit are needed to speed up the insert process.
	o := "BEGIN;\n"
	beginstatement, err := database.Prepare(o)
	if err != nil {
		log.Panicln(err)
	}
	_, err = beginstatement.Exec()
	if err != nil {
		log.Panicln(err)
	}

	// Loop over the slice of players and call the InsertIntoTable function to create the SQL statement.
	for _, v := range Players {
		tmpsql := InsertIntoTable("player", v)
		//fmt.Println(tmpsql)
		statement, _ := database.Prepare(tmpsql)
		statement.Exec()
	}

	o = "COMMIT;\n"
	commitstatement, err := database.Prepare(o)
	if err != nil {
		log.Panicln(err)
	}
	_, err = commitstatement.Exec()
	if err != nil {
		log.Panicln(err)
	}

}

func CreateNewDB() {

	database, _ = sql.Open("sqlite3", "sqldata/emptyspace.db") // this should be in a config file

	statement, _ := database.Prepare("DROP TABLE IF EXISTS system")
	_, err := statement.Exec()
	if err != nil {
		log.Panicln(err)
	}

	statement, _ = database.Prepare("DROP TABLE IF EXISTS planet")
	_, err = statement.Exec()
	if err != nil {
		log.Panicln(err)
	}

	statement, _ = database.Prepare("DROP TABLE IF EXISTS player")
	_, err = statement.Exec()
	if err != nil {
		log.Panicln(err)
	}

	t := CreateTableFromStruct("system", System{})
	//fmt.Println(t)
	statement, _ = database.Prepare(t)
	_, err = statement.Exec()
	if err != nil {
		log.Panicln(err)
	}

	t = CreateTableFromStruct("planet", Planet{})
	statement, _ = database.Prepare(t)
	_, err = statement.Exec()
	if err != nil {
		log.Panicln(err)
	}

	t = CreateTableFromStruct("player", DBPlayer{})
	statement, _ = database.Prepare(t)
	_, err = statement.Exec()
	if err != nil {
		log.Panicln(err)
	}

}
