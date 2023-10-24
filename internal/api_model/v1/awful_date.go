package v1

import "fmt"

// AwfulDate is a version 1 API data object representing date information, serialized in the format MM/DD/YY,
// and will overflow in the year 2100.
type AwfulDate struct {
	// The year of the date.
	Year uint16

	// The month of the date.
	Month uint8

	// The day of the date.
	Day uint8
}

// NewAwfulDate creates a new AwfulDate object with the given year, month, and day.
//
// If year < 2000 or greater than 2099, month is not between 1 and 12, or day is not appropriate for the given month,
// NewNaiveDate will return an error.
func NewAwfulDate(year uint, month uint, day uint) (AwfulDate, error) {
	if year < 2000 || year > 2099 {
		return AwfulDate{}, fmt.Errorf("year must be between 2000 and 2099: %v", year)
	}

	if month < 1 || month > 12 {
		return AwfulDate{}, fmt.Errorf("month must be between 1 and 12: %v", month)
	}

	if month == 1 || month == 3 || month == 5 || month == 7 || month == 8 || month == 10 || month == 12 {
		if day == 0 || day > 31 {
			return AwfulDate{}, fmt.Errorf("day must be between 1 and 31 for month %v: %v", month, day)
		}
	} else if month == 4 || month == 6 || month == 9 || month == 11 {
		if day == 0 || day > 30 {
			return AwfulDate{}, fmt.Errorf("day must be between 1 and 30 for month %v: %v", month, day)
		}
	} else /* month == 2 */ {
		upperBound := uint(28)
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			upperBound = 29
		}

		if day == 0 || day > upperBound {
			return AwfulDate{}, fmt.Errorf("day must be between 1 and 29 for year %v month %v: %v", year, month, day)
		}
	}

	return AwfulDate{
		Year:  uint16(year),
		Month: uint8(month),
		Day:   uint8(day),
	}, nil
}

func (ad AwfulDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%02d/%02d/%02d\"", ad.Month, ad.Day, ad.Year-2000)), nil
}

func (ad *AwfulDate) UnmarshalJSON(b []byte) error {
	var month, day, year uint
	n, err := fmt.Sscanf(string(b), "\"%02u/%02u/%02u\"", &month, &day, &year)
	if err != nil {
		return err
	}

	if n != 3 {
		return fmt.Errorf("expected 3 values, got %v", n)
	}

	newAD, err := NewAwfulDate(year+2000, month, day)
	if err != nil {
		return err
	}

	ad.Year = newAD.Year
	ad.Month = newAD.Month
	ad.Day = newAD.Day
	return nil
}
