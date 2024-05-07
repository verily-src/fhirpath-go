package identifier_test

import (
	"fmt"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/element/identifier"
)

func ExampleQueryIdentifier() {
	ident := &dtpb.Identifier{
		System: &dtpb.Uri{Value: "http://fake.com"},
		Value:  &dtpb.String{Value: "b0459744-b74b-441a-aee4-9dd97c80c642"},
	}

	search := identifier.QueryIdentifier(ident)
	fmt.Printf("identifier:exact=%s", search)
	// Output: identifier:exact=http://fake.com|b0459744-b74b-441a-aee4-9dd97c80c642
}

func ExampleQueryIdentifier_escape() {
	ident := &dtpb.Identifier{
		System: &dtpb.Uri{Value: "http://fake.com"},
		Value:  &dtpb.String{Value: "foo,bar|baz"},
	}

	search := identifier.QueryIdentifier(ident)
	fmt.Printf("identifier:exact=%s", search)
	// Output: identifier:exact=http://fake.com|foo\,bar\|baz
}

func ExampleQueryString() {
	search := identifier.QueryString("http://fake.com", "1234")
	fmt.Printf("identifier:exact=%s", search)
	// Output: identifier:exact=http://fake.com|1234
}
func ExampleQueryString_escape() {
	search := identifier.QueryString("http://fake.com", `$foo|bar\baz`)
	fmt.Printf("identifier:exact=%s", search)
	// Output: identifier:exact=http://fake.com|\$foo\|bar\\baz
}

func ExampleGenerateIfNoneExist() {
	id := &dtpb.Identifier{
		System: &dtpb.Uri{Value: "http://fake.com"},
		Value:  &dtpb.String{Value: "9efbf82d-7a58-4d14-bec1-63f8fda148a8"},
	}

	header := identifier.GenerateIfNoneExist(id)
	fmt.Printf("If-None-Exist: %v", header)
	// Output: If-None-Exist: identifier=http%3A%2F%2Ffake.com%7C9efbf82d-7a58-4d14-bec1-63f8fda148a8
}
