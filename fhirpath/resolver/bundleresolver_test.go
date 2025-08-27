package resolver_test

import (
	"testing"

	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	r4pb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/observation_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/fhirpath/resolver"
	"github.com/verily-src/fhirpath-go/internal/containedresource"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"google.golang.org/protobuf/testing/protocmp"
)

// Helper to create a Bundle for testing
func createBundle(entries []*r4pb.Bundle_Entry) *r4pb.Bundle {
	return &r4pb.Bundle{
		Type:  &r4pb.Bundle_TypeCode{Value: cpb.BundleTypeCode_SEARCHSET},
		Entry: entries,
	}
}

func TestBundleResolver_Resolve(t *testing.T) {
	uuidRef := fhir.UUID("c757873d-ec9a-4326-a141-556f43239520")
	invalidUuidRef := fhir.UUID("invalid-uuid")

	oidRef := fhir.OID("1.2.3.4.5")

	patAbsVersionlessUrl := "https://healthcare.googleapis.com/v1/projects/123/locations/abc/datasets/def/fhirStores/ghi/fhir/Patient/123"
	obsAbsVersionlessUrl := "https://healthcare.googleapis.com/v1/projects/123/locations/abc/datasets/def/fhirStores/ghi/fhir/Observation/123"
	invalidAbsVersionlessUrl := "https://healthcare.googleapis.com/v1/projects/123/locations/abc/datasets/def/fhirStores/ghi/fhir/Whale/123"

	patAbsVersionedUrl := "https://healthcare.googleapis.com/v1/projects/123/locations/abc/datasets/def/fhirStores/ghi/fhir/Patient/123/_history/v1"
	obsAbsVersionedUrl := "https://healthcare.googleapis.com/v1/projects/123/locations/abc/datasets/def/fhirStores/ghi/fhir/Observation/123/_history/v1"
	invalidVersionedAbsUrl := "https://healthcare.googleapis.com/v1/projects/123/locations/abc/datasets/def/fhirStores/ghi/fhir/Patient/123/_history/%#@^%$#"

	obs123 := &observation_go_proto.Observation{
		Id: &dtpb.Id{Value: "123"},
		Subject: &dtpb.Reference{
			Reference: &dtpb.Reference_PatientId{PatientId: &dtpb.ReferenceId{Value: "123"}},
		},
	}

	patient123 := &ppb.Patient{
		Id: fhir.ID("123"),
		Meta: &dtpb.Meta{
			LastUpdated: &dtpb.Instant{ValueUs: 1000},
		},
	}
	crPatient123 := containedresource.Wrap(patient123)

	patient123Latest := &ppb.Patient{
		Id: fhir.ID("123"),
		Meta: &dtpb.Meta{
			LastUpdated: &dtpb.Instant{ValueUs: 500},
		},
	}

	patientMissingMeta := &ppb.Patient{
		Id: fhir.ID("789"),
	}
	crPatientMissingMeta := containedresource.Wrap(patientMissingMeta)

	patientMissingLastUpdated := &ppb.Patient{
		Id:   fhir.ID("456"),
		Meta: &dtpb.Meta{},
	}
	crPatientMissingLastUpdated := containedresource.Wrap(patientMissingLastUpdated)

	patient123VersionedV1 := &ppb.Patient{
		Id: fhir.ID("123"),
		Meta: &dtpb.Meta{
			VersionId: fhir.ID("v1"),
		},
	}

	patient123VersionedV1Copy := &ppb.Patient{
		Id: fhir.ID("123"),
		Meta: &dtpb.Meta{
			VersionId: fhir.ID("v1"),
		},
	}

	testCases := []struct {
		name          string
		bundle        *r4pb.Bundle
		toResolveRefs []string
		wantResources []fhir.Resource
		wantErr       error
	}{
		{
			name:          "Empty Bundle",
			bundle:        createBundle([]*r4pb.Bundle_Entry{}),
			toResolveRefs: []string{"Patient/123"},
			wantResources: []fhir.Resource{},
		},
		{
			name: "Empty toResolve References",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URIFromOID(oidRef)},
			}),
			toResolveRefs: []string{},
			wantResources: []fhir.Resource{},
		},
		{
			name: "Resource not found - returns empty",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URIFromOID(oidRef)},
			}),
			toResolveRefs: []string{uuidRef.GetValue()},
			wantResources: []fhir.Resource{},
		},
		{
			name: "Resolve URN - UUID Reference",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URIFromUUID(uuidRef)},
			}),
			toResolveRefs: []string{uuidRef.GetValue()},
			wantResources: []fhir.Resource{patient123},
		},
		{
			name: "Resolve URN - Invalid UUID Reference; returns empty",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URIFromUUID(invalidUuidRef)},
			}),
			toResolveRefs: []string{invalidUuidRef.GetValue()},
			wantResources: []fhir.Resource{},
		},
		{
			name: "Resolve URN - OID Reference",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URIFromOID(oidRef)},
			}),
			toResolveRefs: []string{oidRef.GetValue()},
			wantResources: []fhir.Resource{patient123},
		},
		{
			name: "Resolve Absolute URL - Versionless - 1 match",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URI(patAbsVersionlessUrl)},
			}),
			toResolveRefs: []string{patAbsVersionlessUrl},
			wantResources: []fhir.Resource{patient123},
		},
		{
			name: "Resolve Absolute URL - Versionless - multiple matches",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URI(patAbsVersionlessUrl)},
				{Resource: containedresource.Wrap(patient123Latest),
					FullUrl: fhir.URI(patAbsVersionlessUrl)},
			}),
			toResolveRefs: []string{patAbsVersionlessUrl},
			wantResources: []fhir.Resource{patient123Latest},
		},
		{
			name: "Resolve Absolute URL - Versionless - multiple matches without Meta; throws",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URI(patAbsVersionlessUrl)},
				{Resource: crPatientMissingMeta,
					FullUrl: fhir.URI(patAbsVersionlessUrl)},
			}),
			toResolveRefs: []string{patAbsVersionlessUrl},
			wantErr:       resolver.ErrMissingMetaOrLastUpdated,
		},
		{
			name: "Resolve Absolute URL - Versionless - multiple matches without Meta.LastUpdated; throws",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URI(patAbsVersionlessUrl)},
				{Resource: crPatientMissingLastUpdated,
					FullUrl: fhir.URI(patAbsVersionlessUrl)},
			}),
			toResolveRefs: []string{patAbsVersionlessUrl},
			wantErr:       resolver.ErrMissingMetaOrLastUpdated,
		},
		{
			name: "Resolve Absolute URL - Versionless - single match without Meta; resolves",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatientMissingMeta,
					FullUrl: fhir.URI(patAbsVersionlessUrl)},
			}),
			toResolveRefs: []string{patAbsVersionlessUrl},
			wantResources: []fhir.Resource{patientMissingMeta},
		},
		{
			name: "Resolve Absolute URL - Versionless - single match without Meta.LastUpdated; resolves",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatientMissingLastUpdated,
					FullUrl: fhir.URI(patAbsVersionlessUrl)},
			}),
			toResolveRefs: []string{patAbsVersionlessUrl},
			wantResources: []fhir.Resource{patientMissingLastUpdated},
		},
		{
			name: "Resolve Absolute URL - Versionless - no match found",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: containedresource.Wrap(obs123),
					FullUrl: fhir.URI(obsAbsVersionlessUrl)},
			}),
			toResolveRefs: []string{patAbsVersionlessUrl},
			wantResources: []fhir.Resource{},
		},
		{
			name: "Resolve Absolute URL - Invalid Versionless URL; returns empty",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URI(invalidAbsVersionlessUrl)},
			}),
			toResolveRefs: []string{invalidAbsVersionlessUrl},
			wantResources: []fhir.Resource{},
		},
		{
			name: "Resolve Absolute URL - Versioned - 1 match",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: containedresource.Wrap(patient123VersionedV1),
					FullUrl: fhir.URI(patAbsVersionedUrl)},
			}),
			toResolveRefs: []string{patAbsVersionedUrl},
			wantResources: []fhir.Resource{patient123VersionedV1},
		},
		{
			name: "Resolve Absolute URL - Versioned - multiple matches with same ID and version - throws",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: containedresource.Wrap(patient123VersionedV1),
					FullUrl: fhir.URI(patAbsVersionedUrl)},
				{Resource: containedresource.Wrap(patient123VersionedV1Copy),
					FullUrl: fhir.URI(patAbsVersionedUrl)},
			}),
			toResolveRefs: []string{patAbsVersionedUrl},
			wantErr:       resolver.ErrMultipleResourcesWithSameIDAndVersion,
		},
		{
			name: "Resolve Absolute URL - Versioned - no match found",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: containedresource.Wrap(obs123),
					FullUrl: fhir.URI(obsAbsVersionedUrl)},
			}),
			toResolveRefs: []string{patAbsVersionedUrl},
			wantResources: []fhir.Resource{},
		},
		{
			name: "Resolve Absolute URL - Versioned - invalid versioned URL; returns empty",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: containedresource.Wrap(patient123VersionedV1),
					FullUrl: fhir.URI(invalidVersionedAbsUrl)},
			}),
			toResolveRefs: []string{invalidVersionedAbsUrl},
			wantResources: []fhir.Resource{},
		},
		{
			name: "Resolve Relative URL Reference - Versionless - match found",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URI(patAbsVersionlessUrl)},
			}),
			toResolveRefs: []string{"Patient/123"},
			wantResources: []fhir.Resource{patient123},
		},
		{
			name: "Resolve Relative URL Reference - Versionless - no match found",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URI(patAbsVersionlessUrl)},
			}),
			toResolveRefs: []string{"Observation/123"},
			wantResources: []fhir.Resource{},
		},
		{
			name: "Resolve Relative URL Reference - Invalid Versionless URL - returns empty",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URI(invalidAbsVersionlessUrl)},
			}),
			toResolveRefs: []string{"Whale/123"},
			wantResources: []fhir.Resource{},
		},
		{
			name: "Resolve Relative URL Reference - invalid relative reference; returns empty",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URI(patAbsVersionlessUrl)},
			}),
			toResolveRefs: []string{"Shark/Patient/Whale/123"},
			wantResources: []fhir.Resource{},
		},
		{
			name: "Resolve Relative URL Reference - empty rootURLs in resolveResources; returns empty",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URIFromOID(oidRef)},
			}),
			toResolveRefs: []string{"Patient/123"},
			wantResources: []fhir.Resource{},
		},
		{
			name: "Resolve Relative URL Reference - Versioned - match found",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: containedresource.Wrap(patient123VersionedV1),
					FullUrl: fhir.URI(patAbsVersionedUrl)},
			}),
			toResolveRefs: []string{"Patient/123/_history/v1"},
			wantResources: []fhir.Resource{patient123VersionedV1},
		},
		{
			name: "Resolve Relative URL Reference - Versioned - no match found",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: containedresource.Wrap(patient123VersionedV1),
					FullUrl: fhir.URI(patAbsVersionedUrl)},
			}),
			toResolveRefs: []string{"Observation/123/_history/v1"},
			wantResources: []fhir.Resource{},
		},
		{
			name: "Resolve Relative URL Reference - Invalid Versioned; returns empty",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: containedresource.Wrap(patient123VersionedV1),
					FullUrl: fhir.URI(invalidVersionedAbsUrl)},
			}),
			toResolveRefs: []string{"Patient/123/_history/%#@^%$#"},
			wantResources: []fhir.Resource{},
		},
		{
			name: "Repeated bundle entries, resolves to a single resource for URNs",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URIFromOID(oidRef)},
				{Resource: crPatient123,
					FullUrl: fhir.URIFromOID(oidRef)},
			}),
			toResolveRefs: []string{oidRef.GetValue()},
			wantResources: []fhir.Resource{patient123},
		},
		{
			name: "Repeated toResolve reference inputs, resolves to one resource match for each reference, if found",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URIFromOID(oidRef)},
			}),
			toResolveRefs: []string{oidRef.GetValue(), oidRef.GetValue()},
			wantResources: []fhir.Resource{patient123, patient123},
		},
		{
			name: "Bundle entries with same ID for different resource type returns the right match, if found",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URI(patAbsVersionlessUrl)},
				{Resource: containedresource.Wrap(obs123),
					FullUrl: fhir.URI(obsAbsVersionlessUrl)},
			}),
			toResolveRefs: []string{patAbsVersionlessUrl},
			wantResources: []fhir.Resource{patient123},
		},
		{
			name: "Nil entry in bundle; returns empty",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: nil,
					FullUrl: fhir.URIFromOID(oidRef)},
			}),
			toResolveRefs: []string{oidRef.GetValue()},
			wantResources: []fhir.Resource{},
		},
		{
			name: "Nil resource in bundle entry; returns empty",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: containedresource.Wrap(nil),
					FullUrl: fhir.URIFromOID(oidRef)},
			}),
			toResolveRefs: []string{oidRef.GetValue()},
			wantResources: []fhir.Resource{},
		},
		{
			name:          "Empty bundle entries; returns empty",
			bundle:        createBundle([]*r4pb.Bundle_Entry{}),
			toResolveRefs: []string{oidRef.GetValue()},
			wantResources: []fhir.Resource{},
		},
		{
			name: "No bundle entries; returns empty",
			bundle: &r4pb.Bundle{
				Type: &r4pb.Bundle_TypeCode{Value: cpb.BundleTypeCode_SEARCHSET},
			},
			toResolveRefs: []string{oidRef.GetValue()},
			wantResources: []fhir.Resource{},
		},
		{
			name: "Nil FullUrl in bundle entry; returns empty",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: nil},
			}),
			toResolveRefs: []string{patAbsVersionlessUrl},
			wantResources: []fhir.Resource{},
		},
		{
			name: "Empty FullUrl in bundle entry; returns empty",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URI("")},
			}),
			toResolveRefs: []string{patAbsVersionlessUrl},
			wantResources: []fhir.Resource{},
		},
		{
			name: "Empty Resolve Reference; returns empty",
			bundle: createBundle([]*r4pb.Bundle_Entry{
				{Resource: crPatient123,
					FullUrl: fhir.URIFromOID(oidRef)},
			}),
			toResolveRefs: []string{""},
			wantResources: []fhir.Resource{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			br, err := resolver.NewBundleResolver(tc.bundle)
			if err != nil {
				t.Fatalf("NewBundleResolver() error = %v, expected none", err)
			}
			gotResources, gotErr := br.Resolve(tc.toResolveRefs)

			if !cmp.Equal(gotErr, tc.wantErr, cmpopts.EquateErrors()) {
				t.Errorf("Resolve() error = %v, wantErr %v", err, tc.wantErr)
			}

			if diff := cmp.Diff(tc.wantResources, gotResources, protocmp.Transform()); diff != "" {
				t.Errorf("Resolve(%s): mismatch (-want, +got):\n%s", tc.name, diff)
			}
		})
	}
}

func TestNewBundleResolver(t *testing.T) {
	testCases := []struct {
		name    string
		bundle  *r4pb.Bundle
		wantErr error
	}{
		{
			name: "Supported Bundle Type - COLLECTION",
			bundle: &r4pb.Bundle{
				Type: &r4pb.Bundle_TypeCode{Value: cpb.BundleTypeCode_COLLECTION},
			},
		},
		{
			name: "Supported Bundle Type - SEARCHSET",
			bundle: &r4pb.Bundle{
				Type: &r4pb.Bundle_TypeCode{Value: cpb.BundleTypeCode_SEARCHSET},
			},
		},
		{
			name: "Unsupported Bundle Type - HISTORY",
			bundle: &r4pb.Bundle{
				Type: &r4pb.Bundle_TypeCode{Value: cpb.BundleTypeCode_HISTORY},
			},
			wantErr: resolver.ErrUnsupportedBundleType,
		},
		{
			name:    "Nil Bundle - throws",
			wantErr: resolver.ErrNilBundleInit,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			br, gotErr := resolver.NewBundleResolver(tc.bundle)
			if gotErr != tc.wantErr {
				t.Errorf("NewBundleResolver() error = %v, wantErr %v", gotErr, tc.wantErr)
			}
			if br == nil && tc.wantErr == nil {
				t.Errorf("NewBundleResolver() got nil resolver, want non-nil")
			}
		})
	}
}
