package patch

import (
	"errors"
	"fmt"
	"strings"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/iancoleman/strcase"
	"github.com/verily-src/fhirpath-go/fhirpath"
	"github.com/verily-src/fhirpath-go/fhirpath/compopts"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/compile"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/opts"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/parser"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"github.com/verily-src/fhirpath-go/internal/slices"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	ErrNotImplemented     = errors.New("not implemented")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidEnum        = errors.New("invalid enum value")
	ErrInvalidField       = fhirpath.ErrInvalidField
	ErrInvalidUnsignedInt = errors.New("invalid value for unsigned int")
	ErrNotSingleton       = expr.ErrNotSingleton
	ErrNotPatchable       = errors.New("result is not patchable")
)

// Options encapsulates all possible FHIRPath options that
// can be used in the underlying FHIRPath evaluation before
// patching. This includes both compile-time and evaluation-time
// options.
type Options struct {
	CompileOpts []opts.CompileOption
	EvalOpts    []opts.EvaluateOption
}

// Expression is the FHIRPath Patch expression that will be
// compiled from a FHIRPath string.
type Expression struct {
	expression expr.Expression
	path       string
}

// String returns the underlying FHIRPath expression.
func (e *Expression) String() string {
	return e.path
}

// Compile parses and compiles the FHIRPath Patch expression down
// to a single Expression object.
//
// If there are any syntax or semantic errors, this will return an
// error indicating the reason for the compilation failure.
func Compile(path string, options ...opts.CompileOption) (*Expression, error) {
	options = append(options, compopts.Transform(func(e expr.Expression) expr.Expression {
		return storeLastExpression{e}
	}))

	config, err := compile.PopulateConfig(options...)
	if err != nil {
		return nil, err
	}

	tree, err := compile.Tree(path)
	if err != nil {
		return nil, err
	}

	visitor := &parser.FHIRPathVisitor{
		Functions:  config.Table,
		Transform:  config.Transform,
		Permissive: config.Permissive,
	}
	vr, ok := visitor.Visit(tree).(*parser.VisitResult)
	if !ok {
		return nil, errors.New("input expression currently unsupported")
	}

	if vr.Error != nil {
		return nil, vr.Error
	}
	return &Expression{
		expression: vr.Result,
		path:       path,
	}, nil
}

// Add appends a value to the given field name in the element.
// Add can be used for non-repeating elements so long as they do not already exist.
// The field name must be in the same case as defined in the FHIRPath spec,
// which is in camelCase
//
// Note: value inputs to Add must be FHIR elements or FHIR resources.
//
// This function will return the following errors in the given conditions:
//
//   - ErrInvalidInput if the input value to add is incorrect for the result of
//     the expression
//   - ErrInvalidField if the specified field name does not exist in the
//     returned element
//   - ErrNotSingleton if the result of the evaluation returns more than one
//     entry
//   - ErrNotPatchable if the result of the returned entry is not a patchable
//     type -- e.g. a non-scalar type like a list, a value that is already
//     set, etc.
//
// See documentation: https://hl7.org/fhir/R4/fhirpatch.html#concept.
func (e *Expression) Add(res fhirpath.Resource, name string, value fhir.Base, options ...fhirpath.EvaluateOption) error {
	if strcase.ToLowerCamel(name) != name {
		// All field names in FHIRPath are in camelCase, but the protos are in
		// snake_case. To avoid accidentally accepting code like "foo_value" instead
		// of "fooValue", we check first that we are already in the correct form,
		// and error if it would never be possible.
		return fmt.Errorf("%w: '%v'", ErrInvalidField, name)
	}
	if res == nil {
		return fmt.Errorf("%w: nil input resource", ErrInvalidInput)
	}
	if value == nil {
		return fmt.Errorf("%w: nil replacement value", ErrInvalidInput)
	}

	_, evalResult, err := e.evaluate(res, options...)
	if err != nil {
		return err
	}

	singleton, err := evalResult.ToSingleton()
	if err != nil {
		return fmt.Errorf("%w: fhirpatch add requires singleton collection", ErrNotSingleton)
	}

	proto, ok := singleton.(proto.Message)
	if !ok {
		return fmt.Errorf("%w: result of type '%T' is not patchable", ErrNotPatchable, singleton)
	}

	ref := proto.ProtoReflect()
	descriptor := ref.Descriptor()
	fieldName := strcase.ToSnake(name)
	field := descriptor.Fields().ByName(protoreflect.Name(fieldName))
	if field == nil {
		fieldName += "_value"
		field = descriptor.Fields().ByName(protoreflect.Name(fieldName))
		if field == nil {
			return fmt.Errorf("%w: '%v'", fhirpath.ErrInvalidField, name)
		}
	}

	if !field.IsList() && ref.Has(field) {
		return fmt.Errorf("%w: unable to add value to populated scalar field '%v' in %v resource", ErrNotPatchable, name, resource.TypeOf(res))
	}

	var update func(m protoreflect.ProtoMessage)
	var valueMessage protoreflect.Message
	if field.IsList() {
		list := ref.Mutable(field).List()
		update = func(m protoreflect.ProtoMessage) {
			list.Append(protoreflect.ValueOfMessage(m.ProtoReflect()))
		}
		valueMessage = list.NewElement().Message()
	} else {
		update = func(m protoreflect.ProtoMessage) {
			ref.Set(field, protoreflect.ValueOfMessage(m.ProtoReflect()))
		}
		valueMessage = ref.Get(field).Message()
	}

	// Special handling for oneof fields, like "Extension", "ContainedResource", etc.
	if e.isSingletonOneof(valueMessage.Interface()) {
		container := e.newSetOneof(valueMessage, value)
		if container == nil {
			return fmt.Errorf(
				"%w: '%v' value provided for field '%v' (which is of type '%v')",
				ErrInvalidInput,
				value.ProtoReflect().Descriptor().Name(),
				name,
				valueMessage.Descriptor().Name(),
			)
		}
		update(container.Interface())
	} else {
		// Normalize data being patched
		value, err = e.normalizeAdd(valueMessage, value)
		if err != nil {
			return err
		}

		if valueMessage.Descriptor() != value.ProtoReflect().Descriptor() {
			return fmt.Errorf(
				"%w: '%v' value provided for field '%v' (which is of type '%v')",
				ErrInvalidInput,
				value.ProtoReflect().Descriptor().Name(),
				name,
				valueMessage.Descriptor().Name(),
			)
		}
		update(value)
	}

	return nil
}

// stringable is an interface to check for a string-valued FHIR type.
// code, markdown and id are all specializations of string that satisfy
// this interface.
// See: https://hl7.org/fhir/r4/datatypes.html
type stringable interface {
	GetValue() string
}

// stringable is an interface to check for a integer-valued FHIR type.
// See: https://hl7.org/fhir/r4/datatypes.html
type intable interface {
	GetValue() int32
}

// normalizeAdd normalizes a value to be patched to the correct type.
// If no normalization is required, the input value will be returned
// unmodified from its original value.
func (e *Expression) normalizeAdd(valueMessage protoreflect.Message, value fhir.Base) (fhir.Base, error) {
	var newVal fhir.Base
	var err error
	switch value := value.(type) {
	case stringable:
		newVal, err = enumFromStringable(valueMessage, value)
		if newVal == nil && err == nil {
			// Check for a reference field - since these are dynamic,
			// we need to patch them in after the reference is created,
			// which is why we're only updating the ID field here.
			valueField := valueMessage.Descriptor().Fields().ByName("value")
			if valueField != nil && valueField.FullName() == "google.fhir.r4.core.ReferenceId.value" {
				newVal = &dtpb.ReferenceId{Value: value.GetValue()}
			}
		}
	case intable:
		newVal, err = intValueFromInt(valueMessage, value)
	default:
	}

	if err != nil {
		return nil, err
	}
	// The value was normalized
	if newVal != nil {
		return newVal, nil
	}

	// Use the original value without normalization
	return value, nil
}

func (e *Expression) evaluate(res fhirpath.Resource, options ...fhirpath.EvaluateOption) (*expr.Context, system.Collection, error) {
	collection := system.Collection{res}
	config := &opts.EvaluateConfig{
		Context: expr.InitializeContext(collection),
	}
	config, err := opts.ApplyOptions(config, options...)
	if err != nil {
		return nil, nil, err
	}

	result, err := e.expression.Evaluate(config.Context, collection)
	return config.Context, result, err
}

func (e *Expression) isSingletonOneof(msg proto.Message) bool {
	message := msg.ProtoReflect()
	descriptor := message.Descriptor()
	oneofs := descriptor.Oneofs()
	return oneofs.Len() == 1 && oneofs.ByName("reference") == nil
}

func (e *Expression) newSetOneof(msg protoreflect.Message, value proto.Message) protoreflect.Message {
	container := msg.New()
	descriptor := container.Descriptor()
	fields := descriptor.Fields()

	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		if value.ProtoReflect().Descriptor() == msg.Get(field).Message().Descriptor() {
			container.Set(field, protoreflect.ValueOfMessage(value.ProtoReflect()))
			return container
		}
	}
	return nil
}

// Delete removes the element of the evaluated expression. It can only remove
// single elements from a resource.
//
// See documentation: https://hl7.org/fhir/R4/fhirpatch.html#concept.
func (e *Expression) Delete(res fhirpath.Resource, options ...fhirpath.EvaluateOption) error {
	if res == nil {
		return fmt.Errorf("%w: nil input resource", ErrInvalidInput)
	}
	ctx, evalResult, err := e.evaluate(res, options...)
	if err != nil {
		return err
	}
	// If we have an empty value, it means the field is already deleted.
	if evalResult.IsEmpty() {
		return nil
	}
	toDelete, err := evalResult.ToSingleton()
	if err != nil {
		return fmt.Errorf("%w: fhirpatch delete can only delete a single element", ErrNotSingleton)
	}

	if err := e.tryDelete(ctx.LastResult, toDelete); err != nil {
		if err := e.tryDelete(ctx.BeforeLastResult, toDelete); err != nil {
			return err
		}
	}

	return nil
}

func (e *Expression) tryDelete(collection system.Collection, toDelete any) error {
	var message protoreflect.Message
	var field protoreflect.FieldDescriptor
	var idx int
	for _, entry := range collection {
		var ok bool
		root, ok := entry.(proto.Message)
		if !ok {
			continue
		}

		field, idx, ok = e.getFieldForCollection(root, system.Collection{toDelete})
		if ok {
			message = root.ProtoReflect()
			break
		}
	}
	if field == nil {
		return fmt.Errorf("%w: field cannot be deleted", ErrNotPatchable)
	}
	if field.IsList() {
		list := message.Get(field).List()
		if idx == -1 {
			if list.Len() <= 1 {
				message.Clear(field)
				return nil
			}
			return fmt.Errorf("%w: list containing more than one element cannot be deleted", ErrNotPatchable)
		}
		newlist := message.NewField(field).List()
		for i := 0; i < idx; i++ {
			newlist.Append(list.Get(i))
		}
		for i := idx + 1; i < list.Len(); i++ {
			newlist.Append(list.Get(i))
		}
		message.Set(field, protoreflect.ValueOfList(newlist))
	} else {
		message.Clear(field)
	}
	return nil
}

// Insert inserts a value into the expression's list, at the 0-based index specified.
// Prefer Add() if you are inserting at the end of a list.
//
// See documentation: https://hl7.org/fhir/R4/fhirpatch.html#concept.
func (e *Expression) Insert(res fhirpath.Resource, value fhir.Base, index int, options ...fhirpath.EvaluateOption) error {
	if res == nil {
		return fmt.Errorf("%w: nil input resource", ErrInvalidInput)
	}
	ctx, evalResult, err := e.evaluate(res, options...)
	if err != nil {
		return err
	}
	last, err := ctx.LastResult.ToSingleton()
	if err != nil {
		return fmt.Errorf("%w: fhirpatch insert requires single element to operate on", ErrNotSingleton)
	}
	root, ok := last.(proto.Message)
	if !ok {
		return fmt.Errorf("%w: %T type is not a FHIR type", ErrNotPatchable, last)
	}
	field, _, ok := e.getFieldForCollection(root, evalResult)
	if !ok {
		return fmt.Errorf("%w: field is empty", ErrNotPatchable)
	}
	if !field.IsList() {
		return fmt.Errorf("%w: named field is not a list", ErrNotPatchable)
	}

	reflect := root.ProtoReflect()
	existing := reflect.Get(field).List()
	if index > existing.Len() || index < 0 {
		return fmt.Errorf("%w: index %v is out of range", ErrNotPatchable, index)
	}
	if elem := existing.NewElement(); elem.Message().Descriptor() != value.ProtoReflect().Descriptor() {
		return fmt.Errorf("%w: Element %T is not assignable to %T", ErrNotPatchable, value, elem)
	}
	// Recreate the list at this field
	list := reflect.NewField(field).List()

	// Rebuild the list in order
	for i := 0; i < existing.Len(); i++ {
		if i == index {
			list.Append(protoreflect.ValueOfMessage(value.ProtoReflect()))
		}
		list.Append(existing.Get(i))
	}
	// Handle insertion at the end
	if index == existing.Len() {
		list.Append(protoreflect.ValueOfMessage(value.ProtoReflect()))
	}
	reflect.Set(field, protoreflect.ValueOfList(list))

	return nil
}

// getFieldForCollection returns the FieldDescriptor that corresponds to the field
// that contains the entries in `collection`.
func (e *Expression) getFieldForCollection(root proto.Message, collection system.Collection) (protoreflect.FieldDescriptor, int, bool) {
	if len(collection) == 0 {
		return nil, -1, false
	}
	ref := root.ProtoReflect()
	descriptor := ref.Descriptor()
	for _, entry := range collection {
		msg, ok := entry.(proto.Message)
		if !ok {
			continue
		}
		fields := descriptor.Fields()
		for i := 0; i < fields.Len(); i++ {
			field := fields.Get(i)
			if !ref.Has(field) {
				continue
			}

			value := ref.Get(field)
			if field.Cardinality() == protoreflect.Repeated {
				list := value.List()
				for i := 0; i < list.Len(); i++ {
					entry := list.Get(i).Message().Interface()
					if e.unwrapOneof(entry) == msg {
						return field, i, true
					}
				}
			} else if field.Kind() == protoreflect.MessageKind {
				if e.unwrapOneof(value.Message().Interface()) == msg {
					return field, -1, true
				}
			} else {
				return nil, -1, false
			}
		}
	}
	return nil, -1, false
}

// Move moves an element within the expression's list from one index to another.
//
// See documentation: https://hl7.org/fhir/R4/fhirpatch.html#concept.
func (e *Expression) Move(resource fhirpath.Resource, sourceIndex, destIndex int, options ...fhirpath.EvaluateOption) error {
	return ErrNotImplemented
}

// Replace replaces the original value of the expression with the provided value.
//
// See documentation: https://hl7.org/fhir/R4/fhirpatch.html#concept.
func (e *Expression) Replace(resource fhirpath.Resource, value fhir.Base, options ...fhirpath.EvaluateOption) error {
	if resource == nil {
		return fmt.Errorf("%w: nil input resource", ErrInvalidInput)
	}
	ctx, evalResult, err := e.evaluate(resource, options...)
	if err != nil {
		return err
	}
	toReplace, err := evalResult.ToSingleton()
	if err != nil {
		return fmt.Errorf("%w: fhirpatch replace requires singleton collection", ErrNotSingleton)
	}

	if err := e.tryReplace(ctx.LastResult, toReplace, value); err != nil {
		if err := e.tryReplace(ctx.BeforeLastResult, toReplace, value); err != nil {
			return err
		}
	}

	return nil
}

// tryReplace attempts to first find the field that contains the element to replace, then replace it
// with the provided value.
func (e *Expression) tryReplace(collection system.Collection, toReplace any, value fhir.Base) error {
	ref, field, idx, err := e.getRefAndFieldForCollection(collection, toReplace)
	if err != nil {
		return err
	}

	var valueMessage protoreflect.Message
	var update func(m protoreflect.ProtoMessage)
	if field.IsList() {
		list := ref.Get(field).List()
		newlist := ref.NewField(field).List()
		valueMessage = newlist.NewElement().Message()
		update = func(m protoreflect.ProtoMessage) {
			for i := 0; i < list.Len(); i++ {
				if i == idx {
					newlist.Append(protoreflect.ValueOfMessage(m.ProtoReflect()))
				} else {
					newlist.Append(list.Get(i))
				}
			}
			ref.Set(field, protoreflect.ValueOfList(newlist))
		}
	} else {
		valueMessage = ref.Get(field).Message()
		update = func(m protoreflect.ProtoMessage) {
			ref.Set(field, protoreflect.ValueOfMessage(m.ProtoReflect()))
		}
	}

	// Special handling for oneof fields, like "Extension", "ContainedResource", etc.
	if e.isSingletonOneof(valueMessage.Interface()) {
		container := e.newSetOneof(valueMessage, value)
		if container == nil {
			return fmt.Errorf(
				"%w: '%v' value provided for field '%v' (which is of type '%v')",
				ErrInvalidInput,
				value.ProtoReflect().Descriptor().Name(),
				e.path,
				valueMessage.Descriptor().Name(),
			)
		}
		update(container.Interface())
	} else {
		// Normalize data being patched
		value, err := e.normalizeAdd(valueMessage, value)
		if err != nil {
			return err
		}

		if valueMessage.Descriptor() != value.ProtoReflect().Descriptor() {
			return fmt.Errorf(
				"%w: '%v' value provided for field '%v' (which is of type '%v')",
				ErrInvalidInput,
				value.ProtoReflect().Descriptor().Name(),
				e.path,
				valueMessage.Descriptor().Name(),
			)
		}
		update(value)
	}
	return nil
}

func (e *Expression) getRefAndFieldForCollection(collection system.Collection, toReplace any) (protoreflect.Message, protoreflect.FieldDescriptor, int, error) {
	for _, entry := range collection {
		var ok bool
		root, ok := entry.(proto.Message)
		if !ok {
			continue
		}

		field, idx, ok := e.getFieldForCollection(root, system.Collection{toReplace})
		if ok {
			ref := root.ProtoReflect()
			return ref, field, idx, nil
		}
	}

	return nil, nil, -1, fmt.Errorf("%w: field cannot be replaced", ErrNotPatchable)
}

// TODO: make this a common function
func (e *Expression) unwrapOneof(obj proto.Message) proto.Message {
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

// Add appends a value to the element identified in the path, using the name specified.
// Add can be used for non-repeating elements so long as they do not already exist.
//
// See documentation: https://hl7.org/fhir/R4/fhirpatch.html#concept.
func Add(resource fhirpath.Resource, path, name string, value fhir.Base, opts *Options) error {
	expr, err := Compile(path, opts.CompileOpts...)
	if err != nil {
		return err
	}
	return expr.Add(resource, name, value, opts.EvalOpts...)
}

// Delete removes the element at the specified path. It can only remove
// single elements from a resource.
//
// See documentation: https://hl7.org/fhir/R4/fhirpatch.html#concept.
func Delete(resource fhirpath.Resource, path string, options ...opts.CompileOption) error {
	expr, err := Compile(path, options...)
	if err != nil {
		return err
	}
	return expr.Delete(resource)
}

// Insert inserts a value into the specified list, at the 0-based index specified.
// Prefer Add() if you are inserting at the end of a list.
//
// See documentation: https://hl7.org/fhir/R4/fhirpatch.html#concept.
func Insert(resource fhirpath.Resource, path string, value fhir.Base, index int, options ...opts.CompileOption) error {
	expr, err := Compile(path, options...)
	if err != nil {
		return err
	}
	return expr.Insert(resource, value, index)
}

// Move moves an element within the specified list from one index to another.
//
// See documentation: https://hl7.org/fhir/R4/fhirpatch.html#concept.
func Move(resource fhirpath.Resource, path string, sourceIndex, destIndex int, options ...opts.CompileOption) error {
	expr, err := Compile(path, options...)
	if err != nil {
		return err
	}
	return expr.Move(resource, sourceIndex, destIndex)
}

// Replace replaces the original value at the specified path with the provided value.
//
// See documentation: https://hl7.org/fhir/R4/fhirpatch.html#concept.
func Replace(resource fhirpath.Resource, path string, value fhir.Base, options ...opts.CompileOption) error {
	expr, err := Compile(path, options...)
	if err != nil {
		return err
	}
	return expr.Replace(resource, value)
}

// storeLastExpression is a simple Expression object that can be used to store
// the last result of an evaluation (e.g. the last returned collection that
// occurs before the last evaluation node).
type storeLastExpression struct {
	delegate expr.Expression
}

func (e storeLastExpression) Evaluate(ctx *expr.Context, in system.Collection) (system.Collection, error) {
	// Only store the last result if the slice is not identical to the previous one.
	// This exists in case an intermediate or final node is a no-op that does not
	// alter the slice, e.g.: `Patient.name.trace('something')` -- which would
	// yield the same output as `Patient.name` would.
	if !slices.IsIdentical(ctx.LastResult, in) {
		ctx.BeforeLastResult = ctx.LastResult
		ctx.LastResult = in
	}
	return e.delegate.Evaluate(ctx, in)
}

// enumFromStringable parses a string value into an enum if
// the value field's type is an enum.
func enumFromStringable(msg protoreflect.Message, val stringable) (fhir.Base, error) {
	strVal := val.GetValue()
	container := msg.New()
	valueField := container.Descriptor().Fields().ByName("value")
	if valueField != nil && valueField.Kind() == protoreflect.EnumKind {
		if strcase.ToKebab(strVal) != strVal {
			return nil, fmt.Errorf("%w: %q", ErrInvalidEnum, strVal)
		}
		enumValueStr := protoreflect.Name(strcase.ToScreamingSnake(strVal))
		enum := valueField.Enum().Values().ByName(enumValueStr)
		if enum == nil {
			return nil, fmt.Errorf("%w: %q", ErrInvalidEnum, enumValueStr)
		}
		enumVal := protoreflect.ValueOfEnum(protoreflect.EnumNumber(enum.Number()))
		container.Set(valueField, enumVal)
		return container.Interface(), nil
	}
	return nil, nil
}

// intValueFromInt converts an integer to the appropriate type.
// The type could be integer, unsignedInt, or positiveInt.
func intValueFromInt(msg protoreflect.Message, val intable) (fhir.Base, error) {
	container := msg.New()
	valueField := container.Descriptor().Fields().ByName("value")
	if valueField != nil {
		var intValue protoreflect.Value
		switch valueField.Kind() {
		case protoreflect.Int32Kind:
			intValue = protoreflect.ValueOfInt32(val.GetValue())
		case protoreflect.Uint32Kind:
			if val.GetValue() < 0 {
				return nil, fmt.Errorf("%w: %v", ErrInvalidUnsignedInt, val.GetValue())
			}
			intValue = protoreflect.ValueOfUint32(uint32(val.GetValue()))
		default:
		}
		container.Set(valueField, intValue)
		return container.Interface(), nil
	}
	return nil, nil
}
