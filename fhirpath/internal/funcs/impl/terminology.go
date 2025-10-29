package impl

import (
	"errors"
	"fmt"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/fhirpath/terminology"
)

var (
	ErrUnconfiguredClient = errors.New("memberOf() function requires a Terminology Service client to be configured in the evaluation context")
	ErrNotSupported       = errors.New("Not Supported, memberOf must be called on a single Coding or a CodeableConcept")
)

// MemberOf takes a Coding or a Codeable concept and a ValueSet internal id
// and determines if the coding of any of the coding is inside the ValueSet
func MemberOf(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	if length := len(args); length != 1 {
		return nil, fmt.Errorf("%w: received %d arguments, expected 1", ErrWrongArity, length)
	}

	if len(input) != 1 {
		return system.Collection{}, ErrNotSupported
	}

	content, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("argument dereference error: %s", err)
	}

	if len(content) == 0 {
		return nil, fmt.Errorf("empty argument content")
	}

	var valueSetId string
	switch content[0].(type) {
	case system.String:
		valueSetId = string(content[0].(system.String))
	case *dtpb.String:
		valueSetId = content[0].(*dtpb.String).GetValue()
	default:
		return nil, fmt.Errorf("unsupported argument type")
	}

	validateResult := false
	for _, item := range input {
		switch res := item.(type) {
		case *dtpb.Coding:
			result, err := validateCoding(ctx, res.GetCode().GetValue(), res.GetSystem().GetValue(), valueSetId)
			if err != nil {
				return system.Collection{system.Boolean(false)}, nil
			}
			validateResult = result

		case *dtpb.CodeableConcept:
			// If it's a Codeable Concept, we will checking the coding inside one by one
			for _, coding := range res.GetCoding() {
				result, err := validateCoding(ctx, coding.GetCode().GetValue(), coding.GetSystem().GetValue(), valueSetId)
				if err != nil {
					continue
				}
				if result {
					validateResult = result
					break
				}
			}
		default:
			return nil, ErrNotSupported
		}
	}
	return system.Collection{system.Boolean(validateResult)}, nil
}

func validateCoding(ctx *expr.Context, code string, system string, valueSet string) (bool, error) {
	// We will not evaluate code without system
	if system == "" {
		return false, nil
	}

	ts := ctx.TermService
	if ts == nil {
		return false, ErrUnconfiguredClient
	}

	opt := &terminology.ValueSetValidateCodeOptions{
		Code:   code,
		System: system,
		ID:     valueSet,
	}

	response, err := ts.ValueSetValidateCode(ctx, opt)
	if err != nil {
		return false, fmt.Errorf("validating valueSet code: %s", err)
	}

	result := false
	for _, item := range response.Parameter {
		switch item.GetName().GetValue() {
		case "result":
			result = item.GetValue().GetBoolean().GetValue()
		default:
			continue
		}
	}

	return result, nil
}
