package sam

type AddressingMode int

const (
	// Addressing modes

	A_IMPL  AddressingMode = iota
	A_ACCU                 // A
	A_IMM                  // #zz
	A_REL                  // Branches
	A_ZERO                 // zz
	A_ZEROX                // zz,x
	A_ZEROY                // zz,y
	A_ABS                  // zzzz
	A_ABSX                 // zzzz,x
	A_ABSY                 // zzzz,y
	A_IND                  // (zzzz)
	A_INDX                 // (zz,x)
	A_INDY                 // (zz),y
)

// Addressing mode for each opcode
var adr_mode = [256]AddressingMode{
	A_IMPL, A_INDX, A_IMPL, A_INDX, A_ZERO, A_ZERO, A_ZERO, A_ZERO, // 00
	A_IMPL, A_IMM, A_ACCU, A_IMM, A_ABS, A_ABS, A_ABS, A_ABS,
	A_REL, A_INDY, A_IMPL, A_INDY, A_ZEROX, A_ZEROX, A_ZEROX, A_ZEROX, // 10
	A_IMPL, A_ABSY, A_IMPL, A_ABSY, A_ABSX, A_ABSX, A_ABSX, A_ABSX,
	A_ABS, A_INDX, A_IMPL, A_INDX, A_ZERO, A_ZERO, A_ZERO, A_ZERO, // 20
	A_IMPL, A_IMM, A_ACCU, A_IMM, A_ABS, A_ABS, A_ABS, A_ABS,
	A_REL, A_INDY, A_IMPL, A_INDY, A_ZEROX, A_ZEROX, A_ZEROX, A_ZEROX, // 30
	A_IMPL, A_ABSY, A_IMPL, A_ABSY, A_ABSX, A_ABSX, A_ABSX, A_ABSX,
	A_IMPL, A_INDX, A_IMPL, A_INDX, A_ZERO, A_ZERO, A_ZERO, A_ZERO, // 40
	A_IMPL, A_IMM, A_ACCU, A_IMM, A_ABS, A_ABS, A_ABS, A_ABS,
	A_REL, A_INDY, A_IMPL, A_INDY, A_ZEROX, A_ZEROX, A_ZEROX, A_ZEROX, // 50
	A_IMPL, A_ABSY, A_IMPL, A_ABSY, A_ABSX, A_ABSX, A_ABSX, A_ABSX,
	A_IMPL, A_INDX, A_IMPL, A_INDX, A_ZERO, A_ZERO, A_ZERO, A_ZERO, // 60
	A_IMPL, A_IMM, A_ACCU, A_IMM, A_IND, A_ABS, A_ABS, A_ABS,
	A_REL, A_INDY, A_IMPL, A_INDY, A_ZEROX, A_ZEROX, A_ZEROX, A_ZEROX, // 70
	A_IMPL, A_ABSY, A_IMPL, A_ABSY, A_ABSX, A_ABSX, A_ABSX, A_ABSX,
	A_IMM, A_INDX, A_IMM, A_INDX, A_ZERO, A_ZERO, A_ZERO, A_ZERO, // 80
	A_IMPL, A_IMM, A_IMPL, A_IMM, A_ABS, A_ABS, A_ABS, A_ABS,
	A_REL, A_INDY, A_IMPL, A_INDY, A_ZEROX, A_ZEROX, A_ZEROY, A_ZEROY, // 90
	A_IMPL, A_ABSY, A_IMPL, A_ABSY, A_ABSX, A_ABSX, A_ABSY, A_ABSY,
	A_IMM, A_INDX, A_IMM, A_INDX, A_ZERO, A_ZERO, A_ZERO, A_ZERO, // a0
	A_IMPL, A_IMM, A_IMPL, A_IMM, A_ABS, A_ABS, A_ABS, A_ABS,
	A_REL, A_INDY, A_IMPL, A_INDY, A_ZEROX, A_ZEROX, A_ZEROY, A_ZEROY, // b0
	A_IMPL, A_ABSY, A_IMPL, A_ABSY, A_ABSX, A_ABSX, A_ABSY, A_ABSY,
	A_IMM, A_INDX, A_IMM, A_INDX, A_ZERO, A_ZERO, A_ZERO, A_ZERO, // c0
	A_IMPL, A_IMM, A_IMPL, A_IMM, A_ABS, A_ABS, A_ABS, A_ABS,
	A_REL, A_INDY, A_IMPL, A_INDY, A_ZEROX, A_ZEROX, A_ZEROX, A_ZEROX, // d0
	A_IMPL, A_ABSY, A_IMPL, A_ABSY, A_ABSX, A_ABSX, A_ABSX, A_ABSX,
	A_IMM, A_INDX, A_IMM, A_INDX, A_ZERO, A_ZERO, A_ZERO, A_ZERO, // e0
	A_IMPL, A_IMM, A_IMPL, A_IMM, A_ABS, A_ABS, A_ABS, A_ABS,
	A_REL, A_INDY, A_IMPL, A_INDY, A_ZEROX, A_ZEROX, A_ZEROX, A_ZEROX, // f0
	A_IMPL, A_ABSY, A_IMPL, A_ABSY, A_ABSX, A_ABSX, A_ABSX, A_ABSX,
}

// Instruction length for each addressing mode
var _adr_length = []byte{1, 1, 2, 2, 2, 2, 2, 3, 3, 3, 3, 2, 2}
