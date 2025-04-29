package fhirpath

import (
	"errors"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/fhirpath/evalopts"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/compile"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/opts"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/parser"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/slices"
)

var (
	ErrInvalidField     = expr.ErrInvalidField
	ErrUnsupportedType  = evalopts.ErrUnsupportedType
	ErrExistingConstant = evalopts.ErrExistingConstant
)

// Resource is a FHIR resource. This is an alias for the
// fhir.Resource type, which is the base type for all FHIR resources.
type Resource = fhir.Resource

// Expression is the FHIRPath expression that will be compiled from a FHIRPath string
type Expression struct {
	expression expr.Expression
	path       string
}

// Compile parses and compiles the FHIRPath expression down to a single
// Expression object.
//
// If there are any syntax or semantic errors, this will return an
// error indicating the compilation failure reason.
func Compile(expr string, options ...CompileOption) (*Expression, error) {
	config, err := compile.PopulateConfig(options...)
	if err != nil {
		return nil, err
	}

	tree, err := compile.Tree(expr)
	if err != nil {
		return nil, err
	}

	visitor := &parser.FHIRPathVisitor{
		Functions:  config.Table,
		Permissive: config.Permissive,
	}
	vr, ok := visitor.Visit(tree).(*parser.VisitResult)
	if !ok {
		return nil, errors.New("input expression currently unsupported")
	}

	if vr.Error != nil {
		return nil, vr.Error
	}
	return &Expression{
		expression: vr.Result,
		path:       expr,
	}, nil
}

// String returns the string representation of this FHIRPath expression.
// This is just the input that initially produced the FHIRPath value.
func (e *Expression) String() string {
	return e.path
}

// MustCompile compiles the FHIRpath expression input, and returns the
// compiled expression. If any compilation error occurs, this function
// will panic.
func MustCompile(expr string, opts ...CompileOption) *Expression {
	result, err := Compile(expr, opts...)
	if err != nil {
		panic(err)
	}
	return result
}

// Evaluate the expression, returning either a collection of elements, or error
func (e *Expression) Evaluate(input []Resource, options ...EvaluateOption) (system.Collection, error) {
	config := &opts.EvaluateConfig{
		Context: expr.InitializeContext(slices.MustConvert[any](input)),
	}
	config, err := opts.ApplyOptions(config, options...)
	if err != nil {
		return nil, err
	}

	collection := slices.MustConvert[any](input)
	return e.expression.Evaluate(config.Context, collection)
}

// EvaluateAsString evaluates the expression, returning a string or error
func (e *Expression) EvaluateAsString(input []Resource, options ...EvaluateOption) (string, error) {
	got, err := e.Evaluate(input, options...)
	if err != nil {
		return "", err
	}
	return got.ToString()
}

// EvaluateAsBool evaluates the expression, returning either a boolean or error
func (e *Expression) EvaluateAsBool(input []Resource, options ...EvaluateOption) (bool, error) {
	got, err := e.Evaluate(input, options...)
	if err != nil {
		return false, err
	}
	return got.ToBool()
}

// EvaluateAsInt32 evaluates the expression, returning either an int32 or error
func (e *Expression) EvaluateAsInt32(input []Resource, options ...EvaluateOption) (int32, error) {
	got, err := e.Evaluate(input, options...)
	if err != nil {
		return 0, err
	}
	return got.ToInt32()
}

// EvaluateAsCanonical evaluates the expression, returning a FHIR canonical or error
func (e *Expression) EvaluateAsCanonical(input []Resource, options ...EvaluateOption) (*dtpb.Canonical, error) {
	got, err := e.Evaluate(input, options...)
	if err != nil {
		return nil, err
	}
	return got.ToCanonical()
}
