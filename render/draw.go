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

func renderCar(carState domain.CarState) *imdraw.IMDraw {
	carRender := imdraw.New(nil)

	posVector := pixel.V(carState.Position.X, carState.Position.Y)
	dirVector := pixel.V(carState.Direction.X, carState.Direction.Y)
	rotation := pixel.Matrix.Rotated(pixel.IM, posVector, dirVector.Angle())

	carRender.SetMatrix(rotation)

	// render car fill
	carRender.Color = carColor
	carRender.Push(
		pixel.V(carState.Position.X-carWidth/2, carState.Position.Y),
		pixel.V(carState.Position.X+carWidth/2, carState.Position.Y+carLength),
	)
	carRender.Rectangle(0)

	// render car outline
	carRender.Color = colornames.Black
	carRender.Push(
		pixel.V(carState.Position.X-carWidth/2, carState.Position.Y),
		pixel.V(carState.Position.X+carWidth/2, carState.Position.Y+carLength),
	)
	carRender.Rectangle(1)

	// render sensors
	// carRender.Color = colornames.Orange
	// for i := 0; i < len(carState.Sensors); i++ {
	// 	sensor := carState.Sensors[i]
	// 	sensorVector := pixel.V(
	// 		carState.Position.X-carWidth/2,
	// 		carState.Position.Y,
	// 	).
	// 	carRender.Push(
	// 		pixel.V(carState.Position.X-carWidth/2, carState.Position.Y),
	// 		pixel.V(sensorVector.X, sensorVector.Y),
	// 	)
	// 	carRender.Line(1)
	// }

	// render wheels
	carRender.Color = colornames.Black
	// top left
	carRender.Push(
		pixel.V(carState.Position.X-carWidth/2, carState.Position.Y+wheelOffset),
		pixel.V(carState.Position.X-carWidth/2, carState.Position.Y+wheelOffset+wheelLength),
	)
	carRender.Line(wheelWidth)
	// top right
	carRender.Push(
		pixel.V(carState.Position.X+carWidth/2, carState.Position.Y+wheelOffset),
		pixel.V(carState.Position.X+carWidth/2, carState.Position.Y+wheelOffset+wheelLength),
	)
	carRender.Line(wheelWidth)
	// bottom left
	carRender.Push(
		pixel.V(carState.Position.X-carWidth/2, carState.Position.Y+carLength-wheelOffset),
		pixel.V(carState.Position.X-carWidth/2, carState.Position.Y+carLength-wheelOffset-wheelLength),
	)
	carRender.Line(wheelWidth)
	// bottom right
	carRender.Push(
		pixel.V(carState.Position.X+carWidth/2, carState.Position.Y+carLength-wheelOffset),
		pixel.V(carState.Position.X+carWidth/2, carState.Position.Y+carLength-wheelOffset-wheelLength),
	)
	carRender.Line(wheelWidth)

	return carRender
}

func drawLine(draw *imdraw.IMDraw, points []domain.Position, color color.RGBA, thickness float64) *imdraw.IMDraw {
	for i := 0; i < len(points)-1; i++ {
		pointA := points[i]
		pointB := points[i+1]

		draw.Color = color
		draw.Push(
			pixel.V(pointA.X, pointA.Y),
			pixel.V(pointB.X, pointB.Y),
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
