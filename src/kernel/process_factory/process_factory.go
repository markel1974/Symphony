package process_factory

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

type ProcessFactory struct {
	kernel interfaces.IKernel
}

func NewProcessFactory(kernel interfaces.IKernel) *ProcessFactory {
	return &ProcessFactory{
		kernel: kernel,
	}
}

func (fp *ProcessFactory) Create(cmd interfaces.ICommand, line string, options *interfaces.ProcessOptions) interfaces.IProcess {
	task := process.NewProcess(fp.kernel, cmd, line)
	if options != nil {
		task.SetOptions(options)
	}
	return task
}
