package system_test

import (
	"testing"

	"github.com/verily-src/fhirpath-go/internal/fhir"
	"google.golang.org/protobuf/proto"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

func TestEqual_ReturnsResult(t *testing.T) {
	patient1 := &ppb.Patient{
		Id: fhir.ID("123"),
	}
	patient2 := &ppb.Patient{
		Id: fhir.ID("234"),
	}

	testCases := []struct {
		name            string
		leftCollection  system.Collection
		rightCollection system.Collection
		equal           bool
	}{
		{
			name:            "mismatched collection lengths",
			leftCollection:  system.Collection{"one"},
			rightCollection: system.Collection{"one", "two"},
			equal:           false,
		},
		{
			name:            "mismatched types (primitive and complex)",
			leftCollection:  system.Collection{system.String("abc")},
			rightCollection: system.Collection{&ppb.Patient{}},
			equal:           false,
		},
		{
			name:            "system types not equal",
			leftCollection:  system.Collection{system.String("abc")},
			rightCollection: system.Collection{system.String("abcd")},
			equal:           false,
		},
		{
			name:            "proto types not equal",
			leftCollection:  system.Collection{patient1},
			rightCollection: system.Collection{patient2},
			equal:           false,
		},
		{
			name:            "equal system types",
			leftCollection:  system.Collection{system.String("sausage")},
			rightCollection: system.Collection{system.String("sausage")},
			equal:           true,
		},
		{
			name:            "equal proto types",
			leftCollection:  system.Collection{fhir.Code("#blessed")},
			rightCollection: system.Collection{fhir.Code("#blessed")},
			equal:           true,
		},
		{
			name:            "equal complex types",
			leftCollection:  system.Collection{patient1},
			rightCollection: system.Collection{patient1},
			equal:           true,
		},
		{
			name:            "full collection not equal",
			leftCollection:  system.Collection{system.String("Not"), system.String("Equal")},
			rightCollection: system.Collection{system.String("Not"), system.String("Equivalent")},
			equal:           false,
		},
		{
			name:            "equal collection",
			leftCollection:  system.Collection{system.String("Is"), system.String("Equal")},
			rightCollection: system.Collection{system.String("Is"), system.String("Equal")},
			equal:           true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := tc.leftCollection.TryEqual(tc.rightCollection)

			if !ok {
				t.Fatalf("Collection.TryEqual: did not return value when one expected")
			}
			if got, want := got, tc.equal; got != want {
				t.Errorf("Collection.TryEqual returned incorrect result, got: %v, want %v", got, want)
			}
		})
	}
}

func TestToSingletonBoolean_ConvertsToBool(t *testing.T) {
	testCases := []struct {
		name            string
		inputCollection system.Collection
		want            []system.Boolean
		shouldError     bool
	}{
		{
			name:            "returns error on non-singleton collection",
			inputCollection: system.Collection{1, 2},
			shouldError:     true,
		},
		{
			name:            "returns contained proto boolean",
			inputCollection: system.Collection{fhir.Boolean(false)},
			want:            []system.Boolean{false},
		},
		{
			name:            "returns contained system boolean",
			inputCollection: system.Collection{system.Boolean(true)},
			want:            []system.Boolean{true},
		},
		{
			name:            "returns true on singleton non-bool collection",
			inputCollection: system.Collection{system.String("1")},
			want:            []system.Boolean{true},
		},
		{
			name:            "propagates empty collection input",
			inputCollection: system.Collection{},
			want:            []system.Boolean{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.inputCollection.ToSingletonBoolean()

			if gotErr, wantErr := err != nil, tc.shouldError; gotErr != wantErr {
				t.Fatalf("Collection.SingletonBoolean() returned unexpected error result: gotErr %v, wantErr %v", gotErr, wantErr)
			}
			if !cmp.Equal(got, tc.want) {
				t.Errorf("Collection.SingletonBoolean() returned unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCollection_ToFloat64(t *testing.T) {
	tests := []struct {
		name    string
		c       system.Collection
		want    float64
		wantErr bool
	}{
		{
			name:    "errors if input is not a number",
			c:       system.Collection{system.String("10.1")},
			want:    0,
			wantErr: true,
		},
		{
			name:    "errors if input length is more than 1",
			c:       system.Collection{system.MustParseDecimal("10.1"), system.String("10.2")},
			want:    0,
			wantErr: true,
		},
		{
			name:    "converts Decimal into float64 successfully",
			c:       system.Collection{system.MustParseDecimal("10.5")},
			want:    10.5,
			wantErr: false,
		},
		{
			name:    "converts Integer into float64 successfully",
			c:       system.Collection{fhir.Integer(-100)},
			want:    -100,
			wantErr: false,
		},
		{
			name:    "converts PositiveInt into float64 successfully",
			c:       system.Collection{fhir.PositiveInt(100)},
			want:    100,
			wantErr: false,
		},
		{
			name:    "converts UnsignedInt into float64 successfully",
			c:       system.Collection{fhir.UnsignedInt(10)},
			want:    10,
			wantErr: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.c.ToFloat64()
			if (err != nil) != tc.wantErr {
				t.Errorf("ToFloat64() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if got != tc.want {
				t.Errorf("ToFloat64() got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCollection_ToInt32(t *testing.T) {
	tests := []struct {
		name    string
		c       system.Collection
		want    int32
		wantErr bool
	}{
		{
			name:    "errors if input is not a number",
			c:       system.Collection{system.String("10.1")},
			want:    0,
			wantErr: true,
		},
		{
			name:    "errors if input length is more than 1",
			c:       system.Collection{system.Integer(1), system.String("10.2")},
			want:    0,
			wantErr: true,
		},
		{
			name:    "converts positive Integer into int32 successfully",
			c:       system.Collection{system.Integer(100)},
			want:    100,
			wantErr: false,
		},
		{
			name:    "converts negative Integer into int32 successfully",
			c:       system.Collection{fhir.Integer(-100)},
			want:    -100,
			wantErr: false,
		},
		{
			name:    "converts PositiveInt into int32 successfully",
			c:       system.Collection{fhir.PositiveInt(1000)},
			want:    1000,
			wantErr: false,
		},
		{
			name:    "converts UnsignedInt into int32 successfully",
			c:       system.Collection{fhir.UnsignedInt(10)},
			want:    10,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.ToInt32()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToInt32() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	patient := &ppb.Patient{
		Name: []*dtpb.HumanName{
			{
				Given:  fhir.Strings("Foo", "Barl"),
				Family: fhir.String("Bazington"),
			},
		},
	}
	testCases := []struct {
		name     string
		haystack system.Collection
		needle   any
		want     bool
	}{
		{
			name:     "Haystack contains exact match",
			haystack: system.Collection{system.Integer(42), system.String("Hello")},
			needle:   system.String("Hello"),
			want:     true,
		}, {
			name:     "Needle can convert to haystack value",
			haystack: system.Collection{system.Integer(42), system.String("Hello")},
			needle:   fhir.String("Hello"),
			want:     true,
		}, {
			name:     "Haystack can convert to needle value",
			haystack: system.Collection{system.Integer(42), fhir.String("Hello")},
			needle:   system.String("Hello"),
			want:     true,
		}, {
			name:     "Haystack does not contain needle",
			haystack: system.Collection{system.Integer(42), fhir.String("Hello")},
			needle:   system.String("World"),
			want:     false,
		}, {
			name:     "Needle is a proto value in haystack",
			haystack: system.Collection{system.Integer(42), patient, system.String("Foo")},
			needle:   proto.Clone(patient),
			want:     true,
		}, {
			name:     "Needle is a proto value not in haystack",
			haystack: system.Collection{system.Integer(42), patient, system.String("Foo")},
			needle:   &ppb.Patient{},
			want:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.haystack.Contains(tc.needle)

			if got != tc.want {
				t.Errorf("Collection.Contains(%v): got %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}
