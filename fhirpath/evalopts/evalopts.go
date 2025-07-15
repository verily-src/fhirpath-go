/*
Package evalopts provides EvaluateOption values for FHIRPath.

This package exists to isolate the options away from the core FHIRPath logic,
since this will simplify discovery of evaluation-specific options.
*/
package evalopts

import (
	"errors"
	"fmt"
	"time"

	"github.com/verily-src/fhirpath-go/fhirpath/internal/opts"
	"github.com/verily-src/fhirpath-go/fhirpath/resolver"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

var (
	ErrUnsupportedType  = errors.New("external constant type not supported")
	ErrExistingConstant = errors.New("constant already exists")
)

// OverrideTime returns an EvaluateOption that can be used to override the time
// that will be used in FHIRPath expressions.
func OverrideTime(t time.Time) opts.EvaluateOption {
	return opts.Transform(func(cfg *opts.EvaluateConfig) error {
		cfg.Context.Now = t
		return nil
	})
}

// EnvVariable returns an EvaluateOption that sets FHIRPath environment variables
// (e.g. %action).
//
// The input must be one of:
//   - A FHIRPath System type,
//   - A FHIR Element or Resource type, or
//   - A FHIRPath Collection, containing the above types.
//
// If an EnvVariable is specified that already exists in the expression, then
// evaluation will yield an ErrExistingConstant error. If an EnvVariable is
// contains a type that is not one of the above valid types, then evaluation
// will yield an ErrUnsupportedType error.
func EnvVariable(name string, value any) opts.EvaluateOption {
	return opts.Transform(func(cfg *opts.EvaluateConfig) error {
		if err := validateType(value); err != nil {
			return err
		}
		if _, ok := cfg.Context.ExternalConstants[name]; !ok {
			cfg.Context.ExternalConstants[name] = value
			return nil
		}
		return fmt.Errorf("%w: %s", ErrExistingConstant, name)
	})
}

// validateType validates that the input type is a supported
// fhir proto or System type. If a system.Collection is passed in,
// recursively checks each element.
func validateType(input any) error {
	var err error
	switch v := input.(type) {
	case fhir.Base, system.Any:
		break
	case system.Collection:
		for _, elem := range v {
			err = errors.Join(err, validateType(elem))
		}
	default:
		err = fmt.Errorf("%w: %T", ErrUnsupportedType, input)
	}
	return err
}

// WithResolver returns an EvaluateOption that sets the FHIR reference resolver.
func WithResolver(resolver resolver.Resolver) opts.EvaluateOption {
	return opts.Transform(func(cfg *opts.EvaluateConfig) error {
		cfg.Context.Resolver = resolver
		return nil
	})
}
