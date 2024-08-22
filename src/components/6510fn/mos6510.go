package mos6510fn

import (
	"fmt"
	"github.com/markel1974/c64emu/src/flag"
	"log"
	"os"
)

//https://dustlayer.com/c64-architecture
//https://www.zimmers.net/cbmpics/cbm/c64/vic-ii.txt

//Notes(cpu *MOS6510) {
//https://codebase64.org/lib/exe/fetch.php?media=base:safely_freezing_the_c64.pdf
/*
 *  - The zFlag variable has the inverse meaning of the 6510 Z flag
 *  - Only the highest bit of the nFlag variable is used
 */

type MOS6510 struct {
	*Core
	id             string
	overflowBranch func() bool
}

func NewMOS6510fn(id string) *MOS6510 {
	cpu := &MOS6510{
		Core:           nil,
		id:             id,
		overflowBranch: nil,
	}
	return cpu
}

func (cpu *MOS6510) Setup(intr IPic, banks IBanks) {
	cpu.Core = NewCore(intr)
	cpu.banks = banks
}

func (cpu *MOS6510) Reset() {
	doReset(cpu)
}

// SetOverflowBranch implement 6502c SO (SOB) Pin
func (cpu *MOS6510) SetOverflowBranch(sob func() bool) {
	cpu.overflowBranch = sob
}

func (cpu *MOS6510) popFlags(data uint8) {
	cpu.nFlag = data
	cpu.vFlag = data & 0x40
	cpu.dFlag = data & 0x08
	cpu.iFlag = data & 0x04
	cpu.zFlag = flag.BoolToUint8((data & 0x02) == 0)
	cpu.cFlag = data & 0x01
}

func (cpu *MOS6510) pushFlags(bFlags bool) uint8 {
	data := 0x20 | (cpu.nFlag & 0x80)
	if cpu.vFlag != 0 {
		data |= 0x40
	}
	if bFlags {
		data |= 0x10
	}
	if cpu.dFlag != 0 {
		data |= 0x08
	}
	if cpu.iFlag != 0 {
		data |= 0x04
	}
	if cpu.zFlag == 0 {
		data |= 0x02
	}
	if cpu.cFlag != 0 {
		data |= 0x01
	}
	return data
}

func (cpu *MOS6510) branch(data uint8) {
	cpu.ar = cpu.pc + uint16(int8(data))
	if (cpu.ar >> 8) != (cpu.pc >> 8) {
		if data&0x80 != 0 {
			cpu.next = fnO_BRANCH_BP
		} else {
			cpu.next = fnO_BRANCH_FP
		}
	} else {
		cpu.next = fnO_BRANCH_NP
	}
}

func (cpu *MOS6510) doADC(data uint8) {
	k := uint8(0)
	if cpu.cFlag != 0 {
		k = 1
	}
	if cpu.dFlag == 0 {
		// Binary mode
		tmp := uint16(cpu.a) + uint16(data) + uint16(k)
		cpu.cFlag = flag.BoolToUint8(tmp > 0xff)
		p1 := (uint16(cpu.a) ^ uint16(data)) & 0x80
		p2 := (uint16(cpu.a) ^ tmp) & 0x80
		cpu.vFlag = flag.BoolToUint8((p1 == 0) && (p2 != 0))
		cpu.a = uint8(tmp)
		cpu.nFlag = uint8(tmp)
		cpu.zFlag = uint8(tmp)
		return
	}
	// Decimal mode
	al := (cpu.a & 0x0f) + (data & 0x0f) + k
	if al > 9 {
		al += 6 // BCD fixup
	}
	ah := (cpu.a >> 4) + (data >> 4)
	if al > 0x0f {
		ah++
	}
	cpu.zFlag = cpu.a + data + k
	cpu.nFlag = ah << 4 // Only highest bit used
	p1 := ((ah << 4) ^ cpu.a) & 0x80
	p2 := (cpu.a ^ data) & 0x80
	cpu.vFlag = flag.BoolToUint8((p1 != 0) && (p2 == 0))
	if ah > 9 {
		ah += 6
	}
	// BCD fixup for upper nybble
	cpu.cFlag = flag.BoolToUint8(ah > 0x0f) // carry flag
	cpu.a = (ah << 4) | (al & 0x0f)         // result
}

func (cpu *MOS6510) doSBC(data uint8) {
	k := uint8(0)
	if cpu.cFlag == 0 {
		k = 1
	}
	tmp := uint16(cpu.a) - uint16(data) - uint16(k)
	if cpu.dFlag == 0 {
		// Binary mode
		cpu.cFlag = flag.BoolToUint8(tmp < 0x100)
		p1 := (uint16(cpu.a) ^ tmp) & 0x80
		p2 := (uint16(cpu.a) ^ uint16(data)) & 0x80
		cpu.vFlag = flag.BoolToUint8((p1 != 0) && (p2 != 0))
		cpu.a = uint8(tmp)
		cpu.nFlag = uint8(tmp)
		cpu.zFlag = uint8(tmp)
		return
	}
	// Decimal mode
	al := (cpu.a & 0x0f) - (data & 0x0f) - k
	ah := (cpu.a >> 4) - (data >> 4)
	if (al & 0x10) != 0 {
		al -= 6 // BCD fixup
		ah--
	}
	if (ah & 0x10) != 0 {
		ah -= 6 // BCD fixup
	}
	cpu.cFlag = flag.BoolToUint8(uint16(tmp) < 0x100)
	p1 := (uint16(cpu.a) ^ tmp) & 0x80
	p2 := (uint16(cpu.a) ^ uint16(data)) & 0x80
	cpu.vFlag = flag.BoolToUint8((p1 != 0) && (p2 != 0))
	cpu.zFlag = uint8(tmp)
	cpu.nFlag = uint8(tmp)
	cpu.a = (ah << 4) | (al & 0x0f)
}

func (cpu *MOS6510) illegalOp(illOp uint8, at uint16) {
	log.Printf("illegal opcode %02x at %04x.", illOp, at)
	//TODO EVENT
	cpu.Reset()
	os.Exit(1)
}

func (cpu *MOS6510) SetAECLow(aecLow bool) {
	cpu.aecLow = aecLow
	if cpu.aecLow {
		cpu.stop = true
	}
}

func (cpu *MOS6510) SetRDYLow(rdyLow bool) {
	cpu.rdyLow = rdyLow
	if !cpu.rdyLow {
		cpu.stop = false
	}
}

func (cpu *MOS6510) Emulate() {
	if cpu.stop {
		return
	}
	cpu.next(cpu)
}

func (cpu *MOS6510) printRegisters(qCycle uint64, baLow bool) {
	fmt.Printf("CPU] %d|%d||%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d\n",
		qCycle,
		//cpu.state,
		flag.BoolToUint8(baLow),
		cpu.nFlag,
		cpu.zFlag,
		cpu.vFlag,
		cpu.dFlag,
		cpu.iFlag,
		cpu.cFlag,
		cpu.a,
		cpu.x,
		cpu.y,
		cpu.sp,
		cpu.pc,
		cpu.op,
		cpu.ar,
		cpu.ar2,
		cpu.rmw)
	//cpu.extConfig)
}
