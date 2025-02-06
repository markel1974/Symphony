package glrender

import (
	"github.com/markel1974/c64emu/src/c64/board"
	"github.com/markel1974/c64emu/src/components/vic"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/pixels"
	"github.com/markel1974/c64emu/src/render/common"
)

type Render struct {
	cfg          *config.Config
	c64Board     *board.Board
	scale        float64
	fullscreen   bool
	showMap      bool
	screenWidth  int
	screenHeight int
	maxW         float64
	maxH         float64
	picture      *pixels.Picture
	matrix       pixels.Matrix
	surface      *pixels.Sprite
	display      *DisplayBuffer
	inputs       *Inputs
	audio        *Audio
}

func New(cfg *config.Config) *Render {
	g := &Render{
		cfg:          cfg,
		fullscreen:   false,
		screenWidth:  mos6569.DisplayX,
		screenHeight: mos6569.DisplayY,
		scale:        3.0,
		inputs:       NewInputs(),
	}
	g.maxW = float64(g.screenWidth) * g.scale
	g.maxH = float64(g.screenHeight) * g.scale
	return g
}

func (g *Render) setup(pos pixels.Vec) {
	g.picture = pixels.NewPicture(pixels.R(float64(0), float64(0), float64(g.screenWidth), float64(g.screenHeight)))
	g.surface = pixels.NewSprite()
	g.surface.SetCachedMode(pixels.CacheModeUpdate)
	g.surface.Set(g.picture, g.picture.Bounds())
	g.matrix = pixels.IM.Moved(pos).Scaled(pos, g.scale)
	g.display = NewDisplayBuffer(g.picture)
	g.audio = NewAudio()
	g.c64Board = board.NewBoard(g.display, g.audio)
	_ = g.c64Board.Setup(g.cfg)
	g.inputs.Setup(g.c64Board, g.maxW, g.maxH)
}

func (g *Render) Start() error {
	return pixels.GLRun(g.run)
}

func (g *Render) run() {
	cfg := pixels.WindowConfig{
		Bounds:      pixels.R(0, 0, g.maxW, g.maxH),
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

	c := win.Bounds().Center()
	g.setup(c)
	dt := common.NewDynamicThrottling(mos6569.FrameInterval)

	run := true
	for run {
		dt.DynamicThrottling()
		if win.MouseInsideWindow() {
			g.inputs.MouseMove(win.MousePositionXY())
		}
		g.inputs.Keys(win.KeysPressed())
		for {
			if vBlank := g.c64Board.Emulate(); vBlank {
				break
			}
		}
		g.surface.Draw(win, g.matrix)
		g.audio.Play()
		win.Update()
		if (dt.Counter() & 0xf) == 0xf {
			run = !win.Closed()
		}
	}
}
