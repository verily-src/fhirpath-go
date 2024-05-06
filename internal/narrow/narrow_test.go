package narrow_test

import (
	"math"
	"testing"

	"github.com/verily-src/fhirpath-go/internal/narrow"
	"golang.org/x/exp/constraints"
)

func wantTruncation[To constraints.Integer, From constraints.Integer](input From) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		var to To

		_, ok := narrow.ToInteger[To](input)

		if got, want := ok, false; got != want {
			t.Errorf("ToInteger[%T]: got %v, want %v", to, got, want)
		}
	}
}

func wantConversion[To constraints.Integer, From constraints.Integer](input From) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		var to To

		got, ok := narrow.ToInteger[To](input)
		if ok == false {
			t.Fatalf("ToInteger[%T]: unexpected truncation %v => %v", to, input, got)
		}

		if got, want := From(got), input; got != want {
			t.Errorf("ToInteger[%T]: got %v, want %v", to, got, want)
		}
	}
}

func wantConversionFunc[To constraints.Integer, From constraints.Integer](fn func(From) (To, bool), from From) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		var to To

		got, ok := fn(from)
		if ok == false {
			t.Fatalf("%T => %T: unexpected truncation %v => %v", from, to, from, got)
		}

		if got, want := From(got), from; got != want {
			t.Errorf("%T => %T: got %v, want %v", from, to, got, want)
		}
	}
}

func TestToInteger_ValueExceedsReceiver_ReturnsFalse(t *testing.T) {
	// Note: nested sub-tests are being done both for organization and test naming.
	// We also can't use table-driven tests here, since the input is a type -- so
	// this forces duplication.
	t.Run("FromUint8", func(t *testing.T) {
		t.Run("ToInt8", wantTruncation[int8](uint8(math.MaxUint8)))
	})

	t.Run("FromUint16", func(t *testing.T) {
		t.Run("ToInt8", wantTruncation[int8](uint16(math.MaxUint16)))
		t.Run("ToInt16", wantTruncation[int16](uint16(math.MaxUint16)))
		t.Run("ToUint8", wantTruncation[uint8](uint16(math.MaxUint16)))
	})

	t.Run("FromUint32", func(t *testing.T) {
		t.Run("ToInt8", wantTruncation[int8](uint32(math.MaxUint32)))
		t.Run("ToInt16", wantTruncation[int16](uint32(math.MaxUint32)))
		t.Run("ToInt32", wantTruncation[int32](uint32(math.MaxUint32)))

		t.Run("ToUint8", wantTruncation[uint8](uint32(math.MaxUint32)))
		t.Run("ToUint16", wantTruncation[uint16](uint32(math.MaxUint32)))
	})

	t.Run("FromUint64", func(t *testing.T) {
		t.Run("ToInt8", wantTruncation[int8](uint64(math.MaxUint64)))
		t.Run("ToInt16", wantTruncation[int16](uint64(math.MaxUint64)))
		t.Run("ToInt32", wantTruncation[int32](uint64(math.MaxUint64)))
		t.Run("ToInt64", wantTruncation[int64](uint64(math.MaxUint64)))
		t.Run("ToInt", wantTruncation[int](uint64(math.MaxUint64)))

		t.Run("ToUint8", wantTruncation[uint8](uint64(math.MaxUint64)))
		t.Run("ToUint16", wantTruncation[uint16](uint64(math.MaxUint64)))
		t.Run("ToUint32", wantTruncation[uint32](uint64(math.MaxUint64)))
	})

	t.Run("FromUint", func(t *testing.T) {
		t.Run("ToInt8", wantTruncation[int8](uint(math.MaxUint)))
		t.Run("ToInt16", wantTruncation[int16](uint(math.MaxUint)))
		t.Run("ToInt32", wantTruncation[int32](uint(math.MaxUint)))

		t.Run("ToUint8", wantTruncation[uint8](uint(math.MaxUint)))
		t.Run("ToUint16", wantTruncation[uint16](uint(math.MaxUint)))
		// We can't reliably test against uint32, since uint may be AT LEAST 32 bits.
	})

	t.Run("FromInt8", func(t *testing.T) {
		t.Run("ToUint8", wantTruncation[uint8](int8(-1)))
		t.Run("ToUint16", wantTruncation[uint16](int8(-1)))
		t.Run("ToUint32", wantTruncation[uint32](int8(-1)))
		t.Run("ToUint32", wantTruncation[uint64](int8(-1)))
		t.Run("ToUintptr", wantTruncation[uintptr](int8(-1)))
		t.Run("ToUint", wantTruncation[uint](int8(-1)))
	})

	t.Run("FromInt16", func(t *testing.T) {
		t.Run("ToInt8", wantTruncation[int8](int16(math.MaxInt16)))

		t.Run("ToUint8", wantTruncation[uint8](int16(-1)))
		t.Run("ToUint16", wantTruncation[uint16](int16(-1)))
		t.Run("ToUint32", wantTruncation[uint32](int16(-1)))
		t.Run("ToUint32", wantTruncation[uint64](int16(-1)))
		t.Run("ToUintptr", wantTruncation[uintptr](int16(-1)))
		t.Run("ToUint", wantTruncation[uint](int16(-1)))
	})

	t.Run("FromInt32", func(t *testing.T) {
		t.Run("ToInt8", wantTruncation[int8](int32(math.MaxInt32)))
		t.Run("ToInt16", wantTruncation[int16](int32(math.MaxInt32)))

		t.Run("ToUint8", wantTruncation[uint8](int32(-1)))
		t.Run("ToUint16", wantTruncation[uint16](int32(-1)))
		t.Run("ToUint32", wantTruncation[uint32](int32(-1)))
		t.Run("ToUint32", wantTruncation[uint64](int32(-1)))
		t.Run("ToUintptr", wantTruncation[uintptr](int32(-1)))
		t.Run("ToUint", wantTruncation[uint](int32(-1)))
	})

	t.Run("FromInt64", func(t *testing.T) {
		t.Run("ToInt8", wantTruncation[int8](int64(math.MaxInt64)))
		t.Run("ToInt16", wantTruncation[int16](int64(math.MaxInt64)))
		t.Run("ToInt32", wantTruncation[int32](int64(math.MaxInt64)))

		t.Run("ToUint8", wantTruncation[uint8](int64(-1)))
		t.Run("ToUint16", wantTruncation[uint16](int64(-1)))
		t.Run("ToUint32", wantTruncation[uint32](int64(-1)))
		t.Run("ToUint32", wantTruncation[uint64](int64(-1)))
		t.Run("ToUintptr", wantTruncation[uintptr](int64(-1)))
		t.Run("ToUint", wantTruncation[uint](int64(-1)))
	})

	t.Run("FromInt", func(t *testing.T) {
		t.Run("ToInt8", wantTruncation[int8](int(math.MaxInt)))
		t.Run("ToInt16", wantTruncation[int16](int(math.MaxInt)))
		// We can't reliably test against int32, since int is at LEAST 32 bits.

		t.Run("ToUint8", wantTruncation[uint8](int(-1)))
		t.Run("ToUint16", wantTruncation[uint16](int(-1)))
		t.Run("ToUint32", wantTruncation[uint32](int(-1)))
		t.Run("ToUint32", wantTruncation[uint64](int(-1)))
		t.Run("ToUintptr", wantTruncation[uintptr](int(-1)))
		t.Run("ToUint", wantTruncation[uint](int(-1)))
	})
}

func TestToInteger_ValueFitsIntoReceiver_ReturnsValueAndTrue(t *testing.T) {
	// Note: nested sub-tests are being done both for organization and test naming.
	// We also can't use table-driven tests here, since the input is a type -- so
	// this forces duplication.
	t.Run("FromUint", func(t *testing.T) {
		t.Run("ToInt8", wantConversion[int8](uint(0)))
		t.Run("ToInt8Max", wantConversion[int8](uint(math.MaxInt8)))
		t.Run("ToInt16", wantConversion[int16](uint(0)))
		t.Run("ToInt16Max", wantConversion[int16](uint(math.MaxInt16)))
		t.Run("ToInt32", wantConversion[int32](uint(0)))
		t.Run("ToInt32Max", wantConversion[int32](uint(math.MaxInt32)))
		t.Run("ToInt64", wantConversion[int64](uint(0)))
		t.Run("ToInt64MaxUint32", wantConversion[int64](uint(math.MaxUint32)))
		t.Run("ToInt", wantConversion[int](uint(0)))
		t.Run("ToIntMaxInt32", wantConversion[int](uint(math.MaxInt32)))
		t.Run("ToUint8", wantConversion[uint8](uint(0)))
		t.Run("ToUint8Max", wantConversion[uint8](uint(math.MaxUint8)))
		t.Run("ToUint16", wantConversion[uint16](uint(0)))
		t.Run("ToUint16Max", wantConversion[uint16](uint(math.MaxUint16)))
		t.Run("ToUint32", wantConversion[uint32](uint(0)))
		t.Run("ToUint32Max", wantConversion[uint32](uint(math.MaxUint32)))
		t.Run("ToUint64", wantConversion[uint64](uint(0)))
		t.Run("ToUint64MaxUint32", wantConversion[uint64](uint(math.MaxUint32)))
		t.Run("ToUint", wantConversion[uint](uint(0)))
		t.Run("ToUintMaxUint32", wantConversion[uint](uint(math.MaxUint32)))
		t.Run("ToUintptr", wantConversion[uintptr](uint(0)))
	})

	t.Run("FromUint8", func(t *testing.T) {
		t.Run("ToInt8", wantConversion[int8](uint8(0)))
		t.Run("ToInt8Max", wantConversion[int8](uint8(math.MaxInt8)))
		t.Run("ToInt16", wantConversion[int16](uint8(0)))
		t.Run("ToInt16MaxInt8", wantConversion[int16](uint8(math.MaxInt8)))
		t.Run("ToInt32", wantConversion[int32](uint8(0)))
		t.Run("ToInt32MaxUint8", wantConversion[int32](uint8(math.MaxUint8)))
		t.Run("ToInt64", wantConversion[int64](uint8(0)))
		t.Run("ToInt64MaxUint8", wantConversion[int64](uint8(math.MaxUint8)))
		t.Run("ToInt", wantConversion[int](uint8(0)))
		t.Run("ToIntMaxInt8", wantConversion[int](uint8(math.MaxInt8)))
		t.Run("ToUint8", wantConversion[uint8](uint8(0)))
		t.Run("ToUint8Max", wantConversion[uint8](uint8(math.MaxUint8)))
		t.Run("ToUint16", wantConversion[uint16](uint8(0)))
		t.Run("ToUint16MaxInt8", wantConversion[uint16](uint8(math.MaxUint8)))
		t.Run("ToUint32", wantConversion[uint32](uint8(0)))
		t.Run("ToUint32MaxUint8", wantConversion[uint32](uint8(math.MaxUint8)))
		t.Run("ToUint64", wantConversion[uint64](uint8(0)))
		t.Run("ToUint64MaxUint8", wantConversion[uint64](uint8(math.MaxUint8)))
		t.Run("ToUint", wantConversion[uint](uint8(0)))
		t.Run("ToUintMaxUint8", wantConversion[uint](uint8(math.MaxUint8)))
		t.Run("ToUintptr", wantConversion[uintptr](uint8(0)))
	})

	t.Run("FromUint16", func(t *testing.T) {
		t.Run("ToInt8", wantConversion[int8](uint16(0)))
		t.Run("ToInt8Max", wantConversion[int8](uint16(math.MaxInt8)))
		t.Run("ToInt16", wantConversion[int16](uint16(0)))
		t.Run("ToInt16Max", wantConversion[int16](uint16(math.MaxInt16)))
		t.Run("ToInt32", wantConversion[int32](uint16(0)))
		t.Run("ToInt32MaxUint16", wantConversion[int32](uint16(math.MaxInt16)))
		t.Run("ToInt64", wantConversion[int64](uint16(0)))
		t.Run("ToInt64MaxUint16", wantConversion[int64](uint16(math.MaxUint16)))
		t.Run("ToInt", wantConversion[int](uint16(0)))
		t.Run("ToIntMaxInt16", wantConversion[int](uint16(math.MaxInt16)))
		t.Run("ToUint8", wantConversion[uint8](uint16(0)))
		t.Run("ToUint8Max", wantConversion[uint8](uint16(math.MaxUint8)))
		t.Run("ToUint16", wantConversion[uint16](uint16(0)))
		t.Run("ToUint16Max", wantConversion[uint16](uint16(math.MaxUint16)))
		t.Run("ToUint32", wantConversion[uint32](uint16(0)))
		t.Run("ToUint32MaxUint16", wantConversion[uint32](uint16(math.MaxUint16)))
		t.Run("ToUint64", wantConversion[uint64](uint16(0)))
		t.Run("ToUint64MaxUint16", wantConversion[uint64](uint16(math.MaxUint16)))
		t.Run("ToUint", wantConversion[uint](uint16(0)))
		t.Run("ToUintMaxUint16", wantConversion[uint](uint16(math.MaxUint16)))
		t.Run("ToUintptr", wantConversion[uintptr](uint16(0)))
	})

	t.Run("FromUint32", func(t *testing.T) {
		t.Run("ToInt8", wantConversion[int8](uint32(0)))
		t.Run("ToInt8Max", wantConversion[int8](uint32(math.MaxInt8)))
		t.Run("ToInt16", wantConversion[int16](uint32(0)))
		t.Run("ToInt16Max", wantConversion[int16](uint32(math.MaxInt16)))
		t.Run("ToInt32", wantConversion[int32](uint32(0)))
		t.Run("ToInt32Max", wantConversion[int32](uint32(math.MaxInt32)))
		t.Run("ToInt64", wantConversion[int64](uint32(0)))
		t.Run("ToInt64MaxUint32", wantConversion[int64](uint32(math.MaxUint32)))
		t.Run("ToInt", wantConversion[int](uint32(0)))
		t.Run("ToIntMaxInt32", wantConversion[int](uint32(math.MaxInt32)))
		t.Run("ToUint8", wantConversion[uint8](uint32(0)))
		t.Run("ToUint8Max", wantConversion[uint8](uint32(math.MaxUint8)))
		t.Run("ToUint16", wantConversion[uint16](uint32(0)))
		t.Run("ToUint16Max", wantConversion[uint16](uint32(math.MaxUint16)))
		t.Run("ToUint32", wantConversion[uint32](uint32(0)))
		t.Run("ToUint32Max", wantConversion[uint32](uint32(math.MaxUint32)))
		t.Run("ToUint64", wantConversion[uint64](uint32(0)))
		t.Run("ToUint64MaxUint32", wantConversion[uint64](uint32(math.MaxUint32)))
		t.Run("ToUint", wantConversion[uint](uint32(0)))
		t.Run("ToUintMaxUint32", wantConversion[uint](uint32(math.MaxUint32)))
		t.Run("ToUintptr", wantConversion[uintptr](uint32(0)))
	})

	t.Run("FromUint64", func(t *testing.T) {
		t.Run("ToInt8", wantConversion[int8](uint64(0)))
		t.Run("ToInt8Max", wantConversion[int8](uint64(math.MaxInt8)))
		t.Run("ToInt16", wantConversion[int16](uint64(0)))
		t.Run("ToInt16Max", wantConversion[int16](uint64(math.MaxInt16)))
		t.Run("ToInt32", wantConversion[int32](uint64(0)))
		t.Run("ToInt32Max", wantConversion[int32](uint64(math.MaxInt32)))
		t.Run("ToInt64", wantConversion[int64](uint64(0)))
		t.Run("ToInt64Max", wantConversion[int64](uint64(math.MaxInt64)))
		t.Run("ToInt", wantConversion[int](uint64(0)))
		t.Run("ToIntMax", wantConversion[int](uint64(math.MaxInt64)))
		t.Run("ToUint8", wantConversion[uint8](uint64(0)))
		t.Run("ToUint8Max", wantConversion[uint8](uint64(math.MaxUint8)))
		t.Run("ToUint16", wantConversion[uint16](uint64(0)))
		t.Run("ToUint16Max", wantConversion[uint16](uint64(math.MaxUint16)))
		t.Run("ToUint32", wantConversion[uint32](uint64(0)))
		t.Run("ToUint32Max", wantConversion[uint32](uint64(math.MaxUint32)))
		t.Run("ToUint64", wantConversion[uint64](uint64(0)))
		t.Run("ToUint64Max", wantConversion[uint64](uint64(math.MaxUint64)))
		t.Run("ToUint", wantConversion[uint](uint64(0)))
		t.Run("ToUintMax", wantConversion[uint](uint64(math.MaxUint64)))
		t.Run("ToUintptr", wantConversion[uintptr](uint64(0)))
	})

	t.Run("FromInt", func(t *testing.T) {
		t.Run("ToInt8", wantConversion[int8](int(0)))
		t.Run("ToInt8Max", wantConversion[int8](int(math.MaxInt8)))
		t.Run("ToInt8Min", wantConversion[int8](int(math.MinInt8)))
		t.Run("ToInt16", wantConversion[int16](int(0)))
		t.Run("ToInt16Max", wantConversion[int16](int(math.MaxInt16)))
		t.Run("ToInt16Min", wantConversion[int16](int(math.MinInt16)))
		t.Run("ToInt32", wantConversion[int32](int(0)))
		t.Run("ToInt32Max", wantConversion[int32](int(math.MaxInt32)))
		t.Run("ToInt32Min", wantConversion[int32](int(math.MinInt32)))
		t.Run("ToInt64", wantConversion[int64](int(0)))
		t.Run("ToInt64MaxInt32", wantConversion[int64](int(math.MaxInt32)))
		t.Run("ToInt64MinInt32", wantConversion[int64](int(math.MinInt32)))
		t.Run("ToInt", wantConversion[int](int(0)))
		t.Run("ToIntMaxInt32", wantConversion[int](int(math.MaxInt32)))
		t.Run("ToIntMinInt32", wantConversion[int](int(math.MinInt32)))
		t.Run("ToUint8", wantConversion[uint8](int(0)))
		t.Run("ToUint8Max", wantConversion[uint8](int(math.MaxUint8)))
		t.Run("ToUint16", wantConversion[uint16](int(0)))
		t.Run("ToUint16Max", wantConversion[uint16](int(math.MaxUint16)))
		t.Run("ToUint32", wantConversion[uint32](int(0)))
		t.Run("ToUint32MaxInt32", wantConversion[uint32](int(math.MaxInt32)))
		t.Run("ToUint64", wantConversion[uint64](int(0)))
		t.Run("ToUint64MaxInt32", wantConversion[uint64](int(math.MaxInt32)))
		t.Run("ToUint", wantConversion[uint](int(0)))
		t.Run("ToUintMaxInt32", wantConversion[uint](int(math.MaxInt32)))
		t.Run("ToUintptr", wantConversion[uintptr](int(0)))
	})

	t.Run("FromInt8", func(t *testing.T) {
		t.Run("ToInt8", wantConversion[int8](int8(0)))
		t.Run("ToInt8Max", wantConversion[int8](int8(math.MaxInt8)))
		t.Run("ToInt8Min", wantConversion[int8](int8(math.MinInt8)))
		t.Run("ToInt16", wantConversion[int16](int8(0)))
		t.Run("ToInt16MaxInt8", wantConversion[int16](int8(math.MaxInt8)))
		t.Run("ToInt16MinInt8", wantConversion[int16](int8(math.MinInt8)))
		t.Run("ToInt32", wantConversion[int32](int8(0)))
		t.Run("ToInt32MaxInt8", wantConversion[int32](int8(math.MaxInt8)))
		t.Run("ToInt32MinInt8", wantConversion[int32](int8(math.MinInt8)))
		t.Run("ToInt64", wantConversion[int64](int8(0)))
		t.Run("ToInt64MaxInt8", wantConversion[int64](int8(math.MaxInt8)))
		t.Run("ToInt64MinInt8", wantConversion[int64](int8(math.MinInt8)))
		t.Run("ToInt", wantConversion[int](int8(0)))
		t.Run("ToIntMaxInt8", wantConversion[int](int8(math.MaxInt8)))
		t.Run("ToIntMinInt8", wantConversion[int](int8(math.MinInt8)))
		t.Run("ToUint8", wantConversion[uint8](int8(0)))
		t.Run("ToUint8MaxInt8", wantConversion[uint8](int8(math.MaxInt8)))
		t.Run("ToUint16", wantConversion[uint16](int8(0)))
		t.Run("ToUint16MaxInt8", wantConversion[uint16](int8(math.MaxInt8)))
		t.Run("ToUint32", wantConversion[uint32](int8(0)))
		t.Run("ToUint32MaxInt9", wantConversion[uint32](int8(math.MaxInt8)))
		t.Run("ToUint64", wantConversion[uint64](int8(0)))
		t.Run("ToUint64MaxInt8", wantConversion[uint64](int8(math.MaxInt8)))
		t.Run("ToUint", wantConversion[uint](int8(0)))
		t.Run("ToUintMaxInt8", wantConversion[uint](int8(math.MaxInt8)))
		t.Run("ToUintptr", wantConversion[uintptr](int8(0)))
	})

	t.Run("FromInt16", func(t *testing.T) {
		t.Run("ToInt8", wantConversion[int8](int16(0)))
		t.Run("ToInt8Max", wantConversion[int8](int16(math.MaxInt8)))
		t.Run("ToInt8Min", wantConversion[int8](int16(math.MinInt8)))
		t.Run("ToInt16", wantConversion[int16](int16(0)))
		t.Run("ToInt16Max", wantConversion[int16](int16(math.MaxInt16)))
		t.Run("ToInt16Min", wantConversion[int16](int16(math.MinInt16)))
		t.Run("ToInt32", wantConversion[int32](int16(0)))
		t.Run("ToInt32MaxInt16", wantConversion[int32](int16(math.MaxInt16)))
		t.Run("ToInt32MinInt16", wantConversion[int32](int16(math.MinInt16)))
		t.Run("ToInt64", wantConversion[int64](int16(0)))
		t.Run("ToInt64MaxInt16", wantConversion[int64](int16(math.MaxInt16)))
		t.Run("ToInt64MinInt16", wantConversion[int64](int16(math.MinInt16)))
		t.Run("ToInt", wantConversion[int](int16(0)))
		t.Run("ToIntMaxInt16", wantConversion[int](int16(math.MaxInt16)))
		t.Run("ToIntMinInt16", wantConversion[int](int16(math.MinInt16)))
		t.Run("ToUint8", wantConversion[uint8](int16(0)))
		t.Run("ToUint8Max", wantConversion[uint8](int16(math.MaxUint8)))
		t.Run("ToUint16", wantConversion[uint16](int16(0)))
		t.Run("ToUint16MaxInt16", wantConversion[uint16](int16(math.MaxInt16)))
		t.Run("ToUint32", wantConversion[uint32](int16(0)))
		t.Run("ToUint32MaxInt16", wantConversion[uint32](int16(math.MaxInt16)))
		t.Run("ToUint64", wantConversion[uint64](int16(0)))
		t.Run("ToUint64MaxInt16", wantConversion[uint64](int16(math.MaxInt16)))
		t.Run("ToUint", wantConversion[uint](int16(0)))
		t.Run("ToUintMaxInt16", wantConversion[uint](int16(math.MaxInt16)))
		t.Run("ToUintptr", wantConversion[uintptr](int16(0)))
	})

	t.Run("FromInt32", func(t *testing.T) {
		t.Run("ToInt8", wantConversion[int8](int32(0)))
		t.Run("ToInt8Max", wantConversion[int8](int32(math.MaxInt8)))
		t.Run("ToInt8Min", wantConversion[int8](int32(math.MinInt8)))
		t.Run("ToInt16", wantConversion[int16](int32(0)))
		t.Run("ToInt16Max", wantConversion[int16](int32(math.MaxInt16)))
		t.Run("ToInt16Min", wantConversion[int16](int32(math.MinInt16)))
		t.Run("ToInt32", wantConversion[int32](int32(0)))
		t.Run("ToInt32Max", wantConversion[int32](int32(math.MaxInt32)))
		t.Run("ToInt32Min", wantConversion[int32](int32(math.MinInt32)))
		t.Run("ToInt64", wantConversion[int64](int32(0)))
		t.Run("ToInt64MaxInt32", wantConversion[int64](int32(math.MaxInt32)))
		t.Run("ToInt64MinInt32", wantConversion[int64](int32(math.MinInt32)))
		t.Run("ToInt", wantConversion[int](int32(0)))
		t.Run("ToIntMaxInt32", wantConversion[int](int32(math.MaxInt32)))
		t.Run("ToIntMinInt32", wantConversion[int](int32(math.MinInt32)))
		t.Run("ToUint8", wantConversion[uint8](int32(0)))
		t.Run("ToUint8Max", wantConversion[uint8](int32(math.MaxUint8)))
		t.Run("ToUint16", wantConversion[uint16](int32(0)))
		t.Run("ToUint16Max", wantConversion[uint16](int32(math.MaxUint16)))
		t.Run("ToUint32", wantConversion[uint32](int32(0)))
		t.Run("ToUint32MaxInt32", wantConversion[uint32](int32(math.MaxInt32)))
		t.Run("ToUint64", wantConversion[uint64](int32(0)))
		t.Run("ToUint64MaxInt32", wantConversion[uint64](int32(math.MaxInt32)))
		t.Run("ToUint", wantConversion[uint](int32(0)))
		t.Run("ToUintMaxInt32", wantConversion[uint](int32(math.MaxInt32)))
		t.Run("ToUintptr", wantConversion[uintptr](int32(0)))
	})

	t.Run("FromInt64", func(t *testing.T) {
		t.Run("ToInt8", wantConversion[int8](int64(0)))
		t.Run("ToInt8Max", wantConversion[int8](int64(math.MaxInt8)))
		t.Run("ToInt8Min", wantConversion[int8](int64(math.MinInt8)))
		t.Run("ToInt16", wantConversion[int16](int64(0)))
		t.Run("ToInt16Max", wantConversion[int16](int64(math.MaxInt16)))
		t.Run("ToInt16Min", wantConversion[int16](int64(math.MinInt16)))
		t.Run("ToInt32", wantConversion[int32](int64(0)))
		t.Run("ToInt32Max", wantConversion[int32](int64(math.MaxInt32)))
		t.Run("ToInt32Min", wantConversion[int32](int64(math.MinInt32)))
		t.Run("ToInt64", wantConversion[int64](int64(0)))
		t.Run("ToInt64Max", wantConversion[int64](int64(math.MaxInt64)))
		t.Run("ToInt64Min", wantConversion[int64](int64(math.MinInt64)))
		t.Run("ToInt", wantConversion[int](int64(0)))
		t.Run("ToIntMaxInt", wantConversion[int](int64(math.MaxInt64)))
		t.Run("ToIntMinInt", wantConversion[int](int64(math.MinInt64)))
		t.Run("ToUint8", wantConversion[uint8](int64(0)))
		t.Run("ToUint8Max", wantConversion[uint8](int64(math.MaxUint8)))
		t.Run("ToUint16", wantConversion[uint16](int64(0)))
		t.Run("ToUint16Max", wantConversion[uint16](int64(math.MaxUint16)))
		t.Run("ToUint32", wantConversion[uint32](int64(0)))
		t.Run("ToUint32Max", wantConversion[uint32](int64(math.MaxUint32)))
		t.Run("ToUint64", wantConversion[uint64](int64(0)))
		t.Run("ToUint64Max", wantConversion[uint64](int64(math.MaxInt64)))
		t.Run("ToUint", wantConversion[uint](int64(0)))
		t.Run("ToUintMaxInt", wantConversion[uint](int64(math.MaxInt64)))
		t.Run("ToUintptr", wantConversion[uintptr](int64(0)))
	})
}

func TestToInt_ValueFitsIntoReceiver_ReturnsValueAndTrue(t *testing.T) {
	t.Run("FromInt", wantConversionFunc(narrow.ToInt[int], 0))
	t.Run("FromIntMax", wantConversionFunc(narrow.ToInt[int], math.MaxInt))
	t.Run("FromIntMin", wantConversionFunc(narrow.ToInt[int], math.MinInt))
	t.Run("FromInt8", wantConversionFunc(narrow.ToInt[int8], 0))
	t.Run("FromInt8Min", wantConversionFunc(narrow.ToInt[int8], math.MinInt8))
	t.Run("FromInt8Max", wantConversionFunc(narrow.ToInt[int8], math.MaxInt8))
	t.Run("FromInt16", wantConversionFunc(narrow.ToInt[int16], 0))
	t.Run("FromInt16Min", wantConversionFunc(narrow.ToInt[int16], math.MinInt16))
	t.Run("FromInt16Max", wantConversionFunc(narrow.ToInt[int16], math.MaxInt16))
	t.Run("FromInt32", wantConversionFunc(narrow.ToInt[int32], 0))
	t.Run("FromInt32Min", wantConversionFunc(narrow.ToInt[int32], math.MinInt32))
	t.Run("FromInt32Max", wantConversionFunc(narrow.ToInt[int32], math.MaxInt32))
	t.Run("FromInt64", wantConversionFunc(narrow.ToInt[int64], 0))
	t.Run("FromInt64MaxInt", wantConversionFunc(narrow.ToInt[int64], math.MaxInt))
	t.Run("FromInt64MinInt", wantConversionFunc(narrow.ToInt[int64], math.MinInt))

	t.Run("FromUint", wantConversionFunc(narrow.ToInt[uint], 0))
	t.Run("FromUintMaxInt", wantConversionFunc(narrow.ToInt[uint], math.MaxInt))
	t.Run("FromUint8", wantConversionFunc(narrow.ToInt[uint8], 0))
	t.Run("FromUint8MaxInt8", wantConversionFunc(narrow.ToInt[uint8], math.MaxInt8))
	t.Run("FromUint16", wantConversionFunc(narrow.ToInt[uint16], 0))
	t.Run("FromUint16MaxInt16", wantConversionFunc(narrow.ToInt[uint16], math.MaxInt16))
	t.Run("FromUint32", wantConversionFunc(narrow.ToInt[uint32], 0))
	t.Run("FromUint32MaxInt32", wantConversionFunc(narrow.ToInt[uint32], math.MaxInt32))
	t.Run("FromUint64", wantConversionFunc(narrow.ToInt[uint64], 0))
	t.Run("FromUint64MaxInt", wantConversionFunc(narrow.ToInt[uint64], math.MaxInt))
	t.Run("FromUintptr", wantConversionFunc(narrow.ToInt[uintptr], 0))
}

func TestToInt8_ValueFitsIntoReceiver_ReturnsValueAndTrue(t *testing.T) {
	t.Run("FromInt", wantConversionFunc(narrow.ToInt8[int], 0))
	t.Run("FromIntMaxInt8", wantConversionFunc(narrow.ToInt8[int], math.MaxInt8))
	t.Run("FromIntMinInt8", wantConversionFunc(narrow.ToInt8[int], math.MinInt8))
	t.Run("FromInt8", wantConversionFunc(narrow.ToInt8[int8], 0))
	t.Run("FromInt8Min", wantConversionFunc(narrow.ToInt8[int8], math.MinInt8))
	t.Run("FromInt8Max", wantConversionFunc(narrow.ToInt8[int8], math.MaxInt8))
	t.Run("FromInt16", wantConversionFunc(narrow.ToInt8[int16], 0))
	t.Run("FromInt16MinInt8", wantConversionFunc(narrow.ToInt8[int16], math.MinInt8))
	t.Run("FromInt16MaxInt8", wantConversionFunc(narrow.ToInt8[int16], math.MaxInt8))
	t.Run("FromInt32", wantConversionFunc(narrow.ToInt8[int32], 0))
	t.Run("FromInt32MinInt8", wantConversionFunc(narrow.ToInt8[int32], math.MinInt8))
	t.Run("FromInt32MaxInt8", wantConversionFunc(narrow.ToInt8[int32], math.MaxInt8))
	t.Run("FromInt64", wantConversionFunc(narrow.ToInt8[int64], 0))
	t.Run("FromInt64MaxInt8", wantConversionFunc(narrow.ToInt8[int64], math.MaxInt8))
	t.Run("FromInt64MinInt8", wantConversionFunc(narrow.ToInt8[int64], math.MinInt8))

	t.Run("FromUint", wantConversionFunc(narrow.ToInt8[uint], 0))
	t.Run("FromUintMaxInt8", wantConversionFunc(narrow.ToInt8[uint], math.MaxInt8))
	t.Run("FromUint8", wantConversionFunc(narrow.ToInt8[uint8], 0))
	t.Run("FromUint8Max", wantConversionFunc(narrow.ToInt8[uint8], math.MaxInt8))
	t.Run("FromUint16", wantConversionFunc(narrow.ToInt8[uint16], 0))
	t.Run("FromUint16MaxInt8", wantConversionFunc(narrow.ToInt8[uint16], math.MaxInt8))
	t.Run("FromUint32", wantConversionFunc(narrow.ToInt8[uint32], 0))
	t.Run("FromUint32MaxInt8", wantConversionFunc(narrow.ToInt8[uint32], math.MaxInt8))
	t.Run("FromUint64", wantConversionFunc(narrow.ToInt8[uint64], 0))
	t.Run("FromUint64MaxInt8", wantConversionFunc(narrow.ToInt8[uint64], math.MaxInt8))
	t.Run("FromUintptr", wantConversionFunc(narrow.ToInt8[uintptr], 0))
}

func TestToInt16_ValueFitsIntoReceiver_ReturnsValueAndTrue(t *testing.T) {
	t.Run("FromInt", wantConversionFunc(narrow.ToInt16[int], 0))
	t.Run("FromIntMaxInt16", wantConversionFunc(narrow.ToInt16[int], math.MaxInt16))
	t.Run("FromIntMinInt16", wantConversionFunc(narrow.ToInt16[int], math.MinInt16))
	t.Run("FromInt8", wantConversionFunc(narrow.ToInt16[int8], 0))
	t.Run("FromInt8Min", wantConversionFunc(narrow.ToInt16[int8], math.MinInt8))
	t.Run("FromInt8Max", wantConversionFunc(narrow.ToInt16[int8], math.MaxInt8))
	t.Run("FromInt16", wantConversionFunc(narrow.ToInt16[int16], 0))
	t.Run("FromInt16Min", wantConversionFunc(narrow.ToInt16[int16], math.MinInt16))
	t.Run("FromInt16Max", wantConversionFunc(narrow.ToInt16[int16], math.MaxInt16))
	t.Run("FromInt32", wantConversionFunc(narrow.ToInt16[int32], 0))
	t.Run("FromInt32MinInt16", wantConversionFunc(narrow.ToInt16[int32], math.MinInt16))
	t.Run("FromInt32MaxInt16", wantConversionFunc(narrow.ToInt16[int32], math.MaxInt16))
	t.Run("FromInt64", wantConversionFunc(narrow.ToInt16[int64], 0))
	t.Run("FromInt64MaxInt16", wantConversionFunc(narrow.ToInt16[int64], math.MaxInt16))
	t.Run("FromInt64MinInt16", wantConversionFunc(narrow.ToInt16[int64], math.MinInt16))

	t.Run("FromUint", wantConversionFunc(narrow.ToInt16[uint], 0))
	t.Run("FromUintMaxInt16", wantConversionFunc(narrow.ToInt16[uint], math.MaxInt16))
	t.Run("FromUint8", wantConversionFunc(narrow.ToInt16[uint8], 0))
	t.Run("FromUint8Max", wantConversionFunc(narrow.ToInt16[uint8], math.MaxInt8))
	t.Run("FromUint16", wantConversionFunc(narrow.ToInt16[uint16], 0))
	t.Run("FromUint16MaxInt16", wantConversionFunc(narrow.ToInt16[uint16], math.MaxInt16))
	t.Run("FromUint32", wantConversionFunc(narrow.ToInt16[uint32], 0))
	t.Run("FromUint32MaxInt16", wantConversionFunc(narrow.ToInt16[uint32], math.MaxInt16))
	t.Run("FromUint64", wantConversionFunc(narrow.ToInt16[uint64], 0))
	t.Run("FromUint64MaxInt16", wantConversionFunc(narrow.ToInt16[uint64], math.MaxInt16))
	t.Run("FromUintptr", wantConversionFunc(narrow.ToInt16[uintptr], 0))
}

func TestToInt32_ValueFitsIntoReceiver_ReturnsValueAndTrue(t *testing.T) {
	t.Run("FromInt", wantConversionFunc(narrow.ToInt32[int], 0))
	t.Run("FromIntMaxInt32", wantConversionFunc(narrow.ToInt32[int], math.MaxInt32))
	t.Run("FromIntMinInt32", wantConversionFunc(narrow.ToInt32[int], math.MinInt32))
	t.Run("FromInt8", wantConversionFunc(narrow.ToInt32[int8], 0))
	t.Run("FromInt8Min", wantConversionFunc(narrow.ToInt32[int8], math.MinInt8))
	t.Run("FromInt8Max", wantConversionFunc(narrow.ToInt32[int8], math.MaxInt8))
	t.Run("FromInt16", wantConversionFunc(narrow.ToInt32[int16], 0))
	t.Run("FromInt16Min", wantConversionFunc(narrow.ToInt32[int16], math.MinInt16))
	t.Run("FromInt16Max", wantConversionFunc(narrow.ToInt32[int16], math.MaxInt16))
	t.Run("FromInt32", wantConversionFunc(narrow.ToInt32[int32], 0))
	t.Run("FromInt32Min", wantConversionFunc(narrow.ToInt32[int32], math.MinInt32))
	t.Run("FromInt32Max", wantConversionFunc(narrow.ToInt32[int32], math.MaxInt32))
	t.Run("FromInt64", wantConversionFunc(narrow.ToInt32[int64], 0))
	t.Run("FromInt64MaxInt32", wantConversionFunc(narrow.ToInt32[int64], math.MaxInt32))
	t.Run("FromInt64MinInt32", wantConversionFunc(narrow.ToInt32[int64], math.MinInt32))

	t.Run("FromUint", wantConversionFunc(narrow.ToInt32[uint], 0))
	t.Run("FromUintMaxInt32", wantConversionFunc(narrow.ToInt32[uint], math.MaxInt32))
	t.Run("FromUint8", wantConversionFunc(narrow.ToInt32[uint8], 0))
	t.Run("FromUint8MaxInt8", wantConversionFunc(narrow.ToInt32[uint8], math.MaxInt8))
	t.Run("FromUint16", wantConversionFunc(narrow.ToInt32[uint16], 0))
	t.Run("FromUint16MaxInt16", wantConversionFunc(narrow.ToInt32[uint16], math.MaxInt16))
	t.Run("FromUint32", wantConversionFunc(narrow.ToInt32[uint32], 0))
	t.Run("FromUint32MaxInt32", wantConversionFunc(narrow.ToInt32[uint32], math.MaxInt32))
	t.Run("FromUint64", wantConversionFunc(narrow.ToInt32[uint64], 0))
	t.Run("FromUint64MaxInt32", wantConversionFunc(narrow.ToInt32[uint64], math.MaxInt32))
	t.Run("FromUintptr", wantConversionFunc(narrow.ToInt32[uintptr], 0))
}

func TestToInt64_ValueFitsIntoReceiver_ReturnsValueAndTrue(t *testing.T) {
	t.Run("FromInt", wantConversionFunc(narrow.ToInt64[int], 0))
	t.Run("FromIntMaxInt32", wantConversionFunc(narrow.ToInt64[int], math.MaxInt32))
	t.Run("FromIntMin", wantConversionFunc(narrow.ToInt64[int], math.MinInt32))
	t.Run("FromInt8", wantConversionFunc(narrow.ToInt64[int8], 0))
	t.Run("FromInt8Min", wantConversionFunc(narrow.ToInt64[int8], math.MinInt8))
	t.Run("FromInt8Max", wantConversionFunc(narrow.ToInt64[int8], math.MaxInt8))
	t.Run("FromInt16", wantConversionFunc(narrow.ToInt64[int16], 0))
	t.Run("FromInt16Min", wantConversionFunc(narrow.ToInt64[int16], math.MinInt16))
	t.Run("FromInt16Max", wantConversionFunc(narrow.ToInt64[int16], math.MaxInt16))
	t.Run("FromInt32", wantConversionFunc(narrow.ToInt64[int32], 0))
	t.Run("FromInt32Min", wantConversionFunc(narrow.ToInt64[int32], math.MinInt32))
	t.Run("FromInt32Max", wantConversionFunc(narrow.ToInt64[int32], math.MaxInt32))
	t.Run("FromInt64", wantConversionFunc(narrow.ToInt64[int64], 0))
	t.Run("FromInt64Max", wantConversionFunc(narrow.ToInt64[int64], math.MaxInt64))
	t.Run("FromInt64Min", wantConversionFunc(narrow.ToInt64[int64], math.MinInt64))

	t.Run("FromUint", wantConversionFunc(narrow.ToInt64[uint], 0))
	t.Run("FromUintMaxUint32", wantConversionFunc(narrow.ToInt64[uint], math.MaxUint32))
	t.Run("FromUint8", wantConversionFunc(narrow.ToInt64[uint8], 0))
	t.Run("FromUint8Max", wantConversionFunc(narrow.ToInt64[uint8], math.MaxUint8))
	t.Run("FromUint16", wantConversionFunc(narrow.ToInt64[uint16], 0))
	t.Run("FromUint16Max", wantConversionFunc(narrow.ToInt64[uint16], math.MaxUint16))
	t.Run("FromUint32", wantConversionFunc(narrow.ToInt64[uint32], 0))
	t.Run("FromUint32Max", wantConversionFunc(narrow.ToInt64[uint32], math.MaxUint32))
	t.Run("FromUint64", wantConversionFunc(narrow.ToInt64[uint64], 0))
	t.Run("FromUint64Max", wantConversionFunc(narrow.ToInt64[uint64], math.MaxInt64))
	t.Run("FromUintptr", wantConversionFunc(narrow.ToInt64[uintptr], 0))
}

func TestToUint_ValueFitsIntoReceiver_ReturnsValueAndTrue(t *testing.T) {
	t.Run("FromInt", wantConversionFunc(narrow.ToUint[int], 0))
	t.Run("FromIntMaxInt32", wantConversionFunc(narrow.ToUint[int], math.MaxInt32))
	t.Run("FromInt8", wantConversionFunc(narrow.ToUint[int8], 0))
	t.Run("FromInt8Max", wantConversionFunc(narrow.ToUint[int8], math.MaxInt8))
	t.Run("FromInt16", wantConversionFunc(narrow.ToUint[int16], 0))
	t.Run("FromInt16Max", wantConversionFunc(narrow.ToUint[int16], math.MaxInt16))
	t.Run("FromInt32", wantConversionFunc(narrow.ToUint[int32], 0))
	t.Run("FromInt32Max", wantConversionFunc(narrow.ToUint[int32], math.MaxInt32))
	t.Run("FromInt64", wantConversionFunc(narrow.ToUint[int64], 0))
	t.Run("FromInt64MaxUint32", wantConversionFunc(narrow.ToUint[int64], math.MaxInt32))

	t.Run("FromUint", wantConversionFunc(narrow.ToUint[uint], 0))
	t.Run("FromUintMaxUint32", wantConversionFunc(narrow.ToUint[uint], math.MaxUint32))
	t.Run("FromUint8", wantConversionFunc(narrow.ToUint[uint8], 0))
	t.Run("FromUint8Max", wantConversionFunc(narrow.ToUint[uint8], math.MaxUint8))
	t.Run("FromUint16", wantConversionFunc(narrow.ToUint[uint16], 0))
	t.Run("FromUint16Max", wantConversionFunc(narrow.ToUint[uint16], math.MaxUint16))
	t.Run("FromUint32", wantConversionFunc(narrow.ToUint[uint32], 0))
	t.Run("FromUint32Max", wantConversionFunc(narrow.ToUint[uint32], math.MaxUint32))
	t.Run("FromUint64", wantConversionFunc(narrow.ToUint[uint64], 0))
	t.Run("FromUint64MaxUint32", wantConversionFunc(narrow.ToUint[uint64], math.MaxUint32))
	t.Run("FromUintptr", wantConversionFunc(narrow.ToUint[uintptr], 0))
}

func TestToUint8_ValueFitsIntoReceiver_ReturnsValueAndTrue(t *testing.T) {
	t.Run("FromInt", wantConversionFunc(narrow.ToUint8[int], 0))
	t.Run("FromIntMaxUint8", wantConversionFunc(narrow.ToUint8[int], math.MaxUint8))
	t.Run("FromInt8", wantConversionFunc(narrow.ToUint8[int8], 0))
	t.Run("FromInt8Max", wantConversionFunc(narrow.ToUint8[int8], math.MaxInt8))
	t.Run("FromInt16", wantConversionFunc(narrow.ToUint8[int16], 0))
	t.Run("FromInt16MaxUint8", wantConversionFunc(narrow.ToUint8[int16], math.MaxInt8))
	t.Run("FromInt32", wantConversionFunc(narrow.ToUint8[int32], 0))
	t.Run("FromInt32MaxUint8", wantConversionFunc(narrow.ToUint8[int32], math.MaxUint8))
	t.Run("FromInt64", wantConversionFunc(narrow.ToUint8[int64], 0))
	t.Run("FromInt64MaxUint8", wantConversionFunc(narrow.ToUint8[int64], math.MaxUint8))

	t.Run("FromUint", wantConversionFunc(narrow.ToUint8[uint], 0))
	t.Run("FromUintMaxUint8", wantConversionFunc(narrow.ToUint8[uint], math.MaxUint8))
	t.Run("FromUint8", wantConversionFunc(narrow.ToUint8[uint8], 0))
	t.Run("FromUint8Max", wantConversionFunc(narrow.ToUint8[uint8], math.MaxUint8))
	t.Run("FromUint16", wantConversionFunc(narrow.ToUint8[uint16], 0))
	t.Run("FromUint16MaxUint8", wantConversionFunc(narrow.ToUint8[uint16], math.MaxUint8))
	t.Run("FromUint32", wantConversionFunc(narrow.ToUint8[uint32], 0))
	t.Run("FromUint32MaxUint8", wantConversionFunc(narrow.ToUint8[uint32], math.MaxUint8))
	t.Run("FromUint64", wantConversionFunc(narrow.ToUint8[uint64], 0))
	t.Run("FromUint64MaxUint8", wantConversionFunc(narrow.ToUint8[uint64], math.MaxUint8))
	t.Run("FromUintptr", wantConversionFunc(narrow.ToUint8[uintptr], 0))
}

func TestToUint16_ValueFitsIntoReceiver_ReturnsValueAndTrue(t *testing.T) {
	t.Run("FromInt", wantConversionFunc(narrow.ToUint16[int], 0))
	t.Run("FromIntMaxUint16", wantConversionFunc(narrow.ToUint16[int], math.MaxUint16))
	t.Run("FromInt8", wantConversionFunc(narrow.ToUint16[int8], 0))
	t.Run("FromInt8Max", wantConversionFunc(narrow.ToUint16[int8], math.MaxInt8))
	t.Run("FromInt16", wantConversionFunc(narrow.ToUint16[int16], 0))
	t.Run("FromInt16Max", wantConversionFunc(narrow.ToUint16[int16], math.MaxInt16))
	t.Run("FromInt32", wantConversionFunc(narrow.ToUint16[int32], 0))
	t.Run("FromInt32MaxUint16", wantConversionFunc(narrow.ToUint16[int32], math.MaxUint16))
	t.Run("FromInt64", wantConversionFunc(narrow.ToUint16[int64], 0))
	t.Run("FromInt64MaxUint16", wantConversionFunc(narrow.ToUint16[int64], math.MaxUint16))

	t.Run("FromUint", wantConversionFunc(narrow.ToUint16[uint], 0))
	t.Run("FromUintMaxUint16", wantConversionFunc(narrow.ToUint16[uint], math.MaxUint16))
	t.Run("FromUint8", wantConversionFunc(narrow.ToUint16[uint8], 0))
	t.Run("FromUint8Max", wantConversionFunc(narrow.ToUint16[uint8], math.MaxUint8))
	t.Run("FromUint16", wantConversionFunc(narrow.ToUint16[uint16], 0))
	t.Run("FromUint16Max", wantConversionFunc(narrow.ToUint16[uint16], math.MaxUint16))
	t.Run("FromUint32", wantConversionFunc(narrow.ToUint16[uint32], 0))
	t.Run("FromUint32MaxUint16", wantConversionFunc(narrow.ToUint16[uint32], math.MaxUint16))
	t.Run("FromUint64", wantConversionFunc(narrow.ToUint16[uint64], 0))
	t.Run("FromUint64MaxUint16", wantConversionFunc(narrow.ToUint16[uint64], math.MaxUint16))
	t.Run("FromUintptr", wantConversionFunc(narrow.ToUint16[uintptr], 0))
}

func TestToUint32_ValueFitsIntoReceiver_ReturnsValueAndTrue(t *testing.T) {
	t.Run("FromInt", wantConversionFunc(narrow.ToUint32[int], 0))
	t.Run("FromIntMaxInt32", wantConversionFunc(narrow.ToUint32[int], math.MaxInt32))
	t.Run("FromInt8", wantConversionFunc(narrow.ToUint32[int8], 0))
	t.Run("FromInt8Max", wantConversionFunc(narrow.ToUint32[int8], math.MaxInt8))
	t.Run("FromInt16", wantConversionFunc(narrow.ToUint32[int16], 0))
	t.Run("FromInt16Max", wantConversionFunc(narrow.ToUint32[int16], math.MaxInt16))
	t.Run("FromInt32", wantConversionFunc(narrow.ToUint32[int32], 0))
	t.Run("FromInt32Max", wantConversionFunc(narrow.ToUint32[int32], math.MaxInt32))
	t.Run("FromInt64", wantConversionFunc(narrow.ToUint32[int64], 0))
	t.Run("FromInt64MaxUint32", wantConversionFunc(narrow.ToUint32[int64], math.MaxInt32))

	t.Run("FromUint", wantConversionFunc(narrow.ToUint32[uint], 0))
	t.Run("FromUintMaxUint32", wantConversionFunc(narrow.ToUint32[uint], math.MaxUint32))
	t.Run("FromUint8", wantConversionFunc(narrow.ToUint32[uint8], 0))
	t.Run("FromUint8Max", wantConversionFunc(narrow.ToUint32[uint8], math.MaxUint8))
	t.Run("FromUint16", wantConversionFunc(narrow.ToUint32[uint16], 0))
	t.Run("FromUint16Max", wantConversionFunc(narrow.ToUint32[uint16], math.MaxUint16))
	t.Run("FromUint32", wantConversionFunc(narrow.ToUint32[uint32], 0))
	t.Run("FromUint32Max", wantConversionFunc(narrow.ToUint32[uint32], math.MaxUint32))
	t.Run("FromUint64", wantConversionFunc(narrow.ToUint32[uint64], 0))
	t.Run("FromUint64MaxUint32", wantConversionFunc(narrow.ToUint32[uint64], math.MaxUint32))
	t.Run("FromUintptr", wantConversionFunc(narrow.ToUint32[uintptr], 0))
}

func TestToUint64_ValueFitsIntoReceiver_ReturnsValueAndTrue(t *testing.T) {
	t.Run("FromInt", wantConversionFunc(narrow.ToUint64[int], 0))
	t.Run("FromIntMaxInt32", wantConversionFunc(narrow.ToUint64[int], math.MaxInt32))
	t.Run("FromInt8", wantConversionFunc(narrow.ToUint64[int8], 0))
	t.Run("FromInt8Max", wantConversionFunc(narrow.ToUint64[int8], math.MaxInt8))
	t.Run("FromInt16", wantConversionFunc(narrow.ToUint64[int16], 0))
	t.Run("FromInt16Max", wantConversionFunc(narrow.ToUint64[int16], math.MaxInt16))
	t.Run("FromInt32", wantConversionFunc(narrow.ToUint64[int32], 0))
	t.Run("FromInt32Max", wantConversionFunc(narrow.ToUint64[int32], math.MaxInt32))
	t.Run("FromInt64", wantConversionFunc(narrow.ToUint64[int64], 0))
	t.Run("FromInt64Max", wantConversionFunc(narrow.ToUint64[int64], math.MaxInt64))

	t.Run("FromUint", wantConversionFunc(narrow.ToUint64[uint], 0))
	t.Run("FromUintMaxUint32", wantConversionFunc(narrow.ToUint64[uint], math.MaxUint32))
	t.Run("FromUint8", wantConversionFunc(narrow.ToUint64[uint8], 0))
	t.Run("FromUint8Max", wantConversionFunc(narrow.ToUint64[uint8], math.MaxUint8))
	t.Run("FromUint16", wantConversionFunc(narrow.ToUint64[uint16], 0))
	t.Run("FromUint16Max", wantConversionFunc(narrow.ToUint64[uint16], math.MaxUint16))
	t.Run("FromUint32", wantConversionFunc(narrow.ToUint64[uint32], 0))
	t.Run("FromUint32Max", wantConversionFunc(narrow.ToUint64[uint32], math.MaxUint32))
	t.Run("FromUint64", wantConversionFunc(narrow.ToUint64[uint64], 0))
	t.Run("FromUint64Max", wantConversionFunc(narrow.ToUint64[uint64], math.MaxUint64))
	t.Run("FromUintptr", wantConversionFunc(narrow.ToUint64[uintptr], 0))
}
