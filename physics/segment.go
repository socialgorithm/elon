package physics

import "math"

// SegIntersections returns the points of intersection between l and q,
// or limits of intersection for colinear segments, returning up to
// 2 points as a slice
func SegIntersections(l [2][2]float64, m [2][2]float64) [][2]float64 {
	p := l[0]
	q := m[0]
	r := Vec2Subtract(l[1], l[0])
	s := Vec2Subtract(m[1], m[0])

	ut0 := Vec2Subtract(q, p)
	ut1 := Vec2CrossProduct(r, s)
	u0 := Vec2CrossProduct(ut0, r)

	t := Vec2CrossProduct(ut0, s) / ut1
	u := u0 / ut1

	// Colinear case
	if ut1 == 0 && u0 == 0 {
		rr := Vec2DotProduct(r, r)
		t0 := Vec2DotProduct(ut0, r) / rr
		t1 := t0 + (Vec2DotProduct(s, r) / rr)
		tr := [2]float64{t0, t1}

		if RangeIntersects(ZeroToOne, tr) {
			res := make([][2]float64, 2)

			if RangeContainsRange(ZeroToOne, tr) {
				res[0] = m[0]
				res[1] = m[1]
			} else if RangeContainsRange(tr, ZeroToOne) {
				res[0] = l[0]
				res[1] = l[1]
			} else {
				if RangeContainsValue(ZeroToOne, t0) {
					res[0] = l[1]
				} else {
					res[0] = p
				}

				if RangeContainsValue(ZeroToOne, t1) {
					res[1] = m[1]
				} else {
					res[1] = q
				}
			}

			return res
		}

		return [][2]float64{}
	}

	// Parallel case
	if ut1 == 0 && u0 != 0 {
		return [][2]float64{}
	}

	// Non-colinear intersecting case
	if ut1 != 0 && RangeContainsValue(ZeroToOne, u) && RangeContainsValue(ZeroToOne, t) {
		return [][2]float64{
			Vec2Add(p, Vec2Scale(r, t)),
		}
	}

	return [][2]float64{}
}

// SegIntersects returns whether l and q intersect
// TODO: deprecate + remove SegIntersects once boolean solution for RectCrosses is available
func SegIntersects(l [2][2]float64, m [2][2]float64) bool {
	p := l[0]
	q := m[0]
	r := Vec2Subtract(l[1], l[0])
	s := Vec2Subtract(m[1], m[0])

	ut0 := Vec2Subtract(q, p)
	ut1 := Vec2CrossProduct(r, s)
	u0 := Vec2CrossProduct(ut0, r)

	t := Vec2CrossProduct(ut0, s) / ut1
	u := u0 / ut1

	// Colinear case
	if ut1 == 0 && u0 == 0 {
		rr := Vec2DotProduct(r, r)
		t0 := Vec2DotProduct(ut0, r) / rr
		t1 := t0 + (Vec2DotProduct(s, r) / rr)
		tr := [2]float64{t0, t1}

		if RangeIntersects(ZeroToOne, tr) {
			return true
		}

		return false
	}

	// Parallel case
	if ut1 == 0 && u0 != 0 {
		return false
	}

	// Non-colinear intersecting case
	if ut1 != 0 && RangeContainsValue(ZeroToOne, u) && RangeContainsValue(ZeroToOne, t) {
		return true
	}

	return false
}

// SegDistanceToSq returns the square distance from the segment to a point
func SegDistanceToSq(l [2][2]float64, v [2]float64) float64 {
	ln := Vec2LenSq(Vec2Subtract(l[0], l[1]))

	if ln == 0 {
		return Vec2LenSq(Vec2Subtract(l[0], v))
	}

	m := Vec2Subtract(l[1], l[0])
	t := math.Max(0, math.Min(1, Vec2DotProduct(Vec2Subtract(v, l[0]), m)/ln))
	p := Vec2Add(l[0], Vec2Scale(m, t))
	return Vec2LenSq(Vec2Subtract(v, p))
}
