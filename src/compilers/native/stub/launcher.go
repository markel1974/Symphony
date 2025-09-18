package stub

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"

	_nativeCompiler "github.com/markel1974/c64emu/src/compilers/native/compiler"
	_nativeLoader "github.com/markel1974/c64emu/src/compilers/native/sdk"
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
		loader, err := _nativeLoader.NewLoader(gk)
		if err != nil {
			return err
		}
		bc, err := Compile(gk, seq, loader, fileName, debug)
		if err != nil {
			return err
		}
		if err = Exec(gk, seq, loader, bc, debug); err != nil {
			return err
		}
	}
	return nil
}

func Compile(gk objects.IGateKeeper, seq core.ISequencer, loader bytecode.ILoader, fileName string, debug bool) (*bytecode.Bytecode, error) {
	comp := _nativeCompiler.New(gk, loader, seq)
	dataFile, err := All.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read file error: %s", err)
	}
	if err = comp.Compile(fileName, dataFile); err != nil {
		return nil, fmt.Errorf("compiler error: %s", err)
	}
	bc := bytecode.NewBytecode(gk, comp.Constants(), comp.Imports(), comp.Globals(), comp.FileSet())
	if debug {
		d := bytecode.NewDisassembler(bc, seq)
		_ = d.Disassemble(log.Writer())
	}
	return bc, nil
}

func Exec(gk objects.IGateKeeper, seq core.ISequencer, loader bytecode.ILoader, bc *bytecode.Bytecode, debug bool) error {
	var args []interface{} = nil
	//args := []interface{}{1, 2}
	machine := core.New(gk, seq)
	if err := seq.Bind(machine); err != nil {
		return err
	}
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
	fmt.Println("RETURN VALUEs", rv)
	return nil
}

func Marshal(bc *bytecode.Bytecode) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	if err := bc.Encode(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Unmarshal(gk objects.IGateKeeper, buf *bytes.Buffer) (*bytecode.Bytecode, error) {
	bc := bytecode.NewBytecodeEmpty(gk)
	if err := bc.Decode(buf); err != nil {
		return nil, err
	}
	return bc, nil
}

func Relocate(gk objects.IGateKeeper, seq core.ISequencer, loader bytecode.ILoader, z []*bytecode.Bytecode) (*bytecode.Bytecode, error) {
	rel := bytecode.NewRelocator(gk, loader, seq, bytecode.PreInitFunction, bytecode.InitFunction)
	bc, err := rel.Relocate(z)
	return bc, err
}
