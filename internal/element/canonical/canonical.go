package canonical

import (
	"errors"
	"fmt"
	"regexp"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/resource"
)

var (
	// ErrNoCanonicalURL is an error returned when an API receives a canonical that
	// does not have a specified URL field.
	ErrNoCanonicalURL = errors.New("canonical does not contain url")

	// very basic regex for URL matching.
	// Fragment portion is from https://build.fhir.org/references.html#literal
	canonicalRegExp = regexp.MustCompile(`^(?P<url>[^|#]+)(\|(?P<version>[A-z0-9-_\.]+))?(#(?P<fragment>[A-z0-9-_\.]{1,64}))?`)
)

// canonicalConfig is an internal struct for holding canonical information that
// can be updated from a CanonicalOpt.
type canonicalConfig struct {
	fragment, version string
}

// Option is an option interface for constructing canonicals from raw
// data.
type Option interface {
	updateCanonical(data *canonicalConfig)
}

// WithFragment adds a "fragment" portion to Canonical references.
func WithFragment(frag string) Option {
	return canonicalFragOpt(frag)
}

// canonicalFragOpt is a simple canonical option for fragment strings.
type canonicalFragOpt string

func (o canonicalFragOpt) updateCanonical(data *canonicalConfig) {
	data.fragment = string(o)
}

// WithVersion adds a "version" portion to Canonical references.
func WithVersion(version string) Option {
	return canonicalVersionOpt(version)
}

// canonicalVersionOpt is a simple canonical option for version strings.
type canonicalVersionOpt string

func (o canonicalVersionOpt) updateCanonical(data *canonicalConfig) {
	data.version = string(o)
}

// New constructs an R4 FHIR New element from the specified
// url string and canonical options.
//
// See: http://hl7.org/fhir/R4/datatypes.html#canonical
func New(url string, opts ...Option) *dtpb.Canonical {
	data := &canonicalConfig{}
	for _, opt := range opts {
		opt.updateCanonical(data)
	}
	if data.version != "" {
		url = fmt.Sprintf("%v|%v", url, data.version)
	}
	if data.fragment != "" {
		url = fmt.Sprintf("%v#%v", url, data.fragment)
	}
	return &dtpb.Canonical{
		Value: url,
	}
}

// FromResource creates an R4 FHIR FromResource element from a
// resource that has a URL, such as a Questionnaire, Device, etc.
//
// If the input resource is nil, or if the resource does not have a URL
// field assigned, this function will return the error `ErrNoCanonicalURL`.
//
// See: https://hl7.org/fhir/R4/datatypes.html#canonical and
// https://hl7.org/fhir/R4/references.html#canonical
func FromResource(resource fhir.CanonicalResource) (*dtpb.Canonical, error) {
	if resource == nil || resource.GetUrl() == nil {
		return nil, ErrNoCanonicalURL
	}
	return New(resource.GetUrl().GetValue()), nil
}

// FragmentFromResource creates an R4 FHIR Canonical element from a resource that
// has a URL, such as a Questionnaire, Device, etc., and will mark it as a
// fragment-reference.
//
// If the input resource is nil, or if the resource does not have a URL
// field assigned, this function will return the error `ErrNoCanonicalURL`.
//
// See: https://hl7.org/fhir/R4/datatypes.html#canonical and
// https://hl7.org/fhir/R4/references.html#canonical
func FragmentFromResource(resource fhir.CanonicalResource) (*dtpb.Canonical, error) {
	if resource == nil || resource.GetUrl() == nil {
		return nil, ErrNoCanonicalURL
	}
	return New(resource.GetUrl().GetValue(), WithFragment(resource.GetId().GetValue())), nil
}

// VersionedFromResource creates an R4 FHIR Canonical element from a resource that
// has a URL, such as a Questionnaire, Device, etc, along with a version string.
//
// If the input resource is nil, or if the resource does not have a URL
// field assigned, this function will return the error `ErrNoCanonicalURL`.
//
// See: https://hl7.org/fhir/R4/datatypes.html#canonical and
// https://hl7.org/fhir/R4/references.html#canonical
func VersionedFromResource(resource fhir.CanonicalResource) (*dtpb.Canonical, error) {
	if resource == nil || resource.GetUrl() == nil {
		return nil, ErrNoCanonicalURL
	}
	url := resource.GetUrl()
	version := resource.GetVersion()
	if version == nil {
		return New(url.GetValue()), nil
	}
	return New(resource.GetUrl().GetValue(), WithVersion(version.GetValue())), nil
}

// IdentityFromReference returns an Identity object from a given canonical reference
// Replaces: ph.ParseCanonical
func IdentityFromReference(c *dtpb.Canonical) (*resource.CanonicalIdentity, error) {
	value := c.GetValue()
	match := canonicalRegExp.FindStringSubmatch(value)
	result := make(map[string]string)
	for i, name := range canonicalRegExp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return resource.NewCanonicalIdentity(result["url"], result["version"], result["fragment"])
}

// IdentityOf returns a canonicalIdentity representing the given canonical resource
func IdentityOf(res fhir.CanonicalResource) (*resource.CanonicalIdentity, error) {
	if res == nil || res.GetUrl() == nil {
		return nil, ErrNoCanonicalURL
	}

	return resource.NewCanonicalIdentity(res.GetUrl().GetValue(), res.GetVersion().GetValue(), "")
}
