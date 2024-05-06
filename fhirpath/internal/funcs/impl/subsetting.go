package impl

import (
	"fmt"

	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"google.golang.org/protobuf/proto"
)

// First Returns a collection containing only the first item in the input collection.
// This function is equivalent to item[0], so it will return an empty collection if the input collection has no items.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#first-collection
func First(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	return system.Collection{input[0]}, nil
}

// Last Returns a collection containing only the last item in the input collection.
// Will return an empty collection if the input collection has no items.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#last-collection
func Last(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	return system.Collection{input[len(input)-1]}, nil
}

// Tail Returns a collection containing all but the first item in the input collection.
// Will return an empty collection if the input collection has no items, or only one item.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#tail-collection
func Tail(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	return input[1:], nil
}

// Skip Returns a collection containing all but the first num items in the input collection.
// Will return an empty collection if there are no items remaining after the indicated number of items have been skipped,
// or if the input collection is empty.
// If num is less than or equal to zero, the input collection is simply returned.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#skipnum-integer-collection
func Skip(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	// Args validation
	if len(args) != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, len(args))
	}
	argValues, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	}
	skip, err := argValues.ToInt32()
	if err != nil {
		return nil, err
	}
	if skip <= 0 {
		return input, nil
	}
	if skip >= int32(len(input)) {
		return system.Collection{}, nil
	}
	return input[skip:], nil
}

// Take Returns a collection containing the first num items in the input collection,
// or less if there are less than num items. If num is less than or equal to 0,
// or if the input collection is empty ({ }), take returns an empty collection.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#takenum-integer-collection
func Take(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	// Args validation
	if len(args) != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, len(args))
	}
	argValues, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	}
	take, err := argValues.ToInt32()
	if err != nil {
		return nil, err
	}
	if take <= 0 {
		return system.Collection{}, nil
	}
	if take >= int32(len(input)) {
		return input, nil
	}
	return input[:take], nil
}

// Intersect Returns the set of elements that are in both collections.
// Duplicate items will be eliminated by this function.
// Order of items is not guaranteed to be preserved in the result of this function.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#intersectother-collection-collection
func Intersect(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	// Args validation
	if len(args) != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, len(args))
	}
	argValues, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	}
	var result system.Collection
	for _, i := range input {
		for _, c := range argValues {
			if checkEquality(i, c) {
				v, _ := system.From(c)
				result = append(result, v)
			}
		}
	}
	if len(result) == 0 {
		return system.Collection{}, nil
	}
	return removeDuplicates(result), nil
}

// Exclude returns the set of elements that are not in the other collection.
// Duplicate items will not be eliminated by this function, and order will be preserved.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#excludeother-collection-collection
func Exclude(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Input validation
	if input.IsEmpty() {
		return system.Collection{}, nil
	}
	// Args validation
	if len(args) != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, len(args))
	}
	argValues, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	}
	var result system.Collection
	for _, val := range input {
		if !argValues.Contains(val) {
			result = append(result, val)
		}
	}
	for _, arg := range argValues {
		if !input.Contains(arg) {
			result = append(result, arg)
		}
	}
	return result, nil
}

// Distinct returns the set of elements that are distinct and unique from the
// input by applying equality-operation tests.
//
// If the input collection is empty ({}), the result is empty.
// Note that the order of elements in the input collection is not guaranteed to
// be preserved in the result.
//
// See the spec for this function for more details:
// https://hl7.org/fhirpath/N1/#distinct-collection
func Distinct(_ *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	var result system.Collection
	for _, v := range input {
		if result.Contains(v) {
			continue
		}
		result = append(result, v)
	}
	return result, nil
}

// IsDistinct queries whether the input collection is a set of fully distinct
// and unique values. This is effectively short-hand for calling:
//
//	v.count() = v.distinct().count()
//
// See the spec for this function for more details:
// https://hl7.org/fhirpath/N1/#isdistinct-boolean
func IsDistinct(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	got, err := Distinct(ctx, input, args...)
	if err != nil {
		return nil, err
	}
	return system.Collection{system.Boolean(len(got) == len(input))}, nil
}

func removeDuplicates(collection system.Collection) system.Collection {
	seen := make(map[any]bool)
	var result system.Collection
	for _, val := range collection {
		if _, ok := seen[val]; !ok {
			seen[val] = true
			result = append(result, val)
		}
	}
	return result
}

func checkEquality(lhs, rhs any) bool {
	return checkSystemEquality(lhs, rhs) || checkProtoEquality(lhs, rhs)
}

func checkSystemEquality(lhs, rhs any) bool {
	l, lerr := system.From(lhs)
	r, rerr := system.From(rhs)
	if lerr == nil && rerr == nil {
		got, ok := system.TryEqual(l, r)
		return got && ok
	}
	return false
}

func checkProtoEquality(lhs, rhs any) bool {
	l, lok := lhs.(proto.Message)
	r, rok := rhs.(proto.Message)
	if lok && rok {
		return proto.Equal(l, r)
	}
	return false
}
