package fhirjson

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/google/fhir/go/jsonformat"
	bcrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"github.com/verily-src/fhirpath-go/internal/containedresource"
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

var (
	errUnmarshal = errors.New("fhirjson.Unmarshal")

	// ErrUnmarshalEncoding is an error that is raised from bad json encoding.
	ErrUnmarshalEncoding = fmt.Errorf("%w: bad json encoding", errUnmarshal)

	// ErrNilUnmarshalOutput is an error raised if the output is `nil` for Unmarshal calls.
	ErrNilUnmarshalOutput = fmt.Errorf("%w: nil output", errUnmarshal)

	// ErrWrongUnmarshalType is an error raised if the specified output type is
	// not the resource that was contained.
	ErrWrongUnmarshalType = fmt.Errorf("%w: incorrect output type", errUnmarshal)
)

var (
	// defaultUnmarshaller is the unmarshaller used by the default fhirjson.Unmarshal.
	defaultUnmarshaller *jsonformat.Unmarshaller
)

func init() {
	tz := time.Local.String()
	unmarshaller, err := jsonformat.NewUnmarshaller(tz, version)
	if err != nil {
		// SAFETY:
		// This can never occur; the errors returned from NewUnmarshaller are due to
		// invalid version argument (for which we are pinned to a valid version), or
		// due to a bad timezone string (which we are using the system local timezone).
		panic(fmt.Sprintf("An error occurred creating JSON unmarshaller: %v", err))
	}
	defaultUnmarshaller = unmarshaller
}

// Unmarshal parses the JSON-encoded FHIR resource and stores the result in the
// value pointed to by out.
//
// This function may return the following errors:
//   - ErrNilUnmarshalOutput if `out` is nil
//   - ErrWrongUnmarshalType if `out` is not the correct resource type
//   - ErrUnmarshalEncoding if there is an encoding error with `data`
func Unmarshal(data []byte, out fhir.Resource) error {
	return unmarshal(defaultUnmarshaller, data, out)
}

// UnmarshalOut parses the JSON-encoded FHIR resource and returns the parsed
// resource.
//
// This function may only return an ErrUnmarshalEncoding if there is an encoding
// error with `data`.
//
// Deprecated: Use UnmarshalNew
func UnmarshalOut(data []byte) (fhir.Resource, error) {
	return UnmarshalNew(data)
}

// UnmarshalNew parses the JSON-encoded FHIR resource and returns the parsed
// resource.
//
// This function may only return an ErrUnmarshalEncoding if there is an encoding
// error with `data`.
func UnmarshalNew(data []byte) (fhir.Resource, error) {
	return unmarshalOut(defaultUnmarshaller, data)
}

// Unmarshaller is a configurable JSON FHIR format unmarshaler.
type Unmarshaller struct {
	// TimeZone determines which timezone time values are deserialized into for
	// time-based values. If nil, this defaults to the time.Local.
	TimeZone *time.Location
}

// Unmarshal parses the JSON-encoded FHIR resource and stores the result in the
// value pointed to by out.
//
// This function may return the following errors:
//   - ErrNilUnmarshalOutput if `out` is nil
//   - ErrWrongUnmarshalType if `out` is not the correct resource type
//   - ErrUnmarshalEncoding if there is an encoding error with `data`
func (o *Unmarshaller) Unmarshal(data []byte, out fhir.Resource) error {
	unmarshaller, err := o.newUnmarshaller()
	if err != nil {
		return err
	}
	return unmarshal(unmarshaller, data, out)
}

// UnmarshalOut parses the JSON-encoded FHIR resource and returns the parsed
// resource.
//
// This function may only return an ErrUnmarshalEncoding if there is an encoding
// error with `data`.
//
// Deprecated: Use UnmarshalNew
func (o *Unmarshaller) UnmarshalOut(data []byte) (fhir.Resource, error) {
	return o.UnmarshalNew(data)
}

// UnmarshalNew parses the JSON-encoded FHIR resource and returns the parsed
// resource.
//
// This function may only return an ErrUnmarshalEncoding if there is an encoding
// error with `data`.
func (o *Unmarshaller) UnmarshalNew(data []byte) (fhir.Resource, error) {
	unmarshaller, err := o.newUnmarshaller()
	if err != nil {
		return nil, err
	}
	return unmarshalOut(unmarshaller, data)
}

// newUnmarshaller creates a new unmarshaller from the configurations.
func (o *Unmarshaller) newUnmarshaller() (*jsonformat.Unmarshaller, error) {
	unmarshaller, err := jsonformat.NewUnmarshaller(o.timeZoneOrDefault(), version)
	if err != nil {
		// This error should never feasibly happen, since time.Location objects
		// should always have valid string names, and the version is pinned.
		return nil, fmt.Errorf("%w: unable to create unmarshaller: %v", errUnmarshal, err)
	}
	return unmarshaller, nil
}

// timeZoneOrDefault returns the timezone string for the unmarshal options. If one
// is not set, it will provide the Local timezone's string.
func (o *Unmarshaller) timeZoneOrDefault() string {
	if o.TimeZone == nil {
		return time.Local.String()
	}
	return o.TimeZone.String()
}

// unmarshalOut performs unmarshalling out to a fhir.Resource object in the
// return type, using the specified unmarshaller.
func unmarshalOut(unmarshaller *jsonformat.Unmarshaller, data []byte) (fhir.Resource, error) {
	message, err := unmarshaller.Unmarshal(data)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnmarshalEncoding, err)
	}

	cr, ok := message.(*bcrpb.ContainedResource)
	if !ok {
		// This should never happen in practice, since the jsonformat package is
		// supposed to be in terms of ContainedResource objects.
		return nil, fmt.Errorf("fhirjson.Unmarshal: unexpected type returned from unmarshal: %T", message)
	}

	resource := containedresource.Unwrap(cr)

	// SAFETY:
	// ContainedResource can only hold fhir.Resource objects; this cast will always
	// be safe for all possible values that are unwrapped.
	return resource.(fhir.Resource), nil
}

// unmarshal will parse json input and write the result back into the output
// fhir.Resource.
func unmarshal(unmarshaller *jsonformat.Unmarshaller, data []byte, out fhir.Resource) error {
	outrv := reflect.ValueOf(out)
	if outrv.Kind() != reflect.Pointer || outrv.IsNil() {
		return ErrNilUnmarshalOutput
	}

	resource, err := unmarshalOut(unmarshaller, data)
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(resource)
	if rv.Type() != outrv.Type() {
		outName := out.ProtoReflect().Descriptor().Name()
		jsonName := resource.ProtoReflect().Descriptor().Name()
		return fmt.Errorf("%w; out is of type *%v, but json contained *%v", ErrWrongUnmarshalType, outName, jsonName)
	}

	outrv.Elem().Set(rv.Elem())
	return nil
}
