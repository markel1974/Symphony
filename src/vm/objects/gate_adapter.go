package objects

import (
	"fmt"
)

// GateAdapter is a type that wraps a GateKeeper and provides functional adapters to map Go functions to Invocable.
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
func (ga *GateAdapter) CreateSlice(frameId int, highIdx int, lowIdx int, target IObject) (IObject, error) {
	numElem := target.Length()
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
		slice := left.CopyRange(uint(lowIdx), uint(highIdx))
		return ga.factory.NewArray(frameId, slice), nil
	case *String:
		slice := left.Value()[lowIdx:highIdx]
		return ga.factory.NewString(frameId, slice), nil
	case *Bytes:
		slice := left.Value()[lowIdx:highIdx]
		return ga.factory.NewBytes(frameId, slice), nil
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

// CreateObjectPointer creates an object pointer for a given source object and frame, returning an error if the operation fails.
func (ga *GateAdapter) CreateObjectPointer(frame int, objSrc IObject) (*ObjectPointer, error) {
	switch objType := objSrc.(type) {
	case *ObjectPointer:
		return objType, nil
	default:
		obj := ga.factory.NewObjectPointer(frame, &objSrc)
		objPtr, ok := obj.(*ObjectPointer)
		if !ok {
			return nil, fmt.Errorf("not a pointer: %s", obj.TypeName())
		}
		return objPtr, nil
	}
}

// Concrete resolves the provided IObject to a concrete implementation and returns whether the resolution was successful.
func (ga *GateAdapter) Concrete(src IObject) IObject {
	switch io := src.(type) {
	case *Struct:
		return io
	case *Interface:
		return ga.Concrete(io.Value())
	case *Any:
		return io
	//case *ObjectPointer:
	//	return ga.Concrete(*io.Value())
	default:
		return src
	}
}
