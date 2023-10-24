package ddbmodel

type IntendedUser int8

const (
	IntendedUserCommercial IntendedUser = 1
	IntendedUserHome       IntendedUser = 2
)

func (iu IntendedUser) Code() byte {
	switch iu {
	case IntendedUserCommercial:
		return 'C'
	case IntendedUserHome:
		return 'H'
	}
	panic("unknown intended user")
}

func (iu IntendedUser) Name() string {
	switch iu {
	case IntendedUserCommercial:
		return "Commercial"
	case IntendedUserHome:
		return "Home"
	}
	panic("unknown intended user")
}

func (iu IntendedUser) String() string {
	switch iu {
	case IntendedUserCommercial:
		return "Commercial"
	case IntendedUserHome:
		return "Home"
	}
	panic("unknown intended user")
}
