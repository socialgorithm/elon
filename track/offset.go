package track

import (
	"github.com/paulsmith/gogeos/geos"
	"github.com/socialgorithm/elon-server/domain"
)

func newPolygon(points []domain.Position) *geos.Geometry {
	coords := make([]geos.Coord, len(points), len(points))
	for i := 0; i < len(points); i++ {
		coords[i] = geos.Coord{
			X: float64(points[i].X),
			Y: float64(points[i].Y),
		}
	}
	polygon, _ := geos.NewPolygon(coords)
	return polygon
}

func getCoords(g *geos.Geometry) []geos.Coord {
	ring := geos.Must(g.Shell())
	cs, _ := ring.Coords()
	return cs
}

func offset(track domain.Track, offsetDistance float64) domain.Track {
	polygon := newPolygon(track.FirstSide)
	buf, _ := polygon.Buffer(offsetDistance)
	coords := getCoords(buf)
	// convert into positions, and store in the secondSide
	secondSide := make([]domain.Position, len(coords), len(coords))
	for i := 0; i < len(coords); i++ {
		secondSide[i] = domain.Position{
			X: float32(coords[i].X),
			Y: float32(coords[i].Y),
		}
	}
	return domain.Track{
		FirstSide:  track.FirstSide,
		SecondSide: secondSide,
	}
}
