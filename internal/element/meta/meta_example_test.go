package meta_test

import (
	"fmt"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/element/canonical"
	"github.com/verily-src/fhirpath-go/internal/element/meta"
)

func ExampleUpdate() {
	m := &dtpb.Meta{}

	meta.Update(m,
		meta.WithTags(fhir.Coding("urn:oid:verily/sample-tag-system", "sample-tag-value")),
		meta.WithProfiles(canonical.New("urn:oid:verily/sample-profile")),
	)

	fmt.Printf("meta.profile: %q\n", m.Profile[0].Value)
	fmt.Printf("meta.tag: {%q, %q}", m.Tag[0].System.Value, m.Tag[0].Code.Value)
	// Output:
	// meta.profile: "urn:oid:verily/sample-profile"
	// meta.tag: {"urn:oid:verily/sample-tag-system", "sample-tag-value"}
}
