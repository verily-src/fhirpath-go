package impl_test

import (
	"testing"

	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr/exprtest"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"google.golang.org/protobuf/testing/protocmp"

	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs/impl"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

func TestConvertsToBoolean(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.String("T"),
				system.String("True")},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args length is more than 0",
			input: system.Collection{system.String("false")},
			args: []expr.Expression{
				exprtest.Return(system.String("200")),
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
			name:    "returns false if input is not convertible",
			input:   system.Collection{system.String("404 Kg")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.Decimal '0.0'",
			input:   system.Collection{system.MustParseDecimal("0.0")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.Decimal '1.0'",
			input:   system.Collection{system.MustParseDecimal("1.0")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.Decimal '3.5'",
			input:   system.Collection{system.MustParseDecimal("3.5")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.Integer '0'",
			input:   system.Collection{system.Integer(0)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.Integer '1'",
			input:   system.Collection{system.Integer(1)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.Integer '3'",
			input:   system.Collection{system.Integer(3)},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.String '2.0'",
			input:   system.Collection{system.String("2.0")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.String '1.0'",
			input:   system.Collection{system.String("1.0")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'y'",
			input:   system.Collection{system.String("y")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'yes'",
			input:   system.Collection{system.String("yes")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'Y'",
			input:   system.Collection{system.String("Y")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'YES'",
			input:   system.Collection{system.String("YES")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String '0.0'",
			input:   system.Collection{system.String("0.0")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'f'",
			input:   system.Collection{system.String("f")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'false'",
			input:   system.Collection{system.String("false")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'F'",
			input:   system.Collection{system.String("F")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'FALSE'",
			input:   system.Collection{system.String("FALSE")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'n",
			input:   system.Collection{system.String("n")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'no'",
			input:   system.Collection{system.String("no")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'N",
			input:   system.Collection{system.String("N")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'NO'",
			input:   system.Collection{system.String("NO")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.Boolean 'true'",
			input:   system.Collection{system.Boolean(true)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.Boolean 'false'",
			input:   system.Collection{system.Boolean(false)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is fhir.Integer '1'",
			input:   system.Collection{fhir.Integer(1)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is fhir.PositiveInt '1'",
			input:   system.Collection{fhir.PositiveInt(1)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is fhir.UnsignedInt '1'",
			input:   system.Collection{fhir.UnsignedInt(1)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.ConvertsToBoolean(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("ConvertsToBoolean() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ConvertsToBoolean() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestConvertsToDate(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.String("2001-09-11"),
				system.String("2011-05-02")},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args length is more than 0",
			input: system.Collection{system.String("false")},
			args: []expr.Expression{
				exprtest.Return(system.String("minutes")),
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
			name:    "returns false if input is not convertible to system.Date",
			input:   system.Collection{system.MustParseQuantity("75", "Kg")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "returns true for a system.Date",
			input:   system.Collection{system.MustParseDate("1993-08-13")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a partial system.Date",
			input:   system.Collection{system.MustParseDate("1993-08")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a system.DateTime",
			input:   system.Collection{system.MustParseDateTime("1993-08-13T14:20:00")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a convertible system.String",
			input:   system.Collection{system.String("1993-08-13")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns false for a non convertible system.String",
			input:   system.Collection{system.String("93.08.13")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "returns false for a ppb.Patient",
			input:   system.Collection{&ppb.Patient{}},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.ConvertsToDate(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("ConvertsToDate() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ConvertsToDate() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestConvertsToDateTime(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.String("2001-09-11"),
				system.String("2011-05-02")},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args length is more than 0",
			input: system.Collection{system.String("2001-09-11")},
			args: []expr.Expression{
				exprtest.Return(system.String("minutes")),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "returns empty if input is empty",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "returns false if input is not convertible to system.DateTime",
			input:   system.Collection{system.MustParseQuantity("75", "Kg")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "returns true for a system.DateTime",
			input:   system.Collection{system.MustParseDateTime("1993-08-13T14:20:00")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a partial system.DateTime",
			input:   system.Collection{system.MustParseDateTime("2012-01-01T10:00")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a system.Date",
			input:   system.Collection{system.MustParseDate("2006-01-02")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a convertible system.String",
			input:   system.Collection{system.String("1993-08-13T14:20:00")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns empty for a non convertible system.String",
			input:   system.Collection{system.String("93.08.13")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "returns empty for a ppb.Patient",
			input:   system.Collection{&ppb.Patient{}},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.ConvertsToDateTime(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("ConvertsToDateTime() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ConvertsToDateTime() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestConvertsToDecimal(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.String("101"),
				system.String("102")},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args length is more than 0",
			input: system.Collection{system.String("100")},
			args: []expr.Expression{
				exprtest.Return(system.String("200")),
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
			name:    "input is system.Decimal 'true'",
			input:   system.Collection{system.MustParseDecimal("13.5")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.Decimal 'false'",
			input:   system.Collection{system.MustParseDecimal("-13.5")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.Integer '13'",
			input:   system.Collection{system.Integer(13)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input system.Integer '-13'",
			input:   system.Collection{system.Integer(-13)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input system.String '+13.5'",
			input:   system.Collection{system.String("+13.5")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String '-13.5'",
			input:   system.Collection{system.String("-13.5")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String '3.1416 cm'",
			input:   system.Collection{system.String("3.1416 cm")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.String '12.99'",
			input:   system.Collection{system.String(" 12.99 ")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is true system.Boolean 'true'",
			input:   system.Collection{system.Boolean(true)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.Boolean 'false'",
			input:   system.Collection{system.Boolean(false)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is fhir.Integer '10'",
			input:   system.Collection{fhir.Integer(10)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is fhir.PositiveInt '11'",
			input:   system.Collection{fhir.PositiveInt(11)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is fhir.UnsignedInt '12'",
			input:   system.Collection{fhir.UnsignedInt(12)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is fhir.Boolean 'true'",
			input:   system.Collection{fhir.Boolean(true)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.ConvertsToDecimal(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("ConvertsToDecimal() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ConvertsToDecimal() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestConvertsToInteger(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.String("101"),
				system.String("102")},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args length is more than 0",
			input: system.Collection{system.String("100")},
			args: []expr.Expression{
				exprtest.Return(system.String("200")),
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
			name:    "input is system.Integer '13'",
			input:   system.Collection{system.Integer(13)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.Integer '-13'",
			input:   system.Collection{system.Integer(-13)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String '13'",
			input:   system.Collection{system.String("13")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String '+13'",
			input:   system.Collection{system.String("+13")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String '404 kg'",
			input:   system.Collection{system.String("404 Kg")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.String ' 12 '",
			input:   system.Collection{system.String(" 12 ")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.String '-13'",
			input:   system.Collection{system.String("-13")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.Boolean 'true'",
			input:   system.Collection{system.Boolean(true)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.Boolean 'false'",
			input:   system.Collection{system.Boolean(false)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is fhir.Integer '10'",
			input:   system.Collection{fhir.Integer(10)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is fhir.PositiveInt '11'",
			input:   system.Collection{fhir.PositiveInt(11)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is fhir.UnsignedInt '12'",
			input:   system.Collection{fhir.UnsignedInt(12)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is fhir.Boolean 'true'",
			input:   system.Collection{fhir.Boolean(true)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.ConvertsToInteger(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("ConvertsToInteger() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ConvertsToInteger() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestConvertsToQuantity(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.String("10 'km/hr'"),
				system.String("10 'mi/hr'")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "returns false if input is not convertible to a system.Quantity",
			input:   system.Collection{system.String("10 km / hr")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:  "errors if args length is more than 1",
			input: system.Collection{system.String("2 days")},
			args: []expr.Expression{
				exprtest.Return(system.String("hours")),
				exprtest.Return(system.String("minutes")),
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
			name:    "input is system.Integer '13'",
			input:   system.Collection{system.Integer(13)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.Decimal '13.5",
			input:   system.Collection{system.MustParseDecimal("13.5")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.Quantity '13.5 lbs'",
			input:   system.Collection{system.MustParseQuantity("13.5", "lbs")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String '100 days'",
			input:   system.Collection{system.String("100 days")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String '100 km'",
			input:   system.Collection{system.String("100 km")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String '100 km/h'",
			input:   system.Collection{system.String("100 km/h")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.String '100 'km/h''",
			input:   system.Collection{system.String("100 'km/h'")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String '100 'km per hour''",
			input:   system.Collection{system.String("100 'km per hour'")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String '100 km per hour'",
			input:   system.Collection{system.String("100 km per hour")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.String '100           km'",
			input:   system.Collection{system.String("100           km")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String '100           km/h'",
			input:   system.Collection{system.String("100           km/h")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.String '10 'km per hr''",
			input:   system.Collection{system.String("10 'km per hr'")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:  "input is system.String '2 years' with arg 'months'",
			input: system.Collection{system.String("2 years")},
			args: []expr.Expression{
				exprtest.Return(system.String("months")),
			},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:  "input is system.String '2 months' with arg 'days'",
			input: system.Collection{system.String("2 months")},
			args: []expr.Expression{
				exprtest.Return(system.String("days")),
			},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:  "input is system.String '2 days' with arg 'hours'",
			input: system.Collection{system.String("2 days")},
			args: []expr.Expression{
				exprtest.Return(system.String("hours")),
			},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:  "input is system.String '2 hours' with arg 'minutes'",
			input: system.Collection{system.String("2 hours")},
			args: []expr.Expression{
				exprtest.Return(system.String("minutes")),
			},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:  "input is system.String '2 minutes' with arg 'seconds'",
			input: system.Collection{system.String("2 minutes")},
			args: []expr.Expression{
				exprtest.Return(system.String("seconds")),
			},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String '100 'km''",
			input:   system.Collection{system.String("100 'km'")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:  "input is system.Integer '100' with arg ''km''",
			input: system.Collection{system.Integer(100)},
			args: []expr.Expression{
				exprtest.Return(system.String("'km'")),
			},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:  "input is system.Integer '100' with arg 'days''",
			input: system.Collection{system.Integer(100)},
			args: []expr.Expression{
				exprtest.Return(system.String("days")),
			},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.Boolean 'true'",
			input:   system.Collection{system.Boolean(true)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.Boolean 'false'",
			input:   system.Collection{system.Boolean(false)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is fhir.Integer '10'",
			input:   system.Collection{fhir.Integer(10)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is fhir.PositiveInt '11'",
			input:   system.Collection{fhir.PositiveInt(11)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is fhir.UnsignedInt '12'",
			input:   system.Collection{fhir.UnsignedInt(12)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.ConvertsToQuantity(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("ConvertsToQuantity() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ConvertsToQuantity() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestConvertsToString(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if input is not a singleton",
			input: system.Collection{
				system.String("101"),
				system.String("102"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "returns empty for and empty input collection",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "returns true for a system.String",
			input:   system.Collection{system.String("100")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a fhir.String",
			input:   system.Collection{fhir.String("100")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a system.Integer",
			input:   system.Collection{system.Integer(100)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a fhir.Integer",
			input:   system.Collection{fhir.Integer(100)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a fhir.PositiveInt",
			input:   system.Collection{fhir.PositiveInt(11)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a fhir.UnsignedInt",
			input:   system.Collection{fhir.UnsignedInt(12)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a system.Decimal",
			input:   system.Collection{system.MustParseDecimal("100.999")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a system.Date",
			input:   system.Collection{system.MustParseDate("1993-08-13")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a system.Time",
			input:   system.Collection{system.MustParseTime("14:01:45.0000001")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a system.DateTime",
			input:   system.Collection{system.MustParseDateTime("1993-08-13T14:01:45.0000001")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a system.Boolean",
			input:   system.Collection{system.Boolean(true)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a fhir.Boolean",
			input:   system.Collection{fhir.Boolean(true)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a system.Quantity",
			input:   system.Collection{system.MustParseQuantity("75", "kg")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns false for a ppb.Patient",
			input:   system.Collection{&ppb.Patient{}},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.ConvertsToString(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("ConvertsToString() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ConvertsToString() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestConvertsToTime(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.String("2001-09-11"),
				system.String("2011-05-02")},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args length is more than 0",
			input: system.Collection{system.String("2001-09-11")},
			args: []expr.Expression{
				exprtest.Return(system.String("minutes")),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "returns empty if input is empty",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "returns an empty if input is not convertible to system.Time",
			input:   system.Collection{system.MustParseQuantity("75", "Kg")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "returns true for a systemTime",
			input:   system.Collection{system.MustParseTime("16:20:59")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a partial system.Time",
			input:   system.Collection{system.MustParseTime("16:20")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns true for a convertible system.String",
			input:   system.Collection{system.String("12:59:59")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "returns false for a non convertible system.String",
			input:   system.Collection{system.String("12/59/99")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "returns empty for a ppb.Patient",
			input:   system.Collection{&ppb.Patient{}},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.ConvertsToTime(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("ConvertsToTime() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ConvertsToTime() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestToBoolean(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.String("T"),
				system.String("True")},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args length is more than 0",
			input: system.Collection{system.String("false")},
			args: []expr.Expression{
				exprtest.Return(system.String("200")),
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
			name:    "returns an empty collection if input is not convertible",
			input:   system.Collection{system.String("404 Kg")},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "input is system.Decimal '0.0'",
			input:   system.Collection{system.MustParseDecimal("0.0")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.Decimal '1.0'",
			input:   system.Collection{system.MustParseDecimal("1.0")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.Decimal '3.5'",
			input:   system.Collection{system.MustParseDecimal("3.5")},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "input is system.Integer '0'",
			input:   system.Collection{system.Integer(0)},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.Integer '1'",
			input:   system.Collection{system.Integer(1)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.Integer '3'",
			input:   system.Collection{system.Integer(3)},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "input is system.String '2.0'",
			input:   system.Collection{system.String("2.0")},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "input is system.String '1.0'",
			input:   system.Collection{system.String("1.0")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'y'",
			input:   system.Collection{system.String("y")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'yes'",
			input:   system.Collection{system.String("yes")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'Y'",
			input:   system.Collection{system.String("Y")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'YES'",
			input:   system.Collection{system.String("YES")},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.String '0.0'",
			input:   system.Collection{system.String("0.0")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'f'",
			input:   system.Collection{system.String("f")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'false'",
			input:   system.Collection{system.String("false")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'F'",
			input:   system.Collection{system.String("F")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'FALSE'",
			input:   system.Collection{system.String("FALSE")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'n",
			input:   system.Collection{system.String("n")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'no'",
			input:   system.Collection{system.String("no")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'N",
			input:   system.Collection{system.String("N")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.String 'NO'",
			input:   system.Collection{system.String("NO")},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is system.Boolean 'true'",
			input:   system.Collection{system.Boolean(true)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is system.Boolean 'false",
			input:   system.Collection{system.Boolean(false)},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
		{
			name:    "input is fhir.Integer '1'",
			input:   system.Collection{fhir.Integer(1)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is fhir.PositiveInt '1'",
			input:   system.Collection{fhir.PositiveInt(1)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
		{
			name:    "input is fhir.UnsignedInt '1'",
			input:   system.Collection{fhir.UnsignedInt(1)},
			want:    system.Collection{system.Boolean(true)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.ToBoolean(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("ToBoolean() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ToBoolean() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestToDate(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.String("2001-09-11"),
				system.String("2011-05-02")},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args length is more than 0",
			input: system.Collection{system.String("2001-09-11")},
			args: []expr.Expression{
				exprtest.Return(system.String("minutes")),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "returns empty if input is empty",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "returns an empty if input is not convertible to system.Date",
			input:   system.Collection{system.MustParseQuantity("75", "Kg")},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "returns system.Date for a system.Date",
			input:   system.Collection{system.MustParseDate("1993-08-13")},
			want:    system.Collection{system.MustParseDate("1993-08-13")},
			wantErr: false,
		},
		{
			name:    "returns system.Date for a partial system.Date",
			input:   system.Collection{system.MustParseDate("1993-08")},
			want:    system.Collection{system.MustParseDate("1993-08")},
			wantErr: false,
		},
		{
			name:    "returns system.Date for a system.DateTime",
			input:   system.Collection{system.MustParseDateTime("1993-08-13T14:20:00")},
			want:    system.Collection{system.MustParseDate("1993-08-13")},
			wantErr: false,
		},
		{
			name:    "returns system.Date for a convertible system.String",
			input:   system.Collection{system.String("1993-08-13")},
			want:    system.Collection{system.MustParseDate("1993-08-13")},
			wantErr: false,
		},
		{
			name:    "returns empty for a non convertible system.String",
			input:   system.Collection{system.String("93.08.13")},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "returns empty for a ppb.Patient",
			input:   system.Collection{&ppb.Patient{}},
			want:    system.Collection{},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.ToDate(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("ToDate() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ToDate() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestToDateTime(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.String("2001-09-11"),
				system.String("2011-05-02")},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args length is more than 0",
			input: system.Collection{system.String("2001-09-11")},
			args: []expr.Expression{
				exprtest.Return(system.String("minutes")),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "returns empty if input is empty",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "returns an empty if input is not convertible to system.DateTime",
			input:   system.Collection{system.MustParseQuantity("75", "Kg")},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "returns system.DateTime for a system.DateTime",
			input:   system.Collection{system.MustParseDateTime("1993-08-13T14:20:00")},
			want:    system.Collection{system.MustParseDateTime("1993-08-13T14:20:00")},
			wantErr: false,
		},
		{
			name:    "returns system.DateTime for a partial system.DateTime",
			input:   system.Collection{system.MustParseDateTime("2012-01-01T10:00")},
			want:    system.Collection{system.MustParseDateTime("2012-01-01T10:00")},
			wantErr: false,
		},
		{
			name:    "returns system.DateTime for a system.Date",
			input:   system.Collection{system.MustParseDate("2006-01-02")},
			want:    system.Collection{system.MustParseDate("2006-01-02").ToDateTime()},
			wantErr: false,
		},
		{
			name:    "returns system.DateTime for a convertible system.String",
			input:   system.Collection{system.String("1993-08-13T14:20:00")},
			want:    system.Collection{system.MustParseDateTime("1993-08-13T14:20:00")},
			wantErr: false,
		},
		{
			name:    "returns empty for a non convertible system.String",
			input:   system.Collection{system.String("93.08.13")},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "returns empty for a ppb.Patient",
			input:   system.Collection{&ppb.Patient{}},
			want:    system.Collection{},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.ToDateTime(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("ToDateTime() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ToDateTime() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestToDecimal(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.String("101"),
				system.String("102")},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args length is more than 0",
			input: system.Collection{system.String("100")},
			args: []expr.Expression{
				exprtest.Return(system.String("200")),
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
			name:    "returns and empty collection if input is not convertible to system.Decimal",
			input:   system.Collection{system.MustParseQuantity("500", "kg")},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "input is system.Decimal '13'",
			input:   system.Collection{system.MustParseDecimal("13")},
			want:    system.Collection{system.MustParseDecimal("13")},
			wantErr: false,
		},
		{
			name:    "input is system.Integer '13'",
			input:   system.Collection{system.Integer(13)},
			want:    system.Collection{system.MustParseDecimal("13")},
			wantErr: false,
		},
		{
			name:    "input is system.String '13'",
			input:   system.Collection{system.String("13")},
			want:    system.Collection{system.MustParseDecimal("13")},
			wantErr: false,
		},
		{
			name:    "input is system.Boolean 'true'",
			input:   system.Collection{system.Boolean(true)},
			want:    system.Collection{system.MustParseDecimal("1")},
			wantErr: false,
		},
		{
			name:    "input is system.Boolean 'false'",
			input:   system.Collection{system.Boolean(false)},
			want:    system.Collection{system.MustParseDecimal("0")},
			wantErr: false,
		},
		{
			name:    "input is fhir.Integer '0'",
			input:   system.Collection{fhir.Integer(0)},
			want:    system.Collection{system.MustParseDecimal("0")},
			wantErr: false,
		},
		{
			name:    "input is fhir.PositiveInt '1'",
			input:   system.Collection{fhir.PositiveInt(1)},
			want:    system.Collection{system.MustParseDecimal("1")},
			wantErr: false,
		},
		{
			name:    "input is fhir.UnsignedInt '2'",
			input:   system.Collection{fhir.UnsignedInt(2)},
			want:    system.Collection{system.MustParseDecimal("2")},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.ToDecimal(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("ToDecimal() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ToDecimal() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestToInteger(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.String("101"),
				system.String("102")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "errors if input is not convertible to system.Integer",
			input:   system.Collection{system.String("404 Kg")},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args length is more than 0",
			input: system.Collection{system.String("100")},
			args: []expr.Expression{
				exprtest.Return(system.String("200")),
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
			name:    "input is system.Integer '13'",
			input:   system.Collection{system.Integer(13)},
			want:    system.Collection{system.Integer(13)},
			wantErr: false,
		},
		{
			name:    "input is system.String '13'",
			input:   system.Collection{system.String("13")},
			want:    system.Collection{system.Integer(13)},
			wantErr: false,
		},
		{
			name:    "input is system.Boolean 'true'",
			input:   system.Collection{system.Boolean(true)},
			want:    system.Collection{system.Integer(1)},
			wantErr: false,
		},
		{
			name:    "input is system.Boolean 'false'",
			input:   system.Collection{system.Boolean(false)},
			want:    system.Collection{system.Integer(0)},
			wantErr: false,
		},
		{
			name:    "input is fhir.Integer '10'",
			input:   system.Collection{fhir.Integer(10)},
			want:    system.Collection{system.Integer(10)},
			wantErr: false,
		},
		{
			name:    "input is fhir.PositiveInt '11'",
			input:   system.Collection{fhir.PositiveInt(11)},
			want:    system.Collection{system.Integer(11)},
			wantErr: false,
		},
		{
			name:    "input is fhir.UnsignedInt '12'",
			input:   system.Collection{fhir.UnsignedInt(12)},
			want:    system.Collection{system.Integer(12)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.ToInteger(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("ToInteger() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ToInteger() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestToQuantity(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.String("10 'km/hr'"),
				system.String("10 'mi/hr'")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "returns an empty collection if input is not convertible to a system.Quantity",
			input:   system.Collection{system.String("10 km / hr")},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:  "errors if args length is more than 1",
			input: system.Collection{system.String("2 days")},
			args: []expr.Expression{
				exprtest.Return(system.String("hours")),
				exprtest.Return(system.String("minutes")),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args is not a valid unit of time",
			input: system.Collection{system.String("100 years")},
			args: []expr.Expression{
				exprtest.Return(system.String("decades")),
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
			name:    "input is system.Integer '13'",
			input:   system.Collection{system.Integer(13)},
			want:    system.Collection{system.MustParseQuantity("13", "1")},
			wantErr: false,
		},
		{
			name:    "input is system.Decimal '13.5",
			input:   system.Collection{system.MustParseDecimal("13.5")},
			want:    system.Collection{system.MustParseQuantity("13.5", "1")},
			wantErr: false,
		},
		{
			name:    "input is system.Quantity '13.5 lbs'",
			input:   system.Collection{system.MustParseQuantity("13.5", "lbs")},
			want:    system.Collection{system.MustParseQuantity("13.5", "lbs")},
			wantErr: false,
		},
		{
			name:    "input is system.String '100 days'",
			input:   system.Collection{system.String("100 days")},
			want:    system.Collection{system.MustParseQuantity("100", "days")},
			wantErr: false,
		},
		{
			name:    "input is system.String '100 km'",
			input:   system.Collection{system.String("100 km")},
			want:    system.Collection{system.MustParseQuantity("100", "km")},
			wantErr: false,
		},
		{
			name:    "input is system.String '100 km/h'",
			input:   system.Collection{system.String("100 km/h")},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "input is system.String '100 'km/h''",
			input:   system.Collection{system.String("100 'km/h'")},
			want:    system.Collection{system.MustParseQuantity("100", "km/h")},
			wantErr: false,
		},
		{
			name:    "input is system.String '100 'km per hour''",
			input:   system.Collection{system.String("100 'km per hour'")},
			want:    system.Collection{system.MustParseQuantity("100", "km per hour")},
			wantErr: false,
		},
		{
			name:    "input is system.String '100 km per hour'",
			input:   system.Collection{system.String("100 km per hour")},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "input is system.String '100           km'",
			input:   system.Collection{system.String("100           km")},
			want:    system.Collection{system.MustParseQuantity("100", "          km")},
			wantErr: false,
		},
		{
			name:    "input is system.String '100           km/h'",
			input:   system.Collection{system.String("100           km/h")},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "input is system.String '10 'km per hr''",
			input:   system.Collection{system.String("10 'km per hr'")},
			want:    system.Collection{system.MustParseQuantity("10", "km per hr")},
			wantErr: false,
		},
		{
			name:  "input is system.String '2 years' with arg 'months'",
			input: system.Collection{system.String("2 years")},
			args: []expr.Expression{
				exprtest.Return(system.String("months")),
			},
			want:    system.Collection{system.MustParseQuantity("24", "months")},
			wantErr: false,
		},
		{
			name:  "input is system.String '2 months' with arg 'days'",
			input: system.Collection{system.String("2 months")},
			args: []expr.Expression{
				exprtest.Return(system.String("days")),
			},
			want:    system.Collection{system.MustParseQuantity("60", "days")},
			wantErr: false,
		},
		{
			name:  "input is system.String '2 days' with arg 'hours'",
			input: system.Collection{system.String("2 days")},
			args: []expr.Expression{
				exprtest.Return(system.String("hours")),
			},
			want:    system.Collection{system.MustParseQuantity("48", "hours")},
			wantErr: false,
		},
		{
			name:  "input is system.String '2 hours' with arg 'minutes'",
			input: system.Collection{system.String("2 hours")},
			args: []expr.Expression{
				exprtest.Return(system.String("minutes")),
			},
			want:    system.Collection{system.MustParseQuantity("120", "minutes")},
			wantErr: false,
		},
		{
			name:  "input is system.String '2 minutes' with arg 'seconds'",
			input: system.Collection{system.String("2 minutes")},
			args: []expr.Expression{
				exprtest.Return(system.String("seconds")),
			},
			want:    system.Collection{system.MustParseQuantity("120", "seconds")},
			wantErr: false,
		},
		{
			name:    "input is system.String '100 'km''",
			input:   system.Collection{system.String("100 'km'")},
			want:    system.Collection{system.MustParseQuantity("100", "km")},
			wantErr: false,
		},
		{
			name:  "input is system.Integer '100' with arg ''km''",
			input: system.Collection{system.Integer(100)},
			args: []expr.Expression{
				exprtest.Return(system.String("'km'")),
			},
			want:    system.Collection{system.MustParseQuantity("100", "km")},
			wantErr: false,
		},
		{
			name:  "input is system.Integer '100' with arg 'days''",
			input: system.Collection{system.Integer(100)},
			args: []expr.Expression{
				exprtest.Return(system.String("days")),
			},
			want:    system.Collection{system.MustParseQuantity("100", "days")},
			wantErr: false,
		},
		{
			name:    "input is system.Boolean 'true'",
			input:   system.Collection{system.Boolean(true)},
			want:    system.Collection{system.MustParseQuantity("1.0", "1")},
			wantErr: false,
		},
		{
			name:    "input is system.Boolean 'false'",
			input:   system.Collection{system.Boolean(false)},
			want:    system.Collection{system.MustParseQuantity("0.0", "1")},
			wantErr: false,
		},
		{
			name:    "input is fhir.Integer '10'",
			input:   system.Collection{fhir.Integer(10)},
			want:    system.Collection{system.MustParseQuantity("10", "1")},
			wantErr: false,
		},
		{
			name:    "input is fhir.PositiveInt '11'",
			input:   system.Collection{fhir.PositiveInt(11)},
			want:    system.Collection{system.MustParseQuantity("11", "1")},
			wantErr: false,
		},
		{
			name:    "input is fhir.UnsignedInt '12'",
			input:   system.Collection{fhir.UnsignedInt(12)},
			want:    system.Collection{system.MustParseQuantity("12", "1")},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.ToQuantity(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("ToQuantity() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ToQuantity() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestToString(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if input is not a singleton",
			input: system.Collection{
				system.String("101"),
				system.String("102"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "returns empty for and empty input collection",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "returns a system.String for a system.String",
			input:   system.Collection{system.String("100")},
			want:    system.Collection{system.String("100")},
			wantErr: false,
		},
		{
			name:    "returns a system.String for a fhir.String",
			input:   system.Collection{fhir.String("100")},
			want:    system.Collection{system.String("100")},
			wantErr: false,
		},
		{
			name:    "returns a system.String for a system.Integer",
			input:   system.Collection{system.Integer(100)},
			want:    system.Collection{system.String("100")},
			wantErr: false,
		},
		{
			name:    "returns system.String for a fhir.Integer",
			input:   system.Collection{fhir.Integer(100)},
			want:    system.Collection{system.String("100")},
			wantErr: false,
		},
		{
			name:    "returns system.String for a fhir.PositiveInt",
			input:   system.Collection{fhir.PositiveInt(11)},
			want:    system.Collection{system.String("11")},
			wantErr: false,
		},
		{
			name:    "returns system.String for a fhir.UnsignedInt",
			input:   system.Collection{fhir.UnsignedInt(12)},
			want:    system.Collection{system.String("12")},
			wantErr: false,
		},
		{
			name:    "returns system.String for a system.Decimal",
			input:   system.Collection{system.MustParseDecimal("100.999")},
			want:    system.Collection{system.String("100.999")},
			wantErr: false,
		},
		{
			name:    "returns system.String for a system.Date",
			input:   system.Collection{system.MustParseDate("1993-08-13")},
			want:    system.Collection{system.String("1993-08-13")},
			wantErr: false,
		},
		{
			name:    "returns system.String for a system.Time",
			input:   system.Collection{system.MustParseTime("14:01:45")},
			want:    system.Collection{system.String("14:01:45")},
			wantErr: false,
		},
		{
			name:    "returns system.String for a system.DateTime",
			input:   system.Collection{system.MustParseDateTime("1993-08-13T14:01:45")},
			want:    system.Collection{system.String("1993-08-13T14:01:45")},
			wantErr: false,
		},
		{
			name:    "returns system.String for a system.Boolean",
			input:   system.Collection{system.Boolean(true)},
			want:    system.Collection{system.String("true")},
			wantErr: false,
		},
		{
			name:    "returns system.String for a fhir.Boolean",
			input:   system.Collection{fhir.Boolean(true)},
			want:    system.Collection{system.String("true")},
			wantErr: false,
		},
		{
			name:    "returns system.String for a system.Quantity",
			input:   system.Collection{system.MustParseQuantity("75", "kg")},
			want:    system.Collection{system.String("75 kg")},
			wantErr: false,
		},
		{
			name:    "returns false for a ppb.Patient",
			input:   system.Collection{&ppb.Patient{}},
			want:    system.Collection{system.Boolean(false)},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.ToString(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("ToString() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ToString() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestToTime(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		args    []expr.Expression
		want    system.Collection
		wantErr bool
	}{
		{
			name: "errors if input length is more than 1",
			input: system.Collection{
				system.String("2001-09-11"),
				system.String("2011-05-02")},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "errors if args length is more than 0",
			input: system.Collection{system.String("2001-09-11")},
			args: []expr.Expression{
				exprtest.Return(system.String("minutes")),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "returns empty if input is empty",
			input:   system.Collection{},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "returns an empty if input is not convertible to system.Time",
			input:   system.Collection{system.MustParseQuantity("75", "Kg")},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "returns system.Time for a systemTime",
			input:   system.Collection{system.MustParseTime("16:20:59")},
			want:    system.Collection{system.MustParseTime("16:20:59")},
			wantErr: false,
		},
		{
			name:    "returns system.Time for a partial system.Time",
			input:   system.Collection{system.MustParseTime("16:20")},
			want:    system.Collection{system.MustParseTime("16:20")},
			wantErr: false,
		},
		{
			name:    "returns system.Time for a convertible system.String",
			input:   system.Collection{system.String("12:59:59")},
			want:    system.Collection{system.MustParseTime("12:59:59")},
			wantErr: false,
		},
		{
			name:    "returns empty for a non convertible system.String",
			input:   system.Collection{system.String("12/59/99")},
			want:    system.Collection{},
			wantErr: false,
		},
		{
			name:    "returns empty for a ppb.Patient",
			input:   system.Collection{&ppb.Patient{}},
			want:    system.Collection{},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.ToTime(&expr.Context{}, tc.input, tc.args...)
			if (err != nil) != tc.wantErr {
				t.Errorf("ToTime() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ToTime() returned unexpected diff (-want, +got)\n%s", diff)
			}
		})
	}
}
