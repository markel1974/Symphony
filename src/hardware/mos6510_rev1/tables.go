package mos6510_rev1

//https://www.c64-wiki.com/wiki/Opcode

//ORA (short for "Logical OR on Accumulator")
//ASL (short for "Arithmetic Shift Left")

// Ap -> Read effective address, no extra cycles
// Ae -> Read effective address, extra cycle on page crossing
// Mp -> Read operand and write it back (for RMW Instructions), no extra cycles
// Oi -> Operation Immediate/Indirect
// Oa -> Operation Accumulator
// Op -> Operation

// CreateModeTable initializes and returns a slice of instruction mode functions for a CPU.
func CreateModeTable() []func(*CPU) {
	modeTable := []func(*CPU){
		InstOpBRK, InstApINDx, InstOpJAM, InstMpINDx, InstApZER, InstApZER, InstMpZER, InstMpZER, // 00
		InstOpPHP, InstOiOPA, InstOaASL, InstOiANC, InstApABS, InstApABS, InstMpABS, InstMpABS,
		InstOpBPL, InstAeINDy, InstOpJAM, InstMpINDy, InstApZERx, InstApZERx, InstMpZERx, InstMpZERx, // 10
		InstOpCLC, InstAeABSy, InstOpNOP, InstMpABSy, InstAeABSx, InstAeABSx, InstMpABSx, InstMpABSx,
		InstOpJSR, InstApINDx, InstOpJAM, InstMpINDx, InstApZER, InstApZER, InstMpZER, InstMpZER, // 20
		InstOpPLP, InstOiAND, InstOaROL, InstOiANC, InstApABS, InstApABS, InstMpABS, InstMpABS,
		InstOpBMI, InstAeINDy, InstOpJAM, InstMpINDy, InstApZERx, InstApZERx, InstMpZERx, InstMpZERx, // 30
		InstOpSEC, InstAeABSy, InstOpNOP, InstMpABSy, InstAeABSx, InstAeABSx, InstMpABSx, InstMpABSx,
		InstOpRTI, InstApINDx, InstOpJAM, InstMpINDx, InstApZER, InstApZER, InstMpZER, InstMpZER, // 40
		InstOpPHA, InstOiEOR, InstOaLSR, InstOiASR, InstOpJMP, InstApABS, InstMpABS, InstMpABS,
		InstOpBVC, InstAeINDy, InstOpJAM, InstMpINDy, InstApZERx, InstApZERx, InstMpZERx, InstMpZERx, // 50
		InstOpCLI, InstAeABSy, InstOpNOP, InstMpABSy, InstAeABSx, InstAeABSx, InstMpABSx, InstMpABSx,
		InstOpRTS, InstApINDx, InstOpJAM, InstMpINDx, InstApZER, InstApZER, InstMpZER, InstMpZER, // 60
		InstOpPLA, InstOiADC, InstOaROR, InstOiARR, InstApABS, InstApABS, InstMpABS, InstMpABS,
		InstOpBVS, InstAeINDy, InstOpJAM, InstMpINDy, InstApZERx, InstApZERx, InstMpZERx, InstMpZERx, // 70
		InstOpSEI, InstAeABSy, InstOpNOP, InstMpABSy, InstAeABSx, InstAeABSx, InstMpABSx, InstMpABSx,
		InstOiNOP, InstApINDx, InstOiNOP, InstApINDx, InstApZER, InstApZER, InstApZER, InstApZER, // 80
		InstOpDEY, InstOiNOP, InstOpTXA, InstOiANE, InstApABS, InstApABS, InstApABS, InstApABS,
		InstOpBCC, InstApINDy, InstOpJAM, InstApINDy, InstApZERx, InstApZERx, InstApZERy, InstApZERy, // 90
		InstOpTYA, InstApABSy, InstOpTXS, InstApABSy, InstApABSx, InstApABSx, InstApABSy, InstApABSy,
		InstOiLDY, InstApINDx, InstOiLDX, InstApINDx, InstApZER, InstApZER, InstApZER, InstApZER, // a0
		InstOpTAY, InstOiLDA, InstOpTAX, InstOiLXA, InstApABS, InstApABS, InstApABS, InstApABS,
		InstOpBCS, InstAeINDy, InstOpJAM, InstAeINDy, InstApZERx, InstApZERx, InstApZERy, InstApZERy, // b0
		InstOpCLV, InstAeABSy, InstOpTSX, InstAeABSy, InstAeABSx, InstAeABSx, InstAeABSy, InstAeABSy,
		InstOiCPY, InstApINDx, InstOiNOP, InstMpINDx, InstApZER, InstApZER, InstMpZER, InstMpZER, // c0
		InstOpINY, InstOiCMP, InstOpDEX, InstOiSBX, InstApABS, InstApABS, InstMpABS, InstMpABS,
		InstOpBNE, InstAeINDy, InstOpJAM, InstMpINDy, InstApZERx, InstApZERx, InstMpZERx, InstMpZERx, // d0
		InstOpCLD, InstAeABSy, InstOpNOP, InstMpABSy, InstAeABSx, InstAeABSx, InstMpABSx, InstMpABSx,
		InstOiCPX, InstApINDx, InstOiNOP, InstMpINDx, InstApZER, InstApZER, InstMpZER, InstMpZER, // e0
		InstOpINX, InstOiSBC, InstOpNOP, InstOiSBC, InstApABS, InstApABS, InstMpABS, InstMpABS,
		InstOpBEQ, InstAeINDy, InstOpJAM, InstMpINDy, InstApZERx, InstApZERx, InstMpZERx, InstMpZERx, // f0
		InstOpSED, InstAeABSy, InstOpNOP, InstMpABSy, InstAeABSx, InstAeABSx, InstMpABSx, InstMpABSx,
	}
	return modeTable
}

// CreateOpTable initializes and returns an operation table containing a list of CPU instruction functions.
func CreateOpTable() []func(cpu *CPU) {
	opTable := []func(*CPU){
		InstOpJAM, InstOpORA, InstOpJAM, InstOpSLO, InstOaNOP, InstOpORA, InstOpASL, InstOpSLO, // 00
		InstOpJAM, InstOpJAM, InstOpJAM, InstOpJAM, InstOaNOP, InstOpORA, InstOpASL, InstOpSLO,
		InstOpJAM, InstOpORA, InstOpJAM, InstOpSLO, InstOaNOP, InstOpORA, InstOpASL, InstOpSLO, // 10
		InstOpJAM, InstOpORA, InstOpJAM, InstOpSLO, InstOaNOP, InstOpORA, InstOpASL, InstOpSLO,
		InstOpJAM, InstOpAND, InstOpJAM, InstOpRLA, InstOpBIT, InstOpAND, InstOpROL, InstOpRLA, // 20
		InstOpJAM, InstOpJAM, InstOpJAM, InstOpJAM, InstOpBIT, InstOpAND, InstOpROL, InstOpRLA,
		InstOpJAM, InstOpAND, InstOpJAM, InstOpRLA, InstOaNOP, InstOpAND, InstOpROL, InstOpRLA, // 30
		InstOpJAM, InstOpAND, InstOpJAM, InstOpRLA, InstOaNOP, InstOpAND, InstOpROL, InstOpRLA,
		InstOpJAM, InstOpEOR, InstOpJAM, InstOpSRE, InstOaNOP, InstOpEOR, InstOpLSR, InstOpSRE, // 40
		InstOpJAM, InstOpJAM, InstOpJAM, InstOpJAM, InstOpJAM, InstOpEOR, InstOpLSR, InstOpSRE,
		InstOpJAM, InstOpEOR, InstOpJAM, InstOpSRE, InstOaNOP, InstOpEOR, InstOpLSR, InstOpSRE, // 50
		InstOpJAM, InstOpEOR, InstOpJAM, InstOpSRE, InstOaNOP, InstOpEOR, InstOpLSR, InstOpSRE,
		InstOpJAM, InstOpADC, InstOpJAM, InstOpRRA, InstOaNOP, InstOpADC, InstOpROR, InstOpRRA, // 60
		InstOpJAM, InstOpJAM, InstOpJAM, InstOpJAM, InstOiJMP, InstOpADC, InstOpROR, InstOpRRA,
		InstOpJAM, InstOpADC, InstOpJAM, InstOpRRA, InstOaNOP, InstOpADC, InstOpROR, InstOpRRA, // 70
		InstOpJAM, InstOpADC, InstOpJAM, InstOpRRA, InstOaNOP, InstOpADC, InstOpROR, InstOpRRA,
		InstOpJAM, InstOpSTA, InstOpJAM, InstOpSAX, InstOpSTY, InstOpSTA, InstOpSTX, InstOpSAX, // 80
		InstOpJAM, InstOpJAM, InstOpJAM, InstOpJAM, InstOpSTY, InstOpSTA, InstOpSTX, InstOpSAX,
		InstOpJAM, InstOpSTA, InstOpJAM, InstOpSHA, InstOpSTY, InstOpSTA, InstOpSTX, InstOpSAX, // 90
		InstOpJAM, InstOpSTA, InstOpJAM, InstOpSHS, InstOpSHY, InstOpSTA, InstOpSHX, InstOpSHA,
		InstOpJAM, InstOpLDA, InstOpJAM, InstOpLAX, InstOpLDY, InstOpLDA, InstOpLDX, InstOpLAX, // a0
		InstOpJAM, InstOpJAM, InstOpJAM, InstOpJAM, InstOpLDY, InstOpLDA, InstOpLDX, InstOpLAX,
		InstOpJAM, InstOpLDA, InstOpJAM, InstOpLAX, InstOpLDY, InstOpLDA, InstOpLDX, InstOpLAX, // b0
		InstOpJAM, InstOpLDA, InstOpJAM, InstOpLAS, InstOpLDY, InstOpLDA, InstOpLDX, InstOpLAX,
		InstOpJAM, InstOpCMP, InstOpJAM, InstOpDCP, InstOpCPY, InstOpCMP, InstOpDEC, InstOpDCP, // c0
		InstOpJAM, InstOpJAM, InstOpJAM, InstOpJAM, InstOpCPY, InstOpCMP, InstOpDEC, InstOpDCP,
		InstOpJAM, InstOpCMP, InstOpJAM, InstOpDCP, InstOaNOP, InstOpCMP, InstOpDEC, InstOpDCP, // d0
		InstOpJAM, InstOpCMP, InstOpJAM, InstOpDCP, InstOaNOP, InstOpCMP, InstOpDEC, InstOpDCP,
		InstOpJAM, InstOpSBC, InstOpJAM, InstOpISB, InstOpCPX, InstOpSBC, InstOpINC, InstOpISB, // e0
		InstOpJAM, InstOpJAM, InstOpJAM, InstOpJAM, InstOpCPX, InstOpSBC, InstOpINC, InstOpISB,
		InstOpJAM, InstOpSBC, InstOpJAM, InstOpISB, InstOaNOP, InstOpSBC, InstOpINC, InstOpISB, // f0
		InstOpJAM, InstOpSBC, InstOpJAM, InstOpISB, InstOaNOP, InstOpSBC, InstOpINC, InstOpISB,
	}
	return opTable
}
