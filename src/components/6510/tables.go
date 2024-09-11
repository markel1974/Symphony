package mos6510

// _modeTable Addressing mode for each opcode (first part of execution)
var _modeTable []func(*CPU)

// _opTable Operation for each opcode (second part of execution)
var _opTable []func(*CPU)

// Ap -> Read effective address, no extra cycles
// Ae -> Read effective address, extra cycle on page crossing
// Mp -> Read operand and write it back (for RMW instructions), no extra cycles
// Oi -> Operation Immediate/Indirect
// OA -> Operation Accumulator
// Op -> Operation

func init() {
	_modeTable = []func(*CPU){
		instOpBRK, instApIndX, instOpIll, instMpIndX, instApZero, instApZero, instMpZero, instMpZero, // 00
		instOpPHP, instOiOPA, instOaASL, instOiAnc, instApABS, instApABS, instMpAbs, instMpAbs,
		instOpBpl, instAeIndy, instOpIll, instMpIndy, instApZeroX, instApZeroX, instMpZeroX, instMpZeroX, // 10
		instOpCLC, instAeAbsY, instOpNOP, instMpAbsY, instAeAbsX, instAeAbsX, instMpAbsX, instMpAbsX,
		instOpJSR, instApIndX, instOpIll, instMpIndX, instApZero, instApZero, instMpZero, instMpZero, // 20
		instOpPLP, instOiAND, instOaROL, instOiAnc, instApABS, instApABS, instMpAbs, instMpAbs,
		instOpBMI, instAeIndy, instOpIll, instMpIndy, instApZeroX, instApZeroX, instMpZeroX, instMpZeroX, // 30
		instOpSEC, instAeAbsY, instOpNOP, instMpAbsY, instAeAbsX, instAeAbsX, instMpAbsX, instMpAbsX,
		instOpRTI, instApIndX, instOpIll, instMpIndX, instApZero, instApZero, instMpZero, instMpZero, // 40
		instOpPHA, instOiEOR, instOaLSR, instOiASR, instOpJMP, instApABS, instMpAbs, instMpAbs,
		instOpBVC, instAeIndy, instOpIll, instMpIndy, instApZeroX, instApZeroX, instMpZeroX, instMpZeroX, // 50
		instOpCLI, instAeAbsY, instOpNOP, instMpAbsY, instAeAbsX, instAeAbsX, instMpAbsX, instMpAbsX,
		instOpRTS, instApIndX, instOpIll, instMpIndX, instApZero, instApZero, instMpZero, instMpZero, // 60
		instOpPLA, instOiADC, instOaROR, instOiARR, instApABS, instApABS, instMpAbs, instMpAbs,
		instOpBVS, instAeIndy, instOpIll, instMpIndy, instApZeroX, instApZeroX, instMpZeroX, instMpZeroX, // 70
		instOpSEI, instAeAbsY, instOpNOP, instMpAbsY, instAeAbsX, instAeAbsX, instMpAbsX, instMpAbsX,
		instOiNOP, instApIndX, instOiNOP, instApIndX, instApZero, instApZero, instApZero, instApZero, // 80
		instOpDEY, instOiNOP, instOpTXA, instOiANE, instApABS, instApABS, instApABS, instApABS,
		instOpBCC, instApIndY, instOpIll, instApIndY, instApZeroX, instApZeroX, instApZeroY, instApZeroY, // 90
		instOpTYA, instApAbsY, instOpTXS, instApAbsY, instApAbsX, instApAbsX, instApAbsY, instApAbsY,
		instOiLDY, instApIndX, instOiLDX, instApIndX, instApZero, instApZero, instApZero, instApZero, // a0
		instOpTAY, instOiLDA, instOpTAX, instOiLXA, instApABS, instApABS, instApABS, instApABS,
		instOpBCS, instAeIndy, instOpIll, instAeIndy, instApZeroX, instApZeroX, instApZeroY, instApZeroY, // b0
		instOpCLV, instAeAbsY, instOpTSX, instAeAbsY, instAeAbsX, instAeAbsX, instAeAbsY, instAeAbsY,
		instOiCPY, instApIndX, instOiNOP, instMpIndX, instApZero, instApZero, instMpZero, instMpZero, // c0
		instOpINY, instOiCMP, instOpDEX, instOiSBX, instApABS, instApABS, instMpAbs, instMpAbs,
		instOpBNE, instAeIndy, instOpIll, instMpIndy, instApZeroX, instApZeroX, instMpZeroX, instMpZeroX, // d0
		instOpCLD, instAeAbsY, instOpNOP, instMpAbsY, instAeAbsX, instAeAbsX, instMpAbsX, instMpAbsX,
		instOiCPX, instApIndX, instOiNOP, instMpIndX, instApZero, instApZero, instMpZero, instMpZero, // e0
		instOpINX, instOiSBC, instOpNOP, instOiSBC, instApABS, instApABS, instMpAbs, instMpAbs,
		instOpBEQ, instAeIndy, instOpIll, instMpIndy, instApZeroX, instApZeroX, instMpZeroX, instMpZeroX, // f0
		instOpSED, instAeAbsY, instOpNOP, instMpAbsY, instAeAbsX, instAeAbsX, instMpAbsX, instMpAbsX,
	}

	_opTable = []func(*CPU){
		instOpIll, instOpORA, instOpIll, instOpSLO, instOaNOP, instOpORA, instOpASL, instOpSLO, // 00
		instOpIll, instOpIll, instOpIll, instOpIll, instOaNOP, instOpORA, instOpASL, instOpSLO,
		instOpIll, instOpORA, instOpIll, instOpSLO, instOaNOP, instOpORA, instOpASL, instOpSLO, // 10
		instOpIll, instOpORA, instOpIll, instOpSLO, instOaNOP, instOpORA, instOpASL, instOpSLO,
		instOpIll, instOpAND, instOpIll, instOpRLA, instOpBIT, instOpAND, instOpROL, instOpRLA, // 20
		instOpIll, instOpIll, instOpIll, instOpIll, instOpBIT, instOpAND, instOpROL, instOpRLA,
		instOpIll, instOpAND, instOpIll, instOpRLA, instOaNOP, instOpAND, instOpROL, instOpRLA, // 30
		instOpIll, instOpAND, instOpIll, instOpRLA, instOaNOP, instOpAND, instOpROL, instOpRLA,
		instOpIll, instOpEOR, instOpIll, instOpSRE, instOaNOP, instOpEOR, instOpLSR, instOpSRE, // 40
		instOpIll, instOpIll, instOpIll, instOpIll, instOpIll, instOpEOR, instOpLSR, instOpSRE,
		instOpIll, instOpEOR, instOpIll, instOpSRE, instOaNOP, instOpEOR, instOpLSR, instOpSRE, // 50
		instOpIll, instOpEOR, instOpIll, instOpSRE, instOaNOP, instOpEOR, instOpLSR, instOpSRE,
		instOpIll, instOpADC, instOpIll, instOpRRA, instOaNOP, instOpADC, instOpROR, instOpRRA, // 60
		instOpIll, instOpIll, instOpIll, instOpIll, instOiJMP, instOpADC, instOpROR, instOpRRA,
		instOpIll, instOpADC, instOpIll, instOpRRA, instOaNOP, instOpADC, instOpROR, instOpRRA, // 70
		instOpIll, instOpADC, instOpIll, instOpRRA, instOaNOP, instOpADC, instOpROR, instOpRRA,
		instOpIll, instOpSTA, instOpIll, instOpSAX, instOpSTY, instOpSTA, instOpSTX, instOpSAX, // 80
		instOpIll, instOpIll, instOpIll, instOpIll, instOpSTY, instOpSTA, instOpSTX, instOpSAX,
		instOpIll, instOpSTA, instOpIll, instOpSHA, instOpSTY, instOpSTA, instOpSTX, instOpSAX, // 90
		instOpIll, instOpSTA, instOpIll, instOpSHS, instOpSHY, instOpSTA, instOpSHX, instOpSHA,
		instOpIll, instOpLDA, instOpIll, instOpLAX, instOpLDY, instOpLDA, instOpLDX, instOpLAX, // a0
		instOpIll, instOpIll, instOpIll, instOpIll, instOpLDY, instOpLDA, instOpLDX, instOpLAX,
		instOpIll, instOpLDA, instOpIll, instOpLAX, instOpLDY, instOpLDA, instOpLDX, instOpLAX, // b0
		instOpIll, instOpLDA, instOpIll, instOpLAS, instOpLDY, instOpLDA, instOpLDX, instOpLAX,
		instOpIll, instOpCMP, instOpIll, instOpDCP, instOpCPY, instOpCMP, instOpDEC, instOpDCP, // c0
		instOpIll, instOpIll, instOpIll, instOpIll, instOpCPY, instOpCMP, instOpDEC, instOpDCP,
		instOpIll, instOpCMP, instOpIll, instOpDCP, instOaNOP, instOpCMP, instOpDEC, instOpDCP, // d0
		instOpIll, instOpCMP, instOpIll, instOpDCP, instOaNOP, instOpCMP, instOpDEC, instOpDCP,
		instOpIll, instOpSBC, instOpIll, instOpISB, instOpCPX, instOpSBC, instOpINC, instOpISB, // e0
		instOpIll, instOpIll, instOpIll, instOpIll, instOpCPX, instOpSBC, instOpINC, instOpISB,
		instOpIll, instOpSBC, instOpIll, instOpISB, instOaNOP, instOpSBC, instOpINC, instOpISB, // f0
		instOpIll, instOpSBC, instOpIll, instOpISB, instOaNOP, instOpSBC, instOpINC, instOpISB,
	}
}
