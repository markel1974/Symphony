package cpu

import (
	"github.com/markel1974/c64emu/src/board/iboard"
	"github.com/markel1974/c64emu/src/preferences"
)

// Interrupts
const (
	IntVic = 0x1
	IntCia = 0x2
	IntNmi = 0x4
	IntRst = 0x8
)

type interrupt uint32

func (i *interrupt) Clear() {
	*i = 0
}

func (i *interrupt) BitSet(n uint32) {
	*i = *i | (1 << n)
}

func (i *interrupt) BitClear(n uint32) {
	*i = *i & ^(1 << n)
}

func (i *interrupt) BitCheck(n uint32) bool {
	v := (*i >> n) & 1
	return v != 0
}

type InterruptHandler struct {
	board         iboard.IBoard
	prefs         *preferences.Prefs
	intr          interrupt
	firstIrqCycle uint64
	firstNMICycle uint64
}

func New() *InterruptHandler {
	return &InterruptHandler{
		board:         nil,
		prefs:         nil,
		firstIrqCycle: 0,
		firstNMICycle: 0,
	}
}

func (i *InterruptHandler) Setup(board iboard.IBoard) {
	i.board = board
}

func (i *InterruptHandler) Reset() {
	i.intr = 0
}

func (i *InterruptHandler) HasInterrupt() bool {
	return i.intr != 0
}

func (i *InterruptHandler) AsyncReset() {
	i.intr.BitSet(IntRst)
}

func (i *InterruptHandler) HasReset() bool {
	return i.intr.BitCheck(IntRst)
}

func (i *InterruptHandler) TriggerVICIRQ() {
	vic := i.intr.BitCheck(IntVic)
	cia := i.intr.BitCheck(IntCia)
	if !(vic || cia) {
		i.firstIrqCycle = i.board.Cycle()
	}
	i.intr.BitSet(IntVic)
}

func (i *InterruptHandler) ClearVICIRQ() {
	i.intr.BitClear(IntVic)
}

func (i *InterruptHandler) HasVIC() bool {
	return i.intr.BitCheck(IntVic)
}

func (i *InterruptHandler) TriggerCIAIRQ() {
	vic := i.intr.BitCheck(IntVic)
	cia := i.intr.BitCheck(IntCia)
	if !(vic || cia) {
		i.firstIrqCycle = i.board.Cycle()
	}
	i.intr.BitSet(IntCia)
}

func (i *InterruptHandler) ClearCIAIRQ() {
	i.intr.BitClear(IntCia)
}

func (i *InterruptHandler) HasCIA() bool {
	return i.intr.BitCheck(IntCia)
}

func (i *InterruptHandler) AsyncNMI() {
	i.intr.BitSet(IntNmi)
}

func (i *InterruptHandler) TriggerNMI() {
	if !i.intr.BitCheck(IntNmi) {
		i.firstNMICycle = i.board.Cycle()
	}
	i.intr.BitSet(IntNmi)
}

func (i *InterruptHandler) ClearNMI() {
	i.intr.BitClear(IntNmi)
}

func (i *InterruptHandler) HasNMI() bool {
	return i.intr.BitCheck(IntNmi)
}

func (i *InterruptHandler) GetNMICycleDistance(delay uint64) uint64 {
	return i.computeDistance(i.firstNMICycle, delay)
}

func (i *InterruptHandler) GetIrqCycleDistance(delay uint64) uint64 {
	return i.computeDistance(i.firstIrqCycle, delay)
}

func (i *InterruptHandler) computeDistance(base uint64, delay uint64) uint64 {
	cycle := i.board.Cycle()
	if base > cycle {
		return 0
	}
	v := cycle - base
	if v < delay {
		return 0
	}
	v -= delay
	return v
}
