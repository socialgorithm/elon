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

func posSliceToPath(ps []domain.Position) Path {
	r := make(Path, len(ps))
	for i, p := range ps {
		r[i] = posToVec2(p)
	}
	return r
}

// NewEngine creates a new physics engine
func NewEngine(track domain.Track, count int) {
	e := Engine{
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

	for i := 0; i < count; i++ {
		e.State[i].Crashed = false
		e.State[i].Direction = domain.Position{X: 0, Y: 1}
		e.State[i].Position = track.Center[0]
		e.State[i].Sensors = make([]domain.Sensor, 5)
		e.State[i].Velocity = 1

		e.State[i].Sensors[0].Angle = -(0.5 * math.Pi)
		e.State[i].Sensors[1].Angle = -dims.Angle()
		e.State[i].Sensors[2].Angle = 0
		e.State[i].Sensors[3].Angle = dims.Angle()
		e.State[i].Sensors[4].Angle = 0.5 * math.Pi

		for _, s := range e.State[i].Sensors {
			s.Distance = sensorRange
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

	pv := posToVec2(state.Position)
	dv := posToVec2(state.Direction)

	if col < state.Velocity {
		np := pv.Add(dv.ScalarMultiple(col))

		return domain.CarState{
			Position:  domain.Position{X: np[0], Y: np[1]},
			Direction: state.Direction,
			Velocity:  0,
			Sensors:   newSensors,
			Crashed:   true,
		}
	}

	na := normaliseAngle(dv.Angle() + steeringRate*ctrl.Steering)

	return domain.CarState{
		Position: domain.Position{
			X: state.Position.X + state.Direction.X*state.Velocity,
			Y: state.Position.Y + state.Direction.Y*state.Velocity,
		},
		Direction: domain.Position{
			X: math.Cos(na),
			Y: math.Sin(na),
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
