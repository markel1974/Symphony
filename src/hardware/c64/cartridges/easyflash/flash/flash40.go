package flash

import (
	"fmt"
	"github.com/markel1974/c64emu/src/hardware/c64/snapshot"
	"github.com/markel1974/c64emu/src/references"
)

// DumpVerMajor represents the major version of the dump format.
// DumpVerMinor represents the minor version of the dump format.
const (
	DumpVerMajor = 2
	DumpVerMinor = 0
)

// EraseMaskSize defines the size of the erase mask array used for handling flash memory operations.
const (
	EraseMaskSize = 8
)

// Kind represents an enumeration type used for categorizing or distinguishing different variants or states.
type Kind int

// KindNormal represents the Kind 29F040.
// KindB represents the Kind 29F040B.
// Kind010 represents the Kind 29F010.
// Kind032BA01Swap represents the Kind 29F032B with A0/1 swapped.
// Kind064 represents the Kind for expansion S29GL064N.
// KindNum represents the latest Kind.
const (
	KindNormal      = Kind(iota) // 29F040
	KindB                        // 29F040B
	Kind010                      // 29F010
	Kind032BA01Swap              // 29F032B, A0/1 swapped
	Kind064                      // Expansion S29GL064N
	KindNum                      // Latest
)

// State represents the operational state of a specific entity, defined as an integer-based enumeration.
type State int

// StateRead represents the state where a read operation is performed.
// StateMagic1 represents the first magic step in the operation sequence.
// StateMagic2 represents the second magic step in the operation sequence.
// StateAutoSelect represents the state where auto-select operation is executed.
// StateByteProgram represents the state for byte program operation.
// StateByteProgramError represents the state where a byte program error occurs.
// StateEraseMagic1 represents the first magic step in the erase operation sequence.
// StateEraseMagic2 represents the second magic step in the erase operation sequence.
// StateEraseSelect represents the state where erase selection is performed.
// StateChipErase represents the state where a chip erase operation is performed.
// StateSectorErase represents the state where a sector erase operation is performed.
// SectorEraseTimeout represents the state of timeout during sector erase.
// StateSectorEraseSuspend represents the state where a sector erase operation is suspended.
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

// Flash040 represents a flash memory module with various operational states, data storage, and functionalities.
type Flash040 struct {
	flashData      []uint8
	flashState     State
	flashBaseState State
	programByte    uint8
	eraseMask      [EraseMaskSize]uint8
	flashDirty     int
	flashType      Kind
	lastRead       uint8
	board          references.IC64Expansion
	eraseAlarm     references.IQuartzAlarm
}

// NewFlash040 initializes and returns a pointer to a Flash040 instance with the provided cartridge, type, and data.
// It configures the flash memory state, sets up an erase alarm, and clears the erase mask.
func NewFlash040(b references.IC64Expansion, kind Kind, data []byte) *Flash040 {
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
		f.eraseAlarmHandler()
	}
	f.eraseAlarm = b.CycleAlarm("Flash040Alarm", callback)
	return f
}

// Shutdown releases resources associated with the Flash040 instance, such as destroying the erase alarm.
func (f *Flash040) Shutdown() {
	f.eraseAlarm.Destroy()
}

// GetFlashState retrieves the current state of the flash. It returns a value of the State type representing the flash's state.
func (f *Flash040) GetFlashState() State {
	return f.flashState
}

// Reset reinitializes the Flash040 state to default values for flashState, flashBaseState, and programByte.
// Reset also clears the erase mask and unsets the erase alarm.
func (f *Flash040) Reset() {
	f.flashState = StateRead
	f.flashBaseState = StateRead
	f.programByte = 0
	f.flashClearEraseMask()
	_ = f.eraseAlarm.Unset()
}

// StoreInterval writes a slice of data starting at the specified index in the flashData array.
// Returns an error if the start or end interval exceeds the bounds of the flashData array.
func (f *Flash040) StoreInterval(start uint, data []uint8) error {
	if start >= uint(len(f.flashData)) {
		return fmt.Errorf("invalid start interval")
	}
	if start+uint(len(data)) >= uint(len(f.flashData)) {
		return fmt.Errorf("invalid end interval")
	}
	copy(f.flashData[start:], data)
	return nil
}

// ReadInterval reads and returns a slice of bytes from the flash memory within the specified start and end bounds.
// It returns an error if the interval is invalid or out of range.
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

// Store writes a byte of data to the specified address in the flash memory, with handling for read-modify-write scenarios.
func (f *Flash040) Store(addr uint, data uint8) {
	dist := uint64(0)
	if rmwFlags := f.board.RmwFlags(); rmwFlags != 0 {
		f.storeInternal(dist, addr, f.lastRead)
		dist++
	}
	f.storeInternal(dist, addr, data)
}

// Read accesses the flash memory at the specified address and returns the corresponding 8-bit data.
func (f *Flash040) Read(addr uint) uint8 {
	var value uint8
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
	f.lastRead = value
	return value
}

// Peek retrieves the byte value stored at the given memory address from the flash memory.
func (f *Flash040) Peek(addr uint32) uint8 {
	return f.flashData[addr]
}

// SnapshotWriteModule writes the state of the Flash040 module into the provided snapshot module with the specified name.
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

// SnapshotReadModule restores the module-specific state from a snapshot using the provided snapshot object and module name.
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
		_ = f.eraseAlarm.Set(uint64(flashTypes[f.flashType].eraseSectorCycles))
	default:
	}
	return nil
}

// flashMagic1 evaluates a given address using the magic1Mask and magic1Addr to determine a match, returning 1 if true.
func (f *Flash040) flashMagic1(addr uint) int {
	p1 := flashTypes[f.flashType].magic1Mask
	p2 := flashTypes[f.flashType].magic1Addr
	if v := (addr & p1) == p2; v {
		return 1
	}
	return 0
}

// flashMagic2 checks if the given address matches a specific pattern defined by `magic2Mask` and `magic2Addr` for the Flash type.
// Returns 1 if the condition is met, otherwise returns 0.
func (f *Flash040) flashMagic2(addr uint) int {
	p1 := flashTypes[f.flashType].magic2Mask
	p2 := flashTypes[f.flashType].magic2Addr
	if v := (addr & p1) == p2; v {
		return 1
	}
	return 0
}

// flashClearEraseMask resets all elements of the eraseMask array in the Flash040 struct to zero.
func (f *Flash040) flashClearEraseMask() {
	for i := 0; i < EraseMaskSize; i++ {
		f.eraseMask[i] = 0
	}
}

// flashSectorToAddr calculates the start address of the specified sector based on the sector size of the Flash type.
func (f *Flash040) flashSectorToAddr(sector uint) uint {
	sectorSize := flashTypes[f.flashType].sectorSize
	return sector * sectorSize
}

// flashAddrToSectorNumber calculates the sector number for a given flash address based on the flash type's parameters.
func (f *Flash040) flashAddrToSectorNumber(addr uint) uint {
	sectorAddr := flashTypes[f.flashType].sectorMask & addr
	sectorShift := flashTypes[f.flashType].sectorShift
	return sectorAddr >> sectorShift
}

// flashAddSectorToEraseMask adds a specific flash memory sector to the erase mask based on the provided address.
func (f *Flash040) flashAddSectorToEraseMask(addr uint) {
	sectorNum := f.flashAddrToSectorNumber(addr)
	f.eraseMask[sectorNum>>3] |= (uint8)(1 << (sectorNum & 0x7))
}

// flashEraseSector erases a specified flash sector by setting its contents to 0xFF and marking the flash as dirty.
func (f *Flash040) flashEraseSector(sector uint) {
	sectorSize := flashTypes[f.flashType].sectorSize
	sectorAddr := f.flashSectorToAddr(sector)
	for x := uint(0); x < sectorSize; x++ {
		f.flashData[sectorAddr+x] = 0xff
	}
	f.flashDirty = 1
}

// flashEraseChip erases all data on the chip by setting each byte to 0xFF and marks the flash memory as dirty.
func (f *Flash040) flashEraseChip() {
	//FLASH_DEBUG(("Erasing chip"));
	for x := uint(0); x < flashTypes[f.flashType].size; x++ {
		f.flashData[x] = 0xff
	}
	f.flashDirty = 1
}

// flashProgramByte programs a single byte of data into the flash memory at the specified address.
// Returns 1 if the programming operation succeeds, otherwise returns 0.
func (f *Flash040) flashProgramByte(addr uint, data uint8) int {
	oldData := f.flashData[addr]
	newData := oldData & data
	f.programByte = data
	f.flashData[addr] = newData
	f.flashDirty = 1
	if newData == data {
		return 1
	}
	return 0
}

// flashWriteOperationStatus determines and returns the status of a flash memory write operation as an integer bit mask.
func (f *Flash040) flashWriteOperationStatus() int {
	mainCpuClk := f.board.Cycle()
	p1 := int((f.programByte ^ 0x80) & 0x80) //DQ7 = inverse of programmed data
	p2 := int(mainCpuClk&2) << 5             //DQ6 = toggle bit (2 us)
	p3 := 1 << 5                             //DQ5 = timeout
	return p1 | p2 | p3
}

// flashEraseOperationStatus retrieves the current status of the flash erase operation, including toggle bits and sector erase state.
func (f *Flash040) flashEraseOperationStatus() int {
	// DQ6 = toggle bit
	v := f.programByte
	// toggle the toggle bit(s)
	// FIXME better toggle bit II emulation
	p1 := uint8(flashTypes[f.flashType].statusToggleBits)
	f.programByte ^= p1
	// DQ3 = sector erase timer
	if f.flashState != SectorEraseTimeout {
		v |= 0x08
	}
	return int(v)
}

// eraseAlarmHandler handles the alarm events during flash erase operations and transitions the flash state accordingly.
func (f *Flash040) eraseAlarmHandler() {
	_ = f.eraseAlarm.Unset()
	switch f.flashState {
	case SectorEraseTimeout:
		_ = f.eraseAlarm.Set(uint64(flashTypes[f.flashType].eraseSectorCycles))
		f.flashState = StateSectorErase
	case StateSectorErase:
		var i, j int
		var m uint8
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
			_ = f.eraseAlarm.Set(uint64(flashTypes[f.flashType].eraseSectorCycles))
		} else {
			f.flashState = f.flashBaseState
		}
	case StateChipErase:
		f.flashEraseChip()
		f.flashState = f.flashBaseState
	default:
	}
}

// storeInternal processes internal flash memory state transitions based on the provided address, data, and distance values.
func (f *Flash040) storeInternal(dist uint64, addr uint, data uint8) {
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
			_ = f.eraseAlarm.Set(dist + uint64(flashTypes[f.flashType].eraseChipCycles))
		} else if data == 0x30 {
			f.flashAddSectorToEraseMask(addr)
			f.programByte = 0
			f.flashState = SectorEraseTimeout
			_ = f.eraseAlarm.Set(dist + uint64(flashTypes[f.flashType].eraseSectorTimeoutCycles))
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
			_ = f.eraseAlarm.Set(dist + uint64(flashTypes[f.flashType].eraseSectorCycles))
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
}
