package reference

import (
	"errors"
	"fmt"
	"path"

	"github.com/google/fhir/go/jsonformat"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/element"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/protofields"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"google.golang.org/protobuf/proto"
)

var (
	ErrNotCanonicalResource = errors.New("resource type is not a Canonical resource")
	ErrNotResource          = errors.New("not a resource type")
	ErrNoResourceID         = errors.New("no resource ID")
	ErrNoResourceVersion    = errors.New("no resource version")
	ErrStrongConversion     = errors.New("error converting to strong reference from weak reference")
)

// ExtractAll finds all references contained in a resource.
func ExtractAll(resource fhir.Resource) ([]*dtpb.Reference, error) {
	return element.ExtractAll[*dtpb.Reference](resource)
}

// Canonical creates a Canonical reference for the given resource type and
// raw reference. This will error if the provided resource type is not
// a Canonical FHIR resource.
func Canonical(resourceType resource.Type, reference string) (*dtpb.Reference, error) {
	res, ok := protofields.Resources[resourceType.String()]
	if !ok {
		return nil, fmt.Errorf("%w: %T", ErrNotResource, resourceType)
	}
	testResource := res.New()
	_, isCanonical := testResource.(fhir.CanonicalResource)
	if !isCanonical {
		return nil, fmt.Errorf("%w: %T", ErrNotCanonicalResource, resourceType)
	}

	return Weak(resourceType, reference), nil
}

// Weak creates an weak reference for the given resource type and
// raw reference. The resource type must be a valid FHIR resource type.
// Generally, it is preferred to use Canonical() for Canonical references.
// Examples:
//
//	{
//	  "type": "Questionnaire",
//	  "reference": "https://example.com/questionnaire"
//	}
//
//	{
//	  "type": "GuidanceResponse",
//	  "reference": "urn:uuid:5a17b7c2-e01c-4bc7-b973-31d4156b11d7"
//	}
//
// For more info on references see:
// - https://www.hl7.org/fhir/references.html#canonical
// - https://www.hl7.org/fhir/bundle.html#references
func Weak(resourceType resource.Type, reference string) *dtpb.Reference {
	return &dtpb.Reference{
		Type: fhir.URI(resourceType.String()),
		Reference: &dtpb.Reference_Uri{
			Uri: fhir.String(reference),
		},
	}
}

// Typed creates a typed, literal FHIR reference for the given resource type and id.
// Returns an error if resourceId is invalid or resourceType is outside
// of the known R4 types (which should be impossible).
func Typed(resourceType resource.Type, resourceId string) (*dtpb.Reference, error) {
	return typedFromURIString(resourceType, path.Join(resourceType.String(), resourceId))
}

func typedFromURIString(resourceType resource.Type, uri string) (*dtpb.Reference, error) {
	weakReference := Weak(resourceType, uri)
	message := proto.Clone(weakReference)
	err := jsonformat.NormalizeReference(message)
	if err != nil {
		return nil, fmt.Errorf("error normalizing reference: %w", err)
	}
	typedReference := message.(*dtpb.Reference)
	// If conversion to a typed reference succeeded then the untyped URI reference
	// will be removed and replaced with a different option in the oneof.
	// https://github.com/google/fhir/blob/master/proto/google/fhir/proto/r4/core/datatypes.proto#L3303-L3304
	if typedReference.GetUri().GetValue() != "" {
		return nil, fmt.Errorf("%w: %v", ErrStrongConversion, typedReference)
	}
	return message.(*dtpb.Reference), nil
}

// Logical creates a logical FHIR reference for the given resource type, identifier
// system, and identifier value.
// Replaces ph.LogicalReference
func Logical(resourceType resource.Type, identifierSystem, identifierValue string) *dtpb.Reference {
	return &dtpb.Reference{
		Type:       fhir.URI(resourceType.String()),
		Identifier: fhir.Identifier(identifierSystem, identifierValue),
	}
}

// LogicalReferenceIdentifier creates a logical FHIR reference for a given resource type and
// Indentifier.
// Replaces: ph.LogicalReferenceIdentifier
func LogicalFromIdentifier(resourceType resource.Type, identifier *dtpb.Identifier) *dtpb.Reference {
	return &dtpb.Reference{
		Type:       fhir.URI(resourceType.String()),
		Identifier: identifier,
	}
}

// TypedFromResource returns a reference to the given resource that
// is strongly typed and without a version.
func TypedFromResource(res fhir.Resource) (*dtpb.Reference, error) {
	return Typed(resource.TypeOf(res), resource.ID(res))
}

// TypedFromIdentity returns a reference to the given identity that
// is always relative, always strongly typed, and will include a version ID
// if-and-only-if the given identity does.
func TypedFromIdentity(identity *resource.Identity) *dtpb.Reference {
	ref, err := typedFromURIString(identity.Type(), identity.PreferRelativeVersionedURIString())
	if err != nil {
		// Impossible to trigger this error (or at least I couldn't find a way).
		panic(err)
	}
	return ref
}

// WeakRelativeVersioned returns a reference to the given resource
// is is weakly typed, relative (to the FHIR service base URL), and versioned.
// It returns an error if either the resource's ID or version is missing.
func WeakRelativeVersioned(res fhir.Resource) (*dtpb.Reference, error) {
	identity, ok := resource.IdentityOf(res)
	if !ok {
		return nil, ErrNoResourceID // Missing ID seems most likely cause...
	}
	uri, ok := identity.RelativeVersionedURIString()
	if !ok {
		return nil, ErrNoResourceVersion
	}
	return Weak(identity.Type(), uri), nil
}

// Is compares two references for referencial equality.
// If lhs can be determined to refer to the same resource as rhs, then this
// returns true -- otherwise this returns false. If the underlying resource
// cannot be determined, then the result is false.
func Is(lhs, rhs *dtpb.Reference) bool {
	// If the two references have the same value layout, we don't need to do
	// more complicated extraction.
	if proto.Equal(lhs, rhs) {
		return true
	}

	// Check for logical referencial equality
	if li, ri := lhs.GetIdentifier(), rhs.GetIdentifier(); !proto.Equal(li, ri) {
		return false
	}

	// Check for literal referencial equality. This requires determining the literal
	// identity of the referenced resource, since protos represent references as
	// either a URI _or_ a strongly-typed reference ID in a oneof field. This will
	// extract the relevant parts to allow for proper equality.
	if lhs.Reference != nil || rhs.Reference != nil {
		left, err := IdentityOf(lhs)
		if err != nil {
			return false
		}
		right, err := IdentityOf(rhs)
		if err != nil {
			return false
		}

		return left.Equal(right)
	}
	return true
}
