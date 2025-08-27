package fhirjson

import (
	"encoding/json"
	"io"
	"time"

	"github.com/verily-src/fhirpath-go/internal/fhir"
)

// A Decoder reads and decodes FHIR JSON values from an input stream.
type Decoder struct {
	decoder      *json.Decoder
	unmarshaller *Unmarshaller
}

// NewDecoder returns a new decoder that reads from r.
//
// The decoder introduces its own buffering and may
// read data from r beyond the JSON values requested.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		decoder:      json.NewDecoder(r),
		unmarshaller: &Unmarshaller{},
	}
}

// TimeZone sets the timezone to be used during decoding.
func (d *Decoder) TimeZone(location *time.Location) {
	d.unmarshaller.TimeZone = location
}

// Decode reads the next JSON-encoded FHIR value from its
// input and stores it in the value pointed to by out.
//
// See the documentation for [Unmarshal] for details about
// the conversion of JSON into a Go value.
func (d *Decoder) Decode(out fhir.Resource) error {
	var bytes json.RawMessage
	if err := d.decoder.Decode(&bytes); err != nil {
		return err
	}
	return d.unmarshaller.Unmarshal(bytes, out)
}

// DecodeNew reads the next JSON-encoded FHIR value from its
// input and returns it.
//
// See the documentation for [UnmarshalNew] for details about
// the conversion of JSON into a Go value.
func (d *Decoder) DecodeNew() (fhir.Resource, error) {
	var bytes json.RawMessage
	if err := d.decoder.Decode(&bytes); err != nil {
		return nil, err
	}
	return d.unmarshaller.UnmarshalNew(bytes)
}

// Buffered returns a reader of the data remaining in the Decoder's
// buffer. The reader is valid until the next call to [Decoder.Decode].
func (d *Decoder) Buffered() io.Reader {
	return d.decoder.Buffered()
}
