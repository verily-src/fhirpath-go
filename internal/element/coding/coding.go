package coding

import (
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
)

// FindBySystem searches a slice of codings for the first coding that
// contains the specified system.
func FindBySystem(codings []*dtpb.Coding, system string) *dtpb.Coding {
	for _, coding := range codings {
		codingSystem := coding.GetSystem()
		if codingSystem == nil {
			continue
		}
		if codingSystem.GetValue() == system {
			return coding
		}
	}
	return nil
}
