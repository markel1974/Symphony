package mechanic

import (
	"io"
	"os"
)

//const numTracks = 35
//const numHalfTracks = numTracks * 2

//see https://sta.c64.org/cbm1541mem.html

type Mechanic struct {
	disk           IDisk
	writeProtected bool
	diskChanged    bool
	filePath       string
	motor          bool
	empty          IDisk
	factory        *Factory
	headPos        uint8
}

func NewMechanic() *Mechanic {
	factory := NewFactory()
	empty, _ := factory.Create(nil)
	j := &Mechanic{
		disk:           empty,
		writeProtected: false,
		diskChanged:    false,
		filePath:       "",
		motor:          false,
		empty:          empty,
		factory:        factory,
		headPos:        2,
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

func (j *Mechanic) SyncFound() bool {
	if j.disk.Read() == 0xff {
		return true
	}
	return false
}

func (j *Mechanic) RotateDisk() {
	j.disk.Rotate()
}

func (j *Mechanic) ReadByte() uint8 {
	return j.disk.Read()
}

func (j *Mechanic) WriteByte(data uint8) {
	j.disk.Write(data)
}

func (j *Mechanic) MoveHeadOut() {
	//todo halfTrack handler
	if j.headPos <= 2 {
		return
	}
	j.headPos--
	track := j.headPos >> 1
	j.disk.SetHeadTrack(track)
}

func (j *Mechanic) MoveHeadIn() {
	//todo halfTrack handler
	halfTrack := j.disk.GetTracksNumber() * 2
	if j.headPos >= halfTrack {
		return
	}
	j.headPos++
	track := j.headPos >> 1
	j.disk.SetHeadTrack(track)
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
