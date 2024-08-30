package cia

type TOD struct {
	id             string
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

func NewTOD(id string) *TOD {
	return &TOD{
		id: id,
	}
}

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

func (m *TOD) Set10ths(alarm bool, d uint8) {
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

func (m *TOD) SetSec(alarm bool, d uint8) {
	if alarm {
		m.almSec = d
	} else {
		m.todSec = d
	}
}

func (m *TOD) SetMin(alarm bool, d uint8) {
	if alarm {
		m.almMin = d
	} else {
		m.todMin = d
	}
}

func (m *TOD) SetHour(alarm bool, d uint8) {
	if alarm {
		m.almHr = d
	} else {
		m.todHr = d
		m.todHalt = true
	}
}

func (m *TOD) GetHour() uint8 {
	v := m.todHr
	m.freeze()
	return v
}

func (m *TOD) GetMin() uint8 {
	if m.todShadowMin >= 0 {
		return uint8(m.todShadowMin)
	}
	return m.todMin
}

func (m *TOD) GetSec() uint8 {
	if m.todShadowSec >= 0 {
		return uint8(m.todShadowSec)
	}
	return m.todSec
}

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

func (m *TOD) freeze() {
	m.todShadow10ths = int(m.tod10ths)
	m.todShadowSec = int(m.todSec)
	m.todShadowMin = int(m.todMin)
}

func (m *TOD) unfreeze() {
	m.todShadow10ths = -1
	m.todShadowSec = -1
	m.todShadowMin = -1
}

func (m *TOD) Update(v uint8) bool {
	if m.todHalt {
		return false
	}
	if m.todDivider > 0 {
		m.todDivider--
		return false
	}
	// Reload divider (50/60 Hz flag)
	if v != 0 {
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
				// AM/PM flag
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
