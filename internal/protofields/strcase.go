package protofields

import (
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// toSnakeCase is a helper function to convert CamelCase names to snake_case.
// This is needed for finding fields in the Proto descriptors, which are snake_case,
// from resource-names that are CamelCase.
//
// Note: strcase.ToSnake does not work for converting Base64Binary to
// base64_binary, so this function exists to do it for us with the semantics we
// want.
func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
