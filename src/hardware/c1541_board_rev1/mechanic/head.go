package mechanic

import "github.com/markel1974/c64emu/src/hardware/c1541_board_rev1/disk"

// directionNone represents no direction.
const (
	directionNone    = 0
	directionInward  = 1  // Verso tracce più alte
	directionOutward = -1 // Verso tracce più basse
)

const (
	vibrationDecayRate = 0.005
	vibrationFactorMax = 1.40 //40 % Max
	vibrationStep      = 0.20 //20%
)

// notReady represents an uninitialized or inactive state for dataRead and dataWrite in the Head struct.
const (
	notReady = 1
)

// Head represents the read/write head of a disk drive, managing position, state, and data operations.
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
	vibrationFactor  float64
}

// NewHead initializes and returns a new instance of Head with the provided default position.
// It sets up the Head struct with default values for its fields.
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
		vibrationFactor:  1.0,
	}
	return h
}

// Reset reinitializes the head's state, resetting position, counters, direction, data states, and vibration factor.
func (h *Head) Reset() {
	h.currentPos = h.defaultPos
	h.consecutiveSteps = 0
	h.seekTime = 0
	h.direction = directionNone
	h.writing = false
	h.syncCounter = 0
	h.dataWrite = notReady
	h.dataRead = notReady
	h.vibrationFactor = 1.0
}

// SetWritingMode sets the writing mode of the Head to the specified boolean value.
func (h *Head) SetWritingMode(w bool) {
	h.writing = w
}

// DecayVibration reduces the vibrationFactor of the head toward a stable value, stopping at a minimum threshold of 1.0.
func (h *Head) DecayVibration() {
	// Valore di decadimento per tick di emulazione. Può essere regolato.

	if h.vibrationFactor > 1.0 {
		h.vibrationFactor -= vibrationDecayRate
		if h.vibrationFactor < 1.0 {
			h.vibrationFactor = 1.0
		}
	}
}

// WriteByte sets the byte value to be written by assigning it to the `dataWrite` field of the Head instance.
func (h *Head) WriteByte(d uint8) {
	h.dataWrite = int(d)
}

// ReadByte reads and returns the next byte from the head if data is available; returns 0 if no data is ready.
func (h *Head) ReadByte() uint8 {
	if h.dataRead == notReady {
		return 0
	}
	v := uint8(h.dataRead)
	//fmt.Printf("ReadByte %d From Track %d\n", v, j.currentPos)
	h.dataRead = notReady
	return v
}

// ReadWrite performs read or write operations with the disk based on the current writing mode of the head.
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

// ByteReady determines if the head is ready to process the next byte in either read or write mode, depending on its state.
func (h *Head) ByteReady() bool {
	if h.writing {
		return h.dataWrite == notReady
	}
	v := h.dataRead != notReady
	return v
}

// HasSync checks whether the sync counter has reached or exceeded the required synchronization tolerance.
func (h *Head) HasSync() bool {
	return h.syncCounter >= syncTolerance
}

// Move adjusts the head's position on the disk to the desired track, considering delays, direction, and vibration effects.
// Returns a boolean that indicates whether the head is in motion or has reached the desired position.
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
	var seekTime float64
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
			if h.vibrationFactor += vibrationStep; h.vibrationFactor > vibrationFactorMax {
				h.vibrationFactor = vibrationFactorMax
			}
			seekTime += float64(headBacklashDelay + polarity)
		}
	}

	seekTime *= h.vibrationFactor

	if h.consecutiveSteps++; h.consecutiveSteps > 1 {
		seekTime += float64(h.consecutiveSteps-1) * headExtraSettlingPerStep
	}
	if seekTime > headMaxDelay {
		seekTime = headMaxDelay
	}

	//fmt.Printf("ASYNC MOVE HEAD OLD %d NEW %d REQ %d: %d\n", h.currentPos, newPos, headPosRequired, disk.MicroSecPerByte())
	if disk.SetHeadHalfTrack(newPos) {
		h.currentPos = newPos
		h.seekTime = int(seekTime)
		h.direction = direction
		h.dataRead = notReady
		h.syncCounter = 0
	}
	return true
}
