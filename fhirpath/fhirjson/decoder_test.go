package fhirjson_test

import (
	"bytes"
	"testing"
	"time"

	codes_go_proto "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	opb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/observation_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"

	"github.com/verily-src/fhirpath-go/fhirpath/fhirjson"
	"github.com/verily-src/fhirpath-go/internal/fhir"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestDecodeNew_Success(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  fhir.Resource
	}{
		{
			name:  "Patient resource",
			input: `{"resourceType": "Patient", "id": "p1"}`,
			want: &ppb.Patient{
				Id: &dtpb.Id{Value: "p1"},
			},
		},
		{
			name:  "Minimal Patient",
			input: `{"resourceType": "Patient"}`,
			want:  &ppb.Patient{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dec := fhirjson.NewDecoder(bytes.NewReader([]byte(tc.input)))

			res, err := dec.DecodeNew()
			if err != nil {
				t.Errorf("DecodeNew error: got %v, want nil", err)
			}

			if got, want := res, tc.want; !cmp.Equal(got, want, protocmp.Transform()) {
				t.Errorf("DecodeNew resource: got %v, want %v", got, want)
			}
		})
	}
}

func TestDecodeNew_Error(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "Invalid JSON",
			input: `{"resourceType": "Patient",`,
		},
		{
			name:  "Empty input",
			input: ``,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dec := fhirjson.NewDecoder(bytes.NewReader([]byte(tc.input)))

			res, err := dec.DecodeNew()
			if err == nil {
				t.Errorf("DecodeNew: expected error, got nil")
			}

			if res != nil {
				t.Errorf("DecodeNew: expected nil resource, got %v", res)
			}
		})
	}
}

func TestDecoderTimeZone(t *testing.T) {
	testCases := []struct {
		name string
		json string
		zone *time.Location
		want fhir.Resource
	}{
		{
			name: "Not setting timezone implies Local",
			json: `{
				"resourceType": "Patient",
				"id": "example-patient",
				"birthDate": "1990-05-27"
			}`,
			zone: nil,
			want: &ppb.Patient{
				Id: &dtpb.Id{Value: "example-patient"},
				BirthDate: &dtpb.Date{
					ValueUs:   time.Date(1990, 5, 27, 0, 0, 0, 0, time.UTC).UnixMicro(),
					Timezone:  "Local",
					Precision: dtpb.Date_DAY,
				},
			},
		},
		{
			name: "Set timezone to UTC",
			json: `{
				"resourceType": "Patient",
				"id": "example-patient",
				"birthDate": "1990-05-27"
			}`,
			zone: time.UTC,
			want: &ppb.Patient{
				Id: &dtpb.Id{Value: "example-patient"},
				BirthDate: &dtpb.Date{
					ValueUs:   time.Date(1990, 5, 27, 0, 0, 0, 0, time.UTC).UnixMicro(),
					Timezone:  "UTC",
					Precision: dtpb.Date_DAY,
				},
			},
		},
		{
			name: "Timezone does not apply to DateTime",
			json: `{
				   "resourceType": "Observation",
				   "id": "example-observation",
				   "status": "final",
				   "code": { "coding": [ { "system": "http://loinc.org", "code": "1234-5" } ] },
				   "effectiveDateTime": "2022-01-02T03:04:05Z"
			   }`,
			zone: time.UTC,
			want: &opb.Observation{
				Id:     &dtpb.Id{Value: "example-observation"},
				Status: &opb.Observation_StatusCode{Value: codes_go_proto.ObservationStatusCode_FINAL},
				Code: &dtpb.CodeableConcept{
					Coding: []*dtpb.Coding{
						{
							System: &dtpb.Uri{Value: "http://loinc.org"},
							Code:   &dtpb.Code{Value: "1234-5"},
						},
					},
				},
				Effective: &opb.Observation_EffectiveX{
					Choice: &opb.Observation_EffectiveX_DateTime{
						DateTime: &dtpb.DateTime{
							ValueUs:   time.Date(2022, 1, 2, 3, 4, 5, 0, time.UTC).UnixMicro(),
							Timezone:  "Z",
							Precision: dtpb.DateTime_SECOND,
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dec := fhirjson.NewDecoder(bytes.NewReader([]byte(tc.json)))
			if tc.zone != nil {
				dec.TimeZone(tc.zone)
			}

			var got fhir.Resource
			switch tc.want.(type) {
			case *ppb.Patient:
				got = &ppb.Patient{}
			case *opb.Observation:
				got = &opb.Observation{}
			default:
				t.Fatalf("unsupported resource type: %T", tc.want)
			}

			err := dec.Decode(got)
			if err != nil {
				t.Errorf("Decode error: got %v, want nil", err)
			}

			if got, want := got, tc.want; !cmp.Equal(got, want, protocmp.Transform()) {
				t.Errorf("Decode resource diff: %s", cmp.Diff(got, want, protocmp.Transform()))
			}
		})
	}
}
