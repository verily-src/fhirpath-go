package fhir

import "strings"

// These characters have special meaning in FHIR Search queries
const SearchSpecialChars = `\,$|`

// Escape values intended for use as a parameter in a FHIR Search.
//
// These characters have special meaning in Search queries and must be backslash escaped:
//
//	`\`, `|`, `,`, `$`
//
// This function assumes that URL-encoding is performed later. (Percent
// encoding is automatically handled by the healthcare client library when
// query params are passed as a map.)
//
// For example, `foo,bar` becomes `foo\,bar`
func EscapeSearchParam(value string) string {
	out := value
	for _, crune := range SearchSpecialChars {
		c := string(crune)
		out = strings.ReplaceAll(out, c, `\`+c)
	}
	return out
}
