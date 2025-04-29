package fhirpathtest_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/verily-src/fhirpath-go/fhirpath"
	"github.com/verily-src/fhirpath-go/fhirpath/fhirpathtest"
	"github.com/verily-src/fhirpath-go/fhirpath/system"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	pb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
)

func ExampleError() {
	want := errors.New("example error")
	expr := fhirpathtest.Error(want)

	_, err := expr.Evaluate([]fhirpath.Resource{})
	if errors.Is(err, want) {
		fmt.Printf("err = '%v'", want)
	}

	// Output: err = 'example error'
}

func ExampleReturn() {
	want := system.Boolean(true)
	expr := fhirpathtest.Return(want)

	got, err := expr.Evaluate([]fhirpath.Resource{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("got = %v", bool(got[0].(system.Boolean)))
	// Output: got = true
}
func ExampleReturnCollection() {
	want := system.Collection{system.Boolean(true)}
	expr := fhirpathtest.ReturnCollection(want)

	got, err := expr.Evaluate([]fhirpath.Resource{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("got = %v", bool(got[0].(system.Boolean)))
	// Output: got = true
}

// TestFHIRPathResource is based on issue can't call Evaluate() because fhir.Resource is internal? #18
// https://github.com/verily-src/fhirpath-go/issues/18
//
// Testing that by using only the fhirpath package, we can still evaluate
// a FHIRPath expression that returns a resource.
func TestFHIRPathResource(t *testing.T) {
	want := system.Collection{
		&dtpb.String{Value: "John"},
	}
	// Create a FHIR Patient resource
	patient := &pb.Patient{
		Name: []*dtpb.HumanName{
			{
				Given: []*dtpb.String{
					{Value: "John"},
				},
				Family: &dtpb.String{Value: "Doe"},
			},
		},
	}

	// Compile the FHIRPath expression
	expr := fhirpath.MustCompile("name.given")

	// Wrap the Patient resource in a FHIRPath Resource
	resource := []fhirpath.Resource{patient}

	// Evaluate the expression against the Patient resource
	got, err := expr.Evaluate(resource)
	if err != nil {
		t.Errorf("Error evaluating FHIRPath expression: %v", err)
	}

	// Check if the results match the expected output
	if len(got) != 1 {
		t.Errorf("Expected 1 result, got %d", len(got))
	}
	gotValue := got[0].(*dtpb.String).GetValue()
	wantValue := want[0].(*dtpb.String).GetValue()
	if gotValue != wantValue {
		t.Errorf("Expected %s, got %s", wantValue, gotValue)
	}
}
