package mos6510fn

import (
	"fmt"
	"github.com/markel1974/c64emu/src/flag"
	"log"
	"os"
)

type Core struct {
	banks          IBanks
	pic            IPic
	nFlag          uint8  // Negative flag
	zFlag          uint8  // Zero flag
	vFlag          uint8  // Overflow flag
	dFlag          uint8  // Decimal mode flag
	iFlag          uint8  // Interrupt disable flag
	cFlag          uint8  // Carry flag
	a              uint8  // Register
	x              uint8  // Register
	y              uint8  // Register
	sp             uint8  // Stack pointer
	pc             uint16 // Program counter
	op             uint8  // Current opcode
	ar             uint16 // Address register
	ar2            uint16 // Address register 2
	rmw            uint8  // Data buffer for RMW instructions
	opFlags        uint8  //
	stop           bool   //
	next           func(*Core)
	rdyLow         bool // current RDY state
	aecLow         bool // current AEC state
	overflowBranch func() bool
}

func NewCore(pic IPic) *Core {
	regs := &Core{
		banks:          nil,
		pic:            pic,
		a:              0,
		x:              0,
		y:              0,
		sp:             0xff,
		nFlag:          0,
		zFlag:          0,
		vFlag:          0,
		dFlag:          0,
		cFlag:          0,
		iFlag:          1,
		opFlags:        0,
		ar:             0,
		ar2:            0,
		rdyLow:         false,
		aecLow:         false,
		stop:           false,
		overflowBranch: nil,
		next:           nil,
	}
	return regs
}

func (cpu *Core) reset() {
	// Read reset vector
	cpu.pc = uint16(cpu.banks.Read(0xfffc)) | (uint16(cpu.banks.Read(0xfffd)) << 8)
	cpu.next = fnInit
	cpu.opFlags = 0
}

func (cpu *Core) popFlags(data uint8) {
	cpu.nFlag = data
	cpu.vFlag = data & 0x40
	cpu.dFlag = data & 0x08
	cpu.iFlag = data & 0x04
	cpu.zFlag = flag.BoolToUint8((data & 0x02) == 0)
	cpu.cFlag = data & 0x01
}

func (cpu *Core) pushFlags(bFlags bool) uint8 {
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

func (cpu *Core) branch(data uint8) {
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

func (cpu *Core) doADC(data uint8) {
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

func (cpu *Core) doSBC(data uint8) {
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

func (cpu *Core) illegalOp(illOp uint8, at uint16) {
	log.Printf("illegal opcode %02x at %04x.", illOp, at)
	//TODO EVENT
	cpu.reset()
	os.Exit(1)
}

func (cpu *Core) printRegisters(qCycle uint64, baLow bool) {
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
