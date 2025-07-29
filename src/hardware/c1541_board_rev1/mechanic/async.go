package mechanic

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/hardware/c1541_board_rev1/disk"
	"github.com/markel1974/c64emu/src/hardware/c1541_board_rev1/disk/void"
	"github.com/markel1974/c64emu/src/references"
)

// Async represents the main handler for managing disk mechanics and operations including reading and writing data.
type Async struct {
	*component.BaseComponent
	empty           disk.IDisk
	disk            disk.IDisk
	diskChanged     bool
	headPosRequired uint8
	head            *Head
	motor           *Motor
}

// NewAsync creates and initializes a new Async instance with a parent component, factory, label, and instance number.
func NewAsync(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Async {
	voidDisk := void.NewVoidDisk()
	j := &Async{
		BaseComponent: component.NewBaseComponent(),
		empty:         voidDisk,
		disk:          voidDisk,
		diskChanged:   false,
	}
	j.BaseComponent.Register(factory, parent, "async", instance, j, references.IdInternalComponent(label, instance, "Async"))

	j.motor = NewMotor(j, j.GetFactory(), label, 0)
	j.head = NewHead(j, j.GetFactory(), label, 0)
	j.headPosRequired = j.head.DefaultPos()

	return j
}

// Connect initializes a connection for the Async object, preparing it for operations. Returns an error if connection fails.
func (j *Async) Connect() error {
	return nil
}

// Internal indicates if the component operates in an internal state. Returns true when in internal mode.
func (j *Async) Internal() bool {
	return true
}

// Reset reinitializes the Mechanic's state, clearing data, resetting counters, and updating the disk head position.
func (j *Async) Reset() {
	j.diskChanged = false
	j.headPosRequired = 2
	j.motor.Reset()
	j.head.Reset()
	j.disk.SetHeadHalfTrack(j.headPosRequired)
}

// Setup resets the Mechanic state to its default values and ensures proper initialization.
func (j *Async) Setup() error {
	j.Reset()
	return nil
}

// InsertDisk inserts a new disk into the Mechanic, resets its state, and marks the disk as changed. Returns an error if any occurs.
func (j *Async) InsertDisk(d disk.IDisk) error {
	j.diskChanged = true
	j.Reset()
	j.disk = d
	return nil
}

// RemoveDisk removes the currently inserted disk by replacing it with a void disk, effectively resetting the drive state.
func (j *Async) RemoveDisk() error {
	return j.InsertDisk(j.empty)
}

// SetWrite sets the writing state of the Mechanic object. When true, the Mechanic operates in writing mode.
func (j *Async) SetWrite(w bool) {
	j.head.SetWritingMode(w)
}

// EmulationRequired returns true if the Mechanic is currently emulating disk operations
func (j *Async) EmulationRequired() bool {
	return true
}

// Emulate performs the main emulation logic for the Mechanic, handling motor operation, head movement, and disk I/O.
func (j *Async) Emulate() {
	if !j.motor.IsActive() {
		return
	}
	if j.motor.SpinUp() {
		return
	}
	if j.head.Move(j.disk, j.headPosRequired) {
		return
	}
	j.head.DecayVibration()
	if !j.motor.TryRotate() {
		return
	}
	j.head.ReadWrite(j.disk)
	j.motor.Rotate(j.disk)
}

// WriteByte sets the byte value to be written to the disk by assigning it to the `dataWrite` field of the Mechanic instance.
func (j *Async) WriteByte(data uint8) {
	j.head.WriteByte(data)
}

// ReadByte retrieves the next byte of data from the Mechanic if the motor is active. Returns 0 if the motor is off.
func (j *Async) ReadByte() uint8 {
	return j.head.ReadByte()
}

// ByteReady returns true if the mechanic's system is ready to read or write the next byte of data.
func (j *Async) ByteReady() bool {
	return j.head.ByteReady()
}

// SyncFound checks if the mechanic has detected a synchronization state based on motor status and sync counter value.
func (j *Async) SyncFound() bool {
	if !j.motor.IsActive() {
		return true
	}
	return j.head.HasSync()
}

// SetMotor controls the state of the motor. Enables spin-up delay if turning on from an off state.
func (j *Async) SetMotor(m bool) {
	j.motor.SetActive(m)
}

// HasDisk checks if a usable disk is currently inserted into the mechanic and returns true if so, or false otherwise.
func (j *Async) HasDisk() bool {
	return j.disk.Usable()
}

// WriteProtectionState returns the write protection state of the disk as a uint8 value.
func (j *Async) WriteProtectionState() uint8 {
	const wp = 0x10
	if !j.diskChanged {
		if !j.disk.WriteProtected() {
			return wp
		}
		return 0
	}
	j.diskChanged = false
	if j.disk.WriteProtected() {
		return wp
	}
	return 0
}

// MoveHeadOut moves the head outward by decrementing the position if it is greater than the minimum limit (2).
func (j *Async) MoveHeadOut() {
	if j.headPosRequired <= headMinHalfStep {
		return
	}
	j.headPosRequired--
}

// MoveHeadIn increments the head position by one step unless it is already at or beyond the maximum position.
func (j *Async) MoveHeadIn() {
	if j.headPosRequired >= headMaxHalfStep {
		return
	}
	j.headPosRequired++
}
