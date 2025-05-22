//go:build js && wasm

package wasm_render

import (
	"syscall/js"

	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// Render manages the display rendering process and user interaction through inputs and a linked board implementation.
// It integrates with a DisplayBuffer for screen updates and relies on IBoard for controlling board operations.
type Render struct {
	board         references.IBoard
	displayBuffer *DisplayBuffer
	vBlank        bool
	w             int
	h             int
	input         *Inputs
}

// NewRender initializes and returns a pointer to a new instance of the Render struct with default values.
func NewRender() *Render {
	return &Render{
		displayBuffer: nil, //NewDisplayBuffer(320, 200),
		board:         nil,
		vBlank:        false,
		input:         NewInputs(),
	}
}

// CreateDisplayBuffer initializes a new display buffer with the specified width and height, stores it, and returns it.
func (g *Render) CreateDisplayBuffer(w int, h int) (references.IDisplayBuffer, error) {
	g.w = w
	g.h = h
	g.displayBuffer = NewDisplayBuffer(w, h)
	return g.displayBuffer, nil
}

// Setup initializes the Render object with the given IBoard and Config instances, setting up the necessary components.
// Returns an error if the board or input setup fails.
func (g *Render) Setup(board references.IBoard, cfg *config.Config) error {
	g.board = board
	if err := g.board.Mount(g); err != nil {
		return err
	}
	if err := g.input.Setup(g.board, cfg); err != nil {
		return err
	}
	if err := g.initWasm(); err != nil {
		return err
	}
	return nil
}

// Start initializes the rendering process and exposes Go functions for JavaScript interaction via the global scope.
// It blocks indefinitely using a channel to keep the Go routine alive.
func (g *Render) Start() error {
	c := make(chan struct{})
	<-c
	return nil
}

// VBlank sets the vBlank flag to true, signaling the start of the vertical blanking interval in the rendering process.
func (g *Render) VBlank() {
	g.vBlank = true
	//if g.win.MouseInsideWindow() {
	//	g.inputs.MouseMove(g.win.MousePositionXY())
	//}
	//g.inputs.Keys(g.win.KeysPressed())
	//g.surface.Draw(g.win, g.surfaceM)
	//if g.led {
	//	g.ledSurface.Draw(g.win, g.ledSurfaceM)
	//}
	//g.win.Update()
	//g.run = !g.win.Closed()
}

// LedActivity toggles the LED state for the specified device number based on the given boolean value.
func (g *Render) LedActivity(deviceNumber uint8, led bool) {
	//g.led = led
	//fmt.Println("LED STATE", deviceNumber, led)
}

func (g *Render) initWasm() error {
	emulateFrame := func(this js.Value, args []js.Value) interface{} {
		for {
			g.board.Emulate()
			if g.vBlank {
				g.vBlank = false
				return nil
			}
		}
	}
	getSurfacePointer := func(this js.Value, args []js.Value) interface{} {
		surfacePtr := g.displayBuffer.GetSurfacePointer()
		return uintptr(surfacePtr)
	}
	getSurfaceLen := func(this js.Value, args []js.Value) interface{} {
		return js.ValueOf(g.displayBuffer.GetSurfaceLen())
	}
	getDisplayBuffer := func(this js.Value, args []js.Value) interface{} {
		surfaceBytes := g.displayBuffer.GetSurface()
		jsBuffer := js.Global().Get("Uint8Array").New(len(surfaceBytes))
		js.CopyBytesToJS(jsBuffer, surfaceBytes)
		return jsBuffer
	}
	getDisplayWidth := func(this js.Value, args []js.Value) interface{} {
		return js.ValueOf(g.w)
	}
	getDisplayHeight := func(this js.Value, args []js.Value) interface{} {
		return js.ValueOf(g.h)
	}
	keyPressed := func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			if args[0].Type() == js.TypeString {
				keyCode := args[0].String()
				g.input.Key(keyCode, true)
			}
		}
		return nil
	}
	keyReleased := func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			if args[0].Type() == js.TypeString {
				keyCode := args[0].String()
				g.input.Key(keyCode, false)
			}
		}
		return nil
	}
	js.Global().Set("getDisplayBuffer", js.FuncOf(getDisplayBuffer))
	js.Global().Set("getDisplayWidth", js.FuncOf(getDisplayWidth))
	js.Global().Set("getDisplayHeight", js.FuncOf(getDisplayHeight))
	js.Global().Set("getSurfacePointer", js.FuncOf(getSurfacePointer))
	js.Global().Set("getSurfaceLen", js.FuncOf(getSurfaceLen))
	js.Global().Set("emulateFrame", js.FuncOf(emulateFrame))
	js.Global().Set("keyPressed", js.FuncOf(keyPressed))
	js.Global().Set("keyReleased", js.FuncOf(keyReleased))
	return nil
}
