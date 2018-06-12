package simulator

import (
	"github.com/socialgorithm/elon-server/domain"
	"github.com/socialgorithm/elon-server/physics"
)

// Simulation represents the current simulation - need car state/control states decoupled
// as they have different sources
type Simulation struct {
	Track                   domain.Track
	CarStatesChannel        chan []domain.CarState
	CarControlStateReceiver chan<- domain.CarControlState
	Engine                  physics.Engine
}
