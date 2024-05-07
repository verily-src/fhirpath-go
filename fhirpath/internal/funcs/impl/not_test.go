package impl_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/internal/fhir"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs/impl"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

func TestNot_InvertsBoolean(t *testing.T) {
	testCases := []struct {
		name    string
		input   system.Collection
		want    system.Collection
		wantErr bool
	}{
		{
			name:  "inverts system boolean",
			input: system.Collection{system.Boolean(true)},
			want:  system.Collection{system.Boolean(false)},
		},
		{
			name:  "inverts proto boolean",
			input: system.Collection{fhir.Boolean(false)},
			want:  system.Collection{system.Boolean(true)},
		},
		{
			name:    "receives non-singleton collection",
			input:   system.Collection{system.Boolean(true), system.Boolean(false)},
			wantErr: true,
		},
		{
			name:  "passes through empty collection",
			input: system.Collection{},
			want:  system.Collection{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := impl.Not(&expr.Context{}, tc.input)

			if gotErr := err != nil; tc.wantErr != gotErr {
				t.Fatalf("Not function got unexpected error result: gotErr %v, wantErr %v, err: %v", gotErr, tc.wantErr, err)
			}
			if !cmp.Equal(tc.want, got) {
				t.Errorf("Not function returned unexpected result: got: %v, want %v", got, tc.want)
			}
		})
	}
}
