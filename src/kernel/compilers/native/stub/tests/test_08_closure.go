package tests

import "fmt"

// Questa funzione restituisce una closure che cattura la variabile 'count'.
func counterFactory(start int) func() int {
	count := start
	return func() int {
		count = count + 1
		return count
	}
}

func main() {
	fmt.Println("--- Running Test: Closures and Free Variables ---")

	// Creiamo due contatori indipendenti dalla stessa factory.
	counterA := counterFactory(10)
	counterB := counterFactory(100)

	// Chiamiamo il primo contatore tre volte.
	counterA()
	counterA()
	v1 := counterA() // v1 dovrebbe essere 13

	// Chiamiamo il secondo contatore due volte.
	counterB()
	v2 := counterB() // v2 dovrebbe essere 102

	finalValue := v1 + v2
	expectedValue := 115 // 13 + 102

	if finalValue == expectedValue {
		fmt.Println("[TEST PASSED] Closures captured and modified free variables correctly.")
	} else {
		fmt.Printf("[TEST FAILED] Error in closure state management. Got: %d, Expected: %d\n", finalValue, expectedValue)
	}
}
