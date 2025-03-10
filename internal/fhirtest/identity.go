package fhirtest

import (
	"testing"

	bcrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"github.com/verily-src/fhirpath-go/internal/bundle"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/resource"
)

// NewIdentity creates a new resource identity. If an error occurs, it logs the error and fails the test.
func NewIdentity(t *testing.T, resourceType, id, versionID string) *resource.Identity {
	t.Helper()
	identity, err := resource.NewIdentity(resourceType, id, versionID)
	if err != nil {
		t.Fatalf("NewIdentity: %v", err)
	}
	return identity
}

// NewIdentityFromURL creates a resource identity from the provided URL. If an error occurs, it logs the error and fails
// the test.
func NewIdentityFromURL(t *testing.T, url string) *resource.Identity {
	t.Helper()
	identity, err := resource.NewIdentityFromURL(url)
	if err != nil {
		t.Fatalf("url %s can't be converted to an identity", url)
	}
	return identity
}

// NewIdentityOf creates an identity from the provided resource. If an error occurs, it logs the error and fails the
// test.
func NewIdentityOf(t *testing.T, res fhir.Resource) *resource.Identity {
	t.Helper()
	identity, ok := resource.IdentityOf(res)
	if !ok {
		t.Fatalf("resource %T does not have identity", res)
	}
	return identity
}

// NewIdentityFromBundle returns the identity of the i-th resource in the given bundle. If the bundle does not contain
// the resource, or if the resource does not contain an ID, the test fails and logs the error.
func NewIdentityFromBundle(t *testing.T, bndl *bcrpb.Bundle, i int) *resource.Identity {
	t.Helper()
	entries := bndl.GetEntry()
	if i < 0 || i >= len(entries) {
		t.Fatalf("failed to get %d-th entry, bundle has %d entries", i, len(entries))
	}
	return NewIdentityOf(t, bundle.UnwrapEntry(bndl.GetEntry()[i]))
}
