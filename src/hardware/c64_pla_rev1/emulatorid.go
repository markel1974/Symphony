package c64_pla_rev1

import (
	"github.com/markel1974/c64emu/src/version"
	"strconv"
)

// EmulatorId represents an identifier for the emulator, including revision, version, key, and alternate values.
type EmulatorId struct {
	revision uint8
	version  uint8
	key      uint8
	alt      uint8
}

// NewEmulatorId creates and initializes an EmulatorId instance using version information and application name.
func NewEmulatorId() *EmulatorId {
	rev, _ := strconv.Atoi(version.BuildVersion)
	ver, _ := strconv.Atoi(version.MajorVersion)
	return &EmulatorId{
		revision: uint8(rev << 4),
		version:  uint8(ver),
		key:      version.AppName[0],
		alt:      0x55,
	}
}

// Read returns the specific byte value based on the provided address by mapping it to a corresponding register.
// The returned value depends on the register identified by the address, which is masked with 0x7f.
// For register 0x7f, the return alternates between 0x55 and 0xaa on consecutive accesses.
func (s *EmulatorId) Read(addr uint16) uint8 {
	reg := addr & 0x7f
	switch reg {
	case 0x7c: // 0xdffc: revision
		return s.revision
	case 0x7d: // 0xdffd: version
		return s.version
	case 0x7e: // 0xdffe
		return s.key
	case 0x7f: // 0xdfff alternates between $55 and $aa
		s.alt = ^s.alt
		return s.alt
	}
	return s.key
}
