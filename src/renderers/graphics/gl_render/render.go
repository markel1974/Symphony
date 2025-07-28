package gl_render

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/renderers/graphics/gl_render/pixels"
)

// Render represents the main rendering structure responsible for managing display, input, and graphical interactions.
type Render struct {
	board       references.IC64Board
	scale       float64
	fullscreen  bool
	maxW        float64
	maxH        float64
	picture     *pixels.Picture
	inputs      *Inputs
	win         *pixels.GLWindow
	surfaceM    pixels.Matrix
	surface     *pixels.Sprite
	ledSurface  *pixels.Sprite
	ledSurfaceM pixels.Matrix
	run         bool
	led         bool
}

// New initializes and returns a new instance of the Render struct with default values.
func New() *Render {
	g := &Render{
		board:      nil,
		win:        nil,
		fullscreen: false,
		scale:      3,
		inputs:     NewInputs(),
		run:        true,
		led:        false,
	}
	return g
}

// CreateDisplayBuffer initializes a display buffer with specified dimensions and returns it along with any potential error.
func (g *Render) CreateDisplayBuffer(w int, h int) (references.IDisplayBuffer, error) {
	g.maxW = float64(w) * g.scale
	g.maxH = float64(h) * g.scale
	g.picture = pixels.NewPicture(pixels.NewRect(float64(0), float64(0), float64(w), float64(h)))
	display := NewDisplayBuffer(g.picture)
	return display, nil
}

// Setup initializes the Render by wiring it to the given IC64Board and setting up inputs using the provided configuration.
func (g *Render) Setup(board references.IC64Board, cfg *config.Config) error {
	g.board = board
	if err := g.board.Wire(g); err != nil {
		return err
	}
	if err := g.inputs.Setup(g.board, cfg); err != nil {
		return err
	}
	return nil
}

// Start initializes and runs the main graphics loop using pixels.GLRun, ensuring proper rendering setup and execution.
func (g *Render) Start() error {
	return pixels.GLRun(g.runner)
}

// runner initializes the rendering loop, creates a window, sets up sprites and matrix transformations, and executes emulation.
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

	g.surfaceM = pixels.IM.Moved(pos).ScaledXY(pos, pixels.Vec{X: g.scale, Y: -g.scale})
	//g.surfaceM = pixels.IM.Moved(pos).ScaledXY(pos, pixels.Vec{X: g.scale, Y: g.scale})

	g.ledSurface = pixels.NewSprite()
	g.ledSurface.SetCachedMode(pixels.CacheModeUpdate)
	g.ledSurface.Set(g.picture, g.picture.Bounds())
	g.ledSurfaceM = pixels.IM.Moved(pos).Scaled(pos, 0.1)
	emulate := g.board.Emulate
	for g.run {
		emulate()
	}
}

// VBlank handles the vertical blanking interval by processing inputs, updating the display surface, and managing window events.
func (g *Render) VBlank() {
	if g.win.MouseInsideWindow() {
		g.inputs.MouseMove(g.win.MousePositionXY())
	}
	g.inputs.Keys(g.win.KeysPressed())
	g.surface.Draw(g.win, g.surfaceM)
	if g.led {
		g.ledSurface.Draw(g.win, g.ledSurfaceM)
	}
	g.win.Update()
	g.run = !g.win.Closed()
}

// LedActivity updates the LED state for a specified device by setting the `led` field to the provided boolean value.
func (g *Render) LedActivity(deviceNumber uint8, led bool) {
	g.led = led
	//fmt.Println("LED STATE", deviceNumber, led)
}
