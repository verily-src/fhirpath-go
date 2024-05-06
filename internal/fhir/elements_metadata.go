package fhir

import dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"

// Metadata Types:
//
// The section below defines types from the "MetaDataTypes" heading in
// http://hl7.org/fhir/R4/datatypes.html#open

// ContactDetail creates an R4 FHIR ContactDetail element from a string value
// and the specified contact-points.
//
// See: http://hl7.org/fhir/R4/metadatatypes.html#ContactDetail
func ContactDetail(name string, telecom ...*dtpb.ContactPoint) *dtpb.ContactDetail {
	return &dtpb.ContactDetail{
		Name:    String(name),
		Telecom: telecom,
	}
}
