package resource

import (
	"fmt"
	"regexp"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"

	"github.com/verily-src/fhirpath-go/internal/fhir"
)

const (
	resourceTypePattern = "[A-Za-z]+"
	idPattern           = "[0-9A-Za-z.-]{1,64}" // https://www.hl7.org/fhir/R4/datatypes.html#id
	resourceIDPattern   = idPattern             // https://www.hl7.org/fhir/R4/resource.html#resource
	versionIDPattern    = idPattern             // https://www.hl7.org/fhir/R4/resource.html#Meta
)

var (
	historyURLRegexp = regexp.MustCompile(fmt.Sprintf("^.*/(%s)/(%s)/_history/(%s)$", resourceTypePattern, resourceIDPattern, versionIDPattern))
)

// Identity is a representation of a FHIR Resource's temporal instance.
//
// This is similar to a FHIR Reference, except without any explicit sementics of
// a referential relationship. Rather, this object simply acts as a carrier for
// the data that may be used for general identity purposes, such as logging.
type Identity struct {
	typeName Type
	id       string
	version  string
}

// Equal implements equality comparison between identity instances.
//
// The Equal() method, when used by the "cmp" package MUST
// support nils: see "...even if x or y is nil" from second bullet
// in https://pkg.go.dev/github.com/google/go-cmp/cmp#Equal.
func (i *Identity) Equal(other *Identity) bool {
	if i == other {
		return true
	}
	if i == nil || other == nil {
		return false
	}
	return i.typeName == other.typeName && i.id == other.id && i.version == other.version
}

// Type returns the resource Identity's underlying type. This is guaranteed to
// always be a valid resource type.
func (i *Identity) Type() Type {
	return i.typeName
}

// ID returns the resource's ID.
func (i *Identity) ID() string {
	return i.id
}

// VersionID returns the explicit version of the resource, if it is known.
func (i *Identity) VersionID() (string, bool) {
	return i.version, i.version != ""
}

// RelativeURI returns a relative URI of this resource.
func (i *Identity) RelativeURI() *dtpb.Uri {
	return fhir.URI(fmt.Sprintf("%v/%v", i.typeName, i.id))
}

// RelativeURIString returns a string representation of the RelativeURI for
// convenience.
func (i *Identity) RelativeURIString() string {
	return i.RelativeURI().GetValue()
}

// RelativeVersionedURI returns a relative URI of this resource including the
// version identifier, if it is known.
func (i *Identity) RelativeVersionedURI() (*dtpb.Uri, bool) {
	if i.version == "" {
		return nil, false
	}
	return fhir.URI(fmt.Sprintf("%v/%v/_history/%v", i.typeName, i.id, i.version)), true
}

// RelativeVersionedURIString returns a string representation of the RelativeVersionedURI
// for convenience.
func (i *Identity) RelativeVersionedURIString() (string, bool) {
	val, ok := i.RelativeVersionedURI()
	if ok {
		return val.GetValue(), true
	}
	return "", false
}

// PreferRelativeVersionURI returns the relative version URI if available,
// otherwise the relative URI only.
func (i *Identity) PreferRelativeVersionedURI() *dtpb.Uri {
	if uri, ok := i.RelativeVersionedURI(); ok {
		return uri
	}
	return i.RelativeURI()
}

// PreferRelativeVersionedURIString returns the relative version URI string if
// available, otherwise the relative URI only.
func (i *Identity) PreferRelativeVersionedURIString() string {
	if uri, ok := i.RelativeVersionedURIString(); ok {
		return uri
	}
	return i.RelativeURIString()
}

// String returns a string representation of this Identity.
//
// The exact representation should not be relied on for any practical purpose;
// the only thing that is guaranteed is that for the unique triple of data
// containing (type, id, version), the String will contain these details -- but
// the exact form is unspecified.
func (i *Identity) String() string {
	if i.version == "" {
		return fmt.Sprintf("%v/%v", i.typeName, i.id)
	}
	return fmt.Sprintf("%v/%v/_history/%v", i.typeName, i.id, i.version)
}

// Unversioned returns a new Identity that does not have a VersionID.
func (i *Identity) Unversioned() *Identity {
	return &Identity{
		typeName: i.typeName,
		id:       i.id,
	}
}

// WithNewVersion returns a new Identity that has the specified VersionID.
func (i *Identity) WithNewVersion(versionID string) *Identity {
	return &Identity{
		typeName: i.typeName,
		id:       i.id,
		version:  versionID,
	}
}

// NewIdentity attempts to create a new Identity object from a runtime-provided
// string resourceType name, and its id/versionID. If the provided resourceType
// does not name a valid resource-type (case-sensitive), this function will
// return an ErrBadType error.
func NewIdentity(resourceType, id, versionID string) (*Identity, error) {
	name, err := NewType(resourceType)
	if err != nil {
		return nil, err
	}

	return &Identity{
		typeName: name,
		id:       id,
		version:  versionID,
	}, nil
}

// NewIdentityFromHistoryURL attempts to create a new Identity object from a
// runtime-provided history URL.
//
// Input: [FHIR Proxy/Store URL]/fhir/[resourceType]/[resourceId]/_history/[versionId]
// Output: [resourceID]
func NewIdentityFromHistoryURL(url string) (*Identity, error) {
	matches := historyURLRegexp.FindStringSubmatch(url)
	if len(matches) != 4 {
		return nil, fmt.Errorf("error parsing history URL: %s", url)
	}
	return NewIdentity(matches[1], matches[2], matches[3])
}

// IdentityOf attempts to form a resource Identity object to the named
// resource. If the specified resource is either nil, or does not contain an
// ID value, no resource identity will be formed and this function will return
// nil.
func IdentityOf(resource fhir.Resource) (*Identity, bool) {
	if resource == nil || resource.GetId() == nil {
		return nil, false
	}

	return &Identity{
		typeName: TypeOf(resource),
		id:       resource.GetId().GetValue(),
		version:  resource.GetMeta().GetVersionId().GetValue(),
	}, true
}
