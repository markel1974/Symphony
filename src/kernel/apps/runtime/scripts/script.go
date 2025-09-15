//go:build nocompile

package scripts

import (
	"kernel"
)

func main(in int) {
	for x := 10; x > 0; x-- {
		kernel.Printf("x:%d -> %s", x, "World Hello")
	}
	for y := 0; y <= 10; y++ {
		kernel.Printf("y:%d -> %s", y, "Hello World")
	}
	var input string = "INPUT"
	results := "RESULTS"
	a := 7
	b := 3
	z := a + b
	y := a - b
	x := a * b
	kernel.Printf("%s %d %d %s %d %d %d %d", input, a, b, results, z, z, y, x)
}
