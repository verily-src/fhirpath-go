package protofields

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Overwrite overwrites a field of the given name for the specified message.
// If `values` is empty, the field is cleared.
// If `values` contains more than one entry for a non-repeated field, this panics.
func Overwrite(in proto.Message, fieldName string, values ...proto.Message) {
	msg := in.ProtoReflect()
	descriptor := msg.Descriptor()
	field := descriptor.Fields().ByName(protoreflect.Name(fieldName))

	// No values -- remove the field entirely
	if len(values) == 0 {
		msg.Clear(field)
		return
	}

	// For lists, append each one after clearing the previously stored value
	if field.IsList() {
		msg.Clear(field)
		list := msg.Mutable(field).List()
		for _, v := range values {
			list.Append(protoreflect.ValueOfMessage(v.ProtoReflect()))
		}
		return
	}

	// For single values on non-repeated fields, just set it.
	if len(values) == 1 {
		msg.Set(field, protoreflect.ValueOfMessage(values[0].ProtoReflect()))
		return
	}

	panic(
		fmt.Sprintf(
			"invalid use of Overwrite; non-repeated field '%v' used with '%v' values",
			fieldName,
			len(values),
		),
	)
}

// AppendList updates a field of the given name in-place for the specified message.
//
// This function will panic if the field is not a repeated-field.
func AppendList(in proto.Message, fieldName string, values ...proto.Message) {
	msg := in.ProtoReflect()
	descriptor := msg.Descriptor()
	field := descriptor.Fields().ByName(protoreflect.Name(fieldName))

	list := msg.Mutable(field).List()
	for _, v := range values {
		list.Append(protoreflect.ValueOfMessage(v.ProtoReflect()))
	}
}
