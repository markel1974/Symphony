package cpu

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/flag"
	"log"
	"os"
)

//https://dustlayer.com/c64-architecture
//https://www.zimmers.net/cbmpics/cbm/c64/vic-ii.txt

//Notes:
//https://codebase64.org/lib/exe/fetch.php?media=base:safely_freezing_the_c64.pdf
/*
 * Notes:
 * ------
 *
 * Opcode execution:
 *  - All opcodes are resolved into single clock cycles.
 *  - The "state" variable specifies the routine to be executed in the
 *    next cycle. Its upper 8 bits contain the current opcode, its lower
 *    8 bits contain the cycle number (0..7) within the opcode.
 *  - Opcodes are fetched in cycle 0 (state = 0)
 *  - The states 0x0010..0x0027 are used for intr
 *  - There is exactly one memory access in each clock cycle
 *
 *  - The possible interrupt sources are:
 *      IntVic: I flag is checked, jump to ($fffe)
 *      IntCia: I flag is checked, jump to ($fffe)
 *      IntNmi: Jump to ($fffa)
 *      IntRst: Jump to ($fffc)
 *  - The zFlag variable has the inverse meaning of the 6510 Z flag
 *  - Only the highest bit of the nFlag variable is used
 */

type MOS6510 struct {
	*Core
	prefs *config.Config
}

func NewMOS6510() *MOS6510 {
	cpu := &MOS6510{
		prefs: nil,
		Core:  nil,
	}
	return cpu
}

func (cpu *MOS6510) Setup(intr IPic, banks IBanks, prefs *config.Config) {
	cpu.Core = NewCore(intr)
	cpu.banks = banks
	cpu.prefs = prefs
}

func (cpu *MOS6510) Reset() {
	// Read reset vector
	cpu.pc = uint16(cpu.banks.Read(0xfffc)) | (uint16(cpu.banks.Read(0xfffd)) << 8)
	cpu.state = STATE_LAST
	cpu.opFlags = 0
}

// Stack

func (cpu *MOS6510) popFlags(data uint8) {
	cpu.nFlag = data
	cpu.vFlag = data & 0x40
	cpu.dFlag = data & 0x08
	cpu.iFlag = data & 0x04
	cpu.zFlag = flag.BoolToUint8((data & 0x02) == 0)
	cpu.cFlag = data & 0x01
}

func (cpu *MOS6510) buildFlags(bFlags bool) uint8 {
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
			cpu.state = O_BRANCH_BP
		} else {
			cpu.state = O_BRANCH_FP
		}
	} else {
		cpu.state = O_BRANCH_NP
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
	al := (cpu.a & 0x0f) + (data & 0x0f) + k // lower nybble
	if al > 9 {
		al += 6 // BCD fixup
	}
	ah := (cpu.a >> 4) + (data >> 4) // upper nybble
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
	al := (cpu.a & 0x0f) - (data & 0x0f) - k // lower nybble
	ah := (cpu.a >> 4) - (data >> 4)         // upper nybble
	if (al & 0x10) != 0 {
		al -= 6 // BCD fixup
		ah--
	}
	if (ah & 0x10) != 0 {
		ah -= 6 // BCD fixup
	}
	// Set flags
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

func (cpu *MOS6510) GetState() uint8 {
	return cpu.state
}

func (cpu *MOS6510) checkPic() {
	if cpu.state != STATE_LAST {
		return
	}
	if !cpu.pic.HasAny() {
		return
	}
	if cpu.pic.HasReset() {
		cpu.Reset()
		return
	}
	if cpu.pic.HasNMI() {
		// Taken branches to the same page delay the NMI
		delay := 0
		if (cpu.opFlags & OpFlagIntDelayed) != 0 {
			delay = 1
		}
		if (cpu.pic.GetNMICycleDistance(delay)) >= 2 {
			// Simulate an edge-triggered input
			cpu.pic.ClearNMI()
			cpu.state = I_NMI_10
			cpu.opFlags = 0
		}
		return
	}
	if (cpu.pic.HasIRQ()) &&
		((cpu.iFlag == 0) || ((cpu.opFlags & OpFlagIrqDisabled) != 0)) && ((cpu.opFlags & OpFlagIrqEnabled) == 0) {
		delay := 0
		if (cpu.opFlags & OpFlagIntDelayed) != 0 {
			delay = 1
		}
		// Taken branches to the same page delay the IRQ
		if (cpu.pic.GetIrqCycleDistance(delay)) >= 2 {
			cpu.state = I_IRQ_8
			cpu.opFlags = 0
		}
		return
	}
}

func (cpu *MOS6510) SetRDYLow(rdyLow bool) {
	cpu.rdyLow = rdyLow
}

func (cpu *MOS6510) Emulate() {
	cpu.checkPic()

	rdyLow := cpu.rdyLow

	switch cpu.state {
	case STATE_LAST:
		if rdyLow {
			return
		}
		cpu.op = cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.state = _modeTable[cpu.op]
		cpu.opFlags = 0

	// IRQ
	case I_IRQ_8:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.state = I_IRQ_9

	case I_IRQ_9:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.state = I_IRQ_A

	case I_IRQ_A:
		cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
		cpu.sp--
		cpu.state = I_IRQ_B

	case I_IRQ_B:
		cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
		cpu.sp--
		cpu.state = I_IRQ_C

	case I_IRQ_C:
		data := cpu.buildFlags(false)
		cpu.banks.Write((uint16(cpu.sp)&0xff)|0x0100, data)
		cpu.sp--
		cpu.iFlag = 1
		cpu.state = I_IRQ_D

	case I_IRQ_D:
		if rdyLow {
			return
		}
		cpu.pc = uint16(cpu.banks.Read(0xfffe))
		cpu.state = I_IRQ_E

	case I_IRQ_E:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(0xffff)
		cpu.pc |= uint16(data) << 8
		cpu.state = STATE_LAST

	// NMI
	case I_NMI_10:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.state = I_NMI_11

	case I_NMI_11:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.state = I_NMI_12

	case I_NMI_12:
		cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
		cpu.sp--
		cpu.state = I_NMI_13

	case I_NMI_13:
		cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
		cpu.sp--
		cpu.state = I_NMI_14

	case I_NMI_14:
		data := cpu.buildFlags(false)
		cpu.banks.Write((uint16(cpu.sp)&0xff)|0x0100, data)
		cpu.sp--
		cpu.iFlag = 1
		cpu.state = I_NMI_15

	case I_NMI_15:
		if rdyLow {
			return
		}
		cpu.pc = uint16(cpu.banks.Read(0xfffa))
		cpu.state = I_NMI_16

	case I_NMI_16:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(0xfffb)
		cpu.pc |= uint16(data) << 8
		cpu.state = STATE_LAST

	// Addressing modes: Fetch effective address, no extra cycles (-> ar)
	case A_ZERO:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = _opTable[cpu.op]

	case A_ZEROX:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = A_ZEROX1

	case A_ZEROX1:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
		cpu.state = _opTable[cpu.op]

	case A_ZEROY:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = A_ZEROY1

	case A_ZEROY1:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
		cpu.state = _opTable[cpu.op]

	case A_ABS:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = A_ABS1

	case A_ABS1:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.ar = cpu.ar | (uint16(data) << 8)
		cpu.state = _opTable[cpu.op]

	case A_ABSX:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = A_ABSX1

	case A_ABSX1:
		// Note: Some undocumented opcodes rely on the value of ar2
		if rdyLow {
			return
		}
		cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		if cpu.ar+uint16(cpu.x) < 0x100 {
			cpu.state = A_ABSX2
		} else {
			cpu.state = A_ABSX3
		}
		cpu.ar = ((cpu.ar + uint16(cpu.x)) & 0xff) | (cpu.ar2 << 8)

	case A_ABSX2:
		// No page crossed
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.state = _opTable[cpu.op]

	case A_ABSX3:
		// Page crossed
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.ar += 0x100
		cpu.state = _opTable[cpu.op]

	case A_ABSY:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = A_ABSY1

	case A_ABSY1:
		// Note: Some undocumented opcodes rely on the value of ar2
		if rdyLow {
			return
		}
		cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		if cpu.ar+uint16(cpu.y) < 0x100 {
			cpu.state = A_ABSY2
		} else {
			cpu.state = A_ABSY3
		}
		cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (cpu.ar2 << 8)

	case A_ABSY2:
		// No page crossed
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.state = _opTable[cpu.op]

	case A_ABSY3:
		// Page crossed
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.ar += 0x100
		cpu.state = _opTable[cpu.op]

	case A_INDX:
		if rdyLow {
			return
		}
		cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = A_INDX1

	case A_INDX1:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar2)
		cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
		cpu.state = A_INDX2

	case A_INDX2:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
		cpu.state = A_INDX3

	case A_INDX3:
		if rdyLow {
			return
		}
		data := cpu.banks.Read((cpu.ar2 + 1) & 0xff)
		cpu.ar = cpu.ar | (uint16(data) << 8)
		cpu.state = _opTable[cpu.op]

	case A_INDY:
		if rdyLow {
			return
		}
		cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = A_INDY1

	case A_INDY1:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
		cpu.state = A_INDY2

	case A_INDY2:
		// Note: Some undocumented opcodes rely on the value of ar2
		if rdyLow {
			return
		}
		cpu.ar2 = uint16(cpu.banks.Read((cpu.ar2 + 1) & 0xff))
		if cpu.ar+uint16(cpu.y) < 0x100 {
			cpu.state = A_INDY3
		} else {
			cpu.state = A_INDY4
		}
		cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (cpu.ar2 << 8)

	case A_INDY3:
		// No page crossed
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.state = _opTable[cpu.op]

	case A_INDY4:
		// Page crossed
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.ar += 0x100
		cpu.state = _opTable[cpu.op]

		// Addressing modes: Fetch effective address, extra cycle on page crossing (-> ar)
	case AE_ABSX:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = AE_ABSX1

	case AE_ABSX1:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		if cpu.ar+uint16(cpu.x) < 0x100 {
			cpu.ar = ((cpu.ar + uint16(cpu.x)) & 0xff) | (uint16(data) << 8)
			cpu.state = _opTable[cpu.op]
		} else {
			cpu.ar = ((cpu.ar + uint16(cpu.x)) & 0xff) | (uint16(data) << 8)
			cpu.state = AE_ABSX2
		}

	case AE_ABSX2:
		// Page crossed
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.ar += 0x100
		cpu.state = _opTable[cpu.op]

	case AE_ABSY:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = AE_ABSY1

	case AE_ABSY1:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		if cpu.ar+uint16(cpu.y) < 0x100 {
			cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (uint16(data) << 8)
			cpu.state = _opTable[cpu.op]
		} else {
			cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (uint16(data) << 8)
			cpu.state = AE_ABSY2
		}

	case AE_ABSY2:
		// Page crossed
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.ar += 0x100
		cpu.state = _opTable[cpu.op]

	case AE_INDY:
		if rdyLow {
			return
		}
		cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = AE_INDY1

	case AE_INDY1:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
		cpu.state = AE_INDY2

	case AE_INDY2:
		if rdyLow {
			return
		}
		data := cpu.banks.Read((cpu.ar2 + 1) & 0xff)
		if z := cpu.ar + uint16(cpu.y); z < 0x100 {
			cpu.ar = (z & 0xff) | (uint16(data) << 8)
			cpu.state = _opTable[cpu.op]
		} else {
			cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (uint16(data) << 8)
			cpu.state = AE_INDY3
		}

	case AE_INDY3: // Page crossed
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.ar += 0x100
		cpu.state = _opTable[cpu.op]

		// Addressing modes: Read operand, write it back, no extra cycles (-> ar, rmw)
	case M_ZERO:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = RMW_DO_IT

	case M_ZEROX:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = M_ZEROX1

	case M_ZEROX1:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
		cpu.state = RMW_DO_IT

	case M_ZEROY:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = M_ZEROY1

	case M_ZEROY1:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
		cpu.state = RMW_DO_IT

	case M_ABS:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = M_ABS1

	case M_ABS1:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.ar = cpu.ar | (uint16(data) << 8)
		cpu.state = RMW_DO_IT

	case M_ABSX:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = M_ABSX1

	case M_ABSX1:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		if cpu.ar+uint16(cpu.x) < 0x100 {
			cpu.state = M_ABSX2
		} else {
			cpu.state = M_ABSX3
		}
		cpu.ar = (cpu.ar + uint16(cpu.x)&0xff) | (uint16(data) << 8)

	case M_ABSX2:
		// No page crossed
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.state = RMW_DO_IT

	case M_ABSX3:
		// Page crossed
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.ar += 0x100
		cpu.state = RMW_DO_IT

	case M_ABSY:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = M_ABSY1

	case M_ABSY1:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		if cpu.ar+uint16(cpu.y) < 0x100 {
			cpu.state = M_ABSY2
		} else {
			cpu.state = M_ABSY3
		}
		cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (uint16(data) << 8)

	case M_ABSY2:
		// No page crossed
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.state = RMW_DO_IT

	case M_ABSY3:
		// Page crossed
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.ar += 0x100
		cpu.state = RMW_DO_IT

	case M_INDX:
		if rdyLow {
			return
		}
		cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = M_INDX1

	case M_INDX1:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar2)
		cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
		cpu.state = M_INDX2

	case M_INDX2:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
		cpu.state = M_INDX3

	case M_INDX3:
		if rdyLow {
			return
		}
		data := cpu.banks.Read((cpu.ar2 + 1) & 0xff)
		cpu.ar = cpu.ar | (uint16(data) << 8)
		cpu.state = RMW_DO_IT

	case M_INDY:
		if rdyLow {
			return
		}
		cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = M_INDY1

	case M_INDY1:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
		cpu.state = M_INDY2

	case M_INDY2:
		if rdyLow {
			return
		}
		data := cpu.banks.Read((cpu.ar2 + 1) & 0xff)
		if cpu.ar+uint16(cpu.y) < 0x100 {
			cpu.state = M_INDY3
		} else {
			cpu.state = M_INDY4
		}
		cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (uint16(data) << 8)

	case M_INDY3:
		// No page crossed
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.state = RMW_DO_IT

	case M_INDY4:
		// Page crossed
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.ar += 0x100
		cpu.state = RMW_DO_IT

	case RMW_DO_IT:
		if rdyLow {
			return
		}
		cpu.rmw = cpu.banks.Read(cpu.ar)
		cpu.state = RMW_DO_IT1

	case RMW_DO_IT1:
		cpu.banks.Write(cpu.ar, cpu.rmw)
		cpu.state = _opTable[cpu.op]

		// Load group
	case O_LDA:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.ar)
		cpu.a = data
		cpu.nFlag = data
		cpu.zFlag = data
		cpu.state = STATE_LAST

	case O_LDA_I:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.a = data
		cpu.nFlag = data
		cpu.zFlag = data
		cpu.state = STATE_LAST

	case O_LDX:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.ar)
		cpu.x = data
		cpu.nFlag = data
		cpu.zFlag = data
		cpu.state = STATE_LAST

	case O_LDX_I:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.x = data
		cpu.nFlag = data
		cpu.zFlag = data
		cpu.state = STATE_LAST

	case O_LDY:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.ar)
		cpu.y = data
		cpu.nFlag = data
		cpu.zFlag = data
		cpu.state = STATE_LAST

	case O_LDY_I:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.y = data
		cpu.nFlag = data
		cpu.zFlag = data
		cpu.state = STATE_LAST

		// Store group
	case O_STA:
		cpu.banks.Write(cpu.ar, cpu.a)
		cpu.state = STATE_LAST

	case O_STX:
		cpu.banks.Write(cpu.ar, cpu.x)
		cpu.state = STATE_LAST

	case O_STY:
		cpu.banks.Write(cpu.ar, cpu.y)
		cpu.state = STATE_LAST

		// Transfer group
	case O_TAX:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.x = cpu.a
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.state = STATE_LAST

	case O_TXA:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.a = cpu.x
		cpu.nFlag = cpu.x
		cpu.zFlag = cpu.x
		cpu.state = STATE_LAST

	case O_TAY:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.y = cpu.a
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.state = STATE_LAST

	case O_TYA:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.a = cpu.y
		cpu.nFlag = cpu.y
		cpu.zFlag = cpu.y
		cpu.state = STATE_LAST

	case O_TSX:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.x = cpu.sp
		cpu.nFlag = cpu.sp
		cpu.zFlag = cpu.sp
		cpu.state = STATE_LAST

	case O_TXS:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.sp = cpu.x
		cpu.state = STATE_LAST

		// Arithmetic group
	case O_ADC:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.ar)
		cpu.doADC(data)
		cpu.state = STATE_LAST

	case O_ADC_I:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.doADC(data)
		cpu.state = STATE_LAST

	case O_SBC:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.ar)
		cpu.doSBC(data)
		cpu.state = STATE_LAST

	case O_SBC_I:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.doSBC(data)
		cpu.state = STATE_LAST

		// Increment/decrement group
	case O_INX:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.x++
		cpu.nFlag = cpu.x
		cpu.zFlag = cpu.x
		cpu.state = STATE_LAST

	case O_DEX:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.x--
		cpu.nFlag = cpu.x
		cpu.zFlag = cpu.x
		cpu.state = STATE_LAST

	case O_INY:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.y++
		cpu.nFlag = cpu.y
		cpu.zFlag = cpu.y
		cpu.state = STATE_LAST

	case O_DEY:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.y--
		cpu.nFlag = cpu.y
		cpu.zFlag = cpu.y
		cpu.state = STATE_LAST

	case O_INC:
		v := cpu.rmw + 1
		cpu.nFlag = v
		cpu.zFlag = v
		cpu.banks.Write(cpu.ar, v)
		cpu.state = STATE_LAST

	case O_DEC:
		v := cpu.rmw - 1
		cpu.nFlag = v
		cpu.zFlag = v
		cpu.banks.Write(cpu.ar, v)
		cpu.state = STATE_LAST

		// Logic group
	case O_AND:
		if rdyLow {
			return
		}
		cpu.a &= cpu.banks.Read(cpu.ar)
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.state = STATE_LAST

	case O_AND_I:
		if rdyLow {
			return
		}
		cpu.a &= cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.state = STATE_LAST

	case O_ORA:
		if rdyLow {
			return
		}
		cpu.a |= cpu.banks.Read(cpu.ar)
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.state = STATE_LAST

	case O_ORA_I:
		if rdyLow {
			return
		}
		cpu.a |= cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.state = STATE_LAST

	case O_EOR:
		if rdyLow {
			return
		}
		cpu.a ^= cpu.banks.Read(cpu.ar)
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.state = STATE_LAST

	case O_EOR_I:
		if rdyLow {
			return
		}
		cpu.a ^= cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.state = STATE_LAST

		// Compare group
	case O_CMP:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.ar)
		cpu.ar = uint16(cpu.a) - uint16(data)
		cpu.nFlag = uint8(cpu.ar)
		cpu.zFlag = uint8(cpu.ar)
		cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
		cpu.state = STATE_LAST

	case O_CMP_I:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.ar = uint16(cpu.a) - uint16(data)
		cpu.nFlag = uint8(cpu.ar)
		cpu.zFlag = uint8(cpu.ar)
		cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
		cpu.state = STATE_LAST

	case O_CPX:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.ar)
		cpu.ar = uint16(cpu.x) - uint16(data)
		cpu.nFlag = uint8(cpu.ar)
		cpu.zFlag = uint8(cpu.ar)
		cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
		cpu.state = STATE_LAST

	case O_CPX_I:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.ar = uint16(cpu.x) - uint16(data)
		cpu.nFlag = uint8(cpu.ar)
		cpu.zFlag = uint8(cpu.ar)
		cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
		cpu.state = STATE_LAST

	case O_CPY:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.ar)
		cpu.ar = uint16(cpu.y) - uint16(data)
		cpu.nFlag = uint8(cpu.ar)
		cpu.zFlag = uint8(cpu.ar)
		cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
		cpu.state = STATE_LAST

	case O_CPY_I:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.ar = uint16(cpu.y) - uint16(data)
		cpu.nFlag = uint8(cpu.ar)
		cpu.zFlag = uint8(cpu.ar)
		cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
		cpu.state = STATE_LAST

		// Bit-test group
	case O_BIT:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.ar)
		cpu.zFlag = cpu.a & data
		cpu.nFlag = data
		cpu.vFlag = data & 0x40
		cpu.state = STATE_LAST

		// Shift/rotate group
	case O_ASL:
		cpu.cFlag = cpu.rmw & 0x80
		v := cpu.rmw << 1
		cpu.nFlag = v
		cpu.zFlag = v
		cpu.banks.Write(cpu.ar, v)
		cpu.state = STATE_LAST

	case O_ASL_A:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.cFlag = cpu.a & 0x80
		cpu.a <<= 1
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.state = STATE_LAST

	case O_LSR:
		cpu.cFlag = cpu.rmw & 0x01
		v := cpu.rmw >> 1
		cpu.nFlag = v
		cpu.zFlag = v
		cpu.banks.Write(cpu.ar, v)
		cpu.state = STATE_LAST

	case O_LSR_A:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.cFlag = cpu.a & 0x01
		cpu.a >>= 1
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.state = STATE_LAST

	case O_ROL:
		var t uint8
		if cpu.cFlag != 0 {
			t = (cpu.rmw << 1) | 0x01
		} else {
			t = cpu.rmw << 1
		}
		cpu.nFlag = t
		cpu.zFlag = t
		cpu.banks.Write(cpu.ar, t)
		cpu.cFlag = cpu.rmw & 0x80
		cpu.state = STATE_LAST

	case O_ROL_A:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		data := cpu.a & 0x80
		if cpu.cFlag != 0 {
			cpu.a = (cpu.a << 1) | 0x01
		} else {
			cpu.a = cpu.a << 1
		}
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.cFlag = data
		cpu.state = STATE_LAST

	case O_ROR:
		var t uint8
		if cpu.cFlag != 0 {
			t = (cpu.rmw >> 1) | 0x80
		} else {
			t = cpu.rmw >> 1
		}
		cpu.nFlag = t
		cpu.zFlag = t
		cpu.banks.Write(cpu.ar, t)
		cpu.cFlag = cpu.rmw & 0x01
		cpu.state = STATE_LAST

	case O_ROR_A:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		data := cpu.a & 0x01
		if cpu.cFlag != 0 {
			cpu.a = (cpu.a >> 1) | 0x80
		} else {
			cpu.a = cpu.a >> 1
		}
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.cFlag = data
		cpu.state = STATE_LAST

		// Stack group
	case O_PHA:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.state = O_PHA1

	case O_PHA1:
		cpu.banks.Write(uint16(cpu.sp)|0x100, cpu.a)
		cpu.sp--
		cpu.state = STATE_LAST

	case O_PLA:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.state = O_PLA1

	case O_PLA1:
		if rdyLow {
			return
		}
		cpu.banks.Read(uint16(cpu.sp) | 0x100)
		cpu.sp++
		cpu.state = O_PLA2

	case O_PLA2:
		if rdyLow {
			return
		}
		cpu.a = cpu.banks.Read(uint16(cpu.sp) | 0x100)
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.state = STATE_LAST

	case O_PHP:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.state = O_PHP1

	case O_PHP1:
		data := cpu.buildFlags(true)
		cpu.banks.Write((uint16(cpu.sp)&0xff)|0x0100, data)
		cpu.sp--
		cpu.state = STATE_LAST

	case O_PLP:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.state = O_PLP1

	case O_PLP1:
		if rdyLow {
			return
		}
		cpu.banks.Read(uint16(cpu.sp) | 0x100)
		cpu.sp++
		cpu.state = O_PLP2

	case O_PLP2:
		iFlagPrev := cpu.iFlag
		if rdyLow {
			return
		}
		data := cpu.banks.Read(uint16(cpu.sp) | 0x0100)
		cpu.popFlags(data)
		if iFlagPrev == 0 && cpu.iFlag != 0 {
			cpu.opFlags |= OpFlagIrqDisabled
		} else if iFlagPrev != 0 && cpu.iFlag == 0 {
			cpu.opFlags |= OpFlagIrqEnabled
		}
		cpu.state = STATE_LAST

		// Jump/branch group
	case O_JMP:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = O_JMP1

	case O_JMP1:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc = (uint16(data) << 8) | cpu.ar
		cpu.state = STATE_LAST

	case O_JMP_I:
		if rdyLow {
			return
		}
		cpu.pc = uint16(cpu.banks.Read(cpu.ar))
		cpu.state = O_JMP_I1

	case O_JMP_I1:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(((cpu.ar + 1) & 0xff) | (cpu.ar & 0xff00))
		cpu.pc |= uint16(data) << 8
		cpu.state = STATE_LAST

	case O_JSR:
		if rdyLow {
			return
		}
		cpu.ar = uint16(cpu.banks.Read(cpu.pc))
		cpu.pc++
		cpu.state = O_JSR1

	case O_JSR1:
		if rdyLow {
			return
		}
		cpu.banks.Read(uint16(cpu.sp) | 0x100)
		cpu.state = O_JSR2

	case O_JSR2:
		cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
		cpu.sp--
		cpu.state = O_JSR3

	case O_JSR3:
		cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
		cpu.sp--
		cpu.state = O_JSR4

	case O_JSR4:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.pc = cpu.ar | (uint16(data) << 8)
		cpu.state = STATE_LAST

	case O_RTS:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.state = O_RTS1

	case O_RTS1:
		if rdyLow {
			return
		}
		cpu.banks.Read(uint16(cpu.sp) | 0x100)
		cpu.sp++
		cpu.state = O_RTS2

	case O_RTS2:
		if rdyLow {
			return
		}
		cpu.pc = uint16(cpu.banks.Read(uint16(cpu.sp) | 0x100))
		cpu.sp++
		cpu.state = O_RTS3

	case O_RTS3:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(uint16(cpu.sp) | 0x100)
		cpu.pc |= uint16(data) << 8
		cpu.state = O_RTS4

	case O_RTS4:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.state = STATE_LAST

	case O_RTI:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.state = O_RTI1

	case O_RTI1:
		if rdyLow {
			return
		}
		cpu.banks.Read(uint16(cpu.sp) | 0x100)
		cpu.sp++
		cpu.state = O_RTI2

	case O_RTI2:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(uint16(cpu.sp) | 0x0100)
		cpu.popFlags(data)
		cpu.sp++
		cpu.state = O_RTI3

	case O_RTI3:
		if rdyLow {
			return
		}
		cpu.pc = uint16(cpu.banks.Read(uint16(cpu.sp) | 0x100))
		cpu.sp++
		cpu.state = O_RTI4

	case O_RTI4:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(uint16(cpu.sp) | 0x100)
		cpu.pc |= uint16(data) << 8
		cpu.state = STATE_LAST

	case O_BRK:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.state = O_BRK1

	case O_BRK1:
		cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
		cpu.sp--
		cpu.state = O_BRK2

	case O_BRK2:
		cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
		cpu.sp--
		cpu.state = O_BRK3

	case O_BRK3:
		data := cpu.buildFlags(true)
		cpu.banks.Write((uint16(cpu.sp)&0xff)|0x0100, data)
		cpu.sp--
		cpu.iFlag = 1
		// BRK interrupted by NMI?
		if cpu.pic.HasNMI() {
			cpu.pic.ClearNMI()   // Simulate an edge-triggered input
			cpu.state = I_NMI_15 // Jump to NMI sequence
		} else {
			cpu.state = O_BRK4
		}

	case O_BRK4:
		if rdyLow {
			return
		}
		cpu.pc = uint16(cpu.banks.Read(0xfffe))
		cpu.state = O_BRK5

	case O_BRK5:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(0xffff)
		cpu.pc |= uint16(data) << 8
		cpu.state = STATE_LAST

	case O_BCS:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		if cpu.cFlag == 0 {
			cpu.state = STATE_LAST
		} else {
			cpu.branch(data)
		}

	case O_BCC:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		if cpu.cFlag != 0 {
			cpu.state = STATE_LAST
		} else {
			cpu.branch(data)
		}

	case O_BEQ:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		if cpu.zFlag != 0 {
			cpu.state = STATE_LAST
		} else {
			cpu.branch(data)
		}

	case O_BNE:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		if cpu.zFlag == 0 {
			cpu.state = STATE_LAST
		} else {
			cpu.branch(data)
		}

	case O_BVS:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		if cpu.vFlag == 0 {
			cpu.state = STATE_LAST
		} else {
			cpu.branch(data)
		}

	case O_BVC:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		if cpu.vFlag != 0 {
			cpu.state = STATE_LAST
		} else {
			cpu.branch(data)
		}

	case O_BMI:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		if (cpu.nFlag & 0x80) == 0 {
			cpu.state = STATE_LAST
		} else {
			cpu.branch(data)
		}

	case O_BPL:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		if (cpu.nFlag & 0x80) != 0 {
			cpu.state = STATE_LAST
		} else {
			cpu.branch(data)
		}

	case O_BRANCH_NP:
		// No page crossed
		cpu.opFlags |= OpFlagIntDelayed
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.pc = cpu.ar
		cpu.state = STATE_LAST

	case O_BRANCH_BP:
		// Page crossed, branch backwards
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.pc = cpu.ar
		cpu.state = O_BRANCH_BP1

	case O_BRANCH_BP1:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc + 0x100)
		cpu.state = STATE_LAST

	case O_BRANCH_FP:
		// Page crossed, branch forwards
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.pc = cpu.ar
		cpu.state = O_BRANCH_FP1

	case O_BRANCH_FP1:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc - 0x100)
		cpu.state = STATE_LAST

		// Flag group
	case O_SEC:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.cFlag = 1
		cpu.state = STATE_LAST

	case O_CLC:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.cFlag = 0
		cpu.state = STATE_LAST

	case O_SED:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.dFlag = 1
		cpu.state = STATE_LAST

	case O_CLD:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.dFlag = 0
		cpu.state = STATE_LAST

	case O_SEI:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		if cpu.iFlag == 0 {
			cpu.opFlags |= OpFlagIrqDisabled
		}
		cpu.iFlag = 1
		cpu.state = STATE_LAST

	case O_CLI:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		if cpu.iFlag == 0 {
			cpu.opFlags |= OpFlagIrqEnabled
		}
		cpu.iFlag = 0
		cpu.state = STATE_LAST

	case O_CLV:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.vFlag = 0
		cpu.state = STATE_LAST

		// NOP group
	case O_NOP:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.state = STATE_LAST

		// Undocumented opcodes start here

		// NOP group
	case O_NOP_I:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.state = STATE_LAST

	case O_NOP_A:
		if rdyLow {
			return
		}
		cpu.banks.Read(cpu.ar)
		cpu.state = STATE_LAST

		// Load A/X group
	case O_LAX:
		if rdyLow {
			return
		}
		cpu.x = cpu.banks.Read(cpu.ar)
		cpu.a = cpu.x
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.state = STATE_LAST

		// Store A/X group
	case O_SAX:
		cpu.banks.Write(cpu.ar, cpu.a&cpu.x)
		cpu.state = STATE_LAST

		// ASL/ORA group
	case O_SLO:
		cpu.cFlag = cpu.rmw & 0x80
		cpu.rmw <<= 1
		cpu.banks.Write(cpu.ar, cpu.rmw)
		cpu.a |= cpu.rmw
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.state = STATE_LAST

		// ROL/AND group
	case O_RLA:
		tmp := cpu.rmw & 0x80
		if cpu.cFlag != 0 {
			cpu.rmw = (cpu.rmw << 1) | 0x01
		} else {
			cpu.rmw = cpu.rmw << 1
		}
		cpu.cFlag = tmp
		cpu.banks.Write(cpu.ar, cpu.rmw)
		cpu.a &= cpu.rmw
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.state = STATE_LAST

		// LSR/EOR group
	case O_SRE:
		cpu.cFlag = cpu.rmw & 0x01
		cpu.rmw >>= 1
		cpu.banks.Write(cpu.ar, cpu.rmw)
		cpu.a ^= cpu.rmw
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.state = STATE_LAST

		// ROR/ADC group
	case O_RRA:
		tmp := cpu.rmw & 0x01
		if cpu.cFlag != 0 {
			cpu.rmw = (cpu.rmw >> 1) | 0x80
		} else {
			cpu.rmw = cpu.rmw >> 1
		}
		cpu.cFlag = tmp
		cpu.banks.Write(cpu.ar, cpu.rmw)
		cpu.doADC(cpu.rmw)
		cpu.state = STATE_LAST

		// DEC/CMP group
	case O_DCP:
		cpu.rmw--
		cpu.banks.Write(cpu.ar, cpu.rmw)
		cpu.ar = uint16(cpu.a) - uint16(cpu.rmw)
		cpu.nFlag = uint8(cpu.ar)
		cpu.zFlag = uint8(cpu.ar)
		cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
		cpu.state = STATE_LAST

		// INC/SBC group
	case O_ISB:
		cpu.rmw++
		cpu.banks.Write(cpu.ar, cpu.rmw)
		cpu.doSBC(cpu.rmw)
		cpu.state = STATE_LAST

		// Complex functions
	case O_ANC_I:
		if rdyLow {
			return
		}
		cpu.a &= cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.cFlag = cpu.nFlag & 0x80
		cpu.state = STATE_LAST

	case O_ASR_I:
		if rdyLow {
			return
		}
		cpu.a &= cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.cFlag = cpu.a & 0x01
		cpu.a >>= 1
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.state = STATE_LAST

	case O_ARR_I:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		data &= cpu.a
		if cpu.cFlag != 0 {
			cpu.a = (data >> 1) | 0x80
		} else {
			cpu.a = data >> 1
		}
		if cpu.dFlag == 0 {
			cpu.nFlag = cpu.a
			cpu.zFlag = cpu.a
			cpu.cFlag = cpu.a & 0x40
			cpu.vFlag = (cpu.a & 0x40) ^ ((cpu.a & 0x20) << 1)
		} else {
			if cpu.cFlag != 0 {
				cpu.nFlag = 0x80
			} else {
				cpu.nFlag = 0
			}
			cpu.zFlag = cpu.a
			cpu.vFlag = (data ^ cpu.a) & 0x40
			if (data&0x0f)+(data&0x01) > 5 {
				cpu.a = (cpu.a & 0xf0) | ((cpu.a + 6) & 0x0f)
			}
			k := uint16((data)+(uint8(data)&0x10)) & 0x1f0
			cpu.cFlag = uint8(k)
			if k > 0x50 {
				cpu.a += 0x60
			}
		}
		cpu.state = STATE_LAST

	case O_ANE_I:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.a = (cpu.a | 0xee) & cpu.x & data
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.state = STATE_LAST

	case O_LXA_I:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.x = (cpu.a | 0xee) & data
		cpu.a = cpu.x
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.state = STATE_LAST

	case O_SBX_I:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.pc)
		cpu.pc++
		cpu.ar = (uint16(cpu.x) & uint16(cpu.a)) - uint16(data)
		cpu.x = uint8(cpu.ar)
		cpu.nFlag = cpu.x
		cpu.zFlag = cpu.x
		cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
		cpu.state = STATE_LAST

	case O_LAS:
		if rdyLow {
			return
		}
		data := cpu.banks.Read(cpu.ar)
		cpu.sp = data & cpu.sp
		cpu.x = cpu.sp
		cpu.a = cpu.x
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.state = STATE_LAST

	case O_SHS: // ar2 contains the high byte of the operand address
		cpu.sp = cpu.a & cpu.x
		cpu.banks.Write(cpu.ar, uint8((cpu.ar2+1)&uint16(cpu.sp)))
		cpu.state = STATE_LAST

	case O_SHY: // ar2 contains the high byte of the operand address
		cpu.banks.Write(cpu.ar, uint8(uint16(cpu.y)&(cpu.ar2+1)))
		cpu.state = STATE_LAST

	case O_SHX: // ar2 contains the high byte of the operand address
		cpu.banks.Write(cpu.ar, uint8(uint16(cpu.x)&(cpu.ar2+1)))
		cpu.state = STATE_LAST

	case O_SHA: // ar2 contains the high byte of the operand address
		cpu.banks.Write(cpu.ar, uint8(uint16(cpu.a)&uint16(cpu.x)&(cpu.ar2+1)))
		cpu.state = STATE_LAST

	default:
		cpu.illegalOp(cpu.op, cpu.pc-1)
	}
}

func (cpu *MOS6510) printRegisters(qCycle uint64, baLow bool) {
	fmt.Printf("CPU] %d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d\n",
		qCycle,
		cpu.state,
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
