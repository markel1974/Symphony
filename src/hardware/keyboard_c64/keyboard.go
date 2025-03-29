package keyboard_c64

import (
	"github.com/markel1974/c64emu/src/common/fifo"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

// matrix computes and returns a combined integer value by shifting `a` to the left by 3 bits and performing a bitwise OR with `b`.
func matrix(a int, b int) int {
	return ((a) << 3) | (b)
}

// Keyboard represents an abstraction for handling virtual and ASCII keyboard states and input storage.
type Keyboard struct {
	*component.BaseComponent
	storage *fifo.StaticFifo
	virtual *Virtual
	ascii   *Ascii
}

// NewKeyboard initializes and returns a new Keyboard instance with default settings and a reset state.
func NewKeyboard(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Keyboard {
	k := &Keyboard{
		BaseComponent: component.NewBaseComponent(),
		storage:       nil,
		virtual:       NewVirtual(),
		ascii:         NewAscii(),
	}
	k.BaseComponent.Register(factory, parent, Identifier(), k, references.IdIKeyboard(k, label, instance))
	return k
}

func (k *Keyboard) Setup() error {
	k.Reset()
	return nil
}

func (k *Keyboard) Bind(_ references.IKeyboardSocket) error {
	return nil
}

func (k *Keyboard) Connect() error {
	return nil
}

func (k *Keyboard) Internal() bool {
	return false
}

// Reset reinitializes the Keyboard by resetting its storage, virtual key states, and ASCII translations.
func (k *Keyboard) Reset() {
	k.storage = fifo.NewStaticFifo(16384)
	k.virtual.Reset()
	k.ascii.Reset()
}

func (k *Keyboard) Emulate() {
	//
}

func (k *Keyboard) EmulationRequired() bool {
	return false
}

// NumLockToggle toggles the current state of the Num Lock key by inverting its value on the virtual keyboard.
func (k *Keyboard) NumLockToggle() {
	k.virtual.NumLockToggle()
}

// CapitalToggle toggles the state of the capital key (caps lock) on the virtual keyboard.
func (k *Keyboard) CapitalToggle() {
	k.virtual.CapitalToggle()
}

// SetKey processes a virtual key press or release by mapping it to a key code and updating the keyboard's storage.
func (k *Keyboard) SetKey(pressed bool, vKey int) {
	if kc := k.virtual.FromVirtual(vKey); kc >= 0 {
		v := keyCodeToC64(uint8(kc), -1, pressed)
		k.storage.Set(int(v))
	}
}

// Poll retrieves the next key from the keyboard storage as a uint32 and indicates if a key was available.
// Returns the key value and true if a key was retrieved, otherwise returns 0 and false.
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

// SetCommand processes a given string command, converting each character to corresponding key codes and storing them.
func (k *Keyboard) SetCommand(cmd string) {
	for _, c := range cmd {
		v := k.ascii.FromAscii(uint8(c))
		p1 := keyCodeToC64(uint8(v.r1), v.shifted, true)
		if !k.storage.Set(int(p1)) {
			return
		}
		p2 := keyCodeToC64(uint8(v.r1), v.shifted, false)
		if !k.storage.Set(int(p2)) {
			return
		}
	}
}

// keyCodeToC64 converts a given key code, shift state, and pressed state into a C64-compatible key representation.
// kc is the key code, shifted indicates shift state (-1 for auto, >0 for shifted), and pressed indicates press state.
// Returns a uint32 value encoding the key and its state (shifted, pressed).
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
