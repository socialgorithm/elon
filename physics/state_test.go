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
	}
	stateSeg = [2][2]float64{{0, 10}, {10, 0}}
)

func TestStateCheckWithSingleEdgeInfront(t *testing.T) {
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
