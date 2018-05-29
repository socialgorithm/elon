package track

import (
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/socialgorithm/elon-server/domain"
)

const points = 50
const difficulty = float64(0.3) // closer to 0 will create sharper turns, exponentially
const maxDisplacement = float64(100)
const margin = float32(100)
const roadWidth = float32(30)

func genTrack(width int32, height int32) domain.Track {
	return getSecondSideAsOffset(
		addCurvesToTrack(
			genInitialConvexTrack(width, height),
		),
	)
}

// Add one point between every 2 points in the track, and displace it away from
// that middle by a certain amount.
// This should create some left/right turns and make it more challenging
func addCurvesToTrack(track domain.Track) domain.Track {
	firstSide := addCurves(track.FirstSide)
	secondSide := addCurves(track.SecondSide)

	return domain.Track{
		FirstSide:    firstSide,
		SecondSide:   secondSide,
		RandomPoints: track.RandomPoints,
	}
}

// Offset the first side of the track by a set amount to get the second
// side of the track
func getSecondSideAsOffset(track domain.Track) domain.Track {
	firstSide := track.FirstSide
	secondSide := make([]domain.Position, len(firstSide), len(firstSide))

	for i := 0; i < len(firstSide); i++ {
		secondSide[i] = domain.Position{
			X: firstSide[i].X + roadWidth,
			Y: firstSide[i].Y + roadWidth,
		}
	}

	return domain.Track{
		FirstSide:    firstSide,
		SecondSide:   secondSide,
		RandomPoints: track.RandomPoints,
	}
}

// Add one point between every 2 points, displaced
func addCurves(points []domain.Position) []domain.Position {
	rPoints := make([]domain.Position, len(points)*2, len(points)*2+1)

	for i := 0; i < len(points)-1; i++ {
		displacement := math.Pow(rand.Float64(), difficulty) * maxDisplacement
		dispVector := pixel.Unit(rand.Float64() * math.Pi).Scaled(displacement)
		rPoints[i*2] = points[i]
		vectorA := pixel.V(float64(points[i].X), float64(points[i].Y))
		vectorB := pixel.V(float64(points[i+1].X), float64(points[i+1].Y))
		midVector := vectorA.Add(vectorB).Scaled(0.5).Add(dispVector)
		rPoints[i*2+1] = domain.Position{
			X: float32(midVector.X),
			Y: float32(midVector.Y),
		}
	}

	return rPoints
}

// Use random points and the convex hull algorithm to get the initial set of points
func genInitialConvexTrack(_width int32, _height int32) domain.Track {
	rand.Seed(time.Now().UnixNano())

	width := float32(_width)
	height := float32(_height)

	var randPoints = [points]domain.Position{}

	for i := 0; i < points; i++ {
		x := float32(0)
		y := float32(0)
		for x == 0 || y == 0 {
			x = rand.Float32()*(width-margin) + margin
			y = rand.Float32()*(height-margin) + margin
		}
		randPoints[i] = domain.Position{
			X: x,
			Y: y,
		}
	}

	firstSide := findConvexHull(randPoints[0:len(randPoints)])

	return domain.Track{
		FirstSide:    firstSide,
		SecondSide:   firstSide,
		RandomPoints: randPoints[0 : len(randPoints)-1],
	}
}
