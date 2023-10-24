package v1

// Ingredient represents a version 1 API data object for pesticide ingredient information.
type Ingredient struct {
	// The unique PICOL identifier for the ingredient.
	Id int

	// The name of the ingredient.
	Name string

	// Six-digit ingredient code. Leading zeros are significant, so this is stored as a string.
	Code string

	// Notes about the ingredient.
	Notes string

	// Resistance information about the ingredient.
	Resistance Resistance
}
