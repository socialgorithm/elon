package track

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	imdraw "github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/socialgorithm/elon-server/domain"
	"golang.org/x/image/colornames"
)

func drawTrack() *imdraw.IMDraw {

	rand.Seed(time.Now().UnixNano())

	points := 20
	maxPoints := 100
	width := int32(1024)
	height := int32(768)

	var pointsLeft = make([]domain.Position, points, maxPoints)
	var pointsRight = make([]domain.Position, points, maxPoints)

	for i := 0; i < points; i++ {
		pointsLeft = append(pointsLeft, domain.Position{X: float32(rand.Int31n(width)), Y: float32(rand.Int31n(height))})
		pointsRight = append(pointsRight, domain.Position{X: float32(rand.Int31n(width)), Y: float32(rand.Int31n(height))})
	}

	track := imdraw.New(nil)
	track.Color = pixel.RGB(1, 0, 0)

	for i := 0; i < len(pointsLeft)-1; i = i + 2 {
		pointA := pointsLeft[i]
		pointB := pointsLeft[i+1]
		track.Push(
			pixel.V(float64(pointA.X), float64(pointA.Y)),
			pixel.V(float64(pointB.X), float64(pointB.Y)),
		)
		track.Line(10)
	}
	track.Color = pixel.RGB(0, 0, 1)

	for i := 0; i < len(pointsRight)-1; i = i + 2 {
		pointA := pointsRight[i]
		pointB := pointsRight[i+1]
		track.Push(
			pixel.V(float64(pointA.X), float64(pointA.Y)),
			pixel.V(float64(pointB.X), float64(pointB.Y)),
		)
		track.Line(10)
	}

	return track
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Elon Self Driving - Socialgorithm",
		Bounds: pixel.R(0, 0, 1024, 768),
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

	track := drawTrack()

	for !win.Closed() {
		win.Clear(colornames.Skyblue)
		track.Draw(win)
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
