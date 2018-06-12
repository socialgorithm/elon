package simulator

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/socialgorithm/elon-server/domain"
	"github.com/socialgorithm/elon-server/physics"
	"github.com/socialgorithm/elon-server/track"
)

var simulation Simulation

// CreateSimulation creates a new simulation
func CreateSimulation(carCount int) *Simulation {
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
	sensorArc := math.Pi
	var sensors [sensorCount + 1]domain.Sensor
	sensorAngleIncrement := sensorArc / sensorCount
	for i := 0; i <= sensorCount; i++ {
		sensors[i] = domain.Sensor{
			Angle:    -sensorArc/2 + sensorAngleIncrement*float64(i),
			Distance: rand.Float64()*maxSensorDistance + minSensorDistance,
		}
	}
	return sensors[0:len(sensors)]
}
