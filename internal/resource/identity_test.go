package resource_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/internal/fhirtest"
	"github.com/verily-src/fhirpath-go/internal/resource"
)

func TestIdentityOf_NilInputs_ReturnsNoValue(t *testing.T) {
	_, got := resource.IdentityOf(nil)

	if got, want := got, false; got != want {
		t.Errorf("IdentityOf: got %v, want %v", got, want)
	}
}

func TestIdentityOf(t *testing.T) {
	for name, res := range fhirtest.Resources {
		t.Run(name, func(t *testing.T) {
			want, err := resource.NewIdentity(
				string(resource.TypeOf(res)),
				res.GetId().GetValue(),
				res.GetMeta().GetVersionId().GetValue(),
			)
			if err != nil {
				t.Fatalf("IdentityOf(%v): got unexpected err: %v", name, err)
			}

			got, ok := resource.IdentityOf(res)
			if !ok {
				t.Fatalf("IdentityOf(%v): got false for ok", name)
			}

			if !cmp.Equal(got, want) {
				t.Errorf("IdentityOf(%v): got %v, want %v", name, got, want)
			}
		})
	}
}

func TestNewIdentity_BadInput_ReturnsErrBadType(t *testing.T) {
	_, err := resource.NewIdentity("", "1234", "5678")

	if got, want := err, resource.ErrBadType; !errors.Is(got, want) {
		t.Errorf("NewIdentity: got err '%v', want err '%v'", got, want)
	}
}

func TestNewIdentityFromURL(t *testing.T) {
	testCases := []struct {
		name         string
		URL          string
		wantIdentity *resource.Identity
	}{
		{
			"URL",
			"https://healthcare.googleapis.com/v1/projects/my-project-name/locations/us-east4/datasets/my-dataset-name/fhirStores/my-fhir-store-name/fhir/Binary/123",
			mustNewIdentity("Binary", "123", ""),
		},
		{
			"URILong",
			"projects/my-project-name/locations/us-east4/datasets/my-dataset-name/fhirStores/my-fhir-store-name/fhir/Patient/abc",
			mustNewIdentity("Patient", "abc", ""),
		},
		{
			"URIShort",
			"Patient/abc",
			mustNewIdentity("Patient", "abc", ""),
		},
		{
			"Invalid",
			"ThisIsNotAValidResourceName",
			nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotIdentity, _ := resource.NewIdentityFromURL(tc.URL)
			if !cmp.Equal(gotIdentity, tc.wantIdentity) {
				t.Errorf("WithNewIdentityFromURL: got %v, want %v", gotIdentity, tc.wantIdentity)
			}
		})
	}
}

func TestNewIdentityFromHistoryURL(t *testing.T) {
	testCases := []struct {
		name          string
		historyUrl    string
		expectedValue *resource.Identity
	}{
		{
			"URL",
			"https://healthcare.googleapis.com/v1/projects/my-project-name/locations/us-east4/datasets/my-dataset-name/fhirStores/my-fhir-store-name/fhir/Binary/123/_history/456",
			mustNewIdentity("Binary", "123", "456"),
		},
		{
			"URI",
			"projects/my-project-name/locations/us-east4/datasets/my-dataset-name/fhirStores/my-fhir-store-name/fhir/Patient/abc/_history/def",
			mustNewIdentity("Patient", "abc", "def"),
		},
		{
			"Invalid",
			"ThisIsNotAValidResourceName",
			nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, _ := resource.NewIdentityFromHistoryURL(tc.historyUrl)
			if ((result != nil) != (tc.expectedValue != nil)) ||
				(result != nil && tc.expectedValue != nil && *result != *tc.expectedValue) {
				t.Errorf("%s: Got = %v, want = %v", tc.name, result, tc.expectedValue)
			}
		})
	}
}

func TestIdentityEqual(t *testing.T) {
	identityA := mustNewIdentity("Patient", "A", "v1")
	testCases := []struct {
		name      string
		lhs       *resource.Identity
		rhs       *resource.Identity
		wantEqual bool
	}{
		{"both nil", nil, nil, true},
		{"lhs nil", nil, identityA, false},
		{"rhs nil", identityA, nil, false},
		{"same", identityA, mustNewIdentity("Patient", "A", "v1"), true},
		{"different type", identityA, mustNewIdentity("Person", "A", "v1"), false},
		{"different id", identityA, mustNewIdentity("Patient", "B", "v1"), false},
		{"different version", identityA, mustNewIdentity("Patient", "A", "v2"), false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotEqual := tc.lhs.Equal(tc.rhs)
			if gotEqual != tc.wantEqual {
				t.Errorf("Equal(%s) got %v want %v", tc.name, gotEqual, tc.wantEqual)
			}
		})
	}
}

func TestIdentity_Unversioned(t *testing.T) {
	withVersion := mustNewIdentity("Patient", "123", "abc")
	got := withVersion.Unversioned()
	want := mustNewIdentity("Patient", "123", "")
	if !cmp.Equal(got, want) {
		t.Errorf("Unversioned: got %v, want %v", got, want)
	}
}

func mustNewIdentity(resourceType, id, versionID string) *resource.Identity {
	identity, err := resource.NewIdentity(resourceType, id, versionID)
	if err != nil {
		panic(err)
	}
	return identity
}

func TestWithNewVersion(t *testing.T) {
	originalIdentity := mustNewIdentity("Patient", "foo", "")
	wantIdentity := mustNewIdentity("Patient", "foo", "bar")

	gotIdentity := originalIdentity.WithNewVersion("bar")

	if !cmp.Equal(gotIdentity, wantIdentity) {
		t.Errorf("WithNewVersion: got %v, want %v", gotIdentity, wantIdentity)
	}
}
