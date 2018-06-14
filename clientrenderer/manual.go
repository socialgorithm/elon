package clientrenderer

import (
	"fmt"
	"time"

	"github.com/socialgorithm/elon-server/domain"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var inputs chan domain.CarControlState

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Elon Manual Driver - Socialgorithm",
		Bounds: pixel.R(0, 0, 400, 200),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetSmooth(true)

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	for !win.Closed() {
		win.Clear(colornames.Beige)

		throttle := 0
		steering := 0

		if win.Pressed(pixelgl.KeyLeft) {
			steering = -1
		}
		if win.Pressed(pixelgl.KeyRight) {
			steering = 1
		}
		if win.Pressed(pixelgl.KeyDown) {
			throttle = -1
		}
		if win.Pressed(pixelgl.KeyUp) {
			throttle = 1
		}

		carControlState := domain.CarControlState{
			Steering: float64(steering),
			Throttle: float64(throttle),
		}

		inputs <- carControlState

		renderCarState(carControlState).Draw(win)

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

// Manual initiates the render loop
func Manual(_inputs chan domain.CarControlState) {
	inputs = _inputs
	pixelgl.Run(run)
}
