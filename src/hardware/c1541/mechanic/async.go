package mechanic

import (
	"github.com/markel1974/c64emu/src/hardware/c1541/disk"
)

const (
	directionNone    = 0
	directionInward  = 1  // Verso tracce più alte
	directionOutward = -1 // Verso tracce più basse
)

const (
	notReady = 1
)

type Motor struct {
	active     bool
	spinUpTime int
}

func NewMotor() *Motor {
	v := &Motor{
		active:     false,
		spinUpTime: 0,
	}
	v.Reset()
	return v
}

func (m *Motor) Reset() {
	m.active = false
	m.spinUpTime = 0
}

func (m *Motor) IsActive() bool {
	return m.active
}

func (m *Motor) SetActive(active bool) {
	if active && !m.active {
		m.spinUpTime = motorSpinUpDelay
	}
	m.active = active
}

func (m *Motor) SpinUp() bool {
	if m.spinUpTime > 0 {
		m.spinUpTime--
		return true
	}
	return false
}

type Head struct {
	defaultPos       uint8
	currentPos       uint8
	consecutiveSteps int
	direction        int
	seekTime         int
	writing          bool
	dataWrite        int
	dataRead         int
	syncCounter      int
}

func NewHead(defaultPos uint8) *Head {
	h := &Head{
		defaultPos:       defaultPos,
		currentPos:       defaultPos,
		consecutiveSteps: 0,
		seekTime:         0,
		direction:        directionNone,
		writing:          false,
		syncCounter:      0,
		dataWrite:        notReady,
		dataRead:         notReady,
	}
	return h
}

func (h *Head) Reset() {
	h.currentPos = h.defaultPos
	h.consecutiveSteps = 0
	h.seekTime = 0
	h.direction = directionNone
	h.writing = false
	h.syncCounter = 0
	h.dataWrite = notReady
	h.dataRead = notReady
}

func (h *Head) SetWritingMode(w bool) {
	h.writing = w
}

func (h *Head) WriteByte(d uint8) {
	h.dataWrite = int(d)
}

func (h *Head) ReadByte() uint8 {
	if h.dataRead == notReady {
		return 0
	}
	v := uint8(h.dataRead)
	//fmt.Printf("ReadByte %d From Track %d\n", v, j.currentPos)
	h.dataRead = notReady
	return v
}

func (h *Head) ReadWrite(disk disk.IDisk) {
	if h.writing {
		if h.dataWrite != notReady {
			disk.Write(uint8(h.dataWrite))
			h.dataWrite = notReady
		}
	} else {
		current := disk.Read()
		if current == syncByte {
			h.syncCounter++
		} else {
			h.syncCounter = 0
		}
		if h.dataRead == notReady {
			h.dataRead = int(current)
		}
	}
}

func (h *Head) ByteReady() bool {
	if h.writing {
		return h.dataWrite == notReady
	}
	v := h.dataRead != notReady
	return v
}

func (h *Head) HasSync() bool {
	return h.syncCounter >= syncTolerance
}

func (h *Head) Move(disk disk.IDisk, headPosRequired uint8) bool {
	if h.seekTime > 0 {
		h.seekTime--
		return true
	}
	if headPosRequired == h.currentPos {
		h.consecutiveSteps = 0
		h.direction = directionNone
		return false
	}
	newPos := h.currentPos
	var direction int
	var polarity int
	var seekTime int
	if isMovingInward := headPosRequired > h.currentPos; isMovingInward {
		newPos++
		direction = directionInward
		seekTime = headInwardDelay + headBaseDamping
		polarity = headInwardPolarityDelay
	} else {
		newPos--
		direction = directionOutward
		seekTime = headOutwardDelay + headBaseDamping
		polarity = headBackwardPolarityDelay
	}

	if h.direction != direction {
		h.consecutiveSteps = 0
		if h.direction != directionNone {
			seekTime += headBacklashDelay + polarity
		}
	}

	if h.consecutiveSteps++; h.consecutiveSteps > 1 {
		seekTime += (h.consecutiveSteps - 1) * headExtraSettlingPerStep
	}
	if seekTime > headMaxDelay {
		seekTime = headMaxDelay
	}

	//fmt.Printf("ASYNC MOVE HEAD OLD %d NEW %d REQ %d: %d\n", h.currentPos, newPos, headPosRequired, disk.MicroSecPerByte())
	if disk.SetHeadHalfTrack(newPos) {
		h.currentPos = newPos
		h.seekTime = seekTime
		h.direction = direction
		h.dataRead = notReady
		h.syncCounter = 0
	}
	return true
}

// Async represents the main handler for managing disk mechanics and operations including reading and writing data.
type Async struct {
	empty           disk.IDisk
	disk            disk.IDisk
	diskChanged     bool
	rotationCycles  int
	headPosRequired uint8
	head            *Head
	motor           *Motor
}

// NewAsync initializes and returns a new instance of Mechanic with default values and a void disk.
func NewAsync() *Async {
	void := NewVoidDisk()
	j := &Async{
		empty:           void,
		disk:            void,
		diskChanged:     false,
		rotationCycles:  0,
		motor:           NewMotor(),
		head:            NewHead(headMinHalfStep),
		headPosRequired: headMinHalfStep,
	}
	return j
}

// Reset reinitializes the Mechanic's state, clearing data, resetting counters, and updating the disk head position.
func (j *Async) Reset() {
	j.diskChanged = false
	j.headPosRequired = 2
	j.rotationCycles = 0
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
	j.rotationCycles--
	if j.rotationCycles > 0 {
		return
	}
	j.rotationCycles += j.disk.MicroSecPerByte()
	j.head.ReadWrite(j.disk)
	j.disk.Rotate()
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

/*
func (h *Head) MoveOld(disk disk.IDisk, headPosRequired uint8) bool {
	if h.seekTime > 0 {
		h.seekTime--
		return true
	}

	if headPosRequired == h.currentPos {
		return false
	}
	newPos := h.currentPos
	var seekTime = 0
	if headPosRequired > h.currentPos {
		newPos++
		seekTime = headInwardDelay
	} else {
		newPos--
		seekTime = headOutwardDelay
	}

	fmt.Printf("ASYNC MOVE HEAD OLD %d NEW %d: %d\n", h.currentPos, newPos, disk.MicroSecPerByte())
	disk.SetHeadHalfTrack(newPos)
	h.currentPos = newPos
	h.seekTime = seekTime
	h.dataRead = notReady
	h.syncCounter = 0
	return true
}
*/
