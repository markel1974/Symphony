package keyboard

type Virtual struct {
	numLock bool
	capital bool
	ext     bool
	menu    bool
}

func NewVirtual() *Virtual {
	return &Virtual{}
}

func (k *Virtual) Reset() {
	k.numLock = false
	k.capital = false
	k.ext = false
	k.menu = false
}

func (k *Virtual) SetExt() {
	k.ext = !k.ext
}

func (k *Virtual) SetNumLock() {
	k.numLock = !k.numLock
}

func (k *Virtual) SetCapital() {
	k.capital = !k.capital
}

func (k *Virtual) SetMenu() {
	k.menu = !k.menu
}

func (k *Virtual) FromVirtual(vKey int) int {
	var result = -1
	switch vKey {
	case VK_return:
		if k.menu {
			return KEY_ALTENTER
		}
		//if k.control {
		//	return KEY_CTRLENTER
		//}
		if k.ext {
			return matrix(0, 1)
		}
		return matrix(0, 1)
	case VK_back:
		return matrix(0, 0)
	case VK_space:
		return matrix(7, 4)
	case VK_escape:
		//RUN/STOP
		return matrix(7, 7)
	case VK_tab:
		return -1
	case VK_delete:
		if k.ext {
			return matrix(0, 0)
		}
		return KEY_KPPERIOD
	case VK_shift:
		//TODO
		//if k.shift {
		//	return matrix(6, 4)
		//}
		return matrix(1, 7)
	case VK_control:
		if k.ext {
			return matrix(7, 5)
		}
		return matrix(7, 2)
	case VK_menu:
		if k.ext {
			return matrix(7, 5)
		}
		matrix(7, 5)
	case VK_insert:
		if k.ext {
			return matrix(0, 0) | 0x80
		}
		return KEY_FIRE
	case VK_home:
		if k.ext {
			return matrix(6, 3)
		}
		return KEY_JUPLF
	case VK_end:
		if k.ext {
			return matrix(6, 0)
		}
		return KEY_JDNLF
	case VK_prior:
		if k.ext {
			return matrix(6, 6)
		}
		return KEY_JUPRT
	case VK_next:
		if k.ext {
			return matrix(6, 5)
		}
		return KEY_JDNRT
	case VK_up:
		if k.ext {
			return matrix(0, 7) | 0x80
		}
		return KEY_JUP
	case VK_down:
		if k.ext {
			return matrix(0, 7)
		}
		return KEY_JDN
	case VK_left:
		if k.ext {
			return matrix(0, 2) | 0x80
		}
		return KEY_JLF
	case VK_right:
		if k.ext {
			return matrix(0, 2)
		}
		return KEY_JRT
	case VK_joy_fire:
		return KEY_FIRE
	case VK_joy_down_left:
		return KEY_JDNLF
	case VK_joy_down:
		return KEY_JDN
	case VK_joy_down_right:
		return KEY_JDNRT
	case VK_joy_left:
		return KEY_JLF
	case VK_joy_center:
		return KEY_CENTER
	case VK_joy_right:
		return KEY_JRT
	case VK_joy_up_left:
		return KEY_JUPLF
	case VK_joy_up:
		return KEY_JUP
	case VK_joy_up_right:
		return KEY_JUPRT
	case VK_numpad0:
		return KEY_FIRE
	case VK_numpad1:
		return KEY_JDNLF
	case VK_numpad2:
		return KEY_JDN
	case VK_numpad3:
		return KEY_JDNRT
	case VK_numpad4:
		return KEY_JLF
	case VK_numpad5:
		return KEY_CENTER
	case VK_numpad6:
		return KEY_JRT
	case VK_numpad7:
		return KEY_JUPLF
	case VK_numpad8:
		return KEY_JUP
	case VK_numpad9:
		return KEY_JUPRT
	case VK_f1:
		return matrix(0, 4)
	case VK_f2:
		return matrix(0, 4) | 0x80
	case VK_f3:
		return matrix(0, 5)
	case VK_f4:
		return matrix(0, 5) | 0x80
	case VK_f5:
		return matrix(0, 6)
	case VK_f6:
		return matrix(0, 6) | 0x80
	case VK_f7:
		return matrix(0, 3)
	case VK_f8:
		return matrix(0, 3) | 0x80
	case VK_f9:
		return KEY_F9
	case VK_f10:
		return KEY_F10
	case VK_f11:
		return KEY_F11
	case VK_f12:
		return KEY_F12
	case VK_clear:
		return KEY_CENTER
	case VK_pause:
		return KEY_PAUSE
	//case VK_numlock: return KEY_NUMLOCK;
	case VK_numlock:
		return -1
	case VK_capital:
		return -1
	case VK_multiply:
		return KEY_KPMULT
	case VK_divide:
		return KEY_KPDIV
	case VK_subtract:
		return KEY_KPMINUS
	case VK_add:
		return KEY_KPPLUS
	case VK_decimal:
		return KEY_KPPERIOD
	case VK_bracketleft:
		return matrix(5, 6)
	case VK_bracketright:
		return matrix(6, 1)
	case VK_slash:
		return matrix(6, 7)
	case VK_semicolon:
		return matrix(5, 5)
	case VK_grave:
		return matrix(7, 1)
	case VK_minus:
		return matrix(5, 3)
	case VK_plus:
		return matrix(5, 0)
	case VK_equal:
		return matrix(6, 5)
	case VK_comma:
		return matrix(5, 7)
	case VK_period:
		return matrix(5, 4)
	case VK_quote:
		return matrix(6, 2)
	case VK_backslash:
		return matrix(6, 6)
	case '0':
		return matrix(4, 3)
	case '1':
		return matrix(7, 0)
	case '2':
		return matrix(7, 3)
	case '3':
		return matrix(1, 0)
	case '4':
		return matrix(1, 3)
	case '5':
		return matrix(2, 0)
	case '6':
		return matrix(2, 3)
	case '7':
		return matrix(3, 0)
	case '8':
		return matrix(3, 3)
	case '9':
		return matrix(4, 0)
	case 'A':
		result = matrix(1, 2)
	case 'B':
		result = matrix(3, 4)
	case 'C':
		result = matrix(2, 4)
	case 'D':
		result = matrix(2, 2)
	case 'E':
		result = matrix(1, 6)
	case 'F':
		result = matrix(2, 5)
	case 'G':
		result = matrix(3, 2)
	case 'H':
		result = matrix(3, 5)
	case 'I':
		result = matrix(4, 1)
	case 'J':
		result = matrix(4, 2)
	case 'K':
		result = matrix(4, 5)
	case 'L':
		result = matrix(5, 2)
	case 'M':
		result = matrix(4, 4)
	case 'N':
		result = matrix(4, 7)
	case 'O':
		result = matrix(4, 6)
	case 'P':
		result = matrix(5, 1)
	case 'Q':
		result = matrix(7, 6)
	case 'R':
		result = matrix(2, 1)
	case 'S':
		result = matrix(1, 5)
	case 'T':
		result = matrix(2, 6)
	case 'U':
		result = matrix(3, 6)
	case 'V':
		result = matrix(3, 7)
	case 'W':
		result = matrix(1, 1)
	case 'X':
		result = matrix(2, 7)
	case 'Y':
		result = matrix(3, 1)
	case 'Z':
		result = matrix(1, 4)
	}
	if result != -1 && k.capital {
		result |= 0x80
	}
	return result
}
