package impl_test

import (
	"errors"
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs/impl"
	"github.com/verily-src/fhirpath-go/fhirpath/resolver"
	"github.com/verily-src/fhirpath-go/fhirpath/resolver/resolvertest"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestResolve(t *testing.T) {
	var patientChu = &ppb.Patient{
		Id: fhir.ID("123"),
	}

	var patientChuRef = &dtpb.Reference{
		Type: fhir.URI("Patient"),
		Reference: &dtpb.Reference_PatientId{
			PatientId: &dtpb.ReferenceId{Value: "123"}},
	}

	var patientChuUri = &dtpb.Uri{
		Value: "Patient/123",
	}

	var patientChuUrl = &dtpb.Url{
		Value: "http://example.com/Patient/123",
	}

	var patientChuCanonical = &dtpb.Canonical{
		Value: "Patient/123",
	}

	var patientChuString = &dtpb.String{
		Value: "Patient/123",
	}

	resolveErr := errors.New("some resolve() error")

	testCases := []struct {
		name            string
		inputCollection system.Collection
		resolverImpl    resolver.Resolver
		args            []expr.Expression
		wantCollection  system.Collection
		wantErr         error
	}{
		{
			name:            "happy path; successful reference resolution",
			inputCollection: system.Collection{patientChuRef},
			resolverImpl:    resolvertest.HappyResolver(patientChu),
			wantCollection:  system.Collection{patientChu},
		},
		{
			name:            "happy path; successful uri resolution",
			inputCollection: system.Collection{patientChuUri},
			resolverImpl:    resolvertest.HappyResolver(patientChu),
			wantCollection:  system.Collection{patientChu},
		},
		{
			name:            "happy path; successful string resolution",
			inputCollection: system.Collection{patientChuString},
			resolverImpl:    resolvertest.HappyResolver(patientChu),
			wantCollection:  system.Collection{patientChu},
		},
		{
			name:            "happy path; successful url resolution",
			inputCollection: system.Collection{patientChuUrl},
			resolverImpl:    resolvertest.HappyResolver(patientChu),
			wantCollection:  system.Collection{patientChu},
		},
		{
			name:            "happy path; successful canonical resolution",
			inputCollection: system.Collection{patientChuCanonical},
			resolverImpl:    resolvertest.HappyResolver(patientChu),
			wantCollection:  system.Collection{patientChu},
		},
		{
			name:            "happy path; empty input",
			inputCollection: system.Collection{},
			resolverImpl:    resolvertest.HappyResolver(patientChu),
			wantCollection:  system.Collection{},
		},
		{
			name:            "too many arguments; throws error",
			inputCollection: system.Collection{patientChuRef},
			resolverImpl:    resolvertest.HappyResolver(patientChu),
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("")},
			},
			wantErr: impl.ErrWrongArity,
		},
		{
			name:            "happy path; input doesn't have a string item - returns empty collection",
			inputCollection: system.Collection{system.Integer(900)},
			resolverImpl:    resolvertest.HappyResolver(patientChu),
			wantCollection:  system.Collection{},
		},
		{
			name:            "resolverImpl resolve() returns an error",
			inputCollection: system.Collection{patientChuRef},
			resolverImpl:    resolvertest.ErroringResolver(resolveErr),
			wantErr:         resolveErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotCollection, gotErr := impl.Resolve(&expr.Context{Resolver: tc.resolverImpl}, tc.inputCollection, tc.args...)

			if !cmp.Equal(gotErr, tc.wantErr, cmpopts.EquateErrors()) {
				t.Errorf("Resolve() gotErr = %v, wantErr = %v", gotErr, tc.wantErr)
			}
			if diff := cmp.Diff(tc.wantCollection, gotCollection, protocmp.Transform()); diff != "" {
				t.Errorf("Resolve() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
