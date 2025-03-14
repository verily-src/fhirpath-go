package fhirconv_test

import (
	"testing"
	"time"

	"github.com/google/fhir/go/fhirversion"
	"github.com/google/fhir/go/jsonformat"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirconv"
	"google.golang.org/protobuf/testing/protocmp"
)

func newMarshaller(t *testing.T) *jsonformat.Marshaller {
	t.Helper()

	marshaller, err := jsonformat.NewMarshaller(false, "", "", fhirversion.R4)
	if err != nil {
		t.Fatalf("Error creating marshaller: %v", err)
	}
	return marshaller
}

func newUnmarshaller(t *testing.T, zone string) *jsonformat.Unmarshaller {
	t.Helper()

	unmarshaller, err := jsonformat.NewUnmarshaller(zone, fhirversion.R4)
	if err != nil {
		t.Fatalf("Error creating unmarshaller: %v", err)
	}
	return unmarshaller
}

func dateTimeExtension(dt *dtpb.DateTime) *dtpb.Extension {
	return &dtpb.Extension{
		Value: &dtpb.Extension_ValueX{
			Choice: &dtpb.Extension_ValueX_DateTime{
				DateTime: dt,
			},
		},
	}
}

func loadLocation(t *testing.T, name string) *time.Location {
	t.Helper()

	tz, err := time.LoadLocation(name)
	if err != nil {
		t.Fatalf("Unable to load location: %v", err)
	}
	return tz
}

func TestDateTimeToTime_InvalidInput_ReturnsError(t *testing.T) {
	testCases := []struct {
		name string
		zone string
	}{
		{"InvalidTimeZoneCode", "XYZBLAH"},
		{"LargePositiveMinuteOffset", "10:-90"},
		{"LargeNegativeMinuteOffset", "10:90"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dt := fhir.DateTimeNow()
			dt.Timezone = tc.zone

			_, err := fhirconv.DateTimeToTime(dt)

			if err == nil {
				t.Errorf("DateTimeToTime(%v): got nil, want err", tc.name)
			}
		})
	}
}

func TestDateTimeToTime_ValidInput_ReturnsTime(t *testing.T) {
	testCases := []struct {
		name     string
		zone     string
		location *time.Location
	}{
		{"EmptyString", "", time.UTC},
		{"UTC", "UTC", time.UTC},
		{"EST", "EST", loadLocation(t, "EST")},
		{"Z", "Z", time.UTC},
		{"SmallPositiveOffset", "+00:01", time.FixedZone("", 60)},
		{"SmallNegativeOffset", "-00:02", time.FixedZone("", -120)},
		{"LargePositiveOffset", "+12:03", time.FixedZone("", 43380)},
		{"LargeNegativeOffset", "-13:04", time.FixedZone("", -47040)},
	}
	const value = 1000

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dt := fhir.DateTime(time.UnixMicro(value))
			dt.Timezone = tc.zone
			want := time.UnixMicro(value).In(tc.location)

			got, err := fhirconv.DateTimeToTime(dt)
			if err != nil {
				t.Fatalf("DateTimeToTime(%v): got err '%v', want nil", tc.name, err)
			}

			if !cmp.Equal(got, want) {
				t.Errorf("DateTimeToTime(%v): got '%v', want '%v'", tc.name, got, want)
			}
		})
	}
}

func TestDateToTime_InvalidInput_ReturnsError(t *testing.T) {
	testCases := []struct {
		name string
		zone string
	}{
		{"InvalidTimeZoneCode", "XYZBLAH"},
		{"LargePositiveMinuteOffset", "10:-90"},
		{"LargeNegativeMinuteOffset", "10:90"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dt := fhir.DateNow()
			dt.Timezone = tc.zone

			_, err := fhirconv.DateToTime(dt)

			if err == nil {
				t.Errorf("DateToTime(%v): got nil, want err", tc.name)
			}
		})
	}
}

func TestTimeToTime_ValidInput_ReturnsTime(t *testing.T) {
	testCases := []struct {
		name     string
		zone     string
		location *time.Location
	}{
		{"EmptyString", "", time.UTC},
		{"UTC", "UTC", time.UTC},
		{"EST", "EST", loadLocation(t, "EST")},
		{"Z", "Z", time.UTC},
		{"Local", "Local", time.Local},
		{"SmallPositiveOffset", "+00:01", time.FixedZone("", 60)},
		{"SmallNegativeOffset", "-00:02", time.FixedZone("", -120)},
		{"LargePositiveOffset", "+12:03", time.FixedZone("", 43380)},
		{"LargeNegativeOffset", "-13:04", time.FixedZone("", -47040)},
	}
	const value = 1000

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dt := fhir.Date(time.UnixMicro(value))
			dt.Timezone = tc.zone
			want := time.UnixMicro(value).In(tc.location)

			got, err := fhirconv.DateToTime(dt)
			if err != nil {
				t.Fatalf("DateToTime(%v): got err '%v', want nil", tc.name, err)
			}

			if !cmp.Equal(got, want) {
				t.Errorf("DateToTime(%v): got '%v', want '%v'", tc.name, got, want)
			}
		})
	}
}

func TestInstantToTime_InvalidInput_ReturnsError(t *testing.T) {
	testCases := []struct {
		name string
		zone string
	}{
		{"InvalidTimeZoneCode", "XYZBLAH"},
		{"LargePositiveMinuteOffset", "10:-90"},
		{"LargeNegativeMinuteOffset", "10:90"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dt := fhir.InstantNow()
			dt.Timezone = tc.zone

			_, err := fhirconv.InstantToTime(dt)

			if err == nil {
				t.Errorf("InstantToTime(%v): got nil, want err", tc.name)
			}
		})
	}
}

func TestInstantToTime_ValidInput_ReturnsTime(t *testing.T) {
	testCases := []struct {
		name     string
		zone     string
		location *time.Location
	}{
		{"EmptyString", "", time.UTC},
		{"UTC", "UTC", time.UTC},
		{"EST", "EST", loadLocation(t, "EST")},
		{"Z", "Z", time.UTC},
		{"SmallPositiveOffset", "+00:01", time.FixedZone("", 60)},
		{"SmallNegativeOffset", "-00:02", time.FixedZone("", -120)},
		{"LargePositiveOffset", "+12:03", time.FixedZone("", 43380)},
		{"LargeNegativeOffset", "-13:04", time.FixedZone("", -47040)},
	}
	const value = 1000

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dt := fhir.Instant(time.UnixMicro(value))
			dt.Timezone = tc.zone
			want := time.UnixMicro(value).In(tc.location)

			got, err := fhirconv.InstantToTime(dt)
			if err != nil {
				t.Fatalf("InstantToTime(%v): got err '%v', want nil", tc.name, err)
			}

			if !cmp.Equal(got, want) {
				t.Errorf("InstantToTime(%v): got '%v', want '%v'", tc.name, got, want)
			}
		})
	}
}

func TestDurationToDuration_BadInput_ReturnsError(t *testing.T) {
	// Creates a base "good" Duration that we can modify for test purposes without
	// coupling to the units package explicitly.
	makeDuration := func(setup func(*dtpb.Duration)) *dtpb.Duration {
		base := fhir.Seconds(time.Second * 10)
		setup(base)
		return base
	}

	testCases := []struct {
		name     string
		duration *dtpb.Duration
	}{
		{"BadFloatValue", makeDuration(func(d *dtpb.Duration) { d.Value.Value = "4e99999" })},
		{"BadUnitSymbol", makeDuration(func(d *dtpb.Duration) { d.Code = fhir.Code("blarg") })},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := fhirconv.DurationToDuration(tc.duration)

			if err == nil {
				t.Errorf("DurationToDuration(%v): got nil, want err", tc.name)
			}
		})
	}
}

func TestDurationToDuration_RoundTrip(t *testing.T) {
	testCases := []struct {
		name     string
		duration *dtpb.Duration
	}{
		{"Nanoseconds", fhir.Nanoseconds(42)},
		{"Microseconds", fhir.Microseconds(time.Microsecond * 3)},
		{"Milliseconds", fhir.Milliseconds(time.Millisecond * 11)},
		{"Seconds", fhir.Seconds(time.Second * 314)},
		{"Minutes", fhir.Minutes(time.Minute * 8)},
		{"Hours", fhir.Hours(time.Hour * 12)},
		{"Days", fhir.Days(time.Hour * 24)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			duration, err := fhirconv.DurationToDuration(tc.duration)
			if err != nil {
				t.Fatalf("DurationToDuration(%v): unexpected error %v", tc.name, err)
			}

			element := fhir.Duration(duration)

			got, want := element, tc.duration
			if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
				t.Errorf("DurationToDuration(%v): (-got +want):\n%v", tc.name, diff)
			}
		})
	}
}
