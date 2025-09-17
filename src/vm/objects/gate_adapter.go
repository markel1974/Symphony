package objects

import "fmt"

// GateAdapter is a type that wraps a GateKeeper and provides functional adapters to map Go functions to FuncCallable.
type GateAdapter struct {
	factory *GateKeeper
}

// NewGateAdapter creates and returns a new instance of GateAdapter, initialized with the provided GateKeeper.
func NewGateAdapter(gk *GateKeeper) *GateAdapter {
	return &GateAdapter{factory: gk}
}

// ArithmeticOpInt64 performs integer arithmetic or bitwise operations on two int64 Code based on the provided operator.
// Returns the result of the operation and an error if an invalid operator is used or division by zero occurs.
func (ga *GateAdapter) ArithmeticOpInt64(op ArithmeticOperator, lhs int64, rhs int64) (int64, error) {
	return arithmeticOpInt64(lhs, op, rhs)
}

// LogicalOpInt64 performs the specified logical operation on two int64 Code and returns a boolean result or an error.
// Supported operations include less than, greater than, less than or equal, and greater than or equal.
// Returns ErrInvalidOperator if an unsupported operator is provided.
func (ga *GateAdapter) LogicalOpInt64(op LogicalOperator, lhs int64, rhs int64) (bool, error) {
	return logicalOpInt64(lhs, op, rhs)
}

// CreateSlice generates a slice of a target object using the given low and high indices and returns the resulting object.
// The method supports slicing Arrays, Strings, and Bytes, returning an error if the target type is unsupported.
func (ga *GateAdapter) CreateSlice(frameId int, lowObj IObject, highObj IObject, target IObject) (IObject, error) {
	numElem := target.Length()
	lowIdx := int(lowObj.AsInt64())
	highIdx := int(highObj.AsInt64())
	if lowIdx > highIdx {
		lowIdx = highIdx
	}
	if lowIdx < 0 {
		lowIdx = 0
	} else if lowIdx > numElem {
		lowIdx = numElem
	}
	if highIdx < 0 {
		highIdx = 0
	} else if highIdx > numElem {
		highIdx = numElem
	}
	switch left := target.(type) {
	case *Array:
		return ga.factory.NewArray(frameId, left.Values()[lowIdx:highIdx]), nil
	case *String:
		return ga.factory.NewString(frameId, left.Value()[lowIdx:highIdx]), nil
	case *Bytes:
		return ga.factory.NewBytes(frameId, left.Value()[lowIdx:highIdx]), nil
	default:
		return nil, fmt.Errorf("unsupported slice: %s", left.TypeName())
	}
}

// IndexAssign assigns a Code to a nested structure, using selectors to determine the target location.
// It navigates through the provided selectors and performs an assignment on the target object at the final index.
// Returns an error if any selector is invalid, the object is not indexable, or the assignment fails.
func (ga *GateAdapter) IndexAssign(frame int, dst IObject, src IObject, selectors []IObject) error {
	sLen := len(selectors)
	if sLen == 0 {
		return ErrSelectorNotProvided
	}
	for sIdx := sLen - 1; sIdx > 0; sIdx-- {
		next, err := dst.IndexGet(frame, selectors[sIdx])
		if err != nil {
			return ComputeIndexGetError(err, dst.TypeName(), selectors[sIdx].TypeName())
		}
		dst = next
	}
	if err := dst.IndexSet(selectors[0], src); err != nil {
		return ComputeIndexSetError(err, dst.TypeName(), src.TypeName())
	}
	return nil
}
