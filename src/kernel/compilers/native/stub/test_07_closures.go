package stub

import "fmt"

func counterFactory(start int) func() int {
	count := start
	return func() int {
		count = count + 1
		return count
	}
}

func main() {
	fmt.Println("--- Running Test 07: Closures and Free Variables ---")

	counterA := counterFactory(10)
	counterB := counterFactory(100)

	counterA()
	counterA()
	v1 := counterA()

	counterB()
	v2 := counterB()

	finalValue := v1 + v2
	expectedValue := 115

	if finalValue == expectedValue {
		fmt.Println("[TEST PASSED] Closures captured and modified free variables correctly.")
	} else {
		fmt.Println("[TEST FAILED] Error in closure state management.")
	}
}
