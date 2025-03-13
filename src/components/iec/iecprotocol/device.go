package iecprotocol

import (
	"github.com/markel1974/c64emu/src/components/board"
	"github.com/markel1974/c64emu/src/components/iec/iecdevice"
)

type DeviceAdapter struct {
	*board.BaseComponent
	iecdevice.IIecProtocolDevice
	deviceNumber  uint8
	state         uint8
	flags         uint8
	primary       uint8
	secondaryPrev uint8
	secondary     uint8
	timeout       uint64
	byte          uint8
	st            [0xff]uint8
	gs            *Global
}

func NewDeviceAdapter(parent board.IComponent, suffix string, deviceNumber uint8, gs *Global, pd iecdevice.IIecProtocolDevice) *DeviceAdapter {
	v := &DeviceAdapter{
		BaseComponent:      board.NewBaseComponent("iec_protocol_device", suffix),
		IIecProtocolDevice: pd,
		deviceNumber:       deviceNumber,
		gs:                 gs,
	}
	board.Register(parent, v)
	return v
}

func (s *DeviceAdapter) Reset() {
	//
}

func (s *DeviceAdapter) GetDeviceNumber() uint8 {
	return s.deviceNumber
}

func (s *DeviceAdapter) SetFlags(v uint8) {
	s.flags = v
}

func (s *DeviceAdapter) GetFlags() uint8 {
	return s.flags
}

func (s *DeviceAdapter) SetByte(v uint8) {
	s.byte = v
}

func (s *DeviceAdapter) GetByte() uint8 {
	return s.byte
}

func (s *DeviceAdapter) SetStateMachine(v uint8) {
	s.state = v
}

func (s *DeviceAdapter) GetStateMachine() uint8 {
	return s.state
}

func (s *DeviceAdapter) SetStateMachineNext() {
	s.state++
}

func (s *DeviceAdapter) SetPrimary(v uint8) {
	s.primary = v
}

func (s *DeviceAdapter) GetPrimary() uint8 {
	return s.primary
}

func (s *DeviceAdapter) SetSecondary(v uint8) {
	s.secondary = v
}

func (s *DeviceAdapter) GetSecondary() uint8 {
	return s.secondary
}

func (s *DeviceAdapter) SetSecondaryPrev(v uint8) {
	s.secondaryPrev = v
}

func (s *DeviceAdapter) GetSecondaryPrev() uint8 {
	return s.secondaryPrev
}

func (s *DeviceAdapter) SetTimeout(v uint64) {
	s.timeout = v
}

func (s *DeviceAdapter) GetTimeout() uint64 {
	return s.timeout
}

func (s *DeviceAdapter) SetState(idx uint8, v uint8) {
	s.st[idx&0x0f] = v
}

func (s *DeviceAdapter) GetState(idx uint8) uint8 {
	return s.st[idx&0x0f]
}
