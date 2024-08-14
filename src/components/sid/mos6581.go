package sid

import (
	"github.com/markel1974/c64emu/src/components/vic"
	"github.com/markel1974/c64emu/src/config"
)

const (
	SampleFreq      = 44100                  // Sample output frequency in Hz
	Frequency       = 985248                 // SID frequency in Hz
	CalcFreq        = 50                     // Frequency at which calc_buffer is called in Hz (should be 50Hz)
	Cycles          = Frequency / SampleFreq // # of SID clocks per sample frame
	SampleBufSize   = 0x138 * 2              // Size of buffer for sampled voice (double buffered)
	RegisterCount   = 32
	RegisterHistory = 1024
	FragFreq        = vic.ScreenFreq             // one frag per frame
	FragSize        = SampleFreq / FragFreq      // samples, not bytes
	FragInterval    = 1000 / FragFreq            // in milliseconds
	BufferFrags     = FragFreq                   // frags the in buffer
	BufferSize      = 2 * FragSize * BufferFrags // bytes, not samples
	MaxLeadAvg      = BufferFrags                // lead average count
)

type MOS6581 struct {
	regs             []uint8
	regsHistory      [][]uint8
	regsHistoryIndex uint32
	cfg              *config.Config
}

func NewMOS6581() *MOS6581 {
	s := &MOS6581{
		regs:             make([]uint8, RegisterCount),
		regsHistory:      make([][]uint8, RegisterHistory),
		regsHistoryIndex: 0,
		cfg:              nil,
	}
	for x := range s.regsHistory {
		s.regsHistory[x] = make([]uint8, RegisterCount)
	}
	return s
}

func (sid *MOS6581) Setup(cfg *config.Config) {
	sid.cfg = cfg
	sid.cfg.Bind(sid.configChanged)
}

func (sid *MOS6581) SetPotXSlot(pot uint8) {
	// PX7 PX6 PX5 PX4 PX3 PX2 PX1 PX0
	sid.regs[25] = pot
}

func (sid *MOS6581) SetPotYSlot(pot uint8) {
	//PY7 PY6 PY5 PY4 PY3 PY2 PY1 PY0
	sid.regs[26] = pot
}

func (sid *MOS6581) configChanged() {
	//TODO
}

func (sid *MOS6581) Reset() {
	for x := range sid.regs {
		sid.regs[x] = 0
	}
	for x := range sid.regsHistory {
		for y := range sid.regsHistory[x] {
			sid.regsHistory[x][y] = 0
		}
	}
	sid.regsHistoryIndex = 0
}

func (sid *MOS6581) ReadRegister(addr uint16) uint8 {
	addr = addr & 0x1f
	return sid.regs[addr]
}

func (sid *MOS6581) WriteRegister(addr uint16, data uint8) {
	addr = addr & 0x1f
	sid.regs[addr] = data
}

func (sid *MOS6581) Emulate() {
	if sid.regsHistoryIndex < RegisterHistory {
		copy(sid.regsHistory[sid.regsHistoryIndex], sid.regs)
		sid.regsHistoryIndex++
	}
}

func (sid *MOS6581) GetRegsHistory() [][]uint8 {
	return sid.regsHistory
}

func (sid *MOS6581) ResetHistoryCounter() uint32 {
	cycle := sid.regsHistoryIndex
	sid.regsHistoryIndex = 0
	return cycle
}
