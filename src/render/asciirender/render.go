package asciirender

import (
	"github.com/markel1974/c64emu/src/c64/board"
	"github.com/markel1974/c64emu/src/components/vic"
	"github.com/markel1974/c64emu/src/config"
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
	display      *DisplayBuffer
	audio        *Audio
}

func New(cfg *config.Config) *Render {
	g := &Render{
		cfg:          cfg,
		fullscreen:   false,
		screenWidth:  mos6569.DisplayX,
		screenHeight: mos6569.DisplayY,
		scale:        2.0,
		display:      NewDisplayBuffer(),
		audio:        NewAudio(),
	}
	g.maxW = float64(g.screenWidth) * g.scale
	g.maxH = float64(g.screenHeight) * g.scale
	return g
}

func (g *Render) setup() {
	g.c64Board = board.NewBoard(g.display, g.audio)
	_ = g.c64Board.Setup(g.cfg)
	//g.inputs.Setup(g.c64Board, g.maxW, g.maxH)
}

func (g *Render) Start() {
	g.setup()
	dt := common.NewDynamicThrottling(mos6569.FrameInterval)

	run := true
	for run {
		dt.DynamicThrottling()
		//g.inputs.Keys(win.KeysPressed())
		for {
			if vBlank := g.c64Board.Emulate(); vBlank {
				break
			}
		}
		//g.surface.Draw(win, g.matrix)
		//g.audio.Play()
		//win.Update()
		//if (dt.Counter() & 0xf) == 0xf {
		//	run = !win.Closed()
		//}
	}
}
