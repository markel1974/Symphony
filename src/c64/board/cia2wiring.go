package board

import (
	"github.com/markel1974/c64emu/src/c64/iec"
	"github.com/markel1974/c64emu/src/signals"
)

type CIA2Wiring struct {
	bus             *iec.IEC
	signalChangedVA *signals.SignalByte
}

func NewCIA2Wiring() *CIA2Wiring {
	c := &CIA2Wiring{
		signalChangedVA: signals.NewSignalByte(),
	}
	return c
}

func (w *CIA2Wiring) Setup(bus *iec.IEC, fn func(uint8)) {
	w.bus = bus
	w.signalChangedVA.Bind(fn)
}

func (w *CIA2Wiring) Reset() {

}

func (w *CIA2Wiring) ReadPortA(prA uint8, ddrA uint8, _ uint8, _ uint8) uint8 {
	data := w.bus.CpuRead()
	ret := ((prA | (^ddrA)) & 0x3f) | data
	return ret
}

func (w *CIA2Wiring) ReadPortB(_ uint8, _ uint8, prB uint8, ddrB uint8) uint8 {
	ret := prB | (^ddrB)
	return ret
}

func (w *CIA2Wiring) WritePortA(prA uint8, ddrA uint8, _ uint8, _ uint8) {
	w.updateVA(prA, ddrA)
	w.bus.CpuWrite(prA)
}

func (w *CIA2Wiring) WritePortB(_ uint8, _ uint8, _ uint8, _ uint8) {
}

func (w *CIA2Wiring) WriteDdrA(prA uint8, ddrA uint8, _ uint8, _ uint8) {
	w.updateVA(prA, ddrA)
}

func (w *CIA2Wiring) WriteDdrB(_ uint8, _ uint8, _ uint8, _ uint8) {
}

func (w *CIA2Wiring) updateVA(prA uint8, ddrA uint8) {
	//Bit 0..1: Select the position of the VIC-memory
	//Bit 2: RS-232: TXD Output, userPort: Data PA 2 (pin M)
	//Bit 3..5: serial bus Output (0=High/Inactive, 1=Low/Active)
	//Bit 6..7: serial bus Input (0=Low/Active, 1=High/Inactive)

	//%00, 0: Bank 3: $C000-$FFFF, 49152-65535
	//%01, 1: Bank 2: $8000-$BFFF, 32768-49151
	//%10, 2: Bank 1: $4000-$7FFF, 16384-32767
	//%11, 3: Bank 0: $0000-$3FFF, 0-16383 (standard)
	va := (^(prA | (^ddrA))) & 3
	w.signalChangedVA.Emit(va)
}
