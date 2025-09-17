package objects

import (
	"reflect"
	"sync"
	"time"
)

// GateAllocator manages the allocation and lifecycle of IObject instances, leveraging object pools for efficiency.
type GateAllocator struct {
	gk *GateKeeper

	trueValue         IObject
	falseValue        IObject
	undefinedValue    IObject
	undefinedIterator IIterator

	// Pools for primitive and common types
	poolBool          sync.Pool
	poolChar          sync.Pool
	poolInt           sync.Pool
	poolFloat         sync.Pool
	poolString        sync.Pool
	poolTime          sync.Pool
	poolObjectPointer sync.Pool
	poolInterface     sync.Pool
	poolError         sync.Pool

	// Pools for containers
	poolBytes  sync.Pool
	poolArray  sync.Pool
	poolMap    sync.Pool
	poolStruct sync.Pool

	// Pools for iterators
	poolArrayIterator  sync.Pool
	poolBytesIterator  sync.Pool
	poolStringIterator sync.Pool
	poolMapIterator    sync.Pool
	poolStructIterator sync.Pool
	poolAny            sync.Pool

	allocatedObjects *AllocatedObjects
}

// NewGateAllocator initializes a new GateAllocator instance using the provided GateKeeper and maximum allocations limit.
func NewGateAllocator(gk *GateKeeper) *GateAllocator {
	ga := &GateAllocator{
		gk: gk,
	}
	ga.allocatedObjects = NewAllocatedObjects(gk)
	// Initialization of static Code
	ga.trueValue = newBool(NewAllocator(gk, FrameStatic), true)
	ga.falseValue = newBool(NewAllocator(gk, FrameStatic), false)
	ga.undefinedValue = newUndefined(NewAllocator(gk, FrameStatic))
	ga.undefinedIterator = newUndefinedIterator(NewAllocator(gk, FrameStatic))

	// Primitive types
	ga.poolBool.New = ga.allocatedObjects.NewBool
	ga.poolChar.New = ga.allocatedObjects.NewChar
	ga.poolInt.New = ga.allocatedObjects.NewInt
	ga.poolFloat.New = ga.allocatedObjects.NewFloat
	ga.poolString.New = ga.allocatedObjects.NewString
	ga.poolTime.New = ga.allocatedObjects.NewTime
	ga.poolObjectPointer.New = ga.allocatedObjects.NewObjectPointer
	ga.poolInterface.New = ga.allocatedObjects.NewInterface
	ga.poolError.New = ga.allocatedObjects.NewError

	// Containers
	ga.poolBytes.New = ga.allocatedObjects.NewBytes
	ga.poolArray.New = ga.allocatedObjects.NewArray
	ga.poolMap.New = ga.allocatedObjects.NewMap
	ga.poolStruct.New = ga.allocatedObjects.NewStruct
	ga.poolAny.New = ga.allocatedObjects.NewAny

	// Iterators
	ga.poolArrayIterator.New = ga.allocatedObjects.NewArrayIterator
	ga.poolBytesIterator.New = ga.allocatedObjects.NewBytesIterator
	ga.poolStringIterator.New = ga.allocatedObjects.NewStringIterator
	ga.poolMapIterator.New = ga.allocatedObjects.NewMapIterator
	ga.poolStructIterator.New = ga.allocatedObjects.NewStructIterator

	return ga
}

func (f *GateAllocator) AssignAllocator(object IObject) {
	object.setAllocator(NewAllocator(f.gk, FrameStatic))
}

// Reset resets the internal counter of the GateAllocator to its initial state, 0.
func (f *GateAllocator) Reset() {
	f.ReleaseAll()
}

// AllocatedObjects returns the current count of allocated objects managed by the GateAllocator.
func (f *GateAllocator) AllocatedObjects() uint64 {
	return uint64(f.allocatedObjects.Counter())
}

// FalseValue retrieves the internally maintained "false" IObject for this GateAllocator instance.
func (f *GateAllocator) FalseValue() IObject {
	return f.falseValue
}

// TrueValue returns the true Code represented by the GateAllocator instance as an IObject.
func (f *GateAllocator) TrueValue() IObject {
	return f.trueValue
}

// Boolean returns one of two predefined IObject Code based on the provided boolean argument.
func (f *GateAllocator) Boolean(v bool) IObject {
	if v {
		return f.trueValue
	}
	return f.falseValue
}

// UndefinedValue returns the undefined object managed by the GateAllocator as an IObject.
func (f *GateAllocator) UndefinedValue() IObject {
	return f.undefinedValue
}

// --- Methods for Non-Poolable objects ---

// NewFunc creates a new compiled function object with the given frame, name, instructions, locals, parameters, and settings.
// It initializes the object, acquires resources if available, and returns an IObject instance.
// If resource acquisition fails, it returns the undefined object from the allocator.
func (f *GateAllocator) NewFunc(frame int, name string, instructions []byte, numLocals int, numParameters int, varArgs bool, source map[int]int, free []*ObjectPointer) IObject {
	return newFunc(NewAllocator(f.gk, frame), name, instructions, numLocals, numParameters, varArgs, source, free)
}

// NewFuncInternal creates a new IObject using the provided frame and CallId, ensuring proper allocation and preparation.
func (f *GateAllocator) NewFuncInternal(frame int, id CallId) IObject {
	return newFuncInternal(NewAllocator(f.gk, frame), id)
}

// NewFuncImport creates a new function import object with the given frame, name, and callable function.
// Returns an IObject instance or a default undefined Code if object acquisition fails.
func (f *GateAllocator) NewFuncImport(frame int, name string, args int, fn FuncCallable) IObject {
	return newFuncImport(NewAllocator(f.gk, frame), name, args, fn)
}

// NewFuncJit creates and returns a new instance of a function-based IObject with the specified frame, name, and data.
func (f *GateAllocator) NewFuncJit(frame int, name string, data []byte) IObject {
	return newFuncJit(NewAllocator(f.gk, frame), name, data)
}

// NewBool creates a new boolean object wrapped in the IObject interface, based on the provided boolean Code.
func (f *GateAllocator) NewBool(_ int, v bool) IObject {
	return f.Boolean(v)
}

// NewInt retrieves an Int object from the pool, sets its frame and Code, and returns it as an IObject.
func (f *GateAllocator) NewInt(frame int, v int64) IObject {
	obj := f.poolInt.Get().(*Int)
	obj.setFrame(frame)
	obj.data = v
	return obj
}

// NewChar retrieves a Char object from the pool, sets its frame and Code, and returns it as an IObject instance.
func (f *GateAllocator) NewChar(frame int, v rune) IObject {
	obj := f.poolChar.Get().(*Char)
	obj.setFrame(frame)
	obj.data = v
	return obj
}

// NewFloat creates and initializes a new Float object with the specified frame and float64 Code. It returns the object.
func (f *GateAllocator) NewFloat(frame int, v float64) IObject {
	obj := f.poolFloat.Get().(*Float)
	obj.setFrame(frame)
	obj.data = v
	return obj
}

// NewString creates a new String object with the specified frame and string Code, truncating the Code if it exceeds MaxStringLen.
func (f *GateAllocator) NewString(frame int, v string) IObject {
	if len(v) > MaxStringLen {
		v = v[0:MaxStringLen]
	}
	obj := f.poolString.Get().(*String)
	obj.setFrame(frame)
	obj.data = v
	return obj
}

// NewTime creates a new Time object from the object pool, initializes it with the given frame and time Code, and returns it.
func (f *GateAllocator) NewTime(frame int, value time.Time) IObject {
	obj := f.poolTime.Get().(*Time)
	obj.setFrame(frame)
	obj.data = value
	return obj
}

// NewBytes creates and returns a new Bytes object with the specified frame and byte slice, truncating if it exceeds maxBytesLen.
func (f *GateAllocator) NewBytes(frame int, v []byte) IObject {
	if len(v) > maxBytesLen {
		v = v[0:maxBytesLen]
	}
	obj := f.poolBytes.Get().(*Bytes)
	obj.setFrame(frame)
	obj.data = v
	return obj
}

// NewArray creates a new array IObject with the specified frame and elements, truncating elements if exceeding maxArrayLen.
func (f *GateAllocator) NewArray(frame int, v []IObject) IObject {
	if len(v) > maxArrayLen {
		v = v[0:maxArrayLen]
	}
	obj := f.poolArray.Get().(*Array)
	obj.setFrame(frame)
	obj.data = v
	return obj
}

// NewMap creates a new Map object from a provided frame and a map of string keys to IObject Code.
// It reuses an object from the pool if available, setting its frame and Code.
// If the provided map exceeds the maximum allowed len, it may truncate the map.
// Returns the newly created IObject.
func (f *GateAllocator) NewMap(frame int, v map[string]IObject) IObject {
	if len(v) > maxMapLen {
		// todo truncate map
	}
	obj := f.poolMap.Get().(*Map)
	obj.setFrame(frame)
	obj.data = v
	return obj
}

// NewStruct creates and initializes a new Struct object with a specified frame, type name, and key-Code pair map.
func (f *GateAllocator) NewStruct(frame int, name string, v map[string]IObject) IObject {
	if len(v) > maxStructLen {
		// todo truncate struct
	}
	obj := f.poolStruct.Get().(*Struct)
	obj.setFrame(frame)
	obj.name = name
	obj.data = v
	return obj
}

// NewAny creates and initializes a new Any object from the pool using the given frame and value parameters.
func (f *GateAllocator) NewAny(frame int, value interface{}) IObject {
	obj := f.poolAny.Get().(*Any)
	obj.setFrame(frame)
	obj.data = value
	obj.valueOf = reflect.ValueOf(value)
	obj.kind = obj.valueOf.Type()
	return obj
}

// NewError creates and returns a new Error object with the specified frame and error message.
func (f *GateAllocator) NewError(frame int, e string) IObject {
	obj := f.poolError.Get().(*Error)
	obj.setFrame(frame)
	obj.data = f.NewString(FrameStatic, e)
	return obj
}

// NewObjectPointer creates and initializes an ObjectPointer with the given frame and IObject reference, returning it.
func (f *GateAllocator) NewObjectPointer(frame int, v *IObject) IObject {
	obj := f.poolObjectPointer.Get().(*ObjectPointer)
	obj.setFrame(frame)
	obj.acquire(v)
	return obj
}

// NewInterface initializes and returns a new Interface object with the given frame, Code, and internal table.
func (f *GateAllocator) NewInterface(frame int, value IObject, iTable map[string]IObject) IObject {
	obj := f.poolInterface.Get().(*Interface)
	obj.setFrame(frame)
	obj.data = value
	obj.iTable = iTable
	return obj
}

// NewArrayIterator creates and initializes a reusable ArrayIterator instance from the allocator pool for traversing an array.
func (f *GateAllocator) NewArrayIterator(frame int, v []IObject, index int) IIterator {
	obj := f.poolArrayIterator.Get().(*ArrayIterator)
	obj.setFrame(frame)
	obj.data = v
	obj.index = index
	obj.length = len(v)
	return obj
}

// NewBytesIterator allocates and initializes a BytesIterator instance for iterating over a provided byte slice.
func (f *GateAllocator) NewBytesIterator(frame int, v []byte, index int) IIterator {
	obj := f.poolBytesIterator.Get().(*BytesIterator)
	obj.setFrame(frame)
	obj.data = v
	obj.index = index
	obj.length = len(v)
	return obj
}

// NewStringIterator creates and initializes a new StringIterator instance from the pool with the given frame, rune slice, and index.
func (f *GateAllocator) NewStringIterator(frame int, v []rune, index int) IIterator {
	obj := f.poolStringIterator.Get().(*StringIterator)
	obj.setFrame(frame)
	obj.data = v
	obj.index = index
	obj.length = len(v)
	return obj
}

// NewMapIterator initializes and returns a new MapIterator instance with the given frame, map Code, and starting index.
func (f *GateAllocator) NewMapIterator(frame int, v map[string]IObject, index int) IIterator {
	obj := f.poolMapIterator.Get().(*MapIterator)
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	obj.setFrame(frame)
	obj.data = v
	obj.keys = keys
	obj.index = index
	obj.length = len(keys)
	return obj
}

// NewStructIterator initializes a StructIterator with the provided frame, Code map, and starting index for iteration.
func (f *GateAllocator) NewStructIterator(frame int, v map[string]IObject, index int) IIterator {
	obj := f.poolStructIterator.Get().(*StructIterator)
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	obj.setFrame(frame)
	obj.data = v
	obj.keys = keys
	obj.index = index
	obj.length = len(keys)
	return obj
}

// SetPointer updates the ObjectPointer to reference a new IObject and manages reference counting and object lifecycle.
func (f *GateAllocator) SetPointer(ptr *ObjectPointer, newValue IObject) {
	oldValue := *ptr.Value()
	if oldValue.Equals(newValue) {
		return
	}
	ptr.acquire(&newValue)
	if oldValue.Frame() != FrameStatic {
		if oldValue.ReleaseRef() <= 0 {
			f.releaseObject(oldValue.Frame(), oldValue)
		}
	}
}

// ReleaseAll releases all allocated objects managed by the GateAllocator and resets the allocatedObjects map.
func (f *GateAllocator) ReleaseAll() {
	if f.allocatedObjects.Counter() == 0 {
		return
	}
	objectsToRelease := f.allocatedObjects.ReleaseCandidates()
	for _, obj := range objectsToRelease {
		for obj.RefCount() > 0 {
			obj.ReleaseRef()
		}
		f.releaseObject(obj.Frame(), obj)
	}
	f.allocatedObjects.Reset()
}

// ReleaseObjects cleans up and removes the provided objects associated with a specific frame from memory allocation.
func (f *GateAllocator) ReleaseObjects(frame int, objects []IObject) {
	garbageCandidates := make(map[IObject]bool)
	for _, obj := range objects {
		if obj != nil && obj.Frame() != FrameStatic && obj.Frame() == frame {
			garbageCandidates[obj] = true
		}
	}
	if len(garbageCandidates) == 0 {
		return
	}
	// Simulate the Removal of Internal References within the Group ---
	for obj := range garbageCandidates {
		switch o := obj.(type) {
		case *ObjectPointer:
			if target := *o.Value(); garbageCandidates[target] {
				target.ReleaseRef()
			}
		case *Array:
			for _, elem := range o.Values() {
				if garbageCandidates[elem] {
					elem.ReleaseRef()
				}
			}
		case *Map:
			for _, val := range o.Values() {
				if garbageCandidates[val] {
					val.ReleaseRef()
				}
			}
		case *Struct:
			for _, val := range o.Values() {
				if garbageCandidates[val] {
					val.ReleaseRef()
				}
			}
		case *Interface:
			for _, val := range o.iTable {
				if garbageCandidates[val] {
					val.ReleaseRef()
				}
			}
		}
	}
	for obj := range garbageCandidates {
		f.releaseObject(frame, obj)
	}
}

// releaseObject reclaims memory for an object and returns it to the appropriate pool when no references remain.
func (f *GateAllocator) releaseObject(frame int, obj IObject) {
	if obj == nil || obj.Frame() == FrameStatic || obj.RefCount() > 0 {
		return
	}
	f.allocatedObjects.Remove(obj)

	//Set static before release
	obj.SetStatic()
	switch o := obj.(type) {
	case *ObjectPointer:
		target := *o.data
		if target.Frame() != FrameStatic {
			if target.ReleaseRef() <= 0 {
				f.releaseObject(target.Frame(), target)
			}
		}
		f.poolObjectPointer.Put(o)
	case *Interface:
		f.poolInterface.Put(o)
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
	case *Time:
		f.poolTime.Put(o)
	case *Bytes:
		o.data = nil
		f.poolBytes.Put(o)
	case *Error:
		o.data = nil
		f.poolError.Put(o)
	case *Array:
		f.ReleaseObjects(frame, o.data)
		o.data = o.data[:0]
		f.poolArray.Put(o)
	case *Map:
		for _, v := range o.data {
			f.releaseObject(frame, v)
		}
		for k := range o.data {
			delete(o.data, k)
		}
		f.poolMap.Put(o)
	case *Struct:
		for _, v := range o.data {
			f.releaseObject(frame, v)
		}
		for k := range o.data {
			delete(o.data, k)
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

// GobDecode decodes the Bool object from a byte slice encoded using Gob into its internal Code.
func (f *GateAllocator) GobDecode() error {
	return nil
}

// GobEncode encodes the Bool instance into a byte slice representation. Returns the byte slice and any encoding error.
func (f *GateAllocator) GobEncode() ([]byte, error) {
	return nil, nil
}
