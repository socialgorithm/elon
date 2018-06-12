package physics

var (
	// ZeroToOne is the value range from 0 to 1
	ZeroToOne = [2]float64{0, 1}
)

// RangeContainsValue checks whether a range contains a value
func RangeContainsValue(r [2]float64, v float64) bool {
	return r[0] <= v && v <= r[1]
}

// RangeContainsRange checks whether a range entirely contains another range
func RangeContainsRange(r [2]float64, s [2]float64) bool {
	return s[0] >= r[0] && s[1] <= r[1]
}

// RangeIntersects checks whether range r intersects range s
func RangeIntersects(r [2]float64, s [2]float64) bool {
	return RangeContainsValue(r, s[0]) || RangeContainsValue(r, s[1]) ||
		s[0] < r[0] && s[1] > r[1] ||
		r[0] < s[0] && r[1] > r[1]
}
