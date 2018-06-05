package simulator

import (
	"log"
	"time"

	"github.com/socialgorithm/elon-server/domain"
	"github.com/socialgorithm/elon-server/physics"
	"github.com/socialgorithm/elon-server/track"
)

var simulation Simulation

// CreateSimulation creates a new simulation
func CreateSimulation(count int) Simulation {
	log.Println("Preparing simulation")
	track := track.GenTrack()
	return Simulation{
		Track:       track,
		Cars:        []domain.Car{},
		CarsChannel: make(chan []domain.Car),
		Engine:      physics.NewEngine(track, count),
	}
}

// Start starts the physics engine (run this in goroutine to async, don't put in the method)
func (simulation Simulation) Start() {
	log.Println("Starting simulation")
	for {
		simulation.CarsChannel <- simulation.Engine.Next()
		time.Sleep(time.Second)
	}
}
