package impl_test

import (
	"context"
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	pgp "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/parameters_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs/impl"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/fhirpath/terminology"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"google.golang.org/protobuf/testing/protocmp"
)

type fakeValueSet struct {
	valueSetId string
	system     string
	codes      []string
}

type fakeTerminologyService struct {
	dbItems []*fakeValueSet
}

func buildParameters(name string, result bool) *pgp.Parameters {
	return &pgp.Parameters{
		Parameter: []*pgp.Parameters_Parameter{
			{
				Name: &dtpb.String{Value: name},
				Value: &pgp.Parameters_Parameter_ValueX{
					Choice: &pgp.Parameters_Parameter_ValueX_Boolean{
						Boolean: &dtpb.Boolean{
							Value: result,
						},
					},
				},
			},
		},
	}
}

func (fts *fakeTerminologyService) ValueSetValidateCode(ctx context.Context, opts *terminology.ValueSetValidateCodeOptions) (*pgp.Parameters, error) {

	targetValueSet := opts.ID
	targetCode := opts.Code
	targetSystem := opts.System

	if targetValueSet == "" || targetSystem == "" {
		return buildParameters("result", false), nil
	}

	for _, item := range fts.dbItems {
		if item.valueSetId == targetValueSet {
			for _, code := range item.codes {
				if targetCode == code && targetSystem == item.system {
					return buildParameters("result", true), nil
				}
			}
		}
	}

	return buildParameters("result", false), nil
}

func TestMemberOf(t *testing.T) {

	var fakeTerminology = []*fakeValueSet{
		{
			valueSetId: "testValueSet",
			system:     "http://terminology.hl7.org/CodeSystem/v3-MaritalStatus",
			codes:      []string{"M", "D", "S", "W", "A", "L", "C", "P", "T", "U", "I"},
			// https://terminology.hl7.org/6.1.0/CodeSystem-v3-MaritalStatus.html
		},
	}

	fakeTerminologyService := &fakeTerminologyService{dbItems: fakeTerminology}

	var includedCoding = &dtpb.Coding{
		Code: &dtpb.Code{
			Value: "M",
		},
		System: fhir.URI("http://terminology.hl7.org/CodeSystem/v3-MaritalStatus"),
	}

	var notIncludedCoding = &dtpb.Coding{
		Code: &dtpb.Code{
			Value: "not included",
		},
		System: fhir.URI("http://terminology.hl7.org/CodeSystem/v3-MaritalStatus"),
	}

	var noneIncludedCodeableConcept = &dtpb.CodeableConcept{
		Coding: []*dtpb.Coding{
			{
				Code: &dtpb.Code{
					Value: "not included",
				},
				System: fhir.URI("http://terminology.hl7.org/CodeSystem/v3-MaritalStatus"),
			},
			{
				Code: &dtpb.Code{
					Value: "not included either",
				},
				System: fhir.URI("http://terminology.hl7.org/CodeSystem/v3-MaritalStatus"),
			},
		},
	}

	var includedCodeableConcept = &dtpb.CodeableConcept{
		Coding: []*dtpb.Coding{
			{
				Code: &dtpb.Code{
					Value: "not included",
				},
				System: fhir.URI("http://terminology.hl7.org/CodeSystem/v3-MaritalStatus"),
			},
			{
				Code: &dtpb.Code{
					Value: "M",
				},
				System: fhir.URI("http://terminology.hl7.org/CodeSystem/v3-MaritalStatus"),
			},
		},
	}

	valueSetExpr := &expr.LiteralExpression{
		Literal: system.String("testValueSet"),
	}

	testCases := []struct {
		name            string
		inputCollection system.Collection
		termServiceImpl terminology.Service
		args            []expr.Expression
		wantCollection  system.Collection
		wantErr         error
	}{
		{
			name:            "Coding is in the value set",
			inputCollection: system.Collection{includedCoding},
			termServiceImpl: fakeTerminologyService,
			args:            []expr.Expression{valueSetExpr},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
		{
			name:            "Coding is not in the value set",
			inputCollection: system.Collection{notIncludedCoding},
			termServiceImpl: fakeTerminologyService,
			args:            []expr.Expression{valueSetExpr},
			wantCollection:  system.Collection{system.Boolean(false)},
		},
		{
			name:            "CodeableConcept is not in the value set",
			inputCollection: system.Collection{noneIncludedCodeableConcept},
			termServiceImpl: fakeTerminologyService,
			args:            []expr.Expression{valueSetExpr},
			wantCollection:  system.Collection{system.Boolean(false)},
		},
		{
			name:            "CodeableConcept is in the value set",
			inputCollection: system.Collection{includedCodeableConcept},
			termServiceImpl: fakeTerminologyService,
			args:            []expr.Expression{valueSetExpr},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotCollection, gotErr := impl.MemberOf(&expr.Context{TermService: tc.termServiceImpl}, tc.inputCollection, tc.args...)

			if !cmp.Equal(gotErr, tc.wantErr, cmpopts.EquateErrors()) {
				t.Errorf("MemberOf() gotErr = %v, wantErr = %v", gotErr, tc.wantErr)
			}
			if diff := cmp.Diff(tc.wantCollection, gotCollection, protocmp.Transform()); diff != "" {
				t.Errorf("MemberOf() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
