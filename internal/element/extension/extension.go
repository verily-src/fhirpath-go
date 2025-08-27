package extension

import (
	"errors"
	"fmt"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/protofields"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ErrInvalidValueX is an error reported if trying to create an
// Extension from an invalid data-type. Most FHIR Elements are supported, with
// only a few notable exceptions. See the definition of the ValueX constraint to
// see which elements are valid inputs.
var ErrInvalidValueX = errors.New("invalid extension ValueX type")

// ValueX is a constraint that enumerates all the valid 'DataType'
// objects that can be used for extensions. This is used to constrain the
// `Extension` function so that it cannot fail.
//
// This type is exported so that consumers may also expose the same generic
// requirements to their clients.
//
// For more information on these types, see the definition of Extension's
// `value` fields here: https://www.hl7.org/fhir/r4/extensibility.html#Extension
type ValueX interface {
	fhir.Element
	*dtpb.Base64Binary |
		*dtpb.Boolean |
		*dtpb.Canonical |
		*dtpb.Code |
		*dtpb.Date |
		*dtpb.DateTime |
		*dtpb.Decimal |
		*dtpb.Id |
		*dtpb.Instant |
		*dtpb.Integer |
		*dtpb.Markdown |
		*dtpb.Oid |
		*dtpb.PositiveInt |
		*dtpb.String |
		*dtpb.Time |
		*dtpb.UnsignedInt |
		*dtpb.Uri |
		*dtpb.Url |
		*dtpb.Uuid |
		*dtpb.Address |
		*dtpb.Age |
		*dtpb.Annotation |
		*dtpb.Attachment |
		*dtpb.CodeableConcept |
		*dtpb.Coding |
		*dtpb.ContactPoint |
		*dtpb.Count |
		*dtpb.Distance |
		*dtpb.Duration |
		*dtpb.HumanName |
		*dtpb.Identifier |
		*dtpb.Money |
		*dtpb.Period |
		*dtpb.Quantity |
		*dtpb.Range |
		*dtpb.Ratio |
		*dtpb.Reference |
		*dtpb.SampledData |
		*dtpb.Signature |
		*dtpb.Timing |
		*dtpb.ContactDetail |
		*dtpb.Contributor |
		*dtpb.DataRequirement |
		*dtpb.Expression |
		*dtpb.ParameterDefinition |
		*dtpb.RelatedArtifact |
		*dtpb.TriggerDefinition |
		*dtpb.UsageContext |
		*dtpb.Dosage
}

// New creates an New extension object from a concrete, and legal,
// extension type. Unlike `FromElement`, this function cannot fail since the
// type has been checked ot be valid extension type with Go-1.18 constraints.
func New[T ValueX](uri string, element T) *dtpb.Extension {
	ext, err := FromElement(uri, element)
	if err != nil {
		// This branch is unreachable due to the type constraint.
		// Tested in extension_test.go
		panic(err)
	}
	return ext
}

// FromElement creates an extension from the specified datatype resource.
// If the datatype is not a valid input, this function returns an error.
//
// For convenience functions that cannot return an error, see the various
// `Extension*` functions.
func FromElement(uri string, element fhir.Element) (*dtpb.Extension, error) {
	if element == nil {
		return nil, fmt.Errorf("no value provided to datatype extension")
	}
	name := protofields.DescriptorName(element)

	fields, ok := protofields.Elements[name]
	if !ok || fields.Extension.ValueX == nil {
		return nil, fmt.Errorf("extension %v: %w", name, ErrInvalidValueX)
	}

	valueX := &dtpb.Extension_ValueX{}
	reflect := valueX.ProtoReflect()
	reflect.Set(fields.Extension.ValueX, protoreflect.ValueOfMessage(element.ProtoReflect()))

	extension := &dtpb.Extension{
		Url:   fhir.URI(uri),
		Value: valueX,
	}
	return extension, nil
}

// Unwrap will return the wrapped Element in this extension, if one
// exists. If there is none, or if extension is nil, this returns nil.
func Unwrap(extension *dtpb.Extension) fhir.Element {
	if extension == nil || extension.Value.GetChoice() == nil {
		return nil
	}
	const choiceField = "choice"

	reflect := extension.Value.ProtoReflect()
	descriptor := reflect.Descriptor()
	oneof := descriptor.Oneofs().ByName(choiceField)
	field := reflect.WhichOneof(oneof)
	message := reflect.Get(field).Message()

	// Extensions can only have a limited number of extension values, all of which
	// satisfy 'Element'. This assertion cannot fail unless an invalid oneof value
	// has been provided, which violates protobuf definitions (and could never
	// legally happen during transport).
	return message.Interface().(fhir.Element)
}

// Clear removes all extensions from the specified FHIR Element
// or Resource.
func Clear(ext fhir.Extendable) {
	if ext == nil {
		return
	}

	message := ext.ProtoReflect()
	field := message.Descriptor().Fields().ByName("extension")
	message.Clear(field)
}

// Upsert always replaces or inserts the extension by the URL.
func Upsert(ext fhir.Extendable, extension *dtpb.Extension) {
	if ext == nil {
		panic("No extendable object specified for Upsert; ext is nil.")
	}

	for _, currExt := range ext.GetExtension() {
		if currExt.GetUrl().GetValue() == extension.GetUrl().GetValue() {
			currMessage := currExt.ProtoReflect()
			valueDesc := currMessage.Descriptor().Fields().ByName("value")
			currMessage.Set(valueDesc, protoreflect.ValueOfMessage(extension.GetValue().ProtoReflect()))
			return
		}
	}

	AppendInto(ext, extension)
}

// SetByURL always remove all extensions with url,
// and creates N extensions with url and values...
func SetByURL[T ValueX](ext fhir.Extendable, url string, values ...T) {
	if ext == nil {
		panic("No extendable object specified for SetByURL; ext is nil.")
	}

	var newExtensionList []*dtpb.Extension

	for _, currExt := range ext.GetExtension() {
		if currExt.GetUrl().GetValue() != url {
			newExtensionList = append(newExtensionList, currExt)
		}
	}

	for _, val := range values {
		newExtensionList = append(newExtensionList, New(url, val))
	}

	Overwrite(ext, newExtensionList...)
}

// ClearByURL removes all extensions with the specified URL from the
// specified FHIR Element or Resource.
func ClearByURL(ext fhir.Extendable, url string) {
	if ext == nil {
		panic("No extendable object specified for ClearByURL; ext is nil.")
	}

	var newExtensionList []*dtpb.Extension
	for _, currExt := range ext.GetExtension() {
		if currExt.GetUrl().GetValue() != url {
			newExtensionList = append(newExtensionList, currExt)
		}
	}

	Overwrite(ext, newExtensionList...)
}

// Overwrite modifies a Resource or Element in-place to overwrite all of the
// extensions in the given object with the provided ones.
//
// This function can only be called with extendable resources or data-types, and
// thus does not have any error-cases to surface back to the caller.
//
// This function will panic if ext is nil.
func Overwrite(ext fhir.Extendable, extensions ...*dtpb.Extension) {
	updateExtensionsIn(ext, protoreflect.Message.NewField, extensions...)
}

// AppendInto appends extensions into the specified Resource or Element.
// Unlike most "append" operations, this mutates the source object rather than
// copying and returning a new changed object. This is done to prevent expensive
// copy-operations on large FHIR resources that are being extended.
//
// This function can only be called with extendable resources or data-types, and
// thus does not have any error-cases to surface back to the caller.
//
// This function will panic if ext is nil.
func AppendInto(ext fhir.Extendable, extensions ...*dtpb.Extension) {
	updateExtensionsIn(ext, protoreflect.Message.Mutable, extensions...)
}

// updateExtensionsIn updates the state of extensions inside of the extendable
// resource, using the specified field accessor to help.
func updateExtensionsIn(ext fhir.Extendable, get protofields.FieldToValueFunc, extensions ...*dtpb.Extension) {
	if ext == nil {
		panic("No extendable object specified for updateExtensionsIn; ext is nil.")
	}
	message := ext.ProtoReflect()
	field := message.Descriptor().Fields().ByName("extension")
	extensionsList := get(message, field).List()
	for _, ext := range extensions {
		val := protoreflect.ValueOfMessage(ext.ProtoReflect())
		extensionsList.Append(val)
	}
	message.Set(field, protoreflect.ValueOfList(extensionsList))
}

// FindByURL returns all extensions whose URL matches the specified system.
// If no matching extensions are found, this returns an empty slice.
func FindByURL(extensions []*dtpb.Extension, url string) []*dtpb.Extension {
	var matches []*dtpb.Extension

	for _, ext := range extensions {
		if ext.GetUrl().GetValue() == url {
			matches = append(matches, ext)
		}
	}

	return matches
}
