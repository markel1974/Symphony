package config

import "github.com/markel1974/c64emu/src/signals"

type Cartridge struct {
	Path string
	Kind string
}

type Config struct {
	cartridges                []Cartridge
	DisableCartridgeAutostart bool
	changed                   *signals.Signal
	prg                       string
	//JoystickSwap              bool
}

func New() *Config {
	return &Config{
		cartridges:                nil,
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
