package mos6510_rev1

// _modeTable is a slice of functions, each defining behavior or processing logic for specific CPU operational modes.
var _modeTable []func(*CPU)

// _opTable defines a slice of CPU instruction implementations, where each function specifies operation logic for an opcode.
var _opTable []func(*CPU)

//https://www.c64-wiki.com/wiki/Opcode

//ORA (short for "Logical OR on Accumulator")
//ASL (short for "Arithmetic Shift Left")

// Ap -> Read effective address, no extra cycles
// Ae -> Read effective address, extra cycle on page crossing
// Mp -> Read operand and write it back (for RMW instructions), no extra cycles
// Oi -> Operation Immediate/Indirect
// Oa -> Operation Accumulator
// Op -> Operation

// init initializes the CPU operation and mode tables by assigning corresponding instruction functions to their positions.
func init() {
	_modeTable = []func(*CPU){
		instOpBRK, instApINDx, instOpJAM, instMpINDx, instApZER, instApZER, instMpZER, instMpZER, // 00
		instOpPHP, instOiOPA, instOaASL, instOiANC, instApABS, instApABS, instMpABS, instMpABS,
		instOpBPL, instAeINDy, instOpJAM, instMpINDy, instApZERx, instApZERx, instMpZERx, instMpZERx, // 10
		instOpCLC, instAeABSy, instOpNOP, instMpABSy, instAeABSx, instAeABSx, instMpABSx, instMpABSx,
		instOpJSR, instApINDx, instOpJAM, instMpINDx, instApZER, instApZER, instMpZER, instMpZER, // 20
		instOpPLP, instOiAND, instOaROL, instOiANC, instApABS, instApABS, instMpABS, instMpABS,
		instOpBMI, instAeINDy, instOpJAM, instMpINDy, instApZERx, instApZERx, instMpZERx, instMpZERx, // 30
		instOpSEC, instAeABSy, instOpNOP, instMpABSy, instAeABSx, instAeABSx, instMpABSx, instMpABSx,
		instOpRTI, instApINDx, instOpJAM, instMpINDx, instApZER, instApZER, instMpZER, instMpZER, // 40
		instOpPHA, instOiEOR, instOaLSR, instOiASR, instOpJMP, instApABS, instMpABS, instMpABS,
		instOpBVC, instAeINDy, instOpJAM, instMpINDy, instApZERx, instApZERx, instMpZERx, instMpZERx, // 50
		instOpCLI, instAeABSy, instOpNOP, instMpABSy, instAeABSx, instAeABSx, instMpABSx, instMpABSx,
		instOpRTS, instApINDx, instOpJAM, instMpINDx, instApZER, instApZER, instMpZER, instMpZER, // 60
		instOpPLA, instOiADC, instOaROR, instOiARR, instApABS, instApABS, instMpABS, instMpABS,
		instOpBVS, instAeINDy, instOpJAM, instMpINDy, instApZERx, instApZERx, instMpZERx, instMpZERx, // 70
		instOpSEI, instAeABSy, instOpNOP, instMpABSy, instAeABSx, instAeABSx, instMpABSx, instMpABSx,
		instOiNOP, instApINDx, instOiNOP, instApINDx, instApZER, instApZER, instApZER, instApZER, // 80
		instOpDEY, instOiNOP, instOpTXA, instOiANE, instApABS, instApABS, instApABS, instApABS,
		instOpBCC, instApINDy, instOpJAM, instApINDy, instApZERx, instApZERx, instApZERy, instApZERy, // 90
		instOpTYA, instApABSy, instOpTXS, instApABSy, instApABSx, instApABSx, instApABSy, instApABSy,
		instOiLDY, instApINDx, instOiLDX, instApINDx, instApZER, instApZER, instApZER, instApZER, // a0
		instOpTAY, instOiLDA, instOpTAX, instOiLXA, instApABS, instApABS, instApABS, instApABS,
		instOpBCS, instAeINDy, instOpJAM, instAeINDy, instApZERx, instApZERx, instApZERy, instApZERy, // b0
		instOpCLV, instAeABSy, instOpTSX, instAeABSy, instAeABSx, instAeABSx, instAeABSy, instAeABSy,
		instOiCPY, instApINDx, instOiNOP, instMpINDx, instApZER, instApZER, instMpZER, instMpZER, // c0
		instOpINY, instOiCMP, instOpDEX, instOiSBX, instApABS, instApABS, instMpABS, instMpABS,
		instOpBNE, instAeINDy, instOpJAM, instMpINDy, instApZERx, instApZERx, instMpZERx, instMpZERx, // d0
		instOpCLD, instAeABSy, instOpNOP, instMpABSy, instAeABSx, instAeABSx, instMpABSx, instMpABSx,
		instOiCPX, instApINDx, instOiNOP, instMpINDx, instApZER, instApZER, instMpZER, instMpZER, // e0
		instOpINX, instOiSBC, instOpNOP, instOiSBC, instApABS, instApABS, instMpABS, instMpABS,
		instOpBEQ, instAeINDy, instOpJAM, instMpINDy, instApZERx, instApZERx, instMpZERx, instMpZERx, // f0
		instOpSED, instAeABSy, instOpNOP, instMpABSy, instAeABSx, instAeABSx, instMpABSx, instMpABSx,
	}

	_opTable = []func(*CPU){
		instOpJAM, instOpORA, instOpJAM, instOpSLO, instOaNOP, instOpORA, instOpASL, instOpSLO, // 00
		instOpJAM, instOpJAM, instOpJAM, instOpJAM, instOaNOP, instOpORA, instOpASL, instOpSLO,
		instOpJAM, instOpORA, instOpJAM, instOpSLO, instOaNOP, instOpORA, instOpASL, instOpSLO, // 10
		instOpJAM, instOpORA, instOpJAM, instOpSLO, instOaNOP, instOpORA, instOpASL, instOpSLO,
		instOpJAM, instOpAND, instOpJAM, instOpRLA, instOpBIT, instOpAND, instOpROL, instOpRLA, // 20
		instOpJAM, instOpJAM, instOpJAM, instOpJAM, instOpBIT, instOpAND, instOpROL, instOpRLA,
		instOpJAM, instOpAND, instOpJAM, instOpRLA, instOaNOP, instOpAND, instOpROL, instOpRLA, // 30
		instOpJAM, instOpAND, instOpJAM, instOpRLA, instOaNOP, instOpAND, instOpROL, instOpRLA,
		instOpJAM, instOpEOR, instOpJAM, instOpSRE, instOaNOP, instOpEOR, instOpLSR, instOpSRE, // 40
		instOpJAM, instOpJAM, instOpJAM, instOpJAM, instOpJAM, instOpEOR, instOpLSR, instOpSRE,
		instOpJAM, instOpEOR, instOpJAM, instOpSRE, instOaNOP, instOpEOR, instOpLSR, instOpSRE, // 50
		instOpJAM, instOpEOR, instOpJAM, instOpSRE, instOaNOP, instOpEOR, instOpLSR, instOpSRE,
		instOpJAM, instOpADC, instOpJAM, instOpRRA, instOaNOP, instOpADC, instOpROR, instOpRRA, // 60
		instOpJAM, instOpJAM, instOpJAM, instOpJAM, instOiJMP, instOpADC, instOpROR, instOpRRA,
		instOpJAM, instOpADC, instOpJAM, instOpRRA, instOaNOP, instOpADC, instOpROR, instOpRRA, // 70
		instOpJAM, instOpADC, instOpJAM, instOpRRA, instOaNOP, instOpADC, instOpROR, instOpRRA,
		instOpJAM, instOpSTA, instOpJAM, instOpSAX, instOpSTY, instOpSTA, instOpSTX, instOpSAX, // 80
		instOpJAM, instOpJAM, instOpJAM, instOpJAM, instOpSTY, instOpSTA, instOpSTX, instOpSAX,
		instOpJAM, instOpSTA, instOpJAM, instOpSHA, instOpSTY, instOpSTA, instOpSTX, instOpSAX, // 90
		instOpJAM, instOpSTA, instOpJAM, instOpSHS, instOpSHY, instOpSTA, instOpSHX, instOpSHA,
		instOpJAM, instOpLDA, instOpJAM, instOpLAX, instOpLDY, instOpLDA, instOpLDX, instOpLAX, // a0
		instOpJAM, instOpJAM, instOpJAM, instOpJAM, instOpLDY, instOpLDA, instOpLDX, instOpLAX,
		instOpJAM, instOpLDA, instOpJAM, instOpLAX, instOpLDY, instOpLDA, instOpLDX, instOpLAX, // b0
		instOpJAM, instOpLDA, instOpJAM, instOpLAS, instOpLDY, instOpLDA, instOpLDX, instOpLAX,
		instOpJAM, instOpCMP, instOpJAM, instOpDCP, instOpCPY, instOpCMP, instOpDEC, instOpDCP, // c0
		instOpJAM, instOpJAM, instOpJAM, instOpJAM, instOpCPY, instOpCMP, instOpDEC, instOpDCP,
		instOpJAM, instOpCMP, instOpJAM, instOpDCP, instOaNOP, instOpCMP, instOpDEC, instOpDCP, // d0
		instOpJAM, instOpCMP, instOpJAM, instOpDCP, instOaNOP, instOpCMP, instOpDEC, instOpDCP,
		instOpJAM, instOpSBC, instOpJAM, instOpISB, instOpCPX, instOpSBC, instOpINC, instOpISB, // e0
		instOpJAM, instOpJAM, instOpJAM, instOpJAM, instOpCPX, instOpSBC, instOpINC, instOpISB,
		instOpJAM, instOpSBC, instOpJAM, instOpISB, instOaNOP, instOpSBC, instOpINC, instOpISB, // f0
		instOpJAM, instOpSBC, instOpJAM, instOpISB, instOaNOP, instOpSBC, instOpINC, instOpISB,
	}
}
