package keyboard

import (
	"container/list"
	"fmt"
	"unicode"
)

/*
C64 keyboard matrix:

Bit   7   6   5   4   3   2   1   0
0    CUD  F5  F3  F1  F7 CLR RET DEL
1    SHL  E   S   Z   4   A   W   3
2     X   T   F   C   6   D   R   5
3     V   U   H   B   8   G   Y   7
4     N   O   K   M   0   J   I   9
5     ,   @   :   .   -   L   P   +
6     /   ^   =  SHR HOM  ;   *   ?
7    R/S  Q   C= SPC  2  CTL  <-  1
*/

const (
	KEY_F9        = 256
	KEY_F10       = 257
	KEY_F11       = 258
	KEY_F12       = 259
	KEY_FIRE      = 260
	KEY_JUP       = 261
	KEY_JDN       = 262
	KEY_JLF       = 263
	KEY_JRT       = 264
	KEY_JUPLF     = 265
	KEY_JUPRT     = 266
	KEY_JDNLF     = 267
	KEY_JDNRT     = 268
	KEY_CENTER    = 269
	KEY_NUMLOCK   = 270
	KEY_KPPLUS    = 271
	KEY_KPMINUS   = 272
	KEY_KPMULT    = 273
	KEY_KPDIV     = 274
	KEY_KPENTER   = 275
	KEY_KPPERIOD  = 276
	KEY_PAUSE     = 277
	KEY_ALTENTER  = 278
	KEY_CTRLENTER = 279
)

const (
	JOYSTICK_SENSITIVITY = 40     // % of live range
	JOYSTICK_MIN         = 0x0000 // min value of range
	JOYSTICK_MAX         = 0xffff // max value of range
	JOYSTICK_RANGE       = JOYSTICK_MAX - JOYSTICK_MIN
)

const (
	MAX_STORAGE_SIZE = 1024
	CMD_COUNT        = 0
)

const (
	SC_ext     = 0x1
	SC_control = 0x2
	SC_menu    = 0x4
	SC_capital = 0x8
	SC_numlock = 0x10
	SC_shift   = 0x20

	VK_return   = -1
	VK_back     = -2
	VK_space    = -3
	VK_escape   = -4
	VK_tab      = -5
	VK_delete   = -6
	VK_shift    = -7
	VK_control  = -8
	VK_menu     = -9
	VK_insert   = -10
	VK_numlock  = -11
	VK_capital  = -12
	VK_home     = -13
	VK_end      = -14
	VK_prior    = -15
	VK_next     = -16
	VK_up       = -17
	VK_down     = -18
	VK_left     = -19
	VK_right    = -20
	VK_numpad0  = -21
	VK_numpad1  = -22
	VK_numpad2  = -23
	VK_numpad3  = -24
	VK_numpad4  = -25
	VK_numpad5  = -26
	VK_numpad6  = -27
	VK_numpad7  = -28
	VK_numpad8  = -29
	VK_numpad9  = -30
	VK_f1       = -31
	VK_f2       = -32
	VK_f3       = -33
	VK_f4       = -34
	VK_f5       = -35
	VK_f6       = -36
	VK_f7       = -37
	VK_f8       = -38
	VK_f9       = -39
	VK_f10      = -40
	VK_f11      = -41
	VK_f12      = -42
	VK_clear    = -43
	VK_pause    = -44
	VK_multiply = -45
	VK_divide   = -46
	VK_subtract = -47
	VK_add      = -48
	VK_decimal  = -49

	VK_joy_fire       = -50
	VK_joy_down_left  = -51
	VK_joy_down       = -52
	VK_joy_down_right = -53
	VK_joy_left       = -54
	VK_joy_center     = -55
	VK_joy_right      = -56
	VK_joy_up_left    = -57
	VK_joy_up         = -58
	VK_joy_up_right   = -59

	VK_semicolon    = 0xba
	VK_equal        = 0xbb
	VK_comma        = 0xbc
	VK_period       = 0xbe
	VK_minus        = 0xbd
	VK_slash        = 0xbf
	VK_bracketleft  = 0xdb
	VK_backslash    = 0xdc
	VK_quote        = 0xde
	VK_bracketright = 0xdd
	VK_grave        = 0xc0
	VK_plus         = 0xc1
)

func MATRIX(a int, b int) int {
	return ((a) << 3) | (b)
}

type InputKeyData struct {
	c64Byte int
	c64Bit  int
	shifted bool
	joyKey  uint8
	pressed bool
	counter uint8
}

func NewInputKeyData(pressed bool, c64Byte int, c64bit int, shifted bool, joyKey uint8) *InputKeyData {
	return &InputKeyData{
		c64Byte: c64Byte,
		c64Bit:  c64bit,
		shifted: shifted,
		joyKey:  joyKey,
		pressed: pressed,
		counter: CMD_COUNT,
	}
}

type Keyboard struct {
	keyDataStorage *list.List
	inputReady     bool
	srcPre         *list.List
	srcPost        *list.List
	poolCounter    int
	joyKey         int
	numLock        bool
	capital        bool
	ext            bool
	menu           bool
	ready          bool
	joystickSwap   bool
}

func NewKeyboard() *Keyboard {
	k := &Keyboard{}
	k.Reset()
	return k
}

func (k *Keyboard) Reset() {
	k.ready = false
	k.keyDataStorage = list.New()
	k.inputReady = false
	k.srcPre = list.New()
	k.srcPost = list.New()
	k.poolCounter = 0
	k.joyKey = 0xff
	k.numLock = false
	k.capital = false
	k.ext = false
	k.menu = false
	k.joystickSwap = true
}

func (k *Keyboard) SetReady() {
	k.ready = true
	if k.srcPre.Len() > 0 {
		for element := k.srcPre.Front(); element != nil; element = element.Next() {
			k.keyDataStorage.PushBack(element.Value)
		}
		k.srcPre.Init()
	}
}

func (k *Keyboard) SetInputReady(ready bool) {
	if !k.inputReady && ready {
		if k.srcPost.Len() > 0 {
			for element := k.srcPost.Front(); element != nil; element = element.Next() {
				k.keyDataStorage.PushBack(element.Value)
			}
			k.srcPost.Init()
		}
	}
	k.inputReady = ready
}

func (k *Keyboard) NumLock() bool {
	return k.numLock
}

func (k *Keyboard) Capital() bool {
	return k.capital
}

func (k *Keyboard) SetExt(ext bool) {
	k.capital = ext
}

func (k *Keyboard) SetNumLock(numLock bool) {
	k.numLock = numLock
}

func (k *Keyboard) SetCapital(capital bool) {
	k.capital = capital
}

func (k *Keyboard) SetMenu(menu bool) {
	k.menu = menu
}

func (k *Keyboard) SwapJoystick() {
	k.joystickSwap = !k.joystickSwap
}

func (k *Keyboard) HasJoystickSwap() bool {
	return k.joystickSwap
}

func (k *Keyboard) SetVirtualKey(pressed bool, vKey int) {
	if k.keyDataStorage.Len() >= MAX_STORAGE_SIZE {
		return
	}
	kc := k.virtualKey2C64(vKey)
	if kc < 0 {
		return
	}
	if pressed {
		if c64Byte, c64Bit, shifted, joyKey, ok := k.keyDown(kc); ok {
			k.keyDataStorage.PushBack(NewInputKeyData(pressed, c64Byte, c64Bit, shifted, joyKey))
		}
	} else {
		if c64Byte, c64Bit, shifted, joyKey, ok := k.keyUp(kc); ok {
			k.keyDataStorage.PushBack(NewInputKeyData(pressed, c64Byte, c64Bit, shifted, joyKey))
		}
	}
}

func (k *Keyboard) PollKeyboard() (int, int, bool, uint8, bool, bool) {
	if k.poolCounter > 0 {
		k.poolCounter--
	}
	if k.poolCounter == 0 && k.keyDataStorage.Len() > 0 {
		e := k.keyDataStorage.Front()
		i := e.Value.(*InputKeyData)
		c64Byte := i.c64Byte
		c64Bit := i.c64Bit
		shifted := i.shifted
		joyKey := i.joyKey
		pressed := i.pressed
		k.poolCounter = int(i.counter)
		k.keyDataStorage.Remove(e)
		return c64Byte, c64Bit, pressed, joyKey, shifted, true
	}
	return 0, 0, false, 0xff, false, false
}

func (k *Keyboard) SetCommand(srcPre string, srcPost string) {
	var pre *list.List
	var post *list.List
	var ret bool
	if pre, ret = k.prepareCommand(srcPre); !ret {
		return
	}
	if len(srcPost) > 0 {
		if post, ret = k.prepareCommand(srcPost); !ret {
			return
		}
	}
	if pre == nil {
		return
	}
	for element := pre.Front(); element != nil; element = element.Next() {
		if k.ready {
			k.keyDataStorage.PushBack(element.Value)
		} else {
			k.srcPre.PushBack(element.Value)
		}
	}
	if post != nil {
		for element := post.Front(); element != nil; element = element.Next() {
			k.srcPost.PushBack(element.Value)
		}
	}
}

/*
func (k *Keyboard) SetVirtualModifier(modifier int) {
	k.ext = flag.IntToBool(modifier & SC_ext)
	k.numLock = flag.IntToBool(modifier & SC_numlock)
	k.capital = flag.IntToBool(modifier & SC_capital)
	k.menu = flag.IntToBool(modifier & SC_menu)
	k.control = flag.IntToBool(modifier & SC_control)
	k.shift = flag.IntToBool(modifier & SC_shift)
}
*/

func (k *Keyboard) keyUp(kc int) (int, int, bool, uint8, bool) {
	ret := false
	joyKey := 0xff
	var c64Byte int
	var c64Bit int
	var shifted bool

	switch kc {
	case KEY_FIRE:
		k.joyKey |= 0x10
		joyKey = k.joyKey
		ret = true
	case KEY_JUP:
		k.joyKey |= 0x01
		joyKey = k.joyKey
		ret = true
	case KEY_JDN:
		k.joyKey |= 0x02
		joyKey = k.joyKey
		ret = true
	case KEY_JLF:
		k.joyKey |= 0x04
		joyKey = k.joyKey
		ret = true
	case KEY_JRT:
		k.joyKey |= 0x08
		joyKey = k.joyKey
		ret = true
	case KEY_JUPLF:
		k.joyKey |= 0x05
		joyKey = k.joyKey
		ret = true
	case KEY_JUPRT:
		k.joyKey |= 0x09
		joyKey = k.joyKey
		ret = true
	case KEY_JDNLF:
		k.joyKey |= 0x06
		joyKey = k.joyKey
		ret = true
	case KEY_JDNRT:
		k.joyKey |= 0x0a
		joyKey = k.joyKey
		ret = true
	default:
		c64Byte, c64Bit, shifted, ret = k.keyCodeToC64(kc)
	}
	return c64Byte, c64Bit, shifted, uint8(joyKey), ret
}

func (k *Keyboard) keyDown(kc int) (int, int, bool, uint8, bool) {
	ret := false
	joyKey := 0xff
	var c64Byte int
	var c64Bit int
	var shifted bool

	switch kc {
	case KEY_PAUSE:
	case KEY_KPPLUS:
		break
	case KEY_KPMINUS:
		break
	case KEY_KPMULT:
		break
	case KEY_KPDIV:
		break
	case KEY_KPPERIOD:
		break
	case KEY_ALTENTER:
		break
	case KEY_CTRLENTER:
		break
	case KEY_F9:
		break
	case KEY_F10:
		break
	case KEY_F11:
		break
	case KEY_F12:
		break
	case KEY_FIRE:
		k.joyKey &= ^0x10
		joyKey = k.joyKey
		ret = true
		break
	case KEY_JUP:
		k.joyKey |= 0x02
		k.joyKey &= ^0x01
		joyKey = k.joyKey
		ret = true
	case KEY_JDN:
		k.joyKey |= 0x01
		k.joyKey &= ^0x02
		joyKey = k.joyKey
		ret = true
	case KEY_JLF:
		k.joyKey |= 0x08
		k.joyKey &= ^0x04
		joyKey = k.joyKey
		ret = true
	case KEY_JRT:
		k.joyKey |= 0x04
		k.joyKey &= ^0x08
		joyKey = k.joyKey
		ret = true
	case KEY_JUPLF:
		k.joyKey |= 0x0a
		k.joyKey &= ^0x05
		joyKey = k.joyKey
		ret = true
	case KEY_JUPRT:
		k.joyKey |= 0x06
		k.joyKey &= ^0x09
		joyKey = k.joyKey
		ret = true
	case KEY_JDNLF:
		k.joyKey |= 0x09
		k.joyKey &= ^0x06
		joyKey = k.joyKey
		ret = true
	case KEY_JDNRT:
		k.joyKey |= 0x05
		k.joyKey &= ^0x0a
		joyKey = k.joyKey
		ret = true
	case KEY_CENTER:
		k.joyKey |= 0x0f
		ret = true
	default:
		c64Byte, c64Bit, shifted, ret = k.keyCodeToC64(kc)
	}

	return c64Byte, c64Bit, shifted, uint8(joyKey), ret
}

func (k *Keyboard) BuildJoystick(x uint, y uint, buttons uint) int {
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

func (k *Keyboard) keyCodeToC64(kc int) (int, int, bool, bool) {
	if kc >= 0 || kc < 256 {
		c64Byte := (kc >> 3) & 7
		c64Bit := kc & 7
		shifted := (kc & 128) != 0
		return c64Byte, c64Bit, shifted, true
	}
	return 0, 0, false, false
}

func (k *Keyboard) prepareCommand(command string) (*list.List, bool) {
	if len(command) == 0 {
		return nil, false
	}
	cmd := func(keyCode int, pressed bool, storage *list.List) bool {
		if c64Byte, c64Bit, shifted, ok := k.keyCodeToC64(keyCode); ok {
			i1 := NewInputKeyData(pressed, c64Byte, c64Bit, shifted, 0xff)
			storage.PushBack(i1)
			return true
		}
		return false
	}
	storage := list.New()
	for _, it := range command {
		if d1, d2, ok := k.ascii2C64(it); ok {
			if d2 != 0 {
				if ret := cmd(d2, true, storage); !ret {
					return nil, false
				}
			}
			if ret := cmd(d1, true, storage); !ret {
				return nil, false
			}
			if ret := cmd(d1, false, storage); !ret {
				return nil, false
			}
			if d2 != 0 {
				if ret := cmd(d2, false, storage); !ret {
					return nil, false
				}
			}
		}
	}
	return storage, true
}

// TODO MOVE IN PIXEL KEYBOARD!
func (k *Keyboard) virtualKey2C64(vKey int) int {
	var result = -1

	fmt.Println("SHIFT STATE", vKey)

	switch vKey {
	case VK_return:
		if k.menu {
			return KEY_ALTENTER
		}
		//if k.control {
		//	return KEY_CTRLENTER
		//}
		if k.ext {
			return MATRIX(0, 1)
		}
		return MATRIX(0, 1)
	case VK_back:
		return MATRIX(0, 0)
	case VK_space:
		return MATRIX(7, 4)
	case VK_escape:
		return MATRIX(7, 7)
	case VK_tab:
		return -1
	case VK_delete:
		if k.ext {
			return MATRIX(0, 0)
		}
		return KEY_KPPERIOD
	case VK_shift:
		//TODO
		//if k.shift {
		//	return MATRIX(6, 4)
		//}
		return MATRIX(1, 7)
	case VK_control:
		if k.ext {
			return MATRIX(7, 5)
		}
		return MATRIX(7, 2)
	case VK_menu:
		if k.ext {
			return MATRIX(7, 5)
		}
		MATRIX(7, 5)
	case VK_insert:
		if k.ext {
			return MATRIX(0, 0) | 0x80
		}
		return KEY_FIRE
	case VK_home:
		if k.ext {
			return MATRIX(6, 3)
		}
		return KEY_JUPLF
	case VK_end:
		if k.ext {
			return MATRIX(6, 0)
		}
		return KEY_JDNLF
	case VK_prior:
		if k.ext {
			return MATRIX(6, 6)
		}
		return KEY_JUPRT
	case VK_next:
		if k.ext {
			return MATRIX(6, 5)
		}
		return KEY_JDNRT
	case VK_up:
		if k.ext {
			return MATRIX(0, 7) | 0x80
		}
		return KEY_JUP
	case VK_down:
		if k.ext {
			return MATRIX(0, 7)
		}
		return KEY_JDN
	case VK_left:
		if k.ext {
			return MATRIX(0, 2) | 0x80
		}
		return KEY_JLF
	case VK_right:
		if k.ext {
			return MATRIX(0, 2)
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
		return MATRIX(0, 4)
	case VK_f2:
		return MATRIX(0, 4) | 0x80
	case VK_f3:
		return MATRIX(0, 5)
	case VK_f4:
		return MATRIX(0, 5) | 0x80
	case VK_f5:
		return MATRIX(0, 6)
	case VK_f6:
		return MATRIX(0, 6) | 0x80
	case VK_f7:
		return MATRIX(0, 3)
	case VK_f8:
		return MATRIX(0, 3) | 0x80
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
		return MATRIX(5, 6)
	case VK_bracketright:
		return MATRIX(6, 1)
	case VK_slash:
		return MATRIX(6, 7)
	case VK_semicolon:
		return MATRIX(5, 5)
	case VK_grave:
		return MATRIX(7, 1)
	case VK_minus:
		return MATRIX(5, 3)
	case VK_plus:
		return MATRIX(5, 0)
	case VK_equal:
		return MATRIX(6, 5)
	case VK_comma:
		return MATRIX(5, 7)
	case VK_period:
		return MATRIX(5, 4)
	case VK_quote:
		return MATRIX(6, 2)
	case VK_backslash:
		return MATRIX(6, 6)
	case '0':
		return MATRIX(4, 3)
	case '1':
		return MATRIX(7, 0)
	case '2':
		return MATRIX(7, 3)
	case '3':
		return MATRIX(1, 0)
	case '4':
		return MATRIX(1, 3)
	case '5':
		return MATRIX(2, 0)
	case '6':
		return MATRIX(2, 3)
	case '7':
		return MATRIX(3, 0)
	case '8':
		return MATRIX(3, 3)
	case '9':
		return MATRIX(4, 0)
	case 'A':
		result = MATRIX(1, 2)
	case 'B':
		result = MATRIX(3, 4)
	case 'C':
		result = MATRIX(2, 4)
	case 'D':
		result = MATRIX(2, 2)
	case 'E':
		result = MATRIX(1, 6)
	case 'F':
		result = MATRIX(2, 5)
	case 'G':
		result = MATRIX(3, 2)
	case 'H':
		result = MATRIX(3, 5)
	case 'I':
		result = MATRIX(4, 1)
	case 'J':
		result = MATRIX(4, 2)
	case 'K':
		result = MATRIX(4, 5)
	case 'L':
		result = MATRIX(5, 2)
	case 'M':
		result = MATRIX(4, 4)
	case 'N':
		result = MATRIX(4, 7)
	case 'O':
		result = MATRIX(4, 6)
	case 'P':
		result = MATRIX(5, 1)
	case 'Q':
		result = MATRIX(7, 6)
	case 'R':
		result = MATRIX(2, 1)
	case 'S':
		result = MATRIX(1, 5)
	case 'T':
		result = MATRIX(2, 6)
	case 'U':
		result = MATRIX(3, 6)
	case 'V':
		result = MATRIX(3, 7)
	case 'W':
		result = MATRIX(1, 1)
	case 'X':
		result = MATRIX(2, 7)
	case 'Y':
		result = MATRIX(3, 1)
	case 'Z':
		result = MATRIX(1, 4)
	}

	if result != -1 && k.capital {
		result |= 0x80
	}
	return result
}

func (k *Keyboard) ascii2C64(ascii rune) (int, int, bool) {
	r1 := 0
	r2 := 0

	switch unicode.ToUpper(ascii) {
	case '!':
		r2 = MATRIX(6, 4)
		r1 = MATRIX(7, 0)
	case '"':
		r2 = MATRIX(6, 4)
		r1 = MATRIX(7, 3)
	case '#':
		r2 = MATRIX(6, 4)
		r1 = MATRIX(1, 0)
	case '$':
		r2 = MATRIX(6, 4)
		r1 = MATRIX(1, 3)
	case '%':
		r2 = MATRIX(6, 4)
		r1 = MATRIX(2, 0)
	case '&':
		r2 = MATRIX(6, 4)
		r1 = MATRIX(2, 3)
	case '\'':
		r2 = MATRIX(6, 4)
		r1 = MATRIX(3, 0)
	case '(':
		r2 = MATRIX(6, 4)
		r1 = MATRIX(3, 3)
	case ')':
		r2 = MATRIX(6, 4)
		r1 = MATRIX(4, 0)
	case '>':
		r2 = MATRIX(6, 4)
		r1 = MATRIX(5, 4)
	case '<':
		r2 = MATRIX(6, 4)
		r1 = MATRIX(5, 7)
	case '?':
		r2 = MATRIX(6, 4)
		r1 = MATRIX(6, 2)
	case '[':
		r2 = MATRIX(6, 4)
		r1 = MATRIX(5, 5)
	case ']':
		r2 = MATRIX(6, 4)
		r1 = MATRIX(6, 2)
	case '\n':
		r1 = MATRIX(0, 1)
	case '\r':
		r1 = MATRIX(0, 1)
	case ' ':
		r1 = MATRIX(7, 4)
	case '/':
		r1 = MATRIX(6, 7)
	case '^':
		r1 = MATRIX(6, 6)
	case '=':
		r1 = MATRIX(6, 5)
	case ';':
		r1 = MATRIX(6, 2)
	case '*':
		r1 = MATRIX(6, 1)
	//case '£':  r1 = MATRIX(6, 0)
	case ',':
		r1 = MATRIX(5, 7)
	case '@':
		r1 = MATRIX(5, 6)
	case ':':
		r1 = MATRIX(5, 5)
	case '.':
		r1 = MATRIX(5, 4)
	case '-':
		r1 = MATRIX(5, 3)
	case '+':
		r1 = MATRIX(5, 0)
	case '0':
		r1 = MATRIX(4, 3)
	case '1':
		r1 = MATRIX(7, 0)
	case '2':
		r1 = MATRIX(7, 3)
	case '3':
		r1 = MATRIX(1, 0)
	case '4':
		r1 = MATRIX(1, 3)
	case '5':
		r1 = MATRIX(2, 0)
	case '6':
		r1 = MATRIX(2, 3)
	case '7':
		r1 = MATRIX(3, 0)
	case '8':
		r1 = MATRIX(3, 3)
	case '9':
		r1 = MATRIX(4, 0)
	case 'A':
		r1 = MATRIX(1, 2)
	case 'B':
		r1 = MATRIX(3, 4)
	case 'C':
		r1 = MATRIX(2, 4)
	case 'D':
		r1 = MATRIX(2, 2)
	case 'E':
		r1 = MATRIX(1, 6)
	case 'F':
		r1 = MATRIX(2, 5)
	case 'G':
		r1 = MATRIX(3, 2)
	case 'H':
		r1 = MATRIX(3, 5)
	case 'I':
		r1 = MATRIX(4, 1)
	case 'J':
		r1 = MATRIX(4, 2)
	case 'K':
		r1 = MATRIX(4, 5)
	case 'L':
		r1 = MATRIX(5, 2)
	case 'M':
		r1 = MATRIX(4, 4)
	case 'N':
		r1 = MATRIX(4, 7)
	case 'O':
		r1 = MATRIX(4, 6)
	case 'P':
		r1 = MATRIX(5, 1)
	case 'Q':
		r1 = MATRIX(7, 6)
	case 'R':
		r1 = MATRIX(2, 1)
	case 'S':
		r1 = MATRIX(1, 5)
	case 'T':
		r1 = MATRIX(2, 6)
	case 'U':
		r1 = MATRIX(3, 6)
	case 'V':
		r1 = MATRIX(3, 7)
	case 'W':
		r1 = MATRIX(1, 1)
	case 'X':
		r1 = MATRIX(2, 7)
	case 'Y':
		r1 = MATRIX(3, 1)
	case 'Z':
		r1 = MATRIX(1, 4)
	}
	return r1, r2, r1 != 0
}
