package stub

import "fmt"

func increment(val *int) {
	*val = *val + 1
}

func main() {
	fmt.Println("--- Running Test 09: Pointers ---")

	x := 10
	ptr_x := &x
	increment(ptr_x)

	y := 20
	ptr_y := &y
	increment(ptr_y)

	finalValue := x + y
	expectedValue := 32

	if finalValue == expectedValue {
		fmt.Println("[TEST PASSED] Pointers and dereferencing worked correctly.")
	} else {
		fmt.Println("[TEST FAILED] Error in pointer manipulation.")
	}
}
