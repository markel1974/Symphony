package mos6510_rev1

// stackAddr defines the base address for the stack memory operations used in various CPU Instructions.
const (
	stackAddr = 0x100
)

// InstOpINI handles the initial opcode fetch and subsequent CPU instruction cycle logic based on the current CPU state.
// It considers interrupt conditions, updates the program counter, and sets the next instruction handler.
// If the RDY line is low, the CPU execution is halted by setting the `stop` flag.
//
//go:nosplit
func InstOpINI(cpu *CPU) {
	// https://www.zimmers.net/cbmpics/cbm/c64/vic-ii.txt
	// Interrupts are only recognized if the RDY line is high
	if !cpu.rdyLow {
		if !cpu.irqBreaker {
			opFlag := cpu.opFlags
			cpu.opFlags = 0
			switch cpu.pic.VerifyIrq(cpu.iFlag, opFlag) {
			case 1:
				cpu.Reset()
				return
			case 2:
				cpu.irqBreaker = true
				cpu.next = InstOpNMI
				cpu.next(cpu)
				return
			case 3:
				cpu.irqBreaker = true
				cpu.next = InstOpIRQ
				cpu.next(cpu)
				return
			}
		}
	} else {
		cpu.setModeHalt()
		return
	}
	cpu.irqBreaker = false
	cpu.op = cpu.bankRead(cpu.pc)
	cpu.pc++
	cpu.next = cpu.modeTable[cpu.op]
}
