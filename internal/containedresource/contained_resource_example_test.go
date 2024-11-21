package containedresource_test

import (
	"fmt"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/verily-src/fhirpath-go/internal/containedresource"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"google.golang.org/protobuf/proto"
)

func ExampleWrap() {
	patient := &patient_go_proto.Patient{
		Id: fhir.ID("12345"),
	}

	cr := containedresource.Wrap(patient)

	fmt.Printf("Patient ID = %v", cr.GetPatient().GetId().GetValue())
	// Output: Patient ID = 12345
}

func ExampleUnwrap() {
	patient := &patient_go_proto.Patient{
		Id: fhir.ID("12345"),
	}
	cr := containedresource.Wrap(patient)

	unwrapped := containedresource.Unwrap(cr).(*patient_go_proto.Patient)

	if proto.Equal(patient, unwrapped) {
		fmt.Printf("Resources match!")
	}
	// Output: Resources match!
}

func ExampleTypeOf() {
	patient := &patient_go_proto.Patient{
		Id: fhir.ID("12345"),
	}
	cr := containedresource.Wrap(patient)

	crType := containedresource.TypeOf(cr)

	fmt.Printf("Contained Resource type = %v", crType)
	// Output: Contained Resource type = Patient
}

func ExampleID() {
	patient := &patient_go_proto.Patient{
		Id: fhir.ID("12345"),
	}
	cr := containedresource.Wrap(patient)

	id := containedresource.ID(cr)

	fmt.Printf("Contained Resource ID = %v", id)
	// Output: Contained Resource ID = 12345
}

func ExampleGenerateIfNoneExist() {
	patient := &patient_go_proto.Patient{
		Id: fhir.ID("12345"),
		Identifier: []*dtpb.Identifier{
			&dtpb.Identifier{
				System: &dtpb.Uri{Value: "http://fake.com"},
				Value:  &dtpb.String{Value: "9efbf82d-7a58-4d14-bec1-63f8fda148a8"},
			},
		},
	}

	cr := containedresource.Wrap(patient)

	value, err := containedresource.GenerateIfNoneExist(cr, "http://fake.com", true)
	if err != nil {
		panic(err)
	}

	headers := map[string]string{}

	if value != "" {
		headers["If-None-Exist"] = value
	}

	fmt.Printf("If-None-Exist: %v", headers["If-None-Exist"])
	// Output: If-None-Exist: identifier=http%3A%2F%2Ffake.com%7C9efbf82d-7a58-4d14-bec1-63f8fda148a8
}
