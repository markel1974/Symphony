package mos6569

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

// borderWidth defines the width of each border segment in pixels.
// borderCountMax determines the maximum number of border segments that can fit horizontally on the display.
// borderCount specifies the predefined number of border segments currently being used.
const (
	borderWidth = 8
	borderCount = 4
)

// Border bit constants represent various positions within a layout or grid, using increasing iota values.
const (
	borderBitLeft = iota
	borderBitMidLeft
	borderBitCenter
	borderBitMidRight
	borderBitRight
)

// sequencerLength defines the size of the sequencer array, calculated as 2^5, providing 32 possible states.
const (
	borderSequencerSize  = 31
	borderSequencerCount = borderSequencerSize + 1 //1 << 5
)

// Borders is a type responsible for managing and updating border data, configurations, and states for a display system.
type Borders struct {
	*component.BaseComponent
	core               *VIC
	setMulti8          func(int, uint8)
	horizontalFlipFlop uint8
	verticalFlipFlop   uint8
	colors             [0xff]uint8
	offset             int
	left               []int
	midLeft            int
	center             []int
	midRight           int
	right              []int
	sequencer          [borderSequencerCount][]func()
	sequencerState     uint8
}

// NewBorder initializes and returns a new Borders object using the provided VIC core and display buffer interface.
// It configures left, right, center, and sequencer states based on display buffer properties.
func NewBorder(parent references.IComponent, factory references.IComponentFactory, label string, instance int, core *VIC, displayBuffer references.IDisplayBuffer, displayX int) *Borders {
	borderCountMax := displayX / borderWidth
	gr := &Borders{
		BaseComponent:      component.NewBaseComponent(),
		setMulti8:          displayBuffer.SetMulti8,
		core:               core,
		horizontalFlipFlop: 0,
		verticalFlipFlop:   0,
		offset:             0,
		sequencerState:     0,
	}
	for x := 0; x < borderCount; x++ {
		gr.left = append(gr.left, x)
	}
	for x := borderCount - 1; x >= 0; x-- {
		gr.right = append(gr.right, (borderCountMax-1)-x)
	}
	gr.midLeft = gr.left[len(gr.left)-1] + 1
	gr.midRight = gr.right[0] - 1
	for x := gr.midLeft + 1; x < gr.midRight; x++ {
		gr.center = append(gr.center, x)
	}
	gr.sequencer = gr.createSequencer()
	gr.BaseComponent.Register(factory, parent, "border", gr, references.IdInternalComponent(label, instance, "Border"))
	return gr
}

// Setup initializes the Borders instance and prepares it for operation. It returns an error if the setup fails.
func (b *Borders) Setup() error {
	return nil
}

// Connect establishes the necessary connections or associations for the Borders instance and returns an error if it fails.
func (b *Borders) Connect() error {
	return nil
}

// EmulationRequired checks if emulation is required for the Borders instance, always returning false.
func (b *Borders) EmulationRequired() bool {
	return false
}

// Emulate performs the main emulation logic for the Borders instance, processing updates and managing its internal state.
func (b *Borders) Emulate() {
}

// Internal returns true if the internal border logic is currently active.
func (b *Borders) Internal() bool {
	return true
}

// Reset clears and reinitializes the internal state of the Borders instance, preparing it for subsequent operations.
func (b *Borders) Reset() {
}

// ColumnInitialize resets and updates the sequencerState based on the horizontalFlipFlop value for the left border bit.
func (b *Borders) ColumnInitialize() {
	const bitNumber = borderBitLeft
	b.sequencerState &^= 1 << bitNumber
	b.sequencerState |= b.horizontalFlipFlop << bitNumber
}

// Column38Update updates the sequencer state for column 38 based on the horizontal and vertical flip-flop states.
// If column mode is not selected, it updates the vertical flip-flop value and adjusts the horizontal flip-flop accordingly.
func (b *Borders) Column38Update() {
	if !b.core.columnSel {
		b.UpdateVerticalFlipFlop()
		if b.verticalFlipFlop == 0 {
			b.horizontalFlipFlop = 0
		}
	}
	const bitNumber = borderBitCenter
	b.sequencerState &^= 1 << bitNumber
	b.sequencerState |= b.horizontalFlipFlop << bitNumber
}

// Column40Update updates the state of the mid-left border column based on the flip-flop and column selection logic.
func (b *Borders) Column40Update() {
	if b.core.columnSel {
		b.UpdateVerticalFlipFlop()
		if b.verticalFlipFlop == 0 {
			b.horizontalFlipFlop = 0
		}
	}
	const bitNumber = borderBitMidLeft
	b.sequencerState &^= 1 << bitNumber
	b.sequencerState |= b.horizontalFlipFlop << bitNumber
}

// Column38Apply updates the sequence state for the mid-right border bit, conditionally setting the horizontal flip-flop.
func (b *Borders) Column38Apply() {
	if !b.core.columnSel {
		b.horizontalFlipFlop = 1
	}
	const bitNumber = borderBitMidRight
	b.sequencerState &^= 1 << bitNumber
	b.sequencerState |= b.horizontalFlipFlop << bitNumber
}

// Column40Apply updates the sequencer state for the right border column, adjusting horizontal flip-flop if conditions are met.
func (b *Borders) Column40Apply() {
	if b.core.columnSel {
		b.horizontalFlipFlop = 1
	}
	const bitNumber = borderBitRight
	b.sequencerState &^= 1 << bitNumber
	b.sequencerState |= b.horizontalFlipFlop << bitNumber
}

// SetOffset sets the offset value for the Borders instance. It determines the starting point for border rendering.
func (b *Borders) SetOffset(offset int) {
	b.offset = offset
}

// AcquireColor assigns a color to the specified index in the borders color array using the current execution context.
func (b *Borders) AcquireColor(idx uint8) {
	b.colors[idx] = _colors[b.core.ec]
}

// UpdateVerticalFlipFlop updates the vertical border flip-flop state based on the current raster Y coordinate and control flags.
func (b *Borders) UpdateVerticalFlipFlop() {
	//3.9. The border unit
	if b.core.dyBottom == b.core.rasterY {
		//2. If the Y coordinate reaches the bottom comparison value in cycle 63, the vertical border flip flop is set.
		b.verticalFlipFlop = 1
	} else if b.core.dyTop == b.core.rasterY && b.core.den {
		//3. If the Y coordinate reaches the top comparison value in cycle 63 and the DEN bit in register $d011 is set, the vertical border flip flop is reset.
		b.verticalFlipFlop = 0
	}
}

// VerticalFlipFlop returns true if the vertical border flip-flop is set; otherwise, it returns false.
func (b *Borders) VerticalFlipFlop() bool {
	return b.verticalFlipFlop != 0
}

// Draw executes a sequence of rendering functions based on the current sequencer state of the Borders instance.
func (b *Borders) Draw() {
	sequence := b.sequencer[b.sequencerState&borderSequencerSize]
	for _, drawFn := range sequence {
		drawFn()
	}
}

// drawLeft iterates over the left border indices, calculates offsets, and applies corresponding colors using setMulti8.
func (b *Borders) drawLeft() {
	for _, v := range b.left {
		offset := v * borderWidth
		b.setMulti8(b.offset+offset, b.colors[v])
	}
}

// drawMidLeft updates the border area corresponding to midLeft by setting the appropriate color and offset values.
func (b *Borders) drawMidLeft() {
	offset := b.midLeft * borderWidth
	b.setMulti8(b.offset+offset, b.colors[b.midLeft])
}

// drawCenter processes the center border segments by calculating their offsets and applying corresponding colors.
func (b *Borders) drawCenter() {
	for _, v := range b.center {
		offset := v * borderWidth
		b.setMulti8(b.offset+offset, b.colors[v])
	}
}

// drawMidRight calculates the offset for the mid-right border and updates its color using the setMulti8 function.
func (b *Borders) drawMidRight() {
	offset := b.midRight * borderWidth
	b.setMulti8(b.offset+offset, b.colors[b.midRight])
}

// drawRight renders the right-hand border by iterating over the `right` field and setting appropriate color values.
func (b *Borders) drawRight() {
	for _, v := range b.right {
		offset := v * borderWidth
		b.setMulti8(b.offset+offset, b.colors[v])
	}
}

// drawEmpty clears or bypasses border drawing by executing no operations during the rendering process.
func (b *Borders) drawEmpty() {

}

// createSequencer initializes and returns a 2D slice of function sequences for border rendering based on state bits.
func (b *Borders) createSequencer() [borderSequencerCount][]func() {
	const left = 1 << borderBitLeft
	const midLeft = 1 << borderBitMidLeft
	const center = 1 << borderBitCenter
	const midRight = 1 << borderBitMidRight
	const right = 1 << borderBitRight

	var sequencer [borderSequencerCount][]func()
	for idx := range sequencer {
		x := uint8(idx)
		var data []func() = nil
		if (x & left) == left {
			data = append(data, b.drawLeft)
		}
		if (x & midLeft) == midLeft {
			data = append(data, b.drawMidLeft)
		}
		if (x & center) == center {
			data = append(data, b.drawCenter)
		}
		if (x & midRight) == midRight {
			data = append(data, b.drawMidRight)
		}
		if (x & right) == right {
			data = append(data, b.drawRight)
		}
		if data == nil {
			data = append(data, b.drawEmpty)
		}
		sequencer[x] = data
	}
	return sequencer
}
