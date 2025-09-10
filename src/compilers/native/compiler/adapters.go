package compiler

import (
	"go/token"

	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// TokenAdapter represents a mapping between a token and its corresponding bytecode operation with optional arguments.
type TokenAdapter struct {
	op        opcodes.OpcodeId
	arguments []int
}

// NewTokenAdapter creates a new TokenAdapter instance with the specified opcode and argument list.
func NewTokenAdapter(op opcodes.OpcodeId, arguments []int) *TokenAdapter {
	return &TokenAdapter{
		op:        op,
		arguments: arguments,
	}
}

// _binaryAdapter maps token types to their corresponding TokenAdapter for processing binary operations.
var _binaryAdapter = map[token.Token]*TokenAdapter{
	//logical operators
	token.EQL:  NewTokenAdapter(opcodes.OpLogical, []int{int(objects.OperatorLogicalEq)}),
	token.NEQ:  NewTokenAdapter(opcodes.OpLogical, []int{int(objects.OperatorLogicalNotEq)}),
	token.LOR:  NewTokenAdapter(opcodes.OpLogical, []int{int(objects.OperatorLogicalOr)}),
	token.LAND: NewTokenAdapter(opcodes.OpLogical, []int{int(objects.OperatorLogicalAnd)}),
	token.LSS:  NewTokenAdapter(opcodes.OpLogical, []int{int(objects.OperatorLogicalLess)}),
	token.GTR:  NewTokenAdapter(opcodes.OpLogical, []int{int(objects.OperatorLogicalGreater)}),
	token.GEQ:  NewTokenAdapter(opcodes.OpLogical, []int{int(objects.OperatorLogicalGreaterEq)}),
	token.LEQ:  NewTokenAdapter(opcodes.OpLogical, []int{int(objects.OperatorLogicalLessEq)}),

	//arithmetic operators
	token.ADD:     NewTokenAdapter(opcodes.OpArithmetic, []int{int(objects.OperatorAdd)}),
	token.SUB:     NewTokenAdapter(opcodes.OpArithmetic, []int{int(objects.OperatorSub)}),
	token.MUL:     NewTokenAdapter(opcodes.OpArithmetic, []int{int(objects.OperatorMul)}),
	token.QUO:     NewTokenAdapter(opcodes.OpArithmetic, []int{int(objects.OperatorQuo)}),
	token.AND:     NewTokenAdapter(opcodes.OpArithmetic, []int{int(objects.OperatorAnd)}),
	token.AND_NOT: NewTokenAdapter(opcodes.OpArithmetic, []int{int(objects.OperatorAndNot)}),
	token.OR:      NewTokenAdapter(opcodes.OpArithmetic, []int{int(objects.OperatorOr)}),
	token.XOR:     NewTokenAdapter(opcodes.OpArithmetic, []int{int(objects.OperatorXor)}),
	token.REM:     NewTokenAdapter(opcodes.OpArithmetic, []int{int(objects.OperatorRem)}),
	token.SHL:     NewTokenAdapter(opcodes.OpArithmetic, []int{int(objects.OperatorShl)}),
	token.SHR:     NewTokenAdapter(opcodes.OpArithmetic, []int{int(objects.OperatorShr)}),
}

// _unaryAdapter maps unary token types to their corresponding TokenAdapter configurations for bytecode operations.
var _unaryAdapter = map[token.Token]*TokenAdapter{
	token.ADD: NewTokenAdapter(opcodes.OpUnaryAdd, nil),
	token.SUB: NewTokenAdapter(opcodes.OpUnarySub, nil),
	token.NOT: NewTokenAdapter(opcodes.OpUnaryNot, nil),
	token.XOR: NewTokenAdapter(opcodes.OpUnaryBitwiseComplement, nil),
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
