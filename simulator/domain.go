package simulator

import (
	"github.com/socialgorithm/elon-server/domain"
	"github.com/socialgorithm/elon-server/physics"
)

// Simulation represents the current simulation - need car state/control states decoupled
// as they have different sources
type Simulation struct {
	Track       domain.Track
	CarsChannel chan []domain.Car
	Engine      physics.Engine
}
