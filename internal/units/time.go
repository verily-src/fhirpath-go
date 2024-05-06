package units

import "fmt"

// Time is a unit of measure for measuring the passage of time.
type Time int

const (
	// Nanoseconds is a Time unit that measures time in nanoseconds.
	Nanoseconds Time = iota

	// Microseconds is a Time unit that measures time in microseconds.
	Microseconds

	// Milliseconds is a Time unit that measures time in milliseconds.
	Milliseconds

	// Seconds is a Time unit that measures time in seconds.
	Seconds

	// Minutes is a Time unit that measures time in minutes.
	Minutes

	// Hours is a Time unit that measures time in hours.
	Hours

	// Days is a Time unit that measures time in days.
	Days
)

const (
	nanosecondSymbol  = "ns"
	microsecondSymbol = "us"
	millisecondSymbol = "ms"
	secondsSymbol     = "s"
	minutesSymbol     = "min"
	hoursSymbol       = "h"
	daysSymbol        = "d"
)

// Symbol returns the symbol used to represent the underlying unit.
func (t Time) Symbol() string {
	switch t {
	case Nanoseconds:
		return nanosecondSymbol
	case Microseconds:
		return microsecondSymbol
	case Milliseconds:
		return millisecondSymbol
	case Seconds:
		return secondsSymbol
	case Minutes:
		return minutesSymbol
	case Hours:
		return hoursSymbol
	case Days:
		return daysSymbol
	}
	// This is a closed enumeration in an internal package. If this panic ever
	// gets reached, it means that a developer is using this package wrong.
	panic(fmt.Sprintf("invalid time value %v", t))
}

// System returns the time system that this unit comes from.
func (t Time) System() string {
	return "http://unitsofmeasure.org"
}

// TimeFromSymbol creates the Time object
func TimeFromSymbol(symbol string) (Time, error) {
	switch symbol {
	case nanosecondSymbol:
		return Nanoseconds, nil
	case microsecondSymbol:
		return Microseconds, nil
	case millisecondSymbol:
		return Milliseconds, nil
	case secondsSymbol:
		return Seconds, nil
	case minutesSymbol:
		return Minutes, nil
	case hoursSymbol:
		return Hours, nil
	case daysSymbol:
		return Days, nil
	}
	return Time(0), fmt.Errorf("unknown Time symbol '%v'", symbol)
}
