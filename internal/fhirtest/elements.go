package fhirtest

import (
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/protofields"
)

var (
	// Elements is a map of all element-names to an instance of that element type.
	//
	// The elements in this map are not guaranteed to contain any specific value;
	// it is only guaranteed to contain a non-nil instance of a concrete element
	// of the associated name.
	Elements map[string]fhir.Element

	// BackboneElements is a map of all backbone element-names to an instance of
	// that element type.
	//
	// The elements in this map are not guaranteed to contain any specific value;
	// it is only guaranteed to contain a non-nil instance of a concrete backbone
	// element of the associated name.
	BackboneElements map[string]fhir.BackboneElement
)

func init() {
	Elements = map[string]fhir.Element{}
	BackboneElements = map[string]fhir.BackboneElement{}

	for _, msg := range protofields.Elements {
		element, ok := msg.New().(fhir.Element)

		// The proto definition of the XHtml type is missing `GetValue()`, and thus
		// fails this check. This is added to avoid errors here that are otherwise
		// correct for all other cases.
		if !ok {
			continue
		}

		name := protofields.DescriptorName(element)

		Elements[name] = element
		if val, ok := any(element).(fhir.BackboneElement); ok {
			BackboneElements[name] = val
		}
	}
}
