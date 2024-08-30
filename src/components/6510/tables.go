package mos6510

// _modeTable Addressing mode for each opcode (first part of execution)
var _modeTable []func(*Core)

// _opTable Operation for each opcode (second part of execution)
var _opTable []func(*Core)

func init() {
	_modeTable = []func(*Core){
		instO_BRK, instA_INDX, instI_ILL_OP, instM_INDX, instA_ZERO, instA_ZERO, instM_ZERO, instM_ZERO, // 00
		instO_PHP, instO_ORA_I, instO_ASL_A, instO_ANC_I, instA_ABS, instA_ABS, instM_ABS, instM_ABS,
		instO_BPL, instAE_INDY, instI_ILL_OP, instM_INDY, instA_ZEROX, instA_ZEROX, instM_ZEROX, instM_ZEROX, // 10
		instO_CLC, instAE_ABSY, instO_NOP, instM_ABSY, instAE_ABSX, instAE_ABSX, instM_ABSX, instM_ABSX,
		instO_JSR, instA_INDX, instI_ILL_OP, instM_INDX, instA_ZERO, instA_ZERO, instM_ZERO, instM_ZERO, // 20
		instO_PLP, instO_AND_I, instO_ROL_A, instO_ANC_I, instA_ABS, instA_ABS, instM_ABS, instM_ABS,
		instO_BMI, instAE_INDY, instI_ILL_OP, instM_INDY, instA_ZEROX, instA_ZEROX, instM_ZEROX, instM_ZEROX, // 30
		instO_SEC, instAE_ABSY, instO_NOP, instM_ABSY, instAE_ABSX, instAE_ABSX, instM_ABSX, instM_ABSX,
		instO_RTI, instA_INDX, instI_ILL_OP, instM_INDX, instA_ZERO, instA_ZERO, instM_ZERO, instM_ZERO, // 40
		instO_PHA, instO_EOR_I, instO_LSR_A, instO_ASR_I, instO_JMP, instA_ABS, instM_ABS, instM_ABS,
		instO_BVC, instAE_INDY, instI_ILL_OP, instM_INDY, instA_ZEROX, instA_ZEROX, instM_ZEROX, instM_ZEROX, // 50
		instO_CLI, instAE_ABSY, instO_NOP, instM_ABSY, instAE_ABSX, instAE_ABSX, instM_ABSX, instM_ABSX,
		instO_RTS, instA_INDX, instI_ILL_OP, instM_INDX, instA_ZERO, instA_ZERO, instM_ZERO, instM_ZERO, // 60
		instO_PLA, instO_ADC_I, instO_ROR_A, instO_ARR_I, instA_ABS, instA_ABS, instM_ABS, instM_ABS,
		instO_BVS, instAE_INDY, instI_ILL_OP, instM_INDY, instA_ZEROX, instA_ZEROX, instM_ZEROX, instM_ZEROX, // 70
		instO_SEI, instAE_ABSY, instO_NOP, instM_ABSY, instAE_ABSX, instAE_ABSX, instM_ABSX, instM_ABSX,
		instO_NOP_I, instA_INDX, instO_NOP_I, instA_INDX, instA_ZERO, instA_ZERO, instA_ZERO, instA_ZERO, // 80
		instO_DEY, instO_NOP_I, instO_TXA, instO_ANE_I, instA_ABS, instA_ABS, instA_ABS, instA_ABS,
		instO_BCC, instA_INDY, instI_ILL_OP, instA_INDY, instA_ZEROX, instA_ZEROX, instA_ZEROY, instA_ZEROY, // 90
		instO_TYA, instA_ABSY, instO_TXS, instA_ABSY, instA_ABSX, instA_ABSX, instA_ABSY, instA_ABSY,
		instO_LDY_I, instA_INDX, instO_LDX_I, instA_INDX, instA_ZERO, instA_ZERO, instA_ZERO, instA_ZERO, // a0
		instO_TAY, instO_LDA_I, instO_TAX, instO_LXA_I, instA_ABS, instA_ABS, instA_ABS, instA_ABS,
		instO_BCS, instAE_INDY, instI_ILL_OP, instAE_INDY, instA_ZEROX, instA_ZEROX, instA_ZEROY, instA_ZEROY, // b0
		instO_CLV, instAE_ABSY, instO_TSX, instAE_ABSY, instAE_ABSX, instAE_ABSX, instAE_ABSY, instAE_ABSY,
		instO_CPY_I, instA_INDX, instO_NOP_I, instM_INDX, instA_ZERO, instA_ZERO, instM_ZERO, instM_ZERO, // c0
		instO_INY, instO_CMP_I, instO_DEX, instO_SBX_I, instA_ABS, instA_ABS, instM_ABS, instM_ABS,
		instO_BNE, instAE_INDY, instI_ILL_OP, instM_INDY, instA_ZEROX, instA_ZEROX, instM_ZEROX, instM_ZEROX, // d0
		instO_CLD, instAE_ABSY, instO_NOP, instM_ABSY, instAE_ABSX, instAE_ABSX, instM_ABSX, instM_ABSX,
		instO_CPX_I, instA_INDX, instO_NOP_I, instM_INDX, instA_ZERO, instA_ZERO, instM_ZERO, instM_ZERO, // e0
		instO_INX, instO_SBC_I, instO_NOP, instO_SBC_I, instA_ABS, instA_ABS, instM_ABS, instM_ABS,
		instO_BEQ, instAE_INDY, instI_ILL_OP, instM_INDY, instA_ZEROX, instA_ZEROX, instM_ZEROX, instM_ZEROX, // f0
		instO_SED, instAE_ABSY, instO_NOP, instM_ABSY, instAE_ABSX, instAE_ABSX, instM_ABSX, instM_ABSX,
	}

	_opTable = []func(*Core){
		instI_ILL_OP, instO_ORA, instI_ILL_OP, instO_SLO, instO_NOP_A, instO_ORA, instO_ASL, instO_SLO, // 00
		instI_ILL_OP, instI_ILL_OP, instI_ILL_OP, instI_ILL_OP, instO_NOP_A, instO_ORA, instO_ASL, instO_SLO,
		instI_ILL_OP, instO_ORA, instI_ILL_OP, instO_SLO, instO_NOP_A, instO_ORA, instO_ASL, instO_SLO, // 10
		instI_ILL_OP, instO_ORA, instI_ILL_OP, instO_SLO, instO_NOP_A, instO_ORA, instO_ASL, instO_SLO,
		instI_ILL_OP, instO_AND, instI_ILL_OP, instO_RLA, instO_BIT, instO_AND, instO_ROL, instO_RLA, // 20
		instI_ILL_OP, instI_ILL_OP, instI_ILL_OP, instI_ILL_OP, instO_BIT, instO_AND, instO_ROL, instO_RLA,
		instI_ILL_OP, instO_AND, instI_ILL_OP, instO_RLA, instO_NOP_A, instO_AND, instO_ROL, instO_RLA, // 30
		instI_ILL_OP, instO_AND, instI_ILL_OP, instO_RLA, instO_NOP_A, instO_AND, instO_ROL, instO_RLA,
		instI_ILL_OP, instO_EOR, instI_ILL_OP, instO_SRE, instO_NOP_A, instO_EOR, instO_LSR, instO_SRE, // 40
		instI_ILL_OP, instI_ILL_OP, instI_ILL_OP, instI_ILL_OP, instI_ILL_OP, instO_EOR, instO_LSR, instO_SRE,
		instI_ILL_OP, instO_EOR, instI_ILL_OP, instO_SRE, instO_NOP_A, instO_EOR, instO_LSR, instO_SRE, // 50
		instI_ILL_OP, instO_EOR, instI_ILL_OP, instO_SRE, instO_NOP_A, instO_EOR, instO_LSR, instO_SRE,
		instI_ILL_OP, instO_ADC, instI_ILL_OP, instO_RRA, instO_NOP_A, instO_ADC, instO_ROR, instO_RRA, // 60
		instI_ILL_OP, instI_ILL_OP, instI_ILL_OP, instI_ILL_OP, instO_JMP_I, instO_ADC, instO_ROR, instO_RRA,
		instI_ILL_OP, instO_ADC, instI_ILL_OP, instO_RRA, instO_NOP_A, instO_ADC, instO_ROR, instO_RRA, // 70
		instI_ILL_OP, instO_ADC, instI_ILL_OP, instO_RRA, instO_NOP_A, instO_ADC, instO_ROR, instO_RRA,
		instI_ILL_OP, instO_STA, instI_ILL_OP, instO_SAX, instO_STY, instO_STA, instO_STX, instO_SAX, // 80
		instI_ILL_OP, instI_ILL_OP, instI_ILL_OP, instI_ILL_OP, instO_STY, instO_STA, instO_STX, instO_SAX,
		instI_ILL_OP, instO_STA, instI_ILL_OP, instO_SHA, instO_STY, instO_STA, instO_STX, instO_SAX, // 90
		instI_ILL_OP, instO_STA, instI_ILL_OP, instO_SHS, instO_SHY, instO_STA, instO_SHX, instO_SHA,
		instI_ILL_OP, instO_LDA, instI_ILL_OP, instO_LAX, instO_LDY, instO_LDA, instO_LDX, instO_LAX, // a0
		instI_ILL_OP, instI_ILL_OP, instI_ILL_OP, instI_ILL_OP, instO_LDY, instO_LDA, instO_LDX, instO_LAX,
		instI_ILL_OP, instO_LDA, instI_ILL_OP, instO_LAX, instO_LDY, instO_LDA, instO_LDX, instO_LAX, // b0
		instI_ILL_OP, instO_LDA, instI_ILL_OP, instO_LAS, instO_LDY, instO_LDA, instO_LDX, instO_LAX,
		instI_ILL_OP, instO_CMP, instI_ILL_OP, instO_DCP, instO_CPY, instO_CMP, instO_DEC, instO_DCP, // c0
		instI_ILL_OP, instI_ILL_OP, instI_ILL_OP, instI_ILL_OP, instO_CPY, instO_CMP, instO_DEC, instO_DCP,
		instI_ILL_OP, instO_CMP, instI_ILL_OP, instO_DCP, instO_NOP_A, instO_CMP, instO_DEC, instO_DCP, // d0
		instI_ILL_OP, instO_CMP, instI_ILL_OP, instO_DCP, instO_NOP_A, instO_CMP, instO_DEC, instO_DCP,
		instI_ILL_OP, instO_SBC, instI_ILL_OP, instO_ISB, instO_CPX, instO_SBC, instO_INC, instO_ISB, // e0
		instI_ILL_OP, instI_ILL_OP, instI_ILL_OP, instI_ILL_OP, instO_CPX, instO_SBC, instO_INC, instO_ISB,
		instI_ILL_OP, instO_SBC, instI_ILL_OP, instO_ISB, instO_NOP_A, instO_SBC, instO_INC, instO_ISB, // f0
		instI_ILL_OP, instO_SBC, instI_ILL_OP, instO_ISB, instO_NOP_A, instO_SBC, instO_INC, instO_ISB,
	}
}

const (
	OpFlagIrqDisabled = 0x01
	OpFlagIrqEnabled  = 0x02
	OpFlagIntDelayed  = 0x04
)
