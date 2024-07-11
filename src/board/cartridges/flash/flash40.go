package flash

import (
	"fmt"
	"github.com/markel1974/c64emu/src/board/iboard"
	"github.com/markel1974/c64emu/src/board/quartz"
	"github.com/markel1974/c64emu/src/board/snapshot"
)

const (
	FLASH040_DUMP_VER_MAJOR = 2
	FLASH040_DUMP_VER_MINOR = 0
)

const (
	FLASH040_ERASE_MASK_SIZE = 8
)

type Flash040Type int

const (
	FLASH040_TYPE_NORMAL         = Flash040Type(iota) /* 29F040 */
	FLASH040_TYPE_B                                   /* 29F040B */
	FLASH040_TYPE_010                                 /* 29F010 */
	FLASH040_TYPE_032B_A0_1_SWAP                      /* 29F032B, A0/1 swapped */
	FLASH040_TYPE_064                                 /* Expansion S29GL064N */
	FLASH040_TYPE_NUM                                 /* This item always needs to be at the end */
)

type Flash040State int

const (
	FLASH040_STATE_READ = Flash040State(iota)
	FLASH040_STATE_MAGIC_1
	FLASH040_STATE_MAGIC_2
	FLASH040_STATE_AUTOSELECT
	FLASH040_STATE_BYTE_PROGRAM
	FLASH040_STATE_BYTE_PROGRAM_ERROR
	FLASH040_STATE_ERASE_MAGIC_1
	FLASH040_STATE_ERASE_MAGIC_2
	FLASH040_STATE_ERASE_SELECT
	FLASH040_STATE_CHIP_ERASE
	FLASH040_STATE_SECTOR_ERASE
	FLASH040_STATE_SECTOR_ERASE_TIMEOUT
	FLASH040_STATE_SECTOR_ERASE_SUSPEND
)

type Flash040Context struct {
	flashData      []uint8
	flashState     Flash040State
	flashBaseState Flash040State
	programByte    uint8
	eraseMask      [FLASH040_ERASE_MASK_SIZE]uint8
	flashDirty     int
	flashType      Flash040Type
	lastRead       uint8
	board          iboard.IBoard
	eraseAlarm     *quartz.Alarm
}

// NewFlash040Context
// func NewFlash040Context(b iboard.IBoard, kind Flash040Type, data []uint8) *Flash040Context {
func NewFlash040Context(b iboard.IBoard, kind Flash040Type, data []byte) *Flash040Context {
	f := &Flash040Context{
		board:          b,
		flashType:      kind,
		flashData:      data, //make([]uint8, 0x80000), //make([]uint8, 0x100000), //data,
		flashState:     FLASH040_STATE_READ,
		flashBaseState: FLASH040_STATE_READ,
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

func (f *Flash040Context) Shutdown() {
	f.eraseAlarm.Destroy()
}

func (f *Flash040Context) GetFlashState() Flash040State {
	return f.flashState
}

func (f *Flash040Context) Reset() {
	f.flashState = FLASH040_STATE_READ
	f.flashBaseState = FLASH040_STATE_READ
	f.programByte = 0
	f.flashClearEraseMask()
	_ = f.eraseAlarm.Unset()
}

func (f *Flash040Context) StoreInterval(start uint, end uint, data []uint8) error {
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

func (f *Flash040Context) ReadInterval(start uint, end uint) ([]byte, error) {
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

func (f *Flash040Context) Store(addr uint, data uint8) {
	clock := f.board.Cycle()
	rmwFlags := f.board.RmwFlags()
	if rmwFlags != 0 {
		clock--
		f.storeInternal(clock, addr, f.lastRead)
		//clock++
	}
	f.storeInternal(clock, addr, data)
}

func (f *Flash040Context) Read(addr uint) uint8 {
	var value uint8
	//#ifdef FLASH_DEBUG_ENABLED
	//Flash040State old_state = flash040_context->flashState;
	//#endif
	switch f.flashState {
	case FLASH040_STATE_AUTOSELECT:
		if f.flashType == FLASH040_TYPE_032B_A0_1_SWAP {
			if (addr & 0xff) < 4 {
				k := []uint{0, 2, 1, 3}
				//TODO
				//addr = "\0\2\1\3"[addr&0x3]
				addr = k[addr&0x3]
			}
		}
		if (addr & 0xff) == 0 {
			value = flashTypes[f.flashType].manufacturer_ID
		} else if (addr & 0xff) == uint(flashTypes[f.flashType].device_ID_addr) {
			value = flashTypes[f.flashType].device_ID
		} else if (addr & 0xff) == 2 {
			value = 0
		} else {
			value = f.flashData[addr]
		}
	case FLASH040_STATE_BYTE_PROGRAM_ERROR:
		value = uint8(f.flashWriteOperationStatus())
	case FLASH040_STATE_SECTOR_ERASE_SUSPEND, FLASH040_STATE_CHIP_ERASE, FLASH040_STATE_SECTOR_ERASE, FLASH040_STATE_SECTOR_ERASE_TIMEOUT:
		value = uint8(f.flashEraseOperationStatus())
	case FLASH040_STATE_READ:
		value = f.flashData[addr]
	default:
		value = f.flashData[addr] // The state doesn't reset if an read occurs during a command sequence
	}

	//#ifdef FLASH_DEBUG_ENABLED
	//if (old_state != FLASH040_STATE_READ) {
	//	FLASH_DEBUG(("Read %02x from %05x, state %i->%i", value, addr, (int)old_state, (int)flash040_context->flashState));
	//}
	//#endif

	f.lastRead = value
	return value
}

func (f *Flash040Context) Peek(addr uint32) uint8 {
	return f.flashData[addr]
}

func (f *Flash040Context) SnapshotWriteModule(s *snapshot.Snapshot, name string) error {
	m := s.NewModule(name, FLASH040_DUMP_VER_MAJOR, FLASH040_DUMP_VER_MINOR)
	state := uint8(f.flashState)
	baseState := uint8(f.flashBaseState)
	m.Add("state", state)
	m.Add("base_state", baseState)
	m.Add("programByte", f.programByte)
	m.Add("eraseMask", f.eraseMask)
	m.Add("lastRead", f.lastRead)
	return nil
}

func (f *Flash040Context) SnapshotReadModule(s *snapshot.Snapshot, name string) error {
	m := s.GetModule(name)
	if m == nil {
		return fmt.Errorf("unknown module")
	}
	if m.Major != FLASH040_DUMP_VER_MAJOR {
		return fmt.Errorf("invalid version")
	}
	v := m.Get("state")
	if t, ok := v.(Flash040State); ok {
		f.flashState = t
	}
	v = m.Get("base_state")
	if t, ok := v.(Flash040State); ok {
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
	case FLASH040_STATE_SECTOR_ERASE_TIMEOUT, FLASH040_STATE_SECTOR_ERASE, FLASH040_STATE_CHIP_ERASE:
		/* the alarm timing is not saved, just use some value for now */
		mainCpuClk := f.board.Cycle()
		_ = f.eraseAlarm.Set(mainCpuClk + uint64(flashTypes[f.flashType].erase_sector_cycles))
		//alarm_set(f.eraseAlarm, mainCpuClk+flashTypes[f.flashType].erase_sector_cycles)
	default:
	}
	return nil
}

func (f *Flash040Context) flashMagic1(addr uint) int {
	p1 := flashTypes[f.flashType].magic_1_mask
	p2 := flashTypes[f.flashType].magic_1_addr
	if v := (addr & p1) == p2; v {
		return 1
	}
	return 0
}

func (f *Flash040Context) flashMagic2(addr uint) int {
	p1 := flashTypes[f.flashType].magic_2_mask
	p2 := flashTypes[f.flashType].magic_2_addr
	if v := (addr & p1) == p2; v {
		return 1
	}
	return 0
}

func (f *Flash040Context) flashClearEraseMask() {
	for i := 0; i < FLASH040_ERASE_MASK_SIZE; i++ {
		f.eraseMask[i] = 0
	}
}

func (f *Flash040Context) flashSectorToAddr(sector uint) uint {
	sectorSize := flashTypes[f.flashType].sector_size
	return sector * sectorSize
}

func (f *Flash040Context) flashAddrToSectorNumber(addr uint) uint {
	sectorAddr := flashTypes[f.flashType].sector_mask & addr
	sectorShift := flashTypes[f.flashType].sector_shift
	return sectorAddr >> sectorShift
}

func (f *Flash040Context) flashAddSectorToEraseMask(addr uint) {
	sectorNum := f.flashAddrToSectorNumber(addr)
	f.eraseMask[sectorNum>>3] |= (uint8)(1 << (sectorNum & 0x7))
}

func (f *Flash040Context) flashEraseSector(sector uint) {
	sectorSize := flashTypes[f.flashType].sector_size
	sectorAddr := f.flashSectorToAddr(sector)
	//FLASH_DEBUG(("Erasing 0x%x - 0x%x", sector_addr, sector_addr + sector_size - 1));
	for x := uint(0); x < sectorSize; x++ {
		f.flashData[sectorAddr+x] = 0xff
	}
	//memset(&(f.flashData[sectorAddr]), 0xff, sectorSize)
	f.flashDirty = 1
}

func (f *Flash040Context) flashEraseChip() {
	//FLASH_DEBUG(("Erasing chip"));
	for x := uint(0); x < flashTypes[f.flashType].size; x++ {
		f.flashData[x] = 0xff
	}
	//memset(f.flashData, 0xff, flashTypes[f.flashType].size)
	f.flashDirty = 1
}

func (f *Flash040Context) flashProgramByte(addr uint, data uint8) int {
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

func (f *Flash040Context) flashWriteOperationStatus() int {
	mainCpuClk := f.board.Cycle()
	p1 := int((f.programByte ^ 0x80) & 0x80) //DQ7 = inverse of programmed data
	p2 := int(mainCpuClk&2) << 5             /* DQ6 = toggle bit (2 us) */
	p3 := 1 << 5                             /* DQ5 = timeout */
	return p1 | p2 | p3
}

func (f *Flash040Context) flashEraseOperationStatus() int {
	/* DQ6 = toggle bit */
	v := f.programByte
	/* toggle the toggle bit(s) */
	/* FIXME better toggle bit II emulation */
	p1 := uint8(flashTypes[f.flashType].status_toggle_bits)
	f.programByte ^= p1
	/* DQ3 = sector erase timer */
	if f.flashState != FLASH040_STATE_SECTOR_ERASE_TIMEOUT {
		v |= 0x08
	}
	return int(v)
}

func (f *Flash040Context) eraseAlarmHandler(clock uint64, offset uint64) {
	var i, j int
	var m uint8
	_ = f.eraseAlarm.Unset()
	//FLASH_DEBUG(("Erase alarm, state %i", (int)f.flashState)

	switch f.flashState {
	case FLASH040_STATE_SECTOR_ERASE_TIMEOUT:
		_ = f.eraseAlarm.Set(clock + uint64(flashTypes[f.flashType].erase_sector_cycles))
		f.flashState = FLASH040_STATE_SECTOR_ERASE
		break
	case FLASH040_STATE_SECTOR_ERASE:
		for i = 0; i < (8 * FLASH040_ERASE_MASK_SIZE); i++ {
			j = i >> 3
			m = uint8(1 << (i & 0x7))
			if (f.eraseMask[j] & m) != 0 {
				f.flashEraseSector(uint(i))
				f.eraseMask[j] &= ^m
				break
			}
		}
		m = 0
		for i = 0; i < FLASH040_ERASE_MASK_SIZE; i++ {
			m |= f.eraseMask[i]
		}
		if m != 0 {
			_ = f.eraseAlarm.Set(clock + uint64(flashTypes[f.flashType].erase_sector_cycles))
		} else {
			f.flashState = f.flashBaseState
		}
		break
	case FLASH040_STATE_CHIP_ERASE:
		f.flashEraseChip()
		f.flashState = f.flashBaseState
		break
	default:
		//FLASH_DEBUG(("Erase alarm - error, state %i unhandled!", (int)f.flashState));
		break
	}
}

func (f *Flash040Context) storeInternal(mainCpuClk uint64, addr uint, data uint8) {
	//#ifdef FLASH_DEBUG_ENABLED
	//	Flash040State old_state = flash040_context->flashState;
	//	Flash040State old_base_state = flash040_context->flashBaseState;
	//	#endif

	switch f.flashState {
	case FLASH040_STATE_READ:
		if f.flashMagic1(addr) != 0 && (data == 0xaa) {
			f.flashState = FLASH040_STATE_MAGIC_1
		}
	case FLASH040_STATE_MAGIC_1:
		if f.flashMagic2(addr) != 0 && (data == 0x55) {
			f.flashState = FLASH040_STATE_MAGIC_2
		} else {
			f.flashState = f.flashBaseState
		}
	case FLASH040_STATE_MAGIC_2:
		if f.flashMagic1(addr) != 0 {
			switch data {
			case 0x90:
				f.flashState = FLASH040_STATE_AUTOSELECT
				f.flashBaseState = FLASH040_STATE_AUTOSELECT
			case 0xf0:
				f.flashState = FLASH040_STATE_READ
				f.flashBaseState = FLASH040_STATE_READ
			case 0xa0:
				f.flashState = FLASH040_STATE_BYTE_PROGRAM
			case 0x80:
				f.flashState = FLASH040_STATE_ERASE_MAGIC_1
			default:
				f.flashState = f.flashBaseState
			}
		} else {
			f.flashState = f.flashBaseState
		}
	case FLASH040_STATE_BYTE_PROGRAM:
		if f.flashProgramByte(addr, data) != 0 {
			/* The byte program time is short enough to ignore */
			f.flashState = f.flashBaseState
		} else {
			f.flashState = FLASH040_STATE_BYTE_PROGRAM_ERROR
		}
	case FLASH040_STATE_ERASE_MAGIC_1:
		if f.flashMagic1(addr) != 0 && (data == 0xaa) {
			f.flashState = FLASH040_STATE_ERASE_MAGIC_2
		} else {
			f.flashState = f.flashBaseState
		}
	case FLASH040_STATE_ERASE_MAGIC_2:
		if f.flashMagic2(addr) != 0 && (data == 0x55) {
			f.flashState = FLASH040_STATE_ERASE_SELECT
		} else {
			f.flashState = f.flashBaseState
		}
	case FLASH040_STATE_ERASE_SELECT:
		if f.flashMagic1(addr) != 0 && (data == 0x10) {
			f.flashState = FLASH040_STATE_CHIP_ERASE
			f.programByte = 0
			_ = f.eraseAlarm.Set(mainCpuClk + uint64(flashTypes[f.flashType].erase_chip_cycles))
		} else if data == 0x30 {
			f.flashAddSectorToEraseMask(addr)
			f.programByte = 0
			f.flashState = FLASH040_STATE_SECTOR_ERASE_TIMEOUT
			_ = f.eraseAlarm.Set(mainCpuClk + uint64(flashTypes[f.flashType].erase_sector_timeout_cycles))
		} else {
			f.flashState = f.flashBaseState
		}
	case FLASH040_STATE_SECTOR_ERASE_TIMEOUT:
		if data == 0x30 {
			f.flashAddSectorToEraseMask(addr)
		} else {
			f.flashState = f.flashBaseState
			f.flashClearEraseMask()
			_ = f.eraseAlarm.Unset()
		}
	case FLASH040_STATE_SECTOR_ERASE:
		/* TODO not all models support suspending */
		if data == 0xb0 {
			f.flashState = FLASH040_STATE_SECTOR_ERASE_SUSPEND
			_ = f.eraseAlarm.Unset()
		}
	case FLASH040_STATE_SECTOR_ERASE_SUSPEND:
		if data == 0x30 {
			f.flashState = FLASH040_STATE_SECTOR_ERASE
			_ = f.eraseAlarm.Set(mainCpuClk + uint64(flashTypes[f.flashType].erase_sector_cycles))
		}
	case FLASH040_STATE_BYTE_PROGRAM_ERROR, FLASH040_STATE_AUTOSELECT:
		if f.flashMagic1(addr) != 0 && (data == 0xaa) {
			f.flashState = FLASH040_STATE_MAGIC_1
		}
		if data == 0xf0 {
			f.flashState = FLASH040_STATE_READ
			f.flashBaseState = FLASH040_STATE_READ
		}
	case FLASH040_STATE_CHIP_ERASE:
	default:
	}

	//FLASH_DEBUG(("Write %02x to %05x, state %i->%i (base state %i->%i)", byte, addr, (int)old_state, (int)f.flashState, (int)old_base_state, (int)f.flashBaseState));
}
