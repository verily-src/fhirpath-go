package impl_test

import (
	"testing"

	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr/exprtest"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs/impl"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

func TestAbs(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:    "errors if input is not a number",
			input:   system.Collection{system.String("1.2")},
			want:    nil,
			wantErr: true,
		},
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.MustParseDecimal("10.1"),
				system.MustParseDecimal("10.5")},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args length is more than 0",
			input: system.Collection{system.MustParseDecimal("10.8")},
			args: []expr.Expression{
				exprtest.Return(system.String("kg")),
				exprtest.Return(system.String("lb")),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "returns an empty collection if input is empty",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "abs a positive Decimal number",
			input:   system.Collection{system.MustParseDecimal("10.5")},
			want:    system.Collection{system.MustParseDecimal("10.5")},
			wantErr: false,
		},
		{
			name:    "abs a negative Decimal number",
			input:   system.Collection{system.MustParseDecimal("-10.5")},
			want:    system.Collection{system.MustParseDecimal("10.5")},
			wantErr: false,
		},
		{
			name:    "abs a positive Integer number",
			input:   system.Collection{system.Integer(11)},
			want:    system.Collection{system.Integer(11)},
			wantErr: false,
		},
		{
			name:    "abs a negative Integer number",
			input:   system.Collection{system.Integer(-11)},
			want:    system.Collection{system.Integer(11)},
			wantErr: false,
		},
		{
			name:    "abs a positive Quantity number",
			input:   system.Collection{system.MustParseQuantity("10.5", "kg")},
			want:    system.Collection{system.MustParseQuantity("10.5", "kg")},
			wantErr: false,
		},
		{
			name:    "abs a negative Quantity number",
			input:   system.Collection{system.MustParseQuantity("-10.5", "kg")},
			want:    system.Collection{system.MustParseQuantity("10.5", "kg")},
			wantErr: false,
		},
		{
			name:    "abs zero number",
			input:   system.Collection{system.Integer(0)},
			want:    system.Collection{system.Integer(0)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Abs(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("Abs() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Abs() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestCeiling(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:    "errors if input is not a number",
			input:   system.Collection{system.String("1.2")},
			want:    nil,
			wantErr: true,
		},
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.MustParseDecimal("10.1"),
				system.MustParseDecimal("10.5")},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args length is more than 0",
			input: system.Collection{system.MustParseDecimal("10")},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("20")),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "returns an empty collection if input is empty",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "ceiling a positive float number",
			input:   system.Collection{system.MustParseDecimal("10.5")},
			want:    system.Collection{system.Integer(11)},
			wantErr: false,
		},
		{
			name:    "ceiling a positive float number with zero decimal",
			input:   system.Collection{system.MustParseDecimal("10.0")},
			want:    system.Collection{system.Integer(10)},
			wantErr: false,
		},
		{
			name:    "ceiling a negative float number",
			input:   system.Collection{system.MustParseDecimal("-10.5")},
			want:    system.Collection{system.Integer(-10)},
			wantErr: false,
		},
		{
			name:    "ceiling a negative float number with zero decimal",
			input:   system.Collection{system.MustParseDecimal("-10.0")},
			want:    system.Collection{system.Integer(-10)},
			wantErr: false,
		},
		{
			name:    "ceiling an Integer number",
			input:   system.Collection{fhir.Integer(10)},
			want:    system.Collection{system.Integer(10)},
			wantErr: false,
		},
		{
			name:    "ceiling a negative Integer number",
			input:   system.Collection{fhir.Integer(-10)},
			want:    system.Collection{system.Integer(-10)},
			wantErr: false,
		},
		{
			name:    "ceiling a PositiveInt number",
			input:   system.Collection{fhir.PositiveInt(100)},
			want:    system.Collection{system.Integer(100)},
			wantErr: false,
		},
		{
			name:    "ceiling an UnsignedInt number",
			input:   system.Collection{fhir.UnsignedInt(1000)},
			want:    system.Collection{system.Integer(1000)},
			wantErr: false,
		},
		{
			name:    "ceiling zero number",
			input:   system.Collection{fhir.Integer(0)},
			want:    system.Collection{system.Integer(0)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Ceiling(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("Ceiling() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Ceiling() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestExp(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:    "errors if input is not a number",
			input:   system.Collection{system.String("1.2")},
			want:    nil,
			wantErr: true,
		},
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.Integer(1),
				system.Integer(2)},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args length is more than 0",
			input: system.Collection{system.Integer(1)},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("2")),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "returns an empty collection if input is empty",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "returns a system.Decimal if input is system.Integer",
			input:   system.Collection{system.Integer(0)},
			want:    system.Collection{system.MustParseDecimal("1")},
			wantErr: false,
		},
		{
			name:    "returns a system.Decimal if input is system.Decimal",
			input:   system.Collection{system.MustParseDecimal("0")},
			want:    system.Collection{system.MustParseDecimal("1")},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Exp(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("Exp() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Exp() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestFloor(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:    "errors if input is not a number",
			input:   system.Collection{system.String("1.2")},
			want:    nil,
			wantErr: true,
		},
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.MustParseDecimal("10.1"),
				system.MustParseDecimal("10.5")},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args length is more than 0",
			input: system.Collection{system.MustParseDecimal("10")},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("20")),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "returns an empty collection if input is empty",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "floors a positive float number",
			input:   system.Collection{system.MustParseDecimal("10.5")},
			want:    system.Collection{system.Integer(10)},
			wantErr: false,
		},
		{
			name:    "floors a positive float number with zero decimal",
			input:   system.Collection{system.MustParseDecimal("10.0")},
			want:    system.Collection{system.Integer(10)},
			wantErr: false,
		},
		{
			name:    "floors a negative float number",
			input:   system.Collection{system.MustParseDecimal("-10.5")},
			want:    system.Collection{system.Integer(-11)},
			wantErr: false,
		},
		{
			name:    "floors a negative float number with zero decimal",
			input:   system.Collection{system.MustParseDecimal("-10.0")},
			want:    system.Collection{system.Integer(-10)},
			wantErr: false,
		},
		{
			name:    "floors an Integer number",
			input:   system.Collection{fhir.Integer(10)},
			want:    system.Collection{system.Integer(10)},
			wantErr: false,
		},
		{
			name:    "floors a negative Integer number",
			input:   system.Collection{fhir.Integer(-10)},
			want:    system.Collection{system.Integer(-10)},
			wantErr: false,
		},
		{
			name:    "floors a PositiveInt number",
			input:   system.Collection{fhir.PositiveInt(100)},
			want:    system.Collection{system.Integer(100)},
			wantErr: false,
		},
		{
			name:    "floors an UnsignedInt number",
			input:   system.Collection{fhir.UnsignedInt(1000)},
			want:    system.Collection{system.Integer(1000)},
			wantErr: false,
		},
		{
			name:    "floors zero number",
			input:   system.Collection{fhir.Integer(0)},
			want:    system.Collection{system.Integer(0)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Floor(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("Floor() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Floor() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestLn(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:    "errors if input is not a number",
			input:   system.Collection{system.String("1.2")},
			want:    nil,
			wantErr: true,
		},
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.MustParseDecimal("10.1"),
				system.MustParseDecimal("10.5")},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args length is more than 0",
			input: system.Collection{system.MustParseDecimal("10")},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("20")),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "returns an empty collection if input is empty",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "returns an empty collection if result is NaN",
			input:   system.Collection{system.MustParseDecimal("-1.0")},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "lns a positive number",
			input:   system.Collection{system.MustParseDecimal("2")},
			want:    system.Collection{system.MustParseDecimal("0.6931471805599453")},
			wantErr: false,
		},
		{
			name:    "lns a positive float number",
			input:   system.Collection{system.MustParseDecimal("0.5")},
			want:    system.Collection{system.MustParseDecimal("-0.6931471805599453")},
			wantErr: false,
		},
		{
			name:    "lns an PositiveInt number",
			input:   system.Collection{fhir.PositiveInt(16)},
			want:    system.Collection{system.MustParseDecimal("2.772588722239781")},
			wantErr: false,
		},
		{
			name:    "lns an UnsignedInt number",
			input:   system.Collection{fhir.UnsignedInt(16)},
			want:    system.Collection{system.MustParseDecimal("2.772588722239781")},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Ln(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("Ln() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Ln() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestLog(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:    "errors if input is not a number",
			input:   system.Collection{system.String("1.2")},
			want:    nil,
			wantErr: true,
		},
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.MustParseDecimal("10.1"),
				system.MustParseDecimal("10.5")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "returns an empty collection if input is empty",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "returns an empty collection if result is NaN",
			input: system.Collection{system.MustParseDecimal("-1.0")},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("0.5")),
			},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "logs a number with base 2",
			input: system.Collection{system.MustParseDecimal("16")},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("2")),
			},
			want:    system.Collection{system.MustParseDecimal("4")},
			wantErr: false,
		},
		{
			name:  "logs a number with base 10",
			input: system.Collection{system.MustParseDecimal("100")},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("10")),
			},
			want:    system.Collection{system.MustParseDecimal("2")},
			wantErr: false,
		},
		{
			name:  "logs a PositiveInt number with base 2",
			input: system.Collection{fhir.PositiveInt(16)},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("2")),
			},
			want:    system.Collection{system.MustParseDecimal("4")},
			wantErr: false,
		},
		{
			name:  "logs an UnsignedInt number with base 2",
			input: system.Collection{fhir.UnsignedInt(16)},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("2")),
			},
			want:    system.Collection{system.MustParseDecimal("4")},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Log(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("Log() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Log() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestPower(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:    "errors if input is not a number",
			input:   system.Collection{system.String("1.2")},
			want:    nil,
			wantErr: true,
		},
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.MustParseDecimal("10.1"),
				system.MustParseDecimal("10.5")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "returns an empty collection if input is empty",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "returns an empty collection if result is NaN",
			input: system.Collection{system.MustParseDecimal("-1.0")},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("0.5")),
			},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "powers a positive float number to a positive float arg",
			input: system.Collection{system.MustParseDecimal("2.5")},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("2")),
			},
			want:    system.Collection{system.MustParseDecimal("6.25")},
			wantErr: false,
		},
		{
			name:  "powers a positive float number to a negative float arg",
			input: system.Collection{system.MustParseDecimal("2.5")},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("-2")),
			},
			want:    system.Collection{system.MustParseDecimal("0.16")},
			wantErr: false,
		},
		{
			name:  "powers a positive float number to a positive int arg",
			input: system.Collection{system.MustParseDecimal("2.5")},
			args: []expr.Expression{
				exprtest.Return(system.Integer(2)),
			},
			want:    system.Collection{system.MustParseDecimal("6.25")},
			wantErr: false,
		},
		{
			name:  "powers a positive float number to a negative int arg",
			input: system.Collection{system.MustParseDecimal("2.5")},
			args: []expr.Expression{
				exprtest.Return(system.Integer(-2)),
			},
			want:    system.Collection{system.MustParseDecimal("0.16")},
			wantErr: false,
		},
		{
			name:  "powers an integer float number to a positive float arg",
			input: system.Collection{system.Integer(4)},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("2.5")),
			},
			want:    system.Collection{system.MustParseDecimal("32")},
			wantErr: false,
		},
		{
			name:  "powers an integer float number to a negative float arg",
			input: system.Collection{system.Integer(4)},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("-2.5")),
			},
			want:    system.Collection{system.MustParseDecimal("0.03125")},
			wantErr: false,
		},
		{
			name:  "powers an integer number to a negative int arg",
			input: system.Collection{system.Integer(4)},
			args: []expr.Expression{
				exprtest.Return(system.Integer(-2)),
			},
			want:    system.Collection{system.Integer(0)},
			wantErr: false,
		},
		{
			name:  "powers a PositiveInt number",
			input: system.Collection{fhir.PositiveInt(10)},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("2.0")),
			},
			want:    system.Collection{system.MustParseDecimal("100")},
			wantErr: false,
		},
		{
			name:  "powers an UnsignedInt number",
			input: system.Collection{fhir.UnsignedInt(20)},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("2.0")),
			},
			want:    system.Collection{system.MustParseDecimal("400")},
			wantErr: false,
		},
		{
			name:  "powers zero number to positive float exp",
			input: system.Collection{fhir.Integer(0)},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("10.5")),
			},
			want:    system.Collection{system.MustParseDecimal("0")},
			wantErr: false,
		},
		{
			name:  "powers a positive integer number to zero exp",
			input: system.Collection{fhir.Integer(10)},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("0")),
			},
			want:    system.Collection{system.MustParseDecimal("1")},
			wantErr: false,
		},
		{
			name:  "powers zero number to zero exp",
			input: system.Collection{fhir.Integer(0)},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("0")),
			},
			want:    system.Collection{system.MustParseDecimal("1")},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Power(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("Power() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Power() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestRound(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:    "errors if input is not a number",
			input:   system.Collection{system.String("1.2")},
			want:    nil,
			wantErr: true,
		},
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.MustParseDecimal("10.1"),
				system.MustParseDecimal("10.5")},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args len is greater than 1",
			input: system.Collection{system.MustParseDecimal("3.141592653589793")},
			args: []expr.Expression{
				exprtest.Return(system.Integer(1)),
				exprtest.Return(system.Integer(2)),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if arg is not an Integer",
			input: system.Collection{system.MustParseDecimal("3.141592653589793")},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("2.2")),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if precision arg is negative",
			input: system.Collection{system.MustParseDecimal("3.141592653589793")},
			args: []expr.Expression{
				exprtest.Return(system.Integer(-2)),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "returns an empty collection if input is empty",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "rounds a positive Decimal",
			input: system.Collection{system.MustParseDecimal("3.141592653589793")},
			args: []expr.Expression{
				exprtest.Return(system.Integer(4)),
			},
			want:    system.Collection{system.MustParseDecimal("3.1416")},
			wantErr: false,
		},
		{
			name:    "rounds a positive Decimal with no precision arg",
			input:   system.Collection{system.MustParseDecimal("3.141592653589793")},
			want:    system.Collection{system.MustParseDecimal("3")},
			wantErr: false,
		},
		{
			name:  "rounds a negative Decimal",
			input: system.Collection{system.MustParseDecimal("-3.141592653589793")},
			args: []expr.Expression{
				exprtest.Return(system.Integer(4)),
			},
			want:    system.Collection{system.MustParseDecimal("-3.1416")},
			wantErr: false,
		},
		{
			name:    "rounds a positive Integer",
			input:   system.Collection{system.Integer(3)},
			want:    system.Collection{system.MustParseDecimal("3")},
			wantErr: false,
		},
		{
			name:  "rounds a positive Integer with precision arg",
			input: system.Collection{system.Integer(3)},
			args: []expr.Expression{
				exprtest.Return(system.Integer(2)),
			},
			want:    system.Collection{system.MustParseDecimal("3")},
			wantErr: false,
		},
		{
			name:    "rounds a negative Integer",
			input:   system.Collection{system.Integer(-7)},
			want:    system.Collection{system.MustParseDecimal("-7")},
			wantErr: false,
		},
		{
			name:    "rounds a fhir Integer",
			input:   system.Collection{fhir.Integer(20)},
			want:    system.Collection{system.MustParseDecimal("20")},
			wantErr: false,
		},
		{
			name:    "rounds a fhir PositiveInt",
			input:   system.Collection{fhir.PositiveInt(30)},
			want:    system.Collection{system.MustParseDecimal("30")},
			wantErr: false,
		},
		{
			name:    "rounds a fhir UnsignedInt",
			input:   system.Collection{fhir.UnsignedInt(10)},
			want:    system.Collection{system.MustParseDecimal("10")},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Round(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("Round() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Round() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestSqrt(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:    "errors if input is not a number",
			input:   system.Collection{system.String("1.2")},
			want:    nil,
			wantErr: true,
		},
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.MustParseDecimal("10.1"),
				system.MustParseDecimal("10.5")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "returns an empty collection if input is empty",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "errors if input is negative",
			input:   system.Collection{system.MustParseDecimal("-16.0")},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args length is more than 0",
			input: system.Collection{system.MustParseDecimal("10")},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("20")),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "sqrt a positive float number",
			input:   system.Collection{system.MustParseDecimal("16.5")},
			want:    system.Collection{system.MustParseDecimal("4.06201920231798")},
			wantErr: false,
		},
		{
			name:    "sqrt a positive float number with zero decimal",
			input:   system.Collection{system.MustParseDecimal("81.0")},
			want:    system.Collection{system.MustParseDecimal("9.0")},
			wantErr: false,
		},
		{
			name:    "sqrt an Integer number",
			input:   system.Collection{fhir.Integer(400)},
			want:    system.Collection{system.MustParseDecimal("20.0")},
			wantErr: false,
		},
		{
			name:    "sqrt a PositiveInt number",
			input:   system.Collection{fhir.PositiveInt(100)},
			want:    system.Collection{system.MustParseDecimal("10.0")},
			wantErr: false,
		},
		{
			name:    "sqrt an UnsignedInt number",
			input:   system.Collection{fhir.UnsignedInt(9)},
			want:    system.Collection{system.MustParseDecimal("3.0")},
			wantErr: false,
		},
		{
			name:    "sqrt zero number",
			input:   system.Collection{fhir.Integer(0)},
			want:    system.Collection{system.MustParseDecimal("0.0")},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Sqrt(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("Sqrt() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Sqrt() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name:    "errors if input is not a number",
			input:   system.Collection{system.String("1.2")},
			want:    nil,
			wantErr: true,
		},
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.MustParseDecimal("10.1"),
				system.MustParseDecimal("10.5")},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args length is more than 0",
			input: system.Collection{system.MustParseDecimal("10")},
			args: []expr.Expression{
				exprtest.Return(system.MustParseDecimal("20")),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "returns an empty collection if input is empty",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "truncates a positive float number",
			input:   system.Collection{system.MustParseDecimal("10.12345")},
			want:    system.Collection{system.Integer(10)},
			wantErr: false,
		},
		{
			name:    "truncates a positive float number with zero decimal",
			input:   system.Collection{system.MustParseDecimal("10.00")},
			want:    system.Collection{system.Integer(10)},
			wantErr: false,
		},
		{
			name:    "truncates a negative float number",
			input:   system.Collection{system.MustParseDecimal("-10.12345")},
			want:    system.Collection{system.Integer(-10)},
			wantErr: false,
		},
		{
			name:    "truncates a negative float number with zero decimal",
			input:   system.Collection{system.MustParseDecimal("-10.000")},
			want:    system.Collection{system.Integer(-10)},
			wantErr: false,
		},
		{
			name:    "truncates an Integer number",
			input:   system.Collection{fhir.Integer(10)},
			want:    system.Collection{system.Integer(10)},
			wantErr: false,
		},
		{
			name:    "truncates a negative Integer number",
			input:   system.Collection{fhir.Integer(-10)},
			want:    system.Collection{system.Integer(-10)},
			wantErr: false,
		},
		{
			name:    "truncates a PositiveInt number",
			input:   system.Collection{fhir.PositiveInt(100)},
			want:    system.Collection{system.Integer(100)},
			wantErr: false,
		},
		{
			name:    "truncates an UnsignedInt number",
			input:   system.Collection{fhir.UnsignedInt(1000)},
			want:    system.Collection{system.Integer(1000)},
			wantErr: false,
		},
		{
			name:    "truncates zero number",
			input:   system.Collection{fhir.Integer(0)},
			want:    system.Collection{system.Integer(0)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Truncate(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("Truncate() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Truncate() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}
