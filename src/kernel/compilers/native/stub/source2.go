package stub

import "fmt"

func test() (bool, bool) {
	fmt.Println("home")
	return true, false
}

func main() {
	y, z := test()
	x := "home"
	a, b := 0, 1000
	for idx, v := range x {
		a++
		fmt.Println(string(v), " ", a, b, y, z, idx)
	}
}
