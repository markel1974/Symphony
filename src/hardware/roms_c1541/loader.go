package roms_c1541

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
	"os"
)

type RomLoader struct {
	*component.BaseComponent
	factory references.IComponentFactory
	cfg     *config.Config
}

// NewRomLoader initializes and returns a new instance of RomLoader configured with the provided Config.
func NewRomLoader(parent references.IComponent, factory references.IComponentFactory, label int) *RomLoader {
	rl := &RomLoader{
		BaseComponent: component.NewBaseComponent(componentId, label, references.IdIROMLoaderC1541),
		factory:       factory,
		cfg:           nil,
	}
	component.Register(parent, rl)
	return rl
}

func (r *RomLoader) Setup(cfg *config.Config) error {
	r.cfg = cfg
	return nil
}

func (r *RomLoader) Reset() {
}

// Load attempts to load the ROM data from a file if a valid name is provided; otherwise, it returns embedded ROM data.
// If useJiffy is true, it returns the _jiffyRom; if not, it defaults to returning the _builtinRom.
func (r *RomLoader) Load() []byte {
	romName := r.cfg.Get1541RomPath()
	if len(romName) > 0 {
		if dat, err := os.ReadFile(romName); err == nil {
			return dat
		}
	}
	if r.cfg.UseJiffy() {
		return _jiffyRom
	}
	return _builtinRom
}
