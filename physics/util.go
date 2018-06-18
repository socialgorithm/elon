package physics

import (
	"math"

	"github.com/socialgorithm/elon-server/domain"
)

const (
	onePercentOfVelocity        = MaxVelocity * 0.005
	oneOverOnePercentOfVelocity = 1 / onePercentOfVelocity
)

// NormaliseRadians normalises radians between 0 and 2Pi
func NormaliseRadians(t float64) float64 {
	return t - (math.Floor(t/TwoPi) * TwoPi)
}

// CapValue caps v between min and max
func CapValue(v, min, max float64) float64 {
	return math.Max(min, math.Min(v, max))
}

// DPosToVec2 converts a domain position to a vector
func DPosToVec2(p domain.Position) [2]float64 {
	return [2]float64{p.X, p.Y}
}

// AngleToUnitVec2 converts an angle to a unit vector
func AngleToUnitVec2(t float64) [2]float64 {
	return [2]float64{math.Cos(t), math.Sin(t)}
}

// RoundVelocity rounds to the nearest 1% of max velocity
func RoundVelocity(v float64) float64 {
	return onePercentOfVelocity * math.Floor(v*oneOverOnePercentOfVelocity)
}
