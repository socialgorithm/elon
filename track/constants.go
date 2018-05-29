package track

import "golang.org/x/image/colornames"

const lineThickness = 4
const radius = 5
const width = 1024
const height = 768
const points = 50
const difficulty = float64(0.2) // closer to 0 will create sharper turns, exponentially
const maxDisplacement = float64(200)
const margin = float64(100)
const roadWidth = float64(25)
const segmentLength = 15

var bgColor = colornames.Green
var roadColor = colornames.Gray
var outerLinesColor = colornames.White
var centerLineColor = colornames.Yellow
var segmentColor = colornames.Red
