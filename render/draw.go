package render

import (
	"image/color"
	"math"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	imdraw "github.com/faiface/pixel/imdraw"
	"github.com/socialgorithm/elon-server/domain"
)

func renderTrack(track domain.Track) *imdraw.IMDraw {
	trackRender := imdraw.New(nil)

	// Draw road
	drawPolygon(trackRender, track.OuterSide, roadColor)
	drawPolygon(trackRender, track.InnerSide, bgColor)
	// Draw outer lines
	drawLine(trackRender, track.OuterSide, outerLinesColor, lineThickness)
	drawLine(trackRender, track.InnerSide, outerLinesColor, lineThickness)
	// Draw line segments
	drawDashedLine(trackRender, track.OuterSide, segmentColor, lineThickness)
	drawDashedLine(trackRender, track.InnerSide, segmentColor, lineThickness)
	drawDashedLine(trackRender, track.Center, centerLineColor, 1)

	return trackRender
}

func renderCar(car domain.Car) *imdraw.IMDraw {
	carRender := imdraw.New(nil)

	carRender.Color = colornames.Peru
	carRender.Push(
		pixel.V(car.CarState.Position.X, car.CarState.Position.Y),
	)
	carRender.Circle(10, 0)

	return carRender
}

func drawLine(draw *imdraw.IMDraw, points []domain.Position, color color.RGBA, thickness float64) *imdraw.IMDraw {
	for i := 0; i < len(points)-1; i++ {
		pointA := points[i]
		pointB := points[i+1]

		draw.Color = color
		draw.Push(
			pixel.V(float64(pointA.X), float64(pointA.Y)),
			pixel.V(float64(pointB.X), float64(pointB.Y)),
		)
		draw.Line(thickness)
	}
	return draw
}

func drawDashedLine(draw *imdraw.IMDraw, points []domain.Position, color color.RGBA, thickness float64) *imdraw.IMDraw {
	draw.Color = color
	for i := 0; i < len(points)-1; i++ {
		pointA := points[i]
		pointB := points[i+1]

		// convert each two points into a vector, then use that to calculate the segment points
		vectorA := pixel.V(pointA.X, pointA.Y)
		vectorB := pixel.V(pointB.X, pointB.Y)
		vector := vectorB.Sub(vectorA)
		distance := vector.Len()
		uVector := vector.Unit()

		if distance < segmentLength {
			continue
		}
		segmentStart := vectorA
		segments := int(math.Floor(distance / segmentLength))
		for a := 0; a < segments; a = a + 2 {
			segmentEnd := segmentStart.Add(uVector.Scaled(segmentLength))
			draw.Push(
				pixel.V(segmentStart.X, segmentStart.Y),
				pixel.V(segmentEnd.X, segmentEnd.Y),
			)
			draw.Line(thickness)
			segmentStart = segmentEnd.Add(uVector.Scaled(segmentLength))
		}
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
