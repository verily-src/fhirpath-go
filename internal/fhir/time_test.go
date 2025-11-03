package fhir_test

import (
	"testing"
	"time"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

func TestTimeOfDay_BadInput_ReturnsError(t *testing.T) {
	testCases := []struct {
		name         string
		hour         int64
		minute       int64
		second       int64
		microseconds int64
	}{
		{"BigHour", 24, 10, 0, 0},
		{"NegativeHour", -1, 10, 0, 0},
		{"BigMinute", 12, 60, 0, 0},
		{"NegativeMinute", 12, -1, 0, 0},
		{"BigSecond", 12, 10, 60, 0},
		{"NegativeSecond", 12, 10, -1, 0},
		{"BigMicrosecond", 12, 10, 0, 1_000_000},
		{"NegativeMicrosecond", 12, 10, 0, -1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := fhir.TimeOfDay(tc.hour, tc.minute, tc.second, tc.microseconds)

			if err == nil {
				t.Errorf("TimeOfDay(%v): want err, got nil", tc.name)
			}
		})
	}
}

func TestTimeOfDay_GoodInput_ReturnsTimeBelow24Hours(t *testing.T) {
	testCases := []struct {
		name         string
		hour         int64
		minute       int64
		second       int64
		microseconds int64
	}{
		{"NormalValue", 4, 20, 6, 9},
		{"MinHour", 0, 20, 6, 9},
		{"MaxHour", 23, 20, 6, 9},
		{"MinMinute", 4, 0, 6, 9},
		{"MaxMinute", 4, 59, 6, 9},
		{"MinSecond", 4, 20, 0, 9},
		{"MaxSecond", 4, 20, 59, 9},
		{"MinMicros", 4, 20, 6, 0},
		{"MaxMicros", 4, 20, 6, 999_999},
		{"AlmostMidnight", 23, 59, 59, 999_999},
		{"Midnight", 0, 0, 0, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			day := (time.Hour * 24).Microseconds()
			got, err := fhir.TimeOfDay(tc.hour, tc.minute, tc.second, tc.microseconds)
			if err != nil {
				t.Fatalf("TimeOfDay(%v): unexpected error: %v", tc.name, got)
			}

			if got, want := got.ValueUs, day; got >= want {
				t.Errorf("TimeOfDay(%v): got us %v, want below %v", tc.name, got, want)
			}
		})
	}
}

func TestParseTime_GoodInput_ReturnsTime(t *testing.T) {
	testCases := []struct {
		name          string
		value         string
		wantPrecision dtpb.Time_Precision
	}{
		{"Second", "01:02:03", dtpb.Time_SECOND},
		{"Milliseconds", "01:02:03.456", dtpb.Time_MILLISECOND},
		{"Microseconds", "01:02:03.456789", dtpb.Time_MICROSECOND},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := fhir.ParseTime(tc.value)
			if err != nil {
				t.Fatalf("ParseTime(%v): got err: %v", tc.name, err)
			}

			t.Run("HasExpectedPrecision", func(t *testing.T) {
				if got, want := got.Precision, tc.wantPrecision; got != want {
					t.Errorf("ParseTime(%v): got %v, want %v", tc.name, got, want)
				}
			})
		})
	}
}

func TestParseTime_BadInput_ReturnsError(t *testing.T) {
	testCases := []struct {
		name  string
		value string
	}{
		{"BadHour", "25:02:03"},
		{"BadMinute", "01:61:03"},
		{"BadSecond", "01:02:61"},
		{"WrongFormat", "01:02:61:72"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := fhir.ParseTime(tc.value)

			if err == nil {
				t.Errorf("ParseTime(%v): want err, got nil", tc.name)
			}
		})
	}
}

func TestMustParseTime_GoodInput_ReturnsTime(t *testing.T) {
	got := fhir.MustParseTime("01:02:03.456789")

	if got, want := got.Precision, dtpb.Time_MICROSECOND; got != want {
		t.Errorf("MustParseTime: got %v, want %v", got, want)
	}
}

func TestMustParseTime_BadInput_Panics(t *testing.T) {
	defer func() { _ = recover() }()

	fhir.MustParseTime("March 26, 1993")

	// If code reaches here, it means we didn't panic
	t.Errorf("MustParseTime: expected panic")
}

func TestParseDate_GoodInput_ReturnsDate(t *testing.T) {
	testCases := []struct {
		name          string
		value         string
		wantPrecision dtpb.Date_Precision
	}{
		{"Year", "2012", dtpb.Date_YEAR},
		{"Month", "2012-10", dtpb.Date_MONTH},
		{"Day", "2012-10-15", dtpb.Date_DAY},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := fhir.ParseDate(tc.value)
			if err != nil {
				t.Fatalf("ParseDate(%v): got err: %v", tc.name, err)
			}

			t.Run("HasExpectedPrecision", func(t *testing.T) {
				if got, want := got.Precision, tc.wantPrecision; got != want {
					t.Errorf("ParseDate(%v): got %v, want %v", tc.name, got, want)
				}
			})
		})
	}
}

func TestParseDate_BadInput_ReturnsError(t *testing.T) {
	testCases := []struct {
		name  string
		value string
	}{
		{"BadMonth", "2012-13"},
		{"BadDay", "2012-10-32"},
		{"WrongFormat", "2012-10-32T10:32:00"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := fhir.ParseDate(tc.value)

			if err == nil {
				t.Fatalf("ParseDate(%v): want err, got nil", tc.name)
			}
		})
	}
}

func TestMustParseDate_GoodInput_ReturnsTime(t *testing.T) {
	got := fhir.MustParseDate("2012")

	if got, want := got.Precision, dtpb.Date_YEAR; got != want {
		t.Errorf("MustParseDate: got %v, want %v", got, want)
	}
}

func TestMustParseDate_BadInput_Panics(t *testing.T) {
	defer func() { _ = recover() }()

	fhir.MustParseDate("March 26, 1993")

	// If code reaches here, it means we didn't panic
	t.Errorf("MustParseDate: expected panic")
}

func TestParseDateTime_GoodInput_ReturnsDateTime(t *testing.T) {
	testCases := []struct {
		name          string
		value         string
		wantPrecision dtpb.DateTime_Precision
	}{
		{"Year", "2012", dtpb.DateTime_YEAR},
		{"Month", "2012-10", dtpb.DateTime_MONTH},
		{"Day", "2012-10-15", dtpb.DateTime_DAY},
		{"Second", "2012-10-15T01:02:03-04:00", dtpb.DateTime_SECOND},
		{"Millisecond", "2012-10-15T01:02:03.123-04:00", dtpb.DateTime_MILLISECOND},
		{"Microsecond", "2012-10-15T01:02:03.123456-04:00", dtpb.DateTime_MICROSECOND},
		{"SecondZ", "2012-10-15T01:02:03Z", dtpb.DateTime_SECOND},
		{"MillisecondZ", "2012-10-15T01:02:03.123Z", dtpb.DateTime_MILLISECOND},
		{"MicrosecondZ", "2012-10-15T01:02:03.123456Z", dtpb.DateTime_MICROSECOND},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := fhir.ParseDateTime(tc.value)
			if err != nil {
				t.Fatalf("ParseDateTime(%v): got err: %v", tc.name, err)
			}

			t.Run("HasExpectedPrecision", func(t *testing.T) {
				if got, want := got.Precision, tc.wantPrecision; got != want {
					t.Errorf("ParseDateTime(%v): got %v, want %v", tc.name, got, want)
				}
			})
		})
	}
}

func TestParseDateTime_BadInput_ReturnsError(t *testing.T) {
	testCases := []struct {
		name  string
		value string
	}{
		{"BadMonth", "2012-13"},
		{"BadDay", "2012-10-32"},
		{"BadHour", "2012-10-15T24:02:03-04:00"},
		{"BadMinute", "2012-10-15T01:61:03-04:00"},
		{"BadSecond", "2012-10-15T01:02:61-04:00"},
		{"BadTimezone", "2012-10-15T01:02:61-61:00"},
		{"WrongFormat", "January 02, 2015"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := fhir.ParseDateTime(tc.value)

			if err == nil {
				t.Errorf("ParseDateTime(%v): want err, got nil", tc.name)
			}
		})
	}
}

func TestMustParseDateTime_GoodInput_ReturnsTime(t *testing.T) {
	got := fhir.MustParseDateTime("2012")

	if got, want := got.Precision, dtpb.DateTime_YEAR; got != want {
		t.Errorf("MustParseDate: got %v, want %v", got, want)
	}
}

func TestMustParseDateTime_BadInput_Panics(t *testing.T) {
	defer func() { _ = recover() }()

	fhir.MustParseDateTime("March 26, 1993")

	// If code reaches here, it means we didn't panic
	t.Errorf("MustParseDateTime: expected panic")
}

func TestParseInstant_GoodInput_ReturnsInstant(t *testing.T) {
	testCases := []struct {
		name          string
		value         string
		wantPrecision dtpb.Instant_Precision
	}{
		{"Second", "2019-01-02T01:02:03-04:00", dtpb.Instant_SECOND},
		{"Millisecond", "2019-01-02T01:02:03.123-04:00", dtpb.Instant_MILLISECOND},
		{"Microsecond", "2019-01-02T01:02:03.123456-04:00", dtpb.Instant_MICROSECOND},
		{"SecondZ", "2019-01-02T01:02:03Z", dtpb.Instant_SECOND},
		{"MillisecondZ", "2019-01-02T01:02:03.123Z", dtpb.Instant_MILLISECOND},
		{"MicrosecondZ", "2019-01-02T01:02:03.123456Z", dtpb.Instant_MICROSECOND},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := fhir.ParseInstant(tc.value)
			if err != nil {
				t.Fatalf("ParseInstant(%v): got err: %v", tc.name, err)
			}

			t.Run("HasExpectedPrecision", func(t *testing.T) {
				if got, want := got.Precision, tc.wantPrecision; got != want {
					t.Errorf("ParseInstant(%v): got %v, want %v", tc.name, got, want)
				}
			})
		})
	}
}

func TestParseInstant_BadInput_ReturnsError(t *testing.T) {
	testCases := []struct {
		name  string
		value string
	}{
		{"BadMonth", "2019-20-02T10:02:03-04:00"},
		{"BadDay", "2019-01-40T10:02:03-04:00"},
		{"BadHour", "2019-01-02T24:02:03-04:00"},
		{"BadMinute", "2019-01-02T01:60:03-04:00"},
		{"BadSecond", "2019-01-02T01:02:60-04:00"},
		{"BadTimeZone", "2019-01-02T01:02:60-64:00"},
		{"WrongFormat", "January 02, 2015"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := fhir.ParseInstant(tc.value)

			if err == nil {
				t.Errorf("ParseInstant(%v): want err, got nil", tc.name)
			}
		})
	}
}

func TestMustParseInstant_GoodInput_ReturnsTime(t *testing.T) {
	got := fhir.MustParseInstant("2019-10-02T01:02:03-04:00")

	if got, want := got.Precision, dtpb.Instant_SECOND; got != want {
		t.Errorf("MustParseInstant: got %v, want %v", got, want)
	}
}

func TestMustParseInstant_BadInput_Panics(t *testing.T) {
	defer func() { _ = recover() }()

	fhir.MustParseInstant("March 26, 1993")

	// If code reaches here, it means we didn't panic
	t.Errorf("MustParseInstant: expected panic")
}

func TestDate(t *testing.T) {
	testCases := []struct {
		name          string
		input         time.Time
		wantValue     int64
		wantPrecision dtpb.Date_Precision
		wantTimezone  string
	}{
		{
			name:      "UTC",
			input:     time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			wantValue: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC).UnixMicro(),
		},
		{
			name:      "PositiveOffset",
			input:     time.Date(2022, 1, 2, 0, 0, 0, 0, time.FixedZone("+0530", 5*3600+30*60)),
			wantValue: time.Date(2022, 1, 2, 0, 0, 0, 0, time.FixedZone("+0530", 5*3600+30*60)).UnixMicro(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := fhir.Date(tc.input)

			if got.ValueUs != tc.wantValue {
				t.Errorf("Date.ValueUs: got %v, want %v", got.ValueUs, tc.wantValue)
			}
		})
	}
}

func TestDateTime(t *testing.T) {
	testCases := []struct {
		name      string
		input     time.Time
		wantValue int64
	}{
		{
			name:      "UTC",
			input:     time.Date(2022, 1, 2, 3, 4, 5, 123456, time.UTC),
			wantValue: time.Date(2022, 1, 2, 3, 4, 5, 123456, time.UTC).UnixMicro(),
		},
		{
			name:      "NegativeOffset",
			input:     time.Date(2022, 1, 2, 3, 4, 5, 654321, time.FixedZone("-0800", -8*3600)),
			wantValue: time.Date(2022, 1, 2, 3, 4, 5, 654321, time.FixedZone("-0800", -8*3600)).UnixMicro(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := fhir.DateTime(tc.input)

			if got.ValueUs != tc.wantValue {
				t.Errorf("DateTime.ValueUs: got %v, want %v", got.ValueUs, tc.wantValue)
			}
		})
	}
}

func TestInstant(t *testing.T) {
	testCases := []struct {
		name      string
		input     time.Time
		wantValue int64
	}{
		{
			name:      "UTC",
			input:     time.Date(2022, 1, 2, 3, 4, 5, 789000, time.UTC),
			wantValue: time.Date(2022, 1, 2, 3, 4, 5, 789000, time.UTC).UnixMicro(),
		},
		{
			name:      "HalfHourOffset",
			input:     time.Date(2022, 1, 2, 3, 4, 5, 789000, time.FixedZone("+0330", 3*3600+30*60)),
			wantValue: time.Date(2022, 1, 2, 3, 4, 5, 789000, time.FixedZone("+0330", 3*3600+30*60)).UnixMicro(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := fhir.Instant(tc.input)

			if got.ValueUs != tc.wantValue {
				t.Errorf("Instant.ValueUs: got %v, want %v", got.ValueUs, tc.wantValue)
			}
		})
	}
}
