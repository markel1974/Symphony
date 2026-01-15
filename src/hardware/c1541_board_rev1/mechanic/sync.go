package mechanic

import (
	"fmt"

	"github.com/markel1974/symphony/src/hardware/c1541_board_rev1/disk"
	"github.com/markel1974/symphony/src/hardware/c1541_board_rev1/disk/void"
	"github.com/markel1974/symphony/src/kernel/component"
	"github.com/markel1974/symphony/src/references"
)

// Sync represents a drive mechanism that emulates disk operations and manages disk state, motor, and head movement.
type Sync struct {
	*component.BaseComponent
	void        disk.IDisk
	disk        disk.IDisk
	diskChanged bool
	motor       bool
	headPos     uint8
	writing     bool
}

// NewSync creates a new instance of Mechanic, initializing its state and factory dependencies.
func NewSync(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Sync {
	voidDisk := void.NewVoidDisk()
	j := &Sync{
		BaseComponent: component.NewBaseComponent(),
		void:          voidDisk,
		disk:          voidDisk,
		diskChanged:   false,
		motor:         false,
		headPos:       2,
		writing:       false,
	}
	j.BaseComponent.Register(factory, parent, "sync", instance, j, references.IdInternalComponent(label, instance, "Sync"))
	return j
}

// Connect establishes a connection to the drive mechanism and prepares it for operation. Returns an error if the operation fails.
func (j *Sync) Connect() error {
	return nil
}

// Internal returns true if the Sync mechanism is in an internal operational state.
func (j *Sync) Internal() bool {
	return true
}

// Reset restores the Mechanic to its initial state, clearing all internal state and resetting properties to their defaults.
func (j *Sync) Reset() {
	j.diskChanged = false
	j.motor = false
	j.headPos = 2
	j.writing = false
	j.updateHeadPos(j.headPos)
}

// Setup initializes the Mechanic instance by invoking the init method with the provided file path.
func (j *Sync) Setup() error {
	j.Reset()
	return nil
}

// InsertDisk initializes the Mechanic object by resetting its state and inserting a disk from the provided file path.
func (j *Sync) InsertDisk(d disk.IDisk) error {
	j.diskChanged = true
	j.Reset()
	j.disk = d
	return nil
}

func (j *Sync) RemoveDisk() error {
	return j.InsertDisk(j.void)
}

func (j *Sync) SetWrite(w bool) {
	j.writing = w
}

// EmulationRequired returns true if the Mechanic is currently emulating disk operations
func (j *Sync) EmulationRequired() bool {
	return false
}

// Emulate advances the disk's rotation and reads data while checking synchronization with the sync byte.
func (j *Sync) Emulate() {
}

// ReadByte reads a byte from the current disk, rotates the disk, and returns the read value.
func (j *Sync) ReadByte() uint8 {
	data := j.disk.Read()
	j.disk.Rotate()
	return data
}

// WriteByte writes a single byte to the disk and rotates the disk to the next position.
func (j *Sync) WriteByte(data uint8) {
	j.disk.Write(data)
	j.disk.Rotate()
}

// SyncFound checks if the disk drive's current position aligns to a synchronization byte sequence and returns true if found.
func (j *Sync) SyncFound() bool {
	if !j.motor {
		return true
	}
	j.disk.Rotate()
	if (j.disk.Read() == syncByte) && (j.disk.Next() != syncByte) {
		return true
	}
	return false
}

func (j *Sync) ByteReady() bool {
	return true
}

// SetMotor toggles the motor state, resets the data to 0, and clears the sync status.
func (j *Sync) SetMotor(m bool) {
	j.motor = m
}

// HasDisk returns true if a usable disk is present, otherwise false.
func (j *Sync) HasDisk() bool {
	return j.disk.Usable()
}

// WriteProtectionState returns the write protection state of the current disk as a uint8 value.
func (j *Sync) WriteProtectionState() uint8 {
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

// MoveHeadOut decreases the `headPosRequired` of the Mechanic by one step, ensuring it doesn't go below the minimum allowable position.
func (j *Sync) MoveHeadOut() {
	if j.headPos <= headMinHalfStep {
		return
	}
	j.updateHeadPos(j.headPos - 1)
}

// MoveHeadIn increments the head position of the Mechanic unless it has reached the maximum limit (headHalfStep).
func (j *Sync) MoveHeadIn() {
	if j.headPos >= headMaxHalfStep {
		return
	}
	j.updateHeadPos(j.headPos + 1)
}

// updateHeadPos adjusts the disk head position by setting the half-track and recalculating rotation intervals.
func (j *Sync) updateHeadPos(headPos uint8) bool {
	fmt.Printf("MOVE HEAD %d: %d\n", headPos/2, j.disk.MicroSecPerByte())
	if j.disk.SetHeadHalfTrack(headPos) {
		j.headPos = headPos
		return true
	}
	return false
}
