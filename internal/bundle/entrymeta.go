package bundle

import (
	"errors"
	"strings"

	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	bcrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"github.com/hashicorp/go-set/v2"

	"github.com/verily-src/fhirpath-go/internal/containedresource"
	"github.com/verily-src/fhirpath-go/internal/element/reference"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"google.golang.org/protobuf/proto"
)

var (
	ErrLiteralInfoTypeFailure           = errors.New("LiteralInfo.Type() failed")
	ErrResourceTypeCouldNotBeDetermined = errors.New("resource type could not be determined")
)

// bundleEntryURIs holds all possible URIs associated with a FHIR bundle entry,
// which can be used for referencing the entry.
type bundleEntryURIs struct {
	// fullURI represents the full URL of the entry, if available.
	// This is typically a URN (example: urn:uuid:12345).
	fullURI string

	// resourceURI represents the un-versioned URL derived from the resource
	// within the entry, if the resource is present and has an ID. The typical
	// format is "<ResourceType>/<ID>", such as "Patient/12345".
	resourceURI string

	// requestURI represents the URL derived from the request within the entry.
	// This is typically a verioned/un-versioned Relative URL.
	requestURI string

	// responseLocationURI represents the URL derived from the response location
	// within the entry, if available. This is typically an absolute and
	// versioned URL.
	responseLocationURI string
}

// EntryMeta holds metadata about a FHIR bundle entry.
//
// WATCHOUT: This is different from dtpb.Meta, which is a metadata struct for
// the entire bundle. EntryMeta is a metadata struct that contains information
// for a single bundle entry.
type EntryMeta struct {
	// The bundle entry itself.
	entry *bcrpb.Bundle_Entry

	// The 0-based index of this entry within its containing bundle.
	entryIdx int

	// The entry's resource, or nil.
	res fhir.Resource

	// The entry's URL set.
	bundleEntryURIs bundleEntryURIs
}

// NewEntryMetasFromEntries constructs an slice of 1:1 metas from a slice of entries.
func NewEntryMetasFromEntries(bundleEntries []*bcrpb.Bundle_Entry) []*EntryMeta {
	metas := make([]*EntryMeta, len(bundleEntries))
	for entryIdx, entry := range bundleEntries {
		metas[entryIdx] = NewEntryMeta(entryIdx, entry)
	}
	return metas
}

func NewEntryMeta(entryIdx int, entry *bcrpb.Bundle_Entry) *EntryMeta {
	bundleEntryURIs := bundleEntryURIs{}
	if entry.GetFullUrl() != nil {
		bundleEntryURIs.fullURI = entry.GetFullUrl().GetValue()
	}

	if res := UnwrapEntry(entry); res != nil && res.GetId() != nil {
		bundleEntryURIs.resourceURI = resource.URIString(res)
	}

	if entry.GetRequest().GetUrl() != nil {
		bundleEntryURIs.requestURI = entry.GetRequest().GetUrl().GetValue()
	}

	if entry.GetResponse().GetLocation() != nil {
		bundleEntryURIs.responseLocationURI = entry.GetResponse().GetLocation().GetValue()
	}

	return &EntryMeta{
		entry:           entry,
		entryIdx:        entryIdx,
		res:             UnwrapEntry(entry),
		bundleEntryURIs: bundleEntryURIs,
	}
}

// Equal returns true if the two metas have the same value.
func (m *EntryMeta) Equal(other *EntryMeta) bool {
	if m == other {
		return true
	}
	return m.entryIdx == other.entryIdx && proto.Equal(m.entry, other.entry)
}

// Index returns the entry itself.
func (m *EntryMeta) Entry() *bcrpb.Bundle_Entry {
	return m.entry
}

// Index returns the index of the entry within its containing bundle.
func (m *EntryMeta) Index() int {
	return m.entryIdx
}

// Resource returns the FHIR resource of the entry, or nil.
func (m *EntryMeta) Resource() fhir.Resource {
	return m.res
}

// EffectiveURIs returns all the possible effective URLs for a FHIR bundle
// entry. Effective URIs can be
// 1. URN (fullURI)
// 2. Relative un-versioned URL (resourceURI, requestURI)
func (m *EntryMeta) EffectiveURIs() []string {
	// We use a set to avoid duplicate effective URIs.
	// An example case when a duplicate effective URL can occur is when a PUT
	// request is trying to update a resource and the resource URL can be the
	// same as the request URL.
	effectiveURISet := set.New[string](0)

	if m.bundleEntryURIs.fullURI != "" {
		effectiveURISet.Insert(m.bundleEntryURIs.fullURI)
	}

	// WATCHOUT: For PATCH requests generated from binary payloads the resource
	// URL is set to Binary/ID instead of the expected <resourceType>/ID. So we
	// skip adding the resource URL for PATCH requests.
	if m.bundleEntryURIs.resourceURI != "" && m.Method() != cpb.HTTPVerbCode_PATCH {
		effectiveURISet.Insert(m.bundleEntryURIs.resourceURI)
	}

	// WATCHOUT: We support this case for PATCH requests and PATCH requests in a
	// bundle are generated from binary payloads with the request URI set as a
	// relative URI. Therefore, we add the relative URI guard check here.
	if m.bundleEntryURIs.requestURI != "" && isUnversionedRelativeURL(m.bundleEntryURIs.requestURI) {
		effectiveURISet.Insert(m.bundleEntryURIs.requestURI)
	}

	return effectiveURISet.Slice()
}

// extractResourceTypeFromURL extracts the resource type from the given URL. The
// function returns nil resource type and corresponding error if the extraction
// fails.
func (m *EntryMeta) extractResourceTypeFromURL(url string) (resource.Type, error) {
	var nilResourceType resource.Type
	lit, err := reference.LiteralInfoFromURI(url)
	if err != nil {
		return nilResourceType, err
	}

	resType, ok := lit.Type()
	if !ok {
		return nilResourceType, ErrLiteralInfoTypeFailure
	}
	return resType, nil
}

// ResourceTypeBeta identifies the resource type being mutated or queried in
// the current bundle entry. This method successfully resolves the resource type
// for the following cases:
// 1. The entry is not a PATCH request and has a resource.
// 2. The entry has a request URL that targets a resource using a relative URL
// ex: Patient/12345.
// 3. The entry has full URL set in entry.response.location.
//
// WATCHOUT:
// 1. The method is labeled "Beta" due to potential unknown cases in the current
// implementation. Once these are resolved, the method will be renamed to
// ResourceType.
// 2. This method has a known limitation (TODO: PHP-32955) for
// conditional-delete entries where the resource is missing, and the request
// url target a resource using identifier ex:
// Patient?identifier=my-code-system|ABC-12345.
func (m *EntryMeta) ResourceTypeBeta() (resource.Type, error) {
	var nilResourceType resource.Type
	if m.res != nil && m.Method() != cpb.HTTPVerbCode_PATCH {
		return containedresource.TypeOf(m.entry.GetResource()), nil
	}

	if m.bundleEntryURIs.requestURI != "" {
		return m.extractResourceTypeFromURL(m.bundleEntryURIs.requestURI)
	}

	if m.bundleEntryURIs.responseLocationURI != "" {
		return m.extractResourceTypeFromURL(m.bundleEntryURIs.responseLocationURI)
	}

	return nilResourceType, ErrResourceTypeCouldNotBeDetermined
}

// Method returns the HTTP method of the entry.
// Returns HTTPVerbCode_INVALID_UNINITIALIZED if the entry is not a request
// or its method is not set.
func (m *EntryMeta) Method() cpb.HTTPVerbCode_Value {
	return m.entry.GetRequest().GetMethod().GetValue()
}

// IsMutate returns true of the entry is a mutation: POST, PUT, PATCH or DELETE.
func (m *EntryMeta) IsMutate() bool {
	switch m.Method() {
	case cpb.HTTPVerbCode_POST, cpb.HTTPVerbCode_PUT, cpb.HTTPVerbCode_PATCH, cpb.HTTPVerbCode_DELETE:
		return true
	}
	return false
}

func IndexOfEntryMeta(m *EntryMeta) int {
	return m.Index()
}

// EntryMetaMultimap is an multimap of metas with keys of type K.
// For each key k, the multimap holds an unordered set of EntryMetas.
//
// The initial use case for EntryMetaMultimap is to build an index from
// the effective URI of "plain" entries to the provenance resource(s)
// that target them. The relationship is many-to-many: a Provenance
// may target multiple plain resources, and a plain resource may be
// targeted by multiple Provenances.
//
// While the type name borrows from Java's Multimap, the implementation
// is a map of hash sets.
type EntryMetaMultimap[K comparable] struct {
	metaSetsByKey map[K]*set.Set[*EntryMeta]
}

// NewEntryMetaMultimap constructs a multimap of metas with keys of type K.
func NewEntryMetaMultimap[K comparable]() *EntryMetaMultimap[K] {
	return &EntryMetaMultimap[K]{
		metaSetsByKey: map[K]*set.Set[*EntryMeta]{},
	}
}

// GetSetForKey returns the set of metas for the given key as a HashSet.
func (m *EntryMetaMultimap[K]) GetSetForKey(key K) *set.Set[*EntryMeta] {
	keySet, ok := m.metaSetsByKey[key]
	if !ok {
		keySet = set.New[*EntryMeta](0)
		m.metaSetsByKey[key] = keySet
	}
	return keySet
}

// GetAll returns the set of metas for the given key as a slice.
// The returned slice is unordered and may be non-deterministic.
func (m *EntryMetaMultimap[K]) GetAllForKey(key K) []*EntryMeta {
	keySet, ok := m.metaSetsByKey[key]
	if !ok {
		return nil
	}
	return keySet.Slice()
}

// Add adds the given meta to the map under the given key.
// Duplicate additions of the same meta to the same key will be
// ignored. The same meta may be added to multiple keys.
func (m *EntryMetaMultimap[K]) Add(key K, meta *EntryMeta) {
	m.GetSetForKey(key).Insert(meta)
}

// isUnversionedRelativeURL determines if the URL is in the <resourceType>/<ID>
// format, for example, "Patient/12345".
//
// TODO (PHP-38264): This function should be moved to a more common location
// since it is not specific to bundle entries.
func isUnversionedRelativeURL(url string) bool {
	urlComponents := strings.Split(url, "/")
	if len(urlComponents) != 2 {
		return false
	}

	// WATCHOUT: The identity is created to ensure the URL components are
	// non-empty, conform to the FHIR resource ID pattern, and the resource type
	// is valid.
	_, err := resource.NewIdentityFromURL(url)
	return err == nil
}
