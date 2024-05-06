package resource

import (
	"errors"
	"fmt"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/protofields"
)

// ErrBadType is an error raised when a bad type is provided.
var ErrBadType = errors.New("bad resource type")

// Type is a FHIR Resource type object. This is similar to a reflect.Type that
// encodes its name identifier
//
// Type objects should never be constructed manually; rather, use the `CheckType`
// or `TypeOf` functions to get a valid type object. Invalid instances of Type
// may lead to unexpected implicit `panic` behavior, as any code consuming this
// is allowed to assume that `Type` always names a valid instance.
type Type string

// TypeOf gets the underlying type of the named resource.
//
// This function panics if resource is nil. Note that this is only an issue if
// the interface `fhir.Resource` is nil, *not* if the underlying resource is a
// pointer that is nil. E.g. the following holds true:
//
//	assert.True(resource.TypeOf((*ppb.Patient)(nil)) == resource.Patient)
func TypeOf(resource fhir.Resource) Type {
	if resource == nil {
		panic("TypeOf provided nil Resource")
	}
	return Type(resource.ProtoReflect().Descriptor().Name())
}

// NewType checks whether the string type name is a valid resource.Type instance.
// If it is, an instance of the type is returned. If the provided type is not a
// valid type, an ErrBadType is returned, and the type result is garbage.
//
// Note: This is case-sensitive, and expects CamelCase, just as the FHIR spec uses.
func NewType(resourceType string) (Type, error) {
	if !IsType(resourceType) {
		return "", fmt.Errorf("%w '%v'", ErrBadType, resourceType)
	}
	return Type(resourceType), nil
}

// String converts this Type into a string.
func (t Type) String() string {
	return string(t)
}

// New returns an instance of the FHIR Resource which this type names, using
// the provided options to toggle.
//
// This function will panic if this does not name a valid Resource Type.
func (t Type) New(opts ...Option) fhir.Resource {
	return New(t, opts...)
}

// URI returns a URI object containing the resource type name.
func (t Type) URI() *dtpb.Uri {
	return &dtpb.Uri{
		Value: string(t),
	}
}

// StructureDefinitionURI returns an absolute URI to the structure-definition
// URL.
func (t Type) StructureDefinitionURI() *dtpb.Uri {
	const baseURL = "http://hl7.org/fhir/StructureDefinition"

	return &dtpb.Uri{
		Value: fmt.Sprintf("%v/%v", baseURL, t),
	}
}

// IsType queries whether the given string names a Resource type.
//
// Note: This is case-sensitive, and expects CamelCase, jus as the FHIR spec uses.
func IsType(name string) bool {
	_, ok := protofields.Resources[name]
	return ok
}
