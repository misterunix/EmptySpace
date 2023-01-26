package main

import (
	"crypto/sha256"
	"database/sql"
	"reflect"

	"fmt"
	"math/rand"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pmcxs/hexgrid"
)

var database *sql.DB   // Database connection
var Systems []System   // All systems in the universe
var Planets []Planet   // All planets in the universe
var Players []DBPlayer // All players in the universe. Including NPCs
var rnd = rand.New(rand.NewSource(1))

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

// Populate items in the database.
// Creates a new database if one does not exist.
// Number of systems and planets are created based on the input 'universesystems'
func CreateUniverse(universesystems int) {
	fmt.Println("Creating systems in memory")

	for i := 0; i < universesystems; i++ {
		s := NewSystem()
		s.ID = i
		for {
			s.Q = rnd.Intn(299) + 1 // 1-300
			s.R = rnd.Intn(299) + 1 // 1-300
			A := hexgrid.NewHex(s.Q, s.R)
			var bad bool
			if i > 0 {
				for j := 0; j < i; j++ {
					if i == j {
						continue
					}
					B := hexgrid.NewHex(Systems[j].Q, Systems[j].R)
					C := hexgrid.HexDistance(A, B)
					if C > 10 {
						bad = false
						break // Might have broken code here
					} else {
						bad = true
					}
				}
			}
			if !bad {
				break
			}
		}
		name := "Unknown System " + strconv.Itoa(i)
		s.Name = name
		Systems = append(Systems, s)
	}

	fmt.Println()

	planetcount := 0
	fmt.Println("Creating planets in memory.")
	fmt.Println("universesystems:", universesystems)
	for i := 0; i < universesystems; i++ {
		np := rnd.Intn(7) + 1
		for j := 0; j < np; j++ {
			pl := NewPlanet()
			pl.ID = planetcount
			pl.SystemID = i
			name := "System " + strconv.Itoa(i) + " Planet " + strconv.Itoa(j)
			pl.Name = name
			pl.PType = rnd.Intn(9) + 1
			Planets = append(Planets, pl)
			planetcount++
		}

	}

	fmt.Println()

	fmt.Println("Planet Count:", planetcount)

	player := NewDBPlayer()
	player.Name = "Unknown"
	player.HomeWorldID = 0
	player.Username = "Unknown"
	player.Password = "Unknown"
	player.Race = 0
	player.AI = false
	Players = append(Players, player)

	fmt.Println("Creating a few NPCs. Player 0 is always human player.")
	pc := rnd.Intn(3) + 1
	for i := 1; i < pc; i++ {
		np := NewDBPlayer()
		np.Name = "AI_ " + strconv.Itoa(i)
		np.HomeWorldID = rnd.Intn(planetcount)
		np.Username = "AI_" + strconv.Itoa(i)
		np.Race = rnd.Intn(9) + 1
		sha := sha256.New()
		sha.Write([]byte("BadPassword"))
		cp := fmt.Sprintf("%x\n", sha.Sum(nil))
		np.Password = cp
		Players = append(Players, np)
		fmt.Println(np.Name, np.Password)
		Planets[np.HomeWorldID].PlayerID = i
	}

}

// Create table based on struct. This is a work in progress.
func CreateTableFromStruct(table string, s interface{}) string {
	var e reflect.Value = reflect.ValueOf(s)
	var sqlstatement string
	sqlstatement1 := "CREATE TABLE IF NOT EXISTS " + table + " ("
	for i := 0; i < e.NumField(); i++ {
		var vt string
		varName := e.Type().Field(i).Name
		sqlstatement += "," + varName + " "
		//sqlstatement += varName + " "
		varType := e.Type().Field(i).Type
		switch varType.Kind() {
		case reflect.Int:
			if varName == "ID" {
				vt = "INTEGER NOT NULL PRIMARY KEY"
			} else {
				vt = "INTEGER"
			}
		case reflect.String:
			vt = "TEXT"
		case reflect.Float64:
			vt = "REAL"
		}
		sqlstatement += vt
	}
	// such a crappy way to do this
	sqlstatement = sqlstatement[1:]
	sqlstatement += ")"
	sqlstatement = sqlstatement1 + sqlstatement

	return sqlstatement

}

func CreateNewDB() {

	database, _ = sql.Open("sqlite3", "sqldata/emptyspace.db") // this should be in a config file

	statement, _ := database.Prepare("DROP TABLE IF EXISTS system")
	statement.Exec()

	statement, _ = database.Prepare("DROP TABLE IF EXISTS planet")
	statement.Exec()

	statement, _ = database.Prepare("DROP TABLE IF EXISTS player")
	statement.Exec()

	t := CreateTableFromStruct("system", System{})
	fmt.Println(t)
	statement, _ = database.Prepare(t)
	statement.Exec()

	t = CreateTableFromStruct("planet", Planet{})
	statement, _ = database.Prepare(t)
	statement.Exec()

	t = CreateTableFromStruct("player", DBPlayer{})
	statement, _ = database.Prepare(t)
	statement.Exec()
	/**
	// this is a hack to create the table based on the struct
	tsystem := NewSystem() // temporary system
	e := reflect.ValueOf(&tsystem).Elem()
	for i := 0; i < e.NumField(); i++ {
		var vt string
		varName := e.Type().Field(i).Name
		sqlstatement += "," + varName + " "
		varType := e.Type().Field(i).Type
		switch varType.Name() {
		case "int":
			if varName == "ID" {
				vt = "INTEGER NOT NULL PRIMARY KEY"
			} else {
				vt = "INTEGER"
			}
		case "string":
			vt = "TEXT"
		}
		sqlstatement += vt
	}
	sqlstatement = sqlstatement[1:]
	//fmt.Println(sqlstatement)
	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS system (" + sqlstatement + ")")
	statement.Exec()
	*/
}
