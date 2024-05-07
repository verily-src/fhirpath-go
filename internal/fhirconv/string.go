package fhirconv

import (
	b64 "encoding/base64"
	"fmt"
	"time"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

// ToString converts a basic FHIR Primitive into a human-readable string
// representation.
func ToString[T fhir.PrimitiveType](val T) string {
	return toString(val)
}

// toString is the implementation of the ToString function.
//
// This is a separate function written in terms of an interface, rather than
// being generic, to avoid unnecessary code-bloat -- since generics will generate
// code for every instantiation with a different type.
func toString(val fhir.Element) string {
	switch val := val.(type) {
	case interface{ GetValue() string }:
		// URI, URL, OID, UUID, Canonical, String, ID, Markdown, Code, Decimal
		return val.GetValue()
	case interface{ GetValue() bool }:
		// Boolean
		return fmt.Sprintf("%v", val.GetValue())
	case interface{ GetValue() uint32 }:
		// PositiveInteger, UnsignedInt
		return fmt.Sprintf("%v", val.GetValue())
	case interface{ GetValue() int32 }:
		// Integer
		return fmt.Sprintf("%v", val.GetValue())
	case interface{ GetValue() []byte }:
		// Base64Binary
		return b64.StdEncoding.EncodeToString(val.GetValue())
	case *dtpb.Instant:
		return InstantToString(val)
	case *dtpb.DateTime:
		return DateTimeToString(val)
	case *dtpb.Time:
		return TimeToString(val)
	case *dtpb.Date:
		return DateToString(val)
	}
	// This can't be reached; the above switch is exhaustive for all possible
	// inputs, which is restricted by the type constraint.
	return ""
}

// InstantToString converts the FHIR Instant element into its string reprsentation
// as defined in http://hl7.org/fhir/R4/datatypes.html#instant.
//
// The level of precision in the output is equivalent to the precision defined
// in the input Instant proto.
func InstantToString(val *dtpb.Instant) string {
	if tm, err := InstantToTime(val); err == nil {
		switch val.GetPrecision() {
		case dtpb.Instant_SECOND:
			return tm.Format("2006-01-02T15:04:05-07:00")
		case dtpb.Instant_MILLISECOND:
			return tm.Format("2006-01-02T15:04:05.000-07:00")
		case dtpb.Instant_MICROSECOND:
			fallthrough
		default:
			return tm.Format("2006-01-02T15:04:05.000000-07:00")
		}
	}
	// Fall-back to a basic representation (this shouldn't happen unless timezone
	// information is garbage, which is a developer-driven issue).
	return fmt.Sprintf("Instant(%v)", val.GetValueUs())
}

// DateTimeToString converts the FHIR DateTime element into its string reprsentation
// as defined in http://hl7.org/fhir/R4/datatypes.html#datetime.
//
// The level of precision in the output is equivalent to the precision defined
// in the input DateTime proto.
func DateTimeToString(val *dtpb.DateTime) string {
	if tm, err := DateTimeToTime(val); err == nil {
		switch val.GetPrecision() {
		case dtpb.DateTime_YEAR:
			return tm.Format("2006")
		case dtpb.DateTime_MONTH:
			return tm.Format("2006-01")
		case dtpb.DateTime_DAY:
			return tm.Format("2006-01-02")
		case dtpb.DateTime_SECOND:
			return tm.Format("2006-01-02T15:04:05-07:00")
		case dtpb.DateTime_MILLISECOND:
			return tm.Format("2006-01-02T15:04:05.000-07:00")
		case dtpb.DateTime_MICROSECOND:
			fallthrough
		default:
			return tm.Format("2006-01-02T15:04:05.000000-07:00")
		}
	}

	// Fall-back to a basic representation (this shouldn't happen unless timezone
	// information is garbage, which is a developer-driven issue).
	return fmt.Sprintf("DateTime(%v)", val.GetValueUs())
}

// DateToString converts the FHIR Date element into its string reprsentation
// as defined in http://hl7.org/fhir/R4/datatypes.html#date.
//
// The level of precision in the output is equivalent to the precision defined
// in the input Date proto.
func DateToString(val *dtpb.Date) string {
	if tm, err := DateToTime(val); err == nil {
		switch val.GetPrecision() {
		case dtpb.Date_YEAR:
			return tm.Format("2006")
		case dtpb.Date_MONTH:
			return tm.Format("2006-01")
		case dtpb.Date_DAY:
			fallthrough
		default:
			return tm.Format("2006-01-02")
		}
	}

	// Fall-back to a basic representation (this shouldn't happen unless timezone
	// information is garbage, which is a developer-driven issue).
	return fmt.Sprintf("Date(%v)", val.GetValueUs())
}

// TimeToString converts the FHIR Time element into its string reprsentation
// as defined in http://hl7.org/fhir/R4/datatypes.html#time.
//
// The level of precision in the output is equivalent to the precision defined
// in the input Time proto.
func TimeToString(val *dtpb.Time) string {
	duration := TimeToDuration(val)

	hours := (duration / time.Hour) % (time.Hour * 24)
	duration %= time.Hour
	minutes := duration / time.Minute
	duration %= time.Minute
	seconds := duration / time.Second
	duration %= time.Second
	micros := duration / time.Microsecond

	switch val.GetPrecision() {
	case dtpb.Time_SECOND:
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	case dtpb.Time_MILLISECOND:
		millis := micros / 1_000
		return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, millis)
	case dtpb.Time_MICROSECOND:
		fallthrough
	default:
		return fmt.Sprintf("%02d:%02d:%02d.%06d", hours, minutes, seconds, micros)
	}
}
