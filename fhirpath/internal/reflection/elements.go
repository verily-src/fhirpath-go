package reflection

import (
	"github.com/iancoleman/strcase"
	"github.com/verily-src/fhirpath-go/internal/protofields"
)

// IsValidFHIRPathElement checks if the input string represents
// a valid element name. This function is importantly case-sensitive,
// which is a distinction that is important for primitive types.
func IsValidFHIRPathElement(name string) bool {
	if isPrimitive(name) || name == "BackboneElement" {
		return true
	}
	return protofields.IsValidElementType(primitiveToLowercase(name))
}

func isPrimitive(name string) bool {
	switch name {
	case "instant", "time", "date", "dateTime", "base64Binary",
		"decimal", "boolean", "url", "code", "string", "integer", "uri",
		"canonical", "markdown", "id", "oid", "uuid", "unsignedInt", "positiveInt":
		return true
	default:
		return false
	}
}

func primitiveToLowercase(name string) string {
	switch name {
	case "Instant", "Time", "Date", "DateTime", "Base64Binary",
		"Decimal", "Boolean", "Url", "Code", "String", "Integer", "Uri",
		"Canonical", "Markdown", "Id", "Oid", "Uuid", "UnsignedInt", "PositiveInt":
		return strcase.ToLowerCamel(name)
	default:
		return name
	}
}
