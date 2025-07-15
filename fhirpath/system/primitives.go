package system

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/shopspring/decimal"
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

// Boolean is a representation for a true or false value.
type Boolean bool

// ParseBoolean parses a boolean string and returns a Boolean value
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#toboolean-boolean
func ParseBoolean(input string) (Boolean, error) {
	switch strings.ToLower(input) {
	case "true", "t", "yes", "y", "1", "1.0":
		return true, nil
	case "false", "f", "no", "n", "0", "0.0":
		return false, nil
	default:
		return false, fmt.Errorf("%v can't be parsed to boolean", input)
	}
}

// Equal returns true if the input value is a System Boolean,
// and contains the same value.
func (b Boolean) Equal(input Any) bool {
	val, ok := input.(Boolean)
	if !ok {
		return false
	}
	return b == val
}

// Less returns error for Boolean comparison.
func (b Boolean) Less(input Any) (Boolean, error) {
	return false, fmt.Errorf("%w: %T, %T", ErrTypeMismatch, b, input)
}

// Name returns the type name.
func (b Boolean) Name() string {
	return booleanType
}

// String represents string values.
type String string

// ParseString parses the input string and replaces FHIRPath
// escape sequences with their Go-equivalent escape characters.
func ParseString(input string) (String, error) {
	escSequences := []string{
		"\\'", "'",
		"\\\"", "\"",
		"\\`", "`",
		"\\r", "\r",
		"\\t", "\t",
		"\\n", "\n",
		"\\f", "\f",
		"\\\\", "\\",
		"\\", "",
		// TODO: Add support for unicode escape sequences.
	}
	input = strings.TrimPrefix(input, "'")
	input = strings.TrimSuffix(input, "'")
	replacer := strings.NewReplacer(escSequences...)
	escapedString := replacer.Replace(input)
	return String(escapedString), nil
}

// Equal returns true if the input value is a System String,
// and contains the same string value.
func (s String) Equal(input Any) bool {
	val, ok := input.(String)
	if !ok {
		return false
	}
	return s == val
}

// Name returns the type name.
func (s String) Name() string {
	return stringType
}

// Less returns true if s is less than input.(String), by lexicographic
// comparison. If input is not a String, returns an error.
func (s String) Less(input Any) (Boolean, error) {
	val, ok := input.(String)
	if !ok {
		return false, fmt.Errorf("%w: %T, %T,", ErrTypeMismatch, s, input)
	}
	return string(s) < string(val), nil
}

// Add concatenates the input String to the right of s.
func (s String) Add(input String) String {
	return s + input
}

// Integer represents integer values in the range 0 to 2^31 - 1.
// Negative integers in FHIRPath are denoted with the
type Integer int32

// ParseInteger parses a string into an int32 value, and returns
// an error if the input does not represent a valid int32.
func ParseInteger(value string) (Integer, error) {
	i, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return Integer(0), err
	}
	return Integer(i), nil
}

// Equal returns true if the input value is a System Integer, and
// contains the same int32 value.
func (i Integer) Equal(input Any) bool {
	val, ok := input.(Integer)
	if !ok {
		return false
	}
	return i == val
}

// Name returns the type name.
func (i Integer) Name() string {
	return integerType
}

// Less returns true if i is less than input.(Integer).
// If input is not an Integer, returns an error.
func (i Integer) Less(input Any) (Boolean, error) {
	val, ok := input.(Integer)
	if !ok {
		return false, fmt.Errorf("%w: %T, %T,", ErrTypeMismatch, i, input)
	}
	return i < val, nil
}

// Add adds i to the input Integer. Returns an error
// if the result overflows.
func (i Integer) Add(input Integer) (Integer, error) {
	result := i + input
	// If i is incremented (result > i), input must be positive.
	// Similarly, if i is decremented (result <= i), input must be 0 or negative.
	// Otherwise, an overflow must have occured.
	if (result > i) == (input > 0) {
		return result, nil
	}
	return 0, ErrIntOverflow
}

// Sub subtracts the input Integer from i. Returns an error
// if the result overflows.
func (i Integer) Sub(input Integer) (Integer, error) {
	result := i - input
	// If i is decremented (result < i), input must be positive.
	// Similarly, if i is incremented (result >= i), input must be 0 or negative.
	// Otherwise, an overflow must have occured.
	if (result < i) == (input > 0) {
		return result, nil
	}
	return 0, ErrIntOverflow
}

// Mul multiplies the two integers together. Returns an
// error if the result overflows.
func (i Integer) Mul(input Integer) (Integer, error) {
	if i == 0 || input == 0 {
		return 0, nil
	}
	result := i * input
	// Check if the sign of the result aligns with what it should be ( -ve == +ve * -ve )
	if (result < 0) == ((i < 0) != (input < 0)) && (result/input) == i {
		return result, nil
	}
	return 0, ErrIntOverflow
}

// Div divides i by input. Returns a Decimal.
func (i Integer) Div(input Integer) Decimal {
	lhs, rhs := decimal.NewFromInt32(int32(i)), decimal.NewFromInt32(int32(input))
	return Decimal(lhs.Div(rhs))
}

// FloorDiv divides i by input and rounds down.
func (i Integer) FloorDiv(input Integer) Integer {
	return i / input
}

// Mod returns i % integer.
func (i Integer) Mod(input Integer) Integer {
	return i % input
}

// ToProtoInteger returns the proto representation of the system integer.
func (i Integer) ToProtoInteger() *dtpb.Integer {
	return fhir.Integer(int32(i))
}

// Decimal represents fixed-point decimals. Must use
// utilities provided by "github.com/shopspring/decimal" to
// perform arithmetic.
type Decimal decimal.Decimal

// ParseDecimal parses a string representing a decimal, and
// returns an error if the input is invalid.
func ParseDecimal(value string) (Decimal, error) {
	d, err := decimal.NewFromString(value)
	if err != nil {
		return Decimal(decimal.Zero), err
	}
	return Decimal(d), nil
}

// MustParseDecimal converts a string into a Decimal type.
// If the string is not parseable it will throw a panic().
func MustParseDecimal(str string) Decimal {
	dec, err := ParseDecimal(str)
	if err != nil {
		panic(err)
	}
	return dec
}

// Equal returns a bool representing whether or not
// the two Decimals being compared are equal. Uses the API
// provided by the decimal library.
func (d Decimal) Equal(input Any) bool {
	val, ok := input.(Decimal)
	if !ok {
		return false
	}
	return decimal.Decimal(d).Equal(decimal.Decimal(val))
}

// Name returns the type name.
func (d Decimal) Name() string {
	return decimalType
}

// Less returns true if d is less than input.(Decimal). If input
// is not a Decimal, returns an error.
func (d Decimal) Less(input Any) (Boolean, error) {
	val, ok := input.(Decimal)
	if !ok {
		return false, fmt.Errorf("%w, %T, %T", ErrTypeMismatch, d, input)
	}
	return Boolean(decimal.Decimal(d).LessThan(decimal.Decimal(val))), nil
}

// Add adds d to the input Decimal
func (d Decimal) Add(input Decimal) Decimal {
	return Decimal(decimal.Decimal(d).Add(decimal.Decimal(input)))
}

// Sub subtracts the input Decimal from d.
func (d Decimal) Sub(input Decimal) Decimal {
	return Decimal(decimal.Decimal(d).Sub(decimal.Decimal(input)))
}

// String returns the string value from d.
func (d Decimal) String() string {
	return decimal.Decimal(d).String()
}

// Mul multiples the two decimal values together.
func (d Decimal) Mul(input Decimal) Decimal {
	return Decimal(decimal.Decimal(d).Mul(decimal.Decimal(input)))
}

// Div divides d by input.
func (d Decimal) Div(input Decimal) Decimal {
	return Decimal(decimal.Decimal(d).Div(decimal.Decimal(input)))
}

// FloorDiv divides d by input and rounds down.
func (d Decimal) FloorDiv(input Decimal) (Integer, error) {
	result := decimal.Decimal(d).Div(decimal.Decimal(input)).IntPart()
	if (result < math.MinInt32) || (result > math.MaxInt32) {
		return 0, ErrIntOverflow
	}
	return Integer(int32(result)), nil
}

// Mod computes d % input.
func (d Decimal) Mod(input Decimal) Decimal {
	return Decimal(decimal.Decimal(d).Mod(decimal.Decimal(input)))
}

// ToProtoDecimal returns the proto Decimal representation of decimal.
func (d Decimal) ToProtoDecimal() *dtpb.Decimal {
	return fhir.Decimal(decimal.Decimal(d).InexactFloat64())
}

// Round rounds a Decimal at the provided precision.
func (d Decimal) Round(precision int32) Decimal {
	return Decimal(decimal.Decimal(d).Round(precision))
}
