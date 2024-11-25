package system_test

import (
	"errors"
	"os"
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

func TestParseTime_ReturnsTime(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"Hour", "14"},
		{"Minute", "14:23"},
		{"Second", "14:23:21"},
		{"Millisecond", "14:23:21.999"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := system.ParseTime(tc.input); err != nil {
				t.Fatalf("ParseTime(%s) returned unexpected error: %v", tc.input, err)
			}
		})
	}
}

func TestParseTime_ReturnsError(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"Bad format", "14-23-21.556"},
		{"Bad hour", "25:23:21.999"},
		{"Bad minute", "15:65:21.999"},
		{"Bad second", "15:23:65.999"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := system.ParseTime(tc.input); err == nil {
				t.Fatalf("ParseDate(%s) didn't raise error when expected", tc.input)
			}
		})
	}
}

func TestTime_Equal(t *testing.T) {
	testCases := []struct {
		name        string
		timeOne     system.Time
		timeTwo     system.Any
		shouldEqual bool
		wantOk      bool
	}{
		{
			name:    "same time different format",
			timeOne: system.MustParseTime("08"),
			timeTwo: system.MustParseTime("08:00:00"),
			wantOk:  false,
		},
		{
			name:        "different time with mismatched precision",
			timeOne:     system.MustParseTime("08:24:30"),
			timeTwo:     system.MustParseTime("08:30"),
			shouldEqual: false,
			wantOk:      true,
		},
		{
			name:        "treats millisecond and second as same precision",
			timeOne:     system.MustParseTime("08:30:10.000"),
			timeTwo:     system.MustParseTime("08:30:10"),
			shouldEqual: true,
			wantOk:      true,
		},
		{
			name:        "different time",
			timeOne:     system.MustParseTime("08:00:00"),
			timeTwo:     system.MustParseTime("08:30:00"),
			shouldEqual: false,
			wantOk:      true,
		},
		{
			name:    "different type",
			timeOne: system.MustParseTime("08:00:00"),
			timeTwo: system.String("08:00:00"),
			wantOk:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := tc.timeOne.TryEqual(tc.timeTwo)

			if ok != tc.wantOk {
				t.Fatalf("Time.Equal: ok got %v, want %v", ok, tc.wantOk)
			}
			if got != tc.shouldEqual {
				t.Errorf("Time.Equal returned unexpected equality: got %v, want %v", got, tc.shouldEqual)
			}
		})
	}
}

func TestTimeFromProto_Converts(t *testing.T) {
	testCases := []struct {
		name      string
		timeProto *dtpb.Time
		wantTime  system.Time
	}{
		{
			name:      "converts microsecond precision",
			timeProto: fhir.MustParseTime("08:30:00.212123"),
			wantTime:  system.MustParseTime("08:30:00.212"),
		},
		{
			name:      "converts second precision",
			timeProto: fhir.MustParseTime("16:45:00"),
			wantTime:  system.MustParseTime("16:45:00"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			changeLocale(t)
			got := system.TimeFromProto(tc.timeProto)

			if diff := cmp.Diff(tc.wantTime, got); diff != "" {
				t.Errorf("TimeFromProto(%s) incorrectly converts Time: (-want, +got)\n%s", tc.timeProto.String(), diff)
			}
		})
	}
}

func TestTimeLess_ReturnsBoolean(t *testing.T) {
	testCases := []struct {
		name    string
		timeOne system.Time
		timeTwo system.Time
		want    system.Boolean
	}{
		{
			name:    "returns true for an earlier time",
			timeOne: system.MustParseTime("08:30:05"),
			timeTwo: system.MustParseTime("12:03:05"),
			want:    true,
		},
		{
			name:    "returns false for a later time",
			timeOne: system.MustParseTime("18:30:05"),
			timeTwo: system.MustParseTime("12:03:05"),
			want:    false,
		},
		{
			name:    "returns true for an earlier time with mismatched precision",
			timeOne: system.MustParseTime("18:30"),
			timeTwo: system.MustParseTime("18:45:05"),
			want:    true,
		},
		{
			name:    "returns false for a later time with mismatched precision",
			timeOne: system.MustParseTime("18:30:34"),
			timeTwo: system.MustParseTime("12:20"),
			want:    false,
		},
		{
			name:    "returns false (doesn't error) for equal times with mismatched millisecond precision",
			timeOne: system.MustParseTime("04:30:25"),
			timeTwo: system.MustParseTime("04:30:25.000"),
			want:    false,
		},
		{
			name:    "returns true for earlier time with mismatched millisecond precision",
			timeOne: system.MustParseTime("18:30:01"),
			timeTwo: system.MustParseTime("18:30:01.001"),
			want:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.timeOne.Less(tc.timeTwo)

			if err != nil {
				t.Fatalf("Time.Less returned unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("Time.Less returned unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestTimeLess_ReturnsError(t *testing.T) {
	testCases := []struct {
		name    string
		timeOne system.Time
		input   system.Any
		wantErr error
	}{
		{
			name:    "incorrect input type",
			timeOne: system.MustParseTime("08:30:30"),
			input:   system.String("08:30:30"),
			wantErr: system.ErrTypeMismatch,
		},
		{
			name:    "equal until precision mismatch",
			timeOne: system.MustParseTime("11:11:11"),
			input:   system.MustParseTime("11:11"),
			wantErr: system.ErrMismatchedPrecision,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.timeOne.Less(tc.input)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Time.Less returned incorrect error: got %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestTimeAdd_ReturnsSum(t *testing.T) {
	testCases := []struct {
		name    string
		time    system.Time
		input   system.Quantity
		want    system.Time
		wantErr error
	}{
		{
			name:  "overflows time hours",
			time:  system.MustParseTime("23:30:00"),
			input: system.MustParseQuantity("2", "hours"),
			want:  system.MustParseTime("01:30:00"),
		},
		{
			name:  "correctly adds minutes",
			time:  system.MustParseTime("19:45:44"),
			input: system.MustParseQuantity("25", "minutes"),
			want:  system.MustParseTime("20:10:44"),
		},
		{
			name:  "correctly adds seconds",
			time:  system.MustParseTime("08:00:00"),
			input: system.MustParseQuantity("23", "seconds"),
			want:  system.MustParseTime("08:00:23"),
		},
		{
			name:  "disregards decimal part of minutes",
			time:  system.MustParseTime("08:00:00"),
			input: system.MustParseQuantity("15.5", "minutes"),
			want:  system.MustParseTime("08:15:00"),
		},
		{
			name:  "adds decimal part of seconds",
			time:  system.MustParseTime("08:00:00.000"),
			input: system.MustParseQuantity("12.24", "seconds"),
			want:  system.MustParseTime("08:00:12.240"),
		},
		{
			name:  "adds partials by rounding down to the highest precision",
			time:  system.MustParseTime("08"),
			input: system.MustParseQuantity("59", "minutes"),
			want:  system.MustParseTime("08"),
		},
		{
			name:    "returns error on non time-valued quantity",
			time:    system.MustParseTime("08:30:00"),
			input:   system.MustParseQuantity("12", "kg"),
			wantErr: system.ErrMismatchedUnit,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.time.Add(tc.input)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Time.Add returned unexpected error: got %v, want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Time.Add returned unexpected result: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestTimeSub_ReturnsResult(t *testing.T) {
	testCases := []struct {
		name    string
		time    system.Time
		input   system.Quantity
		want    system.Time
		wantErr error
	}{
		{
			name:  "overflows time hours",
			time:  system.MustParseTime("00:30:00"),
			input: system.MustParseQuantity("2", "hours"),
			want:  system.MustParseTime("22:30:00"),
		},
		{
			name:  "correctly subtracts minutes",
			time:  system.MustParseTime("19:45:44"),
			input: system.MustParseQuantity("25", "minutes"),
			want:  system.MustParseTime("19:20:44"),
		},
		{
			name:  "correctly subtracts seconds",
			time:  system.MustParseTime("08:00:00"),
			input: system.MustParseQuantity("23", "seconds"),
			want:  system.MustParseTime("07:59:37"),
		},
		{
			name:  "disregards decimal part of minutes",
			time:  system.MustParseTime("07:45:00"),
			input: system.MustParseQuantity("15.5", "minutes"),
			want:  system.MustParseTime("07:30:00"),
		},
		{
			name:  "subtracts decimal part of seconds",
			time:  system.MustParseTime("08:00:12.240"),
			input: system.MustParseQuantity("12.24", "seconds"),
			want:  system.MustParseTime("08:00:00.000"),
		},
		{
			name:  "subtracts partials by rounding down to highest precision",
			time:  system.MustParseTime("08:00"),
			input: system.MustParseQuantity("119", "seconds"),
			want:  system.MustParseTime("07:59"),
		},
		{
			name:    "returns error on non time-valued quantity",
			time:    system.MustParseTime("08:30:00"),
			input:   system.MustParseQuantity("12", "kg"),
			wantErr: system.ErrMismatchedUnit,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.time.Sub(tc.input)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Time.Add returned unexpected error: got %v, want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Time.Add returned unexpected result: (-want, +got)\n%s", diff)
			}
		})
	}
}

func changeLocale(t *testing.T) {
	t.Helper()

	// find a new locale
	tz := os.Getenv("TZ")
	newLocale := "Asia/Tokyo"
	if tz == newLocale {
		newLocale = "Africa/Cairo"
	}
	if err := os.Setenv("TZ", newLocale); err != nil {
		t.Fatalf("error setting locale: %v", err)
	}

	// revert locale back to original
	t.Cleanup(func() {
		if err := os.Setenv("TZ", tz); err != nil {
			t.Fatalf("error setting locale: %v", err)
		}
	})
}
