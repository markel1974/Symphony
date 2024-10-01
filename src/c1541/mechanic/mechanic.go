package mechanic

import (
	"io"
	"os"
)

//see https://sta.c64.org/cbm1541mem.html

type Mechanic struct {
	gcr            IDisk
	writeProtected bool
	diskChanged    bool
	filePath       string
	motor          bool
	empty          IDisk
	factory        *Factory
}

func NewMechanic() *Mechanic {
	factory := NewFactory()
	empty, _ := factory.Create(nil)
	j := &Mechanic{
		gcr:            empty,
		writeProtected: false,
		diskChanged:    false,
		filePath:       "",
		motor:          false,
		empty:          empty,
		factory:        factory,
	}
	return j
}

func (j *Mechanic) Reset() {
	j.gcr = j.empty
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

func (j *Mechanic) SetMotor(m bool) {
	j.motor = m
}

func (j *Mechanic) HasDisk() bool {
	return j.gcr.Usable()
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
	if j.gcr.Read() == 0xff {
		return true
	}
	return false
}

func (j *Mechanic) RotateDisk() {
	j.gcr.Rotate()
}

func (j *Mechanic) ReadByte() uint8 {
	return j.gcr.Read()
}

func (j *Mechanic) WriteByte(data uint8) {
	j.gcr.Write(data)
}

func (j *Mechanic) MoveHeadOut() {
	j.gcr.MoveOut()
}

func (j *Mechanic) MoveHeadIn() {
	j.gcr.MoveIn()
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
	j.gcr = g
	return nil
}
