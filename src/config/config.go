package config

import "github.com/markel1974/c64emu/src/signals"

type Cartridge struct {
	Path string
	Kind string
}

type Config struct {
	Cartridges                []Cartridge
	DisableCartridgeAutostart bool
	changed                   *signals.Signal
	//JoystickSwap              bool
}

func New() *Config {
	return &Config{
		Cartridges:                nil,
		DisableCartridgeAutostart: false,
		changed:                   signals.NewSignal(),
		//JoystickSwap:              true,
	}
}

func (p *Config) Bind(changed func()) {
	p.changed.Bind(changed)
}

//func (p *Config) Emul1541Proc() bool {
//TODO IMPLEMENT
//	return true
//return false
//}

func (p *Config) GetDrivePath(i int) string {
	//TODO IMPLEMENT
	return "."
	//return false
}

//func (p *Config) GetJoystickSwap() bool {
//	return p.JoystickSwap
//}

func (p *Config) AddCartridge(kind string, path string) {
	p.Cartridges = append(p.Cartridges, Cartridge{Kind: kind, Path: path})
}

func (p *Config) GetCartridges() []Cartridge {
	return p.Cartridges
}

func (p *Config) SetDisableCartridgeAutostart(v bool) {
	p.DisableCartridgeAutostart = v
}

func (p *Config) GetDisableCartridgeAutostart() bool {
	return p.DisableCartridgeAutostart
}
