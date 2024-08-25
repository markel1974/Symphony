package keyboard

import (
	"container/list"
)

func matrix(a int, b int) int {
	return ((a) << 3) | (b)
}

type Keyboard struct {
	keyDataStorage    *list.List
	joyKeyDataStorage *list.List
	inputReady        bool
	srcPre            *list.List
	srcPost           *list.List
	poolCounter       int
	joyKey            int
	ready             bool
	joystickSwap      bool
	virtual           *Virtual
	ascii             *Ascii
}

func NewKeyboard() *Keyboard {
	k := &Keyboard{
		virtual: NewVirtual(),
		ascii:   NewAscii(),
	}
	k.Reset()
	return k
}

func (k *Keyboard) Reset() {
	k.ready = false
	k.keyDataStorage = list.New()
	k.joyKeyDataStorage = list.New()
	k.inputReady = false
	k.srcPre = list.New()
	k.srcPost = list.New()
	k.poolCounter = 0
	k.joyKey = 0xff
	k.joystickSwap = true
	k.virtual.Reset()
	k.ascii.Reset()
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

func (k *Keyboard) SetExt() {
	k.virtual.SetExt()
}

func (k *Keyboard) SetNumLock() {
	k.virtual.SetNumLock()
}

func (k *Keyboard) SetCapital() {
	k.virtual.SetCapital()
}

func (k *Keyboard) SetMenu() {
	k.virtual.SetMenu()
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
	kc := k.virtual.FromVirtual(vKey)
	if kc < 0 {
		return
	}
	if pressed {
		if c64Byte, c64Bit, shifted, ok := k.keyDown(kc); ok {
			k.keyDataStorage.PushBack(NewInputKeyData(pressed, c64Byte, c64Bit, shifted))
		}
	} else {
		if c64Byte, c64Bit, shifted, ok := k.keyUp(kc); ok {
			k.keyDataStorage.PushBack(NewInputKeyData(pressed, c64Byte, c64Bit, shifted))
		}
	}
}

func (k *Keyboard) SetJoyKey(pressed bool, jKey int) {
	if k.joyKeyDataStorage.Len() >= MAX_STORAGE_SIZE {
		return
	}
	if pressed {
		if joyKey, ok := k.joyKeyDown(jKey); ok {
			k.joyKeyDataStorage.PushBack(uint8(joyKey))
		}
	} else {
		if joyKey, ok := k.joyKeyUp(jKey); ok {
			k.joyKeyDataStorage.PushBack(uint8(joyKey))
		}
	}
}

func (k *Keyboard) PollJoyKey() (uint8, bool) {
	if k.joyKeyDataStorage.Len() > 0 {
		e := k.joyKeyDataStorage.Front()
		joyKey := e.Value.(uint8)
		k.joyKeyDataStorage.Remove(e)
		return joyKey, true
	}
	return 0xff, false
}

func (k *Keyboard) PollKeyboard() (int, int, bool, bool, bool) {
	if k.poolCounter > 0 {
		k.poolCounter--
	}
	if k.poolCounter == 0 && k.keyDataStorage.Len() > 0 {
		e := k.keyDataStorage.Front()
		i := e.Value.(*InputKeyData)
		c64Byte := i.c64Byte
		c64Bit := i.c64Bit
		shifted := i.shifted
		pressed := i.pressed
		k.poolCounter = int(i.counter)
		k.keyDataStorage.Remove(e)
		return c64Byte, c64Bit, pressed, shifted, true
	}
	return 0, 0, false, false, false
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

func (k *Keyboard) joyKeyUp(kc int) (int, bool) {
	switch kc {
	case KEY_FIRE:
		k.joyKey |= 0x10
		return k.joyKey, true
	case KEY_JUP:
		k.joyKey |= 0x01
		return k.joyKey, true
	case KEY_JDN:
		k.joyKey |= 0x02
		return k.joyKey, true
	case KEY_JLF:
		k.joyKey |= 0x04
		return k.joyKey, true
	case KEY_JRT:
		k.joyKey |= 0x08
		return k.joyKey, true
	case KEY_JUPLF:
		k.joyKey |= 0x05
		return k.joyKey, true
	case KEY_JUPRT:
		k.joyKey |= 0x09
		return k.joyKey, true
	case KEY_JDNLF:
		k.joyKey |= 0x06
		return k.joyKey, true
	case KEY_JDNRT:
		k.joyKey |= 0x0a
		return k.joyKey, true
	}
	return 0xff, false
}

func (k *Keyboard) joyKeyDown(kc int) (int, bool) {
	switch kc {
	case KEY_FIRE:
		k.joyKey &= ^0x10
		return k.joyKey, true
	case KEY_JUP:
		k.joyKey |= 0x02
		k.joyKey &= ^0x01
		return k.joyKey, true
	case KEY_JDN:
		k.joyKey |= 0x01
		k.joyKey &= ^0x02
		return k.joyKey, true
	case KEY_JLF:
		k.joyKey |= 0x08
		k.joyKey &= ^0x04
		return k.joyKey, true
	case KEY_JRT:
		k.joyKey |= 0x04
		k.joyKey &= ^0x08
		return k.joyKey, true
	case KEY_JUPLF:
		k.joyKey |= 0x0a
		k.joyKey &= ^0x05
		return k.joyKey, true
	case KEY_JUPRT:
		k.joyKey |= 0x06
		k.joyKey &= ^0x09
		return k.joyKey, true
	case KEY_JDNLF:
		k.joyKey |= 0x09
		k.joyKey &= ^0x06
		return k.joyKey, true
	case KEY_JDNRT:
		k.joyKey |= 0x05
		k.joyKey &= ^0x0a
		return k.joyKey, true
	case KEY_CENTER:
		k.joyKey |= 0x0f
		return k.joyKey, true
	}
	return 0xff, false
}

func (k *Keyboard) keyUp(kc int) (int, int, bool, bool) {
	return k.keyCodeToC64(kc)
}

func (k *Keyboard) keyDown(kc int) (int, int, bool, bool) {
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
	default:
		return k.keyCodeToC64(kc)
	}
	return 0, 0, false, false
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
	cmd := func(keyCode int, pressed bool, shifted bool, storage *list.List) bool {
		if c64Byte, c64Bit, _, ok := k.keyCodeToC64(keyCode); ok {
			i1 := NewInputKeyData(pressed, c64Byte, c64Bit, shifted)
			storage.PushBack(i1)
			return true
		}
		return false
	}
	storage := list.New()
	for _, it := range command {
		v := k.ascii.FromAscii(uint8(it))
		d1 := v.r1
		shifted := v.shifted
		if ret := cmd(d1, true, shifted, storage); !ret {
			return nil, false
		}
		if ret := cmd(d1, false, shifted, storage); !ret {
			return nil, false
		}
	}
	return storage, true
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
