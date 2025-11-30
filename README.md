# Shipping Calculator V2 - Test Design Analysis

## Part 1: Equivalence Partitioning

### Weight Partitions:
- **P1 (Invalid - Too Small):** weight ≤ 0
  - Example values: -5, 0
  - Expected: Error "invalid weight"

- **P2 (Valid - Standard Package):** 0 < weight ≤ 10
  - Example values: 0.1, 5, 10
  - Expected: No heavy surcharge

- **P3 (Valid - Heavy Package):** 10 < weight ≤ 50
  - Example values: 10.1, 25, 50
  - Expected: $7.50 heavy surcharge added

- **P4 (Invalid - Too Large):** weight > 50
  - Example values: 50.1, 100
  - Expected: Error "invalid weight"

### Zone Partitions:
- **P5 (Valid Zones):** "Domestic", "International", "Express"
  - Expected: Correct base fee applied

- **P6 (Invalid Zones):** Any other string
  - Example values: "Local", "", "domestic" (lowercase)
  - Expected: Error "invalid zone"

### Insured Partitions:
- **P7 (Insured = true):** Insurance cost = 1.5% of subtotal
  - Expected: Additional 1.5% fee added

- **P8 (Insured = false):** No insurance
  - Expected: No additional fee

## Part 2: Boundary Value Analysis

### Weight Boundaries:

#### Lower Boundary (around 0):
- **0** - Last invalid value (should error)
- **0.1** - First valid value (Standard package)

#### Mid Boundary (around 10 - where Standard becomes Heavy):
- **10** - Last value for Standard package (no surcharge)
- **10.1** - First value for Heavy package ($7.50 surcharge)

#### Upper Boundary (around 50):
- **50** - Last valid value (Heavy package)
- **50.1** - First invalid value (should error)

### Why These Boundaries Matter:
These are the exact points where behavior changes. Developers often make "off-by-one" errors at these transitions (e.g., using `>=` instead of `>`), so we must test both sides of each boundary.

## Part 3: Test Case Design Rationale

### Test Strategy:
1. Test all weight partitions with valid zones
2. Test boundary values for weight
3. Test insurance flag with different scenarios
4. Test invalid inputs (weight and zone)
5. Combine conditions to catch interaction bugs

### Example Calculations:

**Example 1: Standard Domestic, Insured**
- Weight: 5 kg, Zone: "Domestic", Insured: true
- Base: $5.00
- Heavy Surcharge: $0 (5 ≤ 10)
- Subtotal: $5.00
- Insurance: $5.00 × 0.015 = $0.075
- **Final: $5.075**

**Example 2: Heavy International, Not Insured**
- Weight: 20 kg, Zone: "International", Insured: false
- Base: $20.00
- Heavy Surcharge: $7.50 (20 > 10)
- Subtotal: $27.50
- Insurance: $0
- **Final: $27.50**

**Example 3: Boundary at 10 kg**
- Weight: 10 kg = Standard (no surcharge)
- Weight: 10.1 kg = Heavy (+$7.50 surcharge)# SWE302p3
