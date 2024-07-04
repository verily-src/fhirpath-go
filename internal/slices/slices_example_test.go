package slices_test

import (
	"fmt"
	"strings"

	"github.com/verily-src/fhirpath-go/internal/slices"
)

func ExampleAll() {
	greaterThanFive := func(i int) bool { return i > 5 }
	arr := []int{1, 2, 3}
	arr2 := []int{7, 8, 9}

	all := slices.All(arr, greaterThanFive)
	fmt.Printf("arr all > 5: %v\n", all)

	all = slices.All(arr2, greaterThanFive)
	fmt.Printf("arr2 all > 5: %v\n", all)

	// Output:
	// arr all > 5: false
	// arr2 all > 5: true
}

func ExampleAny() {
	isEven := func(i int) bool { return i%2 == 0 }
	oddArr := []int{1, 3, 5, 7}
	evenArr := []int{1, 2, 9}

	any := slices.Any(oddArr, isEven)
	fmt.Printf("oddArr any even: %v\n", any)

	any = slices.Any(evenArr, isEven)
	fmt.Printf("evenArr any even: %v\n", any)

	// Output:
	// oddArr any even: false
	// evenArr any even: true
}

func ExampleFilter() {
	isNegative := func(i int) bool { return i < 0 }
	arr := []int{1, 3, -9, -5, 6}

	filtered := slices.Filter(arr, isNegative)
	fmt.Printf("filtered arr: %v\n", filtered)

	// Output:
	// filtered arr: [-9 -5]
}

func ExampleCount() {
	isNegative := func(i int) bool { return i < 0 }
	arr := []int{1, 3, -9, -5, 6}

	countNeg := slices.Count(arr, isNegative)
	fmt.Printf("count of negative nums: %v\n", countNeg)

	// Output:
	// count of negative nums: 2
}

func ExampleIncludes() {
	arr := []string{"a", "c", "e"}

	includesB := slices.Includes(arr, "b")
	fmt.Printf("arr includes b: %v\n", includesB)

	includesA := slices.Includes(arr, "a")
	fmt.Printf("arr includes a: %v\n", includesA)

	// Output:
	// arr includes b: false
	// arr includes a: true
}

func ExampleFind() {
	isEven := func(i int) bool { return i%2 == 0 }
	arrWithOdd := []int{1, 3, 5, 7}
	arrWithEven := []int{1, 3, 5, 6, 7}

	resultOdd, foundOdd := slices.Find(arrWithOdd, isEven)
	fmt.Printf("found odd: %v\n", foundOdd)
	fmt.Printf("result odd: %v\n", resultOdd)

	resultEven, foundEven := slices.Find(arrWithEven, isEven)
	fmt.Printf("found even: %v\n", foundEven)
	fmt.Printf("result even: %v\n", resultEven)

	// Output:
	// found odd: false
	// result odd: 0
	// found even: true
	// result even: 6
}

func ExampleIndexOf() {
	letters := []string{"w", "o", "o", "a", "h"}

	indexM := slices.IndexOf(letters, "m")
	fmt.Printf("index of m: %v\n", indexM)

	indexO := slices.IndexOf(letters, "o")
	fmt.Printf("index of o: %v\n", indexO)

	// Output:
	// index of m: -1
	// index of o: 1
}

func ExampleJoin() {
	arr := []int{0, 0}

	joined := slices.Join(arr, ".")
	fmt.Printf("joined string: %v\n", joined)

	// Output:
	// joined string: 0.0
}

func ExampleMap() {
	arr := []int{1, 2, 3}
	mapper := func(i int) string { return fmt.Sprintf("%v_%v", i, i) }

	mapped := slices.Map(arr, mapper)
	fmt.Printf("mapped arr: %v\n", mapped)

	// Output:
	// mapped arr: [1_1 2_2 3_3]
}

func ExampleReverse() {
	arr := []int{1, 2, 3, 3, 4}

	slices.Reverse(arr)
	fmt.Printf("reversed arr: %v\n", arr)

	// Output:
	// reversed arr: [4 3 3 2 1]
}

func ExampleSort() {
	arr := []int{1, -2, 3, -4, 5}

	slices.Sort(arr)
	fmt.Printf("sorted arr: %v\n", arr)

	// Output:
	// sorted arr: [-4 -2 1 3 5]
}

func ExampleConvert_good_conversion() {
	arr := []any{1, 2, 3}

	intArr, err := slices.Convert[int](arr)
	if err != nil {
		fmt.Printf("Unable to convert slice!")
	} else {
		fmt.Printf("int array: %v", intArr)
	}

	// Output:
	// int array: [1 2 3]
}

func ExampleConvert_bad_conversion() {
	arr := []any{1, 2, 3}

	intArr, err := slices.Convert[string](arr)
	if err != nil {
		fmt.Printf("Unable to convert slice!")
	} else {
		fmt.Printf("int array: %v", intArr)
	}

	// Output:
	// Unable to convert slice!
}

func ExampleMustConvert() {
	// Easy mechanism to cast arrays to interfaces and back
	arr := []any{1, 2, 3}

	intArr := slices.MustConvert[int](arr)
	fmt.Printf("int array: %v", intArr)

	// Output:
	// int array: [1 2 3]
}

func ExampleTransform() {
	arr := []string{"\t hello\n", "\t world\n"}

	slices.Transform(arr, strings.TrimSpace)
	fmt.Printf("transformed arr: %v\n", arr)

	// Output:
	// transformed arr: [hello world]
}

func ExampleIsUnique() {
	strArr := []string{"hello", "world"}
	intArr := []int{1, 2, 3, 2, 1}

	strUnique := slices.IsUnique(strArr)
	fmt.Printf("str arr unique: %v\n", strUnique)

	intUnique := slices.IsUnique(intArr)
	fmt.Printf("int arr unique: %v\n", intUnique)

	// Output:
	// str arr unique: true
	// int arr unique: false
}

func ExampleUnique() {
	arr := []int{1, 2, 3, 2, 1}

	uniqueArr := slices.Unique(arr)
	fmt.Printf("unique arr: %v\n", uniqueArr)

	// Output:
	// unique arr: [1 2 3]
}

func ExampleChunk() {
	strArr := []string{"a", "b", "c", "d", "e"}
	strChunks := slices.Chunk(strArr, 5)
	fmt.Printf("str arr chunks: %v\n", strChunks)

	intArr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	intChunks := slices.Chunk(intArr, 3)
	fmt.Printf("int arr chunks: %v\n", intChunks)

	// Output:
	// str arr chunks: [[a b c d e]]
	// int arr chunks: [[1 2 3] [4 5 6] [7 8 9] [10]]

}
