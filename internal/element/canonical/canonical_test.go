package canonical_test

import (
	"errors"
	"fmt"
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	qpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/questionnaire_go_proto"

	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/internal/element/canonical"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestNew(t *testing.T) {
	const (
		url      = "https://example.com"
		version  = "1.2.3"
		fragment = "frag"
	)
	testCases := []struct {
		name string
		opts []canonical.Option
		want string
	}{
		{
			name: "NoOpts",
			opts: nil,
			want: url,
		},
		{
			name: "WithVersion",
			opts: []canonical.Option{canonical.WithVersion(version)},
			want: fmt.Sprintf("%v|%v", url, version),
		},
		{
			name: "WithFragment",
			opts: []canonical.Option{canonical.WithFragment(fragment)},
			want: fmt.Sprintf("%v#%v", url, fragment),
		},
		{
			name: "WithVersionAndFragment",
			opts: []canonical.Option{canonical.WithVersion(version), canonical.WithFragment(fragment)},
			want: fmt.Sprintf("%v|%v#%v", url, version, fragment),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := canonical.New(url, tc.opts...)

			if got, want := got.GetValue(), tc.want; got != want {
				t.Errorf("Canonical(%v): got %v, want %v", tc.name, got, want)
			}
		})
	}
}

func TestCanonicalFromResource_Nil_ReturnsError(t *testing.T) {
	_, err := canonical.FromResource(nil)

	if got, want := err, canonical.ErrNoCanonicalURL; !errors.Is(got, want) {
		t.Errorf("CanonicalFromResource: got %v, want %v", got, want)
	}
}

func makeCanonicalResource() fhir.CanonicalResource {
	const (
		url = "https://example.com"
	)
	return &qpb.Questionnaire{
		Id:      fhir.ID("0xdeadbeef"),
		Url:     fhir.URI(url),
		Version: fhir.String("1.0.0"),
	}
}

func TestCanonicalFromResource_ReturnsCanonicalWithUrl(t *testing.T) {
	resource := makeCanonicalResource()

	got, err := canonical.FromResource(resource)

	if err != nil {
		t.Fatalf("CanonicalFromResource: unexpected error: %v", err)
	}
	if got, want := got.GetValue(), resource.GetUrl().GetValue(); got != want {
		t.Errorf("CanonicalFromResource: got '%v', want '%v'", got, want)
	}
}

func TestVersionedFromResource_Nil_ReturnsNil(t *testing.T) {
	_, err := canonical.VersionedFromResource(nil)

	if got, want := err, canonical.ErrNoCanonicalURL; !errors.Is(got, want) {
		t.Errorf("VersionedFromResource: got %v, want %v", got, want)
	}
}

func TestVersionedCanonical_ReturnsCanonicalWithUrl(t *testing.T) {
	const (
		url     = "https://example.com"
		version = "1.2.3"
	)
	testCases := []struct {
		name     string
		resource fhir.CanonicalResource
		want     string
	}{
		{
			name: "NoVersion",
			resource: &qpb.Questionnaire{
				Url: fhir.URI(url),
			},
			want: url,
		},
		{
			name: "WithVersion",
			resource: &qpb.Questionnaire{
				Url:     fhir.URI(url),
				Version: fhir.String(version),
			},
			want: fmt.Sprintf("%v|%v", url, version),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := canonical.VersionedFromResource(tc.resource)

			if err != nil {
				t.Fatalf("VersionedCanonical(%v): unexpected error: %v", tc.name, err)
			}
			if got, want := got.GetValue(), tc.want; got != want {
				t.Errorf("VersionedCanonical(%v): got '%v', want '%v'", tc.name, got, want)
			}
		})
	}
}

func TestFragmentFromResourceFromResourceNil_ReturnsNil(t *testing.T) {
	_, err := canonical.FragmentFromResource(nil)

	if got, want := err, canonical.ErrNoCanonicalURL; !errors.Is(got, want) {
		t.Errorf("CanonicalFragmentFromResource: got %v, want %v", got, want)
	}
}

func TestFragmentFromResource(t *testing.T) {
	resource := makeCanonicalResource()
	want := fmt.Sprintf("%v#%v", resource.GetUrl().GetValue(), resource.GetId().GetValue())

	got, err := canonical.FragmentFromResource(resource)

	if err != nil {
		t.Fatalf("FragmentFromResource: unexpected error: %v", err)
	}
	if got := got.GetValue(); got != want {
		t.Errorf("FragmentFromResource: got '%v', want '%v'", got, want)
	}
}

func TestIdentityFromCanonical(t *testing.T) {
	testCases := []struct {
		name, reference, url, version, fragment string
	}{
		{
			name:      "basic",
			url:       "http://someurl/test-value",
			reference: "http://someurl/test-value",
		},
		{
			name:      "long url",
			url:       "https://fhir.acme.com/Questionnaire/example",
			reference: "https://fhir.acme.com/Questionnaire/example",
		},
		{
			name:      "with version",
			url:       "https://fhir.acme.com/Questionnaire/example",
			version:   "1.0.0",
			reference: "https://fhir.acme.com/Questionnaire/example|1.0.0",
		},
		{
			name:      "with fragment",
			url:       "http://hl7.org/fhir/ValueSet/my-valueset",
			fragment:  "vs1",
			reference: "http://hl7.org/fhir/ValueSet/my-valueset#vs1",
		},
		{
			name:      "with version and fragment",
			url:       "http://fhir.acme.com/Questionnaire/example",
			version:   "1.0",
			fragment:  "vs1",
			reference: "http://fhir.acme.com/Questionnaire/example|1.0#vs1",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := &dtpb.Canonical{
				Value: tc.reference,
			}
			got, err := canonical.IdentityFromReference(c)
			if err != nil {
				t.Fatalf("IdentityFromReference(%s): unexpected error: %v", tc.name, err)
			}
			want := &resource.CanonicalIdentity{
				Url:      tc.url,
				Version:  tc.version,
				Fragment: tc.fragment,
			}
			if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
				t.Errorf("IdentityFromReference(%s): %v", tc.name, diff)
			}
			if diff := cmp.Diff(tc.reference, got.String()); diff != "" {
				t.Errorf("IdentityFromReference(%s).String: %v", tc.name, diff)
			}
		})
	}
}
