package v1

// IntendedUser represents a version 1 API data object for the intended user of a pesticide.
type IntendedUser struct {
	// The unique PICOL identifier for the intended user.
	Id int

	// The name of the intended user.
	Name string

	// Single character intended user code.
	Code string
}
