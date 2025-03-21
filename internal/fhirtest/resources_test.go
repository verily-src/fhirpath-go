package fhirtest_test

import (
	"regexp"
	"testing"
	"time"

	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	dpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/device_go_proto"
	drpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/document_reference_go_proto"
	listpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/list_go_proto"
	locpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/location_go_proto"
	opb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/observation_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	pepb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/person_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/questionnaire_go_proto"
	qrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/questionnaire_response_go_proto"
	rstudpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/research_study_go_proto"
	rsubpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/research_subject_go_proto"
	vrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/verification_result_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/internal/element/reference"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirtest"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestNewResource_GeneratesIdentityInformation(t *testing.T) {
	for name := range fhirtest.Resources {
		sut := fhirtest.NewResource(t, resource.Type(name))

		if sut.GetId().GetValue() == "" {
			t.Errorf("NewResource(%v): expected random ID", name)
		}
		if sut.GetMeta().GetVersionId().GetValue() == "" {
			t.Errorf("NewResource(%v): expected random version ID", name)
		}
		if sut.GetMeta().GetLastUpdated().GetValueUs() == 0 {
			t.Errorf("NewResource(%v): expected random last-update time", name)
		}
	}
}

func TestNewResourceFromBase_PreservesInput(t *testing.T) {
	for name := range fhirtest.Resources {
		sut := fhirtest.NewResource(t, resource.Type(name))
		sutDuplicate := proto.Clone(sut).(fhir.Resource)

		fhirtest.NewResourceFromBase(t, sut,
			fhirtest.WithProtoField("meta", &dtpb.Meta{Source: fhir.URI("urn:uuid:A")}))

		if diff := cmp.Diff(sut, sutDuplicate, protocmp.Transform()); diff != "" {
			t.Errorf("NewResourceFromBase() modified input (-got, +want): %s", diff)
		}
	}
}

func TestNewResourceFromBase_AppliesOption(t *testing.T) {
	for name := range fhirtest.Resources {
		sut := fhirtest.NewResource(t, resource.Type(name))
		newMeta := &dtpb.Meta{Source: fhir.URI("urn:uuid:A")}
		wantMsut := proto.Clone(sut).(fhir.Resource)
		fhirtest.ReplaceMeta(wantMsut, newMeta)

		gotMsut := fhirtest.NewResourceFromBase(t, sut,
			fhirtest.WithProtoField("meta", newMeta))

		if diff := cmp.Diff(gotMsut, wantMsut, protocmp.Transform()); diff != "" {
			t.Errorf("NewResourceFromBase() didn't modify properly (-got, +want): %s", diff)
		}
	}
}

func TestNewResourceOf_WithPatient_CreatesPatient(t *testing.T) {
	got := fhirtest.NewResourceOf[*patient_go_proto.Patient](t)

	if got == nil {
		t.Errorf("NewResourceOf: got nil, want patient")
	}
}

func TestWithResourceModification_ChangesPatientName(t *testing.T) {
	pIn := &patient_go_proto.Patient{
		Name: []*dtpb.HumanName{
			{
				Family: fhir.String("VonBatchery"),
				Given:  []*dtpb.String{fhir.String("GivenA"), fhir.String("GivenB")},
			},
		},
	}
	pWant := &patient_go_proto.Patient{
		Name: []*dtpb.HumanName{
			{
				Family: fhir.String("Rosewater"),
				Given:  []*dtpb.String{fhir.String("GivenB")},
			},
		},
	}

	fhirtest.WithResourceModification(func(patient *patient_go_proto.Patient) {
		patient.Name[0].Family = fhir.String("Rosewater")
		patient.Name[0].Given = patient.Name[0].Given[1:]
	})(t, pIn)

	if diff := cmp.Diff(pWant, pIn, protocmp.Transform()); diff != "" {
		t.Errorf("WithResourceModification didn't produce correct patient (-got, +want): %s", diff)
	}
}

func TestWithFieldJSON_SetsJSONField(t *testing.T) {
	for name := range fhirtest.CanonicalResources {
		sut := fhirtest.NewResource(t, resource.Type(name),
			fhirtest.WithJSONField("status", `{"value": "DRAFT"}`),
		)

		// 'NewResource' would fail the test implicitly if this didn't work correctly;
		// so if we make it here, it means we have passed.
		// Since the "status" field is a different type for each resource, we can't
		// generically test this value over an iterative sequence.
		_ = sut
	}
}

func TestWithFieldJSON_SetsAccountStatusField(t *testing.T) {
	sut := fhirtest.NewResource(t, "Questionnaire",
		fhirtest.WithJSONField("status", `{"value": "DRAFT"}`),
	).(*questionnaire_go_proto.Questionnaire)

	if got, want := sut.GetStatus().GetValue(), cpb.PublicationStatusCode_DRAFT; got != want {
		t.Errorf("WithFieldJSON: got code %v, want code %v", got, want)
	}
}

func TestWithFieldCodeJSON_SetsJSONField(t *testing.T) {
	for name := range fhirtest.CanonicalResources {
		sut := fhirtest.NewResource(t, resource.Type(name),
			fhirtest.WithCodeField("status", "DRAFT"),
		)

		// 'NewResource' would fail the test implicitly if this didn't work correctly;
		// so if we make it here, it means we have passed.
		// Since the "status" field is a different type for each resource, we can't
		// generically test this value over an iterative sequence.
		_ = sut
	}
}

func TestWithFieldCodeJSON_SetsAccountStatusField(t *testing.T) {
	sut := fhirtest.NewResource(t, "Questionnaire",
		fhirtest.WithCodeField("status", "DRAFT"),
	).(*questionnaire_go_proto.Questionnaire)

	if got, want := sut.GetStatus().GetValue(), cpb.PublicationStatusCode_DRAFT; got != want {
		t.Errorf("WithFieldCodeJSON: got code %v, want code %v", got, want)
	}
}

func TestWithFieldProto_SetsField(t *testing.T) {
	want := &dtpb.String{
		Value: "test",
	}
	for name := range fhirtest.CanonicalResources {
		sut := fhirtest.NewResource(t, resource.Type(name),
			fhirtest.WithProtoField("name", want),
		).(fhir.CanonicalResource)

		if got := sut.GetName(); got != want {
			t.Errorf("WithFieldProto(%v): got '%v', wanted '%v'", name, got, want)
		}
	}
}

func TestWithRepeatedFieldProto_SetsField(t *testing.T) {
	toProtoMessage := func(in []*dtpb.Identifier) []proto.Message {
		result := make([]proto.Message, 0, len(in))
		for _, v := range in {
			result = append(result, v)
		}
		return result
	}

	want := []proto.Message{
		&dtpb.Identifier{
			Value: &dtpb.String{
				Value: "foo",
			},
			System: &dtpb.Uri{
				Value: "https://example.com/foo",
			},
		},
		&dtpb.Identifier{
			Value: &dtpb.String{
				Value: "bar",
			},
			System: &dtpb.Uri{
				Value: "https://example.com/bar",
			},
		},
	}
	for name := range fhirtest.CanonicalResources {
		sut := fhirtest.NewResource(t, resource.Type(name),
			fhirtest.WithRepeatedProtoField("identifier", want...),
		).(fhir.CanonicalResource)

		if got := sut.GetIdentifier(); !cmp.Equal(toProtoMessage(got), want, protocmp.Transform()) {
			t.Errorf("WithRepeatedFieldProto(%v): got '%v', wanted '%v'", name, got, want)
		}
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

	if got, want := id.GetSystem().GetValue(), "http://example.com/fake-id"; got != want {
		t.Errorf("%v Resource.Identifier[0].System: got %v, want %v", name, got, want)
	}

	value := id.GetValue().GetValue()
	matched, _ := regexp.MatchString("^[a-f0-9-]+$", value)
	if !matched {
		t.Errorf("%v Resource.Identifier[0].Value: got %v, expected uuid", name, value)
	}
}

// All CanonicalResources should implement GetIdentifier()
func TestWithGeneratedIdentifier_CanonicalResources(t *testing.T) {
	for name := range fhirtest.CanonicalResources {
		t.Run(name, func(t *testing.T) {
			res := fhirtest.NewResource(t, resource.Type(name),
				fhirtest.WithGeneratedIdentifier("http://example.com/fake-id"),
			).(fhir.CanonicalResource)

			ids := res.GetIdentifier()
			assertIdentifier(t, ids, name, res)
		})
	}
}

// Test that all resources can either be cast or not
func TestWithGeneratedIdentifier_AllResources(t *testing.T) {
	for name, stockRes := range fhirtest.Resources {
		t.Run(name, func(t *testing.T) {
			if _, ok := stockRes.(resource.HasGetIdentifierList); ok {
				// GetIdentifier() is list

				res := fhirtest.NewResource(t, resource.Type(name),
					fhirtest.WithGeneratedIdentifier("http://example.com/fake-id"),
				).(resource.HasGetIdentifierList)

				ids := res.GetIdentifier()
				assertIdentifier(t, ids, name, res)

			} else if _, ok := stockRes.(resource.HasGetIdentifierSingle); ok {
				// GetIdentifier() is singleton

				res := fhirtest.NewResource(t, resource.Type(name),
					fhirtest.WithGeneratedIdentifier("http://example.com/fake-id"),
				).(resource.HasGetIdentifierSingle)

				id := res.GetIdentifier()
				assertIdentifier(t, []*dtpb.Identifier{id}, name, res)
			} else {
				// not compatible with GetIdentifier() interface
				return
			}
		})
	}
}

// Test a few specific types
func TestWithGeneratedIdentifier_SpotTest(t *testing.T) {
	// resources w/ Identifier as list
	for _, name := range []string{"Patient", "DocumentReference"} {
		t.Run(name, func(t *testing.T) {
			res := fhirtest.NewResource(t, resource.Type(name),
				fhirtest.WithGeneratedIdentifier("http://example.com/fake-id"),
			).(resource.HasGetIdentifierList)

			ids := res.GetIdentifier()
			assertIdentifier(t, ids, name, res)
		})
	}

	// resources w/ Identifier as singleton
	for _, name := range []string{"Bundle", "QuestionnaireResponse"} {
		t.Run(name, func(t *testing.T) {
			res := fhirtest.NewResource(t, resource.Type(name),
				fhirtest.WithGeneratedIdentifier("http://example.com/fake-id"),
			).(resource.HasGetIdentifierSingle)

			id := res.GetIdentifier()
			assertIdentifier(t, []*dtpb.Identifier{id}, name, res)
		})
	}
}

func TestResources_ResourceHasNonEmptyID(t *testing.T) {
	for name, resource := range fhirtest.Resources {
		t.Run(name, func(t *testing.T) {
			if resource.GetId().GetValue() == "" {
				t.Errorf("Resource(%v): got empty id", name)
			}
		})
	}
}

func TestResources_ResourceHasNonEmptyVersionID(t *testing.T) {
	for name, resource := range fhirtest.Resources {
		t.Run(name, func(t *testing.T) {
			if resource.GetMeta().GetVersionId().GetValue() == "" {
				t.Errorf("Resource(%v): got empty version id", name)
			}
		})
	}
}

func TestResources_ResourceHasLastModified(t *testing.T) {
	for name, resource := range fhirtest.Resources {
		t.Run(name, func(t *testing.T) {
			if resource.GetMeta().GetLastUpdated().GetValueUs() == 0 {
				t.Errorf("Resource(%v): got unset last-update time", name)
			}
		})
	}
}

// Tests that all resources can be created using their respective test helpers
func TestNewResources(t *testing.T) {
	patientURN := "urn:uuid:patient"
	researchStudyURN := "urn:uuid:researchStudy"

	patientRef := reference.Weak(resource.Patient, patientURN)
	researchStudyRef := reference.Weak(resource.ResearchStudy, researchStudyURN)

	meta := &dtpb.Meta{
		LastUpdated: fhir.Instant(time.Now()),
		VersionId:   fhir.RandomID(),
	}
	id := fhir.RandomID()

	testCases := []struct {
		name         string
		gotResource  fhir.Resource
		wantResource fhir.Resource
	}{
		{
			name:        "Device",
			gotResource: fhirtest.NewDevice(t),
			wantResource: &dpb.Device{
				Meta: meta,
				Id:   id,
			},
		},
		{
			name:        "DocumentReference",
			gotResource: fhirtest.NewDocumentReference(t),
			wantResource: &drpb.DocumentReference{
				Status: &drpb.DocumentReference_StatusCode{
					Value: cpb.DocumentReferenceStatusCode_CURRENT,
				},
				Meta: meta,
				Id:   id,
			},
		},
		{
			name:        "List",
			gotResource: fhirtest.NewList(t),
			wantResource: &listpb.List{
				Status: &listpb.List_StatusCode{
					Value: cpb.ListStatusCode_CURRENT,
				},
				Mode: &listpb.List_ModeCode{
					Value: cpb.ListModeCode_WORKING,
				},
				Meta: meta,
				Id:   id,
			},
		},
		{
			name:        "Location",
			gotResource: fhirtest.NewLocation(t),
			wantResource: &locpb.Location{
				Meta: meta,
				Id:   id,
			},
		},
		{
			name:        "Observation",
			gotResource: fhirtest.NewObservation(t),
			wantResource: &opb.Observation{
				Status: &opb.Observation_StatusCode{
					Value: cpb.ObservationStatusCode_PRELIMINARY,
				},
				Code: fhir.CodeableConcept("my-code-text"),
				Meta: meta,
				Id:   id,
			},
		},
		{
			name:        "Patient",
			gotResource: fhirtest.NewPatient(t),
			wantResource: &ppb.Patient{
				Meta: meta,
				Id:   id,
			},
		},
		{
			name:        "Person",
			gotResource: fhirtest.NewPerson(t),
			wantResource: &pepb.Person{
				Meta: meta,
				Id:   id,
			},
		},
		{
			name:        "QuestionnaireResponse",
			gotResource: fhirtest.NewQuestionnaireResponse(t),
			wantResource: &qrpb.QuestionnaireResponse{
				Subject: patientRef,
				Status: &qrpb.QuestionnaireResponse_StatusCode{
					Value: cpb.QuestionnaireResponseStatusCode_COMPLETED,
				},
				Meta: meta,
				Id:   id,
			},
		},
		{
			name:        "ResearchStudy",
			gotResource: fhirtest.NewResearchStudy(t),
			wantResource: &rstudpb.ResearchStudy{
				Status: &rstudpb.ResearchStudy_StatusCode{
					Value: cpb.ResearchStudyStatusCode_APPROVED,
				},
				Meta: meta,
				Id:   id,
			},
		},
		{
			name:        "ResearchSubject",
			gotResource: fhirtest.NewResearchSubject(t),
			wantResource: &rsubpb.ResearchSubject{
				Individual: patientRef,
				Status: &rsubpb.ResearchSubject_StatusCode{
					Value: cpb.ResearchSubjectStatusCode_ELIGIBLE,
				},
				Study: researchStudyRef,
				Meta:  meta,
				Id:    id,
			},
		},
		{
			name:        "VerificationResult",
			gotResource: fhirtest.NewVerificationResult(t),
			wantResource: &vrpb.VerificationResult{
				Status: &vrpb.VerificationResult_StatusCode{
					Value: cpb.StatusCode_VALIDATED,
				},
				Meta: meta,
				Id:   id,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requireMatch(t, resource.Type(tc.name), tc.gotResource, tc.wantResource)
		})
	}
}
