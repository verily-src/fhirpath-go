package fhirjson_test

import (
	"errors"
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	bpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/binary_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/fhirpath/fhirjson"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestMarshal_WithResource_ReturnsJSON(t *testing.T) {
	patient, patientJSON := newPatientData(t)
	testCases := []struct {
		name     string
		resource fhir.Resource
		want     string
	}{
		{
			name:     "Patient",
			resource: patient,
			want:     patientJSON,
		},
		{
			name: "Device",
			resource: &bpb.Binary{
				Data: fhir.Base64Binary([]byte{0xde, 0xad, 0xbe, 0xef}),
			},
			want: `{"data":"3q2+7w==","resourceType":"Binary"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := fhirjson.Marshal(tc.resource)
			if err != nil {
				t.Fatalf("Marshal(%v): unexpected error: %v", tc.name, err)
			}

			if got, want := string(got), tc.want; got != want {
				t.Errorf("Marshal(%v): got '%v', want '%v'", tc.name, got, want)
			}
		})
	}
}

func TestMarshal_NilResource_ReturnsError(t *testing.T) {
	_, err := fhirjson.Marshal(nil)

	if got, want := err, fhirjson.ErrNilMarshalResource; !errors.Is(got, want) {
		t.Errorf("Marshal: got err '%v', want err '%v'", got, want)
	}
}

func TestMarshalOptionsMarshal_WithResource_ReturnsJSON(t *testing.T) {
	const (
		enableIndent = true
		indent       = " "
	)
	testCases := []struct {
		name     string
		resource fhir.Resource
		want     string
	}{
		{
			name: "Patient",
			resource: &ppb.Patient{
				Name: []*dtpb.HumanName{
					{
						Given: fhir.Strings("Matt"),
					},
				},
			},
			want: `{
 "name": [
  {
   "given": [
    "Matt"
   ]
  }
 ],
 "resourceType": "Patient"
}`,
		},
		{
			name: "Device",
			resource: &bpb.Binary{
				Data: fhir.Base64Binary([]byte{0xde, 0xad, 0xbe, 0xef}),
			},
			want: `{
 "data": "3q2+7w==",
 "resourceType": "Binary"
}`,
		},
	}
	sut := &fhirjson.Marshaller{
		EnableIndent: enableIndent,
		Indent:       indent,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := sut.Marshal(tc.resource)
			if err != nil {
				t.Fatalf("MarshalOptions.Marshal(%v): unexpected error: %v", tc.name, err)
			}

			if got, want := string(got), tc.want; got != want {
				t.Errorf("MarshalOptions.Marshal(%v): got '%v', want '%v'", tc.name, got, want)
			}
		})
	}
}

func TestMarshalOptionsMarshal_NilResource_ReturnsError(t *testing.T) {
	const (
		enableIndent = true
		indent       = " "
	)
	sut := &fhirjson.Marshaller{
		EnableIndent: enableIndent,
		Indent:       indent,
	}

	_, err := sut.Marshal(nil)

	if got, want := err, fhirjson.ErrNilMarshalResource; !errors.Is(got, want) {
		t.Errorf("MarshalOptions.Marshal: got err '%v', want err '%v'", got, want)
	}
}

// This is just a smoke-test to prove that a round-trip conversion works correctly.
func TestRoundTrip(t *testing.T) {
	want, _ := newPatientData(t)

	json, err := fhirjson.Marshal(want)
	if err != nil {
		t.Fatalf("Round-Trip: unexpected error: %v", err)
	}

	got := &ppb.Patient{}
	err = fhirjson.Unmarshal(json, got)

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("Round-trip: got-, want+\n%v", diff)
	}
}
