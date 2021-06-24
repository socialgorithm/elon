package simulator

import (
	"log"
	"time"

	"github.com/socialgorithm/elon-server/domain"
	"github.com/socialgorithm/elon-server/physics"
)

const (
	delay = (1000 / 30) * time.Millisecond
)

var simulation Simulation
var carCount = 0
var testMode = false
var track domain.Track

// CreateSimulation creates a new simulation
func CreateSimulation(_carCount int) *Simulation {
	track = ReadTrack()
	carCount = _carCount
	return &Simulation{
		Track:       track,
		CarsChannel: make(chan []domain.Car),
		Engine:      physics.NewEngine(track, carCount),
	}
}

// Start starts the physics engine (run this in goroutine to async, don't put in the method)
func (simulation Simulation) Start(_testMode bool) {
	testMode = _testMode
	log.Println("Starting simulation")
	for {
		simulation.CarsChannel <- simulation.Engine.Next()
		time.Sleep(delay)
	}
}

// Restart restarts the simulation and all the cars in it
func (simulation Simulation) Restart() {
	simulation.Engine = physics.NewEngine(track, carCount)
}

// Input add an input to the simulation
func (simulation Simulation) Input(carIndex int, carControlState domain.CarControlState) {
	simulation.Engine.SetCtrl(0, carControlState)
}
