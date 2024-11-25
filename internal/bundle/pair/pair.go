// Pair package holds pair of request and response bundles and offers functions
// to resolve between the request and response entries.
package pair

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	bcrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"

	"github.com/verily-src/fhirpath-go/internal/bundle"
	"github.com/verily-src/fhirpath-go/internal/element/etag"
	"github.com/verily-src/fhirpath-go/internal/element/reference"
	"github.com/verily-src/fhirpath-go/internal/resource"
)

var (
	// These are parse-time errors.
	ErrMismatchBundleLengths          = errors.New("request and response bundles have different entry lengths")
	ErrUnsupportedBundleType          = errors.New("unsupported bundle type")
	ErrDuplicateBundleEntryURI        = errors.New("duplicate effective URI in bundle")
	ErrMissingBundleEntryRequest      = errors.New("missing Bundle_Entry.Request")
	ErrMissingBundleEntryResponse     = errors.New("missing Bundle_Entry.Response")
	ErrRetrieveEtagVersion            = errors.New("unable to retrieve version ID from etag")
	ErrUnexpectedResponseStatusCode   = errors.New("unexpected response status code")
	ErrBadBundleEntryResponseIdentity = errors.New("bad identity of Bundle_Entry.Response")

	// These are lookup-time (searching) errors.
	ErrNotFoundBundleEntry         = errors.New("bundle entry not found on lookup")
	ErrNotFoundBundleEntryIdentity = errors.New("bundle entry identity not found on lookup")
)

// Pair holds an executed transaction-type request and response bundle,
// and offers functions to resolve between the request and response entries.
//
// Pair is general purpose and not at all specific to Provenance.
type Pair struct {
	req *bcrpb.Bundle
	rsp *bcrpb.Bundle

	// A map from each request's "effective" URI to its entry index
	// in the bundles.
	//
	// WATCHOUT: Each request entry may have multiple effective URIs as a result
	// there might be multiple entries in the map pointing to the same index
	// also some entries might not have an effective URI and hence not be indexed.
	entryIdxByReqURI map[string]int

	// The identity of each entry, relative to serviceBaseURL below.
	// The identity is based upon the response entry. The request entry
	// often (eg, POST operations) don't have an identity.
	rspIdentities []*resource.Identity

	// The serviceBaseURL of the store that executed this bundle.
	serviceBaseURL string
}

// NewPair constructs a Pair.
//
// It performs some initial validity/consistency checks and builds
// some internal indexes for later use.
//
// This could be extended to support batch bundles, not just transaction bundles.
// The primary client (Provenance resource resolution) always uses transactions
// (not batches). Thus adding batch support is not a priority right now.
func NewPair(request *bcrpb.Bundle, response *bcrpb.Bundle) (*Pair, error) {
	numEntries := len(request.GetEntry())
	if numEntries != len(response.GetEntry()) {
		return nil, ErrMismatchBundleLengths
	}
	if request.GetType().GetValue() != cpb.BundleTypeCode_TRANSACTION || response.GetType().GetValue() != cpb.BundleTypeCode_TRANSACTION_RESPONSE {
		return nil, fmt.Errorf("%w: request type %v, response type %v", ErrUnsupportedBundleType,
			request.GetType().GetValue(), response.GetType().GetValue())
	}

	bp := &Pair{
		req:              request,
		rsp:              response,
		entryIdxByReqURI: make(map[string]int),
		rspIdentities:    make([]*resource.Identity, numEntries),
	}

	reqEntryMetas := bundle.NewEntryMetasFromEntries(request.GetEntry())
	for entryIdx := 0; entryIdx < numEntries; entryIdx++ {
		reqEntry := request.GetEntry()[entryIdx]
		if reqEntry.GetRequest() == nil {
			return nil, fmt.Errorf("%w: entry %d", ErrMissingBundleEntryRequest, entryIdx)
		}

		reqEntryMeta := reqEntryMetas[entryIdx]
		reqMethod := reqEntry.GetRequest().GetMethod().GetValue()
		// EffectiveURIs is empty for a create operation that is not
		// referenced by any other resources in the bundle.
		for _, reqURI := range reqEntryMeta.EffectiveURIs() {
			// WATCHOUT: Duplicate URIs for GET requests are permitted.
			// This is since when canonicalize operation is called
			// multiple provenance entries will be merged into a single entry
			// and some of them that are merged will be replaced by the same
			// dummy GET request and hence might have the same URI.
			if priorEntryIdx, exists := bp.entryIdxByReqURI[reqURI]; exists && reqMethod != cpb.HTTPVerbCode_GET {
				return nil, fmt.Errorf("%w: uri='%s' at index %d and %d", ErrDuplicateBundleEntryURI,
					reqURI, priorEntryIdx, entryIdx)
			}
			bp.entryIdxByReqURI[reqURI] = entryIdx
		}

		rspEntry := response.GetEntry()[entryIdx]
		if rspEntry.GetResponse() == nil {
			return nil, fmt.Errorf("%w: entry %d", ErrMissingBundleEntryResponse, entryIdx)
		}
		// DELETE responses do not have a Location.
		// So we use the request URL and response Etag to find the identity.

		if reqMethod == cpb.HTTPVerbCode_DELETE {
			// Successful DELETE call returns a `200 OK`, `204 No Content` or `202 Accepted` status code.
			// See https://hl7.org/implement/standards/fhir/R4/http.html#delete
			successfulDeleteStatusCodes := []int{200, 204, 202}
			entryStatusCode, err := bundle.StatusCodeFromEntry(rspEntry)
			if err != nil {
				return nil, fmt.Errorf("failed to extract status code from bundle entry %d: %w", entryIdx, err)
			}
			if !slices.Contains(successfulDeleteStatusCodes, entryStatusCode) {
				return nil, fmt.Errorf("%w: entry %d: received status code %d",
					ErrUnexpectedResponseStatusCode, entryIdx, entryStatusCode)
			}
			// Create identity and set the version from the response etag.
			etagValue := rspEntry.GetResponse().GetEtag().GetValue()
			versionID, err := etag.VersionIDFromEtag(etagValue)
			if err != nil {
				return nil, fmt.Errorf("%w: entry %d (%v, etag=[%s]): %w",
					ErrRetrieveEtagVersion, entryIdx, reqMethod, etagValue, err)
			}
			identity, err := reference.IdentityFromURL(reqEntry.GetRequest().GetUrl().GetValue())
			if err != nil {
				return nil, fmt.Errorf("%w: entry %d: %w", ErrBadBundleEntryResponseIdentity, entryIdx, err)
			}
			versionedIdentity := identity.WithNewVersion(versionID)
			bp.rspIdentities[entryIdx] = versionedIdentity
			// Can not find serviceBaseURL from DELETE response.
		} else {
			locLit, err := reference.LiteralInfoFromURI(rspEntry.GetResponse().GetLocation().GetValue())
			if err != nil {
				return nil, fmt.Errorf("%w: entry %d: %v", ErrBadBundleEntryResponseIdentity, entryIdx, err)
			}
			identity, ok := locLit.Identity()
			if !ok {
				return nil, fmt.Errorf("%w: entry %d: %v", ErrBadBundleEntryResponseIdentity, entryIdx, "missing identity")
			}
			if _, haveVersion := identity.VersionID(); !haveVersion {
				return nil, fmt.Errorf("%w: entry %d: %v", ErrBadBundleEntryResponseIdentity, entryIdx, "missing version")
			}
			if locLit.ServiceBaseURL() == "" {
				return nil, fmt.Errorf("%w: entry %d: %v", ErrBadBundleEntryResponseIdentity, entryIdx, "missing ServiceBaseURL")
			}
			if bp.serviceBaseURL == "" {
				bp.serviceBaseURL = locLit.ServiceBaseURL()
			} else if bp.serviceBaseURL != locLit.ServiceBaseURL() {
				return nil, fmt.Errorf("%w: entry %d: inconsistent ServiceBaseURL '%s' vs '%s'",
					ErrBadBundleEntryResponseIdentity, entryIdx, bp.serviceBaseURL, locLit.ServiceBaseURL())
			}
			bp.rspIdentities[entryIdx] = identity
		}
	}
	return bp, nil
}

// ReqBundle returns the request bundle of the pair.
func (bp *Pair) ReqBundle() *bcrpb.Bundle {
	return bp.req
}

// ReqBundle returns the response bundle of the pair.
func (bp *Pair) RspBundle() *bcrpb.Bundle {
	return bp.rsp
}

// ReqBundle returns the specified entry of the request bundle.
func (bp *Pair) ReqEntryOfIdx(entryIdx int) *bcrpb.Bundle_Entry {
	return bp.req.GetEntry()[entryIdx]
}

// RspBundle returns the specified entry of the response bundle.
func (bp *Pair) RspEntryOfIdx(entryIdx int) *bcrpb.Bundle_Entry {
	return bp.rsp.GetEntry()[entryIdx]
}

// ServiceBaseURL returns the service base URL of the resources in the response
// bundle, based upon their Location attribute.
func (bp *Pair) ServiceBaseURL() string {
	return bp.serviceBaseURL
}

// IdentityOfIdx returns the resource.Identity of the given bundle entry.
//
// The returned resource.Identity is based upon the response bundle.
// It is always versioned. It will be nil (and the "ok" bool false) for
// DELETE operations.
func (bp *Pair) IdentityOfIdx(entryIdx int) (*resource.Identity, bool) {
	return bp.rspIdentities[entryIdx], bp.rspIdentities[entryIdx] != nil
}

// IdentityOfRef returns the versioned identity of the referenced resource.
//
// The reference should be an element within a resource of the request bundle.
// As such it may contain normal FHIR REST URI (like Patient/1234) or
// may be an URN or canonical that only appears in transaction request bundles.
//
// Regardless of the type of reference, the returned identity holds
// the FHIR REST Identity of the reference resource. This identity is obtained
// from the response bundle and includes the resource's version.
//
// Note that the returned identity does not include the serviceBaseURL
// of the containing FHIR store.
func (bp *Pair) IdentityOfRef(ref *dtpb.Reference) (*resource.Identity, error) {
	refLit, err := reference.LiteralInfoOf(ref)
	if err != nil {
		return nil, err
	}
	// Normalize to favor the relative URL.
	if refLit.ServiceBaseURL() == bp.serviceBaseURL {
		newRefLit, err := refLit.WithServiceBaseURL("")
		if err != nil {
			return nil, err
		}
		refLit = newRefLit
	}

	// The reference's URI, expressed as a normalized string, should
	// exactly match the effective URI of the referenced request bundle entry.
	refURI := refLit.URIString()
	entryIdx, ok := bp.entryIdxByReqURI[refURI]
	if !ok {
		return nil, fmt.Errorf("%w: by reference URI [%v]", ErrNotFoundBundleEntry, refURI)
	}
	identity := bp.rspIdentities[entryIdx]
	if identity == nil {
		// Deleted resources in the bundle do not have a Location and thus
		// do not have an Identity. That itself is not an error. But referencing
		// a deleted resource is an error (for now).
		// This isn't reached by any real code because typical DELETE operations don't have
		// an effective URI and aren't indexed into entryIdxByReqURI.
		return nil, fmt.Errorf("%w: entry %d for reference URI [%v]", ErrNotFoundBundleEntryIdentity, entryIdx, refURI)
	}
	return identity, nil
}

// MustIdentityOfIdx returns the resource.Identity of the given bundle entry
// index.
//
// The returned resource.Identity is based upon the response bundle.
// It is always versioned. It will be panic on DELETE operations.
func (bp *Pair) MustIdentityOfIdx(entryIdx int) *resource.Identity {
	identity, ok := bp.IdentityOfIdx(entryIdx)
	if !ok {
		panic(fmt.Errorf("resource at index %v does not have identity", entryIdx))
	}
	return identity
}

// MustTypedRefOfIdx returns the unversioned typed reference at the given bundle
// entry index.
func (bp *Pair) MustTypedRefOfIdx(entryIdx int) *dtpb.Reference {
	return reference.TypedFromIdentity(bp.MustIdentityOfIdx(entryIdx).Unversioned())
}

// MustVersionedTypedRefOfIdx returns the versioned typed reference at the given bundle
// entry index.
func (bp *Pair) MustVersionedTypedRefOfIdx(entryIdx int) *dtpb.Reference {
	return reference.TypedFromIdentity(bp.MustIdentityOfIdx(entryIdx))
}

// MustWeakRefOfIdx returns the unversioned weak reference at the given bundle
// entry index.
func (bp *Pair) MustWeakRefOfIdx(entryIdx int) *dtpb.Reference {
	identity := bp.MustIdentityOfIdx(entryIdx)
	return reference.Weak(identity.Type(), identity.RelativeURIString())
}

// Dump returns a complete dump of the pair's request and response.
// As the bundle very likely contains PHI, which must not be logged,
// this should only be used for testing.
func (bp *Pair) Dump() string {
	var sb strings.Builder

	for eIdx := range bp.req.Entry {
		fmt.Fprintf(&sb, "Bundle Type: REQ=%v RSP=%v\n", bp.req.GetType(), bp.rsp.GetType())
		// Ideally the req and rsp entries should be indented, but we don't seem to have
		// an indent package imported into verily1 yet.
		fmt.Fprintf(&sb, "Bundle Entry #%d:\nREQ=%v\nRSP=%v\n",
			eIdx, bp.ReqEntryOfIdx(eIdx), bp.RspEntryOfIdx(eIdx))
	}
	return sb.String()
}
