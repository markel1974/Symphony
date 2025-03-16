package keyboard_c64

import (
	"github.com/markel1974/c64emu/src/component"
)

// Virtual represents a virtual keyboard with states for num-lock and capital-lock toggles.
type Virtual struct {
	numLock bool
	capital bool
}

// NewVirtual creates and returns a new instance of Virtual with default states for numLock and capital unset.
func NewVirtual() *Virtual {
	return &Virtual{}
}

// Reset resets the Virtual keyboard state by disabling Num Lock and Caps Lock features.
func (k *Virtual) Reset() {
	k.numLock = false
	k.capital = false
}

// NumLockToggle toggles the state of the Num Lock key by inverting its current value on the virtual keyboard.
func (k *Virtual) NumLockToggle() {
	k.numLock = !k.numLock
}

// CapitalToggle toggles the state of the capital (caps lock) key in the Virtual struct by inverting its current value.
func (k *Virtual) CapitalToggle() {
	k.capital = !k.capital
}

// FromVirtual maps a virtual key code to its corresponding matrix value based on predefined key mappings and shifts.
func (k *Virtual) FromVirtual(vKey int) int {
	var result = -1
	switch vKey {
	case component.VKReturn:
		return matrix(0, 1)
	case component.VKBack:
		return matrix(7, 1) //matrix(0, 0)
	case component.VKSpace:
		return matrix(7, 4)
	case component.VKEscape:
		//RUN/STOP
		return matrix(7, 7)
	case component.VKTab:
		return -1
	case component.VKDelete:
		return matrix(0, 0)
	case component.VKShift:
		return matrix(1, 7)
	case component.VKControl:
		return matrix(7, 2)
	case component.VKInsert:
		return matrix(0, 0) | 0x80
	case component.VKHome:
		return matrix(6, 3)
	case component.VKEnd:
		return matrix(6, 0)
	case component.VKPrior:
		return matrix(6, 6)
	case component.VKNext:
		return matrix(6, 5)
	case component.VKUp:
		return matrix(0, 7) | 0x80
	case component.VKDown:
		return matrix(0, 7)
	case component.VKLeft:
		return matrix(0, 2) | 0x80
	case component.VKRight:
		return matrix(0, 2)
	case component.VKF1:
		return matrix(0, 4)
	case component.VKF2:
		return matrix(0, 4) | 0x80
	case component.VKF3:
		return matrix(0, 5)
	case component.VKF4:
		return matrix(0, 5) | 0x80
	case component.VKF5:
		return matrix(0, 6)
	case component.VKF6:
		return matrix(0, 6) | 0x80
	case component.VKF7:
		return matrix(0, 3)
	case component.VKF8:
		return matrix(0, 3) | 0x80
	case component.VKBracketLeft:
		return matrix(5, 6)
	case component.VKAsterisk:
		return matrix(6, 1)
	case component.VKSlash:
		return matrix(6, 7)
	case component.VKSemicolon:
		return matrix(5, 5)
	case component.VKGrave:
		return matrix(7, 1)
	case component.VKMinus:
		return matrix(5, 3)
	case component.VKPlus:
		return matrix(5, 0)
	case component.VKEqual:
		return matrix(6, 5)
	case component.VKComma:
		return matrix(5, 7)
	case component.VKPeriod:
		return matrix(5, 4)
	case component.VKQuote:
		return matrix(6, 2)
	case component.VKBackslash:
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
