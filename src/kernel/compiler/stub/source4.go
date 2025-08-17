package stub

const Source4 = `
package main

import "fmt"

func Beta() (int, bool) {
    return 10, false
}

func main() {
	a, b := Beta()
	fmt.Println("PROVA", a, b)
}
`
