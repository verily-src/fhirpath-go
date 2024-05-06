package system_test

import (
	"errors"
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

func TestParseString_ReplacesEscapeSequences(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "replaces single quote",
			input: `I\'m testing`,
			want:  `I'm testing`,
		},
		{
			name:  "replaces double quote",
			input: `"This is a quote: \""`,
			want:  `"This is a quote: ""`,
		},
		{
			name:  "replaces backtick",
			input: "let foo = \\`${bar}\\`",
			want:  "let foo = `${bar}`",
		},
		{
			name:  "replaces carriage return",
			input: `here's a return \r`,
			want:  "here's a return \r",
		},
		{
			name:  "replaces newline",
			input: `here's a newline \n`,
			want:  "here's a newline \n",
		},
		{
			name:  "replaces a tab",
			input: `\t indent`,
			want:  "\t indent",
		},
		{
			name:  "replaces a form feed",
			input: `\f`,
			want:  "\f",
		},
		{
			name:  "replaces backslashes w/ multiple escapes",
			input: `escape \n\ \p \\p`,
			want:  "escape \n p \\p",
		},
	}

	for _, tc := range testCases {
		str, err := system.ParseString(tc.input)

		if err != nil {
			t.Fatalf("ParseString(%s) raised error when not expected: %v", tc.input, err)
		}
		if got, want := string(str), tc.want; got != want {
			t.Errorf("ParseString(%s) returned mismatch: got %#v, want %#v", tc.input, got, want)
		}
	}
}

func TestParseBoolean_ReturnsBoolean(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		want      system.Boolean
		wantError bool
	}{
		{
			name:  "parse 'true'",
			input: "true",
			want:  system.Boolean(true),
		},
		{
			name:  "parse 'TRUE'",
			input: "TRUE",
			want:  system.Boolean(true),
		},
		{
			name:  "parse 't'",
			input: "t",
			want:  system.Boolean(true),
		},
		{
			name:  "parse 'T'",
			input: "T",
			want:  system.Boolean(true),
		},
		{
			name:  "parse 'yes'",
			input: "yes",
			want:  system.Boolean(true),
		},
		{
			name:  "parse 'YES'",
			input: "YES",
			want:  system.Boolean(true),
		},
		{
			name:  "parse 'y'",
			input: "y",
			want:  system.Boolean(true),
		},
		{
			name:  "parse 'Y'",
			input: "Y",
			want:  system.Boolean(true),
		},
		{
			name:  "parse '1'",
			input: "1",
			want:  system.Boolean(true),
		},
		{
			name:  "parse '1.0'",
			input: "1.0",
			want:  system.Boolean(true),
		},
		{
			name:  "parse 'false'",
			input: "false",
			want:  system.Boolean(false),
		},
		{
			name:  "parse 'FALSE'",
			input: "FALSE",
			want:  system.Boolean(false),
		},
		{
			name:  "parse 'f'",
			input: "f",
			want:  system.Boolean(false),
		},
		{
			name:  "parse 'F'",
			input: "F",
			want:  system.Boolean(false),
		},
		{
			name:  "parse 'no'",
			input: "no",
			want:  system.Boolean(false),
		},
		{
			name:  "parse 'NO'",
			input: "NO",
			want:  system.Boolean(false),
		},
		{
			name:  "parse 'n'",
			input: "n",
			want:  system.Boolean(false),
		},
		{
			name:  "parse 'N'",
			input: "N",
			want:  system.Boolean(false),
		},
		{
			name:  "parse '0'",
			input: "0",
			want:  system.Boolean(false),
		},
		{
			name:  "parse '0.0'",
			input: "0.0",
			want:  system.Boolean(false),
		},
		{
			name:      "parse '2'",
			input:     "2",
			want:      system.Boolean(false),
			wantError: true,
		},
		{
			name:      "parse '2.0'",
			input:     "2",
			want:      system.Boolean(false),
			wantError: true,
		},
		{
			name:      "parse '2.5'",
			input:     "2.5",
			want:      system.Boolean(false),
			wantError: true,
		},
		{
			name:      "parse '2.5 kg'",
			input:     "2.5 kg",
			want:      system.Boolean(false),
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := system.ParseBoolean(tc.input)
			if (err != nil) != tc.wantError {
				t.Errorf("ParseBoolean(%s) error = %v, wantErr %v", tc.input, err, tc.wantError)
				return
			}
			if got != tc.want {
				t.Errorf("ParseBoolean(%s) parsed incorrectly, got: %v, want %v", tc.input, got, tc.want)
				return
			}
		})
	}
}

func TestParseInteger_ReturnsInteger(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  system.Integer
	}{
		{
			name:  "positive edge",
			input: "2147483647",
			want:  system.Integer(2147483647),
		},
		{
			name:  "negative edge",
			input: "-2147483648",
			want:  system.Integer(-2147483648),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			i, err := system.ParseInteger(tc.input)

			if err != nil {
				t.Fatalf("ParseInteger(%s) returns unexpected error: %v", tc.input, err)
			}
			if got, want := i, tc.want; got != want {
				t.Errorf("ParseInteger(%s) parsed incorrectly: got %v, want %v", tc.input, got, want)
			}
		})
	}
}

func TestParseInteger_ReturnsError_IfOutOfRange(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "positive edge",
			input: "2147483648",
		},
		{
			name:  "negative edge",
			input: "-2147483649",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := system.ParseInteger(tc.input)

			if err == nil {
				t.Fatalf("ParseInteger(%s) doesn't return error when expected to", tc.input)
			}
		})
	}
}

func TestIntegerAdd(t *testing.T) {
	testCases := []struct {
		name    string
		left    system.Integer
		right   system.Integer
		want    system.Integer
		wantErr error
	}{
		{
			name:  "adds two integers",
			left:  2000,
			right: 4001,
			want:  6001,
		},
		{
			name:    "returns error when addition overflows positively",
			left:    math.MaxInt32,
			right:   12120,
			wantErr: system.ErrIntOverflow,
		},
		{
			name:    "returns error when addition overflows negatively",
			left:    math.MinInt32,
			right:   -1,
			wantErr: system.ErrIntOverflow,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.left.Add(tc.right)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Integer.Add returned unexpected error: got %v, want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Integer.Add returned unexpected result: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestIntegerSub(t *testing.T) {
	testCases := []struct {
		name    string
		left    system.Integer
		right   system.Integer
		want    system.Integer
		wantErr error
	}{
		{
			name:  "subtracts two integers",
			left:  2000,
			right: 4001,
			want:  -2001,
		},
		{
			name:    "returns error when subtraction overflows negatively",
			left:    math.MinInt32,
			right:   1,
			wantErr: system.ErrIntOverflow,
		},
		{
			name:    "returns error when subtraction overflows positively",
			left:    math.MaxInt32,
			right:   -1,
			wantErr: system.ErrIntOverflow,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.left.Sub(tc.right)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Integer.Sub returned unexpected error: got %v, want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Integer.Sub returned unexpected result: (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestIntegerMul_ReturnsResult(t *testing.T) {
	testCases := []struct {
		name    string
		left    system.Integer
		right   system.Integer
		want    system.Integer
		wantErr error
	}{
		{
			name:  "multiplies two integers",
			left:  14,
			right: 2,
			want:  28,
		},
		{
			name:    "returns error if multiplication causes overflow",
			left:    1312312312,
			right:   10,
			wantErr: system.ErrIntOverflow,
		},
		{
			name:    "returns overflow error if MinInt32 is multiplied",
			left:    math.MinInt32,
			right:   -1,
			wantErr: system.ErrIntOverflow,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.left.Mul(tc.right)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Integer.Mul returned unexpected error: got %v, want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Integer.Mul returned unexpected result: (-want, +got)\n%s", diff)
			}
		})
	}
}
