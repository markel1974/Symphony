package stub

import "fmt"

func add(a int, b int) int {
	return a + b
}

func factorial(n int) int {
	if n == 0 {
		return 1
	}
	f := factorial(n - 1)
	z := n * f
	fmt.Println("factorial", z)
	return z
}

func main() {
	fmt.Println("--- Running Test 05: Functions and Recursion ---")

	sum := add(10, 20)
	fact5 := factorial(5)

	finalValue := fact5 + sum
	expectedValue := 150

	if finalValue == expectedValue {
		fmt.Println("[TEST PASSED] Functions and recursion worked correctly.")
	} else {
		fmt.Println("[TEST FAILED] Mismatch in function call or recursion result.")
	}
}
