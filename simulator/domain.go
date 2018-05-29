package simulator

import (
	"github.com/socialgorithm/elon-server/domain"
)

// Simulation represents the current simulation
type Simulation struct {
	Track                   domain.Track
	CarStatesChannel        chan []domain.CarState
	CarControlStateReceiver chan<- domain.CarControlState
}
