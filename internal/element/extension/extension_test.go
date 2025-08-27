package extension_test

import (
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	prpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/person_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/internal/element/canonical"
	"github.com/verily-src/fhirpath-go/internal/element/extension"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirtest"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func isValueXType(element fhir.Element) bool {
	switch element.(type) {
	case *dtpb.Base64Binary,
		*dtpb.Boolean,
		*dtpb.Canonical,
		*dtpb.Code,
		*dtpb.Date,
		*dtpb.DateTime,
		*dtpb.Decimal,
		*dtpb.Id,
		*dtpb.Instant,
		*dtpb.Integer,
		*dtpb.Markdown,
		*dtpb.Oid,
		*dtpb.PositiveInt,
		*dtpb.String,
		*dtpb.Time,
		*dtpb.UnsignedInt,
		*dtpb.Uri,
		*dtpb.Url,
		*dtpb.Uuid,
		*dtpb.Address,
		*dtpb.Age,
		*dtpb.Annotation,
		*dtpb.Attachment,
		*dtpb.CodeableConcept,
		*dtpb.Coding,
		*dtpb.ContactPoint,
		*dtpb.Count,
		*dtpb.Distance,
		*dtpb.Duration,
		*dtpb.HumanName,
		*dtpb.Identifier,
		*dtpb.Money,
		*dtpb.Period,
		*dtpb.Quantity,
		*dtpb.Range,
		*dtpb.Ratio,
		*dtpb.Reference,
		*dtpb.SampledData,
		*dtpb.Signature,
		*dtpb.Timing,
		*dtpb.ContactDetail,
		*dtpb.Contributor,
		*dtpb.DataRequirement,
		*dtpb.Expression,
		*dtpb.ParameterDefinition,
		*dtpb.RelatedArtifact,
		*dtpb.TriggerDefinition,
		*dtpb.UsageContext,
		*dtpb.Dosage:
		return true
	}
	return false
}

func TestRoundTrip(t *testing.T) {
	for name, element := range fhirtest.Elements {
		if !isValueXType(element) {
			t.Skip()
		}

		t.Run(name, func(t *testing.T) {
			ext, err := extension.FromElement("foo", element)
			if err != nil {
				t.Fatalf("RoundTrip(%v): got unexpected error %v", name, err)
			}

			got := extension.Unwrap(ext)

			if got, want := got, element; got != want {
				t.Errorf("RoundTrip(%v): got %v, want %v", name, got, want)
			}
		})
	}
}

func TestNew_Base64Binary_ReturnsExtension(t *testing.T) {
	const url = "http://example.com"
	input := fhir.Base64Binary([]byte{0xde, 0xad, 0xbe, 0xef})
	want := &dtpb.Extension{
		Url: fhir.URI(url),
		Value: &dtpb.Extension_ValueX{
			Choice: &dtpb.Extension_ValueX_Base64Binary{
				Base64Binary: input,
			},
		},
	}

	got := extension.New(url, input)

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("Extension(Base64Binary): (-got,+want):\n%v", diff)
	}
}

func TestNew_Boolean_ReturnsExtension(t *testing.T) {
	const url = "http://example.com"
	input := fhir.Boolean(true)
	want := &dtpb.Extension{
		Url: fhir.URI(url),
		Value: &dtpb.Extension_ValueX{
			Choice: &dtpb.Extension_ValueX_Boolean{
				Boolean: input,
			},
		},
	}

	got := extension.New(url, input)

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("Extension(Boolean): (-got,+want):\n%v", diff)
	}
}

func TestNew_Canonical_ReturnsExtension(t *testing.T) {
	const url = "http://example.com"
	input := canonical.New(url)
	want := &dtpb.Extension{
		Url: fhir.URI(url),
		Value: &dtpb.Extension_ValueX{
			Choice: &dtpb.Extension_ValueX_Canonical{
				Canonical: input,
			},
		},
	}

	got := extension.New(url, input)

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("Extension(Boolean): (-got,+want):\n%v", diff)
	}
}
func TestNew_Address_ReturnsExtension(t *testing.T) {
	const url = "http://example.com"
	input := &dtpb.Address{
		Text: fhir.String("hello world"),
	}
	want := &dtpb.Extension{
		Url: fhir.URI(url),
		Value: &dtpb.Extension_ValueX{
			Choice: &dtpb.Extension_ValueX_Address{
				Address: input,
			},
		},
	}

	got := extension.New(url, input)

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("Extension(Address): (-got,+want):\n%v", diff)
	}
}

func TestNew_Id_ReturnsExtension(t *testing.T) {
	const url = "http://example.com"
	input := fhir.String("hello world")
	want := &dtpb.Extension{
		Url: fhir.URI(url),
		Value: &dtpb.Extension_ValueX{
			Choice: &dtpb.Extension_ValueX_StringValue{
				StringValue: input,
			},
		},
	}

	got := extension.New(url, input)

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("Extension(String): (-got,+want):\n%v", diff)
	}
}

func TestUnwrap_NilInput_ReturnsNil(t *testing.T) {
	got := extension.Unwrap(nil)

	if got != nil {
		t.Errorf("Unwrap: got %v, want nil", got)
	}
}

func TestUnwrap(t *testing.T) {
	for name, element := range fhirtest.Elements {
		t.Run(name, func(t *testing.T) {
			ext, err := extension.FromElement("url", element)
			if err != nil {
				// Only test using elements that are valid for extensions.
				t.Skip()
			}

			got := extension.Unwrap(ext)

			want := element
			if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
				t.Errorf("Unwrap(%v): (-got,+want):\n%v", name, diff)
			}
		})
	}
}

func TestFromElement_NilInput_ReturnsError(t *testing.T) {
	_, err := extension.FromElement("url", nil)

	if err == nil {
		t.Errorf("Extension: got nil, want err")
	}
}

func TestFromElement(t *testing.T) {
	const url = "http://example.com"
	testCases := []struct {
		name    string
		input   fhir.Element
		want    *dtpb.Extension
		wantErr error
	}{
		{
			name:    "Nil",
			input:   nil,
			want:    nil,
			wantErr: cmpopts.AnyError,
		}, {
			name:  "String",
			input: fhir.String("hello world"),
			want: &dtpb.Extension{
				Url: fhir.URI(url),
				Value: &dtpb.Extension_ValueX{
					Choice: &dtpb.Extension_ValueX_StringValue{
						StringValue: fhir.String("hello world"),
					},
				},
			},
			wantErr: nil,
		}, {
			name:    "Extension",
			input:   &dtpb.Extension{},
			want:    nil,
			wantErr: extension.ErrInvalidValueX,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := extension.FromElement(url, tc.input)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("FromElement(%v): got err %v, want %v", tc.name, got, want)
			}
			if got, want := got, tc.want; !cmp.Equal(got, want, protocmp.Transform()) {
				t.Errorf("FromElement(%v): got %v, want %v", tc.name, got, want)
			}
		})
	}
}

func TestOverwrite_WithNil_Panics(t *testing.T) {
	defer func() { _ = recover() }()

	extension.Overwrite(nil)

	t.Errorf("Overwrite: expected panic")
}

func TestOverwrite_WithExtendableResources_ExtendsResource(t *testing.T) {
	const url = "url"
	ext := extension.New(url, fhir.Boolean(true))

	for name, object := range fhirtest.DomainResources {
		t.Run(name, func(t *testing.T) {
			object := proto.Clone(object).(fhir.DomainResource)
			extension.Overwrite(object, ext)

			if len(object.GetExtension()) != 1 {
				t.Fatalf("Overwrite: got len %v, want 1", len(object.GetExtension()))
			}
			got := object.GetExtension()[0]
			if !proto.Equal(got, ext) {
				t.Errorf("Overwrite: got extension %v, want %v", got, ext)
			}
		})
	}
}

func TestAppendInto_WithNil_Panics(t *testing.T) {
	defer func() { _ = recover() }()

	extension.AppendInto(nil)

	t.Errorf("AppendInto: expected panic")
}

func TestAppendInto_WithResources_ExtendsResource(t *testing.T) {
	const url = "url"
	ext := extension.New(url, fhir.Boolean(true))

	for name, object := range fhirtest.DomainResources {
		t.Run(name, func(t *testing.T) {
			object := proto.Clone(object).(fhir.DomainResource)
			extension.Overwrite(object, extension.New(url, fhir.String("Some other value")))

			extension.AppendInto(object, ext)

			if len(object.GetExtension()) != 2 {
				t.Fatalf("AppendInto: got len %v, want 1", len(object.GetExtension()))
			}
			got := object.GetExtension()[1]
			if !proto.Equal(got, ext) {
				t.Errorf("AppendInto: got extension %v, want %v", got, ext)
			}
		})
	}
}

func TestUpsert_WithNil_Panics(t *testing.T) {
	defer func() { _ = recover() }()

	extension.AppendInto(nil)

	t.Errorf("AppendInto: expected panic")
}

func TestUpsert_WithResources_ExtendsResource(t *testing.T) {
	const (
		urlA = "urlA"
		urlB = "urlB"
	)
	ext := extension.New(urlA, fhir.Boolean(true))

	for name, object := range fhirtest.DomainResources {
		t.Run(name, func(t *testing.T) {
			object := proto.Clone(object).(fhir.DomainResource)
			extension.Overwrite(object, extension.New(urlB, fhir.String("Some other value")))

			extension.Upsert(object, ext)

			if len(object.GetExtension()) != 2 {
				t.Fatalf("Upsert: got len %v, want 2", len(object.GetExtension()))
			}
			got := object.GetExtension()[1]
			if !proto.Equal(got, ext) {
				t.Errorf("AppendInto: got extension %v, want %v", got, ext)
			}
		})
	}
}

func TestUpsert_WithResources_ModifiesResource(t *testing.T) {
	const urlA = "urlA"

	ext := extension.New(urlA, fhir.Boolean(true))

	for name, object := range fhirtest.DomainResources {
		t.Run(name, func(t *testing.T) {
			object := proto.Clone(object).(fhir.DomainResource)
			extension.Overwrite(object, extension.New(urlA, fhir.String("Some other value")))

			extension.Upsert(object, ext)

			if len(object.GetExtension()) != 1 {
				t.Fatalf("Upsert: got len %v, want 1", len(object.GetExtension()))
			}
			got := object.GetExtension()[0]
			if !proto.Equal(got, ext) {
				t.Errorf("AppendInto: got extension %v, want %v", got, ext)
			}
		})
	}
}

func TestSetByURL_WithNil_Panics(t *testing.T) {
	defer func() { _ = recover() }()

	extension.SetByURL(nil, "some-url", fhir.Boolean(true))

	t.Errorf("AppendInto: expected panic")
}

func TestSetByURL_WithResources_ModifiesResource(t *testing.T) {
	const (
		urlA = "urlA"
		urlB = "urlB"
	)
	extB := extension.New(urlB, fhir.String("Some B value"))
	extA := extension.New(urlA, fhir.String("Some A value"))
	wantExtensions := []*dtpb.Extension{extB, extension.New(urlA, fhir.Boolean(true)), extension.New(urlA, fhir.Boolean(false))}
	for name, object := range fhirtest.DomainResources {
		t.Run(name, func(t *testing.T) {
			object := proto.Clone(object).(fhir.DomainResource)
			extension.Overwrite(object, extA)
			extension.AppendInto(object, extB)

			extension.SetByURL(object, urlA, fhir.Boolean(true), fhir.Boolean(false))

			if len(object.GetExtension()) != 3 {
				t.Fatalf("Upsert: got len %v, want 3", len(object.GetExtension()))
			}
			gotExtensions := object.GetExtension()
			if diff := cmp.Diff(gotExtensions, wantExtensions, protocmp.Transform()); diff != "" {
				t.Errorf("AppendInto: (-got, +want) %v", diff)
			}
		})
	}
}

func TestClear_DoesNothing(t *testing.T) {
	extension.Clear(nil)

	// There really isn't anything to assert on here; there is no input or reaction.
	// There is no crash either.
}

func TestClear_ResourceHasExtensions_RemovesExtensions(t *testing.T) {
	const url = "http://example.com"
	for name, resource := range fhirtest.DomainResources {
		t.Run(name, func(t *testing.T) {
			got := proto.Clone(resource).(fhir.DomainResource)
			want := proto.Clone(resource).(fhir.DomainResource)
			extension.Overwrite(got,
				extension.New(url, fhir.String("hello world")),
				extension.New(url, fhir.Boolean(true)),
			)

			extension.Clear(got)

			if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
				t.Errorf("Clear(%v): (-got,+want):\n%v", name, diff)
			}
		})
	}
}

func TestFindByURL(t *testing.T) {
	const (
		urlA = "http://example.com/extA"
		urlB = "http://example.com/extB"
	)
	extA := extension.New(urlA, fhir.String("valueA"))
	extB := extension.New(urlB, fhir.String("valueB"))

	testCases := []struct {
		name       string
		extensions []*dtpb.Extension
		system     string
		want       []*dtpb.Extension
	}{
		{
			name:       "nil slice",
			extensions: nil,
			system:     urlA,
			want:       [](*dtpb.Extension)(nil),
		},
		{
			name:       "empty slice",
			extensions: []*dtpb.Extension{},
			system:     urlA,
			want:       [](*dtpb.Extension)(nil),
		},
		{
			name:       "no match",
			extensions: []*dtpb.Extension{extA},
			system:     urlB,
			want:       [](*dtpb.Extension)(nil),
		},
		{
			name:       "single match",
			extensions: []*dtpb.Extension{extA},
			system:     urlA,
			want:       []*dtpb.Extension{extA},
		},
		{
			name:       "multiple extensions, match is first",
			extensions: []*dtpb.Extension{extA, extB},
			system:     urlA,
			want:       []*dtpb.Extension{extA},
		},
		{
			name:       "multiple extensions, contains match",
			extensions: []*dtpb.Extension{extB, extA},
			system:     urlA,
			want:       []*dtpb.Extension{extA},
		},
		{
			name:       "multiple matches",
			extensions: []*dtpb.Extension{extA, extB, extA},
			system:     urlA,
			want:       []*dtpb.Extension{extA, extA},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := extension.FindByURL(tc.extensions, tc.system)
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("FindByURL(%s): diff (-want,+got):\n%s", tc.name, diff)
			}
		})
	}
}

func TestClearByURL(t *testing.T) {
	const (
		urlA = "http://example.com/extA"
		urlB = "http://example.com/extB"
	)
	extA := extension.New(urlA, fhir.String("valueA"))
	extB := extension.New(urlB, fhir.String("valueB"))

	testCases := []struct {
		name            string
		extensions      []*dtpb.Extension
		clearExtWithUrl string
		want            []*dtpb.Extension
	}{
		{
			name:            "nil slice",
			extensions:      nil,
			clearExtWithUrl: urlA,
			want:            nil,
		},
		{
			name:            "no match",
			extensions:      []*dtpb.Extension{extA},
			clearExtWithUrl: urlB,
			want:            []*dtpb.Extension{extA},
		},
		{
			name:            "single match",
			extensions:      []*dtpb.Extension{extA},
			clearExtWithUrl: urlA,
			want:            nil, // ClearByURL removes the extension.
		},
		{
			name:            "multiple extensions, match is first",
			extensions:      []*dtpb.Extension{extA, extB},
			clearExtWithUrl: urlA,
			want:            []*dtpb.Extension{extB},
		},
		{
			name:            "multiple extensions, contains match in the middle",
			extensions:      []*dtpb.Extension{extB, extA, extB},
			clearExtWithUrl: urlA,
			want:            []*dtpb.Extension{extB, extB}, // Only the B extension remains.
		},
		{
			name:            "multiple extensions, contains match",
			extensions:      []*dtpb.Extension{extB, extA, extA},
			clearExtWithUrl: urlB,
			want:            []*dtpb.Extension{extA, extA}, // Only the A extensions remain.
		},
		{
			name:            "multiple matches",
			extensions:      []*dtpb.Extension{extA, extB, extA},
			clearExtWithUrl: urlA,
			want:            []*dtpb.Extension{extB},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := append([]*dtpb.Extension{}, tc.extensions...) // Copy to avoid modifying original slice.
			person := &prpb.Person{
				Extension: got,
			}
			extension.ClearByURL(person, tc.clearExtWithUrl)
			if diff := cmp.Diff(tc.want, person.GetExtension(), protocmp.Transform()); diff != "" {
				t.Errorf("ClearByURL(%s): diff (-want,+got):\n%s", tc.name, diff)
			}
		})
	}
}
