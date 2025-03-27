package joystick_c64

import (
	"github.com/markel1974/c64emu/src/common/fifo"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// Joystick represents an input device for handling directional and button presses with customizable sensitivity settings.
type Joystick struct {
	*component.BaseComponent
	storage *fifo.StaticFifo
	joy     int
	s1      uint
	s2      uint
}

// NewJoystick initializes and returns a new instance of Joystick with default settings and sensitivity for controls.
func NewJoystick(parent references.IComponent, factory references.IComponentFactory, instance int) *Joystick {
	j := &Joystick{
		BaseComponent: component.NewBaseComponent(),
		storage:       nil,
		joy:           0xff,
		s1:            0,
		s2:            0,
	}
	j.BaseComponent.Register(factory, parent, Identifier(), j, references.IdIJoystick(j, instance))
	return j
}

func (k *Joystick) Setup(_ references.IJoystickSocket, _ *config.Config) error {
	k.Update(0x0000, 0xffff, 40)
	return nil
}

func (k *Joystick) Connect() error {
	return nil
}

func (k *Joystick) Internal() bool {
	return false
}

// Update recalculates the joystick's sensitivity range based on the provided minimum, maximum, and sensitivity values.
func (k *Joystick) Update(min uint16, max uint16, sensitivity uint16) {
	interval := max - min
	k.s1 = uint(min + ((sensitivity * interval) / 100))
	k.s2 = uint(min + (((100 - sensitivity) * interval) / 100))
}

// Reset reinitializes the joystick's storage buffer and sets its state to the default value.
func (k *Joystick) Reset() {
	k.storage = fifo.NewStaticFifo(256)
	k.joy = 0xff
}

func (k *Joystick) Emulate() {
	//
}

func (k *Joystick) EmulationRequired() bool {
	return false
}

// Move updates the joystick state based on x, y positions and button inputs, then stores the updated state in storage.
func (k *Joystick) Move(x uint, y uint, buttons uint) {
	k.joy = 0xff
	if x < k.s1 {
		k.joy &= 0xfb // Left
	} else if x > k.s2 {
		k.joy &= 0xf7 // Right
	}
	if y < k.s1 {
		k.joy &= 0xfe // Up
	} else if y > k.s2 {
		k.joy &= 0xfd // Down
	}
	if (buttons & 1) != 0 {
		k.joy &= 0xef // Button
	}
	if (buttons & 2) != 0 {
		//TODO SID POTX / POTY
	}
	k.storage.Set(k.joy)
}

// SetKey updates the joystick's state based on the key press or release and stores the new state in the storage.
func (k *Joystick) SetKey(pressed bool, jId int) {
	if pressed {
		k.joy = joyKeyDown(k.joy, jId)
		k.storage.Set(k.joy)
	} else {
		k.joy = joyKeyUp(k.joy, jId)
		k.storage.Set(k.joy)
	}
}

// Poll retrieves and returns the next joystick state and its validity. It returns (0, false) if no state is available.
func (k *Joystick) Poll() (uint8, bool) {
	if k.storage.Len() == 0 {
		return 0, false
	}
	joy, ok := k.storage.Next()
	if !ok {
		return 0, false
	}
	return uint8(joy), true
}

// joyKeyUp processes a key code to update the joystick state, turning off specific bits based on provided key codes.
func joyKeyUp(j int, kc int) int {
	switch kc {
	case component.KeyJFire:
		j |= 0x10
		return j
	case component.KeyJUp:
		j |= 0x01
		return j
	case component.KeyJDown:
		j |= 0x02
		return j
	case component.KeyJLeft:
		j |= 0x04
		return j
	case component.KeyJRight:
		j |= 0x08
		return j
	case component.KeyJUpLeft:
		j |= 0x05
		return j
	case component.KeyJUpRight:
		j |= 0x09
		return j
	case component.KeyJDownLeft:
		j |= 0x06
		return j
	case component.KeyJDownRight:
		j |= 0x0a
		return j
	case component.KeyJCenter:
		return 0xff
	}
	return 0xff
}

// joyKeyDown adjusts the joystick state by setting the specified key as pressed and updating directional bits accordingly.
func joyKeyDown(j int, kc int) int {
	switch kc {
	case component.KeyJFire:
		j &= ^0x10
		return j
	case component.KeyJUp:
		j |= 0x02
		j &= ^0x01
		return j
	case component.KeyJDown:
		j |= 0x01
		j &= ^0x02
		return j
	case component.KeyJLeft:
		j |= 0x08
		j &= ^0x04
		return j
	case component.KeyJRight:
		j |= 0x04
		j &= ^0x08
		return j
	case component.KeyJUpLeft:
		j |= 0x0a
		j &= ^0x05
		return j
	case component.KeyJUpRight:
		j |= 0x06
		j &= ^0x09
		return j
	case component.KeyJDownLeft:
		j |= 0x09
		j &= ^0x06
		return j
	case component.KeyJDownRight:
		j |= 0x05
		j &= ^0x0a
		return j
	case component.KeyJCenter:
		j |= 0x0f
		return j
	}
	return 0xff
}
