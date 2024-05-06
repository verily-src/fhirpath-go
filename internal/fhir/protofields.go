package fhir

import (
	"github.com/verily-src/fhirpath-go/internal/protofields"
)

// UnwrapValueX obtains the underlying Message for oneof ValueX
// elements, which use a nested Choice field. Returns nil if the input message
// doesn't have a Choice field, or if the Oneof descriptor is unpopulated.
// See wrapped implementation for more information.
func UnwrapValueX(element Base) Base {
	return protofields.UnwrapOneofField(element, "choice")
}
