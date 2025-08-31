package tests

import "fmt"

func increment(val *int) {
	*val = *val + 1
}

func main() {
	fmt.Println("--- Running Test: Pointers ---")

	x := 10
	ptrX := &x
	increment(ptrX)

	y := 20
	ptrY := &y
	increment(ptrY)

	finalValue := x + y
	expectedValue := 32

	if finalValue == expectedValue {
		fmt.Println("[TEST PASSED] Pointers and dereferencing worked correctly.")
	} else {
		fmt.Println("[TEST FAILED] Error in pointer manipulation.")
	}
}
