package track

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	imdraw "github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/socialgorithm/elon-server/domain"
	"golang.org/x/image/colornames"
)

const lineThickness = 3
const width = 1024
const height = 768

func drawTrack(track domain.Track) *imdraw.IMDraw {

	trackRender := imdraw.New(nil)
	trackRender.Color = pixel.RGB(1, 0, 0)

	for i := 0; i < len(track.FirstSide)-1; i++ {
		pointA := track.FirstSide[i]
		pointB := track.FirstSide[i+1]
		trackRender.Push(
			pixel.V(float64(pointA.X), float64(pointA.Y)),
			pixel.V(float64(pointB.X), float64(pointB.Y)),
		)
		trackRender.Line(lineThickness)
	}
	trackRender.Color = pixel.RGB(0, 0, 1)

	// for i := 0; i < len(track.SecondSide)-1; i = i + 2 {
	// 	pointA := track.SecondSide[i]
	// 	pointB := track.SecondSide[i+1]
	// 	trackRender.Push(
	// 		pixel.V(float64(pointA.X), float64(pointA.Y)),
	// 		pixel.V(float64(pointB.X), float64(pointB.Y)),
	// 	)
	// 	trackRender.Line(lineThickness)
	// }

	return trackRender
}

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

	win.Clear(colornames.Skyblue)

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	track := genTrack(width, height)
	trackRender := drawTrack(track)

	for !win.Closed() {
		win.Clear(colornames.Skyblue)
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
