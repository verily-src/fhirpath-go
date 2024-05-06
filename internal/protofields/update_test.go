package protofields_test

import (
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/verily-src/fhirpath-go/internal/protofields"
	"google.golang.org/protobuf/proto"
)

func TestOverwrite(t *testing.T) {
	value := &dtpb.String{
		Value: "hello world",
	}
	testCases := []struct {
		name   string
		field  string
		values []proto.Message
		input  proto.Message
		want   proto.Message
	}{
		{
			name:   "Solo field",
			field:  "text",
			values: []proto.Message{value},
			input:  &dtpb.HumanName{},
			want: &dtpb.HumanName{
				Text: value,
			},
		}, {
			name:   "Solo field no input",
			field:  "text",
			values: []proto.Message{},
			input: &dtpb.HumanName{
				Text: value,
			},
			want: &dtpb.HumanName{},
		}, {
			name:   "Repeated field with single input",
			field:  "prefix",
			values: []proto.Message{value},
			input:  &dtpb.HumanName{},
			want: &dtpb.HumanName{
				Prefix: []*dtpb.String{value},
			},
		}, {
			name:   "Repeated field with multiple inputs",
			field:  "prefix",
			values: []proto.Message{value, value},
			input:  &dtpb.HumanName{},
			want: &dtpb.HumanName{
				Prefix: []*dtpb.String{value, value},
			},
		}, {
			name:   "Repeated field no input",
			field:  "prefix",
			values: []proto.Message{},
			input: &dtpb.HumanName{
				Prefix: []*dtpb.String{value, value},
			},
			want: &dtpb.HumanName{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			protofields.Overwrite(tc.input, tc.field, tc.values...)

			if got, want := tc.input, tc.want; !proto.Equal(got, want) {
				t.Errorf("Overwrite(%v): got %v, want %v", tc.name, got, want)
			}
		})
	}
}

func TestOverwrite_WrongCardinality_Panics(t *testing.T) {
	defer func() { _ = recover() }()
	value := &dtpb.String{
		Value: "hello world",
	}
	name := &dtpb.HumanName{}

	protofields.Overwrite(name, "text", value, value)

	t.Errorf("Overwrite: expected panic")
}

func TestAppendList(t *testing.T) {
	toAppend := &dtpb.String{
		Value: "hello world",
	}
	value := &dtpb.String{
		Value: "another string",
	}
	testCases := []struct {
		name  string
		field string
		input proto.Message
		want  proto.Message
	}{
		{
			name:  "Repeated field with no inputs",
			field: "prefix",
			input: &dtpb.HumanName{},
			want: &dtpb.HumanName{
				Prefix: []*dtpb.String{toAppend},
			},
		}, {
			name:  "Repeated field with 1 input",
			field: "prefix",
			input: &dtpb.HumanName{
				Prefix: []*dtpb.String{value},
			},
			want: &dtpb.HumanName{
				Prefix: []*dtpb.String{value, toAppend},
			},
		}, {
			name:  "Repeated field with multiple inputs",
			field: "prefix",
			input: &dtpb.HumanName{
				Prefix: []*dtpb.String{value, value},
			},
			want: &dtpb.HumanName{
				Prefix: []*dtpb.String{value, value, toAppend},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			protofields.AppendList(tc.input, tc.field, toAppend)

			if got, want := tc.input, tc.want; !proto.Equal(got, want) {
				t.Errorf("AppendList(%v): got %v, want %v", tc.name, got, want)
			}
		})
	}
}
