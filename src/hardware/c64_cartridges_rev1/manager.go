package c64_cartridges_rev1

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/hardware/c64_cartridges_rev1/catalog"
	"github.com/markel1974/c64emu/src/references"

	_ "github.com/markel1974/c64emu/src/hardware/c64_cartridges_rev1/easyflash"
	_ "github.com/markel1974/c64emu/src/hardware/c64_cartridges_rev1/external_cpu"
	_ "github.com/markel1974/c64emu/src/hardware/c64_cartridges_rev1/final_cartridge_iii"
	_ "github.com/markel1974/c64emu/src/hardware/c64_cartridges_rev1/generic"
	_ "github.com/markel1974/c64emu/src/hardware/c64_cartridges_rev1/magicdesk"
	_ "github.com/markel1974/c64emu/src/hardware/c64_cartridges_rev1/ocean"
	_ "github.com/markel1974/c64emu/src/hardware/c64_cartridges_rev1/reu"
)

//https://www.c64-wiki.com/wiki/Bank_Switching
//https://luigidifraia.wordpress.com/2021/05/08/commodore-64-cartridges-theory-of-operation-and-ocean-bank-switching-described/
//https://github.com/Project-64/reloaded/blob/master/c64/c64prg/C64PRG11.TXT#L13101

// Manager is responsible for managing cartridge interactions, configurations, and hardware registration in the system.
type Manager struct {
	*component.BaseComponent
	idx     int
	board   references.IC64Expansion
	cfg     *config.Config
	carts   []references.IC64Cartridge
	emulate []func()
	label   string
}

// NewManager initializes and returns a new instance of the Manager type, setting up default configurations and maps.
func NewManager(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Manager {
	m := &Manager{
		BaseComponent: component.NewBaseComponent(),
		idx:           0,
		board:         nil,
		cfg:           nil,
		carts:         nil,
		label:         label,
	}
	m.BaseComponent.Register(factory, parent, Identifier(), m, references.IdIC64CartridgeManager(m, label, instance))
	return m
}

// Setup initializes the Manager by setting up the expansion board, configuration preferences, and cartridge hardware mappings.
func (f *Manager) Setup() error {
	f.cfg = f.GetFactory().GetConfig()
	return nil
}

// Bind associates the Manager with an expansion board, enabling integration with the provided hardware references.
func (f *Manager) Bind(_ references.IC64CartridgeManagerSocket, board references.IC64Expansion) error {
	f.board = board
	return nil
}

// Connect establishes the Manager's connection, preparing it for interactions and operations. Returns an error if the connection fails.
func (f *Manager) Connect() error {
	return nil
}

// Internal indicates whether the Manager's operation is internal and returns a boolean result.
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
		return
	}
	for _, cart := range f.carts {
		cart.Reset()
	}
}

func (f *Manager) CreateCartridges() error {
	for _, crt := range f.cfg.Cartridges() {
		if _, err := f.Add(crt.GetKind(), crt.GetName(), crt.GetData()); err != nil {
			return err
		}
	}
	return nil
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
	if len(f.emulate) == 0 {
		return false
	}
	return true
}

// HardwareButton handles the system response to a physical button press event, updating cartridge state as necessary.
func (f *Manager) HardwareButton(pressed bool, value uint8) {
	if f.carts == nil {
		return
	}
	if len(f.carts) == 1 {
		f.carts[0].HardwareButton(pressed, value)
		return
	}
	for _, cart := range f.carts {
		cart.HardwareButton(pressed, value)
	}
}

// Read retrieves a value from the specified ROM interval and address. Returns the value and a boolean indicating success.
func (f *Manager) Read(interval references.C64RomInterval, addr uint16) (uint8, bool) {
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
func (f *Manager) Write(interval references.C64RomInterval, addr uint16, data uint8) bool {
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

func (f *Manager) IRQ(d uint32) {
	if f.carts == nil {
		return
	}
	if len(f.carts) == 1 {
		f.carts[0].IRQ(d)
		return
	}
	for _, cart := range f.carts {
		cart.IRQ(d)
	}
}

func (f *Manager) IRQClear(d uint32) {
	if f.carts == nil {
		return
	}
	if len(f.carts) == 1 {
		f.carts[0].IRQ(d)
		return
	}
	for _, cart := range f.carts {
		cart.IRQCLear(d)
	}
}

// Add registers a new cartridge using the provided hardware type, name, and data, returning an identifier or an error if failed.
func (f *Manager) Add(hardware string, name string, data []byte) (string, error) {
	id := strconv.Itoa(f.idx)
	ldr := catalog.NewLoader(id, catalog.MachineC64)
	if err := ldr.Setup(name, data); err != nil {
		return "", err
	}
	var factory func(references.IComponent, references.IComponentFactory, string, int) references.IC64Cartridge = nil
	if len(hardware) > 0 {
		hardware = strings.ToUpper(strings.TrimSpace(hardware))
		factory = catalog.ByHardware(hardware)
	} else if catalog.Type(ldr.GetType()) == catalog.TypeCrt {
		factory = catalog.ByType(int(ldr.Kind))
	} else if catalog.Type(ldr.GetType()) == catalog.TypeBin {
		if factory = catalog.BySize(len(ldr.GetData())); factory == nil {
			factory = catalog.BySizeDefault()
		}
	}
	if factory == nil {
		return "", fmt.Errorf("unsupported => %d", ldr.Kind)
	}
	cart := factory(f, f.GetFactory(), f.label, f.idx)
	f.idx++
	f.carts = append(f.carts, cart)
	if err := cart.Setup(); err != nil {
		return "", err
	}
	if err := cart.Bind(f.board, ldr); err != nil {
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
