package main

import (
	"math/rand"
	"time"
)

// var rnd = rand.New(rand.NewSource(time.Hour.Microseconds()))
var rnd *rand.Rand

func main() {

	seed := time.Now().UnixNano()
	rnd = rand.New(rand.NewSource(seed))

	CreateNewDB()
	CreateUniverse(1000)

}
