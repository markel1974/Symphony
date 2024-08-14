package cia

type TOD struct {
	tod10ths   uint8 // TOD 10ths
	todSec     uint8 // TOD sec
	todMin     uint8 // TOD min
	todHr      uint8 // TOD hr
	todHalt    bool  // TOD halted
	todDivider int   // TOD frequency divider
	alm10ths   uint8 // Alarm time
	almSec     uint8 // Alarm time
	almMin     uint8 // Alarm time
	almHr      uint8 // Alarm time
}

func NewTOD() *TOD {
	return &TOD{
		tod10ths:   0,
		todSec:     0,
		todMin:     0,
		todHr:      0,
		todHalt:    false,
		todDivider: 0,
		alm10ths:   0,
		almSec:     0,
		almMin:     0,
		almHr:      0,
	}
}

func (m *TOD) Reset() {
	m.todHalt = false
}

func (m *TOD) Count(crA uint8) bool {
	// Decrement frequency divider
	if (m.todDivider) != 0 {
		m.todDivider--
		return false
	}
	// Reload divider according to 50/60 Hz flag
	if (crA & 0x80) != 0 {
		m.todDivider = 4
	} else {
		m.todDivider = 5
	}
	// 1/10 seconds
	m.tod10ths++
	if m.tod10ths > 9 {
		var lo, hi uint8
		m.tod10ths = 0
		// Seconds
		lo = (m.todSec & 0x0f) + 1
		hi = m.todSec >> 4
		if lo > 9 {
			lo = 0
			hi++
		}
		if hi > 5 {
			m.todSec = 0
			// Minutes
			lo = (m.todMin & 0x0f) + 1
			hi = m.todMin >> 4
			if lo > 9 {
				lo = 0
				hi++
			}
			if hi > 5 {
				m.todMin = 0
				// Hours
				lo = (m.todHr & 0x0f) + 1
				hi = (m.todHr >> 4) & 1
				// Keep AM/PM flag
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
	// Alarm time reached? Trigger interrupt if enabled
	if m.tod10ths == m.alm10ths && m.todSec == m.almSec && m.todMin == m.almMin && m.todHr == m.almHr {
		return true
		//TODO triggerInterrupt
		//m.triggerInterrupt(4)
	}
	return false
}
