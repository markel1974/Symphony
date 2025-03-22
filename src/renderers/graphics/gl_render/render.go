package gl_render

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
	pixels2 "github.com/markel1974/c64emu/src/renderers/graphics/gl_render/pixels"
)

type Render struct {
	board      references.IBoard
	dt         references.IThrottle
	scale      float64
	fullscreen bool
	maxW       float64
	maxH       float64
	picture    *pixels2.Picture
	inputs     *Inputs
	win        *pixels2.GLWindow
	matrix     pixels2.Matrix
	surface    *pixels2.Sprite
	run        bool
}

func New() *Render {
	g := &Render{
		board:      nil,
		dt:         nil,
		win:        nil,
		fullscreen: false,
		scale:      3,
		inputs:     NewInputs(),
		run:        true,
	}
	return g
}

func (g *Render) CreateDisplayBuffer(w int, h int) (references.IDisplayBuffer, error) {
	g.maxW = float64(w) * g.scale
	g.maxH = float64(h) * g.scale
	g.picture = pixels2.NewPicture(pixels2.NewRect(float64(0), float64(0), float64(w), float64(h)))
	display := NewDisplayBuffer(g.picture)
	return display, nil
}

func (g *Render) Setup(board references.IBoard, cfg *config.Config) error {
	g.board = board
	g.dt = board.Throttle()
	g.board.VBlankSignal().Bind(g.vBlankSlot)
	g.board.LEDSignal().Bind(g.ledSlot)
	if err := g.inputs.Setup(g.board, cfg); err != nil {
		return err
	}
	return nil
}

func (g *Render) Start() error {
	return pixels2.GLRun(g.runner)
}

func (g *Render) runner() {
	cfg := pixels2.WindowConfig{
		Bounds:      pixels2.NewRect(0, 0, g.maxW, g.maxH),
		VSync:       true,
		Undecorated: false,
		Smooth:      false,
	}
	if g.fullscreen {
		cfg.Monitor = pixels2.PrimaryMonitor()
	}
	var err error
	g.win, err = pixels2.NewGLWindow(cfg)
	if err != nil {
		panic(err)
	}
	pos := g.win.Bounds().Center()
	g.surface = pixels2.NewSprite()
	g.surface.SetCachedMode(pixels2.CacheModeUpdate)
	g.surface.Set(g.picture, g.picture.Bounds())
	g.matrix = pixels2.IM.Moved(pos).Scaled(pos, g.scale)
	for g.run {
		g.board.Emulate()
	}
}

func (g *Render) vBlankSlot() {
	g.dt.Throttle()
	if g.win.MouseInsideWindow() {
		g.inputs.MouseMove(g.win.MousePositionXY())
	}
	g.inputs.Keys(g.win.KeysPressed())
	g.surface.Draw(g.win, g.matrix)
	//g.player.Write(nil, 0, 0)
	g.win.Update()
	if (g.dt.Counter() & 0xf) == 0xf {
		g.run = !g.win.Closed()
	}
}

func (g *Render) ledSlot(state uint32) {
	device := uint8(state & 0xf)
	led := uint8((state >> 8) & 0xf)
	fmt.Println("LED STATE", device, led)
}
