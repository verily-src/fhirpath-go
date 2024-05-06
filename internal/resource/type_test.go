package resource_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/internal/fhirtest"
	"github.com/verily-src/fhirpath-go/internal/resource"
)

func TestTypeOf_ReturnsType(t *testing.T) {
	for name, res := range fhirtest.Resources {
		t.Run(name, func(t *testing.T) {
			want := string(res.ProtoReflect().Descriptor().Name())

			got := resource.TypeOf(res)

			if !cmp.Equal(string(got), want) {
				t.Errorf("TypeOf(%v): got '%v', want '%v'", name, got, want)
			}
		})
	}
}

func TestTypeOf_NilInput_Panics(t *testing.T) {
	defer func() { _ = recover() }()

	resource.TypeOf(nil)

	t.Errorf("TypeOf: expected panic")
}

func TestNewType_ValidTypeName_ReturnsType(t *testing.T) {
	for name, res := range fhirtest.Resources {
		t.Run(name, func(t *testing.T) {
			want := resource.TypeOf(res)

			got, err := resource.NewType(name)
			if err != nil {
				t.Fatalf("NewType: got unexpected err '%v' from NewType", err)
			}

			if !cmp.Equal(got, want) {
				t.Errorf("NewType(%v): got %v, want %v", name, got, want)
			}
		})
	}
}

func TestNewType_InvalidTypeName_ReturnsErrBadType(t *testing.T) {
	testCases := []struct {
		name  string
		value string
	}{
		{"Empty", ""},
		{"NotAnElement", "Bad-Element"},
		{"AnonymousElement", "Bundle_Entry"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := resource.NewType(tc.value)

			if got, want := err, resource.ErrBadType; !errors.Is(got, want) {
				t.Errorf("NewType(%v): got %v, want %v", tc.name, got, want)
			}
		})
	}
}

func TestIsType_ValidTypeName_ReturnsTrue(t *testing.T) {
	for name := range fhirtest.Resources {
		t.Run(name, func(t *testing.T) {
			got := resource.IsType(name)

			if got != true {
				t.Errorf("IsType(%v): got %v, want true", name, got)
			}
		})
	}
}

func TestIsType_InvalidTypeName_ReturnsFalse(t *testing.T) {
	testCases := []struct {
		name  string
		value string
	}{
		{"Empty", ""},
		{"NotAResource", "ContainedResource"},
		{"Element", "String"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := resource.IsType(tc.value)

			if got != false {
				t.Errorf("IsType(%v): got %v, want false", tc.name, got)
			}
		})
	}
}

func TestTypeNew_ReturnsElementOfType(t *testing.T) {
	for name, elem := range fhirtest.Resources {
		t.Run(name, func(t *testing.T) {
			ty := resource.TypeOf(elem)
			want := reflect.TypeOf(elem)

			got := ty.New()

			if reflect.TypeOf(got) != want {
				t.Errorf("Type.New: got %v, want %v", got, want)
			}
		})
	}
}

func TestTypeNew_Unspecified_ReturnsNil(t *testing.T) {
	defer func() { _ = recover() }()

	resource.Type("").New()

	t.Errorf("Type.New: expected panic")
}
