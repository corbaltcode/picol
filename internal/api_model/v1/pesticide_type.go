package v1

// PesticideType represents a version 1 API data object for pesticide type information.
type PesticideType struct {
	// The unique PICOL identifier for the pesticide type.
	Id int

	// The name of the pesticide type.
	Name string

	// The three- or four-character pesticide type code.
	Code string
}
