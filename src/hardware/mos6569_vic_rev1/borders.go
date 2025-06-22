package mos6569

import (
	"github.com/markel1974/c64emu/src/references"
)

const (
	borderWidth    = 8
	borderCountMax = DisplayX / borderWidth
	borderCount    = 4
)

// BorderType represents the type of border used, defined as an integer-based enumeration.
type BorderType int

// BorderTypeLeft represents the left border type.
// BorderTypeMidLeft represents the mid-left border type.
// BorderTypeCenter represents the center border type.
// BorderTypeMidRight represents the mid-right border type.
// BorderTypeRight represents the right border type.
// BorderTypeLast represents the last border type.
const (
	BorderTypeLeft     = BorderType(0)
	BorderTypeMidLeft  = BorderType(1)
	BorderTypeCenter   = BorderType(2)
	BorderTypeMidRight = BorderType(3)
	BorderTypeRight    = BorderType(4)
	BorderTypeLast     = BorderType(5)
)

// Borders represents the structure for handling visual border rendering in a VIC-based display system.
type Borders struct {
	core               *VIC
	setMulti8          func(int, uint8)
	horizontalFlipFlop bool
	verticalFlipFlop   bool
	samples            [BorderTypeLast]bool
	colors             [0xff]uint8
	offset             int

	left     []int
	midLeft  int
	center   []int
	midRight int
	right    []int
}

// NewBorder creates and initializes a new Borders instance with the provided VIC core and display buffer dependencies.
func NewBorder(core *VIC, displayBuffer references.IDisplayBuffer) *Borders {
	gr := &Borders{
		setMulti8:          displayBuffer.SetMulti8,
		core:               core,
		horizontalFlipFlop: false,
		verticalFlipFlop:   false,
		offset:             0,
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
	return gr
}

// ColumnInitialize reinitializes the left border sample to its main flip-flop state.
func (b *Borders) ColumnInitialize() {
	b.samples[BorderTypeLeft] = b.horizontalFlipFlop
}

// Column38Update updates the BorderTypeCenter state based on column selection and vertical flip-flop conditions.
func (b *Borders) Column38Update() {
	if !b.core.columnSel {
		b.UpdateVerticalFlipFlop()
		if !b.verticalFlipFlop {
			b.horizontalFlipFlop = false
		}
	}
	b.samples[BorderTypeCenter] = b.horizontalFlipFlop
}

// Column40Update adjusts the state of the mid-left border based on the column selector and vertical flip-flop status.
func (b *Borders) Column40Update() {
	if b.core.columnSel {
		b.UpdateVerticalFlipFlop()
		if !b.verticalFlipFlop {
			b.horizontalFlipFlop = false
		}
	}
	b.samples[BorderTypeMidLeft] = b.horizontalFlipFlop
}

// Column38Apply sets the horizontalFlipFlop to true if columnSel is false and updates the BorderTypeMidRight sample accordingly.
func (b *Borders) Column38Apply() {
	if !b.core.columnSel {
		b.horizontalFlipFlop = true
	}
	b.samples[BorderTypeMidRight] = b.horizontalFlipFlop
}

// Column40Apply sets the 40-column mode by updating the `horizontalFlipFlop` if `columnSel` is active and adjusts the right border sample.
func (b *Borders) Column40Apply() {
	if b.core.columnSel {
		b.horizontalFlipFlop = true
	}
	b.samples[BorderTypeRight] = b.horizontalFlipFlop
}

// SetOffset updates the offset value for the Borders instance with the given parameter.
func (b *Borders) SetOffset(offset int) {
	b.offset = offset
}

// AcquireColor updates the color at the specified index in the border's color array using the core's current color configuration.
func (b *Borders) AcquireColor(idx uint8) {
	b.colors[idx] = _colors[b.core.ec]
}

// UpdateVerticalFlipFlop updates the vertical border flip-flop based on the raster Y coordinate and border comparison values.
func (b *Borders) UpdateVerticalFlipFlop() {
	//3.9. The border unit
	if b.core.dyBottom == b.core.rasterY {
		//2. If the Y coordinate reaches the bottom comparison value in cycle 63, the vertical border flip flop is set.
		b.verticalFlipFlop = true
	} else if b.core.dyTop == b.core.rasterY && b.core.den {
		//3. If the Y coordinate reaches the top comparison value in cycle 63 and the DEN bit in register $d011 is set, the vertical border flip flop is reset.
		b.verticalFlipFlop = false
	}
}

// GetVerticalFlipFlop returns the current state of the vertical border flip-flop, indicating its activation status.
func (b *Borders) GetVerticalFlipFlop() bool {
	return b.verticalFlipFlop
}

// Draw renders the border regions based on current configuration, samples, and colors within the specified display buffer.
func (b *Borders) Draw() {
	if b.samples[BorderTypeLeft] {
		for _, v := range b.left {
			offset := v * borderWidth
			b.setMulti8(b.offset+offset, b.colors[v])
		}
	}
	if b.samples[BorderTypeMidLeft] {
		offset := b.midLeft * borderWidth
		b.setMulti8(b.offset+offset, b.colors[b.midLeft])
	}
	if b.samples[BorderTypeCenter] {
		for _, v := range b.center {
			offset := v * borderWidth
			b.setMulti8(b.offset+offset, b.colors[v])
		}
	}
	if b.samples[BorderTypeMidRight] {
		offset := b.midRight * borderWidth
		b.setMulti8(b.offset+offset, b.colors[b.midRight])
	}
	if b.samples[BorderTypeRight] {
		for _, v := range b.right {
			offset := v * borderWidth
			b.setMulti8(b.offset+offset, b.colors[v])
		}
	}
}
