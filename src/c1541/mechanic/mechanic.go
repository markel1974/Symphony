package mechanic

import (
	"github.com/markel1974/c64emu/src/c1541/gcr"
	"io"
	"os"
)

//see https://sta.c64.org/cbm1541mem.html

type Mechanic struct {
	gcr            *gcr.GCR
	writeProtected bool
	diskChanged    bool
	filePath       string
	motor          bool
	factory        *gcr.Factory
}

func NewMechanic() *Mechanic {
	j := &Mechanic{
		writeProtected: false,
		diskChanged:    false,
		motor:          false,
		factory:        gcr.NewFactory(),
		gcr:            nil,
	}
	return j
}

func (j *Mechanic) Reset() {
	j.gcr = nil
	j.writeProtected = false
	j.diskChanged = false
}

func (j *Mechanic) Setup(filePath string) {
	if !j.HasDisk() {
		j.filePath = filePath
		_ = j.readFile(j.filePath)
	} else if j.filePath != filePath {
		j.filePath = filePath
		j.Reset()
		_ = j.readFile(j.filePath)
		j.diskChanged = true
	}
}

func (j *Mechanic) RotateDisk() {
	if j.gcr == nil {
		return
	}
	j.gcr.Rotate()
}

func (j *Mechanic) SetMotor(m bool) {
	j.motor = m
}

func (j *Mechanic) HasDisk() bool {
	return j.gcr != nil
}

func (j *Mechanic) WriteProtectionState() uint8 {
	r := uint8(0)
	if j.diskChanged {
		j.diskChanged = false
		if j.writeProtected {
			r = 0x10
		}
	} else {
		if !j.writeProtected {
			r = 0x10
		}
	}
	return r
}

func (j *Mechanic) SyncFound() bool {
	if j.gcr == nil {
		return false
	}
	if j.gcr.Read() == 0xff {
		return true
	}
	return false
}

func (j *Mechanic) ReadByte() uint8 {
	if j.gcr == nil {
		return 0
	}
	return j.gcr.Read()
}

func (j *Mechanic) WriteByte(data uint8) {
	if j.gcr == nil {
		return
	}
	j.gcr.Write(data)
	//track := j.gcrTrackStart
	//sector := _numSectors[j.currentHalfTrack>>1]
	//offset := j.offsetFromTrackSector(track, int(sector))
	//fmt.Println("WRITING ", j.gcrIdx, data)
	//fmt.Println(j.currentHalfTrack, sector, offset)
	//fmt.Println("------------------")
}

func (j *Mechanic) MoveHeadOut() {
	if j.gcr == nil {
		return
	}
	j.gcr.MoveOut()
}

func (j *Mechanic) MoveHeadIn() {
	if j.gcr == nil {
		return
	}
	j.gcr.MoveIn()
}

func (j *Mechanic) readFile(filePath string) error {
	j.Reset()
	fd, err := os.OpenFile(filePath, os.O_RDWR, 0)
	if err != nil {
		j.writeProtected = true
		if fd, err = os.OpenFile(filePath, os.O_RDONLY, 0); err != nil {
			return err
		}
	}
	image, err := io.ReadAll(fd)
	if err != nil {
		return err
	}
	_ = fd.Close()
	g, err := j.factory.Create(image)
	if err != nil {
		return err
	}
	j.gcr = g
	return nil
}

//var _counter = 0
//var _start = int64(0)

//func (j *Mechanics) Emulate() {
//The measurement of the speed of rotation of the Commodore 1541 floppy disk drive is very accurate,
//as long as it is between 265 and 325 revolutions per minute. Outside these ranges,
//the program may return unreliable values, but this does not matter, since the speed must be set at 300 rpm.

//a rotation every 200ms (300rpm)
/*
	if j.motor {
		_counter++
		if _counter > 250000 {
			now := time.Now().UnixMilli()
			fmt.Println(now - _start)
			_start = now
			_counter = 0
			j.rotateDisk()
		}
	}
*/
//}
