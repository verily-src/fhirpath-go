package expr_test

import (
	"errors"
	"math"
	"testing"
	"time"

	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	bcrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/device_go_proto"
	epb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/encounter_go_proto"
	mrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medication_request_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
	"github.com/verily-src/fhirpath-go/fhirpath"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr/exprtest"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/reflection"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/bundle"
	"github.com/verily-src/fhirpath-go/internal/element/extension"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirconv"
	"github.com/verily-src/fhirpath-go/internal/slices"
	"google.golang.org/protobuf/testing/protocmp"
)

var (
	errMock = errors.New("some error")
)

func TestIdentityExpression_Input_EqualsOutput(t *testing.T) {
	identity := &expr.IdentityExpression{}
	testCases := []struct {
		name string
		data system.Collection
	}{
		{
			name: "Empty set",
			data: system.Collection{},
		},
		{
			name: "Mixed type set",
			data: system.Collection{fhir.String("test"), fhir.Integer(1), fhir.Integer(2)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := identity.Evaluate(&expr.Context{}, tc.data)

			if err != nil {
				t.Fatalf("Input: %s error when not expected, err: %v", tc.name, err)
			}
			if got, want := out, tc.data; !cmp.Equal(got, want, protocmp.Transform()) {
				t.Errorf("Input: %s, got: %v, want %v", tc.name, got, want)
			}
		})
	}
}

func TestFieldExpression_Gets_DesiredField(t *testing.T) {
	patientID := "123"
	patientFirstHumanName := &dtpb.HumanName{
		Given: []*dtpb.String{
			fhir.String("IU"),
			fhir.String("Amanda"),
		},
	}
	patientBirthDay := &dtpb.Date{
		ValueUs:   time.Now().UnixMicro(),
		Precision: dtpb.Date_DAY,
	}
	fullName := &dtpb.HumanName{
		Given: []*dtpb.String{
			fhir.String("Julius"),
		},
		Family: fhir.String("Caesar"),
	}
	patientContactPoint := []*dtpb.ContactPoint{
		{
			System: &dtpb.ContactPoint_SystemCode{
				Value: cpb.ContactPointSystemCode_PHONE,
			},
			Value: fhir.String("123-456-7890"),
			Rank:  fhir.PositiveInt(1),
		},
		{
			System: &dtpb.ContactPoint_SystemCode{
				Value: cpb.ContactPointSystemCode_EMAIL,
			},
			Value: fhir.String("example@gmail.com"),
			Rank:  fhir.PositiveInt(2),
		},
	}

	containedPatient := &bcrpb.ContainedResource{
		OneofResource: &bcrpb.ContainedResource_Patient{
			Patient: &ppb.Patient{
				Id:     fhir.ID(patientID),
				Active: fhir.Boolean(true),
				Gender: &ppb.Patient_GenderCode{
					Value: cpb.AdministrativeGenderCode_FEMALE,
				},
				Deceased: &ppb.Patient_DeceasedX{
					Choice: &ppb.Patient_DeceasedX_Boolean{
						Boolean: fhir.Boolean(true),
					},
				},
				MultipleBirth: &ppb.Patient_MultipleBirthX{
					Choice: &ppb.Patient_MultipleBirthX_Integer{
						Integer: fhir.Integer(int32(2)),
					},
				},
				Meta: &dtpb.Meta{
					Tag: []*dtpb.Coding{
						{
							Code: fhir.Code("#blessed"),
						},
					},
				},
				Name: []*dtpb.HumanName{
					patientFirstHumanName,
					fullName,
				},
				Telecom: patientContactPoint,
			},
		},
	}
	patientMissingName := &ppb.Patient{
		Id:        fhir.ID(patientID),
		Active:    fhir.Boolean(true),
		Telecom:   patientContactPoint,
		BirthDate: patientBirthDay,
	}
	patientWithOneName := &ppb.Patient{
		Id:   fhir.ID(patientID),
		Name: []*dtpb.HumanName{fullName},
	}

	testCases := []struct {
		name           string
		fieldExp       *expr.FieldExpression
		input          system.Collection
		wantCollection system.Collection
		wantErr        error
	}{
		{
			name:           "contained resource has collection in field",
			fieldExp:       &expr.FieldExpression{FieldName: "name"},
			input:          system.Collection{containedPatient},
			wantCollection: system.Collection{patientFirstHumanName, fullName},
		},
		{
			name:           "resource has empty list field",
			fieldExp:       &expr.FieldExpression{FieldName: "name"},
			input:          system.Collection{patientMissingName},
			wantCollection: system.Collection{},
		},
		{
			name:           "resource has empty non-list field",
			fieldExp:       &expr.FieldExpression{FieldName: "family"},
			input:          system.Collection{patientFirstHumanName},
			wantCollection: system.Collection{},
		},
		{
			name:           "field is appended with the suffix `value`",
			fieldExp:       &expr.FieldExpression{FieldName: "class"},
			input:          system.Collection{&epb.Encounter{ClassValue: fhir.Coding("class-system", "class-code")}},
			wantCollection: system.Collection{fhir.Coding("class-system", "class-code")},
		},
		{
			name:     "accessing non existent field",
			fieldExp: &expr.FieldExpression{FieldName: "version"},
			input:    system.Collection{containedPatient},
			wantErr:  expr.ErrInvalidField,
		},
		{
			name:           "resource has singleton in field",
			fieldExp:       &expr.FieldExpression{FieldName: "name"},
			input:          system.Collection{patientWithOneName},
			wantCollection: system.Collection{fullName},
		},
		{
			name:           "accessing non-list field",
			fieldExp:       &expr.FieldExpression{FieldName: "id"},
			input:          system.Collection{containedPatient},
			wantCollection: system.Collection{fhir.ID(patientID)},
		},
		{
			name:           "accessing family in HumanName element",
			fieldExp:       &expr.FieldExpression{FieldName: "family"},
			input:          system.Collection{fullName},
			wantCollection: system.Collection{fhir.String("Caesar")},
		},
		{
			name:           "accessing fields from multiple resources",
			fieldExp:       &expr.FieldExpression{FieldName: "name"},
			input:          system.Collection{containedPatient, patientWithOneName},
			wantCollection: system.Collection{patientFirstHumanName, fullName, fullName},
		},
		{
			name:           "accessing field of 2 words",
			fieldExp:       &expr.FieldExpression{FieldName: "birthDate"},
			input:          system.Collection{patientMissingName},
			wantCollection: system.Collection{patientBirthDay},
		},
		{
			name:           "(Legacy) input contains non-resource items",
			fieldExp:       &expr.FieldExpression{FieldName: "birthDate", Permissive: true},
			input:          system.Collection{patientMissingName, "hello"},
			wantCollection: system.Collection{patientBirthDay},
		},
		{
			name:     "input contains non-resource items",
			fieldExp: &expr.FieldExpression{FieldName: "birthDate"},
			input:    system.Collection{patientMissingName, "hello"},
			wantErr:  fhirpath.ErrInvalidField,
		},
		{
			name:           "accessing value field of primitive returns System primitive",
			fieldExp:       &expr.FieldExpression{FieldName: "value"},
			input:          system.Collection{fullName.Family},
			wantCollection: system.Collection{system.String("Caesar")},
		},
		{
			name:           "accessing value of code returns System primitive",
			fieldExp:       &expr.FieldExpression{FieldName: "value"},
			input:          system.Collection{patientContactPoint[0].System},
			wantCollection: system.Collection{system.String("phone")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.fieldExp.Evaluate(&expr.Context{}, tc.input)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("FieldExpression.Evaluate(%s) raised unexpected error: got %v, want %v", tc.input, err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.wantCollection, got, protocmp.Transform()); diff != "" {
				t.Errorf("FieldExpression.Evaluate(%s) returned unexpected diff (-want, +got):\n%s", tc.input, diff)
			}
		})
	}
}

func TestFieldExpression_ValidInput_GetsField(t *testing.T) {
	tm := time.Now()
	device1 := &device_go_proto.Device{
		Id: fhir.RandomID(),
	}
	patient1 := &patient_go_proto.Patient{
		Id: fhir.RandomID(),
	}
	entries := []*bundle_and_contained_resource_go_proto.Bundle_Entry{
		bundle.NewPostEntry(device1), bundle.NewPostEntry(patient1),
	}
	testCases := []struct {
		name  string
		input system.Collection
		field string
		want  system.Collection
	}{
		{
			name:  "Value field from primitive element",
			input: system.Collection{fhir.Integer(32)},
			field: "value",
			want:  system.Collection{system.Integer(32)},
		}, {
			name:  "Value field from valueX type",
			input: system.Collection{extension.New("url", fhir.Boolean(true))},
			field: "value",
			want:  system.Collection{fhir.Boolean(true)},
		}, {
			name: "Repeated field returns collection of repeated elements",
			input: system.Collection{bundle.NewCollection(
				bundle.WithEntries(entries...),
			)},
			field: "entry",
			want: system.Collection(
				slices.MustConvert[any](entries),
			),
		}, {
			name:  "Field on empty input returns empty",
			input: nil,
			field: "value",
			want:  system.Collection{},
		}, {
			name:  "Bundle entry contained resource",
			input: system.Collection{bundle.NewPostEntry(device1)},
			field: "resource",
			want:  system.Collection{device1},
		}, {
			name: "Reference using URI",
			input: system.Collection{&dtpb.Reference{
				Reference: &dtpb.Reference_Uri{
					Uri: fhir.String("https://some-url.com"),
				},
			}},
			field: "reference",
			want:  system.Collection{fhir.String("https://some-url.com")},
		}, {
			name: "Reference Patient without history ID",
			input: system.Collection{&dtpb.Reference{
				Reference: &dtpb.Reference_PatientId{
					PatientId: &dtpb.ReferenceId{
						Value: "12345",
					},
				},
			}},
			field: "reference",
			want:  system.Collection{fhir.String("Patient/12345")},
		}, {
			name: "Reference using fragment",
			input: system.Collection{&dtpb.Reference{
				Reference: &dtpb.Reference_Fragment{
					Fragment: fhir.String("p0"),
				},
			}},
			field: "reference",
			want:  system.Collection{fhir.String("#p0")},
		}, {
			name:  "Time using value_us",
			input: system.Collection{fhir.Time(tm)},
			field: "value",
			want:  system.Collection{system.String(fhirconv.TimeToString(fhir.Time(tm)))},
		}, {
			name:  "Date using value_us",
			input: system.Collection{fhir.Date(tm)},
			field: "value",
			want:  system.Collection{system.String(fhirconv.DateToString(fhir.Date(tm)))},
		}, {
			name:  "DateTime using value_us",
			input: system.Collection{fhir.DateTime(tm)},
			field: "value",
			want:  system.Collection{system.String(fhirconv.DateTimeToString(fhir.DateTime(tm)))},
		}, {
			name:  "Instant using value_us",
			input: system.Collection{fhir.Instant(tm)},
			field: "value",
			want:  system.Collection{system.String(fhirconv.InstantToString(fhir.Instant(tm)))},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut := &expr.FieldExpression{
				FieldName: tc.field,
			}

			got, err := sut.Evaluate(expr.InitializeContext(tc.input), tc.input)
			if err != nil {
				t.Fatalf("FieldExpression.Evaluate(%v): got unexpected err %v", tc.name, err)
			}

			if diff := cmp.Diff(got, tc.want, protocmp.Transform()); diff != "" {
				t.Errorf("FieldExpression.Evaluate(%v): (-got,+want):\n%v", tc.name, diff)
			}
		})
	}
}

func TestTypeExpression_Filters_DesiredType(t *testing.T) {
	medicationRequest := &mrpb.MedicationRequest{
		Reported: &mrpb.MedicationRequest_ReportedX{
			Choice: &mrpb.MedicationRequest_ReportedX_Boolean{
				Boolean: fhir.Boolean(true),
			},
		},
	}
	patient := &ppb.Patient{
		Id: fhir.ID("123"),
		Name: []*dtpb.HumanName{
			{
				Given: []*dtpb.String{
					fhir.String("Julius"),
				},
				Family: fhir.String("Caesar"),
			},
		},
	}

	testCases := []struct {
		name           string
		typeExp        *expr.TypeExpression
		input          system.Collection
		wantCollection system.Collection
		shouldError    bool
	}{
		{
			name:           "input contains only resource",
			typeExp:        &expr.TypeExpression{Type: "Patient"},
			input:          system.Collection{patient},
			wantCollection: system.Collection{patient},
		},
		{
			name:           "input contains more than desired resource",
			typeExp:        &expr.TypeExpression{Type: "Patient"},
			input:          system.Collection{patient, medicationRequest},
			wantCollection: system.Collection{patient},
		},
		{
			name:           "input doesn't contain desired resource",
			typeExp:        &expr.TypeExpression{Type: "MedicationRequest"},
			input:          system.Collection{patient},
			wantCollection: system.Collection{},
		},
		{
			name:           "desired type has 2 words",
			typeExp:        &expr.TypeExpression{Type: "MedicationRequest"},
			input:          system.Collection{medicationRequest, patient},
			wantCollection: system.Collection{medicationRequest},
		},
		{
			name:           "input collection has non-resource types with desired name",
			typeExp:        &expr.TypeExpression{Type: "Patient"},
			input:          system.Collection{"Patient", struct{ Patient string }{Patient: "Peter Griffin"}},
			wantCollection: system.Collection{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.typeExp.Evaluate(&expr.Context{}, tc.input)

			if err != nil {
				t.Fatalf("TypeExpression.Evaluate(%s) got unexpected error: %s", tc.input, err)
			}
			if diff := cmp.Diff(tc.wantCollection, got, protocmp.Transform()); diff != "" {
				t.Errorf("TypeExpression.Evaluate(%s) returned unexpected diff (-want, +got):\n%s", tc.input, diff)
			}
		})
	}
}

func TestExpressionSequence_EvaluatesInSequence(t *testing.T) {
	mock := &exprtest.MockExpression{
		Eval: func(ctx *expr.Context, input system.Collection) (system.Collection, error) {
			wasCalled := input[0].(int)
			return system.Collection{wasCalled + 1}, nil
		},
	}

	sequence := &expr.ExpressionSequence{[]expr.Expression{mock, mock, mock, mock}}
	result, err := sequence.Evaluate(&expr.Context{}, system.Collection{0})

	if err != nil {
		t.Fatalf("ExpressionSequence.Evaluate raised unexpected error: %v", err)
	}
	if got, want := result[0].(int), len(sequence.Expressions); got != want {
		t.Errorf("ExpressionSequence.Evaluate incorrectly accumulated values, got: %v, want: %v", got, want)
	}
}

func TestExpressionSequence_EvaluateRaisesError(t *testing.T) {
	sequence := &expr.ExpressionSequence{[]expr.Expression{exprtest.Return(), exprtest.Return(), exprtest.Error(errMock), exprtest.Return()}}

	_, err := sequence.Evaluate(&expr.Context{}, system.Collection{})

	if err == nil {
		t.Fatal("ExpressionSequence.Evaluate didn't raise error when expected")
	}
}

func TestLiteralExpression_EvaluateReturnsLiteral(t *testing.T) {
	testCases := []struct {
		name            string
		literalExp      *expr.LiteralExpression
		inputCollection system.Collection
		wantCollection  system.Collection
	}{
		{
			name:            "string literal expression without input",
			literalExp:      &expr.LiteralExpression{system.String("string")},
			inputCollection: system.Collection{"some input"},
			wantCollection:  system.Collection{system.String("string")},
		},
		{
			name:            "null literal expression",
			literalExp:      &expr.LiteralExpression{},
			inputCollection: system.Collection{},
			wantCollection:  system.Collection{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.literalExp.Evaluate(&expr.Context{}, tc.inputCollection)

			if err != nil {
				t.Fatalf("LiteralExpression.Evaluate raised unexpected error %v", err)
			}
			if diff := cmp.Diff(tc.wantCollection, got); diff != "" {
				t.Errorf("LiteralExpression.Evaluate returned unexpected diff: (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestIndexExpression_ReturnsIndex(t *testing.T) {
	givenName := []*dtpb.String{
		fhir.String("He"),
		fhir.String("Who"),
		fhir.String("Must"),
		fhir.String("Not"),
		fhir.String("Be"),
	}
	familyName := fhir.String("Named")
	voldemort := []*dtpb.HumanName{
		{
			Given:  givenName,
			Family: familyName,
		},
		{
			Given: []*dtpb.String{
				fhir.String("Tom"),
				fhir.String("Marvolo"),
			},
			Family: fhir.String("Riddle"),
		},
	}

	testCases := []struct {
		name            string
		indexExpr       *expr.IndexExpression
		inputCollection system.Collection
		wantCollection  system.Collection
	}{
		{
			name:            "indexing a proto string array",
			indexExpr:       &expr.IndexExpression{&expr.LiteralExpression{system.Integer(1)}},
			inputCollection: slices.MustConvert[any](givenName),
			wantCollection:  system.Collection{fhir.String("Who")},
		},
		{
			name:            "indexing a proto HumanName",
			indexExpr:       &expr.IndexExpression{&expr.LiteralExpression{system.Integer(0)}},
			inputCollection: slices.MustConvert[any](voldemort),
			wantCollection:  system.Collection{voldemort[0]},
		},
		{
			name:            "index out of bounds",
			indexExpr:       &expr.IndexExpression{&expr.LiteralExpression{system.Integer(5)}},
			inputCollection: slices.MustConvert[any](givenName),
			wantCollection:  system.Collection{},
		},
		{
			name:            "negative index",
			indexExpr:       &expr.IndexExpression{&expr.LiteralExpression{system.Integer(-1)}},
			inputCollection: slices.MustConvert[any](givenName),
			wantCollection:  system.Collection{},
		},
		{
			name:            "empty input collection",
			indexExpr:       &expr.IndexExpression{&expr.LiteralExpression{system.Integer(0)}},
			inputCollection: system.Collection{},
			wantCollection:  system.Collection{},
		},
		{
			name:            "index that evaluates to empty returns empty",
			indexExpr:       &expr.IndexExpression{exprtest.Return()},
			inputCollection: slices.MustConvert[any](givenName),
			wantCollection:  system.Collection{},
		},
	}

	for _, tc := range testCases {
		got, err := tc.indexExpr.Evaluate(&expr.Context{}, tc.inputCollection)

		if err != nil {
			t.Fatalf("IndexExpression.Evaluate raised unexpected error %v", err)
		}
		if diff := cmp.Diff(tc.wantCollection, got, protocmp.Transform()); diff != "" {
			t.Errorf("IndexExpression.Evaluate returned unexpected diff: (-want, +got)\n%s", diff)
		}
	}
}

func TestIndexExpression_EvaluateRaisesError(t *testing.T) {
	testCases := []struct {
		name            string
		indexExpr       *expr.IndexExpression
		inputCollection system.Collection
	}{
		{
			name:            "evaluate index raises an error",
			indexExpr:       &expr.IndexExpression{Index: exprtest.Error(errMock)},
			inputCollection: system.Collection{},
		},
		{
			name:            "index expression doesn't evaluate to a system type",
			indexExpr:       &expr.IndexExpression{Index: exprtest.Return(1)},
			inputCollection: system.Collection{},
		},
		{
			name:            "index doesn't evaluate to an integer",
			indexExpr:       &expr.IndexExpression{Index: exprtest.Return("not an integer")},
			inputCollection: system.Collection{},
		},
		{
			name:            "index expression evaluates to multiple entries",
			indexExpr:       &expr.IndexExpression{Index: exprtest.Return(system.Integer(1), system.Integer(2))},
			inputCollection: system.Collection{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.indexExpr.Evaluate(&expr.Context{}, tc.inputCollection)

			if err == nil {
				t.Fatalf("IndexExpression.Evaluate didn't return error when expected")
			}
		})
	}
}

func TestEqualityExpression_ReturnsResult(t *testing.T) {
	testCases := []struct {
		name            string
		inputCollection system.Collection
		equalityExpr    *expr.EqualityExpression
		wantCollection  system.Collection
	}{
		{
			name:            "one empty collection",
			inputCollection: system.Collection{},
			equalityExpr:    &expr.EqualityExpression{exprtest.Return(), exprtest.Return("one"), false},
			wantCollection:  system.Collection{},
		},
		{
			name:            "comparing with != operator",
			inputCollection: system.Collection{},
			equalityExpr:    &expr.EqualityExpression{exprtest.Return(system.String("abc")), exprtest.Return(system.String("abcd")), true},
			wantCollection:  system.Collection{system.Boolean(true)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.equalityExpr.Evaluate(&expr.Context{}, tc.inputCollection)

			if err != nil {
				t.Fatalf("EqualityExpression.Evaluate raised unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.wantCollection, got); diff != "" {
				t.Errorf("EqualityExpression.Evaluate returned unexpected diff: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestEqualityExpression_RaisesError(t *testing.T) {
	testCases := []struct {
		name         string
		equalityExpr *expr.EqualityExpression
	}{
		{
			name:         "subexpression one errors",
			equalityExpr: &expr.EqualityExpression{exprtest.Error(errMock), exprtest.Return(system.Boolean(true)), false},
		},
		{
			name:         "subexpression two errors",
			equalityExpr: &expr.EqualityExpression{exprtest.Return(system.Boolean(true)), exprtest.Error(errMock), false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.equalityExpr.Evaluate(&expr.Context{}, system.Collection{})

			if err == nil {
				t.Fatalf("EqualityExpression.Evaluate didn't propagate error when it should have")
			}
		})
	}
}

func TestIsExpression_ReturnsResult(t *testing.T) {
	testCases := []struct {
		name           string
		expr           *expr.IsExpression
		wantCollection system.Collection
	}{
		{
			name:           "returns boolean",
			expr:           &expr.IsExpression{exprtest.Return(fhir.String("str")), reflection.MustCreateTypeSpecifier("FHIR", "string")},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name:           "returns empty collection",
			expr:           &expr.IsExpression{exprtest.Return(), reflection.MustCreateTypeSpecifier("FHIR", "string")},
			wantCollection: system.Collection{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.expr.Evaluate(&expr.Context{}, []any{})

			if err != nil {
				t.Fatalf("IsExpression.Evaluate returned unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.wantCollection, got, protocmp.Transform()); diff != "" {
				t.Errorf("IsExpression.Evaluate returned unexpected diff: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestIsExpression_ReturnsError(t *testing.T) {
	testCases := []struct {
		name string
		expr *expr.IsExpression
	}{
		{
			name: "subexpression errors",
			expr: &expr.IsExpression{exprtest.Error(errors.New("some error")), reflection.TypeSpecifier{}},
		},
		{
			name: "subexpression evaluates to non-singleton",
			expr: &expr.IsExpression{exprtest.Return(system.Boolean(true), system.Boolean(true)), reflection.TypeSpecifier{}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.expr.Evaluate(&expr.Context{}, []any{})

			if err == nil {
				t.Fatalf("IsExpression.Evaluate doesn't return error when expected")
			}
		})
	}
}

func TestAsExpression_ReturnsResult(t *testing.T) {
	deceased := &ppb.Patient_DeceasedX{
		Choice: &ppb.Patient_DeceasedX_Boolean{
			Boolean: fhir.Boolean(true),
		},
	}
	testCases := []struct {
		name           string
		expr           *expr.AsExpression
		wantCollection system.Collection
	}{
		{
			name:           "input is of specified type (returns input)",
			expr:           &expr.AsExpression{exprtest.Return(fhir.Code("#blessed")), reflection.MustCreateTypeSpecifier("FHIR", "code")},
			wantCollection: system.Collection{fhir.Code("#blessed")},
		},
		{
			name:           "input is not of specified type (returns empty)",
			expr:           &expr.AsExpression{exprtest.Return(fhir.Integer(12)), reflection.MustCreateTypeSpecifier("FHIR", "string")},
			wantCollection: system.Collection{},
		},
		{
			name:           "input is empty collection (returns empty)",
			expr:           &expr.AsExpression{exprtest.Return(), reflection.MustCreateTypeSpecifier("FHIR", "string")},
			wantCollection: system.Collection{},
		},
		{
			name:           "input is a polymorphic oneOf type",
			expr:           &expr.AsExpression{exprtest.Return(deceased), reflection.MustCreateTypeSpecifier("FHIR", "boolean")},
			wantCollection: system.Collection{fhir.Boolean(true)},
		},
		{
			name:           "input is a system type",
			expr:           &expr.AsExpression{exprtest.Return(system.Boolean(true)), reflection.MustCreateTypeSpecifier("System", "Boolean")},
			wantCollection: system.Collection{system.Boolean(true)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.expr.Evaluate(&expr.Context{}, []any{})

			if err != nil {
				t.Fatalf("AsExpression.Evaluate returned unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.wantCollection, got, protocmp.Transform()); diff != "" {
				t.Errorf("AsExpression.Evaluate returned unexpected diff: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestAsExpression_ReturnsError(t *testing.T) {
	testCases := []struct {
		name string
		expr *expr.AsExpression
	}{
		{
			name: "subexpression errors",
			expr: &expr.AsExpression{exprtest.Error(errors.New("some error")), reflection.TypeSpecifier{}},
		},
		{
			name: "subexpression evaluates to non-singleton",
			expr: &expr.AsExpression{exprtest.Return(system.Boolean(true), system.Boolean(true)), reflection.TypeSpecifier{}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.expr.Evaluate(&expr.Context{}, []any{})

			if err == nil {
				t.Fatalf("AsExpression.Evaluate doesn't return error when expected")
			}
		})
	}
}

func TestBooleanExpression_ReturnsResult(t *testing.T) {
	testCases := []struct {
		name           string
		expr           *expr.BooleanExpression
		wantCollection system.Collection
	}{
		{
			name:           "[and]both expressions return true",
			expr:           &expr.BooleanExpression{exprtest.Return(system.Boolean(true)), exprtest.Return(system.Boolean(true)), expr.And},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name:           "[and]both expressions return booleans (one false)",
			expr:           &expr.BooleanExpression{exprtest.Return(system.Boolean(true)), exprtest.Return(system.Boolean(false)), expr.And},
			wantCollection: system.Collection{system.Boolean(false)},
		},
		{
			name:           "[and]both expressions return empty",
			expr:           &expr.BooleanExpression{exprtest.Return(), exprtest.Return(), expr.And},
			wantCollection: system.Collection{},
		},
		{
			name:           "[and]one expression is false while other is empty",
			expr:           &expr.BooleanExpression{exprtest.Return(system.Boolean(false)), exprtest.Return(), expr.And},
			wantCollection: system.Collection{system.Boolean(false)},
		},
		{
			name:           "[and]one expression is true while other is empty",
			expr:           &expr.BooleanExpression{exprtest.Return(system.Boolean(true)), exprtest.Return(), expr.And},
			wantCollection: system.Collection{},
		},
		{
			name:           "[and]singleton evaluation of collections (expressions are not booleans)",
			expr:           &expr.BooleanExpression{exprtest.Return(system.String("hi")), exprtest.Return(system.String("hello")), expr.And},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name:           "[or]both expressions return true",
			expr:           &expr.BooleanExpression{exprtest.Return(system.Boolean(true)), exprtest.Return(system.Boolean(true)), expr.Or},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name:           "[or]both expressions return booleans (one true, one false)",
			expr:           &expr.BooleanExpression{exprtest.Return(system.Boolean(true)), exprtest.Return(system.Boolean(false)), expr.Or},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name:           "[or]both expressions return empty",
			expr:           &expr.BooleanExpression{exprtest.Return(), exprtest.Return(), expr.Or},
			wantCollection: system.Collection{},
		},
		{
			name:           "[or]one expression is true while other is empty",
			expr:           &expr.BooleanExpression{exprtest.Return(system.Boolean(true)), exprtest.Return(), expr.Or},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name:           "[or]one expression is false while other is empty",
			expr:           &expr.BooleanExpression{exprtest.Return(system.Boolean(false)), exprtest.Return(), expr.Or},
			wantCollection: system.Collection{},
		},
		{
			name:           "[or]correctly compares proto booleans",
			expr:           &expr.BooleanExpression{exprtest.Return(fhir.Boolean(true)), exprtest.Return(fhir.Boolean(false)), expr.Or},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name:           "[or]singleton evaluation of collections (expressions are not booleans)",
			expr:           &expr.BooleanExpression{exprtest.Return(system.String("hi")), exprtest.Return(), expr.Or},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name:           "[xor]both expressions are equal booleans",
			expr:           &expr.BooleanExpression{exprtest.Return(system.Boolean(true)), exprtest.Return(system.Boolean(true)), expr.Xor},
			wantCollection: system.Collection{system.Boolean(false)},
		},
		{
			name:           "[xor] both expressions are inequal booleans",
			expr:           &expr.BooleanExpression{exprtest.Return(system.Boolean(false)), exprtest.Return(system.Boolean(true)), expr.Xor},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name:           "[xor]one collection is empty",
			expr:           &expr.BooleanExpression{exprtest.Return(), exprtest.Return(system.Boolean(true)), expr.Xor},
			wantCollection: system.Collection{},
		},
		{
			name:           "[implies]false implies false",
			expr:           &expr.BooleanExpression{exprtest.Return(system.Boolean(false)), exprtest.Return(system.Boolean(false)), expr.Implies},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name:           "[implies]false implies true",
			expr:           &expr.BooleanExpression{exprtest.Return(system.Boolean(false)), exprtest.Return(system.Boolean(true)), expr.Implies},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name:           "[implies]false implies empty",
			expr:           &expr.BooleanExpression{exprtest.Return(system.Boolean(false)), exprtest.Return(), expr.Implies},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name:           "[implies]true implies false",
			expr:           &expr.BooleanExpression{exprtest.Return(system.Boolean(true)), exprtest.Return(system.Boolean(false)), expr.Implies},
			wantCollection: system.Collection{system.Boolean(false)},
		},
		{
			name:           "[implies]true implies true",
			expr:           &expr.BooleanExpression{exprtest.Return(system.Boolean(true)), exprtest.Return(system.Boolean(true)), expr.Implies},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name:           "[implies]empty implies true",
			expr:           &expr.BooleanExpression{exprtest.Return(), exprtest.Return(system.Boolean(true)), expr.Implies},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name:           "[implies]returns empty",
			expr:           &expr.BooleanExpression{exprtest.Return(true), exprtest.Return(), expr.Implies},
			wantCollection: system.Collection{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.expr.Evaluate(&expr.Context{}, system.Collection{})

			if err != nil {
				t.Fatalf("BooleanExpression.Evaluate returned unexpected error: %v", err)
			}
			if !cmp.Equal(got, tc.wantCollection) {
				t.Errorf("BooleanExpression.Evaluate returned unexpected result: got %v, want %v", got, tc.wantCollection)
			}
		})
	}
}

func TestComparisonExpression_ReturnsResult(t *testing.T) {
	qty, _ := system.ParseQuantity("23.3", "kg")
	testErr := errors.New("some error")
	testCases := []struct {
		name           string
		expr           *expr.ComparisonExpression
		wantCollection system.Collection
		wantErr        error
	}{
		{
			name: "[Lt] evaluates string comparison",
			expr: &expr.ComparisonExpression{
				exprtest.Return(system.String("abc")),
				exprtest.Return(system.String("def")),
				expr.Lt,
			},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name: "returns empty if either collection is empty",
			expr: &expr.ComparisonExpression{
				exprtest.Return(),
				exprtest.Return(system.String("def")),
				expr.Lt,
			},
			wantCollection: system.Collection{},
		},
		{
			name: "[Gt] correctly compares Date values",
			expr: &expr.ComparisonExpression{
				exprtest.Return(system.MustParseDate("2023-07-22")),
				exprtest.Return(system.MustParseDate("2023-07-21")),
				expr.Gt,
			},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name: "[Gte] correctly compares Time value with different precision",
			expr: &expr.ComparisonExpression{
				exprtest.Return(system.MustParseTime("08:31")),
				exprtest.Return(system.MustParseTime("08:30:30")),
				expr.Gte,
			},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name: "[Lte] correctly compares quantity value",
			expr: &expr.ComparisonExpression{
				exprtest.Return(qty),
				exprtest.Return(qty),
				expr.Lte,
			},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name: "[Gte] returns empty for mismatched precision on DateTime",
			expr: &expr.ComparisonExpression{
				exprtest.Return(system.MustParseDateTime("2020-12-20T08:30")),
				exprtest.Return(system.MustParseDateTime("2020-12-20T08:30:05")),
				expr.Gte,
			},
			wantCollection: system.Collection{},
		},
		{
			name: "[Lt] correctly compares a Date with a DateTime",
			expr: &expr.ComparisonExpression{
				exprtest.Return(system.MustParseDate("2020-12-20")),
				exprtest.Return(system.MustParseDateTime("2020-12-21T")),
				expr.Lt,
			},
			wantCollection: system.Collection{system.Boolean(true)},
		},
		{
			name: "returns error for comparison of invalid types",
			expr: &expr.ComparisonExpression{
				exprtest.Return(system.String("100")),
				exprtest.Return(system.Integer(100)),
				expr.Gte,
			},
			wantErr: system.ErrTypeMismatch,
		},
		{
			name: "Propogates error from one of the expressions",
			expr: &expr.ComparisonExpression{
				exprtest.Return(),
				exprtest.Error(testErr),
				expr.Lt,
			},
			wantErr: testErr,
		},
		{
			name: "returns error if either collection is not a singleton",
			expr: &expr.ComparisonExpression{
				exprtest.Return(system.String("abc"), system.String("abc")),
				exprtest.Return(system.String("hi")),
				expr.Lt,
			},
			wantErr: expr.ErrNotSingleton,
		},
		{
			name: "returns error if a non-system type is returned",
			expr: &expr.ComparisonExpression{
				exprtest.Return(123),
				exprtest.Return(system.Integer(234)),
				expr.Lt,
			},
			wantErr: system.ErrCantBeCast,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.expr.Evaluate(&expr.Context{}, system.Collection{})

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("ComparisonExpression.Evaluate returned unexpected error: got %v, want %v", err, tc.wantErr)
			}
			if !cmp.Equal(got, tc.wantCollection) {
				t.Errorf("ComparisonExpression.Evaluate returned unexpected result: got %v, want %v", got, tc.wantCollection)
			}
		})
	}
}

func TestAndExpression_ReturnsError(t *testing.T) {
	testCases := []struct {
		name string
		expr *expr.BooleanExpression
	}{
		{
			name: "subexpression errors",
			expr: &expr.BooleanExpression{exprtest.Error(errMock), exprtest.Return(system.Boolean(true)), expr.And},
		},
		{
			name: "expression returns non-singleton",
			expr: &expr.BooleanExpression{exprtest.Return(system.Boolean(true)), exprtest.Return(system.Boolean(true), system.Boolean(true)), expr.And},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.expr.Evaluate(&expr.Context{}, system.Collection{})

			if err == nil {
				t.Fatalf("AndExpression.Evaluate didn't return error when expected")
			}
		})
	}
}

func TestArithmeticExpression_ReturnsResult(t *testing.T) {
	testCases := []struct {
		name    string
		expr    *expr.ArithmeticExpression
		want    system.Collection
		wantErr error
	}{
		{
			name: "adds two system types",
			expr: &expr.ArithmeticExpression{
				Left:  exprtest.Return(system.String("hello ")),
				Right: exprtest.Return(system.String("world")),
				Op:    expr.EvaluateAdd,
			},
			want: system.Collection{system.String("hello world")},
		},
		{
			name: "adds two proto types",
			expr: &expr.ArithmeticExpression{
				Left:  exprtest.Return(fhir.MustParseDate("2020-02-12")),
				Right: exprtest.Return(fhir.UCUMQuantity(1, "day")),
				Op:    expr.EvaluateAdd,
			},
			want: system.Collection{system.MustParseDate("2020-02-13")},
		},
		{
			name: "adds a integer with a decimal",
			expr: &expr.ArithmeticExpression{
				Left:  exprtest.Return(system.Integer(25)),
				Right: exprtest.Return(system.Decimal(decimal.NewFromFloat(1.24))),
				Op:    expr.EvaluateAdd,
			},
			want: system.Collection{system.Decimal(decimal.NewFromFloat(26.24))},
		},
		{
			name: "subtracts quantity from date",
			expr: &expr.ArithmeticExpression{
				Left:  exprtest.Return(system.MustParseDate("2020-12")),
				Right: exprtest.Return(system.MustParseQuantity("35", "days")),
				Op:    expr.EvaluateSub,
			},
			want: system.Collection{system.MustParseDate("2020-11")},
		},
		{
			name: "returns empty if either collection is empty",
			expr: &expr.ArithmeticExpression{
				Left:  exprtest.Return(),
				Right: exprtest.Return(system.String("hell0")),
				Op:    expr.EvaluateAdd,
			},
			want: system.Collection{},
		},
		{
			name: "returns error if either collection has multiple elements",
			expr: &expr.ArithmeticExpression{
				Left:  exprtest.Return(system.Integer(1), system.Integer(2)),
				Right: exprtest.Return(system.Integer(2)),
				Op:    expr.EvaluateAdd,
			},
			wantErr: expr.ErrNotSingleton,
		},
		{
			name: "returns error if types can not be added",
			expr: &expr.ArithmeticExpression{
				Left:  exprtest.Return(system.Integer(2)),
				Right: exprtest.Return(system.Boolean(true)),
				Op:    expr.EvaluateAdd,
			},
			wantErr: system.ErrTypeMismatch,
		},
		{
			name: "returns empty if integer addition overflows",
			expr: &expr.ArithmeticExpression{
				Left:  exprtest.Return(system.Integer(math.MaxInt32)),
				Right: exprtest.Return(system.Integer(1)),
				Op:    expr.EvaluateAdd,
			},
			want: system.Collection{},
		},
		{
			name: "returns empty if integer subtraction overflows",
			expr: &expr.ArithmeticExpression{
				Left:  exprtest.Return(system.Integer(math.MinInt32)),
				Right: exprtest.Return(system.Integer(1)),
				Op:    expr.EvaluateSub,
			},
			want: system.Collection{},
		},
		{
			name: "multiplies decimals together",
			expr: &expr.ArithmeticExpression{
				Left:  exprtest.Return(system.Decimal(decimal.NewFromFloat(0.25))),
				Right: exprtest.Return(system.Decimal(decimal.NewFromFloat(0.25))),
				Op:    expr.EvaluateMul,
			},
			want: system.Collection{system.Decimal(decimal.NewFromFloat(0.0625))},
		},
		{
			name: "implicitly converts integer to decimal",
			expr: &expr.ArithmeticExpression{
				Left:  exprtest.Return(system.Decimal(decimal.NewFromFloat(1.2))),
				Right: exprtest.Return(system.Integer(2)),
				Op:    expr.EvaluateMul,
			},
			want: system.Collection{system.Decimal(decimal.NewFromFloat(2.4))},
		},
		{
			name: "multiplies integers together",
			expr: &expr.ArithmeticExpression{
				Left:  exprtest.Return(system.Integer(12)),
				Right: exprtest.Return(system.Integer(12)),
				Op:    expr.EvaluateMul,
			},
			want: system.Collection{system.Integer(144)},
		},
		{
			name: "returns empty if multiplication causes integer overflow",
			expr: &expr.ArithmeticExpression{
				Left:  exprtest.Return(system.Integer(math.MaxInt32)),
				Right: exprtest.Return(system.Integer(2)),
				Op:    expr.EvaluateMul,
			},
			want: system.Collection{},
		},
		{
			name: "returns error on a type mismatch",
			expr: &expr.ArithmeticExpression{
				Left:  exprtest.Return(system.Integer(2)),
				Right: exprtest.Return(system.String("a")),
				Op:    expr.EvaluateMul,
			},
			wantErr: system.ErrTypeMismatch,
		},
		{
			name: "returns decimal on integer division",
			expr: &expr.ArithmeticExpression{
				Left:  exprtest.Return(system.Integer(5)),
				Right: exprtest.Return(system.Integer(2)),
				Op:    expr.EvaluateDiv,
			},
			want: system.Collection{system.Decimal(decimal.NewFromFloat(2.5))},
		},
		{
			name: "performs floor division",
			expr: &expr.ArithmeticExpression{
				Left:  exprtest.Return(system.Integer(5)),
				Right: exprtest.Return(system.Integer(2)),
				Op:    expr.EvaluateFloorDiv,
			},
			want: system.Collection{system.Integer(2)},
		},
		{
			name: "performs mod between decimals",
			expr: &expr.ArithmeticExpression{
				Left:  exprtest.Return(system.Decimal(decimal.NewFromFloat(5.5))),
				Right: exprtest.Return(system.Decimal(decimal.NewFromFloat(0.7))),
				Op:    expr.EvaluateMod,
			},
			want: system.Collection{system.Decimal(decimal.NewFromFloat(0.6))},
		},
		{
			name: "performs mod between integers",
			expr: &expr.ArithmeticExpression{
				Left:  exprtest.Return(system.Integer(19)),
				Right: exprtest.Return(system.Integer(9)),
				Op:    expr.EvaluateMod,
			},
			want: system.Collection{system.Integer(1)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.expr.Evaluate(&expr.Context{}, []any{})

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("ArithmeticExpression.Evaluate returned unexpected error: got %v, want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("ArithmeticExpression.Evaluate returned unexpected diff: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestConcatExpression_AddsStrings(t *testing.T) {
	testCases := []struct {
		name    string
		expr    *expr.ConcatExpression
		want    system.Collection
		wantErr error
	}{
		{
			name: "concatenates two strings",
			expr: &expr.ConcatExpression{
				Left:  exprtest.Return(system.String("hello ")),
				Right: exprtest.Return(system.String("world")),
			},
			want: system.Collection{system.String("hello world")},
		},
		{
			name: "concatenates fhir string with system string",
			expr: &expr.ConcatExpression{
				Left:  exprtest.Return(fhir.String("abc")),
				Right: exprtest.Return(system.String("def")),
			},
			want: system.Collection{system.String("abcdef")},
		},
		{
			name: "if left collection is empty, returns the other string",
			expr: &expr.ConcatExpression{
				Left:  exprtest.Return(system.String("Hello")),
				Right: exprtest.Return(),
			},
			want: system.Collection{system.String("Hello")},
		},
		{
			name: "if right collection is empty, returns the other string",
			expr: &expr.ConcatExpression{
				Left:  exprtest.Return(),
				Right: exprtest.Return(system.String("Hello")),
			},
			want: system.Collection{system.String("Hello")},
		},
		{
			name: "returns an error if collection has multiple elements",
			expr: &expr.ConcatExpression{
				Left:  exprtest.Return(system.String("1"), system.String("2")),
				Right: exprtest.Return(system.String("3")),
			},
			wantErr: expr.ErrNotSingleton,
		},
		{
			name: "returns an error if a string is not returned",
			expr: &expr.ConcatExpression{
				Left:  exprtest.Return(system.Integer(1)),
				Right: exprtest.Return(system.String("3")),
			},
			wantErr: expr.ErrInvalidType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.expr.Evaluate(&expr.Context{}, system.Collection{})

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("ConcatExpression.Evaluate returned unexpected error: got %v, want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("ConcatExpression.Evaluate returned unexpected diff: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestExternalConstantExpression(t *testing.T) {
	testCases := []struct {
		name    string
		expr    *expr.ExternalConstantExpression
		context *expr.Context
		want    system.Collection
		wantErr error
	}{
		{
			name: "returns constant",
			expr: &expr.ExternalConstantExpression{Identifier: "value"},
			context: &expr.Context{
				ExternalConstants: map[string]any{"value": system.String("some string")},
			},
			want: system.Collection{system.String("some string")},
		},
		{
			name: "returns error if constant doesn't exist",
			expr: &expr.ExternalConstantExpression{Identifier: "value"},
			context: &expr.Context{
				ExternalConstants: map[string]any{},
			},
			wantErr: expr.ErrConstantNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.expr.Evaluate(tc.context, system.Collection{})

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("ExternalConstantExpression.Evaluate returned unexpected error: got %v, want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("ExternalConstantExpression.Evaluate returned unexpected diff: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestNegationExpression(t *testing.T) {
	testCases := []struct {
		name    string
		expr    *expr.NegationExpression
		want    system.Collection
		wantErr error
	}{
		{
			name: "negates integer",
			expr: &expr.NegationExpression{Expr: exprtest.Return(system.Integer(4))},
			want: system.Collection{system.Integer(-4)},
		},
		{
			name: "negates decimal",
			expr: &expr.NegationExpression{Expr: exprtest.Return(system.Decimal(decimal.NewFromFloat(1.5)))},
			want: system.Collection{system.Decimal(decimal.NewFromFloat(-1.5))},
		},
		{
			name: "negates quantity",
			expr: &expr.NegationExpression{Expr: exprtest.Return(system.MustParseQuantity("2.5", "kg"))},
			want: system.Collection{system.MustParseQuantity("-2.5", "kg")},
		},
		{
			name: "negates proto integer",
			expr: &expr.NegationExpression{Expr: exprtest.Return(fhir.Integer(-1))},
			want: system.Collection{system.Integer(1)},
		},
		{
			name:    "raises error on negating a collection",
			expr:    &expr.NegationExpression{Expr: exprtest.Return(system.Integer(1), system.Integer(2))},
			wantErr: expr.ErrNotSingleton,
		},
		{
			name:    "raises error if a non-number type is negated",
			expr:    &expr.NegationExpression{Expr: exprtest.Return(system.String("1"))},
			wantErr: expr.ErrInvalidType,
		},
		{
			name:    "raises error if complex type is negated",
			expr:    &expr.NegationExpression{Expr: exprtest.Return(fhir.Ratio(20, 10))},
			wantErr: expr.ErrInvalidType,
		},
		{
			name: "passes through empty collection",
			expr: &expr.NegationExpression{Expr: exprtest.Return()},
			want: system.Collection{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.expr.Evaluate(&expr.Context{}, system.Collection{})

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("NegationExpression.Evaluate returned unexpected error: got %v, want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("NegationExpression.Evaluate returned unexpected diff: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestUnionExpression(t *testing.T) {
	testCases := []struct {
		name    string
		expr    *expr.UnionExpression
		want    system.Collection
		wantErr error
	}{
		{
			name: "unions two distinct collections",
			expr: &expr.UnionExpression{
				Left:  exprtest.Return(system.Integer(1), system.Integer(2)),
				Right: exprtest.Return(system.Integer(3), system.Integer(4)),
			},
			want: system.Collection{system.Integer(1), system.Integer(2), system.Integer(3), system.Integer(4)},
		},
		{
			name: "unions two collections with overlap",
			expr: &expr.UnionExpression{
				Left:  exprtest.Return(system.Integer(1), system.Integer(2)),
				Right: exprtest.Return(system.Integer(2), system.Integer(3)),
			},
			want: system.Collection{system.Integer(1), system.Integer(2), system.Integer(3)},
		},
		{
			name: "unions with an empty collection",
			expr: &expr.UnionExpression{
				Left:  exprtest.Return(system.Integer(1), system.Integer(2)),
				Right: exprtest.Return(),
			},
			want: system.Collection{system.Integer(1), system.Integer(2)},
		},
		{
			name: "unions two empty collections",
			expr: &expr.UnionExpression{
				Left:  exprtest.Return(),
				Right: exprtest.Return(),
			},
			want: system.Collection{},
		},
		{
			name: "propagates error from left expression",
			expr: &expr.UnionExpression{
				Left:  exprtest.Error(errMock),
				Right: exprtest.Return(system.Integer(1)),
			},
			wantErr: errMock,
		},
		{
			name: "propagates error from right expression",
			expr: &expr.UnionExpression{
				Left:  exprtest.Return(system.Integer(1)),
				Right: exprtest.Error(errMock),
			},
			wantErr: errMock,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.expr.Evaluate(&expr.Context{}, system.Collection{})

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("UnionExpression.Evaluate returned unexpected error: got %v, want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("UnionExpression.Evaluate returned unexpected diff: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestMembershipExpression(t *testing.T) {
	testCases := []struct {
		name    string
		expr    *expr.MembershipExpression
		want    system.Collection
		wantErr error
	}{
		{
			name: "item is in collection",
			expr: &expr.MembershipExpression{
				Left:     exprtest.Return(system.Integer(1)),
				Right:    exprtest.Return(system.Integer(1), system.Integer(2)),
				Operator: expr.In,
			},
			want: system.Collection{system.Boolean(true)},
		},
		{
			name: "item is not in collection",
			expr: &expr.MembershipExpression{
				Left:     exprtest.Return(system.Integer(3)),
				Right:    exprtest.Return(system.Integer(1), system.Integer(2)),
				Operator: expr.In,
			},
			want: system.Collection{system.Boolean(false)},
		},
		{
			name: "left operand is empty",
			expr: &expr.MembershipExpression{
				Left:     exprtest.Return(),
				Right:    exprtest.Return(system.Integer(1), system.Integer(2)),
				Operator: expr.In,
			},
			want: system.Collection{},
		},
		{
			name: "left operand is not a singleton",
			expr: &expr.MembershipExpression{
				Left:     exprtest.Return(system.Integer(1), system.Integer(2)),
				Right:    exprtest.Return(system.Integer(1), system.Integer(2)),
				Operator: expr.In,
			},
			wantErr: expr.ErrNotSingleton,
		},
		{
			name: "collection contains item",
			expr: &expr.MembershipExpression{
				Left:     exprtest.Return(system.Integer(1), system.Integer(2)),
				Right:    exprtest.Return(system.Integer(1)),
				Operator: expr.Contains,
			},
			want: system.Collection{system.Boolean(true)},
		},
		{
			name: "collection does not contain item",
			expr: &expr.MembershipExpression{
				Left:     exprtest.Return(system.Integer(1), system.Integer(2)),
				Right:    exprtest.Return(system.Integer(3)),
				Operator: expr.Contains,
			},
			want: system.Collection{system.Boolean(false)},
		},
		{
			name: "right operand is empty",
			expr: &expr.MembershipExpression{
				Left:     exprtest.Return(system.Integer(1), system.Integer(2)),
				Right:    exprtest.Return(),
				Operator: expr.Contains,
			},
			want: system.Collection{},
		},
		{
			name: "right operand is not a singleton",
			expr: &expr.MembershipExpression{
				Left:     exprtest.Return(system.Integer(1), system.Integer(2)),
				Right:    exprtest.Return(system.Integer(1), system.Integer(2)),
				Operator: expr.Contains,
			},
			wantErr: expr.ErrNotSingleton,
		},
		{
			name: "propagates error from left expression",
			expr: &expr.MembershipExpression{
				Left:     exprtest.Error(errMock),
				Right:    exprtest.Return(system.Integer(1)),
				Operator: expr.In,
			},
			wantErr: errMock,
		},
		{
			name: "propagates error from right expression",
			expr: &expr.MembershipExpression{
				Left:     exprtest.Return(system.Integer(1)),
				Right:    exprtest.Error(errMock),
				Operator: expr.In,
			},
			wantErr: errMock,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.expr.Evaluate(&expr.Context{}, system.Collection{})

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("MembershipExpression.Evaluate returned unexpected error: got %v, want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("MembershipExpression.Evaluate returned unexpected diff: (-want, +got)\n%s", diff)
			}
		})
	}
}
