package simulator

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/socialgorithm/elon-server/domain"
	"github.com/socialgorithm/elon-server/track"
)

var simulation Simulation

// StartSimulation Starts a simulation, returning state and control channels
func StartSimulation() Simulation {
	track := track.GenTrack()
	carControlStateChannel := make(chan domain.CarControlState)
	carStateChannel := startCarStateStream(track)
	return Simulation{
		Track:                   track,
		CarStateEmitter:         carStateChannel,
		CarControlStateReceiver: carControlStateChannel,
	}
}

// simulation ticks

func startCarStateStream(track domain.Track) <-chan domain.CarState {
	carStateChannel := make(chan domain.CarState)
	go func() {
		for {
			log.Println("Generating new car state")
			carStateChannel <- genRandomCarState()
			time.Sleep(time.Second)
			runtime.Gosched()
		}
	}()

	return carStateChannel
}

func genRandomCarState() domain.CarState {
	return domain.CarState{
		Position: domain.Position{X: rand.Float64() * 1024, Y: rand.Float64() * 740},
		Velocity: 1,
		Sensors: []domain.Sensor{
			domain.Sensor{Angle: 1, Distance: 1},
		},
	}
}
