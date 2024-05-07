package parser

import "github.com/verily-src/fhirpath-go/fhirpath/internal/expr"

// A VisitorTransform is a function which transforms the specified
// expression. This is used in FHIRPath Patch to modify expressions.
type VisitorTransform func(expr.Expression) expr.Expression

// IdentityTransform returns the given expression without any modification.
func IdentityTransform(e expr.Expression) expr.Expression {
	return e
}
