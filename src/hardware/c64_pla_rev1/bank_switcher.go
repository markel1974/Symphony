package c64_pla_rev1

import "log"

// see C64 Bank Switching
// https://www.c64-wiki.com/wiki/Bank_Switching
// https://www.c64-wiki.com/wiki/Memory_Map#Configurations
// https://sta.c64.org/cbm64mem.html
// https://codebase64.org/doku.php?id=base:memory_management

const (
	// RAM represents the memory type or section with a hexadecimal identifier of 0x0.
	RAM = 0x0
	// KER represents the Kernal ROM memory configuration, used to map specific memory regions to the Kernal ROM.
	KER = 0x1
	// BAS represents a memory bank identifier for accessing the BASIC ROM in the system's memory configuration.
	BAS = 0x2
	// CHA represents the memory configuration for character memory mapping within the system.
	CHA = 0x3
	// I_O represents a memory configuration value indicating an I/O port mapping for specific address ranges.
	I_O = 0x4
	// ROL represents the ROM_LO memory type, typically used in memory mappings for addressing cartridge ROM.
	ROL = 0x5
	// ROH is a constant representing the high ROM segment in memory configurations with the value 0x6.
	ROH = 0x6
	// UND represents an undefined or invalid memory type with a constant value of 0xff.
	UND = 0xff
)

// _bankSwitching defines a 2D array representing the system's memory layout, mapping memory regions to specific access types.
var _bankSwitching = [][]uint8{
	//	               0    1    2    3    4    5    6    7    8    9    A    B    C    D    E    F
	/*  0 - 00000 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM},
	/*  1 - 00001 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM},
	/*  2 - 00010 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, ROH, ROH, RAM, CHA, KER, KER},
	/*  3 - 00011 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, ROL, ROL, ROH, ROH, RAM, CHA, KER, KER},
	/*  4 - 00100 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM},
	/*  5 - 00101 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, I_O, RAM, RAM},
	/*  6 - 00110 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, ROH, ROH, RAM, I_O, KER, KER},
	/*  7 - 00111 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, ROL, ROL, ROH, ROH, RAM, I_O, KER, KER},
	/*  8 - 01000 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM},
	/*  9 - 01001 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, CHA, RAM, RAM},
	/* 10 - 01010 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, CHA, KER, KER},
	/* 11 - 01011 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, ROL, ROL, BAS, BAS, RAM, CHA, KER, KER},
	/* 12 - 01100 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM},
	/* 13 - 01101 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, I_O, RAM, RAM},
	/* 14 - 01110 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, I_O, KER, KER},
	/* 15 - 01111 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, ROL, ROL, BAS, BAS, RAM, I_O, KER, KER},
	/* 16 - 10000 */ {RAM, UND, UND, UND, UND, UND, UND, UND, ROL, ROL, UND, UND, UND, I_O, ROH, ROH},
	/* 17 - 10001 */ {RAM, UND, UND, UND, UND, UND, UND, UND, ROL, ROL, UND, UND, UND, I_O, ROH, ROH},
	/* 18 - 10010 */ {RAM, UND, UND, UND, UND, UND, UND, UND, ROL, ROL, UND, UND, UND, I_O, ROH, ROH},
	/* 19 - 10011 */ {RAM, UND, UND, UND, UND, UND, UND, UND, ROL, ROL, UND, UND, UND, I_O, ROH, ROH},
	/* 20 - 10100 */ {RAM, UND, UND, UND, UND, UND, UND, UND, ROL, ROL, UND, UND, UND, I_O, ROH, ROH},
	/* 21 - 10101 */ {RAM, UND, UND, UND, UND, UND, UND, UND, ROL, ROL, UND, UND, UND, I_O, ROH, ROH},
	/* 22 - 10110 */ {RAM, UND, UND, UND, UND, UND, UND, UND, ROL, ROL, UND, UND, UND, I_O, ROH, ROH},
	/* 23 - 10111 */ {RAM, UND, UND, UND, UND, UND, UND, UND, ROL, ROL, UND, UND, UND, I_O, ROH, ROH},
	/* 24 - 11000 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM},
	/* 25 - 11001 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, CHA, RAM, RAM},
	/* 26 - 11010 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, CHA, KER, KER},
	/* 27 - 11011 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, BAS, BAS, RAM, CHA, KER, KER},
	/* 28 - 11100 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM},
	/* 29 - 11101 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, I_O, RAM, RAM},
	/* 30 - 11110 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, I_O, KER, KER},
	/* 31 - 11111 */ {RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, RAM, BAS, BAS, RAM, I_O, KER, KER},
}

// BankSwitcher represents a structure for managing memory configurations with a maximum index for memory mapping.
type BankSwitcher struct {
	mask  uint8
	index int
}

// NewBankSwitcher creates and initializes a new BankSwitcher instance based on the `_bankSwitching` configuration.
// It panics if `_bankSwitching` is empty or its length is not a multiple of 2.
func NewBankSwitcher() *BankSwitcher {
	bsLen := len(_bankSwitching)
	if bsLen == 0 {
		log.Fatal("wrong memory map (len == 0)")
	}
	if bsLen%2 != 0 {
		log.Fatal("wrong memory map (len isn't multiple of 2)")
	}
	return &BankSwitcher{
		mask:  uint8(bsLen - 1),
		index: -1,
	}
}

// GetIndex returns the current memory configuration index managed by the BankSwitcher instance.
func (m *BankSwitcher) GetIndex() int {
	return m.index
}

// Apply updates the memory configuration index and returns the mode table entry and true if the index has changed.
func (m *BankSwitcher) Apply(memConfig int) ([]byte, bool) {
	idx := uint8(memConfig) & m.mask
	if m.index == int(idx) {
		return nil, false
	}
	m.index = int(idx)
	return _bankSwitching[idx], true
}
