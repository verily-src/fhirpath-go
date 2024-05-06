package resource

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrMissingCanonicalURL is thrown when creating a canonical identity without having a URL.
	ErrMissingCanonicalURL = errors.New("missing canonical url")

	delimiter = "/"
)

// CanonicalIdentity is a canonical representation of a FHIR Resource.
//
// This object stores the individual pieces of id used in creating a canonical reference.
type CanonicalIdentity struct {
	Version  string
	Url      string
	Fragment string // only used if a fragment of resource is targetted
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
