package reference_test

import (
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	acpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/account_go_proto"
	appb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/appointment_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/element/reference"
	"github.com/verily-src/fhirpath-go/internal/fhirtest"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"google.golang.org/protobuf/testing/protocmp"
)

func Test_ExtractAll_ValuesCorrect(t *testing.T) {
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
		name       string
		resource   fhir.Resource
		references []*dtpb.Reference
	}{
		{
			name:     "No reference",
			resource: fhirtest.NewResource(t, "Patient"),
		},
		{
			name: "Single reference",
			resource: fhirtest.NewResource(t, "Patient", fhirtest.WithResourceModification(func(p *ppb.Patient) {
				p.ManagingOrganization = refUri
			})),
			references: []*dtpb.Reference{refUri},
		},
		{
			name: "Multiple references",
			resource: fhirtest.NewResource(t, "Appointment", fhirtest.WithResourceModification(func(a *appb.Appointment) {
				a.Participant = []*appb.Appointment_Participant{
					{
						Type:  []*dtpb.CodeableConcept{fhir.CodeableConcept("", fhir.Coding("systest", "code"))},
						Actor: refUri,
					},
					{
						Actor: refRelatedPersonId,
					},
				}
			})),
			references: []*dtpb.Reference{refRelatedPersonId, refUri},
		},
		{
			name: "Repeated field references",
			resource: fhirtest.NewResource(t, "Account", fhirtest.WithResourceModification(func(a *acpb.Account) {
				a.Subject = []*dtpb.Reference{refRelatedPersonId, refUri}
			})),
			references: []*dtpb.Reference{refRelatedPersonId, refUri},
		},
		{
			name: "Repeated identical references",
			resource: fhirtest.NewResource(t, "Account", fhirtest.WithResourceModification(func(a *acpb.Account) {
				a.Subject = []*dtpb.Reference{refRelatedPersonId, refRelatedPersonId}
			})),
			references: []*dtpb.Reference{refRelatedPersonId, refRelatedPersonId},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			gotReferences, _ := reference.ExtractAll(testCase.resource)

			opts := []cmp.Option{
				protocmp.Transform(),
				cmpopts.SortSlices(func(a *dtpb.Reference, b *dtpb.Reference) bool { return a.String() < b.String() }),
				cmpopts.EquateEmpty(),
			}
			if !cmp.Equal(testCase.references, gotReferences, opts...) {
				t.Errorf("ExtractAll(): got '%v', want '%v'", gotReferences, testCase.references)
			}
		})
	}
}

func Test_ExtractAll_Modifiable(t *testing.T) {
	resource := fhirtest.NewResource(t, "Patient", fhirtest.WithResourceModification(func(p *ppb.Patient) {
		p.ManagingOrganization = &dtpb.Reference{
			Reference: &dtpb.Reference_Uri{
				Uri: fhir.String("old-ref"),
			},
		}
	})).(*ppb.Patient)

	references, _ := reference.ExtractAll(resource)

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

func canonicalReference(t resource.Type, ref string) *dtpb.Reference {
	return &dtpb.Reference{
		Type: fhir.URI(t.String()),
		Reference: &dtpb.Reference_Uri{
			Uri: fhir.String(ref),
		},
	}
}

func TestCanonical(t *testing.T) {
	testCases := []struct {
		name         string
		resourceType resource.Type
		reference    string
		wantError    error
	}{
		{"canonical reference", "Questionnaire", "https://example.com/questionnaire", nil},
		{"invalid resource type", "NotAResource", "https://example.com/questionnaire", reference.ErrNotResource},
		{"resource reference errors", "Patient", "Patient/123", reference.ErrNotCanonicalResource},
		{"bundle reference errors", "GuidanceResponse", "urn:uuid:5a17b7c2-e01c-4bc7-b973-31d4156b11d7", reference.ErrNotCanonicalResource},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ref, err := reference.Canonical(tc.resourceType, tc.reference)

			if got, want := err, tc.wantError; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("Canonical(%s) error got %v, want %v", tc.name, got, want)
			}
			if tc.wantError == nil {
				wantRef := canonicalReference(tc.resourceType, tc.reference)
				got, want := ref, wantRef
				if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
					t.Errorf("Canonical(%s) (-got, +want)\n%s\n", tc.name, diff)
				}
			}
		})
	}
}

func TestWeak(t *testing.T) {
	testCases := []struct {
		name         string
		resourceType resource.Type
		reference    string
	}{
		{"canonical reference", "Questionnaire", "https://example.com/questionnaire"},
		{"resource reference", "Patient", "Patient/123"},
		{"bundle reference", "GuidanceResponse", "urn:uuid:5a17b7c2-e01c-4bc7-b973-31d4156b11d7"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wantRef := canonicalReference(tc.resourceType, tc.reference)

			ref := reference.Weak(tc.resourceType, tc.reference)

			got, want := ref, wantRef
			if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
				t.Errorf("Weak(%s) (-got, +want)\n%s\n", tc.name, diff)
			}
		})
	}
}

func TestTyped(t *testing.T) {
	testCases := []struct {
		name         string
		resourceType resource.Type
		resourceId   string
		wantError    error
		want         *dtpb.Reference
	}{
		{"resource reference", "Patient", "test-patient-id", nil, &dtpb.Reference{
			Reference: &dtpb.Reference_PatientId{PatientId: &dtpb.ReferenceId{Value: "test-patient-id"}},
			Type:      &dtpb.Uri{Value: "Patient"},
		}},
		{"invalid ref type", "InvalidResource", "5a17b7c2-e01c-4bc7-b973-31d4156b11d7", reference.ErrStrongConversion, nil},
		{"missing ref id", "InvalidResource", "", reference.ErrStrongConversion, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ref, err := reference.Typed(tc.resourceType, tc.resourceId)

			if !cmp.Equal(err, tc.wantError, cmpopts.EquateErrors()) {
				t.Fatalf("Typed(%s) error got %v, want %v", tc.name, err, tc.wantError)
			}

			if tc.wantError == nil {
				got, want := ref, tc.want
				if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
					t.Errorf("Typed(%s) (-got, +want)\n%s\n", tc.name, diff)
				}
			}
		})
	}
}

func TestLogical(t *testing.T) {
	resourceType := resource.Patient
	identifierSystem := "IdentifierSystem"
	identifierValue := "IdentifierValue"
	reference := reference.Logical(resourceType, identifierSystem, identifierValue)

	if diff := cmp.Diff(reference.GetType().GetValue(), resourceType.String()); diff != "" {
		t.Errorf("logicalReference.type (-got, +want)\n%s\n", diff)
	}
	if diff := cmp.Diff(reference.GetIdentifier().GetSystem().GetValue(), identifierSystem); diff != "" {
		t.Errorf("logicalReference.identifier.system (-got, +want)\n%s\n", diff)
	}
	if diff := cmp.Diff(reference.GetIdentifier().GetValue().GetValue(), identifierValue); diff != "" {
		t.Errorf("logicalReference.identifier.value (-got, +want)\n%s\n", diff)
	}
}

func TestLogicalFromIdentifier(t *testing.T) {
	resourceType := resource.Patient
	identifierSystem := "IdentifierSystem"
	identifierValue := "IdentifierValue"
	identifier := fhir.Identifier(identifierSystem, identifierValue)
	identifierCc := fhir.CodeableConcept("", fhir.Coding("system", "code"))
	identifier.Type = identifierCc

	reference := reference.LogicalFromIdentifier(resourceType, identifier)
	if diff := cmp.Diff(reference.GetType().GetValue(), resourceType.String()); diff != "" {
		t.Errorf("logicalReference.type (-got, +want)\n%s\n", diff)
	}
	if diff := cmp.Diff(reference.GetIdentifier().GetSystem().GetValue(), identifierSystem); diff != "" {
		t.Errorf("logicalReference.identifier.system (-got, +want)\n%s\n", diff)
	}
	if diff := cmp.Diff(reference.GetIdentifier().GetValue().GetValue(), identifierValue); diff != "" {
		t.Errorf("logicalReference.identifier.value (-got, +want)\n%s\n", diff)
	}

	actualCc := reference.GetIdentifier().GetType()
	if diff := cmp.Diff(identifierCc, actualCc, protocmp.Transform()); diff != "" {
		t.Errorf("unexpected proto diff:\n%v", diff)
	}
}

func TestTypedFromResource(t *testing.T) {
	patient := &ppb.Patient{Id: fhir.ID("1234"),
		Meta: &dtpb.Meta{VersionId: fhir.ID("abcd")}, // ignored
	}
	got, err := reference.TypedFromResource(patient)
	if err != nil {
		t.Fatalf("TypedFromResource failed: %v", err)
	}
	want := &dtpb.Reference{Type: fhir.URI("Patient"), Reference: &dtpb.Reference_PatientId{
		PatientId: &dtpb.ReferenceId{
			Value: "1234",
		},
	}}
	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("TypedFromResource (-got, +want)\n%s\n", diff)
	}
}

func TestTypedFromResource_WhenMissingId(t *testing.T) {
	patient := &ppb.Patient{} // no ID
	_, err := reference.TypedFromResource(patient)
	wantErr := reference.ErrStrongConversion
	if !cmp.Equal(err, wantErr, cmpopts.EquateErrors()) {
		t.Errorf("TypedFromResource error got [%v], want [%v]", err, wantErr)
	}
}

func TestTypedFromIdentity(t *testing.T) {
	testCases := []struct {
		name          string
		inputIdentity *resource.Identity
		wantRef       *dtpb.Reference
	}{
		{
			name:          "unversioned",
			inputIdentity: mustNewIdentity("Patient", "1234", ""),
			wantRef: &dtpb.Reference{
				Type:      fhir.URI("Patient"),
				Reference: &dtpb.Reference_PatientId{PatientId: &dtpb.ReferenceId{Value: "1234"}},
			},
		},
		{
			name:          "versioned",
			inputIdentity: mustNewIdentity("Patient", "1234", "abc"),
			wantRef: &dtpb.Reference{
				Type: fhir.URI("Patient"),
				Reference: &dtpb.Reference_PatientId{
					PatientId: &dtpb.ReferenceId{Value: "1234", History: fhir.ID("abc")}},
			},
		},
		{
			name:          "missing-id",
			inputIdentity: mustNewIdentity("Patient", "", ""),
			wantRef: &dtpb.Reference{
				Type:      fhir.URI("Patient"),
				Reference: &dtpb.Reference_PatientId{PatientId: &dtpb.ReferenceId{Value: ""}},
			},
		},
		{
			name:          "missing-id with version",
			inputIdentity: mustNewIdentity("Patient", "", "abc"),
			wantRef: &dtpb.Reference{
				Type: fhir.URI("Patient"),
				Reference: &dtpb.Reference_PatientId{
					PatientId: &dtpb.ReferenceId{Value: "", History: fhir.ID("abc")}},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotRef := reference.TypedFromIdentity(tc.inputIdentity)
			if diff := cmp.Diff(gotRef, tc.wantRef, protocmp.Transform()); diff != "" {
				t.Errorf("TypedFromIdentity(%s) (-got, +want)\n%s\n", tc.name, diff)
			}
		})
	}
}

func TestWeakRelativeVersioned(t *testing.T) {
	testCases := []struct {
		name      string
		inputRes  fhir.Resource
		wantRef   *dtpb.Reference
		wantError error
	}{
		{"normal",
			&ppb.Patient{Id: fhir.ID("1234"), Meta: &dtpb.Meta{VersionId: fhir.ID("abcd")}},
			&dtpb.Reference{Type: fhir.URI("Patient"),
				Reference: &dtpb.Reference_Uri{Uri: fhir.String("Patient/1234/_history/abcd")},
			},
			nil},
		{"missing-resource-id", &ppb.Patient{}, nil, reference.ErrNoResourceID},
		{"missing-version-id", &ppb.Patient{Id: fhir.ID("1234")}, nil, reference.ErrNoResourceVersion},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := reference.WeakRelativeVersioned(tc.inputRes)
			if !cmp.Equal(err, tc.wantError, cmpopts.EquateErrors()) {
				t.Errorf("WeakRelativeVersioned(%s) error got [%v], want [%v]", tc.name, err, tc.wantError)
			}
			if diff := cmp.Diff(got, tc.wantRef, protocmp.Transform()); diff != "" {
				t.Errorf("WeakRelativeVersioned(%s) (-got, +want)\n%s\n", tc.name, diff)
			}
		})
	}
}

func mustNewIdentity(resourceType, id, versionID string) *resource.Identity {
	identity, err := resource.NewIdentity(resourceType, id, versionID)
	if err != nil {
		panic(err)
	}
	return identity
}

func TestIs(t *testing.T) {
	const (
		system      = "FooSystem"
		systemValue = "value"
		id          = "3591dd62-fdcc-438f-882d-50540d3f5c18"
		url         = "https://example.com"
	)
	testCases := []struct {
		name  string
		left  *dtpb.Reference
		right *dtpb.Reference
		want  bool
	}{
		{
			name:  "Two identical logical references",
			left:  reference.Logical(resource.Account, system, systemValue),
			right: reference.Logical(resource.Account, system, systemValue),
			want:  true,
		}, {
			name:  "Two different logical references",
			left:  reference.Logical(resource.Account, system, systemValue),
			right: reference.Logical(resource.Account, system, systemValue+"-different"),
			want:  false,
		}, {
			name:  "Two identical literal references",
			left:  mustNewLiteral(resource.Patient, id),
			right: mustNewLiteral(resource.Patient, id),
			want:  true,
		}, {
			name:  "Two different literal references",
			left:  mustNewLiteral(resource.Patient, id),
			right: mustNewLiteral(resource.Patient, id+"1"),
			want:  false,
		}, {
			name:  "Two identical literal URLs",
			left:  reference.Weak(resource.Account, url),
			right: reference.Weak(resource.Account, url),
			want:  true,
		}, {
			name:  "Two different literal URLs",
			left:  reference.Weak(resource.Account, url),
			right: reference.Weak(resource.Account, url+"/something-else"),
			want:  false,
		}, {
			name:  "Left unknown reference type",
			left:  reference.Weak(resource.Account, url),
			right: mustNewLiteral(resource.Patient, id),
			want:  false,
		}, {
			name:  "Right unknown reference type",
			left:  mustNewLiteral(resource.Patient, id),
			right: reference.Weak(resource.Account, url),
			want:  false,
		}, {
			name:  "Left logical, right literal reference",
			left:  mustNewLiteral(resource.Patient, id),
			right: reference.Logical(resource.Account, system, systemValue),
			want:  false,
		}, {
			name:  "Left literal, right logical reference",
			left:  reference.Logical(resource.Account, system, systemValue),
			right: mustNewLiteral(resource.Patient, id),
			want:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := reference.Is(tc.left, tc.right)

			if got != tc.want {
				t.Errorf("Is(%v): got %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

func mustNewLiteral(res resource.Type, id string) *dtpb.Reference {
	got, err := reference.Typed(res, id)
	if err != nil {
		panic(err)
	}
	return got
}
