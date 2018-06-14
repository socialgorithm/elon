package physics

import (
	"math"
)

const (
	// Pi is pi
	Pi float64 = math.Pi

	// TwoPi is 2 * pi
	TwoPi float64 = 2 * Pi

	// SensorRange is the range of sensors
	SensorRange = 5

	// XDim is the x dimension of cars
	XDim = 4

	// YDim is the y dimension of cars
	YDim = 6

	// SteeringRate is the modifier for steering
	SteeringRate = 0.05

	// AccelerationRate is the modifier for acceleration
	AccelerationRate = 0.05

	// MaxVelocity is the limit on velocity
	MaxVelocity = 5
)

var (
	sinHalfPi = math.Sin(0.5 * Pi)
)