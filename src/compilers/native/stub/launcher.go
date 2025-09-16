package stub

import (
	"embed"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/markel1974/c64emu/src/compilers"
	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/sequencers/native"
)

// All is an embed.FS instance that provides access to embedded files from specified directories.
//
//go:embed sources/*.go
//go:embed tests/*.go
var All embed.FS

// walk traverses the given directory recursively, appending file paths with matching prefixes to the provided slice.
// It returns an error if the directory cannot be read.
func walk(baseDir string, prefix string, files *[]string) error {
	data, err := All.ReadDir(baseDir)
	if err != nil {
		return err
	}
	for _, v := range data {
		target := filepath.Join(baseDir, v.Name())
		if v.IsDir() {
			if err = walk(target, prefix, files); err != nil {
				return err
			}
		}
		if !strings.HasPrefix(v.Name(), prefix) {
			continue
		}
		*(files) = append(*(files), target)
	}
	return nil
}

// Launch processes files with a specified prefix, compiles them, and executes their bytecode in a virtual machine.
// It uses a created GateKeeper and Sequencer for managing execution contexts and handles debug outputs if enabled.
// Returns an error if file walking, compilation, setup, or execution fails.
func Launch(prefix string, debug bool) error {
	var files []string
	if err := walk(".", prefix, &files); err != nil {
		return nil
	}
	sort.Strings(files)

	gk := objects.NewGateKeeper()
	for _, fileName := range files {
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
		dataFile, err := All.ReadFile(fileName)
		if err != nil {
			return fmt.Errorf("compiler error: %s", err)
		}
		if err = comp.Compile(fileName, dataFile); err != nil {
			return fmt.Errorf("compiler error: %s", err)
		}
		bc := bytecode.NewBytecode(comp.Constants(), comp.Imports(), comp.Globals(), comp.FileSet())
		if debug {
			d := bytecode.NewDisassembler(bc, seq)
			_ = d.Disassemble(log.Writer())
		}
		//buf := bytes.NewBuffer([]byte{})
		//if err = bc.Encode(buf); err != nil {
		//	return fmt.Errorf("compiler error: %s", err)
		//}
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
