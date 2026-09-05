package objects

import (
	"reflect"
)

// ChanType represents the constant type name for channel objects.
const (
	ChanType = "chan"
)

// ChanData represents a data structure for managing channel buffers and queues for send and receive operations.
type ChanData struct {
	buffer    []IObject
	capacity  int
	sendQueue []uint // core IDs waiting to send
	recvQueue []uint // core IDs waiting to receive
}

// AddBuffer appends a value to the buffer if capacity allows, returning true on success and false if the buffer is full.
func (data *ChanData) AddBuffer(val IObject) bool {
	if len(data.buffer) >= data.capacity {
		return false
	}
	data.buffer = append(data.buffer, val)
	return true
}

// GetBuffer retrieves the first element from the buffer if available, removing it; returns nil if the buffer is empty.
func (data *ChanData) GetBuffer() IObject {
	if len(data.buffer) == 0 {
		return nil
	}
	val := data.buffer[0]
	data.buffer = data.buffer[1:]
	return val
}

// GetSend retrieves and removes the first sender ID from the send queue.
// It returns the sender ID and true if the queue is not empty, otherwise 0 and false.
func (data *ChanData) GetSend() (uint, bool) {
	if len(data.sendQueue) == 0 {
		return 0, false
	}
	wakeId := data.sendQueue[0]
	data.sendQueue = data.sendQueue[1:]
	return wakeId, true
}

// GetRecv retrieves and removes the first core ID from the receive queue. Returns the ID and a boolean indicating success.
func (data *ChanData) GetRecv() (uint, bool) {
	if len(data.recvQueue) == 0 {
		return 0, false
	}
	wakeId := data.recvQueue[0]
	data.recvQueue = data.recvQueue[1:]
	return wakeId, true
}

// AddRecv appends the given core ID to the recvQueue, indicating it is waiting to receive from the channel.
func (data *ChanData) AddRecv(id uint) {
	data.recvQueue = append(data.recvQueue, id)
}

// AddSend appends the given core ID to the sendQueue, indicating it is waiting to send to the channel.
func (data *ChanData) AddSend(id uint) {
	data.sendQueue = append(data.sendQueue, id)
}

// Chan represents a memory-managed channel structure with send/receive queues and a buffer for object storage.
type Chan struct {
	IAllocator
	data *ChanData
}

// newChan creates and returns a new channel object with the specified allocator and capacity.
func newChan(allocator IAllocator, capacity int) IObject {
	return &Chan{
		IAllocator: allocator,
		data: &ChanData{
			buffer:    make([]IObject, 0, capacity),
			capacity:  capacity,
			sendQueue: make([]uint, 0),
			recvQueue: make([]uint, 0),
		},
	}
}

// setAllocator sets the memory allocator for the channel instance to manage its lifecycle and memory operations.
func (o *Chan) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// Setup initializes or resets the Chan object with the given frame and capacity, preparing it for use or reuse.
func (o *Chan) Setup(frame int, capacity int) {
	o.setFrame(frame)
	if o.data == nil {
		o.data = &ChanData{
			buffer:    make([]IObject, 0, capacity),
			capacity:  capacity,
			sendQueue: make([]uint, 0),
			recvQueue: make([]uint, 0),
		}
	} else {
		// Reset state if pooled
		o.data.buffer = make([]IObject, 0, capacity)
		o.data.capacity = capacity
		o.data.sendQueue = o.data.sendQueue[:0]
		o.data.recvQueue = o.data.recvQueue[:0]
	}
}

// AsInterface returns the current instance as an empty interface, allowing for generic usage and type abstraction.
func (o *Chan) AsInterface() interface{} {
	return o
}

// AsValue converts the Chan object into a reflect.Value of the specified type and returns a boolean indicating success.
func (o *Chan) AsValue(target reflect.Type) (reflect.Value, bool) {
	return _reflect(o, target)
}

// AsBool returns the boolean representation of the Chan object, which is always true.
func (o *Chan) AsBool() bool {
	return true
}

// AsInt64 converts and returns the Chan object as a 64-bit integer.
func (o *Chan) AsInt64() int64 {
	return 0
}

// AsFloat64 converts the Chan object into its 64-bit floating-point representation and returns the result.
func (o *Chan) AsFloat64() float64 {
	return 0
}

// AsBytes converts the Chan object to a byte slice representation and returns it.
func (o *Chan) AsBytes() []byte {
	return nil
}

// AsString returns a string representation of the Chan object, which is always "[chan]".
func (o *Chan) AsString() string {
	return "[chan]"
}

// AssignValue attempts to assign a value to the Chan object but always returns ErrNotAssignable.
func (o *Chan) AssignValue(v IObject) error {
	return ErrNotAssignable
}

// Nil returns false, indicating the channel object is not nil or uninitialized.
func (o *Chan) Nil() bool {
	return false
}

// Call invokes the Chan instance as a callable, returning an error since this operation is not supported.
func (o *Chan) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, ErrInvalidOperator
}

// TypeName returns the type name of the object as a string, indicating it is a channel.
func (o *Chan) TypeName() string {
	return ChanType
}

// LogicalOp performs a logical operation using the specified LogicalOperator and right-hand-side IObject.
// Returns an error if the operator is invalid.
func (o *Chan) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation using the specified operator and right-hand-side operand and returns the result.
func (o *Chan) ArithmeticOp(frame int, op ArithmeticOperator, rhsIn IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// UnaryOp performs a unary operation on the Chan object and returns the resulting IObject or an ErrInvalidOperator.
func (o *Chan) UnaryOp(_ int, _ UnaryOperator) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Copy creates a "copy" of the channel object, returning the same instance as channels are not deeply copied.
func (o *Chan) Copy(frame int, depth int) IObject {
	// Channels are passed by reference, returning self (with correct frame is tricky, normally channels don't deep copy)
	return o
}

// Falsy determines if the channel evaluates to a falsy value; always returns false for Chan instances.
func (o *Chan) Falsy() bool {
	return false
}

// Equals determines if the provided IObject is of type *Chan and has the same underlying ChanData.
func (o *Chan) Equals(in IObject) bool {
	if other, ok := in.(*Chan); ok {
		return o.data == other.data
	}
	return false
}

// IndexGet retrieves an element from the Chan using the provided index and returns an error if the operation is invalid.
func (o *Chan) IndexGet(_ int, index IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// IndexSet attempts to assign a value to the specified index on the channel and returns ErrInvalidOperator on failure.
func (o *Chan) IndexSet(index IObject, value IObject) error {
	return ErrInvalidOperator
}

// Count returns the number of elements currently stored in the channel's buffer.
func (o *Chan) Count() int {
	return len(o.data.buffer)
}

// Iterable checks if the Chan object supports iteration and returns false to indicate it is not iterable.
func (o *Chan) Iterable() bool {
	return false
}

// Iterate returns an IIterator instance for traversing elements, using the specified frame for any contextual operations.
func (o *Chan) Iterate(frame int) IIterator {
	return o.GateKeeper().UndefinedValue().(IIterator)
}

// Length returns the capacity of the channel as an integer.
func (o *Chan) Length() int {
	return o.data.capacity
}

// Data returns the internal data structure (ChanData) associated with the channel instance.
func (o *Chan) Data() *ChanData {
	return o.data
}
