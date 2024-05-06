package system

import "errors"

// Common errors.
var (
	ErrTypeMismatch = errors.New("operation not defined between given types")
	// Date, Time, DateTime, and Quantity have special cases for equality and inequality logic
	// where an empty collection should be returned when their precisions/units are mismatched.
	ErrMismatchedPrecision = errors.New("mismatched precision")
	ErrMismatchedUnit      = errors.New("mismatched unit")
	ErrIntOverflow         = errors.New("operation resulted in integer overflow")
)

// Type names.
const (
	stringType   = "String"
	booleanType  = "Boolean"
	integerType  = "Integer"
	decimalType  = "Decimal"
	dateType     = "Date"
	dateTimeType = "DateTime"
	timeType     = "Time"
	quantityType = "Quantity"
	anyType      = "Any"
)
