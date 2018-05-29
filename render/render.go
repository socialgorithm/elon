package render

import (
	"fmt"
	"time"

	"github.com/socialgorithm/elon-server/domain"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var carStateChannel <-chan domain.CarState
var trackObj domain.Track

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

	trackRender := renderTrack(trackObj)

	for !win.Closed() {
		win.Clear(bgColor)
		// redraw the track
		trackRender.Draw(win)

		// update cars
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

// Render initiates the render loop
func Render(_trackObj domain.Track, _carStateChannel <-chan domain.CarState) {
	trackObj = _trackObj
	carStateChannel = _carStateChannel
	pixelgl.Run(run)
}
