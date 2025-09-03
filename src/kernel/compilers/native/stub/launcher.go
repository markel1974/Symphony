package stub

import (
	"fmt"
	"log"
	"os"

	"github.com/markel1974/c64emu/src/kernel/compilers"
	"github.com/markel1974/c64emu/src/kernel/vm"
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func VMTest(sequencerId string, baseDir string, prefix string, debug bool) error {
	gk := objects.NewGateKeeper(0)
	op := bytecode.NewOpcodes()
	for _, fileName := range Prepare(baseDir, prefix) {
		fmt.Printf("\n\n------------------ %s ------------------\n", fileName)
		comp, loader, err := compilers.NewCompiler(gk, op, sequencerId)
		if err != nil {
			return fmt.Errorf("compiler error: %s", err)
		}
		var args []interface{} = nil
		//args := []interface{}{1, 2}
		dataFile, _ := os.Open(baseDir + string(os.PathSeparator) + fileName)
		if err = comp.Compile(fileName, dataFile); err != nil {
			return fmt.Errorf("compiler error: %s", err)
		}
		dataFile.Close()
		bc := bytecode.NewBytecode(op, comp.Constants(), comp.Imports(), comp.Globals(), comp.FileSet())
		if debug {
			d := bytecode.NewDisassembler(bc)
			d.Disassemble(log.Writer())
		}
		machine, err := vm.NewVM(gk, op, sequencerId)
		if err != nil {
			return fmt.Errorf("vm initialize error: %s", err)
		}
		entryPoints, err := machine.Setup(loader, bc)
		if err != nil {
			machine.Print(log.Writer())
			return fmt.Errorf("vm setup error: %s", err)
		}
		if err = machine.Run(entryPoints["main"], args...); err != nil {
			if debug {
				machine.Print(log.Writer())
			}
			return fmt.Errorf("vm runtime error: %s", err)
		}
		if debug {
			machine.Print(log.Writer())
		}
		fmt.Println("RETURN VALUE", machine.GetReturnValue(0))
	}
	return nil
}
