package objects

import "time"

// AllocatedObjects manages a collection of IObject allocations and tracks the count of allocated objects.
type AllocatedObjects struct {
	gk      IGateKeeper
	objects map[IObject]bool
	counter int
}

// NewAllocatedObjects creates and returns a new instance of AllocatedObjects with initialized internal structures.
func NewAllocatedObjects(gk IGateKeeper) *AllocatedObjects {
	return &AllocatedObjects{
		gk:      gk,
		objects: make(map[IObject]bool),
		counter: 0,
	}
}

func (f *AllocatedObjects) Reset() {
	f.objects = make(map[IObject]bool)
	f.counter = 0
}

// Counter returns the current count of objects being tracked in the AllocatedObjects instance.
func (f *AllocatedObjects) Counter() int {
	return f.counter
}

// Remove deletes the specified object from the AllocatedObjects map and updates the counter to reflect the new size.
func (f *AllocatedObjects) Remove(obj IObject) {
	if f.counter > 0 {
		f.counter--
	}
	if obj.Frame() != FrameStatic {
		delete(f.objects, obj)
	}
}

// ReleaseCandidates releases all objects managed within the AllocatedObjects, resets internal data, and returns the released objects.
func (f *AllocatedObjects) ReleaseCandidates() []IObject {
	if len(f.objects) == 0 {
		return nil
	}
	objectsToRelease := make([]IObject, 0, len(f.objects))
	for obj := range f.objects {
		objectsToRelease = append(objectsToRelease, obj)
	}
	return objectsToRelease
}

// add adds the given IObject to the allocated objects map and updates the counter to reflect the current size.
func (f *AllocatedObjects) add(obj IObject) any {
	f.counter++
	if obj.Frame() != FrameStatic {
		f.objects[obj] = true
	}
	return obj
}

// NewBool creates a new Bool instance with a static frame and a default Code of false, handling object acquisition logic.
func (f *AllocatedObjects) NewBool() any {
	return f.add(newBool(NewAllocator(f.gk, FrameStatic), false))
}

// NewChar creates or retrieves a Char object from the pool, initializing it with the specified gatekeeper and frame Code.
func (f *AllocatedObjects) NewChar() any {
	return f.add(newChar(NewAllocator(f.gk, FrameStatic), 0))
}

// NewInt creates a new integer object or returns an undefined Code if object acquisition fails.
func (f *AllocatedObjects) NewInt() any {
	return f.add(newInt(NewAllocator(f.gk, FrameStatic), 0))
}

// NewFloat creates and retrieves a reusable floating-point object from the pool, or returns a default Code if unavailable.
func (f *AllocatedObjects) NewFloat() any {
	return f.add(newFloat(NewAllocator(f.gk, FrameStatic), 0))
}

// NewString creates a new String object using the IGateKeeper instance, static frame identifier, and an empty Code.
func (f *AllocatedObjects) NewString() any {
	return f.add(newString(NewAllocator(f.gk, FrameStatic), ""))
}

// NewTime attempts to acquire a new reusable time object or returns an undefined Code upon failure.
func (f *AllocatedObjects) NewTime() any {
	return f.add(newTime(NewAllocator(f.gk, FrameStatic), time.Now()))
}

// NewObjectPointer creates a new object pointer using the gatekeeper and returns it or undefined Code on error.
func (f *AllocatedObjects) NewObjectPointer() any {
	return f.add(newObjectPointer(NewAllocator(f.gk, FrameStatic), nil))
}

// NewError attempts to acquire and return a new error object; if unavailable, it falls back to the undefined Code.
func (f *AllocatedObjects) NewError() any {
	return f.add(newError(NewAllocator(f.gk, FrameStatic), ""))
}

// NewBytes creates and returns a new Bytes object, ensuring the byte slice len does not exceed the maximum allowed size.
func (f *AllocatedObjects) NewBytes() any {
	return f.add(newBytes(NewAllocator(f.gk, FrameStatic), []byte{}))
}

// NewArray allocates and returns a new array object, using the internal GateKeeper and predefined configuration.
func (f *AllocatedObjects) NewArray() any {
	return f.add(newArray(NewAllocator(f.gk, FrameStatic), []IObject{}))
}

// NewMap initializes and returns a pooled Map instance or the undefined Code if the acquisition fails.
func (f *AllocatedObjects) NewMap() any {
	return f.add(newMap(NewAllocator(f.gk, FrameStatic), make(map[string]IObject)))
}

// NewStruct acquires and returns a new struct instance from the object pool or an undefined Code if allocation fails.
func (f *AllocatedObjects) NewStruct() any {
	return f.add(newStruct(NewAllocator(f.gk, FrameStatic), make(map[string]IObject)))
}

// NewAny initializes and tracks a new object of any type using the internal allocator and static frame configuration.
func (f *AllocatedObjects) NewAny() any {
	return f.add(newAny(NewAllocator(f.gk, FrameStatic), nil))
}

// NewInterface attempts to allocate a new object; defaults to undefinedValue if allocation fails.
func (f *AllocatedObjects) NewInterface() any {
	return f.add(newInterface(NewAllocator(f.gk, FrameStatic), f.gk.UndefinedValue(), make(map[string]IObject)))
}

// NewChan allocates and returns a new chan object.
func (f *AllocatedObjects) NewChan() any {
	return f.add(newChan(NewAllocator(f.gk, FrameStatic), 0))
}

// NewArrayIterator initializes and returns a new ArrayIterator for iterating over a given slice of IObject elements.
func (f *AllocatedObjects) NewArrayIterator() any {
	return f.add(newArrayIterator(NewAllocator(f.gk, FrameStatic), []IObject{}, 0))
}

// NewBytesIterator creates a new instance of a BytesIterator with a static frame and an empty byte slice.
// It acquires the instance if resources are available; otherwise, returns the undefined Code.
func (f *AllocatedObjects) NewBytesIterator() any {
	return f.add(newBytesIterator(NewAllocator(f.gk, FrameStatic), []byte{}, 0))
}

// NewStringIterator initializes and returns a new StringIterator object for traversing over characters of a string.
func (f *AllocatedObjects) NewStringIterator() any {
	return f.add(newStringIterator(NewAllocator(f.gk, FrameStatic), []rune{}, 0))
}

// NewMapIterator initializes and returns a new map iterator, acquiring necessary resources or returning undefined if unavailable.
func (f *AllocatedObjects) NewMapIterator() any {
	return f.add(newMapIterator(NewAllocator(f.gk, FrameStatic), make(map[string]IObject), 0))
}

// NewStructIterator creates and returns an iterator for a structure if allocation is successful; returns an undefined Code otherwise.
func (f *AllocatedObjects) NewStructIterator() any {
	return f.add(newStructIterator(NewAllocator(f.gk, FrameStatic), make(map[string]IObject), 0))
}
