package physics

import (
	"math"

	"github.com/socialgorithm/elon-server/domain"
)

// NormaliseRadians normalises radians between 0 and 2Pi
func NormaliseRadians(t float64) float64 {
	return t - (TwoPi * math.Floor((t+Pi)/TwoPi))
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
