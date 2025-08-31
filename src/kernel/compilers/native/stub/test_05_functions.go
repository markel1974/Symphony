package stub

import "fmt"

func counterFactory2(start int) func() int {
	a := start
	b := 1
	c := 2
	d := 3
	e := 4
	return func() int {
		count := 2
		fmt.Println("here", count, a, b, c, d, e)
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
