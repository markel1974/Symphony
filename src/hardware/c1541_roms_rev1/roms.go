package c1541_roms_rev1

import (
	"github.com/markel1974/symphony/src/config"
	"github.com/markel1974/symphony/src/kernel/component"
	"github.com/markel1974/symphony/src/references"
)

// Roms represents a component that incorporates base functionality and manages system configurations and data.
type Roms struct {
	*component.BaseComponent
	cfg    *config.Config
	kernal []byte
}

// NewRoms initializes and returns a pointer to a Roms instance, registering it with the provided factory and parent component.
func NewRoms(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Roms {
	rl := &Roms{
		BaseComponent: component.NewBaseComponent(),
		kernal:        nil,
		cfg:           nil,
	}
	rl.BaseComponent.Register(factory, parent, Identifier(), instance, rl, references.IdIC1541Roms(rl, label, instance))
	return rl
}

// Setup initializes the Roms component by loading configuration and system data. Returns an error if loading fails.
func (r *Roms) Setup() error {
	r.cfg = r.GetFactory().GetConfig()
	kernal, err := r.load()
	if err != nil {
		return err
	}
	r.kernal = kernal
	return nil
}

// Bind links the Roms instance to the provided IC1541RomsSocket interface to configure its dependencies.
// Returns an error if binding fails or cannot be established.
func (r *Roms) Bind(_ references.IC1541RomsSocket) error {
	return nil
}

// Connect establishes the connection required for the Roms component and ensures it is ready for operation.
func (r *Roms) Connect() error {
	return nil
}

// Internal determines if the component operates in an internal state configuration, returning false by default.
func (r *Roms) Internal() bool {
	return false
}

// Emulate starts the ROM emulation process using the loaded configuration and system data.
func (r *Roms) Emulate() {

}

// EmulationRequired determines if emulation is required for the ROMs component and returns a boolean value.
func (r *Roms) EmulationRequired() bool {
	return false
}

// Reset clears the internal state of the Roms component, preparing it for reinitialization or reuse.
func (r *Roms) Reset() {
}

// load attempts to load the ROM data based on configuration, favoring Jiffy mode or specified file paths; falls back to builtin ROM.
func (r *Roms) load() ([]byte, error) {
	if r.cfg.UseJiffy() {
		return _jiffyRom, nil
	}
	if name := r.cfg.C1541RomAsset(); len(name) > 0 {
		if dat, err := r.cfg.AssetRead(name); err == nil {
			return dat, nil
		}
	}
	return _builtinRom, nil
}

// KernalRead returns the system ROM data as a byte slice.
func (r *Roms) KernalRead(addr uint16) uint8 {
	return r.kernal[addr&0x3fff]
}
