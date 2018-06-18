package physics

import (
	"math"
)

const (
	xHDim = 0.5 * XDim
	yHDim = 0.5 * YDim
)

var (
	sqrtMaxFloat64 = math.Sqrt(math.MaxFloat64)
)

var (
	sensorAngledRange = sinHalfPi * SensorRange
	minHDim           = math.Min(xHDim, yHDim)
	dims              = [2]float64{XDim, YDim}
	dimsRect          = [2][2]float64{
		{-xHDim, -yHDim},
		{xHDim, yHDim},
	}
	dimsFSeg = [2][2]float64{
		{-xHDim, yHDim},
		{xHDim, yHDim},
	}
	dimsBSeg = [2][2]float64{
		{-xHDim, -yHDim},
		{xHDim, -yHDim},
	}
	sensors = [5][2][2]float64{
		{ // Left
			{-xHDim, 0},
			{-(xHDim + SensorRange), 0},
		},
		{ // Left-Centre
			{-minHDim, minHDim},
			{-(minHDim + sensorAngledRange), minHDim + sensorAngledRange},
		},
		{ // Centre
			{0, yHDim},
			{0, yHDim + SensorRange},
		},
		{ // Right-Centre
			{minHDim, minHDim},
			{minHDim + sensorAngledRange, minHDim + sensorAngledRange},
		},
		{ // Right
			{xHDim, 0},
			{xHDim + SensorRange, 0},
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

func (s State) forwardRect() (rect, rectLSeg, rectRSeg [2][2]float64) {
	rect = [2][2]float64{
		{dimsRect[0][0], dimsRect[1][1]},
		{dimsRect[1][0], dimsRect[1][1] + s.Velocity},
	}
	rectLSeg = [2][2]float64{
		{dimsRect[0][0], dimsRect[1][1]},
		{dimsRect[0][0], dimsRect[1][1] + s.Velocity},
	}
	rectRSeg = [2][2]float64{
		{dimsRect[1][0], dimsRect[1][1]},
		{dimsRect[1][0], dimsRect[1][1] + s.Velocity},
	}
	return
}

func (s State) backwardRect() (rect, rectLSeg, rectRSeg [2][2]float64) {
	rect = [2][2]float64{
		{dimsRect[0][0], dimsRect[0][1] + s.Velocity},
		{dimsRect[1][0], dimsRect[0][1]},
	}
	rectLSeg = [2][2]float64{
		{dimsRect[0][0], dimsRect[0][1]},
		{dimsRect[0][0], dimsRect[0][1] + s.Velocity},
	}
	rectRSeg = [2][2]float64{
		{dimsRect[1][0], dimsRect[0][1]},
		{dimsRect[1][0], dimsRect[0][1] + s.Velocity},
	}
	return
}

// Check checks whether a car has collided and its sensor values
func (s State) Check(p [][2][2]float64) (bool, float64, [5]float64) {
	sin := math.Sin(-s.Angle - (math.Pi * 1.5))
	cos := math.Cos(-s.Angle - (math.Pi * 1.5))

	var rect, rectLSeg, rectRSeg, tSeg [2][2]float64
	if s.Velocity < 0 {
		rect, rectLSeg, rectRSeg = s.backwardRect()
		tSeg = dimsBSeg
	} else {
		rect, rectLSeg, rectRSeg = s.forwardRect()
		tSeg = dimsFSeg
	}

	mn := math.MaxFloat64
	sns := [5]float64{}
	for idx := range sns {
		sns[idx] = math.MaxFloat64
	}

	for _, seg := range p {
		act := [2][2]float64{
			Vec2RotateWithSinAndCos(Vec2Subtract(seg[0], s.Position), sin, cos),
			Vec2RotateWithSinAndCos(Vec2Subtract(seg[1], s.Position), sin, cos),
		}

		// Collision detection
		if RectIntersects(rect, act) {
			cn0 := RectContains(rect, act[0])
			cn1 := RectContains(rect, act[1])

			if cn0 {
				mn = math.Min(mn, SegDistanceToSq(tSeg, act[0]))
			}

			if cn1 {
				mn = math.Min(mn, SegDistanceToSq(tSeg, act[1]))
			}

			if !cn0 || !cn1 {
				for _, in := range SegIntersections(rectLSeg, act) {
					mn = math.Min(mn, SegDistanceToSq(tSeg, in))
				}
				for _, in := range SegIntersections(rectRSeg, act) {
					mn = math.Min(mn, SegDistanceToSq(tSeg, in))
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

	return mn < sqrtMaxFloat64, mn, sns
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
	s.Angle = NormaliseRadians(s.Angle - (s.Velocity * velocityFactoryRange * SteeringRate * s.Steering))

	// Velocity
	s.Velocity = RoundVelocity(CapValue((frictionMultiplier*s.Velocity)+(AccelerationRate*s.Throttle), -MaxReverseVelocity, MaxVelocity))

	res = sensorValues
	return
}
