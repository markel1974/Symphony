package mos6510

import (
	"fmt"
	"github.com/markel1974/c64emu/src/flag"
	"log"
	"os"
)

//https://dustlayer.com/c64-architecture
//https://www.zimmers.net/cbmpics/cbm/c64/vic-ii.txt

//Notes() {
//https://codebase64.org/lib/exe/fetch.php?media=base:safely_freezing_the_c64.pdf
/*
 *  - The zFlag variable has the inverse meaning of the 6510 Z flag
 *  - Only the highest bit of the nFlag variable is used
 */

type MOS6510fn struct {
	*Core
	id             string
	overflowBranch func() bool
	functions      []func()
}

func NewMOS6510fn(id string) *MOS6510fn {
	cpu := &MOS6510fn{
		Core:           nil,
		id:             id,
		overflowBranch: nil,
		functions:      make([]func(), 0xff),
	}
	for x := range cpu.functions {
		cpu.functions[x] = cpu.fnIllegalOp
	}
	return cpu
}

func (cpu *MOS6510fn) Setup(intr IPic, banks IBanks) {
	cpu.Core = NewCore(intr)
	cpu.banks = banks

	cpu.functions[STATE_LAST] = cpu.fnSTATE_LAST
	cpu.functions[I_IRQ_8] = cpu.fnI_IRQ_8
	cpu.functions[I_IRQ_9] = cpu.fnI_IRQ_9
	cpu.functions[I_IRQ_A] = cpu.fnI_IRQ_A
	cpu.functions[I_IRQ_B] = cpu.fnI_IRQ_B
	cpu.functions[I_IRQ_C] = cpu.fnI_IRQ_C
	cpu.functions[I_IRQ_D] = cpu.fnI_IRQ_D
	cpu.functions[I_IRQ_E] = cpu.fnI_IRQ_E
	cpu.functions[I_NMI_10] = cpu.fnI_NMI_10
	cpu.functions[I_NMI_11] = cpu.fnI_NMI_11
	cpu.functions[I_NMI_12] = cpu.fnI_NMI_12
	cpu.functions[I_NMI_13] = cpu.fnI_NMI_13
	cpu.functions[I_NMI_14] = cpu.fnI_NMI_14
	cpu.functions[I_NMI_15] = cpu.fnI_NMI_15
	cpu.functions[I_NMI_16] = cpu.fnI_NMI_16
	cpu.functions[A_ZERO] = cpu.fnA_ZERO
	cpu.functions[A_ZEROX] = cpu.fnA_ZEROX
	cpu.functions[A_ZEROX1] = cpu.fnA_ZEROX1
	cpu.functions[A_ZEROY] = cpu.fnA_ZEROY
	cpu.functions[A_ZEROY1] = cpu.fnA_ZEROY1
	cpu.functions[A_ABS] = cpu.fnA_ABS
	cpu.functions[A_ABS1] = cpu.fnA_ABS1
	cpu.functions[A_ABSX] = cpu.fnA_ABSX
	cpu.functions[A_ABSX1] = cpu.fnA_ABSX1
	cpu.functions[A_ABSX2] = cpu.fnA_ABSX2
	cpu.functions[A_ABSX3] = cpu.fnA_ABSX3
	cpu.functions[A_ABSY] = cpu.fnA_ABSY
	cpu.functions[A_ABSY1] = cpu.fnA_ABSY1
	cpu.functions[A_ABSY2] = cpu.fnA_ABSY2
	cpu.functions[A_ABSY3] = cpu.fnA_ABSY3
	cpu.functions[A_INDX] = cpu.fnA_INDX
	cpu.functions[A_INDX1] = cpu.fnA_INDX1
	cpu.functions[A_INDX2] = cpu.fnA_INDX2
	cpu.functions[A_INDX3] = cpu.fnA_INDX3
	cpu.functions[A_INDY] = cpu.fnA_INDY
	cpu.functions[A_INDY1] = cpu.fnA_INDY1
	cpu.functions[A_INDY2] = cpu.fnA_INDY2
	cpu.functions[A_INDY3] = cpu.fnA_INDY3
	cpu.functions[A_INDY4] = cpu.fnA_INDY4
	cpu.functions[AE_ABSX] = cpu.fnAE_ABSX
	cpu.functions[AE_ABSX1] = cpu.fnAE_ABSX1
	cpu.functions[AE_ABSX2] = cpu.fnAE_ABSX2
	cpu.functions[AE_ABSY] = cpu.fnAE_ABSY
	cpu.functions[AE_ABSY1] = cpu.fnAE_ABSY1
	cpu.functions[AE_ABSY2] = cpu.fnAE_ABSY2
	cpu.functions[AE_INDY] = cpu.fnAE_INDY
	cpu.functions[AE_INDY1] = cpu.fnAE_INDY1
	cpu.functions[AE_INDY2] = cpu.fnAE_INDY2
	cpu.functions[AE_INDY3] = cpu.fnAE_INDY3
	cpu.functions[M_ZERO] = cpu.fnM_ZERO
	cpu.functions[M_ZEROX] = cpu.fnM_ZEROX
	cpu.functions[M_ZEROX1] = cpu.fnM_ZEROX1
	cpu.functions[M_ZEROY] = cpu.fnM_ZEROY
	cpu.functions[M_ZEROY1] = cpu.fnM_ZEROY1
	cpu.functions[M_ABS] = cpu.fnM_ABS
	cpu.functions[M_ABS1] = cpu.fnM_ABS1
	cpu.functions[M_ABSX] = cpu.fnM_ABSX
	cpu.functions[M_ABSX1] = cpu.fnM_ABSX1
	cpu.functions[M_ABSX2] = cpu.fnM_ABSX2
	cpu.functions[M_ABSX3] = cpu.fnM_ABSX3
	cpu.functions[M_ABSY] = cpu.fnM_ABSY
	cpu.functions[M_ABSY1] = cpu.fnM_ABSY1
	cpu.functions[M_ABSY2] = cpu.fnM_ABSY2
	cpu.functions[M_ABSY3] = cpu.fnM_ABSY3
	cpu.functions[M_INDX] = cpu.fnM_INDX
	cpu.functions[M_INDX1] = cpu.fnM_INDX1
	cpu.functions[M_INDX2] = cpu.fnM_INDX2
	cpu.functions[M_INDX3] = cpu.fnM_INDX3
	cpu.functions[M_INDY] = cpu.fnM_INDY
	cpu.functions[M_INDY1] = cpu.fnM_INDY1
	cpu.functions[M_INDY2] = cpu.fnM_INDY2
	cpu.functions[M_INDY3] = cpu.fnM_INDY3
	cpu.functions[M_INDY4] = cpu.fnM_INDY4
	cpu.functions[RMW_DO_IT] = cpu.fnRMW_DO_IT
	cpu.functions[RMW_DO_IT1] = cpu.fnRMW_DO_IT1
	cpu.functions[O_LDA] = cpu.fnO_LDA
	cpu.functions[O_LDA_I] = cpu.fnO_LDA_I
	cpu.functions[O_LDX] = cpu.fnO_LDX
	cpu.functions[O_LDX_I] = cpu.fnO_LDX_I
	cpu.functions[O_LDY] = cpu.fnO_LDY
	cpu.functions[O_LDY_I] = cpu.fnO_LDY_I
	cpu.functions[O_STA] = cpu.fnO_STA
	cpu.functions[O_STX] = cpu.fnO_STX
	cpu.functions[O_STY] = cpu.fnO_STY
	cpu.functions[O_TAX] = cpu.fnO_TAX
	cpu.functions[O_TXA] = cpu.fnO_TXA
	cpu.functions[O_TAY] = cpu.fnO_TAY
	cpu.functions[O_TYA] = cpu.fnO_TYA
	cpu.functions[O_TSX] = cpu.fnO_TSX
	cpu.functions[O_TXS] = cpu.fnO_TXS
	cpu.functions[O_ADC] = cpu.fnO_ADC
	cpu.functions[O_ADC_I] = cpu.fnO_ADC_I
	cpu.functions[O_SBC] = cpu.fnO_SBC
	cpu.functions[O_SBC_I] = cpu.fnO_SBC_I
	cpu.functions[O_INX] = cpu.fnO_INX
	cpu.functions[O_DEX] = cpu.fnO_DEX
	cpu.functions[O_INY] = cpu.fnO_INY
	cpu.functions[O_DEY] = cpu.fnO_DEY
	cpu.functions[O_INC] = cpu.fnO_INC
	cpu.functions[O_DEC] = cpu.fnO_DEC
	cpu.functions[O_AND] = cpu.fnO_AND
	cpu.functions[O_AND_I] = cpu.fnO_AND_I
	cpu.functions[O_ORA] = cpu.fnO_ORA
	cpu.functions[O_ORA_I] = cpu.fnO_ORA_I
	cpu.functions[O_EOR] = cpu.fnO_EOR
	cpu.functions[O_EOR_I] = cpu.fnO_EOR_I
	cpu.functions[O_CMP] = cpu.fnO_CMP
	cpu.functions[O_CMP_I] = cpu.fnO_CMP_I
	cpu.functions[O_CPX] = cpu.fnO_CPX
	cpu.functions[O_CPX_I] = cpu.fnO_CPX_I
	cpu.functions[O_CPY] = cpu.fnO_CPY
	cpu.functions[O_CPY_I] = cpu.fnO_CPY_I
	cpu.functions[O_BIT] = cpu.fnO_BIT
	cpu.functions[O_ASL] = cpu.fnO_ASL
	cpu.functions[O_ASL_A] = cpu.fnO_ASL_A
	cpu.functions[O_LSR] = cpu.fnO_LSR
	cpu.functions[O_LSR_A] = cpu.fnO_LSR_A
	cpu.functions[O_ROL] = cpu.fnO_ROL
	cpu.functions[O_ROL_A] = cpu.fnO_ROL_A
	cpu.functions[O_ROR] = cpu.fnO_ROR
	cpu.functions[O_ROR_A] = cpu.fnO_ROR_A
	cpu.functions[O_PHA] = cpu.fnO_PHA
	cpu.functions[O_PHA1] = cpu.fnO_PHA1
	cpu.functions[O_PLA] = cpu.fnO_PLA
	cpu.functions[O_PLA1] = cpu.fnO_PLA1
	cpu.functions[O_PLA2] = cpu.fnO_PLA2
	cpu.functions[O_PHP] = cpu.fnO_PHP
	cpu.functions[O_PHP1] = cpu.fnO_PHP1
	cpu.functions[O_PLP] = cpu.fnO_PLP
	cpu.functions[O_PLP1] = cpu.fnO_PLP1
	cpu.functions[O_PLP2] = cpu.fnO_PLP2
	cpu.functions[O_JMP] = cpu.fnO_JMP
	cpu.functions[O_JMP1] = cpu.fnO_JMP1
	cpu.functions[O_JMP_I] = cpu.fnO_JMP_I
	cpu.functions[O_JMP_I1] = cpu.fnO_JMP_I1
	cpu.functions[O_JSR] = cpu.fnO_JSR
	cpu.functions[O_JSR1] = cpu.fnO_JSR1
	cpu.functions[O_JSR2] = cpu.fnO_JSR2
	cpu.functions[O_JSR3] = cpu.fnO_JSR3
	cpu.functions[O_JSR4] = cpu.fnO_JSR4
	cpu.functions[O_RTS] = cpu.fnO_RTS
	cpu.functions[O_RTS1] = cpu.fnO_RTS1
	cpu.functions[O_RTS2] = cpu.fnO_RTS2
	cpu.functions[O_RTS3] = cpu.fnO_RTS3
	cpu.functions[O_RTS4] = cpu.fnO_RTS4
	cpu.functions[O_RTI] = cpu.fnO_RTI
	cpu.functions[O_RTI1] = cpu.fnO_RTI1
	cpu.functions[O_RTI2] = cpu.fnO_RTI2
	cpu.functions[O_RTI3] = cpu.fnO_RTI3
	cpu.functions[O_RTI4] = cpu.fnO_RTI4
	cpu.functions[O_BRK] = cpu.fnO_BRK
	cpu.functions[O_BRK1] = cpu.fnO_BRK1
	cpu.functions[O_BRK2] = cpu.fnO_BRK2
	cpu.functions[O_BRK3] = cpu.fnO_BRK3
	cpu.functions[O_BRK4] = cpu.fnO_BRK4
	cpu.functions[O_BRK5] = cpu.fnO_BRK5
	//cpu.functions[O_BRK5NMI] = cpu.fnO_BRK5NMI
	cpu.functions[O_BCS] = cpu.fnO_BCS
	cpu.functions[O_BCC] = cpu.fnO_BCC
	cpu.functions[O_BEQ] = cpu.fnO_BEQ
	cpu.functions[O_BNE] = cpu.fnO_BNE
	cpu.functions[O_BVS] = cpu.fnO_BVS
	cpu.functions[O_BVC] = cpu.fnO_BVC
	cpu.functions[O_BMI] = cpu.fnO_BMI
	cpu.functions[O_BPL] = cpu.fnO_BPL
	cpu.functions[O_BRANCH_NP] = cpu.fnO_BRANCH_NP
	cpu.functions[O_BRANCH_BP] = cpu.fnO_BRANCH_BP
	cpu.functions[O_BRANCH_BP1] = cpu.fnO_BRANCH_BP1
	cpu.functions[O_BRANCH_FP] = cpu.fnO_BRANCH_FP
	cpu.functions[O_BRANCH_FP1] = cpu.fnO_BRANCH_FP1
	cpu.functions[O_SEC] = cpu.fnO_SEC
	cpu.functions[O_CLC] = cpu.fnO_CLC
	cpu.functions[O_SED] = cpu.fnO_SED
	cpu.functions[O_CLD] = cpu.fnO_CLD
	cpu.functions[O_SEI] = cpu.fnO_SEI
	cpu.functions[O_CLI] = cpu.fnO_CLI
	cpu.functions[O_CLV] = cpu.fnO_CLV
	cpu.functions[O_NOP] = cpu.fnO_NOP
	cpu.functions[O_NOP_I] = cpu.fnO_NOP_I
	cpu.functions[O_NOP_A] = cpu.fnO_NOP_A
	cpu.functions[O_LAX] = cpu.fnO_LAX
	cpu.functions[O_SAX] = cpu.fnO_SAX
	cpu.functions[O_SLO] = cpu.fnO_SLO
	cpu.functions[O_RLA] = cpu.fnO_RLA
	cpu.functions[O_SRE] = cpu.fnO_SRE
	cpu.functions[O_RRA] = cpu.fnO_RRA
	cpu.functions[O_DCP] = cpu.fnO_DCP
	cpu.functions[O_ISB] = cpu.fnO_ISB
	cpu.functions[O_ANC_I] = cpu.fnO_ANC_I
	cpu.functions[O_ASR_I] = cpu.fnO_ASR_I
	cpu.functions[O_ARR_I] = cpu.fnO_ARR_I
	cpu.functions[O_ANE_I] = cpu.fnO_ANE_I
	cpu.functions[O_LXA_I] = cpu.fnO_LXA_I
	cpu.functions[O_SBX_I] = cpu.fnO_SBX_I
	cpu.functions[O_LAS] = cpu.fnO_LAS
	cpu.functions[O_SHS] = cpu.fnO_SHS
	cpu.functions[O_SHY] = cpu.fnO_SHY
	cpu.functions[O_SHX] = cpu.fnO_SHX
	cpu.functions[O_SHA] = cpu.fnO_SHA
}

func (cpu *MOS6510fn) Reset() {
	// Read reset vector
	cpu.pc = uint16(cpu.banks.Read(0xfffc)) | (uint16(cpu.banks.Read(0xfffd)) << 8)
	cpu.state = STATE_LAST
	cpu.opFlags = 0
}

// SetOverflowBranch implement 6502c SO (SOB) Pin
func (cpu *MOS6510fn) SetOverflowBranch(sob func() bool) {
	cpu.overflowBranch = sob
}

func (cpu *MOS6510fn) popFlags(data uint8) {
	cpu.nFlag = data
	cpu.vFlag = data & 0x40
	cpu.dFlag = data & 0x08
	cpu.iFlag = data & 0x04
	cpu.zFlag = flag.BoolToUint8((data & 0x02) == 0)
	cpu.cFlag = data & 0x01
}

func (cpu *MOS6510fn) pushFlags(bFlags bool) uint8 {
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

func (cpu *MOS6510fn) branch(data uint8) {
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

func (cpu *MOS6510fn) doADC(data uint8) {
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

func (cpu *MOS6510fn) doSBC(data uint8) {
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

func (cpu *MOS6510fn) illegalOp(illOp uint8, at uint16) {
	log.Printf("illegal opcode %02x at %04x.", illOp, at)
	//TODO EVENT
	cpu.Reset()
	os.Exit(1)
}

func (cpu *MOS6510fn) checkPic() {
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
	//Interrupts are only recognized if the RDY line is high
	if cpu.pic.HasIRQ() && !cpu.rdyLow {
		//if cpu.pic.HasIRQ() {
		if ((cpu.iFlag == 0) || ((cpu.opFlags & OpFlagIrqDisabled) != 0)) && ((cpu.opFlags & OpFlagIrqEnabled) == 0) {
			delay := 0
			if (cpu.opFlags & OpFlagIntDelayed) != 0 {
				delay = 1
			}
			if (cpu.pic.GetIrqCycleDistance(delay)) >= 2 {
				cpu.state = I_IRQ_8
				cpu.opFlags = 0
			}
		}
	}
}

func (cpu *MOS6510fn) SetAECLow(aecLow bool) {
	cpu.aecLow = aecLow
	if cpu.aecLow {
		cpu.stop = true
	}
}

func (cpu *MOS6510fn) SetRDYLow(rdyLow bool) {
	cpu.rdyLow = rdyLow
	if !cpu.rdyLow {
		cpu.stop = false
	}
}

func (cpu *MOS6510fn) Emulate() {
	if cpu.stop {
		return
	}
	if cpu.state == STATE_LAST {
		cpu.checkPic()
	}
	cpu.functions[cpu.state]()
}

func (cpu *MOS6510fn) fnSTATE_LAST() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.op = cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.state = _modeTable[cpu.op]
	cpu.opFlags = 0
}

// IRQ
func (cpu *MOS6510fn) fnI_IRQ_8() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.state = I_IRQ_9
}

func (cpu *MOS6510fn) fnI_IRQ_9() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.state = I_IRQ_A
}

func (cpu *MOS6510fn) fnI_IRQ_A() {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.state = I_IRQ_B
}

func (cpu *MOS6510fn) fnI_IRQ_B() {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
	cpu.sp--
	cpu.state = I_IRQ_C
}

func (cpu *MOS6510fn) fnI_IRQ_C() {
	data := cpu.pushFlags(false)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|0x0100, data)
	cpu.sp--
	cpu.iFlag = 1
	cpu.state = I_IRQ_D
}

func (cpu *MOS6510fn) fnI_IRQ_D() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(0xfffe))
	cpu.state = I_IRQ_E
}

func (cpu *MOS6510fn) fnI_IRQ_E() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(0xffff)
	cpu.pc |= uint16(data) << 8
	cpu.state = STATE_LAST
}

// NMI
func (cpu *MOS6510fn) fnI_NMI_10() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.state = I_NMI_11
}

func (cpu *MOS6510fn) fnI_NMI_11() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.state = I_NMI_12
}

func (cpu *MOS6510fn) fnI_NMI_12() {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.state = I_NMI_13
}

func (cpu *MOS6510fn) fnI_NMI_13() {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
	cpu.sp--
	cpu.state = I_NMI_14
}

func (cpu *MOS6510fn) fnI_NMI_14() {
	data := cpu.pushFlags(false)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|0x0100, data)
	cpu.sp--
	cpu.iFlag = 1
	cpu.state = I_NMI_15
}

func (cpu *MOS6510fn) fnI_NMI_15() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(0xfffa))
	cpu.state = I_NMI_16
}

func (cpu *MOS6510fn) fnI_NMI_16() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(0xfffb)
	cpu.pc |= uint16(data) << 8
	cpu.state = STATE_LAST
}

// Addressing modes: Fetch effective address, no extra cycles (-> ar)
func (cpu *MOS6510fn) fnA_ZERO() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = _opTable[cpu.op]
}

func (cpu *MOS6510fn) fnA_ZEROX() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = A_ZEROX1
}

func (cpu *MOS6510fn) fnA_ZEROX1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
	cpu.state = _opTable[cpu.op]
}

func (cpu *MOS6510fn) fnA_ZEROY() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = A_ZEROY1
}

func (cpu *MOS6510fn) fnA_ZEROY1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
	cpu.state = _opTable[cpu.op]
}

func (cpu *MOS6510fn) fnA_ABS() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = A_ABS1
}

func (cpu *MOS6510fn) fnA_ABS1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.state = _opTable[cpu.op]
}

func (cpu *MOS6510fn) fnA_ABSX() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = A_ABSX1
}

func (cpu *MOS6510fn) fnA_ABSX1() {
	// Note: Some undocumented functions rely on the value of ar2
	if cpu.rdyLow {
		cpu.stop = true
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
}

func (cpu *MOS6510fn) fnA_ABSX2() {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.state = _opTable[cpu.op]
}

func (cpu *MOS6510fn) fnA_ABSX3() {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.state = _opTable[cpu.op]
}

func (cpu *MOS6510fn) fnA_ABSY() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = A_ABSY1
}

func (cpu *MOS6510fn) fnA_ABSY1() {
	// Note: Some undocumented functions rely on the value of ar2
	if cpu.rdyLow {
		cpu.stop = true
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
}

func (cpu *MOS6510fn) fnA_ABSY2() {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.state = _opTable[cpu.op]
}

func (cpu *MOS6510fn) fnA_ABSY3() {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.state = _opTable[cpu.op]
}

func (cpu *MOS6510fn) fnA_INDX() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = A_INDX1
}

func (cpu *MOS6510fn) fnA_INDX1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar2)
	cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
	cpu.state = A_INDX2
}

func (cpu *MOS6510fn) fnA_INDX2() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.state = A_INDX3
}

func (cpu *MOS6510fn) fnA_INDX3() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read((cpu.ar2 + 1) & 0xff)
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.state = _opTable[cpu.op]
}

func (cpu *MOS6510fn) fnA_INDY() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = A_INDY1
}

func (cpu *MOS6510fn) fnA_INDY1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.state = A_INDY2
}

func (cpu *MOS6510fn) fnA_INDY2() {
	// Note: Some undocumented functions rely on the value of ar2
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read((cpu.ar2 + 1) & 0xff))
	if cpu.ar+uint16(cpu.y) < 0x100 {
		cpu.state = A_INDY3
	} else {
		cpu.state = A_INDY4
	}
	cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (cpu.ar2 << 8)
}

func (cpu *MOS6510fn) fnA_INDY3() {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.state = _opTable[cpu.op]
}

func (cpu *MOS6510fn) fnA_INDY4() {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.state = _opTable[cpu.op]
}

func (cpu *MOS6510fn) fnAE_ABSX() {
	// Addressing modes: Fetch effective address, extra cycle on page crossing (-> ar)
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = AE_ABSX1
}

func (cpu *MOS6510fn) fnAE_ABSX1() {
	if cpu.rdyLow {
		cpu.stop = true
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
}

func (cpu *MOS6510fn) fnAE_ABSX2() {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.state = _opTable[cpu.op]
}

func (cpu *MOS6510fn) fnAE_ABSY() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = AE_ABSY1
}

func (cpu *MOS6510fn) fnAE_ABSY1() {
	if cpu.rdyLow {
		cpu.stop = true
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
}

func (cpu *MOS6510fn) fnAE_ABSY2() {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.state = _opTable[cpu.op]
}

func (cpu *MOS6510fn) fnAE_INDY() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = AE_INDY1
}

func (cpu *MOS6510fn) fnAE_INDY1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.state = AE_INDY2
}

func (cpu *MOS6510fn) fnAE_INDY2() {
	if cpu.rdyLow {
		cpu.stop = true
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
}

func (cpu *MOS6510fn) fnAE_INDY3() {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.state = _opTable[cpu.op]
}

func (cpu *MOS6510fn) fnM_ZERO() {
	// Addressing modes: Read operand, write it back, no extra cycles (-> ar, rmw)
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = RMW_DO_IT
}

func (cpu *MOS6510fn) fnM_ZEROX() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = M_ZEROX1
}

func (cpu *MOS6510fn) fnM_ZEROX1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
	cpu.state = RMW_DO_IT
}

func (cpu *MOS6510fn) fnM_ZEROY() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = M_ZEROY1
}

func (cpu *MOS6510fn) fnM_ZEROY1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
	cpu.state = RMW_DO_IT
}

func (cpu *MOS6510fn) fnM_ABS() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = M_ABS1
}

func (cpu *MOS6510fn) fnM_ABS1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.state = RMW_DO_IT
}

func (cpu *MOS6510fn) fnM_ABSX() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = M_ABSX1
}

func (cpu *MOS6510fn) fnM_ABSX1() {
	if cpu.rdyLow {
		cpu.stop = true
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
}

func (cpu *MOS6510fn) fnM_ABSX2() {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.state = RMW_DO_IT
}

func (cpu *MOS6510fn) fnM_ABSX3() {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.state = RMW_DO_IT
}

func (cpu *MOS6510fn) fnM_ABSY() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = M_ABSY1
}

func (cpu *MOS6510fn) fnM_ABSY1() {
	if cpu.rdyLow {
		cpu.stop = true
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
}

func (cpu *MOS6510fn) fnM_ABSY2() {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.state = RMW_DO_IT
}

func (cpu *MOS6510fn) fnM_ABSY3() {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.state = RMW_DO_IT
}

func (cpu *MOS6510fn) fnM_INDX() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = M_INDX1
}

func (cpu *MOS6510fn) fnM_INDX1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar2)
	cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
	cpu.state = M_INDX2
}

func (cpu *MOS6510fn) fnM_INDX2() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.state = M_INDX3
}

func (cpu *MOS6510fn) fnM_INDX3() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read((cpu.ar2 + 1) & 0xff)
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.state = RMW_DO_IT
}

func (cpu *MOS6510fn) fnM_INDY() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar2 = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = M_INDY1
}

func (cpu *MOS6510fn) fnM_INDY1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.ar2))
	cpu.state = M_INDY2
}

func (cpu *MOS6510fn) fnM_INDY2() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read((cpu.ar2 + 1) & 0xff)
	if cpu.ar+uint16(cpu.y) < 0x100 {
		cpu.state = M_INDY3
	} else {
		cpu.state = M_INDY4
	}
	cpu.ar = ((cpu.ar + uint16(cpu.y)) & 0xff) | (uint16(data) << 8)
}

func (cpu *MOS6510fn) fnM_INDY3() {
	// No page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.state = RMW_DO_IT
}

func (cpu *MOS6510fn) fnM_INDY4() {
	// Page crossed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.ar += 0x100
	cpu.state = RMW_DO_IT
}

func (cpu *MOS6510fn) fnRMW_DO_IT() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.rmw = cpu.banks.Read(cpu.ar)
	cpu.state = RMW_DO_IT1
}

func (cpu *MOS6510fn) fnRMW_DO_IT1() {
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.state = _opTable[cpu.op]
}

// Load group
func (cpu *MOS6510fn) fnO_LDA() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.a = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_LDA_I() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.a = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_LDX() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.x = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_LDX_I() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.x = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_LDY() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.y = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_LDY_I() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.y = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.state = STATE_LAST
}

// Store group
func (cpu *MOS6510fn) fnO_STA() {
	cpu.banks.Write(cpu.ar, cpu.a)
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_STX() {
	cpu.banks.Write(cpu.ar, cpu.x)
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_STY() {
	cpu.banks.Write(cpu.ar, cpu.y)
	cpu.state = STATE_LAST
}

// Transfer group
func (cpu *MOS6510fn) fnO_TAX() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.x = cpu.a
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_TXA() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.a = cpu.x
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_TAY() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.y = cpu.a
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_TYA() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.a = cpu.y
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_TSX() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.x = cpu.sp
	cpu.nFlag = cpu.sp
	cpu.zFlag = cpu.sp
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_TXS() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.sp = cpu.x
	cpu.state = STATE_LAST
}

// Arithmetic group
func (cpu *MOS6510fn) fnO_ADC() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.doADC(data)
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_ADC_I() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.doADC(data)
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_SBC() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.doSBC(data)
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_SBC_I() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.doSBC(data)
	cpu.state = STATE_LAST
}

// Increment/decrement group
func (cpu *MOS6510fn) fnO_INX() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.x++
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_DEX() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.x--
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_INY() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.y++
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_DEY() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.y--
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_INC() {
	v := cpu.rmw + 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_DEC() {
	v := cpu.rmw - 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.state = STATE_LAST
}

// Logic group
func (cpu *MOS6510fn) fnO_AND() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a &= cpu.banks.Read(cpu.ar)
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_AND_I() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a &= cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_ORA() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a |= cpu.banks.Read(cpu.ar)
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_ORA_I() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a |= cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_EOR() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a ^= cpu.banks.Read(cpu.ar)
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_EOR_I() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a ^= cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.state = STATE_LAST
}

// Compare group
func (cpu *MOS6510fn) fnO_CMP() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.ar = uint16(cpu.a) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_CMP_I() {
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
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_CPX() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.ar = uint16(cpu.x) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_CPX_I() {
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
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_CPY() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.ar = uint16(cpu.y) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_CPY_I() {
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
	cpu.state = STATE_LAST
}

// Bit-test group
func (cpu *MOS6510fn) fnO_BIT() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.ar)
	cpu.zFlag = cpu.a & data
	cpu.nFlag = data
	cpu.vFlag = data & 0x40
	cpu.state = STATE_LAST
}

// Shift/rotate group
func (cpu *MOS6510fn) fnO_ASL() {
	cpu.cFlag = cpu.rmw & 0x80
	v := cpu.rmw << 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_ASL_A() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.cFlag = cpu.a & 0x80
	cpu.a <<= 1
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_LSR() {
	cpu.cFlag = cpu.rmw & 0x01
	v := cpu.rmw >> 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_LSR_A() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.cFlag = cpu.a & 0x01
	cpu.a >>= 1
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_ROL() {
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
}

func (cpu *MOS6510fn) fnO_ROL_A() {
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
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_ROR() {
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
}

func (cpu *MOS6510fn) fnO_ROR_A() {
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
	cpu.state = STATE_LAST
}

// Stack group
func (cpu *MOS6510fn) fnO_PHA() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.state = O_PHA1
}

func (cpu *MOS6510fn) fnO_PHA1() {
	cpu.banks.Write(uint16(cpu.sp)|0x100, cpu.a)
	cpu.sp--
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_PLA() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.state = O_PLA1
}

func (cpu *MOS6510fn) fnO_PLA1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.sp++
	cpu.state = O_PLA2
}

func (cpu *MOS6510fn) fnO_PLA2() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a = cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_PHP() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.state = O_PHP1
}

func (cpu *MOS6510fn) fnO_PHP1() {
	data := cpu.pushFlags(true)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|0x0100, data)
	cpu.sp--
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_PLP() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.state = O_PLP1
}

func (cpu *MOS6510fn) fnO_PLP1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.sp++
	cpu.state = O_PLP2
}

func (cpu *MOS6510fn) fnO_PLP2() {
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
	cpu.state = STATE_LAST
}

// Jump/branch group
func (cpu *MOS6510fn) fnO_JMP() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = O_JMP1
}

func (cpu *MOS6510fn) fnO_JMP1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc = (uint16(data) << 8) | cpu.ar
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_JMP_I() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(cpu.ar))
	cpu.state = O_JMP_I1
}

func (cpu *MOS6510fn) fnO_JMP_I1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(((cpu.ar + 1) & 0xff) | (cpu.ar & 0xff00))
	cpu.pc |= uint16(data) << 8
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_JSR() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.ar = uint16(cpu.banks.Read(cpu.pc))
	cpu.pc++
	cpu.state = O_JSR1
}

func (cpu *MOS6510fn) fnO_JSR1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.state = O_JSR2
}

func (cpu *MOS6510fn) fnO_JSR2() {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.state = O_JSR3
}

func (cpu *MOS6510fn) fnO_JSR3() {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
	cpu.sp--
	cpu.state = O_JSR4
}

func (cpu *MOS6510fn) fnO_JSR4() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.pc = cpu.ar | (uint16(data) << 8)
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_RTS() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.state = O_RTS1
}

func (cpu *MOS6510fn) fnO_RTS1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.sp++
	cpu.state = O_RTS2
}

func (cpu *MOS6510fn) fnO_RTS2() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(uint16(cpu.sp) | 0x100))
	cpu.sp++
	cpu.state = O_RTS3
}

func (cpu *MOS6510fn) fnO_RTS3() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.pc |= uint16(data) << 8
	cpu.state = O_RTS4
}

func (cpu *MOS6510fn) fnO_RTS4() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_RTI() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.state = O_RTI1
}

func (cpu *MOS6510fn) fnO_RTI1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.sp++
	cpu.state = O_RTI2
}

func (cpu *MOS6510fn) fnO_RTI2() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(uint16(cpu.sp) | 0x0100)
	cpu.popFlags(data)
	cpu.sp++
	cpu.state = O_RTI3
}

func (cpu *MOS6510fn) fnO_RTI3() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(uint16(cpu.sp) | 0x100))
	cpu.sp++
	cpu.state = O_RTI4
}

func (cpu *MOS6510fn) fnO_RTI4() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(uint16(cpu.sp) | 0x100)
	cpu.pc |= uint16(data) << 8
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_BRK() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.state = O_BRK1
}

func (cpu *MOS6510fn) fnO_BRK1() {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.state = O_BRK2
}

func (cpu *MOS6510fn) fnO_BRK2() {
	cpu.banks.Write(uint16(cpu.sp)|0x100, uint8(cpu.pc))
	cpu.sp--
	cpu.state = O_BRK3
}

func (cpu *MOS6510fn) fnO_BRK3() {
	data := cpu.pushFlags(true)
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
}

func (cpu *MOS6510fn) fnO_BRK4() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.pc = uint16(cpu.banks.Read(0xfffe))
	cpu.state = O_BRK5
}

func (cpu *MOS6510fn) fnO_BRK5() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(0xffff)
	cpu.pc |= uint16(data) << 8
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_BCS() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.cFlag == 0 {
		cpu.state = STATE_LAST
	} else {
		cpu.branch(data)
	}
}

func (cpu *MOS6510fn) fnO_BCC() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.cFlag != 0 {
		cpu.state = STATE_LAST
	} else {
		cpu.branch(data)
	}
}

func (cpu *MOS6510fn) fnO_BEQ() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.zFlag != 0 {
		cpu.state = STATE_LAST
	} else {
		cpu.branch(data)
	}
}

func (cpu *MOS6510fn) fnO_BNE() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if cpu.zFlag == 0 {
		cpu.state = STATE_LAST
	} else {
		cpu.branch(data)
	}
}

func (cpu *MOS6510fn) fnO_BVS() {
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
		cpu.state = STATE_LAST
	} else {
		cpu.branch(data)
	}
}

func (cpu *MOS6510fn) fnO_BVC() {
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
		cpu.state = STATE_LAST
	} else {
		cpu.branch(data)
	}
}

func (cpu *MOS6510fn) fnO_BMI() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if (cpu.nFlag & 0x80) == 0 {
		cpu.state = STATE_LAST
	} else {
		cpu.branch(data)
	}
}

func (cpu *MOS6510fn) fnO_BPL() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	if (cpu.nFlag & 0x80) != 0 {
		cpu.state = STATE_LAST
	} else {
		cpu.branch(data)
	}
}

func (cpu *MOS6510fn) fnO_BRANCH_NP() {
	// No page crossed
	cpu.opFlags |= OpFlagIntDelayed
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc = cpu.ar
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_BRANCH_BP() {
	// Page crossed, branch backwards
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc = cpu.ar
	cpu.state = O_BRANCH_BP1
}

func (cpu *MOS6510fn) fnO_BRANCH_BP1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc + 0x100)
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_BRANCH_FP() {
	// Page crossed, branch forwards
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc = cpu.ar
	cpu.state = O_BRANCH_FP1
}

func (cpu *MOS6510fn) fnO_BRANCH_FP1() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc - 0x100)
	cpu.state = STATE_LAST
}

// Flag group
func (cpu *MOS6510fn) fnO_SEC() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.cFlag = 1
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_CLC() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.cFlag = 0
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_SED() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.dFlag = 1
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_CLD() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.dFlag = 0
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_SEI() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	if cpu.iFlag == 0 {
		cpu.opFlags |= OpFlagIrqDisabled
	}
	cpu.iFlag = 1
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_CLI() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	if cpu.iFlag == 0 {
		cpu.opFlags |= OpFlagIrqEnabled
	}
	cpu.iFlag = 0
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_CLV() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.vFlag = 0
	cpu.state = STATE_LAST
}

// NOP group
func (cpu *MOS6510fn) fnO_NOP() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.state = STATE_LAST
}

// Undocumented functions start here

// NOP group
func (cpu *MOS6510fn) fnO_NOP_I() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_NOP_A() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.banks.Read(cpu.ar)
	cpu.state = STATE_LAST
}

// Load A/X group
func (cpu *MOS6510fn) fnO_LAX() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.x = cpu.banks.Read(cpu.ar)
	cpu.a = cpu.x
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.state = STATE_LAST
}

// Store A/X group
func (cpu *MOS6510fn) fnO_SAX() {
	cpu.banks.Write(cpu.ar, cpu.a&cpu.x)
	cpu.state = STATE_LAST
}

// ASL/ORA group
func (cpu *MOS6510fn) fnO_SLO() {
	cpu.cFlag = cpu.rmw & 0x80
	cpu.rmw <<= 1
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.a |= cpu.rmw
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.state = STATE_LAST
}

// ROL/AND group
func (cpu *MOS6510fn) fnO_RLA() {
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
}

// LSR/EOR group
func (cpu *MOS6510fn) fnO_SRE() {
	cpu.cFlag = cpu.rmw & 0x01
	cpu.rmw >>= 1
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.a ^= cpu.rmw
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.state = STATE_LAST
}

// ROR/ADC group
func (cpu *MOS6510fn) fnO_RRA() {
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
}

// DEC/CMP group
func (cpu *MOS6510fn) fnO_DCP() {
	cpu.rmw--
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.ar = uint16(cpu.a) - uint16(cpu.rmw)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = flag.BoolToUint8(cpu.ar < 0x100)
	cpu.state = STATE_LAST
}

// INC/SBC group
func (cpu *MOS6510fn) fnO_ISB() {
	cpu.rmw++
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.doSBC(cpu.rmw)
	cpu.state = STATE_LAST
}

// Complex functions
func (cpu *MOS6510fn) fnO_ANC_I() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	cpu.a &= cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.cFlag = cpu.nFlag & 0x80
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_ASR_I() {
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
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_ARR_I() {
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
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_ANE_I() {
	if cpu.rdyLow {
		cpu.stop = true
		return
	}
	data := cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.a = (cpu.a | 0xee) & cpu.x & data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_LXA_I() {
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
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_SBX_I() {
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
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_LAS() {
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
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_SHS() { // ar2 contains the high byte of the operand address
	cpu.sp = cpu.a & cpu.x
	cpu.banks.Write(cpu.ar, uint8((cpu.ar2+1)&uint16(cpu.sp)))
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_SHY() { // ar2 contains the high byte of the operand address
	cpu.banks.Write(cpu.ar, uint8(uint16(cpu.y)&(cpu.ar2+1)))
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_SHX() { // ar2 contains the high byte of the operand address
	cpu.banks.Write(cpu.ar, uint8(uint16(cpu.x)&(cpu.ar2+1)))
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnO_SHA() {
	// ar2 contains the high byte of the operand address
	cpu.banks.Write(cpu.ar, uint8(uint16(cpu.a)&uint16(cpu.x)&(cpu.ar2+1)))
	cpu.state = STATE_LAST
}

func (cpu *MOS6510fn) fnIllegalOp() {
	cpu.illegalOp(cpu.op, cpu.pc-1)
}

func (cpu *MOS6510fn) printRegisters(qCycle uint64, baLow bool) {
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
