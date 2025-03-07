package cartridges

import (
	"fmt"
	"github.com/markel1974/c64emu/src/c64/cartridges/easyflash"
	"github.com/markel1974/c64emu/src/c64/cartridges/externalcpu"
	"github.com/markel1974/c64emu/src/c64/cartridges/generic"
	"github.com/markel1974/c64emu/src/c64/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/c64/cartridges/loader"
	"github.com/markel1974/c64emu/src/c64/cartridges/magicdesk"
	"github.com/markel1974/c64emu/src/c64/cartridges/ocean"
	"github.com/markel1974/c64emu/src/c64/cartridges/reu"
	"github.com/markel1974/c64emu/src/components/board"
	"github.com/markel1974/c64emu/src/config"
	"strconv"
	"strings"
)

//https://www.c64-wiki.com/wiki/Bank_Switching
//https://luigidifraia.wordpress.com/2021/05/08/commodore-64-cartridges-theory-of-operation-and-ocean-bank-switching-described/
//https://github.com/Project-64/reloaded/blob/master/c64/c64prg/C64PRG11.TXT#L13101

// Manager is responsible for managing cartridge interactions, configurations, and hardware registration in the system.
type Manager struct {
	*board.BaseComponent
	idx                 int
	board               icartridge.IExpansion
	prefs               *config.Config
	carts               []icartridge.ICartridge
	emulate             []icartridge.ICartridge
	registerHardware    map[string]func(*board.Node, string) icartridge.ICartridge
	registerType        map[int]func(*board.Node, string) icartridge.ICartridge
	registerSize        map[int]func(*board.Node, string) icartridge.ICartridge
	registerSizeDefault func(*board.Node, string) icartridge.ICartridge
}

// NewManager initializes and returns a new instance of the Manager type, setting up default configurations and maps.
func NewManager(node *board.Node, suffix string) *Manager {
	return &Manager{
		BaseComponent:       board.NewBaseComponent(node, "cartridgeManager", suffix, nil),
		idx:                 0,
		board:               nil,
		prefs:               nil,
		carts:               nil,
		registerHardware:    make(map[string]func(*board.Node, string) icartridge.ICartridge),
		registerType:        make(map[int]func(*board.Node, string) icartridge.ICartridge),
		registerSize:        make(map[int]func(*board.Node, string) icartridge.ICartridge),
		registerSizeDefault: nil,
	}
}

// Setup initializes the Manager by setting up the expansion board, configuration preferences, and cartridge hardware mappings.
func (f *Manager) Setup(board icartridge.IExpansion, prefs *config.Config) {
	f.board = board
	f.prefs = prefs
	f.registerHardware[externalcpu.Id] = externalcpu.New
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

// Reset resets all cartridges managed by the Manager. If no cartridges exist, it performs no operations.
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

// Config returns the Game and ExROM status along with a boolean indicating whether configuration was successful.
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

// Emulate performs a simulation step for all emulated cartridges in the `emulate` slice.
// If no cartridges exist, the method exits immediately.
// If only one cartridge exists, it directly delegates the emulate call to it.
// For multiple cartridges, it iterates and calls the emulate function for each.
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

// Read retrieves a value from the specified ROM interval and address. Returns the value and a boolean indicating success.
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

// Write attempts to write the given data to the specified address within the provided ROM interval for all managed cartridges.
// Returns true if any cartridge successfully handles the write operation, otherwise returns false.
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

// IOWrite attempts to write the specified data to the given address using all attached cartridges.
// Returns true if at least one cartridge successfully processes the I/O write operation.
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

// IORead reads a byte from the specified I/O address and returns the value along with a flag indicating success.
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

// Add registers a new cartridge using the provided hardware type, name, and data, returning an identifier or an error if failed.
func (f *Manager) Add(hardware string, name string, data []byte) (string, error) {
	id := strconv.Itoa(f.idx)
	ldr := loader.NewLoader(id, loader.MachineC64)
	if err := ldr.Setup(name, data); err != nil {
		return "", err
	}
	var factory func(*board.Node, string) icartridge.ICartridge = nil
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
	cart := factory(f.GetNode(), strconv.Itoa(f.idx))
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

// Remove removes a cartridge and its emulation by the given ID from the Manager's list. Returns an error if ID is not found.
func (f *Manager) Remove(id string) error {
	found := false
	for s, cart := range f.carts {
		if cart.GetLoaderId() == id {
			f.carts = append(f.carts[:s], f.carts[s+1:]...)
			found = true
			break
		}
	}
	for e, c := range f.emulate {
		if c.GetLoaderId() == id {
			f.emulate = append(f.emulate[:e], f.emulate[e+1:]...)
			break
		}
	}
	if !found {
		return fmt.Errorf("can't remove cartridge id %s", id)
	}
	return nil
}
