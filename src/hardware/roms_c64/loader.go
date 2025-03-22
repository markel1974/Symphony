package roms_c64

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
	"os"
)

// BasicRomFile defines the filename for the Basic ROM.
// CharRomFile defines the filename for the Character ROM.
const (
	BasicRomFile = "Basic.rom"
	CharRomFile  = "Char.rom"
)

// RomLoader provides functionality to load ROM files, including kernal, basic, and character ROMs, with optional fallback handling.
type RomLoader struct {
	*component.BaseComponent
	cfg *config.Config
}

// NewRomLoader initializes and returns a new instance of RomLoader configured with the provided Config.
func NewRomLoader(parent references.IComponent, factory references.IComponentFactory, instance int) *RomLoader {
	rl := &RomLoader{
		BaseComponent: component.NewBaseComponent(),
		cfg:           nil,
	}
	rl.BaseComponent.Register(factory, parent, Identifier(), instance, rl, references.IdIROMLoaderC64(rl))
	return rl
}

func (r *RomLoader) Setup(cfg *config.Config) error {
	r.cfg = cfg
	return nil
}

func (r *RomLoader) Emulate() {

}

func (r *RomLoader) Reset() {

}

// LoadKernal loads the Kernal ROM, using the Jiffy ROM if Jiffy mode is enabled, otherwise defaults to the standard ROM.
func (r *RomLoader) LoadKernal() []byte {
	if r.cfg.UseJiffy() {
		return r.load(BuiltinKernalJiffyRom, r.cfg.GetKernalRomPath())
	}
	return r.load(BuiltinKernalRom, r.cfg.GetKernalRomPath())
}

// LoadBasic loads the built-in BASIC ROM if the external BASIC ROM file cannot be read. Returns the loaded ROM data.
func (r *RomLoader) LoadBasic() []byte {
	return r.load(BuiltinBasicRom, BasicRomFile)
}

// LoadChar loads the character ROM using the default built-in ROM or the specified file and returns its byte content.
func (r *RomLoader) LoadChar() []byte {
	return r.load(BuiltinCharRom, CharRomFile)
}

// load reads the specified ROM file from disk and returns its contents, or returns the default ROM if an error occurs.
func (r *RomLoader) load(defaultRom []byte, romName string) []byte {
	dat, err := os.ReadFile(romName)
	if err == nil {
		return dat
	}
	return defaultRom
}
