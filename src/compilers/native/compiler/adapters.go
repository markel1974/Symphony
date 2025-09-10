package compiler

import (
	"go/token"

	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
	"github.com/markel1974/c64emu/src/vm/sequencers/native"
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
	token.EQL:  NewTokenAdapter(native.OpLogicalId, []int{int(objects.OperatorLogicalEq)}),
	token.NEQ:  NewTokenAdapter(native.OpLogicalId, []int{int(objects.OperatorLogicalNotEq)}),
	token.LOR:  NewTokenAdapter(native.OpLogicalId, []int{int(objects.OperatorLogicalOr)}),
	token.LAND: NewTokenAdapter(native.OpLogicalId, []int{int(objects.OperatorLogicalAnd)}),
	token.LSS:  NewTokenAdapter(native.OpLogicalId, []int{int(objects.OperatorLogicalLess)}),
	token.GTR:  NewTokenAdapter(native.OpLogicalId, []int{int(objects.OperatorLogicalGreater)}),
	token.GEQ:  NewTokenAdapter(native.OpLogicalId, []int{int(objects.OperatorLogicalGreaterEq)}),
	token.LEQ:  NewTokenAdapter(native.OpLogicalId, []int{int(objects.OperatorLogicalLessEq)}),

	//arithmetic operators
	token.ADD:     NewTokenAdapter(native.OpArithmeticId, []int{int(objects.OperatorAdd)}),
	token.SUB:     NewTokenAdapter(native.OpArithmeticId, []int{int(objects.OperatorSub)}),
	token.MUL:     NewTokenAdapter(native.OpArithmeticId, []int{int(objects.OperatorMul)}),
	token.QUO:     NewTokenAdapter(native.OpArithmeticId, []int{int(objects.OperatorQuo)}),
	token.AND:     NewTokenAdapter(native.OpArithmeticId, []int{int(objects.OperatorAnd)}),
	token.AND_NOT: NewTokenAdapter(native.OpArithmeticId, []int{int(objects.OperatorAndNot)}),
	token.OR:      NewTokenAdapter(native.OpArithmeticId, []int{int(objects.OperatorOr)}),
	token.XOR:     NewTokenAdapter(native.OpArithmeticId, []int{int(objects.OperatorXor)}),
	token.REM:     NewTokenAdapter(native.OpArithmeticId, []int{int(objects.OperatorRem)}),
	token.SHL:     NewTokenAdapter(native.OpArithmeticId, []int{int(objects.OperatorShl)}),
	token.SHR:     NewTokenAdapter(native.OpArithmeticId, []int{int(objects.OperatorShr)}),
}

// _unaryAdapter maps unary token types to their corresponding TokenAdapter configurations for bytecode operations.
var _unaryAdapter = map[token.Token]*TokenAdapter{
	token.ADD: NewTokenAdapter(native.OpUnaryAddId, nil),
	token.SUB: NewTokenAdapter(native.OpUnarySubId, nil),
	token.NOT: NewTokenAdapter(native.OpUnaryNotId, nil),
	token.XOR: NewTokenAdapter(native.OpUnaryBitwiseComplementId, nil),
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
