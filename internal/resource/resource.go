/*
Package resource contains utilities for working with abstract FHIR Resource
objects.
*/
package resource

import (
	"errors"
	"fmt"
	"strings"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/protofields"
	"github.com/verily-src/fhirpath-go/internal/resourceopt"
)

var (
	ErrGetIdentifierList = errors.New("GetIdentifierList()")
)

// Option is an option that may be supplied to updates or creations of Resource
// types.
type Option = resourceopt.Option

// NewFromString attempts to construct a resource of the specified string name type,
// using the specified options to construct it. This function returns an error
// if 'name' does not name a valid type.
func NewFromString(name string, opts ...Option) (fhir.Resource, error) {
	fields, ok := protofields.Resources[name]
	if !ok {
		return nil, fmt.Errorf("invalid resource name '%v'", name)
	}
	resource := fields.New().(fhir.Resource)

	return Update(resource, opts...), nil
}

// New constructs a new resource from the input type, using the specified
// options to construct it.
//
// This function assumes that Type is a validly-constructed type object.
// Failure to pass a valid type will result in a panic.
func New(name Type, opts ...Option) fhir.Resource {
	resource, err := NewFromString(string(name), opts...)
	if err != nil {
		// This is unreachable with validly-constructed Type objects
		panic(err)
	}
	return resource
}

// NewOf constructs a new resource of the named T resource type, using the
// specified options to construct it.
func NewOf[T fhir.Resource](opts ...Option) fhir.Resource {
	var t T
	return New(TypeOf(t), opts...)
}

// Update modifies the input resource in-place with the specified options.
func Update(res fhir.Resource, opts ...Option) fhir.Resource {
	return resourceopt.ApplyOptions(res, opts...)
}

// ID gets the ID of the specified resource as a string. If `nil` is
// provided, this returns an empty string.
func ID(resource fhir.Resource) string {
	if resource == nil {
		return ""
	}
	return resource.GetId().GetValue()
}

// VersionID gets the version-ID of the specified resource as a string.
// If `nil` is provided, this returns an empty string.
//
// This function on its own just simplifies the need of calling
// `GetMeta().GetVersionId().GetValue()` all the time.
func VersionID(resource fhir.Resource) string {
	if resource == nil {
		return ""
	}
	return resource.GetMeta().GetVersionId().GetValue()
}

// VersionETag pulls the "version" from the resource if it's an existing resource
// that was queried from a FHIR store; the version returned matches the ETag header
// returned by GET fhir-prefix/{resourceType}/{id} for this resource. This is used
// for optimistic locking on resources per https://hl7.org/fhir/http.html#concurrency
func VersionETag(r fhir.Resource) string {
	version := VersionID(r)
	if version == "" {
		return ""
	}
	return fmt.Sprintf(`W/"%s"`, version)
}

// URI is a helper for getting the URI of a resource as a URI object.
// The URI is returned in the format Type/ID, e.g. Patient/123.
//
// If the resource is nil, this will return a nil URI.
func URI(resource fhir.Resource) *dtpb.Uri {
	uri := URIString(resource)
	if uri == "" {
		return nil
	}
	return fhir.URI(uri)
}

// URIString is a helper for getting the URI of a resource in
// string form. The URI is returned in the format Type/ID, e.g. Patient/123.
//
// If the resource is nil, this will return an empty string.
func URIString(resource fhir.Resource) string {
	if resource == nil {
		return ""
	}
	id := resource.GetId().GetValue()
	return fmt.Sprintf("%v/%v", TypeOf(resource), id)
}

// Profiles is a helper for getting the fhir profiles of a resource in
// canonical form.
//
// If the resource is nil, this will return nil.
func Profiles(resource fhir.Resource) []*dtpb.Canonical {
	if resource == nil {
		return nil
	}
	return resource.GetMeta().GetProfile()
}

// ProfileStrings is a helper for getting the fhir profiles of a resource in
// string form.
//
// If the resource is nil, this will return nil.
func ProfileStrings(resource fhir.Resource) []string {
	if resource == nil {
		return nil
	}
	var profiles []string
	for _, profile := range resource.GetMeta().GetProfile() {
		profiles = append(profiles, profile.GetValue())
	}

	return profiles
}

// UnversionedProfileStrings is a helper for getting the unversioned fhir profiles of a resource in
// string form. The order of the profiles are preserved and duplicate profile base urls are removed.
//
// If the resource is nil, this will return nil.
func UnversionedProfileStrings(resource fhir.Resource) []string {
	profiles := ProfileStrings(resource)
	if profiles == nil {
		return nil
	}

	var unversionedProfiles []string
	seen := map[string]bool{}

	for _, profile := range profiles {
		unversionedProfile := strings.Split(profile, "|")[0]
		if !seen[unversionedProfile] {
			seen[unversionedProfile] = true
			unversionedProfiles = append(unversionedProfiles, unversionedProfile)
		}
	}
	return unversionedProfiles
}

// VersionedURI is a helper for getting the URI of a resource as a URI object.
// The URI is returned in the format Type/ID/_history/VERSION.
//
// If the resource is nil, this will return a nil URI.
func VersionedURI(resource fhir.Resource) *dtpb.Uri {
	uri, found := VersionedURIString(resource)
	if !found {
		return nil
	}
	return fhir.URI(uri)
}

// VersionedURIString is a helper for getting the URI of a resource in
// string form. The URI is returned in the format Type/ID/_history/VERSION.
//
// If the resource is nil, this will return an empty string.
func VersionedURIString(resource fhir.Resource) (string, bool) {
	if resource == nil {
		return "", false
	}
	vID := VersionID(resource)
	if vID == "" {
		return "", false
	}
	id := resource.GetId().GetValue()
	return fmt.Sprintf("%v/%v/_history/%v", TypeOf(resource), id, vID), true
}

// RemoveDuplicates finds all duplicates of resources -- determined by the
// same <resource>/<id>/<version-id> -- and removes them, returning an
// updated list of resources.
//
// Nil resources are skipped.
func RemoveDuplicates(resources []fhir.Resource) []fhir.Resource {
	deduper := map[string]struct{}{}

	result := make([]fhir.Resource, 0, len(resources))
	for _, res := range resources {
		if res == nil {
			continue
		}
		// Note: using this instead of VersionedURIString, since whether a version-id
		// exists or not will not affect the behavior here.
		key := fmt.Sprintf("%v/%v/%v", TypeOf(res), ID(res), res.GetMeta().GetVersionId().GetValue())

		if _, ok := deduper[key]; ok {
			continue
		}
		result = append(result, res)
		deduper[key] = struct{}{}
	}
	return result
}

// GroupResources organizes all resources by their underlying resource Type,
// and returns a map of the Type to the list of resources of that given type.
//
// Nil resources are skipped.
// Resources with existing IDs are skipped
func GroupResources(resources []fhir.Resource) map[Type][]fhir.Resource {
	result := map[Type][]fhir.Resource{}
	seen := map[string]struct{}{}
	for _, res := range resources {
		if res == nil {
			continue
		}
		if id := res.GetId(); id != nil {
			uri := URIString(res)
			if _, ok := seen[uri]; ok {
				continue // skip ones we have seen before
			}
			// add in new ones
			seen[uri] = struct{}{}
		}

		key := TypeOf(res)
		result[key] = append(result[key], res)
	}
	return result
}

// HasGetIdentifierList is a custom interface for duck typing resources that
// have a GetIdentifier method that returns a slice of Identifiers.
type HasGetIdentifierList interface {
	GetIdentifier() []*dtpb.Identifier

	// embed Resource since anything with an Identifier is also a Resource
	fhir.Resource
}

// HasGetIdentifierSingle is a custom interface for duck typing resources that
// have a GetIdentifier method that returns a single Identifier.
type HasGetIdentifierSingle interface {
	GetIdentifier() *dtpb.Identifier

	// embed Resource since anything with an Identifier is also a Resource
	fhir.Resource
}

// GetIdentifierList takes a Resource and returns a list of Identifiers.
// It uses duck typing to determine whether the resource has a GetIdentifier()
// method, and if so, whether it returns a list or a single Identifier.
// It returns ErrGetIdentifierList if the resource does not implement GetIdentifier().
// The list may be nil or empty if no identifiers are present.
// See interfaces: fhir.HasGetIdentifierList, fhir.HasGetIdentifierSingle
func GetIdentifierList(res fhir.Resource) ([]*dtpb.Identifier, error) {
	if cast, ok := res.(HasGetIdentifierList); ok {
		// resource implements GetIdentifier() as a list
		return cast.GetIdentifier(), nil
	}

	if cast, ok := res.(HasGetIdentifierSingle); ok {
		// resource implements GetIdentifier() as a single Identifier
		id := cast.GetIdentifier()
		if id == nil {
			return nil, nil
		}
		return []*dtpb.Identifier{id}, nil
	}

	// This is likely a bug / results from passing an unexpected type of resource
	return nil, fmt.Errorf("%w: Resource does not implement GetIdentifier(): %v", ErrGetIdentifierList, res)
}

// GetExtensions returns the list of extensions from a resource, if available.
func GetExtensions(res fhir.Resource) (extensions []*dtpb.Extension, ok bool) {
	type extensionList interface {
		GetExtension() []*dtpb.Extension
	}
	type extensionSingle interface {
		GetExtension() *dtpb.Extension
	}

	// Check if the resource contains an Extension, either as a list or a single item.
	if cast, ok := res.(extensionList); ok {
		return cast.GetExtension(), true
	}

	if cast, ok := res.(extensionSingle); ok {
		ext := cast.GetExtension()
		if ext == nil {
			return nil, true
		}

		return []*dtpb.Extension{ext}, true
	}

	return nil, false
}
