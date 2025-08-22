package stub

const Source5 = `
package main

import "fmt"

func test() int {
	return 10
}

func main(a int, b int) {
	//var input string = "INPUT"
	//results := "RESULTS"
	//for y:=10; y==0; y-- {
	//	fmt.Printf("y:%d -> %s\n", y, "CIAO")
	//}
	y:=len("CIAO")
	k := test()
	fmt.Printf("y:%d -> %s %d\n", y, "CIAO", k)
	for y:=0; y<10; y++ {
		fmt.Printf("y:%d -> %s\n", y, "CIAO")
	}
	//z := a + b
	//y := a - b
	//x := a * b
	//fmt.Println(input, a, b, results, z, z, y, x)
	return 10
}
`
