package track

import (
	"math/rand"
	"time"

	"github.com/socialgorithm/elon-server/domain"
)

func genTrack() domain.Track {
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

	return domain.Track{
		FirstSide:  pointsLeft,
		SecondSide: pointsRight,
	}
}
