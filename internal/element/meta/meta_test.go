package meta_test

import (
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"

	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/element/canonical"
	"github.com/verily-src/fhirpath-go/internal/element/extension"
	"github.com/verily-src/fhirpath-go/internal/element/meta"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestWithTags(t *testing.T) {
	at := fhir.Coding("at-s", "at-v")
	bt := fhir.Coding("b-s", "b-v")
	testCases := []struct {
		name     string
		inMeta   *dtpb.Meta
		inTags   []*dtpb.Coding
		wantTags []*dtpb.Coding
	}{
		{
			name: "Clears Tags",
			inMeta: &dtpb.Meta{
				Tag: []*dtpb.Coding{at},
			},
			inTags:   nil,
			wantTags: nil,
		},
		{
			name: "Replaces Tags",
			inMeta: &dtpb.Meta{
				Tag: []*dtpb.Coding{at},
			},
			inTags:   []*dtpb.Coding{bt},
			wantTags: []*dtpb.Coding{bt},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			inMeta := proto.Clone(testCase.inMeta).(*dtpb.Meta)

			meta.Update(inMeta, meta.WithTags(testCase.inTags...))

			if diff := cmp.Diff(inMeta.GetTag(), testCase.wantTags, protocmp.Transform()); diff != "" {
				t.Errorf("WithTags(): (-want, +got)\n%v", diff)
			}
		})
	}
}

func TestIncludesTags(t *testing.T) {
	at := fhir.Coding("at-s", "at-v")
	bt := fhir.Coding("b-s", "b-v")
	testCases := []struct {
		name     string
		inMeta   *dtpb.Meta
		inTags   []*dtpb.Coding
		wantTags []*dtpb.Coding
	}{
		{
			name: "Maintains Tags",
			inMeta: &dtpb.Meta{
				Tag: []*dtpb.Coding{at},
			},
			inTags:   nil,
			wantTags: []*dtpb.Coding{at},
		},
		{
			name: "Appends Tags",
			inMeta: &dtpb.Meta{
				Tag: []*dtpb.Coding{at},
			},
			inTags:   []*dtpb.Coding{bt},
			wantTags: []*dtpb.Coding{at, bt},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			inMeta := proto.Clone(testCase.inMeta).(*dtpb.Meta)

			meta.Update(inMeta, meta.IncludeTags(testCase.inTags...))

			if diff := cmp.Diff(inMeta.GetTag(), testCase.wantTags, protocmp.Transform()); diff != "" {
				t.Errorf("IncludeTags(): (-want, +got)\n%v", diff)
			}
		})
	}
}

func TestWithMetaProfiles(t *testing.T) {
	ap := canonical.New("ap")
	bp := canonical.New("bp")
	testCases := []struct {
		name         string
		inMeta       *dtpb.Meta
		inProfiles   []*dtpb.Canonical
		wantProfiles []*dtpb.Canonical
	}{
		{
			name: "Clears Profiles",
			inMeta: &dtpb.Meta{
				Profile: []*dtpb.Canonical{ap},
			},
			inProfiles:   nil,
			wantProfiles: nil,
		},
		{
			name: "Replaces Profiles",
			inMeta: &dtpb.Meta{
				Profile: []*dtpb.Canonical{ap},
			},
			inProfiles:   []*dtpb.Canonical{bp},
			wantProfiles: []*dtpb.Canonical{bp},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			inMeta := proto.Clone(testCase.inMeta).(*dtpb.Meta)

			meta.Update(inMeta, meta.WithProfiles(testCase.inProfiles...))

			if diff := cmp.Diff(inMeta.GetProfile(), testCase.wantProfiles, protocmp.Transform()); diff != "" {
				t.Errorf("WithProfiles(): (-want, +got)\n%v", diff)
			}
		})
	}
}

func TestIncludeMetaProfiles(t *testing.T) {
	ap := canonical.New("ap")
	bp := canonical.New("bp")
	testCases := []struct {
		name         string
		inMeta       *dtpb.Meta
		inProfiles   []*dtpb.Canonical
		wantProfiles []*dtpb.Canonical
	}{
		{
			name: "Maintains Profiles",
			inMeta: &dtpb.Meta{
				Profile: []*dtpb.Canonical{ap},
			},
			inProfiles:   nil,
			wantProfiles: []*dtpb.Canonical{ap},
		},
		{
			name: "Appends Profiles",
			inMeta: &dtpb.Meta{
				Profile: []*dtpb.Canonical{ap},
			},
			inProfiles:   []*dtpb.Canonical{bp},
			wantProfiles: []*dtpb.Canonical{ap, bp},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			inMeta := proto.Clone(testCase.inMeta).(*dtpb.Meta)

			meta.Update(inMeta, meta.IncludeProfiles(testCase.inProfiles...))

			if diff := cmp.Diff(inMeta.GetProfile(), testCase.wantProfiles, protocmp.Transform()); diff != "" {
				t.Errorf("IncludeProfiles(): (-want, +got)\n%v", diff)
			}
		})
	}
}

func TestReplaceMeta(t *testing.T) {
	patient := &ppb.Patient{Meta: &dtpb.Meta{}}
	wantMeta := &dtpb.Meta{VersionId: fhir.ID("apple")}

	t.Run("ReplaceMeta", func(t *testing.T) {
		meta.ReplaceInResource(patient, &dtpb.Meta{VersionId: fhir.ID("apple")})
		if diff := cmp.Diff(patient.GetMeta(), wantMeta, protocmp.Transform()); diff != "" {
			t.Errorf("ReplaceMeta(): (-want, +got)\n%v", diff)
		}
	})
}

func TestEnsureMeta(t *testing.T) {
	patient := &ppb.Patient{}
	wantMeta := &dtpb.Meta{}

	t.Run("EnsureMeta", func(t *testing.T) {
		meta.EnsureInResource(patient)
		if diff := cmp.Diff(patient.GetMeta(), wantMeta, protocmp.Transform()); diff != "" {
			t.Errorf("EnsureMeta(): (-want, +got)\n%v", diff)
		}
	})
}

func TestWithExtension(t *testing.T) {
	oldExtension := extension.New("extension-url-old", fhir.String("extension-old"))
	newExtension := extension.New("extension-url-new", fhir.String("extension-new"))
	testCases := []struct {
		name          string
		inMeta        *dtpb.Meta
		inExtension   []*dtpb.Extension
		wantExtension []*dtpb.Extension
	}{
		{
			name: "Clears Extension",
			inMeta: &dtpb.Meta{
				Extension: []*dtpb.Extension{oldExtension},
			},
			inExtension:   nil,
			wantExtension: nil,
		},
		{
			name: "Replaces Extension",
			inMeta: &dtpb.Meta{
				Extension: []*dtpb.Extension{oldExtension},
			},
			inExtension:   []*dtpb.Extension{newExtension},
			wantExtension: []*dtpb.Extension{newExtension},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			inMeta := proto.Clone(testCase.inMeta).(*dtpb.Meta)
			meta.Update(inMeta, meta.WithExtensions(testCase.inExtension...))
			if diff := cmp.Diff(inMeta.GetExtension(), testCase.wantExtension, protocmp.Transform()); diff != "" {
				t.Errorf("WithExtensions(): (-want, +got)\n%v", diff)
			}
		})
	}
}
