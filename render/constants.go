package render

import (
	"math"

	"golang.org/x/image/colornames"
)

const lineThickness = 4
const segmentLength = 15
const carWidth = 15
const carLength = 30
const wheelOffset = carLength / 7
const wheelLength = carLength / 4
const wheelWidth = 3
const maxWheelSteering = math.Pi / 4
const zoom = 3.0

var bgColor = colornames.Green
var roadColor = colornames.Gray
var outerLinesColor = colornames.White
var centerLineColor = colornames.Yellow
var segmentColor = colornames.Red
var carColor = colornames.Cyan
var carCrashedColor = colornames.Red
var startLineColor = colornames.Lightgrey
