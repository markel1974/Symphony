package mos6569

import (
	"github.com/markel1974/c64emu/src/kernel/component"
	"github.com/markel1974/c64emu/src/references"
)

const (
	ioAndCharRomArea = 0x7000
	charRomOffset    = 0x1000
)

// MemoryUnit represents a structure for managing video memory and related base addresses and operations in a system.
type MemoryUnit struct {
	*component.BaseComponent
	reflect        *MemoryUnitReflect
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

// NewMemory initializes and returns a pointer to a new MemoryUnit struct with default values.
func NewMemory(parent references.IComponent, factory references.IComponentFactory, label string, instance int, readRam func(addr uint16) uint8, readColorRam func(addr uint16) uint8, readCharRom func(addr uint16) uint8) *MemoryUnit {
	m := &MemoryUnit{
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
	m.reflect = NewMemoryUnitReflect(m, factory, parent, "memoryUnit", instance, references.IdInternalComponent(label, instance, "MemoryUnit"))
	return m
}

// Setup initializes the MemoryUnit instance, preparing it for use. Returns an error if the setup process fails.
func (m *MemoryUnit) Setup() error {
	return nil
}

// Connect establishes the necessary connections or dependencies for the MemoryUnit instance and returns an error if it fails.
func (m *MemoryUnit) Connect() error {
	return nil
}

// EmulationRequired determines if emulation is currently required for the memory instance. Returns a boolean value.
func (m *MemoryUnit) EmulationRequired() bool {
	return false
}

// Emulate performs the core video memory emulation logic associated with the MemoryUnit structure and its components.
func (m *MemoryUnit) Emulate() {
}

// Internal determines if the MemoryUnit instance is configured for internal operations. Always returns true.
func (m *MemoryUnit) Internal() bool {
	return true
}

// Reset restores the MemoryUnit object to its initial state, resetting all memory-related configurations and properties.
func (m *MemoryUnit) Reset() {
}

// SetCIAVABase updates the CIA VA14/15 video base based on the given newVA value and performs a memory pointer update.
func (m *MemoryUnit) SetCIAVABase(newVA uint8) {
	m.ciaVaBase = uint16(newVA) << 14
	m.memoryPointerUpdate()
}

// GetLastByte retrieves the value of the last byte read by the VIC from memory within the current instance.
func (m *MemoryUnit) GetLastByte() uint8 {
	return m.lastByte
}

// GetVABase retrieves the value of vaBase with the least significant bit forced to 1.
func (m *MemoryUnit) GetVABase() uint8 {
	return m.vaBase | 0x01
}

// GetCharBase retrieves the base address of the character generator from the memory structure.
func (m *MemoryUnit) GetCharBase() uint16 {
	return m.charBase
}

// GetBitmapBase retrieves the base address of the bitmap within the memory structure.
func (m *MemoryUnit) GetBitmapBase() uint16 {
	return m.bitmapBase
}

// GetMatrixBase retrieves the base address of the video matrix from the memory structure.
func (m *MemoryUnit) GetMatrixBase() uint16 {
	return m.matrixBase
}

// SetVABase sets the video address base value and updates memory pointer configurations.
func (m *MemoryUnit) SetVABase(data uint8) {
	m.vaBase = data
	m.memoryPointerUpdate()
}

// ResetRefreshCounter resets the refresh counter to its default value of 0xff.
func (m *MemoryUnit) ResetRefreshCounter() {
	m.refreshCounter = 0xff
}

// ReadByte reads a byte of data from the specified memory address and updates the last byte read.
func (m *MemoryUnit) ReadByte(addr uint16) uint8 {
	va := addr | m.ciaVaBase
	if (va & ioAndCharRomArea) == charRomOffset {
		m.lastByte = m.readCharRom(va)
		return m.lastByte
	}
	m.lastByte = m.readRam(va)
	return m.lastByte
}

// AccessIdle performs idle access by reading from the memory address 0x3FFF using the ReadByte method.
func (m *MemoryUnit) AccessIdle() {
	_ = m.ReadByte(0x3fff)
}

// AccessRefresh performs a DRAM refresh operation by reading from an address constructed from the refresh counter.
// Decrements the refresh counter after each operation.
func (m *MemoryUnit) AccessRefresh() {
	_ = m.ReadByte(0x3f00 | uint16(m.refreshCounter))
	m.refreshCounter--
}

// memoryPointerUpdate recalculates and updates the matrixBase, charBase, and bitmapBase values based on the current vaBase.
func (m *MemoryUnit) memoryPointerUpdate() {
	m.matrixBase = (uint16(m.vaBase) & 0xf0) << 6
	m.charBase = (uint16(m.vaBase) & 0x0e) << 10
	m.bitmapBase = (uint16(m.vaBase) & 0x08) << 10
}
