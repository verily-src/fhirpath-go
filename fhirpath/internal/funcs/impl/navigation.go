package impl

import (
	"fmt"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/protofields"
)

// Children returns a collection with all immediate child nodes of all items in the input collection
// with no specific order.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#children-collection
func Children(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	if length := len(args); length != 0 {
		return nil, fmt.Errorf("%w, received %v arguments, expected 0", ErrWrongArity, length)
	}

	result := system.Collection{}
	for _, item := range input {
		if system.IsPrimitive(item) {
			continue
		}
		base, ok := item.(fhir.Base)
		if !ok {
			return nil, fmt.Errorf("%w: unexpected input of type '%T'", ErrInvalidInput, item)
		}
		if oneOf := protofields.UnwrapOneofField(base, "choice"); oneOf != nil {
			if system.IsPrimitive(oneOf) {
				continue
			}
			base = oneOf
		}

		var fields []string
		if _, ok := base.(*dtpb.Reference); ok {
			fields = append(fields, "reference")
		} else {
			fd := base.ProtoReflect().Descriptor().Fields()
			for i := 0; i < fd.Len(); i++ {
				fields = append(fields, fd.Get(i).JSONName())
			}
		}
		for _, f := range fields {
			fe := expr.FieldExpression{FieldName: f}
			messages, err := fe.Evaluate(ctx, system.Collection{base})
			if err != nil {
				return nil, err
			}
			for _, val := range messages {
				if c, ok := val.(system.Collection); ok {
					result = append(result, c...)
				} else {
					result = append(result, val)
				}
			}
		}
	}
	return result, nil
}

// Descendants returns a collection with all descendant nodes of all items in the input collection.
// FHIRPath docs here: https://hl7.org/fhirpath/N1/#descendants-collection
func Descendants(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	if length := len(args); length != 0 {
		return nil, fmt.Errorf("%w, received %v arguments, expected 0", ErrWrongArity, length)
	}

	result := system.Collection{}
	for !input.IsEmpty() {
		var err error
		input, err = Children(ctx, input)
		if err != nil {
			return nil, err
		}
		result = append(result, input...)
	}
	return result, nil
}
