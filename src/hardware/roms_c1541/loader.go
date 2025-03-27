package roms_c1541

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
	"os"
)

// RomLoader is responsible for loading and managing ROM configurations and their associated resources.
type RomLoader struct {
	*component.BaseComponent
	cfg *config.Config
}

// NewRomLoader creates a new instance of RomLoader, registers it, and initializes its base component with the provided parameters.
func NewRomLoader(parent references.IComponent, factory references.IComponentFactory, instance int) *RomLoader {
	rl := &RomLoader{
		BaseComponent: component.NewBaseComponent(),
		cfg:           nil,
	}
	rl.BaseComponent.Register(factory, parent, Identifier(), rl, references.IdIROMLoaderC1541(rl, instance))
	return rl
}

// Setup initializes the RomLoader with the provided configuration and socket reference.
func (r *RomLoader) Setup(_ references.IROMLoaderC1541Socket, cfg *config.Config) error {
	r.cfg = cfg
	return nil
}

// Connect initializes the ROM loader's connection to required resources or systems, preparing it for further operations.
func (r *RomLoader) Connect() error {
	return nil
}

func (r *RomLoader) Internal() bool {
	return false
}

// Emulate begins the emulation process for the ROM loader, simulating its expected behavior within the system context.
func (r *RomLoader) Emulate() {

}

// EmulationRequired checks if ROM emulation is necessary for the current configuration and returns the result.
func (r *RomLoader) EmulationRequired() bool {
	return false
}

// Reset reinitializes the RomLoader to its default state, clearing any runtime modifications or temporary data.
func (r *RomLoader) Reset() {
}

// Load retrieves the ROM data based on the configuration, prioritizing Jiffy ROM, then custom ROM, and finally builtin ROM.
func (r *RomLoader) Load() []byte {
	if r.cfg.UseJiffy() {
		return _jiffyRom
	}
	if romName := r.cfg.C1541RomPath(); len(romName) > 0 {
		if dat, err := os.ReadFile(romName); err == nil {
			return dat
		}
	}
	return _builtinRom
}
