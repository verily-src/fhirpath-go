/*
Package bundle provides utilities for working with FHIR R4 Bundle proto
definitions. This includes functionality for constructing/wrapping and
unwrapping bundle/entry objects.
*/
package bundle

import (
	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	bcrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"github.com/verily-src/fhirpath-go/internal/slices"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/bundleopt"
	"github.com/verily-src/fhirpath-go/internal/resource"
)

// New creates a new New by building it from the bundle options.
// Use of this function directly is discouraged; prefer to use the various
// *New functions instead for the explicit types.
func New(bundleType cpb.BundleTypeCode_Value, options ...Option) *bcrpb.Bundle {
	bundle := &bcrpb.Bundle{
		Type: &bcrpb.Bundle_TypeCode{
			Value: bundleType,
		},
	}
	return Extend(bundle, options...)
}

// Extend extends an existing bundle with the provided bundle options.
// The options will be extended in-place; the return value is not necessary to be
// looked at, but is available for convenience when used in fluent APIs.
//
// This decision was made to avoid cloning the bundle per-invocation, since in a
// loop this would grow the cost involved with calling this function substantially.
func Extend(bundle *bcrpb.Bundle, opts ...Option) *bcrpb.Bundle {
	bundleopt.Apply(bundle, opts...)
	return bundle
}

// NewTransaction is a helper function for building a transaction bundle.
func NewTransaction(options ...Option) *bcrpb.Bundle {
	return New(cpb.BundleTypeCode_TRANSACTION, options...)
}

// NewCollection is a helper function for building a collection bundle.
func NewCollection(options ...Option) *bcrpb.Bundle {
	return New(cpb.BundleTypeCode_COLLECTION, options...)
}

// NewBatch is a helper function for building a batch bundle.
func NewBatch(options ...Option) *bcrpb.Bundle {
	return New(cpb.BundleTypeCode_BATCH, options...)
}

// NewHistory is a helper function for building a history bundle.
func NewHistory(options ...Option) *bcrpb.Bundle {
	return New(cpb.BundleTypeCode_HISTORY, options...)
}

// NewSearchset is a helper function for building a searchset bundle.
func NewSearchset(options ...Option) *bcrpb.Bundle {
	return New(cpb.BundleTypeCode_SEARCHSET, options...)
}

// Unwrap unwraps a bundle into a slice of resources.
func Unwrap(bundle *bcrpb.Bundle) []fhir.Resource {
	return slices.Map(bundle.GetEntry(), UnwrapEntry)
}

// UnwrapMap unwraps a bundle into a map indexed by resource type.
func UnwrapMap(bundle *bcrpb.Bundle) map[resource.Type][]fhir.Resource {
	resourceMap := map[resource.Type][]fhir.Resource{}
	resources := Unwrap(bundle)
	for _, res := range resources {
		resourceType := resource.TypeOf(res)
		resourceMap[resourceType] = append(resourceMap[resourceType], res)
	}
	return resourceMap
}
