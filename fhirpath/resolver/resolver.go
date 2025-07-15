// Package resolver defines the Resolver interface that implements Resolve FHIRPath  functionality.
package resolver

import (
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

// Resolver interface defines the Resolve() method for resolving a collection of string-based references
// (URIs, canonical URLs, or plain URLs) into a collection of FHIR resources.
type Resolver interface {
	Resolve(input []string) ([]fhir.Resource, error)
}
