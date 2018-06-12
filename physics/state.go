package physics

import (
	"math"
)

const (
	sensorRange      = 5
	xDim             = 4
	yDim             = 6
	xHDim            = 0.5 * xDim
	yHDim            = 0.5 * yDim
	steeringRate     = 0.05
	accelerationRate = 0.05
	maxSpeed         = 1
)

var (
	sensorAngledRange = sinHalfPi * sensorRange
	minHDim           = math.Min(xHDim, yHDim)
	dims              = [2]float64{xDim, yDim}
	dimsRect          = [2][2]float64{
		{-xHDim, -yHDim},
		{xHDim, yHDim},
	}
	dimsFSeg = [2][2]float64{
		{-xHDim, yHDim},
		{xHDim, yHDim},
	}
	sensors = [5][2][2]float64{
		{ // Left
			{-xHDim, 0},
			{-(xHDim + sensorRange), 0},
		},
		{ // Left-Centre
			{-minHDim, minHDim},
			{-(minHDim + sensorAngledRange), minHDim + sensorAngledRange},
		},
		{ // Centre
			{0, yHDim},
			{0, yHDim + sensorRange},
		},
		{ // Right-Centre
			{minHDim, minHDim},
			{minHDim + sensorAngledRange, minHDim + sensorAngledRange},
		},
		{ // Right
			{xHDim, 0},
			{xHDim + sensorRange, 0},
		},
	}
)

// State keeps the state of a car
type State struct {
	Position [2]float64
	Angle    float64
	Velocity float64
	Throttle float64
	Steering float64
	Crashed  bool
}

// Check checks whether a car has collided and its sensor values
func (s State) Check(p [][2][2]float64) (bool, float64, [5]float64) {
	sin := math.Sin(s.Angle)
	cos := math.Cos(s.Angle)
	rot := [2][2]float64{
		{cos, -sin},
		{sin, cos},
	}
	rect := [2][2]float64{
		{dimsRect[0][0], dimsRect[1][1]},
		{dimsRect[1][0], dimsRect[1][1] + s.Velocity},
	}
	rectLSeg := [2][2]float64{
		{dimsRect[0][0], dimsRect[1][1]},
		{dimsRect[0][0], dimsRect[1][1] + s.Velocity},
	}
	rectRSeg := [2][2]float64{
		{dimsRect[1][0], dimsRect[1][1]},
		{dimsRect[1][0], dimsRect[1][1] + s.Velocity},
	}

	mn := math.MaxFloat64
	sns := [5]float64{}
	for idx := range sns {
		sns[idx] = math.MaxFloat64
	}

	for _, seg := range p {
		// Effectively rotates the world around the car (leads to cheaper collision checks)
		act := [2][2]float64{
			Vec2Subtract(Vec2RotateWithMatrix(seg[0], rot), s.Position),
			Vec2Subtract(Vec2RotateWithMatrix(seg[1], rot), s.Position),
		}

		// Collision detection
		if RectIntersects(rect, act) {
			cn0 := RectContains(rect, act[0])
			cn1 := RectContains(rect, act[1])

			if cn0 {
				mn = math.Min(mn, SegDistanceToSq(dimsFSeg, act[0]))
			}

			if cn1 {
				mn = math.Min(mn, SegDistanceToSq(dimsFSeg, act[1]))
			}

			if !cn0 || !cn1 {
				for _, in := range SegIntersections(rectLSeg, act) {
					mn = math.Min(mn, SegDistanceToSq(dimsFSeg, in))
				}
				for _, in := range SegIntersections(rectRSeg, act) {
					mn = math.Min(mn, SegDistanceToSq(dimsFSeg, in))
				}
			}
		}

		// Sensor calculations
		for idx, sn := range sensors {
			for _, in := range SegIntersections(sn, act) {
				sns[idx] = math.Min(sns[idx], Vec2LenSq(Vec2Subtract(in, sn[0])))
			}
		}
	}

	mn = math.Sqrt(mn)
	for idx, sv := range sns {
		sns[idx] = math.Sqrt(sv)
	}

	return mn < math.MaxFloat64, mn, sns
}

// Update updates a state item for a single tick
func (s *State) Update(p [][2][2]float64) (res [5]float64) {
	if s.Crashed {
		return
	}

	crashed, distanceToCrash, sensorValues := s.Check(p)

	updateBy := s.Velocity
	if crashed {
		updateBy = distanceToCrash
	}

	angleVector := AngleToUnitVec2(s.Angle)
	updateVector := Vec2Scale(angleVector, updateBy)

	s.Position[0] += updateVector[0]
	s.Position[1] += updateVector[1]
	if crashed {
		s.Crashed = true
		return
	}

	// Angle
	s.Angle = NormaliseRadians(s.Angle + (steeringRate * s.Steering))

	// Velocity
	s.Velocity = CapValue(s.Velocity+(accelerationRate*s.Throttle), 0, maxSpeed)

	res = sensorValues
	return
}
