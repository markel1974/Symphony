package mos6569

type BorderType int

const (
	BorderTypeLeft     = BorderType(0)
	BorderTypeMidLeft  = BorderType(1)
	BorderTypeCenter   = BorderType(2)
	BorderTypeMidRight = BorderType(3)
	BorderTypeRight    = BorderType(4)
	BorderTypeLast     = BorderType(5)
)

type Borders struct {
	core             *Core
	db               IDisplayBuffer
	mainFlipFlop     bool
	verticalFlipFlop bool
	samples          []bool
	colors           []uint8
	offset           int
}

func NewBorder(core *Core, db IDisplayBuffer) *Borders {
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

func (b *Borders) SetOffset(offset int) {
	b.offset = offset
}

func (b *Borders) AcquireColor(idx uint8) {
	//if b.mainFlipFlop {
	b.colors[idx] = _colors[b.core.ec]
	//}
}

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

func (b *Borders) EnableColumn40() {
	if b.core.columnSel {
		b.mainFlipFlop = true
	}
	b.samples[BorderTypeRight] = b.mainFlipFlop
}

func (b *Borders) EnableColumn38() {
	if !b.core.columnSel {
		b.mainFlipFlop = true
	}
	b.samples[BorderTypeMidRight] = b.mainFlipFlop
}

func (b *Borders) UpdateColumn40() {
	if b.core.columnSel {
		b.UpdateVerticalFlipFlop()
		if !b.verticalFlipFlop {
			b.mainFlipFlop = false
		}
	}
	b.samples[BorderTypeMidLeft] = b.mainFlipFlop
}

func (b *Borders) UpdateColumn38() {
	if !b.core.columnSel {
		b.UpdateVerticalFlipFlop()
		if !b.verticalFlipFlop {
			b.mainFlipFlop = false
		}
	}
	b.samples[BorderTypeCenter] = b.mainFlipFlop
}

func (b *Borders) GetVerticalFlipFlop() bool {
	return b.verticalFlipFlop
}

func (b *Borders) Reset() {
	b.samples[BorderTypeLeft] = b.mainFlipFlop
}

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
