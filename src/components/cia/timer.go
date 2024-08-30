package cia

import (
	"fmt"
	"github.com/markel1974/c64emu/src/signals"
)

const (
	timerStop = iota
	timerWaitThenCount
	timerLoadThenStop
	timerLoadThenCount
	timerLoadThenWaitThenCount
	timerCount
	timerCountThenStop
)

//Bit 0:
//0 = Stop timer
//1 = Start timer
//Bit 1:
//0 = Indicates no timer underflow at port B in bit 6.
//1 = Indicates a timer underflow at port B in bit 6.
//Bit 2:
//0 = Through a timer overflow, bit 6 of port B will get high for one cycle,
//1 = Through a timer underflow, bit 6 of port B will be inverted
//Bit 3:
//0 = Timer-restart after underflow (latch will be reloaded),
//1 = Timer stops after underflow.
//Bit 4:
//0 = Not Load latch
//1 = Load latch into the timer once.
//Bit 5:
//0 = Timer counts system cycles,
//1 = Timer counts positive slope at CNT-pin
//Bit 6: Direction of the serial shift register,
//0 = SP-pin is input (read),
//1 = SP-pin is output (write)
//Bit 7: Real Time Clock,
//0 = 60 Hz
//1 = 50 Hz

const (
	crBitStopped                 = 0x1
	crBitSignalUnderflow         = 0x2
	crBitSignalUnderflowInverted = 0x4
	crBitOneShot                 = 0x8
	crBitForceLoad               = 0x10
	crBitTimerCountSystemCycle   = 0x20
	crBitShiftRegDir             = 0x40
	crBitRTC                     = 0x80
)

type Timer struct {
	id                   string
	cr                   uint8
	crLatch              uint8  // New values for cr
	crPending            bool   // New value for crLatch pending
	timer                uint16 // Timer
	timerLatch           uint16 // Timer latch
	timerState           uint8  // Timer states
	countPhi2            bool   // Timer is counting Phi 2
	checkUnderflowTimerX bool   // Timer is counting underflow's of Timer X
	signalTimerUnderflow *signals.Signal
}

func NewTimer(id string) *Timer {
	m := &Timer{
		id:                   id,
		signalTimerUnderflow: signals.NewSignal(),
	}
	m.Reset()
	return m
}

func (m *Timer) Reset() {
	m.cr = 0
	m.crLatch = 0
	m.crPending = false
	m.timer = 0xffff
	m.timerLatch = 1
	m.timerState = timerStop
	m.countPhi2 = false
	m.checkUnderflowTimerX = false
}

func (m *Timer) GetRTC() bool {
	return (m.cr & crBitRTC) != 0
}

func (m *Timer) GetCR() uint8 {
	return m.cr
}

func (m *Timer) GetTimerLow() uint8 {
	return uint8(m.timer)
}

func (m *Timer) GetTimerHigh() uint8 {
	return uint8(m.timer >> 8)
}

func (m *Timer) SetTimerLow(data uint8) {
	m.timerLatch = (m.timerLatch & 0xff00) | uint16(data)
}

func (m *Timer) SetTimerHigh(data uint8) {
	m.timerLatch = (m.timerLatch & 0xff) | (uint16(data) << 8)
	if (m.cr & crBitStopped) == 0 {
		// Reload timer if stopped
		m.timer = m.timerLatch
	}
}

func (m *Timer) SetTimerControl(data uint8, count bool) {
	//m.printTimerControlData(data)
	m.crPending = true
	m.crLatch = data
	if count {
		m.countPhi2 = (data & 0x60) == 0x00
		m.checkUnderflowTimerX = (data & 0x60) == 0x40
	} else {
		m.countPhi2 = (data & 0x20) == 0x00
		m.checkUnderflowTimerX = false
	}
}

func (m *Timer) SignalTimerUnderflowBind(fn func()) {
	m.signalTimerUnderflow.Bind(fn)
}

func (m *Timer) Emulate(underflowTimerX bool) bool {
	underflow := false
	interrupt := false

	// Timer state machine
	switch m.timerState {
	case timerWaitThenCount:
		// fall through
		m.timerState = timerCount
		goto labelIdle
	case timerStop:
		goto labelIdle
	case timerLoadThenStop:
		m.timerState = timerStop
		// Reload timer
		m.timer = m.timerLatch
		goto labelIdle
	case timerLoadThenCount:
		m.timerState = timerCount
		// Reload timer
		m.timer = m.timerLatch
		goto labelIdle
	case timerLoadThenWaitThenCount:
		m.timerState = timerWaitThenCount
		if m.timer == 1 {
			interrupt = true
			goto labelCount
		} else {
			// Reload timer
			m.timer = m.timerLatch
			goto labelIdle
		}
	case timerCount:
		goto labelCount
	case timerCountThenStop:
		m.timerState = timerStop
		goto labelCount
	}

labelCount:
	if !interrupt {
		if m.countPhi2 || (m.checkUnderflowTimerX && underflowTimerX) {
			timer := m.timer
			m.timer--
			if (timer == 0) || (m.timer == 0) {
				if m.timerState != timerStop {
					interrupt = true
				}
				underflow = true
			}
		}
	}

	if interrupt {
		// Reload timer
		m.timer = m.timerLatch
		m.signalTimerUnderflow.Emit()
		if (m.cr & crBitOneShot) != 0 {
			// stop timer
			m.cr &= 0xfe
			m.crLatch &= 0xfe
			// Reload in next cycle
			m.timerState = timerLoadThenStop
		} else {
			// delay one cycle (and reload)
			m.timerState = timerLoadThenCount
		}
		underflow = true
	}

	// Delayed write to CR?
labelIdle:
	if m.crPending {
		switch m.timerState {
		case timerStop, timerLoadThenStop:
			// Timer started, wasn't running
			if (m.crLatch & crBitStopped) != 0 {
				if (m.crLatch & crBitForceLoad) != 0 {
					m.timerState = timerLoadThenWaitThenCount
				} else {
					m.timerState = timerWaitThenCount
				}
			} else {
				// Timer stopped, was already stopped
				if (m.crLatch & crBitForceLoad) != 0 {
					m.timerState = timerLoadThenStop
				}
			}
		case timerCount:
			if (m.crLatch & crBitStopped) != 0 {
				// Timer started, was already running
				if (m.crLatch & crBitForceLoad) != 0 {
					m.timerState = timerLoadThenWaitThenCount
				}
			} else {
				// Timer stopped, was running
				if (m.crLatch & crBitForceLoad) != 0 {
					m.timerState = timerLoadThenStop
				} else {
					m.timerState = timerCountThenStop
				}
			}
		case timerLoadThenCount, timerWaitThenCount:
			if (m.crLatch & crBitStopped) != 0 {
				if (m.crLatch & crBitOneShot) != 0 {
					// One-shot, stop timer
					m.crLatch &= 0xfe
					m.timerState = timerStop
				} else if (m.crLatch & crBitForceLoad) != 0 {
					// No One-shot, force load
					m.timerState = timerLoadThenWaitThenCount
				}
			} else {
				m.timerState = timerStop
			}
		default:
			fmt.Println("TIMER - UNDEFINED", m.timerState)
		}
		//no force load set
		m.cr = m.crLatch & 0xef
		m.crPending = false
	}
	return underflow
}

func (m *Timer) printTimerControlData(data uint8) {
	fmt.Printf("\n")
	fmt.Printf("%s Timer Control -> crBitStopped: %v\n", m.id, data&crBitStopped != 0)
	fmt.Printf("%s Timer Control -> crBitSignalUnderflow: %v\n", m.id, data&crBitSignalUnderflow != 0)
	fmt.Printf("%s Timer Control -> crBitSignalUnderflowInverted: %v\n", m.id, data&crBitSignalUnderflowInverted != 0)
	fmt.Printf("%s Timer Control -> crBitOneShot: %v\n", m.id, data&crBitOneShot != 0)
	fmt.Printf("%s Timer Control -> crBitForceLoad: %v\n", m.id, data&crBitForceLoad != 0)
	fmt.Printf("%s Timer Control -> crBitTimerCountSystemCycle: %v\n", m.id, data&crBitTimerCountSystemCycle != 0)
	fmt.Printf("%s Timer Control -> crBitShiftRegDir: %v\n", m.id, data&crBitShiftRegDir != 0)
	fmt.Printf("%s Timer Control -> crBitRTC: %v\n", m.id, data&crBitRTC != 0)
	fmt.Printf("\n")
}
