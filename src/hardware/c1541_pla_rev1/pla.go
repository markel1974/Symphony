package c1541_pla_rev1

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

// $0000-$07ff _ram (2K)
// $0800-$0fff _ram mirror
// $1000-$17ff free
// $1800-$1bff VIA 1
// $1c00-$1fff VIA 2
// $2000-$bfff free
// $c000-$ffff _rom (16K)

// c1541RamSize represents the size of RAM in the C1541 disk drive, defined as 2 KB (0x0800).

// PLA represents a programmable logic array that links memory and peripheral devices in a system.
// It embeds BaseComponent and provides RAM, ROM, and connections to two VIA components.
type PLA struct {
	*component.BaseComponent
	ram        references.IRamC1541
	bankRead   func(uint16) uint8
	bankWrite  func(uint16, uint8)
	kernalRead func(uint16) uint8
	via1       references.IVIA
	via2       references.IVIA
}

// NewPLA initializes and returns a new instance of the PLA structure with specified parent, factory, and instance ID.
// It sets up the PLA's RAM with a predefined size and registers the component in a hierarchy via the factory.
func NewPLA(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *PLA {
	p := &PLA{
		BaseComponent: component.NewBaseComponent(),
	}
	p.BaseComponent.Register(factory, parent, Identifier(), p, references.IdIPLAc1541(p, label, instance))
	return p
}

func (r *PLA) Setup() error {
	return nil
}

func (r *PLA) Bind(_ references.IPLAc1541Socket, via1 references.IVIA, via2 references.IVIA, ram references.IRamC1541, roms references.IRomsC1541) error {
	r.via1 = via1
	r.via2 = via2
	r.ram = ram
	r.bankRead = r.ram.Read
	r.bankWrite = r.ram.Write
	r.kernalRead = roms.KernalRead
	return nil
}

// Connect associates the PLA with two VIA interfaces and loads ROM data via the provided ROM loader.
func (r *PLA) Connect() error {
	return nil
}

func (r *PLA) Internal() bool {
	return false
}

// Reset restores the PLA instance to its initial state, clearing relevant data and resetting internal components.
func (r *PLA) Reset() {

}

// Emulate performs a cycle of emulation logic for the PLA, handling memory and I/O operations as required.
func (r *PLA) Emulate() {}

// EmulationRequired determines if emulation is needed for the current component. Always returns false for this implementation.
func (r *PLA) EmulationRequired() bool {
	return false
}

// Read retrieves a byte of data from the specified memory address, accessing either ROM, RAM, or I/O based on the address.
func (r *PLA) Read(addr uint16) uint8 {
	if addr >= 0xc000 {
		return r.kernalRead(addr)
	}
	if addr < 0x1000 {
		return r.bankRead(addr & 0x07ff)
	}
	return r.readByteIO(addr)
}

// Write stores a byte of data at a given memory address, handling RAM or calling IO write methods based on the address range.
func (r *PLA) Write(addr uint16, data uint8) {
	if addr < 0x1000 {
		r.bankWrite(addr&0x7ff, data)
		//if addr == 0x7c {
		//	fmt.Println("--------------------------- ADDR 0x7c", data)
		//}
		return
	}
	r.writeByteIO(addr, data)
}

// readByteIO reads a byte from a specified I/O address.
// Delegates the read operation to via1 or via2 based on the address range.
// Returns the high byte of the address if no match is found.
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

// writeByteIO handles the write operations to I/O devices based on the given address and data.
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
