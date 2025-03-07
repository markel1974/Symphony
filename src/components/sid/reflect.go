package mos6581

import (
	"fmt"
	"github.com/markel1974/c64emu/src/components/board"
	"log"
)

type Reflect struct {
	props *board.Properties
	sid   *SID
}

func NewReflect(s *SID) *Reflect {
	r := &Reflect{
		props: nil,
		sid:   s,
	}
	r.props = board.NewProperties(r.RunCommand)
	r.props.Add(board.CreatePropertyInfo("registers", r.sid.registers, "SID register", false, r.GetRegisters, r.SetRegisters))
	return r
}

func (r *Reflect) GetProperties() *board.Properties {
	return r.props
}

func (r *Reflect) RunCommand(cmd string, args []string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (r *Reflect) SetRegisters(data []byte) {
	if len(data) != len(r.sid.registers) {
		log.Println("invalid register data length")
		return
	}
	copy(r.sid.registers, data)
}

func (r *Reflect) GetRegisters() []byte {
	v := make([]byte, len(r.sid.registers))
	copy(v, r.sid.registers)
	return v
}
