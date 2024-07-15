package preferences

type Prefs struct {
	Cartridge                 string
	DisableCartridgeAutostart bool
	//JoystickSwap              bool
}

func NewPrefs() *Prefs {
	return &Prefs{
		Cartridge:                 "",
		DisableCartridgeAutostart: false,
		//JoystickSwap:              true,
	}
}

//func (p *Prefs) Emul1541Proc() bool {
//TODO IMPLEMENT
//	return true
//return false
//}

func (p *Prefs) GetDrivePath(i int) string {
	//TODO IMPLEMENT
	return "."
	//return false
}

//func (p *Prefs) GetJoystickSwap() bool {
//	return p.JoystickSwap
//}

func (p *Prefs) SetCartridge(c string) {
	p.Cartridge = c
}

func (p *Prefs) GetCartridge() string {
	return p.Cartridge
}

func (p *Prefs) SetDisableCartridgeAutostart(v bool) {
	p.DisableCartridgeAutostart = v
}

func (p *Prefs) GetDisableCartridgeAutostart() bool {
	return p.DisableCartridgeAutostart
}
