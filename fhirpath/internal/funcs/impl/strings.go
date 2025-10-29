package impl

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

var (
	ErrInvalidRegex = errors.New("invalid regex")
)

// StartsWith returns true if the input string starts with the given prefix.
func StartsWith(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Validate single string input
	if length := len(input); length > 1 {
		return nil, fmt.Errorf("%w: input has length %v, expected 1", ErrWrongArity, length)
	} else if length == 0 {
		return system.Collection{}, nil
	}
	fullString, err := input.ToString()
	if err != nil {
		return nil, err
	}

	// Validate single string argument
	if length := len(args); length != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, length)
	}
	output, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	} else if length := len(output); length != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, length)
	}
	prefix, err := output.ToString()
	if err != nil {
		return nil, err
	}

	result := system.Boolean(strings.HasPrefix(string(fullString), string(prefix)))
	return system.Collection{result}, nil
}

// EndsWith returns true if the input string ends with the given prefix.
func EndsWith(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Validate single string input
	if length := len(input); length > 1 {
		return nil, fmt.Errorf("%w: input has length %v, expected 1", ErrWrongArity, length)
	} else if length == 0 {
		return system.Collection{}, nil
	}
	fullString, err := input.ToString()
	if err != nil {
		return nil, err
	}

	// Validate single string argument
	if length := len(args); length != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, length)
	}
	output, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	} else if length := len(output); length != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, length)
	}
	suffix, err := output.ToString()
	if err != nil {
		return nil, err
	}

	result := system.Boolean(strings.HasSuffix(string(fullString), string(suffix)))
	return system.Collection{result}, nil
}

// Length returns the length of the input string.
func Length(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Validate single string input
	if length := len(input); length > 1 {
		return nil, fmt.Errorf("%w: input has length %v, expected 1", ErrWrongArity, length)
	} else if length == 0 {
		return system.Collection{}, nil
	}
	fullString, err := input.ToString()
	if err != nil {
		return nil, err
	}

	if length := len(args); length != 0 {
		return nil, fmt.Errorf("%w, received %v arguments, expected 0", ErrWrongArity, length)
	}

	result := system.Integer(len(fullString))
	return system.Collection{result}, nil
}

// Upper returns the input string with all characters converted to upper case.
func Upper(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Validate single string input
	if length := len(input); length > 1 {
		return nil, fmt.Errorf("%w: input has length %v, expected 1", ErrWrongArity, length)
	} else if length == 0 {
		return system.Collection{}, nil
	}
	fullString, err := input.ToString()
	if err != nil {
		return nil, err
	}

	if length := len(args); length != 0 {
		return nil, fmt.Errorf("%w, received %v arguments, expected 0", ErrWrongArity, length)
	}

	result := system.String(strings.ToUpper(fullString))
	return system.Collection{result}, nil
}

// Lower returns the input string with all characters converted to lower case.
func Lower(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Validate single string input
	if length := len(input); length > 1 {
		return nil, fmt.Errorf("%w: input has length %v, expected 1", ErrWrongArity, length)
	} else if length == 0 {
		return system.Collection{}, nil
	}
	fullString, err := input.ToString()
	if err != nil {
		return nil, err
	}

	if length := len(args); length != 0 {
		return nil, fmt.Errorf("%w, received %v arguments, expected 0", ErrWrongArity, length)
	}

	result := system.String(strings.ToLower(fullString))
	return system.Collection{result}, nil
}

// Contains returns true if the input string contains the given substring.
func Contains(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Validate single string input
	if length := len(input); length > 1 {
		return nil, fmt.Errorf("%w: input has length %v, expected 1", ErrWrongArity, length)
	} else if length == 0 {
		return system.Collection{}, nil
	}
	fullString, err := input.ToString()
	if err != nil {
		return nil, err
	}

	// Validate single string argument
	if length := len(args); length != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, length)
	}
	output, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	} else if length := len(output); length != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, length)
	}
	substring, err := output.ToString()
	if err != nil {
		return nil, err
	}

	result := system.Boolean(strings.Contains(string(fullString), string(substring)))
	return system.Collection{result}, nil
}

// ToChars returns the list of characters in the input string.
func ToChars(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Validate single string input
	if length := len(input); length > 1 {
		return nil, fmt.Errorf("%w: input has length %v, expected 1", ErrWrongArity, length)
	} else if length == 0 {
		return system.Collection{}, nil
	}
	fullString, err := input.ToString()
	if err != nil {
		return nil, err
	}

	if length := len(args); length != 0 {
		return nil, fmt.Errorf("%w, received %v arguments, expected 0", ErrWrongArity, length)
	}

	result := system.Collection{}
	for _, char := range strings.Split(fullString, "") {
		result = append(result, system.String(char))
	}
	return result, nil
}

// Substring returns the part of the string starting at position start (zero-based).
// If length is given, will return at most length number of characters from the input string.
func Substring(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Validate single string input
	if length := len(input); length > 1 {
		return nil, fmt.Errorf("%w: input has length %v, expected 1", ErrWrongArity, length)
	} else if length == 0 {
		return system.Collection{}, nil
	}
	fullString, err := input.ToString()
	if err != nil {
		return nil, err
	}

	argLength := len(args)
	if argLength < 1 || argLength > 2 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1 or 2", ErrWrongArity, argLength)
	}

	// Validate 1st integer argument (start)
	startOutput, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	} else if length := len(startOutput); length != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, length)
	}
	start, err := startOutput.ToInt32()
	if err != nil {
		return nil, err
	}
	if int(start) >= len(fullString) {
		return system.Collection{}, nil
	}

	// Validate optional 2nd integer argument (length).
	var substringLength int32 = -1
	if argLength == 2 {
		lengthOutput, err := args[1].Evaluate(ctx, input)
		if err != nil {
			return nil, err
		} else if length := len(lengthOutput); length != 1 {
			return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, length)
		}
		substringLength, err = lengthOutput.ToInt32()
		if err != nil {
			return nil, err
		}
	}

	var result system.String
	if substringLength > -1 && int(start+substringLength) < len(fullString) {
		// Substring will not go out of bounds
		result = system.String(fullString[start : start+substringLength])
	} else {
		result = system.String(fullString[start:])
	}
	return system.Collection{result}, nil
}

// IndexOf returns the 0-based index of the first position in which the
// substring is found in the input string, or -1 if it is not found.
func IndexOf(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Validate single string input
	if length := len(input); length > 1 {
		return nil, fmt.Errorf("%w: input has length %v, expected 1", ErrWrongArity, length)
	} else if length == 0 {
		return system.Collection{}, nil
	}
	fullString, err := input.ToString()
	if err != nil {
		return nil, err
	}

	// Validate single string argument
	if length := len(args); length != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, length)
	}
	output, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	} else if length := len(output); length == 0 {
		// Return empty for empty argument
		return system.Collection{}, nil
	} else if length > 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, length)
	}
	substring, err := output.ToString()
	if err != nil {
		return nil, err
	}

	result := system.Integer(strings.Index(fullString, substring))
	return system.Collection{result}, nil
}

// Matches returns true when the value matches the given regular expression.
func Matches(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Validate single string input
	if length := len(input); length > 1 {
		return nil, fmt.Errorf("%w: input has length %v, expected 1", ErrWrongArity, length)
	} else if length == 0 {
		return system.Collection{}, nil
	}
	fullString, err := input.ToString()
	if err != nil {
		return nil, err
	}

	// Validate single string argument
	if length := len(args); length != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, length)
	}
	output, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	} else if length := len(output); length == 0 {
		return system.Collection{}, nil
	} else if length != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, length)
	}
	regexString, err := output.ToString()
	if err != nil {
		return nil, err
	}
	re, err := regexp.Compile(regexString)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidRegex, regexString)
	}

	result := system.Boolean(re.Match([]byte(fullString)))
	return system.Collection{result}, nil
}

// Replace returns the input string with all instances of pattern replaced with substitution.
func Replace(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Validate single string input
	if length := len(input); length > 1 {
		return nil, fmt.Errorf("%w: input has length %v, expected 1", ErrWrongArity, length)
	} else if length == 0 {
		return system.Collection{}, nil
	}
	fullString, err := input.ToString()
	if err != nil {
		return nil, err
	}

	if length := len(args); length != 2 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 2", ErrWrongArity, length)
	}

	// Validate 1st string argument (pattern)
	patternOutput, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	} else if length := len(patternOutput); length == 0 {
		// Empty arg
		return system.Collection{}, nil
	} else if length > 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, length)
	}
	pattern, err := patternOutput.ToString()
	if err != nil {
		return nil, err
	}

	// Validate 2nd string argument (substitution)
	subOutput, err := args[1].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	} else if length := len(subOutput); length == 0 {
		// Empty arg
		return system.Collection{}, nil
	} else if length > 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, length)
	}
	substitution, err := subOutput.ToString()
	if err != nil {
		return nil, err
	}

	result := system.String(strings.ReplaceAll(fullString, pattern, substitution))
	return system.Collection{result}, nil
}

// ReplaceMatches matches the input using the regular expression in
// regex and replaces each match with the substitution string. The
// substitution may refer to identified match groups in the regular expression.
func ReplaceMatches(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Validate single string input
	if length := len(input); length > 1 {
		return nil, fmt.Errorf("%w: input has length %v, expected 1", ErrWrongArity, length)
	} else if length == 0 {
		return system.Collection{}, nil
	}
	fullString, err := input.ToString()
	if err != nil {
		return nil, err
	}

	if length := len(args); length != 2 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, length)
	}

	// Validate 1st string argument (regex)
	regexOutput, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	} else if length := len(regexOutput); length == 0 {
		return system.Collection{}, nil
	} else if length > 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, length)
	}
	regexString, err := regexOutput.ToString()
	if err != nil {
		return nil, err
	}
	re, err := regexp.Compile(regexString)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidRegex, regexString)
	}

	// Validate 2nd string argument (substitution)
	subOutput, err := args[1].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	} else if length := len(subOutput); length == 0 {
		return system.Collection{}, nil
	} else if length > 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, length)
	}
	substitution, err := subOutput.ToString()
	if err != nil {
		return nil, err
	}

	result := system.String(re.ReplaceAllString(fullString, substitution))
	return system.Collection{result}, nil
}

// Join takes a collection of strings and joins them into a single string,
// optionally using the given separator. If no separator is specified,
// the strings are directly concatenated.
func Join(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	if length := len(args); length > 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected at most 1", ErrWrongArity, length)
	}
	if length := len(input); length == 0 {
		return system.Collection{}, nil
	}
	delimiter := ""
	if len(args) == 1 {
		argValues, err := args[0].Evaluate(ctx, input)
		if err != nil {
			return nil, err
		}
		delimiter, err = argValues.ToString()
		if err != nil {
			return nil, err
		}
	}
	var strs []string
	for _, item := range input {
		value, err := system.From(item)
		if err != nil {
			return nil, err
		}
		str, ok := value.(system.String)
		if !ok {
			return nil, fmt.Errorf("expected string, got %T", item)
		}
		strs = append(strs, string(str))
	}
	return system.Collection{system.String(strings.Join(strs, delimiter))}, nil
}

// Split splits the input string by the specified separator and returns a collection of strings.
//
// Syntax: input.split(separator)
//
// Examples:
//
//	empty.split(',')               -> []                   (empty input produces empty collection)
//	'a,b,c'.split(empty)           -> []                   (empty separator argument produces empty collection)
//	''.split(',')                  -> ['']                 (empty string produces single empty element)
//	'a,b,c'.split(',')             -> ['a', 'b', 'c']      (single-char separator)
//	'a::b::c'.split('::')          -> ['a', 'b', 'c']      (multi-char separator)
//	patient.reference.split('/')   -> ['Patient', '12345'] (e.g. 'Patient/12345')
//	'abc'.split('')                -> ['a', 'b', 'c']      (empty separator splits to chars)
//	'abc'.split(',')               -> ['abc']              (no split occurs)
func Split(ctx *expr.Context, input system.Collection, args ...expr.Expression) (system.Collection, error) {
	// Validate single argument (separator)
	if len(args) != 1 {
		return nil, fmt.Errorf("%w: received %v arguments, expected 1", ErrWrongArity, len(args))
	}

	if input.IsEmpty() {
		// Empty input returns empty collection
		return system.Collection{}, nil
	}

	// Validate single string input
	if !input.IsSingleton() {
		return nil, fmt.Errorf("%w: input has %d elements", ErrNotSingleton, len(input))
	}
	fullString, err := input.ToString()
	if err != nil {
		return nil, err
	}

	// Evaluate the separator argument
	argValues, err := args[0].Evaluate(ctx, input)
	if err != nil {
		return nil, err
	}
	if argValues.IsEmpty() {
		// Empty separator argument returns empty collection
		return system.Collection{}, nil
	}
	if !argValues.IsSingleton() {
		return nil, fmt.Errorf("%w: separator has %d elements", ErrNotSingleton, len(argValues))
	}
	separator, err := argValues.ToString()
	if err != nil {
		return nil, err
	}

	// Split the string and return collection of strings
	parts := strings.Split(fullString, separator)
	result := make(system.Collection, len(parts))
	for i, part := range parts {
		result[i] = system.String(part)
	}
	return result, nil
}
