package contactable_test

import (
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirtest"
	"github.com/verily-src/fhirpath-go/internal/resource/contactable"
	"google.golang.org/protobuf/proto"
)

func TestWithContact(t *testing.T) {
	want := &dtpb.ContactDetail{
		Name: fhir.String("deadbeef"),
	}

	for name, res := range fhirtest.CanonicalResources {
		t.Run(name, func(t *testing.T) {
			got := proto.Clone(res).(contactable.ContactableResource)

			contactable.Update(got, contactable.WithContacts(want))

			for _, got := range got.(fhir.CanonicalResource).GetContact() {
				if !proto.Equal(got, want) {
					t.Errorf("WithContact(%v): got %v, want %v", name, got, want)
				}
			}
		})
	}
}
