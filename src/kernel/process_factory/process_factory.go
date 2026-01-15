package process_factory

import (
	"github.com/markel1974/symphony/src/kernel/interfaces"
	"github.com/markel1974/symphony/src/kernel/process"
)

// ProcessFactory is responsible for creating and managing process instances within the system.
// It encapsulates interactions with the provided IKernel for process-related operations.
type ProcessFactory struct {
}

// NewProcessFactory creates and returns a new ProcessFactory instance using the provided kernel for process management.
func NewProcessFactory() *ProcessFactory {
	return &ProcessFactory{}
}

// Create initializes a new process using the provided command, line, and optional window settings.
func (fp *ProcessFactory) Create(cmd interfaces.ICommand) interfaces.IUserProcess {
	p := process.NewProcess(cmd)
	return p
}
