package fhir

import (
	"time"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/units"
)

// DurationFromTime converts an R4 FHIR Time Element into an R4 FHIR Duration
// value.
//
// If the underlying time has Second-based precision, the returned Duration will
// also have seconds precision; otherwise this will fallback into nanosecond
// precision.
func DurationFromTime(t *dtpb.Time) *dtpb.Duration {
	// Constrain the time to 24 hours
	duration := time.Microsecond * time.Duration(t.GetValueUs()) % (time.Hour * 24)

	switch t.GetPrecision() {
	case dtpb.Time_SECOND:
		return Seconds(duration)
	case dtpb.Time_MILLISECOND:
		return Milliseconds(duration)
	case dtpb.Time_MICROSECOND:
		fallthrough
	default:
		return Microseconds(duration)
	}
}

// Duration creates a Duration proto with the provided value, computing the
// largest whole-unit of time that can be used to represent the time.
func Duration(d time.Duration) *dtpb.Duration {
	value := d.Nanoseconds()
	if value == 0 {
		return durationValue(float64(value), units.Days)
	}
	unitConversions := []struct {
		unit     units.Time
		duration time.Duration
	}{
		{units.Days, 24 * time.Hour},
		{units.Hours, time.Hour},
		{units.Minutes, time.Minute},
		{units.Seconds, time.Second},
		{units.Milliseconds, time.Millisecond},
		{units.Microseconds, time.Microsecond},
		{units.Nanoseconds, time.Nanosecond},
	}
	for _, conversion := range unitConversions {
		if d >= conversion.duration && d == d.Round(conversion.duration) {
			numUnits := d / conversion.duration
			return durationValue(float64(numUnits), conversion.unit)
		}
	}
	return durationValue(float64(value), units.Nanoseconds)
}

// Nanoseconds creates a Duration proto with the specified time value, rounded
// to nanosecond accuracy.
func Nanoseconds(value time.Duration) *dtpb.Duration {
	return durationValue(float64(value.Nanoseconds()), units.Nanoseconds)
}

// Milliseconds creates a Duration proto with the specified time value, rounded
// to millisecond accuracy.
func Milliseconds(value time.Duration) *dtpb.Duration {
	millis := float64(value.Nanoseconds()) / float64(time.Millisecond.Nanoseconds())
	return durationValue(millis, units.Milliseconds)
}

// Microseconds creates a Duration proto with the specified time value, rounded
// to microsecond accuracy.
func Microseconds(value time.Duration) *dtpb.Duration {
	micros := float64(value.Nanoseconds()) / float64(time.Microsecond.Nanoseconds())
	return durationValue(micros, units.Microseconds)
}

// Seconds creates a Duration proto with the specified time value, rounded
// to second accuracy.
func Seconds(value time.Duration) *dtpb.Duration {
	return durationValue(value.Seconds(), units.Seconds)
}

// Minutes creates a Duration proto with the specified time value, rounded
// to minute accuracy.
func Minutes(value time.Duration) *dtpb.Duration {
	return durationValue(value.Minutes(), units.Minutes)
}

// Hours creates a Duration proto with the specified time value, rounded
// to hour-accuracy.
func Hours(value time.Duration) *dtpb.Duration {
	return durationValue(value.Hours(), units.Hours)
}

// Days creates a Duration proto with the specified time value, rounded
// to day-accuracy.
func Days(value time.Duration) *dtpb.Duration {
	return durationValue(value.Hours()/24, units.Days)
}

func durationValue(value float64, unit units.Time) *dtpb.Duration {
	return &dtpb.Duration{
		Value:  Decimal(value),
		Code:   Code(unit.Symbol()),
		System: URI(unit.System()),
	}
}
