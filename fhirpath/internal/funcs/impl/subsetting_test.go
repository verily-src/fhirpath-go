package impl_test

import (
	"testing"

	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr/exprtest"

	"google.golang.org/protobuf/testing/protocmp"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs/impl"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

func TestFirst(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:  "returns and empty collection if input is empty",
			input: system.Collection{},
			want:  system.Collection{},
		},
		{
			name: "returns first collection element",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5)},
			want: system.Collection{system.Integer(1)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.First(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("First() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("First() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestLast(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:  "returns and empty collection if input is empty",
			input: system.Collection{},
			want:  system.Collection{},
		},
		{
			name: "returns last collection element",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5)},
			want: system.Collection{system.Integer(5)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Last(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("Last() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Last() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestTail(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:  "returns and empty collection if input is empty",
			input: system.Collection{},
			want:  system.Collection{},
		},
		{
			name: "returns collection tail",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5)},
			want: system.Collection{
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Tail(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("Tail() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Tail() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestSkip(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if arg is not provided",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5)},
			wantErr: true,
		},
		{
			name:  "returns an empty collection if input is empty",
			input: system.Collection{},
			want:  system.Collection{},
		},
		{
			name: "returns an empty collection if input arg is greater than or equal to input length",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5)},
			args: []expr.Expression{
				exprtest.Return(system.Integer(5)),
			},
			want: system.Collection{},
		},
		{
			name: "returns the same input collection if arg is <= 0",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5)},
			args: []expr.Expression{
				exprtest.Return(system.Integer(-2)),
			},
			want: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5)},
		},
		{
			name: "returns the skipped input collection",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5)},
			args: []expr.Expression{
				exprtest.Return(system.Integer(3)),
			},
			want: system.Collection{
				system.Integer(4),
				system.Integer(5)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Skip(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("Skip() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Skip() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestTake(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if arg is not provided",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5)},
			wantErr: true,
		},
		{
			name:  "returns an empty collection if input is empty",
			input: system.Collection{},
			want:  system.Collection{},
		},
		{
			name: "returns an empty collection if arg is lower than or equal to 0 ",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5)},
			args: []expr.Expression{
				exprtest.Return(system.Integer(0)),
			},
			want: system.Collection{},
		},
		{
			name: "returns the same collection if arg is greater than or equal to input length",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5)},
			args: []expr.Expression{
				exprtest.Return(system.Integer(7)),
			},
			want: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5)},
		},
		{
			name: "returns the taken input collection",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5)},
			args: []expr.Expression{
				exprtest.Return(system.Integer(3)),
			},
			want: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Take(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("Take() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Take() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if arg is not provided",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5)},
			wantErr: true,
		},
		{
			name:  "returns an empty collection if input is empty",
			input: system.Collection{},
			want:  system.Collection{},
		},
		{
			name: "returns the intersection of both collections",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5),
				system.Integer(6)},
			args: []expr.Expression{
				exprtest.Return(
					system.Integer(3),
					system.Integer(6),
					system.Integer(9),
				),
			},
			want: system.Collection{
				system.Integer(3),
				system.Integer(6)},
		},
		{
			name: "returns the intersection of both collections ignoring duplicates",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(2),
				system.Integer(3)},
			args: []expr.Expression{
				exprtest.Return(
					system.Integer(2),
					system.Integer(2),
					system.Integer(4),
				),
			},
			want: system.Collection{
				system.Integer(2),
			},
		},
		{
			name: "returns the intersection of both collections with fhir.Integer types",
			input: system.Collection{
				fhir.Integer(1),
				fhir.Integer(2),
				fhir.Integer(3),
				fhir.Integer(4)},
			args: []expr.Expression{
				exprtest.Return(
					fhir.Integer(2),
					fhir.Integer(4),
				),
			},
			want: system.Collection{
				system.Integer(2),
				system.Integer(4),
			},
		},
		{
			name: "returns an empty collection if there is no intersection",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2)},
			args: []expr.Expression{
				exprtest.Return(
					system.Integer(3),
					system.Integer(4),
				),
			},
			want: system.Collection{},
		},
		{
			name: "returns intersection of normalized values",
			input: system.Collection{
				system.Integer(1),
				fhir.Integer(1)},
			args: []expr.Expression{
				exprtest.Return(
					fhir.Integer(1),
					system.Integer(1),
				),
			},
			want: system.Collection{
				system.Integer(1),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Intersect(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("Intersect() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Intersect() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestExclude(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if arg is not provided",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4),
				system.Integer(5)},
			wantErr: true,
		},
		{
			name:  "returns an empty collection if input is empty",
			input: system.Collection{},
			want:  system.Collection{},
		},
		{
			name: "returns the exclude of both collections",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(3),
				system.Integer(4)},
			args: []expr.Expression{
				exprtest.Return(
					system.Integer(3),
					system.Integer(4),
					system.Integer(5),
					system.Integer(6),
				),
			},
			want: system.Collection{
				system.Integer(1),
				system.Integer(2),
				system.Integer(5),
				system.Integer(6)},
		},
		{
			name: "returns the exclude of both collections with fhir.Integer types",
			input: system.Collection{
				fhir.Integer(1),
				fhir.Integer(2),
				fhir.Integer(3),
				fhir.Integer(4)},
			args: []expr.Expression{
				exprtest.Return(
					fhir.Integer(2),
					fhir.Integer(4),
				),
			},
			want: system.Collection{
				fhir.Integer(1),
				fhir.Integer(3),
			},
		},
		{
			name: "returns an empty collection if there is no exclude",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2)},
			args: []expr.Expression{
				exprtest.Return(
					system.Integer(1),
					system.Integer(2),
				),
			},
		},
		{
			name: "returns exclude of normalized values",
			input: system.Collection{
				system.Integer(1),
				fhir.Integer(2),
			},
			args: []expr.Expression{
				exprtest.Return(
					fhir.Integer(2),
					system.Integer(3),
				),
			},
			want: system.Collection{
				system.Integer(1),
				system.Integer(3),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Exclude(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("Exclude() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Exclude(%v) returned unexpected diff (-want, +got)\n%s", tc.name, diff)
			}
		})
	}
}

func TestDistinct(t *testing.T) {
	testCases := []struct {
		name  string
		input system.Collection
		want  system.Collection
	}{
		{
			name: "Empty input returns empty output",
		},
		{
			name:  "Inputs are distinct",
			input: system.Collection{system.String("Hello"), system.Integer(1), fhir.Integer(2)},
			want:  system.Collection{system.String("Hello"), system.Integer(1), fhir.Integer(2)},
		},
		{
			name:  "Inputs contain exact duplicate",
			input: system.Collection{system.String("Hello"), system.Integer(1), system.Integer(1)},
			want:  system.Collection{system.String("Hello"), system.Integer(1)},
		}, {
			name:  "Inputs contain system-convertible duplicates",
			input: system.Collection{system.String("Hello"), system.Integer(1), fhir.Integer(1)},
			want:  system.Collection{system.String("Hello"), system.Integer(1)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := expr.Context{}
			got, err := impl.Distinct(&ctx, tc.input)
			if err != nil {
				t.Fatalf("Distinct(%v): got unexpected err: %v", tc.name, err)
			}

			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Distinct(%v) returned unexpected diff (-want, +got)\n%s", tc.name, diff)
			}
		})
	}
}

func TestDistinct_BadInputs(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		wantErr error
	}{
		{
			name:    "Function called with nonzero arguments",
			input:   system.Collection{},
			args:    []expr.Expression{exprtest.Return()},
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := expr.Context{}
			_, err := impl.Distinct(&ctx, tc.input, tc.args...)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("Distinct(%v): got err %v, want  err%v", tc.name, got, want)
			}
		})
	}
}

func TestIsDistinct(t *testing.T) {
	testCases := []struct {
		name  string
		input system.Collection
		want  bool
	}{
		{
			name: "Empty input is distinct",
			want: true,
		},
		{
			name:  "Inputs are distinct",
			input: system.Collection{system.String("Hello"), system.Integer(1), fhir.Integer(2)},
			want:  true,
		},
		{
			name:  "Inputs contain exact duplicate",
			input: system.Collection{system.String("Hello"), system.Integer(1), system.Integer(1)},
			want:  false,
		}, {
			name:  "Inputs contain system-convertible duplicates",
			input: system.Collection{system.String("Hello"), system.Integer(1), fhir.Integer(1)},
			want:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := expr.Context{}
			collection, err := impl.IsDistinct(&ctx, tc.input)
			if err != nil {
				t.Fatalf("IsDistinct(%v): got unexpected err: %v", tc.name, err)
			}
			got, err := collection.ToBool()
			if err != nil {
				t.Fatalf("IsDistinct(%v): got unexpected err: %v", tc.name, err)
			}

			if got != tc.want {
				t.Errorf("IsDistinct(%v): got %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

func TestIsDistinct_BadInputs(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		wantErr error
	}{
		{
			name:    "Function called with nonzero arguments",
			input:   system.Collection{},
			args:    []expr.Expression{exprtest.Return()},
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := expr.Context{}
			_, err := impl.IsDistinct(&ctx, tc.input, tc.args...)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("IsDistinct(%v): got err %v, want  err%v", tc.name, got, want)
			}
		})
	}
}
