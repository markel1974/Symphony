package objects

type Operator int

const (
	OperatorIllegal Operator = iota
	OperatorInt
	OperatorString
	OperatorAdd     // +
	OperatorSub     // -
	OperatorMul     // *
	OperatorQuo     // /
	OperatorRem     // %
	OperatorAnd     // &
	OperatorAndNot  // &^
	OperatorOr      // |
	OperatorXor     // ^
	OperatorShl     // <<
	OperatorShr     // >>
	OperatorLess    // <
	OperatorGreater // >
	OperatorLessEq
	OperatorGreaterEq
)
