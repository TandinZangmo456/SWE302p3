// shipping_v2.go
package shipping

import (
	"errors"
	"fmt"
)

// CalculateShippingFeeV2 calculates the fee based on new tiered logic.
func CalculateShippingFeeV2(weight float64, zone string, insured bool) (float64, error) {
	if weight <= 0 || weight > 50 {
		return 0, errors.New("invalid weight")
	}

	var baseFee float64
	switch zone {
	case "Domestic":
		baseFee = 5.0
	case "International":
		baseFee = 20.0
	case "Express":
		baseFee = 30.0
	default:
		return 0, fmt.Errorf("invalid zone: %s", zone)
	}

	// Add weight-based per-kg cost
	var perKgCost float64
	switch zone {
	case "Domestic":
		perKgCost = weight * 1.0
	case "International":
		perKgCost = weight * 2.5
	case "Express":
		perKgCost = weight * 5.0
	}

	var heavySurcharge float64
	if weight > 10 {
		heavySurcharge = 7.50
	}

	subTotal := baseFee + perKgCost + heavySurcharge

	var insuranceCost float64
	if insured {
		insuranceCost = subTotal * 0.015
	}

	finalTotal := subTotal + insuranceCost

	return finalTotal, nil
}