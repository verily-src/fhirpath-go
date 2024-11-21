# FHIRPath

This package contains a Go implementation of the [FHIRPath] specification, implemented directly with
the [google/fhir] proto definitions.

This package aims to be compliant with both:

- the [N1 Normative Release](http://hl7.org/fhirpath/N1/) specification, and
- the [R4 specifications](http://hl7.org/fhir/R4/fhirpath.html).

Development of this project is tracked under the FHIRPath epic in Jira: [PHP-3997].

[Private fork]: https://github.com/verily-src/fhirpath-go
[Public mirror]: https://github.com/verily-src/fhirpath-go
[FHIRPath]: http://hl7.org/fhirpath/
[google/fhir]: https://github.com/google/fhir
[PHP-3997]: https://verily.atlassian.net/browse/PHP-3997

## Explorer

We built a little web app that allows you to play around with FHIRPath.
You can use it to test your strings, and get familiar with its
functionality.

[FHIRPath Explorer](https://dev.home.example.com/enrichment/fhirpath)

## Installation

No installation is needed, since this is all part of `verily1`. Just use from the correct import
path:

```go
import "github.com/verily-src/fhirpath-go/fhirpath"
```

## Usage

A FHIRPath must be compiled before running it against a resource using the `Compile` method like so:

```go
expression, err := fhirpath.Compile("Patient.name.given")
if err != nil {
    panic("error while compiling FHIRPath")
}
```

The compilation result can then be run against a resource:

```go
inputResources := []fhir.Resource{somePatient, someMedication}

result, err := expression.Evaluate(inputResources)
if err != nil {
    panic("error while running FHIRPath against resource")
}
```

As defined in the FHIRPath specification, the output of evaluation is a **Collection**. So, the
result of Evaluate is of type `[]any`. As such, the result must be unpacked and cast to the desired
type for further processing.

### CompileOptions and EvaluateOptions

Options are provided for optional modification of compilation and evaluation. There is currently
support for:

- adding custom functions during Compile time
- adding custom external constant variables

#### To add a custom function:

The constraints on the custom function are as follows:

- First argument must be `system.Collection`
- Arguments that follow must be either a fhir proto type or primitive system type

```go
customFn := func (input system.Collection, args ...any) (system.Collection error) {
    fmt.Print("called custom fn")
    return input, nil
}
expression, err := fhirpath.Compile("print()", compopts.AddFunction("print", customFn))
```

#### To add external constants:

The constraints on external constants are as follows:

- Must be a fhir proto type, primitive system type, or `system.Collection`
- If you pass in a collection, contained elements must be fhir proto or system type.

```go
customVar := system.String("custom variable")
result, err := expression.Evaluate([]fhir.Resource{someResource}, evalopts.EnvVariable("var", customVar))
```

### System Types

The FHIRPath [spec](http://hl7.org/fhirpath/N1/#literals) defines the following custom System types:

- Boolean
- String
- Integer
- Decimal
- Quantity
- Date
- Time
- DateTime

FHIR Protos get implicitly to the above types according to this
[chart](http://hl7.org/fhir/R4/fhirpath.html#types), when used in some FHIRPath expressions.

Package `system` provides implementations for these types, along with various exposed receiver
functions.

### Things to be aware of

FHIRPath is not the most intuitive language, and there are some quirks. See
[go/fhirpathgotchas](http://go/fhirpathgotchas)

## Dev Guide

### Go Installation

See [go/verily-go](http://go/verily-go)

### Local Development

If modifying the ANTLR parser grammar file, run `go generate ./...`

### Open Source

We maintain a [public version](https://github.com/verily-src/fhirpath-go) of this package.
The process by which we perform the private-to-public transfer is documentedin the
[private fork](https://github.com/verily-src/fhirpath-go/copy-process.md).

We keep track of the parts of the spec that we haven't implemented in 
[this sheet][wip-sheet].

The maintainers' slack channel is called [fhir-path-go-public][slack-chan].

[wip-sheet]: https://docs.google.com/spreadsheets/d/1qQHqOSff5Axn2kKtg4EWodyH5uYNZPHB2Vb03lCFXCg/edit?gid=0#gid=0
[slack-chan]: https://verily.enterprise.slack.com/archives/C071PJ7GCPQ
