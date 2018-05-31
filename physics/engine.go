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
	State        []domain.CarState
	ControlState []domain.CarControlState
	Locks        []sync.Mutex
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
func NewEngine(track domain.Track, count int) {
	engine := Engine{
		State:        make([]domain.CarState, count),
		ControlState: make([]domain.CarControlState, count),
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

	for idx := 0; idx < count; idx++ {
		engine.State[idx].Crashed = false
		engine.State[idx].Direction = domain.Position{X: 0, Y: 1}
		engine.State[idx].Position = track.Center[0]
		engine.State[idx].Sensors = make([]domain.Sensor, 5)
		engine.State[idx].Velocity = 1

		engine.State[idx].Sensors[0].Angle = -(0.5 * math.Pi)
		engine.State[idx].Sensors[1].Angle = -dims.Angle()
		engine.State[idx].Sensors[2].Angle = 0
		engine.State[idx].Sensors[3].Angle = dims.Angle()
		engine.State[idx].Sensors[4].Angle = 0.5 * math.Pi

		for _, sensor := range engine.State[idx].Sensors {
			sensor.Distance = sensorRange
		}
	}
}

// SetCtrl updates the control state for a given index
func (engine Engine) SetCtrl(idx int, ctrl domain.CarControlState) {
	engine.Locks[idx].Lock()
	defer engine.Locks[idx].Unlock()
	engine.ControlState[idx] = ctrl
}

func (engine Engine) nextForIndex(idx int) domain.CarState {
	engine.Locks[idx].Lock()
	defer engine.Locks[idx].Unlock()

	state := engine.State[idx]
	ctrl := engine.ControlState[idx]

	if state.Crashed {
		return state
	}

	col, sens := engine.Execute(posToVec2(state.Position), posToVec2(state.Direction).Angle())

	newSensors := make([]domain.Sensor, 5)
	for idx, sensor := range newSensors {
		sensor.Distance = sens[idx]
		sensor.Angle = state.Sensors[idx].Angle
	}

	posVec := posToVec2(state.Position)
	distVec := posToVec2(state.Direction)

	if col < state.Velocity {
		newPos := posVec.Add(distVec.ScalarMultiple(col))

		return domain.CarState{
			Position:  domain.Position{X: newPos[0], Y: newPos[1]},
			Direction: state.Direction,
			Velocity:  0,
			Sensors:   newSensors,
			Crashed:   true,
		}
	}

	newAngle := normaliseAngle(distVec.Angle() + steeringRate*ctrl.Steering)

	return domain.CarState{
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
	}
}

// Next proceeds to the next state
func (engine Engine) Next() []domain.CarState {
	// TODO: make parallel, is already threadsafe with mutex slice
	nxt := make([]domain.CarState, len(engine.State))
	for idx := range engine.State {
		nxt[idx] = engine.nextForIndex(idx)
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
