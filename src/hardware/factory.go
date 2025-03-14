package hardware

import (
	"fmt"
	"strings"

	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	mos6510 "github.com/markel1974/c64emu/src/hardware/6510"
	c1541board "github.com/markel1974/c64emu/src/hardware/c1541/board"
	c64board "github.com/markel1974/c64emu/src/hardware/c64/board"
	mos6526 "github.com/markel1974/c64emu/src/hardware/cia"
	"github.com/markel1974/c64emu/src/hardware/quartz"
	mos6581 "github.com/markel1974/c64emu/src/hardware/sid"
	mos6522 "github.com/markel1974/c64emu/src/hardware/via"
	mos6569 "github.com/markel1974/c64emu/src/hardware/vic"
	vic20board "github.com/markel1974/c64emu/src/hardware/vic20/board"
)

type Factory struct {
	cfg *config.Config
}

func NewFactory(cfg *config.Config) *Factory {
	return &Factory{cfg: cfg}
}

func (f *Factory) Create(parent component.IComponent, id string, suffix string) (component.IComponent, error) {
	rid := strings.ToLower(strings.TrimSpace(id))
	switch rid {
	case "c64":
		return c64board.NewBoard(parent, f, suffix), nil
	case "vic20":
		return vic20board.NewBoard(parent, f, suffix), nil
	case "c1541":
		return c1541board.New(parent, f, suffix), nil
	case "mos65102":
		return mos6510.NewCPU(parent, f, suffix), nil
	case "mos6510":
		return mos6510.NewCPU(parent, f, suffix), nil
	case "mos6526":
		return mos6526.NewCIA(parent, f, suffix), nil
	case "mos6581":
		return mos6581.NewSID(parent, f, suffix), nil
	case "mos6569":
		return mos6569.NewVIC(parent, f, suffix), nil
	case "mos6522":
		return mos6522.NewVia(parent, f, suffix), nil
	case "quartz":
		return quartz.NewQuartz(parent, f, suffix), nil
	default:
		return nil, fmt.Errorf("unknown component %s", id)
	}
}
