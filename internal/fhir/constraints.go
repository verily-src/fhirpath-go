package fhir

import (
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
)

type primitiveDataType interface {
	*dtpb.Base64Binary |
		*dtpb.Boolean |
		*dtpb.Canonical |
		*dtpb.Code |
		*dtpb.Date |
		*dtpb.DateTime |
		*dtpb.Decimal |
		*dtpb.Id |
		*dtpb.Instant |
		*dtpb.Integer |
		*dtpb.Markdown |
		*dtpb.Oid |
		*dtpb.PositiveInt |
		*dtpb.String |
		*dtpb.Time |
		*dtpb.UnsignedInt |
		*dtpb.Uri |
		*dtpb.Url |
		*dtpb.Uuid
}

type complexDataType interface {
	*dtpb.Address |
		*dtpb.Age |
		*dtpb.Attachment |
		*dtpb.CodeableConcept |
		*dtpb.Coding |
		*dtpb.ContactPoint |
		*dtpb.Count |
		*dtpb.Distance |
		*dtpb.Duration |
		*dtpb.HumanName |
		*dtpb.Identifier |
		*dtpb.Money |
		*dtpb.MoneyQuantity |
		*dtpb.Period |
		*dtpb.Quantity |
		*dtpb.Range |
		*dtpb.Ratio |
		*dtpb.SampledData |
		*dtpb.Signature |
		*dtpb.SimpleQuantity |
		*dtpb.Timing
}

type metaDataType interface {
	*dtpb.ContactDetail |
		*dtpb.Contributor |
		*dtpb.DataRequirement |
		*dtpb.Expression |
		*dtpb.ParameterDefinition |
		*dtpb.RelatedArtifact |
		*dtpb.TriggerDefinition |
		*dtpb.UsageContext
}

type specialPurposeDataType interface {
	*dtpb.Dosage |
		*dtpb.ElementDefinition |
		*dtpb.Extension |
		*dtpb.MarketingStatus |
		*dtpb.Meta |
		*dtpb.Narrative |
		*dtpb.ProductShelfLife |
		*dtpb.Reference
}

// DataType is an constraint-definition of FHIR datatypes, which all support ID
// and Extension fields, in addition to their base values.
//
// Note: "DataType" is also an "Element", so these interfaces are logically
// equivalent -- and so this is represented as a constraint of valid datatypes.
//
// The R4 spec doesn't explicitly refer to "DataType" as a distinction from
// "Element", but the R5 spec does, and its definition is compatible with R4.
// This is retained here so that we can have a proper vernacular and mechanism
// for referring to these types in generic ways through constraints.
//
// See https://www.hl7.org/fhir/r5/types.html#DataType for more details.
type DataType interface {
	Element
	primitiveDataType | complexDataType | metaDataType | specialPurposeDataType
}

// PrimitiveType is a constraint-definition of FHIR datatypes, which all support ID
// and Extension fields, in addition to their base values.
//
// Note: "DataType" is also an "Element", so these interfaces are logically
// equivalent -- and so this is represented as a constraint of valid datatypes.
//
// The R4 spec doesn't explicitly refer to "PrimitiveType" as a distinction from
// "Element", but the R5 spec does, and its definition is compatible with R4.
// This is retained here so that we can have a proper vernacular and mechanism
// for referring to these types in generic ways through constraints.
//
// See https://www.hl7.org/fhir/types.html#PrimitiveType for more details.
type PrimitiveType interface {
	Element
	primitiveDataType
}

// BackboneType is a constraint-definition of FHIR backbone element, which all
// support ID, Extension, and modifier-extension fields, in addition to their
// base values.
//
// Note: "BackboneType" is also an "BackboneElement", so these interfaces are logically
// equivalent -- and so this is represented as a constraint of valid datatypes.
//
// The R4 spec doesn't explicitly refer to "BackboneType" as a distinction from
// "BackboneElement", but the R5 spec does, and its definition is compatible with R4.
// This is retained here so that we can have a proper vernacular and mechanism
// for referring to these types in generic ways through constraints.
//
// See https://www.hl7.org/fhir/r5/types.html#BackboneType for more details.
type BackboneType interface {
	BackboneElement
	*dtpb.Timing |
		*dtpb.ElementDefinition |
		*dtpb.MarketingStatus |
		*dtpb.ProductShelfLife |
		*dtpb.Dosage
}
