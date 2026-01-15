package hardware

import (
	"fmt"
	"github.com/markel1974/symphony/src/config"
	"github.com/markel1974/symphony/src/references"
	"github.com/markel1974/symphony/src/registry"

	_ "github.com/markel1974/symphony/src/hardware/c1541_board_rev1/board"
	_ "github.com/markel1974/symphony/src/hardware/c1541_pla_rev1"
	_ "github.com/markel1974/symphony/src/hardware/c1541_ram_rev1"
	_ "github.com/markel1974/symphony/src/hardware/c1541_roms_rev1"
	_ "github.com/markel1974/symphony/src/hardware/c64_board_rev1/board"
	_ "github.com/markel1974/symphony/src/hardware/c64_cartridges_rev1"
	_ "github.com/markel1974/symphony/src/hardware/c64_cartridges_rev1/easyflash"
	_ "github.com/markel1974/symphony/src/hardware/c64_cartridges_rev1/external_cpu"
	_ "github.com/markel1974/symphony/src/hardware/c64_cartridges_rev1/generic"
	_ "github.com/markel1974/symphony/src/hardware/c64_cartridges_rev1/magicdesk"
	_ "github.com/markel1974/symphony/src/hardware/c64_cartridges_rev1/ocean"
	_ "github.com/markel1974/symphony/src/hardware/c64_cartridges_rev1/reu"
	_ "github.com/markel1974/symphony/src/hardware/c64_color_ram_rev1"
	_ "github.com/markel1974/symphony/src/hardware/c64_joystick_rev1"
	_ "github.com/markel1974/symphony/src/hardware/c64_keyboard_rev1"
	_ "github.com/markel1974/symphony/src/hardware/c64_pla_rev1"
	_ "github.com/markel1974/symphony/src/hardware/c64_ram_rev1"
	_ "github.com/markel1974/symphony/src/hardware/c64_roms_rev1"
	_ "github.com/markel1974/symphony/src/hardware/dynamic_throttle_rev1"
	_ "github.com/markel1974/symphony/src/hardware/iec_rev1"
	_ "github.com/markel1974/symphony/src/hardware/media_drive_rev1"
	_ "github.com/markel1974/symphony/src/hardware/mos6510_rev1"
	_ "github.com/markel1974/symphony/src/hardware/mos6522_via_rev1"
	_ "github.com/markel1974/symphony/src/hardware/mos6526_cia_rev1"
	_ "github.com/markel1974/symphony/src/hardware/mos6569_vic_rev1"
	_ "github.com/markel1974/symphony/src/hardware/mos6581_sid_rev1"
	_ "github.com/markel1974/symphony/src/hardware/quartz_rev1"
)

type Factory struct {
	cfg       *config.Config
	db        references.IDisplayBuffer
	player    references.IAudioRender
	factories []references.IFactory
	byId      map[string]references.IFactory
}

func NewFactory(db references.IDisplayBuffer, player references.IAudioRender, cfg *config.Config) *Factory {
	f := &Factory{
		db:        db,
		player:    player,
		cfg:       cfg,
		factories: registry.ComponentFactories(),
		byId:      make(map[string]references.IFactory),
	}
	for _, v := range f.factories {
		f.byId[v.Identifier()] = v
	}
	return f
}

func (f *Factory) Create(parent references.IComponent, label string, id string, instance int) (references.IComponent, error) {
	val, ok := f.byId[id]
	if !ok {
		return nil, fmt.Errorf("unknown component %s", id)
	}
	ret := val.Create(parent, f, label, instance)
	//component.Register(parent, ret)
	return ret, nil
}

func (f *Factory) GetIDisplayBuffer() references.IDisplayBuffer {
	return f.db
}

func (f *Factory) GetIAudioRender() references.IAudioRender {
	return f.player
}

func (f *Factory) GetConfig() *config.Config {
	return f.cfg
}
