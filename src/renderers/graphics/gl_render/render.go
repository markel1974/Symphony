package gl_render

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/renderers/graphics/gl_render/pixels"
)

type Render struct {
	board      references.IBoard
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

func (g *Render) Setup(board references.IBoard, cfg *config.Config) error {
	g.board = board
	if err := g.board.Mount(g); err != nil {
		return err
	}
	if err := g.inputs.Setup(g.board, cfg); err != nil {
		return err
	}
	return nil
}

func (g *Render) Start() error {
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

func (g *Render) VBlank() {
	if g.win.MouseInsideWindow() {
		g.inputs.MouseMove(g.win.MousePositionXY())
	}
	g.inputs.Keys(g.win.KeysPressed())
	g.surface.Draw(g.win, g.matrix)
	g.win.Update()
	g.run = !g.win.Closed()
}

func (g *Render) LedActivity(deviceNumber uint8, led bool) {
	fmt.Println("LED STATE", deviceNumber, led)
}
