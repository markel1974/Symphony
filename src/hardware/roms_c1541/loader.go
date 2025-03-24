package roms_c1541

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
	"os"
)

// RomLoader is responsible for loading ROM files and managing their lifecycle within the application.
type RomLoader struct {
	*component.BaseComponent
	cfg *config.Config
}

// NewRomLoader initializes a new RomLoader instance and registers it with the given parent and component factory.
func NewRomLoader(parent references.IComponent, factory references.IComponentFactory, instance int) *RomLoader {
	rl := &RomLoader{
		BaseComponent: component.NewBaseComponent(),
		cfg:           nil,
	}
	rl.BaseComponent.Register(factory, parent, Identifier(), rl, references.IdIROMLoaderC1541(rl, instance))
	return rl
}

// Setup configures the RomLoader with the provided configuration instance and initializes its internal settings.
func (r *RomLoader) Setup(cfg *config.Config) error {
	r.cfg = cfg
	return nil
}

// Emulate triggers the emulation process for the ROM loader component.
func (r *RomLoader) Emulate() {

}

func (r *RomLoader) EmulationRequired() bool {
	return false
}

// Reset reinitializes the RomLoader to its default state.
func (r *RomLoader) Reset() {
}

// Load attempts to load the appropriate ROM data based on configuration settings, falling back to a default if necessary.
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
