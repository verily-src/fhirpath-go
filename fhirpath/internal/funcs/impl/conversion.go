package impl

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

// DefaultQuantityUnit is defined by the following FHIRPath rules:
// the item is an Integer, or Decimal, where the resulting quantity will have the default unit ('1')
// the item is a Boolean, where true results in the quantity 1.0 '1', and false results in the quantity 0.0 '1'
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#toquantityunit-string-quantity
const DefaultQuantityUnit = "1"

// Based on the FHIRPath Quantity string validation regexp defined here:
// https://hl7.org/fhirpath/N1/#convertstoquantityunit-string-boolean
const fhirQuantityRegexp = `^(?P<value>(\+|-)?\d+(\.\d+)?)\s*('(?P<unit>[^']+)'|(?P<time>[a-zA-Z]+))?$`

var regex = regexp.MustCompile(fhirQuantityRegexp)

// ConvertsToBoolean checks if the input can be converted to a Boolean
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#convertstoboolean-boolean
func ConvertsToBoolean(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	if !input.IsSingleton() {
		return nil, errors.New("invalid input, is not a singleton")
	}
	// Argument validation
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Conversion validation
	result, err := ToBoolean(ctx, input, args...)
	if result.IsEmpty() || err != nil {
		return system.Collection{system.Boolean(false)}, nil
	}
	return system.Collection{system.Boolean(true)}, nil
}

// ConvertsToDate checks if the input can be converted to a Date
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#convertstodate-boolean
func ConvertsToDate(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	if !input.IsSingleton() {
		return nil, errors.New("invalid input, is not a singleton")
	}
	// Argument validation
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Conversion validation
	result, err := ToDate(ctx, input, args...)
	if result.IsEmpty() || err != nil {
		return system.Collection{system.Boolean(false)}, nil
	}
	return system.Collection{system.Boolean(true)}, nil
}

// ConvertsToDateTime checks if the input can be converted to a Time
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#convertstodatetime-boolean
func ConvertsToDateTime(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	if !input.IsSingleton() {
		return nil, errors.New("invalid input, is not a singleton")
	}
	// Argument validation
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Conversion validation
	result, err := ToDateTime(ctx, input, args...)
	if result.IsEmpty() || err != nil {
		return system.Collection{system.Boolean(false)}, nil
	}
	return system.Collection{system.Boolean(true)}, nil
}

// ConvertsToDecimal checks if the input can be converted to a Decimal
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#convertstodecimal-boolean
func ConvertsToDecimal(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	if !input.IsSingleton() {
		return nil, errors.New("invalid input, is not a singleton")
	}
	// Argument validation
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Conversion validation
	result, err := ToDecimal(ctx, input, args...)
	if result.IsEmpty() || err != nil {
		return system.Collection{system.Boolean(false)}, nil
	}
	return system.Collection{system.Boolean(true)}, nil
}

// ConvertsToInteger checks if the input can be converted to an Integer
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#convertstointeger-boolean
func ConvertsToInteger(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	if !input.IsSingleton() {
		return nil, errors.New("invalid input, is not a singleton")
	}
	// Argument validation
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Conversion validation
	result, err := ToInteger(ctx, input, args...)
	if result.IsEmpty() || err != nil {
		return system.Collection{system.Boolean(false)}, nil
	}
	return system.Collection{system.Boolean(true)}, nil
}

// ConvertsToQuantity checks if the input can be converted to a Quantity
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#convertstoquantityunit-string-boolean
func ConvertsToQuantity(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	if !input.IsSingleton() {
		return nil, errors.New("invalid input, is not a singleton")
	}
	// Argument validation
	if len(args) > 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1 or 0", ErrWrongArity, len(args))
	}
	// Conversion validation
	result, err := ToQuantity(ctx, input, args...)
	if result.IsEmpty() || err != nil {
		return system.Collection{system.Boolean(false)}, nil
	}
	return system.Collection{system.Boolean(true)}, nil
}

// ConvertsToString checks if the input can be converted to a String
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#convertstostring-string
func ConvertsToString(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	if !input.IsSingleton() {
		return nil, errors.New("invalid input, is not a singleton")
	}
	// Argument validation
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Conversion validation
	result, err := ToString(ctx, input, args...)
	if result.IsEmpty() || err != nil {
		return system.Collection{system.Boolean(false)}, nil
	}
	if boolean, _ := result.ToBool(); boolean == false {
		return system.Collection{system.Boolean(false)}, nil
	}
	return system.Collection{system.Boolean(true)}, nil
}

// ConvertsToTime checks if the input can be converted to a Time
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#convertstotime-boolean
func ConvertsToTime(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	if !input.IsSingleton() {
		return nil, errors.New("invalid input, is not a singleton")
	}
	// Argument validation
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Conversion validation
	result, err := ToTime(ctx, input, args...)
	if result.IsEmpty() || err != nil {
		return system.Collection{system.Boolean(false)}, nil
	}
	return system.Collection{system.Boolean(true)}, nil
}

// ToBoolean converts the input to a Boolean
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#toboolean-boolean
func ToBoolean(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	if !input.IsSingleton() {
		return nil, errors.New("invalid input, is not a singleton")
	}
	// Argument validation
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Input reading
	value, err := system.From(input[0])
	if err != nil {
		return nil, err
	}
	// Input conversion
	switch value := value.(type) {
	case system.Decimal:
		result, err := system.ParseBoolean(value.String())
		if err != nil {
			return system.Collection{}, nil
		}
		return system.Collection{result}, nil
	case system.Integer, system.String:
		str := fmt.Sprintf("%v", value)
		result, err := system.ParseBoolean(str)
		if err != nil {
			return system.Collection{}, nil
		}
		return system.Collection{result}, nil
	case system.Boolean:
		return system.Collection{value}, nil
	}
	return system.Collection{}, nil
}

// ToDate converts the input to a Date
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#todate-date
func ToDate(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	if !input.IsSingleton() {
		return nil, errors.New("invalid input, is not a singleton")
	}
	// Argument validation
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Input reading
	value, err := system.From(input[0])
	if err != nil {
		return system.Collection{}, nil
	}
	// Input conversion
	switch value := value.(type) {
	case system.Date:
		return system.Collection{value}, nil
	case system.DateTime:
		dt := value.String()
		result := system.MustParseDate(dt[:10])
		return system.Collection{result}, nil
	case system.String:
		result, err := system.ParseDate(string(value))
		if err != nil {
			return system.Collection{}, nil
		}
		return system.Collection{result}, nil
	}
	return system.Collection{}, nil
}

// ToDateTime converts the input to a Date
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#todatetime-datetime
func ToDateTime(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validations
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	if !input.IsSingleton() {
		return nil, errors.New("invalid input, is not a singleton")
	}
	// Argument validations
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Reading input
	value, err := system.From(input[0])
	if err != nil {
		return system.Collection{}, nil
	}
	// Input conversion
	switch value := value.(type) {
	case system.Date:
		return system.Collection{value.ToDateTime()}, nil
	case system.DateTime:
		return system.Collection{value}, nil
	case system.String:
		result, err := system.ParseDateTime(string(value))
		if err != nil {
			return system.Collection{}, nil
		}
		return system.Collection{result}, nil
	}
	return system.Collection{}, nil
}

// ToDecimal converts the input to a Decimal
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#todecimal-decimal
func ToDecimal(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	if !input.IsSingleton() {
		return nil, errors.New("invalid input, is not a singleton")
	}
	// Argument validation
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Input reading
	value, err := system.From(input[0])
	if err != nil {
		return nil, err
	}
	// Input conversion
	switch value.(type) {
	case system.Decimal:
		return system.Collection{value}, nil
	case system.Integer:
		str := fmt.Sprintf("%v", value)
		result, err := system.ParseDecimal(str)
		if err != nil {
			return system.Collection{}, nil
		}
		return system.Collection{result}, nil
	case system.String:
		str := fmt.Sprintf("%s", value)
		result, err := system.ParseDecimal(str)
		if err != nil {
			return system.Collection{}, nil
		}
		return system.Collection{result}, nil
	case system.Boolean:
		if value.(system.Boolean) {
			return system.Collection{system.MustParseDecimal("1.0")}, nil
		}
		return system.Collection{system.MustParseDecimal("0.0")}, nil
	}
	return system.Collection{}, nil
}

// ToInteger converts the input to an Integer
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#tointeger-integer
func ToInteger(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	if !input.IsSingleton() {
		return nil, errors.New("invalid input, is not a singleton")
	}
	// Argument validation
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Input reading
	value, err := system.From(input[0])
	if err != nil {
		return nil, err
	}
	// Input conversion
	switch value.(type) {
	case system.Integer:
		return system.Collection{value}, nil
	case system.String:
		str := fmt.Sprintf("%s", value)
		result, err := system.ParseInteger(str)
		if err != nil {
			return nil, err
		}
		return system.Collection{result}, nil
	case system.Boolean:
		if value.(system.Boolean) {
			return system.Collection{system.Integer(1)}, nil
		}
		return system.Collection{system.Integer(0)}, nil
	}
	return system.Collection{}, nil
}

// ToQuantity converts the input to a Quantity
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#toquantityunit-string-quantity
func ToQuantity(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	if !input.IsSingleton() {
		return nil, errors.New("invalid input, is not a singleton")
	}
	// Argument validation
	if len(args) > 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1 or 0", ErrWrongArity, len(args))
	}
	argStr := ""
	if len(args) == 1 {
		argValue, err := args[0].Evaluate(ctx, input)
		if err != nil {
			return nil, err
		}
		argStr, err = argValue.ToString()
		if err != nil {
			return nil, err
		}
	}
	// Input reading
	value, err := system.From(input[0])
	if err != nil {
		return nil, err
	}
	// Input conversion
	switch value := value.(type) {
	case system.Integer:
		matches := regex.FindStringSubmatch(fmt.Sprintf("%v %v", value, argStr))
		if matches == nil {
			return system.Collection{}, nil
		}
		if argStr != "" {
			unit := regex.SubexpIndex("unit")
			t := regex.SubexpIndex("time")
			if matches[unit] != "" {
				result := system.MustParseQuantity(fmt.Sprintf("%v", value), matches[unit])
				return system.Collection{result}, nil
			}
			if matches[t] != "" {
				result := system.MustParseQuantity(fmt.Sprintf("%v", value), matches[t])
				return system.Collection{result}, nil
			}
		}
		result := system.MustParseQuantity(fmt.Sprintf("%v", value), DefaultQuantityUnit)
		return system.Collection{result}, nil
	case system.Decimal:
		str := value.String()
		result, err := system.ParseQuantity(string(str), DefaultQuantityUnit)
		if err != nil {
			return nil, err
		}
		return system.Collection{result}, nil
	case system.Quantity:
		return system.Collection{value}, nil
	case system.String:
		matches := regex.FindStringSubmatch(string(value))
		if matches == nil {
			return system.Collection{}, nil
		}
		if argStr != "" {
			if !isValidUnitConversion(argStr) {
				return nil, fmt.Errorf("invalid unit of time: %v", input)
			}
			conversion, err := convertDuration(string(value), argStr)
			if err != nil {
				return nil, err
			}
			res := strings.SplitN(conversion, " ", 2)
			result := system.MustParseQuantity(res[0], res[1])
			return system.Collection{result}, nil
		}
		res := strings.SplitN(string(value), " ", 2)
		unit := strings.Trim(res[1], "'")
		result := system.MustParseQuantity(res[0], unit)
		return system.Collection{result}, nil
	case system.Boolean:
		if value {
			result := system.MustParseQuantity("1.0", DefaultQuantityUnit)
			return system.Collection{result}, nil
		}
		result := system.MustParseQuantity("0.0", DefaultQuantityUnit)
		return system.Collection{result}, nil
	}
	return system.Collection{}, nil
}

// ToString converts the input to a String
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#tostring-string
func ToString(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	if !input.IsSingleton() {
		return nil, errors.New("invalid input, is not a singleton")
	}
	// Argument validation
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Input reading
	value, err := system.From(input[0])
	if err != nil {
		return system.Collection{system.Boolean(false)}, nil
	}
	// Input conversion
	switch value := value.(type) {
	case system.String:
		return system.Collection{value}, nil
	case system.Integer:
		return system.Collection{system.String(fmt.Sprintf("%v", value))}, nil
	case system.Decimal:
		return system.Collection{system.String(value.String())}, nil
	case system.Quantity:
		return system.Collection{system.String(value.String())}, nil
	case system.Date:
		return system.Collection{system.String(value.String())}, nil
	case system.Time:
		return system.Collection{system.String(value.String())}, nil
	case system.DateTime:
		return system.Collection{system.String(value.String())}, nil
	case system.Boolean:
		if value {
			return system.Collection{system.String("true")}, nil
		}
		return system.Collection{system.String("false")}, nil
	}
	return system.Collection{system.Boolean(false)}, nil
}

// ToTime converts the input to a Time
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#totime-time
func ToTime(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	if !input.IsSingleton() {
		return nil, errors.New("invalid input, is not a singleton")
	}
	// Argument validation
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	// Input reading
	value, err := system.From(input[0])
	if err != nil {
		return system.Collection{}, nil
	}
	// Input conversion
	switch value := value.(type) {
	case system.Time:
		return system.Collection{value}, nil
	case system.String:
		result, err := system.ParseTime(fmt.Sprintf("%v", value))
		if err != nil {
			return system.Collection{}, nil
		}
		return system.Collection{result}, nil
	}
	return system.Collection{}, nil
}

// Iif function is an immediate if/conditional operator
// If the criterion expression evaluates to true, the result is the evaluation of the true expression.
// If the criterion expression evaluates to false or is empty, the result is the evaluation of the false expression or
// empty if the false expression is not provided.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#iifcriterion-expression-true-result-collection-otherwise-result-collection-collection
func Iif(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Note for readability: criterion = args[0], true-result = args[1], otherwise-result = args[2]

	// Argument validation
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 2 or 3", ErrWrongArity, len(args))
	}

	// Evaluate the criterion
	criterionResult, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	}

	// Convert the criterion result to a boolean
	criterionBool, err := criterionResult.ToBool()
	if err != nil {
		return nil, err
	}

	// Check if the criterion result is empty or false
	if !criterionBool {
		if len(args) != 3 {
			return system.Collection{}, nil
		}

		return args[2].Evaluate(ctx, input)
	}

	// Return the true-result collection
	return args[1].Evaluate(ctx, input)
}

func isValidUnitConversion(outputFormat string) bool {
	validFormats := map[string]bool{
		"years":   true,
		"months":  true,
		"days":    true,
		"hours":   true,
		"minutes": true,
		"seconds": true,
	}

	return validFormats[outputFormat]
}

func convertDuration(input string, outputFormat string) (string, error) {
	duration, err := parseHumanDuration(input)
	if err != nil {
		return "", err
	}

	var convertedValue float64
	switch outputFormat {
	case "years":
		convertedValue = duration.Hours() / (24 * 365)
	case "months":
		convertedValue = duration.Hours() / (24 * 30)
	case "days":
		convertedValue = duration.Hours() / 24
	case "hours":
		convertedValue = duration.Hours()
	case "minutes":
		convertedValue = duration.Minutes()
	case "seconds":
		convertedValue = duration.Seconds()
	}

	if outputFormat == "years" {
		convertedValue = math.Ceil(convertedValue / 12.0)
	}

	return fmt.Sprintf("%.0f %s", convertedValue, outputFormat), nil
}

func parseHumanDuration(input string) (time.Duration, error) {
	re := regexp.MustCompile(`(\d+)\s*(\w+)`)
	matches := re.FindAllStringSubmatch(input, -1)
	totalSeconds := int64(0)

	for _, match := range matches {
		value, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return 0, err
		}

		unit := strings.ToLower(match[2])
		switch unit {
		case "second", "seconds":
			totalSeconds += value
		case "minute", "minutes":
			totalSeconds += value * 60
		case "hour", "hours":
			totalSeconds += value * 3600
		case "day", "days":
			totalSeconds += value * 86400
		case "month", "months":
			totalSeconds += value * 30 * 86400 // Assuming one month is 30 days
		case "year", "years":
			totalSeconds += value * 365 * 86400 // Assuming one year is 365 days
		default:
			return 0, fmt.Errorf("invalid unit: %s", unit)
		}
	}

	return time.Duration(totalSeconds) * time.Second, nil
}
