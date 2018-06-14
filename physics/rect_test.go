package physics_test

import (
	"testing"

	"github.com/socialgorithm/elon-server/physics"
	"github.com/stretchr/testify/assert"
)

var (
	rectRes interface{}
	rect    = [2][2]float64{{-2, -3}, {2, 3}}
)

func TestRectContainsWhereFalse(t *testing.T) {
	assert := assert.New(t)

	r := physics.RectContains(rect, [2]float64{4, 4})

	assert.False(r)
}

func TestRectContainsWhereTrue(t *testing.T) {
	assert := assert.New(t)

	r := physics.RectContains(rect, [2]float64{1, 1})

	assert.True(r)
}

func BenchmarkRectContains(b *testing.B) {
	var res bool
	for i := 0; i < b.N; i++ {
		res = physics.RectContains(rect, [2]float64{0, 0})
	}
	rectRes = res
}

func TestRectCrossesTrueHorizontally(t *testing.T) {
	assert := assert.New(t)

	r := physics.RectCrosses(rect, [2][2]float64{{-3, 0}, {3, 0}})

	assert.True(r)
}

func TestRectCrossesTrueVertically(t *testing.T) {
	assert := assert.New(t)

	r := physics.RectCrosses(rect, [2][2]float64{{0, -4}, {0, 4}})

	assert.True(r)
}

func TestRectCrossesTrueDiagonallyBLToTR(t *testing.T) {
	assert := assert.New(t)

	r := physics.RectCrosses(rect, [2][2]float64{{-4, -6}, {4, 6}})

	assert.True(r)
}

func TestRectCrossesTrueDiagonallyTLToBR(t *testing.T) {
	assert := assert.New(t)

	r := physics.RectCrosses(rect, [2][2]float64{{-4, 6}, {4, -6}})

	assert.True(r)
}

func TestRectCrossesFalseWhereContains(t *testing.T) {
	assert := assert.New(t)

	r := physics.RectCrosses(rect, [2][2]float64{{0, -1}, {0, 1}})

	assert.False(r)
}

func TestRectCrossesFalseWhereSeparate(t *testing.T) {
	assert := assert.New(t)

	r := physics.RectCrosses(rect, [2][2]float64{{0, -5}, {0, -4}})

	assert.False(r)
}

func BenchmarkRectCrosses(b *testing.B) {
	var res bool
	for i := 0; i < b.N; i++ {
		res = physics.RectCrosses(rect, [2][2]float64{{0, 0}, {0, -6}})
	}
	rectRes = res
}
