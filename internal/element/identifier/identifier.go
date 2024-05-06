package identifier

import (
	"fmt"
	"net/url"

	cpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

// Use is an alias of the Identifier Use-codes for easier access and readability.
type Use = cpb.IdentifierUseCode_Value

const (
	// UseUsual is an alias of the USUAL Identifier use-code for easier access and
	// readability.
	UseUsual Use = cpb.IdentifierUseCode_USUAL

	// UseOfficial is an alias of the OFFICIAL Identifier use-code for easier access
	// and readability.
	UseOfficial Use = cpb.IdentifierUseCode_OFFICIAL

	// UseTemp is an alias of the TEMP Identifier use-code for easier access and
	// readability.
	UseTemp Use = cpb.IdentifierUseCode_TEMP

	// UseSecondary is an alias of the SECONDARY Identifier use-code for easier
	// access and readability.
	UseSecondary Use = cpb.IdentifierUseCode_SECONDARY

	// UseOld is an alias of the OLD Identifier use-code for easier access and
	// readability.
	UseOld Use = cpb.IdentifierUseCode_OLD
)

// New constructs a new Identifier object from the given options.
func New(value, system string, opts ...Option) *dtpb.Identifier {
	identifier := &dtpb.Identifier{
		System: fhir.URI(system),
		Value:  fhir.String(value),
	}
	return Update(identifier, opts...)
}

// Usual is a convenience constructor for forming an identifier with a "Usual"
// Use-code assigned.
func Usual(value, system string, opts ...Option) *dtpb.Identifier {
	return newWithUse(UseUsual, value, system, opts...)
}

// Official is a convenience constructor for forming an identifier with an
// "Official" Use-code assigned.
func Official(value, system string, opts ...Option) *dtpb.Identifier {
	return newWithUse(UseOfficial, value, system, opts...)
}

// Temp is a convenience constructor for forming an identifier with a "Temp"
// Use-code assigned.
func Temp(value, system string, opts ...Option) *dtpb.Identifier {
	return newWithUse(UseTemp, value, system, opts...)
}

// Secondary is a convenience constructor for forming an identifier with a
// "Secondary" Use-code assigned.
func Secondary(value, system string, opts ...Option) *dtpb.Identifier {
	return newWithUse(UseSecondary, value, system, opts...)
}

// Old is a convenience constructor for forming an identifier with an "Old"
// Use-code assigned.
func Old(value, system string, opts ...Option) *dtpb.Identifier {
	return newWithUse(UseOld, value, system, opts...)
}

func newWithUse(use Use, value, system string, opts ...Option) *dtpb.Identifier {
	identifier := &dtpb.Identifier{
		System: fhir.URI(system),
		Value:  fhir.String(value),
		Use: &dtpb.Identifier_UseCode{
			Value: use,
		},
	}
	return Update(identifier, opts...)
}

// Update modifies an identifier in-place with the given identifier options.
//
// This function returns the input identifier to allow for functional chaining.
func Update(identifier *dtpb.Identifier, opts ...Option) *dtpb.Identifier {
	for _, opt := range opts {
		opt.update(identifier)
	}
	return identifier
}

// Equivalent checks if two identifiers are equivalent by comparing the
// system and values of the identifier.
func Equivalent(lhs, rhs *dtpb.Identifier) bool {
	lsystem, rsystem := lhs.GetSystem(), rhs.GetSystem()
	lvalue, rvalue := lhs.GetValue(), rhs.GetValue()

	if lsystem.GetValue() != rsystem.GetValue() {
		return false
	}
	if lvalue.GetValue() != rvalue.GetValue() {
		return false
	}
	return true
}

// FindBySystem searches a slice of identifiers for the first identifier that
// contains the specified system.
func FindBySystem(identifiers []*dtpb.Identifier, system string) *dtpb.Identifier {
	for _, identifier := range identifiers {
		identifierSystem := identifier.GetSystem()
		if identifierSystem == nil {
			continue
		}
		if identifierSystem.GetValue() == system {
			return identifier
		}
	}
	return nil
}

// QueryString formats a system and value for use in a Search query,
// escaping FHIR special characters `,|$\` in the input.
// Use this in a query param as the value with key `identifier`.
func QueryString(system string, value string) string {
	// escape special characters like `|` from identifier
	escapedSystem := fhir.EscapeSearchParam(system)
	escapedValue := fhir.EscapeSearchParam(value)
	return fmt.Sprintf("%s|%s", escapedSystem, escapedValue)
}

// QueryIdentifier formats an Identifier proto for use in a Search query,
// escaping FHIR special characters `,|$\` in the input.
// Use this in a query param as the value with key `identifier`.
func QueryIdentifier(id *dtpb.Identifier) string {
	return QueryString(id.System.Value, id.Value.Value)
}

// GenerateIfNoneExist takes an Identifier and generates a query appropriate for use in an If-None-Exist header.
// This is used for FHIR conditional create or other conditional methods.
//
// Untrusted data in Identifiers is escaped both for FHIR and for URL safety.
//
// Returns an empty string if identifier is nil.
//
// This function only accepts a single identifier due to limitations of the GCP
// FHIR store. Important note:
// The GCP FHIR store only supports conditional queries on a single identifier,
// with no modifiers (so identifier=foo is OK, while identifier:exact=foo is
// invalid). Deviating from this in API v1 will result in an HTTP 400 invalid
// query error. NB: Deviating from this in API v1beta1 results in silent
// fallback to non-transactional search, meaning the conditional queries will
// have race conditions.
func GenerateIfNoneExist(identifier *dtpb.Identifier) string {
	if identifier == nil {
		return ""
	}

	// build up a query string like format for the header
	// data is URL encoded in case the identifier contains special characters
	qs := url.Values{}

	search := QueryIdentifier(identifier)
	// use identifier= rather than identifier:exact=, see limitation above
	qs.Add("identifier", search)

	// URL encode the result
	return qs.Encode()
}
