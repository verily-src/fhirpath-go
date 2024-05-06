package fhir_test

import (
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestCodeableConcept_WithoutText(t *testing.T) {
	myCoding := fhir.Coding("my-system", "my-code")
	yourCoding := fhir.Coding("your-system", "your-code")
	testCases := []struct {
		name    string
		text    string
		codings []*dtpb.Coding
		want    *dtpb.CodeableConcept
	}{
		{"empty", "", nil, &dtpb.CodeableConcept{}},
		{"full", "my-text", []*dtpb.Coding{myCoding, yourCoding},
			&dtpb.CodeableConcept{
				Coding: []*dtpb.Coding{myCoding, yourCoding},
				Text:   fhir.String("my-text"),
			},
		},
		{"without text", "", []*dtpb.Coding{myCoding},
			&dtpb.CodeableConcept{
				Coding: []*dtpb.Coding{myCoding},
				// The key behavior is the absence of the Text element.
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut := fhir.CodeableConcept(tc.text, tc.codings...)
			if diff := cmp.Diff(tc.want, sut, protocmp.Transform()); diff != "" {
				t.Errorf("CodeableConcept mismatch (-want, +got):\n%s", diff)
			}
		})
	}

}
