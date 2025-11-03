package protofields

import (
	apb "github.com/google/fhir/go/proto/google/fhir/proto/annotations_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	bcrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"github.com/iancoleman/strcase"
	"github.com/verily-src/fhirpath-go/internal/slices"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// FieldToValueFunc is a function converting a protoreflect.Message's field-descriptor
// into a Value.
//
// This corresponds to either `protoreflect.Message.Mutable` or `protoreflect.Message.NewField`.
//
// This type is an implementation-detail shared through various APIs that support
// a mix of appending or overwriting sequences.
type FieldToValueFunc func(protoreflect.Message, protoreflect.FieldDescriptor) protoreflect.Value

// ResourceFieldRefs is a container of back-references for a specified Resource.
// This enables finding one-ofs that can work with the resource it represents,
// for easier/faster lookup than waiting on protoreflect APIs.
type ResourceFieldRefs struct {
	// ContainedResource contains back-references for the ContainedResource to
	// this type
	ContainedResource struct {
		// Resource is the ContainedResource.Resource field that corresponds to
		// this resource.
		//
		// This field cannot be nil, as all R4 resources are valid contained-resource
		// types.
		Resource protoreflect.FieldDescriptor
	}

	DummyResource proto.Message

	// New is a function that will create a new instance of this FHIR Resource.
	New func() proto.Message
}

// ElementFieldRefs is a container of back-references for a specified Element.
// This enables finding one-ofs that can work with the resource it represents,
// for easier/faster lookup than waiting on protoreflect APIs.
type ElementFieldRefs struct {
	// Extension contains back-references for the Extension to this type
	Extension struct {
		// ValueX is the Extension.ValueX field that corresponds to
		// this resource.
		//
		// This field cannot be nil, as all R4 resources are valid contained-resource
		// types.
		ValueX protoreflect.FieldDescriptor
	}

	// New is a function that will create a new instance of this FHIR element.
	New func() proto.Message
}

var (
	// Resources is a map of all resource names to their corresponding field references.
	Resources map[string]*ResourceFieldRefs

	// Elements is a map of all element names to their corresponding field references.
	Elements map[string]*ElementFieldRefs
)

// IsValidResourceType checks that the given name is a valid resource name
func IsValidResourceType(name string) bool {
	_, ok := Resources[name]
	return ok
}

// IsValidElementType checks that the given name is a valid element name
func IsValidElementType(name string) bool {
	_, ok := Elements[name]
	return ok
}

// TypeToContainedResourceOneOfFieldName converts a resource type name into their
// respective ContainedResource "OneOf" field.
func TypeToContainedResourceOneOfFieldName(resource string) protoreflect.Name {
	return protoreflect.Name(toSnakeCase(resource))
}

// UnwrapOneofField obtains the underlying Message for "Oneof" elements
// contained in fields with the given fieldName. Returns nil if the input message
// doesn't have the given field, or if the Oneof descriptor is unpopulated.
func UnwrapOneofField(element proto.Message, fieldName string) proto.Message {
	message := element.ProtoReflect()
	oneOfDescriptor := message.Descriptor().Oneofs().ByName(protoreflect.Name(fieldName))
	if oneOfDescriptor == nil {
		return nil
	}
	fd := message.WhichOneof(oneOfDescriptor)
	if fd == nil {
		return nil
	}
	return message.Get(fd).Message().Interface()
}

// IsCodeField returns true if the message represents a FHIR code type.
// Codes with enum values and string values are both considered valid.
func IsCodeField(message proto.Message) bool {
	reflect := message.ProtoReflect()
	if !proto.HasExtension(reflect.Descriptor().Options(), apb.E_FhirValuesetUrl) {
		return false
	}
	field := reflect.Descriptor().Fields().ByName(protoreflect.Name("value"))
	if field != nil {
		allowedKinds := []protoreflect.Kind{protoreflect.EnumKind, protoreflect.StringKind}
		isValidFieldType := slices.Includes(allowedKinds, field.Kind())
		return isValidFieldType
	}
	return false
}

// StringValueFromCodeField gets the Field Descriptor of a message that
// represents a FHIR Code type. Returns the string value of the enum or
// string value of the Code, along with a boolean flag representing
// whether or not the input is a code type.
func StringValueFromCodeField(message proto.Message) (string, bool) {
	if IsCodeField(message) {
		reflect := message.ProtoReflect()
		field := reflect.Descriptor().Fields().ByName(protoreflect.Name("value"))
		if field.Kind() == protoreflect.EnumKind {
			enum := reflect.Get(field).Enum()
			code := string(field.Enum().Values().ByNumber(enum).Name())
			return strcase.ToKebab(code), true
		}
		if field.Kind() == protoreflect.StringKind {
			return reflect.Get(field).String(), true
		}
	}
	return "", false
}

// Field is a struct containing both the Value and FieldDescriptor for a proto field.
type Field struct {
	Value      protoreflect.Value
	Descriptor protoreflect.FieldDescriptor
}

// GetField retrieves the proto field of the specified name from the message.
func GetField(message proto.Message, name string) (*Field, bool) {
	fieldName := strcase.ToSnake(name)
	msg := message.ProtoReflect()
	descriptor := msg.Descriptor()
	field := descriptor.Fields().ByName(protoreflect.Name(fieldName))
	if field == nil {
		return nil, false
	}

	value := msg.Get(field)
	return &Field{
		Value:      value,
		Descriptor: field,
	}, true
}

func getContainedResourceOneOf(message proto.Message) protoreflect.FieldDescriptor {
	cr := (*bcrpb.ContainedResource)(nil)

	name := DescriptorName(message)
	fieldName := TypeToContainedResourceOneOfFieldName(name)
	return cr.ProtoReflect().Descriptor().Fields().ByName(fieldName)
}

// typeToExtensionFieldName converts a data-type name into their expected
// Extension ValueX field name.
func typeToExtensionFieldName(name string) protoreflect.Name {
	fieldName := toSnakeCase(name)
	if fieldName == "string" {
		// The protobufs use "string_value" rather than "string" because "string" is
		// a keyword.
		fieldName = "string_value"
	}
	return protoreflect.Name(fieldName)
}

func getExtensionValueX(message proto.Message) protoreflect.FieldDescriptor {
	valueX := (*dtpb.Extension_ValueX)(nil)

	reflect := valueX.ProtoReflect()
	name := DescriptorName(message)
	fieldName := typeToExtensionFieldName(name)
	return reflect.Descriptor().Fields().ByName(fieldName)
}

func newProto(msg protoreflect.ProtoMessage) func() proto.Message {
	return func() proto.Message {
		return msg.ProtoReflect().New().Interface()
	}
}

func init() {
	Resources = make(map[string]*ResourceFieldRefs)
	Elements = make(map[string]*ElementFieldRefs)

	for _, msg := range dummyResources {
		name := DescriptorName(msg)
		fields := &ResourceFieldRefs{}
		fields.ContainedResource.Resource = getContainedResourceOneOf(msg)
		fields.New = newProto(msg)
		fields.DummyResource = msg
		Resources[name] = fields
	}
	for _, msg := range dummyElements {
		name := DescriptorName(msg)
		fields := &ElementFieldRefs{}
		fields.New = newProto(msg)
		fields.Extension.ValueX = getExtensionValueX(msg)
		Elements[name] = fields
	}
}
