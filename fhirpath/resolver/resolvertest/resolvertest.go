// Package resolvertest provides test utilities for the resolver package.
package resolvertest

import (
	"github.com/verily-src/fhirpath-go/fhirpath/resolver"

	"github.com/verily-src/fhirpath-go/fhirpath"
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

type resolverFunc func(input []string) ([]fhir.Resource, error)

func (rf resolverFunc) Resolve(input []string) ([]fhir.Resource, error) {
	return rf(input)
}

func HappyResolver(resources ...fhirpath.Resource) resolver.Resolver {
	return resolverFunc(func(input []string) ([]fhir.Resource, error) {
		return resources, nil
	})
}

func ErroringResolver(err error) resolver.Resolver {
	return resolverFunc(func(input []string) ([]fhir.Resource, error) {
		return nil, err
	})
}
