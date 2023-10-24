package v1

// Application represents a version 1 API data object for pesticide application information.
type Application struct {
	// The unique PICOL identifier for the application.
	Id int

	// The name of the application.
	Name string

	// Single character application code.
	Code string
}
