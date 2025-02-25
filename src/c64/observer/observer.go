package observer

import (
	"fmt"
	"log"
)

const lineLength = 40

type Observer struct {
	b IAdapter
}

func NewObserver(b IAdapter) *Observer {
	return &Observer{b: b}
}

func (o *Observer) getBasicText() (uint16, uint16) {
	b := o.b.Read(0x2b)
	c := o.b.Read(0x2c)
	d := o.b.Read(0x2d)
	e := o.b.Read(0x2e)
	start := uint16(b) | (uint16(c) << 8)
	end := uint16(d) | (uint16(e) << 8)
	return start, end
}

func (o *Observer) setBasicText(start uint16, end uint16) {
	s1 := uint8(start) & 0xff
	o.b.Write(0xac, s1)
	o.b.Write(0x2b, s1)

	s2 := uint8(start >> 8)
	o.b.Write(0xad, s2)
	o.b.Write(0x2c, s2)

	e1 := uint8(end & 0xff)
	o.b.Write(0xae, e1)
	o.b.Write(0x31, e1)
	o.b.Write(0x2f, e1)
	o.b.Write(0x2d, e1)

	e2 := uint8(end >> 8)
	o.b.Write(0xaf, e2)
	o.b.Write(0x32, e2)
	o.b.Write(0x30, e2)
	o.b.Write(0x2e, e2)
}

func (o *Observer) Inject(autostartBasicLoad bool, startAddr uint16, data []byte) {
	size := uint16(len(data))
	start, _ := o.getBasicText()
	// load to basic start if requested
	if autostartBasicLoad {
		startAddr = start
	}
	for x := uint16(0); x < size; x++ {
		o.b.Write(startAddr+x, data[x])
	}
	// simulate a basic load
	end := startAddr + size
	o.setBasicText(start, end)
}

func (o *Observer) GetCursorParameter(lineOffset int) {
	if lineOffset < 0 {
		log.Printf("Invalid lineOffset: %d\n", lineOffset)
		return
	}
	// Physical Screen Line Length
	screenBase := (int(o.b.Read(0xd1)) + (int(o.b.Read(0xd2)) * 256)) & ^0x3ff // the upper bits will not change
	//blinking := b.ram[0xcc] == 0 //? 0 : 1;
	// Current Screen Line Address
	screenAddr := screenBase + (int(o.b.Read(0xd6)) * lineLength)
	// Cursor Column on Current Line
	cursorColumn := int(o.b.Read(0xd3))
	for cursorColumn >= lineLength {
		cursorColumn -= lineLength
		screenAddr += lineLength
	}
	addr := screenAddr
	addr += lineLength * lineOffset
	for x := 0; x < lineLength; x++ {
		a := uint16(addr+x) & 0xffff
		v := o.b.Read(a)
		if v != 32 && v != 0 && v != 255 && v != 160 {
			fmt.Println(v)
		}
	}
	//fmt.Println(blinking, screenAddr)
}
