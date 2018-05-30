package track

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/socialgorithm/elon-server/domain"
)

// GenTrack generates a track within the given width/height
func GenTrack() domain.Track {
	start := time.Now()
	track := offset(
		//addCurvesToTrack(
		genInitialConvexTrack(),
		//),
		roadWidth,
	)
	elapsed := time.Since(start)
	log.Printf("Track Generated in: %s", elapsed)
	return track
}

// Add one point between every 2 points in the track, and displace it away from
// that middle by a certain amount.
// This should create some left/right turns and make it more challenging
func addCurvesToTrack(track domain.Track) domain.Track {
	return domain.Track{
		Center:    addCurves(track.Center),
		InnerSide: track.InnerSide,
		OuterSide: track.OuterSide,
	}
}

// Add one point between every 2 points, displaced
func addCurves(points []domain.Position) []domain.Position {
	rPoints := make([]domain.Position, len(points)*2, len(points)*2+1)

	for i := 0; i < len(points)-1; i++ {
		displacement := math.Pow(rand.Float64(), difficulty) * maxDisplacement
		dispVector := pixel.Unit(rand.Float64() * math.Pi).Scaled(displacement)
		rPoints[i*2] = points[i]
		vectorA := pixel.V(points[i].X, points[i].Y)
		vectorB := pixel.V(points[i+1].X, points[i+1].Y)
		midVector := vectorA.Add(vectorB).Scaled(0.5).Add(dispVector)
		rPoints[i*2+1] = domain.Position{
			X: midVector.X,
			Y: midVector.Y,
		}
	}

	return rPoints
}

// Use random points and the convex hull algorithm to get the initial set of points
func genInitialConvexTrack() domain.Track {
	rand.Seed(time.Now().UnixNano())

	usableWidth := width - margin
	usableHeight := height - margin

	points := rand.Intn(maxPoints-minPoints) + minPoints
	randPoints := make([]domain.Position, points, points)

	for i := 0; i < points; i++ {
		x := float64(0)
		y := float64(0)
		for x == 0 || y == 0 {
			x = rand.Float64()*(usableWidth-margin) + margin
			y = rand.Float64()*(usableHeight-margin) + margin
		}
		randPoints[i] = domain.Position{
			X: x,
			Y: y,
		}
	}

	center := findConvexHull(randPoints)

	return domain.Track{
		Center: center,
	}
}
