package glrender

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/render/glrender/pixels"
)

type Render struct {
	cfg          *config.Config
	c64Board     references.IBoard
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

func New(board references.IBoard, cfg *config.Config) *Render {
	w, h := board.GetScreenSize()
	g := &Render{
		c64Board:     board,
		cfg:          cfg,
		fullscreen:   false,
		screenWidth:  w,
		screenHeight: h,
		scale:        3.0,
		inputs:       NewInputs(),
	}
	g.maxW = float64(g.screenWidth) * g.scale
	g.maxH = float64(g.screenHeight) * g.scale
	return g
}

func (g *Render) setup(pos pixels.Vec) {
	g.picture = pixels.NewPicture(pixels.NewRect(float64(0), float64(0), float64(g.screenWidth), float64(g.screenHeight)))
	g.surface = pixels.NewSprite()
	g.surface.SetCachedMode(pixels.CacheModeUpdate)
	g.surface.Set(g.picture, g.picture.Bounds())
	g.matrix = pixels.IM.Moved(pos).Scaled(pos, g.scale)
	g.display = NewDisplayBuffer(g.picture)
	g.audio = NewAudio()
	if err := g.c64Board.Setup(g.display, g.audio, g.cfg); err != nil {
		panic(err)
	}
	g.inputs.Setup(g.c64Board, g.maxW, g.maxH)
}

func (g *Render) Start() error {
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

	c := win.Bounds().Center()
	g.setup(c)
	dt := g.c64Board.Throttle()

	run := true
	for run {
		dt.Throttle()
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
