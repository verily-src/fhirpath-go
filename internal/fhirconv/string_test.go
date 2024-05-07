package fhirconv_test

import (
	"fmt"
	"testing"
	"time"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirconv"
)

func TestToString_String_ReturnsValue(t *testing.T) {
	want := "Hello world"
	str := fhir.String(want)

	got := fhirconv.ToString(str)

	if got != want {
		t.Errorf("String[String]: got %v, want %v", got, want)
	}
}

func TestToString_Code_ReturnsValue(t *testing.T) {
	want := "Hello world"
	str := fhir.Code(want)

	got := fhirconv.ToString(str)

	if got != want {
		t.Errorf("String[Code]: got %v, want %v", got, want)
	}
}

func TestToString_ID_ReturnsValue(t *testing.T) {
	want := "0xdeadbeef"
	str := fhir.ID(want)

	got := fhirconv.ToString(str)

	if got != want {
		t.Errorf("String[Id]: got %v, want %v", got, want)
	}
}

func TestToString_Markdown_ReturnsValue(t *testing.T) {
	want := "**This** is markdown"
	str := fhir.Markdown(want)

	got := fhirconv.ToString(str)

	if got != want {
		t.Errorf("String[Markdown]: got %v, want %v", got, want)
	}
}

func TestToString_URI_ReturnsValue(t *testing.T) {
	want := "https://example.com"
	str := fhir.URI(want)

	got := fhirconv.ToString(str)

	if got != want {
		t.Errorf("String[Uri]: got %v, want %v", got, want)
	}
}

func TestToString_URL_ReturnsValue(t *testing.T) {
	want := "https://example.com"
	str := fhir.URL(want)

	got := fhirconv.ToString(str)

	if got != want {
		t.Errorf("String[Url]: got %v, want %v", got, want)
	}
}

func TestToString_OID_ReturnsValue(t *testing.T) {
	want := "123456"
	str := fhir.OID(want)

	got := fhirconv.ToString(str)

	if want := fmt.Sprintf("urn:oid:%v", want); got != want {
		t.Errorf("String[Oid]: got %v, want %v", got, want)
	}
}

func TestToString_UUID_ReturnsValue(t *testing.T) {
	want := "25674bb5-6153-4dd7-9d4e-00fecc9058f1"
	str := fhir.UUID(want)

	got := fhirconv.ToString(str)

	if want := fmt.Sprintf("urn:uuid:%v", want); got != want {
		t.Errorf("String[Uuid]: got %v, want %v", got, want)
	}
}

func TestToString_Boolean_ReturnsValue(t *testing.T) {
	want := "true"
	str := fhir.Boolean(true)

	got := fhirconv.ToString(str)

	if got != want {
		t.Errorf("String[Boolean]: got %v, want %v", got, want)
	}
}

func TestToString_PositiveInt_ReturnsValue(t *testing.T) {
	want := "42"
	str := fhir.PositiveInt(42)

	got := fhirconv.ToString(str)

	if got != want {
		t.Errorf("String[PositiveInt]: got %v, want %v", got, want)
	}
}

func TestToString_UnsignedInt_ReturnsValue(t *testing.T) {
	want := "42"
	str := fhir.UnsignedInt(42)

	got := fhirconv.ToString(str)

	if got != want {
		t.Errorf("String[UnsignedInt]: got %v, want %v", got, want)
	}
}

func TestToString_Integer_ReturnsValue(t *testing.T) {
	want := "42"
	str := fhir.Integer(42)

	got := fhirconv.ToString(str)

	if got != want {
		t.Errorf("String[Integer]: got %v, want %v", got, want)
	}
}

func TestToString_Base64Binary_ReturnsValue(t *testing.T) {
	want := "3q2+7w=="
	str := fhir.Base64Binary([]byte{0xde, 0xad, 0xbe, 0xef})

	got := fhirconv.ToString(str)

	if got != want {
		t.Errorf("String[Base64Binary]: got %v, want %v", got, want)
	}
}

func TestToString_Instant_ReturnsValue(t *testing.T) {
	testCases := []struct {
		name      string
		want      string
		precision dtpb.Instant_Precision
	}{
		{"UnspecifiedPrecision", "1970-01-01T00:00:32.000000+00:00", dtpb.Instant_PRECISION_UNSPECIFIED},
		{"MicrosecondPrecision", "1970-01-01T00:00:32.000000+00:00", dtpb.Instant_MICROSECOND},
		{"MillisecondPrecision", "1970-01-01T00:00:32.000+00:00", dtpb.Instant_MILLISECOND},
		{"SecondPrecision", "1970-01-01T00:00:32+00:00", dtpb.Instant_SECOND},
	}
	timestamp := int64(time.Second * 32 / time.Microsecond)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := fhir.Instant(time.UnixMicro(timestamp).UTC())
			input.Precision = tc.precision

			got := fhirconv.ToString(input)

			if got, want := got, tc.want; got != want {
				t.Errorf("ToString[Instant]: got %v, want %v", got, want)
			}
		})
	}
}

func TestToString_DateTime_ReturnsValue(t *testing.T) {
	testCases := []struct {
		name      string
		want      string
		precision dtpb.DateTime_Precision
	}{
		{"UnspecifiedPrecision", "1970-01-01T00:00:32.000000+00:00", dtpb.DateTime_PRECISION_UNSPECIFIED},
		{"MicrosecondPrecision", "1970-01-01T00:00:32.000000+00:00", dtpb.DateTime_MICROSECOND},
		{"MillisecondPrecision", "1970-01-01T00:00:32.000+00:00", dtpb.DateTime_MILLISECOND},
		{"SecondPrecision", "1970-01-01T00:00:32+00:00", dtpb.DateTime_SECOND},
		{"DayPrecision", "1970-01-01", dtpb.DateTime_DAY},
		{"MonthPrecision", "1970-01", dtpb.DateTime_MONTH},
		{"YearPrecision", "1970", dtpb.DateTime_YEAR},
	}
	timestamp := int64(time.Second * 32 / time.Microsecond)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := fhir.DateTime(time.UnixMicro(timestamp).UTC())
			input.Precision = tc.precision

			got := fhirconv.ToString(input)

			if got, want := got, tc.want; got != want {
				t.Errorf("ToString[DateTime]: got %v, want %v", got, want)
			}
		})
	}
}

func TestToString_Date_ReturnsValue(t *testing.T) {
	testCases := []struct {
		name      string
		want      string
		precision dtpb.Date_Precision
	}{
		{"DayPrecision", "1970-01-01", dtpb.Date_DAY},
		{"MonthPrecision", "1970-01", dtpb.Date_MONTH},
		{"YearPrecision", "1970", dtpb.Date_YEAR},
	}
	timestamp := int64(time.Second * 32 / time.Microsecond)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := fhir.Date(time.UnixMicro(timestamp).UTC())
			input.Precision = tc.precision

			got := fhirconv.ToString(input)

			if got, want := got, tc.want; got != want {
				t.Errorf("ToString[Date]: got %v, want %v", got, want)
			}
		})
	}
}

func TestToString_Time_ReturnsValue(t *testing.T) {
	testCases := []struct {
		name                              string
		hour, minute, second, microsecond int64
		want                              string
	}{
		{"MicrosecondPrecision", 12, 32, 10, 123456, "12:32:10.123456"},
		{"MillisecondPrecision", 12, 32, 10, 123_000, "12:32:10.123"},
		{"SecondPrecision", 12, 32, 10, 0, "12:32:10"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input, err := fhir.TimeOfDay(tc.hour, tc.minute, tc.second, tc.microsecond)
			if err != nil {
				t.Fatalf("ToString[Time]: error setting up time-of-day: %v", err)
			}

			got := fhirconv.ToString(input)

			if got, want := got, tc.want; got != want {
				t.Errorf("ToString[Date]: got %v, want %v", got, want)
			}
		})
	}
}

func TestDateToString_RoundTrip_ReturnsInput(t *testing.T) {
	testCases := []struct {
		name string
		want string
	}{
		{"YearPrecision", "2019"},
		{"MonthPrecision", "2019-01"},
		{"DayPrecision", "2019-01-02"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value := fhir.MustParseDate(tc.want)

			got := fhirconv.ToString(value)

			if got != tc.want {
				t.Errorf("RoundTrip(%v): got %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

func TestTimeToString_RoundTrip_ReturnsInput(t *testing.T) {
	testCases := []struct {
		name string
		want string
	}{
		{"MicrosecondPrecision", "01:02:03.123456"},
		{"MillisecondPrecision", "01:02:03.123"},
		{"SecondPrecision", "01:02:03"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value := fhir.MustParseTime(tc.want)

			got := fhirconv.ToString(value)

			if got != tc.want {
				t.Errorf("RoundTrip(%v): got %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

func TestDateTimeToString_RoundTrip_ReturnsInput(t *testing.T) {
	testCases := []struct {
		name string
		want string
	}{
		{"YearPrecision", "2019"},
		{"MonthPrecision", "2019-01"},
		{"DayPrecision", "2019-01-02"},
		{"SecondPrecision", "2019-01-02T01:02:03-04:00"},
		{"MillisecondPrecision", "2019-01-02T01:02:03.123-04:00"},
		{"MicrosecondPrecision", "2019-01-02T01:02:03.123456-04:00"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value := fhir.MustParseDateTime(tc.want)

			got := fhirconv.ToString(value)

			if got != tc.want {
				t.Errorf("RoundTrip(%v): got %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}
func TestInstantToString_RoundTrip_ReturnsInput(t *testing.T) {
	testCases := []struct {
		name string
		want string
	}{
		{"MicrosecondPrecision", "2019-01-02T01:02:03.123456-04:00"},
		{"MillisecondPrecision", "2019-01-02T01:02:03.123-04:00"},
		{"SecondPrecision", "2019-01-02T01:02:03-04:00"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value := fhir.MustParseInstant(tc.want)

			got := fhirconv.ToString(value)

			if got != tc.want {
				t.Errorf("RoundTrip(%v): got %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}
