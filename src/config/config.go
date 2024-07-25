package config

import "github.com/markel1974/c64emu/src/signals"

type Config struct {
	Cartridges                []string
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

func (p *Config) AddCartridge(c string) {
	p.Cartridges = append(p.Cartridges, c)
}

func (p *Config) GetCartridges() []string {
	return p.Cartridges
}

func (p *Config) SetDisableCartridgeAutostart(v bool) {
	p.DisableCartridgeAutostart = v
}

func (p *Config) GetDisableCartridgeAutostart() bool {
	return p.DisableCartridgeAutostart
}
