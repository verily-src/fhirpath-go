package fhir

import (
	"errors"
	"fmt"
	"math"
	"regexp"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/uuid"
	"github.com/verily-src/fhirpath-go/internal/slices"
)

var (
	// ErrIntegerDataLoss is an error raised in APIs that might unintentionally
	// truncate integral values that the user wouldn't expect.
	ErrIntegerDataLoss = errors.New("data-loss occurred during integer conversion")
)

// Primitive Types:
//
// The section below defines types from the "Primitive Types" heading in
// http://hl7.org/fhir/R4/datatypes.html#open

// Base64Binary creates an R4 FHIR Base64Binary element the specified bytes.
//
// See: https://hl7.org/fhir/R4/datatypes.html#base64Binary
func Base64Binary(value []byte) *dtpb.Base64Binary {
	return &dtpb.Base64Binary{
		Value: value,
	}
}

// Boolean creates a Boolean proto from a primitive value.
//
// See: https://hl7.org/fhir/R4/datatypes.html#boolean
func Boolean(value bool) *dtpb.Boolean {
	return &dtpb.Boolean{
		Value: value,
	}
}

// Canonical defined in canonical.go

// Code creates an R4 FHIR Code element from a string value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#code
func Code(value string) *dtpb.Code {
	return &dtpb.Code{
		Value: value,
	}
}

// Date and DateTime in time.go

// Decimal creates an R4 FHIR Decimal element from a float64 (double) value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#decimal
func Decimal(value float64) *dtpb.Decimal {
	return &dtpb.Decimal{
		Value: fmt.Sprint(value),
	}
}

// ID creates an R4 FHIR ID element from a string value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#id
func ID(value string) *dtpb.Id {
	return &dtpb.Id{
		Value: value,
	}
}

var idRegexp *regexp.Regexp = regexp.MustCompile(`^[A-Za-z0-9\-\.]{1,64}$`)

// IsID returns true if the given string a valid FHIR ID.
// See http://hl7.org/fhir/R4/datatypes.html#id.
func IsID(id string) bool {
	return idRegexp.MatchString(id)
}

// RandomID generates a new random R4 FHIR ID element.
func RandomID() *dtpb.Id {
	return ID(uuid.NewString())
}

// Instant in time.go

// Integer creates an R4 FHIR Integer element from an int32 value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#integer
func Integer(value int32) *dtpb.Integer {
	return &dtpb.Integer{
		Value: value,
	}
}

// IntegerFromInt creates an R4 FHIR Integer element from an int value.
//
// Go's int type is architecture dependent, values greater than int32 will be
// truncated, as described in the go tour: https://tour.golang.org/basics/11
// If this occurs, this function returns an error.
func IntegerFromInt(value int) (*dtpb.Integer, error) {
	if val, ok := tryNarrowInt32(value); ok {
		return &dtpb.Integer{
			Value: val,
		}, nil
	}
	return nil, fmt.Errorf("integer(%v): %w", value, ErrIntegerDataLoss)
}

// tryNarrowInt32 is an implementation function used in the `IntegerFromInt`
// that narrows an int to an in32 and tests if there was any truncation.
func tryNarrowInt32(v int) (int32, bool) {
	v32 := int32(v)
	if int(v32) == v {
		return v32, true
	}
	return 0, false
}

// IntegerFromPositiveInt attempts to create an R4 FHIR Integer element from a
// PositiveInt value. This function may fail of the value stored in the PositiveInt
// exceeds the cardinality of Integer, which may cause a signed integer overflow.
//
// For more information, see the diagram for Primitive Types here:
// https://www.hl7.org/fhir/datatypes.html
func IntegerFromPositiveInt(value *dtpb.PositiveInt) (*dtpb.Integer, error) {
	return integerFrom(value)
}

// IntegerFromUnsignedInt attempts to create an R4 FHIR Integer element from an
// UnsignedInt value. This function may fail of the value stored in the UnsignedInt
// exceeds the cardinality of Integer, which may cause a signed integer overflow.
//
// For more information, see the diagram for Primitive Types here:
// https://www.hl7.org/fhir/datatypes.html
func IntegerFromUnsignedInt(value *dtpb.UnsignedInt) (*dtpb.Integer, error) {
	return integerFrom(value)
}

// integerFrom is an implementation function used in the `IntegerFrom*`
// functions to minimize code duplication.
func integerFrom(value interface{ GetValue() uint32 }) (*dtpb.Integer, error) {
	if isLossyConversionToInt32(value.GetValue()) {
		return nil, fmt.Errorf("integer(%v): %w", value, ErrIntegerDataLoss)
	}
	return &dtpb.Integer{
		Value: int32(value.GetValue()),
	}, nil
}

// isLossyConversionToInt32 checks if a uint32 value can be converted to an int32
// value without losing its original value, which can happen if the sign-bit is
// set in the uint32.
func isLossyConversionToInt32(value uint32) bool {
	// This works by testing that casting a uint32 to an int32 and checking that
	// the value is still positive. If it changes to negative, the cast was never
	// valid.
	return value > math.MaxInt32
}

// Markdown creates an R4 FHIR Markdown element from a string value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#markdown
func Markdown(value string) *dtpb.Markdown {
	return &dtpb.Markdown{
		Value: value,
	}
}

// OID creates an R4 FHIR OID element from a OID-string value, prepending the
// necessary "urn:oid:" to the value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#oid
func OID(value string) *dtpb.Oid {
	return &dtpb.Oid{
		Value: fmt.Sprintf("urn:oid:%v", value),
	}
}

// PositiveInt creates an R4 FHIR PositiveInt element from a uint32 value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#positiveInt
func PositiveInt(value uint32) *dtpb.PositiveInt {
	return &dtpb.PositiveInt{
		Value: value,
	}
}

// String creates an R4 FHIR String element from a string value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#string
func String(value string) *dtpb.String {
	return &dtpb.String{
		Value: value,
	}
}

// Strings creates an array of R4 FHIR String elements from a string value. This
// is offered as a convenience function, since many FHIR protos have arrays of
// FHIR string types, and converting between Go strings and FHIR strings is a
// common and repetitive process for some types.
func Strings(values ...string) []*dtpb.String {
	return slices.Map(values, String)
}

// StringFromCode is a convenience utility for converting a Code to its
// base-class definition of String. If the input is nil, this returns nil.
//
// For more information, see the diagram for Primitive Types here:
// https://www.hl7.org/fhir/datatypes.html
func StringFromCode(code *dtpb.Code) *dtpb.String {
	return stringFrom(code, code == nil)
}

// StringFromMarkdown is a convenience utility for converting Markdown to its
// base-class definition of String. If the input is nil, this returns nil.
//
// For more information, see the diagram for Primitive Types here:
// https://www.hl7.org/fhir/datatypes.html
func StringFromMarkdown(markdown *dtpb.Markdown) *dtpb.String {
	return stringFrom(markdown, markdown == nil)
}

// StringFromID is a convenience utility for converting an Id to its
// base-class definition of String. If the input is nil, this returns nil.
//
// For more information, see the diagram for Primitive Types here:
// https://www.hl7.org/fhir/datatypes.html
func StringFromID(id *dtpb.Id) *dtpb.String {
	return stringFrom(id, id == nil)
}

// stringFrom is an implementation of the `StringFrom*` series of functions to
// cut down on repetition.
//
// This function takes 'isNil' as a boolean argument to work around the fact that
// a nil pointer passed to an interface forms a non-nil interface in Go, and
// reflection is more costly than a bool check.
func stringFrom(value interface{ GetValue() string }, isNil bool) *dtpb.String {
	if isNil {
		return nil
	}
	return &dtpb.String{
		Value: value.GetValue(),
	}
}

// Time in time.go

// UnsignedInt creates an R4 FHIR UnsignedInt element from a uint32 value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#unsignedInt
func UnsignedInt(value uint32) *dtpb.UnsignedInt {
	return &dtpb.UnsignedInt{
		Value: value,
	}
}

// URI creates an R4 FHIR URI element from a string value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#uri
func URI(value string) *dtpb.Uri {
	return &dtpb.Uri{
		Value: value,
	}
}

// URIFromCanonical is a convenience utility for converting a canonical to its
// base-class definition of URI. If the input is nil, this returns nil.
//
// For more information, see the diagram for Primitive Types here:
// https://www.hl7.org/fhir/datatypes.html
func URIFromCanonical(canonical *dtpb.Canonical) *dtpb.Uri {
	return uriFrom(canonical, canonical == nil)
}

// URIFromOID is a convenience utility for converting an OID to its
// base-class definition of URI. If the input is nil, this returns nil.
//
// For more information, see the diagram for Primitive Types here:
// https://www.hl7.org/fhir/datatypes.html
func URIFromOID(oid *dtpb.Oid) *dtpb.Uri {
	return uriFrom(oid, oid == nil)
}

// URIFromURL is a convenience utility for converting a URL to its
// base-class definition of URI. If the input is nil, this returns nil.
//
// For more information, see the diagram for Primitive Types here:
// https://www.hl7.org/fhir/datatypes.html
func URIFromURL(url *dtpb.Url) *dtpb.Uri {
	return uriFrom(url, url == nil)
}

// URIFromUUID is a convenience utility for converting a UUID to its
// base-class definition of URI. If the input is nil, this returns nil.
//
// For more information, see the diagram for Primitive Types here:
// https://www.hl7.org/fhir/datatypes.html
func URIFromUUID(uuid *dtpb.Uuid) *dtpb.Uri {
	return uriFrom(uuid, uuid == nil)
}

// uriFrom is a convenience helper for implementing the various `UriFrom*`
// functions which are all the same repetative logic.
//
// This function takes 'isNil' as a boolean argument to work around the fact that
// a nil pointer passed to an interface forms a non-nil interface in Go, and
// reflection is more costly than a bool check.
func uriFrom(other interface{ GetValue() string }, isNil bool) *dtpb.Uri {
	if isNil {
		return nil
	}
	return URI(other.GetValue())
}

// URL creates an R4 FHIR URL element from a string value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#url
func URL(value string) *dtpb.Url {
	return &dtpb.Url{
		Value: value,
	}
}

// UUID creates an R4 FHIR UUID element from a uuid-string value, prepending the
// necessary "urn:uuid:" to the value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#uuid
func UUID(value string) *dtpb.Uuid {
	return &dtpb.Uuid{
		Value: fmt.Sprintf("urn:uuid:%v", value),
	}
}

// RandomUUID generates a random new UUID.
//
// See: http://hl7.org/fhir/R4/datatypes.html#uuid
func RandomUUID() *dtpb.Uuid {
	return UUID(uuid.NewString())
}
