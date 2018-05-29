package simulator

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/socialgorithm/elon-server/domain"
	"github.com/socialgorithm/elon-server/track"
)

var simulation Simulation

// StartSimulation Starts a simulation, returning state and control channels
func PrepareSimulation() Simulation {
	track := track.GenTrack()
	carControlStateChannel := make(chan domain.CarControlState)
	carStateChannel := make(chan []domain.CarState, 5)
	return Simulation{
		Track:                   track,
		CarStatesChannel:        carStateChannel,
		CarControlStateReceiver: carControlStateChannel,
	}
}

// StartSimulation starts the physics engine
func StartSimulation(simulation Simulation) {
	for {
		simulation.CarStatesChannel <- genRandomCarState()
		time.Sleep(time.Second)
		runtime.Gosched()
	}
}

func genRandomCarState() []domain.CarState {
	carStates := make([]domain.CarState, 1, 1)
	carStates[0] = domain.CarState{
		Position:  domain.Position{X: rand.Float64() * 1024, Y: rand.Float64() * 740},
		Direction: domain.Position{X: rand.Float64(), Y: rand.Float64()},
		Velocity:  1,
		Sensors: []domain.Sensor{
			domain.Sensor{Angle: 1, Distance: 1},
		},
	}
	return carStates
}
