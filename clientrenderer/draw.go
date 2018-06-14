package clientrenderer

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/socialgorithm/elon-server/domain"
	"golang.org/x/image/colornames"
)

const dashWidth = 250.0
const dashHeight = 100.0

var dashTopLeft = pixel.V(20, 150)
var throttleColor = colornames.Blue
var brakeColor = colornames.Red
var steeringColor = colornames.Green

func renderCarState(carState domain.CarControlState) *imdraw.IMDraw {
	stateRender := imdraw.New(nil)

	// Draw carState rectangles
	// Steering
	stateRender.Color = steeringColor
	barWidth := dashWidth / 2

	stateRender.Push(
		dashTopLeft.Add(pixel.V(barWidth, 0)),
		dashTopLeft.Add(pixel.V(barWidth+barWidth*carState.Steering, -dashHeight/2)),
	)

	stateRender.Rectangle(0)

	// Throttle
	if carState.Throttle >= 0 {
		stateRender.Color = throttleColor
	} else {
		stateRender.Color = brakeColor
	}

	stateRender.Push(
		dashTopLeft.Add(pixel.V(barWidth, -dashHeight/2)),
		dashTopLeft.Add(pixel.V(barWidth+barWidth*carState.Throttle, -dashHeight)),
	)

	stateRender.Rectangle(0)

	// Draw Grid
	stateRender.Color = colornames.Gray

	// Draw vertical divider
	stateRender.Push(
		dashTopLeft.Add(pixel.V(dashWidth/2, 0)),
		dashTopLeft.Add(pixel.V(dashWidth/2, -dashHeight)),
	)

	stateRender.Line(1)

	stateRender.Color = colornames.Black

	// Draw enclosing rectangle
	stateRender.Push(
		dashTopLeft,
		dashTopLeft.Add(pixel.V(dashWidth, -dashHeight)),
	)

	stateRender.Rectangle(2)

	// Draw center line
	stateRender.Push(
		dashTopLeft.Add(pixel.V(0, -dashHeight/2)),
		dashTopLeft.Add(pixel.V(dashWidth, -dashHeight/2)),
	)

	stateRender.Line(1)

	return stateRender
}
