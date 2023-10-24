package ddbmodel

type State int8

const (
	StateWashington State = 1
	StateOregon     State = 2
)

func (s State) Name() string {
	switch s {
	case StateWashington:
		return "Washington"
	case StateOregon:
		return "Oregon"
	}
	panic("unknown state")
}

func (s State) String() string {
	return s.Name()
}
