// shipping_v2_test.go
package shipping

import (
	"testing"
)

func TestCalculateShippingFeeV2(t *testing.T) {
	testCases := []struct {
		name        string
		weight      float64
		zone        string
		insured     bool
		expectedFee float64
		expectError bool
	}{
		// ===== INVALID WEIGHT TESTS (P1 & P4) =====
		{"Invalid: Weight too small (negative)", -5, "Domestic", false, 0, true},
		{"Invalid: Weight exactly 0", 0, "Domestic", false, 0, true},
		{"Invalid: Weight too large", 50.1, "International", false, 0, true},
		{"Invalid: Weight way too large", 100, "Express", false, 0, true},

		// ===== BOUNDARY VALUE TESTS =====
		// Lower boundary (around 0)
		{"Boundary: Weight 0.1 (just valid)", 0.1, "Domestic", false, 5.1, false},
		
		// Mid boundary (around 10 - Standard to Heavy transition)
		{"Boundary: Weight exactly 10 (Standard)", 10, "Domestic", false, 15.0, false}, // 5 + 10*1 = 15, no surcharge
		{"Boundary: Weight 10.1 (Heavy)", 10.1, "Domestic", false, 22.6, false}, // 5 + 10.1*1 + 7.5 = 22.6
		
		// Upper boundary (around 50)
		{"Boundary: Weight exactly 50 (valid Heavy)", 50, "International", false, 152.5, false}, // 20 + 125 + 7.5 = 152.5
		{"Boundary: Weight 50.1 (invalid)", 50.1, "Express", false, 0, true},

		// ===== STANDARD PACKAGE TESTS (P2: 0 < weight <= 10) =====
		{"Standard Domestic, Not Insured", 5, "Domestic", false, 10.0, false}, // 5 + 5*1 = 10
		{"Standard International, Not Insured", 8, "International", false, 40.0, false}, // 20 + 8*2.5 = 40
		{"Standard Express, Not Insured", 3, "Express", false, 45.0, false}, // 30 + 3*5 = 45

		// ===== HEAVY PACKAGE TESTS (P3: 10 < weight <= 50) =====
		{"Heavy Domestic, Not Insured", 20, "Domestic", false, 32.5, false}, // 5 + 20*1 + 7.5 = 32.5
		{"Heavy International, Not Insured", 25, "International", false, 90.0, false}, // 20 + 25*2.5 + 7.5 = 90
		{"Heavy Express, Not Insured", 15, "Express", false, 112.5, false}, // 30 + 15*5 + 7.5 = 112.5

		// ===== INSURANCE TESTS (P7: insured = true) =====
		{"Standard Domestic, Insured", 5, "Domestic", true, 10.15, false}, // (5+5)*1.015 = 10.15
		{"Heavy Domestic, Insured", 20, "Domestic", true, 32.9875, false}, // (5+20+7.5)*1.015 = 32.9875
		{"Standard International, Insured", 8, "International", true, 40.6, false}, // (20+20)*1.015 = 40.6
		{"Heavy International, Insured", 25, "International", true, 91.35, false}, // (20+62.5+7.5)*1.015 = 91.35
		{"Standard Express, Insured", 3, "Express", true, 45.675, false}, // (30+15)*1.015 = 45.675
		{"Heavy Express, Insured", 15, "Express", true, 114.1875, false}, // (30+75+7.5)*1.015 = 114.1875

		// ===== INVALID ZONE TESTS (P6) =====
		{"Invalid: Zone Unknown", 10, "Unknown", false, 0, true},
		{"Invalid: Zone empty string", 10, "", false, 0, true},
		{"Invalid: Zone lowercase", 10, "domestic", false, 0, true},
		{"Invalid: Zone Local", 15, "Local", false, 0, true},

		// ===== COMBINED EDGE CASES =====
		{"Edge: Minimum valid weight, cheapest zone, no insurance", 0.1, "Domestic", false, 5.1, false},
		{"Edge: Maximum valid weight, most expensive zone, insured", 50, "Express", true, 291.8125, false},
		// 30 + 250 + 7.5 = 287.5, then 287.5 * 1.015 = 291.8125
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fee, err := CalculateShippingFeeV2(tc.weight, tc.zone, tc.insured)

			// Check error expectation
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				}
				return // Don't check fee if we expected an error
			}

			// Check for unexpected errors
			if err != nil {
				t.Fatalf("Expected no error, but got: %v", err)
			}

			// Check fee calculation (allowing small floating point differences)
			tolerance := 0.01
			if diff := fee - tc.expectedFee; diff < -tolerance || diff > tolerance {
				t.Errorf("Expected fee %.4f, but got %.4f (difference: %.4f)", 
					tc.expectedFee, fee, diff)
			}
		})
	}
}