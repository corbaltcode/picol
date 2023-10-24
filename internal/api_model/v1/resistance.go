package v1

// Resistance represents a version 1 API data object for resistance information.
type Resistance struct {
	// The unique PICOL identifier for the resistance.
	Id int

	// Four-character source code.
	Source string

	// Alphanumeric resistance code.
	Code string

	// The method of action for the resistance.
	MethodOfAction string
}
