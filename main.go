package main

import (
	"database/sql"
	"flag"
	"math/rand"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// var rnd = rand.New(rand.NewSource(time.Hour.Microseconds()))
var rnd *rand.Rand

var globalstartclock time.Time
var database *sql.DB

func main() {

	seed := time.Now().UnixNano()
	rnd = rand.New(rand.NewSource(seed))

	var newuniverse bool

	flag.BoolVar(&newuniverse, "new", false, "Create a new universe.")

	flag.Parse()

	database, _ = sql.Open("sqlite3", "sqldata/emptyspace.db") // this should be in a config file

	if newuniverse {
		CreateNewDB()
		CreateUniverse(1000)
		os.Exit(1)
	}

	globalstartclock = time.Now()

}
