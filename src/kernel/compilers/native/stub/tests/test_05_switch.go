package tests

import "fmt"

func testSwitch(value int) string {
	result := ""
	switch value {
	case 1:
		result = "one"
	case 3:
		result = "three"
	case 5:
		result = "five"
	default:
		result = "other"
	}
	return result
}

func main() {
	fmt.Println("--- Running Test: Control Flow (switch) ---")

	// Test a specific case
	r1 := testSwitch(3) // "three"

	// Test another case
	r2 := testSwitch(1) // "one"

	// Test the default case
	r3 := testSwitch(100) // "other"

	// Test a case that is defined but not explicitly hit in this sequence
	r4 := testSwitch(5) // "five"

	// Concatenate results to check the overall flow
	finalValue := r1 + "-" + r2 + "-" + r3 + "-" + r4
	expectedValue := "three-one-other-five"

	if finalValue == expectedValue {
		fmt.Println("[TEST PASSED] Switch statement worked correctly.")
	} else {
		fmt.Printf("[TEST FAILED] Mismatch in switch logic. Got: %s, Expected: %s\n", finalValue, expectedValue)
	}
}
