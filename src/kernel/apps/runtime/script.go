package runtime

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

const script = `
package main

import "fmt"
import "kernel"

func main(a int) {
	//var input string = "INPUT"
	//results := "RESULTS"
	//for y:=10; y==0; y-- {
	//	fmt.Printf("y:%d -> %s\n", y, "CIAO")
	//}
	for y:=0; y<10; y++ {
		kernel.Printf("y:%d -> %s", y, "CIAO")
	}
	//z := a + b
	//y := a - b
	//x := a * b
	//fmt.Println(input, a, b, results, z, z, y, x)
}
`

// CreateScript initializes and returns a command that triggers garbage collection when executed.
func CreateScript() *process.Command {
	run := func(process interfaces.IUserProcess, args []string) error {
		return nil
	}
	root := process.NewCommand("script", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Script", "Script")
	root.SetScript(script)
	return root
}
