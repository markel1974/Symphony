package stub

import (
	"fmt"
	"log"
	"os"

	"github.com/markel1974/c64emu/src/compilers"
	"github.com/markel1974/c64emu/src/vm"
	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
)

func VMTest(sequencerId string, baseDir string, prefix string, debug bool) error {
	gk := objects.NewGateKeeper()
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
		bc := bytecode.NewBytecode(comp.Constants(), comp.Imports(), comp.Globals(), comp.FileSet())
		if debug {
			d := bytecode.NewDisassembler(bc, op)
			_ = d.Disassemble(log.Writer())
		}

		//rel := bytecode.NewRelocator(gk, loader, op, nil)
		//_, _ = rel.Relocate([]*bytecode.Bytecode{bc, bc})

		machine, err := vm.NewVM(gk, op, sequencerId)
		if err != nil {
			return fmt.Errorf("vm initialize error: %s", err)
		}
		entryPoints, err := machine.Setup(loader, bc)
		if err != nil {
			machine.Print(log.Writer())
			return fmt.Errorf("vm setup error: %s", err)
		}
		machine.EnableRetValues(true)
		rv, err := machine.Run(entryPoints["main"], args...)
		if err != nil {
			if debug {
				machine.Print(log.Writer())
			}
			return fmt.Errorf("vm runtime error: %s", err)
		}
		if debug {
			machine.Print(log.Writer())
		}
		//machine.GetReturnValue(0)
		fmt.Println("RETURN VALUEs", rv)
	}
	return nil
}
