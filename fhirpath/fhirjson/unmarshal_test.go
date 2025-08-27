package fhirjson_test

import (
	"errors"
	"testing"
	"time"

	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/fhirpath/fhirjson"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestUnmarshal_ValidInput_OutputsResource(t *testing.T) {
	want, json := newPatientData(t)
	got := &ppb.Patient{}

	err := fhirjson.Unmarshal([]byte(json), got)
	if err != nil {
		t.Fatalf("Unmarshal: unexpected error: %v", err)
	}

	if !cmp.Equal(got, want, protocmp.Transform()) {
		t.Errorf("Unmarshal: got %v, want %v", got, want)
	}
}

func TestUnmarshal_InvalidInput_ReturnsError(t *testing.T) {
	testCases := []struct {
		name string
		json string
		out  fhir.Resource
		want error
	}{
		{
			name: "NilPointer",
			json: `{ "resourceType": "Patient" }`,
			out:  nil,
			want: fhirjson.ErrNilUnmarshalOutput,
		},
		{
			name: "WrongResourceType",
			json: `{ "resourceType": "Device" }`,
			out:  &ppb.Patient{},
			want: fhirjson.ErrWrongUnmarshalType,
		},
		{
			name: "BadJSON",
			json: `{ "resourceType": `,
			out:  &ppb.Patient{},
			want: fhirjson.ErrUnmarshalEncoding,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := fhirjson.Unmarshal([]byte(tc.json), tc.out)

			if got, want := err, tc.want; !errors.Is(got, want) {
				t.Errorf("Unmarshal(%v): got err '%v', want err '%v'", tc.name, got, want)
			}
		})
	}
}

func TestUnmarshalOut_ValidInput_ReturnsResource(t *testing.T) {
	want, json := newPatientData(t)

	got, err := fhirjson.UnmarshalOut([]byte(json))
	if err != nil {
		t.Fatalf("UnmarshalOut: unexpected error: %v", err)
	}

	if !cmp.Equal(got, want, protocmp.Transform()) {
		t.Errorf("UnmarshalOut: got %v, want %v", got, want)
	}
}

func TestUnmarshalOut_InvalidInput_ReturnsError(t *testing.T) {
	testCases := []struct {
		name string
		json string
		out  fhir.Resource
		want error
	}{
		{
			name: "BadJSON",
			json: `{ "resourceType": `,
			want: fhirjson.ErrUnmarshalEncoding,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := fhirjson.UnmarshalOut([]byte(tc.json))

			if got, want := err, tc.want; !errors.Is(got, want) {
				t.Errorf("UnmarshalOut(%v): got err '%v', want err '%v'", tc.name, got, want)
			}
		})
	}
}

func newUnmarshaller(t *testing.T) *fhirjson.Unmarshaller {
	t.Helper()
	tz, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("UnmarshalOptions.Unmarshal: got unexpected err %v", err)
	}
	return &fhirjson.Unmarshaller{
		TimeZone: tz,
	}
}

func TestUnmarshalOptionsUnmarshal_ValidInput_ReturnsResource(t *testing.T) {
	sut := newUnmarshaller(t)
	want, json := newPatientData(t)
	got := &ppb.Patient{}

	err := sut.Unmarshal([]byte(json), got)
	if err != nil {
		t.Fatalf("UnmarshalOptions.Unmarshal: unexpected error: %v", err)
	}

	if !cmp.Equal(got, want, protocmp.Transform()) {
		t.Errorf("UnmarshalOptions.Unmarshal: got %v, want %v", got, want)
	}
}

func TestUnmarshalOptionsUnmarshalOut_ValidInput_ReturnsResource(t *testing.T) {
	sut := newUnmarshaller(t)
	want, json := newPatientData(t)

	got, err := sut.UnmarshalOut([]byte(json))
	if err != nil {
		t.Fatalf("UnmarshalOptions.UnmarshalOut: unexpected error: %v", err)
	}

	if !cmp.Equal(got, want, protocmp.Transform()) {
		t.Errorf("UnmarshalOptions.UnmarshalOut: got %v, want %v", got, want)
	}
}
