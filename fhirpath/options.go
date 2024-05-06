package fhirpath

import (
	"github.com/verily-src/fhirpath-go/fhirpath/compopts"
	"github.com/verily-src/fhirpath-go/fhirpath/evalopts"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/opts"
)

// CompileOption is a function type that modifies a passed in compileOption.
// Can define mutator functions of this type (see WithLimitation below)
type CompileOption = opts.CompileOption

// EvaluateOption is a function type that mutates the evalOptions type.
type EvaluateOption = opts.EvaluateOption

// WithFunction is a compile option that allows the addition of user-defined
// functions to a FHIRPath expression. Function argument must match the signature
// func(Collection, ...any) (Collection, error), or an error will be raised.
//
// Deprecated: Use compopts.Function instead.
func WithFunction(name string, fn any) CompileOption {
	return compopts.AddFunction(name, fn)
}

// WithConstant is an EvaluateOption that allows the addition of external
// constant variables. An error will be raised if the value passed in is
// neither a fhir proto or system type.
//
// Deprecated: Use evalopts.EnvVariable instead
func WithConstant(name string, value any) EvaluateOption {
	return evalopts.EnvVariable(name, value)
}
