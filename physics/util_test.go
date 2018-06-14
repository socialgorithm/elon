package physics_test

import (
	"math"
	"testing"

	"github.com/socialgorithm/elon-server/domain"
	"github.com/socialgorithm/elon-server/physics"
	"github.com/stretchr/testify/assert"
)

const (
	piOverSix = 0.16667 * math.Pi
)

var (
	utilRes interface{}
)

func TestNormaliseRadiansWhereBelow(t *testing.T) {
	assert := assert.New(t)

	r := physics.NormaliseRadians(-(0.5 * math.Pi))

	assert.InEpsilon(1.5*math.Pi, r, 0.01)
}

func TestNormaliseRadiansWhereAbove(t *testing.T) {
	assert := assert.New(t)

	r := physics.NormaliseRadians(2.5 * math.Pi)

	assert.InEpsilon(0.5*math.Pi, r, 0.01)
}

func TestNormaliseRadiansWhereWithin(t *testing.T) {
	assert := assert.New(t)

	r := physics.NormaliseRadians(math.Pi)

	assert.InEpsilon(math.Pi, r, 0.01)
}

func BenchmarkNormaliseRadians(b *testing.B) {
	var r float64
	for i := 0; i < b.N; i++ {
		r = physics.NormaliseRadians(math.Pi)
	}
	utilRes = r
}

func TestCapValueWhereContains(t *testing.T) {
	assert := assert.New(t)

	r := physics.CapValue(1, 0, 5)

	assert.Equal(1., r)
}

func TestCapValueWhereBelow(t *testing.T) {
	assert := assert.New(t)

	r := physics.CapValue(-1, 0, 5)

	assert.Equal(0., r)
}

func TestCapValueWhereAbove(t *testing.T) {
	assert := assert.New(t)

	r := physics.CapValue(6, 0, 5)

	assert.Equal(5., r)
}

func BenchmarkCapValue(b *testing.B) {
	var r float64
	for i := 0; i < b.N; i++ {
		r = physics.CapValue(1, 0, 5)
	}
	utilRes = r
}

func TestDPosToVec2(t *testing.T) {
	assert := assert.New(t)

	r := physics.DPosToVec2(domain.Position{X: 2, Y: 3})

	assert.Equal([2]float64{2, 3}, r)
}

func BenchmarkDPosToVec2(b *testing.B) {
	var r [2]float64
	for i := 0; i < b.N; i++ {
		r = physics.DPosToVec2(domain.Position{X: 2, Y: 3})
	}
	utilRes = r
}

func TestAngleToUnitVec2(t *testing.T) {
	assert := assert.New(t)

	r := physics.AngleToUnitVec2(piOverSix)

	assert.InEpsilon(0.866, r[0], 0.01)
	assert.InEpsilon(0.5, r[1], 0.01)
}

func BenchmarkAngleToUnitVec2(b *testing.B) {
	var r [2]float64
	for i := 0; i < b.N; i++ {
		r = physics.AngleToUnitVec2(piOverSix)
	}
	utilRes = r
}
