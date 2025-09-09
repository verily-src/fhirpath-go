// Package resolvertest provides test utilities for the resolver package.
package resolvertest

import (
	"github.com/verily-src/fhirpath-go/fhirpath/resolver"

	"github.com/verily-src/fhirpath-go/fhirpath"
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

type resolverFunc func(input []string) ([]fhir.Resource, error)

// Resolve calls the underlying resolverFunc, fulfilling the resolver.Resolver interface.
func (rf resolverFunc) Resolve(input []string) ([]fhir.Resource, error) {
	return rf(input)
}

// HappyResolver returns a resolver.Resolver that always returns the provided resources and no error.
func HappyResolver(resources ...fhirpath.Resource) resolver.Resolver {
	return resolverFunc(func(input []string) ([]fhir.Resource, error) {
		return resources, nil
	})
}

// ErroringResolver returns a resolver.Resolver that always returns the given error and no resources.
func ErroringResolver(err error) resolver.Resolver {
	return resolverFunc(func(input []string) ([]fhir.Resource, error) {
		return nil, err
	})
}
