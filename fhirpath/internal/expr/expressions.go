package expr

import (
	"errors"
	"fmt"
	"strings"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	bcrpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	"github.com/iancoleman/strcase"
	"github.com/shopspring/decimal"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/reflection"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/containedresource"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/fhirconv"
	"github.com/verily-src/fhirpath-go/internal/protofields"
	"github.com/verily-src/fhirpath-go/internal/slices"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	ErrNotSingleton     = errors.New("collection is not a singleton")
	ErrInvalidType      = errors.New("collection evaluates to incorrect type")
	ErrInvalidOperator  = errors.New("received invalid operator")
	ErrToBeImplemented  = errors.New("expression not yet implemented")
	ErrInvalidField     = errors.New("invalid field")
	ErrConstantNotFound = errors.New("external constant not found")
)

// Expression is the abstraction for all FHIRPath expressions,
// ie. taking in some input collection and outputting some collection.
type Expression interface {
	Evaluate(*Context, system.Collection) (system.Collection, error)
}

// ExpressionSequence abstracts the flow of evaluation for
// compound expressions. It consists of a sequence of expressions
// whose outputs flow into the inputs of the next expression.
type ExpressionSequence struct {
	Expressions []Expression
}

// Evaluate iterates through the ExpressionSequence, feeding the output of
// an evaluation to the next Expression.
func (s *ExpressionSequence) Evaluate(ctx *Context, input system.Collection) (system.Collection, error) {
	output := input

	for _, expr := range s.Expressions {
		result, err := expr.Evaluate(ctx, output)
		// raise error as soon as one is encountered
		if err != nil {
			return nil, err
		}
		output = result
	}
	return output, nil
}

var _ Expression = (*ExpressionSequence)(nil)

// IdentityExpression encapsulates the top-level expression, ie. the
// first step in the chain of evaluation. A no-op expression that
// returns itself.
type IdentityExpression struct{}

// Evaluate returns the input collection, without raising
// an error.
func (*IdentityExpression) Evaluate(ctx *Context, input system.Collection) (system.Collection, error) {
	return input, nil
}

var _ Expression = (*IdentityExpression)(nil)

// FieldExpression is the expression that accesses the specified
// FieldName in the input collection.
type FieldExpression struct {
	FieldName  string
	Permissive bool
}

// Evaluate filters the input collections by those that contain
// the FieldName string, and returns the result.
func (e *FieldExpression) Evaluate(ctx *Context, input system.Collection) (system.Collection, error) {
	output := system.Collection{}

	for _, item := range input {
		message, ok := item.(proto.Message)
		if !ok {
			if e.Permissive {
				continue
			}
			return nil, e.errField(item)
		}

		// Date, Time, DateTime, and Instant have "fake" fields 'value_us', 'timezone',
		// and 'precision'. This checks to ensure that such fields aren't being accessed,
		// since they aren't actually real and don't exist in the FHIR spec.
		if !e.isEvaluable(message) {
			return nil, e.errField(message)
		}

		message = e.unpackAny(message)
		// unwrap if a ContainedResource
		if contained, ok := message.(*bcrpb.ContainedResource); ok {
			message = containedresource.Unwrap(contained)
		}

		// Get desired field
		fieldName := strcase.ToSnake(e.FieldName)
		reflect := message.ProtoReflect()
		field := reflect.Descriptor().Fields().ByName(protoreflect.Name(fieldName))

		// extract field and append to output, flattening
		// if the field is a list. Raises error if field doesn't exist
		if field == nil {
			// If the field is a reference, we need to combine the type and
			// ID fields to create a usable reference, e.g. Type/ID. Since
			// ID is a oneof (e.g. questionnaire_id), we need to determine
			// which it is to find the appropriate field.
			if fieldName == "reference" {
				if reference, ok := message.(*dtpb.Reference); ok {
					refString := e.unwrapReference(reference)
					if refString != nil {
						output = append(output, refString)
					}
					continue
				}
			}

			// Attempting to get a "value" field from a Date, DateTime, Time, or Instant
			// needs to convert the value to a System String type.
			// The FHIR Protos model time datatypes using a "value_us" field, which
			// is normalized here, since the FHIR spec models these types as strings
			// with a "value" field.
			if fieldName == "value" {
				switch v := message.(type) {
				case *dtpb.Date:
					output = append(output, system.String(fhirconv.DateToString(v)))
					continue
				case *dtpb.DateTime:
					output = append(output, system.String(fhirconv.DateTimeToString(v)))
					continue
				case *dtpb.Time:
					output = append(output, system.String(fhirconv.TimeToString(v)))
					continue
				case *dtpb.Instant:
					output = append(output, system.String(fhirconv.InstantToString(v)))
					continue
				}
			}

			// Try again with "_value" added because sometimes Google protos do that
			// for primitives like:
			// Observation.ValueX.String --> Observation_ValueX_StringValue
			fieldName = fieldName + "_value"
			field = reflect.Descriptor().Fields().ByName(protoreflect.Name(fieldName))
			if field == nil {
				return nil, fmt.Errorf("%w: %s not a field on %T", ErrInvalidField, fieldName, message)
			}
		}

		// If the field is not a message, it is a primitive (enum or go native type).
		// So, it can be cast to a system type. Otherwise, a field is being accessed that
		// shouldn't be accessed, so the error is returned.
		if field.Kind() != protoreflect.MessageKind {
			primitive, err := system.From(message)
			if err != nil {
				return nil, err
			}
			output = append(output, primitive)
			continue
		}

		unwrap := func(obj protoreflect.ProtoMessage) protoreflect.ProtoMessage {
			obj = e.unpackAny(obj)
			if contained, ok := obj.(*bcrpb.ContainedResource); ok {
				obj = containedresource.Unwrap(contained)
			}
			return e.unwrapOneof(obj)
		}
		if e.Permissive {
			unwrap = func(obj proto.Message) proto.Message { return obj }
		}

		if !field.IsList() {
			message := reflect.Get(field).Message()
			if !message.IsValid() {
				continue
			}
			output = append(output, unwrap(message.Interface()))
			continue
		}
		content := reflect.Get(field).List()
		for i := 0; i < content.Len(); i++ { // flatten out list
			result := content.Get(i).Message().Interface()
			output = append(output, unwrap(result))
		}
	}
	return output, nil
}

var nonEvaluableFields = []string{
	"valueUs", "precision", "timezone",
}

func (e *FieldExpression) isEvaluable(msg proto.Message) bool {
	if e.Permissive {
		return true
	}

	// Prevent snake_case fields, since all FHIRPath fields need to be in
	// camelCase.
	if strcase.ToLowerCamel(e.FieldName) != e.FieldName {
		return false
	}

	// Prevent manually accessing idiosynchratic fields from google/fhir like
	// value_us, precision, and time_zone
	switch msg.(type) {
	case *dtpb.Time, *dtpb.Date, *dtpb.DateTime, *dtpb.Instant:
		return !slices.Includes(nonEvaluableFields, e.FieldName)
	}

	return true
}

func (e *FieldExpression) errField(object any) error {
	return fmt.Errorf("%w: %s not a field on %T", ErrInvalidField, e.FieldName, object)
}

func (e *FieldExpression) unwrapReference(ref *dtpb.Reference) *dtpb.String {
	if ref.GetReference() == nil {
		return nil
	}
	rv := ref.ProtoReflect()
	switch ref := ref.GetReference().(type) {
	case *dtpb.Reference_Uri:
		return fhir.String(ref.Uri.GetValue())
	case *dtpb.Reference_Fragment:
		return fhir.String("#" + ref.Fragment.GetValue())
	default:
		descriptor := rv.Descriptor()
		oneof := descriptor.Oneofs().ByName("reference")
		field := rv.WhichOneof(oneof)
		if field == nil {
			return nil
		}
		refid := rv.Get(field).Message().Interface().(*dtpb.ReferenceId)
		fieldName, ok := strings.CutSuffix(string(field.Name()), "_id")
		if !ok {
			return nil
		}
		fieldName = strcase.ToCamel(fieldName)
		if history := refid.GetHistory(); history != nil {
			return fhir.String(fmt.Sprintf("%v/%v/_history/%v", fieldName, refid.GetValue(), history.GetValue()))
		}
		return fhir.String(fmt.Sprintf("%v/%v", fieldName, refid.GetValue()))
	}
}

func (e *FieldExpression) unwrapOneof(obj proto.Message) proto.Message {
	message := obj.ProtoReflect()
	descriptor := message.Descriptor()
	if name := string(descriptor.Name()); !(strings.HasSuffix(name, "ValueX") || name == "ContainedResource") {
		return obj
	}
	oneofsNum := descriptor.Oneofs().Len()
	if oneofsNum != 1 {
		return obj
	}

	oneof := descriptor.Oneofs().Get(0)
	field := message.WhichOneof(oneof)
	if oneof == nil || field == nil {
		return obj
	}
	if msg := message.Get(field).Message(); msg != nil {
		return msg.Interface()
	}
	return obj
}

func (e *FieldExpression) unpackAny(obj protoreflect.ProtoMessage) protoreflect.ProtoMessage {
	if anyMsg, ok := obj.(*anypb.Any); ok {
		cr := &bcrpb.ContainedResource{}
		if err := anyMsg.UnmarshalTo(cr); err == nil {
			return cr
		}
	}
	return obj
}

var _ Expression = (*FieldExpression)(nil)

// TypeExpression contains the FHIR Type identifier string,
// to be able to filter the items in the input collection that have the
// given type.
type TypeExpression struct {
	Type string
}

// Evaluate filters the messages in the input that are identified by the Type
// defined in the expression.
func (e *TypeExpression) Evaluate(ctx *Context, input system.Collection) (system.Collection, error) {
	output := system.Collection{}

	for _, item := range input {
		message, ok := item.(proto.Message)
		if !ok {
			continue
		}

		// find message name, add to collection only if it matches
		pReflect := message.ProtoReflect()
		name := pReflect.Descriptor().Name()

		if string(name) != e.Type {
			continue
		}

		output = append(output, message)
	}
	return output, nil
}

var _ Expression = (*TypeExpression)(nil)

// LiteralExpression abstracts FHIRPath system types, that
// are returned from parsing literals.
type LiteralExpression struct {
	Literal system.Any
}

// Evaluate returns the contained literal, without
// raising an error. Returns an empty collection if the
// literal is nil, representing a Null literal.
func (e *LiteralExpression) Evaluate(*Context, system.Collection) (system.Collection, error) {
	if e.Literal != nil {
		return system.Collection{e.Literal}, nil
	}
	return system.Collection{}, nil
}

var _ Expression = (*LiteralExpression)(nil)

// IndexExpression allows accessing of an input system.Collection's index.
// Contains an expression, that when evaluated, should return an integer
// that represents the index.
type IndexExpression struct {
	Index Expression
}

// Evaluate indexes the input system.Collection and returns the
// item located at the given index. If the index is negative
// or out of bounds, returns an empty collection. Raises an error if the
// contained expression does not evaluate to an expression, or raises an error
// itself.
func (e *IndexExpression) Evaluate(ctx *Context, input system.Collection) (system.Collection, error) {
	indexResult, err := e.Index.Evaluate(ctx, input)
	if err != nil {
		return nil, err
	}
	length := len(indexResult)
	if length == 0 {
		return system.Collection{}, nil
	}
	if length > 1 {
		return nil, fmt.Errorf("%w: contains %v elements", ErrNotSingleton, length)
	}
	value, err := system.From(indexResult[0])
	if err != nil {
		return nil, err
	}
	index, ok := value.(system.Integer)
	if !ok {
		return nil, fmt.Errorf("%w: want Integer but got %T", ErrInvalidType, index)
	}
	if int(index) >= len(input) || int(index) < 0 {
		return system.Collection{}, nil
	}
	return system.Collection{input[int(index)]}, nil
}

var _ Expression = (*IndexExpression)(nil)

// EqualityExpression allows checking equality of the two contained
// subexpressions. The two expressions should return comparable values
// when evaluated.
type EqualityExpression struct {
	Left  Expression
	Right Expression
	Not   bool
}

// Evaluate evaluates the two subexpressions, and returns true if their
// contents are equal, using the functionality of system.Collection.Equal. If either
// collection is empty, returns an empty collection.
func (e *EqualityExpression) Evaluate(ctx *Context, input system.Collection) (system.Collection, error) {
	leftResult, err := e.Left.Evaluate(ctx.Clone(), input)
	if err != nil {
		return nil, err
	}
	rightResult, err := e.Right.Evaluate(ctx.Clone(), input)
	if err != nil {
		return nil, err
	}
	if len(leftResult) == 0 || len(rightResult) == 0 {
		return system.Collection{}, nil
	}

	result, ok := leftResult.TryEqual(rightResult)
	if !ok {
		return system.Collection{}, nil
	}
	if e.Not {
		result = !result
	}
	return system.Collection{system.Boolean(result)}, nil
}

var _ Expression = (*EqualityExpression)(nil)

// FunctionExpression enables evaluation of Function Invocation expressions.
// It holds the function and function arguments.
type FunctionExpression struct {
	Fn   func(*Context, system.Collection, ...Expression) (system.Collection, error)
	Args []Expression
}

// Evaluate evaluates the function with respect to its arguments. Returns the result
// of the function, or an error if raised.
func (e *FunctionExpression) Evaluate(ctx *Context, input system.Collection) (system.Collection, error) {
	return e.Fn(ctx.Clone(), input, e.Args...)
}

var _ Expression = (*FunctionExpression)(nil)

// IsExpression enables evaluation of an "is" type expression.
type IsExpression struct {
	Expr Expression
	Type reflection.TypeSpecifier
}

// Evaluate evaluates the contained expression with respect to singleton evaluation
// of collections, and determines whether or not it is the given type.
func (e *IsExpression) Evaluate(ctx *Context, input system.Collection) (system.Collection, error) {
	result, err := e.Expr.Evaluate(ctx, input)
	if err != nil {
		return nil, err
	}
	length := len(result)
	if length == 0 {
		return system.Collection{}, nil
	}
	if length > 1 {
		return nil, fmt.Errorf("%w: contains %v elements", ErrNotSingleton, length)
	}
	typeSpecifier, err := reflection.TypeOf(result[0])
	if err != nil {
		return nil, err
	}
	return system.Collection{typeSpecifier.Is(e.Type)}, nil
}

var _ Expression = (*IsExpression)(nil)

// AsExpression enables evaluation of an "as" type expression.
type AsExpression struct {
	Expr Expression
	Type reflection.TypeSpecifier
}

// Evaluate evaluates the contained expression with respect to singleton evaluation
// of collections, returns the singleton if it is of the given type. Returns empty otherwise.
func (e *AsExpression) Evaluate(ctx *Context, input system.Collection) (system.Collection, error) {
	result, err := e.Expr.Evaluate(ctx, input)
	if err != nil {
		return nil, err
	}
	length := len(result)
	if length == 0 {
		return system.Collection{}, nil
	}
	if length > 1 {
		return nil, fmt.Errorf("%w: contains %v elements", ErrNotSingleton, length)
	}
	typeSpecifier, err := reflection.TypeOf(result[0])
	if err != nil {
		return nil, err
	}
	if !typeSpecifier.Is(e.Type) {
		return system.Collection{}, nil
	}
	// attempt to unwrap polymorphic types
	message, ok := result[0].(fhir.Base)
	if !ok {
		return result, nil
	}
	if oneOf := protofields.UnwrapOneofField(message, "choice"); oneOf != nil {
		return system.Collection{oneOf}, nil
	}
	return result, nil
}

var _ Expression = (*AsExpression)(nil)

// BooleanExpression enables evaluation of boolean expressions,
// including "and", "or", "xor", and "implies".
type BooleanExpression struct {
	Left  Expression
	Right Expression
	Op    Operator
}

// Evaluate evaluates the subexpressions with respect to singleton evaluation of
// collections, and performs the respective Boolean operation.
func (e *BooleanExpression) Evaluate(ctx *Context, input system.Collection) (system.Collection, error) {
	leftResult, err := e.Left.Evaluate(ctx.Clone(), input)
	if err != nil {
		return nil, err
	}
	rightResult, err := e.Right.Evaluate(ctx.Clone(), input)
	if err != nil {
		return nil, err
	}

	leftBool, err := leftResult.ToSingletonBoolean()
	if err != nil {
		return nil, err
	}
	rightBool, err := rightResult.ToSingletonBoolean()
	if err != nil {
		return nil, err
	}

	switch e.Op {
	case And:
		return evaluateAnd(leftBool, rightBool), nil
	case Or:
		return evaluateOr(leftBool, rightBool), nil
	case Xor:
		return evaluateXor(leftBool, rightBool), nil
	case Implies:
		return evaluateImplies(leftBool, rightBool), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidOperator, e.Op)
	}
}

var _ Expression = (*BooleanExpression)(nil)

type ComparisonExpression struct {
	Left  Expression
	Right Expression
	Op    Operator
}

// Evaluate evaluates the subexpressions with respect to singleton evaluation of collections,
// and performs the respective comparison operation.
func (e *ComparisonExpression) Evaluate(ctx *Context, input system.Collection) (system.Collection, error) {
	leftResult, err := e.Left.Evaluate(ctx.Clone(), input)
	if err != nil {
		return nil, err
	}
	rightResult, err := e.Right.Evaluate(ctx.Clone(), input)
	if err != nil {
		return nil, err
	}

	if len(leftResult) == 0 || len(rightResult) == 0 {
		return system.Collection{}, nil
	}
	if len(leftResult) != 1 || len(rightResult) != 1 {
		return nil, fmt.Errorf("%w: left contains %v elements, right contains %v elements", ErrNotSingleton, len(leftResult), len(rightResult))
	}

	leftPrimitive, err := system.From(leftResult[0])
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidType, err)
	}
	rightPrimitive, err := system.From(rightResult[0])
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidType, err)
	}

	// Implicitly convert types
	leftPrimitive = system.Normalize(leftPrimitive, rightPrimitive)
	rightPrimitive = system.Normalize(rightPrimitive, leftPrimitive)

	// Calculate both less than and greater than
	lessThan, err := leftPrimitive.Less(rightPrimitive)
	if errors.Is(err, system.ErrMismatchedPrecision) || errors.Is(err, system.ErrMismatchedUnit) {
		return system.Collection{}, nil
	}
	if err != nil {
		return nil, err
	}
	greaterThan, err := rightPrimitive.Less(leftPrimitive)
	if errors.Is(err, system.ErrMismatchedPrecision) {
		return system.Collection{}, nil
	}
	if err != nil {
		return nil, err
	}

	switch e.Op {
	case Lt:
		return system.Collection{lessThan}, nil
	case Gt:
		return system.Collection{greaterThan}, nil // (a > b) = (b < a)
	case Lte:
		return system.Collection{!greaterThan}, nil // (a <= b) = !(a > b)
	case Gte:
		return system.Collection{!lessThan}, nil // (a >= b) = !(a < b)
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidOperator, e.Op)
	}
}

var _ Expression = (*ComparisonExpression)(nil)

// ArithmeticExpression enables mathematical arithmetic operations.
// Includes '+', '-', '*', "/", "div", and 'mod'.
type ArithmeticExpression struct {
	Left  Expression
	Right Expression
	Op    func(system.Any, system.Any) (system.Any, error)
}

// Evaluate evaluates the two subexpressions, with respect to singleton evaluation of collections,
// and performs the respective additive operation.
func (e *ArithmeticExpression) Evaluate(ctx *Context, input system.Collection) (system.Collection, error) {
	leftResult, err := e.Left.Evaluate(ctx.Clone(), input)
	if err != nil {
		return nil, err
	}
	rightResult, err := e.Right.Evaluate(ctx.Clone(), input)
	if err != nil {
		return nil, err
	}

	if len(leftResult) == 0 || len(rightResult) == 0 {
		return system.Collection{}, nil
	}
	if len(leftResult) != 1 || len(rightResult) != 1 {
		return nil, fmt.Errorf("%w: left contains %v elements, right contains %v elements", ErrNotSingleton, len(leftResult), len(rightResult))
	}

	// Cast contents to system types. Addition and subtraction is not supported for protos.
	leftPrimitive, err := system.From(leftResult[0])
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidType, err)
	}
	rightPrimitive, err := system.From(rightResult[0])
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidType, err)
	}

	// Implicitly convert types
	leftPrimitive = system.Normalize(leftPrimitive, rightPrimitive)
	rightPrimitive = system.Normalize(rightPrimitive, leftPrimitive)

	result, err := e.Op(leftPrimitive, rightPrimitive)
	if errors.Is(err, system.ErrIntOverflow) {
		return system.Collection{}, nil // "Operations that cause arithmetic overflow or underflow will result in empty ( { } )".
	}
	if err != nil {
		return nil, err
	}
	return system.Collection{result}, nil
}

var _ Expression = (*ArithmeticExpression)(nil)

// ConcatExpression enables the evaluation of a string concatenation expression.
type ConcatExpression struct {
	Left  Expression
	Right Expression
}

// Evaluate evaluates the two subexpressions with respect to singleton evaluation of
// collections, and attempts to concatenate the two strings. Returns an error if the expressions
// don't resolve to strings. This differs from string addition when either collection is empty. Rather
// than returning empty, it will treat the empty collection as an empty string.
func (e *ConcatExpression) Evaluate(ctx *Context, input system.Collection) (system.Collection, error) {
	leftResult, err := e.Left.Evaluate(ctx.Clone(), input)
	if err != nil {
		return nil, err
	}
	rightResult, err := e.Right.Evaluate(ctx.Clone(), input)
	if err != nil {
		return nil, err
	}

	// Convert empty collection to empty string
	if len(leftResult) == 0 {
		leftResult = append(leftResult, system.String(""))
	}
	if len(rightResult) == 0 {
		rightResult = append(rightResult, system.String(""))
	}

	if len(leftResult) > 1 || len(rightResult) > 1 {
		return nil, fmt.Errorf("%w: left contains %v elements, right contains %v elements", ErrNotSingleton, len(leftResult), len(rightResult))
	}

	// Cast contents to system types
	leftPrimitive, err := system.From(leftResult[0])
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidType, err)
	}
	rightPrimitive, err := system.From(rightResult[0])
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidType, err)
	}

	leftStr, ok := leftPrimitive.(system.String)
	if !ok {
		return nil, fmt.Errorf("%w: expected a string", ErrInvalidType)
	}
	rightStr, ok := rightPrimitive.(system.String)
	if !ok {
		return nil, fmt.Errorf("%w: expected a string", ErrInvalidType)
	}

	return system.Collection{leftStr + rightStr}, nil
}

var _ Expression = (*ConcatExpression)(nil)

// ExternalConstantExpression enables evaluation of external constants.
type ExternalConstantExpression struct {
	Identifier string
}

// Evaluate retrieves the constant from the map located in the Context. Returns an error if the
// constant is not present.
func (e *ExternalConstantExpression) Evaluate(ctx *Context, input system.Collection) (system.Collection, error) {
	constant, ok := ctx.ExternalConstants[e.Identifier]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrConstantNotFound, e.Identifier)
	}
	if collection, ok := constant.(system.Collection); ok {
		return collection, nil
	}
	return system.Collection{constant}, nil
}

var _ Expression = (*ExternalConstantExpression)(nil)

// NegationExpression enables negation of number values (Integer, Decimal, Quantity).
type NegationExpression struct {
	Expr Expression
}

// Evaluate negates the contained expression, that is evaluated with respect to singleton evaluation.
// If the contained value is not a number, returns an error.
func (e *NegationExpression) Evaluate(ctx *Context, input system.Collection) (system.Collection, error) {
	result, err := e.Expr.Evaluate(ctx, input)
	if err != nil {
		return nil, err
	}

	length := len(result)
	if length == 0 {
		return system.Collection{}, nil
	}
	if length != 1 {
		return nil, fmt.Errorf("%w: can't negate a collection", ErrNotSingleton)
	}

	primitive, err := system.From(result[0])
	if err != nil {
		return nil, fmt.Errorf("%w: can't negate complex type %T", ErrInvalidType, result[0])
	}

	// handle negation of value
	switch v := primitive.(type) {
	case system.Integer:
		return system.Collection{system.Integer(-1) * v}, nil
	case system.Decimal:
		negative := system.Decimal(decimal.NewFromInt(-1))
		return system.Collection{v.Mul(negative)}, nil
	case system.Quantity:
		return system.Collection{v.Negate()}, nil
	default:
		return nil, fmt.Errorf("%w: can't negate %T", ErrInvalidType, primitive)
	}
}

var _ Expression = (*NegationExpression)(nil)
