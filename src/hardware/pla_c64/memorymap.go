package pla_c64

// RAM represents the memory type or section with a hexadecimal identifier of 0x0.
const RAM = 0x0

// KER represents the Kernal ROM memory configuration, used to map specific memory regions to the Kernal ROM.
const KER = 0x1

// BAS represents a memory bank identifier for accessing the BASIC ROM in the system's memory configuration.
const BAS = 0x2

// CHA represents the memory configuration for character memory mapping within the system.
const CHA = 0x3

// I_O represents a memory configuration value indicating an I/O port mapping for specific address ranges.
const I_O = 0x4

// ROL represents the ROM_LO memory type, typically used in memory mappings for addressing cartridge ROM.
const ROL = 0x5

// ROH is a constant representing the high ROM segment in memory configurations with the value 0x6.
const ROH = 0x6

// UND represents an undefined or invalid memory type with a constant value of 0xff.
const UND = 0xff

//https://sta.c64.org/cbm64mem.html
//https://www.c64-wiki.com/wiki/Bank_Switching
//%x00: RAM visible in all three areas.
//%x01: RAM visible at 0xA000-0xBFFF and 0xE000-0xFFFF.
//%x10: RAM visible at 0xA000-0xBFFF; KERNAL ROM visible at 0xE000-0xFFFF.
//%x11: BASIC ROM visible at $A000-0xBFFF; KERNAL ROM visible at 0xE000-0xFFFF.
//%0xx: Character ROM visible at $D000-0xDFFF. (Except for the value 0x000, see above.)
//%1xx: I/O area visible at $D000-0xDFFF. (Except for the value 0x100, see above.)

// _memoryMap defines a 2D array representing the system's memory layout, mapping memory regions to specific access types.
var _memoryMap = [][]uint8{
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

// MemoryMap represents a structure for managing memory configurations with a maximum index for memory mapping.
type MemoryMap struct {
	max uint8
}

// NewMemoryMap creates and initializes a new MemoryMap instance based on the `_memoryMap` configuration.
// It panics if `_memoryMap` is empty or its length is not a multiple of 2.
func NewMemoryMap() *MemoryMap {
	memLen := len(_memoryMap)
	if memLen == 0 {
		panic("wrong memory map (len == 0)")
	}
	if memLen%2 != 0 {
		panic("wrong memory map (len isn't multiple of 2)")
	}
	return &MemoryMap{
		max: uint8(memLen - 1),
	}
}

// Get fetches a memory configuration based on the provided memConfig by masking it with the maximum valid value.
func (m *MemoryMap) Get(memConfig uint8) []byte {
	idx := memConfig & m.max //0x1f
	return _memoryMap[idx]
}
