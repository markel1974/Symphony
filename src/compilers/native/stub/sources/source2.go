package sources

import "fmt"

func test() (bool, bool) {
	fmt.Println("home", len("home"))
	return true, false
}

func main() {
	y, z := test()
	x := "home"
	//k, ok := y.(bool)
	a, b := 0, 1000
	for idx, v := range x {
		a++
		fmt.Println(string(v), " ", a, b, y, z, idx)
	}
}
