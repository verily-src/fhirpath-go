package reference

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/resource"
)

var (
	ErrNotLiteral              error = errors.New("Reference is not a literal")
	ErrTypeInvalid             error = errors.New("Reference.type is invalid")
	ErrTypeInconsistent        error = errors.New("Reference.reference and Reference.type are inconsistent")
	ErrWeakInvalid             error = errors.New("Reference.reference is invalid weak URI")
	ErrStrongInvalid           error = errors.New("Reference.reference is invalid strong URI")
	ErrExplicitFragmentInvalid error = errors.New("Reference.reference is invalid explicit fragment")
	ErrServiceBaseURLInvalid   error = errors.New("ServiceBaseURL is invalid")
)

// A LiteralInfo is the result of parsing a FHIR Reference when that reference
// contains a literal reference. It is immutable.
//
// LiteralInfo supports all valid forms a literal reference:
//   - REST URI: absolute and relative, versioned and unversioned,
//     strong and weak.
//   - FragmentID: a reference to a contained resource within an implied
//     containing resource.
//   - non-REST URI: URNs (eg., UUID and OID) and versionless canonicals.
//     These occur most commonly in request bundles (vs response bundles).
//
// A key goal of LiteralInfo is to canonicalize between the "weak" and "strong"
// representations of a literal reference. Note that this distinction is
// specific to the GCP proto mapping of FHIR and does not appear in the standard.
//
// The fields fragmentID, identity, and nonRESTURI are mutually exclusive: exactly
// one will be non-nil.
//
// See FHIR spec: http://hl7.org/fhir/R4/references.html#literal.
type LiteralInfo struct {
	// If non-nil, the type of the referenced resource.
	// If non-nil and identity also present, both guaranteed to have the same type.
	resType *resource.Type

	// If non-nil, the fragment ID that is relative to an implied containing resource.
	// A fragment ID may be an empty string, which indicates the implied
	// containing resource itself (as opposed to a contained resource within
	// the containing resource).
	fragmentID *string

	// The identity (type, ID and optional version) of the referenced resource.
	// Applies to REST URLs
	identity *resource.Identity

	// The service base URL of the FHIR store that holds the referenced resource.
	// Only non-empty if identity is non-nil and the reference is absolute.
	// Generally references to resources in the same store (as the reference itself)
	// are relative, and the service base URL will be empty.
	// See http://hl7.org/fhir/R4/http.html#general
	serviceBaseURL string

	// If non-nil, the non-REST URI of the referenced resource.
	// Typically an URN (UUID or OID). See NonRESTURI() accessor for details.
	nonRESTURI *string
}

// clone returns a shallow copy.
//
// All fields have immutable types, so a shallow copy is safe.
// LiteralInfo itself is immutable; thus, clone() is only used internally
// by the With*() mutators.
func (lit *LiteralInfo) clone() *LiteralInfo {
	return &LiteralInfo{
		resType:        lit.resType,
		fragmentID:     lit.fragmentID,
		identity:       lit.identity,
		serviceBaseURL: lit.serviceBaseURL,
		nonRESTURI:     lit.nonRESTURI,
	}
}

// Type returns the resource type of the reference resource, if known.
func (lit *LiteralInfo) Type() (resource.Type, bool) {
	if lit.resType == nil {
		var emptyType resource.Type
		return emptyType, false
	}
	return *lit.resType, true
}

// FragmentID returns the fragment ID of the reference fragment,
// if the literal reference is to fragment.
func (lit *LiteralInfo) FragmentID() (string, bool) {
	if lit.fragmentID == nil {
		return "", false
	}
	return *lit.fragmentID, true
}

// Identity returns the identity of the reference resource,
// if the referenced resource has REST identity.
func (lit *LiteralInfo) Identity() (*resource.Identity, bool) {
	return lit.identity, lit.identity != nil
}

// Identity returns the identity of the reference resource.
// Only non-empty if the reference resource has a REST identity,
// and that REST identity was absolute (not relative).
func (lit *LiteralInfo) ServiceBaseURL() string {
	return lit.serviceBaseURL
}

// WithServiceBaseURL returns a LiteralInfo with the given serviceBaseURL.
// The given serviceBaseURL must be empty or a valid serviceBaseURL.
func (lit *LiteralInfo) WithServiceBaseURL(serviceBaseURL string) (*LiteralInfo, error) {
	// WATCHOUT: Our regex wants a trailing slash, but our convention is to omit it.
	// We add the slash here rather than change the regex to keep traceability to the spec.
	if serviceBaseURL != "" && !restFHIRServiceBaseURLRegex.MatchString(serviceBaseURL+"/") {
		return nil, fmt.Errorf("%w: '%s'", ErrServiceBaseURLInvalid, serviceBaseURL)
	}
	if serviceBaseURL == lit.serviceBaseURL {
		return lit, nil
	}
	newLit := lit.clone()
	newLit.serviceBaseURL = serviceBaseURL
	return newLit, nil
}

// NonRESTURI returns the non-REST URI of the referenced resource,
// if the referenced resource has a non-REST URI. A non-REST URI
// is typically an URN (UUID or OID) or a versionless canonical URL.
// Th URN form is typically only used in bundle transaction requests
// when the REST URI is not yet known.
//
// See http://hl7.org/fhir/R4/references.html#literal, and specifically:
//   - "... in a bundle during a transaction, reference URLs may actually
//     contain logical URIs (e.g. OIDs or UUIDs) that resolve within the
//     transaction."
//   - "The URL may contain a reference to a canonical URL and applications can
//     use the canonical URL resolution methods they support ..."
func (lit *LiteralInfo) NonRESTURI() (string, bool) {
	if lit.nonRESTURI == nil {
		return "", false
	}
	return *lit.nonRESTURI, true
}

// URI returns the LiteralInfo's URI equivalent as a *dtpb.Uri.
func (lit *LiteralInfo) URI() *dtpb.Uri {
	s := lit.URIString()
	if s == "" {
		return nil
	}
	return fhir.URI(s)
}

// URIString returns the LiteralInfo's URI equivalent as a string.
// The returned URI is suitable for use in Reference.reference.
// This function is the inverse of LiteralInfoFromURI().
func (lit *LiteralInfo) URIString() string {
	if lit == nil {
		return ""
	}
	if lit.fragmentID != nil {
		return "#" + *lit.fragmentID
	}
	if lit.identity != nil {
		uri := lit.identity.PreferRelativeVersionedURIString()
		// Could use net/url.JoinPath but it does validation on serviceBaseURL
		// that we don't want.
		if lit.serviceBaseURL != "" {
			uri = lit.serviceBaseURL + "/" + uri
		}
		return uri
	}
	if lit.nonRESTURI != nil {
		return *lit.nonRESTURI
	}
	// UNREACHED
	return ""
}

// PreferRelativeVersionedURIString returns the relative URI with version if
// available, otherwise just relative URI without version.
func (lit *LiteralInfo) PreferRelativeVersionedURIString() string {
	identity, ok := lit.Identity()
	if !ok {
		return ""
	}
	return identity.PreferRelativeVersionedURIString()
}

// LiteralInfoOf parses the given literal reference.
//
// Returns an error if the given reference is not a valid literal reference.
func LiteralInfoOf(ref *dtpb.Reference) (*LiteralInfo, error) {
	var explicitType *resource.Type
	if resTypeStr := ref.GetType(); resTypeStr != nil {
		t, err := resource.NewType(resTypeStr.GetValue())
		if err != nil {
			// Returned err includes the bad type, so no need to duplicate here.
			return nil, fmt.Errorf("%w: %v", ErrTypeInvalid, err)
		}
		explicitType = &t
	}

	// Fragment
	if frag := ref.GetFragment(); frag != nil {
		fragStr := frag.GetValue()
		// WATCHOUT: an empty fragStr is a valid fragID but not a valid ID.
		if fragStr != "" && !fhir.IsID(fragStr) {
			return nil, fmt.Errorf("%w: invalid fragment ID", ErrExplicitFragmentInvalid)
		}
		return &LiteralInfo{
			resType:    explicitType,
			fragmentID: &fragStr,
		}, nil
	}

	// An absolute or relative weak URI, a weak fragment URI, or a non-REST URI.
	if uri := ref.GetUri(); uri != nil {
		litUri, err := LiteralInfoFromURI(uri.GetValue())
		if err != nil {
			return nil, fmt.Errorf("%w: uri='%v': %v", ErrWeakInvalid, uri.GetValue(), err)
		}
		if litUri.resType == nil {
			// This will occur for a fragment or URN.
			litUri.resType = explicitType
		} else {
			if explicitType != nil && *litUri.resType != *explicitType {
				return nil, fmt.Errorf("%w: weak='%v' vs type='%v'",
					ErrTypeInconsistent, *litUri.resType, explicitType)
			}
		}
		return litUri, nil
	}

	// Strong relative URIs
	identity, err := identityOfStrong(ref)
	if err != nil {
		if errors.Is(err, ErrReferenceOneOfResourceNotSet) {
			return nil, ErrNotLiteral
		}
		return nil, fmt.Errorf("%w: %v", ErrStrongInvalid, err)
	}
	strongType := identity.Type()
	if explicitType != nil && *explicitType != strongType {
		return nil, fmt.Errorf("%w: strong='%v' vs type='%v'",
			ErrTypeInconsistent, strongType, explicitType)
	}
	return &LiteralInfo{
		resType:  &strongType,
		identity: identity,
	}, nil
}

// LiteralInfoFromURI parses a Reference.reference URI string,
// where Reference is the FHIR element. Within the GCP proto mapping
// (where Reference.reference is a oneof), this function parses
// the Reference.uri member of that oneof.
//
// Parsing a Reference URI is easier than parsing general resource URL:
//   - We don't need to worry about canonical URLs. (AFIAK). Canonical URLs
//     are there own datatype, not a Reference element.
//   - Fragments within Reference URI are always relative to the containing
//     resource. That is, Reference.reference cannot specify a containing
//     resource (relative or absolute) reference that is suffixed by a fragment.
//     (Canonical references do allow this).
//
// This supports the special case of an fragment without an ID. This is used
// by a contained resource to reference its containing resource.
//
// Any returned error does NOT include the uri. That is left to the caller.
func LiteralInfoFromURI(uri string) (*LiteralInfo, error) {
	if uri == "" {
		return nil, fmt.Errorf("%w: empty uri", ErrInvalidURI)
	}
	if uri[0] == '#' {
		// Note that "#" alone is a valid fragment.
		fragStr := uri[1:]
		if fragStr != "" && !fhir.IsID(fragStr) {
			return nil, fmt.Errorf("%w: invalid fragment id", ErrInvalidURI)
		}
		return &LiteralInfo{fragmentID: &fragStr}, nil
	}
	if strings.Contains(uri, "#") {
		return nil, fmt.Errorf("%w: non-relative fragment found", ErrInvalidURI)
	}
	if strings.Contains(uri, "|") {
		return nil, fmt.Errorf("%w: canonical version found", ErrInvalidURI)
	}

	// This regexp matches both relative and absolute REST URIs.
	uriIndexes := restFHIRServiceResourceURLRegex.FindStringSubmatchIndex(uri)

	if uriIndexes == nil {
		// Try as a non-REST URI (typically an URN or a versionless canonical).
		parsedUrl, err := url.Parse(uri)
		if err != nil {
			return nil, fmt.Errorf("%w: unparsable: %v", ErrInvalidURI, err)
		}
		if parsedUrl.Scheme == "" {
			return nil, fmt.Errorf("%w: non-REST and missing scheme component", ErrInvalidURI)
		}
		if parsedUrl.Opaque != "" || strings.TrimLeft(parsedUrl.Path, "/") != "" {
			// An URN is parsed as Opaque and a canonical is parsed as a Path.
			return &LiteralInfo{
				nonRESTURI: &uri,
			}, nil
		}
		return nil, fmt.Errorf("%w: non-REST and missing both Opaque and Path", ErrInvalidURI)
	}

	// The relative portion of the URI is the 4th submatch in the regexp.
	relStartIdx := uriIndexes[8]
	baseUrl := strings.TrimRight(uri[0:relStartIdx], "/")
	relUrl := uri[relStartIdx:]

	// The REST regexp could be used to identify all the parts of the relative URI,
	// but easier just to split.
	relParts := strings.Split(relUrl, "/")
	var identity *resource.Identity
	if !resource.IsType(relParts[0]) {
		return nil, fmt.Errorf("%w: resource type is invalid", ErrInvalidURI)
	}
	if !fhir.IsID(relParts[1]) {
		return nil, fmt.Errorf("%w: resource id is invalid", ErrInvalidURI)
	}
	if len(relParts) == 2 {
		identity, _ = resource.NewIdentity(relParts[0], relParts[1], "")
	} else if len(relParts) == 4 && relParts[2] == "_history" {
		if !fhir.IsID(relParts[3]) {
			return nil, fmt.Errorf("%w: version id is invalid", ErrInvalidURI)
		}
		identity, _ = resource.NewIdentity(relParts[0], relParts[1], relParts[3])
	} else {
		return nil, fmt.Errorf("%w: invalid relative component", ErrInvalidURI)
	}

	identityType := identity.Type()
	return &LiteralInfo{
		resType:        &identityType,
		identity:       identity,
		serviceBaseURL: baseUrl,
	}, nil
}
