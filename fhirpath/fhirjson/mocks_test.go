package fhirjson_test

import (
	"testing"

	"github.com/google/fhir/go/fhirversion"
	"github.com/google/fhir/go/jsonformat"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/verily-src/fhirpath-go/internal/containedresource"
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

// newPatientData returns a patient object along with the JSON data that can
// form it.
func newPatientData(t *testing.T) (*ppb.Patient, string) {
	t.Helper()
	marshaller, err := jsonformat.NewMarshaller(false, "", "", fhirversion.R4)
	if err != nil {
		t.Fatalf("unexpected error while creating marshaller: %v", err)
	}
	patient := &ppb.Patient{
		Id: fhir.ID("deadbeef"),
		Name: []*dtpb.HumanName{
			{
				Prefix: fhir.Strings("Dr."),
				Given:  fhir.Strings("Mantis"),
				Family: fhir.String("Taboggan"),
				Suffix: fhir.Strings("Md."),
			},
		},
	}

	cr := containedresource.Wrap(patient)

	json, err := marshaller.Marshal(cr)
	if err != nil {
		t.Fatalf("unexpected error while marshalling resource: %v", err)
	}

	return patient, string(json)
}
