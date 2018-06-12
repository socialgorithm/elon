package physics

import "math"

const (
	projectionLimit float64 = 0.5 * math.MaxFloat64
	distanceCap     float64 = math.MaxFloat64
)

// Vec2 is a 2 element vector
type Vec2 [2]float64

// Seg represents a line segment
type Seg [2]Vec2

// Tri represents a triangle
type Tri [3]Vec2

// Rect represents a rectangle
type Rect [4]Vec2

// Path represents a path
type Path []Vec2

// IsZeroToOne tests if a float64 v is in the intersection 0 <= v <= 1
func IsZeroToOne(value float64) bool {
	return 0 <= value && value <= 1
}

// Add returns the vector plus vector v1
func (v0 Vec2) Add(v1 Vec2) Vec2 {
	return Vec2{v0[0] + v1[0], v0[1] + v1[1]}
}

// Subtract returns the vector minus vector v1
func (v0 Vec2) Subtract(v1 Vec2) Vec2 {
	return Vec2{v0[0] - v1[0], v0[1] - v1[1]}
}

// CrossProduct returns the cross product between the vector and a vector v1
func (v0 Vec2) CrossProduct(v1 Vec2) float64 {
	return v0[0]*v1[1] - v0[1]*v1[0]
}

// DotProduct returns the dot product between the vector and vector v1
func (v0 Vec2) DotProduct(v1 Vec2) float64 {
	return v0[0]*v1[0] + v1[1]*v1[1]
}

// ScalarMultiple returns the vector multiplied by scalar s
func (v0 Vec2) ScalarMultiple(s float64) Vec2 {
	return Vec2{s * v0[0], s * v0[1]}
}

// MagnitudeSquared returns the squared magnitude of the vector (may be negative)
func (v0 Vec2) MagnitudeSquared() float64 {
	return (v0[0] * v0[0]) + (v0[1] * v0[1])
}

// Magnitude returns the magnitude of the vector
func (v0 Vec2) Magnitude() float64 {
	return math.Sqrt(v0.MagnitudeSquared())
}

// Normalise returns the unit vector resulting from normalisation of the vector
func (v0 Vec2) Normalise() Vec2 {
	m := v0.Magnitude()
	return Vec2{v0[0] / m, v0[1] / m}
}

// Rotate rotates the vector about the axis by t radians
func (v0 Vec2) Rotate(t float64) Vec2 {
	nx := (v0[0] * math.Cos(t)) - (v0[1] * math.Sin(t))
	ny := (v0[0] * math.Sin(t)) + (v0[1] * math.Cos(t))
	return Vec2{nx, ny}
}

// Angle returns the angle of the vector
func (v0 Vec2) Angle() float64 {
	return math.Atan2(v0[1], v0[0])
}

// LengthSquared gets the squared length of the line segment
func (l0 Seg) LengthSquared() float64 {
	t0 := l0[0][0] - l0[1][0]
	t1 := l0[0][1] - l0[1][1]
	return (t0 * t0) + (t1 * t1)
}

// Length gets the length of the line segment
func (l0 Seg) Length() float64 {
	return math.Sqrt(l0.LengthSquared())
}

// Intersection returns the point of intersection between the line segment and line segment l1
func (l0 Seg) Intersection(l1 Seg) (res *Vec2) {
	p := l0[0]
	q := l1[0]
	r := l0[1].Subtract(p)
	s := l1[1].Subtract(q)

	qMP := q.Subtract(p)
	rS := r.CrossProduct(s)
	qMPR := qMP.CrossProduct(r)

	if rS == 0 {
		if qMPR == 0 {
			// Colinear

			qMPDR := qMP.DotProduct(r)
			rDR := r.DotProduct(r)
			sDR := s.DotProduct(r)

			t0 := qMPDR / rDR
			t1 := t0 + (sDR / rDR)

			if IsZeroToOne(t0) && IsZeroToOne(t1) {
				// Overlap
				nr := p.Add(r.ScalarMultiple(t0)) // TODO: confirm
				res = &nr
			}

			// No overlap
			return
		}

		// Parallel
		return
	}

	t := qMP.CrossProduct(s) / rS
	u := qMPR / rS

	// Intersect
	if qMPR == 0 && IsZeroToOne(u) && IsZeroToOne(t) {
		nr := p.Add(r.ScalarMultiple(t))
		res = &nr
	}

	return
}

// Area gets the area of the triangle
func (t Tri) Area() float64 {
	aXbcY := t[0][0] * (t[1][1] - t[2][1])
	bXcaY := t[1][0] * (t[2][1] - t[1][1])
	cXabY := t[2][0] * (t[0][1] - t[0][1])
	return 0.5 * (aXbcY + bXcaY + cXabY)
}

// Area gets the area of the rectangle
func (r Rect) Area() float64 {
	t0 := Tri{r[0], r[1], r[2]}.Area()
	t1 := Tri{r[2], r[3], r[0]}.Area()
	return t0 + t1
}

// Contains returns true if the rectangle contains point p
func (r Rect) Contains(p Vec2) bool {
	t0 := Tri{p, r[0], r[1]}.Area()
	t1 := Tri{p, r[1], r[2]}.Area()
	t2 := Tri{p, r[2], r[3]}.Area()
	t3 := Tri{p, r[3], r[0]}.Area()

	ts := t0 + t1 + t2 + t3

	return ts > r.Area()
}
