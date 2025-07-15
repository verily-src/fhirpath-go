package impl

import (
	"fmt"
	"strings"

	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/reflection"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

// Where evaluates the expression args[0] on each input item, collects the items that cause
// the expression to evaluate to true.
func Where(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, len(args))
	}
	e := args[0]
	result := system.Collection{}
	for _, item := range input {
		output, err := e.Evaluate(ctx, system.Collection{item})
		if err != nil {
			return nil, err
		}
		if len(output) == 0 {
			continue
		}
		pass, err := output.ToSingletonBoolean()
		if err != nil {
			return nil, fmt.Errorf("evaluating where condition as boolean resulted in an error: %w", err)
		}
		if pass[0] {
			result = append(result, item)
		}
	}
	return result, nil
}

// OfType filters the input collection to include only items that are of the specified type
// or a subclass thereof. The type is specified as an argument in the form of an identifier.
// If the input collection is empty, the result will also be empty.
func OfType(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, len(args))
	}

	var typeSpecifier reflection.TypeSpecifier
	typeExpr, ok := args[0].(*expr.TypeExpression)
	if !ok {
		return nil, fmt.Errorf("received invalid argument, expected a type")
	}
	var err error
	if parts := strings.Split(typeExpr.Type, "."); len(parts) == 2 {
		if typeSpecifier, err = reflection.NewQualifiedTypeSpecifier(parts[0], parts[1]); err != nil {
			return nil, err
		}
	} else if typeSpecifier, err = reflection.NewTypeSpecifier(typeExpr.Type); err != nil {
		return nil, err
	}
	result := system.Collection{}
	for _, item := range input {
		typ, err := reflection.TypeOf(item)
		if err != nil {
			return nil, err
		}
		if typ.Is(typeSpecifier) {
			result = append(result, item)
		}
	}
	return result, nil
}
