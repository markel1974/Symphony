package objects

import "time"

type GateAllocator struct {
	gk *GateKeeper

	trueValue         IObject
	falseValue        IObject
	undefinedValue    IObject
	counter           int64
	maxAllocations    int64
	undefinedIterator IIterator

	// Aggiungiamo i pool per gli oggetti
	//intPool   sync.Pool
	//floatPool sync.Pool
	//charPool  sync.Pool
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
	//ga.intPool.New = func() interface{} {
	//	return &_newInt(gk, 0) // Crea un Int con valore di default
	//}
	return ga
}

func (f *GateAllocator) Reset() {
	f.counter = 0
}

func (f *GateAllocator) acquireObject() error {
	f.counter++
	if f.maxAllocations > 0 {
		if f.counter > f.maxAllocations {
			return ErrObjectAllocLimit
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

// NewArrayImmutable constructs a new ArrayImmutable instance with the provided slice of IObject, ensuring immutability.
func (f *GateAllocator) NewArrayImmutable(frame int, values []IObject) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newArrayImmutable(f.gk, frame, values)
}

// NewBool creates and returns a new Bool object initialized with the specified boolean value.
func (f *GateAllocator) NewBool(frame int, value bool) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newBool(f.gk, frame, value)
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
	return newChar(f.gk, frame, value)
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
	return newFloat(f.gk, frame, v)
}

// NewInt creates and returns a new instance of Int initialized with the given int64 value.
func (f *GateAllocator) NewInt(frame int, v int64) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	//obj := f.intPool.Get().(*Int)
	//obj.value = v
	//return obj
	return newInt(f.gk, frame, v)
}

//func (f *GateKeeper) ReleaseInt(obj *Int) {
//	// It's good practice to reset the object's state before putting it back in the pool
//	obj.value = 0
//	f.intPool.Put(obj)
//}

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

// NewMapImmutable creates a new immutable map with string keys and IObject values from the provided map.
func (f *GateAllocator) NewMapImmutable(frame int, v map[string]IObject) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newMapImmutable(f.gk, frame, v)
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
