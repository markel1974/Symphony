package stub

import "fmt"

func counterFactory2(start int) func() int {
	//count := start
	return func() int {
		count := 2
		fmt.Println("here", count)
		return count
	}
}

func main() {
	fmt.Println("--- Running Test 07: Closures and Free Variables ---")
	counterA := counterFactory2(100)
	counterA()
	//z := counterA()
	//fmt.Println(z)
}
