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
	cart  icartridge.ICartridge
}

func NewManager() *Manager {
	return &Manager{
		cart: nil,
	}
}

func (f *Manager) Setup(board icartridge.IExpansion, prefs *preferences.Prefs) {
	f.board = board
	f.prefs = prefs
}

func (f *Manager) Cartridge() icartridge.ICartridge {
	return f.cart
}

func (f *Manager) Read(interval icartridge.Interval, addr uint16) (uint8, bool) {
	if f.cart == nil {
		return 0, false
	}
	return f.cart.Read(interval, addr)
}

func (f *Manager) Write(interval icartridge.Interval, addr uint16, data uint8) bool {
	if f.cart == nil {
		return false
	}
	return f.cart.Write(interval, addr, data)
}

func (f *Manager) IOWrite(addr uint16, data uint8) bool {
	if f.cart == nil {
		return false
	}
	return f.cart.IOWrite(addr, data)
}

func (f *Manager) IORead(addr uint16) (uint8, bool) {
	if f.cart == nil {
		return 0, false
	}
	return f.cart.IORead(addr)
}

func (f *Manager) Load(p string) error {
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
	f.cart = cart
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
	f.cart = cart
	return nil
}
