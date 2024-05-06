package system

import (
	"fmt"
	"strings"
	"time"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/shopspring/decimal"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirconv"
)

// DateTime represents date/time values and partial date/time values
// in the format YYYY-MM-DDThh:mm:ss.fff(+|-)hh:mm format.
type DateTime struct {
	dateTime time.Time
	l        layout
}

// ParseDateTime parses a string and returns a DateTime object if
// the string is a valid FHIRPath DateTime string, or an error otherwise.
func ParseDateTime(value string) (DateTime, error) {
	dateTimeLayouts := []string{
		dtMillisecondLayoutTZ,
		dtMillisecondLayout,
		dtSecondLayoutTZ,
		dtSecondLayout,
		dtMinuteLayoutTZ,
		dtMinuteLayout,
		dtHourLayoutTZ,
		dtHourLayout,
		dtDayLayout,
		dtMonthLayout,
		dtYearLayout,
	}

	var t time.Time
	var err error
	value = strings.TrimPrefix(value, "@")
	for _, l := range dateTimeLayouts {
		if t, err = time.Parse(l, value); err == nil {
			return DateTime{t, layout(l)}, nil
		}
	}
	return DateTime{}, fmt.Errorf("unable to parse DateTime '%s': %w", value, err)
}

// MustParseDateTime returns a DateTime if the string represents a
// valid DateTime. Otherwise, panics.
func MustParseDateTime(value string) DateTime {
	dt, err := ParseDateTime(value)
	if err != nil {
		panic(err)
	}
	return dt
}

// DateTimeFromProto takes a proto DateTime as input and returns a System DateTime.
// Note that the highest precision supported by System.DateTime is Millisecond.
func DateTimeFromProto(proto *dtpb.DateTime) (DateTime, error) {
	t, err := fhirconv.DateTimeToTime(proto)
	if err != nil {
		return DateTime{}, err
	}
	var l layout
	switch proto.Precision {
	case dtpb.DateTime_MICROSECOND:
		fallthrough
	case dtpb.DateTime_MILLISECOND:
		l = dtMillisecondLayoutTZ
	case dtpb.DateTime_SECOND:
		l = dtSecondLayoutTZ
	case dtpb.DateTime_DAY:
		l = dtDayLayout
	case dtpb.DateTime_MONTH:
		l = dtMonthLayout
	case dtpb.DateTime_YEAR:
		l = dtYearLayout
	}
	return DateTime{t, l}, nil
}

// ToProtoDateTime returns a proto DateTime based on a system DateTime.
// Note that the highest precision supported by System.DateTime is Millisecond.
func (dt DateTime) ToProtoDateTime() *dtpb.DateTime {
	dateTime := fhir.DateTime(dt.dateTime)
	var p dtpb.DateTime_Precision
	switch dt.l {
	case dtMillisecondLayoutTZ, dtMillisecondLayout:
		p = dtpb.DateTime_MILLISECOND
	case dtSecondLayoutTZ, dtSecondLayout:
		p = dtpb.DateTime_SECOND
	case dtDayLayout:
		p = dtpb.DateTime_DAY
	case dtMonthLayout:
		p = dtpb.DateTime_MONTH
	case dtYearLayout:
		p = dtpb.DateTime_YEAR
	}
	dateTime.Precision = p
	return dateTime
}

// TryEqual returns a boolean representing whether or not
// the value of dt is equal to the value of dt2.
// Not intended to be used for cmp.Equal. The comparison is
// not symmetric and may cause unexpected behaviour.
func (dt DateTime) TryEqual(input Any) (bool, bool) {
	val, ok := input.(DateTime)
	if !ok {
		return false, true
	}
	if dt.l == val.l {
		return dt.dateTime.Equal(val.dateTime), true
	}

	// normalize time zone
	dt.dateTime = dt.dateTime.UTC()
	val.dateTime = val.dateTime.UTC()

	dtComponents := dt.getComponents()
	valComponents := val.getComponents()

	minPrecision := min(int(dateTimeMap[dt.l]), int(dateTimeMap[val.l]))

	for i := 0; i <= minPrecision; i++ {
		if dtComponents[i] == valComponents[i] && i != int(dtSecond) {
			continue
		}
		return dtComponents[i] == valComponents[i], true
	}
	return false, false
}

// Less returns true if the value of dt is less than input.(DateTime).
// Compares component by component, and returns an error if there is a
// precision mismatch. If input is not a Date, returns an error.
func (dt DateTime) Less(input Any) (Boolean, error) {
	val, ok := input.(DateTime)
	if !ok {
		return false, fmt.Errorf("%w, %T, %T", ErrTypeMismatch, dt, input)
	}
	if dt.l == val.l {
		return Boolean(dt.dateTime.Before(val.dateTime)), nil
	}

	// normalize time zone
	dt.dateTime = dt.dateTime.UTC()
	val.dateTime = val.dateTime.UTC()

	dtComponents := dt.getComponents()
	valComponents := val.getComponents()

	minPrecision := min(int(dateTimeMap[dt.l]), int(dateTimeMap[val.l]))

	for i := 0; i <= minPrecision; i++ {
		// precisions below second are irrelevant, and should be treated the same.
		if dtComponents[i] == valComponents[i] && i != int(dtSecond) {
			continue
		}
		return dtComponents[i] < valComponents[i], nil
	}
	return false, ErrMismatchedPrecision
}

// Add returns the result of dt + input. Returns an
// error if input does not represent a valid time valued quantity.
func (dt DateTime) Add(input Quantity) (DateTime, error) {
	var result time.Time
	value := int(decimal.Decimal(input.value).IntPart())
	switch input.unit {
	case "year", "years":
		result = addYear(dt.dateTime, value)
	case "month", "months":
		result = addMonth(dt.dateTime, value)
	case "week", "weeks":
		value = 7 * value
		result = dt.dateTime.AddDate(0, 0, value)
	case "day", "days":
		result = dt.dateTime.AddDate(0, 0, value)
	default:
		duration, err := input.timeDuration()
		if err != nil {
			return DateTime{}, err
		}
		result = dt.dateTime.Add(duration)
	}

	// Reformat to truncate DateTime to initial precision, rounding down to
	// highest precision value.
	result, err := time.Parse(string(dt.l), result.Format(string(dt.l)))
	if err != nil {
		return DateTime{}, err
	}
	return DateTime{result, dt.l}, nil
}

// Sub returns the result of dt - input.(Quantity). Returns an
// error if the input is not a Quantity, or if it does not represent
// a valid time duration.
func (dt DateTime) Sub(input Quantity) (DateTime, error) {
	// Handle partial dates by rounding quantity to appropriate precision.
	// Subtraction is not symmetric with addition, so truncation cannot be
	// applied here.
	if dt.l == dtYearLayout {
		years, err := input.toYears()
		if err != nil {
			return DateTime{}, err
		}
		return DateTime{dt.dateTime.AddDate(-years, 0, 0), dt.l}, nil
	}
	if dt.l == dtMonthLayout {
		months, err := input.toMonths()
		if err != nil {
			return DateTime{}, err
		}
		return DateTime{dt.dateTime.AddDate(0, -months, 0), dt.l}, nil
	}

	// Handles non-partial dates here.
	var result time.Time
	value := -int(decimal.Decimal(input.value).IntPart())
	switch input.unit {
	case "year", "years":
		result = addYear(dt.dateTime, value)
	case "month", "months":
		result = addMonth(dt.dateTime, value)
	case "week", "weeks":
		value = 7 * value
		result = dt.dateTime.AddDate(0, 0, value)
	case "day", "days":
		result = dt.dateTime.AddDate(0, 0, value)
	default:
		// Get time valued duration, and round down to appropriate precision.
		duration, err := input.timeDuration()
		if err != nil {
			return DateTime{}, err
		}
		duration = roundToDateTimePrecision(dateTimeMap[dt.l], duration)
		result = dt.dateTime.Add(-duration)
	}
	return DateTime{result, dt.l}, nil
}

// Name returns the type name.
func (dt DateTime) Name() string {
	return dateTimeType
}

// String returns a formatted DateTime string.
func (dt DateTime) String() string {
	return dt.dateTime.Format(string(dt.l))
}

// Equal method to override cmp.Equal.
func (dt DateTime) Equal(dt2 DateTime) bool {
	return dt.dateTime.Format(string(dt.l)) == dt2.dateTime.Format(string(dt2.l))
}

func (dt DateTime) getComponents() []int {
	return []int{
		dt.dateTime.Year(),
		int(dt.dateTime.Month()),
		dt.dateTime.Day(),
		dt.dateTime.Hour(),
		dt.dateTime.Minute(),
		dt.dateTime.Second()*1000000000 + dt.dateTime.Nanosecond(),
	}
}

// roundToDateTimePrecision rounds the duration down to the appropriate precision.
// Eg. 2012-03-20T + 23 'hours' = 2012-03-20T but 2012-03-20T + 24 'hours' = 2012-03-21T.
func roundToDateTimePrecision(p dateTimePrecision, d time.Duration) time.Duration {
	switch p {
	case dtYear:
		return d / (time.Hour * 24 * 365)
	case dtMonth:
		return d / (time.Hour * 24 * 30)
	case dtDay:
		return d / (time.Hour * 24)
	case dtHour:
		return d / time.Hour
	case dtMinute:
		return d / time.Minute
	default:
		return d
	}
}
