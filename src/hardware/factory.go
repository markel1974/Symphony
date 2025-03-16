package hardware

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/hardware/6510"
	c1541board "github.com/markel1974/c64emu/src/hardware/c1541/board"
	c64board "github.com/markel1974/c64emu/src/hardware/c64/board"
	"github.com/markel1974/c64emu/src/hardware/cartridges_c64"
	"github.com/markel1974/c64emu/src/hardware/cia"
	"github.com/markel1974/c64emu/src/hardware/dynamic_throttle"
	"github.com/markel1974/c64emu/src/hardware/fs_drive"
	"github.com/markel1974/c64emu/src/hardware/iec"
	"github.com/markel1974/c64emu/src/hardware/joystick_c64"
	"github.com/markel1974/c64emu/src/hardware/keyboard_c64"
	"github.com/markel1974/c64emu/src/hardware/pic_6510"
	"github.com/markel1974/c64emu/src/hardware/pla_c1541"
	"github.com/markel1974/c64emu/src/hardware/pla_c64"
	"github.com/markel1974/c64emu/src/hardware/quartz"
	"github.com/markel1974/c64emu/src/hardware/roms_c1541"
	"github.com/markel1974/c64emu/src/hardware/roms_c64"
	"github.com/markel1974/c64emu/src/hardware/sid"
	"github.com/markel1974/c64emu/src/hardware/via"
	"github.com/markel1974/c64emu/src/hardware/vic"
	vic20board "github.com/markel1974/c64emu/src/hardware/vic20/board"
	"github.com/markel1974/c64emu/src/references"
)

type Factory struct {
	cfg       *config.Config
	container map[string]references.IFactory
}

func NewFactory(cfg *config.Config) *Factory {
	f := &Factory{cfg: cfg}
	f.container = make(map[string]references.IFactory)
	var hardware []references.IFactory
	hardware = append(hardware, mos6510.NewFactory())
	hardware = append(hardware, c64board.NewFactory())
	hardware = append(hardware, c1541board.NewFactory())
	hardware = append(hardware, cartridges_c64.NewFactory())
	hardware = append(hardware, dynamic_throttle.NewFactory())
	hardware = append(hardware, mos6526.NewFactory())
	hardware = append(hardware, fs_drive.NewFactory())
	hardware = append(hardware, iec.NewFactory())
	hardware = append(hardware, keyboard_c64.NewFactory())
	hardware = append(hardware, joystick_c64.NewFactory())
	hardware = append(hardware, pic_6510.NewFactory())
	hardware = append(hardware, pla_c64.NewFactory())
	hardware = append(hardware, pla_c1541.NewFactory())
	hardware = append(hardware, quartz.NewFactory())
	hardware = append(hardware, roms_c64.NewFactory())
	hardware = append(hardware, roms_c1541.NewFactory())
	hardware = append(hardware, mos6581.NewFactory())
	hardware = append(hardware, mos6522.NewFactory())
	hardware = append(hardware, mos6569.NewFactory())
	hardware = append(hardware, vic20board.NewFactory())
	for _, h := range hardware {
		f.container[h.Identifier()] = h
	}
	return f
}

func (f *Factory) Create(parent references.IComponent, id string, suffix string) (references.IComponent, error) {
	val, ok := f.container[id]
	if !ok {
		return nil, fmt.Errorf("unknown component %s", id)
	}
	ret := val.Create(parent, f, suffix)
	//component.Register(parent, ret)
	return ret, nil
}

func (f *Factory) CreateIVIA(parent references.IComponent, id string, suffix string) (references.IComponent, references.IVIA, error) {
	component, err := f.Create(parent, id, suffix)
	if err != nil {
		return nil, nil, err
	}
	v, ok := component.(references.IVIA)
	if !ok {
		return nil, nil, fmt.Errorf("component %s is not a via", id)
	}
	return component, v, nil
}

func (f *Factory) CreateICIA(parent references.IComponent, id string, suffix string) (references.IComponent, references.ICIA, error) {
	component, err := f.Create(parent, id, suffix)
	if err != nil {
		return nil, nil, err
	}
	v, ok := component.(references.ICIA)
	if !ok {
		return nil, nil, fmt.Errorf("component %s is not a cia", id)
	}
	return component, v, nil
}

func (f *Factory) CreateI6510(parent references.IComponent, id string, suffix string) (references.IComponent, references.I6510, error) {
	component, err := f.Create(parent, id, suffix)
	if err != nil {
		return nil, nil, err
	}
	v, ok := component.(references.I6510)
	if !ok {
		return nil, nil, fmt.Errorf("component %s is not a 6510", id)
	}
	return component, v, nil
}

func (f *Factory) CreateIPIC6510(parent references.IComponent, id string, suffix string) (references.IComponent, references.IPIC6510, error) {
	component, err := f.Create(parent, id, suffix)
	if err != nil {
		return nil, nil, err
	}
	v, ok := component.(references.IPIC6510)
	if !ok {
		return nil, nil, fmt.Errorf("component %s is not a pic6510", id)
	}
	return component, v, nil
}

func (f *Factory) CreateIPLAc1541(parent references.IComponent, id string, suffix string) (references.IComponent, references.IPLAc1541, error) {
	component, err := f.Create(parent, id, suffix)
	if err != nil {
		return nil, nil, err
	}
	v, ok := component.(references.IPLAc1541)
	if !ok {
		return nil, nil, fmt.Errorf("component %s is not a cia", id)
	}
	return component, v, nil
}

func (f *Factory) CreateIROMLoaderC1541(parent references.IComponent, id string, suffix string) (references.IComponent, references.IROMLoaderC1541, error) {
	component, err := f.Create(parent, id, suffix)
	if err != nil {
		return nil, nil, err
	}
	v, ok := component.(references.IROMLoaderC1541)
	if !ok {
		return nil, nil, fmt.Errorf("component %s is not a cia", id)
	}
	return component, v, nil
}
