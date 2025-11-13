package impl

import (
	"errors"
	"fmt"

	"github.com/google/fhir/go/jsonformat"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"google.golang.org/protobuf/proto"
)

var (
	ErrInvalidReference     = errors.New("invalid reference")
	ErrUnconfiguredResolver = errors.New("resolve() function requires a Resolver to be configured in the evaluation context")
)

func isValidResolveInput(item any) bool {
	switch item.(type) {
	case *dtpb.String, *dtpb.Uri, *dtpb.Url, *dtpb.Canonical, *dtpb.Reference, system.String, string:
		return true
	default:
		return false
	}
}

func stringifyReference(ref *dtpb.Reference) string {
	if ref == nil {
		return ""
	}
	r := proto.Clone(ref).(*dtpb.Reference)
	if err := jsonformat.DenormalizeReference(r); err != nil {
		return ""
	}
	return r.GetUri().GetValue()
}

func toString(item any) string {
	switch item := item.(type) {
	case *dtpb.String:
		return item.GetValue()
	case *dtpb.Uri:
		return item.GetValue()
	case *dtpb.Url:
		return item.GetValue()
	case *dtpb.Canonical:
		return item.GetValue()
	case *dtpb.Reference:
		return stringifyReference(item)
	case system.String:
		return string(item)
	default:
		return ""
	}
}

// Resolve locates the target of each reference in the input collection. For each item that
// is a string representing a URI (or canonical or URL), or a Reference, Resolve locates
// the target resource and adds it to the resulting collection. If an item does not resolve
// to a resource, it is ignored. If the input is empty, the output will be empty.
func Resolve(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	if argLen := len(args); argLen != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, argLen)
	}

	toResolve := []string{}
	for _, item := range input {
		if isValidResolveInput(item) {
			ref := toString(item)
			if ref != "" {
				toResolve = append(toResolve, ref)
			}
		}
	}

	if len(toResolve) == 0 {
		return system.Collection{}, nil
	}

	resolverImpl := ctx.Resolver
	if resolverImpl == nil {
		return nil, ErrUnconfiguredResolver
	}
	resources, err := resolverImpl.Resolve(toResolve)
	if err != nil {
		return nil, err
	}

	resolved := system.Collection{}
	for _, res := range resources {
		resolved = append(resolved, res)
	}
	return resolved, nil
}
