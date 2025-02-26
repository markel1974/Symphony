package banks

import (
	"github.com/markel1974/c64emu/src/components/via"
	"github.com/markel1974/c64emu/src/config"
)

// $0000-$07ff _ram (2K)
// $0800-$0fff _ram mirror
// $1000-$17ff free
// $1800-$1bff VIA 1
// $1c00-$1fff VIA 2
// $2000-$bfff free
// $c000-$ffff _rom (16K)

// c1541RamSize defines the size of RAM in the C1541 disk drive, measured in bytes.
const c1541RamSize = 0x0800

// Banks represents the memory and I/O peripherals for a system, including RAM, ROM, and VIA interfaces.
type Banks struct {
	ram  []uint8
	rom  []uint8
	via1 *mos6522.Via
	via2 *mos6522.Via
}

// New initializes and returns a new Banks instance with its RAM allocated to the size defined by c1541RamSize.
func New() *Banks {
	return &Banks{ram: make([]uint8, c1541RamSize)}
}

// Setup initializes the Banks instance by assigning VIA instances and loading ROM data based on the provided configuration.
func (r *Banks) Setup(via1 *mos6522.Via, via2 *mos6522.Via, cfg *config.Config) {
	r.via1 = via1
	r.via2 = via2
	loader := NewLoader()
	r.rom = loader.Load(cfg.UseJiffy(), cfg.Get1541RomPath())
}

//func (r *Banks) AtnWakeUp() {
//Interrupt by negative edge of ATN on IEC bus
//	r.ram[0x7c] = 1
//}

// ReadInterval returns a slice of bytes from the RAM, starting at the specified address and spanning the given count.
func (r *Banks) ReadInterval(start uint16, count uint16) []byte {
	return r.ram[start : start+count]
}

// Read fetches a byte of data from RAM, ROM, or I/O space based on the specified memory address.
func (r *Banks) Read(addr uint16) uint8 {
	if addr >= 0xc000 {
		return r.rom[addr&0x3fff]
	}
	if addr < 0x1000 {
		return r.ram[addr&0x07ff]
	}
	return r.readByteIO(addr)
}

// Write writes a single byte of data to the given address in the memory or I/O space managed by the Banks instance.
func (r *Banks) Write(addr uint16, data uint8) {
	if addr < 0x1000 {
		r.ram[addr&0x7ff] = data
		//if addr == 0x7c {
		//	fmt.Println("--------------------------- ADDR 0x7c", data)
		//}
		return
	}
	r.writeByteIO(addr, data)
}

// readByteIO determines the appropriate byte to return based on the address range and passes control to VIA instances if applicable.
func (r *Banks) readByteIO(addr uint16) uint8 {
	v := addr & 0xfc00
	if v == 0x1800 {
		return r.via1.ReadByte(addr)
	}
	if v == 0x1c00 {
		return r.via2.ReadByte(addr)
	}
	return uint8(addr >> 8)
}

// writeByteIO writes a byte of data to a specific I/O address space, delegating to VIA instances when applicable.
func (r *Banks) writeByteIO(addr uint16, data uint8) {
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
