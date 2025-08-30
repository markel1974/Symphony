package stub

import "fmt"

func main() {
	fmt.Println("--- Running Test 01: Expressions and Operators ---")

	a := 10
	b := 20
	c := 5

	// Aritmetica e precedenza
	result1 := a + b*c // 110
	//TODO *ast.ParenExpr
	//result2 := (a + b) * c // 150
	tmp := a + b
	result2 := tmp * c

	// Espressioni booleane
	ok := true
	nok := false
	boolResult := false
	if !ok || nok {
		boolResult = true
	}
	if (10 > 5) && (20 == 20) {
		// ok
	} else {
		boolResult = true
	}

	finalValue := (result2 - result1) + 5
	expectedValue := 45

	if finalValue == expectedValue && !boolResult {
		fmt.Println("[TEST PASSED] Expressions evaluated correctly.")
	} else {
		fmt.Println("[TEST FAILED] Mismatch in expression evaluation.")
	}
}
