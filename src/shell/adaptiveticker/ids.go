package adaptiveticker

import (
	"container/list"
	"sync"
)

// UnknownId is a constant representing an invalid or uninitialized identifier with a value of -1.
const UnknownId = -1

// IIds represents an interface requiring a method to set an integer ID for implementing types.
type IIds interface {
	SetId(int)
}

// Ids manages a pool of reusable integer IDs and their associated objects with thread-safe access.
type Ids struct {
	max     int
	kv      map[int]*list.Element
	ll      *list.List
	freeIds []int
	lock    sync.RWMutex
}

// NewIds initializes and returns a pointer to an Ids structure with a specified maximum capacity or a default of 1024.
func NewIds(max int) *Ids {
	if max <= 0 {
		max = 1024
	}
	free := make([]int, max)
	for i := 0; i < max; i++ {
		free[i] = i
	}
	return &Ids{
		max:     max,
		kv:      make(map[int]*list.Element, max),
		ll:      list.New(),
		freeIds: free,
	}
}

// Set assigns a free ID to the provided object and stores it in the internal structures, returning true on success.
func (a *Ids) Set(obj IIds) bool {
	a.lock.Lock()
	defer a.lock.Unlock()
	if len(a.freeIds) == 0 {
		obj.SetId(UnknownId)
		return false
	}
	lastIndex := len(a.freeIds) - 1
	id := a.freeIds[lastIndex]
	a.freeIds = a.freeIds[:lastIndex]
	obj.SetId(id)
	element := a.ll.PushBack(obj)
	a.kv[id] = element
	return true
}

// Get retrieves an element by its ID from the map. Returns the element and true if found, otherwise nil and false.
func (a *Ids) Get(id int) (IIds, bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	element, ok := a.kv[id]
	if !ok {
		return nil, false
	}
	return element.Value.(IIds), true
}

// Unset removes an ID and its associated object, marking the ID as free and resetting the object's ID. Returns true on success.
func (a *Ids) Unset(id int) bool {
	a.lock.Lock()
	defer a.lock.Unlock()
	if id < 0 || id >= a.max {
		return false
	}
	element, ok := a.kv[id]
	if !ok {
		return false
	}
	obj := element.Value.(IIds)
	a.ll.Remove(element)
	delete(a.kv, id)
	a.freeIds = append(a.freeIds, id)
	obj.SetId(UnknownId)
	return true
}

// Len returns the number of elements currently stored in the Ids instance. It is safe for concurrent use.
func (a *Ids) Len() int {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return len(a.kv)
}

// Cap returns the maximum number of IDs that can be managed by the Ids instance.
func (a *Ids) Cap() int {
	return a.max
}

// Range iterates over all elements in the list, applying the provided function to each element.
// Stops iteration if the function returns false.
// The method provides read-locking to ensure thread-safety.
func (a *Ids) Range(f func(obj IIds) bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	for e := a.ll.Front(); e != nil; e = e.Next() {
		if !f(e.Value.(IIds)) {
			break
		}
	}
}

// All returns a slice of all IIds elements currently stored in the Ids structure. It is thread-safe for concurrent use.
func (a *Ids) All() []IIds {
	a.lock.RLock()
	defer a.lock.RUnlock()
	out := make([]IIds, 0, a.ll.Len())
	for e := a.ll.Front(); e != nil; e = e.Next() {
		out = append(out, e.Value.(IIds))
	}
	return out
}
