package mechanic

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/hardware/c1541_board_rev1/disk"
	"github.com/markel1974/c64emu/src/references"
)

// Motor represents a structure for managing the state and operation of a motor in a disk drive system.
// It tracks whether the motor is active, the time required for it to spin up, and the rotation cycles during operation.
type Motor struct {
	*component.BaseComponent
	active         bool
	spinUpTime     int
	rotationCycles int
}

// NewMotor initializes and returns a new Motor instance with default values and resets its state.
func NewMotor(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Motor {
	v := &Motor{
		BaseComponent:  component.NewBaseComponent(),
		active:         false,
		spinUpTime:     0,
		rotationCycles: 0,
	}
	v.BaseComponent.Register(factory, parent, "motor", v, references.IdInternalComponent(label, instance, "Motor"))

	v.Reset()
	return v
}

func (m *Motor) Setup() error {
	return nil
}

func (m *Motor) Connect() error {
	return nil
}

func (m *Motor) EmulationRequired() bool {
	return false
}

func (m *Motor) Emulate() {
}

func (m *Motor) Internal() bool {
	return true
}

// Reset reinitializes the motor's state, setting it to inactive and clearing spin-up time and rotation cycles.
func (m *Motor) Reset() {
	m.active = false
	m.spinUpTime = 0
	m.rotationCycles = 0
}

// IsActive returns true if the motor is currently active, otherwise false.
func (m *Motor) IsActive() bool {
	return m.active
}

// SetActive sets the motor's active state. Initiates spin-up delay if activating the motor from an inactive state.
func (m *Motor) SetActive(active bool) {
	if active && !m.active {
		m.spinUpTime = motorSpinUpDelay
	}
	m.active = active
}

// SpinUp decreases the motor's spin-up time if greater than 0, returning true if the motor is still spinning up.
func (m *Motor) SpinUp() bool {
	if m.spinUpTime > 0 {
		m.spinUpTime--
		return true
	}
	return false
}

// TryRotate decrements the motor's rotation cycles and returns true if rotation is complete, otherwise returns false.
func (m *Motor) TryRotate() bool {
	if m.rotationCycles--; m.rotationCycles > 0 {
		return false
	}
	return true
}

// Rotate advances the motor's rotation cycles based on disk timing and triggers the disk to perform its rotate operation.
func (m *Motor) Rotate(disk disk.IDisk) {
	m.rotationCycles += disk.MicroSecPerByte()
	disk.Rotate()
	return
}
