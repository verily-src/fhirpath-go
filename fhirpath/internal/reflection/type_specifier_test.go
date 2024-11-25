package reflection_test

import (
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/reflection"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

func TestTypeSpecifier_Is(t *testing.T) {
	testCases := []struct {
		name    string
		typeOne reflection.TypeSpecifier
		typeTwo reflection.TypeSpecifier
		want    system.Boolean
	}{
		{
			name:    "mismatched namespaces",
			typeOne: reflection.MustCreateTypeSpecifier("FHIR", "Element"),
			typeTwo: reflection.MustCreateTypeSpecifier("System", "Any"),
			want:    false,
		},
		{
			name:    "same type",
			typeOne: reflection.MustCreateTypeSpecifier("FHIR", "DomainResource"),
			typeTwo: reflection.MustCreateTypeSpecifier("FHIR", "DomainResource"),
			want:    true,
		},
		{
			name:    "child type is parent",
			typeOne: reflection.MustCreateTypeSpecifier("FHIR", "markdown"),
			typeTwo: reflection.MustCreateTypeSpecifier("FHIR", "string"),
			want:    true,
		},
		{
			name:    "child type is base",
			typeOne: reflection.MustCreateTypeSpecifier("FHIR", "markdown"),
			typeTwo: reflection.MustCreateTypeSpecifier("FHIR", "Element"),
			want:    true,
		},
		{
			name:    "check if type name is Resource",
			typeOne: reflection.MustCreateTypeSpecifier("FHIR", "Patient"),
			typeTwo: reflection.MustCreateTypeSpecifier("FHIR", "Resource"),
			want:    true,
		},
		{
			name:    "check if type name is Element",
			typeOne: reflection.MustCreateTypeSpecifier("FHIR", "Timing"),
			typeTwo: reflection.MustCreateTypeSpecifier("FHIR", "BackboneElement"),
			want:    true,
		},
		{
			name:    "Quantity is FHIR Element",
			typeOne: reflection.MustCreateTypeSpecifier("FHIR", "Quantity"),
			typeTwo: reflection.MustCreateTypeSpecifier("FHIR", "Element"),
			want:    true,
		},
		{
			name:    "System type is Any",
			typeOne: reflection.MustCreateTypeSpecifier("System", "Decimal"),
			typeTwo: reflection.MustCreateTypeSpecifier("System", "Any"),
			want:    true,
		},
		{
			name:    "integer types are Integers",
			typeOne: reflection.MustCreateTypeSpecifier("FHIR", "positiveInt"),
			typeTwo: reflection.MustCreateTypeSpecifier("FHIR", "integer"),
			want:    true,
		},
		{
			name:    "Patient is DomainResource",
			typeOne: reflection.MustCreateTypeSpecifier("FHIR", "Patient"),
			typeTwo: reflection.MustCreateTypeSpecifier("FHIR", "DomainResource"),
			want:    true,
		},
		{
			name:    "Mismatched types",
			typeOne: reflection.MustCreateTypeSpecifier("FHIR", "Patient"),
			typeTwo: reflection.MustCreateTypeSpecifier("FHIR", "Practitioner"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.typeOne.Is(tc.typeTwo); got != tc.want {
				t.Errorf("TypeSpecifier.Is returned unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestNewQualifiedTypeSpecifier_CreatesValidTS(t *testing.T) {
	testCases := []struct {
		name      string
		namespace string
		typeName  string
		want      reflection.TypeSpecifier
	}{
		{
			name:      "Valid FHIR primitive",
			namespace: "FHIR",
			typeName:  "decimal",
			want:      reflection.MustCreateTypeSpecifier("FHIR", "decimal"),
		},
		{
			name:      "Valid FHIR Element",
			namespace: "FHIR",
			typeName:  "ContactPoint",
			want:      reflection.MustCreateTypeSpecifier("FHIR", "ContactPoint"),
		},
		{
			name:      "Valid FHIR Resource",
			namespace: "FHIR",
			typeName:  "Medication",
			want:      reflection.MustCreateTypeSpecifier("FHIR", "Medication"),
		},
		{
			name:      "Valid System type",
			namespace: "System",
			typeName:  "DateTime",
			want:      reflection.MustCreateTypeSpecifier("System", "DateTime"),
		},
		{
			name:      "Valid Base type",
			namespace: "FHIR",
			typeName:  "Element",
			want:      reflection.MustCreateTypeSpecifier("FHIR", "Element"),
		},
		{
			name:      "Create DomainResource",
			namespace: "FHIR",
			typeName:  "DomainResource",
			want:      reflection.MustCreateTypeSpecifier("FHIR", "DomainResource"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := reflection.NewQualifiedTypeSpecifier(tc.namespace, tc.typeName)

			if err != nil {
				t.Fatalf("NewQualifiedTypeSpecifier(%s, %s) returned unexpected error: %v", tc.namespace, tc.typeName, err)
			}
			if got != tc.want {
				t.Errorf("NewQualifiedTypeSpecifier(%s, %s) returned incorrect TypeSpecifier: got %v, want %v", tc.namespace, tc.typeName, got, tc.want)
			}
		})
	}
}

func TestNewQualifiedTypeSpecifier_ReturnsError(t *testing.T) {
	testCases := []struct {
		name      string
		namespace string
		typeName  string
	}{
		{
			name:      "FHIR primitive with wrong case",
			namespace: "FHIR",
			typeName:  "Decimal",
		},
		{
			name:      "Non-existent FHIR type",
			namespace: "FHIR",
			typeName:  "Hospital",
		},
		{
			name:      "Mismatched namespace",
			namespace: "System",
			typeName:  "Medication",
		},
		{
			name:      "System type with wrong case",
			namespace: "System",
			typeName:  "dateTime",
		},
		{
			name:      "invalid namespace",
			namespace: "Enrichments",
			typeName:  "Engine",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := reflection.NewQualifiedTypeSpecifier(tc.namespace, tc.typeName); err == nil {
				t.Fatalf("NewQualifiedTypeSpecifier(%s, %s) didn't return error when expected to", tc.namespace, tc.typeName)
			}
		})
	}
}

func TestNewTypeSpecifier_CreatesValidTS(t *testing.T) {
	testCases := []struct {
		name     string
		typeName string
		want     reflection.TypeSpecifier
	}{
		{
			name:     "creates resource specifier from Patient",
			typeName: "Patient",
			want:     reflection.MustCreateTypeSpecifier("FHIR", "Patient"),
		},
		{
			name:     "creates element specifier from Element",
			typeName: "Element",
			want:     reflection.MustCreateTypeSpecifier("FHIR", "Element"),
		},
		{
			name:     "creates system specifier from Decimal",
			typeName: "Decimal",
			want:     reflection.MustCreateTypeSpecifier("System", "Decimal"),
		},
		{
			name:     "creates FHIR specifier from decimal",
			typeName: "decimal",
			want:     reflection.MustCreateTypeSpecifier("FHIR", "decimal"),
		},
		{
			name:     "creates FHIR specifier from Quantity",
			typeName: "Quantity",
			want:     reflection.MustCreateTypeSpecifier("FHIR", "Quantity"),
		},
		{
			name:     "Creates System specifier from Any",
			typeName: "Any",
			want:     reflection.MustCreateTypeSpecifier("System", "Any"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := reflection.NewTypeSpecifier(tc.typeName)

			if err != nil {
				t.Fatalf("NewTypeSpecifier(%s) returned unexpected error: %v", tc.name, err)
			}
			if got != tc.want {
				t.Errorf("NewTypeSpecifier(%s) returned incorrect TypeSpecifier: got %v, want %v", tc.typeName, got, tc.want)
			}
		})
	}
}

func TestNewTypeSpecifier_ReturnsError(t *testing.T) {
	if _, err := reflection.NewTypeSpecifier("notAType"); err == nil {
		t.Fatalf("NewTypeSpecifier(%s) didn't raise error on invalid type name", "notAType")
	}
}

func TestTypeOf_ReturnsTS(t *testing.T) {
	quantity, _ := system.ParseQuantity("123", "kg")

	testCases := []struct {
		name  string
		input any
		want  reflection.TypeSpecifier
	}{
		{
			name:  "Gets correct specifier for Patient type",
			input: (*ppb.Patient)(nil),
			want:  reflection.MustCreateTypeSpecifier("FHIR", "Patient"),
		},
		{
			name:  "Gets lowercase type name for primitive type",
			input: (*dtpb.Decimal)(nil),
			want:  reflection.MustCreateTypeSpecifier("FHIR", "decimal"),
		},
		{
			name:  "Gets correct specifier for Code type",
			input: (*ppb.Patient_GenderCode)(nil),
			want:  reflection.MustCreateTypeSpecifier("FHIR", "code"),
		},
		{
			name:  "Gets correct specifier for Oneof type",
			input: &ppb.Patient_DeceasedX{Choice: &ppb.Patient_DeceasedX_DateTime{DateTime: fhir.DateTimeNow()}},
			want:  reflection.MustCreateTypeSpecifier("FHIR", "dateTime"),
		},
		{
			name:  "Gets correct specifier for system type",
			input: quantity,
			want:  reflection.MustCreateTypeSpecifier("System", "Quantity"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := reflection.TypeOf(tc.input)

			if err != nil {
				t.Fatalf("GetTypeSpecifier returned unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("GetTypeSpecifier returned incorrect type specifier: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestGetTypeSpecifier_ReturnsError(t *testing.T) {
	if _, err := reflection.TypeOf("unsupported type"); err == nil {
		t.Fatalf("GetTypeSpecifier didn't return error for unsupported type")
	}
}
