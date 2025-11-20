package trace_test

import (
	"context"
	"strings"
	"testing"

	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	ppb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"github.com/google/go-cmp/cmp"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/fhirpath/trace"
	"github.com/verily-src/fhirpath-go/internal/fhir"
)

func TestWriterTracer(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		prefix string
		input  system.Collection
		want   string
	}{
		{
			name:  "Input is empty collection",
			input: system.Collection{},
			want:  "",
		}, {
			name:   "Input is system string",
			prefix: "test-prefix",
			input: system.Collection{
				system.String("test string"),
			},
			want: "test-prefix:\n\"test string\"\n",
		}, {
			name:   "Input is system integer",
			prefix: "test-prefix",
			input: system.Collection{
				system.Integer(42),
			},
			want: "test-prefix:\n42\n",
		}, {
			name:   "Input is system boolean",
			prefix: "test-prefix",
			input: system.Collection{
				system.Boolean(true),
			},
			want: "test-prefix:\ntrue\n",
		}, {
			name:   "Input is date",
			prefix: "test-prefix",
			input: system.Collection{
				system.MustParseDate("2024-01-02"),
			},
			want: "test-prefix:\n\"2024-01-02\"\n",
		}, {
			name:   "Input is resource",
			prefix: "test-prefix",
			input: system.Collection{
				&ppb.Patient{
					Name: []*dtpb.HumanName{
						{
							Given: []*dtpb.String{fhir.String("John")},
						},
					},
				},
			},
			want: `test-prefix:
{
  "name": [
    {
      "given": [
        "John"
      ]
    }
  ],
  "resourceType": "Patient"
}
`,
		}, {
			name:   "Multiple inputs",
			prefix: "multiple",
			input: system.Collection{
				system.String("string value"),
				system.Integer(100),
				system.Boolean(false),
			},
			want: `multiple:
"string value"
100
false
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var sb strings.Builder
			sut := &trace.WriterTracer{
				Out: &sb,
			}
			ctx := context.Background()

			sut.Trace(ctx, tc.prefix, tc.input)

			if got, want := sb.String(), tc.want; !cmp.Equal(got, want) {
				t.Errorf("Trace() mismatch (-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}
