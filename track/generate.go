package track

import (
	"math/rand"
	"time"

	"github.com/socialgorithm/elon-server/domain"
)

func genTrack(_width int32, _height int32) domain.Track {
	rand.Seed(time.Now().UnixNano())

	const points = 10
	margin := float32(50)
	width := float32(_width)
	height := float32(_height)

	var randPoints = [points]domain.Position{}

	for i := 0; i < points; i++ {
		x := rand.Float32()*(width-margin) + margin
		y := rand.Float32()*(height-margin) + margin
		randPoints[i] = domain.Position{
			X: x,
			Y: y,
		}
	}

	firstSide := findConvexHull(randPoints[0:len(randPoints)])

	return domain.Track{
		FirstSide:  firstSide,
		SecondSide: firstSide,
	}
}
