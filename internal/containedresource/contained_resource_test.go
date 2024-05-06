package containedresource_test

import (
	"errors"
	"net/url"
	"strings"
	"testing"

	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	bcrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/verily-src/fhirpath-go/internal/containedresource"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirtest"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"google.golang.org/protobuf/proto"
)

func TestWrap_WithNilContainedResource_ReturnsNil(t *testing.T) {
	got := containedresource.Wrap(nil)

	if got != nil {
		t.Errorf("ContainedResource: got %v, want nil", got)
	}
}

func TestWrap_WithResource_WrapsResource(t *testing.T) {
	for name, resource := range fhirtest.Resources {
		t.Run(name, func(t *testing.T) {

			got := containedresource.Wrap(resource)
			if got == nil {
				t.Fatalf("ContainedResource(%v): got nil, want value", name)
			}

			if got, want := containedresource.ID(got), resource.GetId().GetValue(); got != want {
				t.Errorf("ContainedResource(%v): got %v, want %v", name, got, want)
			}
		})
	}
}

func TestUnwrap_WithInvalidresource_ReturnsNil(t *testing.T) {
	testCases := []struct {
		name     string
		resource *bcrpb.ContainedResource
	}{
		{"NilContainedResource", nil},
		// Note: ContainedResource is not a real FHIR type, and thus should never
		// be populated as empty in practice when deserializing or receiving valid
		// FHIR payloads; so this test should never happen to begin with.
		{"EmptyContainedResource", &bcrpb.ContainedResource{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := containedresource.Unwrap(nil)

			if got != nil {
				t.Errorf("Unwrap(%v): got %v, want nil", tc.name, got)
			}
		})
	}
}

func TestUnwrap_WithContainedResource_ReturnsResource(t *testing.T) {
	for name, resource := range fhirtest.Resources {
		t.Run(name, func(t *testing.T) {
			got := containedresource.Unwrap(containedresource.Wrap(resource))

			if !proto.Equal(got, resource) {
				t.Errorf("Unwrap(%v): got %v, want %v", name, got, resource)
			}
		})
	}
}

func TestTypeOf_WithNil_Panics(t *testing.T) {
	defer func() { _ = recover() }()

	containedresource.TypeOf(nil)

	t.Errorf("TypeOf: expected panic")
}

func TestTypeOf_WithEmptyContainedResource_Panics(t *testing.T) {
	defer func() { _ = recover() }()
	cr := &bcrpb.ContainedResource{}

	containedresource.TypeOf(cr)

	t.Errorf("TypeOf: expected panic")
}

func TestTypeOf_WithResource_ReturnsName(t *testing.T) {
	for name, res := range fhirtest.Resources {
		t.Run(name, func(t *testing.T) {
			got := containedresource.TypeOf(containedresource.Wrap(res))

			if got, want := got, resource.Type(name); got != want {
				t.Errorf("ContainedResourceName(%v): got %v, want %v", name, got, want)
			}
		})
	}
}

func TestID_WithNil_ReturnsEmptyString(t *testing.T) {
	got := containedresource.ID(nil)

	if got != "" {
		t.Errorf("ID: got %v, want empty string", got)
	}
}

func TestID_WithEmptyContainedResource_ReturnsEmptyString(t *testing.T) {
	cr := &bcrpb.ContainedResource{}

	got := containedresource.ID(cr)

	if got != "" {
		t.Errorf("ContainedResourceID: got %v, want empty string", got)
	}
}

func TestID_WithResource_ReturnsID(t *testing.T) {
	for name, resource := range fhirtest.Resources {
		t.Run(name, func(t *testing.T) {
			got := containedresource.ID(containedresource.Wrap(resource))

			if got, want := got, resource.GetId().GetValue(); got != want {
				t.Errorf("ContainedResourceID(%v): got %v, want %v", name, got, want)
			}
		})
	}
}

func TestGenerateIfNoneExist_Errors(t *testing.T) {
	patient0 := &patient_go_proto.Patient{
		Id: fhir.ID("12345"),
	}
	patient1 := &patient_go_proto.Patient{
		Id: fhir.ID("12345"),
		Identifier: []*dtpb.Identifier{
			&dtpb.Identifier{
				System: &dtpb.Uri{Value: "http://fake.com"},
				Value:  &dtpb.String{Value: "9efbf82d-7a58-4d14-bec1-63f8fda148a8"},
			},
		},
	}

	testCases := []struct {
		name    string
		input   *bcrpb.ContainedResource
		wantErr string
	}{
		{
			"empty ContainedResource",
			&bcrpb.ContainedResource{},
			"Unwrap() returned nil / no contained resource",
		},
		{
			"nil ContainedResource",
			nil,
			"ContainedResource is nil",
		},
		{
			"No identifier, emptyIsErr=true",
			containedresource.Wrap(patient0),
			"found no Identifiers",
		},
		{
			"No matching identifier, emptyIsErr=true",
			containedresource.Wrap(patient1),
			"found no Identifiers",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			header, err := containedresource.GenerateIfNoneExist(tc.input, "no-such-system", true)

			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("got error %#v, want %#v", err.Error(), tc.wantErr)
			}
			if !errors.Is(err, containedresource.ErrGenerateIfNoneExist) {
				t.Errorf("got error %#v, want errors.Is(..., ErrGenerateIfNoneExist)", err)
			}
			if header != "" {
				t.Errorf("got %v, want empty string", header)
			}
		})
	}
}

func TestGenerateIfNoneExist(t *testing.T) {
	patient0 := &patient_go_proto.Patient{
		Id: fhir.ID("12345"),
	}
	patient1 := &patient_go_proto.Patient{
		Id: fhir.ID("12345"),
		Identifier: []*dtpb.Identifier{
			&dtpb.Identifier{
				System: &dtpb.Uri{Value: "http://fake.com"},
				Value:  &dtpb.String{Value: "9efbf82d-7a58-4d14-bec1-63f8fda148a8"},
			},
		},
	}
	patient2 := &patient_go_proto.Patient{
		Id: fhir.ID("12345"),
		Identifier: []*dtpb.Identifier{
			&dtpb.Identifier{
				System: &dtpb.Uri{Value: "http://fake.com"},
				Value:  &dtpb.String{Value: "9efbf82d-7a58-4d14-bec1-63f8fda148a8"},
			},
			&dtpb.Identifier{
				Use: &dtpb.Identifier_UseCode{
					Value: codes_go_proto.IdentifierUseCode_USUAL,
				},
				System: &dtpb.Uri{Value: "urn:oid:2.16.840.1.113883.2.4.6.3"},
				Value:  &dtpb.String{Value: "12345"},
			},
		},
	}

	patient3 := fhirtest.NewResource(t, "Patient", fhirtest.WithGeneratedIdentifier("http://example.com/fake-id")).(*patient_go_proto.Patient)

	patient4 := &patient_go_proto.Patient{
		Id: fhir.ID("12345"),
		Identifier: []*dtpb.Identifier{
			&dtpb.Identifier{
				System: &dtpb.Uri{Value: "http://fake.com"},
				Value:  &dtpb.String{Value: "foo,bar,baz|omg"},
			},
		},
	}

	patient5 := &patient_go_proto.Patient{
		Id: fhir.ID("12345"),
		Identifier: []*dtpb.Identifier{
			&dtpb.Identifier{
				System: &dtpb.Uri{Value: "http://fake.com"},
				Value:  &dtpb.String{Value: "9efbf82d-7a58-4d14-bec1-63f8fda148a8"},
			},
			&dtpb.Identifier{
				System: &dtpb.Uri{Value: "http://fake.com"},
				Value:  &dtpb.String{Value: "7d541708-b068-4347-a8dc-cce1dcdb5314"},
			},
		},
	}

	testCases := []struct {
		name    string
		patient *patient_go_proto.Patient
		system  string
		want    string
		wantErr string
	}{
		{"Patient with no Identifier, emptyIsErr false", patient0, "system", "", ""},
		{"Patient with single Identifier", patient1, "http://fake.com", "identifier=" + url.QueryEscape("http://fake.com|9efbf82d-7a58-4d14-bec1-63f8fda148a8"), ""},
		{"Patient with two Identifiers but one matching", patient2, "http://fake.com", "identifier=" + url.QueryEscape("http://fake.com|9efbf82d-7a58-4d14-bec1-63f8fda148a8"), ""},
		{"Patient with two Identifiers and both matching", patient5, "http://fake.com", "", "found multiple Identifiers"},
		{"Patient with generated ID", patient3, "http://example.com/fake-id", "identifier=" + url.QueryEscape("http://example.com/fake-id|"+patient3.Identifier[0].Value.Value), ""},
		{"Special chars in Identifier", patient4, "http://fake.com", "identifier=" + url.QueryEscape(`http://fake.com|foo\,bar\,baz\|omg`), ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cr := containedresource.Wrap(tc.patient)

			got, err := containedresource.GenerateIfNoneExist(cr, tc.system, false)

			if tc.wantErr == "" {
				if err != nil {
					t.Errorf("%#v: Bad If-None-Exist:\n  got err %#v\n  want nil", tc.name, err)
				}
			} else {
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Errorf("%#v: got error %#v, want %#v", tc.name, err.Error(), tc.wantErr)
				}
				if !errors.Is(err, containedresource.ErrGenerateIfNoneExist) {
					t.Errorf("%#v: got error %#v, want errors.Is(..., ErrGenerateIfNoneExist)", tc.name, err)
				}
			}

			if got != tc.want {
				t.Errorf("%#v: Bad If-None-Exist:\n  got  %#v\n  want %#v", tc.name, got, tc.want)
			}
		})
	}
}
