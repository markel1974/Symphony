package stub

import (
	"bytes"
	"fmt"
	"log"

	"github.com/markel1974/symphony/src/compilers"
	"github.com/markel1974/symphony/src/vm/bytecode"
	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/sequencers"
)

// Launcher represents a system that facilitates compilation, setup, and execution of bytecode through a virtual machine.
type Launcher struct {
	gk      objects.IGateKeeper
	vm      *handler.VM
	debug   bool
	entries map[string]uint
}

// NewLauncher creates and returns a new instance of Launcher with a pre-initialized GateKeeper.
func NewLauncher(debug bool) *Launcher {
	l := &Launcher{
		gk:    objects.NewGateKeeper(),
		debug: debug,
	}
	return l
}

// Setup initializes the Launcher by compiling provided source files, creating a VM instance, and configuring entry points.
func (l *Launcher) Setup(seqId string, sources []*FileInfo) error {
	if len(sources) == 0 {
		return fmt.Errorf("no sources found")
	}
	seq, err := sequencers.NewSequencers(seqId)
	if err != nil {
		return err
	}
	loader := bytecode.NewLoader(l.gk)
	compiler, err := compilers.NewCompiler(l.gk, seq, loader)
	if err != nil {
		return err
	}
	var bc *bytecode.Bytecode
	if len(sources) == 1 {
		if err = compiler.Compile(sources[0].Name, sources[0].Data); err != nil {
			return err
		}
		bc = bytecode.NewBytecode(l.gk, compiler.Constants(), compiler.Imports(), compiler.Globals(), compiler.FileSet())
	} else {
		rel := bytecode.NewRelocator(l.gk, loader, seq, bytecode.PreInitFunction, bytecode.InitFunction)
		for _, data := range sources {
			if err = compiler.Compile(data.Name, data.Data); err != nil {
				return err
			}
			bCode := bytecode.NewBytecode(l.gk, compiler.Constants(), compiler.Imports(), compiler.Globals(), compiler.FileSet())
			rel.Add(bCode)
		}
		if bc, err = rel.Relocate(); err != nil {
			return err
		}
	}
	if l.debug {
		d := bytecode.NewDisassembler(bc, seq)
		_ = d.Disassemble(log.Writer())
	}
	l.vm = handler.NewVM(l.gk, seq, seq)
	l.entries, err = l.vm.Setup(loader, bc)
	return err
}

// Exec executes a specified entry point with given arguments and returns the results or an error if execution fails.
// The `entry` parameter specifies the function to execute.
// Debug information is optionally printed when `debug` is true.
func (l *Launcher) Exec(entry string, args []interface{}) ([]interface{}, error) {
	mainId, ok := l.entries[entry]
	if !ok {
		return nil, fmt.Errorf("entry point not found: %s", entry)
	}
	l.vm.EnableRetValues(true)
	rv, err := l.vm.Run(mainId, args...)
	if l.debug {
		l.vm.Print(log.Writer())
	}
	return rv, err
}

// Marshal serializes the provided Bytecode instance into a byte slice and returns it, or an error if encoding fails.
func (l *Launcher) Marshal(bc *bytecode.Bytecode) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	if err := bc.Encode(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal reads and decodes bytecode data from a buffer and returns a Bytecode instance or an error if decoding fails.
func (l *Launcher) Unmarshal(gk objects.IGateKeeper, buf *bytes.Buffer) (*bytecode.Bytecode, error) {
	bc := bytecode.NewBytecodeEmpty(gk)
	if err := bc.Decode(buf); err != nil {
		return nil, err
	}
	return bc, nil
}
