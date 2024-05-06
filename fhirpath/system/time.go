package system

import (
	"fmt"
	"strings"
	"time"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirconv"
)

// Time represents a time of day in the range 00:00:00.000
// to 23:59:59.999 with a step size of 1ms. It uses the format
// hh:mm:ss.fff format to parse times.
type Time struct {
	time time.Time
	l    layout
}

// ParseTime takes an input string and returns a Time object if
// formatted correctly according to the FHIRPath spec, otherwise
// an error.
func ParseTime(value string) (Time, error) {
	timeLayouts := []string{
		millisecondLayout,
		secondLayout,
		minuteLayout,
		hourLayout,
	}

	var t time.Time
	var err error
	value = strings.TrimPrefix(value, "@T")
	for _, l := range timeLayouts {
		if t, err = time.Parse(l, value); err == nil {
			return Time{t, layout(l)}, nil
		}
	}
	return Time{}, fmt.Errorf("unable to parse time '%s': %w", value, err)
}

// MustParseTime takes an input string and returns a Time object
// if formatted correctly. Otherwise, panics.
func MustParseTime(value string) Time {
	t, err := ParseTime(value)
	if err != nil {
		panic(err)
	}
	return t
}

// TimeFromProto takes a proto Time and returns a System Time.
func TimeFromProto(proto *dtpb.Time) Time {
	duration := fhirconv.TimeToDuration(proto)
	t := time.UnixMicro(duration.Microseconds()).In(time.UTC)
	var l layout
	switch proto.Precision {
	case dtpb.Time_MICROSECOND:
		fallthrough
	case dtpb.Time_MILLISECOND:
		l = millisecondLayout
	case dtpb.Time_SECOND:
		l = secondLayout
	}
	return Time{t, l}
}

// ToProtoTime returns a proto Time based on a system Time.
func (t Time) ToProtoTime() *dtpb.Time {
	tp := fhir.Time(t.time)
	switch t.l {
	case millisecondLayout:
		tp.Precision = dtpb.Time_MILLISECOND
	default:
		tp.Precision = dtpb.Time_SECOND
	}
	return tp
}

// TryEqual returns a boolean representing whether or not
// the value of t is equal to the value of t2.
// Not intended to be used for cmp.Equal. The comparison is
// not symmetric and may cause unexpected behaviour.
func (t Time) TryEqual(input Any) (bool, bool) {
	val, ok := input.(Time)
	if !ok {
		return false, true
	}
	if t.l == val.l {
		return t.time.Equal(val.time), true
	}

	tComponents := t.getComponents()
	valComponents := val.getComponents()

	minPrecision := min(int(timeMap[t.l]), int(timeMap[val.l]))

	for i := 0; i <= minPrecision; i++ {
		if tComponents[i] == valComponents[i] && i != int(second) {
			continue
		}
		return tComponents[i] == valComponents[i], true
	}
	return false, false
}

// Name returns the type name.
func (t Time) Name() string {
	return timeType
}

// Equal method to override cmp.Equal.
func (t Time) Equal(t2 Time) bool {
	return t.time.Format(string(t.l)) == t2.time.Format(string(t2.l))
}

// String formats the time as a time string.
func (t Time) String() string {
	return t.time.Format(string(t.l))
}

// Less returns true if the value of t is less than input.(Time).
// Compares component by component, and returns an error if there is a
// precision mismatch. If input is not a Time, returns an error.
func (t Time) Less(input Any) (Boolean, error) {
	val, ok := input.(Time)
	if !ok {
		return false, fmt.Errorf("%w: %T, %T", ErrTypeMismatch, t, input)
	}
	if t.l == val.l {
		return Boolean(t.time.Before(val.time)), nil
	}

	tComponents := t.getComponents()
	valComponents := val.getComponents()

	minPrecision := min(int(timeMap[t.l]), int(timeMap[val.l]))

	for i := 0; i <= minPrecision; i++ {
		// precisions below second are irrelevant, and should be treated the same.
		if tComponents[i] == valComponents[i] && i != int(second) {
			continue
		}
		return tComponents[i] < valComponents[i], nil
	}
	return false, ErrMismatchedPrecision
}

// Add returns the result of t with the time-valued quantity added to it.
// Returns an error if the Quantity does not represent a valid duration.
func (t Time) Add(input Quantity) (Time, error) {
	duration, err := input.timeDuration()
	if err != nil {
		return Time{}, err
	}
	duration = roundToTimePrecision(timeMap[t.l], duration)
	return Time{t.time.Add(duration), t.l}, nil
}

// Sub returns the result of the time-valued quantity subtracted from t.
// Returns an error if the Quantity does not represent a valid duration.
func (t Time) Sub(input Quantity) (Time, error) {
	duration, err := input.timeDuration()
	if err != nil {
		return Time{}, err
	}
	duration = roundToTimePrecision(timeMap[t.l], duration)
	return Time{t.time.Add(-duration), t.l}, nil
}

// roundToTimePrecision is used to round down to the highest precision of
// the time value.
// Eg. 08:30 + 59 'seconds' = 08:30, but 08:30 + 60 'seconds' = 08:31
func roundToTimePrecision(p timePrecision, d time.Duration) time.Duration {
	switch p {
	case hour:
		return d / time.Hour
	case minute:
		return d / time.Minute
	default:
		return d
	}
}

func (t Time) getComponents() []int {
	return []int{
		t.time.Hour(),
		t.time.Minute(),
		t.time.Second()*1000000000 + t.time.Nanosecond(),
	}
}
