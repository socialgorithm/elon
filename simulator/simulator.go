package simulator

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
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
		Cars:        genCars(carCount, track),
		CarsChannel: make(chan []domain.Car),
		Engine:      physics.NewEngine(track, carCount),
	}
}

// Start starts the physics engine (run this in goroutine to async, don't put in the method)
func (simulation Simulation) Start(testMode bool) {
	log.Println("Starting simulation")
	for {
		simulation.CarsChannel <- simulation.Engine.Next(simulation.Cars)
		time.Sleep(time.Second)
	}
}

// Input add an input to the simulation
func (simulation Simulation) Input(carIndex int, carControlState domain.CarControlState) {
	simulation.Cars[carIndex].CarControlState = carControlState
}

func genCars(carCount int, track domain.Track) []domain.Car {
	cars := make([]domain.Car, carCount)
	centre0 := pixel.V(track.Center[0].X, track.Center[0].Y)
	centre1 := pixel.V(track.Center[1].X, track.Center[1].Y)
	startAngle := centre1.Sub(centre0).Unit()
	for i := range cars {
		cars[i].CarState = domain.CarState{
			Crashed:   false,
			Position:  track.Center[0],
			Direction: domain.Position{X: startAngle.X, Y: startAngle.Y},
			Velocity:  0,
			// using 4 sensors for now
			Sensors: genRandomSensorData(),
		}
	}

	return cars
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
