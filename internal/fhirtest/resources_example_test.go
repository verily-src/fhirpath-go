package fhirtest_test

import (
	"fmt"
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirtest"
)

func ExampleWithResourceModification() {
	t := &testing.T{}

	patient := fhirtest.NewResource(t, "Patient", fhirtest.WithResourceModification(func(p *ppb.Patient) {
		p.Name = []*dtpb.HumanName{{Family: fhir.String("Ursa")}}
	})).(*ppb.Patient)

	fmt.Printf("patient.Name[0].Family = %v", patient.Name[0].Family)
	// Output: patient.Name[0].Family = value:"Ursa"
}

func ExampleNewResourceFromBase() {
	t := &testing.T{}
	original := &ppb.Patient{
		Id:   fhir.ID("uuid-a"),
		Name: []*dtpb.HumanName{{Family: fhir.String("Ursa")}},
	}

	// Apply options on original.
	modified := fhirtest.NewResourceFromBase(t, original,
		fhirtest.WithResourceModification(func(p *ppb.Patient) {
			p.Name[0].Family = fhir.String("Major")
			p.Name[0].Given = fhir.Strings("Aseem")
		}),
		fhirtest.WithProtoField("id", fhir.ID("uuid-b")),
	).(*ppb.Patient)

	fmt.Printf("ID = '%v', Family = '%v', Given = '%v'",
		modified.Id.Value,
		modified.Name[0].Family.Value,
		modified.Name[0].Given[0].Value,
	)
	// Output: ID = 'uuid-b', Family = 'Major', Given = 'Aseem'
}
