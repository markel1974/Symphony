package keyboard

type Joystick struct {
}

func NewJoystick() *Joystick {
	return &Joystick{}
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
