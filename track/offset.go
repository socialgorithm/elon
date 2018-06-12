package track

import (
	"github.com/paulsmith/gogeos/geos"
	"github.com/socialgorithm/elon-server/domain"
)

// offsets the track to the sides to make a road
func offset(center []domain.Position, offsetDistance float64) domain.Track {
	polygon := newPolygon(center)
	simplified := geos.Must(polygon.Simplify(tolerance))
	bufferOpts := geos.BufferOpts{
		QuadSegs:   8,
		MitreLimit: 5,
		CapStyle:   geos.CapSquare,
		JoinStyle:  geos.JoinBevel,
	}
	innerSideBuf := geos.Must(simplified.BufferWithOpts(-offsetDistance, bufferOpts))
	outerSideBuf := geos.Must(simplified.BufferWithOpts(offsetDistance, bufferOpts))
	innerCoords := getCoords(innerSideBuf)
	outerCoords := getCoords(outerSideBuf)
	return domain.Track{
		Center:    getCoords(simplified),
		InnerSide: innerCoords,
		OuterSide: outerCoords,
		Width:     width,
		Height:    height,
	}
}

// create a polygon from the given positions
func newPolygon(points []domain.Position) *geos.Geometry {
	coords := make([]geos.Coord, len(points), len(points))
	for i := 0; i < len(points); i++ {
		coords[i] = geos.Coord{
			X: float64(points[i].X),
			Y: float64(points[i].Y),
		}
	}
	polygon := geos.Must(geos.NewPolygon(coords))
	return polygon
}

// get the positions from a geometry
func getCoords(g *geos.Geometry) []domain.Position {
	ring := geos.Must(g.Shell())
	cs, _ := ring.Coords()
	return coords2Position(cs)
}

// convert an array of geos.Coord into an array of domain.Position
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
