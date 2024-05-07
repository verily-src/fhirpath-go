package bundle

import (
	"time"

	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	bcrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/bundleopt"
)

// Option is an option interface for constructing bundles
// from raw data.
type Option = bundleopt.BundleOption

// entriesOpt is a bundle option for including bundle entries.
type entriesOpt []*bcrpb.Bundle_Entry

func (o entriesOpt) updateBundle(bundle *bcrpb.Bundle) {
	if len(o) > 0 {
		bundle.Entry = o
		if bundle.GetType().GetValue() == cpb.BundleTypeCode_SEARCHSET || bundle.GetType().GetValue() == cpb.BundleTypeCode_HISTORY {
			bundle.Total = fhir.UnsignedInt(uint32(len(o)))
		}
	}
}

// WithEntries adds bundle entries to a bundle.
func WithEntries(entries ...*bcrpb.Bundle_Entry) Option {
	entriesopt := entriesOpt(entries)
	return bundleopt.Transform(entriesopt.updateBundle)
}

// timeOpt is a bundle option for including a timestamp.
type timeOpt time.Time

func (o timeOpt) updateBundle(bundle *bcrpb.Bundle) {
	bundle.Timestamp = fhir.Instant(time.Time(o))
}

// WithTimestamp adds a given time to the bundle's timestamp.
func WithTimestamp(t time.Time) Option {
	timeOpt := timeOpt(t)
	return bundleopt.Transform(timeOpt.updateBundle)
}

// WithTimestampNow adds the current time to the bundle's timestamp.
func WithTimestampNow() Option {
	timeOpt := timeOpt(time.Now())
	return bundleopt.Transform(timeOpt.updateBundle)
}
