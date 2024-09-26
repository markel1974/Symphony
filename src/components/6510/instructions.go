package mos6510

import (
	"github.com/markel1974/c64emu/src/conversion"
	"log"
	"os"
)

const (
	stackAddr = 0x100
)

func instOpINI(cpu *CPU) {
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
				cpu.next = instOpNMI
				cpu.next(cpu)
				return
			case 3:
				cpu.irqBreaker = true
				cpu.next = instOpIRQ
				cpu.next(cpu)
				return
			}
		}
	} else {
		cpu.stop = true
		return
	}
	cpu.irqBreaker = false
	cpu.op = cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.next = _modeTable[cpu.op]
}

func instOpIRQ(cpu *CPU) {
	//internal operation
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpIRQ1
}

func instOpIRQ1(cpu *CPU) {
	//internal operation
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpIRQ2
}

func instOpIRQ2(cpu *CPU) {
	//push return address high byte onto stack
	cpu.banks.Write(uint16(cpu.sp)|stackAddr, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = instOpIRQ3
}

func instOpIRQ3(cpu *CPU) {
	//push return address low byte onto stack
	cpu.banks.Write(uint16(cpu.sp)|stackAddr, uint8(cpu.pc))
	cpu.sp--
	cpu.next = instOpIRQ4
}

func instOpIRQ4(cpu *CPU) {
	//push status register onto stack
	data := cpu.pushFlags(false)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|stackAddr, data)
	cpu.sp--
	cpu.iFlag = 1
	cpu.next = instOpIRQ5
}

func instOpIRQ5(cpu *CPU) {
	//get irq vector from 0xfffe
	data, ok := cpu.read(0xfffe)
	if !ok {
		return
	}
	cpu.pc = uint16(data)
	cpu.next = instOpIRQ6
}

func instOpIRQ6(cpu *CPU) {
	//get irq vector from 0xffff
	data, ok := cpu.read(0xffff)
	if !ok {
		return
	}
	cpu.pc |= uint16(data) << 8
	cpu.next = instOpINI
}

func instOpNMI(cpu *CPU) {
	//internal operation
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpNMI1
}

func instOpNMI1(cpu *CPU) {
	//internal operation
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpNMI2
}

func instOpNMI2(cpu *CPU) {
	//push return address high byte onto stack
	cpu.banks.Write(uint16(cpu.sp)|stackAddr, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = instOpNMI3
}

func instOpNMI3(cpu *CPU) {
	//push return address low byte onto stack
	cpu.banks.Write(uint16(cpu.sp)|stackAddr, uint8(cpu.pc))
	cpu.sp--
	cpu.next = instOpNMI4
}

func instOpNMI4(cpu *CPU) {
	//push status register onto stack
	data := cpu.pushFlags(false)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|stackAddr, data)
	cpu.sp--
	cpu.iFlag = 1
	cpu.next = instOpNMI5
}

func instOpNMI5(cpu *CPU) {
	//get irq vector from 0xfffa
	data, ok := cpu.read(0xfffa)
	if !ok {
		return
	}
	cpu.pc = uint16(data)
	cpu.next = instOpNMI6
}

func instOpNMI6(cpu *CPU) {
	//get irq vector from 0xfffb
	data, ok := cpu.read(0xfffb)
	if !ok {
		return
	}
	cpu.pc |= uint16(data) << 8
	cpu.next = instOpINI
}

// Addressing modes: Fetch effective address (no extra cycles)
func instApZER(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = _opTable[cpu.op]
}

func instApZERx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instApZERx1
}

func instApZERx1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
	cpu.next = _opTable[cpu.op]
}

func instApZERy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instApZERy1
}

func instApZERy1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
	cpu.next = _opTable[cpu.op]
}

func instApABS(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instApABS1
}

func instApABS1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = _opTable[cpu.op]
}

func instApABSx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instApABSx1
}

func instApABSx1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	z := cpu.ar + uint16(cpu.x)
	cpu.ar = (z & 0xff) | (cpu.ar2 << 8)
	if z < stackAddr {
		cpu.next = instApABSx2
	} else {
		cpu.next = instApABSx3
	}
}

func instApABSx2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = _opTable[cpu.op]
}

func instApABSx3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = _opTable[cpu.op]
}

func instApABSy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instApABSy1
}

func instApABSy1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (cpu.ar2 << 8)
	if z < stackAddr {
		cpu.next = instApABSy2
	} else {
		cpu.next = instApABSy3
	}
}

func instApABSy2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = _opTable[cpu.op]
}

func instApABSy3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = _opTable[cpu.op]
}

func instApINDx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = instApINDx1
}

func instApINDx1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar2); !ok {
		return
	}
	cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
	cpu.next = instApINDx2
}

func instApINDx2(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = instApINDx3
}

func instApINDx3(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = _opTable[cpu.op]
}

func instApINDy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = instApINDy1
}

func instApINDy1(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = instApINDy2
}

func instApINDy2(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	cpu.ar2 = uint16(data)
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (cpu.ar2 << 8)
	if z < stackAddr {
		cpu.next = instApINDy3
	} else {
		cpu.next = instApINDy4
	}
}

func instApINDy3(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = _opTable[cpu.op]
}

func instApINDy4(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = _opTable[cpu.op]
}

func instAeABSx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instAeABSx1
}

func instAeABSx1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.x)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = _opTable[cpu.op]
	} else {
		cpu.next = instAeABSx2
	}
}

func instAeABSx2(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = _opTable[cpu.op]
}

func instAeABSy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instAeABSy1
}

func instAeABSy1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = _opTable[cpu.op]
	} else {
		cpu.next = instAeABSy2
	}
}

func instAeABSy2(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = _opTable[cpu.op]
}

func instAeINDy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = instAeINDy1
}

func instAeINDy1(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = instAeINDy2
}

func instAeINDy2(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = _opTable[cpu.op]
	} else {
		cpu.next = instAeINDy3
	}
}

func instAeINDy3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = _opTable[cpu.op]
}

func instMpZER(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instOpRMW
}

func instMpZERx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instMpZERx1
}

func instMpZERx1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
	cpu.next = instOpRMW
}

//func instMpZERy(cpu *CPU) {
//	data, ok := cpu.read(cpu.pc)
//	if !ok {
//		return
//	}
//	cpu.pc++
//	cpu.ar = uint16(data)
//	cpu.next = instMpZERy1
//}

//func instMpZERy1(cpu *CPU) {
//	data, ok := cpu.read(cpu.ar)
//	if !ok {
//		return
//	}
//	cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
//	cpu.next = instOpRMW
//}

func instMpABS(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instMpABS1
}

func instMpABS1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = instOpRMW
}

func instMpABSx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instMpABSx1
}

func instMpABSx1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.x)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = instMpABSx2
	} else {
		cpu.next = instMpABSx3
	}
}

func instMpABSx2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = instOpRMW
}

func instMpABSx3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = instOpRMW
}

func instMpABSy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instMpABSy1
}

func instMpABSy1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = instMpABSy2
	} else {
		cpu.next = instMpABSy3
	}
}

func instMpABSy2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = instOpRMW
}

func instMpABSy3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = instOpRMW
}

func instMpINDx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = instMpINDx1
}

func instMpINDx1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar2); !ok {
		return
	}
	cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
	cpu.next = instMpINDx2
}

func instMpINDx2(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = instMpINDx3
}

func instMpINDx3(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = instOpRMW
}

func instMpINDy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = instMpINDy1
}

func instMpINDy1(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = instMpINDy2
}

func instMpINDy2(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = instMpINDy3
	} else {
		cpu.next = instMpINDy4
	}
}

func instMpINDy3(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = instOpRMW
}

func instMpINDy4(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = instOpRMW
}

func instOpRMW(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.rmw = data
	cpu.next = instOpRMW1
}

func instOpRMW1(cpu *CPU) {
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.next = _opTable[cpu.op]
}

// Load group
func instOpLDA(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.a = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instOpINI
}

func instOiLDA(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instOpINI
}

func instOpLDX(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.x = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instOpINI
}

func instOiLDX(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.x = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instOpINI
}

func instOpLDY(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.y = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instOpINI
}

func instOiLDY(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.y = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instOpINI
}

// Store

func instOpSTA(cpu *CPU) {
	cpu.banks.Write(cpu.ar, cpu.a)
	cpu.next = instOpINI
}

func instOpSTX(cpu *CPU) {
	cpu.banks.Write(cpu.ar, cpu.x)
	cpu.next = instOpINI
}

func instOpSTY(cpu *CPU) {
	cpu.banks.Write(cpu.ar, cpu.y)
	cpu.next = instOpINI
}

// Transfer

func instOpTAX(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.x = cpu.a
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

func instOpTXA(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.a = cpu.x
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.next = instOpINI
}

func instOpTAY(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.y = cpu.a
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

func instOpTYA(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.a = cpu.y
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.next = instOpINI
}

func instOpTSX(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.x = cpu.sp
	cpu.nFlag = cpu.sp
	cpu.zFlag = cpu.sp
	cpu.next = instOpINI
}

func instOpTXS(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.sp = cpu.x
	cpu.next = instOpINI
}

// Arithmetic

func instOpADC(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.doADC(data)
	cpu.next = instOpINI
}

func instOiADC(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.doADC(data)
	cpu.next = instOpINI
}

func instOpSBC(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.doSBC(data)
	cpu.next = instOpINI
}

func instOiSBC(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.doSBC(data)
	cpu.next = instOpINI
}

// Increment, decrement

func instOpINX(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.x++
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.next = instOpINI
}

func instOpDEX(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.x--
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.next = instOpINI
}

func instOpINY(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.y++
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.next = instOpINI
}

func instOpDEY(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.y--
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.next = instOpINI
}

func instOpINC(cpu *CPU) {
	v := cpu.rmw + 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instOpINI
}

func instOpDEC(cpu *CPU) {
	v := cpu.rmw - 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instOpINI
}

// Logic group
func instOpAND(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.a &= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

func instOiAND(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a &= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

func instOpORA(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.a |= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

func instOiOPA(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a |= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

func instOpEOR(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.a ^= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

func instOiEOR(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a ^= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// Compare group
func instOpCMP(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.ar = uint16(cpu.a) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

func instOiCMP(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(cpu.a) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

func instOpCPX(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.ar = uint16(cpu.x) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

func instOiCPX(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(cpu.x) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

func instOpCPY(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.ar = uint16(cpu.y) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

func instOiCPY(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(cpu.y) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

// Bit-test

func instOpBIT(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.zFlag = cpu.a & data
	cpu.nFlag = data
	cpu.vFlag = data & 0x40
	cpu.next = instOpINI
}

// Shift/rotate group
func instOpASL(cpu *CPU) {
	cpu.cFlag = cpu.rmw & 0x80
	v := cpu.rmw << 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instOpINI
}

func instOaASL(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.cFlag = cpu.a & 0x80
	cpu.a <<= 1
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

func instOpLSR(cpu *CPU) {
	cpu.cFlag = cpu.rmw & 0x1
	v := cpu.rmw >> 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instOpINI
}

func instOaLSR(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.cFlag = cpu.a & 0x1
	cpu.a >>= 1
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

func instOpROL(cpu *CPU) {
	var t uint8
	if cpu.cFlag != 0 {
		t = (cpu.rmw << 1) | 0x1
	} else {
		t = cpu.rmw << 1
	}
	cpu.nFlag = t
	cpu.zFlag = t
	cpu.banks.Write(cpu.ar, t)
	cpu.cFlag = cpu.rmw & 0x80
	cpu.next = instOpINI
}

func instOaROL(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	data := cpu.a & 0x80
	if cpu.cFlag != 0 {
		cpu.a = (cpu.a << 1) | 0x1
	} else {
		cpu.a = cpu.a << 1
	}
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.cFlag = data
	cpu.next = instOpINI
}

func instOpROR(cpu *CPU) {
	var t uint8
	if cpu.cFlag != 0 {
		t = (cpu.rmw >> 1) | 0x80
	} else {
		t = cpu.rmw >> 1
	}
	cpu.nFlag = t
	cpu.zFlag = t
	cpu.banks.Write(cpu.ar, t)
	cpu.cFlag = cpu.rmw & 0x1
	cpu.next = instOpINI
}

func instOaROR(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	data := cpu.a & 0x1
	if cpu.cFlag != 0 {
		cpu.a = (cpu.a >> 1) | 0x80
	} else {
		cpu.a = cpu.a >> 1
	}
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.cFlag = data
	cpu.next = instOpINI
}

// Stack

func instOpPHA(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpPHA1
}

func instOpPHA1(cpu *CPU) {
	cpu.banks.Write(uint16(cpu.sp)|stackAddr, cpu.a)
	cpu.sp--
	cpu.next = instOpINI
}

func instOpPLA(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpPLA1
}

func instOpPLA1(cpu *CPU) {
	if _, ok := cpu.read(uint16(cpu.sp) | stackAddr); !ok {
		return
	}
	cpu.sp++
	cpu.next = instOpPLA2
}

func instOpPLA2(cpu *CPU) {
	data, ok := cpu.read(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	cpu.a = data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

func instOpPHP(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpPHP1
}

func instOpPHP1(cpu *CPU) {
	data := cpu.pushFlags(true)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|stackAddr, data)
	cpu.sp--
	cpu.next = instOpINI
}

func instOpPLP(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpPLP1
}

func instOpPLP1(cpu *CPU) {
	if _, ok := cpu.read(uint16(cpu.sp) | stackAddr); !ok {
		return
	}
	cpu.sp++
	cpu.next = instOpPLP2
}

func instOpPLP2(cpu *CPU) {
	data, ok := cpu.read(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	iFlagPrev := cpu.iFlag
	cpu.popFlags(data)
	if iFlagPrev == 0 && cpu.iFlag != 0 {
		cpu.opFlags |= opFlagIrqDisabled
	} else if iFlagPrev != 0 && cpu.iFlag == 0 {
		cpu.opFlags |= opFlagIrqEnabled
	}
	cpu.next = instOpINI
}

// Jump - Branch

func instOpJMP(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instOpJMP1
}

func instOpJMP1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc = (uint16(data) << 8) | cpu.ar
	cpu.next = instOpINI
}

func instOiJMP(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.pc = uint16(data)
	cpu.next = instOiJMP1
}

func instOiJMP1(cpu *CPU) {
	data, ok := cpu.read(((cpu.ar + 1) & 0xff) | (cpu.ar & 0xff00))
	if !ok {
		return
	}
	cpu.pc |= uint16(data) << 8
	cpu.next = instOpINI
}

func instOpJSR(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instOpJSR1
}

func instOpJSR1(cpu *CPU) {
	if _, ok := cpu.read(uint16(cpu.sp) | stackAddr); !ok {
		return
	}
	cpu.next = instOpJSR2
}

func instOpJSR2(cpu *CPU) {
	cpu.banks.Write(uint16(cpu.sp)|stackAddr, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = instOpJSR3
}

func instOpJSR3(cpu *CPU) {
	cpu.banks.Write(uint16(cpu.sp)|stackAddr, uint8(cpu.pc))
	cpu.sp--
	cpu.next = instOpJSR4
}

func instOpJSR4(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.pc = cpu.ar | (uint16(data) << 8)
	cpu.next = instOpINI
}

func instOpRTS(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpRTS1
}

func instOpRTS1(cpu *CPU) {
	if _, ok := cpu.read(uint16(cpu.sp) | stackAddr); !ok {
		return
	}
	cpu.sp++
	cpu.next = instOpRTS2
}

func instOpRTS2(cpu *CPU) {
	data, ok := cpu.read(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	cpu.pc = uint16(data)
	cpu.sp++
	cpu.next = instOpRTS3
}

func instOpRTS3(cpu *CPU) {
	data, ok := cpu.read(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	cpu.pc |= uint16(data) << 8
	cpu.next = instOpRTS4
}

func instOpRTS4(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.pc++
	cpu.next = instOpINI
}

func instOpRTI(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpRTI1
}

func instOpRTI1(cpu *CPU) {
	if _, ok := cpu.read(uint16(cpu.sp) | stackAddr); !ok {
		return
	}
	cpu.sp++
	cpu.next = instOpRTI2
}

func instOpRTI2(cpu *CPU) {
	data, ok := cpu.read(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	cpu.popFlags(data)
	cpu.sp++
	cpu.next = instOpRTI3
}

func instOpRTI3(cpu *CPU) {
	data, ok := cpu.read(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	cpu.pc = uint16(data)
	cpu.sp++
	cpu.next = instOpRTI4
}

func instOpRTI4(cpu *CPU) {
	data, ok := cpu.read(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	cpu.pc |= uint16(data) << 8
	cpu.next = instOpINI
}

func instOpBRK(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.pc++
	cpu.next = instOpBRK1
}

func instOpBRK1(cpu *CPU) {
	cpu.banks.Write(uint16(cpu.sp)|stackAddr, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = instOpBRK2
}

func instOpBRK2(cpu *CPU) {
	cpu.banks.Write(uint16(cpu.sp)|stackAddr, uint8(cpu.pc))
	cpu.sp--
	cpu.next = instOpBRK3
}

func instOpBRK3(cpu *CPU) {
	data := cpu.pushFlags(true)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|stackAddr, data)
	cpu.sp--
	cpu.iFlag = 1
	if cpu.pic.HasNMI() {
		cpu.pic.ClearNMI()    // Simulate an edge-triggered input
		cpu.next = instOpNMI5 // Jump to NMI sequence
	} else {
		cpu.next = instOpBRK4
	}
}

func instOpBRK4(cpu *CPU) {
	data, ok := cpu.read(0xfffe)
	if !ok {
		return
	}
	cpu.pc = uint16(data)
	cpu.next = instOpBRK5
}

func instOpBRK5(cpu *CPU) {
	data, ok := cpu.read(0xffff)
	if !ok {
		return
	}
	cpu.pc |= uint16(data) << 8
	cpu.next = instOpINI
}

func instOpBCS(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if cpu.cFlag == 0 {
		cpu.next = instOpINI
	} else {
		cpu.branch(data)
	}
}

func instOpBCC(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if cpu.cFlag != 0 {
		cpu.next = instOpINI
	} else {
		cpu.branch(data)
	}
}

func instOpBEQ(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if cpu.zFlag != 0 {
		cpu.next = instOpINI
	} else {
		cpu.branch(data)
	}
}

func instOpBNE(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if cpu.zFlag == 0 {
		cpu.next = instOpINI
	} else {
		cpu.branch(data)
	}
}

func instOpBVS(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if cpu.overflowBranch != nil && cpu.overflowBranch() {
		cpu.vFlag = 1
	}
	if cpu.vFlag == 0 {
		cpu.next = instOpINI
	} else {
		cpu.branch(data)
	}
}

func instOpBVC(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if cpu.overflowBranch != nil && cpu.overflowBranch() {
		cpu.vFlag = 1
	}
	if cpu.vFlag != 0 {
		cpu.next = instOpINI
	} else {
		cpu.branch(data)
	}
}

func instOpBMI(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if (cpu.nFlag & 0x80) == 0 {
		cpu.next = instOpINI
	} else {
		cpu.branch(data)
	}
}

func instOpBPL(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if (cpu.nFlag & 0x80) != 0 {
		cpu.next = instOpINI
	} else {
		cpu.branch(data)
	}
}

func instOpBRAnp(cpu *CPU) {
	// No page crossed
	cpu.opFlags |= opFlagIntDelayed
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.pc = cpu.ar
	cpu.next = instOpINI
}

func instOpBRAbp(cpu *CPU) {
	// Page crossed (branch backwards)
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.pc = cpu.ar
	cpu.next = instOpBRAbp1
}

func instOpBRAbp1(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc + stackAddr); !ok {
		return
	}
	cpu.next = instOpINI
}

func instOpBRAfp(cpu *CPU) {
	// Page crossed (branch forwards)
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.pc = cpu.ar
	cpu.next = instOpBRAfp1
}

func instOpBRAfp1(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc - stackAddr); !ok {
		return
	}
	cpu.next = instOpINI
}

// Flag

func instOpSEC(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.cFlag = 1
	cpu.next = instOpINI
}

func instOpCLC(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.cFlag = 0
	cpu.next = instOpINI
}

func instOpSED(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.dFlag = 1
	cpu.next = instOpINI
}

func instOpCLD(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.dFlag = 0
	cpu.next = instOpINI
}

func instOpSEI(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	if cpu.iFlag == 0 {
		cpu.opFlags |= opFlagIrqDisabled
	}
	cpu.iFlag = 1
	cpu.next = instOpINI
}

func instOpCLI(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	if cpu.iFlag == 0 {
		cpu.opFlags |= opFlagIrqEnabled
	}
	cpu.iFlag = 0
	cpu.next = instOpINI
}

func instOpCLV(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.vFlag = 0
	cpu.next = instOpINI
}

// NOP
func instOpNOP(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpINI
}

// Undocumented functions

// NOP

func instOiNOP(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.pc++
	cpu.next = instOpINI
}

func instOaNOP(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = instOpINI
}

// Load A/X

func instOpLAX(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.x = data
	cpu.a = cpu.x
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// Store A/X

func instOpSAX(cpu *CPU) {
	cpu.banks.Write(cpu.ar, cpu.a&cpu.x)
	cpu.next = instOpINI
}

// ASL/ORA

func instOpSLO(cpu *CPU) {
	cpu.cFlag = cpu.rmw & 0x80
	cpu.rmw <<= 1
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.a |= cpu.rmw
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// ROL/AND

func instOpRLA(cpu *CPU) {
	tmp := cpu.rmw & 0x80
	if cpu.cFlag != 0 {
		cpu.rmw = (cpu.rmw << 1) | 0x1
	} else {
		cpu.rmw = cpu.rmw << 1
	}
	cpu.cFlag = tmp
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.a &= cpu.rmw
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// LSR/EOR

func instOpSRE(cpu *CPU) {
	cpu.cFlag = cpu.rmw & 0x1
	cpu.rmw >>= 1
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.a ^= cpu.rmw
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// ROR/ADC

func instOpRRA(cpu *CPU) {
	tmp := cpu.rmw & 0x1
	if cpu.cFlag != 0 {
		cpu.rmw = (cpu.rmw >> 1) | 0x80
	} else {
		cpu.rmw = cpu.rmw >> 1
	}
	cpu.cFlag = tmp
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.doADC(cpu.rmw)
	cpu.next = instOpINI
}

// DEC/CMP

func instOpDCP(cpu *CPU) {
	cpu.rmw--
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.ar = uint16(cpu.a) - uint16(cpu.rmw)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

// INC/SBC group
func instOpISB(cpu *CPU) {
	cpu.rmw++
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.doSBC(cpu.rmw)
	cpu.next = instOpINI
}

// Complex functions
func instOiANC(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a &= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.cFlag = cpu.nFlag & 0x80
	cpu.next = instOpINI
}

func instOiASR(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a &= data
	cpu.cFlag = cpu.a & 0x1
	cpu.a >>= 1
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

func instOiARR(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
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
		if ((data & 0xf) + (data & 0x1)) > 5 {
			cpu.a = (cpu.a & 0xf0) | ((cpu.a + 6) & 0xf)
		}
		k := uint16((data)+(uint8(data)&0x10)) & 0x1f0
		cpu.cFlag = uint8(k)
		if k > 0x50 {
			cpu.a += 0x60
		}
	}
	cpu.next = instOpINI
}

func instOiANE(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a = (cpu.a | 0xee) & cpu.x & data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

func instOiLXA(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.x = (cpu.a | 0xee) & data
	cpu.a = cpu.x
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

func instOiSBX(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = (uint16(cpu.x) & uint16(cpu.a)) - uint16(data)
	cpu.x = uint8(cpu.ar)
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

func instOpLAS(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.sp = data & cpu.sp
	cpu.x = cpu.sp
	cpu.a = cpu.x
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

func instOpSHS(cpu *CPU) {
	cpu.sp = cpu.a & cpu.x
	d := uint8((cpu.ar2 + 1) & uint16(cpu.sp))
	cpu.banks.Write(cpu.ar, d)
	cpu.next = instOpINI
}

func instOpSHY(cpu *CPU) {
	d := uint8(uint16(cpu.y) & (cpu.ar2 + 1))
	cpu.banks.Write(cpu.ar, d)
	cpu.next = instOpINI
}

func instOpSHX(cpu *CPU) {
	d := uint8(uint16(cpu.x) & (cpu.ar2 + 1))
	cpu.banks.Write(cpu.ar, d)
	cpu.next = instOpINI
}

func instOpSHA(cpu *CPU) {
	d := uint8(uint16(cpu.a) & uint16(cpu.x) & (cpu.ar2 + 1))
	cpu.banks.Write(cpu.ar, d)
	cpu.next = instOpINI
}

func instOpJAM(cpu *CPU) {
	log.Printf("[%s] illegal opcode %02x at %04x.", cpu.id, cpu.op, cpu.pc-1)
	//TODO EVENT
	cpu.Reset()
	os.Exit(1)
}
