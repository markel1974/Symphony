package mos6510_rev1

import (
	"github.com/markel1974/c64emu/src/references"
	"reflect"
	"testing"
)

// MockPIC is a mock implementation of the PIC interface for testing purposes.
type MockPIC struct {
	verifyIrqResult uint8
}

func (m *MockPIC) Setup(quartz references.IQuartz) error { return nil }
func (m *MockPIC) ClearIRQ(u uint32)                     {}
func (m *MockPIC) TriggerIRQ(u uint32)                   {}
func (m *MockPIC) TriggerReset()                         {}
func (m *MockPIC) TriggerNMI()                           {}
func (m *MockPIC) IRQTriggerBind(fn func(uint32))        {}
func (m *MockPIC) IRQClearBind(fn func(uint32))          {}
func (m *MockPIC) VerifyIrq(iFlag uint8, opFlag uint8) uint8 {
	return m.verifyIrqResult
}
func (m *MockPIC) Reset()       {}
func (m *MockPIC) ClearNMI()    {}
func (m *MockPIC) HasNMI() bool { return false }
func (m *MockPIC) HasIRQ() bool { return false }

// MockBank is a mock implementation of the Bank interface for testing purposes.
type MockBank struct {
	readResult byte
}

func (m *MockBank) Write(u uint16, u2 uint8) {
}

func (m *MockBank) Read(addr uint16) byte {
	return m.readResult
}

func TestInstOpINI(t *testing.T) {
	tests := []struct {
		name       string
		cpu        *CPU
		verifyIrq  int
		expPC      uint16
		expStop    bool
		expNext    func(*CPU)
		expBreaker bool
		pic        references.IMos6510Pic
		banks      references.IMos6510Banks
	}{
		{
			name: "CPU_RDY_Low",
			cpu: &CPU{
				rdyLow: true,
			},
			verifyIrq:  0, // Not applicable since rdyLow is true and we exit early
			expPC:      0,
			expStop:    true,
			expBreaker: false,
		},
		{
			name: "IRQ_Verify_Returns_1_Reset",
			cpu: &CPU{
				rdyLow:     false,
				irqBreaker: false,
			},
			verifyIrq:  1,
			expPC:      0,
			expStop:    false,
			expBreaker: false,
			pic:        &MockPIC{verifyIrqResult: 1},
			banks:      &MockBank{readResult: 0},
		},
		{
			name: "IRQ_Verify_Returns_2_NMI",
			cpu: &CPU{
				rdyLow:     false,
				irqBreaker: false,
			},
			verifyIrq:  2,
			expPC:      1,
			expStop:    false,
			expNext:    instOpNMI,
			expBreaker: false,
			pic:        &MockPIC{verifyIrqResult: 2},
			banks:      &MockBank{readResult: 0},
		},
		{
			name: "IRQ_Verify_Returns_3_IRQ",
			cpu: &CPU{
				rdyLow:     false,
				irqBreaker: false,
			},
			verifyIrq:  3,
			expPC:      1,
			expStop:    false,
			expNext:    instOpIRQ,
			expBreaker: false,
			pic:        &MockPIC{verifyIrqResult: 3},
			banks:      &MockBank{readResult: 0},
		},
		{
			name: "Normal_Execution",
			cpu: &CPU{
				rdyLow:     false,
				irqBreaker: true,
				pc:         0x1000,
			},
			verifyIrq:  0, // Not applicable since IRQ is not triggered
			expPC:      0x1002,
			expStop:    false,
			expNext:    _modeTable[0x20],
			expBreaker: false,
			pic:        &MockPIC{verifyIrqResult: 0},
			banks:      &MockBank{readResult: 0x20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.cpu.Setup(&mockSocket{banks: tt.banks, pic: tt.pic})
			// Call instOpINI
			instOpINI(tt.cpu)
			instOpINI(tt.cpu)

			tt.cpu.SetRDYLow(true)

			// Validate outcomes
			if tt.cpu.pc != tt.expPC {
				t.Errorf("Expected PC to be %04X, got %04X", tt.expPC, tt.cpu.pc)
			}
			if tt.cpu.stop != tt.expStop {
				t.Errorf("Expected stop to be %v, got %v", tt.expStop, tt.cpu.stop)
			}

			if tt.expNext != nil {
				sf1 := reflect.ValueOf(tt.cpu.next)
				sf2 := reflect.ValueOf(tt.expNext)
				if sf1.Pointer() != sf2.Pointer() {
					t.Errorf("Expected next function to be %v, got %v", sf2.Pointer(), sf1.Pointer())
				}
			}
			if tt.cpu.irqBreaker != tt.expBreaker {
				t.Errorf("Expected irqBreaker to be %v, got %v", tt.expBreaker, tt.cpu.irqBreaker)
			}
		})
	}
}
