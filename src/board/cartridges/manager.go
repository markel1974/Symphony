package cartridges

import (
	"fmt"
	"github.com/markel1974/c64emu/src/board/cartridges/cartridge16k"
	"github.com/markel1974/c64emu/src/board/cartridges/cartridge8k"
	"github.com/markel1974/c64emu/src/board/cartridges/easyflash"
	"github.com/markel1974/c64emu/src/board/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/board/cartridges/loader"
	"github.com/markel1974/c64emu/src/board/cartridges/ocean"
	"github.com/markel1974/c64emu/src/preferences"
)

//https://www.c64-wiki.com/wiki/Bank_Switching
//https://luigidifraia.wordpress.com/2021/05/08/commodore-64-cartridges-theory-of-operation-and-ocean-bank-switching-described/
//https://github.com/Project-64/reloaded/blob/master/c64/c64prg/C64PRG11.TXT#L13101

type Manager struct {
	board icartridge.IExpansion
	prefs *preferences.Prefs
	carts []icartridge.ICartridge
}

func NewManager() *Manager {
	return &Manager{
		carts: nil,
	}
}

func (f *Manager) Setup(board icartridge.IExpansion, prefs *preferences.Prefs) {
	f.board = board
	f.prefs = prefs
}

func (f *Manager) Config() (uint8, uint8, bool) {
	if f.carts == nil {
		return 0, 0, false
	}
	if len(f.carts) == 1 {
		return f.carts[0].GetGame(), f.carts[0].GetExRom(), true
	}
	g := uint8(1)
	x := uint8(1)
	for _, cart := range f.carts {
		if g == 1 {
			g = cart.GetGame()
		}
		if x == 1 {
			x = cart.GetExRom()
		}
	}
	return g, x, true
}

func (f *Manager) Read(interval icartridge.Interval, addr uint16) (uint8, bool) {
	if f.carts == nil {
		return 0, false
	}
	if len(f.carts) == 1 {
		return f.carts[0].Read(interval, addr)
	}
	val := uint8(0)
	ret := false
	for _, cart := range f.carts {
		if v, ok := cart.Read(interval, addr); ok {
			if !ret {
				val = v
				ret = true
			}
		}
	}
	return val, ret
}

func (f *Manager) Write(interval icartridge.Interval, addr uint16, data uint8) bool {
	if f.carts == nil {
		return false
	}
	if len(f.carts) == 1 {
		return f.carts[0].Write(interval, addr, data)
	}
	ret := false
	for _, cart := range f.carts {
		if ok := cart.Write(interval, addr, data); ok {
			ret = true
		}
	}
	return ret
}

func (f *Manager) IOWrite(addr uint16, data uint8) bool {
	if f.carts == nil {
		return false
	}
	if len(f.carts) == 1 {
		return f.carts[0].IOWrite(addr, data)
	}
	ret := false
	for _, cart := range f.carts {
		if ok := cart.IOWrite(addr, data); ok {
			ret = true
		}
	}
	return ret
}

func (f *Manager) IORead(addr uint16) (uint8, bool) {
	if f.carts == nil {
		return 0, false
	}
	if len(f.carts) == 1 {
		return f.carts[0].IORead(addr)
	}
	val := uint8(0)
	ret := false
	for _, cart := range f.carts {
		if v, ok := cart.IORead(addr); ok {
			if !ret {
				val = v
				ret = true
			}
		}
	}
	return val, ret
}

func (f *Manager) Add(p string) error {
	ldr := loader.NewLoader(loader.MachineC64)
	err := ldr.Setup(p)
	if err != nil {
		return err
	}
	if ldr.GetMode() == loader.ModeCrt {
		return f.loadCrt(ldr)
	}
	return f.loadBin(ldr)
}

func (f *Manager) loadCrt(ldr *loader.CRTLoader) error {
	var cart icartridge.ICartridge
	switch ldr.Kind {
	case loader.CARTRIDGE_OCEAN:
		cart = ocean.New(uint8(ldr.Game), uint8(ldr.ExRom), icartridge.ROM_LO, icartridge.ROM_HI_1)
	case loader.CARTRIDGE_EASYFLASH:
		cart = easyflash.New(uint8(ldr.Game), uint8(ldr.ExRom), icartridge.ROM_LO, icartridge.ROM_HI_1)
	}
	if cart == nil {
		return fmt.Errorf("unsupported")
	}
	if err := cart.Setup(f.board, ldr); err != nil {
		return err
	}
	f.carts = append(f.carts, cart)
	return nil
}

func (f *Manager) loadBin(ldr *loader.CRTLoader) error {
	var cart icartridge.ICartridge = nil
	lCartridge := len(ldr.GetData())
	if lCartridge == 0x2000 {
		cart = cartridge8k.New(0, 1, icartridge.ROM_LO)
	} else if lCartridge == 0x4000 {
		cart = cartridge16k.New(0, 0, icartridge.ROM_LO, icartridge.ROM_HI_1)
	} else if lCartridge > 0x4000 {
		//TODO VERIFICA
		//shadow of the beast funziona con ROM_LO, ROM_HI_1
		//robocop2 funziona con ROM_LO, ROM_HI_1
		cart = ocean.New(0, 0, icartridge.ROM_LO, icartridge.ROM_HI_1)
	}
	if cart == nil {
		return fmt.Errorf("invalid cartridge")
	}
	if err := cart.Setup(f.board, ldr); err != nil {
		return err
	}
	f.carts = append(f.carts, cart)
	return nil
}
