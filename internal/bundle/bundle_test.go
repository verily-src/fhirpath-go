package bundle_test

import (
	"testing"
	"time"

	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	bcrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/internal/bundle"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirtest"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"github.com/verily-src/fhirpath-go/internal/slices"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestBundle(t *testing.T) {
	iuTime := time.Date(1993, time.May, 16, 0, 0, 0, 0, time.UTC)
	iuPatient := &ppb.Patient{Id: fhir.ID("IU")}
	uaenaPatient := &ppb.Patient{Id: fhir.ID("Uaena")}
	firstEntry := bundle.NewPostEntry(iuPatient)
	secondEntry := bundle.NewPostEntry(uaenaPatient)
	baseBundle := &bcrpb.Bundle{Type: &bcrpb.Bundle_TypeCode{Value: cpb.BundleTypeCode_TRANSACTION}}
	timeBundle := proto.Clone(baseBundle).(*bcrpb.Bundle)
	timeBundle.Timestamp = fhir.Instant(iuTime)
	fullTxnBundle := &bcrpb.Bundle{
		Type:  &bcrpb.Bundle_TypeCode{Value: cpb.BundleTypeCode_TRANSACTION},
		Entry: []*bcrpb.Bundle_Entry{firstEntry, secondEntry},
	}

	fullSearchBundle := &bcrpb.Bundle{
		Type:  &bcrpb.Bundle_TypeCode{Value: cpb.BundleTypeCode_SEARCHSET},
		Entry: []*bcrpb.Bundle_Entry{firstEntry, secondEntry},
		Total: fhir.UnsignedInt(2),
	}

	testCases := []struct {
		name           string
		opts           []bundle.Option
		bundleTypeCode cpb.BundleTypeCode_Value
		wantBundle     *bcrpb.Bundle
	}{
		{
			name:           "no options",
			opts:           nil,
			bundleTypeCode: cpb.BundleTypeCode_TRANSACTION,
			wantBundle:     baseBundle,
		},
		{
			name:           "empty entries",
			opts:           []bundle.Option{bundle.WithEntries()},
			bundleTypeCode: cpb.BundleTypeCode_TRANSACTION,
			wantBundle:     baseBundle,
		},
		{
			name:           "multiple entries in a transaction bundle",
			opts:           []bundle.Option{bundle.WithEntries(firstEntry, secondEntry)},
			bundleTypeCode: cpb.BundleTypeCode_TRANSACTION,
			wantBundle:     fullTxnBundle,
		},
		{
			name:           "multiple entries in a search bundle",
			opts:           []bundle.Option{bundle.WithEntries(firstEntry, secondEntry)},
			bundleTypeCode: cpb.BundleTypeCode_SEARCHSET,
			wantBundle:     fullSearchBundle,
		},
		{
			name:           "timestamp",
			opts:           []bundle.Option{bundle.WithTimestamp(iuTime)},
			bundleTypeCode: cpb.BundleTypeCode_TRANSACTION,
			wantBundle:     timeBundle,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bundle := bundle.New(tc.bundleTypeCode, tc.opts...)

			got, want := bundle, tc.wantBundle
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Bundle(%s) mismatch (-want, +got):\n%s", tc.name, diff)
			}
		})
	}
}

func TestTransactionBundle_WithNoArguments_CreatesEmptyTransactionBundle(t *testing.T) {
	want := cpb.BundleTypeCode_TRANSACTION

	bundle := bundle.NewTransaction()

	got := bundle.GetType().GetValue()
	if got != want {
		t.Errorf("TransactionBundle: got %v, want %v", got, want)
	}
}

func TestCollectionBundle_WithNoArguments_CreatesEmptyCollectionBundle(t *testing.T) {
	want := cpb.BundleTypeCode_COLLECTION

	bundle := bundle.NewCollection()

	got := bundle.GetType().GetValue()
	if got != want {
		t.Errorf("CollectionBundle: got %v, want %v", got, want)
	}
}

func TestBatchBundle_WithNoArguments_CreatesEmptyBatchBundle(t *testing.T) {
	want := cpb.BundleTypeCode_BATCH

	bundle := bundle.NewBatch()

	got := bundle.GetType().GetValue()
	if got != want {
		t.Errorf("BatchBundle: got %v, want %v", got, want)
	}
}

func TestHistoryBundle_WithNoArguments_CreatesEmptyHistoryBundle(t *testing.T) {
	want := cpb.BundleTypeCode_HISTORY

	bundle := bundle.NewHistory()

	got := bundle.GetType().GetValue()
	if got != want {
		t.Errorf("HistoryBundle: got %v, want %v", got, want)
	}
}

func TestUnwrapBundle(t *testing.T) {
	testCases := []struct {
		name string
		want []fhir.Resource
	}{
		{"Empty", []fhir.Resource{}},
		{"Contains resources", []fhir.Resource{
			fhirtest.NewResource(t, resource.Patient),
			fhirtest.NewResource(t, resource.Patient),
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			put := func(resource fhir.Resource) *bcrpb.Bundle_Entry {
				return bundle.NewPutEntry(resource)
			}
			sut := bundle.NewCollection(
				bundle.WithEntries(slices.Map(tc.want, put)...),
			)

			got := bundle.Unwrap(sut)

			want := tc.want
			if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
				t.Errorf("UnwrapBundle: diff (-got,+want):\n%v", diff)
			}
		})
	}
}

func TestUnwrapMap(t *testing.T) {
	patientOne := fhirtest.NewResource(t, resource.Patient)
	patientTwo := fhirtest.NewResource(t, resource.Patient)
	taskOne := fhirtest.NewResource(t, resource.Task)
	taskTwo := fhirtest.NewResource(t, resource.Task)
	testCases := []struct {
		name string
		in   []fhir.Resource
		want map[resource.Type][]fhir.Resource
	}{
		{"Empty", []fhir.Resource{}, map[resource.Type][]fhir.Resource{}},
		{"Contains unique resource types", []fhir.Resource{
			patientOne,
			taskOne,
		}, map[resource.Type][]fhir.Resource{
			resource.Patient: {patientOne},
			resource.Task:    {taskOne},
		}},
		{"Contains multiple of a single type", []fhir.Resource{
			patientOne,
			patientTwo,
		}, map[resource.Type][]fhir.Resource{
			resource.Patient: {patientOne, patientTwo},
		}},
		{"Contains multiple of multiple types", []fhir.Resource{
			patientOne,
			patientTwo,
			taskOne,
			taskTwo,
		}, map[resource.Type][]fhir.Resource{
			resource.Task:    {taskOne, taskTwo},
			resource.Patient: {patientOne, patientTwo},
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			put := func(resource fhir.Resource) *bcrpb.Bundle_Entry {
				return bundle.NewPutEntry(resource)
			}
			sut := bundle.NewCollection(
				bundle.WithEntries(slices.Map(tc.in, put)...),
			)

			got := bundle.UnwrapMap(sut)

			want := tc.want
			if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
				t.Errorf("UnwrapBundle: diff (-got,+want):\n%v", diff)
			}
		})
	}
}
