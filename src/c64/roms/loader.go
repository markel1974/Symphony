package roms

import (
	"github.com/markel1974/c64emu/src/config"
	"os"
)

// BasicRomFile represents the filename for the Basic ROM file.
// CharRomFile represents the filename for the Char ROM file.
const (
	BasicRomFile = "Basic.rom"
	CharRomFile  = "Char.rom"
)

type RomLoader struct {
	cfg *config.Config
}

func NewRomLoader(cfg *config.Config) *RomLoader {
	return &RomLoader{
		cfg: cfg,
	}
}

func (r *RomLoader) LoadKernal() []byte {
	if r.cfg.UseJiffy() {
		return r.load(BuiltinKernalJiffyRom, r.cfg.GetKernalRomPath())
	}
	return r.load(BuiltinKernalRom, r.cfg.GetKernalRomPath())
}

func (r *RomLoader) LoadBasic() []byte {
	return r.load(BuiltinBasicRom, BasicRomFile)
}

func (r *RomLoader) LoadChar() []byte {
	return r.load(BuiltinCharRom, CharRomFile)
}

func (r *RomLoader) load(defaultRom []byte, romName string) []byte {
	dat, err := os.ReadFile(romName)
	if err == nil {
		return dat
	}
	return defaultRom
}
