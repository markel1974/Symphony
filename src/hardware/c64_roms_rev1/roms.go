package c64_roms_rev1

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// Roms represents a structure containing ROM data and configuration for components in a computing system.
type Roms struct {
	*component.BaseComponent
	cfg    *config.Config
	kernal []byte
	basic  []byte
	char   []byte
}

// NewRomLoader initializes a new Roms instance with its base component, configuration, and ROM resources.
// It registers the component in the factory with specified parent, label, and instance identifier.
func NewRomLoader(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Roms {
	rl := &Roms{
		BaseComponent: component.NewBaseComponent(),
		cfg:           nil,
		kernal:        nil,
		basic:         nil,
		char:          nil,
	}
	rl.BaseComponent.Register(factory, parent, Identifier(), rl, references.IdIC64Roms(rl, label, instance))
	return rl
}

// Setup initializes the Roms object by loading necessary configuration and ROM data. Returns an error if initialization fails.
func (r *Roms) Setup() error {
	var err error
	r.cfg = r.GetFactory().GetConfig()
	if r.cfg.UseJiffy() {
		r.kernal = _builtinKernalJiffyRom
	} else {
		if r.kernal, err = r.load(_builtinKernalRom, r.cfg.C64RomKernalAsset()); err != nil {
			return err
		}
	}
	r.basic, err = r.load(_builtinBasicRom, r.cfg.C64RomBasicAsset())
	if err != nil {
		return err
	}
	r.char, err = r.load(_builtinCharRom, r.cfg.C64RomCharAsset())
	if err != nil {
		return err
	}
	return nil
}

// Bind associates the Roms instance with the provided IC64RomsSocket interface and initializes necessary connections.
func (r *Roms) Bind(_ references.IC64RomsSocket) error {
	return nil
}

// Connect establishes internal connections and prepares the Roms component for operation. It returns an error if initialization fails.
func (r *Roms) Connect() error {
	return nil
}

// Internal determines whether the ROM uses an internal configuration or default setup. Always returns false.
func (r *Roms) Internal() bool {
	return false
}

// Emulate initiates or performs the ROM emulation logic for the Roms component. This method handles emulation processes.
func (r *Roms) Emulate() {
}

// EmulationRequired determines if emulation is necessary for the component. Always returns false in current implementation.
func (r *Roms) EmulationRequired() bool {
	return false
}

// Reset restores the state of the Roms object to its initial configuration or default settings.
func (r *Roms) Reset() {
}

// KernalRead retrieves a byte from the kernal ROM at the specified memory address.
func (r *Roms) KernalRead(addr uint16) uint8 {
	return r.kernal[addr&0x1fff]
}

// BasicRead retrieves a byte from the BASIC ROM at the specified memory address.
func (r *Roms) BasicRead(addr uint16) uint8 {
	return r.basic[addr&0x1fff]
}

// CharRead reads a byte from the character ROM at the specified 16-bit address and returns it.
func (r *Roms) CharRead(addr uint16) uint8 {
	return r.char[addr&0x0fff]
}

// load attempts to load a ROM from the provided asset path or returns the default ROM if the asset is unavailable.
func (r *Roms) load(defaultRom []byte, asset string) ([]byte, error) {
	if len(asset) == 0 {
		return defaultRom, nil
	}
	if dat, err := r.cfg.Asset(asset); err == nil {
		return dat, err
	}
	return defaultRom, nil
}
