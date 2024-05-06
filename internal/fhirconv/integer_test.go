package fhirconv_test

import (
	"errors"
	"math"
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirconv"
	"golang.org/x/exp/constraints"
)

type integerType interface {
	*dtpb.Integer | *dtpb.UnsignedInt | *dtpb.PositiveInt
}

func wantTruncation[To constraints.Integer, From integerType](value From) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		var to To

		_, err := fhirconv.ToInteger[To](value)

		if got, want := err, fhirconv.ErrIntegerTruncated; !errors.Is(got, want) {
			t.Errorf("ToInteger[%T]: got err %v, want %v", to, got, want)
		}
	}
}

func wantConversion[To constraints.Integer, From integerType, Input constraints.Integer](conv func(Input) From, input Input) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		var to To
		value := conv(input)

		got, err := fhirconv.ToInteger[To](value)
		if err != nil {
			t.Fatalf("ToInteger[%T]: unexpected error '%v'", to, err)
		}

		if got, want := got, To(input); got != want {
			t.Errorf("ToInteger[%T]: got err %v, want %v", to, got, want)
		}
	}
}

func TestToInteger_ValueExceedsReceiver_ReturnsError(t *testing.T) {
	// Note: nested sub-tests are being done both for organization and test naming.
	// We also can't use table-driven tests here, since the input is a type -- so
	// this forces duplication.
	t.Run("FromPositiveInt", func(t *testing.T) {
		t.Run("ToInt8", wantTruncation[int8](fhir.PositiveInt(math.MaxUint32)))
		t.Run("ToInt16", wantTruncation[int16](fhir.PositiveInt(math.MaxUint32)))
		t.Run("ToInt32", wantTruncation[int32](fhir.PositiveInt(math.MaxUint32)))
		t.Run("ToUint8", wantTruncation[uint8](fhir.PositiveInt(math.MaxUint32)))
		t.Run("ToUint16", wantTruncation[uint16](fhir.PositiveInt(math.MaxUint32)))
	})

	t.Run("FromUnsignedInt", func(t *testing.T) {
		t.Run("ToInt8", wantTruncation[int8](fhir.UnsignedInt(math.MaxUint32)))
		t.Run("ToInt16", wantTruncation[int16](fhir.UnsignedInt(math.MaxUint32)))
		t.Run("ToInt32", wantTruncation[int32](fhir.UnsignedInt(math.MaxUint32)))
		t.Run("ToUint8", wantTruncation[uint8](fhir.UnsignedInt(math.MaxUint32)))
		t.Run("ToUint16", wantTruncation[uint16](fhir.UnsignedInt(math.MaxUint32)))
	})

	t.Run("FromInteger", func(t *testing.T) {
		t.Run("ToInt8", wantTruncation[int8](fhir.Integer(math.MaxInt32)))
		t.Run("ToInt16", wantTruncation[int16](fhir.Integer(math.MaxInt32)))
		// int32, int64, and int can't truncate
		t.Run("ToUint8", wantTruncation[uint8](fhir.Integer(-1)))
		t.Run("ToUint16", wantTruncation[uint16](fhir.Integer(-1)))
		t.Run("ToUint32", wantTruncation[uint32](fhir.Integer(-1)))
		t.Run("ToUint32", wantTruncation[uint64](fhir.Integer(-1)))
		t.Run("ToUintptr", wantTruncation[uintptr](fhir.Integer(-1)))
		t.Run("ToUint", wantTruncation[uint](fhir.Integer(-1)))
	})
}

func TestToInteger_ValueWithinRange_ReturnsValue(t *testing.T) {
	// Note: nested sub-tests are being done both for organization and test naming.
	// We also can't use table-driven tests here, since the input is a type -- so
	// this forces duplication.
	t.Run("FromPositiveInt", func(t *testing.T) {
		t.Run("ToInt8", wantConversion[int8](fhir.PositiveInt, 0))
		t.Run("ToInt8Max", wantConversion[int8](fhir.PositiveInt, math.MaxInt8))
		t.Run("ToInt16", wantConversion[int16](fhir.PositiveInt, 0))
		t.Run("ToInt16Max", wantConversion[int16](fhir.PositiveInt, math.MaxInt16))
		t.Run("ToInt32", wantConversion[int32](fhir.PositiveInt, 0))
		t.Run("ToInt32Max", wantConversion[int32](fhir.PositiveInt, math.MaxInt32))
		t.Run("ToInt64", wantConversion[int64](fhir.PositiveInt, 0))
		t.Run("ToInt64MaxUint32", wantConversion[int64](fhir.PositiveInt, math.MaxUint32))
		t.Run("ToInt", wantConversion[int](fhir.PositiveInt, 0))
		t.Run("ToIntMaxInt32", wantConversion[int](fhir.PositiveInt, math.MaxInt32))
		t.Run("ToUint8", wantConversion[uint8](fhir.PositiveInt, 0))
		t.Run("ToUint8Max", wantConversion[uint8](fhir.PositiveInt, math.MaxUint8))
		t.Run("ToUint16", wantConversion[uint16](fhir.PositiveInt, 0))
		t.Run("ToUint16Max", wantConversion[uint16](fhir.PositiveInt, math.MaxUint8))
		t.Run("ToUint32", wantConversion[uint32](fhir.PositiveInt, 0))
		t.Run("ToUint32Max", wantConversion[uint32](fhir.PositiveInt, math.MaxUint32))
		t.Run("ToUint64", wantConversion[uint64](fhir.PositiveInt, 0))
		t.Run("ToUint64MaxUint32", wantConversion[uint64](fhir.PositiveInt, math.MaxUint32))
		t.Run("ToUint", wantConversion[uint](fhir.PositiveInt, 0))
		t.Run("ToUintMaxUint32", wantConversion[uint](fhir.PositiveInt, math.MaxUint32))
		t.Run("ToUintptr", wantConversion[uintptr](fhir.PositiveInt, 0))
	})

	t.Run("FromUnsignedInt", func(t *testing.T) {
		t.Run("ToInt8", wantConversion[int8](fhir.UnsignedInt, 0))
		t.Run("ToInt8Max", wantConversion[int8](fhir.UnsignedInt, math.MaxInt8))
		t.Run("ToInt16", wantConversion[int16](fhir.UnsignedInt, 0))
		t.Run("ToInt16Max", wantConversion[int16](fhir.UnsignedInt, math.MaxInt16))
		t.Run("ToInt32", wantConversion[int32](fhir.UnsignedInt, 0))
		t.Run("ToInt32Max", wantConversion[int32](fhir.UnsignedInt, math.MaxInt32))
		t.Run("ToInt64", wantConversion[int64](fhir.UnsignedInt, 0))
		t.Run("ToInt64MaxUint32", wantConversion[int64](fhir.UnsignedInt, math.MaxUint32))
		t.Run("ToInt", wantConversion[int](fhir.UnsignedInt, 0))
		t.Run("ToIntMaxInt32", wantConversion[int](fhir.UnsignedInt, math.MaxInt32))
		t.Run("ToUint8", wantConversion[uint8](fhir.UnsignedInt, 0))
		t.Run("ToUint8Max", wantConversion[uint8](fhir.UnsignedInt, math.MaxUint8))
		t.Run("ToUint16", wantConversion[uint16](fhir.UnsignedInt, 0))
		t.Run("ToUint16Max", wantConversion[uint16](fhir.UnsignedInt, math.MaxUint8))
		t.Run("ToUint32", wantConversion[uint32](fhir.UnsignedInt, 0))
		t.Run("ToUint32Max", wantConversion[uint32](fhir.UnsignedInt, math.MaxUint32))
		t.Run("ToUint64", wantConversion[uint64](fhir.UnsignedInt, 0))
		t.Run("ToUint64MaxUint32", wantConversion[uint64](fhir.UnsignedInt, math.MaxUint32))
		t.Run("ToUint", wantConversion[uint](fhir.UnsignedInt, 0))
		t.Run("ToUintMaxUint32", wantConversion[uint](fhir.UnsignedInt, math.MaxUint32))
		t.Run("ToUintptr", wantConversion[uintptr](fhir.UnsignedInt, 0))
	})

	t.Run("FromInteger", func(t *testing.T) {
		t.Run("ToInt8", wantConversion[int8](fhir.Integer, 0))
		t.Run("ToInt8Max", wantConversion[int8](fhir.Integer, math.MaxInt8))
		t.Run("ToInt8Min", wantConversion[int8](fhir.Integer, math.MinInt8))
		t.Run("ToInt16", wantConversion[int16](fhir.Integer, 0))
		t.Run("ToInt16Max", wantConversion[int16](fhir.Integer, math.MaxInt16))
		t.Run("ToInt16Min", wantConversion[int16](fhir.Integer, math.MinInt16))
		t.Run("ToInt32", wantConversion[int32](fhir.Integer, 0))
		t.Run("ToInt32Max", wantConversion[int32](fhir.Integer, math.MaxInt32))
		t.Run("ToInt32Min", wantConversion[int32](fhir.Integer, math.MinInt32))
		t.Run("ToInt64", wantConversion[int64](fhir.Integer, 0))
		t.Run("ToInt64MaxInt32", wantConversion[int64](fhir.Integer, math.MaxInt32))
		t.Run("ToInt64MinInt32", wantConversion[int64](fhir.Integer, math.MinInt32))
		t.Run("ToInt", wantConversion[int](fhir.Integer, 0))
		t.Run("ToIntMaxInt32", wantConversion[int](fhir.Integer, math.MaxInt32))
		t.Run("ToIntMinInt32", wantConversion[int](fhir.Integer, math.MinInt32))
		t.Run("ToUint8", wantConversion[uint8](fhir.Integer, 0))
		t.Run("ToUint8Max", wantConversion[uint8](fhir.Integer, math.MaxUint8))
		t.Run("ToUint16", wantConversion[uint16](fhir.Integer, 0))
		t.Run("ToUint16Max", wantConversion[uint16](fhir.Integer, math.MaxUint16))
		t.Run("ToUint32", wantConversion[uint32](fhir.Integer, 0))
		t.Run("ToUint32MaxInt32", wantConversion[uint32](fhir.Integer, math.MaxInt32))
		t.Run("ToUint64", wantConversion[uint64](fhir.Integer, 0))
		t.Run("ToUint64MaxInt32", wantConversion[uint64](fhir.Integer, math.MaxInt32))
		t.Run("ToUint", wantConversion[uint](fhir.Integer, 0))
		t.Run("ToUintMaxInt32", wantConversion[uint](fhir.Integer, math.MaxInt32))
		t.Run("ToUintptr", wantConversion[uintptr](fhir.Integer, 0))
	})
}

func TestMustConvertToInteger_ValueWithinRange_ReturnsValue(t *testing.T) {
	input := int32(42)
	value := fhir.Integer(input)

	result := fhirconv.MustConvertToInteger[int32](value)

	if got, want := result, input; got != want {
		t.Fatalf("MustConvertToInteger[int32]: got %v, want %v", got, want)
	}
}

func TestMustConvertToInteger_ValueOutsideRange_Panics(t *testing.T) {
	defer func() { _ = recover() }()
	value := fhir.Integer(-1)

	fhirconv.MustConvertToInteger[uint8](value)

	// This can't be reached if a panic occurs
	t.Fatalf("MustConvertToInt: expected panic")
}
