package exprtest

import (
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

// MockExpression is a test double expression that calls
// the contained function when evaluated.
type MockExpression struct {
	Eval func(*expr.Context, system.Collection) (system.Collection, error)
}

// Evaluate calls the contained Eval function.
func (e *MockExpression) Evaluate(ctx *expr.Context, input system.Collection) (system.Collection, error) {
	return e.Eval(ctx, input)
}

// Error creates a MockExpression that returns the provided
// error when evaluated.
func Error(input error) *MockExpression {
	return &MockExpression{
		func(*expr.Context, system.Collection) (system.Collection, error) {
			return nil, input
		},
	}
}

// Return creates a MockExpression that returns the provided
// inputs when evaluated.
func Return(out ...any) *MockExpression {
	return &MockExpression{
		func(*expr.Context, system.Collection) (system.Collection, error) {
			return out, nil
		},
	}
}
