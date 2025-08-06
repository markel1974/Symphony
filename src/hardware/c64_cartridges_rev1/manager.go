package c64_cartridges_rev1

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/hardware/c64_cartridges_rev1/catalog"
	"github.com/markel1974/c64emu/src/kernel/component"
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

// Manager represents a central managing unit for IC64 components, holding configuration, state, and operational functions.
type Manager struct {
	*component.BaseComponent
	idx                 int
	board               references.IC64Expansion
	cfg                 *config.Config
	carts               []references.IC64Cartridge
	label               string
	multiDefault        references.IC64Cartridge
	resetFn             func()
	emulationRequiredFn func() bool
	emulateFn           func()
	hardwareButtonFn    func(pressed bool, value uint8)
	configFn            func() (uint8, uint8, bool)
	readFn              func(addr uint16) uint8
	ioWriteFn           func(addr uint16, data uint8) bool
	ioReadFn            func(addr uint16) (uint8, bool)
	irqFn               func(d uint32)
	irqClearFn          func(d uint32)
}

// NewManager creates and initializes a new Manager instance with the specified parent, factory, label, and instance ID.
// It assigns default empty or no-op functions for various hardware-related operations and registers the component.
func NewManager(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Manager {
	m := &Manager{
		BaseComponent: component.NewBaseComponent(),
		idx:           0,
		board:         nil,
		cfg:           nil,
		carts:         nil,
		label:         label,
	}
	m.BaseComponent.Register(factory, parent, Identifier(), instance, m, references.IdIC64CartridgeManager(m, label, instance))
	return m
}

// Setup initializes the configuration for the Manager using its associated factory and returns an error if it fails.
func (f *Manager) Setup() error {
	f.cfg = f.GetFactory().GetConfig()
	f.carts = nil
	f.rebuildEmulation(f.carts)
	return nil
}

// Bind links the Manager to the IC64Expansion implementation, storing the provided board instance for future interactions.
func (f *Manager) Bind(_ references.IC64CartridgeManagerSocket, board references.IC64Expansion) error {
	f.board = board
	return nil
}

// Connect establishes a connection between the Manager and its associated components. Returns an error on failure.
func (f *Manager) Connect() error {
	return nil
}

// Internal determines if the manager operates in internal mode. Returns false by default.
func (f *Manager) Internal() bool {
	return false
}

// Reset invokes the manager's `resetFn` function to reset its state or associated components.
func (f *Manager) Reset() {
	f.resetFn()
}

// Config retrieves configuration values for the Manager by invoking its associated configuration function.
func (f *Manager) Config() (uint8, uint8, bool) {
	return f.configFn()
}

// Emulate invokes the emulation function associated with the current Manager instance.
func (f *Manager) Emulate() {
	f.emulateFn()
}

// EmulationRequired checks if emulation is required by delegating the decision to the configured function.
func (f *Manager) EmulationRequired() bool {
	return f.emulationRequiredFn()
}

// HardwareButton invokes a function to handle a hardware button press with the given state and value.
func (f *Manager) HardwareButton(pressed bool, value uint8) {
	f.hardwareButtonFn(pressed, value)
}

// Read invokes the associated read function with the given memory address and returns the resulting byte value.
func (f *Manager) Read(addr uint16) uint8 {
	return f.readFn(addr)
}

// IOWrite performs an I/O write operation to a specified address with the given data. Returns true on success, false otherwise.
func (f *Manager) IOWrite(addr uint16, data uint8) bool {
	return f.ioWriteFn(addr, data)
}

// IORead reads data from the IO memory at the specified address and returns the value and a success flag.
// It delegates the read operation to the underlying ioReadFn function.
func (f *Manager) IORead(addr uint16) (uint8, bool) {
	return f.ioReadFn(addr)
}

// IRQ triggers an interrupt request in the Manager using the supplied data (d).
func (f *Manager) IRQ(d uint32) {
	f.irqFn(d)
}

// IRQClear clears a specific IRQ signal determined by the given identifier `d`.
func (f *Manager) IRQClear(d uint32) {
	f.irqClearFn(d)
}

// CreateCartridges initializes cartridges defined in the configuration and adds them to the Manager. Returns an error if any fail.
func (f *Manager) CreateCartridges() error {
	for _, crt := range f.cfg.Cartridges() {
		if _, err := f.Add(crt.GetKind(), crt.GetName(), crt.GetData()); err != nil {
			return err
		}
	}
	return nil
}

// Add initializes and adds a new hardware component using the provided hardware type, name, and data, returning its ID or an error.
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
	} else if catalog.Type(ldr.Type()) == catalog.TypeCrt {
		factory = catalog.ByType(int(ldr.Model()))
	} else if catalog.Type(ldr.Type()) == catalog.TypeBin {
		if factory = catalog.BySize(len(ldr.Data())); factory == nil {
			factory = catalog.BySizeDefault()
		}
	}
	if factory == nil {
		return "", fmt.Errorf("unsupported => %d", ldr.Model())
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
	f.rebuildEmulation(f.carts)
	return id, nil
}

// Remove removes a cartridge identified by the given id from the manager. Returns an error if the id is not found.
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
	f.rebuildEmulation(f.carts)
	return nil
}

// rebuildEmulation reconfigures emulation functions for the Manager based on the number and state of loaded cartridges.
func (f *Manager) rebuildEmulation(carts []references.IC64Cartridge) {
	f.multiDefault = nil

	if carts == nil {
		f.resetFn = f.resetEmpty
		f.emulationRequiredFn = f.emulationRequiredEmpty
		f.emulateFn = f.emulateEmpty
		f.hardwareButtonFn = f.hardwareButtonEmpty
		f.configFn = f.configEmpty
		f.readFn = f.readEmpty
		f.ioWriteFn = f.ioWriteEmpty
		f.ioReadFn = f.ioReadEmpty
		f.irqFn = f.irqEmpty
		f.irqClearFn = f.irqClearEmpty
		return
	}
	if len(carts) == 1 {
		f.resetFn = carts[0].Reset
		f.emulationRequiredFn = carts[0].EmulationRequired
		f.emulateFn = carts[0].Emulate
		f.hardwareButtonFn = carts[0].HardwareButton
		f.configFn = carts[0].Config
		f.readFn = carts[0].Read
		f.ioWriteFn = carts[0].IOWrite
		f.ioReadFn = carts[0].IORead
		f.irqFn = carts[0].IRQ
		f.irqClearFn = carts[0].IRQCLear
		return
	}
	f.resetFn = f.resetMulti
	f.emulationRequiredFn = f.emulationRequiredMulti
	f.emulateFn = f.emulateMulti
	f.hardwareButtonFn = f.hardwareButtonMulti
	f.ioWriteFn = f.ioWriteMulti
	f.ioReadFn = f.ioReadMulti
	f.irqFn = f.irqMulti
	f.irqClearFn = f.irqClearMulti
	f.configFn = f.configEmpty
	f.readFn = f.readEmpty
	specOff := references.C64CartridgeSpecOff
	for _, cart := range carts {
		if game, exRom, ok := cart.Config(); ok {
			if game != specOff.Game || exRom != specOff.ExRom {
				f.configFn = cart.Config
				f.readFn = cart.Read
				break
			}
		}
	}
}

// resetEmpty performs a no-operation reset action for the Manager, used when no cartridges are present or configured.
func (f *Manager) resetEmpty() {
}

// resetMulti iterates over all cartridges in the Manager and calls their Reset method.
func (f *Manager) resetMulti() {
	for _, cart := range f.carts {
		cart.Reset()
	}
}

// emulationRequiredEmpty returns false indicating that emulation is not required.
func (f *Manager) emulationRequiredEmpty() bool {
	return false
}

// emulationRequiredMulti determines if emulation is required by checking all attached cartridges. Returns true if required.
func (f *Manager) emulationRequiredMulti() bool {
	for _, cart := range f.carts {
		if cart.EmulationRequired() {
			return true
		}
	}
	return false
}

// emulateEmpty is a no-op method used when no emulation logic is required.
func (f *Manager) emulateEmpty() {
}

// emulateMulti executes the Emulate method on each IC64Cartridge instance stored in the Manager.
func (f *Manager) emulateMulti() {
	for _, c := range f.carts {
		c.Emulate()
	}
}

// configEmpty returns default configuration values with no active cartridges.
func (f *Manager) configEmpty() (uint8, uint8, bool) {
	return 0, 0, false
}

// hardwareButtonEmpty is a no-op placeholder for handling hardware button interactions with the given pressed state and value.
func (f *Manager) hardwareButtonEmpty(pressed bool, value uint8) {
}

// hardwareButtonSingle calls the HardwareButton method of the first cartridge in the carts slice with the given parameters.
func (f *Manager) hardwareButtonSingle(pressed bool, value uint8) {
	f.carts[0].HardwareButton(pressed, value)
}

// hardwareButtonMulti relays the hardware button event to all connected cartridges with the given pressed state and value.
func (f *Manager) hardwareButtonMulti(pressed bool, value uint8) {
	for _, cart := range f.carts {
		cart.HardwareButton(pressed, value)
	}
}

// readEmpty is a placeholder method that always returns 0, providing a default read behavior for unimplemented cases.
func (f *Manager) readEmpty(_ uint16) uint8 {
	return 0
}

// ioWriteEmpty is a no-op implementation for the IOWrite operation, always returning false.
func (f *Manager) ioWriteEmpty(_ uint16, _ uint8) bool {
	return false
}

// ioWriteMulti attempts to write data to the specified address through multiple cartridges and returns true if any succeed.
func (f *Manager) ioWriteMulti(addr uint16, data uint8) bool {
	ret := false
	for _, cart := range f.carts {
		if ok := cart.IOWrite(addr, data); ok {
			ret = true
		}
	}
	return ret
}

// ioReadEmpty returns default values (0, false) indicating no data is read for the provided address.
func (f *Manager) ioReadEmpty(_ uint16) (uint8, bool) {
	return 0, false
}

// ioReadMulti reads an 8-bit value from a given address across multiple cartridges and returns the value and success status.
func (f *Manager) ioReadMulti(addr uint16) (uint8, bool) {
	for _, cart := range f.carts {
		if v, ok := cart.IORead(addr); ok {
			return v, true
		}
	}
	return 0, false
}

// irqEmpty is a no-op implementation for handling IRQ logic, typically used when no IRQ processing is required.
func (f *Manager) irqEmpty(_ uint32) {
}

// irqMulti triggers the IRQ signal with the provided data on all connected cartridges in the manager.
func (f *Manager) irqMulti(d uint32) {
	for _, cart := range f.carts {
		cart.IRQ(d)
	}
}

// irqClearEmpty provides a no-op implementation for clearing IRQ in scenarios where no cartridges are present.
func (f *Manager) irqClearEmpty(d uint32) {
}

// irqClearMulti clears IRQ for multiple cartridges by delegating the operation to each cartridge in the manager's list.
func (f *Manager) irqClearMulti(d uint32) {
	for _, cart := range f.carts {
		cart.IRQCLear(d)
	}
}
