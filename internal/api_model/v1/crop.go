package v1

// Crop represents a version 1 API data object for crop information.
type Crop struct {
	// The unique identifer for the crop.
	Id int

	// The name of the crop.
	Name string

	// Four-character crop code.
	Code string

	// Notes about the crop.
	Notes string
}
