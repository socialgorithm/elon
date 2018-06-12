package physics_test

import (
	"math"
	"testing"

	"github.com/socialgorithm/elon-server/physics"
	"github.com/stretchr/testify/assert"
)

var (
	segRes     interface{}
	sqrt8      = math.Sqrt(8)
	seg        = [2][2]float64{{-3, -3}, {3, 3}}
	segSm      = [2][2]float64{{-2, -2}, {2, 2}}
	segDiagMv1 = [2][2]float64{{-4, -4}, {2, 2}}
	segDiagMv7 = [2][2]float64{{-10, -10}, {-4, -4}}
	segPara    = [2][2]float64{{-4, -3}, {2, 3}}
	segInv     = [2][2]float64{{-3, 3}, {3, -3}}
	segOther   = [2][2]float64{{1, 6}, {2, 9}}
	segPoint   = [2][2]float64{{0, 0}, {0, 0}}
)

func TestSegIntersectionsWhereColinearAndFullyIntersecting(t *testing.T) {
	assert := assert.New(t)

	r0 := physics.SegIntersections(seg, segSm)
	r1 := physics.SegIntersections(segSm, seg)
	e := [][2]float64{segSm[0], segSm[1]}

	assert.Contains(r0, e[0])
	assert.Contains(r0, e[1])
	assert.Contains(r1, e[0])
	assert.Contains(r1, e[1])
}

func BenchmarkSegIntersectionsWhereColinearAndOuterFullyContainsInner(b *testing.B) {
	var r [][2]float64
	for i := 0; i < b.N; i++ {
		r = physics.SegIntersections(seg, segSm)
	}
	segRes = r
}

func BenchmarkSegIntersectionsWhereColinearAndInnerFullContainsOuter(b *testing.B) {
	var r [][2]float64
	for i := 0; i < b.N; i++ {
		r = physics.SegIntersections(segSm, seg)
	}
	segRes = r
}

func TestSegIntersectionsWhereColinearAndPartiallyIntersecting(t *testing.T) {
	assert := assert.New(t)

	e := [][2]float64{{-3, -3}, {2, 2}}

	r0 := physics.SegIntersections(seg, segDiagMv1)
	r1 := physics.SegIntersections(segDiagMv1, seg)

	assert.Contains(r0, e[0])
	assert.Contains(r0, e[1])
	assert.Contains(r1, e[0])
	assert.Contains(r1, e[1])
}

func BenchmarkSegIntersectionsWhereColinearAndPartiallyIntersecting(b *testing.B) {
	var r [][2]float64
	for i := 0; i < b.N; i++ {
		r = physics.SegIntersections(seg, segDiagMv1)
	}
	segRes = r
}

func BenchmarkSegIntersectionsWhereColinearAndPartiallyIntersectingInversely(b *testing.B) {
	var r [][2]float64
	for i := 0; i < b.N; i++ {
		r = physics.SegIntersections(segDiagMv1, seg)
	}
	segRes = r
}

func TestSegIntersectionsWhereColinearAndNonIntersecting(t *testing.T) {
	assert := assert.New(t)

	e := [][2]float64{}
	r := physics.SegIntersections(seg, segDiagMv7)

	assert.Equal(e, r)
}

func BenchmarkSegIntersectionsWhereColinearAndNonIntersecting(b *testing.B) {
	var r [][2]float64
	for i := 0; i < b.N; i++ {
		r = physics.SegIntersections(seg, segDiagMv7)
	}
	segRes = r
}

func TestSegIntersectionsWhereParallel(t *testing.T) {
	assert := assert.New(t)

	e := [][2]float64{}
	r := physics.SegIntersections(seg, segPara)

	assert.Equal(e, r)
}

func BenchmarkSegIntersectionsWhereParallel(b *testing.B) {
	var r [][2]float64
	for i := 0; i < b.N; i++ {
		r = physics.SegIntersections(seg, segPara)
	}
	segRes = r
}

func TestSegIntersectionsWhereIntersecting(t *testing.T) {
	assert := assert.New(t)

	e := [][2]float64{{0, 0}}
	r := physics.SegIntersections(seg, segInv)

	assert.Equal(e, r)
}

func BenchmarkSegIntersectionsWhereIntersecting(b *testing.B) {
	var r [][2]float64
	for i := 0; i < b.N; i++ {
		r = physics.SegIntersections(seg, segInv)
	}
	segRes = r
}

func TestSegIntersectionsWhereNonIntersecting(t *testing.T) {
	assert := assert.New(t)

	e := [][2]float64{}
	r := physics.SegIntersections(seg, segOther)

	assert.Equal(e, r)
}

func BenchmarkSegIntersectionsWhereNonIntersecting(b *testing.B) {
	var r [][2]float64
	for i := 0; i < b.N; i++ {
		r = physics.SegIntersections(seg, segOther)
	}
	segRes = r
}

func TestSegDistanceToSqWithDistance(t *testing.T) {
	assert := assert.New(t)

	r := physics.SegDistanceToSq(seg, [2]float64{-1, 1})

	assert.Equal(2., r)
}

func TestSegDistanceToSqWithNoDistance(t *testing.T) {
	assert := assert.New(t)

	r := physics.SegDistanceToSq(seg, [2]float64{0, 0})

	assert.Equal(0., r)
}

func BenchmarkSegDistanceToSqWithFullSeg(b *testing.B) {
	var r float64
	for i := 0; i < b.N; i++ {
		physics.SegDistanceToSq(seg, [2]float64{-1, -1})
	}
	segRes = r
}

func TestSegDistanceToSqWithPointSeg(t *testing.T) {
	assert := assert.New(t)

	r := physics.SegDistanceToSq(segPoint, [2]float64{-1, -1})

	assert.Equal(2., r)
}

func BenchmarkSegDistanceToSqWithPointSeg(b *testing.B) {
	var r float64
	for i := 0; i < b.N; i++ {
		physics.SegDistanceToSq(segPoint, [2]float64{-1, -1})
	}
	segRes = r
}
