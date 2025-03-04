package sam

type Mnemonic int

// Mnemonics
const (
	M_ADC Mnemonic = iota
	M_AND
	M_ASL
	M_BCC
	M_BCS
	M_BEQ
	M_BIT
	M_BMI
	M_BNE
	M_BPL
	M_BRK
	M_BVC
	M_BVS
	M_CLC
	M_CLD
	M_CLI
	M_CLV
	M_CMP
	M_CPX
	M_CPY

	M_DEC
	M_DEX
	M_DEY
	M_EOR
	M_INC
	M_INX
	M_INY
	M_JMP
	M_JSR
	M_LDA

	M_LDX
	M_LDY
	M_LSR
	M_NOP
	M_ORA
	M_PHA
	M_PHP
	M_PLA
	M_PLP
	M_ROL

	M_ROR
	M_RTI
	M_RTS
	M_SBC
	M_SEC
	M_SED
	M_SEI
	M_STA
	M_STX
	M_STY

	M_TAX
	M_TAY
	M_TSX
	M_TXA
	M_TXS
	M_TYA

	M_ILLEGAL
	// Undocumented opcodes start here

	M_IANC
	M_IANE
	M_IARR
	M_IASR
	M_IDCP
	M_IISB
	M_IJAM
	M_INOP
	M_ILAS
	M_ILAX
	M_ILXA
	M_IRLA
	M_IRRA
	M_ISAX
	M_ISBC
	M_ISBX
	M_ISHA
	M_ISHS

	M_ISHX
	M_ISHY
	M_ISLO
	M_ISRE
	M_MAXIMUM // Highest element
)

// Mnemonic for each opcode
var mnemonic = [256]Mnemonic{
	M_BRK, M_ORA, M_IJAM, M_ISLO, M_INOP, M_ORA, M_ASL, M_ISLO, // 00
	M_PHP, M_ORA, M_ASL, M_IANC, M_INOP, M_ORA, M_ASL, M_ISLO,
	M_BPL, M_ORA, M_IJAM, M_ISLO, M_INOP, M_ORA, M_ASL, M_ISLO, // 10
	M_CLC, M_ORA, M_INOP, M_ISLO, M_INOP, M_ORA, M_ASL, M_ISLO,
	M_JSR, M_AND, M_IJAM, M_IRLA, M_BIT, M_AND, M_ROL, M_IRLA, // 20
	M_PLP, M_AND, M_ROL, M_IANC, M_BIT, M_AND, M_ROL, M_IRLA,
	M_BMI, M_AND, M_IJAM, M_IRLA, M_INOP, M_AND, M_ROL, M_IRLA, // 30
	M_SEC, M_AND, M_INOP, M_IRLA, M_INOP, M_AND, M_ROL, M_IRLA,
	M_RTI, M_EOR, M_IJAM, M_ISRE, M_INOP, M_EOR, M_LSR, M_ISRE, // 40
	M_PHA, M_EOR, M_LSR, M_IASR, M_JMP, M_EOR, M_LSR, M_ISRE,
	M_BVC, M_EOR, M_IJAM, M_ISRE, M_INOP, M_EOR, M_LSR, M_ISRE, // 50
	M_CLI, M_EOR, M_INOP, M_ISRE, M_INOP, M_EOR, M_LSR, M_ISRE,
	M_RTS, M_ADC, M_IJAM, M_IRRA, M_INOP, M_ADC, M_ROR, M_IRRA, // 60
	M_PLA, M_ADC, M_ROR, M_IARR, M_JMP, M_ADC, M_ROR, M_IRRA,
	M_BVS, M_ADC, M_IJAM, M_IRRA, M_INOP, M_ADC, M_ROR, M_IRRA, // 70
	M_SEI, M_ADC, M_INOP, M_IRRA, M_INOP, M_ADC, M_ROR, M_IRRA,
	M_INOP, M_STA, M_INOP, M_ISAX, M_STY, M_STA, M_STX, M_ISAX, // 80
	M_DEY, M_INOP, M_TXA, M_IANE, M_STY, M_STA, M_STX, M_ISAX,
	M_BCC, M_STA, M_IJAM, M_ISHA, M_STY, M_STA, M_STX, M_ISAX, // 90
	M_TYA, M_STA, M_TXS, M_ISHS, M_ISHY, M_STA, M_ISHX, M_ISHA,
	M_LDY, M_LDA, M_LDX, M_ILAX, M_LDY, M_LDA, M_LDX, M_ILAX, // a0
	M_TAY, M_LDA, M_TAX, M_ILXA, M_LDY, M_LDA, M_LDX, M_ILAX,
	M_BCS, M_LDA, M_IJAM, M_ILAX, M_LDY, M_LDA, M_LDX, M_ILAX, // b0
	M_CLV, M_LDA, M_TSX, M_ILAS, M_LDY, M_LDA, M_LDX, M_ILAX,
	M_CPY, M_CMP, M_INOP, M_IDCP, M_CPY, M_CMP, M_DEC, M_IDCP, // c0
	M_INY, M_CMP, M_DEX, M_ISBX, M_CPY, M_CMP, M_DEC, M_IDCP,
	M_BNE, M_CMP, M_IJAM, M_IDCP, M_INOP, M_CMP, M_DEC, M_IDCP, // d0
	M_CLD, M_CMP, M_INOP, M_IDCP, M_INOP, M_CMP, M_DEC, M_IDCP,
	M_CPX, M_SBC, M_INOP, M_IISB, M_CPX, M_SBC, M_INC, M_IISB, // e0
	M_INX, M_SBC, M_NOP, M_ISBC, M_CPX, M_SBC, M_INC, M_IISB,
	M_BEQ, M_SBC, M_IJAM, M_IISB, M_INOP, M_SBC, M_INC, M_IISB, // f0
	M_SED, M_SBC, M_INOP, M_IISB, M_INOP, M_SBC, M_INC, M_IISB,
}

// Chars for each mnemonic
var _mnem_ = []byte("aaabbbbbbbbbbcccccccdddeiiijjllllnopppprrrrssssssstttttt?aaaadijnlllrrsssssssss")
var _mnem_2 = []byte("dnscceimnprvvllllmppeeeonnnmsdddsorhhlloottbeeetttaasxxy?nnrscsaoaaxlrabbhhhhlr")
var _mnem_3 = []byte("cdlcsqtielkcscdivpxycxyrcxypraxyrpaapaplrisccdiaxyxyxasa?cerrpbmpsxaaaxcxasxyoe")
