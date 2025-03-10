package fhirconv

import (
	"fmt"
	"strconv"
	"time"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/units"
)

// DateTimeToTime converts a FHIR DateTime element into a Go time.Time value.
//
// This function will only error if the TimeZone field is invalid.
//
// Note: the error that can be returned from this function is unlikely to actually
// occur in practice. FHIR Timezones are always required to be specified in the
// form `+zz:zz` or `-zz:zz` -- and the JSON conversion translates this into
// exactly this form, an empty string, or "UTC" in some cases. These are all
// valid locations and should never fail as a result.
func DateTimeToTime(dt *dtpb.DateTime) (time.Time, error) {
	tz, err := parseLocation(dt.GetTimezone())
	if err != nil {
		return time.Time{}, fmt.Errorf("fhirconv.TimeFromDateTime: %w", err)
	}
	return time.UnixMicro(dt.GetValueUs()).In(tz), nil
}

// InstantToTime converts a FHIR Instant element into a Go time.Time value.
//
// This function will only error if the TimeZone field is invalid.
//
// Note: the error that can be returned from this function is unlikely to actually
// occur in practice. FHIR Timezones are always required to be specified in the
// form `+zz:zz` or `-zz:zz` -- and the JSON conversion translates this into
// exactly this form, an empty string, or "UTC" in some cases. These are all
// valid locations and should never fail as a result.
func InstantToTime(dt *dtpb.Instant) (time.Time, error) {
	tz, err := parseLocation(dt.GetTimezone())
	if err != nil {
		return time.Time{}, fmt.Errorf("fhirconv.TimeFromInstant: %w", err)
	}
	return time.UnixMicro(dt.GetValueUs()).In(tz), nil
}

// DateToTime converts a FHIR Date element into a Go time.Time value.
//
// This function will only error if the TimeZone field is invalid.
//
// Note: the error that can be returned from this function is unlikely to actually
// occur in practice. FHIR Timezones are always required to be specified in the
// form `+zz:zz` or `-zz:zz` -- and the JSON conversion translates this into
// exactly this form, an empty string, or "UTC" in some cases. These are all
// valid locations and should never fail as a result.
func DateToTime(dt *dtpb.Date) (time.Time, error) {
	tz, err := parseLocation(dt.GetTimezone())
	if err != nil {
		return time.Time{}, fmt.Errorf("fhirconv.Date: %w", err)
	}
	return time.UnixMicro(dt.GetValueUs()).In(tz), nil
}

// TimeToDuration converts a FHIR Time element into a Go time.Duration value.
//
// Despite the name `Time` for the FHIR Element, the time is not actually
// associated to any real date -- and thus does not correspond to a distinct
// chronological point, and thus cannot be converted logically into a `time.Time`
// object.
func TimeToDuration(dt *dtpb.Time) time.Duration {
	return time.Microsecond * time.Duration(dt.GetValueUs())
}

// parseLocation attempts to parse the timezone location from the zone string.
//
// Timezones may be specified in one of 3 formats:
//   - Z
//   - +zz:zz or -zz:zz
//   - UTC (or some name)
//
// Additionally, this function supports empty strings being translated into
// UTC.
func parseLocation(zone string) (*time.Location, error) {
	if zone == "" {
		return time.UTC, nil
	}
	if tm, err := time.Parse("MST", zone); err == nil {
		return tm.Location(), nil
	}
	if tm, err := time.Parse("Z07:00", zone); err == nil {
		return tm.Location(), nil
	}
	if zone == "Local" {
		return time.Local, nil
	}
	return nil, fmt.Errorf("unable to parse time-zone from '%v'", zone)
}

// DurationToDuration converts a FHIR Duration element into a Go native
// time.Duration object.
//
// This function may return an error in the following conditions:
//   - The underlying Decimal value is not able to be parsed into a float64
//   - The unit is not a valid time unit
func DurationToDuration(d *dtpb.Duration) (time.Duration, error) {
	value := d.GetValue().GetValue()
	decimal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return durationToDurationError("bad decimal value '%v': %v", value, err)
	}
	code := d.GetCode().GetValue()

	unit, err := units.TimeFromSymbol(code)
	if err != nil {
		return durationToDurationError("invalid unit symbol '%v'", code)
	}

	symbol := unit.Symbol()

	// Special handling is necessary as days are not supported by
	// time.ParseDuration, as well as the minutes unit being m, not min
	switch unit {
	case units.Minutes:
		symbol = "m"
	case units.Days:
		decimal *= 24
		symbol = units.Hours.Symbol()
	}

	duration, err := time.ParseDuration(fmt.Sprintf("%v%s", decimal, symbol))
	if err != nil {
		// This branch should not be possible to be reached. If we have reached this,
		// something really bad has happened -- because we form the format string
		// manually above.
		return durationToDurationError("%v", err)
	}
	return duration, nil
}

func durationToDurationError(format string, args ...any) (time.Duration, error) {
	message := fmt.Sprintf(format, args...)
	return time.Duration(0), fmt.Errorf("fhirconv.DurationToDuration: %v", message)
}
