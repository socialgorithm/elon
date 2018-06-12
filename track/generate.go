package track

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/paulsmith/gogeos/geos"
	"github.com/socialgorithm/elon-server/domain"
)

// GenTrack generates a track within the given width/height
func GenTrack() domain.Track {
	start := time.Now()
	center := addCurves(genInitialConvexTrack())
	track := offset(
		center,
		RoadWidth,
	)
	elapsed := time.Since(start)
	log.Printf("Track Generated in: %s", elapsed)
	return track
}

// Add one point between every 2 points, displaced
func addCurves(points []domain.Position) []domain.Position {
	// we'll add a displacement point between every two points in the track
	numPoints := len(points)*2 - 1
	rPoints := make([]domain.Position, numPoints, numPoints)

	for i, point := range points {
		rPoints[i*2] = point
		if i+1 == len(points) {
			break
		}
		// for each two points, get the middle, then displace it
		vecA := pixel.Vec{
			X: point.X,
			Y: point.Y,
		}
		vecB := pixel.Vec{
			X: points[i+1].X,
			Y: points[i+1].Y,
		}
		middle := vecA.Add(vecB).Scaled(0.5)
		displacement := math.Pow(rand.Float64(), difficulty) * maxDisplacement
		dispVector := pixel.Unit(rand.Float64() * math.Pi).Scaled(displacement)
		midVector := middle.Add(dispVector)
		rPoints[i*2] = point
		rPoints[i*2+1] = domain.Position{
			X: midVector.X,
			Y: midVector.Y,
		}
	}

	return rPoints
}

// Use random points and the convex hull algorithm to get the initial set of points
func genInitialConvexTrack() []domain.Position {
	rand.Seed(time.Now().UnixNano())

	usableWidth := width - margin
	usableHeight := height - margin

	points := rand.Intn(maxPoints-minPoints) + minPoints
	randPoints := make([]*geos.Geometry, points, points)

	for i := 0; i < points; i++ {
		x := float64(0)
		y := float64(0)
		for x == 0 || y == 0 {
			x = rand.Float64()*(usableWidth-margin) + margin
			y = rand.Float64()*(usableHeight-margin) + margin
		}
		randPoints[i], _ = geos.NewPoint(geos.Coord{
			X: x,
			Y: y,
		})
	}

	geom, _ := geos.NewCollection(geos.MULTIPOINT, randPoints...)
	hull := geos.Must(geom.ConvexHull())

	return getCoords(hull)
}
