package physics

import (
	"math"
	"sync"

	"github.com/socialgorithm/elon-server/domain"
)

const (
	sensorRange  float64 = 2
	steeringRate float64 = 0.05
	accelRate    float64 = 0.05
	maxVelocity  float64 = 5
)

var (
	dims  = Vec2{2, 4}
	ndims = dims.Normalise()
)

// Engine describes a physics engine
type Engine struct {
	Locks []sync.Mutex
	DistCalc
}

func posToVec2(pos domain.Position) Vec2 {
	return Vec2{pos.X, pos.Y}
}

func posSliceToPath(positions []domain.Position) Path {
	res := make(Path, len(positions))
	for idx, position := range positions {
		res[idx] = posToVec2(position)
	}
	return res
}

// NewEngine creates a new physics engine
func NewEngine(track domain.Track, count int) Engine {
	engine := Engine{
		DistCalc: NewDistCalc(
			dims,
			[]Path{
				posSliceToPath(track.InnerSide),
				posSliceToPath(track.OuterSide),
			},
			sensorRange,
		),
		Locks: make([]sync.Mutex, count),
	}

	return engine
}

// SetCtrl updates the control state for a given index
// func (engine Engine) SetCtrl(idx int, ctrl domain.CarControlState) {
// 	engine.Locks[idx].Lock()
// 	defer engine.Locks[idx].Unlock()
// 	engine.ControlState[idx] = ctrl
// }

func (engine Engine) nextForCar(idx int, car domain.Car) domain.Car {
	engine.Locks[idx].Lock()
	defer engine.Locks[idx].Unlock()

	state := car.CarState
	ctrl := car.CarControlState

	if state.Crashed {
		return car
	}

	col, sens := engine.Execute(posToVec2(state.Position), posToVec2(state.Direction).Angle())

	newSensors := make([]domain.Sensor, len(state.Sensors))
	for idx, sensor := range newSensors {
		sensor.Distance = sens[idx]
		sensor.Angle = state.Sensors[idx].Angle
	}

	posVec := posToVec2(state.Position)
	distVec := posToVec2(state.Direction)

	if col < state.Velocity {
		newPos := posVec.Add(distVec.ScalarMultiple(col))

		return domain.Car{
			CarState: domain.CarState{
				Position:  domain.Position{X: newPos[0], Y: newPos[1]},
				Direction: state.Direction,
				Velocity:  0,
				Sensors:   newSensors,
				Crashed:   true,
			},
		}
	}

	newAngle := normaliseAngle(distVec.Angle() + steeringRate*ctrl.Steering)

	return domain.Car{
		CarState: domain.CarState{
			Position: domain.Position{
				X: state.Position.X + state.Direction.X*state.Velocity,
				Y: state.Position.Y + state.Direction.Y*state.Velocity,
			},
			Direction: domain.Position{
				X: math.Cos(newAngle),
				Y: math.Sin(newAngle),
			},
			Velocity: limit(state.Velocity+(ctrl.Throttle*accelRate), 0, maxVelocity),
			Sensors:  newSensors,
			Crashed:  false,
		},
	}
}

// Next proceeds to the next state
func (engine Engine) Next(cars []domain.Car) []domain.Car {
	// TODO: make parallel, is already threadsafe with mutex slice
	nxt := make([]domain.Car, len(cars))
	for idx := range cars {
		nxt[idx] = engine.nextForCar(idx, cars[idx])
	}
	return nxt
}

func limit(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
