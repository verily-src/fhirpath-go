package bundle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	bpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/binary_go_proto"
	bcrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"github.com/google/uuid"
	"github.com/verily-src/fhirpath-go/internal/containedresource"
	"github.com/verily-src/fhirpath-go/internal/element/reference"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/resource"
)

var (
	ErrInvalidIdentity             = fmt.Errorf("invaild resource identity")
	ErrInvalidPayload              = fmt.Errorf("invalid payload")
	ErrMissingPayload              = fmt.Errorf("%w: nil or empty payload", ErrInvalidPayload)
	ErrRspBundleEntryMissingStatus = fmt.Errorf("missing Bundle_Entry.response.status")
	ErrRspBundleEntryShortStatus   = fmt.Errorf("invalid (too short) Bundle_Entry.response.status")
)

// EntryOption is an option interface for constructing bundle entries
// from raw data.
type EntryOption interface {
	updateEntry(entry *bcrpb.Bundle_Entry)
}

// fullURLOpt is a bundle entry option for including a full url.
type fullURLOpt string

func (o fullURLOpt) updateEntry(entry *bcrpb.Bundle_Entry) {
	if o != "" {
		entry.FullUrl = &dtpb.Uri{
			Value: string(o),
		}
	}
}

// ifNoneExistOpt is a bundle entry option for including the If-None-Exist header
type ifNoneExistOpt struct {
	identifier *dtpb.Identifier
}

// ifNoneExistOpt.updateEntry sets the If-None-Exist header of a POST entry
// in the format `identifier=system|value`
// Either system or value can be empty (see https://hl7.org/fhir/R4/search.html#token)
// If system is empty, only `identifier=value` supported; not `identifier=|value`
func (o ifNoneExistOpt) updateEntry(entry *bcrpb.Bundle_Entry) {
	req := entry.Request
	if req == nil || req.GetMethod().GetValue() != cpb.HTTPVerbCode_POST {
		return
	}
	var sysVal string
	if sys := o.identifier.GetSystem().GetValue(); sys != "" {
		sysVal += sys + "|"
	}
	sysVal += o.identifier.GetValue().GetValue()
	if sysVal == "" {
		return
	}
	req.IfNoneExist = fhir.String(
		fmt.Sprintf("identifier=%s", sysVal),
	)
}

// WithFullURL adds a FullUrl field to a BundleEntry.
func WithFullURL(url string) EntryOption {
	return fullURLOpt(url)
}

// WithGeneratedFullURL adds a randomly generated FullUrl
// field to a BundleEntry, taking the form urn:uuid:$UUID.
func WithGeneratedFullURL() EntryOption {
	url := fmt.Sprintf("urn:uuid:%s", uuid.NewString())
	return fullURLOpt(url)
}

// WithIfNoneExist adds an identifier to the If-None-Exist request header of
// a POST BundleEntry, in the format `identifier=system|value`
func WithIfNoneExist(identifier *dtpb.Identifier) EntryOption {
	return ifNoneExistOpt{identifier: identifier}
}

// NewPostEntry wraps a FHIR Resource for a transaction bundle to be written to
// storage via a POST request.
//
// This is based on GCP documentation here:
// https://cloud.google.com/healthcare-api/docs/how-tos/fhir-bundles#resolving_references_to_resources_created_in_a_bundle
func NewPostEntry(res fhir.Resource, opts ...EntryOption) *bcrpb.Bundle_Entry {
	typeString := resource.TypeOf(res)
	entry := bundleEntry(cpb.HTTPVerbCode_POST, res, string(typeString))
	return applyOptions(entry, opts)
}

// NewPutEntry wraps a FHIR Resource for a transaction bundle to be written to
// storage via a PUT request. BundleEntry.fullUrl is not populated.
//
// This is based on GCP documentation here:
// https://cloud.google.com/healthcare-api/docs/how-tos/fhir-bundles#resolving_references_to_resources_created_in_a_bundle
func NewPutEntry(res fhir.Resource, opts ...EntryOption) *bcrpb.Bundle_Entry {
	uri := resource.URIString(res)
	entry := bundleEntry(cpb.HTTPVerbCode_PUT, res, uri)
	return applyOptions(entry, opts)
}

// NewDeleteEntry constructs a delete resource operation via a DELETE request.
//
// For use within a batch or transaction bundle.
func NewDeleteEntry(typeName resource.Type, id string, opts ...EntryOption) *bcrpb.Bundle_Entry {
	url := string(typeName) + "/" + id
	entry := bundleEntry(cpb.HTTPVerbCode_DELETE, nil /*resource*/, url)
	return applyOptions(entry, opts)
}

// NewCollectionEntry takes in a FHIR Resource and creates a BundleEntry
// for the resource.
func NewCollectionEntry(res fhir.Resource) *bcrpb.Bundle_Entry {
	return &bcrpb.Bundle_Entry{
		Resource: containedresource.Wrap(res),
		FullUrl:  resource.URI(res),
	}
}

type patchOps []struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// validatePatch performs basic validation of the patch payload, disallowing
// unknown fields. It doesn't check existence of required fields and their
// values.
func validatePatch(payload []byte) error {
	if !json.Valid(payload) {
		return ErrInvalidPayload
	}
	var patch patchOps
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.DisallowUnknownFields()
	return dec.Decode(&patch)
}

// PatchEntryFromBytes takes in a Resource Identity and a byte array payload,
// and returns a BundleEntry corresponding to a JSON Patch for the resource.
//
// See https://cloud.google.com/healthcare-api/docs/how-tos/fhir-resources#executing_a_patch_request_in_a_fhir_bundle
// for the details of how PATCH payload can be generated for known fields.
//
// For PATCH operations generation consider using one of the Go jsonpatch libraries from https://jsonpatch.com/#go.
// See the ./bundle_example_test.go/TestPatchEntry*() for the list of usage examples.
func PatchEntryFromBytes(identity *resource.Identity, payload []byte) (*bcrpb.Bundle_Entry, error) {
	if identity == nil || identity.Type() == "" || identity.ID() == "" {
		return nil, fmt.Errorf("unable to patch resource with identity (%v): %w", identity, ErrInvalidIdentity)
	}
	if payload == nil || len(payload) == 0 {
		return nil, fmt.Errorf("unable to patch resource: %w", ErrMissingPayload)
	}
	if err := validatePatch(payload); err != nil {
		return nil, fmt.Errorf("unable to patch resource, %w: %w", ErrInvalidPayload, err)
	}

	br := &bcrpb.ContainedResource{
		OneofResource: &bcrpb.ContainedResource_Binary{
			Binary: &bpb.Binary{
				ContentType: &bpb.Binary_ContentTypeCode{Value: "application/json-patch+json"},
				Data:        &dtpb.Base64Binary{Value: payload},
			},
		},
	}
	return &bcrpb.Bundle_Entry{
		Resource: br,
		Request: &bcrpb.Bundle_Entry_Request{
			Method: &bcrpb.Bundle_Entry_Request_MethodCode{
				Value: cpb.HTTPVerbCode_PATCH,
			},
			Url: identity.RelativeURI(),
		},
	}, nil
}

// UnwrapEntry unwraps a bundle entry into a FHIR Resource.
//
// If the bundle entry is nil, or if the entry does not contain a resource, this
// function will return nil.
func UnwrapEntry(entry *bcrpb.Bundle_Entry) fhir.Resource {
	if entry == nil {
		return nil
	}
	return containedresource.Unwrap(entry.GetResource())
}

func bundleEntry(method cpb.HTTPVerbCode_Value, resource fhir.Resource, requestURL string) *bcrpb.Bundle_Entry {
	return &bcrpb.Bundle_Entry{
		Resource: containedresource.Wrap(resource),
		Request: &bcrpb.Bundle_Entry_Request{
			Method: &bcrpb.Bundle_Entry_Request_MethodCode{
				Value: method,
			},
			Url: &dtpb.Uri{
				Value: requestURL,
			},
		},
	}
}

func applyOptions(entry *bcrpb.Bundle_Entry, opts []EntryOption) *bcrpb.Bundle_Entry {
	for _, opt := range opts {
		opt.updateEntry(entry)
	}
	return entry
}

// EntryReference generates a FHIR Reference proto pointing to the Entry's
// resource.
func EntryReference(entry *bcrpb.Bundle_Entry) *dtpb.Reference {
	if entry.GetResource() == nil {
		return nil
	}
	res := UnwrapEntry(entry)
	resourceType := resource.TypeOf(res)
	return reference.Weak(resourceType, entry.GetFullUrl().GetValue())
}

// SetEntryIfMatch sets the entry.Request.IfMatch based on
// entry.Resource.Meta.VersionId.
func SetEntryIfMatch(entry *bcrpb.Bundle_Entry) {
	req := entry.GetRequest()
	// No request or a request with If-Match already set
	if req == nil || req.GetIfMatch() != nil {
		return
	}

	// If-Match is only respected for PUT and PATCH, and PATCH uses a specially-
	// constructed resource so it doesn't make sense to infer anything from it;
	// some day If-Match on DELETE might be respected, update this code then
	if req.GetMethod().GetValue() != cpb.HTTPVerbCode_PUT {
		return
	}

	// Set If-Match should be based on the entry.Resource
	version := resource.VersionETag(UnwrapEntry(entry))
	if version != "" {
		req.IfMatch = fhir.String(version)
	}
}

// StatusCodeFromEntry returns the numeric http code from
// Bundle_Entry.response.status.
func StatusCodeFromEntry(entry *bcrpb.Bundle_Entry) (int, error) {
	status := entry.GetResponse().GetStatus()
	if status == nil {
		return 0, ErrRspBundleEntryMissingStatus
	}
	statusLine := status.GetValue()
	if len(statusLine) < 3 {
		return 0, ErrRspBundleEntryShortStatus
	}
	code, err := strconv.Atoi(statusLine[0:3])
	if err != nil {
		return 0, fmt.Errorf("invalid (atoi) Bundle_Entry.response.status: %v", err)
	}
	return code, nil
}
