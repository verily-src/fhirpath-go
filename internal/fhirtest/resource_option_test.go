package fhirtest_test

import (
	"testing"
	"time"

	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	dpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/device_go_proto"
	drpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/document_reference_go_proto"
	listpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/list_go_proto"
	locpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/location_go_proto"
	opb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/observation_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	pepb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/person_go_proto"
	qrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/questionnaire_response_go_proto"
	rstudpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/research_study_go_proto"
	rsubpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/research_subject_go_proto"
	vrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/verification_result_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/internal/element/reference"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirtest"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"google.golang.org/protobuf/testing/protocmp"
)

// Tests that options are correctly applied to resources
func TestResourceOptions(t *testing.T) {
	patientURN := "urn:uuid:patient"
	researchStudyURN := "urn:uuid:researchStudy"
	locationURN := "urn:uuid:location"

	patientRef := reference.Weak(resource.Patient, patientURN)
	multiPatientRefs := []*dtpb.Reference{patientRef, patientRef, patientRef}
	researchStudy := reference.Weak(resource.ResearchStudy, researchStudyURN)
	locationRef := reference.Weak(resource.Location, locationURN)

	binaryContentTitle := "Uri where the oversized resource can be found"
	binaryURI := "Binary/urn:uuid:binary"

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
			gotResource: fhirtest.NewDevice(t, fhirtest.WithStatus("ACTIVE")),
			wantResource: &dpb.Device{
				Status: &dpb.Device_StatusCode{
					Value: cpb.FHIRDeviceStatusCode_ACTIVE,
				},
				Meta: meta,
				Id:   id,
			},
		},
		{
			name: "DocumentReference",
			gotResource: fhirtest.NewDocumentReference(t,
				fhirtest.WithStatus("SUPERSEDED"),
				fhirtest.WithSubject(patientRef),
				fhirtest.WithContent(binaryContentTitle, binaryURI)),
			wantResource: &drpb.DocumentReference{
				Content: []*drpb.DocumentReference_Content{
					{
						Attachment: &dtpb.Attachment{
							Title: fhir.String(binaryContentTitle),
							Url:   fhir.URL(binaryURI),
						},
					},
				},
				Subject: patientRef,
				Status: &drpb.DocumentReference_StatusCode{
					Value: cpb.DocumentReferenceStatusCode_SUPERSEDED,
				},
				Meta: meta,
				Id:   id,
			},
		},
		{
			name: "List",
			gotResource: fhirtest.NewList(t,
				fhirtest.WithSubject(patientRef),
				fhirtest.WithStatus("RETIRED"),
				fhirtest.WithMode("SNAPSHOT"),
				fhirtest.WithEntry(reference.Weak("Location", "urn:uuid:loc1")),
				fhirtest.WithEntry(reference.Weak("Location", "urn:uuid:loc2"))),
			wantResource: &listpb.List{
				Subject: patientRef,
				Entry: []*listpb.List_Entry{
					{Item: reference.Weak("Location", "urn:uuid:loc1")},
					{Item: reference.Weak("Location", "urn:uuid:loc2")},
				},
				Status: &listpb.List_StatusCode{
					Value: cpb.ListStatusCode_RETIRED,
				},
				Mode: &listpb.List_ModeCode{
					Value: cpb.ListModeCode_SNAPSHOT,
				},
				Meta: meta,
				Id:   id,
			},
		},
		{
			name:        "Location",
			gotResource: fhirtest.NewLocation(t, fhirtest.WithPartOf(locationRef)),
			wantResource: &locpb.Location{
				PartOf: locationRef,
				Meta:   meta,
				Id:     id,
			},
		},
		{
			name:        "Observation",
			gotResource: fhirtest.NewObservation(t, fhirtest.WithSubject(patientRef), fhirtest.WithStatus("FINAL")),
			wantResource: &opb.Observation{
				Subject: patientRef,
				Status: &opb.Observation_StatusCode{
					Value: cpb.ObservationStatusCode_FINAL,
				},
				Code: fhir.CodeableConcept("my-code-text"),
				Meta: meta,
				Id:   id,
			},
		},
		{
			name:        "Patient",
			gotResource: fhirtest.NewPatient(t, fhirtest.WithHumanName("Parker", "Peter"), fhirtest.WithPatientLink(patientRef, cpb.LinkTypeCode_SEEALSO)),
			wantResource: &ppb.Patient{
				Link: []*ppb.Patient_Link{
					{
						Other: patientRef,
						Type: &ppb.Patient_Link_TypeCode{
							Value: cpb.LinkTypeCode_SEEALSO,
						},
					},
				},
				Name: []*dtpb.HumanName{{
					Family: fhir.String("Parker"),
					Given:  []*dtpb.String{fhir.String("Peter")},
				}},
				Meta: meta,
				Id:   id,
			},
		},
		{
			name:        "Person",
			gotResource: fhirtest.NewPerson(t, fhirtest.WithPersonLink(multiPatientRefs...)),
			wantResource: &pepb.Person{
				Link: []*pepb.Person_Link{
					{Target: patientRef},
					{Target: patientRef},
					{Target: patientRef},
				},
				Meta: meta,
				Id:   id,
			},
		},
		{
			name:        "QuestionnaireResponse",
			gotResource: fhirtest.NewQuestionnaireResponse(t, fhirtest.WithSubject(patientRef), fhirtest.WithStatus("COMPLETED")),
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
			gotResource: fhirtest.NewResearchStudy(t, fhirtest.WithStatus("COMPLETED")),
			wantResource: &rstudpb.ResearchStudy{
				Status: &rstudpb.ResearchStudy_StatusCode{
					Value: cpb.ResearchStudyStatusCode_COMPLETED,
				},
				Meta: meta,
				Id:   id,
			},
		},
		{
			name:        "ResearchSubject",
			gotResource: fhirtest.NewResearchSubject(t, fhirtest.WithStatus("ON_STUDY")),
			wantResource: &rsubpb.ResearchSubject{
				Individual: patientRef,
				Status: &rsubpb.ResearchSubject_StatusCode{
					Value: cpb.ResearchSubjectStatusCode_ON_STUDY,
				},
				Study: researchStudy,
				Meta:  meta,
				Id:    id,
			},
		},
		{
			name:        "VerificationResult",
			gotResource: fhirtest.NewVerificationResult(t, fhirtest.WithStatus("ATTESTED"), fhirtest.WithTarget(patientRef)),
			wantResource: &vrpb.VerificationResult{
				Status: &vrpb.VerificationResult_StatusCode{
					Value: cpb.StatusCode_ATTESTED,
				},
				Target: []*dtpb.Reference{patientRef},
				Meta:   meta,
				Id:     id,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requireMatch(t, resource.Type(tc.name), tc.gotResource, tc.wantResource)
		})
	}
}

// requireMatch requires that the provided resources match. Meta and ID fields are ignored.
func requireMatch(t *testing.T, desc resource.Type, gotRes, wantRes fhir.Resource) {
	if diff := cmp.Diff(gotRes, wantRes, protocmp.Transform(),
		protocmp.IgnoreFields(&dtpb.Meta{}, "last_updated", "version_id"),
		protocmp.IgnoreFields(&dtpb.Id{}, "value")); diff != "" {
		t.Errorf("%v resources don't match (-got, +want): %s", desc, diff)
	}
}
