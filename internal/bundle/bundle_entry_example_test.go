package bundle_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/fhir/go/fhirversion"
	"github.com/google/fhir/go/jsonformat"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	r4pb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/mattbaird/jsonpatch"
	"github.com/verily-src/fhirpath-go/internal/bundle"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"google.golang.org/protobuf/proto"
)

var resIdentity, _ = resource.NewIdentity("Patient", "123", "")

func ExamplePatchEntryFromBytes_stringPatch() {
	patch := `[{"op":"add","path":"/active","value":true}]`
	pEntry, err := bundle.PatchEntryFromBytes(resIdentity, []byte(patch))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("PatchEntry: %+v", pEntry)
}

func ExamplePatchEntryFromBytes_mapPatch() {
	patch := []map[string]interface{}{
		{
			"op":    "replace",
			"path":  "/active",
			"value": true,
		},
	}
	payload, err := json.Marshal(patch)
	if err != nil {
		log.Fatal(err)
	}
	pEntry, err := bundle.PatchEntryFromBytes(resIdentity, payload)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("PatchEntry: %+v", pEntry)
}

// ExamplePatchEntryFromBytes_diffPatch creates a patch for the diff of two
// given resources diff.
func ExamplePatchEntryFromBytes_diffPatch() {
	res := &r4pb.ContainedResource{
		OneofResource: &r4pb.ContainedResource_Patient{
			Patient: &ppb.Patient{
				Id: &dtpb.Id{
					Value: "123",
				},
				Active: &dtpb.Boolean{
					Value: false,
				},
			},
		},
	}
	newRes := proto.Clone(res).(*r4pb.ContainedResource)
	newRes.GetPatient().Active = &dtpb.Boolean{Value: true}

	m, err := jsonformat.NewMarshaller(false, "", "", fhirversion.R4)
	if err != nil {
		log.Fatal(err)
	}
	resB, err := m.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}
	newResB, err := m.Marshal(newRes)
	if err != nil {
		log.Fatal(err)
	}

	patch, err := jsonpatch.CreatePatch(resB, newResB)
	if err != nil {
		log.Fatal(err)
	}
	pPayload, err := json.Marshal(patch)
	if err != nil {
		log.Fatal(err)
	}

	pEntry, err := bundle.PatchEntryFromBytes(resIdentity, pPayload)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("PatchEntry: %+v", pEntry)
}
