package mos6581

// Costanti per gli indici dei registri del SID (all'interno della slice 'registers').
const (
	freqLO1 = 0
	freqHI1 = 1
	pwLO1   = 2
	pwHI1   = 3
	cr1     = 4 // Control Register - Voice 1
	ad1     = 5 // Attack/Decay - Voice 1
	sr1     = 6 // Sustain/Release - Voice 1
	freqLO2 = 7
	freqHI2 = 8
	pwLO2   = 9
	pwHI2   = 10
	cr2     = 11
	ad2     = 12
	sr2     = 13
	freqLO3 = 14
	freqHI3 = 15
	pwLO3   = 16
	pwHI3   = 17
	cr3     = 18
	ad3     = 19
	sr3     = 20
	fcLO    = 21 // Filter Cutoff Low
	fcHI    = 22
	resFilt = 23 // Resonance and Filter Control
	modeVol = 24 // Mode and Volume
	potX    = 25 // Potentiometer X (read-only)
	potY    = 26 // Potentiometer Y (read-only)
	osc3    = 27 // Oscillator 3 / Random Number (read-only)
	env3    = 28 // Envelope 3 (read-only)
)

type SidReflect struct {
	sid *SID
}

func NewSidReflect(s *SID) *SidReflect {
	r := &SidReflect{
		sid: s,
	}
	//s.AddCommand("sumRegister", "somma di due registri", r.sumTest)
	// Voice 1
	s.AddProperty("freqLO1", "Voice 1: Frequency Low Byte (0xD400)", false, r.getFreqLO1, r.setFreqLO1)
	s.AddProperty("freqHI1", "Voice 1: Frequency High Byte (0xD401)", false, r.getFreqHI1, r.setFreqHI1)
	s.AddProperty("pwLO1", "Voice 1: Pulse Width Low Byte (0xD402)", false, r.getPwLO1, r.setPwLO1)
	s.AddProperty("pwHI1", "Voice 1: Pulse Width High Byte (0xD403)", false, r.getPwHI1, r.setPwHI1)
	s.AddProperty("cr1", "Voice 1: Control Register (0xD404) - Gate, Sync, Ring, Test, Waveform", false, r.getCr1, r.setCr1)
	s.AddProperty("ad1", "Voice 1: Attack/Decay (0xD405) - Attack Rate, Decay Rate", false, r.getAd1, r.setAd1)
	s.AddProperty("sr1", "Voice 1: Sustain/Release (0xD406) - Sustain Level, Release Rate", false, r.getSr1, r.setSr1)
	// Voice 2
	s.AddProperty("freqLO2", "Voice 2: Frequency Low Byte (0xD407)", false, r.getFreqLO2, r.setFreqLO2)
	s.AddProperty("freqHI2", "Voice 2: Frequency High Byte (0xD408)", false, r.getFreqHI2, r.setFreqHI2)
	s.AddProperty("pwLO2", "Voice 2: Pulse Width Low Byte (0xD409)", false, r.getPwLO2, r.setPwLO2)
	s.AddProperty("pwHI2", "Voice 2: Pulse Width High Byte (0xD40A)", false, r.getPwHI2, r.setPwHI2)
	s.AddProperty("cr2", "Voice 2: Control Register (0xD40B) - Gate, Sync, Ring, Test, Waveform", false, r.getCr2, r.setCr2)
	s.AddProperty("ad2", "Voice 2: Attack/Decay (0xD40C) - Attack Rate, Decay Rate", false, r.getAd2, r.setAd2)
	s.AddProperty("sr2", "Voice 2: Sustain/Release (0xD40D) - Sustain Level, Release Rate", false, r.getSr2, r.setSr2)
	// Voice 3
	s.AddProperty("freqLO3", "Voice 3: Frequency Low Byte (0xD40E)", false, r.getFreqLO3, r.setFreqLO3)
	s.AddProperty("freqHI3", "Voice 3: Frequency High Byte (0xD40F)", false, r.getFreqHI3, r.setFreqHI3)
	s.AddProperty("pwLO3", "Voice 3: Pulse Width Low Byte (0xD410)", false, r.getPwLO3, r.setPwLO3)
	s.AddProperty("pwHI3", "Voice 3: Pulse Width High Byte (0xD411)", false, r.getPwHI3, r.setPwHI3)
	s.AddProperty("cr3", "Voice 3: Control Register (0xD412) - Gate, Sync, Ring, Test, Waveform", false, r.getCr3, r.setCr3)
	s.AddProperty("ad3", "Voice 3: Attack/Decay (0xD413) - Attack Rate, Decay Rate", false, r.getAd3, r.setAd3)
	s.AddProperty("sr3", "Voice 3: Sustain/Release (0xD414) - Sustain Level, Release Rate", false, r.getSr3, r.setSr3)
	// Filtro
	s.AddProperty("fcLO", "Filter Cutoff Frequency Low Byte (0xD415)", false, r.getFcLO, r.setFcLO)
	s.AddProperty("fcHI", "Filter Cutoff Frequency High Byte (0xD416)", false, r.getFcHI, r.setFcHI)
	s.AddProperty("resFilt", "Filter Resonance and Control Register (0xD417) - Resonance, Filter Mode, External Input", false, r.getResFilt, r.setResFilt)
	s.AddProperty("modeVol", "SID Mode and Volume Register (0xD418) - Filter Mode, Voice 3 Mute, Volume", false, r.getModeVol, r.setModeVol)
	// POTX, POTY, OSC3, ENV3 (sola lettura)
	s.AddProperty("potX", "POTX Register (0xD419) - Paddle X Input", false, r.getPotX, r.setPotX)
	s.AddProperty("potY", "POTY Register (0xD41A) - Paddle Y Input", false, r.getPotY, r.setPotY)
	s.AddProperty("osc3", "OSC3 Register (0xD41B) - Oscillator 3 Value", false, r.getOsc3, r.setOsc3)
	s.AddProperty("env3", "ENV3 Register (0xD41C) - Envelope 3 Value", false, r.getEnv3, r.setEnv3)
	return r
}

//func (s *SidReflect) sumTest(a int) int {
//	return int(s.sid.registers[0]) + int(s.sid.registers[1]) + a
//}

func (s *SidReflect) getFreqLO1() uint8        { return s.sid.registers[freqLO1] }
func (s *SidReflect) setFreqLO1(v uint8) error { s.sid.registers[freqLO1] = v; return nil }

func (s *SidReflect) getFreqHI1() uint8        { return s.sid.registers[freqHI1] }
func (s *SidReflect) setFreqHI1(v uint8) error { s.sid.registers[freqHI1] = v; return nil }

func (s *SidReflect) getPwLO1() uint8        { return s.sid.registers[pwLO1] }
func (s *SidReflect) setPwLO1(v uint8) error { s.sid.registers[pwLO1] = v; return nil }

func (s *SidReflect) getPwHI1() uint8        { return s.sid.registers[pwHI1] }
func (s *SidReflect) setPwHI1(v uint8) error { s.sid.registers[pwHI1] = v; return nil }

func (s *SidReflect) getCr1() uint8        { return s.sid.registers[cr1] }
func (s *SidReflect) setCr1(v uint8) error { s.sid.registers[cr1] = v; return nil }

func (s *SidReflect) getAd1() uint8        { return s.sid.registers[ad1] }
func (s *SidReflect) setAd1(v uint8) error { s.sid.registers[ad1] = v; return nil }

func (s *SidReflect) getSr1() uint8        { return s.sid.registers[sr1] }
func (s *SidReflect) setSr1(v uint8) error { s.sid.registers[sr1] = v; return nil }

func (s *SidReflect) getFreqLO2() uint8        { return s.sid.registers[freqLO2] }
func (s *SidReflect) setFreqLO2(v uint8) error { s.sid.registers[freqLO2] = v; return nil }

func (s *SidReflect) getFreqHI2() uint8        { return s.sid.registers[freqHI2] }
func (s *SidReflect) setFreqHI2(v uint8) error { s.sid.registers[freqHI2] = v; return nil }

func (s *SidReflect) getPwLO2() uint8        { return s.sid.registers[pwLO2] }
func (s *SidReflect) setPwLO2(v uint8) error { s.sid.registers[pwLO2] = v; return nil }

func (s *SidReflect) getPwHI2() uint8        { return s.sid.registers[pwHI2] }
func (s *SidReflect) setPwHI2(v uint8) error { s.sid.registers[pwHI2] = v; return nil }

func (s *SidReflect) getCr2() uint8        { return s.sid.registers[cr2] }
func (s *SidReflect) setCr2(v uint8) error { s.sid.registers[cr2] = v; return nil }

func (s *SidReflect) getAd2() uint8        { return s.sid.registers[ad2] }
func (s *SidReflect) setAd2(v uint8) error { s.sid.registers[ad2] = v; return nil }

func (s *SidReflect) getSr2() uint8        { return s.sid.registers[sr2] }
func (s *SidReflect) setSr2(v uint8) error { s.sid.registers[sr2] = v; return nil }

func (s *SidReflect) getFreqLO3() uint8        { return s.sid.registers[freqLO3] }
func (s *SidReflect) setFreqLO3(v uint8) error { s.sid.registers[freqLO3] = v; return nil }

func (s *SidReflect) getFreqHI3() uint8        { return s.sid.registers[freqHI3] }
func (s *SidReflect) setFreqHI3(v uint8) error { s.sid.registers[freqHI3] = v; return nil }

func (s *SidReflect) getPwLO3() uint8        { return s.sid.registers[pwLO3] }
func (s *SidReflect) setPwLO3(v uint8) error { s.sid.registers[pwLO3] = v; return nil }

func (s *SidReflect) getPwHI3() uint8        { return s.sid.registers[pwHI3] }
func (s *SidReflect) setPwHI3(v uint8) error { s.sid.registers[pwHI3] = v; return nil }

func (s *SidReflect) getCr3() uint8        { return s.sid.registers[cr3] }
func (s *SidReflect) setCr3(v uint8) error { s.sid.registers[cr3] = v; return nil }

func (s *SidReflect) getAd3() uint8        { return s.sid.registers[ad3] }
func (s *SidReflect) setAd3(v uint8) error { s.sid.registers[ad3] = v; return nil }

func (s *SidReflect) getSr3() uint8        { return s.sid.registers[sr3] }
func (s *SidReflect) setSr3(v uint8) error { s.sid.registers[sr3] = v; return nil }

func (s *SidReflect) getFcLO() uint8        { return s.sid.registers[fcLO] }
func (s *SidReflect) setFcLO(v uint8) error { s.sid.registers[fcLO] = v; return nil }

func (s *SidReflect) getFcHI() uint8        { return s.sid.registers[fcHI] }
func (s *SidReflect) setFcHI(v uint8) error { s.sid.registers[fcHI] = v; return nil }

func (s *SidReflect) getResFilt() uint8        { return s.sid.registers[resFilt] }
func (s *SidReflect) setResFilt(v uint8) error { s.sid.registers[resFilt] = v; return nil }

func (s *SidReflect) getModeVol() uint8        { return s.sid.registers[modeVol] }
func (s *SidReflect) setModeVol(v uint8) error { s.sid.registers[modeVol] = v; return nil }

func (s *SidReflect) getPotX() uint8        { return s.sid.registers[potX] }
func (s *SidReflect) setPotX(v uint8) error { s.sid.registers[potX] = v; return nil }

func (s *SidReflect) getPotY() uint8        { return s.sid.registers[potY] }
func (s *SidReflect) setPotY(v uint8) error { s.sid.registers[potY] = v; return nil }

func (s *SidReflect) getOsc3() uint8        { return s.sid.registers[osc3] }
func (s *SidReflect) setOsc3(v uint8) error { s.sid.registers[osc3] = v; return nil }

func (s *SidReflect) getEnv3() uint8        { return s.sid.registers[env3] }
func (s *SidReflect) setEnv3(v uint8) error { s.sid.registers[env3] = v; return nil }
