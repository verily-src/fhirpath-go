package fhirtest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	binpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/binary_go_proto"
	conpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/consent_go_proto"
	dpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/device_go_proto"
	drpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/document_reference_go_proto"
	listpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/list_go_proto"
	locpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/location_go_proto"
	opb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/observation_go_proto"
	orgpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/organization_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	pepb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/person_go_proto"
	qrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/questionnaire_response_go_proto"
	rstudpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/research_study_go_proto"
	rsubpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/research_subject_go_proto"
	vrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/verification_result_go_proto"
	"github.com/verily-src/fhirpath-go/internal/element/reference"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/protofields"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Shared constants for test resources.
const (
	patientURN       = "urn:uuid:patient"
	researchStudyURN = "urn:uuid:researchStudy"
)

// Resources is a map of all resource-names to an instance of that resource type.
var Resources map[string]fhir.Resource

// DomainResources is a map of all domain-resource-names to an instance of that
// domain-resource type.
var DomainResources map[string]fhir.DomainResource

// CanonicalResources is a map of all canonical-resource-names to an instance of
// that canonical-resource type.
var CanonicalResources map[string]fhir.CanonicalResource

// MetadataResources is a map of all metadata-resource-names to an instance of
// that metadata-resource type
var MetadataResources map[string]fhir.MetadataResource

// ResourceOption is an option type acting on resources for setting up a resource.
type ResourceOption func(*testing.T, fhir.Resource)

// NewResource creates a dummy resource object for the purposes of testing
// of type `resourceType. If `resourceType` does not name a valid R4 FHIR resource,
// this will fail the calling test.
func NewResource(t *testing.T, resourceType resource.Type, opts ...ResourceOption) fhir.Resource {
	t.Helper()

	if val, ok := Resources[string(resourceType)]; ok {
		object := proto.Clone(val).(fhir.Resource)
		ReplaceMeta(object, StableRandomMeta())
		UpdateID(object, stableRandomID().Value)
		for _, opt := range opts {
			opt(t, object)
		}
		return object
	}

	// This should be unreachable in valid code
	t.Fatalf("No resource of type %v found. Please specify a valid resource, or update this package.", resourceType)
	return nil
}

// NewResourceOf creates a dummy resource object for the purposes of testing
// of type `T`. If `T` is not a valid FHIR R4 resource, this will fail testing.
func NewResourceOf[T fhir.Resource](t *testing.T, opts ...ResourceOption) T {
	t.Helper()
	var res T
	resourceType := resource.TypeOf(res)

	return NewResource(t, resourceType, opts...).(T)
}

// NewResourceFromBase returns a copy of a resource modified by given options. The
// input resource remains unmodified. Useful for tweaking an existing resource
// for test inputs.
func NewResourceFromBase(t *testing.T, resource fhir.Resource, opts ...ResourceOption) fhir.Resource {
	t.Helper()

	modifiedResource := proto.Clone(resource).(fhir.Resource)
	for _, opt := range opts {
		opt(t, modifiedResource)
	}

	return modifiedResource
}

// NewBinary returns a test Binary resource with its required field populated:
//
//	content_type = image/jpeg
func NewBinary(t *testing.T, opts ...ResourceOption) *binpb.Binary {
	finalOpts := []ResourceOption{
		WithCodeField("content_type", "image/jpeg"),
	}
	finalOpts = append(finalOpts, opts...)
	return NewResourceOf[*binpb.Binary](t, finalOpts...)
}

// NewConsent returns a test Consent resource with its required fields populated:
//
//	category = my-code-text
//	policy_rule = my-code-text
//	status = ACTIVE
//	scope = my-code-text
func NewConsent(t *testing.T, opts ...ResourceOption) *conpb.Consent {
	finalOpts := []ResourceOption{
		WithRepeatedProtoField("category", fhir.CodeableConcept("my-code-text")),
		WithProtoField("policy_rule", fhir.CodeableConcept("my-code-text")),
		WithCodeField("status", "ACTIVE"),
		WithProtoField("scope", fhir.CodeableConcept("my-code-text")),
	}
	finalOpts = append(finalOpts, opts...)
	return NewResourceOf[*conpb.Consent](t, finalOpts...)
}

// NewDevice returns a test Device resource.
func NewDevice(t *testing.T, opts ...ResourceOption) *dpb.Device {
	return NewResourceOf[*dpb.Device](t, opts...)
}

// NewDocumentReference returns a test DocumentReference resource with one of its required fields populated:
//
//	status = CURRENT
//
// Note: Content is also required, but is left out because there is no default title or uri that can suffice.
func NewDocumentReference(t *testing.T, opts ...ResourceOption) *drpb.DocumentReference {
	// Default options
	finalOpts := []ResourceOption{
		WithCodeField("status", "CURRENT"),
	}
	finalOpts = append(finalOpts, opts...)
	return NewResourceOf[*drpb.DocumentReference](t, finalOpts...)
}

// NewList returns a test List resource with its required fields populated:
//
//	status = CURRENT
//	mode = WORKING
func NewList(t *testing.T, opts ...ResourceOption) *listpb.List {
	// Default options
	finalOpts := []ResourceOption{
		WithCodeField("status", "CURRENT"),
		WithCodeField("mode", "WORKING"),
	}
	// Add on opts provided, which will override the defaults.
	finalOpts = append(finalOpts, opts...)
	return NewResourceOf[*listpb.List](t, finalOpts...)
}

// NewLocation returns a test Location resource.
func NewLocation(t *testing.T, opts ...ResourceOption) *locpb.Location {
	return NewResourceOf[*locpb.Location](t, opts...)
}

// NewObservation returns a test Observation resource with its required fields populated:
//
//	code = my-code-text
//	status = PRELIMINARY
func NewObservation(t *testing.T, opts ...ResourceOption) *opb.Observation {
	// Default options
	finalOpts := []ResourceOption{
		WithProtoField("code", fhir.CodeableConcept("my-code-text")),
		WithCodeField("status", "PRELIMINARY"),
	}
	// Add on opts provided, which will override the defaults.
	finalOpts = append(finalOpts, opts...)
	return NewResourceOf[*opb.Observation](t, finalOpts...)
}

// NewOrganization returns a test Organization resource.
func NewOrganization(t *testing.T, opts ...ResourceOption) *orgpb.Organization {
	return NewResourceOf[*orgpb.Organization](t, opts...)
}

// NewPatient returns a test Patient resource.
func NewPatient(t *testing.T, opts ...ResourceOption) *ppb.Patient {
	return NewResourceOf[*ppb.Patient](t, opts...)
}

// NewPerson returns a test Person resource.
func NewPerson(t *testing.T, opts ...ResourceOption) *pepb.Person {
	return NewResourceOf[*pepb.Person](t, opts...)
}

// NewQuestionnaireResponse returns a test QuestionnaireResponse resource with its required fields populated:
//
//	subject = <Patient reference>
//	status = COMPLETED
func NewQuestionnaireResponse(t *testing.T, opts ...ResourceOption) *qrpb.QuestionnaireResponse {
	// Default options
	finalOpts := []ResourceOption{
		WithProtoField("subject", reference.Weak("Patient", patientURN)),
		WithCodeField("status", "COMPLETED"),
	}
	// Add on opts provided, which will override the defaults.
	finalOpts = append(finalOpts, opts...)
	return NewResourceOf[*qrpb.QuestionnaireResponse](t, finalOpts...)
}

// NewResearchStudy returns a test ResearchStudy resource with its required field populated:
//
//	status = APPROVED
func NewResearchStudy(t *testing.T, opts ...ResourceOption) *rstudpb.ResearchStudy {
	// Default options
	finalOpts := []ResourceOption{
		WithCodeField("status", "APPROVED"),
	}
	// Add on opts provided, which will override the defaults.
	finalOpts = append(finalOpts, opts...)
	return NewResourceOf[*rstudpb.ResearchStudy](t, finalOpts...)
}

// NewResearchSubject returns a test ResearchSubject resource with its required fields populated:
//
//	individual = <Patient reference>
//	status = ELIGIBLE
//	study = <ResearchStudy reference>
func NewResearchSubject(t *testing.T, opts ...ResourceOption) *rsubpb.ResearchSubject {
	// Default options
	finalOpts := []ResourceOption{
		WithProtoField("individual", reference.Weak("Patient", patientURN)),
		WithCodeField("status", "ELIGIBLE"),
		WithProtoField("study", reference.Weak("ResearchStudy", researchStudyURN)),
	}
	// Add on opts provided, which will override the defaults.
	finalOpts = append(finalOpts, opts...)
	return NewResourceOf[*rsubpb.ResearchSubject](t, finalOpts...)
}

// NewVerificationResult returns a test VerificationResult resource with its required fields populated:
//
//	status = VALIDATED
func NewVerificationResult(t *testing.T, opts ...ResourceOption) *vrpb.VerificationResult {
	// Default options
	finalOpts := []ResourceOption{
		WithCodeField("status", "VALIDATED"),
	}
	// Add on opts provided, which will override the defaults.
	finalOpts = append(finalOpts, opts...)
	return NewResourceOf[*vrpb.VerificationResult](t, finalOpts...)
}

func init() {
	Resources = make(map[string]fhir.Resource, len(protofields.Resources))
	DomainResources = map[string]fhir.DomainResource{}
	CanonicalResources = map[string]fhir.CanonicalResource{}
	MetadataResources = map[string]fhir.MetadataResource{}

	for name, res := range protofields.Resources {
		resource := res.New().(fhir.Resource)

		ReplaceMeta(resource, StableRandomMeta())
		UpdateID(resource, stableRandomID().Value)

		Resources[name] = resource
		// The various resource interfaces grow off of one another, so we can
		// minimize the number of checks here by leveraging this fact.
		if v, ok := resource.(fhir.MetadataResource); ok {
			setCanonicalFields(v)
			MetadataResources[name] = v
			CanonicalResources[name] = v
			DomainResources[name] = v
		} else if v, ok := resource.(fhir.CanonicalResource); ok {
			setCanonicalFields(v)
			CanonicalResources[name] = v
			DomainResources[name] = v
		} else if v, ok := resource.(fhir.DomainResource); ok {
			DomainResources[name] = v
		}
	}
}

// setCanonicalFields will automatically set canonical fields.
func setCanonicalFields(res fhir.CanonicalResource) {
	setCanonicalURL(res, fmt.Sprintf(
		"https://example.com/%v/%v",
		strings.ToLower(string(resource.TypeOf(res))),
		res.GetId().GetValue()),
	)
	setCanonicalVersion(res, "1.0")
	setCanonicalStatus(res, "DRAFT")
}

// setResourceProtoField will set the resource's field to the given proto message type.
func setResourceProtoField(resource fhir.Resource, fieldName string, value proto.Message) {
	reflect := resource.ProtoReflect()
	descriptor := reflect.Descriptor()
	field := descriptor.Fields().ByName(protoreflect.Name(fieldName))
	reflect.Set(field, protoreflect.ValueOfMessage(value.ProtoReflect()))
}

// setCanonicalURL will set a canonical URL for the given resource.
func setCanonicalURL(resource fhir.CanonicalResource, url string) {
	setResourceProtoField(resource, "url", &datatypes_go_proto.Uri{
		Value: url,
	})
}

// setCanonicalVersion will set a canonical version for the given resource.
func setCanonicalVersion(resource fhir.CanonicalResource, version string) {
	setResourceProtoField(resource, "version", &datatypes_go_proto.String{
		Value: version,
	})
}

func setCanonicalStatus(resource fhir.CanonicalResource, status string) {
	msg := resource.ProtoReflect()
	descriptor := msg.Descriptor()
	field := descriptor.Fields().ByName("status")

	s := msg.Get(field).Message().New().Interface()
	payload := fmt.Sprintf(`{ "value": "%v" }`, status)
	if err := protojson.Unmarshal([]byte(payload), s); err != nil {
		// This is tested to not happen, and is only privately called internal to
		// the setup of this library.
		panic(err)
	}

	msg.Set(field, protoreflect.ValueOfMessage(s.ProtoReflect()))
}

// setIdentifier sets the singleton Identifier on the given resource.
func setIdentifier(resource resource.HasGetIdentifierSingle, identifier *datatypes_go_proto.Identifier) {
	setResourceProtoField(resource, "identifier", identifier)
}

// setIdentifierList sets the Identifier list on the given resource.
func setIdentifierList(resource resource.HasGetIdentifierList, identifiers []*datatypes_go_proto.Identifier) {
	msg := resource.ProtoReflect()
	descriptor := msg.Descriptor()
	field := descriptor.Fields().ByName("identifier")

	list := msg.Mutable(field).List()
	for _, identifier := range identifiers {
		list.Append(protoreflect.ValueOfMessage(identifier.ProtoReflect()))
	}
	msg.Set(field, protoreflect.ValueOfList(list))
}

// getProtoField is a helper function for retrieving the named proto field
func getProtoField(t *testing.T, msg protoreflect.Message, fieldName string) protoreflect.FieldDescriptor {
	t.Helper()
	descriptor := msg.Descriptor()
	field := descriptor.Fields().ByName(protoreflect.Name(fieldName))
	if field == nil {
		t.Fatalf("Proto field '%v' not found", field)
	}
	return field
}
