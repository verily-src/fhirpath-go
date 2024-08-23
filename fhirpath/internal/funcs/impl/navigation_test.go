package impl_test

import (
	"testing"

	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/internal/slices"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs/impl"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"google.golang.org/protobuf/testing/protocmp"
)

var names = []*dtpb.HumanName{
	{
		Use:    &dtpb.HumanName_UseCode{Value: codes_go_proto.NameUseCode_OFFICIAL},
		Family: &dtpb.String{Value: "Smith"},
		Given: []*dtpb.String{
			{Value: "Bob"},
			{Value: "Bobbie"},
		},
	},
	{
		Use:  &dtpb.HumanName_UseCode{Value: codes_go_proto.NameUseCode_NICKNAME},
		Text: &dtpb.String{Value: "foo bar"},
		Period: &dtpb.Period{
			Start: &dtpb.DateTime{ValueUs: 10},
			End:   &dtpb.DateTime{ValueUs: 20},
		},
	},
}

func TestChildren_Evaluates(t *testing.T) {
	testCases := []struct {
		name            string
		inputCollection system.Collection
		wantCollection  system.Collection
	}{
		{
			name:            "empty input collection",
			inputCollection: system.Collection{},
			wantCollection:  system.Collection{},
		},
		{
			name:            "children of primitives yields empty result",
			inputCollection: slices.MustConvert[any](system.Collection{system.Integer(1), system.Boolean(true)}),
			wantCollection:  system.Collection{},
		},
		{
			name:            "returns all child nodes",
			inputCollection: slices.MustConvert[any](names),
			wantCollection: slices.MustConvert[any](
				system.Collection{
					names[0].GetUse(),
					names[0].GetFamily(),
					names[0].GetGiven()[0],
					names[0].GetGiven()[1],
					names[1].GetUse(),
					names[1].GetText(),
					names[1].GetPeriod(),
				},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Children(&expr.Context{}, tc.inputCollection)
			if err != nil {
				t.Fatalf("Children function returned unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.wantCollection, got, protocmp.Transform()); diff != "" {
				t.Errorf("Children function returned unexpected diff (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestDescendants_Evaluates(t *testing.T) {
	testCases := []struct {
		name            string
		inputCollection system.Collection
		wantCollection  system.Collection
	}{
		{
			name:            "empty input collection",
			inputCollection: system.Collection{},
			wantCollection:  system.Collection{},
		},
		{
			name:            "descendants of primitives yields empty result",
			inputCollection: slices.MustConvert[any](system.Collection{system.Integer(1), system.Boolean(true)}),
			wantCollection:  system.Collection{},
		},
		{
			name:            "returns all descendant nodes",
			inputCollection: slices.MustConvert[any](names),
			wantCollection: slices.MustConvert[any](
				system.Collection{
					names[0].GetUse(),
					names[0].GetFamily(),
					names[0].GetGiven()[0],
					names[0].GetGiven()[1],
					names[1].GetUse(),
					names[1].GetText(),
					names[1].GetPeriod(),
					names[1].GetPeriod().GetStart(),
					names[1].GetPeriod().GetEnd(),
				},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Descendants(&expr.Context{}, tc.inputCollection)
			if err != nil {
				t.Fatalf("Descendants function returned unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.wantCollection, got, protocmp.Transform()); diff != "" {
				t.Errorf("Descendants function returned unexpected diff (-want, +got):\n%s", diff)
			}
		})
	}
}
