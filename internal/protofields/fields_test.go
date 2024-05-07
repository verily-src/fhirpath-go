package protofields_test

import (
	"testing"

	opb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/observation_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/protofields"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestUnwrapChoiceField_GetsUnderlyingMessage(t *testing.T) {
	dateTime := fhir.DateTimeNow()

	testCases := []struct {
		name  string
		input proto.Message
		want  proto.Message
	}{
		{
			name: "gets boolean of Patient deceased field",
			input: &ppb.Patient_DeceasedX{
				Choice: &ppb.Patient_DeceasedX_Boolean{
					Boolean: fhir.Boolean(true),
				},
			},
			want: fhir.Boolean(true),
		},
		{
			name: "gets date of Patient deceased field",
			input: &ppb.Patient_DeceasedX{
				Choice: &ppb.Patient_DeceasedX_DateTime{
					DateTime: dateTime,
				},
			},
			want: dateTime,
		},
		{
			name: "",
			input: &opb.Observation_Component_ValueX{
				Choice: &opb.Observation_Component_ValueX_StringValue{
					StringValue: fhir.String("some string"),
				},
			},
			want: fhir.String("some string"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := protofields.UnwrapOneofField(tc.input, "choice")

			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("UnwrapChoiceField returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}
