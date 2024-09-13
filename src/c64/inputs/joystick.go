package inputs

import (
	"container/list"
	"github.com/markel1974/c64emu/src/components/quartz"
)

type joyData struct {
	joy         uint8
	persistence uint8
}

type Joystick struct {
	quartz      *quartz.Quartz
	dataStorage *list.List
	joy         int
}

func NewJoystick(quartz *quartz.Quartz) *Joystick {
	return &Joystick{
		quartz: quartz,
	}
}

func (k *Joystick) Reset() {
	k.dataStorage = list.New()
	k.joy = 0xff
}

func (k *Joystick) Build(x uint, y uint, buttons uint) int {
	s1 := JOYSTICK_SENSITIVITY
	s2 := 100 - JOYSTICK_SENSITIVITY
	s1Val := uint(JOYSTICK_MIN + s1*JOYSTICK_RANGE/100)
	s2Val := uint(JOYSTICK_MIN + s2*JOYSTICK_RANGE/100)
	joystick := 0xff
	if x < s1Val {
		joystick &= 0xfb // Left
	} else if x > s2Val {
		joystick &= 0xf7 // Right
	}
	if y < s1Val {
		joystick &= 0xfe // Up
	} else if y > s2Val {
		joystick &= 0xfd // Down
	}
	if (buttons & 1) != 0 {
		joystick &= 0xef // Button
	}
	if (buttons & 2) != 0 {
		//TODO
	}
	return joystick
}

func (k *Joystick) Set(pressed bool, jId int) {
	const persistence = 1
	if k.dataStorage.Len() >= MAX_STORAGE_SIZE {
		return
	}
	if pressed {
		if joy, ok := joyDown(k.joy, jId); ok {
			k.joy = joy
			k.dataStorage.PushBack(&joyData{joy: uint8(joy), persistence: persistence})
		}
	} else {
		if joy, ok := joyUp(k.joy, jId); ok {
			k.joy = joy
			k.dataStorage.PushBack(&joyData{joy: uint8(joy), persistence: persistence})
		}
	}
}

func (k *Joystick) Poll() (uint8, bool) {
	if k.dataStorage.Len() == 0 {
		return 0, false
	}
	e := k.dataStorage.Front()
	i := e.Value.(*joyData)
	if i.persistence--; i.persistence == 0 {
		k.dataStorage.Remove(e)
	}
	return i.joy, true
}

func joyUp(j int, kc int) (int, bool) {
	switch kc {
	case KEY_FIRE:
		j |= 0x10
		return j, true
	case KEY_JUP:
		j |= 0x01
		return j, true
	case KEY_JDN:
		j |= 0x02
		return j, true
	case KEY_JLF:
		j |= 0x04
		return j, true
	case KEY_JRT:
		j |= 0x08
		return j, true
	case KEY_JUPLF:
		j |= 0x05
		return j, true
	case KEY_JUPRT:
		j |= 0x09
		return j, true
	case KEY_JDNLF:
		j |= 0x06
		return j, true
	case KEY_JDNRT:
		j |= 0x0a
		return j, true
	}
	return 0xff, false
}

func joyDown(j int, kc int) (int, bool) {
	switch kc {
	case KEY_FIRE:
		j &= ^0x10
		return j, true
	case KEY_JUP:
		j |= 0x02
		j &= ^0x01
		return j, true
	case KEY_JDN:
		j |= 0x01
		j &= ^0x02
		return j, true
	case KEY_JLF:
		j |= 0x08
		j &= ^0x04
		return j, true
	case KEY_JRT:
		j |= 0x04
		j &= ^0x08
		return j, true
	case KEY_JUPLF:
		j |= 0x0a
		j &= ^0x05
		return j, true
	case KEY_JUPRT:
		j |= 0x06
		j &= ^0x09
		return j, true
	case KEY_JDNLF:
		j |= 0x09
		j &= ^0x06
		return j, true
	case KEY_JDNRT:
		j |= 0x05
		j &= ^0x0a
		return j, true
	case KEY_CENTER:
		j |= 0x0f
		return j, true
	}
	return 0xff, false
}
