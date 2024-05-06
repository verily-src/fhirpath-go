/*
Package fhirpathtest provides an easy way to generate test-doubles within
FHIRPath.
*/
package fhirpathtest

import (
	"github.com/verily-src/fhirpath-go/fhirpath"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

// Return creates a FHIRPath expression that will always return the given
// values.
func Return(args ...any) *fhirpath.Expression {
	return ReturnCollection(system.Collection(args))
}

// ReturnCollection creates a FHIRPath expression that will always return the
// given input collection.
func ReturnCollection(collection system.Collection) *fhirpath.Expression {
	return fhirpath.MustCompile("return()",
		fhirpath.WithFunction("return", func(system.Collection) (system.Collection, error) {
			return collection, nil
		}),
	)
}

// Error creates a FHIRPath expression that will always return the specified error.
func Error(err error) *fhirpath.Expression {
	return fhirpath.MustCompile("return()",
		fhirpath.WithFunction("return", func(system.Collection) (system.Collection, error) {
			return nil, err
		}),
	)
}

var (
	// Empty is a FHIRPath expression that returns an empty collection when
	// evaluated.
	Empty = Return(system.Collection{})

	// True is a FHIRPath expression that returns a collection containing a single
	// system boolean of 'true'. This is useful for testing expected boolean
	// logic in paths.
	True = Return(system.Collection{system.Boolean(true)})

	// False is a FHIRPath expression that returns a collection containing a single
	// system boolean of 'false'. This is useful for testing expected boolean
	// logic in paths.
	False = Return(system.Collection{system.Boolean(false)})
)
