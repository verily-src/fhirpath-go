package fhirpathtest_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/verily-src/fhirpath-go/fhirpath"
	"github.com/verily-src/fhirpath-go/fhirpath/fhirpathtest"
	"github.com/verily-src/fhirpath-go/fhirpath/system"

	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	pb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	qrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/questionnaire_response_go_proto"
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

// Create a FHIR Patient resource
var (
	patientOfficialName = &dtpb.HumanName{
		Given: []*dtpb.String{
			{Value: "John"},
		},
		Family: &dtpb.String{Value: "Doe"},
		Use:    &dtpb.HumanName_UseCode{Value: *cpb.NameUseCode_OFFICIAL.Enum()},
	}

	testPatient = &pb.Patient{
		Name: []*dtpb.HumanName{
			patientOfficialName,
			{
				Given: []*dtpb.String{
					{Value: "Johnny"},
				},
				Family: &dtpb.String{Value: "Doe"},
				Use:    &dtpb.HumanName_UseCode{Value: *cpb.NameUseCode_NICKNAME.Enum()},
			},
		},
	}

	testQuestionnaireResponse = &qrpb.QuestionnaireResponse{
		Item: []*qrpb.QuestionnaireResponse_Item{
			{
				Answer: []*qrpb.QuestionnaireResponse_Item_Answer{
					{Value: &qrpb.QuestionnaireResponse_Item_Answer_ValueX{
						Choice: &qrpb.QuestionnaireResponse_Item_Answer_ValueX_Coding{
							Coding: &dtpb.Coding{Code: &dtpb.Code{Value: "foo"}},
						},
					}},
				},
			},
		},
	}
)

// TestFHIRPathResource is based on issue can't call Evaluate() because fhir.Resource is internal? #18
// https://github.com/verily-src/fhirpath-go/issues/18
//
// Testing that by using only the fhirpath package, we can still evaluate
// a FHIRPath expression that returns a resource.
func TestFHIRPathResource(t *testing.T) {
	want := system.Collection{
		&dtpb.String{Value: "John"},
	}

	// Compile the FHIRPath expression
	expr := fhirpath.MustCompile("name.given")

	// Wrap the Patient resource in a FHIRPath Resource
	resource := []fhirpath.Resource{testPatient}

	// Evaluate the expression against the Patient resource
	got, err := expr.Evaluate(resource)
	if err != nil {
		t.Errorf("Error evaluating FHIRPath expression: %v", err)
	}

	// Check if the results match the expected output
	if len(got) != 2 {
		t.Errorf("Expected 2 results, got %d", len(got))
	}
	gotValue := got[0].(*dtpb.String).GetValue()
	wantValue := want[0].(*dtpb.String).GetValue()
	if gotValue != wantValue {
		t.Errorf("Expected %s, got %s", wantValue, gotValue)
	}
}

func TestFHIRPathWhere(t *testing.T) {
	resource := []fhirpath.Resource{testPatient}

	// Test where: name is official (cpb.NameUseCode_OFFICIAL)
	exprWhere := fhirpath.MustCompile("name.where(use = 'official')")
	gotWhere, err := exprWhere.Evaluate(resource)
	if err != nil {
		t.Errorf("Error evaluating where: %v", err)
	}
	wantWhere := patientOfficialName

	if len(gotWhere) != 1 {
		t.Errorf("Expected 1 result from where, got %d", len(gotWhere))
	}
	if len(gotWhere) == 1 {
		gotValue := gotWhere[0].(*dtpb.HumanName).GetGiven()[0].GetValue()
		wantValue := wantWhere.GetGiven()[0].GetValue()
		if gotValue != wantValue {
			t.Errorf("Expected %s from where, got %q", wantValue, gotValue)
		}
	}
}

func TestFHIRPathCombine(t *testing.T) {
	// Use the testQuestionnaireResponse defined above (empty for now, but can be extended)
	resource := []fhirpath.Resource{testQuestionnaireResponse}

	// Evaluate the FHIRPath expression: QuestionnaireResponse.item.answer.combine(today())
	expr := fhirpath.MustCompile("item.answer.combine(today())")
	got, err := expr.Evaluate(resource)
	if err != nil {
		t.Fatalf("Error evaluating combine: %v", err)
	}

	if got == nil {
		t.Errorf("Expected non-nil result from combine, got nil")
	}
	if len(got) != 2 {
		t.Errorf("Expected 2 results from combine, got %d", len(got))
	}
}
