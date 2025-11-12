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

type simpleResolver map[string][]fhirpath.Resource

// SimpleResolverOption is a function that updates a SimpleResolver.
type SimpleResolverOption interface {
	update(resolver simpleResolver)
}

type simpleResolverOption func(simpleResolver)

func (o simpleResolverOption) update(resolver simpleResolver) {
	o(resolver)
}

// Entry adds entries at the specified key to the SimpleResolver.
func Entry(key string, resources ...fhirpath.Resource) SimpleResolverOption {
	return simpleResolverOption(func(r simpleResolver) {
		r[key] = resources
	})
}

// Resolve returns the resources at the specified keys if they exist.
func (r simpleResolver) Resolve(input []string) ([]fhir.Resource, error) {
	var result []fhir.Resource
	for _, ref := range input {
		if v, ok := r[ref]; ok {
			result = append(result, v...)
		}
	}
	return result, nil
}

// NewSimpleResolver returns a resolver.Resolver that uses the reference as a key to lookup for the
// resource in a map.
func NewSimpleResolver(opts ...SimpleResolverOption) resolver.Resolver {
	r := simpleResolver{}
	for _, opt := range opts {
		opt.update(r)
	}
	return r
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
