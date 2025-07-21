package mos6569

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

// Memory represents a structure for managing video memory and related base addresses and operations in a system.
type Memory struct {
	*component.BaseComponent
	vaBase         uint8  // vaBase
	ciaVaBase      uint16 // CIA VA14/15 video base
	matrixBase     uint16 // Video matrix base
	charBase       uint16 // Character generator base
	bitmapBase     uint16 // Bitmap base
	lastByte       uint8  // Last byte read by VIC
	refreshCounter uint8  // refreshCounter tracks the number of times a refresh operation has been performed.
	readRam        func(addr uint16) uint8
	readColorRam   func(addr uint16) uint8
	readCharRom    func(addr uint16) uint8
}

// NewMemory initializes and returns a pointer to a new Memory struct with default values.
func NewMemory(parent references.IComponent, factory references.IComponentFactory, label string, instance int, readRam func(addr uint16) uint8, readColorRam func(addr uint16) uint8, readCharRom func(addr uint16) uint8) *Memory {
	m := &Memory{
		BaseComponent:  component.NewBaseComponent(),
		readRam:        readRam,
		readColorRam:   readColorRam,
		readCharRom:    readCharRom,
		matrixBase:     0,
		charBase:       0,
		bitmapBase:     0,
		vaBase:         0,
		ciaVaBase:      0,
		lastByte:       0,
		refreshCounter: 0,
	}
	m.BaseComponent.Register(factory, parent, "memory", m, references.IdInternalComponent(label, instance, "Memory"))
	return m
}

// Setup initializes the Memory instance, preparing it for use. Returns an error if the setup process fails.
func (m *Memory) Setup() error {
	return nil
}

// Connect establishes the necessary connections or dependencies for the Memory instance and returns an error if it fails.
func (m *Memory) Connect() error {
	return nil
}

// EmulationRequired determines if emulation is currently required for the memory instance. Returns a boolean value.
func (m *Memory) EmulationRequired() bool {
	return false
}

// Emulate performs the core video memory emulation logic associated with the Memory structure and its components.
func (m *Memory) Emulate() {
}

// Internal determines if the Memory instance is configured for internal operations. Always returns true.
func (m *Memory) Internal() bool {
	return true
}

// Reset restores the Memory object to its initial state, resetting all memory-related configurations and properties.
func (m *Memory) Reset() {
}

// SetCIAVABase updates the CIA VA14/15 video base based on the given newVA value and performs a memory pointer update.
func (m *Memory) SetCIAVABase(newVA uint8) {
	m.ciaVaBase = uint16(newVA) << 14
	m.memoryPointerUpdate()
}

// GetLastByte retrieves the value of the last byte read by the VIC from memory within the current instance.
func (m *Memory) GetLastByte() uint8 {
	return m.lastByte
}

// GetVABase retrieves the value of vaBase with the least significant bit forced to 1.
func (m *Memory) GetVABase() uint8 {
	return m.vaBase | 0x01
}

// GetCharBase retrieves the base address of the character generator from the memory structure.
func (m *Memory) GetCharBase() uint16 {
	return m.charBase
}

// GetBitmapBase retrieves the base address of the bitmap within the memory structure.
func (m *Memory) GetBitmapBase() uint16 {
	return m.bitmapBase
}

// GetMatrixBase retrieves the base address of the video matrix from the memory structure.
func (m *Memory) GetMatrixBase() uint16 {
	return m.matrixBase
}

// SetVABase sets the video address base value and updates memory pointer configurations.
func (m *Memory) SetVABase(data uint8) {
	m.vaBase = data
	m.memoryPointerUpdate()
}

// ResetRefreshCounter resets the refresh counter to its default value of 0xff.
func (m *Memory) ResetRefreshCounter() {
	m.refreshCounter = 0xff
}

// ReadByte reads a byte of data from the specified memory address and updates the last byte read.
func (m *Memory) ReadByte(addr uint16) uint8 {
	va := addr | m.ciaVaBase
	if (va & ioAndCharRomArea) == charRomOffset {
		m.lastByte = m.readCharRom(va)
		return m.lastByte
	}
	m.lastByte = m.readRam(va)
	return m.lastByte
}

// AccessIdle performs an idle access by reading from the memory address 0x3FFF using the ReadByte method.
func (m *Memory) AccessIdle() {
	_ = m.ReadByte(0x3fff)
}

// AccessRefresh performs a DRAM refresh operation by reading from an address constructed from the refresh counter.
// Decrements the refresh counter after each operation.
func (m *Memory) AccessRefresh() {
	_ = m.ReadByte(0x3f00 | uint16(m.refreshCounter))
	m.refreshCounter--
}

// memoryPointerUpdate recalculates and updates the matrixBase, charBase, and bitmapBase values based on the current vaBase.
func (m *Memory) memoryPointerUpdate() {
	m.matrixBase = (uint16(m.vaBase) & 0xf0) << 6
	m.charBase = (uint16(m.vaBase) & 0x0e) << 10
	m.bitmapBase = (uint16(m.vaBase) & 0x08) << 10
}
