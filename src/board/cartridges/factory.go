package cartridges

import (
	"fmt"
	"github.com/markel1974/c64emu/src/board/cartridges/cartridge16k"
	"github.com/markel1974/c64emu/src/board/cartridges/cartridge8k"
	"github.com/markel1974/c64emu/src/board/cartridges/easyflash"
	"github.com/markel1974/c64emu/src/board/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/board/cartridges/loader"
	"github.com/markel1974/c64emu/src/board/cartridges/ocean"
	"github.com/markel1974/c64emu/src/board/iboard"
	"github.com/markel1974/c64emu/src/preferences"
)

//https://www.c64-wiki.com/wiki/Bank_Switching
//https://luigidifraia.wordpress.com/2021/05/08/commodore-64-cartridges-theory-of-operation-and-ocean-bank-switching-described/
//https://github.com/Project-64/reloaded/blob/master/c64/c64prg/C64PRG11.TXT#L13101

type Factory struct {
	board iboard.IBoard
	prefs *preferences.Prefs
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) Setup(board iboard.IBoard, prefs *preferences.Prefs) {
	f.board = board
	f.prefs = prefs
}

func (f *Factory) Load(p string) (icartridge.ICartridge, error) {
	ldr := loader.NewLoader(loader.MachineC64)
	if err := ldr.Setup(p); err != nil {
		return nil, err
	}
	if ldr.GetMode() == loader.ModeCrt {
		return f.loadCrt(ldr)
	}
	return f.loadBin(ldr)
}

func (f *Factory) loadCrt(ldr *loader.CRTLoader) (icartridge.ICartridge, error) {
	var crt icartridge.ICartridge
	switch ldr.Kind {
	case loader.CARTRIDGE_OCEAN:
		crt = ocean.New(uint8(ldr.Game), uint8(ldr.ExRom), icartridge.ROM_LO, icartridge.ROM_HI_1)
	case loader.CARTRIDGE_EASYFLASH:
		crt = easyflash.New(uint8(ldr.Game), uint8(ldr.ExRom), icartridge.ROM_LO, icartridge.ROM_HI_1)
	}
	if crt == nil {
		return nil, fmt.Errorf("unsupported")
	}
	if err := crt.Setup(f.board, ldr); err != nil {
		return nil, err
	}
	return crt, nil
}

func (f *Factory) loadBin(ldr *loader.CRTLoader) (icartridge.ICartridge, error) {
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
		return nil, fmt.Errorf("invalid cartridge")
	}
	if err := cart.Setup(f.board, ldr); err != nil {
		return nil, err
	}
	return cart, nil
}
