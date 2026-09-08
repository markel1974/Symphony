package tests

import (
	"fmt"
)

func worker(c chan int, val int) {
	// fmt.Println("Worker processing value:", val)
	c <- val
}

func main() {
	fmt.Println("--- Running Test: Channels Pattern ---")
	const v1, v2, v3 = 10, 20, 30
	// 1. Test Buffered Channel
	bufChan := make(chan int, 2)
	bufChan <- v1
	bufChan <- v2
	val1 := <-bufChan
	val2 := <-bufChan

	// 2. Test Unbuffered Channel with Goroutine
	unbufChan := make(chan int, 0)
	go worker(unbufChan, v3)
	val3 := <-unbufChan

	finalValue := val1 + val2 + val3
	expectedValue := v1 + v2 + v3

	if finalValue == expectedValue {
		fmt.Println("[TEST PASSED] Channels and Goroutines worked correctly.")
	} else {
		fmt.Printf("[TEST FAILED] Mismatch in channels logic.\nGot: %v\nExpected: %v\n", finalValue, expectedValue)
	}
}
