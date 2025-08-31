package sources

import "fmt"

func counterFactory2(start int) func() int {
	a := start
	b := 1
	c := 2
	d := 3
	e := 4
	return func() int {
		local := "local"
		//fmt.Println("here", local, a, b, e)
		fmt.Println("here", local, start, e, b, c)
		return 1000
	}
}

func main() {
	fmt.Println("--- Running Test 07: Closures and Free Variables ---")
	counterA := counterFactory2(100)
	z := counterA()
	fmt.Println("RESULT", z)
	//z := counterA()
	//fmt.Println(z)
}
