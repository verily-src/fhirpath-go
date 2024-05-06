package impl

import (
	"fmt"

	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

// Extension is syntactic sugar over `extension.where(url = ...)`, and is
// specific to the R4 extensions for FHIRPath (as oppose to being part of the
// N1 normative spec).
//
// For more details, see https://hl7.org/fhir/R4/fhirpath.html#functions
func Extension(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("%w: expected 1 argument", ErrWrongArity)
	}
	arg, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	}
	str, err := arg.ToString()
	if err != nil {
		return nil, err
	}

	var result system.Collection
	for _, entry := range input {
		entry, ok := entry.(fhir.Extendable)
		if !ok {
			continue
		}
		for _, ext := range entry.GetExtension() {
			if url := ext.GetUrl(); url != nil && url.Value == str {
				result = append(result, ext)
			}
		}
	}
	return result, nil
}
