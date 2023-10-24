package v1

// Registrant represents a version 1 API data object for registrant information.
type Registrant struct {
	// The unique PICOL identifier for the registrant.
	Id int

	// The name of the registrant.
	Name string

	// The registrant's website.
	Website string
}
