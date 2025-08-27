package system

import (
	"encoding/base64"
	"errors"
	"fmt"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/shopspring/decimal"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirconv"
	"github.com/verily-src/fhirpath-go/internal/protofields"
)

var (
	ErrCantBeCast = errors.New("value can't be cast to system type")
)

// Any is the root abstraction for all FHIRPath system types.
type Any interface {
	isSystemType()
	Name() string
	Less(input Any) (Boolean, error)
}

// Stub methods on each type to implement interface Any.
func (s String) isSystemType()   {}
func (b Boolean) isSystemType()  {}
func (i Integer) isSystemType()  {}
func (d Decimal) isSystemType()  {}
func (d Date) isSystemType()     {}
func (t Time) isSystemType()     {}
func (d DateTime) isSystemType() {}
func (q Quantity) isSystemType() {}

// IsValid validates whether the input string represents
// a valid system type name.
func IsValid(typeName string) bool {
	switch typeName {
	case stringType, booleanType, integerType, decimalType,
		dateType, timeType, dateTimeType, quantityType, anyType:
		return true
	default:
		return false
	}
}

// IsPrimitive evaluates to check whether or not the input
// is a primitive FHIR type. If so, returns true, otherwise
// returns false
func IsPrimitive(input any) bool {
	switch v := input.(type) {
	case *dtpb.Boolean, *dtpb.String, *dtpb.Uri, *dtpb.Url, *dtpb.Canonical, *dtpb.Code, *dtpb.Oid, *dtpb.Id, *dtpb.Uuid, *dtpb.Markdown,
		*dtpb.Base64Binary, *dtpb.Integer, *dtpb.UnsignedInt, *dtpb.PositiveInt, *dtpb.Decimal, *dtpb.Date,
		*dtpb.Time, *dtpb.DateTime, *dtpb.Instant, *dtpb.Quantity, Any:
		return true
	case fhir.Base:
		return protofields.IsCodeField(v)
	default:
		return false
	}
}

// From converts primitive FHIR types to System types.
// Returns the input if already a System type, and an error
// if the input is not convertible.
func From(input any) (Any, error) {
	switch v := input.(type) {
	case *dtpb.Boolean:
		return Boolean(v.Value), nil
	case *dtpb.String:
		return String(v.Value), nil
	case *dtpb.Uri:
		return String(v.Value), nil
	case *dtpb.Url:
		return String(v.Value), nil
	case *dtpb.Code:
		return String(v.Value), nil
	case *dtpb.Oid:
		return String(v.Value), nil
	case *dtpb.Id:
		return String(v.Value), nil
	case *dtpb.Uuid:
		return String(v.Value), nil
	case *dtpb.Markdown:
		return String(v.Value), nil
	case *dtpb.Base64Binary:
		return String(base64.StdEncoding.EncodeToString(v.Value)), nil
	case *dtpb.Canonical:
		return String(v.Value), nil
	case *dtpb.Integer:
		return Integer(v.Value), nil
	case *dtpb.UnsignedInt:
		return Integer(v.Value), nil
	case *dtpb.PositiveInt:
		return Integer(v.Value), nil
	case *dtpb.Decimal:
		value, err := decimal.NewFromString(v.Value)
		if err != nil {
			return nil, err
		}
		return Decimal(value), nil
	case *dtpb.Date:
		value, err := DateFromProto(v)
		if err != nil {
			return nil, err
		}
		return value, nil
	case *dtpb.Time:
		return TimeFromProto(v), nil
	case *dtpb.DateTime:
		value, err := DateTimeFromProto(v)
		if err != nil {
			return nil, err
		}
		return value, nil
	case *dtpb.Instant:
		value, err := ParseDateTime(fhirconv.InstantToString(v))
		if err != nil {
			return nil, err
		}
		return value, nil
	case *dtpb.Quantity:
		unit := v.GetUnit().GetValue()
		if v.GetValue() == nil {
			return Quantity{unit: unit}, nil
		}
		value, err := decimal.NewFromString(v.GetValue().GetValue())
		if err != nil {
			return nil, err
		}
		return Quantity{Decimal(value), unit}, nil
	case Any:
		return v, nil
	case fhir.Base:
		value, ok := protofields.StringValueFromCodeField(v)
		if !ok {
			return nil, fmt.Errorf("%w: complex type %T", ErrCantBeCast, input)
		}
		return String(value), nil
	default:
		return nil, fmt.Errorf("%w: %T", ErrCantBeCast, input)
	}
}

// Normalize casts the "from" type to the "to" type if implicit casting
// is supported between the types. Otherwise, it returns the from input.
func Normalize(from Any, to Any) Any {
	switch v := from.(type) {
	case Integer:
		if _, ok := to.(Decimal); ok {
			return Decimal(decimal.NewFromInt32(int32(v)))
		}
		if q, ok := to.(Quantity); ok {
			dec := Decimal(decimal.NewFromInt32(int32(v)))
			return Quantity{dec, q.unit}
		}
	case Decimal:
		if q, ok := to.(Quantity); ok {
			return Quantity{v, q.unit}
		}
	case Date:
		if _, ok := to.(DateTime); ok {
			newLayout := v.l + "T"
			return DateTime{v.date, newLayout}
		}
	default:
		return from
	}
	return from
}
