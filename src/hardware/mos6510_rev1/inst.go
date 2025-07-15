package mos6510_rev1

// stackAddr defines the base address for the stack memory operations used in various CPU Instructions.
const (
	stackAddr = 0x100
)

// InstOpHalt pauses the CPU's operation by acting as a no-op function while the CPU remains in the halted state.
//
//go:nosplit
func InstOpHalt(_ *CPU) {
}

// InstOpINI handles the initial opcode fetch and subsequent CPU instruction cycle logic based on the current CPU state.
// It considers interrupt conditions, updates the program counter, and sets the next instruction handler.
// If the RDY line is low, the CPU execution is halted by setting the `stop` flag.
//
//go:nosplit
func InstOpINI(cpu *CPU) {
	// https://www.zimmers.net/cbmpics/cbm/c64/vic-ii.txt
	// Interrupts are only recognized if the RDY line is high
	if !cpu.rdyLow {
		if !cpu.interrupts.HasIrqBreaker() {
			opFlag := cpu.opFlags
			cpu.opFlags = 0
			switch cpu.interrupts.VerifyIrq(cpu.iFlag, opFlag) {
			case 1:
				cpu.Reset()
				return
			case 2:
				cpu.interrupts.EnableIrqBreaker()
				cpu.next = InstOpNMI
				cpu.next(cpu)
				return
			case 3:
				cpu.interrupts.EnableIrqBreaker()
				cpu.next = InstOpIRQ
				cpu.next(cpu)
				return
			}
		}
	} else {
		cpu.setModeHalt()
		return
	}
	cpu.interrupts.DisableIrqBreaker()
	v, ok := cpu.busRead(cpu.pc)
	if !ok {
		return
	}
	cpu.op = v
	cpu.pc++
	cpu.next = cpu.modeTable[cpu.op]
}
