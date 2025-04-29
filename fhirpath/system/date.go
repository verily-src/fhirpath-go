package system

import (
	"fmt"
	"strings"
	"time"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/shopspring/decimal"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirconv"
)

// Date represents Date values in the range 0001-01-01 to
// 9999-12-31 with a 1 day step size. It uses YYYY-MM-DD with
// optional month and day parts.
type Date struct {
	date time.Time
	l    layout
}

// ParseDate takes as input a string and returns a Date
// if the input is a valid FHIRPath Date string, otherwise
// returns an error.
func ParseDate(value string) (Date, error) {
	dateLayouts := []string{
		dayLayout,
		monthLayout,
		yearLayout,
	}

	var t time.Time
	var err error
	value = strings.TrimPrefix(value, "@")
	for _, l := range dateLayouts {
		if t, err = time.Parse(l, value); err == nil {
			return Date{t, layout(l)}, nil
		}
	}
	return Date{}, fmt.Errorf("unable to parse date '%s': %w", value, err)
}

// MustParseDate takes as input a string and returns a Date if
// the input is valid. Panics if the string can't be parsed.
func MustParseDate(value string) Date {
	date, err := ParseDate(value)
	if err != nil {
		panic(err)
	}
	return date
}

// DateFromProto takes a proto Date as input and returns a system Date.
func DateFromProto(proto *dtpb.Date) (Date, error) {
	t, err := fhirconv.DateToTime(proto)
	if err != nil {
		return Date{}, err
	}
	var l layout
	switch proto.Precision {
	case dtpb.Date_DAY:
		l = dayLayout
	case dtpb.Date_MONTH:
		l = monthLayout
	case dtpb.Date_YEAR:
		l = yearLayout
	}
	return Date{t, l}, nil
}

// ToProtoDate returns a proto Date based on a system Date.
func (d Date) ToProtoDate() *dtpb.Date {
	date := fhir.Date(d.date)
	var p dtpb.Date_Precision
	switch d.l {
	case dayLayout:
		p = dtpb.Date_DAY
	case monthLayout:
		p = dtpb.Date_MONTH
	case yearLayout:
		p = dtpb.Date_YEAR
	}
	date.Precision = p
	return date
}

// String formats the time as a date string.
func (d Date) String() string {
	return d.date.Format(string(d.l))
}

// TryEqual returns a bool representing whether or not
// the value of d is equal to input.(Date).
// This function may not return a value, depending on the precision of
// the other value; represented by the second bool return value.
func (d Date) TryEqual(input Any) (bool, bool) {
	val, ok := input.(Date)
	if !ok {
		return false, true
	}
	if d.l == val.l {
		return d.date.Equal(val.date), true
	}

	dComponents := d.getComponents()
	valComponents := val.getComponents()

	minPrecision := min(int(dateMap[d.l]), int(dateMap[val.l]))

	for i := 0; i <= minPrecision; i++ {
		if dComponents[i] == valComponents[i] {
			continue
		}
		return false, true
	}
	return false, false
}

// Less returns true if the value of d is less than input.(Date).
// Compares component by component, and returns an error if there is a
// precision mismatch. If input is not a Date, returns an error.
func (d Date) Less(input Any) (Boolean, error) {
	val, ok := input.(Date)
	if !ok {
		return false, fmt.Errorf("%w: %T, %T", ErrTypeMismatch, d, input)
	}
	if d.l == val.l {
		return Boolean(d.date.Before(val.date)), nil
	}

	dComponents := d.getComponents()
	valComponents := val.getComponents()

	minPrecision := min(int(dateMap[d.l]), int(dateMap[val.l]))

	for i := 0; i <= minPrecision; i++ {
		if dComponents[i] == valComponents[i] {
			continue
		}
		return dComponents[i] < valComponents[i], nil
	}
	return false, ErrMismatchedPrecision
}

// Add returns the result of d + input. Returns an
// error if it is not a valid time-valued quantity.
func (d Date) Add(input Quantity) (Date, error) {
	var result time.Time
	value := int(decimal.Decimal(input.value).IntPart())
	switch input.unit {
	case "year", "years":
		result = addYear(d.date, value)
	case "month", "months":
		result = addMonth(d.date, value)
	case "week", "weeks":
		value = 7 * value
		result = d.date.AddDate(0, 0, value)
	case "day", "days":
		result = d.date.AddDate(0, 0, value)
	default:
		return Date{}, fmt.Errorf("%w: can't add to date", ErrMismatchedUnit)
	}

	// Reformat to truncate date to initial precision. This causes the addition result
	// to round down to the highest precision value.
	result, err := time.Parse(string(d.l), result.Format(string(d.l)))
	if err != nil {
		return Date{}, err
	}
	return Date{result, d.l}, nil
}

// Sub returns the result of d - input. Returns an error if the
// input does not represent a valid time-valued quantity.
func (d Date) Sub(input Quantity) (Date, error) {
	// Handle partial dates by rounding quantity to appropriate precision.
	// Subtraction is not symmetric with addition, so the solution of truncation
	// as done in Add, cannot be applied here.
	if d.l == yearLayout {
		years, err := input.toYears()
		if err != nil {
			return Date{}, err
		}
		return Date{d.date.AddDate(-years, 0, 0), d.l}, nil
	}
	if d.l == monthLayout {
		months, err := input.toMonths()
		if err != nil {
			return Date{}, err
		}
		return Date{d.date.AddDate(0, -months, 0), d.l}, nil
	}

	// subtract appropriate position of date, for non-partial dates.
	var result time.Time
	value := -int(decimal.Decimal(input.value).IntPart())
	switch input.unit {
	case "year", "years":
		result = addYear(d.date, value)
	case "month", "months":
		result = addMonth(d.date, value)
	case "week", "weeks":
		value = 7 * value
		result = d.date.AddDate(0, 0, value)
	case "day", "days":
		result = d.date.AddDate(0, 0, value)
	default:
		return Date{}, fmt.Errorf("%w: can't add to date", ErrMismatchedUnit)
	}

	return Date{result, d.l}, nil
}

// Name returns the type name.
func (d Date) Name() string {
	return dateType
}

// Equal method to override cmp.Equal.
func (d Date) Equal(d2 Date) bool {
	return d.date.Format(string(d.l)) == d2.date.Format(string(d2.l))
}

func (d Date) getComponents() []int {
	return []int{
		d.date.Year(),
		int(d.date.Month()),
		d.date.Day(),
	}
}

// ToDateTime converts a Date to a DateTime, preserving the precision of the original Date.
func (d Date) ToDateTime() DateTime {
	var dateToDateTime = map[layout]layout{
		dayLayout:   dtDayLayout,
		monthLayout: dtMonthLayout,
		yearLayout:  dtYearLayout,
	}

	return DateTime{d.date, dateToDateTime[d.l]}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// addMonth adds m months to the given time, without
// normalizing to the next month.
// eg. 2020-01-30 + 1 month = 2020-02-28.
func addMonth(t time.Time, m int) time.Time {
	added := t.AddDate(0, m, 0)
	if day := added.Day(); day != t.Day() {
		return added.AddDate(0, 0, -day)
	}
	return added
}

// addYear adds y years to the given time, without
// normalizing to the next month in the special case of leap years.
// eg. 2020-02-29 + 1 year = 2020-02-28.
func addYear(t time.Time, y int) time.Time {
	added := t.AddDate(y, 0, 0)
	if day := added.Day(); day != t.Day() {
		return added.AddDate(0, 0, -day)
	}
	return added
}
