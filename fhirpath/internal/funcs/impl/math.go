package impl

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

var (
	ErrInvalidInput = errors.New("invalid input")
)

// Abs returns the absolute value of the input.
// When taking the absolute value of a quantity, the unit is unchanged.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#abs-integer-decimal-quantity
func Abs(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validations
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	// Argument validations
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}

	switch input[0].(type) {
	case system.Integer:
		// Input type conversion to int32
		number, err := input.ToInt32()
		if err != nil {
			return nil, err
		}
		// Absolution number
		res := math.Abs(float64(number))
		return system.Collection{system.Integer(res)}, nil
	case system.Decimal:
		// Input type conversion to float64
		number, err := input.ToFloat64()
		if err != nil {
			return nil, err
		}
		// Absolution number
		res := math.Abs(number)
		result := decimal.NewFromFloat(res)
		return system.Collection{system.Decimal(result)}, nil
	case system.Quantity:
		quantity := strings.Split(input[0].(system.Quantity).String(), " ")
		// Input type conversion
		f, err := strconv.ParseFloat(quantity[0], 64)
		if err != nil {
			return nil, err
		}
		// Absolution number
		res := math.Abs(f)
		return system.Collection{system.MustParseQuantity(fmt.Sprintf("%f", res), quantity[1])}, nil
	}
	return nil, errors.New("input is not a number")
}

// Ceiling returns the first integer greater than or equal to the input.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#ceiling-integer
func Ceiling(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validations
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	// Argument validations
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Input type conversion to float64
	number, err := input.ToFloat64()
	if err != nil {
		return nil, err
	}
	// Ceiling number
	result := math.Ceil(number)
	return system.Collection{system.Integer(result)}, nil
}

// Exp returns e raised to the power of the input.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#exp-decimal
func Exp(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validations
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	// Argument validations
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Reading input
	number, err := input.ToFloat64()
	if err != nil {
		return nil, err
	}
	// Exp number
	res := math.Pow(math.E, number)
	result := system.MustParseDecimal(fmt.Sprintf("%v", res))
	return system.Collection{result}, nil
}

// Floor returns the first integer less than or equal to the input.
// FHIRPath docs here: https://hl7.org/fhirpath/n1/#floor-integer
func Floor(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validations
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	// Argument validations
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Input type conversion to float64
	number, err := input.ToFloat64()
	if err != nil {
		return nil, err
	}
	// Flooring number
	result := math.Floor(number)
	return system.Collection{system.Integer(result)}, nil
}

// Ln returns the natural logarithm of the input number.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#ln-decimal
func Ln(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validations
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	// Argument validations
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Input type conversion to float64
	number, err := input.ToFloat64()
	if err != nil {
		return nil, err
	}
	res := math.Log(number)
	// Validating NaN case
	if math.IsNaN(res) {
		return system.Collection{}, nil
	}
	// Type conversion to system.Decimal
	result := decimal.NewFromFloat(res)
	return system.Collection{system.Decimal(result)}, nil
}

// Log returns the logarithm base of the input number.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#logbase-decimal-decimal
func Log(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validations
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	// Input type conversion to float64
	number, err := input.ToFloat64()
	if err != nil {
		return nil, err
	}
	// Validating args
	if len(args) != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, len(args))
	}
	argValues, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	}
	base, err := argValues.ToFloat64()
	if err != nil {
		return nil, err
	}
	// Log number to base
	res := logToBase(number, base)
	// Validating NaN case
	if math.IsNaN(res) {
		return system.Collection{}, nil
	}
	// Type conversion to system.Decimal
	result := decimal.NewFromFloat(res)
	return system.Collection{system.Decimal(result)}, nil
}

// Power returns a number to the exponent power.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#powerexponent-integer-decimal-integer-decimal
func Power(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Validating input
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	// Validating args
	if len(args) != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, len(args))
	}
	argValues, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	}
	// Validating integers case
	_, ok := input[0].(system.Integer)
	_, ok2 := argValues[0].(system.Integer)
	if ok && ok2 {
		// Input type conversion to int32
		number, err := input.ToInt32()
		if err != nil {
			return nil, err
		}
		// Input type conversion to int32
		exp, err := argValues.ToInt32()
		if err != nil {
			return nil, err
		}
		// Powering ints
		res := powInt32(number, exp)
		return system.Collection{system.Integer(res)}, nil
	}
	// Input type conversion to float64
	number, err := input.ToFloat64()
	if err != nil {
		return nil, err
	}
	// Input type conversion to float64
	exp, err := argValues.ToFloat64()
	if err != nil {
		return nil, err
	}
	// Powering number
	res := math.Pow(number, exp)
	// Validating NaN case
	if math.IsNaN(res) {
		return system.Collection{}, nil
	}
	// Type conversion to system.Decimal
	result := decimal.NewFromFloat(res)
	return system.Collection{system.Decimal(result)}, nil
}

// Round rounds the decimal to the nearest whole number using a traditional round (i.e. 0.5 or higher will round to 1).
// If specified, the precision argument determines the decimal place at which the rounding will occur.
// If not specified, the rounding will default to 0 decimal places.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#roundprecision-integer-decimal
func Round(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Validating input
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	if !input.IsSingleton() {
		return nil, errors.New("invalid input, is not a singleton")
	}
	// Validating args
	if len(args) > 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0 or 1 arguments", ErrWrongArity, len(args))
	}
	precision := int32(0)
	if len(args) == 1 {
		argValues, err := args[0].Evaluate(ctx, input)
		if err != nil {
			return nil, err
		}
		// Arg type conversion to int32
		precision, err = argValues.ToInt32()
		if err != nil {
			return nil, err
		}
		if precision < 0 {
			return nil, errors.New("precision must be greater or equal than 0")
		}
	}
	value, err := system.From(input[0])
	if err != nil {
		return nil, err
	}
	// Rounding number
	switch value.(type) {
	case system.Decimal:
		res, _ := input[0].(system.Decimal)
		result := res.Round(precision)
		return system.Collection{result}, nil
	case system.Integer:
		number, err := input.ToInt32()
		if err != nil {
			return nil, err
		}
		res := system.MustParseDecimal(fmt.Sprintf("%d", number))
		result := res.Round(precision)
		return system.Collection{result}, nil
	}
	return nil, errors.New("input is not a number")
}

// Sqrt returns the square root of the input number as a Decimal.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#sqrt-decimal
func Sqrt(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validations
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	// Argument validations
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Input type conversion to float64
	number, err := input.ToFloat64()
	if err != nil {
		return nil, err
	}
	// Validate negative input
	if number < 0 {
		return nil, fmt.Errorf("%w: unable to sqrt negative value", ErrInvalidInput)
	}
	// Ceiling number
	value := math.Sqrt(number)
	result := decimal.NewFromFloat(value)
	return system.Collection{system.Decimal(result)}, nil
}

// Truncate returns the integer portion of the input.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#truncate-integer
func Truncate(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validations
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	// Argument validations
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Input type conversion to float64
	number, err := input.ToFloat64()
	if err != nil {
		return nil, err
	}
	// Ceiling number
	result := math.Trunc(number)
	return system.Collection{system.Integer(result)}, nil
}

func logToBase(number, base float64) float64 {
	if number <= 0 || base <= 1 {
		return math.NaN() // Return NaN for invalid inputs
	}

	return math.Log(number) / math.Log(base)
}

// powInt32 returns the powering of a number to a given exponential.
func powInt32(base, exp int32) int32 {
	if exp == 0 {
		return 1
	}
	if exp < 0 {
		return 0
	}

	result := base
	for i := int32(2); i <= exp; i++ {
		result *= base
	}
	return result
}
