package impl

import (
	"fmt"

	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

// TimeOfDay returns the current time as a system.Time object.
func TimeOfDay(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	timeString := ctx.Now.Format("15:04:05.000")
	return system.Collection{system.MustParseTime(timeString)}, nil
}

// Today returns the current date as a system.Date object.
func Today(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	dateString := ctx.Now.Format("2006-01-02")
	return system.Collection{system.MustParseDate(dateString)}, nil
}

// Now returns the current time as a system.DateTime object.
func Now(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	dateTimeString := ctx.Now.Format("2006-01-02T15:04:05.000Z07:00")
	return system.Collection{system.MustParseDateTime(dateTimeString)}, nil
}

// Trace implements the fhirpath "trace" function, which enables custom diagnostic
// logging during FHIRPath evaluation.
//
// For full information, see:
// https://hl7.org/fhirpath/N1/#tracename-string-projection-expression-collection
func Trace(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	if length := len(args); length < 1 {
		return nil, fmt.Errorf("%w: received %d arguments, expected minimum 1", ErrWrongArity, length)
	} else if length > 2 {
		return nil, fmt.Errorf("%w: received %d arguments, expected maximum 2", ErrWrongArity, length)
	}

	if ctx.Tracer == nil {
		return input, nil
	}
	nameArg, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	}
	name, err := nameArg.ToString()
	if err != nil {
		return nil, err
	}

	result := input
	if len(args) == 2 {
		result = system.Collection{}
		projection := args[1]
		for _, entry := range input {
			r, err := projection.Evaluate(ctx, system.Collection{entry})
			if err != nil {
				return nil, err
			}
			result = append(result, r...)
		}
	}
	ctx.Tracer.Trace(ctx.GoContext, name, result)

	return input, nil
}
