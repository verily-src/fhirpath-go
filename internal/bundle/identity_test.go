package bundle_test

import (
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	bcrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/internal/bundle"
	"github.com/verily-src/fhirpath-go/internal/resource"
)

func newIdentity(t *testing.T, typeName, id, version string) *resource.Identity {
	t.Helper()
	ident, err := resource.NewIdentity(typeName, id, version)
	if err != nil {
		t.Fatalf("NewIdentity: %v", err)
	}
	return ident
}

func TestIdentityOfResponse_NoLocation_ReturnsError(t *testing.T) {
	response := &bcrpb.Bundle_Entry_Response{}
	_, err := bundle.IdentityOfResponse(response)

	if !cmp.Equal(err, bundle.ErrNoLocation, cmpopts.EquateErrors()) {
		t.Errorf("IdentityOfResponse error got %v, want nil", err)
	}
}

func TestIdentityOfResponse(t *testing.T) {
	wantIdentity := newIdentity(t, "Patient", "123", "abc")
	response := &bcrpb.Bundle_Entry_Response{
		Location: &dtpb.Uri{
			Value: "https://healthcare.googleapis.com/v1/projects/123/locations/abc/datasets/def/fhirStores/ghi/fhir/Patient/123/_history/abc",
		},
	}

	ident, err := bundle.IdentityOfResponse(response)

	if err != nil {
		t.Fatalf("IdentityOfResponse error got %v, want nil", err)
	}
	if got, want := ident, wantIdentity; !got.Equal(want) {
		t.Errorf("IdentityOfResponse got %s, want %s", got, want)
	}
}
