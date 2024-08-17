package render

import (
	"github.com/markel1974/c64emu/src/board"
	"github.com/markel1974/c64emu/src/components/vic"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/pixels"
)

type Point struct {
	X int
	Y int
}

type Render struct {
	prefs        *config.Config
	c64Board     *board.Board
	scale        float64
	fullscreen   bool
	showMap      bool
	screenWidth  int
	screenHeight int
	maxW         float64
	maxH         float64
	mainSurface  *pixels.PictureRGBA
	mainMatrix   pixels.Matrix
	mainSprite   *pixels.Sprite
	db           *DisplayBuffer
	inputs       *Inputs
}

func New(prefs *config.Config) *Render {
	g := &Render{
		prefs:        prefs,
		fullscreen:   false,
		screenWidth:  vic.DisplayX,
		screenHeight: vic.DisplayY,
		scale:        3.0,
		inputs:       NewInputs(),
	}
	g.maxW = float64(g.screenWidth) * g.scale
	g.maxH = float64(g.screenHeight) * g.scale
	return g
}

func (g *Render) setup(pos pixels.Vec) {
	g.mainSurface = pixels.NewPictureRGBA(pixels.R(float64(0), float64(0), float64(g.screenWidth), float64(g.screenHeight)))
	g.mainSprite = pixels.NewSprite()
	g.mainSprite.SetCached(pixels.CacheModeUpdate)
	g.mainSprite.Set(g.mainSurface, g.mainSurface.Bounds())
	g.mainMatrix = pixels.IM.Moved(pos).Scaled(pos, g.scale)
	g.db = NewDisplayBuffer(g.mainSurface)
	g.c64Board = board.NewBoard(g.db)
	_ = g.c64Board.Setup(g.prefs)
	g.inputs.Setup(g.c64Board, g.maxW, g.maxH)
}

func (g *Render) Start() {
	pixels.GLRun(g.run)
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
	dt := NewDynamicThrottling(vic.FrameInterval)

	run := true
	for run {
		dt.DynamicThrottling()
		if win.MouseInsideWindow() {
			mousePos := win.MousePosition()
			g.inputs.MouseMove(mousePos.X, mousePos.Y)
		}
		g.inputs.Keys(win.KeysPressed())
		for {
			if vBlank := g.c64Board.Emulate(); vBlank {
				break
			}
		}
		g.mainSprite.Draw(win, g.mainMatrix)
		//g.drawDirect(win, g.c64Board.GetDisplayBuffer())
		win.Update()
		if dt.Counter()&0xf == 0xf {
			run = !win.Closed()
		}
	}
}

//func (g *Render) drawDirect(win *pixels.GLWindow, chunky []uint8) {
//	for cIdx, cVal := range g.cacheCoords {
//		v := chunky[cIdx]
//		g.mainSurface.SetRGBADirectArray(cVal, _colors[v])
//	}
//	g.mainSprite.Draw(win, g.mainMatrix)
//	return
//}

//func (g *Render) SetPixel(idx int, val uint8) {
//	if idx < len(g.cacheCoords) {
//		g.mainSurface.SetRGBADirectArray(g.cacheCoords[idx], _colors[val])
//	}
//}

//func (g *Render) draw(win *pixels.GLWindow, chunky []uint8) {
//	idx := 0
//	var x int
//	var p uint8
//	for y := 0; y < vic.DisplayY; y++ {
//		for x = 0; x < vic.DisplayX; x++ {
//			p = chunky[idx]
//			idx++
//			g.mainSurface.SetRGBAArray(x, y, _colors[p])
//		}
//	}
//	g.mainSprite.Draw(win, g.mainMatrix)
//}
