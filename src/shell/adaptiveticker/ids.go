package adaptiveticker

import (
	"container/list"
	"sync"
)

// UnknownId is a constant representing an invalid or uninitialized ID, typically used as a default or error value.
const UnknownId = -1

// IIds represents an interface that defines a method for setting an identifier for an object.
type IIds interface {
	SetId(int)
}

// Ids manages a set of unique integer IDs with allocation, retrieval, and deallocation functionality.
type Ids struct {
	max     int
	kv      map[int]*list.Element
	ll      *list.List
	freeIds []int
	lock    sync.RWMutex
}

// NewIds initializes and returns a new instance of Ids with a specified maximum capacity for unique identifiers.
// If the provided max is less than or equal to zero, a default capacity of 1024 is used.
// The returned Ids structure includes pre-initialized maps, lists, and free ID slices for managing identifiers efficiently.
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

// Get retrieves the element associated with the given id. It returns the element and true if found, otherwise nil and false.
func (a *Ids) Get(id int) (IIds, bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	element, ok := a.kv[id]
	if !ok {
		return nil, false
	}
	return element.Value.(IIds), true
}

// Unset removes an object from the Ids collection, marks the id as free, and resets the object's id to UnknownId.
// Returns true if the id was successfully removed; otherwise, returns false.
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

// Len returns the number of elements currently stored in the Ids structure. It is safe for concurrent use.
func (a *Ids) Len() int {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return len(a.kv)
}

// Cap returns the maximum capacity of IDs that can be managed by the Ids structure.
func (a *Ids) Cap() int {
	return a.max
}

// Range iterates over all elements in the Ids list, invoking the provided function for each element until it returns false.
func (a *Ids) Range(f func(obj IIds) bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	for e := a.ll.Front(); e != nil; e = e.Next() {
		if !f(e.Value.(IIds)) {
			break
		}
	}
}

// All returns a slice of all elements currently stored in the Ids list, thread-safe for concurrent access.
func (a *Ids) All() []IIds {
	a.lock.RLock()
	defer a.lock.RUnlock()
	out := make([]IIds, 0, a.ll.Len())
	for e := a.ll.Front(); e != nil; e = e.Next() {
		out = append(out, e.Value.(IIds))
	}
	return out
}
