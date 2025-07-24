package mos6510_rev1

// stackAddr defines the base address for the stack memory operations used in various CPU Instructions.
const (
	stackAddr = 0x100
)

// InstOpHalt pauses the CPU's operation by acting as a no-op function while the CPU remains in the halted state.
//
//go:nosplit
func (er *Executor) InstOpHalt(_ *CPU) {
}

// InstOpINI handles the initial opcode fetch and subsequent CPU instruction cycle logic based on the current CPU state.
// It considers interrupt conditions, updates the program counter, and sets the next instruction handler.
// If the RDY line is low, the CPU execution is halted by setting the `stop` flag.
//
//go:nosplit
func (er *Executor) InstOpINI(cpu *CPU) {
	// https://www.zimmers.net/cbmpics/cbm/c64/vic-ii.txt
	// Interrupts are only recognized if the RDY line is high
	if cpu.rdyLow {
		cpu.setModeHalt()
		return
	}
	opFlags := cpu.opFlags
	cpu.opFlags = 0
	if cpu.interrupts.Compute(cpu.iFlag, opFlags) {
		return
	}
	v, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.op = v
	cpu.pc++
	cpu.next = cpu.modeTable[cpu.op]
}
