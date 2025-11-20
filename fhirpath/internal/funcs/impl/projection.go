package impl

import (
	"errors"
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
	var fieldErrs []error
	for _, item := range input {
		output, err := e.Evaluate(ctx, system.Collection{item})
		// If the error is ErrInvalidField, don't immediately raise it
		if err != nil {
			if errors.Is(err, expr.ErrInvalidField) {
				fieldErrs = append(fieldErrs, err)
				continue
			}
			return nil, err
		}
		result = append(result, output...)
	}
	// Raise field errors if one was raised for each input.
	if len(input) > 0 && len(fieldErrs) == len(input) && !ctx.SkipUnknownFields {
		return nil, errors.Join(fieldErrs...)
	}
	return result, nil
}

// Repeat is a version of select that will repeat the projection and add it to the output
// collection, as long as the projection yields new items (as determined by the = (Equals)
// operator).
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#repeatprojection-expression-collection
func Repeat(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, len(args))
	}

	var result system.Collection

	for !input.IsEmpty() {
		out, err := Select(ctx, input, args...)
		if err != nil {
			return nil, err
		}

		if eq, _ := out.TryEqual(input); eq {
			break
		}

		result = append(result, out...)
		input = out
	}

	return result, nil
}
