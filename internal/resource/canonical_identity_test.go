package resource_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/internal/fhirtest"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestCanonicalIdentity_EmptyURL_ReturnsError(t *testing.T) {
	_, got := resource.NewCanonicalIdentity("", "v1", "")

	if got != resource.ErrMissingCanonicalURL {
		t.Errorf("NewCanonicalIdentity: got %v, want %v", got, resource.ErrMissingCanonicalURL)
	}
}

func TestCanonicalIdentity(t *testing.T) {
	testCases := []struct {
		name, url, version, fragment string
		want                         *resource.CanonicalIdentity
		wantString                   string
		wantType                     resource.Type
		hasType                      bool
	}{
		{
			name:       "basic",
			url:        "http://someurl/test-value",
			wantString: "http://someurl/test-value",
		},
		{
			name:       "long url",
			url:        "https://fhir.acme.com/Questionnaire/example",
			wantString: "https://fhir.acme.com/Questionnaire/example",
			hasType:    true,
			wantType:   resource.Questionnaire,
		},
		{
			name:       "with version",
			url:        "https://fhir.acme.com/PlanDefinition/example",
			version:    "1.0.0",
			wantString: "https://fhir.acme.com/PlanDefinition/example|1.0.0",
			hasType:    true,
			wantType:   resource.PlanDefinition,
		},
		{
			name:       "with fragment",
			url:        "http://hl7.org/fhir/ValueSet/my-valueset",
			fragment:   "vs1",
			wantString: "http://hl7.org/fhir/ValueSet/my-valueset#vs1",
		},
		{
			name:       "with version and fragment",
			url:        "http://fhir.acme.com/ActivityDefinition/example",
			version:    "1.0",
			fragment:   "vs1",
			wantString: "http://fhir.acme.com/ActivityDefinition/example|1.0#vs1",
			hasType:    true,
			wantType:   resource.ActivityDefinition,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := resource.NewCanonicalIdentity(tc.url, tc.version, tc.fragment)
			want := &resource.CanonicalIdentity{
				Url:      tc.url,
				Version:  tc.version,
				Fragment: tc.fragment,
			}

			if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
				t.Errorf("CanonicalIdentity(%s): %v", tc.name, diff)
			}

			if s := got.String(); tc.wantString != s {
				t.Errorf("CanonicalIdentity(%s).String: want: %s, got: %s", tc.name, tc.wantString, s)
			}
			if tc.hasType {
				gt, ok := got.Type()
				if !ok || gt != tc.wantType {
					t.Errorf("CanonicalIdentity(%s).Type: want: %s, got: %s", tc.name, tc.wantType, gt)
				}
			}
		})
	}
}

func TestIsCanonicalResourceType(t *testing.T) {
	wantTypes := []string{
		"ActivityDefinition",
		"CapabilityStatement",
		"ChargeItemDefinition",
		"CodeSystem",
		"CompartmentDefinition",
		"ConceptMap",
		"Contract",
		"EffectEvidenceSynthesis",
		"EventDefinition",
		"Evidence",
		"EvidenceVariable",
		"ExampleScenario",
		"GraphDefinition",
		"ImplementationGuide",
		"Library",
		"Measure",
		"MessageDefinition",
		"OperationDefinition",
		"PlanDefinition",
		"Questionnaire",
		"ResearchDefinition",
		"ResearchElementDefinition",
		"RiskEvidenceSynthesis",
		"SearchParameter",
		"StructureDefinition",
		"StructureMap",
		"TerminologyCapabilities",
		"TestScript",
		"ValueSet",
	}
	gotTypes := []string{}
	for name := range fhirtest.Resources {
		if resource.IsCanonicalType(name) {
			gotTypes = append(gotTypes, name)
		}
	}
	if diff := cmp.Diff(wantTypes, gotTypes, cmpopts.SortSlices(strings.Compare)); diff != "" {
		t.Errorf("IsCanonicalType: %v", diff)
	}
}
