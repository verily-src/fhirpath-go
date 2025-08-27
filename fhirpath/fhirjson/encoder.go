package fhirjson

import (
	"encoding/json"
	"io"

	"github.com/verily-src/fhirpath-go/internal/fhir"
)

// An Encoder writes JSON values to an output stream.
type Encoder struct {
	encoder *json.Encoder
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		encoder: json.NewEncoder(w),
	}
}

// SetEscapeHTML specifies whether problematic HTML characters
// should be escaped inside JSON quoted strings.
// The default behavior is to escape &, <, and > to \u0026, \u003c, and \u003e
// to avoid certain safety problems that can arise when embedding JSON in HTML.
//
// In non-HTML settings where the escaping interferes with the readability
// of the output, SetEscapeHTML(false) disables this behavior.
func (e *Encoder) SetEscapeHTML(on bool) {
	e.encoder.SetEscapeHTML(on)
}

// SetIndent instructs the encoder to format each subsequent encoded
// value as if indented by the package-level function Indent(dst, src, prefix, indent).
// Calling SetIndent("", "") disables indentation.
func (e *Encoder) SetIndent(prefix, indent string) {
	e.encoder.SetIndent(prefix, indent)
}

// Encode writes the JSON encoding of the FHIR resource r to the stream,
// followed by a newline character.
//
// See the documentation for [Marshal] for details about the
// conversion of FHIR values to JSON.
func (e *Encoder) Encode(r fhir.Resource) error {
	bytes, err := Marshal(r)
	if err != nil {
		return err
	}
	// Note: The underlying encoder takes care of indentation by default, even
	// for formatting raw messages -- so there is no need to use a
	// fhirjson.Marshaller here
	return e.encoder.Encode(json.RawMessage(bytes))
}
