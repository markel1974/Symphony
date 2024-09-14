package inputs

import (
	"github.com/markel1974/c64emu/src/fifo"
)

type Joystick struct {
	storage *fifo.Queue
	joy     int
}

func NewJoystick() *Joystick {
	return &Joystick{
		storage: nil,
		joy:     0xff,
	}
}

func (k *Joystick) Reset() {
	k.storage = fifo.NewQueue(256)
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
		//TODO SID POTX / POTY
	}
	return joystick
}

func (k *Joystick) SetKey(pressed bool, jId int) {
	if pressed {
		k.joy = joyKeyDown(k.joy, jId)
		k.storage.Add(k.joy)
	} else {
		k.joy = joyKeyUp(k.joy, jId)
		k.storage.Add(k.joy)
	}
}

func (k *Joystick) Poll() (uint8, bool) {
	if k.storage.Len() == 0 {
		return 0, false
	}
	joy, ok := k.storage.Next()
	if !ok {
		return 0, false
	}
	return uint8(joy), true
}

func joyKeyUp(j int, kc int) int {
	switch kc {
	case KEY_FIRE:
		j |= 0x10
		return j
	case KEY_JUP:
		j |= 0x01
		return j
	case KEY_JDN:
		j |= 0x02
		return j
	case KEY_JLF:
		j |= 0x04
		return j
	case KEY_JRT:
		j |= 0x08
		return j
	case KEY_JUPLF:
		j |= 0x05
		return j
	case KEY_JUPRT:
		j |= 0x09
		return j
	case KEY_JDNLF:
		j |= 0x06
		return j
	case KEY_JDNRT:
		j |= 0x0a
		return j
	}
	return 0xff
}

func joyKeyDown(j int, kc int) int {
	switch kc {
	case KEY_FIRE:
		j &= ^0x10
		return j
	case KEY_JUP:
		j |= 0x02
		j &= ^0x01
		return j
	case KEY_JDN:
		j |= 0x01
		j &= ^0x02
		return j
	case KEY_JLF:
		j |= 0x08
		j &= ^0x04
		return j
	case KEY_JRT:
		j |= 0x04
		j &= ^0x08
		return j
	case KEY_JUPLF:
		j |= 0x0a
		j &= ^0x05
		return j
	case KEY_JUPRT:
		j |= 0x06
		j &= ^0x09
		return j
	case KEY_JDNLF:
		j |= 0x09
		j &= ^0x06
		return j
	case KEY_JDNRT:
		j |= 0x05
		j &= ^0x0a
		return j
	case KEY_CENTER:
		j |= 0x0f
		return j
	}
	return 0xff
}
