package render

import (
	"fmt"
	"log"
	"time"

	"github.com/socialgorithm/elon-server/simulator"

	"github.com/socialgorithm/elon-server/domain"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var simulation *simulator.Simulation
var cars []domain.Car

func update() {
	for {
		select {
		case updatedCars := <-simulation.CarsChannel:
			cars = updatedCars
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

		// follow top car with camera and zoom
		camPos := pixel.V(
			cars[0].CarState.Position.X,
			cars[0].CarState.Position.Y,
		)
		cam := pixel.IM.Scaled(camPos, zoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		// redraw the track
		trackRender.Draw(win)

		// update cars
		if len(cars) > 0 {
			for i := range cars {
				carRender := renderCar(cars[i])
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
func Render(_simulation *simulator.Simulation) {
	log.Println("Rendering simulation")
	simulation = _simulation
	go update()
	pixelgl.Run(run)
}
