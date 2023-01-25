package main

import (
	"datastructs"
	"fmt"
	"math/rand"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pmcxs/hexgrid"
)

var rnd = rand.New(rand.NewSource(1))

func NewPlanet() datastructs.Planet {
	p := datastructs.Planet{}
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

func NewSystem() datastructs.System {
	s := datastructs.System{}
	s.Name = "Unknown"
	s.Q = 0
	s.R = 0
	return s
}

func NewDBPlayer() datastructs.DBPlayer {
	p := datastructs.DBPlayer{}
	return p
}

// CreateUniverse : populate items in the database
func CreateUniverse() {
	fmt.Println("Creating systems in memory")

	for i := 0; i < universesystems; i++ {
		s := es.NewSystem()
		s.ID = i
		for {
			s.Q = rnd.Intn(299) + 1
			s.R = rnd.Intn(299) + 1
			A := hexgrid.NewHex(s.Q, s.R)
			bad := 0
			if i > 0 {
				for j := 0; j < i; j++ {
					if i == j {
						continue
					}
					B := hexgrid.NewHex(datastructs.Systems[j].Q, datastructs.Systems[j].R)
					C := hexgrid.HexDistance(A, B)
					if C > 10 {
						bad = 0
						break // Might have broken code here
					} else {
						bad = 1
					}
				}
			}
			if bad == 0 {
				break
			}
		}
		name := "Unknown System " + strconv.Itoa(i)
		s.Name = name
		datastructs.Systems = append(datastructs.Systems, s)

	}

	fmt.Println()
	/*
		planetcount := 0
		fmt.Println("Creating planets in memory.")
		fmt.Println("universesystems:", universesystems)
		for i := 0; i < universesystems; i++ {
			np := godice.Roll(1, 8)
			for j := 0; j < np; j++ {
				pl := es.NewPlanet()
				pl.Id = planetcount
				pl.SystemId = i
				name := "System " + strconv.Itoa(i) + " Planet " + strconv.Itoa(j)
				pl.Name = name
				pl.PType = godice.Roll(1, 10)
				Planets = append(Planets, pl)
				planetcount++
			}
			bar1.Add(1)
		}
		bar1.Reset()
		fmt.Println()
	*/
	MakePlanets()
	fmt.Println("Planet Count:", planetcount)
	MakePlayers(planetcount)
	/*
		fmt.Println("Creating a few NPCs. Player 0 is always admin.")
		for i := 1; i < godice.Roll(1, 20); i++ {
			np := es.NewDBPlayer()
			np.Name = "Player " + strconv.Itoa(i)
			np.HomeWorldID = godice.Roll(1, planetcount)
			np.Username = "player" + strconv.Itoa(i)
			sha := sha256.New()
			sha.Write([]byte("BadPassword"))
			cp := fmt.Sprintf("%x\n", sha.Sum(nil))
			np.Password = cp
			Players = append(Players, np)
			fmt.Println(np.Name, np.Password)
			Planets[np.HomeWorldID].PlayerID = i
		}
	*/
}

func main() {
}
