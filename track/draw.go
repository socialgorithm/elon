package track

import (
	"image/color"

	"github.com/faiface/pixel"
	imdraw "github.com/faiface/pixel/imdraw"
	"github.com/socialgorithm/elon-server/domain"
	"golang.org/x/image/colornames"
)

func drawTrack(track domain.Track) *imdraw.IMDraw {
	trackRender := imdraw.New(nil)
	// draw sides
	drawPolygon(trackRender, track.SecondSide, colornames.Gray)
	drawPolygon(trackRender, track.FirstSide, colornames.Green)
	drawLine(trackRender, track.FirstSide, colornames.White)
	drawLine(trackRender, track.SecondSide, colornames.White)

	return trackRender
}

func drawLine(draw *imdraw.IMDraw, points []domain.Position, color color.RGBA) *imdraw.IMDraw {
	for i := 0; i < len(points)-1; i++ {
		pointA := points[i]
		pointB := points[i+1]

		draw.Color = color
		draw.Push(
			pixel.V(float64(pointA.X), float64(pointA.Y)),
			pixel.V(float64(pointB.X), float64(pointB.Y)),
		)
		draw.Line(lineThickness)
	}
	return draw
}

func drawPolygon(draw *imdraw.IMDraw, points []domain.Position, color color.RGBA) *imdraw.IMDraw {
	draw.Color = color
	for i := 0; i < len(points); i++ {
		point := points[i]
		draw.Push(
			pixel.V(float64(point.X), float64(point.Y)),
		)
	}
	draw.Polygon(0) // filled
	return draw
}
