/*
Package compopts provides CompileOption values for FHIRPath.

This package exists to isolate the options away from the core FHIRPath logic,
since this will simplify discovery of compile-specific options.
*/
package compopts

import (
	"errors"

	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/opts"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/parser"
)

var (
	ErrMultipleTransforms = errors.New("multiple transforms provided")
)

// AddFunction creates a CompileOption that will register a custom FHIRPath
// function that can be called during evaluation with the given name.
//
// If the function already exists, then compilation will return an error.
func AddFunction(name string, fn any) opts.CompileOption {
	return opts.Transform(func(cfg *opts.CompileConfig) error {
		return cfg.Table.Register(name, fn)
	})
}

// Transform creates a CompileOption that will set a transform
// to be called on each expression returned by the Visitor.
//
// If there is already a Transform set, then compilation will return an error.
func Transform(v parser.VisitorTransform) opts.CompileOption {
	return opts.Transform(func(cfg *opts.CompileConfig) error {
		if cfg.Transform != nil {
			return ErrMultipleTransforms
		}
		cfg.Transform = v
		return nil
	})
}

// Permissive is an option that enables deprecated behavior in FHIRPath field
// navigation. This can be used as a temporary fix for FHIRpaths that have never
// been valid FHIRPaths, but have worked up until this point.
//
// This option is marked Deprecated so that it nags users until the paths can
// be resolved.
//
// Deprecated: Please update FHIRPaths whenever possible.
func Permissive() opts.CompileOption {
	return opts.Transform(func(cfg *opts.CompileConfig) error {
		cfg.Permissive = true
		return nil
	})
}

// WithExperimentalFuncs is an option that enables experimental functions not
// in the N1 Normative specification.
func WithExperimentalFuncs() opts.CompileOption {
	return opts.Transform(func(cfg *opts.CompileConfig) error {
		cfg.Table = funcs.AddExperimentalFuncs(cfg.Table)
		return nil
	})
}
