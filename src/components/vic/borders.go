package vic

type BorderType int

const (
	BorderTypeLeft     = BorderType(0)
	BorderTypeMidLeft  = BorderType(1)
	BorderTypeCenter   = BorderType(2)
	BorderTypeMidRight = BorderType(3)
	BorderTypeRight    = BorderType(4)
)

type Borders struct {
	core             *Core
	db               IDisplayBuffer
	mainFlipFlop     bool
	verticalFlipFlop bool
	samples          []bool
	colors           []uint8
	firstCycle       uint8
}

func NewBorder(core *Core, db IDisplayBuffer, firstCycle uint8) *Borders {
	gr := &Borders{
		db:               db,
		core:             core,
		samples:          make([]bool, 5),
		colors:           make([]uint8, 0xff),
		mainFlipFlop:     false,
		verticalFlipFlop: false,
		firstCycle:       firstCycle,
	}
	return gr
}

func (b *Borders) Sample(cycle uint8) {
	if b.mainFlipFlop {
		b.colors[cycle-b.firstCycle] = b.core.ecColor
	}
}

func (b *Borders) Enable() {
	b.mainFlipFlop = true
}

func (b *Borders) SetSample(t BorderType) {
	b.samples[t] = b.mainFlipFlop
}

func (b *Borders) GetVerticalFlipFlop() bool {
	return b.verticalFlipFlop
}

func (b *Borders) UpdateVerticalFlipFlop() {
	if b.core.rasterY == b.core.dyStop {
		b.verticalFlipFlop = true
	} else if (b.core.cr1&0x10) != 0 && b.core.rasterY == b.core.dyStart {
		b.verticalFlipFlop = false
	}
}

func (b *Borders) Update() {
	if b.core.rasterY == b.core.dyStop {
		b.verticalFlipFlop = true
		return
	}
	if (b.core.cr1 & 0x10) != 0 {
		if b.core.rasterY == b.core.dyStart {
			b.verticalFlipFlop = false
			b.mainFlipFlop = false
		} else if !b.verticalFlipFlop {
			b.mainFlipFlop = false
		}
	} else if !b.verticalFlipFlop {
		b.mainFlipFlop = false
	}
}

func (b *Borders) Draw(lineStart int) {
	const bSize = 8
	const border0Start = 0
	const border1End = 4
	const border1Offset = border1End * bSize
	const border2Start = border1End + 1
	const border2End = 43
	const border3Offset = border2End * bSize
	const border4Start = border2End + 1
	const border4End = DisplayXDiv8

	if b.samples[BorderTypeLeft] {
		for idx, offset := border0Start, border0Start*bSize; idx < border1End; idx, offset = idx+1, offset+bSize {
			b.db.SetMulti8(lineStart+offset, b.colors[idx])
		}
	}
	if b.samples[BorderTypeMidLeft] {
		b.db.SetMulti8(lineStart+(border1Offset), b.colors[border1End])
	}
	if b.samples[BorderTypeCenter] {
		for idx, offset := border2Start, border2Start*bSize; idx < border2End; idx, offset = idx+1, offset+bSize {
			b.db.SetMulti8(lineStart+offset, b.colors[idx])
		}
	}
	if b.samples[BorderTypeMidRight] {
		b.db.SetMulti8(lineStart+(border3Offset), b.colors[border2End])
	}
	if b.samples[BorderTypeRight] {
		for idx, offset := border4Start, border4Start*bSize; idx < border4End; idx, offset = idx+1, offset+bSize {
			b.db.SetMulti8(lineStart+offset, b.colors[idx])
		}
	}
}
