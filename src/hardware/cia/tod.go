package mos6526

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

// TOD represents a Time of Day (TOD) and Alarm functionality with tracking for hours, minutes, seconds, and tenths.
type TOD struct {
	*component.BaseComponent
	tod10ths       uint8 // TOD 10ths
	todSec         uint8 // TOD sec
	todMin         uint8 // TOD min
	todHr          uint8 // TOD hr
	todHalt        bool  // TOD halted
	todDivider     int   // TOD frequency divider
	todShadow10ths int
	todShadowSec   int
	todShadowMin   int
	alm10ths       uint8 // Alarm time
	almSec         uint8 // Alarm time
	almMin         uint8 // Alarm time
	almHr          uint8 // Alarm time
}

// NewTOD creates and returns a new instance of the TOD struct with the specified ID initialized.
func NewTOD(parent references.IComponent, factory references.IComponentFactory, instance int) *TOD {
	t := &TOD{
		BaseComponent: component.NewBaseComponent(),
	}
	t.BaseComponent.Register(factory, parent, "tod", t, references.IdInternalComponent("TOD", instance))
	return t
}

// Reset reinitializes the TOD registers, alarm registers, and control flags to their default states.
func (m *TOD) Reset() {
	m.todHalt = false
	m.todDivider = 0
	m.tod10ths = 0
	m.todSec = 0
	m.todMin = 0
	m.todHr = 0
	m.todShadow10ths = -1
	m.todShadowSec = -1
	m.todShadowMin = -1
	m.alm10ths = 0
	m.almSec = 0
	m.almMin = 0
	m.almHr = 0
}

// Emulate handles the emulation process for the TOD, managing operations based on the current TOD state.
func (m *TOD) Emulate() {
	//
}

// EmulationRequired determines if emulation is necessary for the TOD. Always returns false in its current implementation.
func (m *TOD) EmulationRequired() bool {
	return false
}

// Set10ths updates the 10ths component of the time or alarm based on the `alarm` flag and provided `data`.
// If updating the time and TOD is halted, it resets the divider and unhalts TOD.
func (m *TOD) Set10ths(alarm bool, data uint8) {
	d := data & 0x0f
	if alarm {
		m.alm10ths = d
	} else {
		m.tod10ths = d
		if m.todHalt {
			m.todDivider = 0
			m.todHalt = false
		}
	}
}

// SetSec sets the seconds value for either the TOD clock or the alarm based on the alarm flag.
func (m *TOD) SetSec(alarm bool, data uint8) {
	d := data & 0x7f
	if alarm {
		m.almSec = d
	} else {
		m.todSec = d
	}
}

// SetMin sets the minute value for either the TOD clock or the alarm based on the alarm flag.
func (m *TOD) SetMin(alarm bool, data uint8) {
	d := data & 0x7f
	if alarm {
		m.almMin = d
	} else {
		m.todMin = d
	}
}

// SetHour sets the hour for either the alarm or the time of day based on the `alarm` flag and provided `data` input.
func (m *TOD) SetHour(alarm bool, data uint8) {
	d := data & 0x9f
	if alarm {
		m.almHr = d
	} else {
		m.todHr = d
		m.todHalt = true
	}
}

// GetHour retrieves the current hour value from the TOD device and freezes the shadow registers state.
func (m *TOD) GetHour() uint8 {
	v := m.todHr
	m.freeze()
	return v
}

// GetMin retrieves the current minute value from the TOD. If a shadow minute is available, it returns the shadow value.
func (m *TOD) GetMin() uint8 {
	if m.todShadowMin >= 0 {
		return uint8(m.todShadowMin)
	}
	return m.todMin
}

// GetSec returns the current seconds value from the TOD, checking the shadow register if applicable.
func (m *TOD) GetSec() uint8 {
	if m.todShadowSec >= 0 {
		return uint8(m.todShadowSec)
	}
	return m.todSec
}

// Get10ths retrieves the current 10ths value from the TOD (Time-of-Day) clock, using the shadow value if available.
func (m *TOD) Get10ths() uint8 {
	var v uint8
	if m.todShadow10ths >= 0 {
		v = uint8(m.todShadow10ths)
	} else {
		v = m.tod10ths
	}
	m.unfreeze()
	return v
}

// freeze copies the current time-of-day values (10ths, seconds, and minutes) into their respective shadow variables.
func (m *TOD) freeze() {
	m.todShadow10ths = int(m.tod10ths)
	m.todShadowSec = int(m.todSec)
	m.todShadowMin = int(m.todMin)
}

// unfreeze clears the shadow registers for tenths, seconds, and minutes, restoring real-time tracking for these values.
func (m *TOD) unfreeze() {
	m.todShadow10ths = -1
	m.todShadowSec = -1
	m.todShadowMin = -1
}

// Update processes the TOD counter, increments time fields, and checks if the alarm matches the current TOD time.
// Returns true if an alarm match occurs, otherwise false.
// The rtc parameter determines the TOD frequency divider (50Hz or 60Hz).
// Automatically handles AM/PM and 24-hour format transitions.
// Resets the divider when time fields are updated.
func (m *TOD) Update(rtc bool) bool {
	if m.todHalt {
		return false
	}
	if m.todDivider > 0 {
		m.todDivider--
		return false
	}
	// Divider (50/60 Hz)
	if rtc {
		m.todDivider = 4
	} else {
		m.todDivider = 5
	}
	// 1/10 seconds
	m.tod10ths++
	if m.tod10ths > 9 {
		m.tod10ths = 0
		lo := (m.todSec & 0x0f) + 1
		hi := m.todSec >> 4
		if lo > 9 {
			lo = 0
			hi++
		}
		if hi > 5 {
			m.todSec = 0
			lo = (m.todMin & 0x0f) + 1
			hi = m.todMin >> 4
			if lo > 9 {
				lo = 0
				hi++
			}
			if hi > 5 {
				m.todMin = 0
				lo = (m.todHr & 0x0f) + 1
				hi = (m.todHr >> 4) & 1
				// AM/PM
				m.todHr &= 0x80
				if lo > 9 {
					lo = 0
					hi++
				}
				m.todHr |= (hi << 4) | lo
				if (m.todHr & 0x1f) > 0x11 {
					m.todHr = (m.todHr & 0x80) ^ 0x80
				}
			} else {
				m.todMin = (hi << 4) | lo
			}
		} else {
			m.todSec = (hi << 4) | lo
		}
	}
	if m.tod10ths == m.alm10ths && m.todSec == m.almSec && m.todMin == m.almMin && m.todHr == m.almHr {
		return true
	}
	return false
}
