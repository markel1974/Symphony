package config

import "github.com/markel1974/c64emu/src/signals"

type Cartridge struct {
	Kind string
	Path string
}

type Drive struct {
	Kind string
	Opts string
}

type Config struct {
	cartridges                []Cartridge
	drives                    []Drive
	DisableCartridgeAutostart bool
	changed                   *signals.Signal
	prg                       string
	driveIndex                int
}

func New() *Config {
	return &Config{
		cartridges:                nil,
		DisableCartridgeAutostart: false,
		driveIndex:                0,
		changed:                   signals.NewSignal(),
	}
}

func (p *Config) Bind(changed func()) {
	p.changed.Bind(changed)
}

func (p *Config) AddDrive(kind string, opts string) {
	p.drives = append(p.drives, Drive{Kind: kind, Opts: opts})
}

func (p *Config) GetDrives() []Drive {
	return p.drives
}

func (p *Config) GetDrivesOpt(id uint8) (string, bool) {
	if int(id) < len(p.drives) {
		return p.drives[id].Opts, true
	}
	return "", false
}

func (p *Config) SetDriveOpt(opt string, id uint8) bool {
	if int(id) < len(p.drives) {
		p.drives[id].Opts = opt
		p.changed.Emit()
		return true
	}
	return false
}

func (p *Config) SetPrg(prg string) {
	p.prg = prg
}

func (p *Config) GetPrg() string {
	return p.prg
}

func (p *Config) AddCartridge(kind string, path string) {
	p.cartridges = append(p.cartridges, Cartridge{Kind: kind, Path: path})
}

func (p *Config) GetCartridges() []Cartridge {
	return p.cartridges
}

func (p *Config) SetDisableCartridgeAutostart(v bool) {
	p.DisableCartridgeAutostart = v
}

func (p *Config) GetDisableCartridgeAutostart() bool {
	return p.DisableCartridgeAutostart
}

func (p *Config) Get1541RomPath() string {
	return ""
}

func (p *Config) GetKernalRomPath() string {
	return ""
}

func (p *Config) UseJiffy() bool {
	return true
}
