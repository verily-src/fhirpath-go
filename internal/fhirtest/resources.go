package fhirtest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"github.com/verily-src/fhirpath-go/internal/protofields"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
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

// NewResource creates a dummy resource object for the purposes of testing
// of type `T. If `T` is not a valid FHIR R4 resource, this will fail testing.
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
			ids := []*datatypes_go_proto.Identifier{generateIdentifier(system)}
			setIdentifierList(cast, ids)
		} else if cast, ok := r.(resource.HasGetIdentifierSingle); ok {
			id := generateIdentifier(system)
			setIdentifier(cast, id)
		} else {
			t.Errorf("WithGeneratedIdentifier: invalid resource type %v has no GetIdentifier()", r)
		}
	}
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

// generateIdentifier generates a (stable) random Identifier with the given system
func generateIdentifier(system string) *datatypes_go_proto.Identifier {
	return &datatypes_go_proto.Identifier{
		System: &datatypes_go_proto.Uri{Value: system},
		Value:  &datatypes_go_proto.String{Value: stableRandomID().Value},
	}
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
