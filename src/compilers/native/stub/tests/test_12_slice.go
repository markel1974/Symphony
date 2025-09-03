package tests

import "fmt"

func main() {
	fmt.Println("--- Running Test: Data Structures (Array, Map, Slice) ---")

	arr := []int{10, 20, 30}
	val := arr[1]

	m := map[string]int{"one": 1, "two": 2}
	m["three"] = 3
	mapVal := m["two"] + m["three"]

	slice := arr[0:2]
	sliceSum := slice[0] + slice[1]

	finalValue := val + mapVal + sliceSum
	expectedValue := 55

	if finalValue == expectedValue {
		fmt.Println("[TEST PASSED] Data structures handled correctly.")
	} else {
		fmt.Printf("[TEST FAILED] Error in data structure manipulation. %d\n", finalValue)
	}
}
