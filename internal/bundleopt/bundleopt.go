package bundleopt

import (
	bcrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
)

// BundleOption is an option interface for constructing bundles
// from raw data.
type BundleOption interface {
	updateBundle(bundle *bcrpb.Bundle)
}

type Transform func(b *bcrpb.Bundle)

func (t Transform) updateBundle(entry *bcrpb.Bundle) {
	t(entry)
}

func Apply(bundle *bcrpb.Bundle, opts ...BundleOption) {
	for _, opt := range opts {
		opt.updateBundle(bundle)
	}
}

var _ BundleOption = (*Transform)(nil)
