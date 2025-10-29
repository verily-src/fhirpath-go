package resource

import (
	"errors"
	"fmt"
	"strings"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/protofields"
)

var (
	// ErrMissingCanonicalURL is thrown when creating a canonical identity without having a URL.
	ErrMissingCanonicalURL = errors.New("missing canonical url")

	delimiter = "/"

	canonicalResourceTypes map[string]bool
)

// CanonicalIdentity is a canonical representation of a FHIR Resource.
//
// This object stores the individual pieces of id used in creating a canonical reference.
type CanonicalIdentity struct {
	Version  string
	Url      string
	Fragment string // only used if a fragment of resource is targetted
}

// canonicalTypeMatcher is used to identify if a resource is a canonical resource.
// This
type canonicalTypeMatcher interface {
	GetUrl() *dtpb.Uri
	GetVersion() *dtpb.String
}

// Type attempts to identify the resource type associated with the identity.
func (c *CanonicalIdentity) Type() (Type, bool) {
	for _, r := range strings.Split(c.Url, delimiter) {
		if IsType(r) {
			return Type(r), true
		}
	}
	return Type(""), false
}

// String returns a string representation of this CanonicalIdentity.
func (c *CanonicalIdentity) String() string {
	res := c.Url
	if c.Version != "" {
		res = fmt.Sprintf("%s|%s", res, c.Version)
	}
	if c.Fragment != "" {
		res = fmt.Sprintf("%s#%s", res, c.Fragment)
	}
	return res
}

// NewCanonicalIdentity creates a canonicalIdentity based on the given url, version and fragment
func NewCanonicalIdentity(url, version, fragment string) (*CanonicalIdentity, error) {
	if url == "" {
		return nil, ErrMissingCanonicalURL
	}

	return &CanonicalIdentity{
		Url:      url,
		Version:  version,
		Fragment: fragment,
	}, nil
}

// IsCanonicalType checks if a resource type is a canonical resource.
// https://hl7.org/fhir/R4/references.html#canonical does not list contract as a canonical resource.
// But contract has a url and version, making it look like a canonical resource.
// NamingSystem is on the list, but does not have a url and version, so it is not a canonical resource.
func IsCanonicalType(resourceType string) bool {
	return canonicalResourceTypes[resourceType]
}

func init() {
	canonicalResourceTypes = make(map[string]bool)

	for name, res := range protofields.Resources {
		d := res.DummyResource
		if _, ok := d.(canonicalTypeMatcher); ok {
			canonicalResourceTypes[name] = true
		}
	}
}
