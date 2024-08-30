package mos6510

//https://dustlayer.com/c64-architecture
//https://www.zimmers.net/cbmpics/cbm/c64/vic-ii.txt

//Notes(cpu *MOS6510) {
//https://codebase64.org/lib/exe/fetch.php?media=base:safely_freezing_the_c64.pdf
/*
 *  - The zFlag variable has the inverse meaning of the 6510 Z flag
 *  - Only the highest bit of the nFlag variable is used
 */

type MOS6510 struct {
	core *Core
	id   string
}

func NewMOS6510(id string) *MOS6510 {
	cpu := &MOS6510{
		core: nil,
		id:   id,
	}
	return cpu
}

func (cpu *MOS6510) Setup(intr IPic, banks IBanks) {
	cpu.core = NewCore(intr, banks)
}

func (cpu *MOS6510) Reset() {
	cpu.core.reset()
}

// SetOverflowBranch implement 6502c SO (SOB) Pin
func (cpu *MOS6510) SetOverflowBranch(sob func() bool) {
	cpu.core.overflowBranch = sob
}

func (cpu *MOS6510) SetAECLow(aecLow bool) {
	cpu.core.aecLow = aecLow
	if cpu.core.aecLow {
		cpu.core.stop = true
	}
}

func (cpu *MOS6510) SetRDYLow(rdyLow bool) {
	cpu.core.rdyLow = rdyLow
	if !cpu.core.rdyLow {
		cpu.core.stop = false
	}
}

func (cpu *MOS6510) Emulate() {
	if cpu.core.stop {
		return
	}
	cpu.core.next(cpu.core)
}
