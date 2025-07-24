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
