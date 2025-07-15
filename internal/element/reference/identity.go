package reference

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/iancoleman/strcase"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Source: https://hl7.org/fhir/r4/references.html#literal
// WATCHOUT: See ticket about anchored matches.
var restFHIRServiceBaseURLRegex = regexp.MustCompile(`^(http|https):\/\/([A-Za-z0-9\-\\\.\:\%\$]*\/)+$`)
var restFHIRServiceResourceURLRegex = regexp.MustCompile(`^((http|https):\/\/([A-Za-z0-9\-\\\.\:\%\$\_]*\/)+)?(Account|ActivityDefinition|AdverseEvent|AllergyIntolerance|Appointment|AppointmentResponse|AuditEvent|Basic|Binary|BiologicallyDerivedProduct|BodyStructure|Bundle|CapabilityStatement|CarePlan|CareTeam|CatalogEntry|ChargeItem|ChargeItemDefinition|Claim|ClaimResponse|ClinicalImpression|CodeSystem|Communication|CommunicationRequest|CompartmentDefinition|Composition|ConceptMap|Condition|Consent|Contract|Coverage|CoverageEligibilityRequest|CoverageEligibilityResponse|DetectedIssue|Device|DeviceDefinition|DeviceMetric|DeviceRequest|DeviceUseStatement|DiagnosticReport|DocumentManifest|DocumentReference|EffectEvidenceSynthesis|Encounter|Endpoint|EnrollmentRequest|EnrollmentResponse|EpisodeOfCare|EventDefinition|Evidence|EvidenceVariable|ExampleScenario|ExplanationOfBenefit|FamilyMemberHistory|Flag|Goal|GraphDefinition|Group|GuidanceResponse|HealthcareService|ImagingStudy|Immunization|ImmunizationEvaluation|ImmunizationRecommendation|ImplementationGuide|InsurancePlan|Invoice|Library|Linkage|List|Location|Measure|MeasureReport|Media|Medication|MedicationAdministration|MedicationDispense|MedicationKnowledge|MedicationRequest|MedicationStatement|MedicinalProduct|MedicinalProductAuthorization|MedicinalProductContraindication|MedicinalProductIndication|MedicinalProductIngredient|MedicinalProductInteraction|MedicinalProductManufactured|MedicinalProductPackaged|MedicinalProductPharmaceutical|MedicinalProductUndesirableEffect|MessageDefinition|MessageHeader|MolecularSequence|NamingSystem|NutritionOrder|Observation|ObservationDefinition|OperationDefinition|OperationOutcome|Organization|OrganizationAffiliation|Patient|PaymentNotice|PaymentReconciliation|Person|PlanDefinition|Practitioner|PractitionerRole|Procedure|Provenance|Questionnaire|QuestionnaireResponse|RelatedPerson|RequestGroup|ResearchDefinition|ResearchElementDefinition|ResearchStudy|ResearchSubject|RiskAssessment|RiskEvidenceSynthesis|Schedule|SearchParameter|ServiceRequest|Slot|Specimen|SpecimenDefinition|StructureDefinition|StructureMap|Subscription|Substance|SubstanceNucleicAcid|SubstancePolymer|SubstanceProtein|SubstanceReferenceInformation|SubstanceSourceMaterial|SubstanceSpecification|SupplyDelivery|SupplyRequest|Task|TerminologyCapabilities|TestReport|TestScript|ValueSet|VerificationResult|VisionPrescription)\/[A-Za-z0-9\-\.]{1,64}(\/_history\/[A-Za-z0-9\-\.]{1,64})?$`)

var (
	ErrInvalidAbsoluteURL           = errors.New("invalid absolute uri")
	ErrInvalidRelativeURI           = errors.New("invalid relative uri")
	ErrInvalidURL                   = errors.New("invalid url")
	ErrInvalidURI                   = errors.New("invalid reference uri")
	ErrFragmentMissingType          = errors.New("fragment reference missing type")
	ErrReferenceOneOfResourceNotSet = errors.New("reference.oneof_resource was not set")
)

// oneofReferenceDescriptor returns the FieldDescriptor for the oneof option that has been set
// within the Reference. An error is returned if no option is set.
func oneofReferenceDescriptor(x *dtpb.Reference) (protoreflect.FieldDescriptor, error) {
	msg := x.ProtoReflect()
	oneofDescriptor := msg.Descriptor().Oneofs().ByName("reference")
	fd := msg.WhichOneof(oneofDescriptor)
	if fd == nil {
		return nil, ErrReferenceOneOfResourceNotSet
	}
	return fd, nil
}

// IdentityOf returns a complete Identity (Type and ID always set,
// VersionID set if applicable) representing the given reference.
func IdentityOf(ref *dtpb.Reference) (*resource.Identity, error) {
	// Fragment
	if ref.GetFragment() != nil {
		if refType := ref.GetType(); refType != nil {
			return resource.NewIdentity(refType.GetValue(), ref.GetFragment().GetValue(), "")
		}
		return nil, ErrFragmentMissingType
	}

	// Absolute and Relative URIs
	if uri := ref.GetUri(); uri != nil {
		return IdentityFromURL(uri.GetValue())
	}
	return identityOfStrong(ref)
}

func identityOfStrong(ref *dtpb.Reference) (*resource.Identity, error) {
	fd, err := oneofReferenceDescriptor(ref)
	if err != nil {
		return nil, err
	}

	// All "reference" fields are named "[resource_type]_id" in the FHIR protos.
	// If we chop off "_id" and convert to camel-case, it's the resource type
	// name.
	resType := string(fd.Name())
	resType, _ = strings.CutSuffix(resType, "_id")
	resType = strcase.ToCamel(resType)

	m := ref.ProtoReflect().Get(fd).Message().Interface()
	refID, ok := m.(*dtpb.ReferenceId)
	if !ok {
		return nil, fmt.Errorf("unable to extract refID")
	}
	identID := refID.GetValue()
	identVersion := refID.GetHistory().GetValue()

	return resource.NewIdentity(resType, identID, identVersion)
}

// IdentityFromURL returns a complete Identity (Type and ID always
// set, VersionID set if applicable) representing the resource at the given
// relative or absolute URL, which may optionally include a _history component.
func IdentityFromURL(url string) (*resource.Identity, error) {
	lit, err := LiteralInfoFromURI(url)
	if err != nil {
		return nil, err
	}
	if lit.identity == nil {
		return nil, ErrInvalidURL
	}
	return lit.identity, nil
}

// IdentityFromAbsoluteURL returns a complete Identity (Type and ID always
// set, VersionID set if applicable) representing the resource at the given
// absolute URL, which may optionally include a _history component. It
// validates the url against FHIR-provided regex for FHIR REST servers.
// Example absolute: "https://healthcare.googleapis.com/v1/projects/${project}/locations/${location}/datasets/${dataset}/fhirStores/${fhirStore}/fhir/Patient/123/_history/abc"
// Example relative: "Patient/123/_history/abc"
func IdentityFromAbsoluteURL(url string) (*resource.Identity, error) {
	lit, err := LiteralInfoFromURI(url)
	if err != nil {
		return nil, err
	}
	if lit.identity == nil || lit.serviceBaseURL == "" {
		return nil, ErrInvalidAbsoluteURL
	}
	return lit.identity, nil
}

// IdentityFromRelativeURI returns a complete Identity (Type and ID always
// set, VersionID set if applicable) representing the resource at the given
// relative URI, which may optionally include a _history component.
func IdentityFromRelativeURI(uri string) (*resource.Identity, error) {
	uriParts := strings.Split(uri, "/")
	switch len(uriParts) {
	case 2:
		// e.g. Patient/123
		return resource.NewIdentity(uriParts[0], uriParts[1], "")
	case 4:
		// e.g. Patient/123/_history/abc
		if uriParts[2] != "_history" {
			break
		}
		return resource.NewIdentity(uriParts[0], uriParts[1], uriParts[3])
	}
	return nil, fmt.Errorf("%w: %s", ErrInvalidRelativeURI, uri)
}
