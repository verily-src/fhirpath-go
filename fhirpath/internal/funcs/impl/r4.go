package impl

import (
	"fmt"

	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/protofields"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// Extension is syntactic sugar over `extension.where(url = ...)`, and is
// specific to the R4 extensions for FHIRPath (as oppose to being part of the
// N1 normative spec).
//
// For more details, see https://hl7.org/fhir/R4/fhirpath.html#functions
func Extension(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("%w: expected 1 argument", ErrWrongArity)
	}
	arg, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	}
	str, err := arg.ToString()
	if err != nil {
		return nil, err
	}

	var result system.Collection
	for _, entry := range input {
		entry, ok := entry.(fhir.Extendable)
		if !ok {
			continue
		}
		for _, ext := range entry.GetExtension() {
			if url := ext.GetUrl(); url != nil && url.Value == str {
				result = append(result, ext)
			}
		}
	}
	return result, nil
}

// HasValue returns true if the input collection contains a single value which is a FHIR primitive,
// and it has a primitive value (e.g. as opposed to not having a value and just having extensions).
//
// For more details, see https://hl7.org/fhir/R4/fhirpath.html#functions
func HasValue(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 0", ErrWrongArity, len(args))
	}
	if !input.IsSingleton() {
		return system.Collection{system.Boolean(false)}, nil
	}

	if primitive, ok := input[0].(fhir.Base); ok {
		msg := primitive.ProtoReflect()

		// attempt to unwrap polymorphic types
		oneOf := protofields.UnwrapOneofField(input[0].(fhir.Base), "choice")
		if oneOf != nil {
			msg = oneOf.ProtoReflect()
		} else if !system.IsPrimitive(input[0]) {
			return system.Collection{system.Boolean(false)}, nil
		}

		descriptor := msg.Descriptor()
		field := descriptor.Fields().ByName(protoreflect.Name("value"))
		if field != nil && msg.Has(field) {
			return system.Collection{system.Boolean(true)}, nil
		}
	}

	return system.Collection{system.Boolean(false)}, nil
}
