package reference_test

import (
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/internal/element/reference"
	"github.com/verily-src/fhirpath-go/internal/resource"
)

func newIdentity(t *testing.T, typeName, id, version string) *resource.Identity {
	t.Helper()
	ident, err := resource.NewIdentity(typeName, id, version)
	if err != nil {
		t.Fatalf("NewIdentity: %v", err)
	}
	return ident
}

func TestIdentityOf_BadReference_ReturnsError(t *testing.T) {
	testCases := []struct {
		name      string
		reference *dtpb.Reference
		wantErr   error
	}{
		{
			"invalid absolute uri",
			&dtpb.Reference{
				Reference: &dtpb.Reference_Uri{
					Uri: &dtpb.String{
						Value: "https://example.com",
					},
				},
			},
			reference.ErrInvalidURI,
		},
		{
			"fragment without type",
			&dtpb.Reference{
				Reference: &dtpb.Reference_Fragment{
					Fragment: &dtpb.String{
						Value: "https://example.com",
					},
				},
			},
			reference.ErrFragmentMissingType,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := reference.IdentityOf(tc.reference)

			got, want := err, tc.wantErr
			if !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("IdentityOf(%s) error got '%v', want '%v'", tc.name, got, want)
			}
		})
	}
}

func TestIdentityOf(t *testing.T) {
	testCases := []struct {
		name      string
		reference *dtpb.Reference
		want      *resource.Identity
	}{
		{
			"Reference",
			&dtpb.Reference{
				Reference: &dtpb.Reference_AccountId{
					AccountId: &dtpb.ReferenceId{Value: "123"},
				},
			},
			newIdentity(t, "Account", "123", ""),
		},
		{
			"Reference with history",
			&dtpb.Reference{
				Reference: &dtpb.Reference_PatientId{
					PatientId: &dtpb.ReferenceId{Value: "123", History: &dtpb.Id{Value: "abc"}},
				},
			},
			newIdentity(t, "Patient", "123", "abc"),
		},
		{
			"Relative uri reference",
			&dtpb.Reference{
				Reference: &dtpb.Reference_Uri{
					Uri: &dtpb.String{
						Value: "Patient/123",
					},
				},
			},
			newIdentity(t, "Patient", "123", ""),
		},
		{
			"Relative uri reference with history",
			&dtpb.Reference{
				Reference: &dtpb.Reference_Uri{
					Uri: &dtpb.String{
						Value: "Patient/123/_history/abc",
					},
				},
			},
			newIdentity(t, "Patient", "123", "abc"),
		},
		{
			"Absolute uri reference",
			&dtpb.Reference{
				Reference: &dtpb.Reference_Uri{
					Uri: &dtpb.String{
						Value: "https://healthcare.googleapis.com/v1/projects/123/locations/abc/datasets/def/fhirStores/ghi/fhir/Patient/123",
					},
				},
			},
			newIdentity(t, "Patient", "123", ""),
		},
		{
			"Absolute uri reference with history",
			&dtpb.Reference{
				Reference: &dtpb.Reference_Uri{
					Uri: &dtpb.String{
						Value: "https://healthcare.googleapis.com/v1/projects/123/locations/abc/datasets/def/fhirStores/ghi/fhir/Patient/123/_history/abc",
					},
				},
			},
			newIdentity(t, "Patient", "123", "abc"),
		},
		{
			"Fragment reference",
			&dtpb.Reference{
				Reference: &dtpb.Reference_Fragment{
					Fragment: &dtpb.String{
						Value: "123",
					},
				},
				Type: &dtpb.Uri{
					Value: "Patient",
				},
			},
			newIdentity(t, "Patient", "123", ""),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ident, err := reference.IdentityOf(tc.reference)

			if err != nil {
				t.Fatalf("IdentityOf(%s) error got %v, want nil", tc.name, err)
			}
			if got, want := ident, tc.want; !got.Equal(want) {
				t.Errorf("IdentityOf(%s) got %s, want %s", tc.name, got, want)
			}
		})
	}
}

func TestIdentityFromAbsoluteURL_BadInput_ReturnsError(t *testing.T) {
	testCases := []struct {
		name    string
		url     string
		wantErr error
	}{
		{
			"url with fragment",
			"https://healthcare.googleapis.com/v1/projects/123/locations/abc/datasets/def/fhirStores/ghi/fhir/Patient/123#crID",
			reference.ErrInvalidURI,
		},
		{
			"canonical url",
			"https://healthcare.googleapis.com/v1/projects/123/locations/abc/datasets/def/fhirStores/ghi/fhir/Patient/123|1",
			reference.ErrInvalidURI,
		},
		{
			"invalid server url",
			"https://not/a/fhir/server/url",
			cmpopts.AnyError,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := reference.IdentityFromAbsoluteURL(tc.url)

			got, want := err, tc.wantErr
			if !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("IdentityFromAbsoluteURL(%s) error got %v, want %v", tc.name, got, want)
			}
		})
	}
}

func TestIdentityFromAbsoluteURL(t *testing.T) {
	testCases := []struct {
		name string
		url  string
		want *resource.Identity
	}{
		{
			"no version",
			"https://healthcare.googleapis.com/v1/projects/123/locations/abc/datasets/def/fhirStores/ghi/fhir/Patient/123",
			newIdentity(t, "Patient", "123", ""),
		},
		{
			"versioned",
			"https://healthcare.googleapis.com/v1/projects/123/locations/abc/datasets/def/fhirStores/ghi/fhir/Patient/123/_history/abc",
			newIdentity(t, "Patient", "123", "abc"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ident, err := reference.IdentityFromAbsoluteURL(tc.url)

			if err != nil {
				t.Fatalf("IdentityFromAbsoluteURL(%s) error got %v, want nil", tc.name, err)
			}
			if got, want := ident, tc.want; !cmp.Equal(got, want) {
				t.Errorf("IdentityFromAbsoluteURL(%s) got %v, want %v", tc.name, got, want)
			}
		})
	}
}

func TestIdentityFromRelativeURI_BadURL_ReturnsError(t *testing.T) {
	_, err := reference.IdentityFromRelativeURI("Patient/123/_history")

	if got, want := err, reference.ErrInvalidRelativeURI; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
		t.Fatalf("IdentityFromRelativeURI error got %v, want %v", got, want)
	}
}

func TestIdentityFromRelativeURI(t *testing.T) {
	testCases := []struct {
		name string
		url  string
		want *resource.Identity
	}{
		{
			"no version",
			"Patient/123",
			newIdentity(t, "Patient", "123", ""),
		},
		{
			"versioned",
			"Patient/123/_history/abc",
			newIdentity(t, "Patient", "123", "abc"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ident, err := reference.IdentityFromRelativeURI(tc.url)

			if err != nil {
				t.Fatalf("IdentityFromRelativeURI(%s) error got %v, want nil", tc.name, err)
			}
			if got, want := ident, tc.want; !cmp.Equal(got, want) {
				t.Errorf("IdentityFromRelativeURI(%s) got %v, want %v", tc.name, got, want)
			}
		})
	}
}

func TestIdentityFromURL_BadURL_ReturnsError(t *testing.T) {
	_, err := reference.IdentityFromURL("Patient/123/_history")

	if got, want := err, reference.ErrInvalidURI; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
		t.Fatalf("IdentityFromURL error got %v, want %v", got, want)
	}
}

func TestIdentityFromURL(t *testing.T) {
	testCases := []struct {
		name string
		url  string
		want *resource.Identity
	}{
		{
			"absolute url",
			"https://healthcare.googleapis.com/v1/projects/123/locations/abc/datasets/def/fhirStores/ghi/fhir/Patient/123",
			newIdentity(t, "Patient", "123", ""),
		},
		{
			"relative uri",
			"Patient/123/_history/abc",
			newIdentity(t, "Patient", "123", "abc"),
		},
		{
			"relative uri using RequestGroup resource type",
			"RequestGroup/123/_history/abc",
			newIdentity(t, "RequestGroup", "123", "abc"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ident, err := reference.IdentityFromURL(tc.url)

			if err != nil {
				t.Fatalf("IdentityFromURL(%s) error got %v, want nil", tc.name, err)
			}
			if got, want := ident, tc.want; !cmp.Equal(got, want) {
				t.Errorf("IdentityFromURL(%s) got %v, want %v", tc.name, got, want)
			}
		})
	}
}
