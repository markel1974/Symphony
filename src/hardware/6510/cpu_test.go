package mos6510

import (
	"github.com/markel1974/c64emu/src/references"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockBanks struct {
	data map[uint16]uint8
}

func (m *mockBanks) Read(addr uint16) uint8 {
	return m.data[addr]
}

func (m *mockBanks) Write(addr uint16, value uint8) {
	m.data[addr] = value
}

type mockPic struct{}

func (m *mockPic) Setup(quartz references.IQuartz) error     { return nil }
func (m *mockPic) IRQTriggerBind(fn func(uint32))            {}
func (m *mockPic) IRQClearBind(fn func(uint32))              {}
func (m *mockPic) Reset()                                    {}
func (m *mockPic) ClearNMI()                                 {}
func (m *mockPic) ClearIRQ(u uint32)                         {}
func (m *mockPic) HasNMI() bool                              { return false }
func (m *mockPic) HasIRQ() bool                              { return false }
func (m *mockPic) VerifyIrq(iFlag uint8, opFlag uint8) uint8 { return 0 }
func (m *mockPic) TriggerIRQ(u uint32)                       {}
func (m *mockPic) TriggerReset()                             {}
func (m *mockPic) TriggerNMI()                               {}

type mockSocket struct {
	banks references.I6510Banks
	pic   references.IPIC6510
}

func (m *mockSocket) GetBanks() references.I6510Banks {
	return m.banks
}

func (m *mockSocket) GetPic() references.IPIC6510 {
	return m.pic
}

func TestCPU_Setup(t *testing.T) {
	banks := &mockBanks{data: make(map[uint16]uint8)}
	pic := &mockPic{}
	socket := &mockSocket{banks: banks, pic: pic}

	cpu := NewCPU(nil, nil, 0)
	_ = cpu.Setup(socket)

	assert.Equal(t, banks, cpu.banks)
	assert.Equal(t, pic, cpu.pic)
}

func TestCPU_Reset(t *testing.T) {
	banks := &mockBanks{data: map[uint16]uint8{0xfffc: 0x34, 0xfffd: 0x12}}
	pic := &mockPic{}
	socket := &mockSocket{banks: banks, pic: pic}

	cpu := NewCPU(nil, nil, 0)
	cpu.Setup(socket)
	cpu.Reset()

	assert.Equal(t, uint16(0x1234), cpu.pc)
	assert.NotNil(t, cpu.next)
	assert.False(t, cpu.irqBreaker)
}

func TestCPU_SetOverflowBranch(t *testing.T) {
	cpu := NewCPU(nil, nil, 0)
	branchFn := func() bool { return true }
	cpu.SetOverflowBranch(branchFn)

	assert.NotNil(t, cpu.overflowBranch)
	assert.True(t, cpu.overflowBranch())
}

func TestCPU_SetAECLow(t *testing.T) {
	cpu := NewCPU(nil, nil, 0)

	cpu.SetAECLow(true)
	assert.True(t, cpu.aecLow)
	assert.True(t, cpu.stop)

	cpu.SetAECLow(false)
	assert.False(t, cpu.aecLow)
	assert.True(t, cpu.stop)
}

func TestCPU_SetRDYLow(t *testing.T) {
	cpu := NewCPU(nil, nil, 0)

	cpu.SetRDYLow(true)
	assert.True(t, cpu.rdyLow)
	assert.False(t, cpu.stop)

	cpu.SetRDYLow(false)
	assert.False(t, cpu.rdyLow)
	assert.False(t, cpu.stop)
}

func TestCPU_Emulate(t *testing.T) {
	cpu := NewCPU(nil, nil, 0)
	var nextCalled bool
	cpu.next = func(_ *CPU) { nextCalled = true }

	cpu.stop = true
	cpu.Emulate()
	assert.False(t, nextCalled)

	cpu.stop = false
	cpu.Emulate()
	assert.True(t, nextCalled)
}

func TestCPU_read(t *testing.T) {
	banks := &mockBanks{data: map[uint16]uint8{0x1000: 42}}
	cpu := NewCPU(nil, nil, 0)
	cpu.banks = banks

	cpu.SetRDYLow(true)
	value, ok := cpu.read(0x1000)
	assert.Equal(t, uint8(0), value)
	assert.False(t, ok)
	assert.True(t, cpu.stop)

	cpu.SetRDYLow(false)
	value, ok = cpu.read(0x1000)
	assert.Equal(t, uint8(42), value)
	assert.True(t, ok)
}

func TestCPU_popFlags(t *testing.T) {
	cpu := NewCPU(nil, nil, 0)
	const target = 0xd3

	cpu.popFlags(target) //11010011

	assert.Equal(t, uint8(target), cpu.nFlag)
	assert.Equal(t, uint8(0x40), cpu.vFlag)
	assert.Equal(t, uint8(0x00), cpu.dFlag)
	assert.Equal(t, uint8(0x00), cpu.iFlag)
	assert.Equal(t, uint8(0x00), cpu.zFlag)
	assert.Equal(t, uint8(0x01), cpu.cFlag)
}

func TestCPU_pushFlags(t *testing.T) {
	cpu := NewCPU(nil, nil, 0)
	cpu.nFlag = 0x80
	cpu.vFlag = 0x40
	cpu.dFlag = 0x08
	cpu.iFlag = 0x04
	cpu.zFlag = 0x00
	cpu.cFlag = 0x01

	data := cpu.pushFlags(true)
	assert.Equal(t, uint8(0xff), data)
	data = cpu.pushFlags(false)
	assert.Equal(t, uint8(0xef), data)
}

func TestCPU_branch(t *testing.T) {
	cpu := NewCPU(nil, nil, 0)
	cpu.pc = 0x0100

	cpu.branch(0x7f)
	assert.Equal(t, uint16(0x017f), cpu.ar)
	assert.NotNil(t, cpu.next)

	cpu.branch(0x80)
	assert.Equal(t, uint16(0x0080), cpu.ar)
	assert.NotNil(t, cpu.next)
}
