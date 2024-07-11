package iec

import "github.com/markel1974/c64emu/src/preferences"

type C1541Model struct {
	active       bool
	deviceNumber int
}

func NewC1541Model(deviceNumber int) *C1541Model {
	return &C1541Model{deviceNumber: deviceNumber, active: false}
}

func (m *C1541Model) GetDeviceNumber() int {
	return m.deviceNumber
}

func (m *C1541Model) NewPrefs(prefs *preferences.Prefs) {

}

func (m *C1541Model) IsActive() bool {
	return m.active
}

func (m *C1541Model) AtnStateChanged(state bool) {

}
