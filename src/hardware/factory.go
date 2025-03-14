package hardware

import (
	"fmt"
	"github.com/markel1974/c64emu/src/hardware/iec"
	"github.com/markel1974/c64emu/src/hardware/joystick"
	"github.com/markel1974/c64emu/src/hardware/keyboard"
	"github.com/markel1974/c64emu/src/hardware/throttle"
	"github.com/markel1974/c64emu/src/references"
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

type constructorFn func(component.IComponent, references.IComponentFactory, string) component.IComponent
type Factory struct {
	cfg       *config.Config
	container map[string]constructorFn
}

func NewFactory(cfg *config.Config) *Factory {
	f := &Factory{cfg: cfg}
	f.container = make(map[string]constructorFn)
	f.container["c64"] = c64board.NewBoardComponent
	f.container["vic20"] = vic20board.NewBoardComponent
	f.container["c1541"] = c1541board.NewComponent
	f.container["iec"] = iec.NewDispatcherComponent
	f.container["mos6510"] = mos6510.NewCPUComponent
	f.container["mos6502"] = mos6510.NewCPUComponent
	f.container["mos6526"] = mos6526.NewCIAComponent
	f.container["mos6581"] = mos6581.NewSIDComponent
	f.container["mos6569"] = mos6569.NewVICComponent
	f.container["mos6522"] = mos6522.NewViaComponent
	f.container["quartz"] = quartz.NewQuartzComponent
	f.container["joystick"] = joystick.NewJoystickComponent
	f.container["keyboard"] = keyboard.NewKeyboardComponent
	f.container["throttle"] = throttle.NewDynamicThrottleComponent
	return f
}

func (f *Factory) Create(parent component.IComponent, id string, suffix string) (component.IComponent, error) {
	val, ok := f.container[id]
	if !ok {
		return nil, fmt.Errorf("unknown component %s", id)
	}
	ret := val(parent, f, suffix)
	//component.Register(parent, ret)
	return ret, nil
}

func (f *Factory) CreateComponent(parent component.IComponent, id string, suffix string) (component.IComponent, error) {
	rid := strings.ToLower(strings.TrimSpace(id))
	switch rid {
	case "c64":
		return c64board.NewBoard(parent, f, suffix), nil
	case "vic20":
		return vic20board.NewBoard(parent, f, suffix), nil
	case "c1541":
		return c1541board.NewBoard(parent, f, suffix), nil
	case "iec":
		return iec.NewDispatcher(parent, f, suffix), nil
	case "mos6502":
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
	case "joystick":
		return joystick.NewJoystick(parent, f, suffix), nil
	case "keyboard":
		return keyboard.NewKeyboard(parent, f, suffix), nil
	case "throttle":
		return throttle.NewDynamicThrottle(parent, f, suffix), nil
	default:
		return nil, fmt.Errorf("unknown component %s", id)
	}
}
