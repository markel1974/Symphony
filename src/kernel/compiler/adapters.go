package compiler

import (
	"go/token"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// TokenAdapter represents a mapping between a token and its corresponding bytecode operation with optional arguments.
type TokenAdapter struct {
	op        bytecode.Opcode
	arguments []int
}

// NewTokenAdapter creates a new TokenAdapter instance with the specified opcode and argument list.
func NewTokenAdapter(op bytecode.Opcode, arguments []int) *TokenAdapter {
	return &TokenAdapter{
		op:        op,
		arguments: arguments,
	}
}

// _binaryAdapter maps token types to their corresponding TokenAdapter for processing binary operations.
var _binaryAdapter = map[token.Token]*TokenAdapter{
	token.EQL:     NewTokenAdapter(bytecode.OpEqual, nil),
	token.NEQ:     NewTokenAdapter(bytecode.OpNotEqual, nil),
	token.ADD:     NewTokenAdapter(bytecode.OpBinary, []int{int(objects.OperatorAdd)}),
	token.SUB:     NewTokenAdapter(bytecode.OpBinary, []int{int(objects.OperatorSub)}),
	token.MUL:     NewTokenAdapter(bytecode.OpBinary, []int{int(objects.OperatorMul)}),
	token.QUO:     NewTokenAdapter(bytecode.OpBinary, []int{int(objects.OperatorQuo)}),
	token.LSS:     NewTokenAdapter(bytecode.OpBinary, []int{int(objects.OperatorLess)}),
	token.GTR:     NewTokenAdapter(bytecode.OpBinary, []int{int(objects.OperatorGreater)}),
	token.GEQ:     NewTokenAdapter(bytecode.OpBinary, []int{int(objects.OperatorGreaterEq)}),
	token.LEQ:     NewTokenAdapter(bytecode.OpBinary, []int{int(objects.OperatorLessEq)}),
	token.AND:     NewTokenAdapter(bytecode.OpBinary, []int{int(objects.OperatorAnd)}),
	token.AND_NOT: NewTokenAdapter(bytecode.OpBinary, []int{int(objects.OperatorAndNot)}),
	token.OR:      NewTokenAdapter(bytecode.OpBinary, []int{int(objects.OperatorOr)}),
	token.XOR:     NewTokenAdapter(bytecode.OpBinary, []int{int(objects.OperatorXor)}),
	token.REM:     NewTokenAdapter(bytecode.OpBinary, []int{int(objects.OperatorRem)}),
	token.SHL:     NewTokenAdapter(bytecode.OpBinary, []int{int(objects.OperatorShl)}),
	token.SHR:     NewTokenAdapter(bytecode.OpBinary, []int{int(objects.OperatorShr)}),
}

// _unaryAdapter maps unary token types to their corresponding TokenAdapter configurations for bytecode operations.
var _unaryAdapter = map[token.Token]*TokenAdapter{
	token.SUB: NewTokenAdapter(bytecode.OpMinus, nil),
	token.NOT: NewTokenAdapter(bytecode.OpLNot, nil),
}

// BinaryAdapterFor retrieves the TokenAdapter for the given token operator and indicates if it exists in the mapping.
func BinaryAdapterFor(op token.Token) (*TokenAdapter, bool) {
	v, ok := _binaryAdapter[op]
	return v, ok
}

// UnaryAdapterFor retrieves the TokenAdapter for the given token if it exists in the unary adapter map.
// Returns the TokenAdapter and a boolean indicating if a match was found.
func UnaryAdapterFor(op token.Token) (*TokenAdapter, bool) {
	v, ok := _unaryAdapter[op]
	return v, ok
}
