package fhir

import dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"

// Special Types:
//
// The section below defines types from the "Special Types" heading in
// http://hl7.org/fhir/R4/datatypes.html#open

// Narrative creates a R4 FHIR Narrative element from a string value.
//
// See: http://hl7.org/fhir/R4/narrative.html
func Narrative(value string) *dtpb.Narrative {
	return &dtpb.Narrative{
		Div: XHTML(value),
	}
}

// XHTML creates an R4 FHIR XHTML element from a string value.
//
// See: http://hl7.org/fhir/R4/narrative.html#xhtml
func XHTML(value string) *dtpb.Xhtml {
	return &dtpb.Xhtml{
		Value: value,
	}
}
