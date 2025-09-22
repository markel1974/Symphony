package objects

// logicalOpNil evaluates logical operations when the right-hand operand is Nil.
// It supports equality and inequality operators; others return false and ErrInvalidOperator.
func logicalOpNil(op LogicalOperator) (bool, error) {
	switch op {
	case OperatorLogicalEq:
		return false, nil
	case OperatorLogicalNotEq:
		return true, nil
	default:
		return false, ErrInvalidOperator
	}
}

// logicalOpInt64 performs a logical operation between two int64 values based on the specified LogicalOperator.
// Returns a boolean result if the operation is successful, or an error if the operator is invalid.
func logicalOpInt64(lhsValue int64, op LogicalOperator, rhsValue int64) (bool, error) {
	switch op {
	case OperatorLogicalEq:
		return lhsValue == rhsValue, nil
	case OperatorLogicalNotEq:
		return lhsValue != rhsValue, nil
	case OperatorLogicalOr:
		return lhsValue != 0 || rhsValue != 0, nil
	case OperatorLogicalAnd:
		return lhsValue != 0 && rhsValue != 0, nil
	case OperatorLogicalGreater:
		return lhsValue > rhsValue, nil
	case OperatorLogicalGreaterEq:
		return lhsValue >= rhsValue, nil
	case OperatorLogicalLess:
		return lhsValue < rhsValue, nil
	case OperatorLogicalLessEq:
		return lhsValue <= rhsValue, nil
	default:
		return false, ErrInvalidOperator
	}
}

// logicalOpFloat64 performs a logical operation on two float64 values based on the specified LogicalOperator.
// Returns a boolean result or an error if the operator is invalid.
func logicalOpFloat64(lhsValue float64, op LogicalOperator, rhsValue float64) (bool, error) {
	switch op {
	case OperatorLogicalEq:
		return lhsValue == rhsValue, nil
	case OperatorLogicalNotEq:
		return lhsValue != rhsValue, nil
	case OperatorLogicalOr:
		return lhsValue != 0 || rhsValue != 0, nil
	case OperatorLogicalAnd:
		return lhsValue != 0 && rhsValue != 0, nil
	case OperatorLogicalGreater:
		return lhsValue > rhsValue, nil
	case OperatorLogicalGreaterEq:
		return lhsValue >= rhsValue, nil
	case OperatorLogicalLess:
		return lhsValue < rhsValue, nil
	case OperatorLogicalLessEq:
		return lhsValue <= rhsValue, nil
	default:
		return false, ErrInvalidOperator
	}
}

// arithmeticOpInt64 performs an arithmetic or bitwise operation on two int64 values based on the provided operator.
// It returns the result of the operation or an error if the operator is invalid or division by zero is attempted.
func arithmeticOpInt64(lhsValue int64, op ArithmeticOperator, rhsValue int64) (int64, error) {
	switch op {
	case OperatorAdd:
		return lhsValue + rhsValue, nil
	case OperatorSub:
		return lhsValue - rhsValue, nil
	case OperatorMul:
		return lhsValue * rhsValue, nil
	case OperatorQuo:
		if rhsValue == 0 {
			return -1, ErrDivisionByZero
		}
		return lhsValue / rhsValue, nil
	case OperatorRem:
		if rhsValue == 0 {
			return -1, ErrDivisionByZero
		}
		return lhsValue % rhsValue, nil
	case OperatorAnd:
		return lhsValue & rhsValue, nil
	case OperatorOr:
		return lhsValue | rhsValue, nil
	case OperatorXor:
		return lhsValue ^ rhsValue, nil
	case OperatorAndNot:
		return lhsValue &^ rhsValue, nil
	case OperatorShl:
		return lhsValue << rhsValue, nil
	case OperatorShr:
		return lhsValue >> rhsValue, nil
	default:
		return -1, ErrInvalidOperator
	}
}

// arithmeticOpFloat64 performs the specified arithmetic operation on two float64 operands and returns the result.
// Returns an error if the operator is invalid or division by zero is attempted.
func arithmeticOpFloat64(lhsValue float64, op ArithmeticOperator, rhsValue float64) (float64, error) {
	switch op {
	case OperatorAdd:
		return lhsValue + rhsValue, nil
	case OperatorSub:
		return lhsValue - rhsValue, nil
	case OperatorMul:
		return lhsValue * rhsValue, nil
	case OperatorQuo:
		if rhsValue == 0 {
			return -1, ErrDivisionByZero
		}
		return lhsValue / rhsValue, nil
	default:
		return -1, ErrInvalidOperator
	}
}

// unaryOpInt64 applies a unary operator to a 64-bit integer and returns the result or an error if the operator is invalid.
func unaryOpInt64(op UnaryOperator, value int64) (int64, error) {
	switch op {
	case OperatorUnaryAdd:
		return +value, nil
	case OperatorUnarySub:
		return -value, nil
	case OperatorUnaryNot:
		if value == 0 {
			return 1, nil
		}
		return 0, nil
	case OperatorUnaryBitwiseComplement:
		return ^value, nil
	default:
		return -1, ErrInvalidOperator
	}
}

// unaryOpFloat64 performs a unary operation on a float64 value using the specified UnaryOperator. It returns the result or an error.
func unaryOpFloat64(op UnaryOperator, value float64) (float64, error) {
	switch op {
	case OperatorUnaryAdd:
		return +value, nil
	case OperatorUnarySub:
		return -value, nil
	case OperatorUnaryNot:
		if value == 0 {
			return 1, nil
		}
		return 0, nil
	case OperatorUnaryBitwiseComplement:
		return float64(^int64(value)), nil
	default:
		return -1, ErrInvalidOperator
	}
}

// unaryOpBool performs a unary operation on a boolean value using the provided UnaryOperator and returns the result.
// Returns an error if the operator is invalid.
func unaryOpBool(op UnaryOperator, value bool) (bool, error) {
	switch op {
	case OperatorUnaryAdd:
		return true, nil
	case OperatorUnarySub:
		return false, nil
	case OperatorUnaryNot:
		return !value, nil
	case OperatorUnaryBitwiseComplement:
		return !value, nil
	default:
		return false, ErrInvalidOperator
	}
}
