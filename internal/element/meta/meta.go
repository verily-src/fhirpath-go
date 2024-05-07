package meta

import (
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Option is an option interface for modifying meta in place.
type Option interface {
	updateMeta(meta *dtpb.Meta)
}

// Update updates meta in place with given opts.
func Update(meta *dtpb.Meta, opts ...Option) *dtpb.Meta {
	for _, opt := range opts {
		opt.updateMeta(meta)
	}
	return meta
}

// WithTags replaces meta.tag.
func WithTags(tags ...*dtpb.Coding) Option {
	return withCodingOpt(tags)
}

type withCodingOpt []*dtpb.Coding

func (wco withCodingOpt) updateMeta(meta *dtpb.Meta) {
	meta.Tag = wco
}

// WithExtensions replaces meta.extension.
func WithExtensions(exts ...*dtpb.Extension) Option {
	return withExtensionOpt(exts)
}

type withExtensionOpt []*dtpb.Extension

func (weo withExtensionOpt) updateMeta(meta *dtpb.Meta) {
	meta.Extension = weo
}

// IncludeTags appends to meta.tag.
func IncludeTags(tags ...*dtpb.Coding) Option {
	return includeCodingOpt(tags)
}

type includeCodingOpt []*dtpb.Coding

func (ico includeCodingOpt) updateMeta(meta *dtpb.Meta) {
	meta.Tag = append(meta.Tag, ico...)
}

// WithProfiles replaces meta.profile.
func WithProfiles(profiles ...*dtpb.Canonical) Option {
	return withCanonicalOpt(profiles)
}

type withCanonicalOpt []*dtpb.Canonical

func (wco withCanonicalOpt) updateMeta(meta *dtpb.Meta) {
	meta.Profile = wco
}

// IncludeProfiles appends to meta.profile.
func IncludeProfiles(profiles ...*dtpb.Canonical) Option {
	return includeCanonicalOpt(profiles)
}

type includeCanonicalOpt []*dtpb.Canonical

func (ico includeCanonicalOpt) updateMeta(meta *dtpb.Meta) {
	meta.Profile = append(meta.Profile, ico...)
}

// ReplaceInResource replaces the resource meta field with the provided meta
// object.
func ReplaceInResource(resource fhir.Resource, meta *dtpb.Meta) {
	reflect := resource.ProtoReflect()
	metaField := getMetaField(reflect)

	reflect.Set(metaField, protoreflect.ValueOfMessage(meta.ProtoReflect()))
}

// EnsureInResource ensures that the resource meta field exists.
func EnsureInResource(resource fhir.Resource) {
	if resource.GetMeta() == nil {
		ReplaceInResource(resource, &dtpb.Meta{})
	}
}

func getMetaField(reflect protoreflect.Message) protoreflect.FieldDescriptor {
	return reflect.Descriptor().Fields().ByName("meta")
}
