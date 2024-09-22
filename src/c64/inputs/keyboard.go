package inputs

import (
	"github.com/markel1974/c64emu/src/fifo"
)

func matrix(a int, b int) int {
	return ((a) << 3) | (b)
}

type Keyboard struct {
	storage *fifo.Queue
	virtual *Virtual
	ascii   *Ascii
}

func NewKeyboard() *Keyboard {
	k := &Keyboard{
		storage: nil,
		virtual: NewVirtual(),
		ascii:   NewAscii(),
	}
	k.Reset()
	return k
}

func (k *Keyboard) Reset() {
	k.storage = fifo.NewQueue(16384)
	k.virtual.Reset()
	k.ascii.Reset()
}

func (k *Keyboard) NumLockToggle() {
	k.virtual.NumLockToggle()
}

func (k *Keyboard) CapitalToggle() {
	k.virtual.CapitalToggle()
}

func (k *Keyboard) SetVirtualKey(pressed bool, vKey int) {
	if kc := k.virtual.FromVirtual(vKey); kc >= 0 {
		v := keyCodeToC64(uint8(kc), -1, pressed)
		k.storage.Add(int(v))
	}
}

func (k *Keyboard) Poll() (uint32, bool) {
	if k.storage.Len() == 0 {
		return 0, false
	}
	key, ok := k.storage.Next()
	if !ok {
		return 0, false
	}
	return uint32(key), true
}

func (k *Keyboard) SetCommand(cmd string) {
	for _, c := range cmd {
		v := k.ascii.FromAscii(uint8(c))
		p1 := keyCodeToC64(uint8(v.r1), v.shifted, true)
		if !k.storage.Add(int(p1)) {
			return
		}
		p2 := keyCodeToC64(uint8(v.r1), v.shifted, false)
		if !k.storage.Add(int(p2)) {
			return
		}
	}
}

func keyCodeToC64(kc uint8, shifted int, pressed bool) uint32 {
	keyM := (kc >> 3) & 7
	revM := kc & 7
	out := uint32(keyM)
	out |= uint32(revM) << 8
	if shifted > 0 {
		out |= 0x10000 //shifted
	} else if shifted < 0 {
		if (kc & 128) != 0 {
			out |= 0x10000 //shifted
		}
	}
	if pressed {
		out |= 0x20000
	}
	return out
}

//func (k *Keyboard) SetExt() {
//	k.virtual.SetExt()
//}

//func (k *Keyboard) SetMenu() {
//	k.virtual.SetMenu()
//}
