package funcs

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs/impl"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

var (
	errNotFunc       = errors.New("value is not a function")
	errMissingArgs   = errors.New("missing arguments")
	errInvalidParams = errors.New("invalid input parameters")
	errInvalidReturn = errors.New("invalid function return signature")
)

var notImplemented = Function{Func: unimplemented}

// FHIRPathFunc is the common abstraction for all function types
// supported by FHIRPath.
type FHIRPathFunc func(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error)

// Function contains the FHIRPathFunction, along with metadata
type Function struct {
	Func           FHIRPathFunc
	MinArity       int
	MaxArity       int
	IsTypeFunction bool
}

// ToFunction takes in a function with any arguments and attempts to
// convert it to a functions.Function type. If the conversion is successful,
// the new function will assert the argument expressions resolve to the original
// argument types.
func ToFunction(fn any) (Function, error) {
	rv := reflect.ValueOf(fn)
	if err := validateFunc(rv); err != nil {
		return Function{}, fmt.Errorf("constructing FHIRPathFunction: %w", err)
	}

	arity := rv.Type().NumIn() - 1 // 'True' arity, as the first argument is the input Collection.
	fhirpathFunc := func(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
		if len(args) != arity {
			return nil, fmt.Errorf("%w: function expects %v arguments, received %v", impl.ErrWrongArity, arity, len(args))
		}
		// ensure the arguments match the function signature.
		funcArgs := []reflect.Value{reflect.ValueOf(input)}
		for i, exp := range args {
			result, err := exp.Evaluate(ctx, input)
			if err != nil {
				return nil, err
			}
			if len(result) != 1 {
				return nil, fmt.Errorf("%w: doesn't return singleton", impl.ErrInvalidReturnType)
			}
			if expectedType, gotType := rv.Type().In(i+1), reflect.TypeOf(result[0]); !gotType.AssignableTo(expectedType) {
				return nil, fmt.Errorf("%w: got type '%s' when type '%s' was expected", impl.ErrInvalidReturnType, gotType.String(), expectedType.Name())
			}
			funcArgs = append(funcArgs, reflect.ValueOf(result[0]))
		}
		output := rv.Call(funcArgs)
		if err, ok := output[1].Interface().(error); ok {
			return output[0].Interface().(system.Collection), err
		}
		return output[0].Interface().(system.Collection), nil
	}
	return Function{fhirpathFunc, arity, arity, false}, nil
}

// validateFunc verifies that the input reflect value represents a
// valid FHIRPath function. If not, it returns an error.
func validateFunc(rv reflect.Value) error {
	if rv.Kind() != reflect.Func {
		return errNotFunc
	}
	errs := []error{}
	if rv.Type().NumIn() < 1 {
		errs = append(errs, errMissingArgs)
	} else if rv.Type().In(0) != reflect.TypeOf(system.Collection{}) {
		errs = append(errs, errInvalidParams)
	}
	if rv.Type().NumOut() != 2 || rv.Type().Out(0) != reflect.TypeOf(system.Collection{}) || rv.Type().Out(1).Name() != "error" {
		errs = append(errs, errInvalidReturn)
	}
	return errors.Join(errs...)
}

// unimplemented is a no-op placeholder function that satisfies the FHIRPathFunction contract
func unimplemented(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	return nil, fmt.Errorf("function not yet implemented")
}
