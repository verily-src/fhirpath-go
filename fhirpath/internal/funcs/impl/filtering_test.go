package impl_test

import (
	"errors"
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr/exprtest"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs/impl"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/slices"
	"google.golang.org/protobuf/testing/protocmp"
)

var contact = []*dtpb.ContactDetail{
	{
		Name: fhir.String("Vick"),
		Id:   fhir.String("123"),
	},
	{
		Name: fhir.String("Vick"),
		Id:   fhir.String("234"),
	},
	{
		Name: fhir.String("Matt"),
		Id:   fhir.String("123"),
	},
}

func TestWhere_Evaluates(t *testing.T) {
	nameEquality := &expr.EqualityExpression{
		Left:  &expr.FieldExpression{FieldName: "name"},
		Right: &expr.LiteralExpression{Literal: system.String("Vick")},
	}

	testCases := []struct {
		name            string
		inputCollection system.Collection
		inputArgs       []expr.Expression
		wantCollection  system.Collection
	}{
		{
			name:            "filters those that pass name query",
			inputCollection: slices.MustConvert[any](contact),
			inputArgs:       []expr.Expression{nameEquality},
			wantCollection:  slices.MustConvert[any](contact[0:2]),
		},
		{
			name:            "passes through when expression evaluates to singleton",
			inputCollection: slices.MustConvert[any](contact),
			inputArgs:       []expr.Expression{exprtest.Return(system.String("1"))},
			wantCollection:  slices.MustConvert[any](contact),
		},
		{
			name:            "passes through when expression evaluates to proto boolean true",
			inputCollection: slices.MustConvert[any](contact),
			inputArgs:       []expr.Expression{exprtest.Return(fhir.Boolean(true))},
			wantCollection:  slices.MustConvert[any](contact),
		},
		{
			name:            "filters out when expression evaluates to empty",
			inputCollection: slices.MustConvert[any](contact),
			inputArgs:       []expr.Expression{exprtest.Return()},
			wantCollection:  system.Collection{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Where(&expr.Context{}, tc.inputCollection, tc.inputArgs...)
			if err != nil {
				t.Fatalf("Where function returned unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.wantCollection, got, protocmp.Transform()); diff != "" {
				t.Errorf("Where function returned unexpected diff (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestWhere_RaisesError(t *testing.T) {
	testCases := []struct {
		name            string
		inputArgs       []expr.Expression
		inputCollection system.Collection
	}{
		{
			name:            "multiple arguments",
			inputArgs:       []expr.Expression{exprtest.Return(1), exprtest.Return(1)},
			inputCollection: slices.MustConvert[any](contact),
		},
		{
			name:            "argument expression raises error",
			inputArgs:       []expr.Expression{exprtest.Error(errors.New("some error"))},
			inputCollection: slices.MustConvert[any](contact),
		},
		{
			name:            "argument expression returns multiple items",
			inputArgs:       []expr.Expression{exprtest.Return(1, 2)},
			inputCollection: slices.MustConvert[any](contact),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := impl.Where(&expr.Context{}, tc.inputCollection, tc.inputArgs...); err == nil {
				t.Fatalf("evaluating Where function didn't return error when expected")
			}
		})
	}
}

func TestOfType_Evaluates(t *testing.T) {
	testCases := []struct {
		name            string
		inputCollection system.Collection
		inputArgs       expr.Expression
		wantCollection  system.Collection
	}{
		{
			name:            "filter on base node",
			inputCollection: slices.MustConvert[any](contact),
			inputArgs:       &expr.TypeExpression{Type: "FHIR.ContactDetail"},
			wantCollection:  slices.MustConvert[any](contact),
		},
		{
			name:            "filter on base node without namespace",
			inputCollection: slices.MustConvert[any](contact[0:3]),
			inputArgs:       &expr.TypeExpression{Type: "ContactDetail"},
			wantCollection:  slices.MustConvert[any](contact),
		},
		{
			name:            "returns empty collection",
			inputCollection: slices.MustConvert[any](contact[0:2]),
			inputArgs:       &expr.TypeExpression{Type: "string"},
			wantCollection:  system.Collection{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.OfType(&expr.Context{}, tc.inputCollection, tc.inputArgs)
			if err != nil {
				t.Fatalf("OfType function returned unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.wantCollection, got, protocmp.Transform()); diff != "" {
				t.Errorf("OfType function returned unexpected diff (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestOfType_RaisesError(t *testing.T) {
	testCases := []struct {
		name            string
		inputArgs       []expr.Expression
		inputCollection system.Collection
	}{
		{
			name:            "multiple arguments",
			inputArgs:       []expr.Expression{exprtest.Return(1), exprtest.Return(1)},
			inputCollection: slices.MustConvert[any](contact),
		},
		{
			name:            "argument expression raises error",
			inputArgs:       []expr.Expression{exprtest.Error(errors.New("some error"))},
			inputCollection: slices.MustConvert[any](contact),
		},
		{
			name:            "invalid argument expression",
			inputArgs:       []expr.Expression{exprtest.Return(1, 2)},
			inputCollection: slices.MustConvert[any](contact),
		},
		{
			name:            "invalid type",
			inputArgs:       []expr.Expression{&expr.TypeExpression{Type: "foo"}},
			inputCollection: slices.MustConvert[any](contact),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := impl.OfType(&expr.Context{}, tc.inputCollection, tc.inputArgs...); err == nil {
				t.Fatalf("evaluating OfType function didn't return error when expected")
			}
		})
	}
}
