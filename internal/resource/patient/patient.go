package patient

import (
	"errors"
	"fmt"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	appb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/appointment_go_proto"
	arpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/appointment_response_go_proto"
	cppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/care_plan_go_proto"
	clpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/claim_go_proto"
	commpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/communication_go_proto"
	crpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/communication_request_go_proto"
	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/condition_go_proto"
	cerpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/coverage_eligibility_request_go_proto"
	dpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/device_go_proto"
	drpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/device_request_go_proto"
	epb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/encounter_go_proto"
	erpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/enrollment_request_go_proto"
	eobpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/explanation_of_benefit_go_proto"
	irpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/immunization_recommendation_go_proto"
	lpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/list_go_proto"
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
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/resource"
)

var (
	ErrExtractingPatientID = errors.New("extracting patient id")
	ErrUnsupportedType     = errors.New("extraction unsupported for type")
)

// IDFromResource returns the patient's ID from a FHIR Resource.
// This function provides a mapping from an input resource to the patient ID it references.
func IDFromResource(res fhir.Resource) (string, error) {
	idOrError := func(patientRef *dtpb.Reference) (string, error) {
		if id := patientRef.GetPatientId().GetValue(); id != "" {
			return id, nil
		}
		return "", fmt.Errorf("%w from %T", ErrExtractingPatientID, res)
	}

	switch res := res.(type) {

	//---------------------------------------------------------------------------
	// Unpatterned Resources
	//---------------------------------------------------------------------------

	case *ppb.Patient:
		if id := res.GetId().GetValue(); id != "" {
			return id, nil
		}
		return "", fmt.Errorf("%w from %T", ErrExtractingPatientID, res)
	case *epb.Encounter:
		return idOrError(res.GetSubject())
	case *dpb.Device:
		return idOrError(res.GetPatient())
	case *eobpb.ExplanationOfBenefit:
		return idOrError(res.GetPatient())
	case *rspb.ResearchSubject:
		return idOrError(res.GetIndividual())
	case *rppb.RelatedPerson:
		return idOrError(res.GetPatient())
	case *lpb.List:
		return idOrError(res.GetSubject())

	//---------------------------------------------------------------------------
	// Event Pattern Resources
	//---------------------------------------------------------------------------

	case *qrpb.QuestionnaireResponse:
		return idOrError(res.GetSubject())
	case *rapb.RiskAssessment:
		return idOrError(res.GetSubject())
	case *cpb.Condition:
		return idOrError(res.GetSubject())
	case *procpb.Procedure:
		return idOrError(res.GetSubject())
	case *opb.Observation:
		return idOrError(res.GetSubject())
	case *tpb.Task:
		// TODO(b/254654059): Remove usage of Task.for once decision engine supports event-patient extraction override.
		if id := res.GetFocus().GetPatientId().GetValue(); id != "" {
			return id, nil
		} else if id := res.GetForValue().GetPatientId().GetValue(); id != "" {
			return id, nil
		}
		return "", fmt.Errorf("%w from %T", ErrExtractingPatientID, res)
	case *commpb.Communication:
		return idOrError(res.GetSubject())

	//---------------------------------------------------------------------------
	// Request Pattern Resources
	//---------------------------------------------------------------------------

	case *appb.Appointment:
		// TODO(b/240690479): Appointments are a request-pattern type that supports
		// multiple participants, which may be of type Practitioner, Patient, etc.
		// This means we may have to support multiple patient IDs. Currently we
		// assume only 1 patient by searching for and returning the first patient we
		// discover.
		for _, participant := range res.GetParticipant() {
			if id := participant.GetActor().GetPatientId().GetValue(); id != "" {
				return id, nil
			}
		}
		return "", fmt.Errorf("%w from %T", ErrExtractingPatientID, res)
	case *arpb.AppointmentResponse:
		return idOrError(res.GetActor())
	case *cppb.CarePlan:
		return idOrError(res.GetSubject())
	case *clpb.Claim:
		return idOrError(res.GetPatient())
	case *crpb.CommunicationRequest:
		return idOrError(res.GetSubject())
	case *cerpb.CoverageEligibilityRequest:
		return idOrError(res.GetPatient())
	case *drpb.DeviceRequest:
		return idOrError(res.GetSubject())
	case *erpb.EnrollmentRequest:
		return idOrError(res.GetCandidate())
	case *irpb.ImmunizationRecommendation:
		return idOrError(res.GetPatient())
	case *mrpb.MedicationRequest:
		return idOrError(res.GetSubject())
	case *nopb.NutritionOrder:
		return idOrError(res.GetPatient())
	case *rgpb.RequestGroup:
		return idOrError(res.GetSubject())
	case *srpb.ServiceRequest:
		return idOrError(res.GetSubject())
	case *surpb.SupplyRequest:
		// Supply requests may contain PatientID in two possible locations: either as
		// the source of the request, or as the destination for the request
		if id := res.GetDeliverTo().GetPatientId().GetValue(); id != "" {
			return id, nil
		} else if id := res.GetRequester().GetPatientId().GetValue(); id != "" {
			return id, nil
		}
		return "", fmt.Errorf("%w from %T", ErrExtractingPatientID, res)
	case *vppb.VisionPrescription:
		return idOrError(res.GetPatient())
	default:
		return "", fmt.Errorf("%w: %T", ErrUnsupportedType, res)
	}
}

// Reference creates a literal Patient reference.
// This replaces verily-go-fhir/protohelpers.PatientReference.
func Reference(patientID string) *dtpb.Reference {
	return &dtpb.Reference{
		Type: fhir.URI(resource.Patient.String()),
		Reference: &dtpb.Reference_PatientId{
			PatientId: &dtpb.ReferenceId{
				Value: patientID,
			},
		},
	}
}
