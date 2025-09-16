package objects

func logicalOpNil(gk IGateAllocator, op LogicalOperator) (IObject, error) {
	switch op {
	case OperatorLogicalEq:
		return gk.FalseValue(), nil
	case OperatorLogicalNotEq:
		return gk.TrueValue(), nil
	default:
		return gk.FalseValue(), ErrInvalidOperator
	}
}

// logicalOpInt64 performs a logical comparison between two int64 Code using the specified LogicalOperator.
// Returns the boolean result of the operation or an error if the operator is invalid.
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

// arithmeticOpFloat64 performs the specified arithmetic operation on two float64 operands and returns the result or an error.
// Supported operations are addition, subtraction, multiplication, and division. Returns errors for invalid operators or division by zero.
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

// logicalOpFloat64 performs a logical operation between two float64 Code using the specified LogicalOperator.
// Returns the result as a boolean and an error if the operator is invalid.
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

// arithmeticOpInt64 performs the specified arithmetic or bitwise operation on two int64 operands and returns the result.
// Returns an error if the operation is invalid or division by zero occurs in applicable cases.
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
