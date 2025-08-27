/*
Package fhirjson provides JSON marshalling and unmarshalling support for the
FHIR Resource types defined in package fhir.

This exposes both an idiomatic `fhirjson.Marshal`/`fhirjson.Unmarshal` written
in terms of the `fhir.Resource` interface, unlike the `jsonformat` package's
approach of operating only with ContainedResource objects.

This package also aims to set opinionated defaults to the formatted outputs,
while also remaining more visible to consumers than the inappropriately named
`jsonformat` package, which is actually *required* to form valid FHIR JSON
since it remaps various `Element` types to the required fields in FHIR.
*/
package fhirjson
