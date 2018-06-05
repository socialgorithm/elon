package simulator

import (
	"log"
	"math/rand"
	"time"

	"github.com/socialgorithm/elon-server/domain"
	"github.com/socialgorithm/elon-server/physics"
	"github.com/socialgorithm/elon-server/track"
)

var simulation Simulation

// CreateSimulation creates a new simulation
func CreateSimulation(carCount int) *Simulation {
	log.Println("Preparing simulation")
	track := track.GenTrack()
	return &Simulation{
		Track:       track,
		Cars:        make([]domain.Car, carCount),
		CarsChannel: make(chan []domain.Car),
		Engine:      physics.NewEngine(track, carCount),
	}
}

// Start starts the physics engine (run this in goroutine to async, don't put in the method)
func (simulation Simulation) Start(testMode bool) {
	log.Println("Starting simulation")
	for {
		if testMode {
			for idx := range simulation.Engine.State {
				simulation.Engine.SetCtrl(
					idx,
					domain.CarControlState{Throttle: rand.Float64(), Steering: rand.Float64()},
				)
			}
		}
		simulation.CarsChannel <- simulation.Engine.Next()
		time.Sleep(time.Second)
	}
}
