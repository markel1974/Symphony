package hardware

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/registry"

	_ "github.com/markel1974/c64emu/src/hardware/6510"
	_ "github.com/markel1974/c64emu/src/hardware/c1541/board"
	_ "github.com/markel1974/c64emu/src/hardware/c64/board"
	_ "github.com/markel1974/c64emu/src/hardware/cartridges_c64"
	_ "github.com/markel1974/c64emu/src/hardware/cartridges_c64/easyflash"
	_ "github.com/markel1974/c64emu/src/hardware/cartridges_c64/external_cpu"
	_ "github.com/markel1974/c64emu/src/hardware/cartridges_c64/generic"
	_ "github.com/markel1974/c64emu/src/hardware/cartridges_c64/magicdesk"
	_ "github.com/markel1974/c64emu/src/hardware/cartridges_c64/ocean"
	_ "github.com/markel1974/c64emu/src/hardware/cartridges_c64/reu"
	_ "github.com/markel1974/c64emu/src/hardware/cia"
	_ "github.com/markel1974/c64emu/src/hardware/dynamic_throttle"
	_ "github.com/markel1974/c64emu/src/hardware/fs_drive"
	_ "github.com/markel1974/c64emu/src/hardware/iec"
	_ "github.com/markel1974/c64emu/src/hardware/joystick_c64"
	_ "github.com/markel1974/c64emu/src/hardware/keyboard_c64"
	_ "github.com/markel1974/c64emu/src/hardware/pic_6510"
	_ "github.com/markel1974/c64emu/src/hardware/pla_c1541"
	_ "github.com/markel1974/c64emu/src/hardware/pla_c64"
	_ "github.com/markel1974/c64emu/src/hardware/quartz"
	_ "github.com/markel1974/c64emu/src/hardware/roms_c1541"
	_ "github.com/markel1974/c64emu/src/hardware/roms_c64"
	_ "github.com/markel1974/c64emu/src/hardware/sid"
	_ "github.com/markel1974/c64emu/src/hardware/via"
	_ "github.com/markel1974/c64emu/src/hardware/vic"
)

type Factory struct {
	cfg       *config.Config
	db        references.IDisplayBuffer
	player    references.IAudioRender
	container map[string]references.IFactory
}

func NewFactory(db references.IDisplayBuffer, player references.IAudioRender, cfg *config.Config) *Factory {
	f := &Factory{
		db:        db,
		player:    player,
		cfg:       cfg,
		container: registry.ComponentFactories(),
	}
	return f
}

func (f *Factory) Create(parent references.IComponent, id string, label int) (references.IComponent, error) {
	val, ok := f.container[id]
	if !ok {
		return nil, fmt.Errorf("unknown component %s", id)
	}
	ret := val.Create(parent, f, label)
	//component.Register(parent, ret)
	return ret, nil
}

func (f *Factory) GetIDisplayBuffer() references.IDisplayBuffer {
	return f.db
}

func (f *Factory) GetIAudioRender() references.IAudioRender {
	return f.player
}
