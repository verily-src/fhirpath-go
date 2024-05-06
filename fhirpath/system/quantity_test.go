package system_test

import (
	"testing"

	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

func TestParseQuantity_ErrorsOnInvalidString(t *testing.T) {
	testCases := []struct {
		name   string
		number string
		unit   string
	}{
		{
			name:   "non-number",
			number: "word",
			unit:   "kg",
		},
		{
			name:   "empty strings",
			number: "",
			unit:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := system.ParseQuantity(tc.number, tc.unit); err == nil {
				t.Fatalf("ParseQuantity(%s, %s) didn't raise error when expected", tc.number, tc.unit)
			}
		})
	}
}

func TestParseQuantity_ReturnsQuantity(t *testing.T) {
	testCases := []struct {
		name   string
		number string
		unit   string
	}{
		{
			name:   "unit quantity",
			number: "100",
			unit:   "lbs",
		},
		{
			name:   "time quantity",
			number: "3",
			unit:   "minutes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := system.ParseQuantity(tc.number, tc.unit); err != nil {
				t.Fatalf("ParseQuantity(%s, %s) returned unexpected error: %v", tc.number, tc.unit, err)
			}
		})
	}
}

func TestQuantity_Equal(t *testing.T) {
	onePound, _ := system.ParseQuantity("1", "lbs")
	oneLb, _ := system.ParseQuantity("1", "lbs")
	oneKg, _ := system.ParseQuantity("1", "kg")
	twoLbs, _ := system.ParseQuantity("2", "lbs")

	testCases := []struct {
		name        string
		quantityOne system.Quantity
		quantityTwo system.Any
		shouldEqual bool
		wantOk      bool
	}{
		{
			name:        "same quantity different objects",
			quantityOne: onePound,
			quantityTwo: oneLb,
			shouldEqual: true,
			wantOk:      true,
		},
		{
			name:        "same number different unit",
			quantityOne: oneLb,
			quantityTwo: oneKg,
			wantOk:      false,
		},
		{
			name:        "same unit different number",
			quantityOne: oneLb,
			quantityTwo: twoLbs,
			shouldEqual: false,
			wantOk:      true,
		},
		{
			name:        "different type",
			quantityOne: oneLb,
			quantityTwo: system.String("1 lbs"),
			shouldEqual: false,
			wantOk:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := tc.quantityOne.TryEqual(tc.quantityTwo)

			if ok != tc.wantOk {
				t.Fatalf("Quantity.Equal: ok got %v, want %v", ok, tc.wantOk)
			}
			if got != tc.shouldEqual {
				t.Errorf("Quantity.Equal returned unexpected equality: got %v, want %v", got, tc.shouldEqual)
			}
		})
	}
}
