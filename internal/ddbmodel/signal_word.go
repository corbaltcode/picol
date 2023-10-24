package ddbmodel

type SignalWord int8

const (
	SignalWordCaution      SignalWord = 1
	SignalWordDanger       SignalWord = 2
	SignalWordDangerPoison SignalWord = 3
	SignalWordWarning      SignalWord = 4
	SignalWordNone         SignalWord = 5
)

func (sw SignalWord) Code() byte {
	switch sw {
	case SignalWordCaution:
		return 'C'
	case SignalWordDanger:
		return 'D'
	case SignalWordDangerPoison:
		return 'T'
	case SignalWordWarning:
		return 'W'
	case SignalWordNone:
		return 'N'
	}
	panic("unknown signal word")
}

func (sw SignalWord) Name() string {
	switch sw {
	case SignalWordCaution:
		return "CAUTION"
	case SignalWordDanger:
		return "DANGER"
	case SignalWordDangerPoison:
		return "DANGER/POISON"
	case SignalWordWarning:
		return "WARNING"
	case SignalWordNone:
		return "NO SIGNAL WORD GIVEN"
	}
	panic("unknown signal word")
}

func (sw SignalWord) String() string {
	switch sw {
	case SignalWordCaution:
		return "Caution"
	case SignalWordDanger:
		return "Danger"
	case SignalWordDangerPoison:
		return "Danger/Poison"
	case SignalWordWarning:
		return "Warning"
	case SignalWordNone:
		return "None"
	}
	panic("unknown signal word")
}
