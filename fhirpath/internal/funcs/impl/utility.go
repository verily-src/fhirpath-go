package impl

import (
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
