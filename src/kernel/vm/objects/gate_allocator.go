package objects

import (
	"log"
	"sync"
	"time"
)

type GateAllocator struct {
	gk *GateKeeper

	trueValue         IObject
	falseValue        IObject
	undefinedValue    IObject
	counter           int64
	maxAllocations    int64
	undefinedIterator IIterator

	poolBool  sync.Pool
	poolChar  sync.Pool
	poolInt   sync.Pool
	poolFloat sync.Pool
}

func NewGateAllocator(gk *GateKeeper, maxAllocations int64) *GateAllocator {
	ga := &GateAllocator{
		gk:             gk,
		maxAllocations: maxAllocations,
	}
	ga.trueValue = newBool(gk, FrameStatic, true)
	ga.falseValue = newBool(gk, FrameStatic, false)
	ga.undefinedValue = newUndefined(gk, FrameStatic)
	ga.undefinedIterator = newUndefinedIterator(gk, FrameStatic)
	ga.maxAllocations = maxAllocations

	ga.poolBool.New = func() any {
		return newBool(gk, FrameStatic, false)
	}
	ga.poolChar.New = func() any {
		return newChar(gk, FrameStatic, 0)
	}
	ga.poolInt.New = func() any {
		return newInt(gk, FrameStatic, 0)
	}
	ga.poolFloat.New = func() any {
		return newFloat(gk, FrameStatic, 0)
	}
	return ga
}

func (f *GateAllocator) Reset() {
	f.counter = 0
}

func (f *GateAllocator) acquireObject() error {
	f.counter++
	if f.maxAllocations > 0 {
		if f.counter > f.maxAllocations {
			return ErrAllocationLimit
		}
	}
	return nil
}

// FalseValue returns the false representation as an IObject from the GateKeeper instance.
func (f *GateAllocator) FalseValue() IObject {
	return f.falseValue
}

// TrueValue returns the IObject instance representing the true value from the GateKeeper.
func (f *GateAllocator) TrueValue() IObject {
	return f.trueValue
}

// UndefinedValue returns the undefined value of the GateKeeper as an IObject.
func (f *GateAllocator) UndefinedValue() IObject {
	return f.undefinedValue
}

// ReleaseObjects releases the given objects back to the pool.
func (f *GateAllocator) ReleaseObjects(obj []IObject) {
	for _, o := range obj {
		f.ReleaseObject(o)
	}
}

// ReleaseObject releases the given object back to the pool.
func (f *GateAllocator) ReleaseObject(obj IObject) {
	if obj.Frame() == FrameStatic {
		return
	}
	switch o := obj.(type) {
	case *Bool:
		o.value = false
		o.frame = FrameStatic
		f.poolBool.Put(obj)
	case *Char:
		o.value = 0
		o.frame = FrameStatic
		f.poolChar.Put(obj)
	case *Int:
		o.value = 0
		o.frame = FrameStatic
		f.poolInt.Put(obj)
	case *Float:
		o.value = 0
		o.frame = FrameStatic
		f.poolFloat.Put(obj)
	case *String:
		//v.factory.ReleaseString(o)
	case *Array:
		//v.factory.ReleaseArray(o)
	}
}

// NewInt creates and returns a new instance of Int initialized with the given int64 value.
func (f *GateAllocator) NewInt(frame int, v int64) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	//return newInt(f.gk, frame, v)
	i := f.poolInt.Get()
	obj, ok := i.(*Int)
	if !ok {
		log.Printf("Error: poolInt.Get() returned wrong type: %T\n", i)
		return f.undefinedValue
	}
	obj.value = v
	obj.frame = frame
	return obj
}

// NewBytesIterator creates a new BytesIterator for iterating over the provided byte slice `v` using the specified GateKeeper.
func (f *GateAllocator) NewBytesIterator(frame int, v []byte, index int) IIterator {
	if err := f.acquireObject(); err != nil {
		return f.undefinedIterator
	}
	return newBytesIterator(f.gk, frame, v, index)
}

// NewArrayIterator creates a new ArrayIterator for iterating over the provided slice of IObject values.
func (f *GateAllocator) NewArrayIterator(frame int, values []IObject, index int) IIterator {
	if err := f.acquireObject(); err != nil {
		return f.undefinedIterator
	}
	return newArrayIterator(f.gk, frame, values, index)
}

// NewStringIterator creates a new StringIterator instance for a given slice of runes, enabling character traversal.
func (f *GateAllocator) NewStringIterator(frame int, v []rune, index int) IIterator {
	if err := f.acquireObject(); err != nil {
		return f.undefinedIterator
	}
	return newStringIterator(f.gk, frame, v, index)
}

// NewStructIterator creates a new StructIterator instance for iterating over a map with string keys and IObject values.
func (f *GateAllocator) NewStructIterator(frame int, v map[string]IObject, index int) IIterator {
	if err := f.acquireObject(); err != nil {
		return f.undefinedIterator
	}
	return newStructIterator(f.gk, frame, v, index)
}

// NewMapIterator creates and returns a new MapIterator for the provided map of string keys and IObject values.
func (f *GateAllocator) NewMapIterator(frame int, v map[string]IObject, index int) IIterator {
	if err := f.acquireObject(); err != nil {
		return f.undefinedIterator
	}
	return newMapIterator(f.gk, frame, v, index)
}

// NewFuncCompiled creates and returns a new FuncCompiled instance using the provided function metadata and bytecode.
func (f *GateAllocator) NewFuncCompiled(frame int, name string, instructions []byte, numLocals int, numParameters int, varArgs bool, sourceMap map[int]int, free []*ObjectPointer) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newFuncCompiled(f.gk, frame, name, instructions, numLocals, numParameters, varArgs, sourceMap, free)
}

// NewFuncPackage creates a new instance of FuncPackage with the specified kind, name, and callable function.
func (f *GateAllocator) NewFuncPackage(kind string, name string, fn FuncCallable) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newFuncPackage(f.gk, FrameStatic, kind, name, fn)
}

// NewFuncPackageFrame creates a new instance of FuncPackage with the specified kind, name, and callable function.
func (f *GateAllocator) NewFuncPackageFrame(frame int, kind string, name string, fn FuncCallable) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newFuncPackage(f.gk, frame, kind, name, fn)
}

// NewFuncJit creates a new JIT-compiled function object with the specified kind, name, and bytecode data.
func (f *GateAllocator) NewFuncJit(kind string, name string, data []byte) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newFuncJit(f.gk, FrameStatic, kind, name, data)
}

// NewFuncJitFrame creates a new JIT-compiled function frame with provided frame ID, kind, name, and associated data.
func (f *GateAllocator) NewFuncJitFrame(frame int, kind string, name string, data []byte) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newFuncJit(f.gk, frame, kind, name, data)
}

// NewBuiltin creates a new Builtin object with the specified name and index using the GateKeeper.
func (f *GateAllocator) NewBuiltin(frame int, name string, index int) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newBuiltin(f.gk, frame, name, index)
}

// NewArray creates and returns a new Array populated with the provided slice of IObject elements.
func (f *GateAllocator) NewArray(frame int, values []IObject) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newArray(f.gk, frame, values)
}

// NewBool creates and returns a new Bool object initialized with the specified boolean value.
func (f *GateAllocator) NewBool(frame int, value bool) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	//return newBool(f.gk, frame, value)
	i := f.poolBool.Get()
	obj, ok := i.(*Bool)
	if !ok {
		log.Printf("Error: poolBool.Get() returned wrong type: %T\n", i)
		return f.undefinedValue
	}
	obj.value = value
	obj.frame = frame
	return obj
}

// NewBytes creates and returns a new instance of Bytes initialized with the provided byte slice and gk context.
func (f *GateAllocator) NewBytes(frame int, value []byte) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newBytes(f.gk, frame, value)
}

// NewChar creates a new Char instance associated with the GateKeeper, initialized with the given rune value.
func (f *GateAllocator) NewChar(frame int, value rune) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	//return newChar(f.gk, frame, value)
	i := f.poolChar.Get()
	obj, ok := i.(*Char)
	if !ok {
		log.Printf("Error: poolChar.Get() returned wrong type: %T\n", i)
		return f.undefinedValue
	}
	obj.value = value
	obj.frame = frame
	return obj
}

// NewError creates and returns a new Error instance based on the provided IObject value and the associated GateKeeper.
func (f *GateAllocator) NewError(frame int, e string) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newError(f.gk, frame, e)
}

// NewFloat creates a new Float instance with the given float64 value, using the GateKeeper for initialization.
func (f *GateAllocator) NewFloat(frame int, v float64) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	i := f.poolFloat.Get()
	obj, ok := i.(*Float)
	if !ok {
		log.Printf("Error: poolFloat.Get() returned wrong type: %T\n", i)
		return f.undefinedValue
	}
	obj.value = v
	obj.frame = frame
	return newFloat(f.gk, frame, v)
}

// NewObjectPointer creates a new ObjectPointer instance wrapping the provided IObject pointer.
func (f *GateAllocator) NewObjectPointer(frame int, value *IObject) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newObjectPointer(f.gk, frame, value)
}

// NewMap creates and returns a new instance of Map initialized with the provided map of string keys and IObject values.
func (f *GateAllocator) NewMap(frame int, v map[string]IObject) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newMap(f.gk, frame, v)
}

// NewString creates a new instance of String with the given value, utilizing the GateKeeper for initialization.
func (f *GateAllocator) NewString(frame int, value string) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newString(f.gk, frame, value)
}

// NewStruct creates and returns a new instance of Struct using the provided map of string keys and IObject values.
func (f *GateAllocator) NewStruct(frame int, value map[string]IObject) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newStruct(f.gk, frame, value)
}

// NewTime creates a new instance of Time using the provided time.Time value and initializes it with the gk instance.
func (f *GateAllocator) NewTime(frame int, value time.Time) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newTime(f.gk, frame, value)
}

// NewInterface creates a new instance of Interface using the provided IObject value and itable map.
func (f *GateAllocator) NewInterface(frame int, value IObject, itable map[string]IObject) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newInterface(f.gk, frame, value, itable)
}
