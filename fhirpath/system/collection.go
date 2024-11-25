package system

import (
	"errors"
	"fmt"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/shopspring/decimal"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/narrow"
	"google.golang.org/protobuf/proto"
)

var (
	// ErrNotConvertible is an error raised when attempting to call Collection.To*
	// to a type that is not convertible.
	ErrNotConvertible = errors.New("not convertible")
)

// Collection abstracts the input and output type for
// FHIRPath expressions as a collection that can contain anything.
type Collection []any

// IsSingleton is a utility function to determine if a collection is a
// "Singleton" collection (i.e. contains 1 value).
func (c Collection) IsSingleton() bool {
	return len(c) == 1
}

// IsEmpty is a utility function to determine if a collection is empty.
func (c Collection) IsEmpty() bool {
	return len(c) == 0
}

// TryEqual compares this collection to the supplied collection for, comparing each entry
// with the corresponding entry in c2. Returns true if every entry
// in c1 is equivalent to that of c2. If the lengths are mismatched,
// returns false.
func (c Collection) TryEqual(other Collection) (bool, bool) {
	if len(c) != len(other) {
		return false, true
	}

	for i := range other {
		okOne := IsPrimitive(c[i])
		okTwo := IsPrimitive(other[i])
		if okOne != okTwo {
			return false, true
		}
		if !okOne && !proto.Equal(c[i].(fhir.Base), other[i].(fhir.Base)) {
			return false, true
		}
		if !okOne {
			return true, true
		}
		primitiveOne, err := From(c[i])
		if err != nil {
			return false, true
		}
		primitiveTwo, err := From(other[i])
		if err != nil {
			return false, true
		}
		primitiveOne = Normalize(primitiveOne, primitiveTwo)
		primitiveTwo = Normalize(primitiveTwo, primitiveOne)
		equal, ok := TryEqual(primitiveOne, primitiveTwo)
		if !ok {
			return false, false
		}
		if !equal {
			return false, true
		}
	}
	return true, true
}

// ToSingletonBoolean evaluates a collection as a boolean with singleton evaluation of
// collection rules. Returns a collection containing a single Boolean, or empty if the
// input is empty.
func (c Collection) ToSingletonBoolean() ([]Boolean, error) {
	length := len(c)
	if length == 0 {
		return []Boolean{}, nil
	}
	if length > 1 {
		return nil, fmt.Errorf("collection can't evaluate to bool, contains %v elements", length)
	}
	val, _ := From(c[0])
	if boolean, ok := val.(Boolean); ok {
		return []Boolean{boolean}, nil
	}
	return []Boolean{true}, nil
}

// ToSingleton evaluates a collection as a single value, returning an error if the
// collection contains 0 or more than 1 entry.
func (c Collection) ToSingleton() (any, error) {
	if !c.IsSingleton() {
		return nil, fmt.Errorf("collection is not singleton")
	}
	return c[0], nil
}

// ToBool converts this Collection into a Go native 'bool' type, following the
// logic of singleton evaluation of booleans. If this collection is empty,
// it returns false, if it contains more than 1 entry, it will return an error.
func (c Collection) ToBool() (bool, error) {
	if c.IsEmpty() {
		return false, nil
	}
	v, err := c.ToSingleton()
	if err != nil {
		return false, err
	}
	switch val := v.(type) {
	case Boolean:
		return bool(val), nil
	case *dtpb.Boolean:
		return val.GetValue(), nil
	}
	// A single value that is not bool is always implicitly 'true'
	return true, nil
}

// ToInt32 converts this Collection into a Go native 'int32' type.
// If this collection is empty, or contains more than 1 entry, it will return
// an error. If the type in the collection is not a System.Integer, or something
// derived from a FHIR.Integer, this will raise an ErrNotConvertible.
func (c Collection) ToInt32() (int32, error) {
	v, err := c.ToSingleton()
	if err != nil {
		return 0, err
	}
	switch val := v.(type) {
	case Integer:
		return int32(val), nil
	case *dtpb.Integer:
		return val.GetValue(), nil
	case *dtpb.PositiveInt:
		if val, ok := narrow.ToInt32(val.GetValue()); ok {
			return val, nil
		}
		return 0, c.convertErr(val.GetValue(), "int32")
	case *dtpb.UnsignedInt:
		if val, ok := narrow.ToInt32(val.GetValue()); ok {
			return val, nil
		}
		return 0, c.convertErr(val.GetValue(), "int32")
	}
	return 0, c.convertErr(v, "int32")
}

// ToFloat64 converts this Collection into a Go native 'float64' type.
// If this collection is empty, or contains more than 1 entry, it will return
// an error. If the type in the collection is not a System.Integer, or something
// derived from a FHIR.Integer, this will raise an ErrNotConvertible.
func (c Collection) ToFloat64() (float64, error) {
	v, err := c.ToSingleton()
	if err != nil {
		return 0, err
	}
	switch val := v.(type) {
	case Decimal:
		return decimal.Decimal(val).InexactFloat64(), nil
	case Integer:
		return float64(val), nil
	case *dtpb.Integer:
		return float64(val.GetValue()), nil
	case *dtpb.PositiveInt:
		return float64(val.GetValue()), nil
	case *dtpb.UnsignedInt:
		return float64(val.GetValue()), nil
	}
	return 0, c.convertErr(v, "float64")
}

// ToString converts this Collection into a Go native 'string' type.
// If this collection is empty, or contains more than 1 entry, it will return
// an error. If the type in the collection is not a System.String, or something
// derived from a FHIR.String, this will raise an ErrNotConvertible.
func (c Collection) ToString() (string, error) {
	v, err := c.ToSingleton()
	if err != nil {
		return "", err
	}
	result, err := From(v)
	if err != nil {
		return "", c.convertErr(v, "string")
	}
	if str, ok := result.(String); ok {
		return string(str), nil
	}
	return "", c.convertErr(v, "string")
}

// Contains returns true if the specified value is contained within this
// collection. This will normalize types to system-types if necessary to check
// for containment.
func (c Collection) Contains(value any) bool {
	sys, err := From(value)
	if err == nil {
		return c.containsSystem(sys)
	}
	msg, ok := value.(proto.Message)
	if ok {
		return c.containsProto(msg)
	}
	return false
}

func (c Collection) containsSystem(value Any) bool {
	for _, v := range c {
		sys, err := From(v)
		if err != nil {
			continue
		}
		if Equal(sys, value) {
			return true
		}
	}
	return false
}

func (c Collection) containsProto(value proto.Message) bool {
	for _, v := range c {
		msg, ok := v.(proto.Message)
		if !ok {
			continue
		}
		if proto.Equal(msg, value) {
			return true
		}
	}
	return false
}

func (c Collection) convertErr(got any, want string) error {
	return fmt.Errorf("type %T %w to %v", got, ErrNotConvertible, want)
}
