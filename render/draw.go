package render

import (
	"image/color"
	"math"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	imdraw "github.com/faiface/pixel/imdraw"
	"github.com/socialgorithm/elon-server/domain"
)

const (
	halfPi float64 = 0.5 * math.Pi
)

func renderTrack(track domain.Track) *imdraw.IMDraw {
	trackRender := imdraw.New(nil)

	// Rotation matrix for the car rendering
	//rotation := pixel.IM.Rotated(pixel.Vec{X: track.Center[0].X, Y: track.Center[0].Y}, halfPi)
	//trackRender.SetMatrix(rotation)

	// Draw road (not working for concave polygons)
	//drawPolygon(trackRender, track.OuterSide, roadColor)
	//drawPolygon(trackRender, track.InnerSide, bgColor)

	// Draw track using triangulation
	//drawTrack(trackRender, track)
	// Draw outer lines
	drawLine(trackRender, track.OuterSide, outerLinesColor, lineThickness)
	drawLine(trackRender, track.InnerSide, outerLinesColor, lineThickness)
	// Draw line segments
	drawDashedLine(trackRender, track.OuterSide, segmentColor, lineThickness)
	drawDashedLine(trackRender, track.InnerSide, segmentColor, lineThickness)
	drawDashedLine(trackRender, track.Center, centerLineColor, 1)
	// Draw start
	drawStart(trackRender, track)

	return trackRender
}

// draw a line representing the start of the track
func drawStart(draw *imdraw.IMDraw, track domain.Track) {
	draw.Color = startLineColor
	// the line will be drawn in an angle disecting the start & end segments
	startVec := pixel.V(track.Center[0].X, track.Center[0].Y)
	first := pixel.V(track.Center[1].X, track.Center[1].Y)
	last := pixel.V(track.Center[len(track.Center)-2].X, track.Center[len(track.Center)-2].Y)
	// now get the angle between the first point and the last point
	vecA := first.Sub(startVec).Unit()
	vecB := last.Sub(startVec).Unit()
	angle := math.Acos(vecA.Dot(vecB))
	lineVec := vecB.Rotated(angle / 2).Scaled(25)
	draw.Push(
		startVec.Add(lineVec),
		startVec.Sub(lineVec),
	)
	draw.Line(5)
}

func renderCar(carState domain.CarState) *imdraw.IMDraw {
	carRender := imdraw.New(nil)

	// Prepare some vectors
	posVector := pixel.V(carState.Position.X, carState.Position.Y)
	dirVector := pixel.V(carState.Direction.X, carState.Direction.Y)
	sensorUnitVector := pixel.Unit(-math.Pi / 2)

	// Rotation matrix for the car rendering
	rotation := pixel.IM.Rotated(posVector, dirVector.Angle()+halfPi)
	carRender.SetMatrix(rotation)

	// render car fill
	if carState.Crashed != true {
		carRender.Color = carColor
	} else {
		carRender.Color = carCrashedColor
	}

	carRender.Push(
		pixel.V(carState.Position.X-carWidth*.5, carState.Position.Y-carLength*.5),
		pixel.V(carState.Position.X+carWidth*.5, carState.Position.Y+carLength*.5),
	)
	carRender.Rectangle(0)

	// render car outline
	carRender.Color = colornames.Black
	carRender.Push(
		pixel.V(carState.Position.X-carWidth*.5, carState.Position.Y-carLength*.5),
		pixel.V(carState.Position.X+carWidth*.5, carState.Position.Y+carLength*.5),
	)
	carRender.Rectangle(1)

	// render car middle point
	carRender.Color = colornames.Yellow
	carRender.Push(
		pixel.V(carState.Position.X, carState.Position.Y),
	)
	carRender.Circle(2, 0)

	// render sensors
	carRender.Color = colornames.Orange
	for i := 0; i < len(carState.Sensors); i++ {
		sensor := carState.Sensors[i]
		sensorVector := sensorUnitVector.Scaled(sensor.Distance).Rotated(sensor.Angle + math.Pi)
		carRender.Push(
			pixel.V(carState.Position.X, carState.Position.Y),
			pixel.V(carState.Position.X+sensorVector.X, carState.Position.Y+sensorVector.Y),
		)
		carRender.Line(1)
	}

	// render wheels
	carRender.Color = colornames.Black
	// top left
	carRender.Push(
		pixel.V(carState.Position.X-carWidth*.5, carState.Position.Y-carLength*.5+wheelOffset),
		pixel.V(carState.Position.X-carWidth*.5, carState.Position.Y-carLength*.5+wheelOffset+wheelLength),
	)
	carRender.Line(wheelWidth)
	// top right
	carRender.Push(
		pixel.V(carState.Position.X+carWidth*.5, carState.Position.Y-carLength*.5+wheelOffset),
		pixel.V(carState.Position.X+carWidth*.5, carState.Position.Y-carLength*.5+wheelOffset+wheelLength),
	)
	carRender.Line(wheelWidth)
	// bottom left
	carRender.Push(
		pixel.V(carState.Position.X-carWidth*.5, carState.Position.Y+carLength*.5-wheelOffset),
		pixel.V(carState.Position.X-carWidth*.5, carState.Position.Y+carLength*.5-wheelOffset-wheelLength),
	)
	carRender.Line(wheelWidth)
	// bottom right
	carRender.Push(
		pixel.V(carState.Position.X+carWidth*.5, carState.Position.Y+carLength*.5-wheelOffset),
		pixel.V(carState.Position.X+carWidth*.5, carState.Position.Y+carLength*.5-wheelOffset-wheelLength),
	)
	carRender.Line(wheelWidth)
	return carRender
}

func drawLine(draw *imdraw.IMDraw, points []domain.Position, color color.RGBA, thickness float64) {
	for i := 0; i < len(points)-1; i++ {
		pointA := points[i]
		pointB := points[i+1]

		draw.Color = color
		draw.Push(
			pixel.V(pointA.X, pointA.Y),
			pixel.V(pointB.X, pointB.Y),
		)
	}
	draw.Line(thickness)
}

func drawDashedLine(draw *imdraw.IMDraw, points []domain.Position, color color.RGBA, thickness float64) {
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
}

func drawPolygon(draw *imdraw.IMDraw, points []domain.Position, color color.RGBA) {
	draw.Color = color
	for i := 0; i < len(points); i++ {
		point := points[i]
		draw.Push(
			pixel.V(point.X, point.Y),
		)
	}
	draw.Polygon(0) // filled
}
