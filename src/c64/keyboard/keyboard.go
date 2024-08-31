package keyboard

import (
	"container/list"
)

func matrix(a int, b int) int {
	return ((a) << 3) | (b)
}

type Keyboard struct {
	keyDataStorage      *list.List
	joystickDataStorage *list.List
	inputReady          bool
	srcPre              *list.List
	srcPost             *list.List
	joyKey              int
	ready               bool
	joystickSwap        bool
	virtual             *Virtual
	ascii               *Ascii
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
	k.joystickDataStorage = list.New()
	k.inputReady = false
	k.srcPre = list.New()
	k.srcPost = list.New()
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
	k.joyKey = 0xff
}

func (k *Keyboard) SetVirtualKey(pressed bool, vKey int) {
	if k.keyDataStorage.Len() >= MAX_STORAGE_SIZE {
		return
	}
	kc := k.virtual.FromVirtual(vKey)
	if kc < 0 {
		return
	}
	if keyM, revM, shifted, ok := k.keyCodeToC64(kc); ok {
		k.keyDataStorage.PushBack(NewInputKeyData(pressed, keyM, revM, shifted))
	}
	return
}

func (k *Keyboard) SetJoystick(pressed bool, jId int) {
	if k.joystickDataStorage.Len() >= MAX_STORAGE_SIZE {
		return
	}
	if pressed {
		if joy, ok := k.joyDown(k.joyKey, jId); ok {
			k.joyKey = joy
			k.joystickDataStorage.PushBack(uint8(joy))
		}
	} else {
		if joy, ok := k.joyUp(k.joyKey, jId); ok {
			k.joyKey = joy
			k.joystickDataStorage.PushBack(uint8(joy))
		}
	}
}

func (k *Keyboard) PollJoysticks() (uint8, uint8, bool) {
	if k.joystickDataStorage.Len() == 0 {
		//TODO RETURN JOY
		return 0xff, 0xff, false
	}
	e := k.joystickDataStorage.Front()
	j := e.Value.(uint8)
	k.joystickDataStorage.Remove(e)
	if k.joystickSwap {
		return 0xff, j, true
	}
	return j, 0xff, true
}

func (k *Keyboard) PollKeyboard() (int, int, bool, bool, bool) {
	if k.keyDataStorage.Len() == 0 {
		return 0, 0, false, false, false
	}
	e := k.keyDataStorage.Front()
	i := e.Value.(*InputKeyData)
	k.keyDataStorage.Remove(e)
	return i.keyM, i.revM, i.pressed, i.shifted, true
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

func (k *Keyboard) joyUp(j int, kc int) (int, bool) {
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

func (k *Keyboard) joyDown(j int, kc int) (int, bool) {
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

func (k *Keyboard) keyCodeToC64(kc int) (int, int, bool, bool) {
	if kc >= 0 || kc < 256 {
		keyM := (kc >> 3) & 7
		revM := kc & 7
		shifted := (kc & 128) != 0
		return keyM, revM, shifted, true
	}
	return 0, 0, false, false
}

func (k *Keyboard) prepareCommand(command string) (*list.List, bool) {
	if len(command) == 0 {
		return nil, false
	}
	cmd := func(keyCode int, pressed bool, shifted bool, storage *list.List) bool {
		if keyM, revM, _, ok := k.keyCodeToC64(keyCode); ok {
			i1 := NewInputKeyData(pressed, keyM, revM, shifted)
			storage.PushBack(i1)
			return true
		}
		return false
	}
	storage := list.New()
	for _, it := range command {
		v := k.ascii.FromAscii(uint8(it))
		if ret := cmd(v.r1, true, v.shifted, storage); !ret {
			return nil, false
		}
		if ret := cmd(v.r1, false, v.shifted, storage); !ret {
			return nil, false
		}
	}
	return storage, true
}
