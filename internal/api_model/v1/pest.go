package v1

// Pest represents a version 1 API data object for pest information.
type Pest struct {
	// The unique PICOL identifier for the pest.
	Id int

	// The name of the pest.
	Name string

	// The four- or five-character pest code.
	Code string

	// Notes about the pest.
	Notes string
}
