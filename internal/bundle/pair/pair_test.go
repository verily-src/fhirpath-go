package pair

import (
	"fmt"
	"testing"
	"time"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	bcrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	epb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/encounter_go_proto"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/internal/bundle"
	"github.com/verily-src/fhirpath-go/internal/element/reference"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestPair_New(t *testing.T) {
	encounter1 := &epb.Encounter{}
	encounter2 := &epb.Encounter{}
	modifiedTime, _ := time.Parse(time.RFC3339, "2019-01-02T01:02:03.123456-04:00")
	testCases := []struct {
		name         string
		reqBundle    *bcrpb.Bundle
		rspBundle    *bcrpb.Bundle
		wantIdentity *resource.Identity
		wantErr      error
	}{
		{
			name:         "create",
			reqBundle:    bundle.NewTransaction(bundle.WithEntries(bundle.NewPostEntry(encounter1))),
			rspBundle:    bundle.NewTransactionResponse(newResponseEntry("Encounter/123/_history/abc", modifiedTime)),
			wantIdentity: mustNewIdentity("Encounter", "123", "abc"),
		},
		{
			name:         "delete",
			reqBundle:    bundle.NewTransaction(bundle.WithEntries(bundle.NewDeleteEntry("Encounter", "123"))),
			rspBundle:    bundle.NewTransactionResponse(newOkDeleteResponseEntry("foo", modifiedTime)),
			wantIdentity: mustNewIdentity("Encounter", "123", "foo"),
		},
		{
			name:      "delete returns no etag",
			reqBundle: bundle.NewTransaction(bundle.WithEntries(bundle.NewDeleteEntry("Encounter", "123"))),
			rspBundle: bundle.NewTransactionResponse(&bcrpb.Bundle_Entry{
				Response: &bcrpb.Bundle_Entry_Response{
					Status: fhir.String("200 OK"),
				},
			}),
			wantErr: ErrRetrieveEtagVersion,
		},
		{
			name:      "delete returns unsuccessful status code",
			reqBundle: bundle.NewTransaction(bundle.WithEntries(bundle.NewDeleteEntry("Encounter", "123"))),
			rspBundle: bundle.NewTransactionResponse(&bcrpb.Bundle_Entry{
				Response: &bcrpb.Bundle_Entry_Response{
					Status: fhir.String("400 Client Error"),
					Etag:   fhir.String(`W/"foo"`),
				},
			}),
			wantErr: ErrUnexpectedResponseStatusCode,
		},
		{
			name:      "delete returns bad etag format",
			reqBundle: bundle.NewTransaction(bundle.WithEntries(bundle.NewDeleteEntry("Encounter", "123"))),
			rspBundle: bundle.NewTransactionResponse(&bcrpb.Bundle_Entry{
				Response: &bcrpb.Bundle_Entry_Response{
					Status: fhir.String("200 OK"),
					Etag:   fhir.String("foo"),
				},
			}),
			wantErr: ErrRetrieveEtagVersion,
		},
		{
			name:         "delete entry with old version in request URL",
			reqBundle:    bundle.NewTransaction(bundle.WithEntries(bundle.NewDeleteEntry("Encounter", "123/_history/foo"))),
			rspBundle:    bundle.NewTransactionResponse(newOkDeleteResponseEntry("bar", modifiedTime)),
			wantIdentity: mustNewIdentity("Encounter", "123", "bar"),
		},
		{
			name:      "mismatch lengths",
			reqBundle: bundle.NewTransaction(bundle.WithEntries(bundle.NewPostEntry(encounter1))),
			rspBundle: bundle.NewTransactionResponse(),
			wantErr:   ErrMismatchBundleLengths,
		},
		{
			name:      "bad bundle type",
			reqBundle: bundle.NewBatch(),
			rspBundle: bundle.NewTransactionResponse(),
			wantErr:   ErrUnsupportedBundleType,
		},
		{
			name: "dup URI",
			reqBundle: bundle.NewTransaction(bundle.WithEntries(
				bundle.NewPostEntry(encounter1, bundle.WithFullURL("urn:uuid:1111")),
				bundle.NewPostEntry(encounter2, bundle.WithFullURL("urn:uuid:1111")))),
			rspBundle: bundle.NewTransactionResponse(
				newResponseEntry("Encounter/123/_history/abc", modifiedTime),
				newResponseEntry("Encounter/456/_history/abc", modifiedTime),
			),
			wantErr: ErrDuplicateBundleEntryURI,
		},
		{
			name: "dup GET URI",
			reqBundle: bundle.NewTransaction(bundle.WithEntries(
				bundle.NewGetEntry("Encounter", "123"),
				bundle.NewGetEntry("Encounter", "123"))),
			rspBundle: bundle.NewTransactionResponse(
				newResponseEntry("Encounter/123/_history/abc", modifiedTime),
				newResponseEntry("Encounter/123/_history/abc", modifiedTime),
			),
			wantIdentity: mustNewIdentity("Encounter", "123", "abc"),
			wantErr:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pair, gotErr := NewPair(tc.reqBundle, tc.rspBundle)
			if !cmp.Equal(gotErr, tc.wantErr, cmpopts.EquateErrors()) {
				t.Fatalf("NewPair(%s) error mismatch: got [%v], want [%v]", tc.name, gotErr, tc.wantErr)
			}
			if gotErr == nil && !cmp.Equal(pair.rspIdentities[0], tc.wantIdentity) {
				t.Errorf("NewPair.rspIdentities[0](%s) identity mismatch: got [%v], want [%v]",
					tc.name, pair.rspIdentities[0], tc.wantIdentity)
			}
		})
	}
}

func TestPair_Getters(t *testing.T) {
	modifiedTime, _ := time.Parse(time.RFC3339, "2019-01-02T01:02:03.123456-04:00")
	reqBundle := bundle.NewTransaction(bundle.WithEntries(bundle.NewPostEntry(&epb.Encounter{})))
	rspBundle := bundle.NewTransactionResponse(newResponseEntry("Encounter/123/_history/abc", modifiedTime))
	rr, err := NewPair(Clone(reqBundle), Clone(rspBundle))
	if err != nil {
		t.Fatalf("NewPair: error [%v]", err)
	}
	if diff := cmp.Diff(rr.ReqBundle(), reqBundle, protocmp.Transform()); diff != "" {
		t.Errorf("ReqBundle() mismatch (-want, +got):\n%s", diff)
	}
	if diff := cmp.Diff(rr.RspBundle(), rspBundle, protocmp.Transform()); diff != "" {
		t.Errorf("RspBundle() mismatch (-want, +got):\n%s", diff)
	}
	if diff := cmp.Diff(rr.ReqEntryOfIdx(0), reqBundle.GetEntry()[0], protocmp.Transform()); diff != "" {
		t.Errorf("ReqEntryOfIdx(0) mismatch (-want, +got):\n%s", diff)
	}
	if diff := cmp.Diff(rr.RspEntryOfIdx(0), rspBundle.GetEntry()[0], protocmp.Transform()); diff != "" {
		t.Errorf("RspEntryOfIdx(0) mismatch (-want, +got):\n%s", diff)
	}
	if got, want := rr.ServiceBaseURL(), myServiceBaseURL; !cmp.Equal(got, want) {
		t.Fatalf("ServiceBaseURL mismatch: got [%v], want [%v]", got, want)
	}
}

// Tests the Pair.IdentityOfRef method.
// Algoritmically, the code-under-test:
//  1. Starts with a Reference
//  2. Maps that to a request entry in the request bundle
//  3. Maps that to its sibling response entry
//  4. Maps that to the resource.Identity of that response entry.
//
// Most of this steps have their own complexites and can fail in interesting ways.
// Thus the complete set of possible permutions is quite large.
func TestPair_IdentityOfRef(t *testing.T) {
	modifiedTime, _ := time.Parse(time.RFC3339, "2019-01-02T01:02:03.123456-04:00")
	uuidReqEntry := bundle.NewPostEntry(&epb.Encounter{}, bundle.WithFullURL("urn:uuid:1111"))
	uuidRspEntry := newResponseEntry("Encounter/123/_history/abc", modifiedTime)

	testCases := []struct {
		name         string
		reqEntry     *bcrpb.Bundle_Entry
		rspEntry     *bcrpb.Bundle_Entry
		inputRef     *dtpb.Reference
		wantIdentity *resource.Identity
		wantErr      error
	}{
		{
			name:         "POST UUID",
			reqEntry:     uuidReqEntry,
			rspEntry:     uuidRspEntry,
			inputRef:     reference.Weak("Encounter", "urn:uuid:1111"),
			wantIdentity: mustNewIdentity("Encounter", "123", "abc"),
		},
		{
			name:         "POST REST and relative ref",
			reqEntry:     bundle.NewPostEntry(&epb.Encounter{}, bundle.WithFullURL("Encounter/456")),
			rspEntry:     newResponseEntry("Encounter/456/_history/def", modifiedTime),
			inputRef:     reference.Weak("Encounter", "Encounter/456"),
			wantIdentity: mustNewIdentity("Encounter", "456", "def"),
		},
		{
			name:         "POST REST and absolute ref",
			reqEntry:     bundle.NewPostEntry(&epb.Encounter{}, bundle.WithFullURL("Encounter/456")),
			rspEntry:     newResponseEntry("Encounter/456/_history/def", modifiedTime),
			inputRef:     reference.Weak("Encounter", myServiceBaseURL+"/Encounter/456"),
			wantIdentity: mustNewIdentity("Encounter", "456", "def"),
		},
		{
			name:     "POST REST and other absolute ref",
			reqEntry: bundle.NewPostEntry(&epb.Encounter{}, bundle.WithFullURL("Encounter/456")),
			rspEntry: newResponseEntry("Encounter/456/_history/def", modifiedTime),
			inputRef: reference.Weak("Encounter", "http://someother/store"+"/Encounter/456"),
			wantErr:  ErrNotFoundBundleEntry,
		},
		{
			name:         "PUT no URL and relative ref",
			reqEntry:     bundle.NewPutEntry(&epb.Encounter{Id: fhir.ID("789")}),
			rspEntry:     newResponseEntry("Encounter/789/_history/ghi", modifiedTime),
			inputRef:     reference.Weak("Encounter", "Encounter/789"),
			wantIdentity: mustNewIdentity("Encounter", "789", "ghi"),
		},
		{
			name:     "malformed ref",
			reqEntry: uuidReqEntry,
			rspEntry: uuidRspEntry,
			inputRef: reference.Weak("Encounter", "NotAThing/333"),
			wantErr:  reference.ErrWeakInvalid,
		},
		{
			name:     "missing ref",
			reqEntry: uuidReqEntry,
			rspEntry: uuidRspEntry,
			inputRef: reference.Weak("Encounter", "Encounter/333"),
			wantErr:  ErrNotFoundBundleEntry,
		},
		{
			name:         "DELETE no fullURL",
			reqEntry:     bundle.NewDeleteEntry("Encounter", "666"),
			rspEntry:     newOkDeleteResponseEntry("foo", modifiedTime),
			inputRef:     reference.Weak("Encounter", "Encounter/666"),
			wantIdentity: mustNewIdentity("Encounter", "666", "foo"),
		},
		{
			name:         "DELETE REST",
			reqEntry:     bundle.NewDeleteEntry("Encounter", "666", bundle.WithFullURL("Encounter/666")),
			rspEntry:     newOkDeleteResponseEntry("foo", modifiedTime),
			inputRef:     reference.Weak("Encounter", "Encounter/666"),
			wantIdentity: mustNewIdentity("Encounter", "666", "foo"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			bp, err := NewPair(
				bundle.NewTransaction(bundle.WithEntries(tc.reqEntry)),
				bundle.NewTransactionResponse(tc.rspEntry))
			if err != nil {
				t.Fatalf("NewPair(%s) setup error: %v", tc.name, err)
			}
			gotIdentity, gotErr := bp.IdentityOfRef(tc.inputRef)
			if !cmp.Equal(gotErr, tc.wantErr, cmpopts.EquateErrors()) {
				t.Fatalf("Pair.IdentityOfRef(%s) error mismatch: got [%v], want [%v]",
					tc.name, gotErr, tc.wantErr)
			}
			if !cmp.Equal(gotIdentity, tc.wantIdentity) {
				t.Errorf("Pair.IdentityOfRef(%s) identity mismatch: got [%v], want [%v]",
					tc.name, gotIdentity, tc.wantIdentity)
			}
		})
	}
}

const myServiceBaseURL = "http://myfhir/store"

func newResponseEntry(versionedURI string, lastModified time.Time) *bcrpb.Bundle_Entry {
	// WATCHOUT: we don't populate entry.fullURL. I'm not sure what should go here.
	return &bcrpb.Bundle_Entry{Response: &bcrpb.Bundle_Entry_Response{
		Status:       fhir.String("200 OK"),
		Location:     fhir.URI(myServiceBaseURL + "/" + versionedURI),
		LastModified: fhir.Instant(lastModified),
	}}

}

func newOkDeleteResponseEntry(versionID string, lastModified time.Time) *bcrpb.Bundle_Entry {
	return &bcrpb.Bundle_Entry{
		Response: &bcrpb.Bundle_Entry_Response{
			Status:       fhir.String("200 OK"),
			Etag:         fhir.String(fmt.Sprintf(`W/"%v"`, versionID)),
			LastModified: fhir.Instant(lastModified),
		},
	}
}

func mustNewIdentity(resourceType, id, versionID string) *resource.Identity {
	identity, err := resource.NewIdentity(resourceType, id, versionID)
	if err != nil {
		panic(err)
	}
	return identity
}

// Clone returns a deep copy of the
func Clone(bundle *bcrpb.Bundle) *bcrpb.Bundle {
	return proto.Clone(bundle).(*bcrpb.Bundle)
}
