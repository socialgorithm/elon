package render

import (
	"github.com/aurbano/triangulate"
	"github.com/faiface/pixel"
	imdraw "github.com/faiface/pixel/imdraw"
	"github.com/socialgorithm/elon-server/domain"
)

// triangulate convert a given track into a list of triangles that can be rendered
// as primitives on the screen (used to fill in the track)
// TODO Figure out why its not working...
func drawTrack(draw *imdraw.IMDraw, track domain.Track) {
	polygon := triangulate.Polygon{
		Exterior:  positions2Ring(track.OuterSide),
		Interiors: []triangulate.Ring{positions2Ring(track.InnerSide)},
	}
	triangulated := polygon.Triangulate()
	draw.Color = roadColor
	for i := 0; i < len(triangulated); i++ {
		triangle := triangulated[i]
		draw.Push(
			pixel.V(triangle.A.X, triangle.A.Y),
			pixel.V(triangle.B.X, triangle.B.Y),
			pixel.V(triangle.C.X, triangle.C.Y),
		)
		draw.Polygon(3)
	}
}

func positions2Ring(positions []domain.Position) triangulate.Ring {
	ring := make(triangulate.Ring, len(positions), len(positions))
	for i, pos := range positions {
		ring[i] = triangulate.Point{
			X: pos.X,
			Y: pos.Y,
		}
	}
	return ring
}
