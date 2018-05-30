package simulator

import (
	"math"
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
		Direction: domain.Position{X: rand.Float64() - 0.5, Y: rand.Float64() - 0.5},
		Velocity:  1,
		// using 4 sensors for now
		Sensors: genRandomSensorData(),
	}
	return carStates
}

func genRandomSensorData() []domain.Sensor {
	const sensorCount = 4
	minSensorDistance := 10.0
	maxSensorDistance := 50.0
	var sensors [sensorCount + 1]domain.Sensor
	sensorAngleIncrement := math.Pi / sensorCount
	for i := 0; i <= sensorCount; i++ {
		sensors[i] = domain.Sensor{
			Angle:    -math.Pi/2 + sensorAngleIncrement*float64(i),
			Distance: rand.Float64()*maxSensorDistance + minSensorDistance,
		}
	}
	return sensors[0:len(sensors)]
}
