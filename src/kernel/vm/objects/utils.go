package objects

// toInt64 converts an IObject to an int64 value with a success flag.
// Returns the converted int64 value if successful and a boolean indicating success or failure.
func toInt64(rhsIn IObject) (int64, error) {
	switch rhs := rhsIn.(type) {
	case *Bool:
		if rhs.value {
			return 1, nil
		}
		return 0, nil
	case *Int:
		return rhs.value, nil
	case *Float:
		return int64(rhs.value), nil
	case *Char:
		return int64(rhs.value), nil
	case *Time:
		return rhs.value.Unix(), nil
	case *String:
		return int64(rhs.Length()), nil
	case *Map:
		return int64(rhs.Length()), nil
	case *Array:
		return int64(rhs.Length()), nil
	default:
		return 0, ErrInvalidOperator
	}
}

// toFloat64 converts an IObject to a float64 if possible, returning the converted value and a boolean indicating success.
func toFloat64(rhsIn IObject) (float64, error) {
	switch rhs := rhsIn.(type) {
	case *Bool:
		if rhs.value {
			return 1, nil
		}
		return 0, nil
	case *Int:
		return float64(rhs.value), nil
	case *Float:
		return rhs.value, nil
	case *Char:
		return float64(rhs.value), nil
	case *Time:
		return float64(rhs.value.Unix()), nil
	case *String:
		return float64(rhs.Length()), nil
	case *Map:
		return float64(rhs.Length()), nil
	case *Array:
		return float64(rhs.Length()), nil
	default:
		return 0, ErrInvalidOperator
	}
}

// logicalOpInt64 performs a logical comparison between two int64 values using the specified LogicalOperator.
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

// logicalOpFloat64 performs a logical operation between two float64 values using the specified LogicalOperator.
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
