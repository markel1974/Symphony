package mos6569

import (
	"github.com/markel1974/c64emu/src/kernel/component"
	"github.com/markel1974/c64emu/src/references"
)

// LightPen represents a programmable light pen interface for graphical systems, capturing X and Y coordinates on trigger.
// x stores the X coordinate of the light pen for the current frame.
// y stores the Y coordinate of the light pen for the current frame.
// triggered is a flag indicating if the light pen was triggered in the current frame.
// interruptsEmit is a function used to emit interrupts when the light pen is triggered.
type LightPen struct {
	*component.BaseComponent
	reflect        *LightPenReflect
	x              uint8
	y              uint8
	triggered      bool
	interruptsEmit func(uint8)
}

// NewLightPen initializes and returns a new LightPen instance with the provided interrupt emission function.
func NewLightPen(parent references.IComponent, factory references.IComponentFactory, label string, instance int, interruptsEmit func(uint8)) *LightPen {
	lp := &LightPen{
		BaseComponent:  component.NewBaseComponent(),
		interruptsEmit: interruptsEmit,
		x:              0,
		y:              0,
		triggered:      false,
	}
	lp.reflect = NewLightPenReflect(lp, factory, parent, "lightPen", instance, references.IdInternalComponent(label, instance, "LightPen"))
	return lp
}

// Setup initializes the LightPen and prepares it for operation, returning an error if initialization fails.
func (lp *LightPen) Setup() error {
	return nil
}

// Connect establishes the connection for the LightPen component, preparing it for use within the system's context.
func (lp *LightPen) Connect() error {
	return nil
}

// EmulationRequired determines whether emulation routines are necessary for the LightPen instance.
func (lp *LightPen) EmulationRequired() bool {
	return false
}

// Emulate processes the light pen's behavior for the current emulation frame, checking its state and coordinating updates.
func (lp *LightPen) Emulate() {
}

// Internal determines if the LightPen is an internal device, always returning true.
func (lp *LightPen) Internal() bool {
	return true
}

// Reset reinitializes the LightPen's state, clearing its coordinates and triggered status to defaults.
func (lp *LightPen) Reset() {
}

// Trigger sets the light pen as triggered and records the given raster coordinates (rasterX, rasterY).
// It also emits an interrupt using the specified light pen interrupt bit.
// If the light pen is already triggered, this method does nothing.
func (lp *LightPen) Trigger(rasterX uint16, rasterY uint16) {
	if !lp.triggered {
		lp.triggered = true
		lp.x = uint8(rasterX >> 1)
		lp.y = uint8(rasterY)
		lp.interruptsEmit(irqLightPenBit)
	}
}

// TriggerClear resets the light pen's triggered state to false, indicating it has not been activated.
func (lp *LightPen) TriggerClear() {
	lp.triggered = false
}

// ReadX retrieves the current X-coordinate value of the light pen.
func (lp *LightPen) ReadX() uint8 {
	return lp.x
}

// ReadY returns the Y-coordinate stored by the LightPen, representing the last triggered position on the raster.
func (lp *LightPen) ReadY() uint8 {
	return lp.y
}

// WriteX sets the X-coordinate of the LightPen by updating the internal lpx register with the provided data.
func (lp *LightPen) WriteX(data uint8) {
	lp.x = data
}

// WriteY sets the LightPen's internal lpx register to the specified data value.
func (lp *LightPen) WriteY(data uint8) {
	lp.x = data
}
