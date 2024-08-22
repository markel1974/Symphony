package mos6510fn

import "github.com/markel1974/c64emu/src/flag"

func instInit(cpu *Core) {
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
				cpu.next = instNMI
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
					cpu.next = instIRQ
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
func instIRQ(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instIRQ1
}

func instIRQ1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instIRQ2
}

func instIRQ2(cpu *Core) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = instIRQ3
}

func instIRQ3(cpu *Core) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
	cpu.sp--
	cpu.next = instIRQ4
}

func instIRQ4(cpu *Core) {
	data := cpu.pushFlags(false)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|0x0100, data)
	cpu.sp--
	cpu.iFlag = 1
	cpu.next = instIRQ5
}

func instIRQ5(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(0xfffe))
	cpu.next = instIRQ6
}

func instIRQ6(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(0xffff)
	cpu.pc |= uint16(data) << 8
	cpu.next = instInit
}

// NMI
func instNMI(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instNMI1
}

func instNMI1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instNMI2
}

func instNMI2(cpu *Core) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = instNMI3
}

func instNMI3(cpu *Core) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
	cpu.sp--
	cpu.next = instNMI4
}

func instNMI4(cpu *Core) {
	data := cpu.pushFlags(false)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|0x0100, data)
	cpu.sp--
	cpu.iFlag = 1
	cpu.next = instNMI5
}

func instNMI5(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(0xfffa))
	cpu.next = instNMI6
}

func instNMI6(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(0xfffb)
	cpu.pc |= uint16(data) << 8
	cpu.next = instInit
}

// Addressing modes: Fetch effective address, no extra cycles (-> ar)
func instA_ZERO(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = _opTable[cpu.op]
}

func instA_ZEROX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instA_ZEROX1
}

func instA_ZEROX1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
	cpu.next = _opTable[cpu.op]
}

func instA_ZEROY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instA_ZEROY1
}

func instA_ZEROY1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
	cpu.next = _opTable[cpu.op]
}

func instA_ABS(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instA_ABS1
}

func instA_ABS1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = _opTable[cpu.op]
}

func instA_ABSX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instA_ABSX1
}

func instA_ABSX1(cpu *Core) {
	// Note: Some undocumented functions rely on the value of ar2
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	if cpu.ar+uint16(cpu.x) < 0x100 {
		cpu.next = instA_ABSX2
	} else {
		cpu.next = instA_ABSX3
	}
	cpu.ar = ((cpu.ar + uint16(cpu.x)) & 0xff) | (cpu.ar2 << 8)
}

func instA_ABSX2(cpu *Core) {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.next = _opTable[cpu.op]
}

func instA_ABSX3(cpu *Core) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = _opTable[cpu.op]
}

func instA_ABSY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instA_ABSY1
}

func instA_ABSY1(cpu *Core) {
	// Note: Some undocumented functions rely on the value of ar2
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	if cpu.ar+uint16(cpu.y) < 0x100 {
		cpu.next = instA_ABSY2
	} else {
		cpu.next = instA_ABSY3
	}
	cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (cpu.ar2 << 8)
}

func instA_ABSY2(cpu *Core) {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.next = _opTable[cpu.op]
}

func instA_ABSY3(cpu *Core) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = _opTable[cpu.op]
}

func instA_INDX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instA_INDX1
}

func instA_INDX1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar2)
	cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
	cpu.next = instA_INDX2
}

func instA_INDX2(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.next = instA_INDX3
}

func instA_INDX3(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read((cpu.ar2 + 1) & 0xff)
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = _opTable[cpu.op]
}

func instA_INDY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instA_INDY1
}

func instA_INDY1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.next = instA_INDY2
}

func instA_INDY2(cpu *Core) {
	// Note: Some undocumented functions rely on the value of ar2
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read((cpu.ar2 + 1) & 0xff))
	if cpu.ar+uint16(cpu.y) < 0x100 {
		cpu.next = instA_INDY3
	} else {
		cpu.next = instA_INDY4
	}
	cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (cpu.ar2 << 8)
}

func instA_INDY3(cpu *Core) {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.next = _opTable[cpu.op]
}

func instA_INDY4(cpu *Core) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = _opTable[cpu.op]
}

func instAE_ABSX(cpu *Core) {
	// Addressing modes: Fetch effective address, extra cycle on page crossing (-> ar)
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instAE_ABSX1
}

func instAE_ABSX1(cpu *Core) {
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
		cpu.next = instAE_ABSX2
	}
}

func instAE_ABSX2(cpu *Core) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = _opTable[cpu.op]
}

func instAE_ABSY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instAE_ABSY1
}

func instAE_ABSY1(cpu *Core) {
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
		cpu.next = instAE_ABSY2
	}
}

func instAE_ABSY2(cpu *Core) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = _opTable[cpu.op]
}

func instAE_INDY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instAE_INDY1
}

func instAE_INDY1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.next = instAE_INDY2
}

func instAE_INDY2(cpu *Core) {
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
		cpu.next = instAE_INDY3
	}
}

func instAE_INDY3(cpu *Core) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = _opTable[cpu.op]
}

func instM_ZERO(cpu *Core) {
	// Addressing modes: Read operand, write it back, no extra cycles (-> ar, rmw)
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instRMW
}

func instM_ZEROX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instM_ZEROX1
}

func instM_ZEROX1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
	cpu.next = instRMW
}

/*
func instM_ZEROY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instM_ZEROY1
}

func instM_ZEROY1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
	cpu.next = instRMW_DO_IT
}
*/

func instM_ABS(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instM_ABS1
}

func instM_ABS1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = instRMW
}

func instM_ABSX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instM_ABSX1
}

func instM_ABSX1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.ar+uint16(cpu.x) < 0x100 {
		cpu.next = instM_ABSX2
	} else {
		cpu.next = instM_ABSX3
	}
	cpu.ar = (cpu.ar + uint16(cpu.x)&0xff) | (uint16(data) << 8)
}

func instM_ABSX2(cpu *Core) {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.next = instRMW
}

func instM_ABSX3(cpu *Core) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = instRMW
}

func instM_ABSY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instM_ABSY1
}

func instM_ABSY1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.ar+uint16(cpu.y) < 0x100 {
		cpu.next = instM_ABSY2
	} else {
		cpu.next = instM_ABSY3
	}
	cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (uint16(data) << 8)
}

func instM_ABSY2(cpu *Core) {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.next = instRMW
}

func instM_ABSY3(cpu *Core) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = instRMW
}

func instM_INDX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instM_INDX1
}

func instM_INDX1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar2)
	cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
	cpu.next = instM_INDX2
}

func instM_INDX2(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.next = instM_INDX3
}

func instM_INDX3(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read((cpu.ar2 + 1) & 0xff)
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = instRMW
}

func instM_INDY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instM_INDY1
}

func instM_INDY1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.next = instM_INDY2
}

func instM_INDY2(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read((cpu.ar2 + 1) & 0xff)
	if cpu.ar+uint16(cpu.y) < 0x100 {
		cpu.next = instM_INDY3
	} else {
		cpu.next = instM_INDY4
	}
	cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (uint16(data) << 8)
}

func instM_INDY3(cpu *Core) {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	_ = cpu.banks.Read(cpu.ar)
	cpu.next = instRMW
}

func instM_INDY4(cpu *Core) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = instRMW
}

func instRMW(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.rmw = cpu.banks.Read(cpu.ar)
	cpu.next = instRMW1
}

func instRMW1(cpu *Core) {
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.next = _opTable[cpu.op]
}

// Load group
func instO_LDA(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.a = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instInit
}

func instO_LDA_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.a = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instInit
}

func instO_LDX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.x = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instInit
}

func instO_LDX_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.x = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instInit
}

func instO_LDY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.y = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instInit
}

func instO_LDY_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.y = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instInit
}

// Store group
func instO_STA(cpu *Core) {
	cpu.banks.Write(cpu.ar, cpu.a)
	cpu.next = instInit
}

func instO_STX(cpu *Core) {
	cpu.banks.Write(cpu.ar, cpu.x)
	cpu.next = instInit
}

func instO_STY(cpu *Core) {
	cpu.banks.Write(cpu.ar, cpu.y)
	cpu.next = instInit
}

// Transfer group
func instO_TAX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.x = cpu.a
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

func instO_TXA(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.a = cpu.x
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.next = instInit
}

func instO_TAY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.y = cpu.a
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

func instO_TYA(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.a = cpu.y
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.next = instInit
}

func instO_TSX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.x = cpu.sp
	cpu.nFlag = cpu.sp
	cpu.zFlag = cpu.sp
	cpu.next = instInit
}

func instO_TXS(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.sp = cpu.x
	cpu.next = instInit
}

// Arithmetic group
func instO_ADC(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.doADC(data)
	cpu.next = instInit
}

func instO_ADC_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.doADC(data)
	cpu.next = instInit
}

func instO_SBC(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.doSBC(data)
	cpu.next = instInit
}

func instO_SBC_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.doSBC(data)
	cpu.next = instInit
}

// Increment/decrement group
func instO_INX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.x++
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.next = instInit
}

func instO_DEX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.x--
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.next = instInit
}

func instO_INY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.y++
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.next = instInit
}

func instO_DEY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.y--
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.next = instInit
}

func instO_INC(cpu *Core) {
	v := cpu.rmw + 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instInit
}

func instO_DEC(cpu *Core) {
	v := cpu.rmw - 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instInit
}

// Logic group
func instO_AND(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a &= cpu.banks.Read(cpu.ar)
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

func instO_AND_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a &= cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

func instO_ORA(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a |= cpu.banks.Read(cpu.ar)
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

func instO_ORA_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a |= cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

func instO_EOR(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a ^= cpu.banks.Read(cpu.ar)
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

func instO_EOR_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a ^= cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

// Compare group
func instO_CMP(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.ar = uint16(cpu.a) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
	cpu.next = instInit
}

func instO_CMP_I(cpu *Core) {
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
	cpu.next = instInit
}

func instO_CPX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.ar = uint16(cpu.x) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
	cpu.next = instInit
}

func instO_CPX_I(cpu *Core) {
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
	cpu.next = instInit
}

func instO_CPY(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.ar = uint16(cpu.y) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
	cpu.next = instInit
}

func instO_CPY_I(cpu *Core) {
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
	cpu.next = instInit
}

// Bit-test group
func instO_BIT(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.zFlag = cpu.a & data
	cpu.nFlag = data
	cpu.vFlag = data & 0x40
	cpu.next = instInit
}

// Shift/rotate group
func instO_ASL(cpu *Core) {
	cpu.cFlag = cpu.rmw & 0x80
	v := cpu.rmw << 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instInit
}

func instO_ASL_A(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.cFlag = cpu.a & 0x80
	cpu.a <<= 1
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

func instO_LSR(cpu *Core) {
	cpu.cFlag = cpu.rmw & 0x01
	v := cpu.rmw >> 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instInit
}

func instO_LSR_A(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.cFlag = cpu.a & 0x01
	cpu.a >>= 1
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

func instO_ROL(cpu *Core) {
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
	cpu.next = instInit
}

func instO_ROL_A(cpu *Core) {
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
	cpu.next = instInit
}

func instO_ROR(cpu *Core) {
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
	cpu.next = instInit
}

func instO_ROR_A(cpu *Core) {
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
	cpu.next = instInit
}

// Stack group
func instO_PHA(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instO_PHA1
}

func instO_PHA1(cpu *Core) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, cpu.a)
	cpu.sp--
	cpu.next = instInit
}

func instO_PLA(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instO_PLA1
}

func instO_PLA1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.sp++
	cpu.next = instO_PLA2
}

func instO_PLA2(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a = cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

func instO_PHP(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instO_PHP1
}

func instO_PHP1(cpu *Core) {
	data := cpu.pushFlags(true)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|0x0100, data)
	cpu.sp--
	cpu.next = instInit
}

func instO_PLP(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instO_PLP1
}

func instO_PLP1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.sp++
	cpu.next = instO_PLP2
}

func instO_PLP2(cpu *Core) {
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
	cpu.next = instInit
}

// Jump/branch group
func instO_JMP(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instO_JMP1
}

func instO_JMP1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc = (uint16(data) << 8) | cpu.ar
	cpu.next = instInit
}

func instO_JMP_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(cpu.ar))
	cpu.next = instO_JMP_I1
}

func instO_JMP_I1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(((cpu.ar + 1) & 0xff) | (cpu.ar & 0xff00))
	cpu.pc |= uint16(data) << 8
	cpu.next = instInit
}

func instO_JSR(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instO_JSR1
}

func instO_JSR1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.next = instO_JSR2
}

func instO_JSR2(cpu *Core) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = instO_JSR3
}

func instO_JSR3(cpu *Core) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
	cpu.sp--
	cpu.next = instO_JSR4
}

func instO_JSR4(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.pc = cpu.ar | (uint16(data) << 8)
	cpu.next = instInit
}

func instO_RTS(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instO_RTS1
}

func instO_RTS1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.sp++
	cpu.next = instO_RTS2
}

func instO_RTS2(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(uint16(cpu.sp) | 0x100))
	cpu.sp++
	cpu.next = instO_RTS3
}

func instO_RTS3(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.pc |= uint16(data) << 8
	cpu.next = instO_RTS4
}

func instO_RTS4(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.next = instInit
}

func instO_RTI(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instO_RTI1
}

func instO_RTI1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.sp++
	cpu.next = instO_RTI2
}

func instO_RTI2(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(uint16(cpu.sp) | 0x0100)
	cpu.popFlags(data)
	cpu.sp++
	cpu.next = instO_RTI3
}

func instO_RTI3(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(uint16(cpu.sp) | 0x100))
	cpu.sp++
	cpu.next = instO_RTI4
}

func instO_RTI4(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.pc |= uint16(data) << 8
	cpu.next = instInit
}

func instO_BRK(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.next = instO_BRK1
}

func instO_BRK1(cpu *Core) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = instO_BRK2
}

func instO_BRK2(cpu *Core) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
	cpu.sp--
	cpu.next = instO_BRK3
}

func instO_BRK3(cpu *Core) {
	data := cpu.pushFlags(true)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|0x0100, data)
	cpu.sp--
	cpu.iFlag = 1
	// BRK interrupted by NMI?
	if cpu.pic.HasNMI() {
		cpu.pic.ClearNMI()  // Simulate an edge-triggered input
		cpu.next = instNMI5 // Jump to NMI sequence
	} else {
		cpu.next = instO_BRK4
	}
}

func instO_BRK4(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(0xfffe))
	cpu.next = instO_BRK5
}

func instO_BRK5(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(0xffff)
	cpu.pc |= uint16(data) << 8
	cpu.next = instInit
}

func instO_BCS(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.cFlag == 0 {
		cpu.next = instInit
	} else {
		cpu.branch(data)
	}
}

func instO_BCC(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.cFlag != 0 {
		cpu.next = instInit
	} else {
		cpu.branch(data)
	}
}

func instO_BEQ(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.zFlag != 0 {
		cpu.next = instInit
	} else {
		cpu.branch(data)
	}
}

func instO_BNE(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.zFlag == 0 {
		cpu.next = instInit
	} else {
		cpu.branch(data)
	}
}

func instO_BVS(cpu *Core) {
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
		cpu.next = instInit
	} else {
		cpu.branch(data)
	}
}

func instO_BVC(cpu *Core) {
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
		cpu.next = instInit
	} else {
		cpu.branch(data)
	}
}

func instO_BMI(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if (cpu.nFlag & 0x80) == 0 {
		cpu.next = instInit
	} else {
		cpu.branch(data)
	}
}

func instO_BPL(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if (cpu.nFlag & 0x80) != 0 {
		cpu.next = instInit
	} else {
		cpu.branch(data)
	}
}

func instO_BRANCH_NP(cpu *Core) {
	// No page crossed
	cpu.opFlags |= OpFlagIntDelayed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc = cpu.ar
	cpu.next = instInit
}

func instO_BRANCH_BP(cpu *Core) {
	// Page crossed, branch backwards
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc = cpu.ar
	cpu.next = instO_BRANCH_BP1
}

func instO_BRANCH_BP1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc + 0x100)
	cpu.next = instInit
}

func instO_BRANCH_FP(cpu *Core) {
	// Page crossed, branch forwards
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc = cpu.ar
	cpu.next = instO_BRANCH_FP1
}

func instO_BRANCH_FP1(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc - 0x100)
	cpu.next = instInit
}

// Flag group
func instO_SEC(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.cFlag = 1
	cpu.next = instInit
}

func instO_CLC(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.cFlag = 0
	cpu.next = instInit
}

func instO_SED(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.dFlag = 1
	cpu.next = instInit
}

func instO_CLD(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.dFlag = 0
	cpu.next = instInit
}

func instO_SEI(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	if cpu.iFlag == 0 {
		cpu.opFlags |= OpFlagIrqDisabled
	}
	cpu.iFlag = 1
	cpu.next = instInit
}

func instO_CLI(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	if cpu.iFlag == 0 {
		cpu.opFlags |= OpFlagIrqEnabled
	}
	cpu.iFlag = 0
	cpu.next = instInit
}

func instO_CLV(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.vFlag = 0
	cpu.next = instInit
}

// NOP group
func instO_NOP(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instInit
}

// Undocumented functions start here

// NOP group
func instO_NOP_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.next = instInit
}

func instO_NOP_A(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.next = instInit
}

// Load A/X group
func instO_LAX(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.x = cpu.banks.Read(cpu.ar)
	cpu.a = cpu.x
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

// Store A/X group
func instO_SAX(cpu *Core) {
	cpu.banks.Write(cpu.ar, cpu.a&cpu.x)
	cpu.next = instInit
}

// ASL/ORA group
func instO_SLO(cpu *Core) {
	cpu.cFlag = cpu.rmw & 0x80
	cpu.rmw <<= 1
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.a |= cpu.rmw
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

// ROL/AND group
func instO_RLA(cpu *Core) {
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
	cpu.next = instInit
}

// LSR/EOR group
func instO_SRE(cpu *Core) {
	cpu.cFlag = cpu.rmw & 0x01
	cpu.rmw >>= 1
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.a ^= cpu.rmw
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

// ROR/ADC group
func instO_RRA(cpu *Core) {
	tmp := cpu.rmw & 0x01
	if cpu.cFlag != 0 {
		cpu.rmw = (cpu.rmw >> 1) | 0x80
	} else {
		cpu.rmw = cpu.rmw >> 1
	}
	cpu.cFlag = tmp
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.doADC(cpu.rmw)
	cpu.next = instInit
}

// DEC/CMP group
func instO_DCP(cpu *Core) {
	cpu.rmw--
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.ar = uint16(cpu.a) - uint16(cpu.rmw)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
	cpu.next = instInit
}

// INC/SBC group
func instO_ISB(cpu *Core) {
	cpu.rmw++
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.doSBC(cpu.rmw)
	cpu.next = instInit
}

// Complex functions
func instO_ANC_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a &= cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.cFlag = cpu.nFlag & 0x80
	cpu.next = instInit
}

func instO_ASR_I(cpu *Core) {
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
	cpu.next = instInit
}

func instO_ARR_I(cpu *Core) {
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
	cpu.next = instInit
}

func instO_ANE_I(cpu *Core) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.a = (cpu.a | 0xee) & cpu.x & data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

func instO_LXA_I(cpu *Core) {
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
	cpu.next = instInit
}

func instO_SBX_I(cpu *Core) {
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
	cpu.next = instInit
}

func instO_LAS(cpu *Core) {
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
	cpu.next = instInit
}

func instO_SHS(cpu *Core) {
	// ar2 contains the high byte of the operand address
	cpu.sp = cpu.a & cpu.x
	cpu.banks.Write(cpu.ar, uint8((cpu.ar2+1)&uint16(cpu.sp)))
	cpu.next = instInit
}

func instO_SHY(cpu *Core) {
	// ar2 contains the high byte of the operand address
	cpu.banks.Write(cpu.ar, uint8(uint16(cpu.y)&(cpu.ar2+1)))
	cpu.next = instInit
}

func instO_SHX(cpu *Core) {
	// ar2 contains the high byte of the operand address
	cpu.banks.Write(cpu.ar, uint8(uint16(cpu.x)&(cpu.ar2+1)))
	cpu.next = instInit
}

func instO_SHA(cpu *Core) {
	// ar2 contains the high byte of the operand address
	cpu.banks.Write(cpu.ar, uint8(uint16(cpu.a)&uint16(cpu.x)&(cpu.ar2+1)))
	cpu.next = instInit
}

func instI_ILL_OP(cpu *Core) {
	cpu.illegalOp(cpu.op, cpu.pc-1)
}
