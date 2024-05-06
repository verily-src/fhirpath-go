package narrow_test

import (
	"fmt"
	"math"

	"github.com/verily-src/fhirpath-go/internal/narrow"
)

func ExampleToInteger_narrows() {
	from := -1

	val, ok := narrow.ToInteger[uint8](from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Unable to convert to uint8!
}

func ExampleToInteger_no_narrowing() {
	from := uint32(42)

	val, ok := narrow.ToInteger[uint8](from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Converted 42 to uint8!
}

func ExampleToInt8_narrows() {
	from := 10_000

	val, ok := narrow.ToInt8(from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Unable to convert to int8!
}

func ExampleToInt8_no_narrowing() {
	from := 42

	val, ok := narrow.ToInt8(from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Converted 42 to int8!
}

func ExampleToUint8_narrows() {
	from := 10_000

	val, ok := narrow.ToUint8(from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Unable to convert to uint8!
}

func ExampleToUint8_no_narrowing() {
	from := 42

	val, ok := narrow.ToUint8(from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Converted 42 to uint8!
}

func ExampleToInt16_narrows() {
	from := math.MaxInt32

	val, ok := narrow.ToInt16(from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Unable to convert to int16!
}

func ExampleToInt16_no_narrowing() {
	from := 42

	val, ok := narrow.ToInt16(from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Converted 42 to int16!
}

func ExampleToUint16_narrows() {
	from := math.MaxInt32

	val, ok := narrow.ToUint16(from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Unable to convert to uint16!
}

func ExampleToUint16_no_narrowing() {
	from := 42

	val, ok := narrow.ToUint16(from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Converted 42 to uint16!
}

func ExampleToInt32_narrows() {
	from := math.MaxInt64

	val, ok := narrow.ToInt32(from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Unable to convert to int32!
}

func ExampleToInt32_no_narrowing() {
	from := 42

	val, ok := narrow.ToInt32(from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Converted 42 to int32!
}

func ExampleToUint32_narrows() {
	from := math.MaxInt64

	val, ok := narrow.ToUint32(from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Unable to convert to uint32!
}

func ExampleToUint32_no_narrowing() {
	from := 42

	val, ok := narrow.ToUint32(from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Converted 42 to uint32!
}

func ExampleToInt64_no_narrowing() {
	from := 42

	val, ok := narrow.ToInt64(from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Converted 42 to int64!
}

func ExampleToUint64_narrows() {
	from := -1

	val, ok := narrow.ToUint64(from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Unable to convert to uint64!
}

func ExampleToUint64_no_narrowing() {
	from := 42

	val, ok := narrow.ToUint64(from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Converted 42 to uint64!
}

func ExampleToInt_no_narrowing() {
	from := 42

	val, ok := narrow.ToInt32(from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Converted 42 to int32!
}

func ExampleToUint_narrows() {
	from := -1

	val, ok := narrow.ToUint(from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Unable to convert to uint!
}

func ExampleToUint_no_narrowing() {
	from := 42

	val, ok := narrow.ToUint(from)
	if !ok {
		fmt.Printf("Unable to convert to %T!", val)
	} else {
		fmt.Printf("Converted %v to %T!", from, val)
	}
	// Output:
	// Converted 42 to uint!
}
