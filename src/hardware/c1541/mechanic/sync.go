package mechanic

import (
	"fmt"
	"github.com/markel1974/c64emu/src/hardware/c1541/disk"
)

// Sync represents a drive mechanism that emulates disk operations and manages disk state, motor, and head movement.
type Sync struct {
	void        disk.IDisk
	disk        disk.IDisk
	diskChanged bool
	motor       bool
	headPos     uint8
	writing     bool
}

// NewSync creates a new instance of Mechanic, initializing its state and factory dependencies.
func NewSync() *Sync {
	void := NewVoidDisk()
	j := &Sync{
		void:        void,
		disk:        void,
		diskChanged: false,
		motor:       false,
		headPos:     2,
		writing:     false,
	}
	return j
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
	if j.headPos <= 2 {
		return
	}
	j.headPos--
	j.updateHeadPos(j.headPos)
}

// MoveHeadIn increments the head position of the Mechanic unless it has reached the maximum limit (headHalfStep).
func (j *Sync) MoveHeadIn() {
	if j.headPos >= headHalfStep {
		return
	}
	j.headPos++
	j.updateHeadPos(j.headPos)
}

// updateHeadPos adjusts the disk head position by setting the half-track and recalculating rotation intervals.
func (j *Sync) updateHeadPos(headPos uint8) {
	j.disk.SetHeadHalfTrack(headPos)
	fmt.Printf("MOVE HEAD %d: %f\n", headPos/2, j.disk.MicroSecPerByte())
}
