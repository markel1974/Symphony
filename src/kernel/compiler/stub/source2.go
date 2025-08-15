package stub

const Source2 = `
package main

import "fmt"

func main() {
	x := "home"
	a := 0
	for idx, v := range x {
		a++
		fmt.Println(v, " ", a)
	}
	return 0
}
`
