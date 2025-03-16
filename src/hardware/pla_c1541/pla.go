package pla_c1541

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// $0000-$07ff _ram (2K)
// $0800-$0fff _ram mirror
// $1000-$17ff free
// $1800-$1bff VIA 1
// $1c00-$1fff VIA 2
// $2000-$bfff free
// $c000-$ffff _rom (16K)

// c1541RamSize defines the size of the RAM for a C1541 emulator, set to 2 KB (0x0800 bytes).
const c1541RamSize = 0x0800

// PLA represents a memory management and I/O coordination unit, including RAM, ROM, and communication with VIAs.
type PLA struct {
	*component.BaseComponent
	factory references.IComponentFactory
	ram     []uint8
	rom     []uint8
	via1    references.IVIA
	via2    references.IVIA
}

func NewPLAComponent(parent references.IComponent, factory references.IComponentFactory, suffix string) references.IComponent {
	return NewPLA(parent, factory, suffix)
}

// NewPLA initializes and returns a pointer to a new instance of PLA with default memory and configurations set.
func NewPLA(parent references.IComponent, factory references.IComponentFactory, suffix string) *PLA {
	p := &PLA{
		BaseComponent: component.NewBaseComponent(componentId, suffix),
		factory:       factory,
		ram:           make([]uint8, c1541RamSize),
	}
	component.Register(parent, p)
	return p
}

// New creates and returns a new instance of PLA with initialized RAM of the specified size.
//func New() *PLA {
//	return &PLA{ram: make([]uint8, c1541RamSize)}
//}

// Setup initializes the PLA instance by configuring VIA components and loading required ROM based on the provided configuration.
func (r *PLA) Setup(via1 references.IVIA, via2 references.IVIA, roms references.IROMLoaderC1541, cfg *config.Config) error {
	r.via1 = via1
	r.via2 = via2
	r.rom = roms.Load()
	return nil
}

func (r *PLA) Reset() {

}

//func (r *PLA) AtnWakeUp() {
//Interrupt by negative edge of ATN on IEC bus
//	r.ram[0x7c] = 1
//}

// ReadInterval returns a slice of bytes from the RAM starting at the provided start address for the specified count.
func (r *PLA) ReadInterval(start uint16, count uint16) []byte {
	return r.ram[start : start+count]
}

// Read retrieves a byte at the specified memory address using prioritized access to ROM, RAM, or I/O.
func (r *PLA) Read(addr uint16) uint8 {
	if addr >= 0xc000 {
		return r.rom[addr&0x3fff]
	}
	if addr < 0x1000 {
		return r.ram[addr&0x07ff]
	}
	return r.readByteIO(addr)
}

// Write writes an 8-bit unsigned data value to the specified 16-bit memory address in the PLA instance.
func (r *PLA) Write(addr uint16, data uint8) {
	if addr < 0x1000 {
		r.ram[addr&0x7ff] = data
		//if addr == 0x7c {
		//	fmt.Println("--------------------------- ADDR 0x7c", data)
		//}
		return
	}
	r.writeByteIO(addr, data)
}

// readByteIO reads a byte from a specified I/O address, delegating to via1 or via2 if within their address ranges.
func (r *PLA) readByteIO(addr uint16) uint8 {
	v := addr & 0xfc00
	if v == 0x1800 {
		return r.via1.ReadByte(addr)
	}
	if v == 0x1c00 {
		return r.via2.ReadByte(addr)
	}
	return uint8(addr >> 8)
}

// writeByteIO handles the writing of a byte to the VIA1 or VIA2 interfaces based on the provided address range conditions.
func (r *PLA) writeByteIO(addr uint16, data uint8) {
	v := addr & 0xfc00
	if v == 0x1800 {
		r.via1.WriteByte(addr, data)
		return
	}
	if v == 0x1c00 {
		r.via2.WriteByte(addr, data)
		return
	}
}
