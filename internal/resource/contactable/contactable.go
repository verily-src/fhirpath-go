// Package contactable contains utilities for working with FHIR Resource objects that
// include a contact field
package contactable

import (
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"github.com/verily-src/fhirpath-go/internal/resourceopt"
)

// Option is an option that may be supplied to updates of ContactableResource types
type Option = resourceopt.Option

// ContactableResource is the interface for FHIR resources that include a contact field
type ContactableResource interface {
	GetContact() []*dtpb.ContactDetail
	fhir.DomainResource
}

// WithContacts returns a resource Option for setting the ContactableResource
// Contact with the specified contact entry.
func WithContacts(contact ...*dtpb.ContactDetail) Option {
	return resourceopt.WithProtoField("contact", contact...)
}

// Update modifies the input resource in-place with the specified options.
func Update(cr ContactableResource, opts ...Option) {
	resource.Update(cr.(fhir.Resource), opts...)
}
