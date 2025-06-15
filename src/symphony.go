package symphony

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/hardware"
	"github.com/markel1974/c64emu/src/references"
	"log"
)

var _c64DefaultHardware = []struct {
	label    string
	id       string
	instance int
}{
	{"c64", "c64_board", 0},
	{"c64", "iec", 0},
	{"c64", "mos6569", 0},
	{"c64", "mos6526", 0},
	{"c64", "mos6526", 1},
	{"c64", "cartridges_c64", 0},
	{"c64", "mos6510", 0},
	{"c64", "pic_6510", 0},
	{"c64", "dynamic_throttle", 0},
	{"c64", "mos6581", 0},
	{"c64", "ram_c64", 0},
	{"c64", "c64_color_ram", 0},
	{"c64", "roms_c64", 0},
	{"c64", "pla_c64", 0},
	{"c64", "keyboard_c64", 0},
	{"c64", "joystick_c64", 0},
	{"c64", "joystick_c64", 1},
	{"c64", "quartz", 0},

	/*
		{"c1541_8", "roms_c1541", 0},
		{"c1541_8", "pla_c1541", 0},
		{"c1541_8", "pic_6510", 0},
		{"c1541_8", "mos6510", 0},
		{"c1541_8", "mos6522", 0},
		{"c1541_8", "mos6522", 1},

		{"c1541_9", "roms_c1541", 0},
		{"c1541_9", "pla_c1541", 0},
		{"c1541_9", "pic_6510", 0},
		{"c1541_9", "mos6510", 0},
		{"c1541_9", "mos6522", 0},
		{"c1541_9", "mos6522", 1},

		{"c1541_10", "roms_c1541", 0},
		{"c1541_10", "pla_c1541", 0},
		{"c1541_10", "pic_6510", 0},
		{"c1541_10", "mos6510", 0},
		{"c1541_10", "mos6522", 0},
		{"c1541_10", "mos6522", 1},

		{"c1541_11", "roms_c1541", 0},
		{"c1541_11", "pla_c1541", 0},
		{"c1541_11", "pic_6510", 0},
		{"c1541_11", "mos6510", 0},
		{"c1541_11", "mos6522", 0},
		{"c1541_11", "mos6522", 1},

	*/
}

type Symphony struct {
	boardComponent references.IComponent
	displayRender  references.IDisplayRender
	cfg            *config.Config
}

func New() *Symphony {
	return &Symphony{
		boardComponent: nil,
		displayRender:  nil,
		cfg:            nil,
	}
}

func (s *Symphony) Setup(opt *Options) error {
	gFactory := opt.GetGraphicsFactory()
	if gFactory == nil {
		return fmt.Errorf("graphics factory is nil")
	}
	aFactory := opt.GetAudioFactory()
	if aFactory == nil {
		return fmt.Errorf("audio factory is nil")
	}
	s.cfg = config.New()
	if len(opt.Prg) > 0 {
		s.cfg.SetPrg(opt.Prg)
	}
	for _, c := range opt.Cartridges {
		if err := s.cfg.AddCartridge(c); err != nil {
			return err
		}
	}
	for _, c := range opt.Drives {
		if err := s.cfg.AddDrive(c); err != nil {
			return err
		}
	}
	for _, c := range opt.Disks {
		if err := s.cfg.AddSpareDisk(c); err != nil {
			return err
		}
	}
	if opt.NoJiffy {
		s.cfg.DisableJiffy()
	}

	s.displayRender = gFactory.Create(opt.RenderId)
	display, err := s.displayRender.CreateDisplayBuffer(0x180, 0x110) //mos6569.DisplayX, mos6569.DisplayY
	if err != nil {
		return err
	}
	audioRender := aFactory.Create(opt.PlayerId)
	if err = audioRender.Setup(s.cfg); err != nil {
		return err
	}
	hwFactory := hardware.NewFactory(display, audioRender, s.cfg)
	s.boardComponent = nil
	var hw []references.IComponent
	components := make(map[string]references.IComponent)
	for _, h := range _c64DefaultHardware {
		var comp references.IComponent
		if comp, err = hwFactory.Create(s.boardComponent, h.label, h.id, h.instance); err != nil {
			log.Fatal(err.Error())
		}
		components[comp.HardwareId()] = comp
		if s.boardComponent == nil {
			s.boardComponent = comp
		}
		hw = append(hw, comp)
	}
	if s.boardComponent == nil {
		return fmt.Errorf("board component is nil")
	}

	//TODO REMOVE WHEN TREE IS READY
	//BEGIN
	for _, c := range hw {
		if err = c.Setup(); err != nil {
			return err
		}
	}
	for _, c := range hw {
		if err = c.Connect(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Symphony) Start() error {
	if s.boardComponent == nil {
		return fmt.Errorf("board component is nil")
	}
	board, ok := s.boardComponent.(references.IC64Board)
	if !ok {
		return fmt.Errorf("board component is not a board")
	}
	if err := s.displayRender.Setup(board, s.cfg); err != nil {
		return err
	}
	if err := board.Start(); err != nil {
		return err
	}
	if err := s.displayRender.Start(); err != nil {
		return err
	}
	return nil
}

func (s *Symphony) GetBoard() references.IComponent {
	return s.boardComponent
}
