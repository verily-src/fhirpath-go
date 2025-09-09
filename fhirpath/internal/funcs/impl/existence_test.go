package impl_test

import (
	"errors"
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr/exprtest"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs/impl"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/reflection"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/slices"
	"google.golang.org/protobuf/testing/protocmp"
)

var coding = []*dtpb.Coding{
	fhir.Coding("loinc-system", "loinc-code"),
	fhir.Coding("loinc-system", "generic-code"),
	fhir.Coding("snomed-system", "snomed-code"),
	fhir.Coding("snomed-system", "snomed-code"),
	fhir.Coding("icd10-system", "icd10-code"),
	{},
}

func TestAllTrue(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:    "returns true if input is empty",
			input:   system.Collection{},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name: "returns false if input contains a false value",
			input: system.Collection{
				system.Boolean(true),
				system.Boolean(true),
				system.Boolean(false)},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name: "returns true if input contains only true values",
			input: system.Collection{system.Boolean(true),
				system.Boolean(true),
				system.Boolean(true)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name: "returns true if input (fhir.Boolean) contains only true values",
			input: system.Collection{fhir.Boolean(true),
				fhir.Boolean(true),
				fhir.Boolean(true)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.AllTrue(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("AllTrue() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("AllTrue() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestAnyTrue(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:    "returns false if input is empty",
			input:   system.Collection{},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name: "returns true if input contains a true value",
			input: system.Collection{
				system.Boolean(false),
				system.Boolean(false),
				system.Boolean(true)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name: "returns false if input contains only false values",
			input: system.Collection{system.Boolean(false),
				system.Boolean(false),
				system.Boolean(false)},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name: "returns false if input (fhir.Boolean) contains only false values",
			input: system.Collection{fhir.Boolean(false),
				fhir.Boolean(false),
				fhir.Boolean(false)},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.AnyTrue(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("AnyTrue() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("AnyTrue() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestAllFalse(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:    "returns true if input is empty",
			input:   system.Collection{},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name: "returns false if input contains a true value",
			input: system.Collection{
				system.Boolean(false),
				system.Boolean(true),
				system.Boolean(false)},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name: "returns true if input contains only false values",
			input: system.Collection{system.Boolean(false),
				system.Boolean(false),
				system.Boolean(false)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name: "returns true if input (fhir.Boolean) contains only false values",
			input: system.Collection{fhir.Boolean(false),
				fhir.Boolean(false),
				fhir.Boolean(false)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.AllFalse(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("AllFalse() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("AllFalse() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestAnyFalse(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:    "returns false if input is empty",
			input:   system.Collection{},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name: "returns true if input contains a false value",
			input: system.Collection{
				system.Boolean(true),
				system.Boolean(true),
				system.Boolean(false)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name: "returns true if input contains only false values",
			input: system.Collection{system.Boolean(false),
				system.Boolean(false),
				system.Boolean(false)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name: "returns true if input (fhir.Boolean) contains only false values",
			input: system.Collection{fhir.Boolean(false),
				fhir.Boolean(false),
				fhir.Boolean(false)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name: "returns false if input is not boolean",
			input: system.Collection{fhir.Integer(5),
				fhir.Integer(6)},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.AnyFalse(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("AnyFalse() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("AnyFalse() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestAll(t *testing.T) {
	testCases := []struct {
		name            string
		inputCollection system.Collection
		inputArgs       []expr.Expression
		wantCollection  system.Collection
		wantErr         bool
		wantErrMsg      string
	}{
		{
			name:            "returns true if input is empty",
			inputCollection: system.Collection{},
			inputArgs:       []expr.Expression{exprtest.Return(system.Boolean(true))},
			wantCollection:  system.Collection{system.Boolean(true)},
			wantErr:         false,
		},
		{
			name: "returns false if expression returns false",
			inputCollection: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3)},
			inputArgs:      []expr.Expression{exprtest.Return(system.Boolean(false))},
			wantCollection: system.Collection{system.Boolean(false)},
			wantErr:        false,
		},
		{
			name: "returns true if expression returns true",
			inputCollection: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3)},
			inputArgs:      []expr.Expression{exprtest.Return(system.Boolean(true))},
			wantCollection: system.Collection{system.Boolean(true)},
			wantErr:        false,
		},
		{
			name: "returns true if all elements are integers",
			inputCollection: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3)},
			inputArgs: []expr.Expression{&expr.IsExpression{
				Expr: &expr.IdentityExpression{},
				Type: reflection.MustCreateTypeSpecifier("System", "Integer"),
			}},
			wantCollection: system.Collection{system.Boolean(true)},
			wantErr:        false,
		},
		{
			name: "returns false if not all elements are integers",
			inputCollection: system.Collection{
				system.Integer(1),
				system.String("test"),
				system.Integer(3)},
			inputArgs: []expr.Expression{&expr.IsExpression{
				Expr: &expr.IdentityExpression{},
				Type: reflection.MustCreateTypeSpecifier("System", "Integer"),
			}},
			wantCollection: system.Collection{system.Boolean(false)},
			wantErr:        false,
		},
		{
			name: "returns error if criteria expression raises error",
			inputCollection: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3)},
			inputArgs:      []expr.Expression{exprtest.Error(errors.New("some error"))},
			wantCollection: nil,
			wantErr:        true,
			wantErrMsg:     "evaluating criteria expression resulted in an error: some error",
		},
		{
			name: "returns error if args length is not 1",
			inputCollection: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3)},
			inputArgs:      []expr.Expression{},
			wantCollection: nil,
			wantErr:        true,
			wantErrMsg:     "incorrect function arity: received 0 arguments, expected 1",
		},
		{
			name: "returns true if criteria expression returns non-boolean values",
			inputCollection: system.Collection{
				system.Integer(1),
				system.Integer(2)},
			inputArgs:      []expr.Expression{exprtest.Return(system.String("non-boolean"))},
			wantCollection: system.Collection{system.Boolean(true)},
			wantErr:        false,
		},
		{
			name: "returns false if criteria expression returns an empty collection",
			inputCollection: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3)},
			inputArgs:      []expr.Expression{exprtest.Return()},
			wantCollection: system.Collection{system.Boolean(false)},
			wantErr:        false,
		},
		{
			name: "returns false if criteria expression returns multiple values",
			inputCollection: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3)},
			inputArgs:      []expr.Expression{exprtest.Return(system.Boolean(true), system.Boolean(false))},
			wantCollection: nil,
			wantErr:        true,
			wantErrMsg:     "collection is not singleton",
		},
		{
			name: "returns true if input collection contains mixed types",
			inputCollection: system.Collection{
				system.Integer(1),
				system.String("test"),
				system.Boolean(true)},
			inputArgs: []expr.Expression{&expr.IsExpression{
				Expr: &expr.IdentityExpression{},
				Type: reflection.MustCreateTypeSpecifier("System", "Any"),
			}},
			wantCollection: system.Collection{system.Boolean(true)},
			wantErr:        false,
		},
		{
			name: "returns true if input collection contains only strings",
			inputCollection: system.Collection{
				system.String("1"),
				system.String("test"),
				system.String("true")},
			inputArgs: []expr.Expression{&expr.IsExpression{
				Expr: &expr.IdentityExpression{},
				Type: reflection.MustCreateTypeSpecifier("System", "String"),
			}},
			wantCollection: system.Collection{system.Boolean(true)},
			wantErr:        false,
		},
		{
			name: "returns false if input collection not contains only strings",
			inputCollection: system.Collection{
				system.Boolean(true),
				system.String("test"),
				system.String("true")},
			inputArgs: []expr.Expression{&expr.IsExpression{
				Expr: &expr.IdentityExpression{},
				Type: reflection.MustCreateTypeSpecifier("System", "String"),
			}},
			wantCollection: system.Collection{system.Boolean(false)},
			wantErr:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.All(&expr.Context{}, tc.inputCollection, tc.inputArgs...)
			if (err != nil) != tc.wantErr {
				t.Errorf("All() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if err != nil && err.Error() != tc.wantErrMsg {
				t.Errorf("All() error message = %v, wantErrMsg %v", err.Error(), tc.wantErrMsg)
			}
			if diff := cmp.Diff(tc.wantCollection, got, protocmp.Transform()); diff != "" {
				t.Errorf("All() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestCount(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:    "returns 0 if input is empty",
			input:   system.Collection{},
			want:    system.Collection{system.Integer(0)},
			wantErr: false,
		},
		{
			name:    "input 1 if input length is 1 ",
			input:   system.Collection{system.Integer(1)},
			want:    system.Collection{system.Integer(1)},
			wantErr: false,
		},
		{
			name: "input 5 if input length is 5 ",
			input: system.Collection{
				system.Integer(2),
				system.Integer(4),
				system.Integer(6),
				system.Integer(8),
				system.Integer(10)},
			want:    system.Collection{system.Integer(5)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Count(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("Count() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Count() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestExists_Evaluates(t *testing.T) {
	testCases := []struct {
		name            string
		inputCollection system.Collection
		inputArgs       []expr.Expression
		wantCollection  system.Collection
	}{
		{
			name:            "exists loinc system + loinc code",
			inputCollection: slices.MustConvert[any](coding),
			inputArgs: []expr.Expression{&expr.BooleanExpression{
				Left: &expr.EqualityExpression{
					Left:  &expr.FieldExpression{FieldName: "system"},
					Right: &expr.LiteralExpression{Literal: system.String("loinc-system")},
				},
				Right: &expr.EqualityExpression{
					Left:  &expr.FieldExpression{FieldName: "code"},
					Right: &expr.LiteralExpression{Literal: system.String("loinc-code")},
				},
				Op: expr.And,
			}},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name:            "exists loinc system",
			inputCollection: slices.MustConvert[any](coding),
			inputArgs: []expr.Expression{&expr.EqualityExpression{
				Left:  &expr.FieldExpression{FieldName: "system"},
				Right: &expr.LiteralExpression{Literal: system.String("loinc-system")},
			}},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name:            "exists icd10 system + icd10 code",
			inputCollection: slices.MustConvert[any](coding),
			inputArgs: []expr.Expression{&expr.BooleanExpression{
				Left: &expr.EqualityExpression{
					Left:  &expr.FieldExpression{FieldName: "system"},
					Right: &expr.LiteralExpression{Literal: system.String("icd10-system")},
				},
				Right: &expr.EqualityExpression{
					Left:  &expr.FieldExpression{FieldName: "code"},
					Right: &expr.LiteralExpression{Literal: system.String("icd10-code")},
				},
				Op: expr.And,
			}},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name:            "exists snomed system + snomed code",
			inputCollection: slices.MustConvert[any](coding),
			inputArgs: []expr.Expression{&expr.BooleanExpression{
				Left: &expr.EqualityExpression{
					Left:  &expr.FieldExpression{FieldName: "system"},
					Right: &expr.LiteralExpression{Literal: system.String("snomed-system")},
				},
				Right: &expr.EqualityExpression{
					Left:  &expr.FieldExpression{FieldName: "code"},
					Right: &expr.LiteralExpression{Literal: system.String("snomed-code")},
				},
				Op: expr.And,
			}},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name:            "does not exist snomed system + loinc code",
			inputCollection: slices.MustConvert[any](coding),
			inputArgs: []expr.Expression{&expr.BooleanExpression{
				Left: &expr.EqualityExpression{
					Left:  &expr.FieldExpression{FieldName: "system"},
					Right: &expr.LiteralExpression{Literal: system.String("snomed-system")},
				},
				Right: &expr.EqualityExpression{
					Left:  &expr.FieldExpression{FieldName: "code"},
					Right: &expr.LiteralExpression{Literal: system.String("loic-code")},
				},
				Op: expr.And,
			}},
			wantCollection: system.Collection{system.Boolean(false)},
		},
		{
			name:            "non empty inputs with empty args",
			inputCollection: slices.MustConvert[any](coding),
			inputArgs:       nil,
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "empty inputs with empty args",
			inputCollection: system.Collection{},
			inputArgs:       nil,
			wantCollection:  system.Collection{system.Boolean(false)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Exists(&expr.Context{}, tc.inputCollection, tc.inputArgs...)
			if err != nil {
				t.Fatalf("Exists function returned unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.wantCollection, got, protocmp.Transform()); diff != "" {
				t.Errorf("Exists function returned unexpected diff (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestExists_RaisesError(t *testing.T) {
	testCases := []struct {
		name            string
		inputArgs       []expr.Expression
		inputCollection system.Collection
	}{
		{
			name:            "multiple arguments",
			inputArgs:       []expr.Expression{exprtest.Return(1), exprtest.Return(1)},
			inputCollection: slices.MustConvert[any](coding),
		},
		{
			name:            "argument expression raises error",
			inputArgs:       []expr.Expression{exprtest.Error(errors.New("some error"))},
			inputCollection: slices.MustConvert[any](coding),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := impl.Exists(&expr.Context{}, tc.inputCollection, tc.inputArgs...); err == nil {
				t.Fatalf("evaluating Exists function didn't return error when expected")
			}
		})
	}
}

func TestEmpty_Evaluates(t *testing.T) {
	testCases := []struct {
		name            string
		inputCollection system.Collection
		inputArgs       []expr.Expression
		wantCollection  system.Collection
	}{
		{
			name:            "empty inputs with empty args",
			inputCollection: system.Collection{},
			inputArgs:       nil,
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "non empty inputs",
			inputCollection: slices.MustConvert[any](coding),
			inputArgs:       nil,
			wantCollection:  system.Collection{system.Boolean(false)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Empty(&expr.Context{}, tc.inputCollection, tc.inputArgs...)
			if err != nil {
				t.Fatalf("Empty function returned unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.wantCollection, got, protocmp.Transform()); diff != "" {
				t.Errorf("Empty function returned unexpected diff (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestEmpty_RaisesError(t *testing.T) {
	testCases := []struct {
		name            string
		inputArgs       []expr.Expression
		inputCollection system.Collection
	}{
		{
			name:            "multiple arguments",
			inputArgs:       []expr.Expression{exprtest.Return(1)},
			inputCollection: slices.MustConvert[any](coding),
		},
		{
			name:            "argument expression raises error",
			inputArgs:       []expr.Expression{exprtest.Error(errors.New("some error"))},
			inputCollection: slices.MustConvert[any](coding),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := impl.Empty(&expr.Context{}, tc.inputCollection, tc.inputArgs...); err == nil {
				t.Fatalf("evaluating Empty function didn't return error when expected")
			}
		})
	}
}

func TestSubsetOf_Success(t *testing.T) {
	testCases := []struct {
		name  string
		input system.Collection
		args  []expr.Expression
		want  system.Collection
	}{
		{
			name:  "returns true if input is empty",
			input: system.Collection{},
			args:  []expr.Expression{exprtest.Return(system.Integer(1), system.Integer(2))},
			want:  system.Collection{system.Boolean(true)},
		},
		{
			name: "returns false if argument collection is empty but input is not",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
			},
			args: []expr.Expression{exprtest.Return()},
			want: system.Collection{system.Boolean(false)},
		},
		{
			name: "returns true if input is subset of argument collection",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
			},
			args: []expr.Expression{exprtest.Return(
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
			)},
			want: system.Collection{system.Boolean(true)},
		},
		{
			name: "returns false if input is not subset of argument collection",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(5),
			},
			args: []expr.Expression{exprtest.Return(
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
			)},
			want: system.Collection{system.Boolean(false)},
		},
		{
			name: "returns true if input equals argument collection",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
			},
			args: []expr.Expression{exprtest.Return(
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
			)},
			want: system.Collection{system.Boolean(true)},
		},
		{
			name: "returns true with fhir.Integer types - subset case",
			input: system.Collection{
				fhir.Integer(1),
				fhir.Integer(2),
			},
			args: []expr.Expression{exprtest.Return(
				fhir.Integer(1),
				fhir.Integer(2),
				fhir.Integer(3),
			)},
			want: system.Collection{system.Boolean(true)},
		},
		{
			name: "returns false with fhir.Integer types - not subset case",
			input: system.Collection{
				fhir.Integer(1),
				fhir.Integer(4),
			},
			args: []expr.Expression{exprtest.Return(
				fhir.Integer(1),
				fhir.Integer(2),
				fhir.Integer(3),
			)},
			want: system.Collection{system.Boolean(false)},
		},
		{
			name: "returns true with mixed system and fhir types",
			input: system.Collection{
				system.Integer(1),
				fhir.Integer(2),
			},
			args: []expr.Expression{exprtest.Return(
				fhir.Integer(1),
				system.Integer(2),
				system.Integer(3),
			)},
			want: system.Collection{system.Boolean(true)},
		},
		{
			name: "returns true with string collections - subset case",
			input: system.Collection{
				system.String("a"),
				system.String("b"),
			},
			args: []expr.Expression{exprtest.Return(
				system.String("a"),
				system.String("b"),
				system.String("c"),
			)},
			want: system.Collection{system.Boolean(true)},
		},
		{
			name: "returns false with string collections - not subset case",
			input: system.Collection{
				system.String("a"),
				system.String("d"),
			},
			args: []expr.Expression{exprtest.Return(
				system.String("a"),
				system.String("b"),
				system.String("c"),
			)},
			want: system.Collection{system.Boolean(false)},
		},
		{
			name: "returns true with coding collections - subset case",
			input: system.Collection{
				fhir.Coding("loinc-system", "loinc-code"),
				fhir.Coding("snomed-system", "snomed-code"),
			},
			args: []expr.Expression{exprtest.Return(
				fhir.Coding("loinc-system", "loinc-code"),
				fhir.Coding("snomed-system", "snomed-code"),
				fhir.Coding("icd10-system", "icd10-code"),
			)},
			want: system.Collection{system.Boolean(true)},
		},
		{
			name: "returns false with coding collections - not subset case",
			input: system.Collection{
				fhir.Coding("loinc-system", "loinc-code"),
				fhir.Coding("missing-system", "missing-code"),
			},
			args: []expr.Expression{exprtest.Return(
				fhir.Coding("loinc-system", "loinc-code"),
				fhir.Coding("snomed-system", "snomed-code"),
			)},
			want: system.Collection{system.Boolean(false)},
		},
		{
			name: "treats FHIR messages with different field orders as equivalent",
			input: system.Collection{
				&dtpb.HumanName{
					Given:  []*dtpb.String{fhir.String("John")},
					Family: fhir.String("Doe"),
				},
			},
			args: []expr.Expression{exprtest.Return(
				&dtpb.HumanName{
					Family: fhir.String("Doe"),
					Given:  []*dtpb.String{fhir.String("John")},
				},
				&dtpb.HumanName{
					Given:  []*dtpb.String{fhir.String("Jane")},
					Family: fhir.String("Smith"),
				},
			)},
			want: system.Collection{system.Boolean(true)},
		},
		{
			name: "returns false for large collections - element not in superset",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5),
			},
			args: []expr.Expression{exprtest.Return(
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(6),
				system.Integer(7),
			)},
			want: system.Collection{system.Boolean(false)},
		},
		{
			name: "returns true for large collections - proper subset",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5),
			},
			args: []expr.Expression{exprtest.Return(
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5),
				system.Integer(6),
			)},
			want: system.Collection{system.Boolean(true)},
		},
		{
			name: "returns false for multiset - duplicate elements in input",
			input: system.Collection{
				system.Integer(1),
				system.Integer(1),
			},
			args: []expr.Expression{exprtest.Return(
				system.Integer(1),
				system.Integer(2),
			)},
			want: system.Collection{system.Boolean(false)},
		},
		{
			name: "returns false for multiset - input has more duplicates than other collection",
			input: system.Collection{
				system.Integer(1),
				system.Integer(1),
				system.Integer(1),
			},
			args: []expr.Expression{exprtest.Return(
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
			)},
			want: system.Collection{system.Boolean(false)},
		},
		{
			name: "returns true for multiset - input has fewer duplicates than other collection",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
			},
			args: []expr.Expression{exprtest.Return(
				system.Integer(1),
				system.Integer(1),
				system.Integer(2),
				system.Integer(2),
			)},
			want: system.Collection{system.Boolean(true)},
		},
		{
			name:  "returns true for single element subset",
			input: system.Collection{system.Integer(1)},
			args:  []expr.Expression{exprtest.Return(system.Integer(1), system.Integer(2))},
			want:  system.Collection{system.Boolean(true)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.SubsetOf(&expr.Context{}, tc.input, tc.args...)
			if err != nil {
				t.Errorf("SubsetOf() error = %v", err)
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("SubsetOf() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestSubsetOf_Errors(t *testing.T) {
	testCases := []struct {
		name            string
		inputArgs       []expr.Expression
		inputCollection system.Collection
		wantErrMsg      string
	}{
		{
			name:            "no arguments provided",
			inputArgs:       []expr.Expression{},
			inputCollection: system.Collection{system.Integer(1), system.Integer(2)},
			wantErrMsg:      "incorrect function arity: received 0 arguments, expected 1",
		},
		{
			name: "multiple arguments provided",
			inputArgs: []expr.Expression{
				exprtest.Return(system.Integer(1)),
				exprtest.Return(system.Integer(2)),
			},
			inputCollection: system.Collection{system.Integer(1)},
			wantErrMsg:      "incorrect function arity: received 2 arguments, expected 1",
		},
		{
			name:            "argument expression raises error",
			inputArgs:       []expr.Expression{exprtest.Error(errors.New("some error"))},
			inputCollection: system.Collection{system.Integer(1)},
			wantErrMsg:      "some error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := impl.SubsetOf(&expr.Context{}, tc.inputCollection, tc.inputArgs...)
			if err == nil {
				t.Fatalf("evaluating SubsetOf function didn't return error when expected")
			}

			if err.Error() != tc.wantErrMsg {
				t.Errorf("SubsetOf() error = %v, wantErrMsg %v", err, tc.wantErrMsg)
			}
		})
	}
}
