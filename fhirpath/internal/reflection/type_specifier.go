package reflection

import (
	"errors"
	"fmt"
	"strings"

	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/protofields"
)

var (
	errInvalidType      = errors.New("invalid type name")
	errInvalidNamespace = errors.New("invalid namespace")
	errInvalidInput     = errors.New("invalid input type")
)

// TypeSpecifier is a FHIRPath type that enables use of
// is and as operators. Provides a namespace and type name
type TypeSpecifier struct {
	namespace string
	typeName  string
}

// NewQualifiedTypeSpecifier constructs a Qualified Type Specifier given a namespace and typeName.
// Returns an error if the typeName is not found within the namespace, or the namespace is invalid.
func NewQualifiedTypeSpecifier(namespace string, typeName string) (TypeSpecifier, error) {
	switch namespace {
	case FHIR:
		if IsValidFHIRPathElement(typeName) || protofields.IsValidResourceType(typeName) || isBaseType(typeName) {
			return TypeSpecifier{namespace: namespace, typeName: typeName}, nil
		}
		return TypeSpecifier{}, fmt.Errorf("%w: %s", errInvalidType, typeName)
	case System:
		if system.IsValid(typeName) {
			return TypeSpecifier{namespace: namespace, typeName: typeName}, nil
		}
		return TypeSpecifier{}, fmt.Errorf("%w: %s", errInvalidType, typeName)
	default:
		return TypeSpecifier{}, fmt.Errorf("%w: %s", errInvalidNamespace, namespace)
	}
}

// NewTypeSpecifier constructs a Qualified Type Specifier given a typeName. The namespace
// is inferred with the priority rules of FHIRPath. Returns an error if the typeName cannot
// be resolved.
func NewTypeSpecifier(typeName string) (TypeSpecifier, error) {
	if parts := strings.Split(typeName, "."); len(parts) == 2 {
		return NewQualifiedTypeSpecifier(parts[0], parts[1])
	}

	if IsValidFHIRPathElement(typeName) || protofields.IsValidResourceType(typeName) || isBaseType(typeName) {
		return TypeSpecifier{FHIR, typeName}, nil
	}
	if system.IsValid(typeName) {
		return TypeSpecifier{System, typeName}, nil
	}
	return TypeSpecifier{}, fmt.Errorf("%w: %s", errInvalidType, typeName)
}

// TypeOf retrieves the Type Specifier of the input, given that it is
// a supported FHIRPath type. Otherwise, returns an error.
func TypeOf(input any) (TypeSpecifier, error) {
	if item, ok := input.(system.Any); ok {
		return TypeSpecifier{System, item.Name()}, nil
	}
	item, ok := input.(fhir.Base)
	if !ok {
		return TypeSpecifier{}, fmt.Errorf("%w: no type specifier available", errInvalidInput)
	}
	if oneOf := protofields.UnwrapOneofField(item, "choice"); oneOf != nil {
		item = oneOf
	}
	name := string(item.ProtoReflect().Descriptor().Name())
	if protofields.IsCodeField(item) {
		return TypeSpecifier{FHIR, "code"}, nil
	}
	return TypeSpecifier{FHIR, primitiveToLowercase(name)}, nil
}

// String returns a string representation of the type specifier in the format:
// <namespace>.<type>
func (ts TypeSpecifier) String() string {
	if ts.namespace != "" {
		return ts.namespace + "." + ts.typeName
	}
	return ts.typeName
}

// Is returns a boolean representing whether or not the receiver type is equivalent to the
// input type, or if it's a valid subtype.
func (ts TypeSpecifier) Is(input TypeSpecifier) system.Boolean {
	if ts.namespace != input.namespace {
		return false
	}
	// If the root type has been reached and the equality is still false, they are not equal
	if ts == ts.parent() && ts.typeName != input.typeName {
		return false
	}
	if ts.typeName == input.typeName {
		return true
	}
	return ts.parent().Is(input) // Recursively compare the parent type
}

// MustCreateTypeSpecifier creates a qualified type specifier and panics if the
// provided namespace or typeName is invalid. Returns the created TypeSpecifier
func MustCreateTypeSpecifier(namespace string, typeName string) TypeSpecifier {
	typeSpecifier, err := NewQualifiedTypeSpecifier(namespace, typeName)
	if err != nil {
		panic(err)
	}
	return typeSpecifier
}

func (ts TypeSpecifier) parent() TypeSpecifier {
	if ts.namespace == System {
		return TypeSpecifier{"System", "Any"}
	}

	switch ts.typeName {
	case "code", "markdown", "id":
		return TypeSpecifier{ts.namespace, "string"}
	case "unsignedInt", "positiveInt":
		return TypeSpecifier{ts.namespace, "integer"}
	case "url", "canonical", "uuid", "oid":
		return TypeSpecifier{ts.namespace, "uri"}
	case "Duration", "MoneyQuantity", "Age", "Count", "Distance", "SimpleQuantity":
		return TypeSpecifier{ts.namespace, "Quantity"}
	case "Timing", "Dosage", "ElementDefinition":
		return TypeSpecifier{ts.namespace, "BackboneElement"}
	case "Bundle", "Binary", "Parameters", "DomainResource":
		return TypeSpecifier{ts.namespace, "Resource"}
	case "Element", "Resource":
		return ts
	default:
		if IsValidFHIRPathElement(ts.typeName) {
			return TypeSpecifier{ts.namespace, "Element"}
		}
		return TypeSpecifier{ts.namespace, "DomainResource"}
	}
}

func isBaseType(name string) bool {
	switch name {
	case "Element", "Resource", "DomainResource":
		return true
	}
	return false
}
