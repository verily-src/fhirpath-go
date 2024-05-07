/*
Package opts is an internal package that exists for setting configuration
settings for FHIRPath. This is an internal package so that only parts of this
may be publicly re-exported, while the implementation has access to the full
thing.
*/
package opts

import (
	"errors"

	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/parser"
)

// CompileConfig provides the configuration values for the Compile command.
type CompileConfig struct {
	// Table is the current function table to be called.
	Table     funcs.FunctionTable
	Transform parser.VisitorTransform

	// Permissive is a legacy option to allow FHIRpaths with *invalid* fields to be
	// compiled (to reduce breakages).
	Permissive bool
}

// EvaluateConfig provides the configuration values for the Evaluate command.
type EvaluateConfig struct {
	// Context is the current context information.
	Context *expr.Context
}

// Option is the base interface for FHIRPath options.
type Option[T any] interface {
	updateConfig(*T) error
}

// CompileOption is an Option that sets CompileConfig.
type CompileOption = Option[CompileConfig]

// EvaluateOption is an Option that sets EvaluateConfig.
type EvaluateOption = Option[EvaluateConfig]

// Transform creates either an Evaluate or Compile configuration option, done
// as a function callback.
func Transform[T any](callback func(cfg *T) error) Option[T] {
	return callbackOption[T]{callback: callback}
}

// ApplyOptions applies all the options to the given configuration.
func ApplyOptions[T any](cfg *T, opts ...Option[T]) (*T, error) {
	var errs []error
	for _, opt := range opts {
		errs = append(errs, opt.updateConfig(cfg))
	}
	return cfg, errors.Join(errs...)
}

type callbackOption[T any] struct {
	callback func(*T) error
}

func (o callbackOption[T]) updateConfig(cfg *T) error {
	return o.callback(cfg)
}
