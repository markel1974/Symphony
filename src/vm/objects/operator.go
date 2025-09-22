package objects

// ArithmeticOperator represents a numeric type used to define arithmetic operators like addition, subtraction, etc.
type ArithmeticOperator int

// LogicalOperator represents an enumerated type for logical operators.
// It is primarily used to define logical operation constants.
type LogicalOperator int

// OperatorAdd represents the addition operator (+).
// OperatorSub represents the subtraction operator (-).
// OperatorMul represents the multiplication operator (*).
// OperatorQuo represents the division operator (/).
// OperatorRem represents the remainder operator (%).
// OperatorAnd represents the bitwise AND operator (&).
// OperatorAndNot represents the bitwise AND NOT operator (&^).
// OperatorOr represents the bitwise OR operator (|).
// OperatorXor represents the bitwise XOR operator (^).
// OperatorShl represents the left shift operator (<<).
// OperatorShr represents the right shift operator (>>).
const (
	OperatorAdd    ArithmeticOperator = iota // +
	OperatorSub                              // -
	OperatorMul                              // *
	OperatorQuo                              // /
	OperatorRem                              // %
	OperatorAnd                              // &
	OperatorAndNot                           // &^
	OperatorOr                               // |
	OperatorXor                              // ^
	OperatorShl                              // <<
	OperatorShr                              // >>
)

// OperatorLogicalOr represents the logical OR operator.
// OperatorLogicalAnd represents the logical AND operator.
// OperatorLogicalLess represents the logical LESS THAN (<) operator.
// OperatorLogicalGreater represents the logical GREATER THAN (>) operator.
// OperatorLogicalLessEq represents the logical LESS THAN OR EQUAL TO operator.
// OperatorLogicalGreaterEq represents the logical GREATER THAN OR EQUAL TO operator.
const (
	OperatorLogicalEq LogicalOperator = iota
	OperatorLogicalNotEq
	OperatorLogicalOr
	OperatorLogicalAnd
	OperatorLogicalLess    // <
	OperatorLogicalGreater // >
	OperatorLogicalLessEq
	OperatorLogicalGreaterEq
)

// UnaryOperator represents an enumerated type for unary operators.
type UnaryOperator int

// OperatorUnary... represents the unary operators.
const (
	OperatorUnaryAdd UnaryOperator = iota // +
	OperatorUnarySub                      // -
	OperatorUnaryNot
	OperatorUnaryBitwiseComplement
)
