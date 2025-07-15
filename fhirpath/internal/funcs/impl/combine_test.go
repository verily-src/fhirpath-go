package impl

import (
	"math"
	"testing"

	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr/exprtest"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestCombine(t *testing.T) {
	ctx := &expr.Context{}

	testCases := []struct {
		name     string
		input    system.Collection
		args     []expr.Expression
		expected system.Collection
	}{
		{
			name:     "Combine two non-empty collections",
			input:    system.Collection{1, 2, 3},
			args:     []expr.Expression{exprtest.Return(4, 5), exprtest.Return(6)},
			expected: system.Collection{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "Combine empty input with non-empty collection",
			input:    system.Collection{},
			args:     []expr.Expression{exprtest.Return(7, 8, 9)},
			expected: system.Collection{7, 8, 9},
		},
		{
			name:     "Combine non-empty input with empty collection",
			input:    system.Collection{10, 11},
			args:     []expr.Expression{exprtest.Return()},
			expected: system.Collection{10, 11},
		},
		{
			name:     "Combine multiple empty collections",
			input:    system.Collection{},
			args:     []expr.Expression{exprtest.Return(), exprtest.Return()},
			expected: nil,
		},
		{
			name:     "Combine with no other collections",
			input:    system.Collection{12, 13},
			args:     []expr.Expression{},
			expected: system.Collection{12, 13},
		},
		{
			name:     "Combine with duplicate values",
			input:    system.Collection{14, 15},
			args:     []expr.Expression{exprtest.Return(15, 16), exprtest.Return(14)},
			expected: system.Collection{14, 15, 15, 16, 14},
		},
		{
			name:     "Combine with nil collections",
			input:    system.Collection{17, 18},
			args:     []expr.Expression{exprtest.Return(), exprtest.Return(19)},
			expected: system.Collection{17, 18, 19},
		},
		{
			name:     "Combine with different types",
			input:    system.Collection{20, 21, fhir.UnsignedInt(math.MaxUint32)},
			args:     []expr.Expression{exprtest.Return(22.0, "23"), exprtest.Return(24)},
			expected: system.Collection{20, 21, fhir.UnsignedInt(math.MaxUint32), 22.0, "23", 24},
		},
		{
			name:     "Combine FHIR elements",
			input:    system.Collection{fhir.ID("123"), &ppb.Patient_GenderCode{Value: cpb.AdministrativeGenderCode_FEMALE}},
			args:     []expr.Expression{exprtest.Return(fhir.String("456"), fhir.UnsignedInt(789))},
			expected: system.Collection{fhir.ID("123"), &ppb.Patient_GenderCode{Value: cpb.AdministrativeGenderCode_FEMALE}, fhir.String("456"), fhir.UnsignedInt(789)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Combine(ctx, tc.input, tc.args...)
			if err != nil {
				t.Fatalf("Combine() returned an error: %v", err)
			}
			if !cmp.Equal(result, tc.expected, protocmp.Transform()) {
				t.Errorf("Combine() result diff (-got, +want):\n%s", cmp.Diff(tc.expected, result, protocmp.Transform()))
			}
		})
	}
}
