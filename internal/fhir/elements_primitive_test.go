package fhir_test

import (
	"math"
	"strings"
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/internal/element/canonical"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/slices"
)

func TestBase64Binary(t *testing.T) {
	want := []byte{0xde, 0xad, 0xbe, 0xef}

	sut := fhir.Base64Binary(want)

	if got := sut.GetValue(); !cmp.Equal(got, want) {
		t.Errorf("Base64Binary: got %v, want %v", got, want)
	}
}

func TestBoolean(t *testing.T) {
	want := true

	sut := fhir.Boolean(want)

	if got := sut.GetValue(); !cmp.Equal(got, want) {
		t.Errorf("Boolean: got %v, want %v", got, want)
	}
}

func TestCode(t *testing.T) {
	want := "value"

	sut := fhir.Code(want)

	if got := sut.GetValue(); !cmp.Equal(got, want) {
		t.Errorf("Code: got %v, want %v", got, want)
	}
}

func TestID(t *testing.T) {
	want := "id"

	sut := fhir.ID(want)

	if got := sut.GetValue(); !cmp.Equal(got, want) {
		t.Errorf("ID: got %v, want %v", got, want)
	}
}

func TestIsID(t *testing.T) {
	testCases := []struct {
		name    string
		inputId string
		wantOk  bool
	}{
		{"empty", "", false},
		{"typical", "NormalId-.123", true},
		{"valid inside bogus", "&&&hello", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotOk := fhir.IsID(tc.inputId)
			if gotOk != tc.wantOk {
				t.Errorf("IsID(%s) ok mismatch: got %v, want %v", tc.name, gotOk, tc.wantOk)
			}
		})
	}
}

func TestInteger(t *testing.T) {
	want := int32(42)

	sut := fhir.Integer(want)

	if got := sut.GetValue(); !cmp.Equal(got, want) {
		t.Errorf("Integer: got %v, want %v", got, want)
	}
}

func TestIntegerFromInt_Truncates_ReturnsError(t *testing.T) {
	input := math.MaxInt64

	_, err := fhir.IntegerFromInt(input)

	if got, want := err, fhir.ErrIntegerDataLoss; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
		t.Errorf("IntegerFromInt: got %v, err %v", got, want)
	}
}

func TestIntegerFromInt_ValidValue_ReturnsInteger(t *testing.T) {
	testCases := []struct {
		name  string
		value int
	}{
		{"Zero", 0},
		{"MaxInt32", math.MaxInt32},
		{"MinInt32", math.MinInt32},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut, err := fhir.IntegerFromInt(tc.value)
			if err != nil {
				t.Fatalf("IntegerFromInt(%v): unexpected error '%v'", tc.name, err)
			}

			if got, want := sut.GetValue(), tc.value; got != int32(want) {
				t.Errorf("IntegerFromInt(%v): got %v, want %v", tc.name, got, want)
			}
		})
	}
}

func TestIntegerFromPositiveInt_Truncates_ReturnsError(t *testing.T) {
	testCases := []struct {
		name  string
		value uint32
	}{
		{"OneOverInt32Max", math.MaxInt32 + 1},
		{"UInt32Max", math.MaxUint32},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := fhir.IntegerFromPositiveInt(fhir.PositiveInt(tc.value))

			if got, want := err, fhir.ErrIntegerDataLoss; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("IntegerFromPositiveInt(%v): got %v, err %v", tc.name, got, want)
			}
		})
	}
}

func TestIntegerFromPositiveInt_ValidValue_ReturnsInteger(t *testing.T) {
	testCases := []struct {
		name  string
		value uint32
	}{
		{"Zero", 0},
		{"MaxInt32", math.MaxInt32},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut, err := fhir.IntegerFromPositiveInt(fhir.PositiveInt(tc.value))
			if err != nil {
				t.Fatalf("IntegerFromPositiveInt(%v): unexpected error '%v'", tc.name, err)
			}

			if got, want := sut.GetValue(), tc.value; got != int32(want) {
				t.Errorf("IntegerFromPositiveInt(%v): got %v, want %v", tc.name, got, want)
			}
		})
	}
}

func TestIntegerFromUnsignedInt_Truncates_ReturnsError(t *testing.T) {
	testCases := []struct {
		name  string
		value uint32
	}{
		{"OneOverInt32Max", math.MaxInt32 + 1},
		{"UInt32Max", math.MaxUint32},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := fhir.IntegerFromUnsignedInt(fhir.UnsignedInt(tc.value))

			if got, want := err, fhir.ErrIntegerDataLoss; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("IntegerFromUnsignedInt(%v): got %v, err %v", tc.name, got, want)
			}
		})
	}
}

func TestIntegerFromUnsignedInt_ValidValue_ReturnsInteger(t *testing.T) {
	testCases := []struct {
		name  string
		value uint32
	}{
		{"Zero", 0},
		{"MaxInt32", math.MaxInt32},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut, err := fhir.IntegerFromUnsignedInt(fhir.UnsignedInt(tc.value))
			if err != nil {
				t.Fatalf("IntegerFromUnsignedInt(%v): unexpected error '%v'", tc.name, err)
			}

			if got, want := sut.GetValue(), tc.value; got != int32(want) {
				t.Errorf("IntegerFromUnsignedInt(%v): got %v, want %v", tc.name, got, want)
			}
		})
	}
}

func TestMarkdown(t *testing.T) {
	want := "This is **basically** just for code _coverage_; let's be honest."

	sut := fhir.Markdown(want)

	if got := sut.GetValue(); !cmp.Equal(got, want) {
		t.Errorf("Markdown: got %v, want %v", got, want)
	}
}

func TestOID(t *testing.T) {
	const wantPrefix = "urn:oid:"
	want := "foobar"

	sut := fhir.OID(want)

	if got := sut.GetValue(); !strings.HasPrefix(got, wantPrefix) {
		t.Errorf("OID: got value '%v', want prefix '%v'", wantPrefix, got)
	}
	if got := sut.GetValue(); !strings.HasSuffix(got, want) {
		t.Errorf("OID: got value '%v', want suffix '%v'", want, got)
	}
}

func TestString(t *testing.T) {
	want := "Lorem ipsum"

	sut := fhir.String(want)

	if got := sut.GetValue(); !cmp.Equal(got, want) {
		t.Errorf("String: got %v, want %v", got, want)
	}
}

func TestStrings(t *testing.T) {
	want := []string{"Lorem", "ipsum", "dalor", "sit", "amet"}
	toString := func(s *dtpb.String) string { return s.GetValue() }

	got := fhir.Strings(want...)

	if got := slices.Map(got, toString); !cmp.Equal(got, want) {
		t.Errorf("String: got %v, want %v", got, want)
	}
}

func TestStringFromCode(t *testing.T) {
	testCases := []struct {
		name  string
		value *dtpb.Code
	}{
		{"Nil", nil},
		{"WithValue", fhir.Code("foobar")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := fhir.StringFromCode(tc.value)

			if got, want := got.GetValue(), tc.value.GetValue(); got != want {
				t.Errorf("StringFromCode(%v): got %v, want %v", tc.name, got, want)
			}
		})
	}
}

func TestStringFromMarkdown(t *testing.T) {
	testCases := []struct {
		name  string
		value *dtpb.Markdown
	}{
		{"Nil", nil},
		{"WithValue", fhir.Markdown("foobar")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := fhir.StringFromMarkdown(tc.value)

			if got, want := got.GetValue(), tc.value.GetValue(); got != want {
				t.Errorf("StringFromMarkdown(%v): got %v, want %v", tc.name, got, want)
			}
		})
	}
}

func TestStringFromID(t *testing.T) {
	testCases := []struct {
		name  string
		value *dtpb.Id
	}{
		{"Nil", nil},
		{"WithValue", fhir.ID("foobar")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := fhir.StringFromID(tc.value)

			if got, want := got.GetValue(), tc.value.GetValue(); got != want {
				t.Errorf("StringFromID(%v): got %v, want %v", tc.name, got, want)
			}
		})
	}
}

func TestURIFromCanonical(t *testing.T) {
	testCases := []struct {
		name  string
		value *dtpb.Canonical
	}{
		{"Nil", nil},
		{"WithValue", canonical.New("https://some-uri")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := fhir.URIFromCanonical(tc.value)

			if got, want := got.GetValue(), tc.value.GetValue(); got != want {
				t.Errorf("URIFromCanonical(%v): got %v, want %v", tc.name, got, want)
			}
		})
	}
}

func TestURIFromOID(t *testing.T) {
	testCases := []struct {
		name  string
		value *dtpb.Oid
	}{
		{"Nil", nil},
		{"WithValue", fhir.OID("https://some-uri")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := fhir.URIFromOID(tc.value)

			if got, want := got.GetValue(), tc.value.GetValue(); got != want {
				t.Errorf("URIFromOID(%v): got %v, want %v", tc.name, got, want)
			}
		})
	}
}

func TestURIFromURL(t *testing.T) {
	testCases := []struct {
		name  string
		value *dtpb.Url
	}{
		{"Nil", nil},
		{"WithValue", fhir.URL("https://some-uri.com")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := fhir.URIFromURL(tc.value)

			if got, want := got.GetValue(), tc.value.GetValue(); got != want {
				t.Errorf("URIFromURL(%v): got %v, want %v", tc.name, got, want)
			}
		})
	}
}

func TestURIFromUUID(t *testing.T) {
	testCases := []struct {
		name  string
		value *dtpb.Uuid
	}{
		{"Nil", nil},
		{"WithValue", fhir.RandomUUID()},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := fhir.URIFromUUID(tc.value)

			if got, want := got.GetValue(), tc.value.GetValue(); got != want {
				t.Errorf("URIFromUUID(%v): got %v, want %v", tc.name, got, want)
			}
		})
	}
}
