package identifier

import (
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

// Option is an abstraction for options to construct or modify Identifier elements.
type Option interface {
	update(*dtpb.Identifier)
}

// WithUse returns an Identifier Option that sets the Identifier.Use to the
// specified use.
func WithUse(use Use) Option {
	return withCallback(func(i *dtpb.Identifier) {
		i.Use = &dtpb.Identifier_UseCode{
			Value: use,
		}
	})
}

// WithExtensions return an Identifier Option that sets the Identifier.Extension
// field to the specified extensions.
func WithExtensions(ext ...*dtpb.Extension) Option {
	return withCallback(func(i *dtpb.Identifier) {
		i.Extension = ext
	})
}

// IncludeExtensions return an Identifier Option that appends the specified
// extensions to the Identifier.Extension field.
func IncludeExtensions(ext ...*dtpb.Extension) Option {
	return withCallback(func(i *dtpb.Identifier) {
		i.Extension = append(i.Extension, ext...)
	})
}

// WithType returns an Identifier Option that sets the Identifier.Type to the
// specified type.
func WithType(ty *dtpb.CodeableConcept) Option {
	return withCallback(func(i *dtpb.Identifier) {
		i.Type = ty
	})
}

// WithSystem returns an Identifier Option that sets the Identifier.System to the
// specified system.
func WithSystem(system *dtpb.Uri) Option {
	return withCallback(func(i *dtpb.Identifier) {
		i.System = system
	})
}

// WithSystemString returns an Identifier Option that sets the Identifier.System
// to the specified system string.
func WithSystemString(system string) Option {
	return WithSystem(fhir.URI(system))
}

// WithValue returns an Identifier Option that sets the Identifier.Value to the
// specified value.
func WithValue(value string) Option {
	return withCallback(func(i *dtpb.Identifier) {
		i.Value = fhir.String(value)
	})
}

// WithPeriod returns an Identifier Option that sets the Identifier.Period to the
// specified period.
func WithPeriod(period *dtpb.Period) Option {
	return withCallback(func(i *dtpb.Identifier) {
		i.Period = period
	})
}

// WithAssigner returns an Identifier Option that sets the Identifier.Assigner to the
// specified assigner reference.
func WithAssigner(assigner *dtpb.Reference) Option {
	return withCallback(func(i *dtpb.Identifier) {
		i.Assigner = assigner
	})
}

// WithID returns an Identifier Option that sets the Identifier.Id to the
// specified ID.
func WithID(id string) Option {
	return withCallback(func(i *dtpb.Identifier) {
		i.Id = fhir.String(id)
	})
}

type callbackOpt struct {
	callback func(*dtpb.Identifier)
}

func (o callbackOpt) update(i *dtpb.Identifier) {
	o.callback(i)
}

func withCallback(callback func(*dtpb.Identifier)) Option {
	return callbackOpt{callback}
}
