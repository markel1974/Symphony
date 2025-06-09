package mechanic

import (
	"github.com/markel1974/c64emu/src/hardware/c1541/disk"
)

const (
	notReady = 1
)

// Async represents the main handler for managing disk mechanics and operations including reading and writing data.
type Async struct {
	void            disk.IDisk
	disk            disk.IDisk
	diskChanged     bool
	filePath        string
	motor           bool
	headPos         uint8
	writing         bool
	motorSpinUpTime int
	headSeekTime    int     // Contatore per il delay del movimento testina (in µs)
	timeToNextByte  float64 // Contatore per il tempo che manca alla lettura/scrittura del prossimo byte (in µs)
	dataWrite       int
	dataRead        int
	syncCounter     int
}

// NewAsync initializes and returns a new instance of Mechanic with default values and a void disk.
func NewAsync() *Async {
	void := NewVoidDisk()
	j := &Async{
		void:            void,
		disk:            void,
		diskChanged:     false,
		filePath:        "",
		motor:           false,
		headPos:         2,
		writing:         false,
		syncCounter:     0,
		motorSpinUpTime: 0,
		headSeekTime:    0,
		timeToNextByte:  0,
		dataWrite:       notReady,
		dataRead:        notReady,
	}
	return j
}

// Reset reinitializes the Mechanic's state, clearing data, resetting counters, and updating the disk head position.
func (j *Async) Reset() {
	j.diskChanged = false
	j.filePath = ""
	j.motor = false
	j.headPos = 2
	j.writing = false
	j.syncCounter = 0
	j.motorSpinUpTime = 0
	j.headSeekTime = 0
	j.timeToNextByte = 0
	j.dataWrite = notReady
	j.dataRead = notReady
	j.disk.SetHeadHalfTrack(j.headPos)
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
	return j.InsertDisk(j.void)
}

// SetWrite sets the writing state of the Mechanic object. When true, the Mechanic operates in writing mode.
func (j *Async) SetWrite(w bool) {
	j.writing = w
}

// EmulationRequired returns true if the Mechanic is currently emulating disk operations
func (j *Async) EmulationRequired() bool {
	return true
}

// Emulate performs the main emulation logic for the Mechanic, handling motor operation, head movement, and disk I/O.
func (j *Async) Emulate() {
	if !j.motor {
		return
	}
	isBusy := false
	if j.motorSpinUpTime > 0 {
		j.motorSpinUpTime--
		isBusy = true
	}
	if j.headSeekTime > 0 {
		j.headSeekTime--
		isBusy = true
	}
	if isBusy {
		return
	}
	j.timeToNextByte--
	if j.timeToNextByte > 0 {
		return
	}
	j.timeToNextByte += j.disk.MicroSecPerByte()
	if j.writing {
		if j.dataWrite != notReady {
			j.disk.Write(uint8(j.dataWrite))
			j.dataWrite = notReady
		}
	} else {
		current := j.disk.Read()
		if current == syncByte {
			j.syncCounter++
		} else {
			j.syncCounter = 0
		}
		j.dataRead = int(current)
	}
	j.disk.Rotate()
}

// WriteByte sets the byte value to be written to the disk by assigning it to the `dataWrite` field of the Mechanic instance.
func (j *Async) WriteByte(data uint8) {
	j.dataWrite = int(data)
}

// ReadByte retrieves the next byte of data from the Mechanic if the motor is active. Returns 0 if the motor is off.
func (j *Async) ReadByte() uint8 {
	if !j.motor {
		return 0
	}
	if j.dataRead == notReady {
		return 0
	}
	v := uint8(j.dataRead)
	j.dataRead = notReady
	return v
}

// ByteReady returns true if the mechanic's system is ready to read or write the next byte of data.
func (j *Async) ByteReady() bool {
	if j.writing {
		return j.dataWrite == notReady
	}
	return j.dataRead != notReady
}

// SyncFound checks if the mechanic has detected a synchronization state based on motor status and sync counter value.
func (j *Async) SyncFound() bool {
	if !j.motor {
		return true
	}
	if j.syncCounter >= syncTolerance {
		return true
	}
	return false
}

// SetMotor controls the state of the motor. Enables spin-up delay if turning on from an off state.
func (j *Async) SetMotor(m bool) {
	if m && !j.motor {
		j.motorSpinUpTime = motorSpinUpDelay
	}
	j.motor = m
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
	if j.headPos <= 2 {
		return
	}
	j.headPos--
	j.updateHeadPos(j.headPos)
}

// MoveHeadIn increments the head position by one step unless it is already at or beyond the maximum position.
func (j *Async) MoveHeadIn() {
	if j.headPos >= headHalfStep {
		return
	}
	j.headPos++
	j.updateHeadPos(j.headPos)
}

// updateHeadPos updates the disk head to the given half-track position and resets related state variables.
func (j *Async) updateHeadPos(headPos uint8) {
	j.disk.SetHeadHalfTrack(headPos)
	j.headSeekTime = stepDelay / 2
	j.dataRead = notReady
	j.syncCounter = 0
	//fmt.Printf("MOVE HEAD %d: %f\n", headPos/2, j.disk.MicroSecPerByte())
}
