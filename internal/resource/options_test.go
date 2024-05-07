package resource_test

import (
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirtest"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"google.golang.org/protobuf/proto"
)

func TestWithMeta(t *testing.T) {
	want := &dtpb.Meta{
		VersionId: &dtpb.Id{
			Value: "deadbeef",
		},
	}

	for name, res := range fhirtest.Resources {
		t.Run(name, func(t *testing.T) {
			got := proto.Clone(res).(fhir.Resource)

			resource.Update(got, resource.WithMeta(want))

			if got, want := got.GetMeta(), want; !proto.Equal(got, want) {
				t.Errorf("WithMeta(%v): got %v, want %v", name, got, want)
			}
		})
	}
}

func TestWithID(t *testing.T) {
	const id = "123456789"
	want := fhir.ID(id)

	for name, res := range fhirtest.Resources {
		t.Run(name, func(t *testing.T) {
			got := proto.Clone(res).(fhir.Resource)

			resource.Update(got, resource.WithID(id))

			if got, want := got.GetId(), want; !proto.Equal(got, want) {
				t.Errorf("WithID(%v): got %v, want %v", name, got, want)
			}
		})
	}
}

func TestWithImplicitRules(t *testing.T) {
	const rules = "https://example.com/some/rules"
	want := fhir.URI(rules)

	for name, res := range fhirtest.Resources {
		t.Run(name, func(t *testing.T) {
			got := proto.Clone(res).(fhir.Resource)

			resource.Update(got, resource.WithImplicitRules(rules))

			if got, want := got.GetImplicitRules(), want; !proto.Equal(got, want) {
				t.Errorf("WithImplicitRules(%v): got %v, want %v", name, got, want)
			}
		})
	}
}

func TestWithLanguage(t *testing.T) {
	const language = "en-gb"
	want := fhir.Code(language)

	for name, res := range fhirtest.Resources {
		t.Run(name, func(t *testing.T) {
			got := proto.Clone(res).(fhir.Resource)

			resource.Update(got, resource.WithLanguage(language))

			if got, want := got.GetLanguage(), want; !proto.Equal(got, want) {
				t.Errorf("WithLanguage(%v): got %v, want %v", name, got, want)
			}
		})
	}
}
