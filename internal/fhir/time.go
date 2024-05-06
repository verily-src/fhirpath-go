package fhir

import (
	"fmt"
	"time"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
)

// extractTimezone gets the string representation of a UTC offset.
//
// This source is taken from:
// https://github.com/google/fhir/blob/1d1e7189749fcdbbececc1c70e00dd498bfb33d1/go/jsonformat/internal/jsonpbhelper/fhirutil.go#L453-L463
func extractTimezone(t time.Time) string {
	_, offset := t.Zone()
	sign := "+"
	if offset < 0 {
		sign = "-"
		offset = -offset
	}
	hour := offset / 3600
	minute := offset % 3600 / 60
	return fmt.Sprintf("%s%02d:%02d", sign, hour, minute)
}

// Date creates an R4 FHIR Date element from a Time value, accurate to a given
// day.
//
// See: http://hl7.org/fhir/R4/datatypes.html#date
func Date(t time.Time) *dtpb.Date {
	return &dtpb.Date{
		ValueUs:   t.UnixMicro(),
		Precision: dtpb.Date_DAY,
		Timezone:  extractTimezone(t),
	}
}

// DateNow creates an R4 FHIR Date element at the current time using the
// highest available precision.
func DateNow() *dtpb.Date {
	return Date(time.Now())
}

// DateTime creates an R4 FHIR DateTime element from a Time value, accurate
// to the microsecond.
//
// See: http://hl7.org/fhir/R4/datatypes.html#datetime
func DateTime(t time.Time) *dtpb.DateTime {
	return &dtpb.DateTime{
		ValueUs:   t.UnixMicro(),
		Precision: dtpb.DateTime_MICROSECOND,
		Timezone:  extractTimezone(t),
	}
}

// DateTimeNow creates an R4 FHIR DateTime element at the current time using the
// highest available precision.
func DateTimeNow() *dtpb.DateTime {
	return DateTime(time.Now())
}

// Instant creates an R4 FHIR Instant element from a Time value, accurate to the
// microsecond.
//
// See: http://hl7.org/fhir/R4/datatypes.html#instant
func Instant(t time.Time) *dtpb.Instant {
	return &dtpb.Instant{
		ValueUs:   t.UnixMicro(),
		Precision: dtpb.Instant_MICROSECOND,
		Timezone:  extractTimezone(t),
	}
}

// InstantNow creates an R4 FHIR Instant element at the current time using the
// highest available precision.
func InstantNow() *dtpb.Instant {
	return Instant(time.Now())
}

// Time creates an R4 FHIR Time element from a Time value.
//
// FHIR Time elements represent a time of day, disconnected from any date.
// As a result, the value stored in this proto will be modulo 24-hours to keep
// it within that 1 day time. Put differently, this will only ever be populated
// with the number of microseconds since the start of the unix epoch, modulo one
// day in microseconds.
//
// See: http://hl7.org/fhir/R4/datatypes.html#time
func Time(t time.Time) *dtpb.Time {
	return &dtpb.Time{
		ValueUs:   t.UnixMicro() % (time.Hour * 24).Microseconds(),
		Precision: dtpb.Time_MICROSECOND,
	}
}

// TimeNow creates an R4 FHIR Time element at the current time using the
// highest available precision.
func TimeNow() *dtpb.Time {
	return Time(time.Now())
}

// TimeOfDay creates a Time proto at the specified time.
//
// This function will return an error if any of the values exceed the valid
// range for their time unit (e.g. if 'hour' exceeds 24, or minute exceeds 60, etc).
//
// The precision is determine by the value set for the micro second parameter;
// if the value is 0, the precision is set to seconds. If the value is
// a multiple of 1000, the value is set to millisecond precision; otherwise, its
// set to microsecond precision.
func TimeOfDay(hour, minute, second, micros int64) (*dtpb.Time, error) {
	const (
		maxHour   = 24
		maxMinute = 60
		maxSecond = 60
		maxMicros = 1_000_000
	)
	if hour >= maxHour || hour < 0 {
		return nil, fmt.Errorf("invalid hour '%v'; expected range is [0, %v)", hour, maxHour)
	}
	if minute >= maxMinute || minute < 0 {
		return nil, fmt.Errorf("invalid minute '%v'; expected range is [0, %v)", minute, maxMinute)
	}
	if second >= maxSecond || second < 0 {
		return nil, fmt.Errorf("invalid second '%v'; expected range is [0, %v)", second, maxSecond)
	}
	if micros >= maxMicros || micros < 0 {
		return nil, fmt.Errorf("invalid micros '%v'; expected range is [0, %v)", micros, maxMicros)
	}

	precision := dtpb.Time_MICROSECOND
	if micros == 0 {
		precision = dtpb.Time_SECOND
	} else if (micros % 1_000) == 0 {
		precision = dtpb.Time_MILLISECOND
	}

	return &dtpb.Time{
		ValueUs: hour*time.Hour.Microseconds() +
			minute*time.Minute.Microseconds() +
			second*time.Second.Microseconds() +
			micros,
		Precision: precision,
	}, nil
}

// ParseDate converts the input string into a FHIR Date element.
// The format of the input string must follow the FHIR Date format as defined
// in http://hl7.org/fhir/R4/datatypes.html#date, e.g.
//
//   - YYYY,
//   - YYYY-MM, or
//   - YYYY-MM-DD
//
// The returned Date will have a precision equal to what was specified in the
// input string.
func ParseDate(value string) (*dtpb.Date, error) {
	dateFormats := []struct {
		format    string
		precision dtpb.Date_Precision
	}{
		{"2006-01-02", dtpb.Date_DAY},
		{"2006-01", dtpb.Date_MONTH},
		{"2006", dtpb.Date_YEAR},
	}

	var t time.Time
	var err error
	for _, format := range dateFormats {
		t, err = time.Parse(format.format, value)
		if err == nil {
			return &dtpb.Date{
				ValueUs:   t.UnixMicro(),
				Timezone:  extractTimezone(t),
				Precision: format.precision,
			}, nil
		}
	}
	return nil, fmt.Errorf("unable to parse date '%v': %w", value, err)
}

// MustParseDate parses a date as according to ParseDate, but panics if the date
// is invalid.
func MustParseDate(value string) *dtpb.Date {
	result, err := ParseDate(value)
	if err != nil {
		panic(err)
	}
	return result
}

// ParseDateTime converts the input string into a FHIR DateTime element.
// The format of the input string must follow the FHIR DateTime format as defined
// in http://hl7.org/fhir/R4/datatypes.html#datetime, e.g.
//
//   - YYYY,
//   - YYYY-MM,
//   - YYYY-MM-DD, or
//   - YYYY-MM-DDThh:mm:ss+zz:zz (with optional milli/micro precision)
//
// The returned DateTime will have a precision equal to what was specified in the
// input string.
func ParseDateTime(value string) (*dtpb.DateTime, error) {
	dateFormats := []struct {
		format    string
		precision dtpb.DateTime_Precision
	}{
		{"2006-01-02T15:04:05.000000-07:00", dtpb.DateTime_MICROSECOND},
		{"2006-01-02T15:04:05.000000Z", dtpb.DateTime_MICROSECOND},
		{"2006-01-02T15:04:05.000-07:00", dtpb.DateTime_MILLISECOND},
		{"2006-01-02T15:04:05.000Z", dtpb.DateTime_MILLISECOND},
		{"2006-01-02T15:04:05-07:00", dtpb.DateTime_SECOND},
		{"2006-01-02T15:04:05Z", dtpb.DateTime_SECOND},
		{"2006-01-02", dtpb.DateTime_DAY},
		{"2006-01", dtpb.DateTime_MONTH},
		{"2006", dtpb.DateTime_YEAR},
	}

	var t time.Time
	var err error
	for _, format := range dateFormats {
		t, err = time.Parse(format.format, value)
		if err == nil {
			return &dtpb.DateTime{
				ValueUs:   t.UnixMicro(),
				Timezone:  extractTimezone(t),
				Precision: format.precision,
			}, nil
		}
	}
	return nil, fmt.Errorf("unable to parse datetime '%v': %w", value, err)
}

// MustParseDateTime parses a date as according to ParseDateTime, but panics if
// the date is invalid.
func MustParseDateTime(value string) *dtpb.DateTime {
	result, err := ParseDateTime(value)
	if err != nil {
		panic(err)
	}
	return result
}

// ParseInstant converts the input string into a FHIR Instant element.
// The format of the input string must follow the FHIR Instant format as defined
// in http://hl7.org/fhir/R4/datatypes.html#instant, e.g.
//
//   - yyyy-mm-ddThh:mm:ss+zz:zz,
//   - yyyy-mm-ddThh:mm:ss.000+zz:zz, or
//   - yyyy-mm-ddThh:mm:ss.000000+zz:zz,
//
// The returned Instant will have a precision equal to what was specified in the
// input string.
func ParseInstant(value string) (*dtpb.Instant, error) {
	dateFormats := []struct {
		format    string
		precision dtpb.Instant_Precision
	}{
		{"2006-01-02T15:04:05.000000-07:00", dtpb.Instant_MICROSECOND},
		{"2006-01-02T15:04:05.000000Z", dtpb.Instant_MICROSECOND},
		{"2006-01-02T15:04:05.000-07:00", dtpb.Instant_MILLISECOND},
		{"2006-01-02T15:04:05.000Z", dtpb.Instant_MILLISECOND},
		{"2006-01-02T15:04:05-07:00", dtpb.Instant_SECOND},
		{"2006-01-02T15:04:05Z", dtpb.Instant_SECOND},
	}

	var t time.Time
	var err error
	for _, format := range dateFormats {
		t, err = time.Parse(format.format, value)
		if err == nil {
			return &dtpb.Instant{
				ValueUs:   t.UnixMicro(),
				Timezone:  extractTimezone(t),
				Precision: format.precision,
			}, nil
		}
	}
	return nil, fmt.Errorf("unable to parse instant '%v': %w", value, err)
}

// MustParseInstant parses a date as according to ParseInstant, but panics if
// the time is invalid.
func MustParseInstant(value string) *dtpb.Instant {
	result, err := ParseInstant(value)
	if err != nil {
		panic(err)
	}
	return result
}

// yearZeroBase is the time set for 0000-01-01.
//
// The time.Parse has the opinionated choice that any parsed time without any
// year associated to it _must_ be offset from this date. Thus, we store this
// so that parsed times will be offset from the unix timestamp instead.
var yearZeroBase time.Time

func init() {
	tm, err := time.Parse("2006-01-02", "0000-01-01")
	if err != nil {
		panic(fmt.Sprintf("Unexpected error while creating base time: %v", err))
	}
	yearZeroBase = tm
}

// ParseTime converts the input string into a FHIR Time element.
// The format of the input string must follow the FHIR Time format as defined
// in http://hl7.org/fhir/R4/datatypes.html#time, e.g.
//
//   - hh:mm:ss,
//   - hh:mm:ss.000, or
//   - hh:mm:ss.000000,
//
// The returned Time will have a precision equal to what was specified in the
// input string.
func ParseTime(value string) (*dtpb.Time, error) {
	dateFormats := []struct {
		format    string
		precision dtpb.Time_Precision
	}{
		{"15:04:05.000000", dtpb.Time_MICROSECOND},
		{"15:04:05.000", dtpb.Time_MILLISECOND},
		{"15:04:05", dtpb.Time_SECOND},
	}

	var t time.Time
	var err error
	for _, format := range dateFormats {
		t, err = time.Parse(format.format, value)
		if err != nil {
			continue
		}

		// time.Parse without any date context forms a time offset from `0000-01-01`
		// for some strange reason. This causes a large negative unix timestamp
		// which is hard to make work for FHIR Time types. To compensate for this,
		// we subtract the 0000-01-01 timestamp from our time, so that 00:00:00 will
		// form a time of '0', e.g. it forces all times to be relative to the
		// unix epoch.
		value := t.UnixMicro() - yearZeroBase.UnixMicro()
		return &dtpb.Time{
			ValueUs:   value,
			Precision: format.precision,
		}, nil
	}
	return nil, fmt.Errorf("unable to parse time '%v': %w", value, err)
}

// MustParseTime parses a date as according to ParseTime, but panics if
// the time is invalid.
func MustParseTime(time string) *dtpb.Time {
	result, err := ParseTime(time)
	if err != nil {
		panic(err)
	}
	return result
}
