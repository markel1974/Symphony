package stub

const Source7 = `
package main

import "fmt"

func test() {
	return 10
}

func main() {
	x := test()
	fmt.Println(x)
}
`
