package objects

// Operator represents an enumerated type for various binary operation identifiers.
type Operator int

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
// OperatorLess represents the less than operator (<).
// OperatorGreater represents the greater than operator (>).
// OperatorLessEq represents the less than or equal to operator (<=).
// OperatorGreaterEq represents the greater than or equal to operator (>=).
const (
	OperatorAdd     Operator = iota // +
	OperatorSub                     // -
	OperatorMul                     // *
	OperatorQuo                     // /
	OperatorRem                     // %
	OperatorAnd                     // &
	OperatorAndNot                  // &^
	OperatorOr                      // |
	OperatorXor                     // ^
	OperatorShl                     // <<
	OperatorShr                     // >>
	OperatorLess                    // <
	OperatorGreater                 // >
	OperatorLessEq
	OperatorGreaterEq
)
