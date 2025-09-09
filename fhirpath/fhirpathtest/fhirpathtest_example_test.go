package fhirpathtest_test

import (
	"errors"
	"fmt"

	"github.com/verily-src/fhirpath-go/fhirpath"
	"github.com/verily-src/fhirpath-go/fhirpath/fhirpathtest"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

func ExampleError() {
	want := errors.New("example error")
	expr := fhirpathtest.Error(want)

	_, err := expr.Evaluate([]fhirpath.Resource{})
	if errors.Is(err, want) {
		fmt.Printf("err = '%v'", want)
	}

	// Output: err = 'example error'
}

func ExampleReturn() {
	want := system.Boolean(true)
	expr := fhirpathtest.Return(want)

	got, err := expr.Evaluate([]fhirpath.Resource{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("got = %v", bool(got[0].(system.Boolean)))
	// Output: got = true
}
func ExampleReturnCollection() {
	want := system.Collection{system.Boolean(true)}
	expr := fhirpathtest.ReturnCollection(want)

	got, err := expr.Evaluate([]fhirpath.Resource{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("got = %v", bool(got[0].(system.Boolean)))
	// Output: got = true
}
