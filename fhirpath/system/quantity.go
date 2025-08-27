package system

import (
	"fmt"
	"time"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/shopspring/decimal"
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

// Quantity type represents a decimal value along with a UCUM unit or
// a calendar duration keyword.
type Quantity struct {
	value Decimal
	unit  string
}

// ParseQuantity takes as input a number string and a unit string and constructs
// a Quantity object. Returns error if the input does not fit into a valid Quantity
func ParseQuantity(number string, unit string) (Quantity, error) {
	d, err := ParseDecimal(number)
	if err != nil {
		return Quantity{}, err
	}
	q, err := newQuantity(d, unit)
	if err != nil {
		return Quantity{}, err
	}
	return q, nil
}

// MustParseQuantity constructs a Quantity object from a number string and unit string.
// Panics if the inputs do not fit into a valid Quantity.
func MustParseQuantity(number string, unit string) Quantity {
	q, err := ParseQuantity(number, unit)
	if err != nil {
		panic(err)
	}
	return q
}

// newQuantity constructs a system Quantity type, given a decimal
// value and a UCUM unit identifier.
func newQuantity(value Decimal, unit string) (Quantity, error) {
	// check isValidUnit(unit) *To be implemented
	return Quantity{value, unit}, nil
}

// TryEqual returns a bool representing whether or not the
// value represented by q is equal to the value of q2.
// The comparison is not symmetric and may not return a value, represented by
// the second boolean being set to false.
func (q Quantity) TryEqual(input Any) (bool, bool) {
	val, ok := input.(Quantity)
	if !ok {
		return false, true
	}
	if q.unit != val.unit {
		return false, false
	}
	return q.value.Equal(val.value), true
}

// Less returns true if q is less than input.(Quantity). If the units
// are mismatched, returns an error. If input is not a Quantity, returns
// an error.
func (q Quantity) Less(input Any) (Boolean, error) {
	val, ok := input.(Quantity)
	if !ok {
		return false, fmt.Errorf("%w: %T, %T", ErrTypeMismatch, q, input)
	}
	if q.unit != val.unit {
		return false, ErrMismatchedUnit
	}
	return q.value.Less(val.value)
}

// Add returns q + input. Returns an error if the units are mismatched.
func (q Quantity) Add(input Quantity) (Quantity, error) {
	if q.unit != input.unit {
		return Quantity{}, ErrMismatchedUnit
	}
	value := Decimal(decimal.Decimal(q.value).Add(decimal.Decimal(input.value)))
	return Quantity{value, q.unit}, nil
}

// Sub returns q - input. Returns an error if the units are mismatched.
func (q Quantity) Sub(input Quantity) (Quantity, error) {
	if q.unit != input.unit {
		return Quantity{}, ErrMismatchedUnit
	}
	value := Decimal(decimal.Decimal(q.value).Sub(decimal.Decimal(input.value)))
	return Quantity{value, q.unit}, nil
}

// Name returns the type name.
func (q Quantity) Name() string {
	return quantityType
}

// String returns the quantity value formatted as a string.
func (q Quantity) String() string {
	if q.unit != "" {
		return fmt.Sprintf("%s %s", decimal.Decimal(q.value).String(), q.unit)
	}
	return fmt.Sprintf("%s", decimal.Decimal(q.value).String())
}

// ToProtoQuantity returns a proto Quantity based on a system Quantity.
func (q Quantity) ToProtoQuantity() *dtpb.Quantity {
	res := &dtpb.Quantity{
		Value: q.value.ToProtoDecimal(),
	}

	if q.unit != "" {
		res.Unit = fhir.String(q.unit)
	}

	return res
}

// TODO: According to our initial investigations, Equal should return empty when value is not set.
// Equal method to override cmp.Equal.
func (q Quantity) Equal(q2 Quantity) bool {
	return q.value.Equal(q2.value) && q.unit == q2.unit
}

// Negate returns the quantity multiplied by -1.
func (q Quantity) Negate() Quantity {
	negative := Decimal(decimal.NewFromInt(-1))
	return Quantity{q.value.Mul(negative), q.unit}
}

// timeDuration returns the time.Duration represented by
// a time-valued Quantity. Returns an error if the Quantity
// doesn't represent a valid time duration.
func (q Quantity) timeDuration() (time.Duration, error) {
	value := decimal.Decimal(q.value).IntPart()

	var duration time.Duration
	switch q.unit {
	case "hour", "hours":
		duration = time.Hour * time.Duration(value)
	case "minute", "minutes":
		duration = time.Minute * time.Duration(value)
	case "second", "seconds":
		milliseconds := decimal.Decimal(q.value).Round(3).Shift(3).IntPart() // Keep decimal precision below seconds
		duration = time.Millisecond * time.Duration(milliseconds)
	case "millisecond":
		duration = time.Millisecond * time.Duration(value)
	default:
		return time.Duration(0), fmt.Errorf("%w: not a time-valued unit", ErrMismatchedUnit)
	}
	return duration, nil
}

// Converts valid time based quantities to a number of years,
// by rounding down.
func (q Quantity) toYears() (int, error) {
	value := int(decimal.Decimal(q.value).IntPart())

	switch q.unit {
	case "year", "years":
		return value, nil
	case "month", "months":
		return value / 12, nil
	case "week", "weeks":
		value = value * 7
		fallthrough
	case "day", "days":
		return value / 365, nil
	case "hour", "hours":
		return value / (365 * 24), nil
	case "minute", "minutes":
		return value / (365 * 24 * 60), nil
	case "second", "seconds":
		return value / (365 * 24 * 60 * 60), nil
	case "millisecond", "milliseconds":
		return value / (365 * 24 * 60 * 60) / 1000, nil
	default:
		return 0, fmt.Errorf("%w: not a time-valued unit", ErrMismatchedUnit)
	}
}

// Converts a valid time based quantity to a number of months,
// by rounding down.
func (q Quantity) toMonths() (int, error) {
	value := int(decimal.Decimal(q.value).IntPart())

	switch q.unit {
	case "year", "years":
		return value * 12, nil
	case "month", "months":
		return value, nil
	case "week", "weeks":
		value = value * 7
		fallthrough
	case "day", "days":
		return value / 30, nil
	case "hour", "hours":
		return value / (30 * 24), nil
	case "minute", "minutes":
		return value / (30 * 24 * 60), nil
	case "second", "seconds":
		return value / (30 * 24 * 60 * 60), nil
	case "millisecond", "milliseconds":
		return value / (30 * 24 * 60 * 60) / 1000, nil
	default:
		return 0, fmt.Errorf("%w: not a time-valued unit", ErrMismatchedUnit)
	}
}
