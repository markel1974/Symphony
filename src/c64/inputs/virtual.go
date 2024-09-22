package inputs

type Virtual struct {
	numLock bool
	capital bool
	//ext     bool
	//menu    bool
}

func NewVirtual() *Virtual {
	return &Virtual{}
}

func (k *Virtual) Reset() {
	k.numLock = false
	k.capital = false
	//k.ext = false
	//k.menu = false
}

//func (k *Virtual) SetExt() {
//	k.ext = !k.ext
//}

//func (k *Virtual) SetMenu() {
//	k.menu = !k.menu
//}

func (k *Virtual) NumLockToggle() {
	k.numLock = !k.numLock
}

func (k *Virtual) CapitalToggle() {
	k.capital = !k.capital
}

func (k *Virtual) FromVirtual(vKey int) int {
	var result = -1
	switch vKey {
	case VKReturn:
		//if k.menu {
		//	return KEY_ALTENTER
		//}
		return matrix(0, 1)
	case VKBack:
		return matrix(0, 0)
	case VKSpace:
		return matrix(7, 4)
	case VKEscape:
		//RUN/STOP
		return matrix(7, 7)
	case VKTab:
		return -1
	case VKDelete:
		return matrix(0, 0)
		//return KEY_KPPERIOD
	case VKShift:
		return matrix(1, 7)
	case VKControl:
		//if k.ext {
		//	return matrix(7, 5)
		//}
		return matrix(7, 2)
	case VKMenu:
		matrix(7, 5)
	case VKInsert:
		return matrix(0, 0) | 0x80
	case VKHome:
		return matrix(6, 3)
	case VKEnd:
		return matrix(6, 0)
	case VKPrior:
		return matrix(6, 6)
	case VKNext:
		return matrix(6, 5)
	case VKUp:
		return matrix(0, 7) | 0x80
	case VKDown:
		return matrix(0, 7)
	case VKLeft:
		return matrix(0, 2) | 0x80
	case VKRight:
		return matrix(0, 2)
	case VKF1:
		return matrix(0, 4)
	case VKF2:
		return matrix(0, 4) | 0x80
	case VKF3:
		return matrix(0, 5)
	case VKF4:
		return matrix(0, 5) | 0x80
	case VKF5:
		return matrix(0, 6)
	case VKF6:
		return matrix(0, 6) | 0x80
	case VKF7:
		return matrix(0, 3)
	case VKF8:
		return matrix(0, 3) | 0x80
	case VKBracketLeft:
		return matrix(5, 6)
	case VKAsterisk:
		return matrix(6, 1)
	case VKSlash:
		return matrix(6, 7)
	case VKSemicolon:
		return matrix(5, 5)
	case VKGrave:
		return matrix(7, 1)
	case VKMinus:
		return matrix(5, 3)
	case VKPlus:
		return matrix(5, 0)
	case VKEqual:
		return matrix(6, 5)
	case VKComma:
		return matrix(5, 7)
	case VKPeriod:
		return matrix(5, 4)
	case VKQuote:
		return matrix(6, 2)
	case VKBackslash:
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

/*
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

	case VK_numpad0:
		return KEY_FIRE
	case VK_numpad1:
		return KeyJDownLeft
	case VK_numpad2:
		return KeyJDown
	case VK_numpad3:
		return KeyJDownRight
	case VK_numpad4:
		return KeyJLeft
	case VK_numpad5:
		return KEY_CENTER
	case VK_numpad6:
		return KeyJRight
	case VK_numpad7:
		return KeyJUpLeft
	case VK_numpad8:
		return KeyJUp
	case VK_numpad9:
		return KeyJUpRight
*/
