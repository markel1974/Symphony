package cartridges_c64

import (
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/hardware/cartridges_c64/easyflash"
	"github.com/markel1974/c64emu/src/hardware/cartridges_c64/external_cpu"
	"github.com/markel1974/c64emu/src/hardware/cartridges_c64/generic"
	"github.com/markel1974/c64emu/src/hardware/cartridges_c64/loader"
	"github.com/markel1974/c64emu/src/hardware/cartridges_c64/magicdesk"
	"github.com/markel1974/c64emu/src/hardware/cartridges_c64/ocean"
	"github.com/markel1974/c64emu/src/hardware/cartridges_c64/reu"
	"github.com/markel1974/c64emu/src/references"
	"strconv"
	"strings"
)

//https://www.c64-wiki.com/wiki/Bank_Switching
//https://luigidifraia.wordpress.com/2021/05/08/commodore-64-cartridges-theory-of-operation-and-ocean-bank-switching-described/
//https://github.com/Project-64/reloaded/blob/master/c64/c64prg/C64PRG11.TXT#L13101

// Manager is responsible for managing cartridge interactions, configurations, and hardware registration in the system.
type Manager struct {
	*component.BaseComponent
	idx                 int
	board               references.IExpansionC64
	cfg                 *config.Config
	carts               []references.ICartridgeC64
	emulate             []func()
	registerHardware    map[string]func(references.IComponent, references.IComponentFactory, int) references.ICartridgeC64
	registerType        map[int]func(references.IComponent, references.IComponentFactory, int) references.ICartridgeC64
	registerSize        map[int]func(references.IComponent, references.IComponentFactory, int) references.ICartridgeC64
	registerSizeDefault func(references.IComponent, references.IComponentFactory, int) references.ICartridgeC64
}

// NewManager initializes and returns a new instance of the Manager type, setting up default configurations and maps.
func NewManager(parent references.IComponent, factory references.IComponentFactory, instance int) *Manager {
	m := &Manager{
		BaseComponent:       component.NewBaseComponent(),
		idx:                 0,
		board:               nil,
		cfg:                 nil,
		carts:               nil,
		registerHardware:    make(map[string]func(references.IComponent, references.IComponentFactory, int) references.ICartridgeC64),
		registerType:        make(map[int]func(references.IComponent, references.IComponentFactory, int) references.ICartridgeC64),
		registerSize:        make(map[int]func(references.IComponent, references.IComponentFactory, int) references.ICartridgeC64),
		registerSizeDefault: nil,
	}
	m.BaseComponent.Register(factory, parent, Identifier(), m, references.IdICartridgeManagerC64(m, instance))
	return m
}

// Setup initializes the Manager by setting up the expansion board, configuration preferences, and cartridge hardware mappings.
func (f *Manager) Setup(_ references.ICartridgeManagerC64Socket, cfg *config.Config, board references.IExpansionC64) error {
	f.cfg = cfg
	f.board = board
	f.registerHardware[external_cpu.Id] = external_cpu.New
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

	return nil
}

func (f *Manager) Connect() error {
	return nil
}

func (f *Manager) Internal() bool {
	return false
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
		f.emulate[0]()
		return
	}
	for _, c := range f.emulate {
		c()
	}
}

func (f *Manager) EmulationRequired() bool {
	return true
}

// Read retrieves a value from the specified ROM interval and address. Returns the value and a boolean indicating success.
func (f *Manager) Read(interval references.RomInterval, addr uint16) (uint8, bool) {
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
func (f *Manager) Write(interval references.RomInterval, addr uint16, data uint8) bool {
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
	var factory func(references.IComponent, references.IComponentFactory, int) references.ICartridgeC64 = nil
	if len(hardware) > 0 {
		hardware = strings.ToUpper(strings.TrimSpace(hardware))
		factory = f.registerHardware[hardware]
	} else if loader.Type(ldr.GetType()) == loader.TypeCrt {
		factory = f.registerType[int(ldr.Kind)]
	} else if loader.Type(ldr.GetType()) == loader.TypeBin {
		if factory = f.registerSize[len(ldr.GetData())]; factory == nil {
			factory = f.registerSizeDefault
		}
	}
	if factory == nil {
		return "", fmt.Errorf("unsupported => %d", ldr.Kind)
	}
	cart := factory(f, f.GetFactory(), f.idx)
	f.idx++
	f.carts = append(f.carts, cart)
	if err := cart.Setup(f.board, ldr, f.cfg); err != nil {
		return "", err
	}
	if cart.EmulationRequired() {
		f.emulate = append(f.emulate, cart.Emulate)
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
	if !found {
		return fmt.Errorf("can't remove cartridge id %s", id)
	}
	f.emulate = []func(){}
	for _, c := range f.carts {
		if c.EmulationRequired() {
			f.emulate = append(f.emulate, c.Emulate)
		}
	}
	return nil
}
