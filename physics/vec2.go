package physics

import "math"

// Vec2DotProduct gets the dot product of p and q
func Vec2DotProduct(p [2]float64, q [2]float64) float64 {
	return (p[0] * q[0]) + (p[1] * q[1])
}

// Vec2CrossProduct gets the cross product of p and q
func Vec2CrossProduct(p [2]float64, q [2]float64) float64 {
	return (p[0] * q[1]) - (p[1] * q[0])
}

// Vec2Add adds q to p
func Vec2Add(p [2]float64, q [2]float64) [2]float64 {
	return [2]float64{p[0] + q[0], p[1] + q[1]}
}

// Vec2Subtract subtracts q from p
func Vec2Subtract(p [2]float64, q [2]float64) [2]float64 {
	return [2]float64{p[0] - q[0], p[1] - q[1]}
}

// Vec2Scale scales p by s
func Vec2Scale(p [2]float64, s float64) [2]float64 {
	return [2]float64{p[0] * s, p[1] * s}
}

// Vec2LenSq gets the squared length of the vector
func Vec2LenSq(p [2]float64) float64 {
	return (p[0] * p[0]) + (p[1] * p[1])
}

// Vec2Len gets the length of the vector
func Vec2Len(p [2]float64) float64 {
	return math.Sqrt(Vec2LenSq(p))
}

// Vec2Unit normalises the vector to a unit vector
func Vec2Unit(p [2]float64) [2]float64 {
	return Vec2Scale(p, 1/Vec2Len(p))
}

// Vec2RotateWithMatrix rotates a vector with a rotation matrix
func Vec2RotateWithMatrix(p [2]float64, r [2][2]float64) [2]float64 {
	return [2]float64{
		p[0]*r[0][0] + p[1]*r[0][1],
		p[1]*r[0][1] + p[1]*r[1][1],
	}
}

// Vec2Normalise normalises a vector
func Vec2Normalise(p [2]float64) [2]float64 {
	m := Vec2Len(p)
	return [2]float64{p[0] / m, p[1] / m}
}

// Vec2UnitToAngle turns a unit angle to a vector
func Vec2UnitToAngle(p [2]float64) float64 {
	return math.Atan2(p[1], p[0])
}
