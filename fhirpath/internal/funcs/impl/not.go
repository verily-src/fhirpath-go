package impl

import (
	"fmt"

	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

// Not returns the boolean inverse of the singleton input collection.
func Not(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	if length := len(args); length != 0 {
		return nil, fmt.Errorf("%w, received %v arguments, expected 0", ErrWrongArity, length)
	}

	boolean, err := input.ToSingletonBoolean()
	if err != nil {
		return nil, err
	}
	if len(boolean) == 0 {
		return system.Collection{}, nil
	}
	result := system.Boolean(!boolean[0])
	return system.Collection{result}, nil
}
