package track

import (
	"github.com/paulsmith/gogeos/geos"
	"github.com/socialgorithm/elon-server/domain"
)

// offsets the track to the sides to make a road
func offset(track domain.Track, offsetDistance float64) domain.Track {
	polygon := newPolygon(track.Center)
	innerSideBuf, _ := polygon.Buffer(-offsetDistance)
	outerSideBuf, _ := polygon.Buffer(offsetDistance)
	innerCoords := coords2Position(getCoords(innerSideBuf))
	outerCoords := coords2Position(getCoords(outerSideBuf))
	return domain.Track{
		Center:    track.Center,
		InnerSide: innerCoords,
		OuterSide: outerCoords,
	}
}

func coords2Position(coords []geos.Coord) []domain.Position {
	positions := make([]domain.Position, len(coords), len(coords))
	for i := 0; i < len(coords); i++ {
		positions[i] = domain.Position{
			X: coords[i].X,
			Y: coords[i].Y,
		}
	}
	return positions
}

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
