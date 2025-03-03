package pla

import (
	"github.com/markel1974/c64emu/src/version"
	"strconv"
)

type EmulatorId struct {
	rev uint8
	ver uint8
	key uint8
	alt uint8
}

func NewEmulatorId() *EmulatorId {
	rev, _ := strconv.Atoi(version.BuildVersion)
	ver, _ := strconv.Atoi(version.MajorVersion)
	return &EmulatorId{
		rev: uint8(rev << 4),
		ver: uint8(ver),
		key: version.AppName[0],
		alt: 0x55,
	}
}

func (s *EmulatorId) Read(addr uint16) uint8 {
	reg := addr & 0x7f
	switch reg {
	case 0x7c: // 0xdffc: revision
		return s.rev
	case 0x7d: // 0xdffd: version
		return s.ver
	case 0x7e: // 0xdffe
		return s.key
	case 0x7f: // 0xdfff alternates between $55 and $aa
		s.alt = ^s.alt
		return s.alt
	}
	return s.key
}
