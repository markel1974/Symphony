package inputs

import (
	"github.com/markel1974/c64emu/src/fifo"
)

type Joystick struct {
	storage *fifo.StaticFifo
	joy     int
	s1      uint
	s2      uint
}

func NewJoystick() *Joystick {
	j := &Joystick{
		storage: nil,
		joy:     0xff,
		s1:      0,
		s2:      0,
	}
	j.Update(0x0000, 0xffff, 40)
	return j
}

func (k *Joystick) Update(min uint16, max uint16, sensitivity uint16) {
	interval := max - min
	k.s1 = uint(min + ((sensitivity * interval) / 100))
	k.s2 = uint(min + (((100 - sensitivity) * interval) / 100))
}

func (k *Joystick) Reset() {
	k.storage = fifo.NewStaticFifo(256)
	k.joy = 0xff
}

func (k *Joystick) Move(x uint, y uint, buttons uint) {
	k.joy = 0xff
	if x < k.s1 {
		k.joy &= 0xfb // Left
	} else if x > k.s2 {
		k.joy &= 0xf7 // Right
	}
	if y < k.s1 {
		k.joy &= 0xfe // Up
	} else if y > k.s2 {
		k.joy &= 0xfd // Down
	}
	if (buttons & 1) != 0 {
		k.joy &= 0xef // Button
	}
	if (buttons & 2) != 0 {
		//TODO SID POTX / POTY
	}
	k.storage.Set(k.joy)
}

func (k *Joystick) SetKey(pressed bool, jId int) {
	if pressed {
		k.joy = joyKeyDown(k.joy, jId)
		k.storage.Set(k.joy)
	} else {
		k.joy = joyKeyUp(k.joy, jId)
		k.storage.Set(k.joy)
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
	case KeyJFire:
		j |= 0x10
		return j
	case KeyJUp:
		j |= 0x01
		return j
	case KeyJDown:
		j |= 0x02
		return j
	case KeyJLeft:
		j |= 0x04
		return j
	case KeyJRight:
		j |= 0x08
		return j
	case KeyJUpLeft:
		j |= 0x05
		return j
	case KeyJUpRight:
		j |= 0x09
		return j
	case KeyJDownLeft:
		j |= 0x06
		return j
	case KeyJDownRight:
		j |= 0x0a
		return j
	case KeyJCenter:
		return 0xff
	}
	return 0xff
}

func joyKeyDown(j int, kc int) int {
	switch kc {
	case KeyJFire:
		j &= ^0x10
		return j
	case KeyJUp:
		j |= 0x02
		j &= ^0x01
		return j
	case KeyJDown:
		j |= 0x01
		j &= ^0x02
		return j
	case KeyJLeft:
		j |= 0x08
		j &= ^0x04
		return j
	case KeyJRight:
		j |= 0x04
		j &= ^0x08
		return j
	case KeyJUpLeft:
		j |= 0x0a
		j &= ^0x05
		return j
	case KeyJUpRight:
		j |= 0x06
		j &= ^0x09
		return j
	case KeyJDownLeft:
		j |= 0x09
		j &= ^0x06
		return j
	case KeyJDownRight:
		j |= 0x05
		j &= ^0x0a
		return j
	case KeyJCenter:
		j |= 0x0f
		return j
	}
	return 0xff
}
