package preferences

type Prefs struct {
	Cartridges                []string
	DisableCartridgeAutostart bool
	//JoystickSwap              bool
}

func NewPrefs() *Prefs {
	return &Prefs{
		Cartridges:                nil,
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

func (p *Prefs) AddCartridge(c string) {
	p.Cartridges = append(p.Cartridges, c)
}

func (p *Prefs) GetCartridges() []string {
	return p.Cartridges
}

func (p *Prefs) SetDisableCartridgeAutostart(v bool) {
	p.DisableCartridgeAutostart = v
}

func (p *Prefs) GetDisableCartridgeAutostart() bool {
	return p.DisableCartridgeAutostart
}
