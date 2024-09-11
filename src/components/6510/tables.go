package mos6510

// _modeTable Addressing mode for each opcode (first part of execution)
var _modeTable []func(*CPU)

// _opTable Operation for each opcode (second part of execution)
var _opTable []func(*CPU)

// A Read effective address, no extra cycles -> Ap
// AE Read effective address, extra cycle on page crossing -> Ae
// M Read operand and write it back (for RMW instructions), no extra cycles -> Mp
// O Operations (_I = Immediate/Indirect, _A = Accumulator) --> Oi | Oa
// O Operation => Op

func init() {
	_modeTable = []func(*CPU){
		instOpBrk, instApIndX, instOpIll, instMpIndX, instApZero, instApZero, instMpZero, instMpZero, // 00
		instOpPhp, instOiOpa, instOaAsl, instOiAnc, instApAbs, instApAbs, instMpAbs, instMpAbs,
		instOpBpl, instAeIndy, instOpIll, instMpIndy, instApZeroX, instApZeroX, instMpZeroX, instMpZeroX, // 10
		instO_CLC, instAeAbsY, instO_NOP, instMpAbsY, instAeAbsX, instAeAbsX, instMpAbsX, instMpAbsX,
		instO_JSR, instApIndX, instOpIll, instMpIndX, instApZero, instApZero, instMpZero, instMpZero, // 20
		instO_PLP, instO_AND_I, instO_ROL_A, instOiAnc, instApAbs, instApAbs, instMpAbs, instMpAbs,
		instO_BMI, instAeIndy, instOpIll, instMpIndy, instApZeroX, instApZeroX, instMpZeroX, instMpZeroX, // 30
		instO_SEC, instAeAbsY, instO_NOP, instMpAbsY, instAeAbsX, instAeAbsX, instMpAbsX, instMpAbsX,
		instO_RTI, instApIndX, instOpIll, instMpIndX, instApZero, instApZero, instMpZero, instMpZero, // 40
		instO_PHA, instO_EOR_I, instO_LSR_A, instO_ASR_I, instO_JMP, instApAbs, instMpAbs, instMpAbs,
		instO_BVC, instAeIndy, instOpIll, instMpIndy, instApZeroX, instApZeroX, instMpZeroX, instMpZeroX, // 50
		instO_CLI, instAeAbsY, instO_NOP, instMpAbsY, instAeAbsX, instAeAbsX, instMpAbsX, instMpAbsX,
		instO_RTS, instApIndX, instOpIll, instMpIndX, instApZero, instApZero, instMpZero, instMpZero, // 60
		instO_PLA, instO_ADC_I, instO_ROR_A, instO_ARR_I, instApAbs, instApAbs, instMpAbs, instMpAbs,
		instO_BVS, instAeIndy, instOpIll, instMpIndy, instApZeroX, instApZeroX, instMpZeroX, instMpZeroX, // 70
		instO_SEI, instAeAbsY, instO_NOP, instMpAbsY, instAeAbsX, instAeAbsX, instMpAbsX, instMpAbsX,
		instO_NOP_I, instApIndX, instO_NOP_I, instApIndX, instApZero, instApZero, instApZero, instApZero, // 80
		instO_DEY, instO_NOP_I, instO_TXA, instO_ANE_I, instApAbs, instApAbs, instApAbs, instApAbs,
		instO_BCC, instApIndY, instOpIll, instApIndY, instApZeroX, instApZeroX, instApZeroY, instApZeroY, // 90
		instO_TYA, instApAbsY, instO_TXS, instApAbsY, instApAbsX, instApAbsX, instApAbsY, instApAbsY,
		instO_LDY_I, instApIndX, instO_LDX_I, instApIndX, instApZero, instApZero, instApZero, instApZero, // a0
		instO_TAY, instO_LDA_I, instO_TAX, instO_LXA_I, instApAbs, instApAbs, instApAbs, instApAbs,
		instO_BCS, instAeIndy, instOpIll, instAeIndy, instApZeroX, instApZeroX, instApZeroY, instApZeroY, // b0
		instO_CLV, instAeAbsY, instO_TSX, instAeAbsY, instAeAbsX, instAeAbsX, instAeAbsY, instAeAbsY,
		instO_CPY_I, instApIndX, instO_NOP_I, instMpIndX, instApZero, instApZero, instMpZero, instMpZero, // c0
		instO_INY, instO_CMP_I, instO_DEX, instO_SBX_I, instApAbs, instApAbs, instMpAbs, instMpAbs,
		instO_BNE, instAeIndy, instOpIll, instMpIndy, instApZeroX, instApZeroX, instMpZeroX, instMpZeroX, // d0
		instO_CLD, instAeAbsY, instO_NOP, instMpAbsY, instAeAbsX, instAeAbsX, instMpAbsX, instMpAbsX,
		instO_CPX_I, instApIndX, instO_NOP_I, instMpIndX, instApZero, instApZero, instMpZero, instMpZero, // e0
		instO_INX, instO_SBC_I, instO_NOP, instO_SBC_I, instApAbs, instApAbs, instMpAbs, instMpAbs,
		instO_BEQ, instAeIndy, instOpIll, instMpIndy, instApZeroX, instApZeroX, instMpZeroX, instMpZeroX, // f0
		instO_SED, instAeAbsY, instO_NOP, instMpAbsY, instAeAbsX, instAeAbsX, instMpAbsX, instMpAbsX,
	}

	_opTable = []func(*CPU){
		instOpIll, instO_ORA, instOpIll, instO_SLO, instO_NOP_A, instO_ORA, instO_ASL, instO_SLO, // 00
		instOpIll, instOpIll, instOpIll, instOpIll, instO_NOP_A, instO_ORA, instO_ASL, instO_SLO,
		instOpIll, instO_ORA, instOpIll, instO_SLO, instO_NOP_A, instO_ORA, instO_ASL, instO_SLO, // 10
		instOpIll, instO_ORA, instOpIll, instO_SLO, instO_NOP_A, instO_ORA, instO_ASL, instO_SLO,
		instOpIll, instO_AND, instOpIll, instO_RLA, instO_BIT, instO_AND, instO_ROL, instO_RLA, // 20
		instOpIll, instOpIll, instOpIll, instOpIll, instO_BIT, instO_AND, instO_ROL, instO_RLA,
		instOpIll, instO_AND, instOpIll, instO_RLA, instO_NOP_A, instO_AND, instO_ROL, instO_RLA, // 30
		instOpIll, instO_AND, instOpIll, instO_RLA, instO_NOP_A, instO_AND, instO_ROL, instO_RLA,
		instOpIll, instO_EOR, instOpIll, instO_SRE, instO_NOP_A, instO_EOR, instO_LSR, instO_SRE, // 40
		instOpIll, instOpIll, instOpIll, instOpIll, instOpIll, instO_EOR, instO_LSR, instO_SRE,
		instOpIll, instO_EOR, instOpIll, instO_SRE, instO_NOP_A, instO_EOR, instO_LSR, instO_SRE, // 50
		instOpIll, instO_EOR, instOpIll, instO_SRE, instO_NOP_A, instO_EOR, instO_LSR, instO_SRE,
		instOpIll, instO_ADC, instOpIll, instO_RRA, instO_NOP_A, instO_ADC, instO_ROR, instO_RRA, // 60
		instOpIll, instOpIll, instOpIll, instOpIll, instO_JMP_I, instO_ADC, instO_ROR, instO_RRA,
		instOpIll, instO_ADC, instOpIll, instO_RRA, instO_NOP_A, instO_ADC, instO_ROR, instO_RRA, // 70
		instOpIll, instO_ADC, instOpIll, instO_RRA, instO_NOP_A, instO_ADC, instO_ROR, instO_RRA,
		instOpIll, instO_STA, instOpIll, instO_SAX, instO_STY, instO_STA, instO_STX, instO_SAX, // 80
		instOpIll, instOpIll, instOpIll, instOpIll, instO_STY, instO_STA, instO_STX, instO_SAX,
		instOpIll, instO_STA, instOpIll, instO_SHA, instO_STY, instO_STA, instO_STX, instO_SAX, // 90
		instOpIll, instO_STA, instOpIll, instO_SHS, instO_SHY, instO_STA, instO_SHX, instO_SHA,
		instOpIll, instO_LDA, instOpIll, instO_LAX, instO_LDY, instO_LDA, instO_LDX, instO_LAX, // a0
		instOpIll, instOpIll, instOpIll, instOpIll, instO_LDY, instO_LDA, instO_LDX, instO_LAX,
		instOpIll, instO_LDA, instOpIll, instO_LAX, instO_LDY, instO_LDA, instO_LDX, instO_LAX, // b0
		instOpIll, instO_LDA, instOpIll, instO_LAS, instO_LDY, instO_LDA, instO_LDX, instO_LAX,
		instOpIll, instO_CMP, instOpIll, instO_DCP, instO_CPY, instO_CMP, instO_DEC, instO_DCP, // c0
		instOpIll, instOpIll, instOpIll, instOpIll, instO_CPY, instO_CMP, instO_DEC, instO_DCP,
		instOpIll, instO_CMP, instOpIll, instO_DCP, instO_NOP_A, instO_CMP, instO_DEC, instO_DCP, // d0
		instOpIll, instO_CMP, instOpIll, instO_DCP, instO_NOP_A, instO_CMP, instO_DEC, instO_DCP,
		instOpIll, instO_SBC, instOpIll, instO_ISB, instO_CPX, instO_SBC, instO_INC, instO_ISB, // e0
		instOpIll, instOpIll, instOpIll, instOpIll, instO_CPX, instO_SBC, instO_INC, instO_ISB,
		instOpIll, instO_SBC, instOpIll, instO_ISB, instO_NOP_A, instO_SBC, instO_INC, instO_ISB, // f0
		instOpIll, instO_SBC, instOpIll, instO_ISB, instO_NOP_A, instO_SBC, instO_INC, instO_ISB,
	}
}
