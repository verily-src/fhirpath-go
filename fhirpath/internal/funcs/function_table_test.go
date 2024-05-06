package funcs_test

import (
	"testing"

	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

func TestRegister_RaisesError(t *testing.T) {
	testCases := []struct {
		name     string
		funcName string
		fn       any
	}{
		{
			name:     "raises error when trying to override existing function",
			funcName: "where",
			fn:       func(system.Collection) (system.Collection, error) { return nil, nil },
		},
		{
			name:     "raises error when adding invalid function",
			funcName: "someFn",
			fn:       func() {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			table := funcs.Clone()

			if err := table.Register(tc.funcName, tc.fn); err == nil {
				t.Fatalf("FunctionTable.Register(%s) doesn't raise error when expected", tc.funcName)
			}
		})
	}
}

func TestRegister_AddsToMap(t *testing.T) {
	table := funcs.Clone()
	fn := func(system.Collection) (system.Collection, error) { return nil, nil }

	if err := table.Register("someFn", fn); err != nil {
		t.Fatalf("FunctionTable.Register raised unexpected error: %v", err)
	}
	if _, ok := table["someFn"]; !ok {
		t.Errorf("FunctionTable.Register did not successfully add function to map")
	}
}
