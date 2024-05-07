package fhir

import (
	"time"

	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/slices"
)

const ucumUnitSystem = "http://unitsofmeasure.org"

// General-Purpose Data Types:
//
// The section below defines types from the "Data Types" heading in
// http://hl7.org/fhir/R4/datatypes.html#open

// Annotation creates an R4 FHIR Annotation element with the specified text,
// author, and time of creation.
//
// See: http://hl7.org/fhir/R4/datatypes.html#annotation
func Annotation(text, author string, when time.Time) *dtpb.Annotation {
	return &dtpb.Annotation{
		Text: Markdown(text),
		Author: &dtpb.Annotation_AuthorX{
			Choice: &dtpb.Annotation_AuthorX_StringValue{
				StringValue: String(author),
			},
		},
		Time: DateTime(when),
	}
}

// AnnotationReference creates an R4 FHIR Annotation element with the specified
// text, a reference to the author, and the time of creation.
//
// See: http://hl7.org/fhir/R4/datatypes.html#annotation
func AnnotationReference(text string, author *dtpb.Reference, when time.Time) *dtpb.Annotation {
	return &dtpb.Annotation{
		Text: Markdown(text),
		Author: &dtpb.Annotation_AuthorX{
			Choice: &dtpb.Annotation_AuthorX_Reference{
				Reference: author,
			},
		},
		Time: DateTime(when),
	}
}

// Coding creates an R4 FHIR Coding element with the provided system and code.
//
// See: http://hl7.org/fhir/R4/datatypes.html#coding
func Coding(system, code string) *dtpb.Coding {
	return &dtpb.Coding{
		System: URI(system),
		Code:   Code(code),
	}
}

// CodingWithVersion creates an R4 FHIR Coding element with the provided system,
// code, and version.
//
// See: http://hl7.org/fhir/R4/datatypes.html#coding
func CodingWithVersion(system, code, version string) *dtpb.Coding {
	return &dtpb.Coding{
		System:  URI(system),
		Code:    Code(code),
		Version: String(version),
	}
}

// CodeableConcept creates an R4 FHIR CodeableConcept with the specified codings,
// and with the Text element if the given text argument is non-empty.
//
// Providing a non-empty Text element is good practice but not required.
// See: http://hl7.org/fhir/R4/datatypes.html#codeableconcept
func CodeableConcept(text string, coding ...*dtpb.Coding) *dtpb.CodeableConcept {
	concept := &dtpb.CodeableConcept{
		Coding: coding,
	}
	if text != "" {
		concept.Text = String(text)
	}
	return concept
}

// ContactPoint creates an R4 FHIR ContactPoint element from the system and value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#contactpoint
func ContactPoint(system cpb.ContactPointSystemCode_Value, value string) *dtpb.ContactPoint {
	return &dtpb.ContactPoint{
		System: &dtpb.ContactPoint_SystemCode{
			Value: system,
		},
		Value: String(value),
	}
}

// EmailContactPoint creates an R4 FHIR ContactPoint element for the Email
// system given a value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#contactpoint
func EmailContactPoint(value string) *dtpb.ContactPoint {
	return ContactPoint(cpb.ContactPointSystemCode_EMAIL, value)
}

// PhoneContactPoint creates an R4 FHIR ContactPoint element for the Phone
// system given a value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#contactpoint
func PhoneContactPoint(value string) *dtpb.ContactPoint {
	return ContactPoint(cpb.ContactPointSystemCode_PHONE, value)
}

// SmsContactPoint creates an R4 FHIR ContactPoint element for the SMS system
// given a value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#contactpoint
func SmsContactPoint(value string) *dtpb.ContactPoint {
	return ContactPoint(cpb.ContactPointSystemCode_SMS, value)
}

// PagerContactPoint creates an R4 FHIR ContactPoint element for the Pager
// system given a value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#contactpoint
func PagerContactPoint(value string) *dtpb.ContactPoint {
	return ContactPoint(cpb.ContactPointSystemCode_PAGER, value)
}

// FaxContactPoint creates an R4 FHIR ContactPoint element for the Fax system
// given a value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#contactpoint
func FaxContactPoint(value string) *dtpb.ContactPoint {
	return ContactPoint(cpb.ContactPointSystemCode_FAX, value)
}

// OtherContactPoint creates an R4 FHIR ContactPoint element for the Other
// system given a value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#contactpoint
func OtherContactPoint(value string) *dtpb.ContactPoint {
	return ContactPoint(cpb.ContactPointSystemCode_OTHER, value)
}

// Identifier creates an R4 FHIR Identifier element with the provided system
// and value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#identifier
func Identifier(system, value string) *dtpb.Identifier {
	return &dtpb.Identifier{
		System: URI(system),
		Value:  String(value),
	}
}

// Money creates an R4 FHIR Money element from the money value.
//
// See: http://hl7.org/fhir/R4/datatypes.html#Money
func Money(value float64) *dtpb.Money {
	return &dtpb.Money{
		Value: Decimal(value),
	}
}

// MoneyQuantity creates an R4 FHIR MoneyQuantity element from the value and units.
//
// See: http://hl7.org/fhir/R4/datatypes.html#MoneyQuantity
func MoneyQuantity(value float64, unit string) *dtpb.MoneyQuantity {
	return &dtpb.MoneyQuantity{
		Value: Decimal(value),
		Unit:  String(unit),
	}
}

// Period creates an R4 FHIR Period element with the provided start and end times.
//
// See: http://hl7.org/fhir/R4/datatypes.html#period
func Period(start, end time.Time) *dtpb.Period {
	return &dtpb.Period{
		Start: DateTime(start),
		End:   DateTime(end),
	}
}

// Quantity creates an R4 FHIR Quantity element from the given value and units.
//
// See: http://hl7.org/fhir/R4/datatypes.html#quantity
func Quantity(value float64, unit string) *dtpb.Quantity {
	return &dtpb.Quantity{
		Value: Decimal(value),
		Unit:  String(unit),
	}
}

// UCUMQuantity creates an R4 FHIR Quantity element representing a
// value and UCUM unit.
//
// See: http://hl7.org/fhir/R4/datatypes.html#quantity
// TODO(PHP-9521): Add a unit package to validate against UCUM units.
func UCUMQuantity(value float64, unit string) *dtpb.Quantity {
	return &dtpb.Quantity{
		Value:  Decimal(value),
		Unit:   String(unit),
		Code:   Code(unit),
		System: URI(ucumUnitSystem),
	}
}

// QuantityFromSimpleQuantity is a convenience utility for converting a
// SimpleQuantity to its base-class definition of Quantity.
// If the input is nil, this returns nil.
//
// For more information, see the diagram for Primitive Types here:
// https://www.hl7.org/fhir/datatypes.html
func QuantityFromSimpleQuantity(value *dtpb.SimpleQuantity) *dtpb.Quantity {
	return quantityFrom(value, value == nil)
}

// QuantityFromDuration is a convenience utility for converting a
// Duration to its base-class definition of Quantity.
// If the input is nil, this returns nil.
//
// For more information, see the diagram for Primitive Types here:
// https://www.hl7.org/fhir/datatypes.html
func QuantityFromDuration(value *dtpb.Duration) *dtpb.Quantity {
	return quantityFrom(value, value == nil)
}

// QuantityFromMoneyQuantity is a convenience utility for converting a
// MoneyQuantity to its base-class definition of Quantity.
// If the input is nil, this returns nil.
//
// For more information, see the diagram for Primitive Types here:
// https://www.hl7.org/fhir/datatypes.html
func QuantityFromMoneyQuantity(value *dtpb.MoneyQuantity) *dtpb.Quantity {
	return quantityFrom(value, value == nil)
}

// quantityLike is a helper interface for implementing the QuantityFrom*
// functions. This defines the common interface for all Quantity objects.
type quantityLike interface {
	GetValue() *dtpb.Decimal
	GetUnit() *dtpb.String
	GetSystem() *dtpb.Uri
	GetCode() *dtpb.Code
}

// quantityFrom is a helper function for implementing the QuantityFrom* functions
// which all have the same implementation.
//
// This function takes 'isNil' as a boolean argument to work around the fact that
// a nil pointer passed to an interface forms a non-nil interface in Go, and
// reflection is more costly than a bool check.
func quantityFrom(value quantityLike, isNil bool) *dtpb.Quantity {
	if isNil {
		return nil
	}
	return &dtpb.Quantity{
		Value:  value.GetValue(),
		Unit:   value.GetUnit(),
		System: value.GetSystem(),
		Code:   value.GetCode(),
	}
}

// Range creates an R4 FHIR Range element with the given low and high end of the
// range, using the specified units.
//
// See: http://hl7.org/fhir/R4/datatypes.html#range
func Range(low, high float64, unit string) *dtpb.Range {
	return &dtpb.Range{
		Low:  SimpleQuantity(low, unit),
		High: SimpleQuantity(high, unit),
	}
}

// Ratio creates an R4 FHIR Ratio element with the given numerator and denominator.
//
// See: http://hl7.org/fhir/R4/datatypes.html#ratio
func Ratio(numerator, denominator float64) *dtpb.Ratio {
	return &dtpb.Ratio{
		Numerator:   &dtpb.Quantity{Value: Decimal(numerator)},
		Denominator: &dtpb.Quantity{Value: Decimal(denominator)},
	}
}

// SimpleQuantity creates an R4 FHIR SimpleQuantity element from the given value
// and units.
//
// See: http://hl7.org/fhir/R4/datatypes.html#SimpleQuantity
func SimpleQuantity(value float64, unit string) *dtpb.SimpleQuantity {
	return &dtpb.SimpleQuantity{
		Value: Decimal(value),
		Unit:  String(unit),
	}
}

// Timing creates an R4 FHIR Timing element observing the events specified in `times`.
//
// See: http://hl7.org/fhir/R4/datatypes.html#timing
func Timing(times ...time.Time) *dtpb.Timing {
	return &dtpb.Timing{
		Event: slices.Map(times, DateTime),
	}
}
