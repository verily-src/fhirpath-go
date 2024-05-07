package element_test

// For extraction, we test Reference and CodeableConcept elements.
// There shouldn't be any need to exhastively test all possible element types,
// since they they are all proto messages (even scalar-ish elements like
// strings and URIs).

import (
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	acpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/account_go_proto"
	appb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/appointment_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/internal/slices"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/bundle"
	"github.com/verily-src/fhirpath-go/internal/element"
	"github.com/verily-src/fhirpath-go/internal/element/extension"
	"github.com/verily-src/fhirpath-go/internal/fhirtest"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func Test_ExtractAll_OfReference_ValuesCorrect(t *testing.T) {
	refUri := &dtpb.Reference{
		Reference: &dtpb.Reference_Uri{
			Uri: fhir.String("uri-ref"),
		},
	}
	refRelatedPersonId := &dtpb.Reference{
		Reference: &dtpb.Reference_RelatedPersonId{
			RelatedPersonId: &dtpb.ReferenceId{
				Id: fhir.String("related-ref"),
			},
		},
	}
	testCases := []struct {
		name      string
		resource  fhir.Resource
		wantRefs  []*dtpb.Reference
		wantPaths []string
	}{
		{
			name:     "No reference",
			resource: fhirtest.NewResource(t, "Patient"),
			wantRefs: []*dtpb.Reference{},
		},
		{
			name: "Single reference",
			resource: fhirtest.NewResource(t, "Patient", fhirtest.WithResourceModification(func(p *ppb.Patient) {
				p.ManagingOrganization = refUri
			})),
			wantRefs:  []*dtpb.Reference{refUri},
			wantPaths: []string{"Patient.managingOrganization"},
		},
		{
			name: "Multiple references",
			resource: fhirtest.NewResource(t, "Appointment", fhirtest.WithResourceModification(func(a *appb.Appointment) {
				a.Participant = []*appb.Appointment_Participant{
					{
						Actor: refRelatedPersonId,
					},
					{
						Type:  []*dtpb.CodeableConcept{fhir.CodeableConcept("", fhir.Coding("systest", "code"))},
						Actor: refUri,
					},
				}
			})),
			wantRefs:  []*dtpb.Reference{refRelatedPersonId, refUri},
			wantPaths: []string{"Appointment.participant[0].actor", "Appointment.participant[1].actor"},
		},
		{
			name: "Repeated field references",
			resource: fhirtest.NewResource(t, "Account", fhirtest.WithResourceModification(func(a *acpb.Account) {
				a.Subject = []*dtpb.Reference{refRelatedPersonId, refUri}
			})),
			wantRefs:  []*dtpb.Reference{refRelatedPersonId, refUri},
			wantPaths: []string{"Account.subject[0]", "Account.subject[1]"},
		},
		{
			name: "Repeated identical references",
			resource: fhirtest.NewResource(t, "Account", fhirtest.WithResourceModification(func(a *acpb.Account) {
				a.Subject = []*dtpb.Reference{refRelatedPersonId, refRelatedPersonId}
			})),
			wantRefs:  []*dtpb.Reference{refRelatedPersonId, refRelatedPersonId},
			wantPaths: []string{"Account.subject[0]", "Account.subject[1]"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			gotReferences, _ := element.ExtractAll[*dtpb.Reference](testCase.resource)

			opts := []cmp.Option{
				protocmp.Transform(),
				cmpopts.SortSlices(func(a *dtpb.Reference, b *dtpb.Reference) bool { return a.String() < b.String() }),
			}
			if !cmp.Equal(testCase.wantRefs, gotReferences, opts...) {
				t.Errorf("ExtractAll(): got '%v', want '%v'", gotReferences, testCase.wantRefs)
			}

			gotRefsWithPath, err := element.ExtractAllWithPath[*dtpb.Reference](testCase.resource)
			if err != nil {
				t.Fatalf("ExtractAllElementWithPath(%s): got error %v", testCase.name, err)
			}
			wantRefsWithPath := zipElementsAndPaths(testCase.wantRefs, testCase.wantPaths)
			if !cmp.Equal(gotRefsWithPath, wantRefsWithPath, protocmp.Transform(), cmpopts.EquateEmpty()) {
				t.Errorf("ExtractAllElementWithPath(%s): got '%v', want '%v'", testCase.name, gotRefsWithPath, wantRefsWithPath)
			}
		})
	}
}

func Test_ExtractAll_OfCodeableConcept_ValuesCorrect(t *testing.T) {
	concept1 := fhir.CodeableConcept("thing1")
	concept2 := fhir.CodeableConcept("thing2")
	concept3 := fhir.CodeableConcept("thing2")
	testCases := []struct {
		name         string
		resource     fhir.Resource
		wantConcepts []*dtpb.CodeableConcept
		wantPaths    []string
		wantError    error
	}{
		{
			name: "single",
			resource: fhirtest.NewResource(t, "Patient", fhirtest.WithResourceModification(func(p *ppb.Patient) {
				p.MaritalStatus = concept1
			})),
			wantConcepts: []*dtpb.CodeableConcept{concept1},
			wantPaths:    []string{"Patient.maritalStatus"},
		},
		{
			name: "repeated nested repeated with dup",
			resource: fhirtest.NewResource(t, "Patient", fhirtest.WithResourceModification(func(p *ppb.Patient) {
				p.Contact = []*ppb.Patient_Contact{
					{Relationship: []*dtpb.CodeableConcept{concept1, concept2}},
					{Relationship: []*dtpb.CodeableConcept{concept2, concept3}},
				}
			})),
			wantConcepts: []*dtpb.CodeableConcept{concept1, concept2, concept2, concept3},
			wantPaths: []string{
				"Patient.contact[0].relationship[0]",
				"Patient.contact[0].relationship[1]",
				"Patient.contact[1].relationship[0]",
				"Patient.contact[1].relationship[1]",
			},
		},
		{
			name: "repeated extension",
			resource: fhirtest.NewResource(t, "Patient", fhirtest.WithResourceModification(func(p *ppb.Patient) {
				p.Extension = []*dtpb.Extension{
					extension.New("my-extension-url", concept1),
					extension.New("my-extension-url", concept2),
				}
			})),
			wantConcepts: []*dtpb.CodeableConcept{concept1, concept2},
			wantPaths: []string{
				"Patient.extension[0].valueCodeableConcept",
				"Patient.extension[1].valueCodeableConcept",
			},
		},
		{
			name: "inside ContainedResource inside Bundle",
			resource: bundle.NewCollection(bundle.WithEntries(bundle.NewCollectionEntry(
				fhirtest.NewResource(t, "Patient", fhirtest.WithResourceModification(func(p *ppb.Patient) {
					p.MaritalStatus = concept1
				}))))),
			wantError: element.ErrFhirPathNotImplemented,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotElementsWithPath, gotErr := element.ExtractAllWithPath[*dtpb.CodeableConcept](tc.resource)
			if !cmp.Equal(gotErr, tc.wantError, cmpopts.EquateErrors()) {
				t.Fatalf("ExtractAllWithPath(%s) error got '%v', want '%v'", tc.name, gotErr, tc.wantError)
			}
			wantElementsWithPath := zipElementsAndPaths(tc.wantConcepts, tc.wantPaths)
			if !cmp.Equal(gotElementsWithPath, wantElementsWithPath, protocmp.Transform(), cmpopts.EquateEmpty()) {
				t.Errorf("ExtractAllElementWithPath(%s): got '%v', want '%v'",
					tc.name, gotElementsWithPath, wantElementsWithPath)
			}
		})
	}
}

func Test_ExtractAll_OfSpecial(t *testing.T) {
	date1 := fhir.DateNow()
	datetime1 := fhir.DateTimeNow()
	instant1 := fhir.InstantNow()
	instant2 := fhir.InstantNow()
	instant3 := fhir.InstantNow()
	time1 := fhir.TimeNow()

	testCases := []struct {
		name         string
		resource     fhir.Resource
		extractFunc  func(fhir.Resource) ([]element.ElementWithPath[fhir.Element], error)
		wantElements []fhir.Element
		wantPaths    []string
	}{
		{name: "simple Date",
			resource: fhirtest.NewResource(t, "Patient", fhirtest.WithResourceModification(func(p *ppb.Patient) {
				p.BirthDate = date1
			})),
			extractFunc: func(res fhir.Resource) ([]element.ElementWithPath[fhir.Element], error) {
				asDates, err := element.ExtractAllWithPath[*dtpb.Date](res)
				return asElements(asDates), err
			},
			wantElements: []fhir.Element{date1},
			wantPaths:    []string{"Patient.birthDate"},
		},
		{name: "DateTime inside choice",
			resource: fhirtest.NewResource(t, "Patient", fhirtest.WithResourceModification(func(p *ppb.Patient) {
				p.Deceased = &ppb.Patient_DeceasedX{
					Choice: &ppb.Patient_DeceasedX_DateTime{DateTime: datetime1},
				}
			})),
			extractFunc: func(res fhir.Resource) ([]element.ElementWithPath[fhir.Element], error) {
				asDateTimes, err := element.ExtractAllWithPath[*dtpb.DateTime](res)
				return asElements(asDateTimes), err
			},
			wantElements: []fhir.Element{datetime1},
			wantPaths:    []string{"Patient.deceasedDateTime"},
		},
		{
			name: "Instant",
			resource: fhirtest.NewResource(t, "Appointment", fhirtest.WithResourceModification(func(a *appb.Appointment) {
				a.Meta.LastUpdated = instant1
				a.Start = instant2
				a.End = instant3
			})),
			extractFunc: func(res fhir.Resource) ([]element.ElementWithPath[fhir.Element], error) {
				asInstants, err := element.ExtractAllWithPath[*dtpb.Instant](res)
				return asElements(asInstants), err
			},
			wantElements: []fhir.Element{instant3, instant1, instant2},
			wantPaths:    []string{"Appointment.end", "Appointment.meta.lastUpdated", "Appointment.start"},
		},
		{name: "Time inside extension",
			resource: fhirtest.NewResource(t, "Patient", fhirtest.WithResourceModification(func(p *ppb.Patient) {
				p.Extension = []*dtpb.Extension{
					extension.New("my-extension-url", time1),
				}
			})),
			extractFunc: func(res fhir.Resource) ([]element.ElementWithPath[fhir.Element], error) {
				asTimes, err := element.ExtractAllWithPath[*dtpb.Time](res)
				return asElements(asTimes), err
			},
			wantElements: []fhir.Element{time1},
			wantPaths:    []string{"Patient.extension[0].valueTime"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotElementsWithPath, _ := tc.extractFunc(tc.resource)
			wantElementsWithPath := zipElementsAndPaths(tc.wantElements, tc.wantPaths)
			if !cmp.Equal(gotElementsWithPath, wantElementsWithPath, protocmp.Transform(), cmpopts.EquateEmpty()) {
				t.Errorf("ExtractAllElementWithPath(%s): got '%v', want '%v'",
					tc.name, gotElementsWithPath, wantElementsWithPath)
			}
		})
	}
}

// asElements maps a slice of ElementWithPath[elementT]
// to a new slice of [ElementWithPath[fhir.Element].
func asElements[elementT fhir.Element](elements []element.ElementWithPath[elementT]) []element.ElementWithPath[fhir.Element] {
	return slices.Map(elements,
		func(e element.ElementWithPath[elementT]) element.ElementWithPath[fhir.Element] {
			return element.ElementWithPath[fhir.Element]{Element: e.Element, FHIRPath: e.FHIRPath}
		},
	)
}

func zipElementsAndPaths[elementT proto.Message](elements []elementT, paths []string) []element.ElementWithPath[elementT] {
	zipped := []element.ElementWithPath[elementT]{}
	for idx, ele := range elements {
		var path string
		if idx < len(paths) {
			path = paths[idx]
		}
		zipped = append(zipped, element.ElementWithPath[elementT]{Element: ele, FHIRPath: path})
	}
	return zipped
}

func Test_ExtractAll_Modifiable(t *testing.T) {
	resource := fhirtest.NewResource(t, "Patient", fhirtest.WithResourceModification(func(p *ppb.Patient) {
		p.ManagingOrganization = &dtpb.Reference{
			Reference: &dtpb.Reference_Uri{
				Uri: fhir.String("old-ref"),
			},
		}
	})).(*ppb.Patient)

	references, _ := element.ExtractAll[*dtpb.Reference](resource)

	if len(references) != 1 {
		t.Fatalf("Expected single reference")
	}

	references[0].Reference = &dtpb.Reference_Uri{
		Uri: fhir.String("new-ref"),
	}

	if got, want := resource.ManagingOrganization, references[0]; got != want {
		t.Errorf("ExtractAll() reference update failed, got '%v', want '%v'", got, want)
	}
}
