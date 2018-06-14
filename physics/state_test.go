package physics_test

import (
	"math"
	"testing"

	"github.com/socialgorithm/elon-server/physics"
	"github.com/stretchr/testify/assert"
)

var (
	stateRes       interface{}
	sqrtMaxFloat64 = math.Sqrt(math.MaxFloat64)
	state          = physics.State{
		Position: [2]float64{2, 2},
		Angle:    0.25 * math.Pi,
		Velocity: 2,
		Throttle: 2,
		Steering: math.Pi,
	}
	stateDown = physics.State{
		Position: [2]float64{2, -2},
		Angle:    0.75 * math.Pi,
		Velocity: 2,
	}
	stateSeg           = [2][2]float64{{0, 10}, {10, 0}}
	stateDownSeg       = [2][2]float64{{0, -10}, {10, 0}}
	stateSegOut        = [2][2]float64{{-10, -10}, {0, -20}}
	stateSegFirstCont  = [2][2]float64{{4.9, 5.1}, {10, 0}}
	stateSegSecondCont = [2][2]float64{{0, 10}, {5.1, 4.9}}
	largeSegSet        = MakeLargeSegmentSet()
)

func TestStateCheckWithSingleEdgeInFront(t *testing.T) {
	assert := assert.New(t)

	hasCrashed, crashDistance, sensorValues := state.Check([][2][2]float64{stateSeg})

	assert.True(hasCrashed)
	assert.InEpsilon(1.243, crashDistance, 0.01)

	assert.InEpsilon(sqrtMaxFloat64, sensorValues[0], 0.01)
	assert.InEpsilon(3.171, sensorValues[1], 0.01)
	assert.InEpsilon(1.243, sensorValues[2], 0.01)
	assert.InEpsilon(3.171, sensorValues[3], 0.01)
	assert.InEpsilon(sqrtMaxFloat64, sensorValues[4], 0.01)
}

func TestStateCheckWithSingleEdgeInFrontPointingDown(t *testing.T) {
	assert := assert.New(t)

	hasCrashed, crashDistance, sensorValues := stateDown.Check([][2][2]float64{stateDownSeg})

	assert.True(hasCrashed)
	assert.InEpsilon(1.243, crashDistance, 0.01)

	assert.InEpsilon(sqrtMaxFloat64, sensorValues[0], 0.01)
	assert.InEpsilon(3.171, sensorValues[1], 0.01)
	assert.InEpsilon(1.243, sensorValues[2], 0.01)
	assert.InEpsilon(3.171, sensorValues[3], 0.01)
	assert.InEpsilon(sqrtMaxFloat64, sensorValues[4], 0.01)
}

func TestStateCheckWithFirstPointContained(t *testing.T) {
	assert := assert.New(t)

	hasCrashed, crashDistance, sensorValues := state.Check([][2][2]float64{stateSegFirstCont})

	assert.True(hasCrashed)
	assert.InEpsilon(1.243, crashDistance, 0.01)

	assert.InEpsilon(sqrtMaxFloat64, sensorValues[0], 0.01)
	assert.InEpsilon(sqrtMaxFloat64, sensorValues[1], 0.01)
	assert.InEpsilon(1.243, sensorValues[2], 0.01)
	assert.InEpsilon(3.171, sensorValues[3], 0.01)
	assert.InEpsilon(sqrtMaxFloat64, sensorValues[4], 0.01)
}

func TestStateCheckWithSecondPointContained(t *testing.T) {
	assert := assert.New(t)

	hasCrashed, crashDistance, sensorValues := state.Check([][2][2]float64{stateSegSecondCont})

	assert.True(hasCrashed)
	assert.InEpsilon(1.243, crashDistance, 0.01)

	assert.InEpsilon(sqrtMaxFloat64, sensorValues[0], 0.01)
	assert.InEpsilon(3.171, sensorValues[1], 0.01)
	assert.InEpsilon(1.243, sensorValues[2], 0.01)
	assert.InEpsilon(sqrtMaxFloat64, sensorValues[3], 0.01)
	assert.InEpsilon(sqrtMaxFloat64, sensorValues[4], 0.01)
}

func BenchmarkStateCheckWithSingleEdgeInfront(b *testing.B) {
	var res0 bool
	var res1 float64
	var res2 [5]float64
	for i := 0; i < b.N; i++ {
		res0, res1, res2 = state.Check([][2][2]float64{stateSeg})
	}
	stateRes = res0
	stateRes = res1
	stateRes = res2
}

func BenchmarkStateCheckWithLargeSegmentSet(b *testing.B) {
	var res0 bool
	var res1 float64
	var res2 [5]float64
	for i := 0; i < b.N; i++ {
		res0, res1, res2 = state.Check(largeSegSet)
	}
	stateRes = res0
	stateRes = res1
	stateRes = res2
}

func TestStateUpdateWhereAlreadyCrashed(t *testing.T) {
	assert := assert.New(t)
	stateCpy := state
	stateCpy.Crashed = true
	stateCpy2 := stateCpy

	stateCpy.Update([][2][2]float64{stateSegOut})

	assert.Equal(stateCpy, stateCpy2)
}

func TestStateUpdateWithNoCollision(t *testing.T) {
	assert := assert.New(t)
	stateCpy := state

	res := stateCpy.Update([][2][2]float64{stateSegOut})

	for _, v := range res {
		assert.Equal(sqrtMaxFloat64, v)
	}

	assert.False(stateCpy.Crashed)
	assert.InEpsilon(2+math.Sqrt(2), stateCpy.Position[0], 0.01)
	assert.InEpsilon(2+math.Sqrt(2), stateCpy.Position[1], 0.01)
	assert.InEpsilon(2+(2*physics.AccelerationRate), stateCpy.Velocity, 0.01)
	assert.InEpsilon((0.25*math.Pi)+(physics.SteeringRate*math.Pi), stateCpy.Angle, 0.01)
}

func TestStateUpdateWithCollision(t *testing.T) {
	assert := assert.New(t)
	stateCpy := state

	res := stateCpy.Update([][2][2]float64{stateSeg})

	for _, v := range res {
		assert.Equal(0., v)
	}

	assert.True(stateCpy.Crashed)
	assert.InEpsilon(2.88, stateCpy.Position[0], 0.01)
	assert.InEpsilon(2.88, stateCpy.Position[1], 0.01)
	assert.InEpsilon(2, stateCpy.Velocity, 0.01)
	assert.InEpsilon(0.25*math.Pi, stateCpy.Angle, 0.01)
}

func MakeLargeSegmentSet() [][2][2]float64 {
	o := 5
	s := 1000
	r := make([][2][2]float64, s)
	for i := 0; i < o; i++ {
		r[i] = stateSeg
	}
	for i := o; i < s; i++ {
		r[i] = stateSegOut
	}
	return r
}
