package sam

type Token int

const (
	T_NULL Token = iota
	T_END
	T_LPAREN
	T_RPAREN
	T_ADD
	T_SUB
	T_MUL
	T_DIV
	T_COMMA
	T_IMMED
	T_X
	T_Y
	T_PC
	T_SP
	T_NUMBER
)
