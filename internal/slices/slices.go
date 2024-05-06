// Package slices provides helpful functions
// for searching, sorting and manipulating slices.
package slices

import (
	"fmt"
	"reflect"
	"sort"

	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/proto"
)

// IsIdentical checks if two slices refer to the same underlying slice object.
func IsIdentical[S2 ~[]T, S1 ~[]T, T any](lhs S1, rhs S2) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	if len(lhs) == 0 {
		return true
	}
	return &lhs[0] == &rhs[0]
}

// All returns true if every element in a slice satisfies the
// given condition. Defaults to true for an empty slice.
func All[S ~[]T, T any](vals S, comp func(T) bool) bool {
	for _, val := range vals {
		if !comp(val) {
			return false
		}
	}
	return true
}

// Any returns true if at least one element in a slice satisfies
// the given condition. Defaults to false for an empty slice.
func Any[S ~[]T, T any](vals S, comp func(T) bool) bool {
	for _, val := range vals {
		if comp(val) {
			return true
		}
	}
	return false
}

// Filter returns a subset of the original slice, consisting of all the
// elements in the original slice which satisfy the given condition.
func Filter[S ~[]T, T any](vals S, comp func(T) bool) S {
	matches := make(S, 0)
	for _, val := range vals {
		if comp(val) {
			matches = append(matches, val)
		}
	}
	return matches
}

// Count returns the number of elements in a slice that match the given
// condition.
func Count[S ~[]T, T any](vals S, comp func(T) bool) int {
	matches := 0
	for _, val := range vals {
		if comp(val) {
			matches += 1
		}
	}
	return matches
}

// Includes returns true if a slice contains the target value.
func Includes[S ~[]T, T any](vals S, target T) bool {
	return IndexOf(vals, target) > -1
}

// IndexOf returns the index of the target value in a slice.
// Returns the first index found, and -1 if the value is not found.
func IndexOf[S ~[]T, T any](vals S, target T) int {
	_, isProto := any(target).(proto.Message)
	for index, val := range vals {
		if isProto && proto.Equal(any(val).(proto.Message), any(target).(proto.Message)) {
			return index
		} else if reflect.DeepEqual(val, target) {
			return index
		}
	}
	return -1
}

// Join returns a string consisting of all elements in a slice,
// separated by the given delimiter.
func Join[S ~[]T, T any](vals S, delimiter string) string {
	if len(vals) == 0 {
		return ""
	}
	combined := fmt.Sprintf("%v", vals[0])
	for _, val := range vals[1:] {
		combined += fmt.Sprintf("%s%v", delimiter, val)
	}
	return combined
}

// Map returns a new slice with the elements of the original
// slice transformed according to the provided function.
func Map[T any, U any](vals []T, mapper func(T) U) []U {
	mapped := make([]U, 0, len(vals))
	for _, val := range vals {
		mapped = append(mapped, mapper(val))
	}
	return mapped
}

// Reverse performs an in-place reversal of the elements in a slice.
// Performance testing has not been done to compare to alternatives.
func Reverse[S ~[]T, T any](t S) {
	sort.SliceStable(t, func(i, j int) bool {
		return i > j
	})
}

// Sort performs an in-place sort of a slice.
// Performance testing has not been done to compare to alternatives.
func Sort[S ~[]T, T constraints.Ordered](vals S) {
	sort.SliceStable(vals, func(i, j int) bool {
		return vals[i] < vals[j]
	})
}

// Convert returns an array of objects of type To from an array of objects of
// type From. If any conversion fails, this will return an error.
//
// This function, along with the non-failing `MustConvert` equivalent, are
// useful for converting one array type []T to an array of another type []U,
// which requires manual iteration in Go. This provides a more convenient
// mechanism for casting between arrays of interfaces and concrete types,
// provided all elements can safely be casted to.
//
// Note: this function only works with language-level type-casts (e.g. `t.(U)`).
// For converting between concrete struct types, use `Map`.
func Convert[To any, S ~[]From, From any](from S) ([]To, error) {
	result := make([]To, 0, len(from))
	for _, val := range from {
		if to, ok := any(val).(To); ok {
			result = append(result, to)
		} else {
			return nil, fmt.Errorf("slices.Convert[%T](from): unable to convert from %T", to, from)
		}
	}
	return result, nil
}

// MustConvert returns an array of To objects by converting every element in
// From to To. This will panic on failure to convert.
//
// This function, along with the failing `Convert` equivalent, are
// useful for converting one array type []T to an array of another type []U,
// which requires manual iteration in Go. This provides a more convenient
// mechanism for casting between arrays of interfaces and concrete types,
// provided all elements can safely be casted to.
//
// This function in particular should only be used if the 'To' type is guaranteed
// to always be convertible -- such as converting a concrete type to its base
// interface, or to any.
func MustConvert[To any, S ~[]From, From any](from S) []To {
	return Map(from, func(from From) To { return any(from).(To) })
}

// Transform performs an in-place transformation of a slice.
func Transform[S ~[]T, T any](vals S, transform func(T) T) {
	for i := range vals {
		vals[i] = transform(vals[i])
	}
}

// IsUnique returns true if all elements of a slice are unique.
func IsUnique[T comparable](vals []T) bool {
	seen := map[T]struct{}{}
	for _, val := range vals {
		if _, found := seen[val]; found {
			return false
		}
		seen[val] = struct{}{}
	}
	return true
}

// Unique returns a new slice of the unique elements in a given slice.
// It keeps the first instance of any duplicate values and ignores any
// subsequent instances of the value.
func Unique[S ~[]T, T comparable](vals S) S {
	set := map[T]struct{}{}
	uniqueSlice := make(S, 0)
	for _, val := range vals {
		if _, found := set[val]; !found {
			set[val] = struct{}{}
			uniqueSlice = append(uniqueSlice, val)
		}
	}
	return uniqueSlice
}

// Chunk divides the given slice into multiple new slices, each at most
// having lengths of the given size. The last chunk may have a smaller
// length if the length of vals is not evenly divisible by size.
func Chunk[T any](vals []T, size int) [][]T {
	var chunks [][]T
	for i := 0; i < len(vals); {
		chunkLen := size
		remainingLen := len(vals[i:])
		if remainingLen < size {
			chunkLen = remainingLen
		}
		chunk := make([]T, chunkLen)
		copy(chunk, vals[i:i+chunkLen])
		chunks = append(chunks, chunk)
		i += chunkLen
	}
	return chunks
}
