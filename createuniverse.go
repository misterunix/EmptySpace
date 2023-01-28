package main

import (
	"crypto/sha256"
	"database/sql"
	"log"
	"reflect"

	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pmcxs/hexgrid"
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

// Populate items in the database.
// Creates a new database if one does not exist.
// Number of systems and planets are created based on the input 'universesystems'
func CreateUniverse(universesystems int) {
	fmt.Println("Creating systems in memory")
	fmt.Println("universesystems:", universesystems)
	universesize := int(float64(universesystems) * 0.5)
	fmt.Println("universesize:", universesize)
	universefog := int(float64(universesize) * 0.03)
	fmt.Println("universefog:", universefog)

	for i := 0; i < universesystems; i++ {
		s := NewSystem()
		s.ID = i
		for {
			s.Q = rnd.Intn(universesize) + 1
			s.R = rnd.Intn(universesize) + 1
			A := hexgrid.NewHex(s.Q, s.R)
			var bad bool
			if i > 0 {
				for j := 0; j < i; j++ {
					if i == j {
						continue
					}
					B := hexgrid.NewHex(Systems[j].Q, Systems[j].R)
					C := hexgrid.HexDistance(A, B)
					if C > universefog {
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
			pl.PlayerID = -1
			Planets = append(Planets, pl)
			planetcount++
		}

	}

	fmt.Println()

	fmt.Println("Planet Count:", planetcount)

	player := NewDBPlayer()
	player.ID = 0
	player.Name = "Unknown"
	player.HomeWorldID = rnd.Intn(planetcount)
	player.Username = "Unknown"
	player.Password = "Unknown"
	player.Race = 0
	player.Email = "no@nope.com"
	player.AI = false
	Players = append(Players, player)
	Planets[player.HomeWorldID].PlayerID = player.ID

	fmt.Println("Creating a few NPCs. Player 0 is always human player.")
	pc := rnd.Intn(3) + 5
	fmt.Println("Number of NPCs:", pc)
	for i := 1; i < pc; i++ {
		fmt.Print(i, " ")
		np := NewDBPlayer()
		np.ID = i
		np.Name = "AI_" + strconv.Itoa(i)
		np.HomeWorldID = rnd.Intn(planetcount)
		np.Username = "AI_" + strconv.Itoa(i)
		np.Race = rnd.Intn(9) + 1
		np.Email = "ai@skynet.net"
		sha := sha256.New()
		sha.Write([]byte("BadPassword"))
		cp := fmt.Sprintf("%x", sha.Sum(nil))
		np.Password = cp
		//fmt.Println(np)
		Players = append(Players, np)
		//fmt.Println(np.Name, np.Password)
		Planets[np.HomeWorldID].PlayerID = i
	}
	fmt.Println()
	fmt.Println("Saving systems to DB")
	InsertSystems()
	InsertPlanets()
	InsertPlayers()

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
		case reflect.Bool:
			vt = "INTEGER"
		}
		sqlstatement += vt
	}
	// such a crappy way to do this
	sqlstatement = sqlstatement[1:]
	sqlstatement += ")"
	sqlstatement = sqlstatement1 + sqlstatement

	return sqlstatement

}

func InsertSystems() {
	fmt.Println("InsertSystems")
	o := "BEGIN;\n"
	beginstatement, err := database.Prepare(o)
	if err != nil {
		log.Panicln(err)
	}
	_, err = beginstatement.Exec()
	if err != nil {
		log.Panicln(err)
	}

	for _, v := range Systems {
		//fmt.Println("Inserting system", ind)
		var ssql, vsql string
		e := reflect.ValueOf(&v).Elem()
		for i := 0; i < e.NumField(); i++ {
			varName := e.Type().Field(i).Name
			ssql += "," + varName
			varValue := e.Field(i).Interface()
			vv := fmt.Sprintf("%v", varValue)
			varType := e.Type().Field(i).Type
			switch varType.Name() {
			case "int":
				vsql += "," + vv
			case "string":
				vsql += "," + "'" + vv + "'"
			}
		}
		ssql = ssql[1:]
		vsql = vsql[1:]
		sqlstatement := "INSERT INTO system( " + ssql + " ) VALUES(" + vsql + ")"
		//fmt.Println(sqlstatement)
		statement, _ := database.Prepare(sqlstatement)
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

func InsertPlanets() {
	fmt.Println("Inserting planets")
	o := "BEGIN;\n"
	beginstatement, err := database.Prepare(o)
	if err != nil {
		log.Panicln(err)
	}
	_, err = beginstatement.Exec()
	if err != nil {
		log.Panicln(err)
	}

	for _, v := range Planets {
		//fmt.Println("Inserting plantes", ind)
		var ssql, vsql string
		e := reflect.ValueOf(&v).Elem()
		for i := 0; i < e.NumField(); i++ {
			varName := e.Type().Field(i).Name
			ssql += "," + varName
			varValue := e.Field(i).Interface()
			vv := fmt.Sprintf("%v", varValue)
			varType := e.Type().Field(i).Type
			switch varType.Name() {
			case "int":
				vsql += "," + vv
			case "string":
				vsql += "," + "'" + vv + "'"
			}
		}
		ssql = ssql[1:]
		vsql = vsql[1:]
		sqlstatement := "INSERT INTO planet( " + ssql + " ) VALUES(" + vsql + ")"
		//fmt.Println(sqlstatement)
		statement, err := database.Prepare(sqlstatement)
		if err != nil {
			log.Panicln(err)
		}
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
func InsertPlayers() {
	fmt.Println("Inserting players")
	fmt.Println("Count:", len(Players))

	o := "BEGIN;\n"
	beginstatement, err := database.Prepare(o)
	if err != nil {
		log.Panicln(err)
	}
	_, err = beginstatement.Exec()
	if err != nil {
		log.Panicln(err)
	}

	for _, v := range Players {
		//fmt.Println("Inserting players", ind)
		//fmt.Println("v:", v.ID)

		var ssql, vsql string
		e := reflect.ValueOf(&v).Elem()
		for i := 0; i < e.NumField(); i++ {
			varName := e.Type().Field(i).Name
			ssql += "," + varName
			varValue := e.Field(i).Interface()
			vv := fmt.Sprintf("%v", varValue)
			varType := e.Type().Field(i).Type
			//fmt.Println(vsql, vv, varType.Name())
			switch varType.Name() {
			case "bool":
				vsql += "," + vv
			case "int":
				vsql += "," + vv
			case "string":
				vsql += "," + "'" + vv + "'"
			}
		}
		ssql = ssql[1:]
		vsql = vsql[1:]
		sqlstatement := "INSERT INTO player( " + ssql + " ) VALUES(" + vsql + ")"
		//fmt.Println(sqlstatement)
		statement, err := database.Prepare(sqlstatement)
		if err != nil {
			log.Panicln(err)
		}
		_, err = statement.Exec()
		if err != nil {
			log.Panicln(err)
		}

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
