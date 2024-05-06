package impl

import (
	"fmt"

	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
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
