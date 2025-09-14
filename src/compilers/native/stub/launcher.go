package stub

import (
	"fmt"
	"log"
	"os"

	"github.com/markel1974/c64emu/src/compilers"
	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/sequencers/native"
)

func VMTest(sequencerId string, baseDir string, prefix string, debug bool) error {
	gk := objects.NewGateKeeper()
	for _, fileName := range Prepare(baseDir, prefix) {
		fmt.Printf("\n\n------------------ %s ------------------\n", fileName)
		seq := native.NewSequencer()
		if err := seq.Setup(); err != nil {
			return err
		}
		comp, loader, err := compilers.NewCompiler(gk, seq)
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
			d := bytecode.NewDisassembler(bc, seq)
			_ = d.Disassemble(log.Writer())
		}
		machine := core.New(gk, seq)
		if err = seq.Bind(machine); err != nil {
			return err
		}
		//rel := bytecode.NewRelocator(gk, loader, op, nil)
		//_, _ = rel.Relocate([]*bytecode.Bytecode{bc, bc})

		//machine, err := vm.NewVM(gk, seq)
		//if err != nil {
		//	return fmt.Errorf("vm initialize error: %s", err)
		//}

		entryPoints, err := machine.Setup(loader, seq.Executors(), bc)
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
