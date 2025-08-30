package stub

import "fmt"

func checkValue(v int) int {
	if v > 100 {
		return 1
	} else if v < 100 && v > 50 {
		return 2
	} else if v == 100 {
		return 3
	} else {
		return 4
	}
}

func main() {
	fmt.Println("--- Running Test 03: Control Flow (if-else) ---")

	r1 := checkValue(200)
	r2 := checkValue(75)
	r3 := checkValue(100)
	r4 := checkValue(10)

	finalValue := r1*1000 + r2*100 + r3*10 + r4
	expectedValue := 1234

	if finalValue == expectedValue {
		fmt.Println("[TEST PASSED] All control flow branches worked correctly.")
	} else {
		fmt.Println("[TEST FAILED] Mismatch in control flow logic.")
	}
}
