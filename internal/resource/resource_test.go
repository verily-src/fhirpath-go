package resource_test

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	dpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/device_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/document_reference_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/questionnaire_response_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirtest"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestVersionETag(t *testing.T) {
	testCases := []struct {
		name string
		res  fhir.Resource
		want string
	}{
		{
			"nil resource",
			nil,
			"",
		},
		{
			"failure: VersionId is the empty string",
			&ppb.Patient{Meta: &dtpb.Meta{VersionId: fhir.ID("")}},
			"",
		},
		{
			"extracted the VersionId from a Patient",
			&ppb.Patient{Meta: &dtpb.Meta{VersionId: fhir.ID("abc")}},
			`W/"abc"`,
		},
		{
			"extracted the VersionId from a Device",
			&dpb.Device{Meta: &dtpb.Meta{VersionId: fhir.ID("xyz")}},
			`W/"xyz"`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := resource.VersionETag(tc.res)

			if got != tc.want {
				t.Errorf("VersionETag(%s) version got = %v, want = %v", tc.name, got, tc.want)
			}
		})
	}
}

func TestVersionedURI(t *testing.T) {
	testCases := []struct {
		name string
		res  fhir.Resource
		want *dtpb.Uri
	}{
		{
			"nil resource",
			nil,
			nil,
		},
		{
			"nil resource",
			&ppb.Patient{Meta: &dtpb.Meta{VersionId: fhir.ID("")}},
			nil,
		},
		{
			"no version",
			&ppb.Patient{Id: fhir.ID("abc")},
			nil,
		},
		{
			"versioned resource",
			&dpb.Device{Id: fhir.ID("123"), Meta: &dtpb.Meta{VersionId: fhir.ID("abc")}},
			&dtpb.Uri{Value: "Device/123/_history/abc"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := resource.VersionedURI(tc.res)

			if diff := cmp.Diff(got, tc.want, protocmp.Transform()); diff != "" {
				t.Fatalf("VersionedURI(%s): (-got, +want):\n%s", tc.name, diff)
			}
		})
	}
}

func TestVersionedURIString(t *testing.T) {
	testCases := []struct {
		name      string
		res       fhir.Resource
		want      string
		wantFound bool
	}{
		{
			"nil resource",
			&ppb.Patient{Meta: &dtpb.Meta{VersionId: fhir.ID("")}},
			"",
			false,
		},
		{
			"no version",
			&ppb.Patient{Id: fhir.ID("abc")},
			"",
			false,
		},
		{
			"versioned resource",
			&dpb.Device{Id: fhir.ID("123"), Meta: &dtpb.Meta{VersionId: fhir.ID("abc")}},
			"Device/123/_history/abc",
			true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, found := resource.VersionedURIString(tc.res)

			if found != tc.wantFound {
				t.Fatalf("VersionedURIString(%s) found got = %v, want = %v", tc.name, got, tc.wantFound)
			}
			if got != tc.want {
				t.Errorf("VersionedURIString(%s) got = %v, want = %v", tc.name, got, tc.want)
			}
		})
	}
}

func TestProfileStrings(t *testing.T) {
	testCases := []struct {
		name string
		res  fhir.Resource
		want []string
	}{
		{
			"nil resource",
			nil,
			nil,
		},
		{
			"nil resource",
			&ppb.Patient{Meta: &dtpb.Meta{VersionId: fhir.ID("")}},
			nil,
		},
		{
			"no profiles",
			&ppb.Patient{Id: fhir.ID("abc")},
			nil,
		},
		{
			"empty profiles",
			&dpb.Device{Id: fhir.ID("123"), Meta: &dtpb.Meta{Profile: []*dtpb.Canonical{}}},
			nil,
		},
		{
			"single profiles",
			&dpb.Device{Id: fhir.ID("123"), Meta: &dtpb.Meta{Profile: []*dtpb.Canonical{{Value: "profileA"}}}},
			[]string{"profileA"},
		},
		{
			"multiple profiles",
			&dpb.Device{Id: fhir.ID("123"), Meta: &dtpb.Meta{Profile: []*dtpb.Canonical{{Value: "profileA"}, {Value: "profileB"}}}},
			[]string{"profileA", "profileB"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := resource.ProfileStrings(tc.res)

			if diff := cmp.Diff(got, tc.want, protocmp.Transform()); diff != "" {
				t.Fatalf("ProfileStrings(%s): (-got, +want):\n%s", tc.name, diff)
			}
		})
	}
}

func TestUnversionedProfileStrings(t *testing.T) {
	testCases := []struct {
		name string
		res  fhir.Resource
		want []string
	}{
		{
			name: "nil resource has no profiles",
			res:  nil,
			want: nil,
		},
		{
			name: "empty resource has no profiles",
			res:  &ppb.Patient{},
			want: nil,
		},
		{
			name: "resource with no profile has no profiles",
			res:  &ppb.Patient{Meta: &dtpb.Meta{Profile: []*dtpb.Canonical{}}},
			want: nil,
		},
		{
			name: "resource with single unversioned profile returns profile",
			res:  &ppb.Patient{Meta: &dtpb.Meta{Profile: []*dtpb.Canonical{{Value: "Profile1"}}}},
			want: []string{"Profile1"},
		},
		{
			name: "resource with single versioned profile returns unversioned profile",
			res:  &ppb.Patient{Meta: &dtpb.Meta{Profile: []*dtpb.Canonical{{Value: "Profile1|1.2.3"}}}},
			want: []string{"Profile1"},
		},
		{
			name: "resource with multiple unversioned profiles returns profiles",
			res: &ppb.Patient{Meta: &dtpb.Meta{Profile: []*dtpb.Canonical{
				{Value: "Profile1"}, {Value: "Profile2"}, {Value: "Profile3"}, {Value: "Profile4"}, {Value: "Profile5"},
			}}},
			want: []string{"Profile1", "Profile2", "Profile3", "Profile4", "Profile5"},
		},
		{
			name: "resource with multiple versioned profiles returns unversioned profiles",
			res: &ppb.Patient{Meta: &dtpb.Meta{Profile: []*dtpb.Canonical{
				{Value: "Profile1|1.2.3"}, {Value: "Profile2|4.5.6"}, {Value: "Profile3|7.8.9"},
			}}},
			want: []string{"Profile1", "Profile2", "Profile3"},
		},
		{
			name: "resource with duplicate profiles returns unique profiles",
			res: &ppb.Patient{Meta: &dtpb.Meta{Profile: []*dtpb.Canonical{
				{Value: "Profile1"}, {Value: "Profile1"}, {Value: "Profile2"}, {Value: "Profile2|1.2.3"}, {Value: "Profile2|4.5.6"}, {Value: "Profile2|7.8.9"},
			}}},
			want: []string{"Profile1", "Profile2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := resource.UnversionedProfileStrings(tc.res)

			if diff := cmp.Diff(got, tc.want, protocmp.Transform()); diff != "" {
				t.Fatalf("UnversionedProfileStrings(%s): (-got, +want):\n%s", tc.name, diff)
			}
		})
	}
}

func TestProfiles(t *testing.T) {
	testCases := []struct {
		name string
		res  fhir.Resource
		want []*dtpb.Canonical
	}{
		{
			"nil resource",
			nil,
			nil,
		},
		{
			"nil resource",
			&ppb.Patient{Meta: &dtpb.Meta{VersionId: fhir.ID("")}},
			nil,
		},
		{
			"no profiles",
			&ppb.Patient{Id: fhir.ID("abc")},
			nil,
		},
		{
			"empty profiles",
			&dpb.Device{Id: fhir.ID("123"), Meta: &dtpb.Meta{Profile: []*dtpb.Canonical{}}},
			[]*dtpb.Canonical{},
		},
		{
			"single profiles",
			&dpb.Device{Id: fhir.ID("123"), Meta: &dtpb.Meta{Profile: []*dtpb.Canonical{{Value: "profileA"}}}},
			[]*dtpb.Canonical{{Value: "profileA"}},
		},
		{
			"multiple profiles",
			&dpb.Device{Id: fhir.ID("123"), Meta: &dtpb.Meta{Profile: []*dtpb.Canonical{{Value: "profileA"}, {Value: "profileB"}}}},
			[]*dtpb.Canonical{{Value: "profileA"}, {Value: "profileB"}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := resource.Profiles(tc.res)

			if diff := cmp.Diff(got, tc.want, protocmp.Transform()); diff != "" {
				t.Fatalf("ProfileStrings(%s): (-got, +want):\n%s", tc.name, diff)
			}
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	patient := fhirtest.NewResource(t, resource.Patient)
	device := fhirtest.NewResource(t, resource.Device)
	account := fhirtest.NewResource(t, resource.Account)

	testCases := []struct {
		name  string
		input []fhir.Resource
		want  []fhir.Resource
	}{
		{
			name:  "Inputs are unique",
			input: []fhir.Resource{patient, device, account},
			want:  []fhir.Resource{patient, device, account},
		},
		{
			name:  "Duplicates are removed",
			input: []fhir.Resource{patient, device, account, device, account, patient},
			want:  []fhir.Resource{patient, device, account},
		},
		{
			name:  "Removes nil",
			input: []fhir.Resource{nil, patient, nil},
			want:  []fhir.Resource{patient},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := resource.RemoveDuplicates(tc.input)

			opts := []cmp.Option{
				cmpopts.SortSlices(func(lhs, rhs fhir.Resource) bool {
					return resource.URIString(lhs) < resource.URIString(rhs)
				}),
				protocmp.Transform(),
			}
			if want := tc.want; !cmp.Equal(got, want, opts...) {
				t.Errorf("RemoveDuplicates(%v): got '%v', want '%v'", tc.name, got, want)
			}
		})
	}
}

func TestGroupResources(t *testing.T) {
	patient := fhirtest.NewResource(t, resource.Patient)
	device := fhirtest.NewResource(t, resource.Device)
	account := fhirtest.NewResource(t, resource.Account)

	testCases := []struct {
		name  string
		input []fhir.Resource
		want  map[resource.Type][]fhir.Resource
	}{
		{
			name:  "Inputs are sorted",
			input: []fhir.Resource{patient, device, account},
			want: map[resource.Type][]fhir.Resource{
				resource.Patient: {patient},
				resource.Device:  {device},
				resource.Account: {account},
			},
		},
		{
			name:  "Removes nil",
			input: []fhir.Resource{nil, patient, nil},
			want: map[resource.Type][]fhir.Resource{
				resource.Patient: {patient},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := resource.GroupResources(tc.input)

			opts := []cmp.Option{
				protocmp.Transform(),
				cmpopts.SortMaps(func(lhs, rhs resource.Type) bool {
					return lhs < rhs
				}),
			}
			if want := tc.want; !cmp.Equal(got, want, opts...) {
				t.Errorf("GroupResources(%v): got '%v', want '%v'", tc.name, got, want)
			}
		})
	}
}

// Assert that an identifier list is present and has expected values
func assertIdentifier(t *testing.T, identifiers []*dtpb.Identifier, name string, resource fhir.Resource) {
	if identifiers == nil {
		t.Errorf("Nil list of ids for %v: %v", name, resource)
		return
	}
	if len(identifiers) == 0 {
		t.Errorf("Empty list of ids for %v: %v", name, resource)
		return
	}

	id := identifiers[0]

	if got, want := id.GetSystem().Value, "http://example.com/fake-id"; got != want {
		t.Errorf("%v Resource.Identifier[0].System: got %v, want %v", name, got, want)
	}

	value := id.GetValue().Value
	matched, _ := regexp.MatchString("^[a-f0-9-]+$", value)
	if !matched {
		t.Errorf("%v Resource.Identifier[0].Value: got %v, expected uuid", name, value)
	}
}

func TestGetIdentifier(t *testing.T) {
	// Test that all resources return nil by default
	for name := range fhirtest.Resources {
		t.Run("Resources/"+name, func(t *testing.T) {
			res := fhirtest.NewResource(
				t,
				resource.Type(name),
			)

			ids, _ := resource.GetIdentifierList(res)

			if ids != nil {
				t.Errorf("%v Resource.Identifier: got %v, want nil -- not supposed to have Identifier", name, ids)
			}
		})
	}

	// Test that all types compatible with CanonicalResource actually return an identifier
	for name := range fhirtest.CanonicalResources {
		t.Run("CanonicalResources/"+name, func(t *testing.T) {
			res := fhirtest.NewResource(
				t,
				resource.Type(name),
				fhirtest.WithGeneratedIdentifier("http://example.com/fake-id"),
			)
			ids, err := resource.GetIdentifierList(res)
			if err != nil {
				t.Errorf("got %v, want nil -- unexpected error", err)
			}

			assertIdentifier(t, ids, name, res)
		})
	}

	// Sanity check a few specific types that we know have Identifier
	for _, name := range []string{"Patient", "DocumentReference", "AdverseEvent", "Bundle"} {
		t.Run("CanonicalResources/"+name, func(t *testing.T) {
			res := fhirtest.NewResource(
				t,
				resource.Type(name),
				fhirtest.WithGeneratedIdentifier("http://example.com/fake-id"),
			)

			ids, err := resource.GetIdentifierList(res)
			if err != nil {
				t.Errorf("got %v, want nil -- unexpected error", err)
			}

			assertIdentifier(t, ids, name, res)
		})
	}

}

func TestGetIdentifier_single(t *testing.T) {
	// Sanity check a few specific types that have a singleton Identifier
	// Bundle, QuestionnaireResponse

	testIds := []*dtpb.Identifier{
		{
			System: &dtpb.Uri{Value: "http://example.com/fake-id"},
			Value:  &dtpb.String{Value: "35c423fc-0651-4c83-b63f-9008e0c96445"},
		},
		{
			System: &dtpb.Uri{Value: "http://example.com/fake-id"},
			Value:  &dtpb.String{Value: "ddec1b6e-4539-4aae-becf-b4dced32189f"},
		},
	}

	testCases := []struct {
		name string
		res  fhir.Resource
		want *dtpb.Identifier
	}{
		{
			"Bundle",
			&bundle_and_contained_resource_go_proto.Bundle{
				Identifier: testIds[0],
			},
			testIds[0],
		},
		{
			"QuestionnaireResponse",
			&questionnaire_response_go_proto.QuestionnaireResponse{
				Identifier: testIds[1],
			},
			testIds[1],
		},
	}
	for _, tc := range testCases {

		ids, err := resource.GetIdentifierList(tc.res)
		want := []*dtpb.Identifier{tc.want}

		if err != nil {
			t.Errorf("got %v, want nil", err)
			return
		}

		assertIdentifier(t, ids, tc.name, tc.res)

		if len(ids) != len(want) {
			t.Errorf("got %v, want %v", ids, want)
			continue
		}
		if ids[0] != want[0] {
			t.Errorf("got %v, want %v", ids[0], want[0])
		}
	}
}

func TestGetIdentifier_nil(t *testing.T) {
	// Sanity check a few specific types that we know do NOT have Identifier
	resourcesWithoutIdentifiers := []string{
		"Provenance",
		"Linkage",
	}
	for _, name := range resourcesWithoutIdentifiers {
		t.Run("CanonicalResources/"+name, func(t *testing.T) {
			res := fhirtest.NewResource(t, resource.Type(name))

			ids, err := resource.GetIdentifierList(res)

			if ids != nil {
				t.Errorf("%v Resource.Identifier: got %v, want nil -- not supposed to have Identifier", name, ids)
			}

			if err == nil {
				t.Errorf("got nil, want error")
			}

			wanterr := "Resource does not implement GetIdentifier()"
			if !strings.Contains(err.Error(), wanterr) {
				t.Errorf("got %#v, want %#v", err.Error(), wanterr)
			}
			if !errors.Is(err, resource.ErrGetIdentifierList) {
				t.Errorf("got error %#v, want errors.Is(..., ErrGenerateIfNoneExist)", err)
			}
		})
	}
}

// Sanity check that a few resources have GetIdentifier() as a list. This is not a complete list.
var _ resource.HasGetIdentifierList = (*ppb.Patient)(nil)
var _ resource.HasGetIdentifierList = (*document_reference_go_proto.DocumentReference)(nil)

// Sanity check that a few resources have GetIdentifier() as a single ID. This is not a complete list.
var _ resource.HasGetIdentifierSingle = (*bundle_and_contained_resource_go_proto.Bundle)(nil)
var _ resource.HasGetIdentifierSingle = (*questionnaire_response_go_proto.QuestionnaireResponse)(nil)

func TestGetExtension_List(t *testing.T) {
	ext1 := &dtpb.Extension{Url: &dtpb.Uri{Value: "http://example.com/ext1"}}
	ext2 := &dtpb.Extension{Url: &dtpb.Uri{Value: "http://example.com/ext2"}}

	testCases := []struct {
		name string
		res  fhir.Resource
		ok   bool
		want []*dtpb.Extension
	}{
		{
			name: "Patient with multiple extensions",
			res: &ppb.Patient{
				Extension: []*dtpb.Extension{ext1, ext2},
			},
			ok:   true,
			want: []*dtpb.Extension{ext1, ext2},
		},
		{
			name: "Patient with no extensions",
			res:  &ppb.Patient{},
			want: nil,
			ok:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := resource.GetExtensions(tc.res)
			if ok != tc.ok {
				t.Errorf("GetExtensions(%s) got ok = %v, want %v", tc.name, ok, tc.ok)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("GetExtensions(%s) mismatch (-want, +got):\n%s", tc.name, diff)
			}
		})
	}
}

func TestGetExtension_Single(t *testing.T) {
	ext := &dtpb.Extension{Url: &dtpb.Uri{Value: "http://example.com/ext-single"}}

	testCases := []struct {
		name string
		res  fhir.Resource
		ok   bool
		want []*dtpb.Extension
	}{
		{
			name: "QuestionnaireResponse with single extension",
			res: &questionnaire_response_go_proto.QuestionnaireResponse{
				Extension: []*dtpb.Extension{ext},
			},
			ok:   true,
			want: []*dtpb.Extension{ext},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := resource.GetExtensions(tc.res)
			if ok != tc.ok {
				t.Errorf("GetExtensions(%s) got ok = %v, want %v", tc.name, ok, tc.ok)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("GetExtensions(%s) mismatch (-want, +got):\n%s", tc.name, diff)
			}
		})
	}
}

func TestGetExtension_nil(t *testing.T) {
	testCases := []struct {
		name string
		res  fhir.Resource
		ok   bool
	}{
		{
			name: "Resource does not implement GetExtension",
			res:  &bundle_and_contained_resource_go_proto.Bundle{},
			ok:   false,
		},
		{
			name: "Nil resource",
			res:  nil,
			ok:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := resource.GetExtensions(tc.res)
			if ok != tc.ok {
				t.Errorf("GetExtensions(%s): got ok = %v, want %v", tc.name, ok, tc.ok)
			}
			if got != nil {
				t.Errorf("GetExtensions(%s): got %v, want nil", tc.name, got)
			}
		})
	}
}
