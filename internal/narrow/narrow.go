/*
Package narrow provides conversion functionality for narrowing
integer types in a safe and generic manner.

These operations will return not only the casted integer, but also a boolean
indicating whether the value was safely casted to the receiver type -- e.g.
that the sign didn't change through casting, and that the value was not unexpectedly
truncated. This can be really useful for cases where data requires explicit
integer-widths, and the incoming data may exceed that integer width.

Conversion is done through one of the two following function types:

  - `narrow.ToInteger[To](from) (To, bool)` -- converts to the specified int-type
  - `narrow.To*(from) (*, bool)` -- convenience functions for converting
    to explicit integer-width types to avoid needing generics (just calls
    into the above).

The former's generic-interface allows better composability with other generic
functions, whereas the latter allows a more idiomatic and visible call for
explicit cases (e.g. `ints.ToInt32(...)`)
*/
package narrow

import (
	"math"

	"golang.org/x/exp/constraints"
)

// isSigned checks that the integer input is a Signed integer type.
func isSigned[T constraints.Integer]() bool {
	// If T is unsigned, -1 will form a large positive number when casted to
	// T -- which makes this a simple way to check for signedness.
	val := -1
	return T(val) < 0
}

// ToInteger converts `From` to a `To` type, additionally returning whether the
// value was converted successfully without any truncation and loss-of-data occurring.
func ToInteger[To constraints.Integer, From constraints.Integer](from From) (To, bool) {
	if isSigned[From]() {
		return To(from), canSafelyNarrowSigned[To](int64(from))
	}
	return To(from), canSafelyNarrowUnsigned[To](uint64(from))
}

// ToInt converts an integer type into an `int`, additionally returning whether the
// value was converted successfully without any truncation and loss-of-data occurring.
func ToInt[From constraints.Integer](from From) (int, bool) {
	return ToInteger[int](from)
}

// ToInt8 converts an integer type into an `int8`, additionally returning whether the
// value was converted successfully without any truncation and loss-of-data occurring.
func ToInt8[From constraints.Integer](from From) (int8, bool) {
	return ToInteger[int8](from)
}

// ToInt16 converts an integer type into an `int16`, additionally returning whether the
// value was converted successfully without any truncation and loss-of-data occurring.
func ToInt16[From constraints.Integer](from From) (int16, bool) {
	return ToInteger[int16](from)
}

// ToInt32 converts an integer type into an `int32`, additionally returning whether the
// value was converted successfully without any truncation and loss-of-data occurring.
func ToInt32[From constraints.Integer](from From) (int32, bool) {
	return ToInteger[int32](from)
}

// ToInt64 converts an integer type into an `int64`, additionally returning whether the
// value was converted successfully without any truncation and loss-of-data occurring.
func ToInt64[From constraints.Integer](from From) (int64, bool) {
	return ToInteger[int64](from)
}

// ToUint converts an integer type into an `uint`, additionally returning whether the
// value was converted successfully without any truncation and loss-of-data occurring.
func ToUint[From constraints.Integer](from From) (uint, bool) {
	return ToInteger[uint](from)
}

// ToUint8 converts an integer type into an `uint8`, additionally returning whether the
// value was converted successfully without any truncation and loss-of-data occurring.
func ToUint8[From constraints.Integer](from From) (uint8, bool) {
	return ToInteger[uint8](from)
}

// ToUint16 converts an integer type into an `uint16`, additionally returning whether the
// value was converted successfully without any truncation and loss-of-data occurring.
func ToUint16[From constraints.Integer](from From) (uint16, bool) {
	return ToInteger[uint16](from)
}

// ToUint32 converts an integer type into an `uint32`, additionally returning whether the
// value was converted successfully without any truncation and loss-of-data occurring.
func ToUint32[From constraints.Integer](from From) (uint32, bool) {
	return ToInteger[uint32](from)
}

// ToUint64 converts an integer type into an `uint64`, additionally returning whether the
// value was converted successfully without any truncation and loss-of-data occurring.
func ToUint64[From constraints.Integer](from From) (uint64, bool) {
	return ToInteger[uint64](from)
}

// ToUintptr converts an integer type into an `uintptr`, as well as whether the value was
// truncated.
func ToUintptr[From constraints.Integer](from From) (uintptr, bool) {
	return ToInteger[uintptr](from)
}

// canSafelyNarrowSigned tests whether the supplied val can be converted into 'To'
// without truncation or loss-of-data.
func canSafelyNarrowSigned[To constraints.Integer](val int64) bool {
	var v To
	switch any(v).(type) {
	case uint:
		// Cast here is to avoid possible "overflow" errors if uint is 64-bits
		return uint64(val) <= math.MaxUint && val >= 0
	case uintptr:
		return int64(uintptr(val)) == val && val >= 0
	case uint8:
		return val <= math.MaxUint8 && val >= 0
	case uint16:
		return val <= math.MaxUint16 && val >= 0
	case uint32:
		return val <= math.MaxUint32 && val >= 0
	case uint64:
		return val >= 0
	case int:
		return val <= math.MaxInt && val >= math.MinInt
	case int8:
		return val <= math.MaxInt8 && val >= math.MinInt8
	case int16:
		return val <= math.MaxInt16 && val >= math.MinInt16
	case int32:
		return val <= math.MaxInt32 && val >= math.MinInt32
	case int64:
		return true // trivially true (identity conversion)
	}
	return false // this should be unreachable
}

// canSafelyNarrowUnsigned tests whether the supplied val can be converted into 'To'
// without truncation or loss-of-data.
func canSafelyNarrowUnsigned[To constraints.Integer](val uint64) bool {
	var v To
	switch any(v).(type) {
	case uint:
		// Cast here is to avoid possible "overflow" errors if uint is 64-bits
		return uint64(val) <= math.MaxUint
	case uintptr:
		return uint64(uintptr(val)) == val
	case uint8:
		return val <= math.MaxUint8
	case uint16:
		return val <= math.MaxUint16
	case uint32:
		return val <= math.MaxUint32
	case uint64:
		return true // trivially true (identity conversion)
	case int:
		// Cast here is to avoid possible "overflow" errors if int is 64-bits
		return uint64(val) <= math.MaxInt
	case int8:
		return val <= math.MaxInt8
	case int16:
		return val <= math.MaxInt16
	case int32:
		return val <= math.MaxInt32
	case int64:
		return val <= math.MaxInt64
	}
	return false // this should be unreachable
}
