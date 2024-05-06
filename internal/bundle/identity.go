package bundle

import (
	"errors"

	bcrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"github.com/verily-src/fhirpath-go/internal/element/reference"
	"github.com/verily-src/fhirpath-go/internal/resource"
)

var (
	ErrNoLocation = errors.New("bundle entry response missing location")
)

// IdentityOfResponse returns a complete Identity
// (Type and ID always set, VersionID set if applicable) representing the
// location contained in the given bundle entry response.
func IdentityOfResponse(response *bcrpb.Bundle_Entry_Response) (*resource.Identity, error) {
	// Per the FHIR spec, location may be a relative or absolute URI, which may
	// include a _history component.
	location := response.GetLocation().GetValue()
	if location == "" {
		return nil, ErrNoLocation
	}

	return reference.IdentityFromURL(location)
}
