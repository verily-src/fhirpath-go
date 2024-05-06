package resource

import (
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/resourceopt"
)

// WithMeta returns a resource Option for setting the Resource Meta with the
// specified meta entry.
func WithMeta(meta *dtpb.Meta) Option {
	return resourceopt.WithProtoField("meta", meta)
}

// WithID returns a resource Option for setting the Resourec ID with the id of
// the provided string.
func WithID(id string) Option {
	return resourceopt.WithProtoField("id", fhir.ID(id))
}

// WithImplicitRules returns a resource Option for setting the Resource implicit
// rules with the provided string.
func WithImplicitRules(rules string) Option {
	return resourceopt.WithProtoField("implicit_rules", fhir.URI(rules))
}

// WithLanguage returns a resource Option for setting the Resource language to
// the code of the provided string.
func WithLanguage(language string) Option {
	return resourceopt.WithProtoField("language", fhir.Code(language))
}
