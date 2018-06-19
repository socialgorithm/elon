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

func renderCar(car domain.Car) *imdraw.IMDraw {
	carRender := imdraw.New(nil)

	// Prepare some vectors
	carCenter := pixel.V(car.CarState.Position.X-carWidth*.5, car.CarState.Position.Y-carLength*.5)
	posVector := pixel.V(car.CarState.Position.X, car.CarState.Position.Y)
	dirVector := pixel.V(car.CarState.Direction.X, car.CarState.Direction.Y)

	// Rotation matrix for the car rendering
	rotationAngle := dirVector.Angle() + halfPi
	rotation := pixel.IM.Rotated(posVector, dirVector.Angle()+halfPi)
	carRender.SetMatrix(rotation)

	// Sensor vectors
	sensorUnitVector := pixel.Unit(math.Pi / 2)            // unit vector pointing forward from the car
	sensorCenter := carCenter.Add(pixel.V(carWidth*.5, 0)) // middle of the front of the car

	// render car fill
	if car.CarState.Crashed != true {
		carRender.Color = carColor
	} else {
		carRender.Color = carCrashedColor
	}

	carRender.Push(
		carCenter,
		pixel.V(car.CarState.Position.X+carWidth*.5, car.CarState.Position.Y+carLength*.5),
	)
	carRender.Rectangle(0)

	// render car outline
	carRender.Color = colornames.Black
	carRender.Push(
		carCenter,
		pixel.V(car.CarState.Position.X+carWidth*.5, car.CarState.Position.Y+carLength*.5),
	)
	carRender.Rectangle(1)

	// render car middle point
	carRender.Color = colornames.Yellow
	carRender.Push(sensorCenter)
	carRender.Circle(2, 0)

	// render sensors
	carRender.Color = colornames.Orange
	for i := 0; i < len(car.CarState.Sensors); i++ {
		sensor := car.CarState.Sensors[i]
		sensorVector := sensorUnitVector.Scaled(sensor.Distance).Rotated(sensor.Angle)
		carRender.Push(
			sensorCenter,
			sensorCenter.Add(sensorVector),
		)
		carRender.Line(1)
	}

	renderWheels(car, rotation, rotationAngle, carRender)

	return carRender
}

func renderWheels(car domain.Car, initialMatrix pixel.Matrix, rotationAngle float64, carRender *imdraw.IMDraw) {
	carRender.Color = colornames.Black

	wheelPositions := [4][3]float64{
		// carWidth, carLength, wheelOffset
		[3]float64{-1, -1, +1}, // top left
		[3]float64{+1, -1, +1}, // top right
		[3]float64{-1, +1, -1}, // bottom left
		[3]float64{+1, +1, -1}, // bottom right
	}

	for i := 0; i < len(wheelPositions); i++ {
		wheelData := wheelPositions[i]
		wheelCenter := pixel.V(
			car.CarState.Position.X+wheelData[0]*carWidth*.5,
			car.CarState.Position.Y+wheelData[1]*carLength*.5+wheelData[2]*wheelOffset,
		).Add(pixel.V(0, wheelData[2]*wheelLength/2))
		carRender.Color = colornames.Black
		wheelMatrix := pixel.IM.Moved(wheelCenter).Rotated(wheelCenter, rotationAngle)
		absoluteCenter := pixel.V(0, 0)
		if i < 2 {
			// front 2 wheels - ugly, whether the wheels turn should be a param somewhere
			wheelMatrix = wheelMatrix.Rotated(wheelCenter, car.CarControlState.Steering*maxWheelSteering)
		}
		carRender.SetMatrix(wheelMatrix)
		lengthVec := pixel.V(0, wheelLength/2)
		carRender.Push(
			absoluteCenter.Add(lengthVec),
			absoluteCenter.Sub(lengthVec),
		)
		carRender.Line(wheelWidth)
		carRender.Color = colornames.Red
		carRender.Push(absoluteCenter)
		carRender.Circle(1, 0)
		// reset the matrix
		carRender.SetMatrix(initialMatrix)
	}
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
