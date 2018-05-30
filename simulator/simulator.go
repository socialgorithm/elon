package simulator

import (
	"log"
	"math/rand"
	"time"

	"github.com/socialgorithm/elon-server/domain"
	"github.com/socialgorithm/elon-server/track"
)

// CreateSimulation creates and starts a new simulation with state channels for communication
func CreateSimulation() Simulation {
	log.Println("Creating simulation")
	return Simulation{
		Track:       track.GenTrack(),
		Cars:        []domain.Car{},
		CarsChannel: make(chan []domain.Car),
	}
}

// AddCars the specified number of cars to the simulation
func AddCars(simulation *Simulation, carCount int) {
	for i := 0; i < carCount; i++ {
		simulation.Cars = append(simulation.Cars, domain.Car{})
	}
}

// StartSimulation starts time within the simulation
func StartSimulation(simulation *Simulation) {
	log.Println("Starting simulation")
	go func() {
		for {
			advanceTick(simulation)
			time.Sleep(time.Second)
		}
	}()
}

func advanceTick(simulation *Simulation) {
	simulation.Cars = randomizeCarStates(simulation.Cars)
	simulation.CarsChannel <- simulation.Cars
}

func randomizeCarStates(cars []domain.Car) []domain.Car {
	newCars := make([]domain.Car, len(cars))
	for i := range cars {
		newCars[i] = domain.Car{
			CarState: domain.CarState{
				Position:  domain.Position{X: rand.Float64() * 1024, Y: rand.Float64() * 740},
				Direction: domain.Position{X: rand.Float64(), Y: rand.Float64()},
				Velocity:  1,
				Sensors: []domain.Sensor{
					domain.Sensor{Angle: 1, Distance: 1},
				},
			},
		}
	}
	return newCars
}
