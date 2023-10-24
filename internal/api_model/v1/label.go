package v1

// Label represents a version 1 API data object for pesticide label information.
type Label struct {
	// The unique PICOL identifier for the pesticide label.
	Id int

	// The name of the label.
	Name string

	// The EPA number.
	EpaNumber string

	// The intended user of the pesticide.
	IntendedUser IntendedUser

	// Ingredients in the pesticide.
	Ingredients []Ingredient

	// The type(s) of this pesticide.
	PesticideTypes []PesticideType

	// The registrant of the pesticide.
	Registrant Registrant

	// The specialized local need (SLN) registration number.
	Sln string

	// The name of the specialized local need (SLN).
	SlnName string

	// The SLN expiration.
	SlnExpiration *AwfulDate `json:",omitempty"`

	// State records related to the pesticide.
	StateRecords []StateRecord

	// Supplemental code.
	Supplemental string

	// The name of the supplemental.
	SupplementalName string

	// The supplemental expiration.
	SupplementalExpiration *AwfulDate `json:",omitempty"`

	// The formulation code.
	Formulation string

	// The signal word.
	SignalWord string

	// Intended usage.
	Usage string

	// Whether the label is Organic Materials Research Institute (OMRI)-certified organic.
	Organic *bool `json:",omitempty"`

	// Whether the label has an Endangered Species Act (ESA) notice.
	EsaNotice *bool `json:",omitempty"`

	// EPA Section 18 emergency exemption.
	Section18 string
}
