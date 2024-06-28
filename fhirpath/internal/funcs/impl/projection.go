package impl

import (
	"fmt"

	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

// Select evaluates the expression args[0] on each input item. The result of each evaluation is
// added to the output collection.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#selectprojection-expression-collection
func Select(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
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
		for _, val := range output {
			if c, ok := val.(system.Collection); ok && c.IsEmpty() {
				continue
			}
			result = append(result, val)
		}
	}
	return result, nil
}
