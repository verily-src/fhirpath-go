package fhirjson

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/fhir/go/jsonformat"
	"github.com/verily-src/fhirpath-go/internal/containedresource"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"google.golang.org/protobuf/proto"
)

var (
	// defaultMarshaller is the marshaller used by the default fhirjson.Marshal.
	defaultMarshaller *jsonformat.Marshaller
)

func init() {
	marshaller, err := jsonformat.NewMarshaller(enableIndent, prefix, indent, version)
	if err != nil {
		// SAFETY:
		// This can never occur; the only error returned from NewMarshaller
		// is from an invalid fhirversion value.
		panic(fmt.Sprintf("An error occurred creating JSON marshaller: %v", err))
	}
	defaultMarshaller = marshaller
}

// Resource is a 0-cost convience wrapper type that implements the json.Marshaler
// and json.Unmarshaler API, so that the underlying proto FHIR Resource can
// properly be converted to and from JSON with the standard json API.
type Resource struct {
	fhir.Resource
}

// MarshalJSON implements json.Marshaler and returns the JSON encoding of the
// this resource.
func (r Resource) MarshalJSON() ([]byte, error) {
	return Marshal(r.Resource)
}

// UnmarshalJSON implements json.Unmarshaler and decode the json data, storing
// it in this resource on success.
func (r *Resource) UnmarshalJSON(data []byte) error {
	var err error
	r.Resource, err = UnmarshalNew(data)
	return err
}

var _ json.Marshaler = (*Resource)(nil)
var _ json.Unmarshaler = (*Resource)(nil)

// Marshal returns serialized JSON object of a FHIR Resource protobuf message.
func Marshal(resource fhir.Resource) ([]byte, error) {
	return marshal(defaultMarshaller, resource)
}

// MarshalIndent is like [Marshal] but applies [Indent] to format the output.
// Each JSON element in the output will begin on a new line beginning with prefix
// followed by one or more copies of indent according to the indentation nesting.
func MarshalIndent(resource fhir.Resource, prefix, indent string) ([]byte, error) {
	marshaller := Marshaller{Prefix: prefix, Indent: indent}
	return marshaller.Marshal(resource)
}

// Marshaller is a configurable JSON FHIR format marshaler.
type Marshaller struct {
	// EnableIndent determines whether indents should be used during formatting
	// the JSON. If this is set, the "Indent" field will be used as the indent
	// text.
	//
	// Deprecated: This field is now intuited from the "Indent" field setting
	EnableIndent bool

	// Prefix determines a prefix that will be included on each line of the formatted
	// output. This will only appear if EnableIndent is set to true.
	Prefix string

	// Indent determines the string that will be used for formatted indenting if
	// EnableIndent is set.
	Indent string
}

var (
	errMarshal = errors.New("fhirjson.Marshal")

	// ErrNilMarshalResource is an error raised for bad resource inputs for marshalling.
	ErrNilMarshalResource = fmt.Errorf("%w: nil resource", errMarshal)
)

// Marshal returns serialized JSON object of a FHIR Resource protobuf message.
// This returns ErrInvalidMarshal if resource is nil.
func (o *Marshaller) Marshal(resource fhir.Resource) ([]byte, error) {
	marshaller, err := jsonformat.NewMarshaller(o.EnableIndent || o.Indent != "", o.Prefix, o.Indent, version)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMarshal, err)
	}
	return marshal(marshaller, resource)
}

func marshal(marshaller *jsonformat.Marshaller, resource fhir.Resource) ([]byte, error) {
	if resource == nil {
		return nil, ErrNilMarshalResource
	}
	// Somehow the default marshaller mutates the resource, so we need to clone it
	// first.
	resource = proto.Clone(resource).(fhir.Resource)
	cr := containedresource.Wrap(resource)
	return marshaller.Marshal(cr)
}
