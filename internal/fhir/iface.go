package fhir

import (
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// Base is the interface-definition of the FHIR abstract base type which is the
// ancestor of all FHIR objects (both resources and elements).
//
// Represented in Go, this simply embeds the proto.Message interface, since this
// is a utility for the google/fhir proto definitions.
type Base interface {
	proto.Message
}

// Resource is the interface-definition of the FHIR Abstract base type which is
// the ancestor of all FHIR resources.
// See https://www.hl7.org/fhir/r4/resource.html#Resource for more details.
//
// This interface is defined by embedding the proto.Message interface, since all
// FHIR resources in this library must also be proto.Message types
type Resource interface {
	GetId() *dtpb.Id
	GetImplicitRules() *dtpb.Uri
	GetMeta() *dtpb.Meta
	GetLanguage() *dtpb.Code
	Base
}

// Extendable is an interface for abstraction resources or data-types that have
// extension properties.
//
// This is not an official FHIR abstract class; this is something simply named
// here for the general convenience, since not all FHIR types are extendable.
//
// This embeds the proto.Message interface into this interface to help distinguish
// that this still refers to protos in the process.
type Extendable interface {
	GetExtension() []*dtpb.Extension
	Base
}

// DomainResource is the interface-definition of the FHIR Abstract base type
// which is the ancestor of all FHIR domain resource objects (effectively
// everything that is not a datatype or bundle/contained-resource).
// See https://www.hl7.org/fhir/r4/domainresource.html for more details.
//
// This interface extends from the `Resource` interface by embedding it. Any
// `DomainResource` is also a `Resource`.
type DomainResource interface {
	GetText() *dtpb.Narrative
	GetContained() []*anypb.Any
	GetModifierExtension() []*dtpb.Extension
	Extendable
	Resource
}

// CanonicalResource represents resources that have a canonical URL:
//
//   - They have a canonical URL (note: all resources with a canonical URL are
//     specializations of this type)
//   - They have version, status, and data properties to help manage their publication
//   - They carry some additional metadata about their use, including copyright information
//
// CanonicalResource objects may be the logical target of Canonical references.
//
// Note: This is technically an "R5" interface type that is not officially part
// of the R4 spec, however its definition is still applicable and applies to "R4"
// resource types. Using this still provides us with a proper vernacular for
// referring to these resources.
//
// See https://www.hl7.org/fhir/r5/canonicalresource.html for more details.
type CanonicalResource interface {
	GetUrl() *dtpb.Uri
	GetIdentifier() []*dtpb.Identifier
	GetVersion() *dtpb.String
	GetName() *dtpb.String
	GetTitle() *dtpb.String
	GetExperimental() *dtpb.Boolean
	GetDate() *dtpb.DateTime
	GetPublisher() *dtpb.String
	GetContact() []*dtpb.ContactDetail
	GetDescription() *dtpb.Markdown
	GetUseContext() []*dtpb.UsageContext
	GetJurisdiction() []*dtpb.CodeableConcept
	GetPurpose() *dtpb.Markdown
	GetCopyright() *dtpb.Markdown
	// This interface should technically also have 'GetStatus()', however the
	// return type differs based on resource type in the proto definitions -- and
	// so this can't be referred to in a homogeneous way.
	// GetStatus() interface{}

	DomainResource
}

// MetadataResource represents resources that carry additional publication
// metadata over other CanonicalResources, describing their review and use in
// more details.
//
// As an interface, this type is never created directly.
//
// Note: This is technically an "R5" interface type that is not officially part
// of the R4 spec, however its definition is still applicable and applies to "R4"
// resource types. Using this still provides us with a proper vernacular for
// referring to these resources.
//
// See https://www.hl7.org/fhir/r5/metadataresource.html for more details.
type MetadataResource interface {
	GetApprovalDate() *dtpb.Date
	GetLastReviewDate() *dtpb.Date
	GetEffectivePeriod() *dtpb.Period
	GetTopic() []*dtpb.CodeableConcept
	GetAuthor() []*dtpb.ContactDetail
	GetEditor() []*dtpb.ContactDetail
	GetReviewer() []*dtpb.ContactDetail
	GetEndorser() []*dtpb.ContactDetail
	GetRelatedArtifact() []*dtpb.RelatedArtifact

	CanonicalResource
}

// Element is the base definition for all elements in a resource.
//
// See https://www.hl7.org/fhir/r4/element.html for more details.
//
// This interface is defined by embedding the proto.Message interface, since all
// FHIR elements in this library must also be proto.Message types.
type Element interface {
	GetId() *dtpb.String
	Extendable
	Base
}

// BackboneElement is the base definition for all elements that are defined
// inside a resource - but not those in a data type.
//
// See https://www.hl7.org/fhir/r4/backboneelement.html for more details.
//
// This interface is defined by embedding the proto.Message interface, since all
// FHIR backbone elements in this library must also be proto.Message types.
type BackboneElement interface {
	GetModifierExtension() []*dtpb.Extension
	Element
}
