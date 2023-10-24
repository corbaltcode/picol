package v1

// StateRecord represents a version 1 API data object for state records.
type StateRecord struct {
	// The unique PICOL identifier for the state record.
	Id int

	// The PICOL identifier for the state
	StateId int

	// The name of the state.
	Name string

	// The agency identifier.
	AgencyId string

	// The version of the state registration.
	Version string

	// The registration year.
	Year int

	// Indicates whether this is approved for use on cannabis production under WA I-502.
	I502 bool

	// Indicates whether this is approved for use on industrial hemp production under WA ESSB 6206.
	Essb6206 bool
}
