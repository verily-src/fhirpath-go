package fhirpath_test

import (
	"errors"
	"math"
	"testing"

	"github.com/verily-src/fhirpath-go/fhirpath"
	"github.com/verily-src/fhirpath-go/fhirpath/fhirpathtest"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

func TestExpressionString(t *testing.T) {
	const want = "Patient.name"
	expr, err := fhirpath.Compile(want)
	if err != nil {
		t.Fatalf("Expression.String(): got unexpected err: %v", err)
	}

	got := expr.String()

	if got != want {
		t.Errorf("Expression.String(): got %v, want %v", got, want)
	}
}

func TestExpressionEvaluateAsBool_EvaluationError_ReturnsError(t *testing.T) {
	want := errors.New("some error")
	path := fhirpathtest.Error(want)

	_, err := path.EvaluateAsBool(nil)

	if got, want := err, want; !errors.Is(got, want) {
		t.Errorf("EvaluateAsBool: want err %v, got %v", want, got)
	}
}

func TestExpressionEvaluateAsBool_NonConvertibleResult_ReturnsError(t *testing.T) {
	testCases := []struct {
		name  string
		input system.Collection
	}{
		{
			name:  "Collection of more than 1",
			input: system.Collection{system.String("1"), system.String("2")},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := fhirpathtest.ReturnCollection(tc.input)

			_, err := path.EvaluateAsBool(nil)

			if err == nil {
				t.Fatalf("EvaluateAsBool: Expected error")
			}
		})
	}
}

func TestExpressionEvaluateAsBoo_EmptyCollection_ReturnsFalse(t *testing.T) {
	testCases := []struct {
		name  string
		input system.Collection
	}{
		{
			name:  "Empty Collection",
			input: system.Collection{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := fhirpathtest.ReturnCollection(tc.input)

			result, err := path.EvaluateAsBool(nil)

			if err != nil {
				t.Fatalf("EvaluateAsBool: Expected no error")
			}

			if result {
				t.Fatalf("EvaluateAsBool: Expected false")
			}
		})
	}
}

func TestExpressionEvaluateAsBool_ConvertibleResult_ReturnsBool(t *testing.T) {
	testCases := []struct {
		name  string
		input any
		want  bool
	}{
		{
			name:  "system boolean",
			input: system.Boolean(true),
			want:  true,
		}, {
			name:  "FHIR boolean",
			input: fhir.Boolean(true),
			want:  true,
		}, {
			name:  "Random singleton type",
			input: system.String("Hello world"),
			want:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := fhirpathtest.Return(tc.input)

			got, err := path.EvaluateAsBool(nil)
			if err != nil {
				t.Fatalf("EvaluateAsBool: Unexpected error %v", err)
			}

			if got != tc.want {
				t.Errorf("EvaluateAsBool: want %v, got %v", tc.want, got)
			}
		})
	}
}

func TestExpressionEvaluateAsString_EvaluationError_ReturnsError(t *testing.T) {
	want := errors.New("some error")
	path := fhirpathtest.Error(want)

	_, err := path.EvaluateAsString(nil)

	if got, want := err, want; !errors.Is(got, want) {
		t.Errorf("EvaluateAsString: want err %v, got %v", want, got)
	}
}

func TestExpressionEvaluateAsString_NonConvertibleResult_ReturnsError(t *testing.T) {
	testCases := []struct {
		name  string
		input system.Collection
	}{
		{
			name:  "Empty Collection",
			input: system.Collection{},
		}, {
			name:  "Collection of more than 1",
			input: system.Collection{system.String("1"), system.String("2")},
		}, {
			name:  "Invalid type",
			input: system.Collection{fhir.Integer(42)},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := fhirpathtest.ReturnCollection(tc.input)

			_, err := path.EvaluateAsString(nil)

			if err == nil {
				t.Fatalf("EvaluateAsString: Expected error")
			}
		})
	}
}

func TestExpressionEvaluateAsString_ConvertibleResult_ReturnsString(t *testing.T) {
	const str = "hello world"
	testCases := []struct {
		name  string
		input any
		want  string
	}{
		{
			name:  "system String",
			input: system.String(str),
			want:  str,
		}, {
			name:  "FHIR String",
			input: fhir.String(str),
			want:  str,
		}, {
			name:  "FHIR Code",
			input: fhir.Code(str),
			want:  str,
		}, {
			name:  "FHIR ID",
			input: fhir.ID(str),
			want:  str,
		}, {
			name:  "FHIR Markdown",
			input: fhir.Markdown(str),
			want:  str,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := fhirpathtest.Return(tc.input)

			got, err := path.EvaluateAsString(nil)
			if err != nil {
				t.Fatalf("EvaluateAsString: Unexpected error %v", err)
			}

			if got != tc.want {
				t.Errorf("EvaluateAsString: want %v, got %v", tc.want, got)
			}
		})
	}
}

func TestExpressionEvaluateAsInt32_EvaluationError_ReturnsError(t *testing.T) {
	want := errors.New("some error")
	path := fhirpathtest.Error(want)

	_, err := path.EvaluateAsInt32(nil)

	if got, want := err, want; !errors.Is(got, want) {
		t.Errorf("EvaluateAsInt32: want err %v, got %v", want, got)
	}
}

func TestExpressionEvaluateAsInt32_NonConvertibleResult_ReturnsError(t *testing.T) {
	testCases := []struct {
		name  string
		input system.Collection
	}{
		{
			name:  "Empty Collection",
			input: system.Collection{},
		}, {
			name:  "Collection of more than 1",
			input: system.Collection{system.Integer(1), system.Integer(2)},
		}, {
			name:  "Integer not representable",
			input: system.Collection{fhir.UnsignedInt(math.MaxUint32)},
		}, {
			name:  "Invalid Type",
			input: system.Collection{fhir.String("hello world")},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := fhirpathtest.ReturnCollection(tc.input)

			_, err := path.EvaluateAsInt32(nil)

			if err == nil {
				t.Fatalf("EvaluateAsInt32: Expected error")
			}
		})
	}
}

func TestExpressionEvaluateAsInt32_ConvertibleResult_ReturnsInt32(t *testing.T) {
	const val = 42
	testCases := []struct {
		name  string
		input any
		want  int32
	}{
		{
			name:  "system Integer",
			input: system.Integer(val),
			want:  val,
		}, {
			name:  "FHIR Integer",
			input: fhir.Integer(val),
			want:  val,
		}, {
			name:  "FHIR PositiveInt",
			input: fhir.PositiveInt(val),
			want:  val,
		}, {
			name:  "FHIR ID",
			input: fhir.UnsignedInt(val),
			want:  val,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := fhirpathtest.Return(tc.input)

			got, err := path.EvaluateAsInt32(nil)
			if err != nil {
				t.Fatalf("EvaluateAsInt32: Unexpected error %v", err)
			}

			if got != tc.want {
				t.Errorf("EvaluateAsInt32: want %v, got %v", tc.want, got)
			}
		})
	}
}
