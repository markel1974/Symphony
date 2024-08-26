package banks

import "fmt"

type Observer struct {
	b *Banks
}

func NewObserver(b *Banks) *Observer {
	return &Observer{b: b}
}

func (o *Observer) getBasicText() (uint16, uint16) {
	start := uint16(o.b.ram[0x2b]) | (uint16(o.b.ram[0x2c]) << 8)
	end := uint16(o.b.ram[0x2d]) | (uint16(o.b.ram[0x2e]) << 8)
	return start, end
}

func (o *Observer) setBasicText(start uint16, end uint16) {
	s1 := uint8(start) & 0xff
	o.b.ram[0xac] = s1
	o.b.ram[0x2b] = s1

	s2 := uint8(start >> 8)
	o.b.ram[0xad] = s2
	o.b.ram[0x2c] = s2

	e1 := uint8(end & 0xff)
	o.b.ram[0xae] = e1
	o.b.ram[0x31] = e1
	o.b.ram[0x2f] = e1
	o.b.ram[0x2d] = e1

	e2 := uint8(end >> 8)
	o.b.ram[0xaf] = e2
	o.b.ram[0x32] = e2
	o.b.ram[0x30] = e2
	o.b.ram[0x2e] = e2
}

func (o *Observer) Inject(autostartBasicLoad bool, startAddr uint16, data []byte) error {
	size := uint16(len(data))
	start, _ := o.getBasicText()
	// load to basic start if requested
	if autostartBasicLoad {
		startAddr = start
	}
	for x := uint16(0); x < size; x++ {
		o.b.ram[startAddr+x] = data[x]
	}
	// simulate a basic load
	end := startAddr + size
	o.setBasicText(start, end)
	return nil
}

// used by autostart to locate and "read" kernal output on the current screen
// this function should return whatever the kernal currently uses, regardless
// what is currently visible/active in the UI
// static CHECKYESNO check2(

func (o *Observer) GetCursorParameter(lineOffset int) {
	//uint16_t *screen_addr, uint8_t *cursor_column, uint8_t *line_length, int *blinking
	// CAUTION: this function can be called at any time when the emulation (KERNAL)
	// is in the middle of a screen update. we must make sure that all
	// values are being looked up in an "atomic" way so we don't use a low-
	// and high- byte from before and after an update, leading to invalid values

	// Physical Screen Line Length
	const lineLength = 40
	screenBase := (int(o.b.ram[0xd1]) + (int(o.b.ram[0xd2]) * 256)) & ^0x3ff // the upper bits will not change
	//blinking := b.ram[0xcc] == 0 //? 0 : 1;
	// Current Screen Line Address
	screenAddr := screenBase + (int(o.b.ram[0xd6]) * lineLength)
	// Cursor Column on Current Line
	cursorColumn := int(o.b.ram[0xd3])
	for cursorColumn >= lineLength {
		cursorColumn -= lineLength
		screenAddr += lineLength
	}
	addr := screenAddr
	addr += lineLength * lineOffset
	for x := 0; x < lineLength; x++ {
		a := uint16(addr+x) & 0xffff
		v := o.b.ram[a]
		if v != 32 && v != 0 && v != 255 && v != 160 {
			fmt.Println(v)
		}
	}
	//fmt.Println(blinking, screenAddr)
}
