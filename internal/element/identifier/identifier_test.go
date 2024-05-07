package identifier_test

import (
	"net/url"
	"testing"

	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/uuid"
	"github.com/verily-src/fhirpath-go/internal/element/identifier"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"google.golang.org/protobuf/proto"
)

// New also tests the behaviors of various opts.
func TestNew(t *testing.T) {
	const system = "system"
	const value = "value"
	period := &dtpb.Period{
		Start: fhir.DateTimeNow(),
	}
	ref := &dtpb.Reference{
		Reference: &dtpb.Reference_Uri{
			Uri: fhir.String("uri"),
		},
	}
	cc := &dtpb.CodeableConcept{
		Text: fhir.String("example"),
	}
	ext := &dtpb.Extension{
		Url: fhir.URI("url"),
	}
	testCases := []struct {
		name string
		opts []identifier.Option
		want *dtpb.Identifier
	}{
		{
			name: "No options",
			want: &dtpb.Identifier{
				System: fhir.URI(system),
				Value:  fhir.String(value),
			},
		}, {
			name: "WithUse",
			opts: []identifier.Option{identifier.WithUse(identifier.UseOfficial)},
			want: &dtpb.Identifier{
				System: fhir.URI(system),
				Value:  fhir.String(value),
				Use: &dtpb.Identifier_UseCode{
					Value: identifier.UseOfficial,
				},
			},
		}, {
			name: "WithAssigner",
			opts: []identifier.Option{identifier.WithAssigner(ref)},
			want: &dtpb.Identifier{
				System:   fhir.URI(system),
				Value:    fhir.String(value),
				Assigner: ref,
			},
		}, {
			name: "WithPeriod",
			opts: []identifier.Option{identifier.WithPeriod(period)},
			want: &dtpb.Identifier{
				System: fhir.URI(system),
				Value:  fhir.String(value),
				Period: period,
			},
		}, {
			name: "WithType",
			opts: []identifier.Option{identifier.WithType(cc)},
			want: &dtpb.Identifier{
				System: fhir.URI(system),
				Value:  fhir.String(value),
				Type:   cc,
			},
		}, {
			name: "WithExtensions",
			opts: []identifier.Option{identifier.WithExtensions(ext)},
			want: &dtpb.Identifier{
				System:    fhir.URI(system),
				Value:     fhir.String(value),
				Extension: []*dtpb.Extension{ext},
			},
		}, {
			name: "IncludeExtensions",
			opts: []identifier.Option{identifier.IncludeExtensions(ext)},
			want: &dtpb.Identifier{
				System:    fhir.URI(system),
				Value:     fhir.String(value),
				Extension: []*dtpb.Extension{ext},
			},
		}, {
			name: "WithID",
			opts: []identifier.Option{identifier.WithID("id")},
			want: &dtpb.Identifier{
				System: fhir.URI(system),
				Value:  fhir.String(value),
				Id:     fhir.String("id"),
			},
		},
	}

	for _, tc := range testCases {
		got := identifier.New(value, system, tc.opts...)

		if want := tc.want; !proto.Equal(got, want) {
			t.Errorf("New(%v): got %v, want %v", tc.name, got, want)
		}
	}
}

func TestUpdate(t *testing.T) {
	const system = "system"
	const value = "value"

	ext := &dtpb.Extension{
		Url: fhir.URI("url"),
	}
	testCases := []struct {
		name string
		opts []identifier.Option
		want *dtpb.Identifier
	}{
		{
			name: "IncludeExtensions",
			opts: []identifier.Option{identifier.IncludeExtensions(ext)},
			want: &dtpb.Identifier{
				Extension: []*dtpb.Extension{ext},
			},
		}, {
			name: "WithSystemString",
			opts: []identifier.Option{identifier.WithSystemString(system)},
			want: &dtpb.Identifier{
				System: fhir.URI(system),
			},
		}, {
			name: "WithValue",
			opts: []identifier.Option{identifier.WithValue(value)},
			want: &dtpb.Identifier{
				Value: fhir.String(value),
			},
		},
	}

	for _, tc := range testCases {

		got := identifier.Update(&dtpb.Identifier{}, tc.opts...)

		if want := tc.want; !proto.Equal(got, want) {
			t.Errorf("Update(%v): got %v, want %v", tc.name, got, want)
		}
	}
}

func TestNewWithUse(t *testing.T) {
	const system = "system"
	const value = "value"

	testCases := []struct {
		name string
		fn   func(string, string, ...identifier.Option) *dtpb.Identifier
		use  identifier.Use
	}{
		{"Usual", identifier.Usual, identifier.UseUsual},
		{"Official", identifier.Official, identifier.UseOfficial},
		{"Temp", identifier.Temp, identifier.UseTemp},
		{"Secondary", identifier.Secondary, identifier.UseSecondary},
		{"Old", identifier.Old, identifier.UseOld},
	}

	for _, tc := range testCases {
		got := tc.fn(value, system)

		want := &dtpb.Identifier{
			System: fhir.URI(system),
			Value:  fhir.String(value),
			Use: &dtpb.Identifier_UseCode{
				Value: tc.use,
			},
		}
		if !proto.Equal(got, want) {
			t.Errorf("%v: got %v, want %v", tc.name, got, want)
		}
	}
}

func TestEquivalent(t *testing.T) {
	const value = "value"
	const system = "system"
	testCases := []struct {
		name string
		lhs  *dtpb.Identifier
		rhs  *dtpb.Identifier
		want bool
	}{
		{
			name: "Systems don't match",
			lhs:  identifier.New(value, system),
			rhs:  identifier.New(value, system+"1"),
			want: false,
		}, {
			name: "Values don't match",
			lhs:  identifier.New(value, system),
			rhs:  identifier.New(value+"1", system),
			want: false,
		}, {
			name: "Both system and values don't match",
			lhs:  identifier.New(value, system),
			rhs:  identifier.New(value+"1", system+"1"),
			want: false,
		}, {
			name: "Both systems and values match",
			lhs:  identifier.New(value, system),
			rhs:  identifier.New(value, system),
			want: true,
		},
	}

	for _, tc := range testCases {
		got := identifier.Equivalent(tc.lhs, tc.rhs)

		if got != tc.want {
			t.Errorf("Equivalent(%v): got %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestFindBySystem(t *testing.T) {
	want := identifier.New("value", "system")

	testCases := []struct {
		name     string
		haystack []*dtpb.Identifier
		want     *dtpb.Identifier
	}{
		{
			name:     "Input empty",
			haystack: nil,
			want:     nil,
		}, {
			name: "Input not found",
			haystack: []*dtpb.Identifier{
				identifier.New("a", "a-system"),
				identifier.New("b", "b-system"),
			},
			want: nil,
		}, {
			name: "Input is first value",
			haystack: []*dtpb.Identifier{
				want,
				identifier.New("a", "a-system"),
				identifier.New("b", "b-system"),
			},
			want: want,
		}, {
			name: "Input is last value",
			haystack: []*dtpb.Identifier{
				identifier.New("a", "a-system"),
				identifier.New("b", "b-system"),
				want,
			},
			want: want,
		}, {
			name: "Input is middle entry",
			haystack: []*dtpb.Identifier{
				identifier.New("a", "a-system"),
				want,
				identifier.New("b", "b-system"),
			},
			want: want,
		},
	}

	for _, tc := range testCases {
		got := identifier.FindBySystem(tc.haystack, want.System.Value)

		if got != tc.want {
			t.Errorf("FindBySystem(%v): got %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestGenerateIfNoneExist(t *testing.T) {

	id1 := &dtpb.Identifier{
		System: &dtpb.Uri{Value: "http://fake.com"},
		Value:  &dtpb.String{Value: "9efbf82d-7a58-4d14-bec1-63f8fda148a8"},
	}
	id2 := &dtpb.Identifier{
		Use: &dtpb.Identifier_UseCode{
			Value: codes_go_proto.IdentifierUseCode_USUAL,
		},
		System: &dtpb.Uri{Value: "urn:oid:2.16.840.1.113883.2.4.6.3"},
		Value:  &dtpb.String{Value: "12345"},
	}

	id3 := &dtpb.Identifier{
		System: &dtpb.Uri{Value: "http://example.com/fake-id"},
		Value:  &dtpb.String{Value: uuid.NewString()},
	}

	id4 := &dtpb.Identifier{
		System: &dtpb.Uri{Value: "http://fake.com"},
		Value:  &dtpb.String{Value: "foo,bar,baz|omg"},
	}

	testCases := []struct {
		name  string
		input *dtpb.Identifier
		want  string
	}{
		{"nil Identifier", nil, ""},
		{"single Identifier", id1, "identifier=" + url.QueryEscape("http://fake.com|9efbf82d-7a58-4d14-bec1-63f8fda148a8")},
		{"Identifier with use code", id2, "identifier=" + url.QueryEscape("urn:oid:2.16.840.1.113883.2.4.6.3|12345")},
		{"generated ID", id3, "identifier=" + url.QueryEscape("http://example.com/fake-id|"+id3.Value.Value)},
		{"Special chars in Identifier", id4, "identifier=" + url.QueryEscape(`http://fake.com|foo\,bar\,baz\|omg`)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := identifier.GenerateIfNoneExist(tc.input)
			if got != tc.want {
				t.Errorf("%#v: Bad If-None-Exist:\n  got  %#v\n  want %#v", tc.name, got, tc.want)
			}
		})
	}
}
