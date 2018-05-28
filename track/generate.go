package track

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/socialgorithm/elon-server/domain"
)

func genTrack(_width int32, _height int32) domain.Track {
	rand.Seed(time.Now().UnixNano())

	points := 20
	margin := float32(50)
	maxPoints := 200
	width := float32(_width)
	height := float32(_height)

	var randPoints = make([]domain.Position, points, maxPoints)

	for i := 0; i < points; i++ {
		x := rand.Float32()*(width-margin) + margin
		y := rand.Float32()*(height-margin) + margin
		fmt.Println(x, y)
		randPoints = append(
			randPoints,
			domain.Position{
				X: x,
				Y: y,
			},
		)
	}

	fmt.Println(randPoints)

	firstSide := findConvexHull(randPoints)

	return domain.Track{
		FirstSide:  firstSide,
		SecondSide: firstSide,
	}
}
