package track

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const lineThickness = 2
const radius = 5
const width = 1024
const height = 768

// func drawTrack(track domain.Track) *imdraw.IMDraw {

// 	trackRender := imdraw.New(nil)

// 	for i := 0; i < len(track.FirstSide)-1; i++ {
// 		pointA := track.FirstSide[i]
// 		pointB := track.FirstSide[i+1]
// 		trackRender.Color = colornames.Orange
// 		trackRender.Push(
// 			pixel.V(float64(pointA.X), float64(pointA.Y)),
// 		)
// 		trackRender.Circle(radius, 0)

// 		trackRender.Color = colornames.Gray
// 		trackRender.Push(
// 			pixel.V(float64(pointA.X), float64(pointA.Y)),
// 			pixel.V(float64(pointB.X), float64(pointB.Y)),
// 		)
// 		trackRender.Line(lineThickness)
// 	}

// 	// draw random points
// 	trackRender.Color = colornames.Whitesmoke
// 	for i := 0; i < len(track.RandomPoints); i++ {
// 		point := track.RandomPoints[i]
// 		trackRender.Push(
// 			pixel.V(float64(point.X), float64(point.Y)),
// 		)
// 		trackRender.Circle(radius, 0)
// 	}

// 	// Draw second side
// 	trackRender.Color = colornames.Gray
// 	for i := 0; i < len(track.SecondSide)-1; i++ {
// 		pointA := track.SecondSide[i]
// 		pointB := track.SecondSide[i+1]
// 		// figure out why we have some 0
// 		if pointA.X == 0 || pointA.Y == 0 || pointB.X == 0 || pointB.Y == 0 {
// 			continue
// 		}
// 		trackRender.Push(
// 			pixel.V(float64(pointA.X), float64(pointA.Y)),
// 			pixel.V(float64(pointB.X), float64(pointB.Y)),
// 		)
// 		trackRender.Line(lineThickness)
// 	}

// 	return trackRender
// }

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Elon Self Driving - Socialgorithm",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	track := genTrack(width, height)
	trackRender := drawTrack(track)

	for !win.Closed() {
		win.Clear(colornames.Green)
		trackRender.Draw(win)
		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

// Main Render a track
func Main() {
	pixelgl.Run(run)
}
