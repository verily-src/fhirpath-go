package extension_test

import (
	"fmt"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/task_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/element/extension"
)

func ExampleOverwrite() {
	const urlBase = "https://verily-src.github.io/vhp-hds-vvs-fhir-ig/StructureDefinitions"
	task := &task_go_proto.Task{
		Extension: []*dtpb.Extension{
			extension.New(fmt.Sprintf("%v/%v", urlBase, "my-int"), fhir.Integer(42)),
		},
	}

	extension.Overwrite(task,
		extension.New(fmt.Sprintf("%v/%v", urlBase, "my-string"), fhir.String("hello world")),
		extension.New(fmt.Sprintf("%v/%v", urlBase, "my-bool"), fhir.Boolean(true)),
	)
	fmt.Printf("%v extensions in Task!", len(task.GetExtension()))
	// Output: 2 extensions in Task!
}

func ExampleAppendInto() {
	const urlBase = "http://example.com/StructureDefinitions"
	task := &task_go_proto.Task{
		Extension: []*dtpb.Extension{
			extension.New(fmt.Sprintf("%v/%v", urlBase, "my-int"), fhir.Integer(42)),
		},
	}

	extension.AppendInto(task,
		extension.New(fmt.Sprintf("%v/%v", urlBase, "my-string"), fhir.String("hello world")),
		extension.New(fmt.Sprintf("%v/%v", urlBase, "my-bool"), fhir.Boolean(true)),
	)
	fmt.Printf("%v extensions in Task!", len(task.GetExtension()))
	// Output: 3 extensions in Task!
}
