package impl_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr/exprtest"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs/impl"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/fhirpath/trace"
	"github.com/verily-src/fhirpath-go/fhirpath/trace/tracetest"
)

func TestTimeOfDay(t *testing.T) {
	ctx := &expr.Context{Now: time.Date(0, time.January, 1, 19, 30, 5, 1000000, time.UTC)}
	wantCollection := system.Collection{system.MustParseTime("19:30:05.001")}

	got, err := impl.TimeOfDay(ctx, []any{})
	if err != nil {
		t.Fatalf("impl.TimeOfDay() returned unexpected error: %v", err)
	}
	if !cmp.Equal(got, wantCollection) {
		t.Errorf("impl.TimeOfDay() returned unexpected result: got %v, want %v", got, wantCollection)
	}
}

func TestToday(t *testing.T) {
	ctx := &expr.Context{Now: time.Date(2010, time.February, 12, 0, 0, 0, 0, time.UTC)}
	wantCollection := system.Collection{system.MustParseDate("2010-02-12")}

	got, err := impl.Today(ctx, []any{})
	if err != nil {
		t.Fatalf("impl.Today() returned unexpected error: %v", err)
	}
	if !cmp.Equal(got, wantCollection) {
		t.Errorf("impl.Today() returned unexpected result: got %v, want %v", got, wantCollection)
	}
}

func TestNow(t *testing.T) {
	ctx := &expr.Context{Now: time.Date(2010, time.February, 12, 12, 30, 34, 2000000, time.UTC)}
	wantCollection := system.Collection{system.MustParseDateTime("2010-02-12T12:30:34.002Z")}

	got, err := impl.Now(ctx, []any{})
	if err != nil {
		t.Fatalf("impl.Now() returned unexpected error: %v", err)
	}
	if !cmp.Equal(got, wantCollection) {
		t.Errorf("impl.Now() returned unexpected result: got %v, want %v", got, wantCollection)
	}
}

func TestTrace(t *testing.T) {
	t.Parallel()

	testErr := errors.New("test error")
	testName := "test-name"
	testCases := []struct {
		name      string
		tracer    trace.Tracer
		input     system.Collection
		args      []expr.Expression
		want      system.Collection
		wantTrace system.Collection
		wantErr   error
	}{
		{
			name:   "Too few arguments",
			tracer: trace.NoopTracer(),
			input: system.Collection{
				system.String("input value"),
			},
			args:    []expr.Expression{},
			wantErr: impl.ErrWrongArity,
		}, {
			name:   "Too many arguments",
			tracer: trace.NoopTracer(),
			input: system.Collection{
				system.String("input value"),
			},
			args: []expr.Expression{
				exprtest.Return(system.String("test")),
				exprtest.Return(system.String("test")),
				exprtest.Return(system.String("test")),
			},
			wantErr: impl.ErrWrongArity,
		}, {
			name:   "Tracer is not set in context, returns input collection",
			tracer: nil,
			input: system.Collection{
				system.String("input value"),
			},
			args: []expr.Expression{
				exprtest.Return(system.String("test")),
			},
			want: system.Collection{
				system.String("input value"),
			},
		}, {
			name:   "Name argument does not evaluate to single value",
			tracer: trace.NoopTracer(),
			input: system.Collection{
				system.String("input value"),
			},
			args: []expr.Expression{
				exprtest.Return(
					system.String("one"),
					system.String("two"),
				),
			},
			wantErr: cmpopts.AnyError,
		}, {
			name:   "Name argument does not evaluate to string",
			tracer: trace.NoopTracer(),
			input: system.Collection{
				system.String("input value"),
			},
			args: []expr.Expression{
				exprtest.Return(system.Integer(42)),
			},
			wantErr: cmpopts.AnyError,
		}, {
			name:   "Name argument fails to evaluate",
			tracer: trace.NoopTracer(),
			input: system.Collection{
				system.String("input value"),
			},
			args: []expr.Expression{
				exprtest.Error(testErr),
			},
			wantErr: testErr,
		}, {
			name:   "Projection argument fails to evaluate",
			tracer: trace.NoopTracer(),
			input: system.Collection{
				system.String("input value"),
			},
			args: []expr.Expression{
				exprtest.Return(system.String("test-name")),
				exprtest.Error(testErr),
			},
			wantErr: testErr,
		}, {
			name:   "Projection argument returns new value",
			tracer: trace.NoopTracer(),
			input: system.Collection{
				system.String("input value"),
			},
			args: []expr.Expression{
				exprtest.Return(system.String("test-name")),
				exprtest.Return(system.String("new value")),
			},
			want: system.Collection{
				system.String("input value"),
			},
		}, {
			name:   "Valid trace without projection",
			tracer: &tracetest.Recorder{},
			input: system.Collection{
				system.String("input value"),
			},
			args: []expr.Expression{
				exprtest.Return(system.String(testName)),
			},
			want: system.Collection{
				system.String("input value"),
			},
			wantTrace: system.Collection{
				system.String("input value"),
			},
		}, {
			name:   "Valid trace with projection",
			tracer: &tracetest.Recorder{},
			input: system.Collection{
				system.String("input value"),
			},
			args: []expr.Expression{
				exprtest.Return(system.String(testName)),
				exprtest.Return(system.String("projected value")),
			},
			want: system.Collection{
				system.String("input value"),
			},
			wantTrace: system.Collection{
				system.String("projected value"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := &expr.Context{Tracer: tc.tracer}

			collection, err := impl.Trace(ctx, tc.input, tc.args...)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("impl.Trace() got error %v, want %v", got, want)
			}
			if got, want := collection, tc.want; !cmp.Equal(got, want, cmpopts.EquateEmpty()) {
				t.Errorf("impl.Trace() got collection %v, want %v", got, want)
			}
			if reporter, ok := tc.tracer.(*tracetest.Recorder); ok {
				if got, want := reporter.Collection(testName), tc.wantTrace; !cmp.Equal(got, want, cmpopts.EquateEmpty()) {
					t.Errorf("impl.Trace() tracer got collection %v, want %v", got, want)
				}
			}
		})
	}
}
