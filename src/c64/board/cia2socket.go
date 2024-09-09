package board

type CIA2Socket struct {
	board  *Board
	intrId uint32
}

func NewCIA2Socket() *CIA2Socket {
	c := &CIA2Socket{}
	return c
}

func (w *CIA2Socket) Setup(board *Board, intrId uint32) {
	w.board = board
	w.intrId = intrId
}

func (w *CIA2Socket) Reset() {

}

func (w *CIA2Socket) ReadPortA(prA uint8, ddrA uint8, _ uint8, _ uint8) uint8 {
	data := w.board.iec.CpuRead()
	ret := ((prA | (^ddrA)) & 0x3f) | data
	return ret
}

func (w *CIA2Socket) ReadPortB(_ uint8, _ uint8, prB uint8, ddrB uint8) uint8 {
	ret := prB | (^ddrB)
	return ret
}

func (w *CIA2Socket) WritePortA(prA uint8, ddrA uint8, _ uint8, _ uint8) {
	w.updateVA(prA, ddrA)
	w.board.iec.CpuWrite(prA)
}

func (w *CIA2Socket) WritePortB(_ uint8, _ uint8, _ uint8, _ uint8) {
}

func (w *CIA2Socket) WriteDdrA(prA uint8, ddrA uint8, _ uint8, _ uint8) {
	w.updateVA(prA, ddrA)
}

func (w *CIA2Socket) WriteDdrB(_ uint8, _ uint8, _ uint8, _ uint8) {
}

func (w *CIA2Socket) updateVA(prA uint8, ddrA uint8) {
	//Bit 0..1: Select the position of the VIC-memory
	//Bit 2: RS-232: TXD Output, userPort: Data PA 2 (pin M)
	//Bit 3..5: serial bus Output (0=High/Inactive, 1=Low/Active)
	//Bit 6..7: serial bus Input (0=Low/Active, 1=High/Inactive)

	//%00, 0: Bank 3: $C000-$FFFF, 49152-65535
	//%01, 1: Bank 2: $8000-$BFFF, 32768-49151
	//%10, 2: Bank 1: $4000-$7FFF, 16384-32767
	//%11, 3: Bank 0: $0000-$3FFF, 0-16383 (standard)
	va := (^(prA | (^ddrA))) & 3
	w.board.vic.ChangedVA(va)
}

func (w *CIA2Socket) IRQTrigger() {
	w.board.nmiTriggerSlot()
}

func (w *CIA2Socket) IRQClear() {
	w.board.nmiClearSlot()
}
