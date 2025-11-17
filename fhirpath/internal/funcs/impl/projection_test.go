package impl_test

import (
	"errors"
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	qrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/questionnaire_response_go_proto"

	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr/exprtest"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs/impl"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/slices"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

var address = []*dtpb.Address{
	{
		Line: []*dtpb.String{
			fhir.String("123 Main St"),
			fhir.String("Apt 1"),
		},
		State: fhir.String("CA"),
	},
	{
		Line: []*dtpb.String{
			fhir.String("456 Main St"),
			fhir.String("Apt 2"),
		},
		State: fhir.String("TX"),
	},
}

func TestSelect_Evaluates(t *testing.T) {
	testCases := []struct {
		name            string
		inputCollection system.Collection
		inputArgs       []expr.Expression
		wantCollection  system.Collection
	}{
		{
			name:            "projection on empty collection",
			inputCollection: system.Collection{},
			inputArgs:       []expr.Expression{exprtest.Return(system.Boolean(true))},
			wantCollection:  system.Collection{},
		},
		{
			name:            "projection yields empty result",
			inputCollection: slices.MustConvert[any](address),
			inputArgs:       []expr.Expression{exprtest.Return()},
			wantCollection:  system.Collection{},
		},
		{
			name:            "project state field",
			inputCollection: slices.MustConvert[any](address),
			inputArgs:       []expr.Expression{&expr.FieldExpression{FieldName: "state"}},
			wantCollection:  system.Collection{address[0].GetState(), address[1].GetState()},
		},
		{
			name:            "projection flattens output collections",
			inputCollection: slices.MustConvert[any](address),
			inputArgs:       []expr.Expression{&expr.FieldExpression{FieldName: "line"}},
			wantCollection: system.Collection{
				address[0].GetLine()[0],
				address[0].GetLine()[1],
				address[1].GetLine()[0],
				address[1].GetLine()[1],
			},
		},
		{
			name:            "does not raise error if field is valid for at least one input",
			inputCollection: system.Collection{address[0], fhir.String("string")},
			inputArgs:       []expr.Expression{&expr.FieldExpression{FieldName: "state"}},
			wantCollection:  system.Collection{address[0].GetState()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Select(&expr.Context{}, tc.inputCollection, tc.inputArgs...)
			if err != nil {
				t.Fatalf("Select function returned unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.wantCollection, got, protocmp.Transform()); diff != "" {
				t.Errorf("Select function returned unexpected diff (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestSelect_RaisesError(t *testing.T) {
	testCases := []struct {
		name            string
		inputArgs       []expr.Expression
		inputCollection system.Collection
	}{
		{
			name:            "multiple arguments",
			inputArgs:       []expr.Expression{exprtest.Return(1), exprtest.Return(1)},
			inputCollection: slices.MustConvert[any](address),
		},
		{
			name:            "argument expression raises error",
			inputArgs:       []expr.Expression{exprtest.Error(errors.New("some error"))},
			inputCollection: slices.MustConvert[any](address),
		},
		{
			name:            "invalid field as argument expression",
			inputArgs:       []expr.Expression{&expr.FieldExpression{FieldName: "invalid"}},
			inputCollection: slices.MustConvert[any](address),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := impl.Select(&expr.Context{}, tc.inputCollection, tc.inputArgs...); err == nil {
				t.Fatalf("evaluating Select function didn't return error when expected")
			}
		})
	}
}

func TestRepeat_Evaluates(t *testing.T) {
	questionnaire := &qrpb.QuestionnaireResponse{
		Item: []*qrpb.QuestionnaireResponse_Item{
			{
				Answer: []*qrpb.QuestionnaireResponse_Item_Answer{{
					Value: &qrpb.QuestionnaireResponse_Item_Answer_ValueX{
						Choice: &qrpb.QuestionnaireResponse_Item_Answer_ValueX_StringValue{StringValue: fhir.String("foo")},
					},
					Item: []*qrpb.QuestionnaireResponse_Item{
						{
							Answer: []*qrpb.QuestionnaireResponse_Item_Answer{{
								Value: &qrpb.QuestionnaireResponse_Item_Answer_ValueX{
									Choice: &qrpb.QuestionnaireResponse_Item_Answer_ValueX_StringValue{StringValue: fhir.String("baz")},
								},
							}},
						},
					},
				}},
				Item: []*qrpb.QuestionnaireResponse_Item{{
					Answer: []*qrpb.QuestionnaireResponse_Item_Answer{
						{
							Value: &qrpb.QuestionnaireResponse_Item_Answer_ValueX{
								Choice: &qrpb.QuestionnaireResponse_Item_Answer_ValueX_StringValue{StringValue: fhir.String("bar")},
							},
						},
					},
				}},
			},
		},
	}

	testCases := []struct {
		name            string
		inputCollection system.Collection
		inputArgs       []expr.Expression
		wantCollection  system.Collection
	}{
		{
			name:            "repeat on empty collection",
			inputCollection: system.Collection{},
			inputArgs:       []expr.Expression{exprtest.Return(system.Boolean(true))},
			wantCollection:  nil,
		},
		{
			name:            "repeat literal",
			inputCollection: system.Collection{questionnaire},
			inputArgs:       []expr.Expression{&expr.LiteralExpression{Literal: system.String("foo")}},
			wantCollection:  system.Collection{system.String("foo")},
		},
		{
			name:            "repeat yields empty result",
			inputCollection: slices.MustConvert[any](address),
			inputArgs:       []expr.Expression{exprtest.Return()},
			wantCollection:  nil,
		},
		{
			name:            "repeat field projection",
			inputCollection: system.Collection{questionnaire},
			inputArgs:       []expr.Expression{&expr.FieldExpression{FieldName: "item"}},
			wantCollection:  system.Collection{questionnaire.GetItem()[0], questionnaire.GetItem()[0].GetItem()[0]},
		},
		{
			name:            "repeat path projection",
			inputCollection: system.Collection{questionnaire},
			inputArgs:       []expr.Expression{&expr.ExpressionSequence{Expressions: []expr.Expression{&expr.FieldExpression{FieldName: "item"}, &expr.FieldExpression{FieldName: "answer"}}}},
			wantCollection:  system.Collection{questionnaire.GetItem()[0].GetAnswer()[0], questionnaire.GetItem()[0].GetAnswer()[0].GetItem()[0].GetAnswer()[0]},
		},
		{
			name:            "descendants using repeat + children",
			inputCollection: system.Collection{questionnaire.GetItem()[0].GetAnswer()[0]},
			inputArgs:       []expr.Expression{&expr.FunctionExpression{Fn: impl.Children}},
			wantCollection: system.Collection{
				questionnaire.Item[0].Answer[0].Value.GetStringValue(),
				questionnaire.Item[0].Answer[0].Item[0],
				questionnaire.Item[0].Answer[0].Item[0].Answer[0],
				questionnaire.Item[0].Answer[0].Item[0].Answer[0].GetValue().GetStringValue(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Repeat(&expr.Context{}, tc.inputCollection, tc.inputArgs...)
			if err != nil {
				t.Fatalf("Repeat function returned unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.wantCollection, got, protocmp.Transform()); diff != "" {
				t.Errorf("Repeat function returned unexpected diff (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestRepeat_RaisesError(t *testing.T) {
	testCases := []struct {
		name            string
		inputArgs       []expr.Expression
		inputCollection system.Collection
	}{
		{
			name:            "multiple arguments",
			inputArgs:       []expr.Expression{exprtest.Return(1), exprtest.Return(1)},
			inputCollection: slices.MustConvert[any](address),
		},
		{
			name:            "argument expression raises error",
			inputArgs:       []expr.Expression{exprtest.Error(errors.New("some error"))},
			inputCollection: slices.MustConvert[any](address),
		},
		{
			name:            "invalid field as argument expression",
			inputArgs:       []expr.Expression{&expr.FieldExpression{FieldName: "invalid"}},
			inputCollection: slices.MustConvert[any](address),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := impl.Repeat(&expr.Context{}, tc.inputCollection, tc.inputArgs...); err == nil {
				t.Fatalf("evaluating Repeat function didn't return error when expected")
			}
		})
	}
}
