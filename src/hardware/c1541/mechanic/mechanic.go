package mechanic

import (
	"io"
	"os"
)

//see
//http://www.baltissen.org/newhtm/1541c.htm
//https://sta.c64.org/cbm1541mem.html
//https://c64os.com/post/howdoes1541work

// syncByte represents a constant value used for synchronization detection in disk emulation processes.
const syncByte = 0xff

// headStep defines the step value used for calculating head movements in a specific context.
const headStep = 35
const headHalfStep = headStep * 2 // headHalfStep represents the number of half-tracks the head can move, calculated as twice the value of headStep.

// Mechanic represents a drive mechanism that emulates disk operations and manages disk state, motor, and head movement.
type Mechanic struct {
	disk              IDisk
	writeProtected    bool
	diskChanged       bool
	filePath          string
	motor             bool
	empty             IDisk
	factory           *Factory
	headPos           uint8
	rotationIntervals int
	rotationCounter   int
	data              uint8
	sync              bool
}

// NewMechanic creates a new instance of Mechanic, initializing its state and factory dependencies.
func NewMechanic() *Mechanic {
	factory := NewFactory()
	empty, _ := factory.Create(nil)
	j := &Mechanic{
		disk:              empty,
		writeProtected:    false,
		diskChanged:       false,
		filePath:          "",
		motor:             false,
		empty:             empty,
		factory:           factory,
		headPos:           2,
		rotationIntervals: 0,
		rotationCounter:   0,
		data:              0,
		sync:              false,
	}
	return j
}

// Reset restores the Mechanic to its initial state, clearing all internal state and resetting properties to their defaults.
func (j *Mechanic) Reset() {
	j.disk = j.empty
	j.writeProtected = false
	j.diskChanged = false
	j.filePath = ""
	j.motor = false
	j.data = 0
	j.sync = false
	j.headPos = 2
	j.rotationIntervals = 0
	j.rotationCounter = 0
	j.updateHeadPos()
}

// init initializes the Mechanic object by resetting its state and inserting a disk from the provided file path.
func (j *Mechanic) init(fp string) error {
	j.Reset()
	if err := j.insertDisk(fp); err != nil {
		return err
	}
	j.filePath = fp
	return nil
}

// Setup initializes the Mechanic instance by invoking the init method with the provided file path.
func (j *Mechanic) Setup(fp string) {
	if err := j.init(fp); err != nil {
		return
	}
}

// Emulate advances the disk's rotation and reads data while checking synchronization with the sync byte.
func (j *Mechanic) Emulate() {
	if j.motor {
		j.rotationCounter++
		if j.rotationCounter >= j.rotationIntervals {
			j.rotationCounter = 0
			j.disk.Rotate()
			j.data = j.disk.Read()
			next := j.disk.Next()
			if j.data == syncByte && next != syncByte {
				j.sync = true
			} else {
				j.sync = false
			}
		}
	}
}

// ReadByte reads a byte from the current disk, rotates the disk, and returns the read value.
func (j *Mechanic) ReadByte() uint8 {
	data := j.disk.Read()
	j.disk.Rotate()
	return data

	//return j.data
}

// WriteByte writes a single byte to the disk and rotates the disk to the next position.
func (j *Mechanic) WriteByte(data uint8) {
	j.disk.Write(data)
	j.disk.Rotate()
}

// SyncFound checks if the disk drive's current position aligns to a synchronization byte sequence and returns true if found.
func (j *Mechanic) SyncFound() bool {
	if !j.motor {
		return true
	}
	j.disk.Rotate()
	if (j.disk.Read() == syncByte) && (j.disk.Next() != syncByte) {
		return true
	}
	return false

	//if j.sync {
	//	return true
	//}
	//return false
}

// SetMotor toggles the motor state, resets the data to 0, and clears the sync status.
func (j *Mechanic) SetMotor(m bool) {
	j.motor = m
	j.data = 0
	j.sync = false
}

// HasDisk returns true if a usable disk is present, otherwise false.
func (j *Mechanic) HasDisk() bool {
	return j.disk.Usable()
}

// WriteProtectionState returns the write protection state of the current disk as a uint8 value.
func (j *Mechanic) WriteProtectionState() uint8 {
	const wp = 0x10
	if !j.diskChanged {
		if !j.writeProtected {
			return wp
		}
		return 0
	}
	j.diskChanged = false
	if j.writeProtected {
		return wp
	}
	return 0
}

// updateHeadPos adjusts the disk head position by setting the half-track and recalculating rotation intervals.
func (j *Mechanic) updateHeadPos() {
	j.disk.SetHeadHalfTrack(j.headPos)
	j.updateRotationIntervals()
}

// updateRotationIntervals recalculates and updates the rotation intervals based on the disk's microseconds per byte rate.
func (j *Mechanic) updateRotationIntervals() {
	j.rotationIntervals = int(j.disk.MicroSecPerByte())
	//log.Printf("halfTrack: %d, track: %d => rotation intervals: %d", j.headPos, j.headPos>>1, j.rotationIntervals)
}

// MoveHeadOut decreases the `headPos` of the Mechanic by one step, ensuring it doesn't go below the minimum allowable position.
func (j *Mechanic) MoveHeadOut() {
	if j.headPos <= 2 {
		return
	}
	j.headPos--
	j.updateHeadPos()
}

// MoveHeadIn increments the head position of the Mechanic unless it has reached the maximum limit (headHalfStep).
func (j *Mechanic) MoveHeadIn() {
	if j.headPos >= headHalfStep {
		return
	}
	j.headPos++
	j.updateHeadPos()
}

// insertDisk loads a disk image from the specified file path and initializes it for use, setting write protection if necessary.
// Returns an error if the file cannot be opened, read, or the disk cannot be created.
func (j *Mechanic) insertDisk(filePath string) error {
	fd, err := os.OpenFile(filePath, os.O_RDWR, 0)
	if err != nil {
		if fd, err = os.OpenFile(filePath, os.O_RDONLY, 0); err != nil {
			return err
		}
		j.writeProtected = true
	}
	defer fd.Close()
	image, err := io.ReadAll(fd)
	if err != nil {
		return err
	}
	g, err := j.factory.Create(image)
	if err != nil {
		return err
	}
	j.disk = g
	j.updateRotationIntervals()
	return nil
}

//func (j *Mechanic) Load(fp string) error {
//	if !j.HasDisk() {
//		return j.init(fp)
//	} else if j.filePath != fp {
//		if err := j.init(fp); err != nil {
//			return err
//		}
//		j.diskChanged = true
//		return nil
//	}
//	return nil
//}
