package protofields

import "google.golang.org/protobuf/proto"

// DescriptorName gets the type name of a proto Message. If value is nil, this
// returns an empty string.
func DescriptorName(value proto.Message) string {
	if value == nil {
		return ""
	}
	return string(value.ProtoReflect().Descriptor().Name())
}
