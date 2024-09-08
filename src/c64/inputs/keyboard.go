package inputs

import (
	"container/list"
)

func matrix(a int, b int) int {
	return ((a) << 3) | (b)
}

type Keyboard struct {
	keyDataStorage *list.List
	inputReady     bool
	srcPre         *list.List
	srcPost        *list.List
	joySwap        bool
	joy1           *Joystick
	joy2           *Joystick
	ready          bool
	virtual        *Virtual
	ascii          *Ascii
}

func NewKeyboard() *Keyboard {
	k := &Keyboard{
		virtual: NewVirtual(),
		ascii:   NewAscii(),
		joy1:    NewJoystick(),
		joy2:    NewJoystick(),
	}
	k.Reset()
	return k
}

func (k *Keyboard) Reset() {
	k.ready = false
	k.keyDataStorage = list.New()
	k.inputReady = false
	k.srcPre = list.New()
	k.srcPost = list.New()
	k.joy1.Reset()
	k.joy2.Reset()
	k.joySwap = true
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
	k.joySwap = !k.joySwap
	k.joy1.Reset()
	k.joy2.Reset()
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
		k.keyDataStorage.PushBack(NewInputKeyData(pressed, keyM, revM, shifted, 1))
	}
	return
}

func (k *Keyboard) SetJoystick1(pressed bool, c int) {
	if k.joySwap {
		k.joy2.Set(pressed, c)
	} else {
		k.joy1.Set(pressed, c)
	}
}

func (k *Keyboard) SetJoystick2(pressed bool, c int) {
	if k.joySwap {
		k.joy1.Set(pressed, c)
	} else {
		k.joy2.Set(pressed, c)
	}
}

func (k *Keyboard) PollJoystick1() (uint8, bool) {
	return k.joy1.Poll()
}

func (k *Keyboard) PollJoystick2() (uint8, bool) {
	return k.joy2.Poll()
}

func (k *Keyboard) PollKeyboard() (int, int, bool, bool, bool) {
	if k.keyDataStorage.Len() == 0 {
		return 0, 0, false, false, false
	}
	e := k.keyDataStorage.Front()
	//if e.Value == nil {
	//	return 0, 0, false, false, false
	//}
	i := e.Value.(*InputKeyData)
	if i.persistence--; i.persistence == 0 {
		k.keyDataStorage.Remove(e)
	}
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
	const persistenceCycles = 20
	if len(command) == 0 {
		return nil, false
	}
	cmd := func(keyCode int, pressed bool, shifted bool, storage *list.List) bool {
		if keyM, revM, _, ok := k.keyCodeToC64(keyCode); ok {
			i1 := NewInputKeyData(pressed, keyM, revM, shifted, persistenceCycles)
			storage.PushBack(i1)
			return true
		}
		/*
			if keyM, revM, _, ok := k.keyCodeToC64(keyCode); ok {
				for x := 0; x < persistence; x++ {
					i1 := NewInputKeyData(pressed, keyM, revM, shifted)
					storage.PushBack(i1)
				}
				return true
			}
		*/
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
