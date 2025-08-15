package stub

const Source1 = `
package main

import "fmt"

func test() int {
	fmt.println("Hello, world!")
	return 10
}

func main() {
	a := test()
	for x:=10; x==0; x-- {
		fmt.println(x)
	}
	fmt.println("PROVA")
	var x = 4
	var y = 15
	z := x+y
	b := len("Hello, world!,Hello, world!, Hello,world!,Hello, world!,Hello, world!")
	fmt.println(b)
	fmt.println(z)
	return z
}
`
