package patient_test

import (
	"fmt"
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	apb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/account_go_proto"
	appb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/appointment_go_proto"
	arpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/appointment_response_go_proto"
	cppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/care_plan_go_proto"
	clpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/claim_go_proto"
	crpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/communication_request_go_proto"
	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/condition_go_proto"
	cerpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/coverage_eligibility_request_go_proto"
	dpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/device_go_proto"
	drpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/device_request_go_proto"
	epb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/encounter_go_proto"
	erpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/enrollment_request_go_proto"
	eobpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/explanation_of_benefit_go_proto"
	irpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/immunization_recommendation_go_proto"
	mrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medication_request_go_proto"
	nopb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/nutrition_order_go_proto"
	opb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/observation_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	procpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/procedure_go_proto"
	qrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/questionnaire_response_go_proto"
	rppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/related_person_go_proto"
	rgpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/request_group_go_proto"
	rspb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/research_subject_go_proto"
	rapb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/risk_assessment_go_proto"
	srpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/service_request_go_proto"
	surpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/supply_request_go_proto"
	tpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/task_go_proto"
	vppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/vision_prescription_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/resource/patient"
	"google.golang.org/protobuf/testing/protocmp"
)

// TODO(PHP-7652): Update testing
func TestIDFromResource(t *testing.T) {
	mockPatientID := "patient-id"
	mockPatientRef := patient.Reference(mockPatientID)
	mockAccount := &apb.Account{}

	type TestCase struct {
		name          string
		resource      fhir.Resource
		expected      string
		expectedError error
	}

	// Helpers to make pass and fail cases consistent
	pass := func(name string, res fhir.Resource) TestCase {
		return TestCase{
			name:          fmt.Sprintf("%v as input, contains patient id, returns id", name),
			resource:      res,
			expected:      mockPatientID,
			expectedError: nil,
		}
	}
	fail := func(name string, res fhir.Resource) TestCase {
		return TestCase{
			name:          fmt.Sprintf("%v as input, does not contain id, returns err", name),
			resource:      res,
			expected:      "",
			expectedError: patient.ErrExtractingPatientID,
		}
	}

	testCases := []TestCase{

		// General Error Conditions
		{"unsupported resource type, returns err", mockAccount, "", patient.ErrUnsupportedType},

		// Resources
		pass("patient", &ppb.Patient{Id: fhir.ID(mockPatientID)}),
		fail("patient", &ppb.Patient{}),
		pass("encounter", &epb.Encounter{Subject: mockPatientRef}),
		fail("encounter", &epb.Encounter{}),
		pass("device", &dpb.Device{Patient: mockPatientRef}),
		fail("device", &dpb.Device{}),
		pass("explanation of benefit", &eobpb.ExplanationOfBenefit{Patient: mockPatientRef}),
		fail("explanation of benefit", &eobpb.ExplanationOfBenefit{}),
		pass("research subject", &rspb.ResearchSubject{Individual: mockPatientRef}),
		fail("research subject", &rspb.ResearchSubject{}),
		pass("related person", &rppb.RelatedPerson{Patient: mockPatientRef}),
		fail("related person", &rppb.RelatedPerson{}),

		// Event Patterns
		pass("questionnaire response", &qrpb.QuestionnaireResponse{Subject: mockPatientRef}),
		fail("questionnaire response", &qrpb.QuestionnaireResponse{}),
		pass("risk assessment", &rapb.RiskAssessment{Subject: mockPatientRef}),
		fail("risk assessment", &rapb.RiskAssessment{}),
		pass("condition", &cpb.Condition{Subject: mockPatientRef}),
		fail("condition", &cpb.Condition{}),
		pass("procedure", &procpb.Procedure{Subject: mockPatientRef}),
		fail("procedure", &procpb.Procedure{}),
		pass("observation", &opb.Observation{Subject: mockPatientRef}),
		fail("observation", &opb.Observation{}),
		pass("task", &tpb.Task{ForValue: mockPatientRef}),
		pass("task", &tpb.Task{Focus: mockPatientRef}),
		fail("task", &tpb.Task{}),

		// Request Patterns
		pass("appointment response", &arpb.AppointmentResponse{Actor: mockPatientRef}),
		fail("appointment response", &arpb.AppointmentResponse{}),
		pass("appointment", &appb.Appointment{Participant: []*appb.Appointment_Participant{
			{Actor: mockPatientRef},
		}}),
		fail("appointment", &appb.Appointment{}),
		pass("care plan", &cppb.CarePlan{Subject: mockPatientRef}),
		fail("care plan", &cppb.CarePlan{}),
		pass("claim", &clpb.Claim{Patient: mockPatientRef}),
		fail("claim", &clpb.Claim{}),
		pass("communication request", &crpb.CommunicationRequest{Subject: mockPatientRef}),
		fail("communication reuqest", &crpb.CommunicationRequest{}),
		pass("coverage eligibility request", &cerpb.CoverageEligibilityRequest{Patient: mockPatientRef}),
		fail("coverage eligibility request", &cerpb.CoverageEligibilityRequest{}),
		pass("device request", &drpb.DeviceRequest{Subject: mockPatientRef}),
		fail("device request", &drpb.DeviceRequest{}),
		pass("enrollment request", &erpb.EnrollmentRequest{Candidate: mockPatientRef}),
		fail("enrollment request", &erpb.EnrollmentRequest{}),
		pass("immunization recommendation", &irpb.ImmunizationRecommendation{Patient: mockPatientRef}),
		fail("immunization recommendation", &irpb.ImmunizationRecommendation{}),
		pass("medication request", &mrpb.MedicationRequest{Subject: mockPatientRef}),
		fail("medication request", &mrpb.MedicationRequest{}),
		pass("nutrition order", &nopb.NutritionOrder{Patient: mockPatientRef}),
		fail("nutrition order", &nopb.NutritionOrder{}),
		pass("request group", &rgpb.RequestGroup{Subject: mockPatientRef}),
		fail("request group", &rgpb.RequestGroup{}),
		pass("service request", &srpb.ServiceRequest{Subject: mockPatientRef}),
		fail("service request", &srpb.ServiceRequest{}),
		pass("supply request", &surpb.SupplyRequest{DeliverTo: mockPatientRef}),
		pass("supply request", &surpb.SupplyRequest{Requester: mockPatientRef}),
		fail("supply request", &surpb.SupplyRequest{}),
		pass("vision prescription", &vppb.VisionPrescription{Patient: mockPatientRef}),
		fail("vision prescription", &vppb.VisionPrescription{}),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			patientID, err := patient.IDFromResource(tc.resource)
			if got, want := err, tc.expectedError; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("IDFromResource(%s) error got = %v, want = %v", tc.name, got, want)
			}
			if got, want := patientID, tc.expected; got != want {
				t.Errorf("IDFromResource(%s) got = %v, want = %v", tc.name, got, want)
			}
		})
	}
}

func TestReference(t *testing.T) {
	want := &dtpb.Reference{
		Type: fhir.URI("Patient"),
		Reference: &dtpb.Reference_PatientId{
			PatientId: &dtpb.ReferenceId{
				Value: "123",
			},
		},
	}

	got := patient.Reference("123")

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("Reference mismatch (-want, +got)\n%s", diff)
	}
}
