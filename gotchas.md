# FHIRPath Gotcha’s

## Empty collections are propagated

* In FHIRPath, whenever an empty collection is encountered, rather than raising an error it gets propagated throughout the rest of the expression. This may make some issues difficult to catch.
* Eg. given `Patient.name` -> `{}`,  `Patient.name.family + ' MD'` -> `{}`

## Equality sometimes returns an empty collection { }, rather than false

* If either collection is empty
* If the **precision_ _**of Date, Time, or DateTime objects are mismatched
* If the **dimension** of a Quantity unit is mismatched

## FHIR type specifiers are case-sensitive

* **Primitive** types are denoted with lower case specifiers.
* **Primitive** types that are written as upper case will be resolved as **System** types, not **FHIR** types.
* Eg. `Patient.birthDate is date = **true**` but `Patient.birthDate is Date = **false**`
* Case should match what’s listed [here](https://www.hl7.org/fhir/r4/datatypes.html)
* System types always begin with an uppercase letter

## `As` Expression is _not_ a filter, expects singleton input

* The as expression (`Observation.value as integer`) expects a singleton as input. For example, if you pass in a resource with multiple value fields, it will raise an error.
* It doesn’t filter out things that don’t match the type. For this purpose, the `where()` function should be used -> `where(value is integer)`
