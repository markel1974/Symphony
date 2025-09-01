package tests

import (
	"errors"
	"fmt"
)

// Una funzione che può restituire un valore o un errore.
func processData(succeed bool) (string, error) {
	if succeed {
		return "SuccessData", nil
	}
	return "", errors.New("processing failed")
}

func main() {
	fmt.Println("--- Running Test: Error Handling Pattern ---")

	var resultLog string

	// Caso 1: Successo
	data1, err1 := processData(true)
	if err1 != nil {
		resultLog = "Test failed on success case"
	} else {
		resultLog = data1
	}
	// Caso 2: Fallimento
	_, err2 := processData(false)
	if err2 != nil {
		fmt.Printf("Error: %s\n", err2)
		resultLog = resultLog + " | " + err2 //.Error()
	} else {
		resultLog = "Test failed on error case"
	}

	finalValue := resultLog
	expectedValue := "SuccessData | processing failed"

	if finalValue == expectedValue {
		fmt.Println("[TEST PASSED] Error handling pattern (if err != nil) worked correctly.")
	} else {
		fmt.Printf("[TEST FAILED] Mismatch in error handling logic.\nGot: %s\nExpected: %s\n", finalValue, expectedValue)
	}
}
