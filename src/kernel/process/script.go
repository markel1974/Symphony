package process

import (
	"fmt"

	"github.com/markel1974/c64emu/src/compilers"
	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/sequencers/native"
)

// Script represents a virtual machine capable of managing and executing bytecode, with components for compilation and execution.
type Script struct {
	gk          objects.IGateKeeper
	loader      bytecode.ILoader
	seq         handler.ISequencer
	compiler    bytecode.ICompiler
	vm          *handler.VM
	entryPoints map[string]uint
	initialized bool
}

// NewScript creates and returns a new instance of the Scripts structure with its fields initialized to their default values.
func NewScript() *Script {
	return &Script{
		initialized: false,
		entryPoints: make(map[string]uint),
	}
}

func (v *Script) IsInitialized() bool {
	return v.initialized
}

// Setup initializes the Scripts with necessary components and binds the sequencer. Returns an error if setup fails.
func (v *Script) Setup(p *Process) error {
	if v.initialized {
		return nil
	}
	var err error
	v.seq = native.NewSequencer()
	if err = v.seq.Setup(); err != nil {
		return err
	}
	v.gk = objects.NewGateKeeper()
	v.loader = bytecode.NewLoader(v.gk)
	v.compiler, err = compilers.NewCompiler(v.gk, v.seq, v.loader)
	if err != nil {
		return err
	}
	pkg := NewLibrary(v.gk, p)
	if err = v.loader.AddPackage(pkg); err != nil {
		return err
	}
	v.vm = handler.NewVM(v.gk, v.seq, v.seq)
	//if err = v.seq.Bind(v.vm); err != nil {
	//	return err
	//}
	v.initialized = true
	return nil
}

// Compile compiles the provided code with the specified name into bytecode and sets up entry points for execution.
func (v *Script) Compile(name string, code string) error {
	if !v.initialized {
		return fmt.Errorf("vm not initialized")
	}
	if err := v.compiler.Compile(name, []byte(code)); err != nil {
		return err
	}
	bc := bytecode.NewBytecode(v.gk, v.compiler.Constants(), v.compiler.Imports(), v.compiler.Globals(), v.compiler.FileSet())
	entryPoints, err := v.vm.Setup(v.loader, bc)
	if err != nil {
		return err
	}
	if entryPoints != nil {
		v.entryPoints = entryPoints
	}
	return nil
}

// Exec runs the specified entry point in the virtual machine with the provided arguments and returns the result or an error.
func (v *Script) Exec(entryPoint string, hasRet bool, args ...interface{}) ([]interface{}, error) {
	if !v.initialized {
		return nil, fmt.Errorf("vm not initialized")
	}
	idx, ok := v.entryPoints[entryPoint]
	if !ok {
		return nil, fmt.Errorf("entry point '%s' not found", entryPoint)
	}
	v.vm.EnableRetValues(hasRet)
	ret, err := v.vm.Run(idx, args...)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
