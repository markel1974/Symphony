package inputs

import (
	"container/list"
	"github.com/markel1974/c64emu/src/components/quartz"
)

func matrix(a int, b int) int {
	return ((a) << 3) | (b)
}

type Keyboard struct {
	keyDataStorage *list.List
	virtual        *Virtual
	ascii          *Ascii
	quartz         *quartz.Quartz
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
	k.keyDataStorage = list.New()
	k.virtual.Reset()
	k.ascii.Reset()
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

func (k *Keyboard) SetVirtualKey(pressed bool, vKey int) {
	if k.keyDataStorage.Len() >= MAX_STORAGE_SIZE {
		return
	}
	kc := k.virtual.FromVirtual(vKey)
	if kc < 0 {
		return
	}
	if v, ok := k.keyCodeToC64(kc, 0, pressed); ok {
		k.keyDataStorage.PushBack(v)
	}
	return
}

func (k *Keyboard) Poll() (uint32, bool) {
	if k.keyDataStorage.Len() == 0 {
		return 0, false
	}
	e := k.keyDataStorage.Front()
	i := e.Value.(uint32)
	k.keyDataStorage.Remove(e)
	return i, true
}

func (k *Keyboard) SetCommand(cmd string) {
	for _, it := range cmd {
		v := k.ascii.FromAscii(uint8(it))
		p1, ok1 := k.keyCodeToC64(v.r1, v.shifted, true)
		p2, ok2 := k.keyCodeToC64(v.r1, v.shifted, false)
		if ok1 && ok2 {
			k.keyDataStorage.PushBack(p1)
			k.keyDataStorage.PushBack(p2)
		}
	}
}

func (k *Keyboard) keyCodeToC64(kc int, shifted uint8, pressed bool) (uint32, bool) {
	if kc >= 0 || kc < 256 {
		keyM := uint8((kc >> 3) & 7)
		revM := uint8(kc & 7)
		out := uint32(keyM)
		out |= uint32(revM) << 8
		if shifted == 1 {
			out |= 0x10000 //shifted
		} else if (kc & 128) != 0 {
			out |= 0x10000 //shifted
		}
		if pressed {
			out |= 0x20000
		}
		return out, true
	}
	return 0, false
}
