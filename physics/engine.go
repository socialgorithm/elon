package physics

import (
	"sync"

	"github.com/socialgorithm/elon-server/domain"
)

var (
	sensorAngles = [5]float64{0, 0.5 * Pi, Pi, 1.5 * Pi, 2 * Pi}
)

// Engine describes a physics engine
type Engine struct {
	state    []State
	locks    []sync.Mutex
	segments [][2][2]float64
}

// NewEngine creates a new physics engine
func NewEngine(track domain.Track, count int) Engine {
	engine := Engine{
		state: make([]State, count),
		locks: make([]sync.Mutex, count),
	}

	centre0 := DPosToVec2(track.Center[0])
	centre1 := DPosToVec2(track.Center[1])
	startAngle := Vec2UnitToAngle(Vec2Normalise(Vec2Subtract(centre1, centre0)))

	for idx := 0; idx < count; idx++ {
		engine.state[idx].Crashed = false
		engine.state[idx].Angle = startAngle
		engine.state[idx].Position = centre0
		engine.state[idx].Velocity = 1
	}

	return engine
}

func (engine Engine) nextForIndex(idx int) domain.Car {
	engine.locks[idx].Lock()
	defer engine.locks[idx].Unlock()

	sv := engine.state[idx].Update(engine.segments)
	s := engine.state[idx]
	uv := AngleToUnitVec2(s.Angle)

	r := domain.CarState{
		Position:  domain.Position{X: s.Position[0], Y: s.Position[1]},
		Direction: domain.Position{X: uv[0], Y: uv[1]},
		Velocity:  s.Velocity,
		Sensors:   make([]domain.Sensor, 5),
		Crashed:   false,
	}

	for idx, sa := range sensorAngles {
		r.Sensors[idx].Angle = sa
		r.Sensors[idx].Distance = sv[idx]
	}

	return domain.Car{
		CarState: r,
		CarControlState: domain.CarControlState{
			Throttle: s.Throttle,
			Steering: s.Steering,
		},
	}
}

// SetCtrl updates the control state for a given index
func (engine Engine) SetCtrl(idx int, ctrl domain.CarControlState) {
	engine.locks[idx].Lock()
	defer engine.locks[idx].Unlock()
	engine.state[idx].Throttle = ctrl.Throttle
	engine.state[idx].Steering = ctrl.Steering
}

// Next proceeds to the next state
func (engine Engine) Next() []domain.Car {
	nxt := make([]domain.Car, len(engine.state))
	for idx := range engine.state {
		nxt[idx] = engine.nextForIndex(idx)
	}
	return nxt
}
