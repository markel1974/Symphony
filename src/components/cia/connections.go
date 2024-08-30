package cia

/*
 * Notes:
 * ------
 *
 *  - The Emulate() function is called for every emulated Phi2 clock cycle.
 *    It counts down the timers and triggers interrupts if necessary.
 *  - The TOD clocks are counted by TODUpdate() during the VBlank, so the input frequency is 50Hz
 *  - The fields keyMatrix and revMatrix contain one bit for each
 *    key on the C64 keyboard (0: key pressed, 1: key released).
 *    keyMatrix is used for normal keyboard polling (PRA->PRB),
 *    revMatrix for reversed polling (PRB->PRA).
 *
 * Incompatibilities:
 * ------------------
 *  - The SDR interrupt is faked
 *  - Some small incompatibilities with the timers
 */

const (
	IRQUnderflowTimerA = 0x1
	IRQUnderflowTimerB = 0x2
	IRQTODAlarmEqual   = 0x4
	IRQSDRFullOrEmpty  = 0x8
	IRQFlagPin         = 0x10
	IRQOccurred        = 0x80
)

const intrCia1Id = 4

type IBus interface {
	CpuRead() uint8
	CpuWrite(uint8)
}

type IPort interface {
	ReadPortA(prA uint8, ddrA uint8, prb uint8, ddrB uint8) uint8
	ReadPortB(prA uint8, ddrA uint8, prb uint8, ddrB uint8) uint8
	ReadDdrA(prA uint8, ddrA uint8, prb uint8, ddrB uint8) uint8
	ReadDdrB(prA uint8, ddrA uint8, prb uint8, ddrB uint8) uint8
	WritePortA(prA uint8, ddrA uint8, prb uint8, ddrB uint8)
	WritePortB(prA uint8, ddrA uint8, prb uint8, ddrB uint8)
	WriteDdrA(prA uint8, ddrA uint8, prb uint8, ddrB uint8)
	WriteDdrB(prA uint8, ddrA uint8, prb uint8, ddrB uint8)
}
