package reference

import (
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/resource"
)

var patientType, _ = resource.NewType("Patient")

func TestLiteralInfoGetters_WhenEmpty(t *testing.T) {
	// WATCHOUT: This empty literal isn't realistic.
	lit := &LiteralInfo{}
	_, ok := lit.Type()
	if ok {
		t.Errorf("Expected Type not ok")
	}
	_, ok = lit.FragmentID()
	if ok {
		t.Errorf("Expected FragmentID not ok")
	}
	_, ok = lit.Identity()
	if ok {
		t.Errorf("Expected Identity not ok")
	}
	url := lit.ServiceBaseURL()
	if url != "" {
		t.Errorf("Expected ServiceBaseURL empty")
	}
	_, ok = lit.NonRESTURI()
	if ok {
		t.Errorf("Expected NonRESTURI not ok")
	}
}

func TestLiteralInfoGetters_WhenFull(t *testing.T) {
	wantFrag := "frag1"
	wantIdent := mustNewIdentity("Patient", "123", "")
	wantBaseURL := "https://my.site/fhir"
	wantNonRESTURI := "nonrest"
	// WATCHOUT: This full literal isn't realistic.
	lit := &LiteralInfo{resType: &patientType, fragmentID: &wantFrag,
		identity: wantIdent, serviceBaseURL: wantBaseURL,
		nonRESTURI: &wantNonRESTURI}

	gotType, ok := lit.Type()
	if !ok {
		t.Errorf("Expected Type ok")
	}
	if gotType != patientType {
		t.Errorf("Expected Patient type: %s", gotType)
	}

	gotFrag, ok := lit.FragmentID()
	if !ok {
		t.Errorf("Expected FragmentID ok")
	}
	if gotFrag != wantFrag {
		t.Errorf("FragmentID(): got %s, want %s", gotFrag, wantFrag)
	}

	gotIdent, ok := lit.Identity()
	if !ok {
		t.Errorf("Expected Identity ok")
	}
	if gotIdent != wantIdent {
		t.Errorf("Identity(): got %s, want %s", gotIdent, wantIdent)
	}

	gotBaseURL := lit.ServiceBaseURL()
	if gotBaseURL != wantBaseURL {
		t.Errorf("Wrong ServiceBaseURL got %s, want %s", gotBaseURL, wantBaseURL)
	}

	gotNonRESTURI, ok := lit.NonRESTURI()
	if !ok {
		t.Errorf("Expected NonRESTURI ok")
	}
	if gotNonRESTURI != wantNonRESTURI {
		t.Errorf("Wrong NonRESTURI got %s, want %s", gotNonRESTURI, wantNonRESTURI)
	}
}

func TestLiteralInfoSetters_WithServiceBaseURL(t *testing.T) {
	baseUrl := "http://fhir.my.com/scope/teststore"

	testCases := []struct {
		name       string
		startLit   *LiteralInfo
		setBaseUrl string
		wantLit    *LiteralInfo
		wantErr    error
	}{
		{
			name:       "clear serviceBaseURL",
			startLit:   &LiteralInfo{serviceBaseURL: baseUrl},
			setBaseUrl: "",
			wantLit:    &LiteralInfo{serviceBaseURL: ""},
		},
		{
			name:       "set serviceBaseURL",
			startLit:   &LiteralInfo{},
			setBaseUrl: baseUrl,
			wantLit:    &LiteralInfo{serviceBaseURL: baseUrl},
		},
		{
			name:       "bad serviceBaseURL",
			startLit:   &LiteralInfo{},
			setBaseUrl: "this is not valid URL",
			wantErr:    ErrServiceBaseURLInvalid,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotLit, gotErr := tc.startLit.WithServiceBaseURL(tc.setBaseUrl)
			if !cmp.Equal(gotErr, tc.wantErr, cmpopts.EquateErrors()) {
				t.Fatalf("SetServiceBaseURL(%s) error mismatch: got '%v', want '%v'",
					tc.name, gotErr, tc.wantErr)
			}
			if diff := cmp.Diff(gotLit, tc.wantLit, cmp.AllowUnexported(LiteralInfo{})); diff != "" {
				t.Errorf("SetServiceBaseURL(%s) literal (-got, +want)\n%s\n", tc.name, diff)
			}
		})
	}
}

func TestLiteralInfo_PreferRelativeVersionedURIString(t *testing.T) {
	baseUrl := "http://fhir.my.com/scope/teststore"

	testCases := []struct {
		name            string
		startLit        *LiteralInfo
		wantRelativeURI string
	}{
		{
			name:            "Get relativeURI with version",
			startLit:        &LiteralInfo{serviceBaseURL: baseUrl, identity: mustNewIdentity("Patient", "1234", "abc")},
			wantRelativeURI: "Patient/1234/_history/abc",
		},
		{
			name:            "Get relativeURI without version",
			startLit:        &LiteralInfo{serviceBaseURL: baseUrl, identity: mustNewIdentity("Patient", "1234", "")},
			wantRelativeURI: "Patient/1234",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotRelativeURI := tc.startLit.PreferRelativeVersionedURIString()
			if diff := cmp.Diff(gotRelativeURI, tc.wantRelativeURI); diff != "" {
				t.Errorf("relativeURI(%s) of literal (-got, +want)\n%s\n", tc.name, diff)
			}
		})
	}
}

// Tests converting a Reference -> Literal -> URI.
// This isn't strictly a round trip, but close enough.
func TestLiteralInfo_Reference_RoundTrip(t *testing.T) {
	strongRef, _ := Typed("Patient", "123")
	strongVersionedRef := &dtpb.Reference{
		Type: fhir.URI("Patient"),
		Reference: &dtpb.Reference_PatientId{
			PatientId: &dtpb.ReferenceId{Value: "123", History: fhir.ID("abc")},
		},
	}
	inconsistentStrongRef, _ := Typed("Patient", "123")
	inconsistentStrongRef.Type = fhir.URI("Person")
	testCases := []struct {
		name     string
		inputRef *dtpb.Reference
		wantLit  *LiteralInfo // vs result of calling LiteralInfoOf(inputRef)
		wantErr  error
		wantUri  string // vs result of calling gotLit.URI()
	}{
		// equivalent weak and strong relative unversioned
		{
			name:     "weak basic",
			inputRef: Weak("Patient", "Patient/123"),
			wantLit:  &LiteralInfo{resType: &patientType, identity: mustNewIdentity("Patient", "123", "")},
			wantUri:  "Patient/123",
		},
		{
			name:     "weak basic no type",
			inputRef: &dtpb.Reference{Reference: &dtpb.Reference_Uri{Uri: fhir.String("Patient/123")}},
			wantLit:  &LiteralInfo{resType: &patientType, identity: mustNewIdentity("Patient", "123", "")},
			wantUri:  "Patient/123",
		},
		{
			name:     "strong basic",
			inputRef: strongRef,
			wantLit:  &LiteralInfo{resType: &patientType, identity: mustNewIdentity("Patient", "123", "")},
			wantUri:  "Patient/123",
		},

		// equivalent weak and strong relative versioned
		{
			name:     "weak versioned",
			inputRef: Weak("Patient", "Patient/123/_history/abc"),
			wantLit:  &LiteralInfo{resType: &patientType, identity: mustNewIdentity("Patient", "123", "abc")},
			wantUri:  "Patient/123/_history/abc",
		},
		{
			name:     "strong versioned",
			inputRef: strongVersionedRef,
			wantLit:  &LiteralInfo{resType: &patientType, identity: mustNewIdentity("Patient", "123", "abc")},
			wantUri:  "Patient/123/_history/abc",
		},

		// equivalent good uri (weak-like) and explicit fragments
		{
			name:     "uri fragment",
			inputRef: &dtpb.Reference{Type: fhir.URI("Patient"), Reference: &dtpb.Reference_Uri{Uri: fhir.String("#hello")}},
			wantLit:  &LiteralInfo{resType: &patientType, fragmentID: ptrString("hello")},
			wantUri:  "#hello",
		},
		{
			name: "explicit fragment",
			inputRef: &dtpb.Reference{Type: fhir.URI("Patient"),
				Reference: &dtpb.Reference_Fragment{Fragment: fhir.String("hello")}},
			wantLit: &LiteralInfo{resType: &patientType, fragmentID: ptrString("hello")},
			wantUri: "#hello",
		},

		// various valid edge cases of fragments
		{
			name:     "empty uri fragment",
			inputRef: &dtpb.Reference{Type: fhir.URI("Patient"), Reference: &dtpb.Reference_Uri{Uri: fhir.String("#")}},
			wantLit:  &LiteralInfo{resType: &patientType, fragmentID: ptrString("")},
			wantUri:  "#",
		},
		{
			name: "empty explicit fragment",
			inputRef: &dtpb.Reference{Type: fhir.URI("Patient"),
				Reference: &dtpb.Reference_Fragment{Fragment: fhir.String("")}},
			wantLit: &LiteralInfo{resType: &patientType, fragmentID: ptrString("")},
			wantUri: "#",
		},
		{
			name:     "uri fragment no type",
			inputRef: &dtpb.Reference{Reference: &dtpb.Reference_Uri{Uri: fhir.String("#hello")}},
			wantLit:  &LiteralInfo{fragmentID: ptrString("hello")},
			wantUri:  "#hello",
		},

		// equivalent invalid uri and explicit fragments
		{
			name:     "bad uri fragment",
			inputRef: &dtpb.Reference{Type: fhir.URI("Patient"), Reference: &dtpb.Reference_Uri{Uri: fhir.String("#@@@@")}},
			wantErr:  ErrWeakInvalid,
		},
		{
			name: "bad explicit fragment",
			inputRef: &dtpb.Reference{Type: fhir.URI("Patient"),
				Reference: &dtpb.Reference_Fragment{Fragment: fhir.String("@@@@")}},
			wantErr: ErrExplicitFragmentInvalid,
		},

		{"nil ref", nil, nil, ErrNotLiteral, ""},
		{"empty", &dtpb.Reference{}, nil, ErrNotLiteral, ""},
		{"type only", &dtpb.Reference{Type: fhir.URI("Patient")}, nil, ErrNotLiteral, ""},
		{"bad type", &dtpb.Reference{Type: fhir.URI("Blah")}, nil, ErrTypeInvalid, ""},
		{"bad weak", Weak("Patient", "bogus-uri"), nil, ErrWeakInvalid, ""},
		{"weak w/frag", Weak("Patient", "Patient/124#blah"), nil, ErrWeakInvalid, ""},
		{"weak w/other type", Weak("Person", "Patient/124"), nil, ErrTypeInconsistent, ""},
		{"strong w/other type", inconsistentStrongRef, nil, ErrTypeInconsistent, ""},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotLit, gotErr := LiteralInfoOf(tc.inputRef)
			if !cmp.Equal(gotErr, tc.wantErr, cmpopts.EquateErrors()) {
				t.Fatalf("RoundTrip(%s) LiteralInfoOf() error got '%v', want '%v'",
					tc.name, gotErr, tc.wantErr)
			}
			if diff := cmp.Diff(gotLit, tc.wantLit, cmp.AllowUnexported(LiteralInfo{})); diff != "" {
				t.Errorf("RoundTrip(%s) LiteralInfoOf() mismatch (-got, +want)\n%s\n", tc.name, diff)
			}
			gotUri := gotLit.URI().GetValue()
			if gotUri != tc.wantUri {
				t.Errorf("RoundTrip(%s) URI() mismatch: got '%s', want '%s'", tc.name, gotUri, tc.wantUri)
			}
		})
	}
}

func TestLiteralInfo_URI_RoundTrip(t *testing.T) {
	baseUrl := "http://fhir.my.com/scope/teststore"
	testCases := []struct {
		name     string
		inputUri string
		wantLit  *LiteralInfo
		wantErr  error
	}{
		{"rel-uri", "Patient/1234",
			&LiteralInfo{resType: &patientType, identity: mustNewIdentity("Patient", "1234", "")}, nil},
		{"rel-vers-uri", "Patient/1234/_history/abc",
			&LiteralInfo{resType: &patientType, identity: mustNewIdentity("Patient", "1234", "abc")}, nil},
		{"abs-uri", baseUrl + "/Patient/1234",
			&LiteralInfo{resType: &patientType, serviceBaseURL: baseUrl, identity: mustNewIdentity("Patient", "1234", "")}, nil},
		{"abs-vers-uri", baseUrl + "/Patient/1234/_history/abc",
			&LiteralInfo{resType: &patientType, serviceBaseURL: baseUrl, identity: mustNewIdentity("Patient", "1234", "abc")}, nil},

		{"rel-missing-vid", "Patient/1234/_history", nil, ErrInvalidURI},
		{"rel-bad-type", "NotAPatient/1234", nil, ErrInvalidURI},
		{"rel-bad-resource-id", "Patient/@@@@", nil, ErrInvalidURI},
		{"rel-bad-version-id", "Patient/1234/_history/@@@@", nil, ErrInvalidURI},

		{"frag-normal", "#1234", &LiteralInfo{fragmentID: ptrString("1234")}, nil},
		{"frag-bad-id", "#@@@@", nil, ErrInvalidURI},

		// This is special case. Note that the frag string is empty but present (not nil).
		{"frag-container", "#", &LiteralInfo{fragmentID: ptrString("")}, nil},

		{"bad-nonlocal-frag", "Patient/1234#frag1", nil, ErrInvalidURI},
		{"bad-canonical-version", "Patient/1234|v12", nil, ErrInvalidURI},
		{"bad-history-token", "Patient/1234/nothistory/abc", nil, ErrInvalidURI},
		{"bad-abs-uri", "my.site.com/Patient/1234", nil, ErrInvalidURI},

		{"urn-uuid", "urn:uuid:1234",
			&LiteralInfo{nonRESTURI: ptrString("urn:uuid:1234")}, nil},
		{"urn-bad", "urn:", nil, ErrInvalidURI},
		{"canonical-http", "http://example.com/my-thing",
			&LiteralInfo{nonRESTURI: ptrString("http://example.com/my-thing")}, nil},
		{"canonical-bad", "http://example.com/", nil, ErrInvalidURI},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotLit, gotErr := LiteralInfoFromURI(tc.inputUri)
			if !cmp.Equal(gotErr, tc.wantErr, cmpopts.EquateErrors()) {
				t.Fatalf("RoundTrip(%s): LiteralInfoFromURI() error got '%v', want '%v'. Got LiteralInfo='%v'",
					tc.name, gotErr, tc.wantErr, gotLit)
			}
			if diff := cmp.Diff(gotLit, tc.wantLit, cmp.AllowUnexported(LiteralInfo{})); diff != "" {
				t.Errorf("RoundTrip(%s): LiteralInfoFromURI() literal (-got, +want)\n%s\n", tc.name, diff)
			}
			gotUri := gotLit.URI().GetValue()
			wantUri := tc.inputUri
			if tc.wantErr != nil {
				wantUri = ""
			}
			if gotUri != wantUri {
				t.Errorf("RoundTrip(%s): URI() got '%s', want '%s'", tc.name, gotUri, wantUri)
			}
		})
	}
}

func mustNewIdentity(resType, id, version string) *resource.Identity {
	identity, err := resource.NewIdentity(resType, id, version)
	if err != nil {
		panic(err)
	}
	return identity
}

func ptrString(s string) *string {
	return &s
}
