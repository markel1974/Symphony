package mos6510fn

// _modeTable Addressing mode for each opcode (first part of execution)
var _modeTable []func(mos6510 *MOS6510)

// _opTable Operation for each opcode (second part of execution)
var _opTable []func(*MOS6510)

func init() {
	_modeTable = []func(mos6510 *MOS6510){
		fnO_BRK, fnA_INDX, fnI_ILL_OP, fnM_INDX, fnA_ZERO, fnA_ZERO, fnM_ZERO, fnM_ZERO, // 00
		fnO_PHP, fnO_ORA_I, fnO_ASL_A, fnO_ANC_I, fnA_ABS, fnA_ABS, fnM_ABS, fnM_ABS,
		fnO_BPL, fnAE_INDY, fnI_ILL_OP, fnM_INDY, fnA_ZEROX, fnA_ZEROX, fnM_ZEROX, fnM_ZEROX, // 10
		fnO_CLC, fnAE_ABSY, fnO_NOP, fnM_ABSY, fnAE_ABSX, fnAE_ABSX, fnM_ABSX, fnM_ABSX,
		fnO_JSR, fnA_INDX, fnI_ILL_OP, fnM_INDX, fnA_ZERO, fnA_ZERO, fnM_ZERO, fnM_ZERO, // 20
		fnO_PLP, fnO_AND_I, fnO_ROL_A, fnO_ANC_I, fnA_ABS, fnA_ABS, fnM_ABS, fnM_ABS,
		fnO_BMI, fnAE_INDY, fnI_ILL_OP, fnM_INDY, fnA_ZEROX, fnA_ZEROX, fnM_ZEROX, fnM_ZEROX, // 30
		fnO_SEC, fnAE_ABSY, fnO_NOP, fnM_ABSY, fnAE_ABSX, fnAE_ABSX, fnM_ABSX, fnM_ABSX,
		fnO_RTI, fnA_INDX, fnI_ILL_OP, fnM_INDX, fnA_ZERO, fnA_ZERO, fnM_ZERO, fnM_ZERO, // 40
		fnO_PHA, fnO_EOR_I, fnO_LSR_A, fnO_ASR_I, fnO_JMP, fnA_ABS, fnM_ABS, fnM_ABS,
		fnO_BVC, fnAE_INDY, fnI_ILL_OP, fnM_INDY, fnA_ZEROX, fnA_ZEROX, fnM_ZEROX, fnM_ZEROX, // 50
		fnO_CLI, fnAE_ABSY, fnO_NOP, fnM_ABSY, fnAE_ABSX, fnAE_ABSX, fnM_ABSX, fnM_ABSX,
		fnO_RTS, fnA_INDX, fnI_ILL_OP, fnM_INDX, fnA_ZERO, fnA_ZERO, fnM_ZERO, fnM_ZERO, // 60
		fnO_PLA, fnO_ADC_I, fnO_ROR_A, fnO_ARR_I, fnA_ABS, fnA_ABS, fnM_ABS, fnM_ABS,
		fnO_BVS, fnAE_INDY, fnI_ILL_OP, fnM_INDY, fnA_ZEROX, fnA_ZEROX, fnM_ZEROX, fnM_ZEROX, // 70
		fnO_SEI, fnAE_ABSY, fnO_NOP, fnM_ABSY, fnAE_ABSX, fnAE_ABSX, fnM_ABSX, fnM_ABSX,
		fnO_NOP_I, fnA_INDX, fnO_NOP_I, fnA_INDX, fnA_ZERO, fnA_ZERO, fnA_ZERO, fnA_ZERO, // 80
		fnO_DEY, fnO_NOP_I, fnO_TXA, fnO_ANE_I, fnA_ABS, fnA_ABS, fnA_ABS, fnA_ABS,
		fnO_BCC, fnA_INDY, fnI_ILL_OP, fnA_INDY, fnA_ZEROX, fnA_ZEROX, fnA_ZEROY, fnA_ZEROY, // 90
		fnO_TYA, fnA_ABSY, fnO_TXS, fnA_ABSY, fnA_ABSX, fnA_ABSX, fnA_ABSY, fnA_ABSY,
		fnO_LDY_I, fnA_INDX, fnO_LDX_I, fnA_INDX, fnA_ZERO, fnA_ZERO, fnA_ZERO, fnA_ZERO, // a0
		fnO_TAY, fnO_LDA_I, fnO_TAX, fnO_LXA_I, fnA_ABS, fnA_ABS, fnA_ABS, fnA_ABS,
		fnO_BCS, fnAE_INDY, fnI_ILL_OP, fnAE_INDY, fnA_ZEROX, fnA_ZEROX, fnA_ZEROY, fnA_ZEROY, // b0
		fnO_CLV, fnAE_ABSY, fnO_TSX, fnAE_ABSY, fnAE_ABSX, fnAE_ABSX, fnAE_ABSY, fnAE_ABSY,
		fnO_CPY_I, fnA_INDX, fnO_NOP_I, fnM_INDX, fnA_ZERO, fnA_ZERO, fnM_ZERO, fnM_ZERO, // c0
		fnO_INY, fnO_CMP_I, fnO_DEX, fnO_SBX_I, fnA_ABS, fnA_ABS, fnM_ABS, fnM_ABS,
		fnO_BNE, fnAE_INDY, fnI_ILL_OP, fnM_INDY, fnA_ZEROX, fnA_ZEROX, fnM_ZEROX, fnM_ZEROX, // d0
		fnO_CLD, fnAE_ABSY, fnO_NOP, fnM_ABSY, fnAE_ABSX, fnAE_ABSX, fnM_ABSX, fnM_ABSX,
		fnO_CPX_I, fnA_INDX, fnO_NOP_I, fnM_INDX, fnA_ZERO, fnA_ZERO, fnM_ZERO, fnM_ZERO, // e0
		fnO_INX, fnO_SBC_I, fnO_NOP, fnO_SBC_I, fnA_ABS, fnA_ABS, fnM_ABS, fnM_ABS,
		fnO_BEQ, fnAE_INDY, fnI_ILL_OP, fnM_INDY, fnA_ZEROX, fnA_ZEROX, fnM_ZEROX, fnM_ZEROX, // f0
		fnO_SED, fnAE_ABSY, fnO_NOP, fnM_ABSY, fnAE_ABSX, fnAE_ABSX, fnM_ABSX, fnM_ABSX,
	}

	_opTable = []func(*MOS6510){
		fnI_ILL_OP, fnO_ORA, fnI_ILL_OP, fnO_SLO, fnO_NOP_A, fnO_ORA, fnO_ASL, fnO_SLO, // 00
		fnI_ILL_OP, fnI_ILL_OP, fnI_ILL_OP, fnI_ILL_OP, fnO_NOP_A, fnO_ORA, fnO_ASL, fnO_SLO,
		fnI_ILL_OP, fnO_ORA, fnI_ILL_OP, fnO_SLO, fnO_NOP_A, fnO_ORA, fnO_ASL, fnO_SLO, // 10
		fnI_ILL_OP, fnO_ORA, fnI_ILL_OP, fnO_SLO, fnO_NOP_A, fnO_ORA, fnO_ASL, fnO_SLO,
		fnI_ILL_OP, fnO_AND, fnI_ILL_OP, fnO_RLA, fnO_BIT, fnO_AND, fnO_ROL, fnO_RLA, // 20
		fnI_ILL_OP, fnI_ILL_OP, fnI_ILL_OP, fnI_ILL_OP, fnO_BIT, fnO_AND, fnO_ROL, fnO_RLA,
		fnI_ILL_OP, fnO_AND, fnI_ILL_OP, fnO_RLA, fnO_NOP_A, fnO_AND, fnO_ROL, fnO_RLA, // 30
		fnI_ILL_OP, fnO_AND, fnI_ILL_OP, fnO_RLA, fnO_NOP_A, fnO_AND, fnO_ROL, fnO_RLA,
		fnI_ILL_OP, fnO_EOR, fnI_ILL_OP, fnO_SRE, fnO_NOP_A, fnO_EOR, fnO_LSR, fnO_SRE, // 40
		fnI_ILL_OP, fnI_ILL_OP, fnI_ILL_OP, fnI_ILL_OP, fnI_ILL_OP, fnO_EOR, fnO_LSR, fnO_SRE,
		fnI_ILL_OP, fnO_EOR, fnI_ILL_OP, fnO_SRE, fnO_NOP_A, fnO_EOR, fnO_LSR, fnO_SRE, // 50
		fnI_ILL_OP, fnO_EOR, fnI_ILL_OP, fnO_SRE, fnO_NOP_A, fnO_EOR, fnO_LSR, fnO_SRE,
		fnI_ILL_OP, fnO_ADC, fnI_ILL_OP, fnO_RRA, fnO_NOP_A, fnO_ADC, fnO_ROR, fnO_RRA, // 60
		fnI_ILL_OP, fnI_ILL_OP, fnI_ILL_OP, fnI_ILL_OP, fnO_JMP_I, fnO_ADC, fnO_ROR, fnO_RRA,
		fnI_ILL_OP, fnO_ADC, fnI_ILL_OP, fnO_RRA, fnO_NOP_A, fnO_ADC, fnO_ROR, fnO_RRA, // 70
		fnI_ILL_OP, fnO_ADC, fnI_ILL_OP, fnO_RRA, fnO_NOP_A, fnO_ADC, fnO_ROR, fnO_RRA,
		fnI_ILL_OP, fnO_STA, fnI_ILL_OP, fnO_SAX, fnO_STY, fnO_STA, fnO_STX, fnO_SAX, // 80
		fnI_ILL_OP, fnI_ILL_OP, fnI_ILL_OP, fnI_ILL_OP, fnO_STY, fnO_STA, fnO_STX, fnO_SAX,
		fnI_ILL_OP, fnO_STA, fnI_ILL_OP, fnO_SHA, fnO_STY, fnO_STA, fnO_STX, fnO_SAX, // 90
		fnI_ILL_OP, fnO_STA, fnI_ILL_OP, fnO_SHS, fnO_SHY, fnO_STA, fnO_SHX, fnO_SHA,
		fnI_ILL_OP, fnO_LDA, fnI_ILL_OP, fnO_LAX, fnO_LDY, fnO_LDA, fnO_LDX, fnO_LAX, // a0
		fnI_ILL_OP, fnI_ILL_OP, fnI_ILL_OP, fnI_ILL_OP, fnO_LDY, fnO_LDA, fnO_LDX, fnO_LAX,
		fnI_ILL_OP, fnO_LDA, fnI_ILL_OP, fnO_LAX, fnO_LDY, fnO_LDA, fnO_LDX, fnO_LAX, // b0
		fnI_ILL_OP, fnO_LDA, fnI_ILL_OP, fnO_LAS, fnO_LDY, fnO_LDA, fnO_LDX, fnO_LAX,
		fnI_ILL_OP, fnO_CMP, fnI_ILL_OP, fnO_DCP, fnO_CPY, fnO_CMP, fnO_DEC, fnO_DCP, // c0
		fnI_ILL_OP, fnI_ILL_OP, fnI_ILL_OP, fnI_ILL_OP, fnO_CPY, fnO_CMP, fnO_DEC, fnO_DCP,
		fnI_ILL_OP, fnO_CMP, fnI_ILL_OP, fnO_DCP, fnO_NOP_A, fnO_CMP, fnO_DEC, fnO_DCP, // d0
		fnI_ILL_OP, fnO_CMP, fnI_ILL_OP, fnO_DCP, fnO_NOP_A, fnO_CMP, fnO_DEC, fnO_DCP,
		fnI_ILL_OP, fnO_SBC, fnI_ILL_OP, fnO_ISB, fnO_CPX, fnO_SBC, fnO_INC, fnO_ISB, // e0
		fnI_ILL_OP, fnI_ILL_OP, fnI_ILL_OP, fnI_ILL_OP, fnO_CPX, fnO_SBC, fnO_INC, fnO_ISB,
		fnI_ILL_OP, fnO_SBC, fnI_ILL_OP, fnO_ISB, fnO_NOP_A, fnO_SBC, fnO_INC, fnO_ISB, // f0
		fnI_ILL_OP, fnO_SBC, fnI_ILL_OP, fnO_ISB, fnO_NOP_A, fnO_SBC, fnO_INC, fnO_ISB,
	}
}

const (
	OpFlagIrqDisabled = 0x01
	OpFlagIrqEnabled  = 0x02
	OpFlagIntDelayed  = 0x04
)
