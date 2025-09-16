package tests

import "fmt"

func main() {
	fmt.Println("--- Running Test: Code Structures (Array, Map, Slice) ---")
	arr := []int{10, 20, 30}
	val := arr[1]
	m := map[string]int{"one": 1, "two": 2, "three": 3}
	mapVal := m["two"] + m["three"]
	sliceSum := arr[0] + arr[1]
	finalValue := val + mapVal + sliceSum
	expectedValue := 55
	if finalValue == expectedValue {
		fmt.Println("[TEST PASSED] Code structures handled correctly.")
	} else {
		fmt.Println("[TEST FAILED] Error in data structure manipulation.", finalValue, expectedValue)
	}
}
