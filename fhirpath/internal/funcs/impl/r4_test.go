package impl_test

import (
	"errors"
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr/exprtest"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs/impl"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/element/extension"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirtest"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestExtension_ValidInput(t *testing.T) {
	extURL := "1234"
	patient := fhirtest.NewResourceOf[*ppb.Patient](t)
	ext := extension.New(extURL, fhir.String("some value"))
	patient.Extension = append(patient.Extension, &dtpb.Extension{})
	patient.Extension = append(patient.Extension, ext, ext)
	testCases := []struct {
		name  string
		input system.Collection
		arg   string
		want  system.Collection
	}{
		{
			name:  "empty input result in empty output",
			input: system.Collection{},
			arg:   "some-url",
			want:  nil,
		},
		{
			name: "entries have extensions not matched by url",
			input: system.Collection{
				fhirtest.NewResourceOf[*ppb.Patient](t),
			},
			arg:  "some-url",
			want: nil,
		},
		{
			name: "input does not have extension field",
			input: system.Collection{
				system.String("hello world"),
			},
			arg:  "some-url",
			want: nil,
		},
		{
			name:  "entries have extensions matched by url",
			input: system.Collection{patient},
			arg:   extURL,
			want:  system.Collection{ext, ext},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Extension(&expr.Context{}, tc.input, exprtest.Return(system.String(tc.arg)))
			if err != nil {
				t.Fatalf("Extension function returned unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Extension function returned unexpected diff (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestExtension_InvalidInput_RaisesError(t *testing.T) {
	testErr := errors.New("test error")
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		wantErr error
	}{
		{
			name:  "too many arguments",
			input: system.Collection{},
			args: []expr.Expression{
				exprtest.Return(system.String("")),
				exprtest.Return(system.String("")),
			},
			wantErr: impl.ErrWrongArity,
		}, {
			name:    "too few arguments",
			input:   system.Collection{},
			args:    []expr.Expression{},
			wantErr: impl.ErrWrongArity,
		}, {
			name:    "invalid argument type",
			input:   system.Collection{},
			args:    []expr.Expression{exprtest.Return(system.Integer(42))},
			wantErr: cmpopts.AnyError,
		}, {
			name:    "argument errors",
			input:   system.Collection{},
			args:    []expr.Expression{exprtest.Error(testErr)},
			wantErr: testErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := impl.Extension(&expr.Context{}, tc.input, tc.args...)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("Extension(%v): got err %v, want err %v", tc.name, got, want)
			}
		})
	}
}

func TestHasValue(t *testing.T) {
	deceased := &ppb.Patient_DeceasedX{
		Choice: &ppb.Patient_DeceasedX_Boolean{
			Boolean: fhir.Boolean(true),
		},
	}

	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:  "non-singleton collection returns false",
			input: system.Collection{fhir.String("a"), fhir.String("b")},
			want:  system.Collection{system.Boolean(false)},
		},
		{
			name:  "empty collection returns false",
			input: system.Collection{},
			want:  system.Collection{system.Boolean(false)},
		},
		{
			name:  "non-primitive type returns false",
			input: system.Collection{coding[0]},
			want:  system.Collection{system.Boolean(false)},
		},
		{
			name:  "primitive with string value returns true",
			input: system.Collection{fhir.String("value")},
			want:  system.Collection{system.Boolean(true)},
		},
		{
			name:  "primitive with boolean value returns true",
			input: system.Collection{fhir.Boolean(true)},
			want:  system.Collection{system.Boolean(true)},
		},
		{
			name:  "primitive with non-zero integer value returns true",
			input: system.Collection{fhir.Integer(123)},
			want:  system.Collection{system.Boolean(true)},
		},
		{
			name:  "primitive as polymorphic oneOf type returns true",
			input: system.Collection{deceased},
			want:  system.Collection{system.Boolean(true)},
		},
		{
			name: "primitive with value and extension returns true",
			input: system.Collection{&dtpb.String{
				Value: "hello",
				Extension: []*dtpb.Extension{
					{Url: &dtpb.Uri{Value: "http://example.com"}},
				},
			}},
			want: system.Collection{system.Boolean(true)},
		},
		{
			name: "primitive with only extension returns false",
			input: system.Collection{&dtpb.String{
				Extension: []*dtpb.Extension{
					{Url: &dtpb.Uri{Value: "http://example.com"}},
				},
			}},
			want: system.Collection{system.Boolean(false)},
		},
		{
			name: "primitive with only id returns false",
			input: system.Collection{&dtpb.String{
				Id: fhir.String("some-id"),
			}},
			want: system.Collection{system.Boolean(false)},
		},
		{
			name: "primitive with id and extension but no value returns false",
			input: system.Collection{&dtpb.String{
				Id: fhir.String("some-id"),
				Extension: []*dtpb.Extension{
					{Url: &dtpb.Uri{Value: "http://example.com"}},
				},
			}},
			want: system.Collection{system.Boolean(false)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.HasValue(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Fatalf("HasValue() error = %v, wantErr %v", err, tc.wantErr)
			}

			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("HasValue() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}
