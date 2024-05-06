package fhir_test

import (
	"testing"
	"time"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/units"
	"google.golang.org/protobuf/testing/protocmp"
)

func newDuration(value int64, unit units.Time) *dtpb.Duration {
	return &dtpb.Duration{
		Value:  fhir.Decimal(float64(value)),
		Code:   fhir.Code(unit.Symbol()),
		System: fhir.URI(unit.System()),
	}
}

func TestDurationFromTime(t *testing.T) {
	testCases := []struct {
		name       string
		time       string
		value      int64
		multiplier int64
		unit       units.Time
	}{
		{"Seconds", "00:01:00", 1, int64(time.Minute) / int64(time.Second), units.Seconds},
		{"Milliseconds", "00:01:00.000", 1, int64(time.Minute) / int64(time.Millisecond), units.Milliseconds},
		{"Microseconds", "00:01:00.000000", 1, int64(time.Minute) / int64(time.Microsecond), units.Microseconds},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := fhir.MustParseTime(tc.time)
			value := tc.value * tc.multiplier
			want := newDuration(value, tc.unit)

			got := fhir.DurationFromTime(input)

			if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
				t.Errorf("DurationFromTime(%v): (-got +want):\n%v", tc.name, diff)
			}
		})
	}
}

func TestDuration(t *testing.T) {
	testCases := []struct {
		name       string
		value      int64
		multiplier int64
		unit       units.Time
	}{
		{"Zero", 0, 1, units.Days},
		{"Nanoseconds", 805, 1, units.Nanoseconds},
		{"Microseconds", 15, int64(time.Microsecond), units.Microseconds},
		{"Milliseconds", 32, int64(time.Millisecond), units.Milliseconds},
		{"Seconds", 42, int64(time.Second), units.Seconds},
		{"Minutes", 1, int64(time.Minute), units.Minutes},
		{"Hours", 12, int64(time.Hour), units.Hours},
		{"Days", 3, 24 * int64(time.Hour), units.Days},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value := tc.value * tc.multiplier
			timeDuration := time.Duration(value)
			want := newDuration(tc.value, tc.unit)

			got := fhir.Duration(timeDuration)

			if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
				t.Errorf("Duration(%v): (-got +want):\n%v", tc.name, diff)
			}
		})
	}
}

func TestNanoseconds(t *testing.T) {
	value := time.Duration(400)
	want := newDuration(int64(value), units.Nanoseconds)

	got := fhir.Nanoseconds(value)

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("Nanoseconds: (-got +want):\n%v", diff)
	}
}

func TestMicroseconds(t *testing.T) {
	value := time.Duration(400)
	duration := time.Microsecond * value
	want := newDuration(int64(value), units.Microseconds)

	got := fhir.Microseconds(duration)

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("Microseconds: (-got +want):\n%v", diff)
	}
}

func TestMilliseconds(t *testing.T) {
	value := 400
	duration := time.Millisecond * 400
	want := newDuration(int64(value), units.Milliseconds)

	got := fhir.Milliseconds(duration)

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("Milliseconds: (-got +want):\n%v", diff)
	}
}

func TestSeconds(t *testing.T) {
	value := 400
	duration := time.Second * 400
	want := newDuration(int64(value), units.Seconds)

	got := fhir.Seconds(duration)

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("Seconds: (-got +want):\n%v", diff)
	}
}

func TestMinutes(t *testing.T) {
	value := 400
	duration := time.Minute * 400
	want := newDuration(int64(value), units.Minutes)

	got := fhir.Minutes(duration)

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("Minutes: (-got +want):\n%v", diff)
	}
}

func TestHours(t *testing.T) {
	value := 400
	duration := time.Hour * 400
	want := newDuration(int64(value), units.Hours)

	got := fhir.Hours(duration)

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("Hours: (-got +want):\n%v", diff)
	}
}

func TestDays(t *testing.T) {
	value := 400
	duration := time.Hour * 24 * 400
	want := newDuration(int64(value), units.Days)

	got := fhir.Days(duration)

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("Days: (-got +want):\n%v", diff)
	}
}
