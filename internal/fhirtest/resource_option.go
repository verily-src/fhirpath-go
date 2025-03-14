package fhirtest

import (
	"fmt"
	"testing"

	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	drpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/document_reference_go_proto"
	listpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/list_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	pepb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/person_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// WithResourceModification applies a modifier on a resource. Resource and
// modifier must be the same type.
func WithResourceModification[T fhir.Resource](modifier func(T)) ResourceOption {
	return func(t *testing.T, res fhir.Resource) {
		t.Helper()
		val, ok := res.(T)
		if !ok {
			var modifierType T
			t.Fatalf("Modifier type != applied resource type: (%s) != (%s)", resource.TypeOf(modifierType), resource.TypeOf(res))
		}
		modifier(val)
	}
}

// WithJSONField creates a ResourceOpt for setting up the specified proto-field
// with the given value in JSON format.
//
// This will fail tests if the field is invalid, or if the JSON does not
// unmarshal correctly.
func WithJSONField(field, value string) ResourceOption {
	return func(t *testing.T, r fhir.Resource) {
		t.Helper()
		msg := r.ProtoReflect()
		field := getProtoField(t, msg, field)

		s := msg.Get(field).Message().New().Interface()
		if err := protojson.Unmarshal([]byte(value), s); err != nil {
			t.Fatalf("Invalid JSON '%v': %v", value, err)
		}

		msg.Set(field, protoreflect.ValueOfMessage(s.ProtoReflect()))
	}
}

// WithCodeField creates a ResourceOpt for setting up the specified proto-field
// with the value of a resource code.
//
// This can be used as a short-hand for constructing arbitrary strongly-typed
// code-fields from a simple string, such as `WithCodeField("status", "DRAFT")`
// for producing the "Draft" state for a CanonicalResource status.
//
// This will fail tests if the field is invalid, or if the JSON does not
// unmarshal correctly.
func WithCodeField(field, value string) ResourceOption {
	// All "Code" objects are serialized in JSON as just { "value": "<value>" };
	// using this property enables us to set different codes for the fhir protos
	// which otherwise use strong-types (e.g. Questionnaire.Status is a different
	// type from Account.Status, but both use the underlying PublicationCode string.)
	return WithJSONField(field, fmt.Sprintf(`{ "value": "%v" }`, value))
}

// WithProtoField creates a ResourceOpt for setting up the specified proto-field
// with the value set in the proto message.
//
// This will fail tests if the field is invalid, or panic if message is the
// wrong input type.
func WithProtoField(field string, message proto.Message) ResourceOption {
	return func(t *testing.T, r fhir.Resource) {
		t.Helper()
		msg := r.ProtoReflect()
		field := getProtoField(t, msg, field)

		msg.Set(field, protoreflect.ValueOfMessage(message.ProtoReflect()))
	}
}

// WithRepeatedProtoField creates a ResourceOpt for setting up the specified
// repeated proto-field with the values set in the supplied proto messages.
//
// This will fail tests if the field is invalid, or panic if message is the
// wrong input type.
func WithRepeatedProtoField(field string, messages ...proto.Message) ResourceOption {
	return func(t *testing.T, r fhir.Resource) {
		t.Helper()
		msg := r.ProtoReflect()
		field := getProtoField(t, msg, field)
		list := msg.Mutable(field).List()
		for _, message := range messages {
			list.Append(protoreflect.ValueOfMessage(message.ProtoReflect()))
		}
		msg.Set(field, protoreflect.ValueOfList(list))
	}
}

// WithGeneratedIdentifier creates a ResourceOpt for automatically
// generating an Identifier for the given resource with the specified system.
// If the resource implements GetIdentifier(), an Identifier will be generated
// and added.
// If the resource does not, then we will fail the test.
func WithGeneratedIdentifier(system string) ResourceOption {
	return func(t *testing.T, r fhir.Resource) {
		// set Identifier if available
		if cast, ok := r.(resource.HasGetIdentifierList); ok {
			ids := []*dtpb.Identifier{generateIdentifier(system)}
			setIdentifierList(cast, ids)
		} else if cast, ok := r.(resource.HasGetIdentifierSingle); ok {
			id := generateIdentifier(system)
			setIdentifier(cast, id)
		} else {
			t.Errorf("WithGeneratedIdentifier: invalid resource type %v has no GetIdentifier()", r)
		}
	}
}

// WithCode creates a ResourceOpt for updating the code field in the resource.
func WithCode(code string) ResourceOption {
	return WithProtoField("code", fhir.CodeableConcept(code))
}

// WithContent creates a ResourceOpt for updating the content field in the resource.
func WithContent(attachmentTitle, attachmentURL string) ResourceOption {
	return WithRepeatedProtoField("content",
		&drpb.DocumentReference_Content{
			Attachment: &dtpb.Attachment{
				Title: fhir.String(attachmentTitle),
				Url:   fhir.URL(attachmentURL),
			},
		},
	)
}

// WithDerivedFrom creates a ResourceOpt for updating the derived_from field in the resource.
func WithDerivedFrom(ref *dtpb.Reference) ResourceOption {
	return WithRepeatedProtoField("derived_from", ref)
}

// WithEntry creates a ResourceOpt for updating the entry field in the resource.
func WithEntry(ref *dtpb.Reference) ResourceOption {
	return WithRepeatedProtoField("entry", &listpb.List_Entry{Item: ref})
}

// WithHumanName creates a ResourceOpt for updating the name field in the resource with a human name.
func WithHumanName(family, given string) ResourceOption {
	return WithRepeatedProtoField("name",
		&dtpb.HumanName{
			Family: fhir.String(family),
			Given:  []*dtpb.String{fhir.String(given)}})
}

// WithId creates a ResourceOpt for updating the id field in the resource.
func WithId(id *dtpb.Id) ResourceOption {
	return WithProtoField("id", id)
}

// WithIndividual creates a ResourceOpt for updating the individual field in the resource.
func WithIndividual(ref *dtpb.Reference) ResourceOption {
	return WithProtoField("individual", ref)
}

// WithMode creates a ResourceOpt for updating the mode field in the resource.
func WithMode(code string) ResourceOption {
	return WithCodeField("mode", code)
}

// WithName creates a ResourceOpt for updating the name field in the resource.
func WithName(name string) ResourceOption {
	return WithCodeField("name", name)
}

// WithPartOf creates a ResourceOpt for updating the partof field in the resource.
func WithPartOf(ref *dtpb.Reference) ResourceOption {
	return WithProtoField("part_of", ref)
}

// WithPatient creates a ResourceOpt for updating the patient field in the resource.
func WithPatient(ref *dtpb.Reference) ResourceOption {
	return WithProtoField("patient", ref)
}

// WithPatientLink creates a ResourceOpt for updating the link.target field in the resource.
func WithPatientLink(ref *dtpb.Reference, linkTypeCode cpb.LinkTypeCode_Value) ResourceOption {
	return WithRepeatedProtoField("link",
		&ppb.Patient_Link{
			Other: ref,
			Type: &ppb.Patient_Link_TypeCode{
				Value: linkTypeCode,
			},
		})
}

// WithPersonLink creates a ResourceOpt for updating the link.target field in the resource.
func WithPersonLink(refs ...*dtpb.Reference) ResourceOption {
	personLinks := []proto.Message{}
	for _, ref := range refs {
		personLinks = append(personLinks, &pepb.Person_Link{Target: ref})
	}
	return WithRepeatedProtoField("link", personLinks...)
}

// WithStatus creates a ResourceOpt for updating the status field in the resource.
func WithStatus(code string) ResourceOption {
	return WithCodeField("status", code)
}

// WithStudy creates a ResourceOpt for updating the study field in the resource.
func WithStudy(ref *dtpb.Reference) ResourceOption {
	return WithProtoField("study", ref)
}

// WithSubject creates a ResourceOpt for updating the subject field in the resource.
func WithSubject(ref *dtpb.Reference) ResourceOption {
	return WithProtoField("subject", ref)
}

// WithTarget creates a ResourceOpt for updating the target field in the resource.
func WithTarget(ref *dtpb.Reference) ResourceOption {
	return WithRepeatedProtoField("target", ref)
}

// WithType creates a ResourceOpt for updating the target field in the resource.
func WithType(code *dtpb.CodeableConcept) ResourceOption {
	return WithRepeatedProtoField("type", code)
}

// generateIdentifier generates a (stable) random Identifier with the given system
func generateIdentifier(system string) *dtpb.Identifier {
	return &dtpb.Identifier{
		System: &dtpb.Uri{Value: system},
		Value:  &dtpb.String{Value: stableRandomID().Value},
	}
}
