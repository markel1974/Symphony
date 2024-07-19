package flash

import (
	"fmt"
	"github.com/markel1974/c64emu/src/board/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/board/quartz"
	"github.com/markel1974/c64emu/src/board/snapshot"
)

const (
	DumpVerMajor = 2
	DumpVerMinor = 0
)

const (
	EraseMaskSize = 8
)

type Kind int

const (
	KindNormal      = Kind(iota) // 29F040
	KindB                        // 29F040B
	Kind010                      // 29F010
	Kind032BA01Swap              // 29F032B, A0/1 swapped
	Kind064                      // Expansion S29GL064N
	KindNum                      // Latest
)

type State int

const (
	StateRead = State(iota)
	StateMagic1
	StateMagic2
	StateAutoSelect
	StateByteProgram
	StateByteProgramError
	StateEraseMagic1
	StateEraseMagic2
	StateEraseSelect
	StateChipErase
	StateSectorErase
	SectorEraseTimeout
	StateSectorEraseSuspend
)

type Flash040 struct {
	flashData      []uint8
	flashState     State
	flashBaseState State
	programByte    uint8
	eraseMask      [EraseMaskSize]uint8
	flashDirty     int
	flashType      Kind
	lastRead       uint8
	board          icartridge.IExpansion
	eraseAlarm     *quartz.Alarm
}

func NewFlash040(b icartridge.IExpansion, kind Kind, data []byte) *Flash040 {
	f := &Flash040{
		board:          b,
		flashType:      kind,
		flashData:      data,
		flashState:     StateRead,
		flashBaseState: StateRead,
		programByte:    0,
		flashDirty:     0,
	}
	f.flashClearEraseMask()
	callback := func(clock uint64, offset uint64) {
		f.eraseAlarmHandler(clock, offset)
	}
	f.eraseAlarm = b.CreateAlarm("Flash040Alarm", callback)
	return f
}

func (f *Flash040) Shutdown() {
	f.eraseAlarm.Destroy()
}

func (f *Flash040) GetFlashState() State {
	return f.flashState
}

func (f *Flash040) Reset() {
	f.flashState = StateRead
	f.flashBaseState = StateRead
	f.programByte = 0
	f.flashClearEraseMask()
	_ = f.eraseAlarm.Unset()
}

func (f *Flash040) StoreInterval(start uint, end uint, data []uint8) error {
	if start >= end {
		return fmt.Errorf("invalid interval")
	}
	if start >= uint(len(f.flashData)) {
		return fmt.Errorf("invalid start interval")
	}
	if end >= uint(len(f.flashData)) {
		return fmt.Errorf("invalid end interval")
	}
	copy(f.flashData[start:end], data)
	return nil
}

func (f *Flash040) ReadInterval(start uint, end uint) ([]byte, error) {
	if start >= end {
		return nil, fmt.Errorf("invalid interval")
	}
	if start >= uint(len(f.flashData)) {
		return nil, fmt.Errorf("invalid start interval")
	}
	if end >= uint(len(f.flashData)) {
		return nil, fmt.Errorf("invalid end interval")
	}
	ret := make([]uint8, end-start)
	copy(ret, f.flashData[start:end])
	return ret, nil
}

func (f *Flash040) Store(addr uint, data uint8) {
	clock := f.board.Cycle()
	rmwFlags := f.board.RmwFlags()
	if rmwFlags != 0 {
		clock--
		f.storeInternal(clock, addr, f.lastRead)
		//clock++
	}
	f.storeInternal(clock, addr, data)
}

func (f *Flash040) Read(addr uint) uint8 {
	var value uint8
	//#ifdef FLASH_DEBUG_ENABLED
	//State old_state = flash040_context->flashState;
	//#endif
	switch f.flashState {
	case StateAutoSelect:
		if f.flashType == Kind032BA01Swap {
			if (addr & 0xff) < 4 {
				k := []uint{0, 2, 1, 3}
				//TODO
				//addr = "\0\2\1\3"[addr&0x3]
				addr = k[addr&0x3]
			}
		}
		if (addr & 0xff) == 0 {
			value = flashTypes[f.flashType].manufacturerId
		} else if (addr & 0xff) == uint(flashTypes[f.flashType].deviceIdAddr) {
			value = flashTypes[f.flashType].deviceId
		} else if (addr & 0xff) == 2 {
			value = 0
		} else {
			value = f.flashData[addr]
		}
	case StateByteProgramError:
		value = uint8(f.flashWriteOperationStatus())
	case StateSectorEraseSuspend, StateChipErase, StateSectorErase, SectorEraseTimeout:
		value = uint8(f.flashEraseOperationStatus())
	case StateRead:
		value = f.flashData[addr]
	default:
		value = f.flashData[addr] // The state doesn't reset if an read occurs during a command sequence
	}

	//#ifdef FLASH_DEBUG_ENABLED
	//if (old_state != StateRead) {
	//	FLASH_DEBUG(("Read %02x from %05x, state %i->%i", value, addr, (int)old_state, (int)flash040_context->flashState));
	//}
	//#endif

	f.lastRead = value
	return value
}

func (f *Flash040) Peek(addr uint32) uint8 {
	return f.flashData[addr]
}

func (f *Flash040) SnapshotWriteModule(s *snapshot.Snapshot, name string) error {
	m := s.NewModule(name, DumpVerMajor, DumpVerMinor)
	state := uint8(f.flashState)
	baseState := uint8(f.flashBaseState)
	m.Add("state", state)
	m.Add("base_state", baseState)
	m.Add("programByte", f.programByte)
	m.Add("eraseMask", f.eraseMask)
	m.Add("lastRead", f.lastRead)
	return nil
}

func (f *Flash040) SnapshotReadModule(s *snapshot.Snapshot, name string) error {
	m := s.GetModule(name)
	if m == nil {
		return fmt.Errorf("unknown module")
	}
	if m.Major != DumpVerMajor {
		return fmt.Errorf("invalid version")
	}
	v := m.Get("state")
	if t, ok := v.(State); ok {
		f.flashState = t
	}
	v = m.Get("base_state")
	if t, ok := v.(State); ok {
		f.flashBaseState = t
	}
	v = m.Get("programByte")
	if t, ok := v.(uint8); ok {
		f.programByte = t
	}
	v = m.Get("eraseMask")
	if t, ok := v.([8]uint8); ok {
		f.eraseMask = t
	}
	v = m.Get("lastRead")
	if t, ok := v.(uint8); ok {
		f.lastRead = t
	}
	/* Restore alarm if needed */
	switch f.flashState {
	case SectorEraseTimeout, StateSectorErase, StateChipErase:
		/* the alarm timing is not saved, just use some value for now */
		mainCpuClk := f.board.Cycle()
		_ = f.eraseAlarm.Set(mainCpuClk + uint64(flashTypes[f.flashType].eraseSectorCycles))
		//alarm_set(f.eraseAlarm, mainCpuClk+flashTypes[f.flashType].eraseSectorCycles)
	default:
	}
	return nil
}

func (f *Flash040) flashMagic1(addr uint) int {
	p1 := flashTypes[f.flashType].magic1Mask
	p2 := flashTypes[f.flashType].magic1Addr
	if v := (addr & p1) == p2; v {
		return 1
	}
	return 0
}

func (f *Flash040) flashMagic2(addr uint) int {
	p1 := flashTypes[f.flashType].magic2Mask
	p2 := flashTypes[f.flashType].magic2Addr
	if v := (addr & p1) == p2; v {
		return 1
	}
	return 0
}

func (f *Flash040) flashClearEraseMask() {
	for i := 0; i < EraseMaskSize; i++ {
		f.eraseMask[i] = 0
	}
}

func (f *Flash040) flashSectorToAddr(sector uint) uint {
	sectorSize := flashTypes[f.flashType].sectorSize
	return sector * sectorSize
}

func (f *Flash040) flashAddrToSectorNumber(addr uint) uint {
	sectorAddr := flashTypes[f.flashType].sectorMask & addr
	sectorShift := flashTypes[f.flashType].sectorShift
	return sectorAddr >> sectorShift
}

func (f *Flash040) flashAddSectorToEraseMask(addr uint) {
	sectorNum := f.flashAddrToSectorNumber(addr)
	f.eraseMask[sectorNum>>3] |= (uint8)(1 << (sectorNum & 0x7))
}

func (f *Flash040) flashEraseSector(sector uint) {
	sectorSize := flashTypes[f.flashType].sectorSize
	sectorAddr := f.flashSectorToAddr(sector)
	//FLASH_DEBUG(("Erasing 0x%x - 0x%x", sector_addr, sector_addr + sectorSize - 1));
	for x := uint(0); x < sectorSize; x++ {
		f.flashData[sectorAddr+x] = 0xff
	}
	f.flashDirty = 1
}

func (f *Flash040) flashEraseChip() {
	//FLASH_DEBUG(("Erasing chip"));
	for x := uint(0); x < flashTypes[f.flashType].size; x++ {
		f.flashData[x] = 0xff
	}
	f.flashDirty = 1
}

func (f *Flash040) flashProgramByte(addr uint, data uint8) int {
	oldData := f.flashData[addr]
	newData := oldData & data
	//FLASH_DEBUG(("Programming 0x%05x with 0x%02x (%02x->%02x)", addr, byte, old_data, old_data & byte));
	f.programByte = data
	f.flashData[addr] = newData
	f.flashDirty = 1
	if newData == data {
		return 1
	}
	return 0
}

func (f *Flash040) flashWriteOperationStatus() int {
	mainCpuClk := f.board.Cycle()
	p1 := int((f.programByte ^ 0x80) & 0x80) //DQ7 = inverse of programmed data
	p2 := int(mainCpuClk&2) << 5             /* DQ6 = toggle bit (2 us) */
	p3 := 1 << 5                             /* DQ5 = timeout */
	return p1 | p2 | p3
}

func (f *Flash040) flashEraseOperationStatus() int {
	/* DQ6 = toggle bit */
	v := f.programByte
	/* toggle the toggle bit(s) */
	/* FIXME better toggle bit II emulation */
	p1 := uint8(flashTypes[f.flashType].statusToggleBits)
	f.programByte ^= p1
	/* DQ3 = sector erase timer */
	if f.flashState != SectorEraseTimeout {
		v |= 0x08
	}
	return int(v)
}

func (f *Flash040) eraseAlarmHandler(clock uint64, offset uint64) {
	var i, j int
	var m uint8
	_ = f.eraseAlarm.Unset()
	//FLASH_DEBUG(("Erase alarm, state %i", (int)f.flashState)

	switch f.flashState {
	case SectorEraseTimeout:
		_ = f.eraseAlarm.Set(clock + uint64(flashTypes[f.flashType].eraseSectorCycles))
		f.flashState = StateSectorErase
		break
	case StateSectorErase:
		for i = 0; i < (8 * EraseMaskSize); i++ {
			j = i >> 3
			m = uint8(1 << (i & 0x7))
			if (f.eraseMask[j] & m) != 0 {
				f.flashEraseSector(uint(i))
				f.eraseMask[j] &= ^m
				break
			}
		}
		m = 0
		for i = 0; i < EraseMaskSize; i++ {
			m |= f.eraseMask[i]
		}
		if m != 0 {
			_ = f.eraseAlarm.Set(clock + uint64(flashTypes[f.flashType].eraseSectorCycles))
		} else {
			f.flashState = f.flashBaseState
		}
		break
	case StateChipErase:
		f.flashEraseChip()
		f.flashState = f.flashBaseState
		break
	default:
		//FLASH_DEBUG(("Erase alarm - error, state %i unhandled!", (int)f.flashState));
		break
	}
}

func (f *Flash040) storeInternal(mainCpuClk uint64, addr uint, data uint8) {
	//#ifdef FLASH_DEBUG_ENABLED
	//	State old_state = flash040_context->flashState;
	//	State old_base_state = flash040_context->flashBaseState;
	//	#endif

	switch f.flashState {
	case StateRead:
		if f.flashMagic1(addr) != 0 && (data == 0xaa) {
			f.flashState = StateMagic1
		}
	case StateMagic1:
		if f.flashMagic2(addr) != 0 && (data == 0x55) {
			f.flashState = StateMagic2
		} else {
			f.flashState = f.flashBaseState
		}
	case StateMagic2:
		if f.flashMagic1(addr) != 0 {
			switch data {
			case 0x90:
				f.flashState = StateAutoSelect
				f.flashBaseState = StateAutoSelect
			case 0xf0:
				f.flashState = StateRead
				f.flashBaseState = StateRead
			case 0xa0:
				f.flashState = StateByteProgram
			case 0x80:
				f.flashState = StateEraseMagic1
			default:
				f.flashState = f.flashBaseState
			}
		} else {
			f.flashState = f.flashBaseState
		}
	case StateByteProgram:
		if f.flashProgramByte(addr, data) != 0 {
			/* The byte program time is short enough to ignore */
			f.flashState = f.flashBaseState
		} else {
			f.flashState = StateByteProgramError
		}
	case StateEraseMagic1:
		if f.flashMagic1(addr) != 0 && (data == 0xaa) {
			f.flashState = StateEraseMagic2
		} else {
			f.flashState = f.flashBaseState
		}
	case StateEraseMagic2:
		if f.flashMagic2(addr) != 0 && (data == 0x55) {
			f.flashState = StateEraseSelect
		} else {
			f.flashState = f.flashBaseState
		}
	case StateEraseSelect:
		if f.flashMagic1(addr) != 0 && (data == 0x10) {
			f.flashState = StateChipErase
			f.programByte = 0
			_ = f.eraseAlarm.Set(mainCpuClk + uint64(flashTypes[f.flashType].eraseChipCycles))
		} else if data == 0x30 {
			f.flashAddSectorToEraseMask(addr)
			f.programByte = 0
			f.flashState = SectorEraseTimeout
			_ = f.eraseAlarm.Set(mainCpuClk + uint64(flashTypes[f.flashType].eraseSectorTimeoutCycles))
		} else {
			f.flashState = f.flashBaseState
		}
	case SectorEraseTimeout:
		if data == 0x30 {
			f.flashAddSectorToEraseMask(addr)
		} else {
			f.flashState = f.flashBaseState
			f.flashClearEraseMask()
			_ = f.eraseAlarm.Unset()
		}
	case StateSectorErase:
		// TODO not all models support suspending
		if data == 0xb0 {
			f.flashState = StateSectorEraseSuspend
			_ = f.eraseAlarm.Unset()
		}
	case StateSectorEraseSuspend:
		if data == 0x30 {
			f.flashState = StateSectorErase
			_ = f.eraseAlarm.Set(mainCpuClk + uint64(flashTypes[f.flashType].eraseSectorCycles))
		}
	case StateByteProgramError, StateAutoSelect:
		if f.flashMagic1(addr) != 0 && (data == 0xaa) {
			f.flashState = StateMagic1
		}
		if data == 0xf0 {
			f.flashState = StateRead
			f.flashBaseState = StateRead
		}
	case StateChipErase:
	default:
	}

	//FLASH_DEBUG(("Write %02x to %05x, state %i->%i (base state %i->%i)", byte, addr, (int)old_state, (int)f.flashState, (int)old_base_state, (int)f.flashBaseState));
}
