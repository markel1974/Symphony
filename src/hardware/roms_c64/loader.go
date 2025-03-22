package roms_c64

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
	"os"
)

// RomLoader is responsible for loading and managing ROM files within the application context.
type RomLoader struct {
	*component.BaseComponent
	cfg *config.Config
}

// NewRomLoader initializes and registers a new RomLoader component with the given parent, factory, and instance settings.
func NewRomLoader(parent references.IComponent, factory references.IComponentFactory, instance int) *RomLoader {
	rl := &RomLoader{
		BaseComponent: component.NewBaseComponent(),
		cfg:           nil,
	}
	rl.BaseComponent.Register(factory, parent, Identifier(), instance, rl, references.IdIROMLoaderC64(rl))
	return rl
}

// Setup initializes the RomLoader with the provided configuration. Returns an error if the setup fails.
func (r *RomLoader) Setup(cfg *config.Config) error {
	r.cfg = cfg
	return nil
}

// Emulate initializes and begins the emulation process for the ROM data loaded in the RomLoader instance.
func (r *RomLoader) Emulate() {

}

// Reset clears the state of the RomLoader and initializes it to default settings.
func (r *RomLoader) Reset() {

}

// LoadKernal loads the Kernal ROM, using the Jiffy variant if enabled, or the default/builtin ROM otherwise.
func (r *RomLoader) LoadKernal() []byte {
	if r.cfg.UseJiffy() {
		return _builtinKernalJiffyRom
	}
	return r.load(_builtinKernalRom, r.cfg.C64RomKernalPath())
}

// LoadBasic loads the BASIC ROM, using a default built-in ROM or a custom ROM specified in the configuration.
func (r *RomLoader) LoadBasic() []byte {
	return r.load(_builtinBasicRom, r.cfg.C64RomBasicPath())
}

// LoadChar loads the character ROM data, defaulting to a built-in ROM if no valid external path is configured.
func (r *RomLoader) LoadChar() []byte {
	return r.load(_builtinCharRom, r.cfg.C64RomCharPath())
}

// load reads a ROM file if the provided romName is valid; otherwise, it returns the defaultRom.
func (r *RomLoader) load(defaultRom []byte, romName string) []byte {
	if len(romName) == 0 {
		return defaultRom
	}
	if dat, err := os.ReadFile(romName); err == nil {
		return dat
	}
	return defaultRom
}
