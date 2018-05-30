package simulator

import (
	"log"
	"time"

	"github.com/socialgorithm/elon-server/domain"
	"github.com/socialgorithm/elon-server/physics"
	"github.com/socialgorithm/elon-server/track"
)

var simulation Simulation

// PrepareSimulation Starts a simulation, returning state and control channels
func PrepareSimulation() Simulation {
	log.Println("Preparing simulation")
	track := track.GenTrack()
	carControlStateChannel := make(chan domain.CarControlState)
	carStateChannel := make(chan []domain.CarState, 5)
	return Simulation{
		Track:                   track,
		CarStatesChannel:        carStateChannel,
		CarControlStateReceiver: carControlStateChannel,
		Engine:                  physics.NewEngine(track, 5),
	}
}

// StartSimulation starts the physics engine (run this in goroutine to async, don't put in the method)
func StartSimulation(simulation Simulation) {
	log.Println("Starting simulation")
	for {
		simulation.CarStatesChannel <- simulation.Engine.Next()
		time.Sleep(time.Second)
	}
}
