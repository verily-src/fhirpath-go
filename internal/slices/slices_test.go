package slices_test

import (
	"fmt"
	"strings"
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/internal/slices"
	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/testing/protocmp"
)

type StrongSlice[T any] []T

type arrTestCase[T constraints.Ordered] struct {
	name     string
	arr      StrongSlice[T]
	expected StrongSlice[T]
}

type boolTestCase[T any] struct {
	name     string
	arr      StrongSlice[T]
	expected bool
}

type includesTestCase[T any] struct {
	name     string
	arr      StrongSlice[T]
	target   T
	expected bool
}

type indexTestCase[T any] struct {
	name     string
	arr      StrongSlice[T]
	target   T
	expected int
}

type intTestCase[T any] struct {
	name     string
	arr      StrongSlice[T]
	expected int
}

type mapTestCase[T any, U any] struct {
	name     string
	arr      StrongSlice[T]
	expected []U
}

type stringTestCase[T any] struct {
	name      string
	arr       StrongSlice[T]
	delimiter string
	expected  string
}

var (
	firstCoding = &dtpb.CodeableConcept{
		Coding: []*dtpb.Coding{
			{
				System: &dtpb.Uri{Value: "test-system"},
				Code:   &dtpb.Code{Value: "1"},
			},
		},
	}
	secondCoding = &dtpb.CodeableConcept{
		Coding: []*dtpb.Coding{
			{
				System: &dtpb.Uri{Value: "test-system"},
				Code:   &dtpb.Code{Value: "2"},
			},
		},
	}
)
func TestIsIdentical(t *testing.T) {
	type s1 []int
	type s2 []int

	base := []int{1, 2, 3, 4}
	testCases := []struct {
		name string
		lhs  s1
		rhs  s2
		want bool
	}{
		{
			name: "Same reference",
			lhs:  s1(base),
			rhs:  s2(base),
			want: true,
		}, {
			name: "Different size",
			lhs:  s1(base),
			rhs:  append(s2(base), 5),
			want: false,
		}, {
			name: "Same size, empty",
			lhs:  nil,
			rhs:  nil,
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := slices.IsIdentical(tc.lhs, tc.rhs)

			if got != tc.want {
				t.Errorf("IsIdentical(%v): got %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

func TestAll(t *testing.T) {
	lengthFive := func(a string) bool { return len(a) == 5 }
	testCases := []boolTestCase[string]{
		{
			"empty slice",
			[]string{},
			true,
		}, {
			"no matching element",
			[]string{"cat", "dog"},
			false,
		}, {
			"some matching elements",
			[]string{"apple", "orange", "banana"},
			false,
		}, {
			"all matching elements",
			[]string{"daisy", "tulip"},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := slices.All(tc.arr, lengthFive)
			if got, want := result, tc.expected; got != want {
				t.Errorf("All(%s) want = %v, got = %v", tc.name, want, got)
			}
		})
	}
}

func TestAny(t *testing.T) {
	multOfFour := func(a int) bool { return a%4 == 0 }
	testCases := []boolTestCase[int]{
		{
			"empty slice",
			[]int{},
			false,
		}, {
			"no matching element",
			[]int{1, 2, 3},
			false,
		}, {
			"matching element",
			[]int{1, 2, 3, 4, 5},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := slices.Any(tc.arr, multOfFour)
			if got, want := result, tc.expected; got != want {
				t.Errorf("Any(%s) want = %v, got = %v", tc.name, want, got)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	lessThanPi := func(a float32) bool { return a < 3.14159 }
	testCases := []arrTestCase[float32]{
		{
			"empty slice",
			[]float32{},
			[]float32{},
		}, {
			"no matching elements",
			[]float32{4.1, 9.12},
			[]float32{},
		}, {
			"matching elements",
			[]float32{1.0, 2.718, 5.56},
			[]float32{1.0, 2.718},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := slices.Filter(tc.arr, lessThanPi)
			got, want := result, tc.expected
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Filter(%s) mismatch (-want, +got):\n%s", tc.name, diff)
			}
		})
	}
}

func TestCount(t *testing.T) {
	lessThanPi := func(a float32) bool { return a < 3.14159 }
	testCases := []intTestCase[float32]{
		{
			"empty slice",
			[]float32{},
			0,
		}, {
			"no matching elements",
			[]float32{4.1, 9.12},
			0,
		}, {
			"matching elements",
			[]float32{1.0, 2.718, 5.56, -3.0},
			3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := slices.Count(tc.arr, lessThanPi)
			got, want := result, tc.expected
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Count(%s) mismatch (-want, +got):\n%s", tc.name, diff)
			}

			// Count() should equal len(Filter())
			filtered := slices.Filter(tc.arr, lessThanPi)
			got = len(filtered)
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Count / len(Filter(%s)) mismatch (-want, +got):\n%s", tc.name, diff)
			}
		})
	}
}

func TestIncludes(t *testing.T) {
	testCases := []includesTestCase[int]{
		{
			"empty slice",
			[]int{},
			1,
			false,
		}, {
			"no matching elements",
			[]int{0, 1},
			2,
			false,
		}, {
			"found element",
			[]int{5, 7, 8},
			7,
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := slices.Includes(tc.arr, tc.target)
			if got, want := result, tc.expected; got != want {
				t.Errorf("Includes(%s) want = %v, got = %v", tc.name, want, got)
			}
		})
	}
}

func TestIndexOf(t *testing.T) {
	firstProto := proto.Message(firstCoding)
	secondProto := proto.Message(secondCoding)
	testCases := []indexTestCase[any]{
		{
			"empty slice",
			StrongSlice[any]{},
			firstCoding,
			-1,
		}, {
			"no matching elements",
			StrongSlice[any]{secondCoding},
			firstCoding,
			-1,
		}, {
			"found element - struct pointer",
			StrongSlice[any]{firstCoding, secondCoding},
			secondCoding,
			1,
		}, {
			"found element - proto pointer",
			StrongSlice[any]{firstProto, secondProto},
			secondProto,
			1,
		}, {
			"multiple matching elements",
			StrongSlice[any]{secondCoding, firstCoding, firstCoding},
			firstCoding,
			1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := slices.IndexOf(tc.arr, tc.target)
			if got, want := result, tc.expected; got != want {
				t.Errorf("IndexOf(%s) want = %v, got = %v", tc.name, want, got)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	testCases := []stringTestCase[uint]{
		{
			"empty slice",
			StrongSlice[uint]{},
			",",
			"",
		}, {
			"one element",
			StrongSlice[uint]{0},
			",",
			"0",
		}, {
			"multiple elements no delimiter",
			StrongSlice[uint]{0, 3},
			"",
			"03",
		}, {
			"multiple elements with delimiter",
			StrongSlice[uint]{0, 3, 7, 8},
			",",
			"0,3,7,8",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := slices.Join(tc.arr, tc.delimiter)
			if got, want := result, tc.expected; got != want {
				t.Errorf("Join(%s) want = %v, got = %v", tc.name, want, got)
			}
		})
	}
}

func TestMap(t *testing.T) {
	strLen := func(s string) int { return len(s) }
	testCases := []mapTestCase[string, int]{
		{
			"empty slice",
			StrongSlice[string]{},
			[]int{},
		}, {
			"one element",
			StrongSlice[string]{"four"},
			[]int{4},
		}, {
			"multiple elements",
			StrongSlice[string]{"one", "two", "three"},
			[]int{3, 3, 5},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := slices.Map(tc.arr, strLen)
			got, want := result, tc.expected
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Map(%s) mismatch (-want, +got):\n%s", tc.name, diff)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	testCases := []arrTestCase[int]{
		{
			"empty slice",
			StrongSlice[int]{},
			StrongSlice[int]{},
		}, {
			"one element",
			StrongSlice[int]{4},
			StrongSlice[int]{4},
		}, {
			"multiple elements",
			StrongSlice[int]{1, 2, 3, 4, 5},
			StrongSlice[int]{5, 4, 3, 2, 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			slices.Reverse(tc.arr)
			got, want := tc.arr, tc.expected
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Reverse(%s) mismatch (-want, +got):\n%s", tc.name, diff)
			}
		})
	}
}

func TestSort(t *testing.T) {
	testCases := []arrTestCase[int]{
		{
			"empty slice",
			StrongSlice[int]{},
			StrongSlice[int]{},
		}, {
			"one element",
			StrongSlice[int]{4},
			StrongSlice[int]{4},
		}, {
			"multiple elements",
			StrongSlice[int]{5, 4, 3, 2, 1},
			StrongSlice[int]{1, 2, 3, 4, 5},
		}, {
			"duplicate values",
			StrongSlice[int]{5, 4, 3, 2, 5, 1, 3},
			StrongSlice[int]{1, 2, 3, 3, 4, 5, 5},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			slices.Sort(tc.arr)
			got, want := tc.arr, tc.expected
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Sort(%s) mismatch (-want, +got):\n%s", tc.name, diff)
			}
		})
	}
}

type testType int

func (t testType) String() string {
	return fmt.Sprintf("%v", int(t))
}

type stringer interface {
	String() string
}

func TestConvert_Upcast_ReturnsConvertedType(t *testing.T) {
	input := StrongSlice[testType]{1, 2, 3}
	want := []stringer{input[0], input[1], input[2]}

	got, err := slices.Convert[stringer](input)

	if err != nil {
		t.Fatalf("Convert: got unexpected err %v", err)
	}
	if !cmp.Equal(got, want) {
		t.Errorf("Convert: got %v, want %v", got, want)
	}
}

func TestConvert_Downcast_ReturnsConvertedType(t *testing.T) {
	input := StrongSlice[testType]{1, 2, 3}
	want := []any{input[0], input[1], input[2]}

	got, err := slices.Convert[any](input)

	if err != nil {
		t.Fatalf("Convert: got unexpected err %v", err)
	}
	if !cmp.Equal(got, want) {
		t.Errorf("Convert: got %v, want %v", got, want)
	}
}

func TestConvert_Identity_ReturnsSelf(t *testing.T) {
	want := []int{1, 2, 3}

	got, err := slices.Convert[int](want)

	if err != nil {
		t.Fatalf("Convert: got unexpected err %v", err)
	}
	if !cmp.Equal(got, want) {
		t.Errorf("Convert: got %v, want %v", got, want)
	}
}

func TestConvert_InvalidCast_ReturnsErr(t *testing.T) {
	want := StrongSlice[int]{1, 2, 3}

	_, err := slices.Convert[string](want)

	if err == nil {
		t.Fatalf("Convert: expected err, got nil")
	}
}

func TestMustConvert_OnSuccess_ReturnsConversion(t *testing.T) {
	want := []int{4, 5, 6}
	input := StrongSlice[any]{4, 5, 6}

	result := slices.MustConvert[int](input)

	if got := result; !cmp.Equal(got, want) {
		t.Errorf("MustConvert: got %v, want %v", got, want)
	}
}

func TestMustConvert_OnBadConversion_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	input := []any{4, 5, 6}

	slices.MustConvert[string](input)
}

func TestTransform(t *testing.T) {
	testCases := []struct {
		name      string
		slice     StrongSlice[string]
		transform func(string) string
		want      StrongSlice[string]
	}{
		{
			name:      "append",
			slice:     []string{"hello", "world"},
			transform: func(s string) string { return "_" + s + "_" },
			want:      []string{"_hello_", "_world_"},
		}, {
			name:      "trim",
			slice:     []string{" \thello ", "\t world\n"},
			transform: strings.TrimSpace,
			want:      []string{"hello", "world"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			slices.Transform(tc.slice, tc.transform)

			got, want := tc.slice, tc.want
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("Transform(%s) mismatch (-want, +got):\n%s", tc.name, diff)
			}
		})
	}
}

func TestIsUnique_SimpleType(t *testing.T) {
	testCases := []struct {
		name  string
		slice StrongSlice[string]
		want  bool
	}{
		{
			name:  "unique",
			slice: []string{"hello", "world"},
			want:  true,
		}, {
			name:  "duplicate",
			slice: []string{"hello", "hello", "world"},
			want:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			unique := slices.IsUnique(tc.slice)

			if got, want := unique, tc.want; got != want {
				t.Errorf("IsUnique(%s): got %v, want %v", tc.name, got, want)
			}
		})
	}
}

func TestIsUnique_Proto(t *testing.T) {
	firstProto := proto.Message(firstCoding)
	secondProto := proto.Message(secondCoding)

	testCases := []struct {
		name  string
		slice StrongSlice[protoreflect.ProtoMessage]
		want  bool
	}{
		{
			name:  "unique",
			slice: []protoreflect.ProtoMessage{firstProto, secondProto},
			want:  true,
		}, {
			name:  "duplicate",
			slice: []protoreflect.ProtoMessage{secondProto, firstProto, secondProto},
			want:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			unique := slices.IsUnique(tc.slice)

			if got, want := unique, tc.want; got != want {
				t.Errorf("IsUnique_Proto(%s): got %v, want %v", tc.name, got, want)
			}
		})
	}
}

func TestUnique_SimpleType(t *testing.T) {
	testCases := []struct {
		name  string
		slice StrongSlice[string]
		want  StrongSlice[string]
	}{
		{
			name:  "already unique",
			slice: []string{"hello", "world"},
			want:  []string{"hello", "world"},
		}, {
			name:  "duplicates removed",
			slice: []string{"hello", "hello", "world"},
			want:  []string{"hello", "world"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uniqueSlice := slices.Unique(tc.slice)

			got, want := uniqueSlice, tc.want
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("Unique(%s) mismatch (-want, +got):\n%s", tc.name, diff)
			}
		})
	}
}

func TestUnique_Proto(t *testing.T) {
	firstProto := proto.Message(firstCoding)
	secondProto := proto.Message(secondCoding)

	testCases := []struct {
		name  string
		slice StrongSlice[protoreflect.ProtoMessage]
		want  StrongSlice[protoreflect.ProtoMessage]
	}{
		{
			name:  "already unique",
			slice: []protoreflect.ProtoMessage{firstProto, secondProto},
			want:  []protoreflect.ProtoMessage{firstProto, secondProto},
		}, {
			name:  "duplicates removed",
			slice: []protoreflect.ProtoMessage{firstProto, firstProto, secondProto},
			want:  []protoreflect.ProtoMessage{firstProto, secondProto},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uniqueSlice := slices.Unique(tc.slice)

			got, want := uniqueSlice, tc.want
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Unique(%s) mismatch (-want, +got):\n%s", tc.name, diff)
			}
		})
	}
}

func TestChunk(t *testing.T) {
	testCases := []struct {
		name string
		vals []int
		size int
		want [][]int
	}{
		{
			name: "less than 1 chunk size",
			vals: []int{1, 2, 3, 4},
			size: 5,
			want: [][]int{{1, 2, 3, 4}},
		},
		{
			name: "exactly 1 chunk size",
			vals: []int{1, 2, 3, 4, 5},
			size: 5,
			want: [][]int{{1, 2, 3, 4, 5}},
		},
		{
			name: "over 1 chunk size",
			vals: []int{1, 2, 3, 4, 5, 6},
			size: 5,
			want: [][]int{{1, 2, 3, 4, 5}, {6}},
		},
		{
			name: "4 chunks",
			vals: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			size: 3,
			want: [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := slices.Chunk(tc.vals, tc.size)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Chunk(%s) mismatch (-want, +got):\n%s", tc.name, diff)
			}
		})
	}
}
