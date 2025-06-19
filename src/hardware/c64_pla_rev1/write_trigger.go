package c64_pla_rev1

// TriggerData is a data structure encapsulating an ID, a 16-bit address, and a write function of type WriteFn.
type TriggerData struct {
	id   int
	addr uint16
	w    WriteFn
}

// NewTriggerData creates a new TriggerData instance with the provided ID, address, and WriteFn.
func NewTriggerData(id int, addr uint16, fn WriteFn) *TriggerData {
	return &TriggerData{
		id:   id,
		addr: addr,
		w:    fn,
	}
}

// Exec executes the WriteFn associated with the TriggerData, passing its address and the provided data as arguments.
func (td *TriggerData) Exec(data uint8) {
	td.w(td.addr, data)
}

// GetId returns the identifier of the TriggerData instance.
func (td *TriggerData) GetId() int {
	return td.id
}

// Trigger represents a write trigger mechanism that holds multiple write functions and an associated address.
// It manages adding, removing, and executing those functions with a specific data payload.
type Trigger struct {
	container []*TriggerData
	addr      uint16
	idx       int
}

// NewTrigger creates and initializes a new Trigger instance with the specified address and an empty container.
func NewTrigger(addr uint16) *Trigger {
	return &Trigger{
		idx:       0,
		addr:      addr,
		container: nil,
	}
}

// Add registers a WriteFn to the Trigger and returns a unique identifier for the registered function.
func (wt *Trigger) Add(fn WriteFn) int {
	id := wt.idx
	wt.container = append(wt.container, NewTriggerData(id, wt.addr, fn))
	wt.idx++
	return id
}

// Remove deletes a TriggerData object from the container by its unique ID, if it exists.
func (wt *Trigger) Remove(id int) {
	for idx, f := range wt.container {
		if id == f.GetId() {
			wt.container = append(wt.container[:idx], wt.container[idx+1:]...)
			break
		}
	}
}

// Exec processes the provided data by invoking the Exec method on each TriggerData in the Trigger's container.
func (wt *Trigger) Exec(data uint8) {
	for _, f := range wt.container {
		f.Exec(data)
	}
}

// WriteTriggers is a structure that manages a collection of write triggers, enabling dynamic execution of write operations.
// Each trigger is associated with a specific address and executes a function when data is written to that address.
type WriteTriggers struct {
	triggers []*Trigger
}

// NewWriteTriggers initializes a WriteTriggers instance with a specified capacity for triggers and returns a pointer to it.
func NewWriteTriggers(r int) *WriteTriggers {
	wt := &WriteTriggers{triggers: nil}
	wt.triggers = make([]*Trigger, r)
	return wt
}

// Add registers a WriteFn for the specified address and returns a unique ID for the added trigger.
func (wt *WriteTriggers) Add(addr uint16, fn WriteFn) int {
	t := wt.triggers[addr]
	if t == nil {
		t = NewTrigger(addr)
		wt.triggers[addr] = t
	}
	return t.Add(fn)
}

// Remove removes a trigger identified by the given id for the specified address from the WriteTriggers collection.
func (wt *WriteTriggers) Remove(addr uint16, id int) {
	if t := wt.triggers[addr]; t != nil {
		t.Remove(id)
	}
}

// Exec triggers the execution of all write functions registered for the specified memory address with the provided data.
func (wt *WriteTriggers) Exec(addr uint16, data uint8) {
	if t := wt.triggers[addr]; t != nil {
		t.Exec(data)
	}
}
