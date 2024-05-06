/*
Package resourceopt is an internal package that provides helper utilities
for forming resource-options in resource packages.
*/
package resourceopt

import (
	"github.com/verily-src/fhirpath-go/internal/slices"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/protofields"
	"google.golang.org/protobuf/proto"
)

// Option is the definition of a resource Option used for creating and updating
// FHIR Resources.
type Option interface {
	update(fhir.Resource)
}

// ApplyOptions applies the specified options to the input resource.
//
// This function is defined here due to the Option interface providing an
// unexported field. This is needed so that the other packages using this can
// accumulate the options without having access to the unexported call.
func ApplyOptions[T fhir.Resource, O Option](r T, opts ...O) T {
	for _, opt := range opts {
		opt.update(r)
	}
	return r
}

// WithProtoField is a resource Option that sets the specified 'field' in the proto
// to the values. If values is empty, the field is cleared. If values is
// not 1, and a field is not repeated, this functon will panic.
//
// Note: This is an internal function intended to be used to form generic
// resource options that will work with all FHIR resources.
func WithProtoField[T proto.Message](fieldName string, values ...T) Option {
	// SAFETY:
	//   MustConvert cannot fail here, since the 'T' constraint above ensures that
	//   all inputs will be valid proto.Message types.
	return withProtoFieldImpl(fieldName, slices.MustConvert[proto.Message](values)...)
}

func withProtoFieldImpl(fieldName string, values ...proto.Message) Option {
	return WithCallback(func(r fhir.Resource) {
		protofields.Overwrite(r, fieldName, values...)
	})
}

// IncludeProtoField is a resource Option that appends the specified entries to
// the given 'field' in the proto. This function will panic if the given field
// is not a repeated field in the proto.
//
// Note: This is an internal function intended to be used to form generic
// resource options that will work with all FHIR resources.
func IncludeProtoField[T proto.Message](fieldName string, values ...T) Option {
	// SAFETY:
	//   MustConvert cannot fail here, since the 'T' constraint above ensures that
	//   all inputs will be valid proto.Message types.
	return includeProtoFieldImpl(fieldName, slices.MustConvert[proto.Message](values)...)
}

func includeProtoFieldImpl(fieldName string, values ...proto.Message) Option {
	return WithCallback(func(r fhir.Resource) {
		protofields.AppendList(r, fieldName, values...)
	})
}

// WithCallback returns a resource Option that simply passes the resource being
// created back into the specified callback. This exists to be built into
// larger, more strongly-typed options.
func WithCallback[T fhir.Resource](callback func(T)) Option {
	return &callbackOpt[T]{callback}
}

type callbackOpt[T fhir.Resource] struct {
	callback func(T)
}

func (o *callbackOpt[T]) update(r fhir.Resource) {
	o.callback(r.(T))
}
