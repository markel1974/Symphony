package stub

const Source7 = `
package main

import "fmt"

func test() {
	return len("prova")
}

func main() {
	x := test()
	fmt.Println(x)
}
`
