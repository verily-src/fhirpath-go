package bundle

import (
	"encoding/json"
	"fmt"
	"sort"
	"testing"

	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	bcrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	epb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/encounter_go_proto"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/internal/element/reference"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/resource"
)

func TestEntryMeta_Getters(t *testing.T) {
	encounterType, err := resource.NewType("Encounter")
	if err != nil {
		t.Fatalf("Failed to create encounter resource type %v", err)
	}

	patientType, err := resource.NewType("Patient")
	if err != nil {
		t.Fatalf("Failed to create patient resource type %v", err)
	}

	deleteResourceURL := "https://dg-integration-test-fhir-store/fhir/Encounter/ServerId-00001/_history/v00001"
	deleteResponse := &bcrpb.Bundle_Entry{
		Response: &bcrpb.Bundle_Entry_Response{
			Location: fhir.URI(deleteResourceURL),
		},
	}

	postResponse := &bcrpb.Bundle_Entry{
		Resource: &bcrpb.ContainedResource{
			OneofResource: &bcrpb.ContainedResource_Encounter{
				Encounter: &epb.Encounter{},
			},
		},
	}

	patientIdentity, _ := resource.NewIdentity("Patient", "123", "")
	patch := []map[string]interface{}{
		{
			"op":    "replace",
			"path":  "/active",
			"value": true,
		},
	}
	patchPayload, err := json.Marshal(patch)
	if err != nil {
		t.Fatalf("Failed to marshal patch payload: %v", err)
	}
	patchEntry, err := PatchEntryFromBytes(patientIdentity, patchPayload)
	if err != nil {
		t.Fatalf("Failed to create patch entry from bytes: %v", err)
	}

	testCases := []struct {
		name                string
		entry               *bcrpb.Bundle_Entry
		wantResourceType    resource.Type
		wantResourceTypeErr error
		wantEffectiveURIs   []string
		wantMethod          cpb.HTTPVerbCode_Value
	}{
		{
			name:                "Empty Entry",
			entry:               &bcrpb.Bundle_Entry{},
			wantResourceType:    "",
			wantResourceTypeErr: ErrResourceTypeCouldNotBeDetermined,
			wantEffectiveURIs:   []string{},
			wantMethod:          cpb.HTTPVerbCode_INVALID_UNINITIALIZED,
		},
		{
			name:                "Invalid Request",
			entry:               NewGetEntry("", ""),
			wantResourceType:    "",
			wantResourceTypeErr: fmt.Errorf("%w: non-REST and missing scheme component", reference.ErrInvalidURI),
			wantEffectiveURIs:   []string{"/"},
			wantMethod:          cpb.HTTPVerbCode_GET,
		},
		{
			name:                "Put Request",
			entry:               NewPutEntry(&epb.Encounter{Id: &dpb.Id{Value: "00012"}}),
			wantResourceType:    encounterType,
			wantResourceTypeErr: nil,
			wantEffectiveURIs:   []string{"Encounter/00012"},
			wantMethod:          cpb.HTTPVerbCode_PUT,
		},
		{
			name:                "Post Request",
			entry:               NewPostEntry(&epb.Encounter{Id: &dpb.Id{Value: "00013"}}, WithFullURL("urn:uuid:5551")),
			wantResourceType:    encounterType,
			wantResourceTypeErr: nil,
			wantEffectiveURIs:   []string{"urn:uuid:5551", "Encounter/00013"},
			wantMethod:          cpb.HTTPVerbCode_POST,
		},
		{
			name:                "Get Request",
			entry:               NewGetEntry("Patient", "123", WithFullURL("urn:uuid:5555")),
			wantResourceType:    patientType,
			wantResourceTypeErr: nil,
			wantEffectiveURIs:   []string{"urn:uuid:5555", "Patient/123"},
			wantMethod:          cpb.HTTPVerbCode_GET,
		},
		{
			name:                "Delete Request",
			entry:               NewDeleteEntry("Patient", "123"),
			wantResourceType:    patientType,
			wantResourceTypeErr: nil,
			wantEffectiveURIs:   []string{"Patient/123"},
			wantMethod:          cpb.HTTPVerbCode_DELETE,
		},
		{
			name:                "Patch Request",
			entry:               patchEntry,
			wantResourceType:    patientType,
			wantResourceTypeErr: nil,
			wantEffectiveURIs:   []string{"Patient/123"},
			wantMethod:          cpb.HTTPVerbCode_PATCH,
		},
		{
			name:                "Post Response",
			entry:               postResponse,
			wantResourceType:    encounterType,
			wantResourceTypeErr: nil,
			wantEffectiveURIs:   []string{},
			wantMethod:          cpb.HTTPVerbCode_INVALID_UNINITIALIZED,
		},
		{
			name:                "Delete Response",
			entry:               deleteResponse,
			wantResourceType:    encounterType,
			wantResourceTypeErr: nil,
			wantEffectiveURIs:   []string{},
			wantMethod:          cpb.HTTPVerbCode_INVALID_UNINITIALIZED,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test ResourceType and its error
			meta := NewEntryMeta(0, tc.entry)
			got, err := meta.ResourceTypeBeta()
			if want := tc.wantResourceType; got != want {
				t.Errorf("ResourceType() mismatch got [%v], want [%v]", got, want)
			}
			if want := tc.wantResourceTypeErr; want != nil && err.Error() != want.Error() {
				t.Errorf("ResourceType() error mismatch got [%v], want [%v]", err, want)
			}

			// Test Method
			if got, want := meta.Method(), tc.wantMethod; got != want {
				t.Errorf("Method() mismatch got [%v], want [%v]", got, want)
			}

			// Test EffectiveURLs
			gotURLs := meta.EffectiveURIs()

			// WATCHOUT: The order of the effective URLs is unspecified and may
			// be nondeterministic. We sort the slices before comparing them.
			sort.Strings(gotURLs)
			sort.Strings(tc.wantEffectiveURIs)
			if diff := cmp.Diff(gotURLs, tc.wantEffectiveURIs); diff != "" {
				t.Errorf("EffectiveURLs() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

func TestEntryMetaMultipmap(t *testing.T) {
	emptyEntry := &bcrpb.Bundle_Entry{}
	mm := NewEntryMetaMultimap[string]()

	meta4 := NewEntryMeta(4, emptyEntry)
	meta5 := NewEntryMeta(6, emptyEntry)
	meta6 := NewEntryMeta(5, emptyEntry)

	// A key may hold multiple metas.
	mm.Add("even", meta4)
	mm.Add("even", meta6)

	// Duplicate adds will be ingored.
	mm.Add("odd", meta5)
	mm.Add("odd", meta5)
	mm.Add("odd", meta5)

	// A meta may be indexed under multiple keys.
	mm.Add("all", meta4)
	mm.Add("all", meta5)
	mm.Add("all", meta6)

	testCases := []struct {
		name      string
		wantMetas []*EntryMeta
	}{
		{"even", []*EntryMeta{meta4, meta6}},
		{"odd", []*EntryMeta{meta5}},
		{"all", []*EntryMeta{meta4, meta5, meta6}},
		{"does-not-exist", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertEqualUnorderedEntryMetaSlices(t, mm.GetAllForKey(tc.name), tc.wantMetas,
				fmt.Sprintf("GetAllForKey(%s)", tc.name))
		})
	}
}

func assertEqualUnorderedEntryMetaSlices(t *testing.T, got, want []*EntryMeta, message string) {
	if diff := cmp.Diff(got, want, cmpopts.SortSlices(entryMetaLess)); diff != "" {
		t.Errorf("%s (-want, +got):\n%s", message, diff)
	}
}

func entryMetaLess(a, b *EntryMeta) bool {
	return a.Index() < b.Index()
}
