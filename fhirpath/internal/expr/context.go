package expr

import (
	"context"
	"time"

	"github.com/verily-src/fhirpath-go/fhirpath/resolver"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/fhirpath/terminology"
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

	// Service is an optional mechanism for providing a terminology service
	// which can be used to validate code in valueSet
	TermService terminology.Service

	// GoContext is a context from the calling main function
	GoContext context.Context
}

// Deadline wraps the Deadline() method of context.Context. More information available at https://pkg.go.dev/context
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.GoContext.Deadline()
}

// Done wraps the Done() method of context.Context. More information available at https://pkg.go.dev/context
func (c *Context) Done() <-chan struct{} {
	return c.GoContext.Done()
}

// Err wraps the Err() method of context.Context. More information available at https://pkg.go.dev/context
func (c *Context) Err() error {
	return c.GoContext.Err()
}

// Value wraps the Value() method of context.Context. More information available at https://pkg.go.dev/context
func (c *Context) Value(key any) any {
	return c.GoContext.Value(key)
}

// Clone copies this Context object to produce a new instance.
func (c *Context) Clone() *Context {
	return &Context{
		Now:               c.Now,
		ExternalConstants: c.ExternalConstants,
		LastResult:        c.LastResult,
		Resolver:          c.Resolver,
		TermService:       c.TermService,
		GoContext:         c.GoContext,
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
