package fhirtest

import (
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/element/reference"
	"github.com/verily-src/fhirpath-go/internal/resource"
)

// NewReferenceTyped constructs a new strong typed un-versioned FHIR reference. If an error occurs, this function will
// log it and fail the test.
func NewReferenceTyped(t *testing.T, resourceType resource.Type, resourceId string) *dtpb.Reference {
	ref, err := reference.Typed(resourceType, resourceId)
	if err != nil {
		t.Fatalf("NewReferenceTyped: %v", err)
	}
	return ref
}

// NewReferenceVersionedTyped constructs a new strong typed versioned FHIR reference. If an error occurs, this function
// will log it and fail the test.
func NewReferenceVersionedTyped(t *testing.T, resourceType resource.Type, resourceId string, versionId string) *dtpb.Reference {
	ref, err := reference.Typed(resourceType, resourceId+"/_history/"+versionId)
	if err != nil {
		t.Fatalf("NewReferenceVersionedTyped: %v", err)
	}
	return ref
}
