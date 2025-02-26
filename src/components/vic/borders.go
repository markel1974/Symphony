package mos6569

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
	core             *VIC
	db               IDisplayBuffer
	mainFlipFlop     bool
	verticalFlipFlop bool
	samples          []bool
	colors           []uint8
	offset           int
}

// NewBorder creates and initializes a new Borders instance with the provided VIC core and display buffer dependencies.
func NewBorder(core *VIC, db IDisplayBuffer) *Borders {
	gr := &Borders{
		db:               db,
		core:             core,
		samples:          make([]bool, BorderTypeLast),
		colors:           make([]uint8, 0xff),
		mainFlipFlop:     false,
		verticalFlipFlop: false,
		offset:           0,
	}
	return gr
}

// SetOffset updates the offset value for the Borders instance with the given parameter.
func (b *Borders) SetOffset(offset int) {
	b.offset = offset
}

// AcquireColor updates the color at the specified index in the border's color array using the core's current color configuration.
func (b *Borders) AcquireColor(idx uint8) {
	//if b.mainFlipFlop {
	b.colors[idx] = _colors[b.core.ec]
	//}
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

// EnableColumn40 sets the 40-column mode by updating the `mainFlipFlop` if `columnSel` is active and adjusts the right border sample.
func (b *Borders) EnableColumn40() {
	if b.core.columnSel {
		b.mainFlipFlop = true
	}
	b.samples[BorderTypeRight] = b.mainFlipFlop
}

// EnableColumn38 sets the mainFlipFlop to true if columnSel is false and updates the BorderTypeMidRight sample accordingly.
func (b *Borders) EnableColumn38() {
	if !b.core.columnSel {
		b.mainFlipFlop = true
	}
	b.samples[BorderTypeMidRight] = b.mainFlipFlop
}

// UpdateColumn40 adjusts the state of the mid-left border based on the column selector and vertical flip-flop status.
func (b *Borders) UpdateColumn40() {
	if b.core.columnSel {
		b.UpdateVerticalFlipFlop()
		if !b.verticalFlipFlop {
			b.mainFlipFlop = false
		}
	}
	b.samples[BorderTypeMidLeft] = b.mainFlipFlop
}

// UpdateColumn38 updates the BorderTypeCenter state based on column selection and vertical flip-flop conditions.
func (b *Borders) UpdateColumn38() {
	if !b.core.columnSel {
		b.UpdateVerticalFlipFlop()
		if !b.verticalFlipFlop {
			b.mainFlipFlop = false
		}
	}
	b.samples[BorderTypeCenter] = b.mainFlipFlop
}

// GetVerticalFlipFlop returns the current state of the vertical border flip-flop, indicating its activation status.
func (b *Borders) GetVerticalFlipFlop() bool {
	return b.verticalFlipFlop
}

// Reset reinitializes the left border sample to its main flip-flop state.
func (b *Borders) Reset() {
	b.samples[BorderTypeLeft] = b.mainFlipFlop
}

// Draw renders the border regions based on current configuration, samples, and colors within the specified display buffer.
func (b *Borders) Draw( /*lineStart int*/ ) {
	const bSize = 8
	const border0Start = 0
	const border1End = 4
	const border1Offset = border1End * bSize
	const border2Start = border1End + 1
	const border2End = 43
	const border3Offset = border2End * bSize
	const border4Start = border2End + 1
	const border4End = DisplayXDiv8

	const border0StartSize = border0Start * bSize
	const border2StartSize = border2Start * bSize
	const border4StartSize = border4Start * bSize

	if b.samples[BorderTypeLeft] {
		for idx, offset := border0Start, border0StartSize; idx < border1End; idx, offset = idx+1, offset+bSize {
			b.db.SetMulti8(b.offset+offset, b.colors[idx])
		}
	}
	if b.samples[BorderTypeMidLeft] {
		b.db.SetMulti8(b.offset+(border1Offset), b.colors[border1End])
	}
	if b.samples[BorderTypeCenter] {
		for idx, offset := border2Start, border2StartSize; idx < border2End; idx, offset = idx+1, offset+bSize {
			b.db.SetMulti8(b.offset+offset, b.colors[idx])
		}
	}
	if b.samples[BorderTypeMidRight] {
		b.db.SetMulti8(b.offset+(border3Offset), b.colors[border2End])
	}
	if b.samples[BorderTypeRight] {
		for idx, offset := border4Start, border4StartSize; idx < border4End; idx, offset = idx+1, offset+bSize {
			b.db.SetMulti8(b.offset+offset, b.colors[idx])
		}
	}
}
