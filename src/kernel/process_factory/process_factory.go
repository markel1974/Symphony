package process_factory

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

// ProcessFactory is responsible for creating and managing process instances within the system.
// It encapsulates interactions with the provided IKernel for process-related operations.
type ProcessFactory struct {
	kernel interfaces.IKernel
}

// NewProcessFactory creates and returns a new ProcessFactory instance using the provided kernel for process management.
func NewProcessFactory(kernel interfaces.IKernel) *ProcessFactory {
	return &ProcessFactory{
		kernel: kernel,
	}
}

// Create initializes a new process using the provided command, line, and optional window settings.
func (fp *ProcessFactory) Create(parent interfaces.IProcess, user string, cmd interfaces.ICommand, line string) interfaces.IProcess {
	p := process.NewProcess(fp.kernel, parent, user, cmd, line)
	return p
}
