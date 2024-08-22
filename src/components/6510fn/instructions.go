package mos6510fn

import "github.com/markel1974/c64emu/src/flag"

func fnInit(cpu *Core) {
	if cpu.pic.HasAny() {
		if cpu.pic.HasReset() {
			cpu.reset()
		} else if cpu.pic.HasNMI() {
			delay := 0
			if (cpu.opFlags & OpFlagIntDelayed) != 0 {
				delay = 1
			}
			if (cpu.pic.GetNMICycleDistance(delay)) >= 2 {
				// Edge-triggered
				cpu.pic.ClearNMI()
				cpu.opFlags = 0
				cpu.next = fnNMI
				cpu.next(cpu)
				return
			}
		} else if cpu.pic.HasIRQ() && !cpu.rdyLow {
			// Interrupts are recognized if the RDY line is high
			if ((cpu.iFlag == 0) || ((cpu.opFlags & OpFlagIrqDisabled) != 0)) && ((cpu.opFlags & OpFlagIrqEnabled) == 0) {
				delay := 0
				if (cpu.opFlags & OpFlagIntDelayed) != 0 {
					delay = 1
				}
				if (cpu.pic.GetIrqCycleDistance(delay)) >= 2 {
					cpu.opFlags = 0
					cpu.next = fnIRQ
					cpu.next(cpu)
					return
				}
			}
		}
	}
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.op = cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.next = _modeTable[cpu.op]
	cpu.opFlags = 0
}

// IRQ
func fnIRQ(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = fnIRQ1
}

func fnIRQ1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = fnIRQ2
}

func fnIRQ2(cpu *Core) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = fnIRQ3
}

func fnIRQ3(cpu *Core) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
	cpu.sp--
	cpu.next = fnIRQ4
}

func fnIRQ4(cpu *Core) {
	data := cpu.pushFlags(false)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|0x0100, data)
	cpu.sp--
	cpu.iFlag = 1
	cpu.next = fnIRQ5
}

func fnIRQ5(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(0xfffe))
	cpu.next = fnIRQ6
}

func fnIRQ6(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(0xffff)
	cpu.pc |= uint16(data) << 8
	cpu.next = fnInit
}

// NMI
func fnNMI(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = fnNMI1
}

func fnNMI1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = fnNMI2
}

func fnNMI2(cpu *Core) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = fnNMI3
}

func fnNMI3(cpu *Core) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
	cpu.sp--
	cpu.next = fnNMI4
}

func fnNMI4(cpu *Core) {
	data := cpu.pushFlags(false)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|0x0100, data)
	cpu.sp--
	cpu.iFlag = 1
	cpu.next = fnNMI5
}

func fnNMI5(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(0xfffa))
	cpu.next = fnNMI6
}

func fnNMI6(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(0xfffb)
	cpu.pc |= uint16(data) << 8
	cpu.next = fnInit
}

// Addressing modes: Fetch effective address, no extra cycles (-> ar)
func fnA_ZERO(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = _opTable[cpu.op]
}

func fnA_ZEROX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnA_ZEROX1
}

func fnA_ZEROX1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
	cpu.next = _opTable[cpu.op]
}

func fnA_ZEROY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnA_ZEROY1
}

func fnA_ZEROY1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
	cpu.next = _opTable[cpu.op]
}

func fnA_ABS(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnA_ABS1
}

func fnA_ABS1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = _opTable[cpu.op]
}

func fnA_ABSX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnA_ABSX1
}

func fnA_ABSX1(cpu *Core) {
	// Note: Some undocumented functions rely on the value of ar2
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	if cpu.ar+uint16(cpu.x) < 0x100 {
		cpu.next = fnA_ABSX2
	} else {
		cpu.next = fnA_ABSX3
	}
	cpu.ar = ((cpu.ar + uint16(cpu.x)) & 0xff) | (cpu.ar2 << 8)
}

func fnA_ABSX2(cpu *Core) {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.next = _opTable[cpu.op]
}

func fnA_ABSX3(cpu *Core) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = _opTable[cpu.op]
}

func fnA_ABSY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnA_ABSY1
}

func fnA_ABSY1(cpu *Core) {
	// Note: Some undocumented functions rely on the value of ar2
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	if cpu.ar+uint16(cpu.y) < 0x100 {
		cpu.next = fnA_ABSY2
	} else {
		cpu.next = fnA_ABSY3
	}
	cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (cpu.ar2 << 8)
}

func fnA_ABSY2(cpu *Core) {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.next = _opTable[cpu.op]
}

func fnA_ABSY3(cpu *Core) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = _opTable[cpu.op]
}

func fnA_INDX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnA_INDX1
}

func fnA_INDX1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar2)
	cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
	cpu.next = fnA_INDX2
}

func fnA_INDX2(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.next = fnA_INDX3
}

func fnA_INDX3(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read((cpu.ar2 + 1) & 0xff)
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = _opTable[cpu.op]
}

func fnA_INDY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnA_INDY1
}

func fnA_INDY1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.next = fnA_INDY2
}

func fnA_INDY2(cpu *Core) {
	// Note: Some undocumented functions rely on the value of ar2
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read((cpu.ar2 + 1) & 0xff))
	if cpu.ar+uint16(cpu.y) < 0x100 {
		cpu.next = fnA_INDY3
	} else {
		cpu.next = fnA_INDY4
	}
	cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (cpu.ar2 << 8)
}

func fnA_INDY3(cpu *Core) {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.next = _opTable[cpu.op]
}

func fnA_INDY4(cpu *Core) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = _opTable[cpu.op]
}

func fnAE_ABSX(cpu *Core) {
	// Addressing modes: Fetch effective address, extra cycle on page crossing (-> ar)
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnAE_ABSX1
}

func fnAE_ABSX1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.ar+uint16(cpu.x) < 0x100 {
		cpu.ar = ((cpu.ar + uint16(cpu.x)) & 0xff) | (uint16(data) << 8)
		cpu.next = _opTable[cpu.op]
	} else {
		cpu.ar = ((cpu.ar + uint16(cpu.x)) & 0xff) | (uint16(data) << 8)
		cpu.next = fnAE_ABSX2
	}
}

func fnAE_ABSX2(cpu *Core) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = _opTable[cpu.op]
}

func fnAE_ABSY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnAE_ABSY1
}

func fnAE_ABSY1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.ar+uint16(cpu.y) < 0x100 {
		cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (uint16(data) << 8)
		cpu.next = _opTable[cpu.op]
	} else {
		cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (uint16(data) << 8)
		cpu.next = fnAE_ABSY2
	}
}

func fnAE_ABSY2(cpu *Core) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = _opTable[cpu.op]
}

func fnAE_INDY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnAE_INDY1
}

func fnAE_INDY1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.next = fnAE_INDY2
}

func fnAE_INDY2(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read((cpu.ar2 + 1) & 0xff)
	if z := cpu.ar + uint16(cpu.y); z < 0x100 {
		cpu.ar = (z & 0xff) | (uint16(data) << 8)
		cpu.next = _opTable[cpu.op]
	} else {
		cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (uint16(data) << 8)
		cpu.next = fnAE_INDY3
	}
}

func fnAE_INDY3(cpu *Core) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = _opTable[cpu.op]
}

func fnM_ZERO(cpu *Core) {
	// Addressing modes: Read operand, write it back, no extra cycles (-> ar, rmw)
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnRMW
}

func fnM_ZEROX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnM_ZEROX1
}

func fnM_ZEROX1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
	cpu.next = fnRMW
}

/*
func fnM_ZEROY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnM_ZEROY1
}

func fnM_ZEROY1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
	cpu.next = fnRMW_DO_IT
}
*/

func fnM_ABS(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnM_ABS1
}

func fnM_ABS1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = fnRMW
}

func fnM_ABSX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnM_ABSX1
}

func fnM_ABSX1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.ar+uint16(cpu.x) < 0x100 {
		cpu.next = fnM_ABSX2
	} else {
		cpu.next = fnM_ABSX3
	}
	cpu.ar = (cpu.ar + uint16(cpu.x)&0xff) | (uint16(data) << 8)
}

func fnM_ABSX2(cpu *Core) {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.next = fnRMW
}

func fnM_ABSX3(cpu *Core) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = fnRMW
}

func fnM_ABSY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnM_ABSY1
}

func fnM_ABSY1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.ar+uint16(cpu.y) < 0x100 {
		cpu.next = fnM_ABSY2
	} else {
		cpu.next = fnM_ABSY3
	}
	cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (uint16(data) << 8)
}

func fnM_ABSY2(cpu *Core) {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.next = fnRMW
}

func fnM_ABSY3(cpu *Core) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = fnRMW
}

func fnM_INDX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnM_INDX1
}

func fnM_INDX1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar2)
	cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
	cpu.next = fnM_INDX2
}

func fnM_INDX2(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.next = fnM_INDX3
}

func fnM_INDX3(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read((cpu.ar2 + 1) & 0xff)
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = fnRMW
}

func fnM_INDY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnM_INDY1
}

func fnM_INDY1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.next = fnM_INDY2
}

func fnM_INDY2(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read((cpu.ar2 + 1) & 0xff)
	if cpu.ar+uint16(cpu.y) < 0x100 {
		cpu.next = fnM_INDY3
	} else {
		cpu.next = fnM_INDY4
	}
	cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (uint16(data) << 8)
}

func fnM_INDY3(cpu *Core) {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	_ = cpu.banks.Read(cpu.ar)
	cpu.next = fnRMW
}

func fnM_INDY4(cpu *Core) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = fnRMW
}

func fnRMW(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.rmw = cpu.banks.Read(cpu.ar)
	cpu.next = fnRMW1
}

func fnRMW1(cpu *Core) {
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.next = _opTable[cpu.op]
}

// Load group
func fnO_LDA(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.a = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = fnInit
}

func fnO_LDA_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.a = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = fnInit
}

func fnO_LDX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.x = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = fnInit
}

func fnO_LDX_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.x = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = fnInit
}

func fnO_LDY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.y = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = fnInit
}

func fnO_LDY_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.y = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = fnInit
}

// Store group
func fnO_STA(cpu *Core) {
	cpu.banks.Write(cpu.ar, cpu.a)
	cpu.next = fnInit
}

func fnO_STX(cpu *Core) {
	cpu.banks.Write(cpu.ar, cpu.x)
	cpu.next = fnInit
}

func fnO_STY(cpu *Core) {
	cpu.banks.Write(cpu.ar, cpu.y)
	cpu.next = fnInit
}

// Transfer group
func fnO_TAX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.x = cpu.a
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = fnInit
}

func fnO_TXA(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.a = cpu.x
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.next = fnInit
}

func fnO_TAY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.y = cpu.a
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = fnInit
}

func fnO_TYA(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.a = cpu.y
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.next = fnInit
}

func fnO_TSX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.x = cpu.sp
	cpu.nFlag = cpu.sp
	cpu.zFlag = cpu.sp
	cpu.next = fnInit
}

func fnO_TXS(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.sp = cpu.x
	cpu.next = fnInit
}

// Arithmetic group
func fnO_ADC(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.doADC(data)
	cpu.next = fnInit
}

func fnO_ADC_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.doADC(data)
	cpu.next = fnInit
}

func fnO_SBC(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.doSBC(data)
	cpu.next = fnInit
}

func fnO_SBC_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.doSBC(data)
	cpu.next = fnInit
}

// Increment/decrement group
func fnO_INX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.x++
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.next = fnInit
}

func fnO_DEX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.x--
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.next = fnInit
}

func fnO_INY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.y++
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.next = fnInit
}

func fnO_DEY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.y--
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.next = fnInit
}

func fnO_INC(cpu *Core) {
	v := cpu.rmw + 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = fnInit
}

func fnO_DEC(cpu *Core) {
	v := cpu.rmw - 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = fnInit
}

// Logic group
func fnO_AND(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a &= cpu.banks.Read(cpu.ar)
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = fnInit
}

func fnO_AND_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a &= cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = fnInit
}

func fnO_ORA(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a |= cpu.banks.Read(cpu.ar)
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = fnInit
}

func fnO_ORA_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a |= cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = fnInit
}

func fnO_EOR(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a ^= cpu.banks.Read(cpu.ar)
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = fnInit
}

func fnO_EOR_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a ^= cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = fnInit
}

// Compare group
func fnO_CMP(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.ar = uint16(cpu.a) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
	cpu.next = fnInit
}

func fnO_CMP_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.ar = uint16(cpu.a) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
	cpu.next = fnInit
}

func fnO_CPX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.ar = uint16(cpu.x) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
	cpu.next = fnInit
}

func fnO_CPX_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.ar = uint16(cpu.x) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
	cpu.next = fnInit
}

func fnO_CPY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.ar = uint16(cpu.y) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
	cpu.next = fnInit
}

func fnO_CPY_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.ar = uint16(cpu.y) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
	cpu.next = fnInit
}

// Bit-test group
func fnO_BIT(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.zFlag = cpu.a & data
	cpu.nFlag = data
	cpu.vFlag = data & 0x40
	cpu.next = fnInit
}

// Shift/rotate group
func fnO_ASL(cpu *Core) {
	cpu.cFlag = cpu.rmw & 0x80
	v := cpu.rmw << 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = fnInit
}

func fnO_ASL_A(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.cFlag = cpu.a & 0x80
	cpu.a <<= 1
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = fnInit
}

func fnO_LSR(cpu *Core) {
	cpu.cFlag = cpu.rmw & 0x01
	v := cpu.rmw >> 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = fnInit
}

func fnO_LSR_A(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.cFlag = cpu.a & 0x01
	cpu.a >>= 1
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = fnInit
}

func fnO_ROL(cpu *Core) {
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
	cpu.next = fnInit
}

func fnO_ROL_A(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
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
	cpu.next = fnInit
}

func fnO_ROR(cpu *Core) {
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
	cpu.next = fnInit
}

func fnO_ROR_A(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
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
	cpu.next = fnInit
}

// Stack group
func fnO_PHA(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = fnO_PHA1
}

func fnO_PHA1(cpu *Core) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, cpu.a)
	cpu.sp--
	cpu.next = fnInit
}

func fnO_PLA(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = fnO_PLA1
}

func fnO_PLA1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.sp++
	cpu.next = fnO_PLA2
}

func fnO_PLA2(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a = cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = fnInit
}

func fnO_PHP(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = fnO_PHP1
}

func fnO_PHP1(cpu *Core) {
	data := cpu.pushFlags(true)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|0x0100, data)
	cpu.sp--
	cpu.next = fnInit
}

func fnO_PLP(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = fnO_PLP1
}

func fnO_PLP1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.sp++
	cpu.next = fnO_PLP2
}

func fnO_PLP2(cpu *Core) {
	iFlagPrev := cpu.iFlag
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(uint16(cpu.sp) | 0x0100)
	cpu.popFlags(data)
	if iFlagPrev == 0 && cpu.iFlag != 0 {
		cpu.opFlags |= OpFlagIrqDisabled
	} else if iFlagPrev != 0 && cpu.iFlag == 0 {
		cpu.opFlags |= OpFlagIrqEnabled
	}
	cpu.next = fnInit
}

// Jump/branch group
func fnO_JMP(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnO_JMP1
}

func fnO_JMP1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc = (uint16(data) << 8) | cpu.ar
	cpu.next = fnInit
}

func fnO_JMP_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(cpu.ar))
	cpu.next = fnO_JMP_I1
}

func fnO_JMP_I1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(((cpu.ar + 1) & 0xff) | (cpu.ar & 0xff00))
	cpu.pc |= uint16(data) << 8
	cpu.next = fnInit
}

func fnO_JSR(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = fnO_JSR1
}

func fnO_JSR1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.next = fnO_JSR2
}

func fnO_JSR2(cpu *Core) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = fnO_JSR3
}

func fnO_JSR3(cpu *Core) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
	cpu.sp--
	cpu.next = fnO_JSR4
}

func fnO_JSR4(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.pc = cpu.ar | (uint16(data) << 8)
	cpu.next = fnInit
}

func fnO_RTS(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = fnO_RTS1
}

func fnO_RTS1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.sp++
	cpu.next = fnO_RTS2
}

func fnO_RTS2(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(uint16(cpu.sp) | 0x100))
	cpu.sp++
	cpu.next = fnO_RTS3
}

func fnO_RTS3(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.pc |= uint16(data) << 8
	cpu.next = fnO_RTS4
}

func fnO_RTS4(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.next = fnInit
}

func fnO_RTI(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = fnO_RTI1
}

func fnO_RTI1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.sp++
	cpu.next = fnO_RTI2
}

func fnO_RTI2(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(uint16(cpu.sp) | 0x0100)
	cpu.popFlags(data)
	cpu.sp++
	cpu.next = fnO_RTI3
}

func fnO_RTI3(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(uint16(cpu.sp) | 0x100))
	cpu.sp++
	cpu.next = fnO_RTI4
}

func fnO_RTI4(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.pc |= uint16(data) << 8
	cpu.next = fnInit
}

func fnO_BRK(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.next = fnO_BRK1
}

func fnO_BRK1(cpu *Core) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = fnO_BRK2
}

func fnO_BRK2(cpu *Core) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
	cpu.sp--
	cpu.next = fnO_BRK3
}

func fnO_BRK3(cpu *Core) {
	data := cpu.pushFlags(true)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|0x0100, data)
	cpu.sp--
	cpu.iFlag = 1
	// BRK interrupted by NMI?
	if cpu.pic.HasNMI() {
		cpu.pic.ClearNMI() // Simulate an edge-triggered input
		cpu.next = fnNMI5  // Jump to NMI sequence
	} else {
		cpu.next = fnO_BRK4
	}
}

func fnO_BRK4(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(0xfffe))
	cpu.next = fnO_BRK5
}

func fnO_BRK5(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(0xffff)
	cpu.pc |= uint16(data) << 8
	cpu.next = fnInit
}

func fnO_BCS(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.cFlag == 0 {
		cpu.next = fnInit
	} else {
		cpu.branch(data)
	}
}

func fnO_BCC(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.cFlag != 0 {
		cpu.next = fnInit
	} else {
		cpu.branch(data)
	}
}

func fnO_BEQ(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.zFlag != 0 {
		cpu.next = fnInit
	} else {
		cpu.branch(data)
	}
}

func fnO_BNE(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.zFlag == 0 {
		cpu.next = fnInit
	} else {
		cpu.branch(data)
	}
}

func fnO_BVS(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.overflowBranch != nil && cpu.overflowBranch() {
		cpu.vFlag = 1
	}
	if cpu.vFlag == 0 {
		cpu.next = fnInit
	} else {
		cpu.branch(data)
	}
}

func fnO_BVC(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.overflowBranch != nil && cpu.overflowBranch() {
		cpu.vFlag = 1
	}
	if cpu.vFlag != 0 {
		cpu.next = fnInit
	} else {
		cpu.branch(data)
	}
}

func fnO_BMI(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if (cpu.nFlag & 0x80) == 0 {
		cpu.next = fnInit
	} else {
		cpu.branch(data)
	}
}

func fnO_BPL(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if (cpu.nFlag & 0x80) != 0 {
		cpu.next = fnInit
	} else {
		cpu.branch(data)
	}
}

func fnO_BRANCH_NP(cpu *Core) {
	// No page crossed
	cpu.opFlags |= OpFlagIntDelayed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc = cpu.ar
	cpu.next = fnInit
}

func fnO_BRANCH_BP(cpu *Core) {
	// Page crossed, branch backwards
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc = cpu.ar
	cpu.next = fnO_BRANCH_BP1
}

func fnO_BRANCH_BP1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc + 0x100)
	cpu.next = fnInit
}

func fnO_BRANCH_FP(cpu *Core) {
	// Page crossed, branch forwards
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc = cpu.ar
	cpu.next = fnO_BRANCH_FP1
}

func fnO_BRANCH_FP1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc - 0x100)
	cpu.next = fnInit
}

// Flag group
func fnO_SEC(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.cFlag = 1
	cpu.next = fnInit
}

func fnO_CLC(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.cFlag = 0
	cpu.next = fnInit
}

func fnO_SED(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.dFlag = 1
	cpu.next = fnInit
}

func fnO_CLD(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.dFlag = 0
	cpu.next = fnInit
}

func fnO_SEI(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	if cpu.iFlag == 0 {
		cpu.opFlags |= OpFlagIrqDisabled
	}
	cpu.iFlag = 1
	cpu.next = fnInit
}

func fnO_CLI(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	if cpu.iFlag == 0 {
		cpu.opFlags |= OpFlagIrqEnabled
	}
	cpu.iFlag = 0
	cpu.next = fnInit
}

func fnO_CLV(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.vFlag = 0
	cpu.next = fnInit
}

// NOP group
func fnO_NOP(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = fnInit
}

// Undocumented functions start here

// NOP group
func fnO_NOP_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.next = fnInit
}

func fnO_NOP_A(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.next = fnInit
}

// Load A/X group
func fnO_LAX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.x = cpu.banks.Read(cpu.ar)
	cpu.a = cpu.x
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = fnInit
}

// Store A/X group
func fnO_SAX(cpu *Core) {
	cpu.banks.Write(cpu.ar, cpu.a&cpu.x)
	cpu.next = fnInit
}

// ASL/ORA group
func fnO_SLO(cpu *Core) {
	cpu.cFlag = cpu.rmw & 0x80
	cpu.rmw <<= 1
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.a |= cpu.rmw
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = fnInit
}

// ROL/AND group
func fnO_RLA(cpu *Core) {
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
	cpu.next = fnInit
}

// LSR/EOR group
func fnO_SRE(cpu *Core) {
	cpu.cFlag = cpu.rmw & 0x01
	cpu.rmw >>= 1
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.a ^= cpu.rmw
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = fnInit
}

// ROR/ADC group
func fnO_RRA(cpu *Core) {
	tmp := cpu.rmw & 0x01
	if cpu.cFlag != 0 {
		cpu.rmw = (cpu.rmw >> 1) | 0x80
	} else {
		cpu.rmw = cpu.rmw >> 1
	}
	cpu.cFlag = tmp
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.doADC(cpu.rmw)
	cpu.next = fnInit
}

// DEC/CMP group
func fnO_DCP(cpu *Core) {
	cpu.rmw--
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.ar = uint16(cpu.a) - uint16(cpu.rmw)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
	cpu.next = fnInit
}

// INC/SBC group
func fnO_ISB(cpu *Core) {
	cpu.rmw++
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.doSBC(cpu.rmw)
	cpu.next = fnInit
}

// Complex functions
func fnO_ANC_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a &= cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.cFlag = cpu.nFlag & 0x80
	cpu.next = fnInit
}

func fnO_ASR_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a &= cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.cFlag = cpu.a & 0x01
	cpu.a >>= 1
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = fnInit
}

func fnO_ARR_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
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
	cpu.next = fnInit
}

func fnO_ANE_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.a = (cpu.a | 0xee) & cpu.x & data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = fnInit
}

func fnO_LXA_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.x = (cpu.a | 0xee) & data
	cpu.a = cpu.x
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = fnInit
}

func fnO_SBX_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.ar = (uint16(cpu.x) & uint16(cpu.a)) - uint16(data)
	cpu.x = uint8(cpu.ar)
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
	cpu.next = fnInit
}

func fnO_LAS(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.sp = data & cpu.sp
	cpu.x = cpu.sp
	cpu.a = cpu.x
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = fnInit
}

func fnO_SHS(cpu *Core) {
	// ar2 contains the high byte of the operand address
	cpu.sp = cpu.a & cpu.x
	cpu.banks.Write(cpu.ar, uint8((cpu.ar2+1)&uint16(cpu.sp)))
	cpu.next = fnInit
}

func fnO_SHY(cpu *Core) {
	// ar2 contains the high byte of the operand address
	cpu.banks.Write(cpu.ar, uint8(uint16(cpu.y)&(cpu.ar2+1)))
	cpu.next = fnInit
}

func fnO_SHX(cpu *Core) {
	// ar2 contains the high byte of the operand address
	cpu.banks.Write(cpu.ar, uint8(uint16(cpu.x)&(cpu.ar2+1)))
	cpu.next = fnInit
}

func fnO_SHA(cpu *Core) {
	// ar2 contains the high byte of the operand address
	cpu.banks.Write(cpu.ar, uint8(uint16(cpu.a)&uint16(cpu.x)&(cpu.ar2+1)))
	cpu.next = fnInit
}

func fnI_ILL_OP(cpu *Core) {
	cpu.illegalOp(cpu.op, cpu.pc-1)
}
