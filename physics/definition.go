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
	SensorRange = 1000

	// XDim is the x dimension of cars
	XDim = 19

	// YDim is the y dimension of cars
	YDim = 34

	// SteeringRate is the modifier for steering
	SteeringRate = 0.04

	// AccelerationRate is the modifier for acceleration
	AccelerationRate = 0.03

	// MaxVelocity is the limit on velocity
	MaxVelocity = 1.5

	// MaxReverseVelocity is the limit on velocity in reverse
	MaxReverseVelocity = 0.5

	// Friction is the proportion of velocity lost per tick (throttle counteracts)
	Friction = 0.005

	frictionMultiplier = 1 - Friction
)

var (
	sinHalfPi = math.Sin(0.5 * Pi)

	velocityFactoryRange = (1 / math.Max(MaxVelocity, MaxReverseVelocity))
)
