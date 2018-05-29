package domain

// Car represents each car in the simulation
type Car struct {
	CarControlState CarControlState
	CarState        CarState
}

// CarControlState represents a command to alter the behaviour of a car
type CarControlState struct {
	Throttle float64
	Steering float64
}

// CarState represents the current reportable state of a car
type CarState struct {
	Position  Position
	Direction Position // normalized vector
	Velocity  float64
	Sensors   []Sensor
	Crashed   bool
}

// Position represents an (X, Y) position
type Position struct {
	X float64
	Y float64
}

// Sensor represents the current state of a sensor
type Sensor struct {
	Angle    float64 // radians [0, 2 * pi]
	Distance float64
}

// Track represents a track that can be traversed by cars
type Track struct {
	InnerSide []Position
	OuterSide []Position
	Center    []Position

	Width  float64
	Height float64
}
