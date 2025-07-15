// Package impl provides implementations of FHIRPath functions.
package impl

import (
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

// Combine merges input and other collections into a single collection without eliminating duplicate values.
// Combining an empty collection with a non-empty collection will return the non-empty collection.
// There is no expectation of order in the resulting collection.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#combineother-collection-collection
func Combine(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Create a new collection with the size of the input collection
	var result system.Collection

	// If the input collection is nil, we don't need to append it
	if input != nil {
		result = append(result, input...)
	}

	// Iterate over the args and append each collection to the result
	for _, arg := range args {
		coll, err := arg.Evaluate(ctx, input)
		if err != nil {
			return nil, err
		}
		if coll != nil {
			result = append(result, coll...)
		}
	}

	return result, nil
}
