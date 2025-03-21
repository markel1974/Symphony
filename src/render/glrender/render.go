package glrender

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/render/glrender/pixels"
)

type Render struct {
	board      references.IBoard
	scale      float64
	fullscreen bool
	maxW       float64
	maxH       float64
	picture    *pixels.Picture
	inputs     *Inputs
	//player     references.IAudioRender
}

func New() *Render {
	g := &Render{
		board:      nil,
		fullscreen: false,
		scale:      3,
		inputs:     NewInputs(),
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
	g.inputs.Setup(g.board)
	return pixels.GLRun(g.run)
}

func (g *Render) run() {
	cfg := pixels.WindowConfig{
		Bounds:      pixels.NewRect(0, 0, g.maxW, g.maxH),
		VSync:       true,
		Undecorated: false,
		Smooth:      false,
	}
	if g.fullscreen {
		cfg.Monitor = pixels.PrimaryMonitor()
	}
	win, err := pixels.NewGLWindow(cfg)
	if err != nil {
		panic(err)
	}
	pos := win.Bounds().Center()
	surface := pixels.NewSprite()
	surface.SetCachedMode(pixels.CacheModeUpdate)
	surface.Set(g.picture, g.picture.Bounds())
	matrix := pixels.IM.Moved(pos).Scaled(pos, g.scale)
	dt := g.board.Throttle()
	run := true
	for run {
		dt.Throttle()
		if win.MouseInsideWindow() {
			g.inputs.MouseMove(win.MousePositionXY())
		}
		g.inputs.Keys(win.KeysPressed())
		for {
			if vBlank := g.board.Emulate(); vBlank {
				break
			}
		}
		surface.Draw(win, matrix)
		//g.player.Write(nil, 0, 0)
		win.Update()
		if (dt.Counter() & 0xf) == 0xf {
			run = !win.Closed()
		}
	}
}
