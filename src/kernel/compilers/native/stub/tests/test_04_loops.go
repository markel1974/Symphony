package tests

import "fmt"

func main() {
	fmt.Println("--- Running Test: Loops (for, range) ---")

	sum := 0
	for i := 0; i < 10; i = i + 1 {
		sum = sum + i
	}

	arr := []int{10, 20, 30}
	arrSum := 0
	for _, v := range arr {
		arrSum = arrSum + v
	}

	m := map[string]int{"a": 1, "b": 10, "c": 100}
	mapSum := 0
	for _, val := range m {
		mapSum = mapSum + val
	}

	finalValue := sum + arrSum + mapSum
	expectedValue := 216

	if finalValue == expectedValue {
		fmt.Println("[TEST PASSED] Loops executed and aggregated correctly.")
	} else {
		fmt.Println("[TEST FAILED] Error during loop execution.")
	}
}
