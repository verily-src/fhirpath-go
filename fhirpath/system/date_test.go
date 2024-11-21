package system_test

import (
	"errors"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"google.golang.org/protobuf/testing/protocmp"
	"testing"
)

func TestParseDate_ReturnsDate(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"Year", "2012"},
		{"Month", "2012-05"},
		{"Day", "2012-05-21"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := system.ParseDate(tc.input); err != nil {
				t.Fatalf("ParseDate(%s) returned unexpected error: %v", tc.input, err)
			}
		})
	}
}

func TestParseDate_ReturnsError(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"Bad format", "05-01-2007"},
		{"Bad month", "2023-13-21"},
		{"Bad day", "2023-12-34"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := system.ParseDate(tc.input); err == nil {
				t.Fatalf("ParseDate(%s) didn't return error when expected", tc.input)
			}
		})
	}
}

func TestDate_Equal(t *testing.T) {
	testCases := []struct {
		name        string
		dateOne     system.Date
		dateTwo     system.Any
		shouldEqual bool
		wantOk      bool
	}{
		{
			name:    "same date different precision",
			dateOne: system.MustParseDate("2023"),
			dateTwo: system.MustParseDate("2023-01-01"),
			wantOk:  false,
		},
		{
			name:        "different date",
			dateOne:     system.MustParseDate("2023-01-01"),
			dateTwo:     system.MustParseDate("2023-01-02"),
			shouldEqual: false,
			wantOk:      true,
		},
		{
			name:        "different type",
			dateOne:     system.MustParseDate("2023-01-01"),
			dateTwo:     system.String("2023-01-01"),
			shouldEqual: false,
			wantOk:      true,
		},
		{
			name:        "mismatched precision but not equal",
			dateOne:     system.MustParseDate("2023-02"),
			dateTwo:     system.MustParseDate("2023-01-29"),
			shouldEqual: false,
			wantOk:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := tc.dateOne.TryEqual(tc.dateTwo)

			if ok != tc.wantOk {
				t.Errorf("Date.Equal: ok got %v, want %v", ok, tc.wantOk)
			}
			if got != tc.shouldEqual {
				t.Errorf("Date.Equal returned unexpected equality: got %v, want %v", got, tc.shouldEqual)
			}
		})
	}
}

func TestDateFromProto_Converts(t *testing.T) {
	testCases := []struct {
		name      string
		dateProto *dtpb.Date
		wantDate  system.Date
	}{
		{
			name:      "converts day precision",
			dateProto: fhir.MustParseDate("2022-05-23"),
			wantDate:  system.MustParseDate("2022-05-23"),
		},
		{
			name:      "converts month precision",
			dateProto: fhir.MustParseDate("2012-12"),
			wantDate:  system.MustParseDate("2012-12"),
		},
		{
			name:      "converts year precision",
			dateProto: fhir.MustParseDate("2002"),
			wantDate:  system.MustParseDate("2002"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := system.DateFromProto(tc.dateProto)
			if err != nil {
				t.Fatalf("DateFromProto(%s) raises unexpected error: %v", tc.dateProto.String(), err)
			}
			if diff := cmp.Diff(tc.wantDate, got); diff != "" {
				t.Errorf("DateFromProto(%s) incorrectly converts Date: (-want, +got)\n%s", tc.dateProto.String(), diff)
			}
		})
	}
}

func TestDateToProtoDate_Converts(t *testing.T) {
	testCases := []struct {
		name string
		date system.Date
		want *dtpb.Date
	}{
		{
			name: "day precision",
			date: system.MustParseDate("2025-01-30"),
			want: fhir.MustParseDate("2025-01-30"),
		},
		{
			name: "month precision",
			date: system.MustParseDate("2025-01"),
			want: fhir.MustParseDate("2025-01"),
		},
		{
			name: "year precision",
			date: system.MustParseDate("2025"),
			want: fhir.MustParseDate("2025"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.date.ToProtoDate()

			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Date.ToProtoDate returned unexpected result: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestDateLess_ReturnsBoolean(t *testing.T) {
	testCases := []struct {
		name    string
		dateOne system.Date
		dateTwo system.Date
		want    system.Boolean
	}{
		{
			name:    "returns true for an earlier date",
			dateOne: system.MustParseDate("2020-05-01"),
			dateTwo: system.MustParseDate("2020-06-01"),
			want:    true,
		},
		{
			name:    "returns true for an earlier date with less precision",
			dateOne: system.MustParseDate("2020-05"),
			dateTwo: system.MustParseDate("2020-07-28"),
			want:    true,
		},
		{
			name:    "returns false for a later date",
			dateOne: system.MustParseDate("2023-05-01"),
			dateTwo: system.MustParseDate("2022-02-02"),
			want:    false,
		},
		{
			name:    "returns false for a later date with less precision",
			dateOne: system.MustParseDate("2020"),
			dateTwo: system.MustParseDate("2019-05"),
			want:    false,
		},
		{
			name:    "returns false for equivalent dates",
			dateOne: system.MustParseDate("2008-02-09"),
			dateTwo: system.MustParseDate("2008-02-09"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.dateOne.Less(tc.dateTwo)

			if err != nil {
				t.Fatalf("Date.Less returned unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("Date.Less returned unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDateLt_ReturnsError(t *testing.T) {
	testCases := []struct {
		name    string
		dateOne system.Date
		input   system.Any
		wantErr error
	}{
		{
			name:    "wrong input type",
			dateOne: system.MustParseDate("2020-12-31"),
			input:   system.String("asdf"),
			wantErr: system.ErrTypeMismatch,
		},
		{
			name:    "date precision shorter than input",
			dateOne: system.MustParseDate("2020-04"),
			input:   system.MustParseDate("2020-04-01"),
			wantErr: system.ErrMismatchedPrecision,
		},
		{
			name:    "input precision shorter than date",
			dateOne: system.MustParseDate("2023-05-31"),
			input:   system.MustParseDate("2023-05"),
			wantErr: system.ErrMismatchedPrecision,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.dateOne.Less(tc.input)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Date.Less returned incorrect error: got %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestDateAdd_ReturnsSum(t *testing.T) {
	testCases := []struct {
		name    string
		date    system.Date
		input   system.Quantity
		want    system.Date
		wantErr error
	}{
		{
			name:  "overflows months",
			date:  system.MustParseDate("2025-11-30"),
			input: system.MustParseQuantity("2", "months"),
			want:  system.MustParseDate("2026-01-30"),
		},
		{
			name:  "correctly adds days",
			date:  system.MustParseDate("2021-02-28"),
			input: system.MustParseQuantity("2", "days"),
			want:  system.MustParseDate("2021-03-02"),
		},
		{
			name:  "correctly adds years",
			date:  system.MustParseDate("2002-04-19"),
			input: system.MustParseQuantity("18", "years"),
			want:  system.MustParseDate("2020-04-19"),
		},
		{
			name:  "adds weeks as days",
			date:  system.MustParseDate("1974-10-30"),
			input: system.MustParseQuantity("18", "weeks"),
			want:  system.MustParseDate("1975-03-05"),
		},
		{
			name:  "disregards decimal part of quantity",
			date:  system.MustParseDate("2002-04-19"),
			input: system.MustParseQuantity("2.5", "months"),
			want:  system.MustParseDate("2002-06-19"),
		},
		{
			name:  "returns correct result when input month is longer than result month",
			date:  system.MustParseDate("2021-01-31"),
			input: system.MustParseQuantity("1", "month"),
			want:  system.MustParseDate("2021-02-28"),
		},
		{
			name:  "adds partials by rounding down to the highest precision",
			date:  system.MustParseDate("1997"),
			input: system.MustParseQuantity("23", "months"),
			want:  system.MustParseDate("1998"),
		},
		{
			name:    "returns error on unsupported time-valued quantity",
			date:    system.MustParseDate("2019-03-31"),
			input:   system.MustParseQuantity("12", "hours"),
			wantErr: system.ErrMismatchedUnit,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.date.Add(tc.input)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Date.Add returned unexpected error: got %v, want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Date.Add returned unexpected result: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestDateSub_ReturnsSum(t *testing.T) {
	testCases := []struct {
		name    string
		date    system.Date
		input   system.Quantity
		want    system.Date
		wantErr error
	}{
		{
			name:  "overflows months",
			date:  system.MustParseDate("2025-01-30"),
			input: system.MustParseQuantity("2", "months"),
			want:  system.MustParseDate("2024-11-30"),
		},
		{
			name:  "correctly subtracts days",
			date:  system.MustParseDate("2021-03-02"),
			input: system.MustParseQuantity("2", "days"),
			want:  system.MustParseDate("2021-02-28"),
		},
		{
			name:  "correctly subtracts years",
			date:  system.MustParseDate("2020-04-19"),
			input: system.MustParseQuantity("18", "years"),
			want:  system.MustParseDate("2002-04-19"),
		},
		{
			name:  "disregards decimal part of quantity",
			date:  system.MustParseDate("2002-04-19"),
			input: system.MustParseQuantity("2.5", "months"),
			want:  system.MustParseDate("2002-02-19"),
		},
		{
			name:  "subtracts year partial by rounding down to the highest precision",
			date:  system.MustParseDate("1997"),
			input: system.MustParseQuantity("23", "months"),
			want:  system.MustParseDate("1996"),
		},
		{
			name:  "subtracts month partial by rounding down to multiples of 30 days",
			date:  system.MustParseDate("1997-04"),
			input: system.MustParseQuantity("29", "days"),
			want:  system.MustParseDate("1997-04"),
		},
		{
			name:  "returns correct result when input month is longer than result month",
			date:  system.MustParseDate("2019-03-31"),
			input: system.MustParseQuantity("1", "month"),
			want:  system.MustParseDate("2019-02-28"),
		},
		{
			name:    "returns error on non time-valued quantity",
			date:    system.MustParseDate("2019-03-31"),
			input:   system.MustParseQuantity("12", "kg"),
			wantErr: system.ErrMismatchedUnit,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.date.Sub(tc.input)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Date.Sub returned unexpected error: got %v, want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Date.Sub returned unexpected result: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestDate_ToDateTime(t *testing.T) {
	tests := []struct {
		name  string
		input system.Date
		want  system.DateTime
	}{
		{
			name:  "returns a DateTime for a Date with day layout",
			input: system.MustParseDate("2006-01-02"),
			want:  system.MustParseDate("2006-01-02").ToDateTime(),
		},
		{
			name:  "returns a DateTime for a Date with month layout",
			input: system.MustParseDate("2006-01"),
			want:  system.MustParseDate("2006-01").ToDateTime(),
		},
		{
			name:  "returns a DateTime for a Date with year layout",
			input: system.MustParseDate("2006"),
			want:  system.MustParseDate("2006").ToDateTime(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.ToDateTime()
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Date_ToDateTime() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}
