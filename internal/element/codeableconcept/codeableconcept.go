package codeableconcept

import (
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
)

// FindBySystem searches the slice of Codings within a CodeableConcept for the
// first Coding that contains the given system.
func FindCodingBySystem(codeableconcept *dtpb.CodeableConcept, system string) *dtpb.Coding {
	for _, coding := range codeableconcept.GetCoding() {
		codingSystem := coding.GetSystem()
		if codingSystem.GetValue() == system {
			return coding
		}
	}
	return nil
}
