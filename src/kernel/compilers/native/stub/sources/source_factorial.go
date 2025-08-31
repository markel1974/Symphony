package sources

import "fmt"

func factorial(n int) int {
	if n == 0 {
		return 1
	}
	k := n - 1
	fmt.Println("Step A", k)
	f := factorial(k)
	fmt.Println("Step B", f)
	z := n * f
	return z
}

func main() {
	fact5 := factorial(10)
	fmt.Println(fact5)
}
