package impl

import (
	"fmt"

	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"google.golang.org/protobuf/proto"
)

// AllTrue Takes a collection of Boolean values and returns true if all the items are true.
// If any items are false, the result is false. If the input is empty, the result is true.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#alltrue-boolean
func AllTrue(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validations
	if input.IsEmpty() {
		return system.Collection{system.Boolean(true)}, nil
	}
	for _, v := range input {
		value, _ := system.From(v)
		if value == system.Boolean(false) {
			return system.Collection{system.Boolean(false)}, nil
		}
	}
	return system.Collection{system.Boolean(true)}, nil
}

// AnyTrue takes a collection of Boolean values and returns true if any of the items are true.
// If all the items are false, or if the input is empty ({ }), the result is false.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#anytrue-boolean
func AnyTrue(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validations
	if input.IsEmpty() {
		return system.Collection{system.Boolean(false)}, nil
	}
	for _, v := range input {
		value, _ := system.From(v)
		if value == system.Boolean(true) {
			return system.Collection{system.Boolean(true)}, nil
		}
	}
	return system.Collection{system.Boolean(false)}, nil
}

// AllFalse takes a collection of Boolean values and returns true if all the items are false.
// If any items are true, the result is false. If the input is empty, the result is true.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#allfalse-boolean
func AllFalse(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validations
	if input.IsEmpty() {
		return system.Collection{system.Boolean(true)}, nil
	}
	for _, v := range input {
		value, _ := system.From(v)
		if value == system.Boolean(true) {
			return system.Collection{system.Boolean(false)}, nil
		}
	}
	return system.Collection{system.Boolean(true)}, nil
}

// AnyFalse takes a collection of Boolean values and returns true if any of the items are false.
// If all the items are true, or if the input is empty ({ }), the result is false.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#anyfalse-boolean
func AnyFalse(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validations
	if input.IsEmpty() {
		return system.Collection{system.Boolean(false)}, nil
	}
	for _, v := range input {
		value, _ := system.From(v)
		if value == system.Boolean(false) {
			return system.Collection{system.Boolean(true)}, nil
		}
	}
	return system.Collection{system.Boolean(false)}, nil
}

// All returns true if for every element in the input collection, criteria evaluates to true.
// Otherwise, the result is false. If the input collection is empty ({}), the result is true.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#allcriteria-expression-boolean
func All(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// If the input collection is empty, return true
	if input.IsEmpty() {
		return system.Collection{system.Boolean(true)}, nil
	}

	// Validate that exactly one argument (the criteria expression) is provided
	if len(args) != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, len(args))
	}

	// Evaluate the criteria expression for each element in the input collection
	for _, element := range input {
		// Evaluate the criteria expression
		output, err := args[0].Evaluate(ctx, system.Collection{element})
		if err != nil {
			return nil, fmt.Errorf("evaluating criteria expression resulted in an error: %w", err)
		}

		// Check that the output for false
		// if 'err' is non-nil, `v` is false
		if v, err := output.ToBool(); err != nil {
			return nil, err
		} else if !v {
			return system.Collection{system.Boolean(false)}, nil
		}
	}

	// Return true if all elements satisfy the criteria
	return system.Collection{system.Boolean(true)}, nil
}

// Count returns the integer count of the number of items in the input collection.
// returns 0 when the input collection is empty.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#convertstodate-boolean
func Count(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	return system.Collection{system.Integer(len(input))}, nil
}

// Exists evaluates the expression args[0] on each input item, returns whether
// there exists at least one item that cause the expression to evaluate to true.
// http://hl7.org/fhirpath/N1/#existscriteria-expression-boolean
func Exists(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	if len(args) == 0 {
		return system.Collection{system.Boolean(len(input) > 0)}, nil
	}
	whereOutput, err := Where(ctx, input, args...)
	if err != nil {
		return nil, fmt.Errorf("calling Where(): %w", err)
	}
	return system.Collection{system.Boolean(len(whereOutput) > 0)}, nil
}

// Empty evaluates the expression args[0] on each input item, returns whether
// none of the items causes the expression to evaluate to true.
// http://hl7.org/fhirpath/N1/#empty-boolean
func Empty(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	if len(args) == 0 {
		return system.Collection{system.Boolean(len(input) == 0)}, nil
	}
	return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, len(args))
}

// SubsetOf returns true if all items in the input collection are members of the
// collection passed as the other argument. Membership is determined using the
// = (Equals) operation.
//
// If the input collection is empty, the result is true.
// If the other collection is empty but input is not, the result is false.
//
// https://hl7.org/fhirpath/N1/index.html#subsetofother-collection-boolean
func SubsetOf(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Validate exactly one argument is provided
	if len(args) != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, len(args))
	}

	// Evaluate the other collection argument in the current context
	otherCollection, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	}

	// Empty input collection is always a subset (conceptually true)
	if input.IsEmpty() {
		return system.Collection{system.Boolean(true)}, nil
	}

	// Non-empty input with empty other collection means false
	if otherCollection.IsEmpty() {
		return system.Collection{system.Boolean(false)}, nil
	}

	// Implement bag semantics using a boolean slice to track used elements.
	used := make([]bool, len(otherCollection))

	// Check each input element for membership in the other collection
	for _, inputItem := range input {
		found := false

		// Search for an unused matching element in the other collection
		for i, otherItem := range otherCollection {
			// Skip if this element has already been used (bag semantics)
			if used[i] {
				continue
			}

			// Convert both items to system types and check equality
			sysInput, errInput := system.From(inputItem)
			sysOther, errOther := system.From(otherItem)

			// Use FHIRPath equality if both conversions succeeded
			if errInput == nil && errOther == nil {
				if eq, hasResult := system.TryEqual(sysInput, sysOther); hasResult && eq {
					// Element found - mark it as used to prevent reuse (bag semantics)
					used[i] = true
					found = true
					break
				}
			} else {
				// Fallback to proto.Equal comparison
				inputMsg, inputOK := inputItem.(proto.Message)
				otherMsg, otherOK := otherItem.(proto.Message)
				if inputOK && otherOK && proto.Equal(inputMsg, otherMsg) {
					// Proto equality check succeeded, mark it as used
					used[i] = true
					found = true
					break
				}
			}
		}

		// If no matching unused element found, input is not a subset
		if !found {
			return system.Collection{system.Boolean(false)}, nil
		}
	}

	return system.Collection{system.Boolean(true)}, nil
}
