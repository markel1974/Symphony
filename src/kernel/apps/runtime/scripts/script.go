//go:build nocompile

package scripts

import "kernel"

func main(a int) {
	//var input string = "INPUT"
	//results := "RESULTS"
	//for y:=10; y==0; y-- {
	//	fmt.Printf("y:%d -> %s\n", y, "Hello World")
	//}
	for y := 0; y < 10; y++ {
		kernel.Printf("y:%d -> %s", y, "Hello World")
	}
	//z := a + b
	//y := a - b
	//x := a * b
	//fmt.Println(input, a, b, results, z, z, y, x)
}
