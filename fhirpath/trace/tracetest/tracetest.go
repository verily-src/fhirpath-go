package tracetest

import (
	"context"

	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

// Recorder is a simple implementation of [trace.Tracer] that records
// trace events in memory for later inspection.
type Recorder struct {
	entries map[string]system.Collection
}

// Trace records a trace event with the given name and value.
func (r *Recorder) Trace(_ context.Context, name string, value system.Collection) {
	if r.entries == nil {
		r.entries = make(map[string]system.Collection)
	}
	r.entries[name] = value
}

// Collection returns the recorded collection for the given trace name.
func (r *Recorder) Collection(name string) system.Collection {
	return r.entries[name]
}
