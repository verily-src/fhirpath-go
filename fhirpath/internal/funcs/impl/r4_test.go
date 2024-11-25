package impl_test

import (
	"errors"
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr/exprtest"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs/impl"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/element/extension"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirtest"
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
