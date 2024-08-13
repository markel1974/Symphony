package banks

import (
	"github.com/markel1974/c64emu/src/board/iec/drives/c1541/via"
	"github.com/markel1974/c64emu/src/config"
)

const c1541RamSize = 0x0800

type Banks struct {
	ram  []uint8
	rom  []uint8
	via1 *via.Via1
	via2 *via.Via2
	cfg  *config.Config
}

func New() *Banks {
	return &Banks{ram: make([]uint8, c1541RamSize)}
}

func (r *Banks) Setup(via1 *via.Via1, via2 *via.Via2, cfg *config.Config) {
	r.via1 = via1
	r.via2 = via2
	r.cfg = cfg
	loader := NewLoader()
	r.rom = loader.Load(cfg.UseJiffy(), cfg.Get1541RomPath())
}

//func (r *Banks) AtnWakeUp() {
//Interrupt by negative edge of ATN on IEC bus
//	r.ram[0x7c] = 1
//}

func (r *Banks) ReadInterval(start uint16, count uint16) []byte {
	return r.ram[start : start+count]
}

func (r *Banks) Read(addr uint16) uint8 {
	if addr >= 0xc000 {
		return r.rom[addr&0x3fff]
	}
	if addr < 0x1000 {
		return r.ram[addr&0x07ff]
	}
	return r.readByteIO(addr)
}

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
