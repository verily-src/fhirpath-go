package system_test

import (
	"testing"

	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	mrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/medication_request_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	qpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/questionnaire_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/element/canonical"
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

var (
	date, _         = system.ParseDate("2012-12-31")
	time, _         = system.ParseTime("08:30:05")
	dateTime, _     = system.ParseDateTime("2002-04-20T08:30:00Z")
	quantity, _     = system.ParseQuantity("1.234", "m")
	zeroQuantity, _ = system.ParseQuantity("0", "m")
)

type testCase struct {
	name       string
	input      any
	want       system.Any
	shouldCast bool
}

var testCases []testCase = []testCase{
	{
		name:       "converts Boolean",
		input:      fhir.Boolean(true),
		want:       system.Boolean(true),
		shouldCast: true,
	},
	{
		name:       "converts String",
		input:      fhir.String("string"),
		want:       system.String("string"),
		shouldCast: true,
	},
	{
		name:       "converts Canonical",
		input:      canonical.New("string"),
		want:       system.String("string"),
		shouldCast: true,
	},
	{
		name:       "converts Uri",
		input:      fhir.URI("Uri"),
		want:       system.String("Uri"),
		shouldCast: true,
	},
	{
		name:       "converts Url",
		input:      fhir.URL("Url"),
		want:       system.String("Url"),
		shouldCast: true,
	},
	{
		name:       "converts Oid",
		input:      fhir.OID("Oid"),
		want:       system.String("urn:oid:Oid"),
		shouldCast: true,
	},
	{
		name:       "converts Id",
		input:      fhir.ID("Id"),
		want:       system.String("Id"),
		shouldCast: true,
	},
	{
		name:       "converts code",
		input:      fhir.Code("Code"),
		want:       system.String("Code"),
		shouldCast: true,
	},
	{
		name:       "converts uuid",
		input:      fhir.UUID("Uuid"),
		want:       system.String("urn:uuid:Uuid"),
		shouldCast: true,
	},
	{
		name:       "converts markdown",
		input:      fhir.Markdown("Markdown"),
		want:       system.String("Markdown"),
		shouldCast: true,
	},
	{
		name:       "converts base64 binary",
		input:      fhir.Base64Binary([]byte("hello world")),
		want:       system.String("aGVsbG8gd29ybGQ="),
		shouldCast: true,
	},
	{
		name:       "converts integer",
		input:      fhir.Integer(123),
		want:       system.Integer(123),
		shouldCast: true,
	},
	{
		name:       "converts positive integer",
		input:      fhir.PositiveInt(212),
		want:       system.Integer(212),
		shouldCast: true,
	},
	{
		name:       "converts unsigned integer",
		input:      fhir.UnsignedInt(10000),
		want:       system.Integer(10000),
		shouldCast: true,
	},
	{
		name:       "converts decimal",
		input:      fhir.Decimal(1.234),
		want:       system.Decimal(decimal.NewFromFloat(1.234)),
		shouldCast: true,
	},
	{
		name:       "converts date",
		input:      fhir.MustParseDate("2012-12-31"),
		want:       date,
		shouldCast: true,
	},
	{
		name:       "converts time",
		input:      fhir.MustParseTime("08:30:05"),
		want:       time,
		shouldCast: true,
	},
	{
		name:       "converts dateTime",
		input:      fhir.MustParseDateTime("2002-04-20T08:30:00Z"),
		want:       dateTime,
		shouldCast: true,
	},
	{
		name:       "converts instant",
		input:      fhir.MustParseInstant("2002-04-20T08:30:00Z"),
		want:       dateTime,
		shouldCast: true,
	},
	{
		name:       "converts quantity",
		input:      fhir.UCUMQuantity(float64(1.234), "m"),
		want:       quantity,
		shouldCast: true,
	},
	{
		name: "converts quantity without value",
		input: &dtpb.Quantity{
			Unit: fhir.String("m"),
		},
		want:       zeroQuantity,
		shouldCast: true,
	},
	{
		name:       "passes through system type",
		input:      system.String("pass through"),
		want:       system.String("pass through"),
		shouldCast: true,
	},
	{
		name:       "doesn't cast complex type",
		input:      &ppb.Patient{},
		shouldCast: false,
	},
	{
		name:       "converts gender code",
		input:      &ppb.Patient_GenderCode{Value: cpb.AdministrativeGenderCode_MALE},
		want:       system.String("male"),
		shouldCast: true,
	},
	{
		name:       "converts PublicationStatus",
		input:      &qpb.Questionnaire_StatusCode{Value: cpb.PublicationStatusCode_RETIRED},
		want:       system.String("retired"),
		shouldCast: true,
	},
	{
		name:       "converts IntentCode",
		input:      &mrpb.MedicationRequest_IntentCode{Value: cpb.MedicationRequestIntentCode_ORIGINAL_ORDER},
		want:       system.String("original-order"),
		shouldCast: true,
	},
	{
		name:       "converts a code with a string value",
		input:      &dtpb.Attachment_ContentTypeCode{Value: "image"},
		want:       system.String("image"),
		shouldCast: true,
	},
	{
		name:       "converts Priority",
		input:      &mrpb.MedicationRequest_PriorityCode{Value: cpb.RequestPriorityCode_URGENT},
		want:       system.String("urgent"),
		shouldCast: true,
	},
	{
		name:       "doesn't cast non-code type with value field",
		input:      &dtpb.ContactPoint{Value: fhir.String("123-456-7890")},
		shouldCast: false,
	},
	{
		name:       "doesn't cast non-fhir type",
		input:      12,
		shouldCast: false,
	},
}

func TestIsPrimitive_ReturnsBoolean(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ok := system.IsPrimitive(tc.input)

			if ok != tc.shouldCast {
				t.Errorf("system.IsPrimitive(%v) returns unexpected result, casts: %v, shouldCast: %v", tc.input, ok, tc.shouldCast)
			}
		})
	}
}

func TestFrom_CastsCorrectly(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.shouldCast {
				got, err := system.From(tc.input)

				if err != nil {
					t.Fatalf("system.From(%v) returns unexpected error: %v", tc.input, err)
				}
				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("system.From(%v) returns unexpected diff: (-want, +got)\n%s", tc.input, diff)
				}
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	wantQuantity, _ := system.ParseQuantity("4", "m")
	wantDateTime, _ := system.ParseDateTime("2012-12-31T")
	testCases := []struct {
		name string
		from system.Any
		to   system.Any
		want system.Any
	}{
		{
			name: "converts integer to decimal",
			from: system.Integer(16),
			to:   system.Decimal(decimal.NewFromInt32(20)),
			want: system.Decimal(decimal.NewFromInt32(16)),
		},
		{
			name: "converts decimal to quantity",
			from: system.Decimal(decimal.NewFromFloat(1.234)),
			to:   quantity,
			want: quantity,
		},
		{
			name: "converts integer to quantity",
			from: system.Integer(4),
			to:   quantity,
			want: wantQuantity,
		},
		{
			name: "converts Date to DateTime",
			from: date,
			to:   dateTime,
			want: wantDateTime,
		},
		{
			name: "passes through types that can't be converted",
			from: system.String("2012-12-31"),
			to:   date,
			want: system.String("2012-12-31"),
		},
	}

	for _, tc := range testCases {
		got := system.Normalize(tc.from, tc.to)

		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("system.Normalize(%T, %T) returns unexpected diff: (-want, +got)\n%s", tc.from, tc.to, diff)
		}
	}
}
