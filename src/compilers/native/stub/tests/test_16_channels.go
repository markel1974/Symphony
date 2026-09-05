package tests

import (
	"fmt"
)

func worker(c chan int, val int) {
	c <- val * 2
}

func main() {
	fmt.Println("--- Running Test: Channels Pattern ---")

	// 1. Test Buffered Channel
	bufChan := make(chan int, 2)
	bufChan <- 10
	bufChan <- 20
	val1 := <-bufChan
	val2 := <-bufChan

	// 2. Test Unbuffered Channel with Goroutine
	unbufChan := make(chan int, 0)
	go worker(unbufChan, 21)
	val3 := <-unbufChan

	finalValue := val1 + val2 + val3
	expectedValue := 10 + 20 + 42

	if finalValue == expectedValue {
		fmt.Println("[TEST PASSED] Channels and Goroutines worked correctly.")
	} else {
		fmt.Printf("[TEST FAILED] Mismatch in channels logic.\nGot: %v\nExpected: %v\n", finalValue, expectedValue)
	}
}
