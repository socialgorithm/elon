package domain

// CarControlState represents a command to alter the behaviour of a car
type CarControlState struct {
	Throttle float32
	Steering float32
}

// CarState represents the current reportable state of a car
type CarState struct {
	Position Position
	Velocity float32
	Sensors  []Sensor
}

// Position represents an (X, Y) position
type Position struct {
	X float32
	Y float32
}

// Sensor represents the current state of a sensor
type Sensor struct {
	Angle    int8
	Distance float32
}

// Track represents a track that can be traversed by cars
type Track struct {
	FirstSide  []Position
	SecondSide []Position
}
