package mos6510_rev1

//https://www.c64-wiki.com/wiki/Opcode

//ORA (short for "Logical OR on Accumulator")
//ASL (short for "Arithmetic Shift Left")

// Ap -> Operand effective address, no extra cycles
// Ae -> Operand effective address, extra cycle on page crossing
// Mp -> Operand operand and write it back (for RMW Instructions), no extra cycles
// Oi -> Operation Immediate/Indirect
// Oa -> Operation Accumulator
// Op -> Operation
