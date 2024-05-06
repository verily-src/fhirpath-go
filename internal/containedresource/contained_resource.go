package containedresource

import (
	"errors"
	"fmt"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	bcrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"github.com/verily-src/fhirpath-go/internal/slices"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/element/identifier"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"github.com/verily-src/fhirpath-go/internal/protofields"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	ErrGenerateIfNoneExist error = errors.New("GenerateIfNoneExist()")
)

// Wrap creates a ContainedResource proto based on an existing FHIR proto.
// Usage:
//
//	patient := &Patient{...}
//	cr := Wrap(patient)
func Wrap(res fhir.Resource) *bcrpb.ContainedResource {
	if res == nil {
		return nil
	}
	name := resource.TypeOf(res)

	cr := &bcrpb.ContainedResource{}

	// This field exists for ALL valid FHIR R4 "Resource" types
	field := protofields.Resources[string(name)].ContainedResource.Resource

	// The only way this can fail is if `resource` was not a valid FHIR resource
	// to begin with (e.g. a custom mock, the r4 "Resource" proto (which is meant
	// to only be a baste type), or some other invalid resource). This has been
	// extensively tested against all valid r4 protos, and thus if this happens, it
	// is best to fail early and panic, rather than make this API virally return an
	// error that is impossible to experience in practice, and very difficult to
	// actually trigger.
	if field == nil {
		panic(fmt.Sprintf("Invalid resource with name %v specified in ContainResource", name))
	}

	cr.ProtoReflect().Set(field, protoreflect.ValueOfMessage(res.ProtoReflect()))
	return cr
}

// Unwrap will extract the underlying value contained in this
// resource, if there is one, and return it. This enables downstream callers to
// switch off of the resource type, or to perform type-conversions.
//
// This function is effectively the inverse of `ContainedResource`, such that the
// following assertion will always hold:
// `proto.Equal(fhirutil.Unwrap(fhirutil.ContainedResource(resource)), resource)`
func Unwrap(cr *bcrpb.ContainedResource) fhir.Resource {
	field := getContainedResourceOneOfField(cr)

	// If field is nil, it means we have no contained resource value set -- so
	// the only valid value to return while unwrapping is `nil`.
	if field == nil {
		return nil
	}
	ref := cr.ProtoReflect()
	message := ref.Get(field).Message()

	// All ContainedType values valid for the OneOf definition MUST satisfy the
	// resource interface. This assertion can not fail.
	return message.Interface().ProtoReflect().Interface().(fhir.Resource)
}

// TypeOf is a helper for getting the type-name of a contained resource.
//
// If the contained resource is nil, or the contained resource does not contain
// any resource, this function will panic.
func TypeOf(cr *bcrpb.ContainedResource) resource.Type {
	return resource.TypeOf(Unwrap(cr))
}

// ID is a helper for getting the ID of a contained resource.
//
// If the contained resource is nil, or the contained resource does not contain
// any resource, this will return an empty string.
func ID(cr *bcrpb.ContainedResource) string {
	return resource.ID(Unwrap(cr))
}

// VersionID gets the version-ID of the specified contained-resource as a string.
// If `nil` is provided, this returns an empty string.
func VersionID(cr *bcrpb.ContainedResource) string {
	return resource.VersionID(Unwrap(cr))
}

// URI is a helper for getting the URI of a contained-resource as a FHIR URI object.
// The URI is returned in the format Type/ID, e.g. Patient/123.
func URI(cr *bcrpb.ContainedResource) *dtpb.Uri {
	return resource.URI(Unwrap(cr))
}

// URIString is a helper for getting the URI of a contained-resource in
// string form. The URI is returned in the format Type/ID, e.g. Patient/123.
func URIString(cr *bcrpb.ContainedResource) string {
	return resource.URIString(Unwrap(cr))
}

// VersionedURI is a helper for getting the URI of a contained-resource as a
// FHIR URI object. The URI is returned in the format Type/ID/_history/VERSION.
func VersionedURI(cr *bcrpb.ContainedResource) *dtpb.Uri {
	return resource.VersionedURI(Unwrap(cr))
}

// VersionedURIString is a helper for getting the URI of a contained-resource in
// string form. The URI is returned in the format Type/ID/_history/VERSION.
func VersionedURIString(cr *bcrpb.ContainedResource) (string, bool) {
	return resource.VersionedURIString(Unwrap(cr))
}

// getContainedResourceOneOfField gets the field for the OneOf entry in the
// ContainedResource. This function returns nil if either the contained-resource
// is nil, or the contained-resource does not contain any resource.
func getContainedResourceOneOfField(cr *bcrpb.ContainedResource) protoreflect.FieldDescriptor {
	if cr == nil {
		return nil
	}
	if cr.GetOneofResource() == nil {
		return nil
	}

	// Get the active OneOf field, and return that value as an interface
	const oneofField = "oneof_resource"

	reflect := cr.ProtoReflect()
	descriptor := reflect.Descriptor()
	oneof := descriptor.Oneofs().ByName(oneofField)
	return reflect.WhichOneof(oneof)
}

// GenerateIfNoneExist generates an If-None-Exist header value using a single
// Identifier from the contained resource. The provided system is used to
// filter identifiers to only an identifier with a matching system.
//
// If no matching Identifier is found, return error if emptyIsErr is true, or
// return empty string and no error if emptyIsErr is false.
//
// The GCP FHIR store only supports atomic conditional operations on a single
// identifier, so this function returns an error if there are multiple
// identifiers matching the query.
//
// This is used for FHIR conditional create or other conditional methods.
// Untrusted data in Identifiers is escaped both for FHIR and for URL safety.
func GenerateIfNoneExist(cr *bcrpb.ContainedResource, system string, emptyIsErr bool) (string, error) {
	if cr == nil {
		return "", fmt.Errorf("%w: ContainedResource is nil", ErrGenerateIfNoneExist)
	}

	res := Unwrap(cr)
	if res == nil {
		return "", fmt.Errorf("%w: Unwrap() returned nil / no contained resource", ErrGenerateIfNoneExist)
	}

	identifiers, err := resource.GetIdentifierList(res)
	if err != nil {
		return "", err
	}

	found := slices.Filter(identifiers, func(id *dtpb.Identifier) bool {
		return id.GetSystem().GetValue() == system
	})

	if len(found) == 0 {
		if emptyIsErr {
			return "", fmt.Errorf("%w: found no Identifiers with system=%#v", ErrGenerateIfNoneExist, system)
		} else {
			return "", nil
		}
	}

	if len(found) > 1 {
		return "", fmt.Errorf("%w: found multiple Identifiers with system=%#v, want just one", ErrGenerateIfNoneExist, system)
	}

	return identifier.GenerateIfNoneExist(found[0]), nil
}
