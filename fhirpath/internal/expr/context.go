package expr

import (
	"time"

	"github.com/verily-src/fhirpath-go/fhirpath/resolver"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

// Context holds the global time and external constant
// variable map, to enable deterministic evaluation.
type Context struct {
	Now               time.Time
	ExternalConstants map[string]any

	// LastResult is required for implementing most FHIRPatch operations, since
	// a reference to the node before the one being (inserted, replaced, moved) is
	// necessary in order to alter the containing object.
	LastResult system.Collection

	// BeforeLastResult is necessary for implementing FHIRPatch delete due to an
	// edge-case, where deleting a specific element from a list requires a pointer
	// to the container that holds the list. In a path like `Patient.name.given[0]`,
	// the 'LastResult' will be the unwrapped list from 'given', but we need the
	// 'name' element that contains the 'given' list in order to alter the list.
	BeforeLastResult system.Collection

	// Resolver is an optional mechanism for resolving FHIR Resources that
	// is used in the 'resolve()' FHIRPath function.
	Resolver resolver.Resolver
}

// Clone copies this Context object to produce a new instance.
func (c *Context) Clone() *Context {
	return &Context{
		Now:               c.Now,
		ExternalConstants: c.ExternalConstants,
		LastResult:        c.LastResult,
		Resolver:          c.Resolver,
	}
}

// InitializeContext returns a base context, initialized with current time and initial
// constant variables set.
func InitializeContext(input system.Collection) *Context {
	return &Context{
		Now: time.Now().Local().UTC(),
		ExternalConstants: map[string]any{
			"context": input,
			"ucum":    system.String("http://unitsofmeasure.org"),
		},
	}
}
