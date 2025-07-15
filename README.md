# FHIRPath

Developed and maintained by the nice folks at [Verily](https://verily.com/).

This package contains a Go implementation of the [FHIRPath][fhirpath] specification, implemented directly with
the [google/fhir][google-fhir] proto definitions.

This package aims to be compliant with both:

- the [N1 Normative Release](http://hl7.org/fhirpath/N1/) specification, and
- the [R4 specifications](http://hl7.org/fhir/R4/fhirpath.html).

## Import

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
inputResources := []fhirpath.Resource{somePatient, someMedication}

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

#### To add a custom function

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

#### To add external constants

The constraints on external constants are as follows:

- Must be a fhir proto type, primitive system type, or `system.Collection`
- If you pass in a collection, contained elements must be fhir proto or system type.

```go
customVar := system.String("custom variable")
result, err := expression.Evaluate([]fhirpath.Resource{someResource}, evalopts.EnvVariable("var", customVar))
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

FHIR Protos get implicitly converted to the above types according to this
[chart](http://hl7.org/fhir/R4/fhirpath.html#types), when used in some FHIRPath expressions.

### Things to be aware of

FHIRPath is not the most intuitive language, and there are some quirks. See [gotchas](gotchas.md).

[fhirpath]: http://hl7.org/fhirpath/
[google-fhir]: https://github.com/google/fhir
