package glrender

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/render/glrender/pixels"
)

type Render struct {
	board      references.IBoard
	dt         references.IThrottle
	scale      float64
	fullscreen bool
	maxW       float64
	maxH       float64
	picture    *pixels.Picture
	inputs     *Inputs
	win        *pixels.GLWindow
	matrix     pixels.Matrix
	surface    *pixels.Sprite
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
	g.picture = pixels.NewPicture(pixels.NewRect(float64(0), float64(0), float64(w), float64(h)))
	display := NewDisplayBuffer(g.picture)
	return display, nil
}

func (g *Render) Start(board references.IBoard) error {
	g.board = board
	g.dt = board.Throttle()
	g.board.SetVBlankEmitter(g.vBlank)
	g.inputs.Setup(g.board)
	return pixels.GLRun(g.runner)
}

func (g *Render) runner() {
	cfg := pixels.WindowConfig{
		Bounds:      pixels.NewRect(0, 0, g.maxW, g.maxH),
		VSync:       true,
		Undecorated: false,
		Smooth:      false,
	}
	if g.fullscreen {
		cfg.Monitor = pixels.PrimaryMonitor()
	}
	var err error
	g.win, err = pixels.NewGLWindow(cfg)
	if err != nil {
		panic(err)
	}
	pos := g.win.Bounds().Center()
	g.surface = pixels.NewSprite()
	g.surface.SetCachedMode(pixels.CacheModeUpdate)
	g.surface.Set(g.picture, g.picture.Bounds())
	g.matrix = pixels.IM.Moved(pos).Scaled(pos, g.scale)
	for g.run {
		g.board.Emulate()
	}
}

func (g *Render) vBlank() {
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
