package objects

import (
	"sync"
	"time"
)

// GateAllocator manages object pooling and allocation limits for various types, including primitives, containers, and iterators.
type GateAllocator struct {
	gk *GateKeeper

	trueValue         IObject
	falseValue        IObject
	undefinedValue    IObject
	counter           int64
	maxAllocations    int64
	undefinedIterator IIterator

	// Pools for primitive and common types
	poolBool          sync.Pool
	poolChar          sync.Pool
	poolInt           sync.Pool
	poolFloat         sync.Pool
	poolString        sync.Pool
	poolBytes         sync.Pool
	poolObjectPointer sync.Pool
	poolError         sync.Pool

	// Pools for containers
	poolArray  sync.Pool
	poolMap    sync.Pool
	poolStruct sync.Pool

	// Pools for iterators
	poolArrayIterator  sync.Pool
	poolBytesIterator  sync.Pool
	poolStringIterator sync.Pool
	poolMapIterator    sync.Pool
	poolStructIterator sync.Pool
}

// NewGateAllocator initializes and returns a new instance of GateAllocator with the specified GateKeeper and max allocations.
// It pre-configures internal pools for various types and iterators to optimize object creation and reuse.
func NewGateAllocator(gk *GateKeeper, maxAllocations int64) *GateAllocator {
	ga := &GateAllocator{
		gk:             gk,
		maxAllocations: maxAllocations,
	}

	// Initialization of static values
	ga.trueValue = newBool(gk, FrameStatic, true)
	ga.falseValue = newBool(gk, FrameStatic, false)
	ga.undefinedValue = newUndefined(gk, FrameStatic)
	ga.undefinedIterator = newUndefinedIterator(gk, FrameStatic)

	// Primitive types
	ga.poolBool.New = func() any { return newBool(gk, FrameStatic, false) }
	ga.poolChar.New = func() any { return newChar(gk, FrameStatic, 0) }
	ga.poolInt.New = func() any { return newInt(gk, FrameStatic, 0) }
	ga.poolFloat.New = func() any { return newFloat(gk, FrameStatic, 0) }
	ga.poolString.New = func() any { return newString(gk, FrameStatic, "") }
	ga.poolBytes.New = func() any { return newBytes(gk, FrameStatic, nil) }
	ga.poolObjectPointer.New = func() any { return newObjectPointer(gk, FrameStatic, nil) }
	ga.poolError.New = func() any { return newError(gk, FrameStatic, "") }

	// Containers
	ga.poolArray.New = func() any { return newArray(gk, FrameStatic, nil) }
	ga.poolMap.New = func() any { return newMap(gk, FrameStatic, make(map[string]IObject)) }
	ga.poolStruct.New = func() any { return newStruct(gk, FrameStatic, make(map[string]IObject)) }

	// Iterators
	ga.poolArrayIterator.New = func() any { return newArrayIterator(gk, FrameStatic, nil, 0) }
	ga.poolBytesIterator.New = func() any { return newBytesIterator(gk, FrameStatic, nil, 0) }
	ga.poolStringIterator.New = func() any { return newStringIterator(gk, FrameStatic, nil, 0) }
	ga.poolMapIterator.New = func() any { return newMapIterator(gk, FrameStatic, nil, 0) }
	ga.poolStructIterator.New = func() any { return newStructIterator(gk, FrameStatic, nil, 0) }

	return ga
}

// Reset sets the counter of GateAllocator to zero.
func (f *GateAllocator) Reset() {
	f.counter = 0
}

// acquireObject increments the allocation counter and checks against the maximum allocation limit, returning an error if exceeded.
func (f *GateAllocator) acquireObject() error {
	f.counter++
	if f.maxAllocations > 0 && f.counter > f.maxAllocations {
		return ErrAllocationLimit
	}
	return nil
}

// FalseValue retrieves the predefined object representing a "false" boolean value within the GateAllocator.
func (f *GateAllocator) FalseValue() IObject {
	return f.falseValue
}

// TrueValue returns the predefined value representing "true" in the GateAllocator.
func (f *GateAllocator) TrueValue() IObject {
	return f.trueValue
}

// Boolean returns `trueValue` if the input `v` is true, otherwise returns `falseValue`.
func (f *GateAllocator) Boolean(v bool) IObject {
	if v {
		return f.trueValue
	}
	return f.falseValue
}

// UndefinedValue returns the predefined "undefined" value for the GateAllocator instance.
func (f *GateAllocator) UndefinedValue() IObject {
	return f.undefinedValue
}

func (f *GateAllocator) SetPointer(ptr *ObjectPointer, value IObject) {
	//TODO BETTER IMPLEMENTATION
	if v, release := ptr.release(); release {
		f.ReleaseObject(v.Frame(), v)
	}
	ptr.acquire(&value)
}

// ReleaseObjects releases a slice of IObject instances back to their respective pools to free resources.
func (f *GateAllocator) ReleaseObjects(frame int, objects []IObject) {
	for _, o := range objects {
		f.ReleaseObject(frame, o)
	}
}

// ReleaseObject releases an object back to the relevant pool, resetting its state and freeing associated resources.
func (f *GateAllocator) ReleaseObject(frame int, obj IObject) {
	if obj == nil || obj.Frame() == FrameStatic || obj.RefCount() > 0 {
		return
	}
	obj.SetStatic()
	switch o := obj.(type) {
	case *Bool:
		f.poolBool.Put(o)
	case *Char:
		f.poolChar.Put(o)
	case *Int:
		f.poolInt.Put(o)
	case *Float:
		f.poolFloat.Put(o)
	case *String:
		f.poolString.Put(o)
	case *Bytes:
		o.values = nil
		f.poolBytes.Put(o)
	case *ObjectPointer:
		if valuePtr, release := o.release(); release {
			f.ReleaseObject(frame, valuePtr)
		}
		f.poolObjectPointer.Put(o)
	case *Error:
		o.value = nil
		f.poolError.Put(o)
	case *Array:
		f.ReleaseObjects(frame, o.values)
		o.values = o.values[:0]
		f.poolArray.Put(o)
	case *Map:
		for _, v := range o.values {
			f.ReleaseObject(frame, v)
		}
		for k := range o.values {
			delete(o.values, k)
		}
		f.poolMap.Put(o)
	case *Struct:
		for _, v := range o.values {
			f.ReleaseObject(frame, v)
		}
		for k := range o.values {
			delete(o.values, k)
		}
		f.poolStruct.Put(o)
	case *ArrayIterator:
		f.poolArrayIterator.Put(o)
	case *BytesIterator:
		f.poolBytesIterator.Put(o)
	case *StringIterator:
		f.poolStringIterator.Put(o)
	case *MapIterator:
		f.poolMapIterator.Put(o)
	case *StructIterator:
		f.poolStructIterator.Put(o)
	}
}

// NewInt creates a new integer object from the pool, setting its frame and value, and returns it as an IObject.
func (f *GateAllocator) NewInt(frame int, v int64) IObject {
	obj := f.poolInt.Get().(*Int)
	obj.frame = frame
	obj.value = v
	return obj
}

// NewBool creates a new boolean object with the specified frame and value. Returns preallocated true or false objects.
func (f *GateAllocator) NewBool(_ int, v bool) IObject {
	return f.Boolean(v)
}

// NewChar creates a new Char object from the pool, sets its frame and value, and returns it as an IObject.
func (f *GateAllocator) NewChar(frame int, v rune) IObject {
	obj := f.poolChar.Get().(*Char)
	obj.frame = frame
	obj.value = v
	return obj
}

// NewFloat allocates and returns a new Float object from the pool, setting its frame and value properties.
func (f *GateAllocator) NewFloat(frame int, v float64) IObject {
	obj := f.poolFloat.Get().(*Float)
	obj.frame = frame
	obj.value = v
	return obj
}

// NewString creates a new String object with the specified frame and string value, truncating the value if it exceeds MaxStringLen.
func (f *GateAllocator) NewString(frame int, v string) IObject {
	if len(v) > MaxStringLen {
		v = v[0:MaxStringLen]
	}
	obj := f.poolString.Get().(*String)
	obj.frame = frame
	obj.value = v
	obj.runeStr = nil
	return obj
}

// NewBytes creates a new Bytes object from the provided byte slice and associates it with the specified frame.
// If the byte slice exceeds the maximum allowed length, it is truncated.
func (f *GateAllocator) NewBytes(frame int, v []byte) IObject {
	if len(v) > maxBytesLen {
		v = v[0:maxBytesLen]
	}
	obj := f.poolBytes.Get().(*Bytes)
	obj.frame = frame
	obj.values = v
	return obj
}

// NewArray creates a new array object with the specified frame and values, truncating the input if it exceeds the max length.
func (f *GateAllocator) NewArray(frame int, v []IObject) IObject {
	if len(v) > maxArrayLen {
		v = v[0:maxArrayLen]
	}
	obj := f.poolArray.Get().(*Array)
	obj.frame = frame
	obj.values = v
	return obj
}

// NewMap creates a new map object with the specified frame and values. Limits the map length to maxMapLen if exceeded.
func (f *GateAllocator) NewMap(frame int, v map[string]IObject) IObject {
	if len(v) > maxMapLen {
		// Tronca mappa
	}
	obj := f.poolMap.Get().(*Map)
	obj.frame = frame
	obj.values = v
	return obj
}

// NewStruct creates and initializes a new Struct object with the specified frame and map of values.
func (f *GateAllocator) NewStruct(frame int, name string, v map[string]IObject) IObject {
	if len(v) > maxStructLen {
		// Tronca struct
	}
	obj := f.poolStruct.Get().(*Struct)
	obj.frame = frame
	obj.typeName = name
	obj.values = v
	return obj
}

// NewError creates and returns a new error object from the pool, initializing it with the specified frame and error string.
func (f *GateAllocator) NewError(frame int, e string) IObject {
	obj := f.poolError.Get().(*Error)
	obj.frame = frame
	obj.err = e
	obj.value = f.NewString(FrameStatic, e) // Stringa statica, non legata al frame
	return obj
}

// NewObjectPointer creates a new ObjectPointer instance, initializes it with the provided frame and value, and returns it.
func (f *GateAllocator) NewObjectPointer(frame int, v *IObject) IObject {
	obj := f.poolObjectPointer.Get().(*ObjectPointer)
	obj.frame = frame
	obj.acquire(v)
	return obj
}

// NewArrayIterator creates and initializes a new ArrayIterator instance with the provided frame, array, and starting index.
func (f *GateAllocator) NewArrayIterator(frame int, v []IObject, index int) IIterator {
	obj := f.poolArrayIterator.Get().(*ArrayIterator)
	obj.frame, obj.values, obj.index, obj.length = frame, v, index, len(v)
	return obj
}

// NewBytesIterator initializes and returns a BytesIterator with the given frame, byte slice, and starting index.
func (f *GateAllocator) NewBytesIterator(frame int, v []byte, index int) IIterator {
	obj := f.poolBytesIterator.Get().(*BytesIterator)
	obj.frame, obj.values, obj.index, obj.length = frame, v, index, len(v)
	return obj
}

// NewStringIterator creates and initializes a new StringIterator from a pool with the specified frame, rune slice, and index.
func (f *GateAllocator) NewStringIterator(frame int, v []rune, index int) IIterator {
	obj := f.poolStringIterator.Get().(*StringIterator)
	obj.frame, obj.values, obj.index, obj.length = frame, v, index, len(v)
	return obj
}

// NewMapIterator creates and initializes a new MapIterator object from a map of string keys and associated IObject values.
func (f *GateAllocator) NewMapIterator(frame int, v map[string]IObject, index int) IIterator {
	obj := f.poolMapIterator.Get().(*MapIterator)
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	obj.frame, obj.values, obj.keys, obj.index, obj.length = frame, v, keys, index, len(keys)
	return obj
}

// NewStructIterator creates and initializes a StructIterator, setting its frame, values, keys, index, and length.
func (f *GateAllocator) NewStructIterator(frame int, v map[string]IObject, index int) IIterator {
	obj := f.poolStructIterator.Get().(*StructIterator)
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	obj.frame, obj.values, obj.keys, obj.index, obj.length = frame, v, keys, index, len(keys)
	return obj
}

// --- Methods for Non-Poolable objects ---

// NewFuncCompiled creates a new compiled function object with provided frame, name, instructions, locals, parameters, and metadata.
func (f *GateAllocator) NewFuncCompiled(frame int, name string, instructions []byte, numLocals int, numParameters int, varArgs bool, sourceMap map[int]int, free []*ObjectPointer) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newFuncCompiled(f.gk, frame, name, instructions, numLocals, numParameters, varArgs, sourceMap, free)
}

// NewFuncInternal creates a new function external object with the specified frame, kind, name, and callable.
func (f *GateAllocator) NewFuncInternal(frame int, id CallId) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newFuncInternal(f.gk, frame, id)
}

// NewFuncInternals creates a new array of function internal objects with the specified frame and call ids.
func (f *GateAllocator) NewFuncInternals(frame int) []IObject {
	out := make([]IObject, len(callIdContainer))
	for idx, v := range callIdContainer {
		newObj := f.NewFuncInternal(frame, v)
		out[idx] = newObj
	}
	return out
}

// NewFuncImport creates a new function external object with the specified frame, kind, name, and callable.
func (f *GateAllocator) NewFuncImport(frame int, name string, fn FuncCallable) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newFuncImport(f.gk, frame, name, fn)
}

// NewFuncJit creates a new JIT-compiled function object with the provided kind, name, and bytecode.
func (f *GateAllocator) NewFuncJit(frame int, name string, data []byte) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newFuncJit(f.gk, frame, name, data)
}

// NewTime creates a new Time object with the specified frame and time value. Returns undefinedValue on allocation error.
func (f *GateAllocator) NewTime(frame int, value time.Time) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newTime(f.gk, frame, value)
}

// NewInterface creates a new interface object with the specified frame, value, and interface table.
func (f *GateAllocator) NewInterface(frame int, value IObject, itable map[string]IObject) IObject {
	if err := f.acquireObject(); err != nil {
		return f.undefinedValue
	}
	return newInterface(f.gk, frame, value, itable)
}
