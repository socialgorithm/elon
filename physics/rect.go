package physics

// RectContains determines whether a [[min_x, min_y], [max_x, max_y]] rect contains point v
func RectContains(r [2][2]float64, v [2]float64) bool {
	return v[0] >= r[0][0] &&
		v[1] >= r[0][1] &&
		v[0] <= r[1][0] &&
		v[1] <= r[1][1]
}

// RectCrosses determines whether a line segment crosses the rect
func RectCrosses(r [2][2]float64, s [2][2]float64) bool {
	// TODO: Should optimise with a boolean replacement, which should be ~10-15ns/op vs ~50-200ns/op

	ls := [4][2][2]float64{
		{{r[0][0], r[0][1]}, {r[0][0], r[1][1]}},
		{{r[0][0], r[1][1]}, {r[1][0], r[1][1]}},
		{{r[1][0], r[1][1]}, {r[1][0], r[0][1]}},
		{{r[1][0], r[0][1]}, {r[0][0], r[0][1]}},
	}

	for _, l := range ls {
		if SegIntersects(l, s) {
			return true
		}
	}

	return false
}

// RectIntersects determines whether a line segment intersects the rect
func RectIntersects(r [2][2]float64, s [2][2]float64) bool {
	return RectContains(r, s[0]) || RectContains(r, s[1]) || RectCrosses(r, s)
}
