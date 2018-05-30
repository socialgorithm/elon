package simulator

import (
	"github.com/socialgorithm/elon-server/domain"
)

// Simulation represents the current simulation
type Simulation struct {
	Track       domain.Track
	Cars        []domain.Car
	CarsChannel chan []domain.Car
	Started     bool
}
