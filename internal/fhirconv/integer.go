package fhirconv

import (
	"errors"
	"fmt"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/narrow"
	"golang.org/x/exp/constraints"
)

var (
	// ErrIntegerTruncated is an error raised when an integer truncation occurs
	// during integer conversion
	ErrIntegerTruncated = errors.New("integer truncation")
)

// integerType is a constraint for FHIR integer types.
type integerType interface {
	*dtpb.Integer | *dtpb.UnsignedInt | *dtpb.PositiveInt
}

// ToInt8 converts a FHIR Integer type into a Go native int8.
func ToInt8[From integerType](v From) (int8, error) {
	return ToInteger[int8](v)
}

// ToInt16 converts a FHIR Integer type into a Go native int16.
func ToInt16[From integerType](v From) (int16, error) {
	return ToInteger[int16](v)
}

// ToInt32 converts a FHIR Integer type into a Go native int32.
func ToInt32[From integerType](v From) (int32, error) {
	return ToInteger[int32](v)
}

// ToInt64 converts a FHIR Integer type into a Go native int64.
func ToInt64[From integerType](v From) (int64, error) {
	return ToInteger[int64](v)
}

// ToInt converts a FHIR Integer type into a Go native int.
func ToInt[From integerType](v From) (int, error) {
	return ToInteger[int](v)
}

// ToUint8 converts a FHIR Integer type into a Go native uint8.
func ToUint8[From integerType](v From) (uint8, error) {
	return ToInteger[uint8](v)
}

// ToUint16 converts a FHIR Integer type into a Go native uint16.
func ToUint16[From integerType](v From) (uint16, error) {
	return ToInteger[uint16](v)
}

// ToUint32 converts a FHIR Integer type into a Go native uint32.
func ToUint32[From integerType](v From) (uint32, error) {
	return ToInteger[uint32](v)
}

// ToUint64 converts a FHIR Integer type into a Go native uint64.
func ToUint64[From integerType](v From) (uint64, error) {
	return ToInteger[uint64](v)
}

// ToUint converts a FHIR Integer type into a Go native uint.
func ToUint[From integerType](v From) (uint, error) {
	return ToInteger[uint](v)
}

// ToInteger converts a FHIR Integer type into a Go native integer type.
//
// If the value of the integer does not fit into the receiver integer type,
// this function will return an ErrIntegerTruncated.
func ToInteger[To constraints.Integer, From integerType](v From) (To, error) {
	var result To
	if val, ok := any(v).(interface{ GetValue() uint32 }); ok {
		if result, ok := narrow.ToInteger[To](uint64(val.GetValue())); ok {
			return result, nil
		}
		return 0, truncationError[To](val.GetValue())
	} else if val, ok := any(v).(interface{ GetValue() int32 }); ok {

		if result, ok := narrow.ToInteger[To](int64(val.GetValue())); ok {
			return result, nil
		}
		return 0, truncationError[To](val.GetValue())
	}
	// This cannot be reached because this function is constrained to only
	// take FHIR Elements that fit one of the above two types.
	return result, ErrIntegerTruncated
}

// MustConvertToInteger converts a FHIR Integer type into a Go native integer type.
//
// If the value stored in the integer type cannot fit into the receiver type,
// this function will panic.
func MustConvertToInteger[To constraints.Integer, From integerType](v From) To {
	result, err := ToInteger[To](v)
	if err != nil {
		panic(err)
	}
	return result
}

// truncationError forms an Error type for truncation errors.
func truncationError[To constraints.Integer, From constraints.Integer](value From) error {
	var result To
	return fmt.Errorf("%w: type %T with value %v does not fit into receiver %T", ErrIntegerTruncated, value, value, result)
}
