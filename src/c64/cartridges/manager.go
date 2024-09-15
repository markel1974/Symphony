package cartridges

import (
	"fmt"
	"github.com/markel1974/c64emu/src/c64/cartridges/easyflash"
	"github.com/markel1974/c64emu/src/c64/cartridges/generic"
	"github.com/markel1974/c64emu/src/c64/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/c64/cartridges/loader"
	"github.com/markel1974/c64emu/src/c64/cartridges/magicdesk"
	"github.com/markel1974/c64emu/src/c64/cartridges/ocean"
	"github.com/markel1974/c64emu/src/c64/cartridges/reu"
	"github.com/markel1974/c64emu/src/c64/cartridges/supercpu"
	"github.com/markel1974/c64emu/src/config"
	"strconv"
	"strings"
)

//https://www.c64-wiki.com/wiki/Bank_Switching
//https://luigidifraia.wordpress.com/2021/05/08/commodore-64-cartridges-theory-of-operation-and-ocean-bank-switching-described/
//https://github.com/Project-64/reloaded/blob/master/c64/c64prg/C64PRG11.TXT#L13101

type Manager struct {
	idx                 int
	board               icartridge.IExpansion
	prefs               *config.Config
	carts               []icartridge.ICartridge
	emulate             []icartridge.ICartridge
	registerHardware    map[string]func() icartridge.ICartridge
	registerType        map[int]func() icartridge.ICartridge
	registerSize        map[int]func() icartridge.ICartridge
	registerSizeDefault func() icartridge.ICartridge
}

func NewManager() *Manager {
	return &Manager{
		idx:                 0,
		board:               nil,
		prefs:               nil,
		carts:               nil,
		registerHardware:    make(map[string]func() icartridge.ICartridge),
		registerType:        make(map[int]func() icartridge.ICartridge),
		registerSize:        make(map[int]func() icartridge.ICartridge),
		registerSizeDefault: nil,
	}
}

func (f *Manager) Setup(board icartridge.IExpansion, prefs *config.Config) {
	f.board = board
	f.prefs = prefs
	f.registerHardware[supercpu.Id] = supercpu.New
	f.registerHardware[reu.Id128K] = reu.New128K
	f.registerHardware[reu.Id256K] = reu.New256K
	f.registerHardware[reu.Id512K] = reu.New512K
	f.registerHardware[reu.Id1M] = reu.New1M
	f.registerHardware[reu.Id2M] = reu.New2M
	f.registerHardware[reu.Id4M] = reu.New4M
	f.registerHardware[reu.Id8M] = reu.New8M
	f.registerHardware[reu.Id16M] = reu.New16M

	f.registerType[ocean.GetType()] = ocean.New
	f.registerType[magicdesk.GetType()] = magicdesk.New
	f.registerType[easyflash.GetType()] = easyflash.New
	f.registerType[generic.GetType()] = generic.New
	//f.registerType[cartridge16k.GetType()] = cartridge16k.New

	f.registerSize[0x2000] = generic.New //cartridge8k.New
	f.registerSize[0x4000] = generic.New //cartridge16k.New
	f.registerSizeDefault = ocean.New
}

func (f *Manager) Reset() {
	if f.carts == nil {
		return
	}
	if len(f.carts) == 1 {
		f.carts[0].Reset()
	}
	for _, cart := range f.carts {
		cart.Reset()
	}
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

func (f *Manager) Emulate() {
	if len(f.emulate) == 0 {
		return
	}
	if len(f.emulate) == 1 {
		f.emulate[0].Emulate()
		return
	}
	for _, c := range f.emulate {
		c.Emulate()
	}
}

func (f *Manager) Read(interval icartridge.RomInterval, addr uint16) (uint8, bool) {
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

func (f *Manager) Write(interval icartridge.RomInterval, addr uint16, data uint8) bool {
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

func (f *Manager) Add(hardware string, name string, data []byte) (string, error) {
	id := strconv.Itoa(f.idx)
	ldr := loader.NewLoader(id, loader.MachineC64)
	if err := ldr.Setup(name, data); err != nil {
		return "", err
	}
	var factory func() icartridge.ICartridge = nil
	if len(hardware) > 0 {
		hardware = strings.ToUpper(strings.TrimSpace(hardware))
		factory = f.registerHardware[hardware]
	} else if ldr.GetType() == loader.TypeCrt {
		factory = f.registerType[int(ldr.Kind)]
	} else if ldr.GetType() == loader.TypeBin {
		if factory = f.registerSize[len(ldr.GetData())]; factory == nil {
			factory = f.registerSizeDefault
		}
	}
	if factory == nil {
		return "", fmt.Errorf("unsupported => %d", ldr.Kind)
	}
	cart := factory()
	f.idx++
	f.carts = append(f.carts, cart)
	if err := cart.Setup(f.board, ldr); err != nil {
		return "", err
	}
	if cart.EmulationRequired() {
		f.emulate = append(f.emulate, cart)
	}
	return id, nil
}

func (f *Manager) Remove(id string) error {
	found := false
	for s, cart := range f.carts {
		if cart.GetId() == id {
			f.carts = append(f.carts[:s], f.carts[s+1:]...)
			found = true
			break
		}
	}
	for e, c := range f.emulate {
		if c.GetId() == id {
			f.emulate = append(f.emulate[:e], f.emulate[e+1:]...)
			break
		}
	}
	if !found {
		return fmt.Errorf("can't remove cartridge id %s", id)
	}
	return nil
}
