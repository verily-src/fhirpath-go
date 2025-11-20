/*
Package trace provides utilities for customizing the FHIRPath "trace" function.

This enables callers to define custom mechanisms for logging or recording
the evaluation steps of FHIRPath expressions, which can be useful for debugging
or analysis purposes.
*/
package trace

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/verily-src/fhirpath-go/fhirpath/fhirjson"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

// Tracer is an abstraction over the FHIRPath "trace" functionality.
//
// Implementations of Tracer objects are provided the name of the trace (if
// specified in the function call), along with the value(s) that exist at that
// point in the evaluation.
//
// Tracer operations are not allowed to fail or modify the evaluation state in
// any way.
type Tracer interface {
	// Trace is called to record a trace event during FHIRPath evaluation.
	//
	// The name parameter is the optional name provided in the trace function
	// call; it may be an empty string if no name was specified.
	Trace(ctx context.Context, name string, value system.Collection)
}

// TracerFunc is an adapter to allow the use of ordinary functions as Tracer.
type TracerFunc func(ctx context.Context, name string, value system.Collection)

// Trace calls f(ctx, name, value).
func (f TracerFunc) Trace(ctx context.Context, name string, value system.Collection) {
	f(ctx, name, value)
}

// NoopTracer returns a Tracer implementation that performs no operations.
func NoopTracer() Tracer {
	return TracerFunc(func(context.Context, string, system.Collection) {
		// No-op
	})
}

// WriterTracer is a Tracer implementation that writes trace information
// to the provided [io.Writer].
//
// Output is logged in non-minified JSON format, with one JSON object per line.
// If a name was provided to the Trace functionality, the start of the collection
// of entries is prefixed with the name followed by a colon.
type WriterTracer struct {
	Out io.Writer

	m sync.Mutex
}

// Trace writes the trace information to the configured output writer.
func (wt *WriterTracer) Trace(_ context.Context, name string, value system.Collection) {
	wt.m.Lock()
	defer wt.m.Unlock()
	var out io.Writer = os.Stdout
	if wt.Out != nil {
		out = wt.Out
	}

	if name != "" {
		_, _ = fmt.Fprintln(out, name+":")
	}
	encoder := json.NewEncoder(out)
	fhirEncoder := fhirjson.NewEncoder(out)
	fhirEncoder.SetIndent("", "  ")
	encoder.SetIndent("", "  ")
	for _, entry := range value {
		switch v := entry.(type) {
		case fhir.Resource:
			_ = fhirEncoder.Encode(v)
		case fmt.Stringer:
			_ = encoder.Encode(v.String())
		default:
			_ = encoder.Encode(v)
		}
	}
}
