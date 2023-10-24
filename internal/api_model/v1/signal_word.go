package v1

// SignalWord represents a version 1 API data object for signal word information.
type SignalWord struct {
	// The unique PICOL identifier for the signal word.
	Id int

	// The name of the signal word.
	Name string

	// The single-character signal word code.
	Code string
}
