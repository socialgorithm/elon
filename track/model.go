package track

import "github.com/socialgorithm/elon-server/domain"

// Polygon to model the race track
type Polygon struct {
	vertices []domain.Position
	minX     int
	minY     int
	maxX     int
	maxY     int
	closed   bool
}

type Edge struct {
	vertex1 domain.Position
	vertex2 domain.Position
}
