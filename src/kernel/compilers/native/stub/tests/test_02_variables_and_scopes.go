package tests

import "fmt"

var globalVar int = 100

func main() {
	fmt.Println("--- Running Test: Variables and Scopes ---")

	res := globalVar
	localVar := 20
	res = res + localVar

	if true {
		innerVar := 5
		res = res - innerVar
	}

	res = res + localVar

	expectedValue := 135

	if res == expectedValue {
		fmt.Println("[TEST PASSED] Scopes and variables handled correctly.")
	} else {
		fmt.Println("[TEST FAILED] Error in scope or variable access.")
	}
}
