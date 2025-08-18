package stub

const Source5 = `
package main

import "fmt"

func main(a int, b int) {
	//var input string = "INPUT"
	//results := "RESULTS"
	//for y:=10; y==0; y-- {
	//	fmt.Printf("y:%d -> %s\n", y, "CIAO")
	//}
	for y:=0; y<10; y++ {
		fmt.Printf("y:%d -> %s\n", y, "CIAO")
	}
	//z := a + b
	//y := a - b
	//x := a * b
	//fmt.Println(input, a, b, results, z, z, y, x)
}
`
