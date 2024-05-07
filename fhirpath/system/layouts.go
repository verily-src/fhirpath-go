package system

// datePrecision enumerates date precision constants.
type datePrecision int

// timePrecision enumerates time precision constants.
type timePrecision int

// dateTimePrecision enumerates dateTime precision constants.
type dateTimePrecision int

// layout represents a layout string for parsing.
type layout string

// Date precision constants.
const (
	year datePrecision = iota
	month
	day
)

// Time precision constants.
const (
	hour timePrecision = iota
	minute
	second
)

// DateTime precision constants.
const (
	dtYear dateTimePrecision = iota
	dtMonth
	dtDay
	dtHour
	dtMinute
	dtSecond
)

// layout constants.
const (
	yearLayout            = "2006"
	monthLayout           = "2006-01"
	dayLayout             = "2006-01-02"
	hourLayout            = "15"
	minuteLayout          = "15:04"
	secondLayout          = "15:04:05"
	millisecondLayout     = "15:04:05.000"
	dtMillisecondLayoutTZ = "2006-01-02T15:04:05.000Z07:00"
	dtMillisecondLayout   = "2006-01-02T15:04:05.000"
	dtSecondLayoutTZ      = "2006-01-02T15:04:05Z07:00"
	dtSecondLayout        = "2006-01-02T15:04:05"
	dtMinuteLayoutTZ      = "2006-01-02T15:04Z07:00"
	dtMinuteLayout        = "2006-01-02T15:04"
	dtHourLayoutTZ        = "2006-01-02T15Z07:00"
	dtHourLayout          = "2006-01-02T15"
	dtDayLayout           = "2006-01-02T"
	dtMonthLayout         = "2006-01T"
	dtYearLayout          = "2006T"
)

var dateMap = map[layout]datePrecision{
	dayLayout:   day,
	monthLayout: month,
	yearLayout:  year,
}

var timeMap = map[layout]timePrecision{
	millisecondLayout: second,
	secondLayout:      second,
	minuteLayout:      minute,
	hourLayout:        hour,
}

var dateTimeMap = map[layout]dateTimePrecision{
	dtMillisecondLayoutTZ: dtSecond,
	dtMillisecondLayout:   dtSecond,
	dtSecondLayoutTZ:      dtSecond,
	dtSecondLayout:        dtSecond,
	dtMinuteLayoutTZ:      dtMinute,
	dtMinuteLayout:        dtMinute,
	dtHourLayoutTZ:        dtHour,
	dtHourLayout:          dtHour,
	dtDayLayout:           dtDay,
	dtMonthLayout:         dtMonth,
	dtYearLayout:          dtYear,
}
