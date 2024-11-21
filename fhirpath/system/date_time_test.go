package system_test

import (
	"errors"
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestParseDateTime_ReturnsTime(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"Full DateTime with offset", "2009-03-06T12:09:45.556-04:30"},
		{"DateTime without offset", "2009-03-06T12:09:45.556"},
		{"Second with offset", "2006-01-02T15:04:05Z"},
		{"Minute", "2010-02-04T14:05"},
		{"Year", "2021T"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := system.ParseDateTime(tc.input); err != nil {
				t.Fatalf("ParseTime(%s) returned unexpected error: %v", tc.input, err)
			}
		})
	}
}

func TestParseDateTime_ReturnsError(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"Bad format", "2010-12-21T08-05"},
		{"Bad hour", "2010-12-25T25"},
		{"Bad minute", "2010-12-25T23:61"},
		{"Offset without minutes", "2010-12-25T23:59+04"},
		{"Date without T", "2010-12-25"},
		{"Time without date", "08:30Z"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := system.ParseDateTime(tc.input); err == nil {
				t.Fatalf("ParseDate(%s) didn't return error when expected", tc.input)
			}
		})
	}
}

func TestDateTime_Equal(t *testing.T) {
	testCases := []struct {
		name        string
		dateTimeOne system.DateTime
		dateTimeTwo system.Any
		shouldEqual bool
		wantOk      bool
	}{
		{
			name:        "same time different format",
			dateTimeOne: system.MustParseDateTime("2023-01-01T"),
			dateTimeTwo: system.MustParseDateTime("2023-01-01T00:00:00.000"),
			wantOk:      false,
		},
		{
			name:        "different date",
			dateTimeOne: system.MustParseDateTime("2023-01-01T00:00:00.000"),
			dateTimeTwo: system.MustParseDateTime("2023-01-02T00:00:00.000Z"),
			shouldEqual: false,
			wantOk:      true,
		},
		{
			name:        "different type",
			dateTimeOne: system.MustParseDateTime("2023-01-01T00:00:00.000"),
			dateTimeTwo: system.String("2023-01-01T00:00:00.000"),
			shouldEqual: false,
			wantOk:      true,
		},
		{
			name:        "same time different offset format",
			dateTimeOne: system.MustParseDateTime("2023-01-02T00:00:00.000Z"),
			dateTimeTwo: system.MustParseDateTime("2023-01-02T00:00:00.000-00:00"),
			shouldEqual: true,
			wantOk:      true,
		},
		{
			name:        "same time different time zone",
			dateTimeOne: system.MustParseDateTime("2023-01-01T08:30:00+03:00"),
			dateTimeTwo: system.MustParseDateTime("2023-01-01T05:30:00Z"),
			shouldEqual: true,
			wantOk:      true,
		},
		{
			name:        "not equal with mismatched precision",
			dateTimeOne: system.MustParseDateTime("2023-02-01T08:30"),
			dateTimeTwo: system.MustParseDateTime("2023-01-01T08:30:45"),
			shouldEqual: false,
			wantOk:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := tc.dateTimeOne.TryEqual(tc.dateTimeTwo)

			if ok != tc.wantOk {
				t.Fatalf("DateTime.Equal: ok got %v, want %v", ok, tc.wantOk)
			}
			if got != tc.shouldEqual {
				t.Errorf("DateTime.Equal returned unexpected equality: got %v, want %v", got, tc.shouldEqual)
			}
		})
	}
}

func TestDateTimeFromProto_Converts(t *testing.T) {
	almostY2K, _ := system.ParseDateTime("1999-12-31T23:59:59.999Z")
	dec2004, _ := system.ParseDateTime("2004-12T")
	pandemic, _ := system.ParseDateTime("2020-03-15T08:30:05Z")

	testCases := []struct {
		name    string
		dtProto *dtpb.DateTime
		wantDT  system.DateTime
	}{
		{
			name:    "converts microsecond precision",
			dtProto: fhir.MustParseDateTime("1999-12-31T23:59:59.999999Z"),
			wantDT:  almostY2K,
		},
		{
			name:    "converts partial date precision",
			dtProto: fhir.MustParseDateTime("2004-12"),
			wantDT:  dec2004,
		},
		{
			name:    "converts second precision",
			dtProto: fhir.MustParseDateTime("2020-03-15T08:30:05Z"),
			wantDT:  pandemic,
		},
		{
			name:    "converts with alternate offset format",
			dtProto: fhir.MustParseDateTime("2020-03-15T08:30:05-00:00"),
			wantDT:  pandemic,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := system.DateTimeFromProto(tc.dtProto)

			if err != nil {
				t.Fatalf("DateTimeFromProto(%s) raised unexpected error: %v", tc.dtProto.String(), err)
			}
			if diff := cmp.Diff(tc.wantDT, got); diff != "" {
				t.Errorf("DateTimeFromProto(%s) incorrectly converts DateTime: (-want, +got)\n%s", tc.dtProto.String(), diff)
			}
		})
	}
}

func TestDateTimeToProtoDateTime_Converts(t *testing.T) {
	testCases := []struct {
		name     string
		dateTime system.DateTime
		want     *dtpb.DateTime
	}{
		{
			name:     "converts millisecond precision with timezone",
			dateTime: system.MustParseDateTime("1999-12-31T23:59:59.999Z"),
			want:     fhir.MustParseDateTime("1999-12-31T23:59:59.999Z"),
		},
		{
			name:     "converts to second precision with timezone",
			dateTime: system.MustParseDateTime("1999-12-31T23:59:59Z"),
			want:     fhir.MustParseDateTime("1999-12-31T23:59:59Z"),
		},
		{
			name:     "converts to day precision",
			dateTime: system.MustParseDateTime("1999-12-31T"),
			want:     fhir.MustParseDateTime("1999-12-31"),
		},
		{
			name:     "converts to month precision",
			dateTime: system.MustParseDateTime("1999-12T"),
			want:     fhir.MustParseDateTime("1999-12"),
		},
		{
			name:     "converts to year precision",
			dateTime: system.MustParseDateTime("1999T"),
			want:     fhir.MustParseDateTime("1999"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.dateTime.ToProtoDateTime()

			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("DateTime.ToProtoDateTime returned unexpected result: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestDateTimeLess_ReturnsBoolean(t *testing.T) {
	testCases := []struct {
		name        string
		dateTimeOne system.DateTime
		dateTimeTwo system.DateTime
		want        system.Boolean
	}{
		{
			name:        "returns true for an earlier date",
			dateTimeOne: system.MustParseDateTime("2012-11-15T08:12:12"),
			dateTimeTwo: system.MustParseDateTime("2012-11-15T08:30:12"),
			want:        true,
		},
		{
			name:        "returns false for a later time",
			dateTimeOne: system.MustParseDateTime("2013-11-15T00:30:00.000Z"),
			dateTimeTwo: system.MustParseDateTime("2012-11-15T00:30:00.000Z"),
			want:        false,
		},
		{
			name:        "returns true for an earlier date with mismatched precision",
			dateTimeOne: system.MustParseDateTime("2022-11-11T18"),
			dateTimeTwo: system.MustParseDateTime("2022-11-12T18:30"),
			want:        true,
		},
		{
			name:        "returns false for a later time with mismatched precision",
			dateTimeOne: system.MustParseDateTime("2021-11-12T18:30"),
			dateTimeTwo: system.MustParseDateTime("2021-11-04T14:30:25"),
			want:        false,
		},
		{
			name:        "returns false (doesn't error) for equal times with mismatched millisecond precision",
			dateTimeOne: system.MustParseDateTime("2021-11-04T14:30:25.000"),
			dateTimeTwo: system.MustParseDateTime("2021-11-04T14:30:25"),
		},
		{
			name:        "returns true for earlier time with mismatched millisecond precision",
			dateTimeOne: system.MustParseDateTime("2000-01-01T18:30:01"),
			dateTimeTwo: system.MustParseDateTime("2000-01-01T18:30:01.001"),
			want:        true,
		},
		{
			name:        "respects time zone offset",
			dateTimeOne: system.MustParseDateTime("2000-12-19T18:30:01+05:00"),
			dateTimeTwo: system.MustParseDateTime("2000-12-19T13:30:02"),
			want:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.dateTimeOne.Less(tc.dateTimeTwo)

			if err != nil {
				t.Fatalf("DateTime.Less returned unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("DateTime.Less returned unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDateTimeLess_ReturnsError(t *testing.T) {
	testCases := []struct {
		name        string
		dateTimeOne system.DateTime
		input       system.Any
		wantErr     error
	}{
		{
			name:        "incorrect input type",
			dateTimeOne: system.MustParseDateTime("2020-09-02T08:30:30"),
			input:       system.String("2020-09-02T08:30:30"),
			wantErr:     system.ErrTypeMismatch,
		},
		{
			name:        "equal until precision mismatch",
			dateTimeOne: system.MustParseDateTime("2020-09-02T08:30:30"),
			input:       system.MustParseDateTime("2020-09-02T08:30"),
			wantErr:     system.ErrMismatchedPrecision,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.dateTimeOne.Less(tc.input)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("DateTime.Less returned incorrect error: got %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestDateTimeAdd_ReturnsSum(t *testing.T) {
	testCases := []struct {
		name        string
		dateTimeOne system.DateTime
		input       system.Quantity
		want        system.DateTime
		wantErr     error
	}{
		{
			name:        "overflows hours",
			dateTimeOne: system.MustParseDateTime("2012-11-15T08:12:12"),
			input:       system.MustParseQuantity("24", "hours"),
			want:        system.MustParseDateTime("2012-11-16T08:12:12"),
		},
		{
			name:        "adds years correctly when start date is a leap year",
			dateTimeOne: system.MustParseDateTime("2020-02-29T08:12:12"),
			input:       system.MustParseQuantity("1", "year"),
			want:        system.MustParseDateTime("2021-02-28T08:12:12"),
		},
		{
			name:        "adds decimal of seconds",
			dateTimeOne: system.MustParseDateTime("2013-11-15T00:30:00.000Z"),
			input:       system.MustParseQuantity("1.232", "seconds"),
			want:        system.MustParseDateTime("2013-11-15T00:30:01.232Z"),
		},
		{
			name:        "rounds quantity down to highest precision",
			dateTimeOne: system.MustParseDateTime("2022-11-11T18"),
			input:       system.MustParseQuantity("59", "minutes"),
			want:        system.MustParseDateTime("2022-11-11T18"),
		},
		{
			name:        "respects adding months when the result month is of a different length",
			dateTimeOne: system.MustParseDateTime("2021-10-31T18:30"),
			input:       system.MustParseQuantity("1", "month"),
			want:        system.MustParseDateTime("2021-11-30T18:30"),
		},
		{
			name:        "disregards decimal part of quantity",
			dateTimeOne: system.MustParseDateTime("2021-11-04T14:30:25.000"),
			input:       system.MustParseQuantity("3.23", "months"),
			want:        system.MustParseDateTime("2022-02-04T14:30:25.000"),
		},
		{
			name:        "adds days correctly",
			dateTimeOne: system.MustParseDateTime("2021-11-04T14:30:25.000"),
			input:       system.MustParseQuantity("32", "days"),
			want:        system.MustParseDateTime("2021-12-06T14:30:25.000"),
		},
		{
			name:        "adds weeks correctly",
			dateTimeOne: system.MustParseDateTime("2021-11-01T14:30:25.000"),
			input:       system.MustParseQuantity("18", "weeks"),
			want:        system.MustParseDateTime("2022-03-07T14:30:25.000"),
		},
		{
			name:        "returns error when adding non time-valued quantity",
			dateTimeOne: system.MustParseDateTime("2021-11-01T14:30:25.000"),
			input:       system.MustParseQuantity("18", "lbs"),
			wantErr:     system.ErrMismatchedUnit,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.dateTimeOne.Add(tc.input)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("DateTime.Add returned unexpected error: got: %v, want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("DateTime.Add returned unexpected diff: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestDateTimeSub_ReturnsResults(t *testing.T) {
	testCases := []struct {
		name        string
		dateTimeOne system.DateTime
		input       system.Quantity
		want        system.DateTime
		wantErr     error
	}{
		{
			name:        "overflows hours",
			dateTimeOne: system.MustParseDateTime("2012-11-15T08:12:12"),
			input:       system.MustParseQuantity("24", "hours"),
			want:        system.MustParseDateTime("2012-11-14T08:12:12"),
		},
		{
			name:        "subtracts years correctly when start date is a leap year",
			dateTimeOne: system.MustParseDateTime("2020-02-29T08:12:12"),
			input:       system.MustParseQuantity("1", "year"),
			want:        system.MustParseDateTime("2019-02-28T08:12:12"),
		},
		{
			name:        "correctly subtracts from a month partial",
			dateTimeOne: system.MustParseDateTime("2020-02T"),
			input:       system.MustParseQuantity("12", "years"),
			want:        system.MustParseDateTime("2008-02T"),
		},
		{
			name:        "correctly subtracts weeks",
			dateTimeOne: system.MustParseDateTime("2020-02-03T08"),
			input:       system.MustParseQuantity("12", "weeks"),
			want:        system.MustParseDateTime("2019-11-11T08"),
		},
		{
			name:        "subtracts decimal of seconds",
			dateTimeOne: system.MustParseDateTime("2013-11-15T00:30:01.232Z"),
			input:       system.MustParseQuantity("1.232", "seconds"),
			want:        system.MustParseDateTime("2013-11-15T00:30:00.000Z"),
		},
		{
			name:        "subtracts years correctly if in a leap year",
			dateTimeOne: system.MustParseDateTime("2020-02-29T"),
			input:       system.MustParseQuantity("1", "year"),
			want:        system.MustParseDateTime("2019-02-28T"),
		},
		{
			name:        "rounds quantity down to highest precision for hour partial",
			dateTimeOne: system.MustParseDateTime("2022-11-11T18"),
			input:       system.MustParseQuantity("59", "minutes"),
			want:        system.MustParseDateTime("2022-11-11T18"),
		},
		{
			name:        "rounds quantity down to to highest precision for year partial",
			dateTimeOne: system.MustParseDateTime("2021T"),
			input:       system.MustParseQuantity("23", "months"),
			want:        system.MustParseDateTime("2020T"),
		},
		{
			name:        "respects subtracting months when the result month is of a different length",
			dateTimeOne: system.MustParseDateTime("2021-10-31T18:30"),
			input:       system.MustParseQuantity("1", "month"),
			want:        system.MustParseDateTime("2021-09-30T18:30"),
		},
		{
			name:        "disregards decimal part of quantity",
			dateTimeOne: system.MustParseDateTime("2021-11-04T14:30:25.000"),
			input:       system.MustParseQuantity("3.23", "months"),
			want:        system.MustParseDateTime("2021-08-04T14:30:25.000"),
		},
		{
			name:        "returns error when subtracting non time-valued quantity",
			dateTimeOne: system.MustParseDateTime("2021-11-01T14:30:25.000"),
			input:       system.MustParseQuantity("18", "lbs"),
			wantErr:     system.ErrMismatchedUnit,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.dateTimeOne.Sub(tc.input)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("DateTime.Sub returned unexpected error: got: %v, want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("DateTime.Sub returned unexpected diff: (-want, +got)\n%s", diff)
			}
		})
	}
}
