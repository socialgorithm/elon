package physics

import (
	"math"
)

const (
	// Pi is pi
	Pi float64 = math.Pi

	// TwoPi is 2 * pi
	TwoPi float64 = 2 * Pi

	// Core definitions
	sensorRange      = 5
	xDim             = 4
	yDim             = 6
	steeringRate     = 0.05
	accelerationRate = 0.05
	maxSpeed         = 1
)

var (
	sinHalfPi = math.Sin(0.5 * Pi)
)
