package impl_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs/impl"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

func TestStartsWith(t *testing.T) {
	fullString := system.String("Lee Jieun")

	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:  "returns empty for empty input",
			input: system.Collection{},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("")},
			},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "returns true for empty prefix",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("")},
			},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:  "returns true for match",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("Lee")},
			},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:  "returns false for no match",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:  "errors if input length is more than 1",
			input: system.Collection{fullString, fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if input is not a string",
			input: system.Collection{system.Integer(516)},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "errors if args length is not 1",
			input:   system.Collection{fullString},
			args:    []expr.Expression{},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if arg is not a string",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.Integer(516)},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.StartsWith(&expr.Context{}, tc.input, tc.args...)

			if gotErr := err != nil; tc.wantErr != gotErr {
				t.Fatalf("StartsWith got unexpected error result: gotErr %v, wantErr %v, err: %v", gotErr, tc.wantErr, err)
			}
			if !cmp.Equal(tc.want, got) {
				t.Errorf("StartsWith returned unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestEndsWith(t *testing.T) {
	fullString := system.String("Lee Jieun")

	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:  "returns empty for empty input",
			input: system.Collection{},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("")},
			},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "returns true for empty suffix",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("")},
			},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:  "returns true for match",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("Jieun")},
			},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:  "returns false for no match",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:  "errors if input length is more than 1",
			input: system.Collection{fullString, fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if input is not a string",
			input: system.Collection{system.Integer(516)},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "errors if args length is not 1",
			input:   system.Collection{fullString},
			args:    []expr.Expression{},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if arg is not a string",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.Integer(516)},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.EndsWith(&expr.Context{}, tc.input, tc.args...)

			if gotErr := err != nil; tc.wantErr != gotErr {
				t.Fatalf("EndsWith got unexpected error result: gotErr %v, wantErr %v, err: %v", gotErr, tc.wantErr, err)
			}
			if !cmp.Equal(tc.want, got) {
				t.Errorf("EndsWith returned unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestLength(t *testing.T) {
	fullString := system.String("Lee Jieun")

	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:    "returns empty for empty input",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "returns length of string",
			input:   system.Collection{fullString},
			want:    system.Collection{system.Integer(9)},
			wantErr: false,
		},
		{
			name:    "errors if input length is more than 1",
			input:   system.Collection{fullString, fullString},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "errors if input is not a string",
			input:   system.Collection{system.Integer(516)},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args is not empty",
			input: system.Collection{system.String("IU")},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("516")},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Length(&expr.Context{}, tc.input, tc.args...)

			if gotErr := err != nil; tc.wantErr != gotErr {
				t.Fatalf("Length got unexpected error result: gotErr %v, wantErr %v, err: %v", gotErr, tc.wantErr, err)
			}
			if !cmp.Equal(tc.want, got) {
				t.Errorf("Length returned unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestUpper(t *testing.T) {
	fullString := system.String("Lee Jieun")

	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:    "returns empty for empty input",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "returns upper string",
			input:   system.Collection{fullString},
			want:    system.Collection{system.String("LEE JIEUN")},
			wantErr: false,
		},
		{
			name:    "errors if input length is more than 1",
			input:   system.Collection{fullString, fullString},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "errors if input is not a string",
			input:   system.Collection{system.Integer(516)},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args is not empty",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Upper(&expr.Context{}, tc.input, tc.args...)

			if gotErr := err != nil; tc.wantErr != gotErr {
				t.Fatalf("Upper got unexpected error result: gotErr %v, wantErr %v, err: %v", gotErr, tc.wantErr, err)
			}
			if !cmp.Equal(tc.want, got) {
				t.Errorf("Upper returned unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestLower(t *testing.T) {
	fullString := system.String("Lee Jieun")

	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:    "returns empty for empty input",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "returns lower string",
			input:   system.Collection{fullString},
			want:    system.Collection{system.String("lee jieun")},
			wantErr: false,
		},
		{
			name:    "errors if input length is more than 1",
			input:   system.Collection{fullString, fullString},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "errors if input is not a string",
			input:   system.Collection{system.Integer(516)},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args is not empty",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Lower(&expr.Context{}, tc.input, tc.args...)

			if gotErr := err != nil; tc.wantErr != gotErr {
				t.Fatalf("Lower got unexpected error result: gotErr %v, wantErr %v, err: %v", gotErr, tc.wantErr, err)
			}
			if !cmp.Equal(tc.want, got) {
				t.Errorf("Lower returned unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	fullString := system.String("Lee Jieun")

	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:  "returns empty for empty input",
			input: system.Collection{},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("")},
			},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "returns true for empty substring",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("")},
			},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:  "returns true for match",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("Jie")},
			},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:  "returns false for no match",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:  "errors if input length is more than 1",
			input: system.Collection{fullString, fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if input is not a string",
			input: system.Collection{system.Integer(516)},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "errors if args length is not 1",
			input:   system.Collection{fullString},
			args:    []expr.Expression{},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if arg is not a string",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.Integer(516)},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Contains(&expr.Context{}, tc.input, tc.args...)

			if gotErr := err != nil; tc.wantErr != gotErr {
				t.Fatalf("Contains got unexpected error result: gotErr %v, wantErr %v, err: %v", gotErr, tc.wantErr, err)
			}
			if !cmp.Equal(tc.want, got) {
				t.Errorf("Contains returned unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestToChars(t *testing.T) {
	fullString := system.String("IU")

	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:    "returns empty for empty input",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "returns chars of string",
			input:   system.Collection{fullString},
			want:    system.Collection{system.String("I"), system.String("U")},
			wantErr: false,
		},
		{
			name:    "errors if input length is more than 1",
			input:   system.Collection{fullString, fullString},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "errors if input is not a string",
			input:   system.Collection{system.Integer(516)},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args is not empty",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.Integer(516)},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.ToChars(&expr.Context{}, tc.input, tc.args...)

			if gotErr := err != nil; tc.wantErr != gotErr {
				t.Fatalf("ToChars got unexpected error result: gotErr %v, wantErr %v, err: %v", gotErr, tc.wantErr, err)
			}
			if !cmp.Equal(tc.want, got) {
				t.Errorf("ToChars returned unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSubstring(t *testing.T) {
	fullString := system.String("Lee Jieun")

	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:  "returns empty for empty input",
			input: system.Collection{},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.Integer(0)},
			},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "returns empty is start is bigger than input string",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.Integer(50)},
			},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "returns substring with no length",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.Integer(4)},
			},
			want:    system.Collection{system.String("Jieun")},
			wantErr: false,
		},
		{
			name:  "returns substring with length",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.Integer(4)},
				&expr.LiteralExpression{Literal: system.Integer(2)},
			},
			want:    system.Collection{system.String("Ji")},
			wantErr: false,
		},
		{
			name:  "returns remaining chars if length overflows",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.Integer(4)},
				&expr.LiteralExpression{Literal: system.Integer(50)},
			},
			want:    system.Collection{system.String("Jieun")},
			wantErr: false,
		},
		{
			name:  "errors if input length is more than 1",
			input: system.Collection{fullString, fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.Integer(4)},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if input is not a string",
			input: system.Collection{system.Integer(516)},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.Integer(4)},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "errors if args length is not 1 or 2",
			input:   system.Collection{fullString},
			args:    []expr.Expression{},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if arg 1 is not an integer",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if arg 2 is not an integer",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.Integer(1)},
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Substring(&expr.Context{}, tc.input, tc.args...)

			if gotErr := err != nil; tc.wantErr != gotErr {
				t.Fatalf("Substring got unexpected error result: gotErr %v, wantErr %v, err: %v", gotErr, tc.wantErr, err)
			}
			if !cmp.Equal(tc.want, got) {
				t.Errorf("Substring returned unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIndexOf(t *testing.T) {
	fullString := system.String("Lee Jieun")

	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:  "returns empty for empty input",
			input: system.Collection{},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("")},
			},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "returns empty for empty arg",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{},
			},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "returns 0 for empty substring",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("")},
			},
			want:    system.Collection{system.Integer(0)},
			wantErr: false,
		},
		{
			name:  "returns proper index for match",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("Jieun")},
			},
			want:    system.Collection{system.Integer(4)},
			wantErr: false,
		},
		{
			name:  "returns -1 for no match",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    system.Collection{system.Integer(-1)},
			wantErr: false,
		},
		{
			name:  "errors if input length is more than 1",
			input: system.Collection{fullString, fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if input is not a string",
			input: system.Collection{system.Integer(516)},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "errors if args length is not 1",
			input:   system.Collection{fullString},
			args:    []expr.Expression{},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if arg is not a string",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.Integer(516)},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.IndexOf(&expr.Context{}, tc.input, tc.args...)

			if gotErr := err != nil; tc.wantErr != gotErr {
				t.Fatalf("IndexOf got unexpected error result: gotErr %v, wantErr %v, err: %v", gotErr, tc.wantErr, err)
			}
			if !cmp.Equal(tc.want, got) {
				t.Errorf("IndexOf returned unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestMatches(t *testing.T) {
	fullString := system.String("Lee Jieun")

	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:  "returns empty for empty input",
			input: system.Collection{},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("")},
			},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "returns empty for empty regex arg",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{},
			},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "returns true for empty regex",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("")},
			},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:  "returns true for match",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("^Lee[A-Za-z ]*$")},
			},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:  "returns false for no match",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("^Lee[0-9]*$")},
			},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:  "errors if input length is more than 1",
			input: system.Collection{fullString, fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if input is not a string",
			input: system.Collection{system.Integer(516)},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "errors if args length is not 1",
			input:   system.Collection{fullString},
			args:    []expr.Expression{},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if arg is not a string",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.Integer(516)},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if regex is invalid",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("^[$")},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Matches(&expr.Context{}, tc.input, tc.args...)

			if gotErr := err != nil; tc.wantErr != gotErr {
				t.Fatalf("Matches got unexpected error result: gotErr %v, wantErr %v, err: %v", gotErr, tc.wantErr, err)
			}
			if !cmp.Equal(tc.want, got) {
				t.Errorf("Matches returned unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestReplace(t *testing.T) {
	fullString := system.String("Lee Jieun")

	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:  "returns empty for empty input",
			input: system.Collection{},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("")},
				&expr.LiteralExpression{Literal: system.String("")},
			},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "returns empty for empty pattern",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{},
				&expr.LiteralExpression{Literal: system.String("Jieun")},
			},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "returns empty for empty substitution",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("Jieun")},
				&expr.LiteralExpression{},
			},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "returns replaced string",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("Jieun")},
				&expr.LiteralExpression{Literal: system.String("Uaena")},
			},
			want:    system.Collection{system.String("Lee Uaena")},
			wantErr: false,
		},
		{
			name:  "returns original string if both args are empty string",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("")},
				&expr.LiteralExpression{Literal: system.String("")},
			},
			want:    system.Collection{fullString},
			wantErr: false,
		},
		{
			name:  "removes pattern if substitution is empty string",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("e")},
				&expr.LiteralExpression{Literal: system.String("")},
			},
			want:    system.Collection{system.String("L Jiun")},
			wantErr: false,
		},
		{
			name:  "surrounds all characters with subtitution if pattern is empty string",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("")},
				&expr.LiteralExpression{Literal: system.String("!")},
			},
			want:    system.Collection{system.String("!L!e!e! !J!i!e!u!n!")},
			wantErr: false,
		},
		{
			name:  "errors if input length is more than 1",
			input: system.Collection{fullString, fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if input is not a string",
			input: system.Collection{system.Integer(516)},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "errors if args length is not 2",
			input:   system.Collection{fullString},
			args:    []expr.Expression{},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if arg 1 is not a string",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.Integer(516)},
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if arg 2 is not a string",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
				&expr.LiteralExpression{Literal: system.Integer(516)},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Replace(&expr.Context{}, tc.input, tc.args...)

			if gotErr := err != nil; tc.wantErr != gotErr {
				t.Fatalf("Replace got unexpected error result: gotErr %v, wantErr %v, err: %v", gotErr, tc.wantErr, err)
			}
			if !cmp.Equal(tc.want, got) {
				t.Errorf("Replace returned unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestReplaceMatches(t *testing.T) {
	fullString := system.String("Woo Young Woo")
	dateRegex := `\b(?P<month>\d{1,2})/(?P<day>\d{1,2})/(?P<year>\d{2,4})\b`

	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:  "returns empty for empty input",
			input: system.Collection{},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("")},
				&expr.LiteralExpression{Literal: system.String("")},
			},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "returns empty for empty pattern arg",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{},
				&expr.LiteralExpression{Literal: system.String(" to the ")},
			},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "returns empty for empty substitution arg",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String(`[ ]`)},
				&expr.LiteralExpression{},
			},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "returns replaced string",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String(`[ ]`)},
				&expr.LiteralExpression{Literal: system.String(" to the ")},
			},
			want:    system.Collection{system.String("Woo to the Young to the Woo")},
			wantErr: false,
		},
		// This test case comes directly from the FHIRPath spec.
		// https://build.fhir.org/ig/HL7/FHIRPath/#replacematchesregex-string-substitution-string-string
		{
			name:  "returns replaced string with named patterns",
			input: system.Collection{system.String("5/16/1993")},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String(dateRegex)},
				&expr.LiteralExpression{Literal: system.String("${day}/${month}/${year}")},
			},
			want:    system.Collection{system.String("16/5/1993")},
			wantErr: false,
		},
		{
			name:  "errors if input length is more than 1",
			input: system.Collection{fullString, fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if input is not a string",
			input: system.Collection{system.Integer(516)},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "errors if args length is not 2",
			input:   system.Collection{fullString},
			args:    []expr.Expression{},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if arg 1 is not a string",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.Integer(516)},
				&expr.LiteralExpression{Literal: system.String("IU")},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if arg 2 is not a string",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("IU")},
				&expr.LiteralExpression{Literal: system.Integer(516)},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if regex is invalid",
			input: system.Collection{fullString},
			args: []expr.Expression{
				&expr.LiteralExpression{Literal: system.String("^[$")},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.ReplaceMatches(&expr.Context{}, tc.input, tc.args...)

			if gotErr := err != nil; tc.wantErr != gotErr {
				t.Fatalf("ReplaceMatches got unexpected error result: gotErr %v, wantErr %v, err: %v", gotErr, tc.wantErr, err)
			}
			if !cmp.Equal(tc.want, got) {
				t.Errorf("ReplaceMatches returned unexpected result: got %v, want %v", got, tc.want)
			}
		})
	}
}
