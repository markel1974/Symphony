package mos6510

import "github.com/markel1974/c64emu/src/flag"

func instInit(cpu *CPU) {
	if cpu.pic.HasAny() {
		if cpu.pic.HasReset() {
			cpu.Reset()
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
func instIRQ(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instIRQ1
}

func instIRQ1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instIRQ2
}

func instIRQ2(cpu *CPU) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = instIRQ3
}

func instIRQ3(cpu *CPU) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
	cpu.sp--
	cpu.next = instIRQ4
}

func instIRQ4(cpu *CPU) {
	data := cpu.pushFlags(false)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|0x0100, data)
	cpu.sp--
	cpu.iFlag = 1
	cpu.next = instIRQ5
}

func instIRQ5(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(0xfffe))
	cpu.next = instIRQ6
}

func instIRQ6(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(0xffff)
	cpu.pc |= uint16(data) << 8
	cpu.next = instInit
}

// NMI
func instNMI(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instNMI1
}

func instNMI1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instNMI2
}

func instNMI2(cpu *CPU) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = instNMI3
}

func instNMI3(cpu *CPU) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
	cpu.sp--
	cpu.next = instNMI4
}

func instNMI4(cpu *CPU) {
	data := cpu.pushFlags(false)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|0x0100, data)
	cpu.sp--
	cpu.iFlag = 1
	cpu.next = instNMI5
}

func instNMI5(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(0xfffa))
	cpu.next = instNMI6
}

func instNMI6(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(0xfffb)
	cpu.pc |= uint16(data) << 8
	cpu.next = instInit
}

// Addressing modes: Fetch effective address, no extra cycles (-> ar)
func instA_ZERO(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = _opTable[cpu.op]
}

func instA_ZEROX(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instA_ZEROX1
}

func instA_ZEROX1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
	cpu.next = _opTable[cpu.op]
}

func instA_ZEROY(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instA_ZEROY1
}

func instA_ZEROY1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
	cpu.next = _opTable[cpu.op]
}

func instA_ABS(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instA_ABS1
}

func instA_ABS1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = _opTable[cpu.op]
}

func instA_ABSX(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instA_ABSX1
}

func instA_ABSX1(cpu *CPU) {
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

func instA_ABSX2(cpu *CPU) {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.next = _opTable[cpu.op]
}

func instA_ABSX3(cpu *CPU) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = _opTable[cpu.op]
}

func instA_ABSY(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instA_ABSY1
}

func instA_ABSY1(cpu *CPU) {
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

func instA_ABSY2(cpu *CPU) {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.next = _opTable[cpu.op]
}

func instA_ABSY3(cpu *CPU) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = _opTable[cpu.op]
}

func instA_INDX(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instA_INDX1
}

func instA_INDX1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar2)
	cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
	cpu.next = instA_INDX2
}

func instA_INDX2(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.next = instA_INDX3
}

func instA_INDX3(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read((cpu.ar2 + 1) & 0xff)
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = _opTable[cpu.op]
}

func instA_INDY(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instA_INDY1
}

func instA_INDY1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.next = instA_INDY2
}

func instA_INDY2(cpu *CPU) {
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

func instA_INDY3(cpu *CPU) {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.next = _opTable[cpu.op]
}

func instA_INDY4(cpu *CPU) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = _opTable[cpu.op]
}

func instAE_ABSX(cpu *CPU) {
	// Addressing modes: Fetch effective address, extra cycle on page crossing (-> ar)
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instAE_ABSX1
}

func instAE_ABSX1(cpu *CPU) {
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

func instAE_ABSX2(cpu *CPU) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = _opTable[cpu.op]
}

func instAE_ABSY(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instAE_ABSY1
}

func instAE_ABSY1(cpu *CPU) {
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

func instAE_ABSY2(cpu *CPU) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = _opTable[cpu.op]
}

func instAE_INDY(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instAE_INDY1
}

func instAE_INDY1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.next = instAE_INDY2
}

func instAE_INDY2(cpu *CPU) {
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

func instAE_INDY3(cpu *CPU) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = _opTable[cpu.op]
}

func instM_ZERO(cpu *CPU) {
	// Addressing modes: Read operand, write it back, no extra cycles (-> ar, rmw)
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instRMW
}

func instM_ZEROX(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instM_ZEROX1
}

func instM_ZEROX1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
	cpu.next = instRMW
}

/*
func instM_ZEROY(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instM_ZEROY1
}

func instM_ZEROY1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
	cpu.next = instRMW_DO_IT
}
*/

func instM_ABS(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instM_ABS1
}

func instM_ABS1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = instRMW
}

func instM_ABSX(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instM_ABSX1
}

func instM_ABSX1(cpu *CPU) {
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

func instM_ABSX2(cpu *CPU) {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.next = instRMW
}

func instM_ABSX3(cpu *CPU) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = instRMW
}

func instM_ABSY(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instM_ABSY1
}

func instM_ABSY1(cpu *CPU) {
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

func instM_ABSY2(cpu *CPU) {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.next = instRMW
}

func instM_ABSY3(cpu *CPU) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = instRMW
}

func instM_INDX(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instM_INDX1
}

func instM_INDX1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar2)
	cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
	cpu.next = instM_INDX2
}

func instM_INDX2(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.next = instM_INDX3
}

func instM_INDX3(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read((cpu.ar2 + 1) & 0xff)
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = instRMW
}

func instM_INDY(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instM_INDY1
}

func instM_INDY1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.next = instM_INDY2
}

func instM_INDY2(cpu *CPU) {
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

func instM_INDY3(cpu *CPU) {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	_ = cpu.banks.Read(cpu.ar)
	cpu.next = instRMW
}

func instM_INDY4(cpu *CPU) {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.next = instRMW
}

func instRMW(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.rmw = cpu.banks.Read(cpu.ar)
	cpu.next = instRMW1
}

func instRMW1(cpu *CPU) {
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.next = _opTable[cpu.op]
}

// Load group
func instO_LDA(cpu *CPU) {
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

func instO_LDA_I(cpu *CPU) {
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

func instO_LDX(cpu *CPU) {
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

func instO_LDX_I(cpu *CPU) {
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

func instO_LDY(cpu *CPU) {
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

func instO_LDY_I(cpu *CPU) {
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
func instO_STA(cpu *CPU) {
	cpu.banks.Write(cpu.ar, cpu.a)
	cpu.next = instInit
}

func instO_STX(cpu *CPU) {
	cpu.banks.Write(cpu.ar, cpu.x)
	cpu.next = instInit
}

func instO_STY(cpu *CPU) {
	cpu.banks.Write(cpu.ar, cpu.y)
	cpu.next = instInit
}

// Transfer group
func instO_TAX(cpu *CPU) {
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

func instO_TXA(cpu *CPU) {
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

func instO_TAY(cpu *CPU) {
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

func instO_TYA(cpu *CPU) {
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

func instO_TSX(cpu *CPU) {
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

func instO_TXS(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.sp = cpu.x
	cpu.next = instInit
}

// Arithmetic group
func instO_ADC(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.doADC(data)
	cpu.next = instInit
}

func instO_ADC_I(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.doADC(data)
	cpu.next = instInit
}

func instO_SBC(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.doSBC(data)
	cpu.next = instInit
}

func instO_SBC_I(cpu *CPU) {
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
func instO_INX(cpu *CPU) {
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

func instO_DEX(cpu *CPU) {
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

func instO_INY(cpu *CPU) {
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

func instO_DEY(cpu *CPU) {
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

func instO_INC(cpu *CPU) {
	v := cpu.rmw + 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instInit
}

func instO_DEC(cpu *CPU) {
	v := cpu.rmw - 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instInit
}

// Logic group
func instO_AND(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a &= cpu.banks.Read(cpu.ar)
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

func instO_AND_I(cpu *CPU) {
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

func instO_ORA(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a |= cpu.banks.Read(cpu.ar)
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

func instO_ORA_I(cpu *CPU) {
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

func instO_EOR(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a ^= cpu.banks.Read(cpu.ar)
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

func instO_EOR_I(cpu *CPU) {
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
func instO_CMP(cpu *CPU) {
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

func instO_CMP_I(cpu *CPU) {
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

func instO_CPX(cpu *CPU) {
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

func instO_CPX_I(cpu *CPU) {
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

func instO_CPY(cpu *CPU) {
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

func instO_CPY_I(cpu *CPU) {
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
func instO_BIT(cpu *CPU) {
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
func instO_ASL(cpu *CPU) {
	cpu.cFlag = cpu.rmw & 0x80
	v := cpu.rmw << 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instInit
}

func instO_ASL_A(cpu *CPU) {
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

func instO_LSR(cpu *CPU) {
	cpu.cFlag = cpu.rmw & 0x01
	v := cpu.rmw >> 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instInit
}

func instO_LSR_A(cpu *CPU) {
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

func instO_ROL(cpu *CPU) {
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

func instO_ROL_A(cpu *CPU) {
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

func instO_ROR(cpu *CPU) {
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

func instO_ROR_A(cpu *CPU) {
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
func instO_PHA(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instO_PHA1
}

func instO_PHA1(cpu *CPU) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, cpu.a)
	cpu.sp--
	cpu.next = instInit
}

func instO_PLA(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instO_PLA1
}

func instO_PLA1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.sp++
	cpu.next = instO_PLA2
}

func instO_PLA2(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a = cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

func instO_PHP(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instO_PHP1
}

func instO_PHP1(cpu *CPU) {
	data := cpu.pushFlags(true)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|0x0100, data)
	cpu.sp--
	cpu.next = instInit
}

func instO_PLP(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instO_PLP1
}

func instO_PLP1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.sp++
	cpu.next = instO_PLP2
}

func instO_PLP2(cpu *CPU) {
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
func instO_JMP(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instO_JMP1
}

func instO_JMP1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc = (uint16(data) << 8) | cpu.ar
	cpu.next = instInit
}

func instO_JMP_I(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(cpu.ar))
	cpu.next = instO_JMP_I1
}

func instO_JMP_I1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(((cpu.ar + 1) & 0xff) | (cpu.ar & 0xff00))
	cpu.pc |= uint16(data) << 8
	cpu.next = instInit
}

func instO_JSR(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.next = instO_JSR1
}

func instO_JSR1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.next = instO_JSR2
}

func instO_JSR2(cpu *CPU) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = instO_JSR3
}

func instO_JSR3(cpu *CPU) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
	cpu.sp--
	cpu.next = instO_JSR4
}

func instO_JSR4(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.pc = cpu.ar | (uint16(data) << 8)
	cpu.next = instInit
}

func instO_RTS(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instO_RTS1
}

func instO_RTS1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.sp++
	cpu.next = instO_RTS2
}

func instO_RTS2(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(uint16(cpu.sp) | 0x100))
	cpu.sp++
	cpu.next = instO_RTS3
}

func instO_RTS3(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.pc |= uint16(data) << 8
	cpu.next = instO_RTS4
}

func instO_RTS4(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.next = instInit
}

func instO_RTI(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instO_RTI1
}

func instO_RTI1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.sp++
	cpu.next = instO_RTI2
}

func instO_RTI2(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(uint16(cpu.sp) | 0x0100)
	cpu.popFlags(data)
	cpu.sp++
	cpu.next = instO_RTI3
}

func instO_RTI3(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(uint16(cpu.sp) | 0x100))
	cpu.sp++
	cpu.next = instO_RTI4
}

func instO_RTI4(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.pc |= uint16(data) << 8
	cpu.next = instInit
}

func instO_BRK(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.next = instO_BRK1
}

func instO_BRK1(cpu *CPU) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = instO_BRK2
}

func instO_BRK2(cpu *CPU) {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
	cpu.sp--
	cpu.next = instO_BRK3
}

func instO_BRK3(cpu *CPU) {
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

func instO_BRK4(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(0xfffe))
	cpu.next = instO_BRK5
}

func instO_BRK5(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(0xffff)
	cpu.pc |= uint16(data) << 8
	cpu.next = instInit
}

func instO_BCS(cpu *CPU) {
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

func instO_BCC(cpu *CPU) {
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

func instO_BEQ(cpu *CPU) {
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

func instO_BNE(cpu *CPU) {
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

func instO_BVS(cpu *CPU) {
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

func instO_BVC(cpu *CPU) {
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

func instO_BMI(cpu *CPU) {
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

func instO_BPL(cpu *CPU) {
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

func instO_BRANCH_NP(cpu *CPU) {
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

func instO_BRANCH_BP(cpu *CPU) {
	// Page crossed, branch backwards
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc = cpu.ar
	cpu.next = instO_BRANCH_BP1
}

func instO_BRANCH_BP1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc + 0x100)
	cpu.next = instInit
}

func instO_BRANCH_FP(cpu *CPU) {
	// Page crossed, branch forwards
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc = cpu.ar
	cpu.next = instO_BRANCH_FP1
}

func instO_BRANCH_FP1(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc - 0x100)
	cpu.next = instInit
}

// Flag group
func instO_SEC(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.cFlag = 1
	cpu.next = instInit
}

func instO_CLC(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.cFlag = 0
	cpu.next = instInit
}

func instO_SED(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.dFlag = 1
	cpu.next = instInit
}

func instO_CLD(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.dFlag = 0
	cpu.next = instInit
}

func instO_SEI(cpu *CPU) {
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

func instO_CLI(cpu *CPU) {
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

func instO_CLV(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.vFlag = 0
	cpu.next = instInit
}

// NOP group
func instO_NOP(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.next = instInit
}

// Undocumented functions start here

// NOP group
func instO_NOP_I(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.next = instInit
}

func instO_NOP_A(cpu *CPU) {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.next = instInit
}

// Load A/X group
func instO_LAX(cpu *CPU) {
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
func instO_SAX(cpu *CPU) {
	cpu.banks.Write(cpu.ar, cpu.a&cpu.x)
	cpu.next = instInit
}

// ASL/ORA group
func instO_SLO(cpu *CPU) {
	cpu.cFlag = cpu.rmw & 0x80
	cpu.rmw <<= 1
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.a |= cpu.rmw
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

// ROL/AND group
func instO_RLA(cpu *CPU) {
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
func instO_SRE(cpu *CPU) {
	cpu.cFlag = cpu.rmw & 0x01
	cpu.rmw >>= 1
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.a ^= cpu.rmw
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instInit
}

// ROR/ADC group
func instO_RRA(cpu *CPU) {
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
func instO_DCP(cpu *CPU) {
	cpu.rmw--
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.ar = uint16(cpu.a) - uint16(cpu.rmw)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
	cpu.next = instInit
}

// INC/SBC group
func instO_ISB(cpu *CPU) {
	cpu.rmw++
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.doSBC(cpu.rmw)
	cpu.next = instInit
}

// Complex functions
func instO_ANC_I(cpu *CPU) {
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

func instO_ASR_I(cpu *CPU) {
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

func instO_ARR_I(cpu *CPU) {
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

func instO_ANE_I(cpu *CPU) {
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

func instO_LXA_I(cpu *CPU) {
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

func instO_SBX_I(cpu *CPU) {
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

func instO_LAS(cpu *CPU) {
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

func instO_SHS(cpu *CPU) {
	// ar2 contains the high byte of the operand address
	cpu.sp = cpu.a & cpu.x
	cpu.banks.Write(cpu.ar, uint8((cpu.ar2+1)&uint16(cpu.sp)))
	cpu.next = instInit
}

func instO_SHY(cpu *CPU) {
	// ar2 contains the high byte of the operand address
	cpu.banks.Write(cpu.ar, uint8(uint16(cpu.y)&(cpu.ar2+1)))
	cpu.next = instInit
}

func instO_SHX(cpu *CPU) {
	// ar2 contains the high byte of the operand address
	cpu.banks.Write(cpu.ar, uint8(uint16(cpu.x)&(cpu.ar2+1)))
	cpu.next = instInit
}

func instO_SHA(cpu *CPU) {
	// ar2 contains the high byte of the operand address
	cpu.banks.Write(cpu.ar, uint8(uint16(cpu.a)&uint16(cpu.x)&(cpu.ar2+1)))
	cpu.next = instInit
}

func instI_ILL_OP(cpu *CPU) {
	cpu.illegalOp(cpu.op, cpu.pc-1)
}
