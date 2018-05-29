package render

import (
	"fmt"
	"time"

	"github.com/socialgorithm/elon-server/simulator"

	"github.com/socialgorithm/elon-server/domain"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var simulation simulator.Simulation
var carStates []domain.CarState

func update() {
	for {
		select {
		case updatedCarStates := <-simulation.CarStatesChannel:
			carStates = updatedCarStates
		}
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Elon Self Driving - Socialgorithm",
		Bounds: pixel.R(0, 0, simulation.Track.Width, simulation.Track.Height),
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

	trackRender := renderTrack(simulation.Track)

	for !win.Closed() {
		win.Clear(bgColor)
		// redraw the track
		trackRender.Draw(win)

		// update cars
		if len(carStates) > 0 {
			for i := 0; i < len(carStates); i++ {
				carRender := renderCar(carStates[i])
				carRender.Draw(win)
			}
		}

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
func Render(_simulation simulator.Simulation) {
	simulation = _simulation
	go update()
	pixelgl.Run(run)
}
