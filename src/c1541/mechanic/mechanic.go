package mechanic

import (
	"io"
	"os"
)

//see https://sta.c64.org/cbm1541mem.html

//0.985mhz / RPM 300 => (300/60) * track.Len() => es 5 * 7434 => 37170

const headStep = 35
const headHalfStep = headStep * 2 //half track

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
}

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
	}
	return j
}

func (j *Mechanic) Reset() {
	j.disk = j.empty
	j.writeProtected = false
	j.diskChanged = false
	j.filePath = ""
	j.motor = false
}

func (j *Mechanic) init(fp string) error {
	j.Reset()
	if err := j.insertDisk(fp); err != nil {
		return err
	}
	j.filePath = fp
	return nil
}

func (j *Mechanic) Setup(fp string) {
	if err := j.init(fp); err != nil {
		return
	}
}

func (j *Mechanic) RotateDisk() {
	j.disk.Rotate()
}

/*
var _counter uint64 = 0
var _readCounter uint64 = 0
func (j *Mechanic) Emulate() {
	//if j.motor {
	//	_counter++
	//} else {
	//	_counter = 0
	//}
	//return
	if j.motor {
		j.rotationCounter++
		if j.rotationCounter >= j.rotationIntervals {
			j.disk.Rotate()
			j.rotationCounter = 0
		}
	}
}
*/

func (j *Mechanic) ReadByte() uint8 {
	//fmt.Println(_counter - _readCounter)
	//_counter = _readCounter
	return j.disk.Read()
}

func (j *Mechanic) WriteByte(data uint8) {
	j.disk.Write(data)
}

func (j *Mechanic) SyncFound() bool {
	if !j.motor {
		return true
	}
	if j.disk.Read() == 0xff {
		return true
	}
	return false
}

func (j *Mechanic) SetMotor(m bool) {
	j.motor = m
}

func (j *Mechanic) HasDisk() bool {
	return j.disk.Usable()
}

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

func (j *Mechanic) updateHeadPos() {
	track := j.headPos >> 1
	j.disk.SetHeadTrack(track)
	j.updateRotationIntervals()

}

func (j *Mechanic) updateRotationIntervals() {
	j.rotationIntervals = int(j.disk.MicroSecPerByte())
	//log.Printf("track: %d => rotation intervals: %d", j.headPos>>1, j.rotationIntervals)
}

func (j *Mechanic) MoveHeadOut() {
	//todo halfTrack handler
	if j.headPos <= 2 {
		return
	}
	j.headPos--
	j.updateHeadPos()
}

func (j *Mechanic) MoveHeadIn() {
	//todo halfTrack handler
	if j.headPos >= headHalfStep {
		return
	}
	j.headPos++
	j.updateHeadPos()
}

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
