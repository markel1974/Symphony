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
	cartridges []Cartridge
	drives     []Drive
	disks      []Drive
	changed    *signals.Signal
	prg        string
	diskIndex  int
	jiffy      bool
}

func New() *Config {
	return &Config{
		cartridges: nil,
		diskIndex:  0,
		changed:    signals.NewSignal(),
		jiffy:      true,
	}
}

func (p *Config) Bind(changed func()) {
	p.changed.Bind(changed)
}

func (p *Config) AddDrive(kind string, opts string) {
	p.drives = append(p.drives, Drive{Kind: kind, Opts: opts})
}

func (p *Config) AddDisk(kind string, opts string) {
	p.disks = append(p.disks, Drive{Kind: kind, Opts: opts})
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

func (p *Config) SwitchDisk() {
	if len(p.disks) == 0 {
		return
	}
	p.diskIndex++
	driveIndex := p.diskIndex % len(p.disks)
	p.SetDriveOpt(p.disks[driveIndex].Opts, 0)
	p.changed.Emit()
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

func (p *Config) Get1541RomPath() string {
	return ""
}

func (p *Config) GetKernalRomPath() string {
	return ""
}

func (p *Config) DisableJiffy() {
	p.jiffy = false
}

func (p *Config) UseJiffy() bool {
	return p.jiffy
}
